package service

import (
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
	resp := new(pb.Model)
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
	} else if utils.FillWithMap(resp, res[0]); false {
		return nil, nil
	} else {
		return resp, nil
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
			mdl := new(pb.Model)
			utils.FillWithMap(mdl, entry)
			resp.Models = append(resp.Models, mdl)
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
	if err := s.genKratosProtoFile(ctx, pjName, pjPath); err != nil {
		return fmt.Errorf("生成Proto文件失败：%v", err)
	}
}

type funcInfo struct {
	name string
}

func (s *Service) genKratosProtoFile(ctx context.Context, pjName string, pjPath string) error {
	if pjName[len(pjName) - 4:] == ".zip" || pjName[len(pjName) - 4:] == ".ZIP" {
		pjName = pjName[:len(pjName) - 4]
	}
	pkgName := utils.CamelToPascal(pjName)
	// 添加proto文件并根据数据库添加message和service
	protoData := "syntax = \"proto3\";\n\n"
	protoData += fmt.Sprintf("package %s.service.v1;\n\n", pkgName)
	protoData += "import \"gogoproto/gogo.proto\";\n"
	protoData += "import \"google/api/annotations.proto\";\n\n"
	protoData += "option go_package = \"api\";\n"
	protoData += "option (gogoproto.goproto_getters_all) = false;\n\n"

	res, err := s.dao.Query(ctx, model.MODELS_TABLE, "", []interface{}{})
	if err != nil {
		return err
	}
	type HttpAPI struct {
		Model string
		Func string
		Path string
		Method  string
	}
	actMap := map[string]string {
		"POST": "Insert",
		"DELETE": "Delete",
		"PUT": "Update",
		"GET": "Select",
		"ALL": "SelectAll"
	}
	var modelApis []HttpAPI
	for _, mdl := range res {
		mname := mdl["name"].(string)
		protoData += fmt.Sprintf("message %s {\n", mname)
		for i, prop := range mdl["props"].([]map[string]interface{}) {
			protoData += fmt.Sprintf("\t%s %s=%d;\n", prop["type"], prop["name"], i+1)
		}
		protoData += "}\n\n"

		for _, method := range mdl["methods"].([]map[string]interface{}) {
			m := method["method"].(string)
			aname, exs := actMap[m]
			if !exs {
				aname = "Select"
			}
			modelApis = append(modelApis, HttpAPI{
				Model: mname,
				Func: fmt.Sprintf("%s%s", aname, utils.Capital(mname)),
				Path: fmt.Sprintf("/api/v1/%s.%s", strings.ToLower(mname), strings.ToLower(aname)),
				Method: strings.ToLower(m),
			})
		}
	}

	if len(modelApis) != 0 {
		protoData += fmt.Sprintf("service %s {\n", utils.Capital(pjName))
	}
	for _, api := range modelApis {
		protoData += fmt.Sprintf("\trpc %s(%s) returns (%s) {\n", api.Func, api.Model, api.Model)
		protoData += "\t\toption (google.api.http) = {\n"
		protoData += fmt.Sprintf("\t\t\t%s: \"%s\"\n\t\t};\n\t};\n", api.Method, api.Path)
	}
	protoData += "}\n\n"

	protoPath := path.Join(pjPath, "api", "api.proto")
	protoFile, err := os.OpenFile(protoPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer protoFile.Close()
	protoFile.WriteString(protoData)
	return nil
}

func (s *Service) genKratosServiceFile(ctx context.Context, pjName string, pjPath string) error {
	code := ""
	servicePath := path.Join(pjPath, "internal", "service", "service.go")
	serviceFile, err := os.OpenFile(servicePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer serviceFile.Close()
	serviceFile.WriteString(code)
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
