// Code generated by protoc-gen-bm v0.1, DO NOT EDIT.
// source: api.proto

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
var PathBackendManagerStructsSelectAllBases = "/backend-dashboard/backend/structs.selectAllBases"
var PathBackendManagerLinksInsert = "/backend-dashboard/backend/links.insert"
var PathBackendManagerLinksSelectAll = "/backend-dashboard/backend/links.selectAll"
var PathBackendManagerLinksDeleteBySymbol = "/backend-dashboard/backend/links.deleteBySymbol"
var PathBackendManagerApisSelectByName = "/backend-dashboard/backend/apis.selectByName"
var PathBackendManagerApisSelectAll = "/backend-dashboard/backend/apis.selectAll"
var PathBackendManagerApisInsert = "/backend-dashboard/backend/apis.insert"
var PathBackendManagerApisDeleteByName = "/backend-dashboard/backend/apis.deleteByName"
var PathBackendManagerStepsInsert = "/backend-dashboard/backend/steps.insert"
var PathBackendManagerStepsDelete = "/backend-dashboard/backend/steps.delete"
var PathBackendManagerTempStepsSelectAll = "/backend-dashboard/backend/temp.steps.selectAll"
var PathBackendManagerTempStepsSelectByKey = "/backend-dashboard/backend/temp.steps.selectByKey"
var PathBackendManagerTempStepsInsert = "/backend-dashboard/backend/temp.steps.insert"
var PathBackendManagerTempStepsInsertMany = "/backend-dashboard/backend/temp.steps.insertMany"
var PathBackendManagerTempStepsDeleteByKey = "/backend-dashboard/backend/temp.steps.deleteByKey"
var PathBackendManagerDaoGroupsSelectAll = "/backend-dashboard/backend/dao.groups.selectAll"
var PathBackendManagerDaoGroupSelectByName = "/backend-dashboard/backend/dao.groups.selectByName"
var PathBackendManagerDaoGroupsInsert = "/backend-dashboard/backend/dao.groups.insert"
var PathBackendManagerDaoGroupDeleteByName = "/backend-dashboard/backend/dao.groups.deleteByName"
var PathBackendManagerDaoGroupUpdateImplement = "/backend-dashboard/backend/dao.groups.updateImplement"
var PathBackendManagerDaoInterfaceInsert = "/backend-dashboard/backend/dao.interface.insert"
var PathBackendManagerDaoInterfaceDelete = "/backend-dashboard/backend/dao.interface.delete"
var PathBackendManagerDaoConfigInsert = "/backend-dashboard/backend/dao.config.insert"
var PathBackendManagerExport = "/backend-dashboard/backend/export"
var PathBackendManagerSpecialSymbols = "/backend-dashboard/backend/specialSymbols"
var PathBackendManagerModuleSignSelectAll = "/backend-dashboard/backend/mod.sign.selectAll"
var PathBackendManagerModuleInfoSelectBySignId = "/backend-dashboard/backend/mod.info.selectBySignId"

