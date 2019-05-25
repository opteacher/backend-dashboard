package http

import (
	"net/http"

	"backend/internal/model"
	globalmw "backend/internal/server/middleware"
	"backend/internal/service"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/log"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

var (
	svc *service.Service
	api *service.ApiService
)

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine) {
	// 读取服务配置
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

	// 初始化服务
	svc = s
	api = service.NewApiService()
	// 使用默认blademaster管理网关
	engine = bm.DefaultServer(hc.Server)
	// 注册中间件
	engine.Use(
		globalmw.SetupCORS(),
		globalmw.ParseJSON(),
	)
	// 初始化路由
	initRouter(engine)
	// 开启监听
	if err := engine.Start(); err != nil {
		panic(err)
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/backend-dashboard/backend")
	{
		service.InsApiService().AddModelAPI(g, model.MODELS_NAME, []string{
			service.INSERT,
		})
	}
}

func ping(ctx *bm.Context) {
	if err := svc.Ping(ctx); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
