package service

import (
	"context"
	"path"
	"strings"
	"os"
	"reflect"
	"fmt"
	"bufio"
	"io"

	pb "backend/api"
	"backend/internal/utils"
	"backend/internal/dao"
	"backend/internal/model"
)

type KratosProjGen struct {
	operSteps []*pb.OperStep
	info *GenInfo
	dao *dao.Dao
}

func NewKratosProjGenerator(dao *dao.Dao, gi *GenInfo) (*KratosProjGen, error) {
	return &KratosProjGen{
		info: gi,
		dao: dao,
	}, nil
}

func (kpg *KratosProjGen) Adjust(ctx context.Context) error {
	if err := kpg.readOperFromDB(ctx); err != nil {
		return fmt.Errorf("读取服务流程项目失败：%v", err)
	} else if apis, err := kpg.genKratosProtoFile(ctx); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := kpg.chgKratosConfig(ctx); err != nil {
		return fmt.Errorf("修改配置文件失败：%v", err)
	} else if err := kpg.chgKratosServiceFile(ctx, apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	} else if err := kpg.switchKratosMicoServ(ctx); err != nil {
		return fmt.Errorf("开启/关闭微服务功能失败:%v", err)
	} else if err := kpg.chgKratosProjName(ctx); err != nil {
		return fmt.Errorf("修改项目名称失败：%v", err)
	}
	return nil
}

