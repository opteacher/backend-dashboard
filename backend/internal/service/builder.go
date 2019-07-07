package service

import (
	"sync"
	"context"
	"reflect"
	"strings"
	"fmt"
	"strconv"

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
	steps []*pb.OperStep
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
		if err := dao.CreateTx(tx, model.API_INFO_TABLE, reflect.TypeOf((*struct {
			Name   string `orm:",PRIMARY_KEY|UNIQUE_KEY"`
			Model  string
			Table  string
			Params string
			Route  string
			Method string
			Return string
			Flows  string
		})(nil)).Elem()); err != nil {
			panic(err)
		} else if err := dao.CreateTx(tx, model.OPER_STEP_TABLE, reflect.TypeOf((*struct {
			OperKey  string `orm:"oper_key,PRIMARY_KEY"`
			Requires string
			Desc     string
			Inputs   string
			Outputs  string
			BlockIn  bool
			BlockOut bool
			Code     string
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
	actMap := map[string]string{
		"POST":   "Insert",
		"DELETE": "Delete",
		"PUT":    "Update",
		"GET":    "Select",
		"ALL":    "SelectAll",
	}
	res, err := pgb.dao.QueryTx(tx, model.MODELS_TABLE, "", []interface{}{})
	if err != nil {
		return nil, err
	}
	var modelApis []*pb.ApiInfo
	for _, mdl := range res {
		if !reflect.TypeOf(mdl["methods"]).ConvertibleTo(reflect.TypeOf(([]interface{})(nil))) {
			continue
		}
		mname := mdl["name"].(string)
		mmname := mname + "Array"
		mmfname := utils.ToPlural(strings.ToLower(mname))
		for _, method := range mdl["methods"].([]interface{}) {
			m := method.(string)
			aname, exs := actMap[m]
			if !exs {
				aname = "Select"
			}
			modelApi := &pb.ApiInfo{
				Name:   fmt.Sprintf("%s%s", aname, mname),
				Model:  mname,
				Table:  utils.ToPlural(utils.CamelToPascal(mname)),
				Params: make(map[string]string),
				Route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				Method: strings.ToLower(m),
			}
			genModelApiInfo(modelApi, mname, mmname, mmfname)
			modelApis = append(modelApis, modelApi)
			// 接口信息存入数据库
			ApiInfoToDB(pgb.dao, tx, modelApi)
		}
	}
	return nil
}

func (pgb *ProjGenBuilder) genModelApiInfo(modelApi *pb.ApiInfo, mname string, mmname string, mmfname string) {
	mtypeInCode := "pb." + mname
	mmtypeInCode := "pb." + mmname
	switch modelApi.Method {
	case "post":
		modelApi.Params["entry"] = mtypeInCode
		modelApi.Return = mtypeInCode
		modelApi.Flows = []*pb.OperStep{
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "entry",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"PACKAGE":  kpg.info.option.Name,
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("database_beginTx", nil),
			kpg.copyStep("database_insertTx", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME": modelApi.Table,
					"OBJ_MAP":    "omap",
				},
			}),
			kpg.copyStep("database_commitTx", nil),
			kpg.copyStep("assignment", map[string]interface{}{
				"Desc": "将改记录的数据库id赋予请求参数",
				"Inputs": map[string]string{
					"SOURCE": "id",
					"TARGET": "entry.Id",
				},
			}),
			kpg.copyStep("return_succeed", map[string]interface{}{
				"Inputs": map[string]string{
					"RETURN": "entry",
				},
			}),
		}
	case "delete":
		modelApi.Params["iden"] = "pb.IdenReqs"
		modelApi.Return = mtypeInCode
		modelApi.Flows = []*pb.OperStep{
			kpg.copyStep("database_beginTx", nil),
			kpg.copyStep("database_queryTx", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "`id`=?",
					"QUERY_ARGUS": "iden.Id",
				},
			}),
			kpg.copyStep("database_deleteTx", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "`id`=?",
					"QUERY_ARGUS": "iden.Id",
				},
			}),
			kpg.copyStep("database_commitTx", nil),
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "res",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("return_succeed", map[string]interface{}{
				"Inputs": map[string]string{
					"RETURN": fmt.Sprintf("omap.(*%s)", mtypeInCode),
				},
			}),
		}
	case "put":
		modelApi.Params["iden"] = "IdenReqs"
		modelApi.Params["entry"] = mtypeInCode
		modelApi.Return = mtypeInCode
		modelApi.Flows = []*pb.OperStep{
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "entry",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"PACKAGE":  kpg.info.option.Name,
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("database_beginTx", nil),
			kpg.copyStep("database_updateTx", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "`id`=?",
					"QUERY_ARGUS": "iden.Id",
					"OBJ_MAP":     "omap",
				},
			}),
			kpg.copyStep("database_queryTx", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "`id`=?",
					"QUERY_ARGUS": "iden.Id",
				},
			}),
			kpg.copyStep("database_commitTx", nil),
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "res",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("return_succeed", map[string]interface{}{
				"Inputs": map[string]string{
					"RETURN": fmt.Sprintf("omap.(*%s)", mtypeInCode),
				},
			}),
		}
	case "get":
		modelApi.Params["iden"] = "IdenReqs"
		modelApi.Return = mtypeInCode
		modelApi.Flows = []*pb.OperStep{
			kpg.copyStep("database_query", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "`id`=?",
					"QUERY_ARGUS": "iden.Id",
				},
			}),
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "res",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("return_succeed", map[string]interface{}{
				"Inputs": map[string]string{
					"RETURN": fmt.Sprintf("omap.(*%s)", mtypeInCode),
				},
			}),
		}
	case "all":
		modelApi.Method = "get"
		modelApi.Params["params"] = "Nil"
		modelApi.Return = mtypeInCode
		modelApi.Flows = []*pb.OperStep{
			kpg.copyStep("database_query", map[string]interface{}{
				"Inputs": map[string]string{
					"TABLE_NAME":  modelApi.Table,
					"QUERY_CONDS": "\"\"",
					"QUERY_ARGUS": "nil",
				},
			}),
			kpg.copyStep("assignment_create", map[string]interface{}{
				"Inputs": map[string]string{
					"SOURCE": fmt.Sprintf("new(%s)", mmtypeInCode),
					"TARGET": "resp",
				},
			}),
			kpg.copyStep("for_each", map[string]interface{}{
				"Inputs": map[string]string{
					"KEY": "_",
					"VALUE": "entry",
					"SET": "res",
				},
			}),
			kpg.copyStep("json_marshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJECT": "entry",
				},
			}),
			kpg.copyStep("json_unmarshal", map[string]interface{}{
				"Inputs": map[string]string{
					"OBJ_TYPE": mtypeInCode,
				},
			}),
			kpg.copyStep("assignment_append", map[string]interface{}{
				"Inputs": map[string]string{
					"ARRAY": "resp." + mmfname,
					"NEW_ADD": fmt.Sprintf("omap.(*%s)", mmtypeInCode),
				},
			}),
			kpg.copyStep("return_succeed", map[string]interface{}{
				"Inputs": map[string]string{
					"RETURN": "resp",
				},
				"BlockOut": true,
			}),
		}
	}
}

