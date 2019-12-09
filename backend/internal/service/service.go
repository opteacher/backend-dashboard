package service

import (
	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/server"
	"backend/internal/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path"
	"path/filepath"
	"reflect"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"go.mongodb.org/mongo-driver/bson"
)

// Service service.
type Service struct {
	ac *paladin.Map
	cc struct {
		Qiniu *utils.StorageConfig
	}
	mongo *dao.Dao
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		mongo: dao.New(),
	}
	if err := paladin.Get("cdn.toml").UnmarshalTOML(&s.cc); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := s.setupModel(ctx); err != nil {
		panic(err)
	}
	if err := s.setupApiInfo(ctx); err != nil {
		panic(err)
	}
	if err := s.setupDaoGroup(ctx); err != nil {
		panic(err)
	}
	return s
}

func (s *Service) setupModel(ctx context.Context) error {
	// 创建模型集合
	if err := s.mongo.Create(ctx, model.MODELS_TABLE); err != nil {
		return err
	}
	projPath, err := s.ac.Get("projPath").String()
	if err != nil {
		return err
	}
	if err := s.mongo.Source(ctx, filepath.Join(projPath, "backend", "datas", model.MODELS_TABLE + ".json")); err != nil {
		return err
	}
	// 创建链接集合
	if err := s.mongo.Create(ctx, model.LINKS_TABLE); err != nil {
		return err
	}
	return nil
}

