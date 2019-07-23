// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: backend/api/api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

It is generated from these files:
	backend/api/api.proto
*/
package api

import (
	"context"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"github.com/bilibili/kratos/pkg/net/http/blademaster/binding"
)

// to suppressed 'imported but not used warning'
var _ *bm.Context
var _ context.Context
var _ binding.StructValidator

var PathBackendManagerModelsInsert = "/backend-dashboard/backend/models.insert"
var PathBackendManagerModelsDelete = "/backend-dashboard/backend/models.delete"
var PathBackendManagerModelsUpdate = "/backend-dashboard/backend/models.update"
var PathBackendManagerModelsSelectAll = "/backend-dashboard/backend/models.selectAll"
var PathBackendManagerModelsSelectByName = "/backend-dashboard/backend/models.selectByName"
var PathBackendManagerLinksInsert = "/backend-dashboard/backend/links.insert"
var PathBackendManagerLinksSelectAll = "/backend-dashboard/backend/links.selectAll"
var PathBackendManagerLinksDeleteBySymbol = "/backend-dashboard/backend/links.deleteBySymbol"
var PathBackendManagerApisSelectByName = "/backend-dashboard/backend/apis.selectByName"
var PathBackendManagerApisSelectAll = "/backend-dashboard/backend/apis.selectAll"
var PathBackendManagerApisInsert = "/backend-dashboard/backend/apis.insert"
var PathBackendManagerFlowInsert = "/backend-dashboard/backend/flows.insert"
var PathBackendManagerOperStepsSelectTemp = "/backend-dashboard/backend/steps.selectTemp"
var PathBackendManagerOperStepsInsert = "/backend-dashboard/backend/steps.insert"
var PathBackendManagerExport = "/backend-dashboard/backend/export"
var PathBackendManagerSpecialSymbols = "/backend-dashboard/backend/specialSymbols"

// BackendManagerBMServer is the server API for BackendManager service.
type BackendManagerBMServer interface {
	ModelsInsert(ctx context.Context, req *Model) (resp *Model, err error)

	ModelsDelete(ctx context.Context, req *NameID) (resp *Model, err error)

	ModelsUpdate(ctx context.Context, req *Model) (resp *Empty, err error)

	ModelsSelectAll(ctx context.Context, req *Empty) (resp *ModelArray, err error)

	ModelsSelectByName(ctx context.Context, req *NameID) (resp *Model, err error)

	LinksInsert(ctx context.Context, req *Link) (resp *Link, err error)

	LinksSelectAll(ctx context.Context, req *Empty) (resp *LinkArray, err error)

	LinksDeleteBySymbol(ctx context.Context, req *SymbolID) (resp *Link, err error)

	ApisSelectByName(ctx context.Context, req *NameID) (resp *ApiInfo, err error)

	ApisSelectAll(ctx context.Context, req *Empty) (resp *ApiInfoArray, err error)

	ApisInsert(ctx context.Context, req *ApiInfo) (resp *ApiInfo, err error)

	FlowInsert(ctx context.Context, req *FlowReqs) (resp *Empty, err error)

	OperStepsSelectTemp(ctx context.Context, req *Empty) (resp *OperStepArray, err error)

	OperStepsInsert(ctx context.Context, req *OperStep) (resp *OperStep, err error)

	Export(ctx context.Context, req *ExpOptions) (resp *UrlResp, err error)

	SpecialSymbols(ctx context.Context, req *Empty) (resp *SymbolsResp, err error)
}

var BackendManagerSvc BackendManagerBMServer

func backendManagerModelsInsert(c *bm.Context) {
	p := new(Model)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerModelsDelete(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsDelete(c, p)
	c.JSON(resp, err)
}

func backendManagerModelsUpdate(c *bm.Context) {
	p := new(Model)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsUpdate(c, p)
	c.JSON(resp, err)
}

func backendManagerModelsSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerModelsSelectByName(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsSelectByName(c, p)
	c.JSON(resp, err)
}

func backendManagerLinksInsert(c *bm.Context) {
	p := new(Link)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.LinksInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerLinksSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.LinksSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerLinksDeleteBySymbol(c *bm.Context) {
	p := new(SymbolID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.LinksDeleteBySymbol(c, p)
	c.JSON(resp, err)
}

func backendManagerApisSelectByName(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ApisSelectByName(c, p)
	c.JSON(resp, err)
}

func backendManagerApisSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ApisSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerApisInsert(c *bm.Context) {
	p := new(ApiInfo)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ApisInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerFlowInsert(c *bm.Context) {
	p := new(FlowReqs)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.FlowInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerOperStepsSelectTemp(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.OperStepsSelectTemp(c, p)
	c.JSON(resp, err)
}

func backendManagerOperStepsInsert(c *bm.Context) {
	p := new(OperStep)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.OperStepsInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerExport(c *bm.Context) {
	p := new(ExpOptions)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.Export(c, p)
	c.JSON(resp, err)
}

func backendManagerSpecialSymbols(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.SpecialSymbols(c, p)
	c.JSON(resp, err)
}

// RegisterBackendManagerBMServer Register the blademaster route
func RegisterBackendManagerBMServer(e *bm.Engine, server BackendManagerBMServer) {
	BackendManagerSvc = server
	e.POST("/backend-dashboard/backend/models.insert", backendManagerModelsInsert)
	e.POST("/backend-dashboard/backend/models.delete", backendManagerModelsDelete)
	e.POST("/backend-dashboard/backend/models.update", backendManagerModelsUpdate)
	e.POST("/backend-dashboard/backend/models.selectAll", backendManagerModelsSelectAll)
	e.POST("/backend-dashboard/backend/models.selectByName", backendManagerModelsSelectByName)
	e.POST("/backend-dashboard/backend/links.insert", backendManagerLinksInsert)
	e.POST("/backend-dashboard/backend/links.selectAll", backendManagerLinksSelectAll)
	e.POST("/backend-dashboard/backend/links.deleteBySymbol", backendManagerLinksDeleteBySymbol)
	e.POST("/backend-dashboard/backend/apis.selectByName", backendManagerApisSelectByName)
	e.POST("/backend-dashboard/backend/apis.selectAll", backendManagerApisSelectAll)
	e.POST("/backend-dashboard/backend/apis.insert", backendManagerApisInsert)
	e.POST("/backend-dashboard/backend/flows.insert", backendManagerFlowInsert)
	e.POST("/backend-dashboard/backend/steps.selectTemp", backendManagerOperStepsSelectTemp)
	e.POST("/backend-dashboard/backend/steps.insert", backendManagerOperStepsInsert)
	e.POST("/backend-dashboard/backend/export", backendManagerExport)
	e.POST("/backend-dashboard/backend/specialSymbols", backendManagerSpecialSymbols)
}
