package service

import (
	"bufio"
	"context"
	"crypto/md5"
	"reflect"
	"time"

	"github.com/bilibili/kratos/pkg/database/sql"

	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/server"

	"fmt"

	"backend/internal/utils"
	"encoding/json"
	"os"
	"path"
	"strings"

	"io"

	"github.com/bilibili/kratos/pkg/conf/paladin"
)

// Service service.
type Service struct {
	ac *paladin.Map
	cc struct {
		Qiniu *utils.StorageConfig
	}
	dao       *dao.Dao
	operSteps []pb.OperStep
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	if err := paladin.Get("cdn.toml").UnmarshalTOML(&s.cc); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		panic(err)
	} else if err := s.dao.CreateTx(tx, model.MODELS_TABLE, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
		panic(err)
	} else if err := s.dao.CreateTx(tx, model.API_INFO_TABLE, reflect.TypeOf((*struct {
		Name   string `orm:",PRIMARY_KEY|UNIQUE_KEY"`
		Model  string
		Params string
		Route  string
		Method string
		Return string
		Flows  string
	})(nil)).Elem()); err != nil {
		panic(err)
	} else if err := s.dao.CreateTx(tx, model.OPER_STEP_TABLE, reflect.TypeOf((*struct {
		OperKey  string `orm:"oper_key,PRIMARY_KEY"`
		Requires string
		Desc     string
		Inputs   string
		Outputs  string
		Code     string
	})(nil)).Elem()); err != nil {
		panic(err)
	} else if err := s.initOperSteps(tx); err != nil {
		panic(err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		panic(err)
	} else {

	}
	return s
}

func (s *Service) initOperSteps(tx *sql.Tx) error {
	if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_marshal",
		"desc":     "将收到的请求参数编码成JSON字节数组",
		"inputs":   "OBJECT:",
		"outputs":  "bytes",
		"requires": "encoding/json",
		"code":     "bytes, err := json.Marshal(%OBJECT%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"转JSON失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "json_unmarshal",
		"desc":     "将JSON字节数组转成Map键值对",
		"inputs":   "OBJ_TYPE:",
		"outputs":  "omap",
		"requires": "%PACKAGE%/internal/utils",
		"code":     "omap, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*%OBJ_TYPE%)(nil)).Elem())\nif err != nil {\n\treturn nil, fmt.Errorf(\"从JSON转回失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_beginTx",
		"desc":     "开启数据库事务",
		"outputs":  "tx",
		"code":     "tx, err := s.dao.BeginTx(ctx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"开启事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_commitTx",
		"desc":     "提交数据库事务",
		"code":     "err := s.dao.CommitTx(tx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"提交事务失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "assignment",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% = %SOURCE%\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "return_succeed",
		"inputs":   "RETURN:",
		"code":     "return %RETURN%, nil\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_insertTx",
		"desc":     "做数据库插入操作",
		"inputs":   "TABLE_NAME:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.InsertTx(tx, \"%TABLE_NAME%\", %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"插入数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_queryTx",
		"desc":     "做数据库查询操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.QueryTx(ctx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_deleteTx",
		"desc":     "做数据库删除操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"code":     "_, err := s.dao.DeleteTx(tx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"删除数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_updateTx",
		"desc":     "做数据库更新操作",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:,OBJ_MAP:",
		"outputs":  "id",
		"code":     "id, err := s.dao.SaveTx(tx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%, %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"更新数据表记录失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Service) AppID() string {
	appID, _ := s.ac.Get("appID").String()
	return appID
}

func (s *Service) SwaggerFile() string {
	swagger, _ := s.ac.Get("swaggerFile").String()
	return swagger
}

func (s *Service) ModelsInsert(ctx context.Context, req *pb.Model) (*pb.Model, error) {
	mm := make(map[string]interface{})
	if bm, err := json.Marshal(req); err != nil {
		return nil, fmt.Errorf("转JSON失败：%v", err)
	} else if err := json.Unmarshal(bm, &mm); err != nil {
		return nil, fmt.Errorf("提交的请求参数格式有误：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if id, err := s.dao.InsertTx(tx, model.MODELS_TABLE, mm); err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交插入事务失败：%v", err)
	} else {
		req.Id = id
		return req, nil
	}
}

func (s *Service) ModelsDelete(ctx context.Context, req *pb.NameID) (*pb.Model, error) {
	conds := "`name`=?"
	argus := []interface{}{req.Name}
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if res, err := s.dao.QueryTx(tx, model.MODELS_TABLE, conds, argus); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else if _, err := s.dao.DeleteTx(tx, model.MODELS_TABLE, conds, argus); err != nil {
		return nil, fmt.Errorf("删除数据库记录失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("数据库提交失败：%v", err)
	} else if bd, err := json.Marshal(res[0]); err != nil {
		return nil, fmt.Errorf("转JSON字节码失败：%v", err)
	} else if resp, err := utils.UnmarshalJSON(bd, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
		return nil, fmt.Errorf("转Model对象失败：%v", err)
	} else {
		return resp.(*pb.Model), nil
	}
}

func (s *Service) ModelsUpdate(context.Context, *pb.Model) (*pb.Empty, error) {
	// TODO:
	return nil, nil
}

func (s *Service) ModelsSelectAll(ctx context.Context, req *pb.Empty) (*pb.ModelArray, error) {
	if res, err := s.dao.Query(ctx, model.MODELS_TABLE, "", []interface{}{}); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.ModelArray)
		for _, entry := range res {
			if bdata, err := json.Marshal(entry); err != nil {
				return nil, fmt.Errorf("转JSON字节码失败：%v", err)
			} else if mdl, err := utils.UnmarshalJSON(bdata, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
				return nil, fmt.Errorf("转Model对象失败：%v", err)
			} else {
				resp.Models = append(resp.Models, mdl.(*pb.Model))
			}
		}
		return resp, nil
	}
}

