package service

import (
	"bufio"
	"context"
	"crypto/md5"
	"reflect"
	"time"

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

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"io"
)

// Service service.
type Service struct {
	ac *paladin.Map
	cc struct {
		Qiniu *utils.StorageConfig
	}
	dao *dao.Dao
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
	s.dao.Create(ctx, model.MODELS_TABLE, reflect.TypeOf((*pb.Model)(nil)).Elem())
	return s
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
	if apis, err := s.genKratosProtoFile(ctx, pjName, pjPath); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := s.chgKratosServiceFile(ctx, pjPath, apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	}
	return nil
}


// 所有事务流都是函数调用，而且所有的函数返回值都是err结尾
type operStep struct {
	imports []string
	proc string
	desc string
	code string
	params map[string]string
	genVars []string
}

// TODO: 以后应该迁到数据库
var operMapper = map[string]string {
	"return_one": "return %RET%, nil\n",
	"assignment": "%TARGET% = %SOURCE%\n",
	"json_marshal": "bytes, err := json.Marshal(%OBJECT%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"转JSON失败：%v\", err)\n}\n",
	"json_unmarshal": "omap, utils.UnmarshalJSON(bytes, reflect.TypeOf((*%OBJ_TYPE%)(nil)).Elem())\nif err != nil {\n\treturn nil, fmt.Errorf(\"从JSON转回失败：%v\", err)\n}\n",
	"database_beginTx": "tx, err := s.dao.BeginTx(ctx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"开启事务失败：%v\", err)\n}\n",
	"database_commitTx": "err := s.dao.CommitTx(tx)\nif err != nil {\n\treturn nil, fmt.Errorf(\"提交事务失败：%v\", err)\n}\n",
	"database_queryTx": "res, err := s.dao.QueryTx(ctx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"查询数据表失败：%v\", err)\n}\n",
	"database_deleteTx": "_, err := s.dao.DeleteTx(tx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"删除数据表记录失败：%v\", err)\n}\n",
	"database_insertTx": "id, err := s.dao.InsertTx(tx, \"%TABLE_NAME%\", %OBJ_MAP%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"插入数据表失败：%v\", err)\n}\n",
	"database_updateTx": "id, err := s.dao.SaveTx(tx, \"%TABLE_NAME%\", %QUERY_CONDS%, %QUERY_ARGUS%)\nif err != nil {\n\treturn nil, fmt.Errorf(\"更新数据表记录失败：%v\", err)\n}\n",
}

type apiInfo struct {
	model  string
	name   string
	params map[string]string
	route  string
	method string
	ret    string
	flows []operStep
}

