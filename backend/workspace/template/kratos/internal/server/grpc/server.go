package grpc

import (
	"strings"
	"template/internal/service"
	pb "template/api"
	svr "template/internal/server"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc *service.Service) *warden.Server {
	var rc struct {
		Server *warden.ServerConfig
	}
	if err := paladin.Get("grpc.toml").UnmarshalTOML(&rc); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	ws := warden.NewServer(rc.Server)
	pb.RegisterDemoServer(ws.Server(), svc)
	svr.RegisterGRPCService(svc.AppID(), []string{strings.Replace(rc.Server.Addr, "0.0.0.0", "127.0.0.1", 1)})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}
