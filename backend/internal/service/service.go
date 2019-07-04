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
	info struct{
		projName string
		pathName string
		pkgName string
		isMicoSv bool
	}
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
	} else if err := s.dao.CreateTx(tx, model.LINKS_TABLE, reflect.TypeOf((*pb.Link)(nil)).Elem()); err != nil {
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
		"oper_key": "assignment_create",
		"inputs":   "SOURCE:,TARGET:",
		"code":     "%TARGET% := %SOURCE%\n",
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
		"desc":     "做数据库查询操作（事务）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.QueryTx(tx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	}); err != nil {
		return err
	} else if _, err := s.dao.InsertTx(tx, model.OPER_STEP_TABLE, map[string]interface{}{
		"oper_key": "database_query",
		"desc":     "做数据库查询操作（会话）",
		"inputs":   "TABLE_NAME:,QUERY_CONDS:,QUERY_ARGUS:",
		"outputs":  "res",
		"code":     "res, err := s.dao.Query(ctx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
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

func (s *Service) LinkInsert(ctx context.Context, req *pb.Link) (*pb.Link, error) {
	if bm, err := json.Marshal(req); err != nil {
		return nil, fmt.Errorf("转JSON失败：%v", err)
	} else if mm, err := utils.UnmarshalJSON(bm, reflect.TypeOf((*map[string]interface{})(nil)).Elem()); err != nil {
		return nil, fmt.Errorf("提交的请求参数格式有误：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if id, err := s.dao.InsertTx(tx, model.LINKS_TABLE, *(mm.(*map[string]interface{}))); err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交插入事务失败：%v", err)
	} else {
		req.Id = id
		return req, nil
	}
	return nil, nil
}

func (s *Service) LinkSelectAll(context.Context, *pb.Empty) (*pb.LinkArray, error) {
	// TODO:
	return nil, nil
}

func (s *Service) LinksDeleteBySymbol(context.Context, *pb.SymbolID) (*pb.Link, error) {
	// TODO:
	return nil, nil
}

func (s *Service) Export(ctx context.Context, req *pb.ExpOptions) (*pb.UrlResp, error) {
	if s.info.isMicoSv = req.IsMicoServ; false {
		return nil, nil
	} else if s.info.projName = strings.TrimRight(req.Name, ".zip"); false {
		return nil, nil
	} else if s.info.projName = strings.TrimRight(s.info.projName, ".ZIP"); false {
		return nil, nil
	} else if s.info.pkgName = utils.CamelToPascal(s.info.projName); false {
		return nil, nil
	} else if wsPath, err := s.ac.Get("workspace").String(); err != nil {
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
	} else if s.info.pathName = path.Join(cchPath, req.Type); false {
		return nil, nil
	} else if utils.CopyFolder(tmpPath, s.info.pathName); false {
		return nil, nil
	} else if err := s.editProject(ctx); err != nil {
		return nil, fmt.Errorf("编辑项目失败：%v", err)
	} else if wsFile, err := os.Open(s.info.pathName); err != nil {
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

func (s *Service) editProject(ctx context.Context) error {
	if err := s.readOperFromDB(ctx); err != nil {
		return fmt.Errorf("读取服务流程项目失败：%v", err)
	} else if apis, err := s.genKratosProtoFile(ctx); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := s.chgKratosServiceFile(ctx, apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	} else if err := s.switchKratosMicoServ(ctx); err != nil {
		return fmt.Errorf("开启/关闭微服务功能失败:%v", err)
	} else if err := s.chgKratosProjName(ctx); err != nil {
		return fmt.Errorf("修改项目名称失败：%v", err)
	}
	return nil
}

func (s *Service) adjServerFile(pathName string, regSvr string, regSvc string) error {
	fpath := path.Join(s.info.pathName, "internal", "server", pathName, "server.go")
	if err := utils.InsertTxt(fpath, func(line string, doProc *bool) (string, bool, error) {
		if !*doProc {
			return line, false, nil
		}
		if strings.Contains(line, "svr \"") && !s.info.isMicoSv {
			return "", false, nil
		} else if strings.Contains(line, fmt.Sprintf("pb.%s", regSvr)) {
			regName := fmt.Sprintf("Register%sServer", utils.Capital(s.info.projName))
			return strings.Replace(line, regSvr, regName, -1), false, nil
		} else if strings.Contains(line, fmt.Sprintf("func %s", regSvc)) {
			*doProc = s.info.isMicoSv
		} else if strings.Contains(line, regSvc) && !s.info.isMicoSv {
			return "", false, nil
		}
		return line, false, nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) switchKratosMicoServ(ctx context.Context) error {
	// 删除api/register.*
	if !s.info.isMicoSv {
		for _, p := range []string{
			path.Join(s.info.pathName, "api", "register.proto"),
			path.Join(s.info.pathName, "api", "register.bm.go"),
			path.Join(s.info.pathName, "api", "register.pb.go"),
			path.Join(s.info.pathName, "api", "register.swagger.json"),
		} {
			if err := os.Remove(p); err != nil {
				return err
			}
		}
	}
	// 删除internal/server/common.go
	if !s.info.isMicoSv {
		if err := os.Remove(path.Join(s.info.pathName, "internal", "server", "common.go")); err != nil {
			return err
		}
	}
	// 调整internal/server/grpc/server.go的逻辑
	if err := s.adjServerFile("grpc", "RegisterDemoServer", "RegisterGRPCService"); err != nil {
		return err
	}
	// 调整internal/server/http/server.go的逻辑
	if err := s.adjServerFile("http", "RegisterDemoBMServer", "RegisterHTTPService"); err != nil {
		return err
	}
	return nil
}

func (s *Service) chgKratosProjName(ctx context.Context) error {
	fixLst := []string{
		path.Join(s.info.pathName, "internal", "dao", "dao.go"),
		path.Join(s.info.pathName, "internal", "server", "grpc", "server.go"),
		path.Join(s.info.pathName, "internal", "server", "http", "server.go"),
	}
	if s.info.isMicoSv {
		fixLst = append(fixLst, path.Join(s.info.pathName, "internal", "server", "common.go"))
	}
	impPkg := fmt.Sprintf("\"%s/", s.info.pkgName)
	for _, p := range fixLst {
		if err := utils.InsertTxt(p, func(line string, doProc *bool) (string, bool, error) {
			if strings.Contains(line, ")") {
				*doProc = false
			}
			if !*doProc {
				return line, false, nil
			}
			return strings.Replace(line, "\"template/", impPkg, -1), false, nil
		}); err != nil {
			return err
		}
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
func (s *Service) genKratosProtoFile(ctx context.Context) ([]*pb.ApiInfo, error) {
	// 添加proto文件并根据数据库添加message和service
	code := "syntax = \"proto3\";\n\n"
	code += fmt.Sprintf("package %s.service.v1;\n\n", s.info.pkgName)
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
				Model:  "pb." + mname,
				Table: utils.ToPlural(utils.CamelToPascal(mname)),
				Params: make(map[string]string),
				Route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				Method: strings.ToLower(m),
			}
			switch modelApi.Method {
			case "post":
				modelApi.Params["entry"] = modelApi.Model
				modelApi.Return = modelApi.Model
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "entry",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"PACKAGE":  s.info.projName,
							"OBJ_TYPE": modelApi.Model,
						},
					}),
					s.copyStep("database_beginTx", nil),
					s.copyStep("database_insertTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME": modelApi.Table,
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
				modelApi.Return = modelApi.Model
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("database_beginTx", nil),
					s.copyStep("database_queryTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  modelApi.Table,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("database_deleteTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  modelApi.Table,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("database_commitTx", nil),
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "res",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJ_TYPE": modelApi.Model,
						},
					}),
					s.copyStep("return_succeed", map[string]interface{}{
						"Inputs": map[string]string{
							"RETURN": fmt.Sprintf("omap.(*%s)", modelApi.Model),
						},
					}),
				}
			case "put":
				modelApi.Params["iden"] = "IdenReqs"
				modelApi.Params["entry"] = modelApi.Model
				modelApi.Return = modelApi.Model
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "entry",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"PACKAGE":  s.info.projName,
							"OBJ_TYPE": modelApi.Model,
						},
					}),
					s.copyStep("database_beginTx", nil),
					s.copyStep("database_updateTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  modelApi.Table,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
							"OBJ_MAP": "omap",
						},
					}),
					s.copyStep("database_queryTx", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  modelApi.Table,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("database_commitTx", nil),
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "res",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJ_TYPE": modelApi.Model,
						},
					}),
					s.copyStep("return_succeed", map[string]interface{}{
						"Inputs": map[string]string{
							"RETURN": fmt.Sprintf("omap.(*%s)", modelApi.Model),
						},
					}),
				}
			case "get":
				modelApi.Params["iden"] = "IdenReqs"
				modelApi.Return = modelApi.Model
				modelApi.Flows = []*pb.OperStep{
					s.copyStep("database_query", map[string]interface{}{
						"Inputs": map[string]string{
							"TABLE_NAME":  modelApi.Table,
							"QUERY_CONDS": "`id`=?",
							"QUERY_ARGUS": "iden.Id",
						},
					}),
					s.copyStep("json_marshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJECT": "res",
						},
					}),
					s.copyStep("json_unmarshal", map[string]interface{}{
						"Inputs": map[string]string{
							"OBJ_TYPE": modelApi.Model,
						},
					}),
					s.copyStep("return_succeed", map[string]interface{}{
						"Inputs": map[string]string{
							"RETURN": fmt.Sprintf("omap.(*%s)", modelApi.Model),
						},
					}),
				}
			case "all":
				modelApi.Method = "get"
				modelApi.Params["params"] = "Nil"
				modelApi.Return = modelApi.Model
			}
			modelApis = append(modelApis, modelApi)
		}
	}

	if len(modelApis) != 0 {
		code += fmt.Sprintf("service %s {\n", utils.Capital(s.info.projName))
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

	protoPath := path.Join(s.info.pathName, "api", "api.proto")
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
func (s *Service) chgKratosServiceFile(ctx context.Context, apis []*pb.ApiInfo) error {
	svcPath := path.Join(s.info.pathName, "internal", "service", "service.go")
	svcFile, err := os.Open(svcPath)
	if err != nil {
		return err
	}
	defer svcFile.Close()
	reader := bufio.NewReader(svcFile)
	// 收集import文件
	requires := make(map[string]interface{})
	models := make(map[string]string)
	for _, a := range apis {
		models[a.Model] = a.Table
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
		case strings.Contains(string(line), "\"template/"):
			// 替换自身包名
			code += strings.Replace(string(line), "\"template/", fmt.Sprintf("\"%s/", s.info.pkgName), -1) + "\n"
		case strings.Contains(string(line), "[APIS]"):
			for _, ai := range apis {
				aparams := make([]string, 0)
				for pname, ptype := range ai.Params {
					aparams = append(aparams, fmt.Sprintf("%s *%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				code += fmt.Sprintf("func (s *Service) %s(ctx context.Context, %s) (*%s, error) {\n", ai.Name, sparams, ai.Return)
				for _, step := range ai.Flows {
					// 添加注释
					if len(step.Desc) != 0 {
						code += utils.AddSpacesBeforeRow(fmt.Sprintf("// %s\n", step.Desc), 1)
					}
					// 提取步骤操作的代码
					cd := step.Code
					// 替换步骤中的占位符
					for o, n := range step.Inputs {
						cd = strings.Replace(cd, fmt.Sprintf("%%%s%%", o), n, -1)
					}
					code += utils.AddSpacesBeforeRow(cd, 1)
				}
				code += "}\n\n"
			}
		case strings.Contains(string(line), "[IMPORTS]"):
			for require, _ := range requires {
				code += fmt.Sprintf("\t\"%s\"\n", strings.Replace(require, "%PACKAGE%", s.info.pkgName, -1))
			}
		case strings.Contains(string(line), "[INIT]"):
			code += "\tif err := s.dao.BeginTx(ctx); err != nil {\n\t\tpanic(err)\n\t}"
			for mdl, tbl := range models {
				str := fmt.Sprintf("s.dao.CreateTx(tx, \"%s\", reflect.TypeOf((*%s)(nil)).Elem())", tbl, mdl)
				code += fmt.Sprintf(" else if err := %s; err != nil {\n\t\tpanic(err)\n\t}", str)
			}
			code += " else if err := s.dao.CommitTx(tx); err != nil {\n\t\tpanic(err)\n\t}\n"
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
