// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: api.proto

/*
Package api is a generated blademaster stub package.
This code was generated with kratos/tool/protobuf/protoc-gen-bm v0.1.

It is generated from these files:
	api.proto
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
var PathBackendManagerLinkInsert = "/backend-dashboard/backend/links.insert"
var PathBackendManagerLinkSelectAll = "/backend-dashboard/backend/links.selectAll"
var PathBackendManagerLinksDeleteBySymbol = "/backend-dashboard/backend/links.deleteBySymbol"
var PathBackendManagerExport = "/backend-dashboard/backend/export"

// BackendManagerBMServer is the server API for BackendManager service.
type BackendManagerBMServer interface {
	ModelsInsert(ctx context.Context, req *Model) (resp *Model, err error)

	ModelsDelete(ctx context.Context, req *NameID) (resp *Model, err error)

	ModelsUpdate(ctx context.Context, req *Model) (resp *Empty, err error)

	ModelsSelectAll(ctx context.Context, req *Empty) (resp *ModelArray, err error)

	ModelsSelectByName(ctx context.Context, req *NameID) (resp *Model, err error)

	LinkInsert(ctx context.Context, req *Link) (resp *Link, err error)

	LinkSelectAll(ctx context.Context, req *Empty) (resp *LinkArray, err error)

	LinksDeleteBySymbol(ctx context.Context, req *SymbolID) (resp *Link, err error)

	Export(ctx context.Context, req *ExpOptions) (resp *UrlResp, err error)
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

func backendManagerLinkInsert(c *bm.Context) {
	p := new(Link)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.LinkInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerLinkSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.LinkSelectAll(c, p)
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

func backendManagerExport(c *bm.Context) {
	p := new(ExpOptions)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.Export(c, p)
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
	e.POST("/backend-dashboard/backend/links.insert", backendManagerLinkInsert)
	e.POST("/backend-dashboard/backend/links.selectAll", backendManagerLinkSelectAll)
	e.POST("/backend-dashboard/backend/links.deleteBySymbol", backendManagerLinksDeleteBySymbol)
	e.POST("/backend-dashboard/backend/export", backendManagerExport)
}