func (s *Service) setupApiInfo(ctx context.Context) error {
	if err := s.mongo.Create(ctx, model.OPER_STEP_TABLE); err != nil {
		return err
	}
	if err := s.mongo.Create(ctx, model.API_INFO_TABLE); err != nil {
		return err
	}
	projPath, err := s.ac.Get("projPath").String()
	if err != nil {
		return err
	}
	if err := s.mongo.Source(ctx, filepath.Join(projPath, "backend", "datas", model.OPER_STEP_TABLE + ".json")); err != nil {
		return err
	}
	if _, err := s.mongo.Insert(ctx, model.API_INFO_TABLE, pb.ApiInfo{
		Name: "SelectOne%MODEL%",
		Model: "%MODEL%",
		Table: "%TABLE_NAME%",
		Route: "/api/v1/%MODEL%.SelectOne",
		Method: "GET",
		Params: map[string]string{
			"req": "StrIdenReqs",
		},
		Returns:[]string{"*%MODEL%"},
		Steps: []*pb.OperStep{{
			ApiName: "SelectOne%MODEL%",
			OperKey: "database_query",
			Inputs: map[string]string{
				"TABLE_NAME": "%TABLE_NAME%",
				"CONDITIONS": "bson.D{\"_id\": req.Id}",
			},
		}, {
			ApiName: "SelectOne%MODEL%",
			OperKey: "return_succeed",
			Inputs: map[string]string{"RETURN": "res.(*%MODEL%)"},
		}},
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) setupDaoGroup(ctx context.Context) error {
	if err := s.mongo.Create(ctx, model.DAO_GROUPS_TABLE); err != nil {
		return err
	}
	if _, err := s.mongo.Insert(ctx, model.DAO_GROUPS_TABLE, pb.DaoGroup{
		Name:                 "MySQL数据库",
		Category:             "databases",
		Language:             "golang",
		Interfaces:           []*pb.DaoInterface{},
	}); err != nil {
		return err
	}
	return nil
}

func (s *Service) AppID() string {
	appID, _ := s.ac.Get("appID").String()
	return appID
}

func (s *Service) SwaggerFile() string {
	pjPath, _ := s.ac.Get("projPath").String()
	swagger, _ := s.ac.Get("swaggerFile").String()
	return path.Join(pjPath, swagger)
}

func (s *Service) ModelsInsert(ctx context.Context, req *pb.Model) (*pb.Model, error) {
	id, err := s.mongo.Insert(ctx, model.MODELS_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	}
	req.Id = id
	return req, nil
}

func (s *Service) ModelsDelete(ctx context.Context, req *pb.NameID) (*pb.Model, error) {
	conds := bson.D{{"name", req.Name }}

	res, err := s.mongo.QueryOne(ctx, model.MODELS_TABLE, conds)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	}

	if _, err := s.mongo.Delete(ctx, model.MODELS_TABLE, conds); err != nil {
		return nil, fmt.Errorf("删除数据库记录失败：%v", err)
	}

	resp, err := utils.ToObj(res, reflect.TypeOf((*pb.Model)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转Model对象失败：%v", err)
	}
	return resp.(*pb.Model), nil
}

func (s *Service) ModelsUpdate(ctx context.Context, req *pb.Model) (*pb.Empty, error) {
	conds := bson.D{{}}
	if len(req.Id) != 0 {
		oid, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, fmt.Errorf("错误的行id：%v", err)
		}
		conds = bson.D{{"_id", oid}}
	} else if len(req.Name) != 0 {
		conds = bson.D{{"name", req.Name}}
	}

	// NOTE: 只能更新x, y, width, height
	entry := bson.D{{"$set", bson.D{
		{"x", req.X},
		{"y", req.Y},
		{"width", req.Width},
		{"height", req.Height},
	}}}

	if _, err := s.mongo.Update(ctx, model.MODELS_TABLE, conds, entry); err != nil {
		return nil, fmt.Errorf("更新数据库记录失败：%v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Service) ModelsSelectAll(ctx context.Context, req *pb.TypeIden) (*pb.ModelArray, error) {
	conds := bson.D{{}}
	if len(req.Type) != 0 {
		conds = bson.D{{"type", req.Type}}
	}

	ress, err := s.mongo.Query(ctx, model.MODELS_TABLE, conds)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	}

	resp := new(pb.ModelArray)
	for _, mobj := range ress {
		obj, err := utils.ToObj(mobj, reflect.TypeOf((*pb.Model)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转为Model类型时发生错误：%v", err)
		}
		resp.Models = append(resp.Models, obj.(*pb.Model))
	}
	return resp, nil
}

func (s *Service) StructsSelectAllBases(context.Context, *pb.Empty) (*pb.NameArray, error) {
	return &pb.NameArray{Names: []string{"Nil", "IdenReqs", "NameReqs"}}, nil
}

func (s *Service) LinksInsert(ctx context.Context, req *pb.Link) (*pb.Link, error) {
	id, err := s.mongo.Insert(ctx, model.LINKS_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	}
	req.Id = id
	return req, nil
}

func (s *Service) LinksSelectAll(ctx context.Context, req *pb.Empty) (*pb.LinkArray, error) {
	ress, err := s.mongo.Query(ctx, model.LINKS_TABLE, bson.D{{}})
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	}

	resp := new(pb.LinkArray)
	for _, mobj := range ress {
		obj, err := utils.ToObj(mobj, reflect.TypeOf((*pb.Link)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转为Link类型时发生错误：%v", err)
		}
		resp.Links = append(resp.Links, obj.(*pb.Link))
	}
	return resp, nil
}

func (s *Service) LinksDeleteBySymbol(ctx context.Context, req *pb.SymbolID) (*pb.Link, error) {
	conds := bson.D{{"symbol", req.Symbol }}

	res, err := s.mongo.QueryOne(ctx, model.LINKS_TABLE, conds)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	}

	if _, err := s.mongo.Delete(ctx, model.LINKS_TABLE, conds); err != nil {
		return nil, fmt.Errorf("删除数据库记录失败：%v", err)
	}

	resp, err := utils.ToObj(res, reflect.TypeOf((*pb.Link)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转Link对象失败：%v", err)
	}
	return resp.(*pb.Link), nil
}

func (s *Service) ApisSelectByName(ctx context.Context, req *pb.NameID) (*pb.ApiInfo, error) {
	res, err := s.mongo.QueryOne(ctx, model.API_INFO_TABLE, bson.D{{"name", req.Name}})
	if err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	}

	resp, err := utils.ToObj(res, reflect.TypeOf((*pb.ApiInfo)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转ApiInfo对象失败：%v", err)
	}

	// 调整步骤的序列号
	apiInfo := resp.(*pb.ApiInfo)
	for idx := range apiInfo.Steps {
		apiInfo.Steps[idx].Index = int32(idx)
	}
	return apiInfo, nil
}

func (s *Service) ApisSelectAll(ctx context.Context, req *pb.Empty) (*pb.ApiInfoArray, error) {
	ress, err := s.mongo.Query(ctx, model.API_INFO_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	}

	resp := new(pb.ApiInfoArray)
	for _, res := range ress {
		api, err := utils.ToObj(res, reflect.TypeOf((*pb.ApiInfo)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转ApiInfo对象失败：%v", err)
		}

		// 调整步骤的序列号
		apiInfo :=  api.(*pb.ApiInfo)
		for idx := range apiInfo.Steps {
			apiInfo.Steps[idx].Index = int32(idx)
		}

		resp.Infos = append(resp.Infos, apiInfo)
	}
	return resp, nil
}

func (s *Service) ApisInsert(ctx context.Context, req *pb.ApiInfo) (*pb.ApiInfo, error) {
	_, err := s.mongo.Insert(ctx, model.API_INFO_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入接口信息失败：%v", err)
	}
	return req, nil
}

func (s *Service) ApisDeleteByName(ctx context.Context, req *pb.NameID) (*pb.ApiInfo, error) {
	res, err := s.ApisSelectByName(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = s.mongo.Delete(ctx, model.API_INFO_TABLE, bson.D{{"name", req.Name}})
	if err != nil {
		return nil, fmt.Errorf("删除接口信息失败：%v", err)
	}
	return res, nil
}

func (s *Service) StepsInsert(ctx context.Context, req *pb.StepReqs) (*pb.Empty, error) {
	apiName := req.OperStep.ApiName
	apiInfo, err := s.ApisSelectByName(ctx, &pb.NameID{Name: apiName})
	if err != nil {
		return nil, err
	}

	rear := append([]*pb.OperStep{}, apiInfo.Steps[req.Index:]...)
	apiInfo.Steps = append(append(apiInfo.Steps[:req.Index], req.OperStep), rear...)
	_, err = s.mongo.Update(ctx, model.API_INFO_TABLE, bson.D{{"name", apiName}}, bson.D{{
		"$set", bson.D{{"steps", apiInfo.Steps}},
	}})
	if err != nil {
		return nil, fmt.Errorf("插入步骤失败：%v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Service) StepsDelete(ctx context.Context, req *pb.DelStepReqs) (*pb.Empty, error) {
	apiName := req.ApiName
	apiInfo, err := s.ApisSelectByName(ctx, &pb.NameID{Name: apiName})
	if err != nil {
		return nil, err
	}

	apiInfo.Steps = append(apiInfo.Steps[:req.StepId], apiInfo.Steps[req.StepId + 1:]...)
	_, err = s.mongo.Update(ctx, model.API_INFO_TABLE, bson.D{{"name", apiName}}, bson.D{{
		"$set", bson.D{{"steps", apiInfo.Steps}},
	}})
	if err != nil {
		return nil, fmt.Errorf("删除步骤失败：%v", err)
	}
	return &pb.Empty{}, nil
}

// 这是添加步骤模板，可以通过设置apiName来指定要插入的接口，但只能追加到api流程的最后
// 如果需要插入到流程中间，则需要使用StepsInsert
func (s *Service) OperStepsInsert(context.Context, *pb.OperStep) (*pb.OperStep, error) {
	// TODO:
	return nil, nil
}

func (s *Service) OperStepsSelectTemp(ctx context.Context, req *pb.Empty) (*pb.OperStepArray, error) {
	// TODO:
	return nil, nil
}

func (s *Service) DaoGroupsSelectAll(ctx context.Context, req *pb.Empty) (*pb.DaoGroupArray, error) {
	ress, err := s.mongo.Query(ctx, model.DAO_GROUPS_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询DAO组失败：%v", err)
	}

	resp := new(pb.DaoGroupArray)
	for _, res := range ress {
		daoGroup, err := utils.ToObj(res, reflect.TypeOf((*pb.DaoGroup)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转DaoGroup对象失败：%v", err)
		}
		resp.Groups = append(resp.Groups, daoGroup.(*pb.DaoGroup))
	}
	return resp, nil
}

func (s *Service) DaoGroupsInsert(ctx context.Context, req *pb.DaoGroup) (*pb.DaoGroup, error) {
	_, err := s.mongo.Insert(ctx, model.DAO_GROUPS_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入DAO组失败：%v", err)
	}
	return req, nil
}

func (s *Service) Export(ctx context.Context, req *pb.ExpOptions) (*pb.UrlResp, error) {
	// TODO:
	return nil, nil
}

func (s *Service) SpecialSymbols(context.Context, *pb.Empty) (*pb.SymbolsResp, error) {
	return &pb.SymbolsResp{
		Values: pb.SpcSymbol_value,
		Names:  pb.SpcSymbol_name,
	}, nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.mongo.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.mongo.Close()
	// 注销服务
	if cli, err := server.RegisterService(); err != nil {
		panic(err)
	} else if _, err := cli.Cancel(context.Background(), &pb.IdenSvcReqs{AppID: s.AppID()}); err != nil {
		panic(err)
	}
}
