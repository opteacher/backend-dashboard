package service

import (
	pb "backend/api"
	"backend/internal/utils"
	"context"
	"fmt"
	"reflect"
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

type newCollector struct {}
var _newCollector = new(newCollector)

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
	fwks, err := svc.TempFrameworkSelect(context.TODO(), &pb.Empty{})
	if err != nil {
		panic(err)
	}
	generator := ""
	for _, framework := range fwks.Frameworks {
		if framework.Id == option.Type {
			generator = framework.Generator
			break
		}
	}
	if len(generator) == 0 {
		panic(fmt.Errorf("为找到指定导出框架：%s", option.Type))
	}
	rets := reflect.ValueOf(_newCollector).MethodByName(generator).Call([]reflect.Value{
		reflect.ValueOf(svc),  reflect.ValueOf(info),
	})
	if !rets[1].IsNil() {
		panic(rets[1].Interface().(error))
	}
	return rets[0].Interface().(BaseBuilder)
}