func (pgb *ProjGenBuilder) copyStep(operKey string, values map[string]interface{}) *pb.OperStep {
	for _, step := range kpg.operSteps {
		if step.OperKey == operKey {
			if values == nil {
				return step
			}
			ret, _ := utils.Copy(step)
			val := reflect.ValueOf(ret).Elem()
			for fname, fvalue := range values {
				val.FieldByName(fname).Set(reflect.ValueOf(fvalue))
			}
			return val.Addr().Interface().(*pb.OperStep)
		}
	}
	return nil
}

func (pgb *ProjGenBuilder) readStepsFromDB(tx *sql.Tx) error {
	if amap, err := kpg.dao.QueryTx(tx, model.OPER_STEP_TABLE, "", nil); err != nil {
		return err
	} else {
		for _, rmap := range amap {
			step := new(pb.OperStep)
			step.OperKey = rmap["operKey"].(string)
			step.Requires = strings.Split(rmap["requires"].(string), ",")
			step.Desc = rmap["desc"].(string)
			step.Inputs = StrToStrMap(rmap["inputs"].(string))
			step.Outputs = strings.Split(rmap["outputs"].(string), ",")
			step.BlockIn = rmap["blockIn"].(int64) == 1
			step.BlockOut = rmap["blockOut"].(int64) == 1
			step.Code = rmap["code"].(string)
			pgb.steps = append(pgb.steps, step)
		}
		return nil
	}
}