// 根据数据库中模型的定义，生成proto的message和service
func (s *Service) genKratosProtoFile(ctx context.Context, pjName string, pjPath string) ([]apiInfo, error) {
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
	var modelApis []apiInfo
	for _, mdl := range res {
		mname := mdl["name"].(string)
		code += fmt.Sprintf("message %s {\n", mname)
		for i, prop := range mdl["props"].([]map[string]interface{}) {
			code += fmt.Sprintf("\t%s %s=%d;\n", prop["type"], prop["name"], i+1)
		}
		code += "}\n\n"
		// 复数message
		code += fmt.Sprintf("message %sArray {\n\t%s %s = 1;\n}\n\n", mname, mname, utils.ToPlural(strings.ToLower(mname)))

		if !reflect.TypeOf(mdl["methods"]).ConvertibleTo(reflect.TypeOf(([]interface{})(nil))) {
			continue
		}
		for _, method := range mdl["methods"].([]interface{}) {
			m := method.(string)
			aname, exs := actMap[m]
			if !exs {
				aname = "Select"
			}
			modelApi := apiInfo{
				model:  mname,
				name:   fmt.Sprintf("%s%s", aname, mname),
				params: make(map[string]string),
				route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				method: strings.ToLower(m),
			}
			switch modelApi.method {
			case "post":
				modelApi.params["entry"] = mname
				modelApi.ret = mname
				modelApi.flows = []operStep{
					{
						desc: "先将收到的请求参数编码成JSON字节数组",
						proc: "json_marshal",
						params: map[string]string{
							"OBJECT": "entry",
						},
					},
					{
						desc: "再将JSON字节数组转成Map键值对",
						// TODO: 自身的包要加包名前缀
						impo: []string{
							"template/internal/utils",
						},
						proc: "json_unmarshal",
						params: map[string]string{
							"OBJ_TYPE": mname,
						},
					},
					{
						desc: "开启数据库事务",
						proc: "database_beginTx",
					},
					{
						desc: "做数据库插入操作",
						proc: "database_insertTx",
						params: map[string]string{
							"TABLE_NAME": utils.CamelToPascal(mname),
							"OBJ_MAP": "omap",
						},
					},
					{
						desc: "提交数据库事务",
						proc: "database_commitTx",
					},
					{
						desc: "将改记录的数据库id赋予请求参数",
						proc: "assignment",
						params: map[string]string{
							"SOURCE": "id",
							"TARGET": "entry.Id",
						},
					},
					{
						proc: "return_one",
						params: map[string]string{
							"RET": "entry",
						},
					},
				}
			case "delete":
				modelApi.params["iden"] = "IdenReqs"
				modelApi.ret = mname
			case "put":
				modelApi.params["iden"] = "IdenReqs"
				modelApi.params["entry"] = mname
				modelApi.ret = mname
			case "get":
				modelApi.params["iden"] = "IdenReqs"
				modelApi.ret = mname
			case "all":
				modelApi.method = "get"
				modelApi.params["params"] = "Nil"
				modelApi.ret = mname + "Array"
			}
			modelApis = append(modelApis, modelApi)
		}
	}

	if len(modelApis) != 0 {
		code += fmt.Sprintf("service %s {\n", utils.Capital(pjName))
	}
	for _, api := range modelApis {
		sparams := ""
		for _, ptyp := range api.params {
			sparams += ptyp + ","
		}
		sparams = strings.TrimRight(sparams, ",")
		code += fmt.Sprintf("\trpc %s(%s) returns (%s)", api.name, sparams, api.ret)
		if len(api.route) != 0 && len(api.method) != 0 {
			code += " {\n\t\toption (google.api.http) = {\n"
			code += fmt.Sprintf("\t\t\t%s: \"%s\"\n\t\t};\n\t};\n", api.method, api.route)
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

// 根据抽取的接口信息，生成完整的service
func (s *Service) chgKratosServiceFile(ctx context.Context, pjPath string, apis []apiInfo) error {
	svcPath := path.Join(pjPath, "internal", "service", "service.go")
	svcFile, err := os.Open(svcPath)
	if err != nil {
		return err
	}
	defer svcFile.Close()
	reader := bufio.NewReader(svcFile)
	// 收集import文件
	imports := make(map[string]interface{})
	for _, a := range apis {
		for _, step := range a.flows {
			for _, i := range step.imports {
				imports[i] = nil
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
			for _, apiInfo := range apis {
				aparams := make([]string, 0)
				for pname, ptype := range apiInfo.params {
					aparams = append(aparams, fmt.Sprintf("%s *pb.%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				code += fmt.Sprintf("func (s *Service) %s(ctx context.Context, %s) (*pb.%s, error) {\n", apiInfo.name, sparams, apiInfo.ret)
				for idx, step := range apiInfo.flows {
					if idx == 0 {
						code += "\t"
					}
					// 添加注释
					if len(step.desc) != 0 {
						code += fmt.Sprintf("// %s\n", step.desc)
					}
					// 提取步骤操作的代码
					cd, exs := operMapper[step.proc]
					if !exs {
						return fmt.Errorf("无法找到指定操作步骤：%s", step.proc)
					}
					// 替换步骤中的占位符
					for o, n := range step.params {
						cd = strings.Replace(cd, fmt.Sprintf("%%%s%%", o), n, -1)
					}
					code += "\t" + utils.AddSpacesBeforeRow(cd, 1)
				}
				code += "}\n\n"
			}
		case strings.Contains(string(line), "[IMPORTS]"):
			for _, impo := range imports {
				code += "\t" + impo
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