func (s *Service) ModelsSelectByName(context.Context, *pb.NameID) (*pb.Model, error) {
	// TODO:
	return nil, nil
}

func (s *Service) Export(ctx context.Context, req *pb.ExpOptions) (*pb.UrlResp, error) {
	if wsPath, err := s.ac.Get("workspace").String(); err != nil {
		return nil, fmt.Errorf("配置文件中未定义工作区目录：%v", err)
	} else if bin, err := time.Now().MarshalBinary(); err != nil {
		return nil, fmt.Errorf("生成临时文件夹名失败：%v", err)
	} else if cchName := fmt.Sprintf("%x", md5.Sum(bin)); false {
		return nil, nil
	} else if cchPath := path.Join(wsPath, "cache", cchName); false {
		return nil, nil
	} else if err := os.MkdirAll(cchPath, 0755); err != nil {
		return nil, fmt.Errorf("创建临时文件夹：%s失败：%v", cchPath, err)
	} else if tmpPath := path.Join(wsPath, "template", req.Type); false {
		return nil, nil
	} else if wsPath := path.Join(cchPath, req.Type); false {
		return nil, nil
	} else if utils.CopyFolder(tmpPath, wsPath); false {
		return nil, nil
	} else if err := s.editProject(ctx, req.Name, wsPath); err != nil {
		return nil, fmt.Errorf("编辑项目失败：%v", err)
	} else if wsFile, err := os.Open(wsPath); err != nil {
		return nil, fmt.Errorf("工作区目录有误，打开失败：%v", err)
	} else if zipPath := path.Join(cchPath, req.Name); false {
		return nil, nil
	} else if err := utils.Compress([]*os.File{wsFile}, zipPath); err != nil {
		wsFile.Close()
		return nil, fmt.Errorf("压缩项目失败：%v", err)
	} else if url, err := utils.Upload(zipPath, *s.cc.Qiniu); err != nil {
		wsFile.Close()
		return nil, fmt.Errorf("上传项目包失败：%v", err)
		//} else if err := os.RemoveAll(cchPath); err != nil {
		//	wsFile.Close()
		//	return nil, fmt.Errorf("删除临时文件夹失败：%v", err)
	} else {
		wsFile.Close()
		return &pb.UrlResp{Url: url}, nil
	}
}

