package service

import (
	pb "backend/api"
	"backend/internal/utils"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Koa struct {
	info *ProjInfo
	svc  *Service
}

func (_ *newCollector) NewKoa(svc *Service, pi *ProjInfo) (*Koa, error) {
	return &Koa{info: pi, svc: svc}, nil
}

func (koa *Koa) Adjust(ctx context.Context) error {
	// 调整前端支持
	if err := koa.adjFrontendSupt(ctx); err != nil {
		return fmt.Errorf("调整前端环境支持失败：%v", err)
	}
	// 调整控制层和服务层
	if err := koa.genApis(ctx); err != nil {
		return fmt.Errorf("生成API路径失败：%v", err)
	}
	// 调整DAO层
	return nil
}

func  (koa *Koa) chkSuptFrontend() bool {
	return koa.info.option.Components["frontendEnv"] != nil && koa.info.option.Components["frontendEnv"].Enable
}

func (koa *Koa) adjFrontendSupt(ctx context.Context) error {
	suptFrontend := koa.chkSuptFrontend()

	if err := utils.DelByTagInFile(path.Join(koa.info.pathName, "app.js"), "FRONTEND_ENV", !suptFrontend); err != nil {
		return fmt.Errorf("调整app.js的前端环境支持逻辑时发生错误：%v", err)
	}

	if suptFrontend {
		return nil
	}

	for _, delFile := range []string{
		filepath.Join(koa.info.pathName, "public"),
		filepath.Join(koa.info.pathName, ".babelrc"),
		filepath.Join(koa.info.pathName, "index.html"),
		filepath.Join(koa.info.pathName, "webpack.config.js"),
	} {
		if err := os.RemoveAll(delFile); err != nil {
			return fmt.Errorf("删除文件失败：%v", err)
		}
	}
	return nil
}

func (koa *Koa) genApis(ctx context.Context) error {
	apis, err := koa.svc.ApisSelectAll(ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	for _, api := range apis.Infos {
		fmt.Println(api.Model)
	}
	return nil
}