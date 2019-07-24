package service

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/bilibili/kratos/pkg/database/sql"

	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/utils"
)

type ProjGenBuilder struct {
	sync.Once
	gen *BaseGen
	dao *dao.Dao
}

type GenInfo struct {
	option   *pb.ExpOptions
	pathName string
	pkgName  string
}

type BaseGen interface {
	Adjust(context.Context) error
}

func NewProjGenBuilder(dao *dao.Dao, tx *sql.Tx) (*ProjGenBuilder, error) {
	pgb := new(ProjGenBuilder)
	pgb.dao = dao
	pgb.Once = sync.Once{}
	pgb.Once.Do(func() {
		if err := dao.DropTx(tx, model.OPER_STEP_TABLE); err != nil {
			panic(err)
		} else if err := dao.DropTx(tx, model.API_INFO_TABLE); err != nil {
			panic(err)
		} else if err := dao.CreateTx(tx, model.API_INFO_TABLE, reflect.TypeOf((*struct {
			Name    string `orm:",PRIMARY_KEY|UNIQUE_KEY"`
			Model   string `orm:",FOREIGN_KEY(models.name)"`
			Table   string
			Params  string
			Route   string
			Method  string
			Returns string
			Flows   string
		})(nil)).Elem()); err != nil {
			panic(err)
		} else if err := dao.CreateTx(tx, model.OPER_STEP_TABLE, reflect.TypeOf((*struct {
			OperKey  string `orm:"oper_key,PRIMARY_KEY"`
			Requires string
			Desc     string
			Inputs   string
			Outputs  string
			Code     string
			ApiName  string `orm:",FOREIGN_KEY(api_infos.name)"`
			Symbol   int32
		})(nil)).Elem()); err != nil {
			panic(err)
		} else if err := pgb.initOperSteps(tx); err != nil {
			panic(err)
		} else if err := pgb.initApiInfos(tx); err != nil {
			panic(err)
		}
	})
	return pgb, nil
}

func (pgb *ProjGenBuilder) initApiInfos(tx *sql.Tx) error {
	if res, err := pgb.dao.QueryTx(tx, model.MODELS_TABLE, "", []interface{}{}); err != nil {
		return err
	} else if len(res) == 0 {
		return nil
	} else if steps, err := ReadStepsFromDB(pgb.dao, tx); err != nil {
		return err
	} else {
		for _, mm := range res {
			if _, err := GenModelApiInfo(pgb.dao, tx, mm, steps); err != nil {
				return err
			}
		}
		return nil
	}
}

var actMap = map[string]string{
	"POST":   "Insert",
	"DELETE": "Delete",
	"PUT":    "Update",
	"GET":    "Select",
	"ALL":    "SelectAll",
}

