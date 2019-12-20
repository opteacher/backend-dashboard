package http

import (
	"net/http"
	"strings"

	svr "template/internal/server"
	"template/internal/service"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	pb "template/api"
	"template/internal/server/mws"
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
	pb.RegisterDemoBMServer(engine, svc)
	svr.RegisterHTTPService(svc, []string{strings.Replace(hc.Server.Addr, "0.0.0.0", "127.0.0.1", 1)})
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
