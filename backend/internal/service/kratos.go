package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	pb "backend/api"
	"backend/internal/utils"
)

type Kratos struct {
	info *ProjInfo
	svc  *Service
}

func NewKratos(svc *Service, pi *ProjInfo) (*Kratos, error) {
	return &Kratos{info: pi, svc: svc}, nil
}

func (kratos *Kratos) Adjust(ctx context.Context) error {
	if apis, err := kratos.genKratosProtoFile(ctx); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := kratos.chgKratosConfig(ctx); err != nil {
		return fmt.Errorf("修改配置文件失败：%v", err)
	} else if err := kratos.chgKratosDaoFile(ctx); err != nil {
		return fmt.Errorf("修改DAO文件失败：%v", err)
	} else if err := kratos.chgKratosServiceFile(ctx, apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	} else if err := kratos.switchKratosMicoServ(ctx); err != nil {
		return fmt.Errorf("开启/关闭微服务功能失败:%v", err)
	} else if err := kratos.chgKratosProjName(ctx); err != nil {
		return fmt.Errorf("修改项目名称失败：%v", err)
	}
	return nil
}

type DaoImplInfo struct {
	Def string `json:"def"`
	New string `json:"new"`
	Init string `json:"init"`
	Files []struct {
		Href string `json:"href"`
		Location string `json:"location"`
	} `json:"files"`
}

type DaoConfInfo struct {
	FileName string `json:"fileName"`
}

