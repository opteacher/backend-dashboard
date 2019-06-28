package server

import (
	"context"
	"fmt"
	pb "template/api"
	"template/internal/utils"

	"github.com/bilibili/kratos/pkg/net/rpc/warden"
	"github.com/bilibili/kratos/pkg/log"
	"template/internal/service"
)

func RegisterService() (pb.RegisterClient, error) {
	cli := warden.NewClient(nil)
	conn, err := cli.Dial(context.Background(), "discovery://default/register.service")
	if err != nil {
		return nil, fmt.Errorf("Register center is unready: %v\n", err)
	}
	return pb.NewRegisterClient(conn), nil
}

func RegisterGRPCService(appID string, addrs []string) {
	if cli, err := RegisterService(); err != nil {
		panic(err)
	} else if resp, err := cli.RegAsGRPC(context.Background(), &pb.RegSvcReqs{
		AppID: appID,
		Urls:  addrs,
	}); err != nil {
		panic(err)
	} else {
		fmt.Println(resp)
	}
}

func RegisterHTTPService(svc *service.Service, addrs []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	appID := svc.AppID()
	if cli, err := RegisterService(); err != nil {
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