func (s *Service) editProject(ctx context.Context, pjName string, pjPath string) error {
	if err := s.readOperFromDB(ctx); err != nil {
		return fmt.Errorf("读取服务流程项目失败：%v", err)
	} else if apis, err := s.genKratosProtoFile(ctx, pjName, pjPath); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := s.chgKratosServiceFile(ctx, pjPath, apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	}
	return nil
}

func (s *Service) readOperFromDB(ctx context.Context) error {
	if rmap, err := s.dao.Query(ctx, model.OPER_STEP_TABLE, "", nil); err != nil {
		return err
	} else if oary, err := dao.ConvertQueryResultToObj(rmap, reflect.TypeOf((*pb.OperStep)(nil)).Elem()); err != nil {
		return err
	} else if s.operSteps = oary.([]pb.OperStep); false {
		return nil
	} else {
		return nil
	}
}

// 根据数据库中模型的定义，生成proto的message和service
func (s *Service) genKratosProtoFile(ctx context.Context, pjName string, pjPath string) ([]*pb.ApiInfo, error) {
	if pjName[len(pjName)-4:] == ".zip" || pjName[len(pjName)-4:] == ".ZIP" {
		pjName = pjName[:len(pjName)-4]
	}
	pkgName := utils.CamelToPascal(pjName)
	// 添加proto文件并根据数据库添加message和service
	code := "syntax = \"proto3\";\n\n"
	code += fmt.Sprintf("package %s.service.v1;\n\n", pkgName)
	code += "import \"gogoproto/gogo.proto\";\n"
	code += "import \"google/api/annotations.proto\";\n\n"
	code += "option go_package = \"api\";\n"
	code += "option (gogoproto.goproto_getters_all) = false;\n\n"
	code += "message Nil {\n}\n\n"
	code += "message IdenReqs {\n\tint64 id = 1;\n}\n\n"

	res, err := s.dao.Query(ctx, model.MODELS_TABLE, "", []interface{}{})
	if err != nil {
		return nil, err
	}
	actMap := map[string]string{
		"POST":   "Insert",
		"DELETE": "Delete",
		"PUT":    "Update",
		"GET":    "Select",
		"ALL":    "SelectAll",
	}
	var modelApis []*pb.ApiInfo
	for _, mdl := range res {
		mname := mdl["name"].(string)
		code += fmt.Sprintf("message %s {\n", mname)
		for i, prop := range mdl["props"].([]map[string]interface{}) {
			code += fmt.Sprintf("\t%s %s=%d;\n", prop["type"], prop["name"], i+1)
		}
		code += "}\n\n"
		// 复数message
		mmname := mname + "Array"
		code += fmt.Sprintf("message %s {\n\t%s %s = 1;\n}\n\n", mmname, mname, utils.ToPlural(strings.ToLower(mname)))

		if !reflect.TypeOf(mdl["methods"]).ConvertibleTo(reflect.TypeOf(([]interface{})(nil))) {
			continue
		}
		for _, method := range mdl["methods"].([]interface{}) {
			m := method.(string)
			aname, exs := actMap[m]
			if !exs {
				aname = "Select"
			}
			modelApi := &pb.ApiInfo{
				Name:   fmt.Sprintf("%s%s", aname, mname),
				Model:  mname,
				Params: make(map[string]string),
				Route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				Method: strings.ToLower(m),
			}
			mtype := "pb." + mname
			mtable := utils.CamelToPascal(mname)
			switch modelApi.Method {
			case "post":
				modelApi.Params["entry"] = mtype
				modelApi.Return = mtype
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "entry",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"PACKAGE":  pjName,
							"OBJ_TYPE": mtype,
						},
					}),
					s.copyStep("database_beginTx", nil),
					s.copyStep("database_insertTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME": mtable,
							"OBJ_MAP":    "omap",
						},
					}),
					s.copyStep("database_commitTx", nil),
					s.copyStep("assignment", map[string]interface{}{
						"Desc": "将改记录的数据库id赋予请求参数",
						"Inputs": map[string]string{
							"SOURCE": "id",
							"TARGET": "entry.Id",
						},
					}),
					s.copyStep("return_succeed", map[string]interface{}{
						"Inputs": map[string]string{
							"RETURN": "entry",
						},
					}),
				}
			case "delete":
				modelApi.Params["iden"] = "pb.IdenReqs"
				modelApi.Return = mtype
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("database_beginTx", nil),
					s.copyStep("database_queryTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  mtable,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("database_deleteTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  mtable,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("database_commitTx", nil),
				}
			case "put":
				modelApi.Params["iden"] = "IdenReqs"
				modelApi.Params["entry"] = mtype
				modelApi.Return = mtype
			case "get":
				modelApi.Params["iden"] = "IdenReqs"
				modelApi.Return = mtype
			case "all":
				modelApi.Method = "get"
				modelApi.Params["params"] = "Nil"
				modelApi.Return = mtype
			}
			modelApis = append(modelApis, modelApi)
		}
	}

	if len(modelApis) != 0 {
		code += fmt.Sprintf("service %s {\n", utils.Capital(pjName))
	}
	for _, api := range modelApis {
		sparams := ""
		for _, ptyp := range api.Params {
			sparams += ptyp + ","
		}
		sparams = strings.TrimRight(sparams, ",")
		code += fmt.Sprintf("\trpc %s(%s) returns (%s)", api.Name, sparams, api.Return)
		if len(api.Route) != 0 && len(api.Method) != 0 {
			code += " {\n\t\toption (google.api.http) = {\n"
			code += fmt.Sprintf("\t\t\t%s: \"%s\"\n\t\t};\n\t};\n", api.Method, api.Route)
		} else {
			code += ";\n"
		}
	}
	code += "}\n\n"

	protoPath := path.Join(pjPath, "api", "api.proto")
	protoFile, err := os.OpenFile(protoPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer protoFile.Close()
	protoFile.WriteString(code)
	return modelApis, nil
}