func (kratos *Kratos) chgKratosDaoFile(ctx context.Context) error {
	daoGroups, err := kratos.svc.DaoGroupsSelectAll(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	daoImplInfos := make([]DaoImplInfo, 0)
	for _, group := range daoGroups.Groups {
		modInfo, err := kratos.svc.ModuleInfoSelectBySignId(ctx, &pb.StrID{Id: group.Implement})
		if err != nil {
			return err
		}
		resp, err := http.Get(modInfo.DaoImplHref)
		if err != nil {
			return fmt.Errorf("获取DAO实例化信息失败：%v", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		var daoImplInfo DaoImplInfo
		if err := json.Unmarshal(body, &daoImplInfo); err != nil {
			return fmt.Errorf("解析DAO实例化信息失败：%v", err)
		}
		daoImplInfos = append(daoImplInfos, daoImplInfo)
		if err := resp.Body.Close(); err != nil {
			return fmt.Errorf("关闭响应体失败：%v", err)
		}

		// 生成配置文件
		resp, err = http.Get(modInfo.DaoConfHref)
		if err != nil {
			return fmt.Errorf("获取DAO实例化信息失败：%v", err)
		}
		body, err = ioutil.ReadAll(resp.Body)
		var daoConfInfo DaoConfInfo
		if err := json.Unmarshal(body, &daoConfInfo); err != nil {
			return fmt.Errorf("解析DAO实例化信息失败：%v", err)
		}
		if err := resp.Body.Close(); err != nil {
			return fmt.Errorf("关闭响应体失败：%v", err)
		}
		daoConfPath := filepath.Join(kratos.info.pathName, "configs", daoConfInfo.FileName)
		daoConf, err := kratos.svc.DaoConfigSelectByImpl(ctx, &pb.DaoConfImplIden{Implement: group.Implement})
		if err != nil {
			return err
		}
		strConf := ""
		for cfgKey, cfgVal := range daoConf.Configs {
			strConf += fmt.Sprintf("%s=\"%s\"\n", cfgKey, cfgVal)
		}
		if err := ioutil.WriteFile(daoConfPath, []byte(strConf), os.ModePerm); err != nil {
			return fmt.Errorf("写入DAO配置失败：%v", err)
		}


		// 下载模块下所有文件
		for _, fileInfo := range daoImplInfo.Files {
			if err := utils.Download(fileInfo.Href, filepath.Join(kratos.info.pathName, fileInfo.Location)); err != nil {
				return fmt.Errorf("下载文件失败：%v", err)
			}
		}
	}

	daoPath := path.Join(kratos.info.pathName, "internal", "dao", "dao.go")
	return utils.ReplaceContentInFile(daoPath, map[string]utils.ReplaceProcFunc{
		"[DEFINITION]": func(line string) (string, error) {
			code := ""
			for _, diInfo := range daoImplInfos {
				code += fmt.Sprintf("\t%s\n", diInfo.Def)
			}
			return code, nil
		},
		"[NEW]": func(line string) (string, error) {
			code := ""
			for _, diInfo := range daoImplInfos {
				code += fmt.Sprintf("\t\t%s,\n", diInfo.New)
			}
			return code, nil
		},
		"[INIT]": func(line string) (string, error) {
			// TODO:
			return "", nil
		},
	})
}

func (kratos *Kratos) chgKratosConfig(ctx context.Context) error {
	// 修改configs/application.toml
	appCfgPath := path.Join(kratos.info.pathName, "configs", "application.toml")
	if err := utils.InsertTxt(appCfgPath, func(line string, doProc *bool) (string, bool, error) {
		if strings.Contains(line, "appID") {
			if kratos.info.option.IsMicoServ {
				return fmt.Sprintf("appID = \"%s.service\"", kratos.info.pkgName), false, nil
			} else {
				return "", false, nil
			}
		}
		if strings.Contains(line, "swaggerFile") {
			if kratos.info.option.IsMicoServ {
				return strings.Replace(line, "template", kratos.info.option.Name, -1), false, nil
			} else {
				return "", false, nil
			}
		}
		return line, false, nil
	}); err != nil {
		return err
	}
	return nil
}

func (kratos *Kratos) adjServerFile(pathName string, regSvr string, regSvc string) error {
	fpath := path.Join(kratos.info.pathName, "internal", "server", pathName, "server.go")
	if err := utils.InsertTxt(fpath, func(line string, doProc *bool) (string, bool, error) {
		if !*doProc {
			return line, false, nil
		}
		if strings.Contains(line, "svr \"") && !kratos.info.option.IsMicoServ {
			return "", false, nil
		} else if strings.Contains(line, fmt.Sprintf("pb.%s", regSvr)) {
			regName := fmt.Sprintf("Register%sServer", utils.Capital(kratos.info.option.Name))
			return strings.Replace(line, regSvr, regName, -1), false, nil
		} else if strings.Contains(line, fmt.Sprintf("func %s", regSvc)) {
			*doProc = kratos.info.option.IsMicoServ
		} else if strings.Contains(line, regSvc) && !kratos.info.option.IsMicoServ {
			return "", false, nil
		}
		return line, false, nil
	}); err != nil {
		return err
	}
	return nil
}

func (kratos *Kratos) switchKratosMicoServ(ctx context.Context) error {
	// 调整cmd/main.go
	if !kratos.info.option.IsMicoServ {
		if err := utils.InsertTxt(path.Join(kratos.info.pathName, "cmd", "main.go"), func(line string, doProc *bool) (string, bool, error) {
			if strings.Contains(line, "resolver") || strings.Contains(line, "discovery") {
				return "", false, nil
			}
			return line, false, nil
		}); err != nil {
			return err
		}
	}
	// 删除api/register.*
	if !kratos.info.option.IsMicoServ {
		for _, p := range []string{
			path.Join(kratos.info.pathName, "api", "register.proto"),
			path.Join(kratos.info.pathName, "api", "register.bm.go"),
			path.Join(kratos.info.pathName, "api", "register.pb.go"),
			path.Join(kratos.info.pathName, "api", "register.swagger.json"),
		} {
			if err := os.Remove(p); err != nil {
				return err
			}
		}
	}
	// 删除internal/server/common.go
	if !kratos.info.option.IsMicoServ {
		if err := os.Remove(path.Join(kratos.info.pathName, "internal", "server", "common.go")); err != nil {
			return err
		}
	}
	// 调整internal/server/grpc/server.go的逻辑
	if err := kratos.adjServerFile("grpc", "RegisterDemoServer", "RegisterGRPCService"); err != nil {
		return err
	}
	// 调整internal/server/http/server.go的逻辑
	if err := kratos.adjServerFile("http", "RegisterDemoBMServer", "RegisterHTTPService"); err != nil {
		return err
	}
	return nil
}

func (kratos *Kratos) chgKratosProjName(ctx context.Context) error {
	fixLst := []string{
		path.Join(kratos.info.pathName, "cmd", "main.go"),
		path.Join(kratos.info.pathName, "internal", "server", "grpc", "server.go"),
		path.Join(kratos.info.pathName, "internal", "server", "http", "server.go"),
	}
	if kratos.info.option.IsMicoServ {
		fixLst = append(fixLst, path.Join(kratos.info.pathName, "internal", "server", "common.go"))
	}
	impPkg := fmt.Sprintf("\"%s/", kratos.info.pkgName)
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

// 根据数据库中模型的定义，生成proto的message和service
func (kratos *Kratos) genKratosProtoFile(ctx context.Context) ([]*pb.ApiInfo, error) {
	// 添加proto文件并根据数据库添加message和service
	code := "syntax = \"proto3\";\n\n"
	code += fmt.Sprintf("package %s.service.v1;\n\n", kratos.info.pkgName)
	code += "import \"gogoproto/gogo.proto\";\n"
	code += "import \"google/api/annotationkpg.proto\";\n\n"
	code += "option go_package = \"api\";\n"
	code += "option (gogoproto.goproto_getters_all) = false;\n\n"

	// 收集接口信息
	mdls, err := kratos.svc.ModelsSelectAll(ctx, &pb.TypeIden{Type: "model"})
	if err != nil {
		return nil, err
	}
	for _, mdl := range mdls.Models {
		code += fmt.Sprintf("message %s {\n", mdl.Name)
		for i, prop := range mdl.Props {
			code += fmt.Sprintf("\t%s %s=%d;\n", prop.Type, prop.Name, i+1)
		}
		code += "}\n\n"
	}

	// 生成proto文件
	apis, err := kratos.svc.ApisSelectAll(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	if len(apis.Infos) != 0 {
		code += fmt.Sprintf("service %s {\n", utils.Capital(kratos.info.option.Name))
	}
	for _, api := range apis.Infos {
		sparams := ""
		for _, ptyp := range api.Params {
			sparams += ptyp + ","
		}
		sparams = strings.TrimRight(sparams, ",")
		code += fmt.Sprintf("\trpc %s(%s) returns (%s)", api.Name, sparams, strings.Join(api.Returns, ","))
		if len(api.Route) != 0 && len(api.Method) != 0 {
			fixedRoute := strings.Replace(api.Route, "%PROJ_NAME%", kratos.info.pkgName, -1)
			code += " {\n\t\toption (google.api.http) = {\n"
			code += fmt.Sprintf("\t\t\t%s: \"%s\"\n\t\t};\n\t};\n", api.Method, fixedRoute)
		} else {
			code += ";\n"
		}
	}
	code += "}\n\n"

	protoPath := path.Join(kratos.info.pathName, "api", "api.proto")
	protoFile, err := os.OpenFile(protoPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer protoFile.Close()
	if _, err := protoFile.WriteString(code); err != nil {
		return nil, fmt.Errorf("写入proto文件失败：%v", err)
	}
	return apis.Infos, nil
}

// 根据抽取的接口信息，生成完整的service
func (kratos *Kratos) chgKratosServiceFile(ctx context.Context, apis []*pb.ApiInfo) error {
	svcPath := path.Join(kratos.info.pathName, "internal", "service", "service.go")
	// 收集import文件
	requires := make(map[string]interface{})
	models := make(map[string]string)
	for _, ai := range apis {
		models[ai.Model] = ai.Table
		for _, step := range ai.Steps {
			for _, i := range step.Requires {
				requires[i] = nil
			}
		}
	}
	return utils.ReplaceContentInFile(svcPath, map[string]utils.ReplaceProcFunc{
		"\"template/": func(line string) (string, error) {
			return strings.Replace(line, "\"template/", fmt.Sprintf("\"%s/", kratos.info.pkgName), -1) + "\n", nil
		},
		"[APIS]": func(line string) (string, error) {
			code := ""
			for _, ai := range apis {
				aparams := make([]string, 0)
				for pname, ptype := range ai.Params {
					// 参数都是api包下的类型，所以要附上pb前缀
					aparams = append(aparams, fmt.Sprintf("%s *pb.%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				areturns := make([]string, 0)
				for _, ret := range ai.Returns {
					areturns = append(areturns, "*pb." + ret)
				}
				sreturns := strings.Join(areturns, ", ")
				code += fmt.Sprintf("func (s *Service) %s(ctx context.Context, %s) (%s, error) {\n", ai.Name, sparams, sreturns)
				preSpaces := 1
				for _, step := range ai.Steps {
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
					if step.Symbol & pb.SpcSymbol_BLOCK_IN != 0 {
						cd = cd + " {\n"
						code += utils.AddSpacesBeforeRow(cd, preSpaces)
						preSpaces++
						continue
					} else if step.Symbol & pb.SpcSymbol_BLOCK_OUT != 0 {
						cd = "}\n" + cd
						preSpaces--
					}
					code += utils.AddSpacesBeforeRow(cd, preSpaces)
				}
				code += "}\n\n"
			}
			return code, nil
		},
		"[IMPORTS]": func(line string) (string, error) {
			code := ""
			for require := range requires {
				code += fmt.Sprintf("\t\"%s\"\n", strings.Replace(require, "%PACKAGE%", kratos.info.pkgName, -1))
			}
			return code, nil
		},
		"[INIT]": func(line string) (string, error) {
			// TODO:
			return "", nil
		},
	})
}
