package service

import (
	"context"
	"reflect"

	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/server"

	"encoding/json"
	"fmt"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"backend/internal/utils"
)

// Service service.
type Service struct {
	ac  *paladin.Map
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

func (s *Service) Models(ctx context.Context, req *pb.IdenReqs) (models *pb.ModelArray, err error) {
	methodID, exs := pb.MethodType_value[req.Method]
	if !exs {
		return nil, fmt.Errorf("指定的Method必须在%v中", pb.MethodType_name)
	}
	models = new(pb.ModelArray)
	switch pb.MethodType(methodID) {
	case pb.MethodType_INSERT:
		tx, err := s.dao.BeginTx(ctx)
		if err != nil {
			return nil, fmt.Errorf("开启事务失败：%v", err)
		}
		for _, strEntry := range req.Params {
			entry := make(map[string]interface{})
			err = json.Unmarshal([]byte(strEntry), &entry)
			if err != nil {
				return nil, fmt.Errorf("提交的请求参数格式有误：%v", err)
			}

			_, err = s.dao.SaveTx(tx, model.MODELS_TABLE, "", nil, entry, false)
			if err != nil {
				return nil, fmt.Errorf("插入数据库失败：%v", err)
			}
		}
		err = s.dao.CommitTx(tx)
		if err != nil {
			return nil, fmt.Errorf("提交插入事务失败：%v", err)
		}
	case pb.MethodType_SELECT:
		conds := string(req.Params[0])
		var argus []interface{}
		for i := 1; i < len(req.Params); i++ {
			argus = append(argus, string(req.Params[i]))
		}
		entries, err := s.dao.Query(ctx, model.MODELS_TABLE, conds, argus)
		if err != nil {
			return nil, fmt.Errorf("查询数据库失败：%v", err)
		}
		for _, entry := range entries {
			mdl := new(pb.Model)
			utils.FillWithMap(mdl, entry)
			if mdl.X == 0 {
				mdl.X = 1
			}
			if mdl.Y == 0 {
				mdl.Y = 1
			}
			models.Models = append(models.Models, mdl)
		}
	default:

	}
	return
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