func GenModelApiInfo(dao *dao.Dao, tx *sql.Tx, mdl map[string]interface{}, steps []*pb.OperStep) ([]*pb.ApiInfo, error) {
	if !reflect.TypeOf(mdl["methods"]).ConvertibleTo(reflect.TypeOf(([]interface{})(nil))) {
		// 该模块没有开启HTTP接口，直接返回
		return nil, nil
	}
	// 获取所有模板步骤记录
	if steps == nil || len(steps) == 0 {
		var err error
		if steps, err = ReadStepsFromDB(dao, tx); err != nil {
			return nil, err
		}
	}
	copyStep := func(src *pb.OperStep) *pb.OperStep {
		for _, step := range steps {
			if step.OperKey == src.OperKey {
				tgt := &pb.OperStep{
					OperKey: step.OperKey,
					Desc:    step.Desc,
					Code:     step.Code,
					Symbol:   step.Symbol,
				}
				// 步骤所需模块
				if len(step.Requires) != 0 {
					tgt.Requires = step.Requires
				}
				// 步骤的输入标识
				if len(src.Inputs) != 0 {
					tgt.Inputs = src.Inputs
				} else if len(step.Inputs) != 0 {
					tgt.Inputs = step.Inputs
				}
				// 步骤的输出标识
				if len(step.Outputs) != 0 {
					tgt.Outputs = step.Outputs
				}
				// 步骤描述
				if len(src.Desc) != 0 {
					tgt.Desc = src.Desc
				}
				// 将描述中的需替换文字用输入标识替换之
				for k, v := range src.Inputs {
					tgt.Desc = strings.Replace(tgt.Desc, fmt.Sprintf("%%%s%%", k), v, -1)
				}
				// 特殊标识
				if src.Symbol != pb.SpcSymbol_NONE {
					tgt.Symbol = src.Symbol
				}
				return tgt
			}
		}
		return nil
	}
	mdlApis := make([]*pb.ApiInfo, 0)
	mname := mdl["name"].(string)
	mmname := mname + "Array"
	mmfname := utils.ToPlural(strings.ToLower(mname))
	for _, method := range mdl["methods"].([]interface{}) {
		m := method.(string)
		aname, exs := actMap[m]
		if !exs {
			aname = "Select"
		}
		mdlApi := &pb.ApiInfo{
			Name:   fmt.Sprintf("%s%s", aname, mname),
			Model:  mname,
			Table:  utils.ToPlural(utils.CamelToPascal(mname)),
			Params: make(map[string]string),
			Route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
			Method: strings.ToLower(m),
		}
		switch mdlApi.Method {
		case "post":
			mdlApi.Params["entry"] = mname
			mdlApi.Returns = []string{mname}
			mdlApi.Flows = []*pb.OperStep{
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "entry",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_beginTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_insertTx",
					Inputs: map[string]string{
						"TABLE_NAME": mdlApi.Table,
						"OBJ_MAP":    "omap",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_commitTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "assignment",
					Desc:    "将改记录的数据库id赋予请求参数",
					Inputs: map[string]string{
						"SOURCE": "id",
						"TARGET": "entry.Id",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "return_succeed",
					Inputs: map[string]string{
						"RETURN": "entry",
					},
				}),
			}
		case "delete":
			mdlApi.Params["iden"] = "IdenReqs"
			mdlApi.Returns = []string{mname}
			mdlApi.Flows = []*pb.OperStep{
				copyStep(&pb.OperStep{
					OperKey: "database_beginTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_queryTx",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "`id`=?",
						"QUERY_ARGUS": "iden.Id",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_deleteTx",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "`id`=?",
						"QUERY_ARGUS": "iden.Id",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_commitTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "res",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "return_succeed",
					Inputs: map[string]string{
						"RETURN": fmt.Sprintf("omap.(*%s)", mname),
					},
				}),
			}
		case "put":
			mdlApi.Params["iden"] = "IdenReqs"
			mdlApi.Params["entry"] = mname
			mdlApi.Returns = []string{mname}
			mdlApi.Flows = []*pb.OperStep{
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "entry",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_beginTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_updateTx",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "`id`=?",
						"QUERY_ARGUS": "iden.Id",
						"OBJ_MAP":     "omap",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_queryTx",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "`id`=?",
						"QUERY_ARGUS": "iden.Id",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "database_commitTx",
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "res",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "return_succeed",
					Inputs: map[string]string{
						"RETURN": fmt.Sprintf("omap.(*%s)", mname),
					},
				}),
			}
		case "get":
			mdlApi.Params["iden"] = "IdenReqs"
			mdlApi.Returns = []string{mname}
			mdlApi.Flows = []*pb.OperStep{
				copyStep(&pb.OperStep{
					OperKey: "database_query",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "`id`=?",
						"QUERY_ARGUS": "iden.Id",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "res",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "return_succeed",
					Inputs: map[string]string{
						"RETURN": fmt.Sprintf("omap.(*%s)", mname),
					},
				}),
			}
		case "all":
			mdlApi.Method = "get"
			mdlApi.Params["params"] = "Nil"
			mdlApi.Returns = []string{mmname}
			mdlApi.Flows = []*pb.OperStep{
				copyStep(&pb.OperStep{
					OperKey: "database_query",
					Inputs: map[string]string{
						"TABLE_NAME":  mdlApi.Table,
						"QUERY_CONDS": "\"\"",
						"QUERY_ARGUS": "nil",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "assignment_create",
					Inputs: map[string]string{
						"SOURCE": fmt.Sprintf("new(%s)", mmname),
						"TARGET": "resp",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "for_each",
					Inputs: map[string]string{
						"KEY":   "_",
						"VALUE": "entry",
						"SET":   "res",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_marshal",
					Inputs: map[string]string{
						"OBJECT": "entry",
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "json_unmarshal",
					Inputs: map[string]string{
						"OBJ_TYPE": mname,
					},
				}),
				copyStep(&pb.OperStep{
					OperKey: "assignment_append",
					Inputs: map[string]string{
						"ARRAY":   "resp." + mmfname,
						"NEW_ADD": fmt.Sprintf("omap.(*%s)", mmname),
					},
					Symbol: pb.SpcSymbol_BLOCK_OUT,
				}),
				copyStep(&pb.OperStep{
					OperKey: "return_succeed",
					Inputs: map[string]string{
						"RETURN": "resp",
					},
					Symbol: pb.SpcSymbol_BLOCK_OUT,
				}),
			}
		}
		// 接口信息存入数据库
		if _, err := ApiInfoToDbByTx(dao, tx, mdlApi); err != nil {
			return nil, fmt.Errorf("接口存入数据库发生异常：%v", err)
		}
		mdlApis = append(mdlApis, mdlApi)
	}
	return mdlApis, nil
}

func ReadStepsFromDB(dao *dao.Dao, tx *sql.Tx) ([]*pb.OperStep, error) {
	steps := make([]*pb.OperStep, 0)
	if amap, err := dao.QueryTx(tx, model.OPER_STEP_TABLE, "`api_name` IS NULL", nil); err != nil {
		return steps, err
	} else {
		for _, rmap := range amap {
			step := new(pb.OperStep)
			step.OperKey = rmap["oper_key"].(string)
			step.Requires = strings.Split(rmap["requires"].(string), ",")
			step.Desc = rmap["desc"].(string)
			step.Inputs = StrToStrMap(rmap["inputs"].(string))
			step.Outputs = strings.Split(rmap["outputs"].(string), ",")
			step.Code = rmap["code"].(string)
			step.ApiName = rmap["api_name"].(string)
			n, _ := rmap["symbol"].(int64)
			step.Symbol = pb.SpcSymbol(n)
			steps = append(steps, step)
		}
		return steps, nil
	}
}

func (pgb *ProjGenBuilder) initOperSteps(tx *sql.Tx) error {
	if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_marshal",
		"desc":     "将收到的请求参数编码成JSON字节数组",
		"inputs":   "OBJECT:",
		"outputs":  "bytes",
		"requires": "encoding/json",
		"code":     "bytes, err := json.Marshal(%OBJECT%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"转JSON失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_unmarshal",
		"desc":     "将JSON字节数组转成Map键值对",
		"inputs":   "OBJ_TYPE:",
		"outputs":  "omap",
		"requires": "%PACKAGE%/internal/utils",
		"code":     "omap, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*%OBJ_TYPE%)(nil)).Elem())\nif err != nil {\n\treturn nil, fmt.Errorf(\"从JSON转回失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_beginTx",
		"desc":     "开启数据库事务",
		"outputs":  "tx",
		"code":     "tx, err := s.dao.BeginTx(ctx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"开启事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_commitTx",
		"desc":     "提交数据库事务",
		"code":     "err := s.dao.CommitTx(tx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"提交事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% = %SOURCE%\n",
		"desc":     "将%SOURCE%赋值给%TARGET%",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment_append",
		"inputs":   "ARRAY:,NEW_ADD:",
		"code":     "%ARRAY% = append(%ARRAY%, %NEW_ADD%)\n",
		"desc":     "将%NEW_ADD%添加进%ARRAY%",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment_create",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% := %SOURCE%\n",
		"desc":     "创建%TARGET%并用%SOURCE%初始化",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "for_each",
		"inputs":   "KEY:,VALUE:,SET:",
		"code":     "for %KEY%, %VALUE% := range %SET%",
		"desc":     "循环遍历%SET%",
		"symbol":   pb.SpcSymbol_BLOCK_IN,
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "return_succeed",
		"inputs":   "RETURN:",
		"code":     "return %RETURN%, nil\n",
		"desc":     "成功返回%RETURN%",
		"symbol":   pb.SpcSymbol_END,
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_insertTx",
		"desc":     "做数据库插入操作",
		"inputs":   "TABLE_NAME:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.InsertTx(tx, \"%TABLE_NAME%\", %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"插入数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_queryTx",
		"desc":     "做数据库查询操作（事务）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.QueryTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_query",
		"desc":     "做数据库查询操作（会话）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.Query(ctx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_deleteTx",
		"desc":     "做数据库删除操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"code":     "_, err := s.dao.DeleteTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"删除数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := pgb.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_updateTx",
		"desc":     "做数据库更新操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.SaveTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%, %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"更新数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (pgb *ProjGenBuilder) Build(ctx context.Context, option *pb.ExpOptions, pathName string) BaseGen {
	info := new(GenInfo)
	info.option = option
	info.option.Name = strings.TrimRight(info.option.Name, ".zip")
	info.option.Name = strings.TrimRight(info.option.Name, ".ZIP")
	info.pkgName = utils.CamelToPascal(info.option.Name)
	info.pathName = pathName
	if gen, err := NewKratosProjGenerator(pgb.dao, info); err != nil {
		panic(err)
	} else {
		return gen
	}
}

func StrMapToStr(mp map[string]string) string {
	ary := make([]string, 0)
	for k, v := range mp {
		if len(k) == 0 && len(v) == 0 {
			continue
		}
		ary = append(ary, k+":"+v)
	}
	return strings.Join(ary, ",")
}

func StrToStrMap(str string) map[string]string {
	mp := make(map[string]string)
	for _, s := range strings.Split(str, ",") {
		if len(s) == 0 {
			return mp
		}
		as := strings.Split(s, ":")
		mp[as[0]] = as[1]
	}
	return mp
}

func OperStepToDbByTx(dao *dao.Dao, tx *sql.Tx, step *pb.OperStep) (int64, error) {
	mstep := make(map[string]interface{})
	mstep["oper_key"] = step.OperKey
	mstep["requires"] = strings.Join(step.Requires, ",")
	mstep["desc"] = step.Desc
	mstep["inputs"] = StrMapToStr(step.Inputs)
	mstep["outputs"] = strings.Join(step.Outputs, ",")
	mstep["code"] = step.Code
	mstep["api_name"] = step.ApiName
	mstep["symbol"] = step.Symbol
	if id, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, mstep); err != nil {
		return -1, err
	} else {
		return id, nil
	}
}

func OperStepFmDbByTx(dao *dao.Dao, tx *sql.Tx, id int64) (*pb.OperStep, error) {
	if mstep, err := dao.QueryTxByID(tx, model.OPER_STEP_TABLE, id); err != nil {
		return nil, err
	} else {
		step := new(pb.OperStep)
		step.Id = id
		step.OperKey = mstep["oper_key"].(string)
		step.Requires = strings.Split(mstep["requires"].(string), ",")
		step.Desc = mstep["desc"].(string)
		step.Inputs = StrToStrMap(mstep["inputs"].(string))
		step.Outputs = strings.Split(mstep["outputs"].(string), ",")
		step.Code = mstep["code"].(string)
		step.ApiName = mstep["api_name"].(string)
		n, _ := mstep["symbol"].(int64)
		step.Symbol = pb.SpcSymbol(n)
		return step, nil
	}
}

func ApiInfoToDbByTx(dao *dao.Dao, tx *sql.Tx, info *pb.ApiInfo) (int64, error) {
	minfo := make(map[string]interface{})
	minfo["name"] = info.Name
	if len(info.Model) != 0 {
		minfo["model"] = info.Model
	}
	minfo["table"] = info.Table
	minfo["params"] = StrMapToStr(info.Params)
	minfo["route"] = info.Route
	minfo["method"] = info.Method
	minfo["returns"] = strings.Join(info.Returns, ",")
	if id, err := dao.InsertTx(tx, model.API_INFO_TABLE, minfo); err != nil {
		return -1, err
	} else {
		sflows := ""
		for _, flow := range info.Flows {
			flow.ApiName = info.Name
			if id, err := OperStepToDbByTx(dao, tx, flow); err != nil {
				return -1, err
			} else {
				sflows += strconv.Itoa(int(id)) + ","
			}
		}
		return dao.SaveTx(tx, model.API_INFO_TABLE, "`id`=?", []interface{}{id}, map[string]interface{}{
			"flows": strings.TrimRight(sflows, ","),
		}, false)
	}
}

func CvtApiInfoFmMap(dao *dao.Dao, tx *sql.Tx, mapi map[string]interface{}) (*pb.ApiInfo, error) {
	info := new(pb.ApiInfo)
	info.Name = mapi["name"].(string)
	info.Model = mapi["model"].(string)
	info.Table = mapi["table"].(string)
	info.Params = StrToStrMap(mapi["params"].(string))
	info.Route = mapi["route"].(string)
	info.Method = mapi["method"].(string)
	info.Returns = strings.Split(mapi["returns"].(string), ",")
	sflows := mapi["flows"].(string)
	if len(sflows) == 0 {
		return info, nil
	}
	mps, err := dao.QueryTx(tx, model.OPER_STEP_TABLE, fmt.Sprintf("`id` IN (%s)", sflows), nil)
	if err != nil {
		return nil, err
	}
	for _, mp := range mps {
		info.Flows = append(info.Flows, CvtOperStepFmMap(mp))
	}
	return info, nil
}

func chkStrEmpty(str interface{}) bool {
	return str != nil && reflect.TypeOf(str).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) && str.(string) != ""
}

func CvtOperStepFmMap(mstep map[string]interface{}) *pb.OperStep {
	step := new(pb.OperStep)
	step.OperKey = mstep["oper_key"].(string)
	if chkStrEmpty(mstep["requires"]) {
		step.Requires = strings.Split(mstep["requires"].(string), ",")
	}
	if chkStrEmpty(mstep["desc"]) {
		step.Desc = mstep["desc"].(string)
	}
	if chkStrEmpty(mstep["inputs"]) {
		step.Inputs = StrToStrMap(mstep["inputs"].(string))
	}
	if chkStrEmpty(mstep["outputs"]) {
		step.Outputs = strings.Split(mstep["outputs"].(string), ",")
	}
	step.Code = mstep["code"].(string)
	if chkStrEmpty(mstep["api_name"]) {
		step.ApiName = mstep["api_name"].(string)
	}
	return step
}

func ApiInfoFmDbByTx(dao *dao.Dao, tx *sql.Tx, name string) (*pb.ApiInfo, error) {
	if apis, err := dao.QueryTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{name}); err != nil {
		return nil, err
	} else if len(apis) == 0 {
		return nil, fmt.Errorf("没有找到指定（name=%s）的接口", name)
	} else {
		return CvtApiInfoFmMap(dao, tx, apis[0])
	}
}

func AllApiInfoFmDB(dao *dao.Dao, ctx context.Context) ([]*pb.ApiInfo, error) {
	if tx, err := dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启数据库事务失败：%v", err)
	} else if apis, err := dao.QueryTx(tx, model.API_INFO_TABLE, "", nil); err != nil {
		return nil, fmt.Errorf("从数据库查找接口失败：%v", err)
	} else {
		infos := make([]*pb.ApiInfo, 0)
		for _, api := range apis {
			if info, err := CvtApiInfoFmMap(dao, tx, api); err != nil {
				return nil, fmt.Errorf("组装接口失败：%v", err)
			} else {
				infos = append(infos, info)
			}
		}
		return infos, nil
	}
}
