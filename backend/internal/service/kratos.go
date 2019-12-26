package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
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
	if apis, apiTyps, err := kratos.genKratosProtoFile(ctx); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	} else if err := kratos.chgKratosConfig(); err != nil {
		return fmt.Errorf("修改配置文件失败：%v", err)
	} else if files, err := kratos.chgKratosDaoFile(ctx); err != nil {
		return fmt.Errorf("修改DAO文件失败：%v", err)
	} else if err := kratos.chgKratosServiceFile(apis); err != nil {
		return fmt.Errorf("修改Service文件失败：%v", err)
	} else if err := kratos.switchKratosMicoServ(); err != nil {
		return fmt.Errorf("开启/关闭微服务功能失败:%v", err)
	} else if err := kratos.chgKratosProjName(files); err != nil {
		return fmt.Errorf("修改项目名称失败：%v", err)
	} else if err := kratos.chgKratosTypeName(apiTyps); err != nil {
		return fmt.Errorf("修改类型名称失败：%v", err)
	}
	return nil
}

type DaoImplInfo struct {
	Def string `json:"def"`
	New string `json:"new"`
	Init string `json:"init"`
	Requires map[string]string `json:"requires"`
	Files []struct {
		Href string `json:"href"`
		Location string `json:"location"`
	} `json:"files"`
}

func (kratos *Kratos) chgKratosTypeName(apiTyps []string) error {
	// 为所有API类型添加pb.前缀
	svcPath := path.Join(kratos.info.pathName, "internal", "service", "service.go")
	svcFile, err := os.Open(svcPath)
	if err != nil {
		return fmt.Errorf("打开service.go文件失败：%v", err)
	}
	defer svcFile.Close()
	svcBytes, err := ioutil.ReadAll(svcFile)
	if err != nil {
		return fmt.Errorf("读取service.go文件失败：%v", err)
	}
	for _, tname := range apiTyps {
		if strings.Index(string(svcBytes), tname) == -1 {
			continue
		}
		pattern := regexp.MustCompile(fmt.Sprintf("\\W%s\\W", tname))
		svcBytes = pattern.ReplaceAllFunc(svcBytes, func(matched []byte) []byte {
			return []byte(strings.Replace(string(matched), tname, "pb." + tname, 1))
		})
	}
	return ioutil.WriteFile(svcPath, svcBytes, os.ModePerm)
}

