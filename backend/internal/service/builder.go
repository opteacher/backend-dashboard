package service

import (
	pb "backend/api"
	"backend/internal/utils"
	"context"
	"strings"
)

type ProjBuilder struct {
}

func NewProjBuilder() *ProjBuilder {
	return new(ProjBuilder)
}

type ProjInfo struct {
	option   *pb.ExpOptions
	pathName string
	pkgName  string
}

type BaseBuilder interface {
	Adjust(context.Context) error
}

func (pgb *ProjBuilder) Build(svc *Service, option *pb.ExpOptions, pathName string) BaseBuilder {
	info := new(ProjInfo)
	info.option = option
	info.option.Name = strings.TrimRight(info.option.Name, ".zip")
	info.option.Name = strings.TrimRight(info.option.Name, ".ZIP")
	info.pkgName = utils.CamelToPascal(info.option.Name)
	info.pathName = pathName
	gen, err := NewKratos(svc, info)
	if err != nil {
		panic(err)
	}
	return gen
}
