package service

import (
	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/utils"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

// Service service.
type Service struct {
	ac *paladin.Map
	cc struct {
		Qiniu *utils.StorageConfig
	}
	mongo *dao.MongoDao
	pBuilder *ProjBuilder
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:       ac,
		mongo:    dao.New(),
		pBuilder: NewProjBuilder(),
	}
	if err := paladin.Get("cdn.toml").UnmarshalTOML(&s.cc); err != nil {
		panic(err)
	}
	return s
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

func (s *Service) ModelsInsertMany(ctx context.Context, req *pb.ModelArray) (*pb.ModelArray, error) {
	var entries []interface{}
	for _, mdl := range req.Models {
		entries = append(entries, mdl)
	}
	if _, err := s.mongo.InsertMany(ctx, model.MODELS_TABLE, entries); err != nil {
		return nil, fmt.Errorf("批量插入模块失败：%v", err)
	}
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

func (s *Service) LinksSelectAll(ctx context.Context, _ *pb.Empty) (*pb.LinkArray, error) {
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

func (s *Service) complApiSteps(ctx context.Context, minfo map[string]interface{}) (*pb.ApiInfo, error) {
	// 获取模板步骤作为API步骤
	tempSteps := make(map[string]*pb.Step)
	jsonSteps := minfo["steps"].([]interface{})
	for _, obj := range jsonSteps {
		mstep := obj.(map[string]interface{})
		tempSteps[mstep["key"].(string)] = nil
	}
	for stepName := range tempSteps {
		step, err := s.TempStepsSelectByKey(ctx, &pb.StrKey{Key: stepName})
		if err != nil {
			return nil, fmt.Errorf("查询模板步骤（%s）失败：%v", stepName, err)
		}
		tempSteps[stepName] = step
	}

	for idx, obj := range jsonSteps {
		mstep := obj.(map[string]interface{})
		tempStep := tempSteps[mstep["key"].(string)]
		if mstep["desc"] == nil || mstep["desc"] == "" {
			mstep["desc"] = tempStep.Desc
		}
		requires := make([]string, len(tempStep.Requires))
		if err := utils.Clone(&tempStep.Requires, &requires); err != nil {
			return nil, fmt.Errorf("复制步骤（%s）依赖失败：%v", tempStep.Key, err)
		}
		mstep["requires"] = requires
		if mstep["inputs"] == nil {
			var inputs interface{}
			if err := utils.Clone(tempStep.Inputs, &inputs); err != nil {
				return nil, fmt.Errorf("复制步骤（%s）输入失败：%v", tempStep.Key, err)
			}
			mstep["inputs"] = inputs
		}
		outputs := make([]string, len(tempStep.Outputs))
		if err := utils.Clone(&tempStep.Outputs, &outputs); err != nil {
			return nil, fmt.Errorf("复制步骤（%s）输出失败：%v", tempStep.Key, err)
		}
		mstep["outputs"] = outputs
		if mstep["code"] == nil || mstep["code"] == "" {
			mstep["code"] = tempStep.Code
		}
		mstep["symbol"] = tempStep.Symbol
		jsonSteps[idx] = mstep
	}
	minfo["steps"] = jsonSteps

	// 转成ApiInfo对象
	obj, err := utils.ToObj(minfo, reflect.TypeOf((*pb.ApiInfo)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转ApiInfo对象失败：%v", err)
	}

	// 调整步骤的序列号
	apiInfo := obj.(*pb.ApiInfo)
	for idx := range apiInfo.Steps {
		apiInfo.Steps[idx].Index = int32(idx)
	}
	return apiInfo, nil
}

func (s *Service) ApisSelectByName(ctx context.Context, req *pb.NameID) (*pb.ApiInfo, error) {
	res, err := s.mongo.QueryOne(ctx, model.API_INFO_TABLE, bson.D{{"name", req.Name}})
	if err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	}
	return s.complApiSteps(ctx, res)
}

func (s *Service) ApisSelectAll(ctx context.Context, _ *pb.Empty) (*pb.ApiInfoArray, error) {
	ress, err := s.mongo.Query(ctx, model.API_INFO_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	}

	resp := new(pb.ApiInfoArray)
	for _, res := range ress {
		apiInfo, err := s.complApiSteps(ctx, res)
		if err != nil {
			return nil, err
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

func (s *Service) ApisInsertByTemp(ctx context.Context, req *pb.AddTmpApiToMdlReq) (*pb.ApiInfo, error) {
	mname := req.ModelName.Name
	fmt.Printf("%v", req.TempApi)
	req.TempApi.Name = strings.Replace(req.TempApi.Name, "%MODEL%", mname, -1)
	req.TempApi.Model = strings.Replace(req.TempApi.Model, "%MODEL%", mname, -1)
	for pname, ptype := range req.TempApi.Params {
		req.TempApi.Params[pname] = strings.Replace(ptype, "%MODEL%", mname, -1)
	}
	req.TempApi.Route = strings.Replace(req.TempApi.Route, "%MODEL%", mname, -1)
	for index, ret := range req.TempApi.Returns {
		req.TempApi.Returns[index] = strings.Replace(ret, "%MODEL%", mname, -1)
	}
	for index, step := range req.TempApi.Steps {
		for idx, require := range step.Requires {
			step.Requires[idx] = strings.Replace(require, "%MODEL%", mname, -1)
		}
		step.Desc = strings.Replace(step.Desc, "%MODEL%", mname, -1)
		for slot, in := range step.Inputs {
			step.Inputs[slot] = strings.Replace(in, "%MODEL%", mname, -1)
		}
		req.TempApi.Steps[index] = step
	}

	if _, err := s.mongo.Insert(ctx, model.API_INFO_TABLE, req.TempApi); err != nil {
		return nil, fmt.Errorf("插入接口信息失败：%v", err)
	}
	return req.TempApi, nil
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

func (s *Service) TempApiSelectAll(ctx context.Context, _ *pb.Empty) (*pb.ApiInfoArray, error) {
	ress, err := s.mongo.Query(ctx, model.TEMP_API_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询所有模板接口失败：%v", err)
	}

	apiInfoType := reflect.TypeOf((*pb.ApiInfo)(nil)).Elem()
	resp := new(pb.ApiInfoArray)
	for _, res := range ress {
		obj, err := utils.ToObj(res, apiInfoType)
		if err != nil {
			return nil, fmt.Errorf("转ApiInfo对象失败：%v", err)
		}
		resp.Infos = append(resp.Infos, obj.(*pb.ApiInfo))
	}
	return resp, nil
}

func (s *Service) TempApiSelectByCategory(ctx context.Context, req *pb.CategoryIden) (*pb.ApiInfoArray, error) {
	ress, err := s.mongo.Query(ctx, model.TEMP_API_TABLE, bson.D{{"category", req.Category}})
	if err != nil {
		return nil, fmt.Errorf("查询类别为（%s）模板接口失败：%v", req.Category, err)
	}

	apiInfoType := reflect.TypeOf((*pb.ApiInfo)(nil)).Elem()
	resp := new(pb.ApiInfoArray)
	for _, res := range ress {
		obj, err := utils.ToObj(res, apiInfoType)
		if err != nil {
			return nil, fmt.Errorf("转ApiInfo对象失败：%v", err)
		}
		resp.Infos = append(resp.Infos, obj.(*pb.ApiInfo))
	}
 	return resp, nil
}

func (s *Service) TempApiInsert(ctx context.Context, req *pb.ApiInfo) (*pb.ApiInfo, error) {
	if _, err := s.mongo.Insert(ctx, model.TEMP_API_TABLE, req); err != nil {
		return nil, fmt.Errorf("插入模板接口失败：%v", err)
	}
	return req, nil
}

func (s *Service) TempApiInsertMany(ctx context.Context, req *pb.ApiInfoArray) (*pb.ApiInfoArray, error) {
	var entries []interface{}
	for _, api := range req.Infos {
		entries = append(entries, api)
	}
	if _, err := s.mongo.InsertMany(ctx, model.TEMP_API_TABLE, entries); err != nil {
		return nil, fmt.Errorf("批量插入模板接口失败：%v", err)
	}
	return req, nil
}

func (s *Service) StepsInsert(ctx context.Context, req *pb.StepReqs) (*pb.Empty, error) {
	apiName := req.Step.Apiname
	apiInfo, err := s.ApisSelectByName(ctx, &pb.NameID{Name: apiName})
	if err != nil {
		return nil, err
	}

	rear := append([]*pb.Step{}, apiInfo.Steps[req.Index:]...)
	apiInfo.Steps = append(append(apiInfo.Steps[:req.Index], req.Step), rear...)
	_, err = s.mongo.Update(ctx, model.API_INFO_TABLE, bson.D{{"name", apiName}}, bson.D{{
		"$set", bson.D{{"steps", apiInfo.Steps}},
	}})
	if err != nil {
		return nil, fmt.Errorf("插入步骤失败：%v", err)
	}
	return &pb.Empty{}, nil
}

func (s *Service) StepsDelete(ctx context.Context, req *pb.DelStepReqs) (*pb.Empty, error) {
	apiName := req.Apiname
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
func (s *Service) TempStepsInsert(ctx context.Context, req *pb.Step) (*pb.Step, error) {
	_, err := s.mongo.Insert(ctx, model.TEMP_STEP_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入模板步骤失败：%v", err)
	} else {
		return req, nil
	}
}

func (s *Service) TempStepsSelectAll(ctx context.Context, _ *pb.Empty) (*pb.StepArray, error) {
	ress, err := s.mongo.Query(ctx, model.TEMP_STEP_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询所有模板步骤失败：%v", err)
	}

	resp := new(pb.StepArray)
	for _, res := range ress {
		obj, err := utils.ToObj(res, reflect.TypeOf((*pb.Step)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转成步骤对象失败：%v", err)
		}
		resp.Steps = append(resp.Steps, obj.(*pb.Step))
	}
	return resp, nil
}

func (s *Service) TempStepsSelectByKey(ctx context.Context, req *pb.StrKey) (*pb.Step, error) {
	res, err := s.mongo.QueryOne(ctx, model.TEMP_STEP_TABLE, bson.D{{"key", req.Key}})
	if err != nil {
		return nil, fmt.Errorf("查询模板步骤失败：%v", err)
	}

	obj, err := utils.ToObj(res, reflect.TypeOf((*pb.Step)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转成模板步骤对象失败：%v", err)
	}
	return obj.(*pb.Step), nil
}

func (s *Service) TempStepsInsertMany(ctx context.Context, req *pb.StepArray) (*pb.StepArray, error) {
	var entries []interface{}
	for _, step := range req.Steps {
		entries = append(entries, step)
	}
	if _, err := s.mongo.InsertMany(ctx, model.TEMP_STEP_TABLE, entries); err != nil {
		return nil, fmt.Errorf("批量插入模板步骤失败：%v", err)
	}
	return req, nil
}

func (s *Service) TempStepsDeleteByKey(ctx context.Context, req *pb.StrKey) (*pb.Step, error) {
	step, err := s.TempStepsSelectByKey(ctx, req)
	if err != nil {
		return nil, err
	}

	if _, err = s.mongo.Delete(ctx, model.TEMP_STEP_TABLE, bson.D{{"key", req.Key}}); err != nil {
		return nil, fmt.Errorf("删除模板步骤失败：%v", err)
	}
	return step, nil
}

func (s *Service) DaoGroupsSelectAll(ctx context.Context, _ *pb.Empty) (*pb.DaoGroupArray, error) {
	ress, err := s.mongo.Query(ctx, model.DAO_GROUPS_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询DAO组失败：%v", err)
	}

	daoGroupType := reflect.TypeOf((*pb.DaoGroup)(nil)).Elem()
	resp := new(pb.DaoGroupArray)
	for _, res := range ress {
		daoGroup, err := utils.ToObj(res, daoGroupType)
		if err != nil {
			return nil, fmt.Errorf("转DaoGroup对象失败：%v", err)
		}
		resp.Groups = append(resp.Groups, daoGroup.(*pb.DaoGroup))
	}
	return resp, nil
}

func (s *Service) DaoGroupSelectByName(ctx context.Context, req *pb.NameID) (*pb.DaoGroup, error) {
	res, err := s.mongo.QueryOne(ctx, model.DAO_GROUPS_TABLE, bson.D{{"name", req.Name}})
	if err != nil {
		return nil, fmt.Errorf("查询DAO组失败：%v", err)
	}

	obj, err := utils.ToObj(res, reflect.TypeOf((*pb.DaoGroup)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转成DAO组失败：%v", err)
	}
	return obj.(*pb.DaoGroup), nil
}

func (s *Service) DaoGroupsInsert(ctx context.Context, req *pb.DaoGroup) (*pb.DaoGroup, error) {
	_, err := s.mongo.Insert(ctx, model.DAO_GROUPS_TABLE, req)
	if err != nil {
		return nil, fmt.Errorf("插入DAO组失败：%v", err)
	}
	return req, nil
}

func (s *Service) DaoGroupDeleteByName(ctx context.Context, req *pb.NameID) (*pb.DaoGroup, error) {
	resp, err := s.DaoGroupSelectByName(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("查询DAO组失败：%v", err)
	}

	_, err = s.mongo.Delete(ctx, model.DAO_GROUPS_TABLE, bson.D{{"name", req.Name}})
	if err != nil {
		return nil, fmt.Errorf("删除DAO组失败：%v", err)
	}
	return resp, nil
}

func (s *Service) TempDaoGroupsSelectAll(ctx context.Context, _ *pb.Empty) (*pb.DaoGroupArray, error) {
	ress, err := s.mongo.Query(ctx, model.TEMP_DAO_GROUPS_TABLE, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("查询所有模板DAO组失败：%v", err)
	}

	daoGroupType := reflect.TypeOf((*pb.DaoGroup)(nil)).Elem()
	resp := new(pb.DaoGroupArray)
	for _, res := range ress {
		obj, err := utils.ToObj(res, daoGroupType)
		if err != nil {
			return nil, fmt.Errorf("转成DAO组类型失败：%v", err)
		}
		resp.Groups = append(resp.Groups, obj.(*pb.DaoGroup))
	}
	return resp, nil
}

func (s *Service) DaoGroupUpdateImplement(ctx context.Context, req *pb.DaoGrpSetImpl) (*pb.DaoGroup, error) {
	modSign, err := s.ModuleInfoSelectBySignId(ctx, &pb.StrID{Id: req.ImplId})
	if err != nil {
		return nil, fmt.Errorf("查询MOD信息失败：%v", err)
	}

	group, err := s.DaoGroupSelectByName(ctx, &pb.NameID{Name: req.Gpname})
	if err != nil {
		return nil, err
	}

	group.Implement = req.ImplId
	_, err = s.mongo.Update(ctx, model.DAO_GROUPS_TABLE, bson.D{{"name", req.Gpname}}, bson.D{
		{"$set", bson.D{{"implement", group.Implement}}},
	})
	if err != nil {
		return nil, fmt.Errorf("配置DAO组实例化失败：%v", err)
	}

	// 生成配置
	if len(modSign.DaoConfHref) != 0 {
		jsonMap, err := utils.HttpGetJsonMap(modSign.DaoConfHref)
		if err != nil {
			return nil, fmt.Errorf("获取MOD配置信息失败：%v", err)
		}
		var configs []*pb.DaoConfItem
		for cfgKey, obj := range jsonMap["configs"].(map[string]interface{}) {
			cfgInfo := obj.(map[string]interface{})
			configs = append(configs, &pb.DaoConfItem{
				Key:                  cfgKey,
				Type:                 cfgInfo["type"].(string),
				Value:                "",
			})
		}
		if _, err := s.DaoConfigInsert(ctx, &pb.DaoConfig{Implement: req.ImplId, Configs: configs}); err != nil {
			return nil, fmt.Errorf("生成DAO配置失败：%v" , err)
		}
	}

	// 检查依赖
	for _, reqModId := range modSign.Requires {
		reqMod, err := s.ModuleInfoSelectBySignId(ctx, &pb.StrID{Id: reqModId})
		if err != nil {
			return nil, fmt.Errorf("查询依赖MOD信息失败：%v", err)
		}
		// 插入模板步骤
		if len(reqMod.TmpStepHref) != 0 {
			jsonMap, err := utils.HttpGetJsonMap(reqMod.TmpStepHref)
			if err != nil {
				return nil, fmt.Errorf("获取MOD模板步骤失败：%v", err)
			}
			if _, err := s.mongo.InsertMany(ctx, model.TEMP_STEP_TABLE, jsonMap["steps"].([]interface{})); err != nil {
				return nil, fmt.Errorf("添加MOD依赖的模板步骤失败：%v", err)
			}
		}
	}
	return group, nil
}

func (s *Service) DaoInterfaceInsert(ctx context.Context, req *pb.DaoItfcIst) (*pb.DaoInterface, error) {
	daoGroup, err := s.DaoGroupSelectByName(ctx, &pb.NameID{Name: req.Gpname})
	if err != nil {
		return nil, err
	}

	daoGroup.Interfaces = append(daoGroup.Interfaces, req.Interface)
	_, err = s.mongo.Update(ctx, model.DAO_GROUPS_TABLE, bson.D{{"name", req.Gpname}}, bson.D{{
		"$set", bson.D{{"interfaces", daoGroup.Interfaces}},
	}})
	if err != nil {
		return nil, fmt.Errorf("插入DAO接口失败：%v", err)
	}
	return req.Interface, nil
}

func (s *Service) DaoInterfaceDelete(ctx context.Context, req *pb.DaoItfcIden) (*pb.DaoInterface, error) {
	daoGroup, err := s.DaoGroupSelectByName(ctx, &pb.NameID{Name: req.Gpname})
	if err != nil {
		return nil, err
	}

	var resp *pb.DaoInterface
	for i, itfc := range daoGroup.Interfaces {
		if itfc.Name != req.Ifname {
			continue
		}
		resp = itfc
		interfaces := append(daoGroup.Interfaces[:i], daoGroup.Interfaces[i + 1:]...)
		_, err := s.mongo.Update(ctx, model.DAO_GROUPS_TABLE, bson.D{{"name", req.Gpname}}, bson.D{{
			"$set", bson.D{{"interfaces", interfaces}},
		}})
		if err != nil {
			return nil, fmt.Errorf("删除DAO接口失败：%v", err)
		}
		break
	}
	return resp, nil
}

func (s *Service) DaoConfigInsert(ctx context.Context, req *pb.DaoConfig) (*pb.DaoConfig, error) {
	_, err := s.DaoConfigSelectByImpl(ctx, &pb.DaoConfImplIden{Implement: req.Implement})
	if err != nil {
		_, err := s.mongo.Insert(ctx, model.DAO_CONFIGS_TABLE, req)
		if err != nil {
			return nil, fmt.Errorf("插入DAO配置失败：%v", err)
		}
	}
	if _, err := s.mongo.Update(ctx, model.DAO_CONFIGS_TABLE,
		bson.D{{"implement", req.Implement}},
		bson.D{{"$set", bson.D{{"configs", req.Configs}}}},
	); err != nil {
		return nil, fmt.Errorf("更新DAO配置失败：%v", err)
	}
	return req, nil
}

func (s *Service) DaoConfigSelectByImpl(ctx context.Context, req *pb.DaoConfImplIden) (*pb.DaoConfig, error) {
	res, err := s.mongo.QueryOne(ctx, model.DAO_CONFIGS_TABLE, bson.D{{"implement", req.Implement}})
	if err != nil {
		return nil, fmt.Errorf("查询DAO配置失败：%v", err)
	}

	obj, err := utils.ToObj(res, reflect.TypeOf((*pb.DaoConfig)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转成DAO配置失败：%v", err)
	}
	return obj.(*pb.DaoConfig), nil
}

func (s *Service) Export(ctx context.Context, req *pb.ExpOptions) (*pb.UrlResp, error) {
	if pjPath, err := s.ac.Get("projPath").String(); err != nil {
		return nil, fmt.Errorf("配置文件中未定义项目目录：%v", err)
	} else if wsPath, err := s.ac.Get("workspace").String(); err != nil {
		return nil, fmt.Errorf("配置文件中未定义工作区目录：%v", err)
	} else if wsPath = path.Join(pjPath, wsPath); false {
		return nil, nil
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
	} else if pathName := path.Join(cchPath, req.Name); false {
		return nil, nil
	} else if utils.CopyFolder(tmpPath, pathName); false {
		return nil, nil
	} else if err := s.pBuilder.Build(s, req, pathName).Adjust(ctx); err != nil {
		return nil, fmt.Errorf("编辑项目失败：%v", err)
	} else if wsFile, err := os.Open(pathName); err != nil {
		return nil, fmt.Errorf("工作区目录有误，打开失败：%v", err)
	} else if zipPath := path.Join(cchPath, req.Name + ".zip"); false {
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

func (s *Service) SpecialSymbols(context.Context, *pb.Empty) (*pb.SymbolsResp, error) {
	return &pb.SymbolsResp{
		Values: pb.SpcSymbol_value,
		Names:  pb.SpcSymbol_name,
	}, nil
}

func (s *Service) ModuleSignSelectAll(ctx context.Context, req *pb.TypeIden) (*pb.ModuleSignArray, error) {
	ress, err := s.mongo.Query(ctx, model.MOD_SIGN_TABLE, bson.D{{"type", req.Type}})
	if err != nil {
		return nil, fmt.Errorf("查询模块标牌失败：%v", err)
	}

	resp := new(pb.ModuleSignArray)
	for _, res := range ress {
		ms, err := utils.ToObj(res, reflect.TypeOf((*pb.ModuleSign)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("转成ModuleSign对象失败：%v", err)
		}
		resp.ModSigns = append(resp.ModSigns, ms.(*pb.ModuleSign))
	}
	return resp, nil
}

func (s *Service) ModuleInfoSelectBySignId(ctx context.Context, req *pb.StrID) (*pb.ModuleSign, error) {
	res, err := s.mongo.QueryOne(ctx, model.MOD_SIGN_TABLE, bson.D{{"id", req.Id}})
	if err != nil {
		return nil, fmt.Errorf("根据模块ID获取模块信息失败：%v", err)
	}

	obj, err := utils.ToObj(res, reflect.TypeOf((*pb.ModuleSign)(nil)).Elem())
	if err != nil {
		return nil, fmt.Errorf("转成ModuleSign对象失败：%v", err)
	}
	return obj.(*pb.ModuleSign), nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.mongo.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.mongo.Close()
}