func (kratos *Kratos) chgKratosDaoFile(ctx context.Context) ([]string, error) {
	daoGroups, err := kratos.svc.DaoGroupsSelectAll(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	modFiles := make(map[string]interface{})
	daoImplInfos := make([]*DaoImplInfo, 0)
	daoImplInfoType := reflect.TypeOf((*DaoImplInfo)(nil)).Elem()
	for _, group := range daoGroups.Groups {
		modInfo, err := kratos.svc.ModuleInfoSelectBySignId(ctx, &pb.StrID{Id: group.Implement})
		if err != nil {
			return nil, err
		}
		jsonMap, err := utils.HttpGetJsonMap(modInfo.DaoImplHref)
		if err != nil {
			return nil, fmt.Errorf("获取DAO实例化信息失败：%v", err)
		}
		obj, err := utils.ToObj(jsonMap, daoImplInfoType)
		if err != nil {
			return nil, fmt.Errorf("DAO实例化信息转换失败：%v", err)
		}
		daoImplInfo := obj.(*DaoImplInfo)
		// 检查是否有依赖MOD，有则添加关联文件至下载列表中
		for _, daoImplId := range modInfo.Requires {
			reqModInfo, err := kratos.svc.ModuleInfoSelectBySignId(ctx, &pb.StrID{Id: daoImplId})
			if err != nil {
				return nil, fmt.Errorf("查询依赖MOD信息失败：%v", err)
			}
			jsonMap, err = utils.HttpGetJsonMap(reqModInfo.DaoImplHref)
			if err != nil {
				return nil, fmt.Errorf("获取依赖MOD实例化信息失败：%v", err)
			}
			obj, err = utils.ToObj(jsonMap, daoImplInfoType)
			if err != nil {
				return nil, fmt.Errorf("DAO实例化信息转换失败：%v", err)
			}
			for _, reqFile := range obj.(*DaoImplInfo).Files {
				daoImplInfo.Files = append(daoImplInfo.Files, reqFile)
			}
		}
		daoImplInfos = append(daoImplInfos, daoImplInfo)

		// 生成配置文件
		jsonMap, err = utils.HttpGetJsonMap(modInfo.DaoConfHref)
		if err != nil {
			return nil, fmt.Errorf("获取DAO配置信息失败：%v", err)
		}
		// 配置文件根据文件自身自定
		daoConfPath := filepath.Join(kratos.info.pathName, jsonMap["fileName"].(string))
		daoConf, err := kratos.svc.DaoConfigSelectByImpl(ctx, &pb.DaoConfImplIden{Implement: group.Implement})
		if err != nil {
			return nil, err
		}
		strConf := ""
		for _, config := range daoConf.Configs {
			strConf += fmt.Sprintf("%s=\"%s\"\n", config.Key, config.Value)
		}
		if err := ioutil.WriteFile(daoConfPath, []byte(strConf), os.ModePerm); err != nil {
			return nil, fmt.Errorf("写入DAO配置失败：%v", err)
		}

		// 下载模块下所有文件
		for _, fileInfo := range daoImplInfo.Files {
			flPath := filepath.Join(kratos.info.pathName, fileInfo.Location)
			if err := utils.Download(fileInfo.Href, flPath); err != nil {
				return nil, fmt.Errorf("下载文件失败：%v", err)
			}
			modFiles[flPath] = nil
		}
	}

	daoPath := path.Join(kratos.info.pathName, "internal", "dao", "dao.go")
	return utils.GenMapKeys(modFiles), utils.ReplaceContentInFile(daoPath, map[string]utils.ReplaceProcFunc{
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

func (kratos *Kratos) chgKratosConfig() error {
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
	}, true); err != nil {
		return err
	}
	return nil
}

func (kratos *Kratos) adjServerFile(pathName string, regSvc string) error {
	fpath := path.Join(kratos.info.pathName, "internal", "server", pathName, "server.go")
	if err := utils.InsertTxt(fpath, func(line string, doProc *bool) (string, bool, error) {
		if !*doProc {
			return line, false, nil
		}
		if strings.Trim(strings.TrimSpace(line), "\t") == "\"strings\"" && !kratos.info.option.IsMicoServ {
			return "", false, nil
		} else if strings.Contains(line, "svr \"") && !kratos.info.option.IsMicoServ {
			return "", false, nil
		} else if strings.Contains(line, "Demo") {
			return strings.Replace(line, "Demo", utils.Capital(kratos.info.option.Name), -1), false, nil
		} else if strings.Contains(line, fmt.Sprintf("func %s", regSvc)) {
			*doProc = kratos.info.option.IsMicoServ
		} else if strings.Contains(line, regSvc) && !kratos.info.option.IsMicoServ {
			return "", false, nil
		}
		return line, false, nil
	}, true); err != nil {
		return err
	}
	return nil
}

func (kratos *Kratos) switchKratosMicoServ() error {
	// 调整cmd/main.go
	if !kratos.info.option.IsMicoServ {
		if err := utils.InsertTxt(path.Join(kratos.info.pathName, "cmd", "main.go"), func(line string, doProc *bool) (string, bool, error) {
			if strings.Contains(line, "resolver") || strings.Contains(line, "discovery") {
				return "", false, nil
			}
			return line, false, nil
		}, true); err != nil {
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
	if err := kratos.adjServerFile("grpc", "RegisterGRPCService"); err != nil {
		return err
	}
	// 调整internal/server/http/server.go的逻辑
	if err := kratos.adjServerFile("http", "RegisterHTTPService"); err != nil {
		return err
	}
	// 调整internal/service/service.go的逻辑
	if kratos.info.option.IsMicoServ {
		return nil
	}
	adjLst := []string{
		path.Join(kratos.info.pathName, "cmd", "main.go"),
		path.Join(kratos.info.pathName, "internal", "service", "service.go"),
	}
	for _, flPath := range adjLst {
		foundEnd := false
		if err := utils.InsertTxt(flPath, func(line string, doProc *bool) (string, bool, error) {
			if foundEnd {
				foundEnd = false
				*doProc = false
			}
			if strings.Index(line, "// [MICO_SERV_BEG]") != -1 {
				*doProc = true
			}
			if strings.Index(line, "// [MICO_SERV_END]") != -1 {
				foundEnd = true
			}
			return "", false, nil
		}, false); err != nil {
			return fmt.Errorf("调整service.go的为服务逻辑时发生错误：%v", err)
		}
	}
	return nil
}

func (kratos *Kratos) chgKratosProjName(addedFiles []string) error {
	fixLst := []string{
		path.Join(kratos.info.pathName, "go.mod"),
		path.Join(kratos.info.pathName, "cmd", "main.go"),
		path.Join(kratos.info.pathName, "internal", "server", "grpc", "server.go"),
		path.Join(kratos.info.pathName, "internal", "server", "http", "server.go"),
		path.Join(kratos.info.pathName, "internal", "service", "service.go"),
	}
	if addedFiles != nil && len(addedFiles) != 0 {
		fixLst = append(fixLst, addedFiles...)
	}
	if kratos.info.option.IsMicoServ {
		fixLst = append(fixLst, path.Join(kratos.info.pathName, "internal", "server", "common.go"))
	}
	for _, p := range fixLst {
		if err := utils.InsertTxt(p, func(line string, doProc *bool) (string, bool, error) {
			if strings.Contains(line, ")") {
				*doProc = false
			}
			if !*doProc {
				return line, false, nil
			}
			return strings.Replace(line, "template", kratos.info.pkgName, -1), false, nil
		}, true); err != nil {
			return err
		}
	}
	return nil
}

// 根据数据库中模型的定义，生成proto的message和service
func (kratos *Kratos) genKratosProtoFile(ctx context.Context) ([]*pb.ApiInfo, []string, error) {
	// 添加proto文件并根据数据库添加message和service
	code := "syntax = \"proto3\";\n\n"
	code += fmt.Sprintf("package %s.service.v1;\n\n", kratos.info.pkgName)
	code += "import \"gogoproto/gogo.proto\";\n"
	code += "import \"google/api/annotations.proto\";\n\n"
	code += "option go_package = \"api\";\n"
	code += "option (gogoproto.goproto_getters_all) = false;\n\n"

	// 收集接口信息
	mdls, err := kratos.svc.ModelsSelectAll(ctx, &pb.TypeIden{Type: "model"})
	if err != nil {
		return nil, nil, fmt.Errorf("查询所有模块失败：%v", err)
	}
	stts, err := kratos.svc.ModelsSelectAll(ctx, &pb.TypeIden{Type: "struct"})
	if err != nil {
		return nil, nil, fmt.Errorf("查询所有结构失败：%v", err)
	}
	mdls.Models = append(mdls.Models, stts.Models...)
	var apiTyps []string
	for _, mdl := range mdls.Models {
		apiTyps = append(apiTyps, mdl.Name)
		code += fmt.Sprintf("message %s {\n", mdl.Name)
		for i, prop := range mdl.Props {
			code += fmt.Sprintf("\t%s %s=%d;\n", prop.Type, prop.Name, i+1)
		}
		code += "}\n\n"
	}

	// 生成proto文件
	apis, err := kratos.svc.ApisSelectAll(ctx, &pb.Empty{})
	if err != nil {
		return nil, nil, err
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
		if len(api.GetHttp().Route) != 0 && len(api.GetHttp().Method) != 0 {
			fixedRoute := strings.Replace(api.GetHttp().Route, "%PROJ_NAME%", kratos.info.pkgName, -1)
			code += " {\n\t\toption (google.api.http) = {\n"
			code += fmt.Sprintf("\t\t\t%s: \"%s\"\n\t\t};\n\t};\n", api.GetHttp().Method, fixedRoute)
		} else {
			code += ";\n"
		}
	}
	code += "}\n\n"

	protoPath := path.Join(kratos.info.pathName, "api", "api.proto")
	protoFile, err := os.OpenFile(protoPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, nil, err
	}
	defer protoFile.Close()
	if _, err := protoFile.WriteString(code); err != nil {
		return nil, nil, fmt.Errorf("写入proto文件失败：%v", err)
	}

	// 生成proto文件
	cmd := exec.CommandContext(ctx, "kratos", []string{
		"tool", "protoc", protoPath,
	}...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, errors.New(err.Error() + "\n" + stderr.String())
	}
	return apis.Infos, apiTyps, nil
}

// 根据抽取的接口信息，生成完整的service
func (kratos *Kratos) chgKratosServiceFile(apis []*pb.ApiInfo) error {
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
					aparams = append(aparams, fmt.Sprintf("%s *%s", pname, ptype))
				}
				sparams := strings.Join(aparams, ", ")
				areturns := make([]string, 0)
				for _, ret := range ai.Returns {
					areturns = append(areturns, "*" + ret)
				}
				sreturns := strings.Join(areturns, ", ")
				code += fmt.Sprintf("func (s *Service) %s(ctx context.Context, %s) (%s, error) {\n", utils.Capital(ai.Name), sparams, sreturns)
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
