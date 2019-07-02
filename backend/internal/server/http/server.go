package http

import (
	"net/http"
	"strings"

	pb "backend/api"
	"backend/internal/service"

	"context"
	"backend/internal/server"
	"backend/internal/utils"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"backend/internal/server/mws"
)

var (
	svc *service.Service
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine) {
	var (
		hc struct {
			Server *bm.ServerConfig
		}
	)
	if err := paladin.Get("http.toml").UnmarshalTOML(&hc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	svc = s
	engine = bm.DefaultServer(hc.Server)
	engine.Use(mws.SetupCORS())
	pb.RegisterBackendManagerBMServer(engine, svc)
	RegisterHTTPService(svc, []string{strings.Replace(hc.Server.Addr, "0.0.0.0", "127.0.0.1", 1)})
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func RegisterHTTPService(svc *service.Service, addrs []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	appID := svc.AppID()
	if cli, err := server.RegisterService(); err != nil {
		log.Error("Fetch discovery service error: %v", err)
	} else if _, err := cli.RegAsHTTP(ctx, &pb.RegSvcReqs{
		AppID: appID,
		Urls:  addrs,
	}); err != nil {
	} else if data, err := utils.PickPathsFromSwaggerJSON(svc.SwaggerFile()); err != nil {
		log.Error("API swagger file open failed: %v", err)
	} else if _, err := cli.AddRoutes(ctx, &pb.AddRoutesReqs{
		ServiceID: appID,
		Paths:     data,
	}); err != nil {
		panic(err)
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
