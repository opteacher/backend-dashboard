package service

import (
	"context"
	"path"
	"fmt"
	pb "template/api"
	"template/internal/dao"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	// [MICO_SERV_BEG]
	"template/internal/server"
	// [MICO_SERV_END]

	// [IMPORTS]
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
	// [INIT]
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

// [APIS]

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
	// [MICO_SERV_BEG]
	if cli, err := server.RegisterService(); err != nil {
		panic(err)
	} else if _, err := cli.Cancel(context.Background(), &pb.IdenSvcReqs{AppID: s.AppID()}); err != nil {
		panic(err)
	}
	// [MICO_SERV_END]
}