func (pgb *ProjGenBuilder) initOperSteps(dao *dao.Dao, tx *sql.Tx) error {
	if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_marshal",
		"desc":     "将收到的请求参数编码成JSON字节数组",
		"inputs":   "OBJECT:",
		"outputs":  "bytes",
		"requires": "encoding/json",
		"code":     "bytes, err := json.Marshal(%OBJECT%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"转JSON失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_unmarshal",
		"desc":     "将JSON字节数组转成Map键值对",
		"inputs":   "OBJ_TYPE:",
		"outputs":  "omap",
		"requires": "%PACKAGE%/internal/utils",
		"code":     "omap, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*%OBJ_TYPE%)(nil)).Elem())\nif err != nil {\n\treturn nil, fmt.Errorf(\"从JSON转回失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_beginTx",
		"desc":     "开启数据库事务",
		"outputs":  "tx",
		"code":     "tx, err := s.dao.BeginTx(ctx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"开启事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_commitTx",
		"desc":     "提交数据库事务",
		"code":     "err := s.dao.CommitTx(tx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"提交事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% = %SOURCE%\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment_append",
		"inputs": "ARRAY:,NEW_ADD:",
		"code": "%ARRAY% = append(%ARRAY%, %NEW_ADD%)\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment_create",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% := %SOURCE%\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "for_each",
		"inputs": "KEY:,VALUE:,SET:",
		"code": "for %KEY%, %VALUE% := range %SET%",
		"block_in": true,
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "return_succeed",
		"inputs":   "RETURN:",
		"code":     "return %RETURN%, nil\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_insertTx",
		"desc":     "做数据库插入操作",
		"inputs":   "TABLE_NAME:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.InsertTx(tx, \"%TABLE_NAME%\", %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"插入数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_queryTx",
		"desc":     "做数据库查询操作（事务）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.QueryTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_query",
		"desc":     "做数据库查询操作（会话）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.Query(ctx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_deleteTx",
		"desc":     "做数据库删除操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"code":     "_, err := s.dao.DeleteTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"删除数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_updateTx",
		"desc":     "做数据库更新操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.SaveTx(tx, \"%TABLE_NAME%\", \"%QUERY_CONDS%\", %QUERY_ARGUS%, %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"更新数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if err := pgb.readStepsFromDB(); err != nil {
		return err
	} else  {
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
	ary := make([]string, len(mp))
	for k, v := range mp {
		if len(k) == 0 && len(v) == 0 {
			continue
		}
		ary = append(ary, k + ":" + v)
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

func OperStepToDB(dao *dao.Dao, tx *sql.Tx, step *pb.OperStep) (int64, error) {
	mstep := make(map[string]interface{})
	mstep["oper_key"] = step.OperKey
	mstep["requires"] = strings.Join(step.Requires, ",")
	mstep["desc"] = step.Desc
	mstep["inputs"] = StrMapToStr(step.Inputs)
	mstep["outputs"] = strings.Join(step.Outputs, ",")
	mstep["block_in"] = step.BlockIn
	mstep["block_out"] = step.BlockIn
	mstep["code"] = step.Code
	if id, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, mstep); err != nil {
		return -1, err
	} else {
		return id, nil
	}
}

func OperStepFmDB(dao *dao.Dao, tx *sql.Tx, id int64) (*pb.OperStep, error) {
	if mstep, err := dao.QueryTxByID(tx, model.OPER_STEP_TABLE, id); err != nil {
		return nil, err
	} else {
		step := new(pb.OperStep)
		step.OperKey = mstep["oper_key"].(string)
		step.Requires = strings.Split(mstep["requires"].(string), ",")
		step.Desc = mstep["desc"].(string)
		step.Inputs = StrToStrMap(mstep["inputs"].(string))
		step.Outputs = strings.Split(mstep["outputs"].(string), ",")
		step.BlockIn = mstep["block_in"].(int32) == 1
		step.BlockOut = mstep["block_out"].(int32) == 1
		step.Code = mstep["code"].(string)
		return step, nil
	}
}

func ApiInfoToDB(dao *dao.Dao, tx *sql.Tx, info *pb.ApiInfo) (int64, error) {
	minfo := make(map[string]interface{})
	minfo["name"] = info.Name
	minfo["model"] = info.Model
	minfo["table"] = info.Table
	minfo["params"] = StrMapToStr(info.Params)
	minfo["route"] = info.Route
	minfo["method"] = info.Method
	minfo["return"] = info.Return
	sflows := ""
	for _, flow := range info.Flows {
		if id, err := OperStepToDB(dao, tx, flow); err != nil {
			return -1, err
		} else {
			sflows += strconv.Itoa(int(id)) + ","
		}
	}
	minfo["flows"] = strings.TrimRight(sflows, ",")
	return dao.InsertTx(tx, model.API_INFO_TABLE, minfo)
}

func CvtApiInfoFmMap(dao *dao.Dao, tx *sql.Tx, mapi map[string]interface{}) (*pb.ApiInfo, error) {
	info := new(pb.ApiInfo)
	info.Name = mapi["name"].(string)
	info.Model = mapi["model"].(string)
	info.Table = mapi["table"].(string)
	info.Params = StrToStrMap(mapi["params"].(string))
	info.Route = mapi["route"].(string)
	info.Method = mapi["method"].(string)
	info.Return = mapi["return"].(string)
	for _, flowID := range strings.Split(mapi["flows"].(string), ",") {
		if iflowID, err := strconv.Atoi(flowID); err != nil {
			return nil, err
		} else if flow, err := OperStepFmDB(dao, tx, int64(iflowID)); err != nil {
			return nil, err
		} else {
			info.Flows = append(info.Flows, flow)
		}
	}
	return info, nil
}

func ApiInfoFmDB(dao *dao.Dao, tx *sql.Tx, name string) (*pb.ApiInfo, error) {
	if apis, err := dao.QueryTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{name}); err != nil {
		return nil, err
	} else if len(apis) == 0 {
		return nil, fmt.Errorf("没有找到指定（name=%s）的接口", name)
	} else  {
		return CvtApiInfoFmMap(dao, tx, apis[0])
	}
}

// // 处理集合
// // 如果目标字段类型是数组，但源map中是字符串，则用逗号分隔之后再填入
// // 如果目标字段类型是map，但源map中是字符串，则用逗号分隔之后再用冒号分隔键名和键值
// func FillMapToObj(maps []map[string]interface{}, tgtTyp reflect.Type) (interface{}, error) {
// 	ret := reflect.MakeSlice(reflect.SliceOf(tgtTyp), len(maps), len(maps))
// 	if len(maps) == 0 {
// 		return ret.Interface(), nil
// 	}
// 	// 挑出存在于目标对象的属性分量
// 	tstEle := maps[0]
// 	fmap := make(map[string]string)
// 	for i := 0; i < tgtTyp.NumField(); i++ {
// 		field := tgtTyp.Field(i)
// 		mkey := utils.Uncapital(field.Name)
// 		if _, exs := tstEle[mkey]; exs {
// 			fmap[mkey] = field.Name
// 		}
// 	}
// 	// 从map填充进对象
// 	for i, mele := range maps {
// 		oele := ret.Index(i)
// 		for mfname, ofname := range fmap {
// 			ofield := oele.FieldByName(ofname)
// 			mfield := reflect.ValueOf(mele[mfname])
// 			ofkind := ofield.Type().Kind()
// 			mfkind := mfield.Type().Kind()
// 			if ofkind == mfkind {
// 				ofield.Set(mfield)
// 			} else if ofkind == reflect.Bool {
// 				ofield.SetBool(mele[mfname].(int64) == 1)
// 			} else if mfkind == reflect.String {
// 				str := mfield.String()
// 				if len(str) == 0 {
// 					continue
// 				}
// 				switch ofkind {
// 				case reflect.Map:
// 					kvs := strings.Split(str, ",")
// 					mapstr := make(map[string]string)
// 					for _, kv := range kvs {
// 						kvAry := strings.Split(kv, ":")
// 						mapstr[kvAry[0]] = kvAry[1]
// 					}
// 					mfield = reflect.ValueOf(mapstr)
// 					ofield.Set(mfield)
// 				case reflect.Array:
// 					fallthrough
// 				case reflect.Slice:
// 					mfield = reflect.ValueOf(strings.Split(str, ","))
// 					ofield.Set(mfield)
// 				}
// 			}
// 		}
// 	}
// 	return ret.Interface(), nil
// }

// // 与上面组成一队
// func FillObjToMap(objs []interface{}, setToStr bool) ([]map[string]interface{}) {
// 	var ret []map[string]interface{}
// 	if len(objs) == 0 {
// 		return ret
// 	}
// 	tstObj := objs[0]
// 	otype := reflect.TypeOf(tstObj)
// 	fnameMap := make(map[string]string)
// 	for i := 0; i < otype.NumField(); i++ {
// 		oname := otype.Name()
// 		fnameMap[oname] = utils.CamelToPascal(oname)
// 	}
// 	for _, obj := range objs {
// 		ovalue := reflect.ValueOf(obj)
// 		mvalue := make(map[string]interface{})
// 		for ofname, mfname := range fnameMap {
// 			field := ovalue.FieldByName(ofname)
// 			switch field.Type().Kind() {
// 			case reflect.Map:
// 				amap := make([]string, field.Len())
// 				for _, k := range field.MapKeys() {
// 					v := field.MapIndex(k)
// 					amap = append(amap, fmt.Sprintf("%s:%s", k.String(), v.String()))
// 				}
// 				mvalue[mfname] = strings.Join(amap, ",")
// 			case reflect.Array:
// 				fallthrough
// 			case reflect.Slice:
// 				str := ""
// 				for i := 0; i < field.Len(); i++ {
// 					str += field.Index(i).String() + ","
// 				}
// 				mvalue[mfname] = strings.TrimRight(str, ",")
// 			default:
// 				mvalue[mfname] = field.Interface()
// 			}
// 		}
// 		ret = append(ret, mvalue)
// 	}
// 	return ret
// }

// func StoragePbApi(ctx context.Context, dao *dao.Dao, api *pb.ApiInfo) error {
// 	// aparams := make([]string, len(api.Params))
// 	// for k, v := range api.Params {
// 	// 	aparams = append(aparams, fmt.Sprintf("%s:%s", k, v))
// 	// }
// 	// sparams := strings.Join(aparams, ",")

// 	if tx, err := dao.BeginTx(ctx); err != nil {
// 		return err
// 	} else {
// 		for _, flow := range api.Flows {
// 			flow.Code = ""
// 			if bytes, err := json.Marshal(flow); err != nil {
// 				return err
// 			} else if omap, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*map[string]interface{})(nil)).Elem()); err != nil {
// 				return err
// 			} else if id, err := dao.InsertTx(tx, model.OPER_STEP_TABLE, *(omap.(*map[string]interface{}))); err != nil {
// 				return err
// 			} else {
// 				fmt.Println(id)
// 			}
// 		}
// 	}
// 	return nil
// }