// BackendManagerBMServer is the server API for BackendManager service.
type BackendManagerBMServer interface {
	ModelsInsert(ctx context.Context, req *Model) (resp *Model, err error)

	ModelsDelete(ctx context.Context, req *NameID) (resp *Model, err error)

	ModelsUpdate(ctx context.Context, req *Model) (resp *Empty, err error)

	ModelsSelectAll(ctx context.Context, req *TypeIden) (resp *ModelArray, err error)

	StructsSelectAllBases(ctx context.Context, req *Empty) (resp *NameArray, err error)

	LinksInsert(ctx context.Context, req *Link) (resp *Link, err error)

	LinksSelectAll(ctx context.Context, req *Empty) (resp *LinkArray, err error)

	LinksDeleteBySymbol(ctx context.Context, req *SymbolID) (resp *Link, err error)

	ApisSelectByName(ctx context.Context, req *NameID) (resp *ApiInfo, err error)

	ApisSelectAll(ctx context.Context, req *Empty) (resp *ApiInfoArray, err error)

	ApisInsert(ctx context.Context, req *ApiInfo) (resp *ApiInfo, err error)

	ApisDeleteByName(ctx context.Context, req *NameID) (resp *ApiInfo, err error)

	StepsInsert(ctx context.Context, req *StepReqs) (resp *Empty, err error)

	StepsDelete(ctx context.Context, req *DelStepReqs) (resp *Empty, err error)

	TempStepsSelectAll(ctx context.Context, req *Empty) (resp *StepArray, err error)

	TempStepsSelectByKey(ctx context.Context, req *StrKey) (resp *Step, err error)

	TempStepsInsert(ctx context.Context, req *Step) (resp *Step, err error)

	TempStepsInsertMany(ctx context.Context, req *StepArray) (resp *StepArray, err error)

	TempStepsDeleteByKey(ctx context.Context, req *StrKey) (resp *Step, err error)

	DaoGroupsSelectAll(ctx context.Context, req *Empty) (resp *DaoGroupArray, err error)

	DaoGroupSelectByName(ctx context.Context, req *NameID) (resp *DaoGroup, err error)

	DaoGroupsInsert(ctx context.Context, req *DaoGroup) (resp *DaoGroup, err error)

	DaoGroupDeleteByName(ctx context.Context, req *NameID) (resp *DaoGroup, err error)

	DaoGroupUpdateImplement(ctx context.Context, req *DaoGrpSetImpl) (resp *DaoGroup, err error)

	DaoInterfaceInsert(ctx context.Context, req *DaoItfcIst) (resp *DaoInterface, err error)

	DaoInterfaceDelete(ctx context.Context, req *DaoItfcIden) (resp *DaoInterface, err error)

	DaoConfigInsert(ctx context.Context, req *DaoConfig) (resp *DaoConfig, err error)

	Export(ctx context.Context, req *ExpOptions) (resp *UrlResp, err error)

	SpecialSymbols(ctx context.Context, req *Empty) (resp *SymbolsResp, err error)

	ModuleSignSelectAll(ctx context.Context, req *TypeIden) (resp *ModuleSignArray, err error)

	ModuleInfoSelectBySignId(ctx context.Context, req *StrID) (resp *ModuleSign, err error)
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
	p := new(TypeIden)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModelsSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerStructsSelectAllBases(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.StructsSelectAllBases(c, p)
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

func backendManagerApisDeleteByName(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ApisDeleteByName(c, p)
	c.JSON(resp, err)
}

func backendManagerStepsInsert(c *bm.Context) {
	p := new(StepReqs)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.StepsInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerStepsDelete(c *bm.Context) {
	p := new(DelStepReqs)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.StepsDelete(c, p)
	c.JSON(resp, err)
}

func backendManagerTempStepsSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.TempStepsSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerTempStepsSelectByKey(c *bm.Context) {
	p := new(StrKey)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.TempStepsSelectByKey(c, p)
	c.JSON(resp, err)
}

func backendManagerTempStepsInsert(c *bm.Context) {
	p := new(Step)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.TempStepsInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerTempStepsInsertMany(c *bm.Context) {
	p := new(StepArray)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.TempStepsInsertMany(c, p)
	c.JSON(resp, err)
}

func backendManagerTempStepsDeleteByKey(c *bm.Context) {
	p := new(StrKey)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.TempStepsDeleteByKey(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoGroupsSelectAll(c *bm.Context) {
	p := new(Empty)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoGroupsSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoGroupSelectByName(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoGroupSelectByName(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoGroupsInsert(c *bm.Context) {
	p := new(DaoGroup)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoGroupsInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoGroupDeleteByName(c *bm.Context) {
	p := new(NameID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoGroupDeleteByName(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoGroupUpdateImplement(c *bm.Context) {
	p := new(DaoGrpSetImpl)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoGroupUpdateImplement(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoInterfaceInsert(c *bm.Context) {
	p := new(DaoItfcIst)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoInterfaceInsert(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoInterfaceDelete(c *bm.Context) {
	p := new(DaoItfcIden)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoInterfaceDelete(c, p)
	c.JSON(resp, err)
}

func backendManagerDaoConfigInsert(c *bm.Context) {
	p := new(DaoConfig)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.DaoConfigInsert(c, p)
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

func backendManagerModuleSignSelectAll(c *bm.Context) {
	p := new(TypeIden)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModuleSignSelectAll(c, p)
	c.JSON(resp, err)
}

func backendManagerModuleInfoSelectBySignId(c *bm.Context) {
	p := new(StrID)
	if err := c.BindWith(p, binding.Default(c.Request.Method, c.Request.Header.Get("Content-Type"))); err != nil {
		return
	}
	resp, err := BackendManagerSvc.ModuleInfoSelectBySignId(c, p)
	c.JSON(resp, err)
}

// RegisterBackendManagerBMServer Register the blademaster route
func RegisterBackendManagerBMServer(e *bm.Engine, server BackendManagerBMServer) {
	BackendManagerSvc = server
	e.POST("/backend-dashboard/backend/models.insert", backendManagerModelsInsert)
	e.POST("/backend-dashboard/backend/models.delete", backendManagerModelsDelete)
	e.POST("/backend-dashboard/backend/models.update", backendManagerModelsUpdate)
	e.POST("/backend-dashboard/backend/models.selectAll", backendManagerModelsSelectAll)
	e.POST("/backend-dashboard/backend/structs.selectAllBases", backendManagerStructsSelectAllBases)
	e.POST("/backend-dashboard/backend/links.insert", backendManagerLinksInsert)
	e.POST("/backend-dashboard/backend/links.selectAll", backendManagerLinksSelectAll)
	e.POST("/backend-dashboard/backend/links.deleteBySymbol", backendManagerLinksDeleteBySymbol)
	e.POST("/backend-dashboard/backend/apis.selectByName", backendManagerApisSelectByName)
	e.POST("/backend-dashboard/backend/apis.selectAll", backendManagerApisSelectAll)
	e.POST("/backend-dashboard/backend/apis.insert", backendManagerApisInsert)
	e.POST("/backend-dashboard/backend/apis.deleteByName", backendManagerApisDeleteByName)
	e.POST("/backend-dashboard/backend/steps.insert", backendManagerStepsInsert)
	e.POST("/backend-dashboard/backend/steps.delete", backendManagerStepsDelete)
	e.POST("/backend-dashboard/backend/temp.steps.selectAll", backendManagerTempStepsSelectAll)
	e.POST("/backend-dashboard/backend/temp.steps.selectByKey", backendManagerTempStepsSelectByKey)
	e.POST("/backend-dashboard/backend/temp.steps.insert", backendManagerTempStepsInsert)
	e.POST("/backend-dashboard/backend/temp.steps.insertMany", backendManagerTempStepsInsertMany)
	e.POST("/backend-dashboard/backend/temp.steps.deleteByKey", backendManagerTempStepsDeleteByKey)
	e.POST("/backend-dashboard/backend/dao.groups.selectAll", backendManagerDaoGroupsSelectAll)
	e.POST("/backend-dashboard/backend/dao.groups.selectByName", backendManagerDaoGroupSelectByName)
	e.POST("/backend-dashboard/backend/dao.groups.insert", backendManagerDaoGroupsInsert)
	e.POST("/backend-dashboard/backend/dao.groups.deleteByName", backendManagerDaoGroupDeleteByName)
	e.POST("/backend-dashboard/backend/dao.groups.updateImplement", backendManagerDaoGroupUpdateImplement)
	e.POST("/backend-dashboard/backend/dao.interface.insert", backendManagerDaoInterfaceInsert)
	e.POST("/backend-dashboard/backend/dao.interface.delete", backendManagerDaoInterfaceDelete)
	e.POST("/backend-dashboard/backend/dao.config.insert", backendManagerDaoConfigInsert)
	e.POST("/backend-dashboard/backend/export", backendManagerExport)
	e.POST("/backend-dashboard/backend/specialSymbols", backendManagerSpecialSymbols)
	e.POST("/backend-dashboard/backend/mod.sign.selectAll", backendManagerModuleSignSelectAll)
	e.POST("/backend-dashboard/backend/mod.info.selectBySignId", backendManagerModuleInfoSelectBySignId)
}