func (kpg *KratosProjGen) chgKratosConfig(ctx context.Context) error {
	// 修改configs/application.toml
	appCfgPath := path.Join(kpg.info.pathName, "configs", "application.toml")
	if err := utils.InsertTxt(appCfgPath, func(line string, doProc *bool) (string, bool, error) {
		if strings.Contains(line, "appID") {
			if kpg.info.option.IsMicoServ {
				return fmt.Sprintf("appID = \"%s.service\"", kpg.info.pkgName), false, nil
			} else {
				return "", false, nil
			}
		}
		if strings.Contains(line, "swaggerFile") {
			if kpg.info.option.IsMicoServ {
				return strings.Replace(line, "template", kpg.info.option.Name, -1), false, nil
			} else {
				return "", false, nil
			}
		}
		return line, false, nil
	}); err != nil {
		return err
	}
	// 修改configs/mysql.toml
	if len(kpg.info.option.Database.Type) != 0 {
		switch kpg.info.option.Database.Type {
		case "mysql":
			mysqlCfgPath := path.Join(kpg.info.pathName, "configs", "mysql.toml")
			if err := utils.InsertTxt(mysqlCfgPath, func(line string, doProc *bool) (string, bool, error) {
				if strings.Contains(line, "{user}") {
					line = strings.Replace(line, "{user}", kpg.info.option.Database.Username, -1)
				}
				if strings.Contains(line, "{password}") {
					line = strings.Replace(line, "{password}", kpg.info.option.Database.Password, -1)
				}
				if strings.Contains(line, "{host}") {
					line = strings.Replace(line, "{host}", kpg.info.option.Database.Host, -1)
				}
				if strings.Contains(line, "{port}") {
					line = strings.Replace(line, "{port}", kpg.info.option.Database.Port, -1)
				}
				if strings.Contains(line, "{database}") {
					line = strings.Replace(line, "{database}", kpg.info.option.Database.Name, -1)
				}
				return line, false, nil
			}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("目前暂时不支持出MySQL以外的数据源系统，%s", kpg.info.option.Database.Type)
		}
	}
	return nil
}

func (kpg *KratosProjGen) adjServerFile(pathName string, regSvr string, regSvc string) error {
	fpath := path.Join(kpg.info.pathName, "internal", "server", pathName, "server.go")
	if err := utils.InsertTxt(fpath, func(line string, doProc *bool) (string, bool, error) {
		if !*doProc {
			return line, false, nil
		}
		if strings.Contains(line, "svr \"") && !kpg.info.option.IsMicoServ {
			return "", false, nil
		} else if strings.Contains(line, fmt.Sprintf("pb.%s", regSvr)) {
			regName := fmt.Sprintf("Register%sServer", utils.Capital(kpg.info.option.Name))
			return strings.Replace(line, regSvr, regName, -1), false, nil
		} else if strings.Contains(line, fmt.Sprintf("func %s", regSvc)) {
			*doProc = kpg.info.option.IsMicoServ
		} else if strings.Contains(line, regSvc) && !kpg.info.option.IsMicoServ {
			return "", false, nil
		}
		return line, false, nil
	}); err != nil {
		return err
	}
	return nil
}

func (kpg *KratosProjGen) switchKratosMicoServ(ctx context.Context) error {
	// 调整cmd/main.go
	if !kpg.info.option.IsMicoServ {
		if err := utils.InsertTxt(path.Join(kpg.info.pathName, "cmd", "main.go"), func(line string, doProc *bool) (string, bool, error) {
			if strings.Contains(line, "resolver") || strings.Contains(line, "discovery") {
				return "", false, nil
			}
			return line, false, nil
		}); err != nil {
			return err
		}
	}
	// 删除api/register.*
	if !kpg.info.option.IsMicoServ {
		for _, p := range []string{
			path.Join(kpg.info.pathName, "api", "register.proto"),
			path.Join(kpg.info.pathName, "api", "register.bm.go"),
			path.Join(kpg.info.pathName, "api", "register.pb.go"),
			path.Join(kpg.info.pathName, "api", "register.swagger.json"),
		} {
			if err := os.Remove(p); err != nil {
				return err
			}
		}
	}
	// 删除internal/server/common.go
	if !kpg.info.option.IsMicoServ {
		if err := os.Remove(path.Join(kpg.info.pathName, "internal", "server", "common.go")); err != nil {
			return err
		}
	}
	// 调整internal/server/grpc/server.go的逻辑
	if err := kpg.adjServerFile("grpc", "RegisterDemoServer", "RegisterGRPCService"); err != nil {
		return err
	}
	// 调整internal/server/http/server.go的逻辑
	if err := kpg.adjServerFile("http", "RegisterDemoBMServer", "RegisterHTTPService"); err != nil {
		return err
	}
	return nil
}

func (kpg *KratosProjGen) chgKratosProjName(ctx context.Context) error {
	fixLst := []string{
		path.Join(kpg.info.pathName, "internal", "dao", "dao.go"),
		path.Join(kpg.info.pathName, "internal", "server", "grpc", "server.go"),
		path.Join(kpg.info.pathName, "internal", "server", "http", "server.go"),
	}
	if kpg.info.option.IsMicoServ {
		fixLst = append(fixLst, path.Join(kpg.info.pathName, "internal", "server", "common.go"))
	}
	impPkg := fmt.Sprintf("\"%s/", kpg.info.pkgName)
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

func (kpg *KratosProjGen) readOperFromDB(ctx context.Context) error {
	if amap, err := kpg.dao.Query(ctx, model.OPER_STEP_TABLE, "", nil); err != nil {
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
			kpg.operSteps = append(kpg.operSteps, step)
		}
		return nil
	}
}

// 根据数据库中模型的定义，生成proto的message和service
func (kpg *KratosProjGen) genKratosProtoFile(ctx context.Context) ([]*pb.ApiInfo, error) {
	// 添加proto文件并根据数据库添加message和service
	code := "syntax = \"proto3\";\n\n"
	code += fmt.Sprintf("package %kpg.service.v1;\n\n", kpg.info.pkgName)
	code += "import \"gogoproto/gogo.proto\";\n"
	code += "import \"google/api/annotationkpg.proto\";\n\n"
	code += "option go_package = \"api\";\n"
	code += "option (gogoproto.goproto_getters_all) = false;\n\n"
	code += "message Nil {\n}\n\n"
	code += "message IdenReqs {\n\tint64 id = 1;\n}\n\n"

	res, err := kpg.dao.Query(ctx, model.MODELS_TABLE, "", []interface{}{})
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
	// 收集接口信息
	tx, err := kpg.dao.BeginTx(ctx)
	if err != nil {
		return nil, err
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
		mmfname := utils.ToPlural(strings.ToLower(mname))
		code += fmt.Sprintf("message %s {\n\t%s %s = 1;\n}\n\n", mmname, mname, mmfname)

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
				Table:  utils.ToPlural(utils.CamelToPascal(mname)),
				Params: make(map[string]string),
				Route:  fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				Method: strings.ToLower(m),
			}
			kpg.genModelApiInfo(modelApi, mname, mmname, mmfname)
			modelApis = append(modelApis, modelApi)
			// 接口信息存入数据库
			ApiInfoToDB(kpg.dao, tx, modelApi)
		}
	}
	if err := kpg.dao.CommitTx(tx); err != nil {
		return nil, err
	}

	// 生成proto文件
	if len(modelApis) != 0 {
		code += fmt.Sprintf("service %s {\n", utils.Capital(kpg.info.option.Name))
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

	protoPath := path.Join(kpg.info.pathName, "api", "api.proto")
	protoFile, err := os.OpenFile(protoPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer protoFile.Close()
	protoFile.WriteString(code)
	return modelApis, nil
}

func (kpg *KratosProjGen) genModelApiInfo(modelApi *pb.ApiInfo, mname string, mmname string, mmfname string) {
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

func (kpg *KratosProjGen) copyStep(operKey string, values map[string]interface{}) *pb.OperStep {
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

// 根据抽取的接口信息，生成完整的service
func (kpg *KratosProjGen) chgKratosServiceFile(ctx context.Context, apis []*pb.ApiInfo) error {
	svcPath := path.Join(kpg.info.pathName, "internal", "service", "service.go")
	svcFile, err := os.Open(svcPath)
	if err != nil {
		return err
	}
	defer svcFile.Close()
	reader := bufio.NewReader(svcFile)
	// 收集import文件
	requires := make(map[string]interface{})
	models := make(map[string]string)
	for _, ai := range apis {
		models[ai.Model] = ai.Table
		for _, step := range ai.Flows {
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
			code += strings.Replace(string(line), "\"template/", fmt.Sprintf("\"%s/", kpg.info.pkgName), -1) + "\n"
		case strings.Contains(string(line), "[APIS]"):
			for _, ai := range apis {
				aparams := make([]string, 0)
				for pname, ptype := range ai.Params {
					aparams = append(aparams, fmt.Sprintf("%s *%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				code += fmt.Sprintf("func (kpg *KratosProjGen) %s(ctx context.Context, %s) (*%s, error) {\n", ai.Name, sparams, ai.Return)
				preSpaces := 1
				for _, step := range ai.Flows {
					// 添加注释
					if len(step.Desc) != 0 {
						code += utils.AddSpacesBeforeRow(fmt.Sprintf("// %s\n", step.Desc), preSpaces)
					}
					// 提取步骤操作的代码
					cd := step.Code
					// 替换步骤中的占位符
					for o, n := range step.Inputs {
						cd = strings.Replace(cd, fmt.Sprintf("%%%s%%", o), n, -1)
					}
					// 跳进或者跳出块段落
					if step.BlockIn {
						cd = cd + " {\n"
						code += utils.AddSpacesBeforeRow(cd, preSpaces)
						preSpaces++
						continue
					} else if step.BlockOut {
						cd = "}\n" + cd
						preSpaces--
					}
					code += utils.AddSpacesBeforeRow(cd, preSpaces)
				}
				code += "}\n\n"
			}
		case strings.Contains(string(line), "[IMPORTS]"):
			for require, _ := range requires {
				code += fmt.Sprintf("\t\"%s\"\n", strings.Replace(require, "%PACKAGE%", kpg.info.pkgName, -1))
			}
		case strings.Contains(string(line), "[INIT]"):
			code += "\tif err := kpg.dao.BeginTx(ctx); err != nil {\n\t\tpanic(err)\n\t}"
			for mdl, tbl := range models {
				str := fmt.Sprintf("kpg.dao.CreateTx(tx, \"%s\", reflect.TypeOf((*%s)(nil)).Elem())", tbl, mdl)
				code += fmt.Sprintf(" else if err := %s; err != nil {\n\t\tpanic(err)\n\t}", str)
			}
			code += " else if err := kpg.dao.CommitTx(tx); err != nil {\n\t\tpanic(err)\n\t}\n"
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