func (s *Service) copyStep(operKey string, values map[string]interface{}) *pb.OperStep {
	for _, step := range s.operSteps {
		if step.OperKey == operKey {
			if values == nil {
				return &step
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

// 根据抽取的接口信息，生成完整的service
func (s *Service) chgKratosServiceFile(ctx context.Context, pjPath string, apis []*pb.ApiInfo) error {
	svcPath := path.Join(pjPath, "internal", "service", "service.go")
	svcFile, err := os.Open(svcPath)
	if err != nil {
		return err
	}
	defer svcFile.Close()
	reader := bufio.NewReader(svcFile)
	// 收集import文件
	requires := make(map[string]interface{})
	for _, a := range apis {
		for _, step := range a.Flows {
			for _, i := range step.Requires {
				requires[i] = nil
			}
		}
	}
	code := ""
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		switch {
		case strings.Contains(string(line), "[APIS]"):
			for _, ai := range apis {
				aparams := make([]string, 0)
				for pname, ptype := range ai.Params {
					aparams = append(aparams, fmt.Sprintf("%s *%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				code += fmt.Sprintf("func (s *Service) %s(ctx context.Context, %s) (*%s, error) {\n", ai.Name, sparams, ai.Return)
				for idx, step := range ai.Flows {
					if idx == 0 {
						code += "\t"
					}
					// 添加注释
					if len(step.Desc) != 0 {
						code += fmt.Sprintf("// %s\n", step.Desc)
					}
					// 提取步骤操作的代码
					cd := step.Code
					// 替换步骤中的占位符
					for o, n := range step.Inputs {
						cd = strings.Replace(cd, fmt.Sprintf("%%%s%%", o), n, -1)
					}
					code += utils.AddSpacesBeforeRow(cd, 1)
				}
				code += "\n}\n\n"
			}
		case strings.Contains(string(line), "[IMPORTS]"):
			for require, _ := range requires {
				code += fmt.Sprintf("\t\"%s\"\n", require)
			}
		default:
			code += string(line) + "\n"
		}
	}
	svcFile, err = os.OpenFile(svcPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer svcFile.Close()
	svcFile.WriteString(code)
	return nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
	// 注销服务
	if cli, err := server.RegisterService(); err != nil {
		panic(err)
	} else if _, err := cli.Cancel(context.Background(), &pb.IdenSvcReqs{AppID: s.AppID()}); err != nil {
		panic(err)
	}
}
