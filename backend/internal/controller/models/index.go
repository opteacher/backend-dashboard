package modelsctl

import (
	"backend/internal/model"
	"encoding/json"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type ModelsController struct {
}

func New() (ctl *ModelsController) {
	return &ModelsController{}
}

func (ctl *ModelsController) HandlePost(ctx *bm.Context) {
	if b, exists := ctx.Get("body"); !exists {
		ctx.String(400, "未给出参数")
	} else {
		var m model.Model
		if str, err := json.Marshal(b); err != nil {
			ctx.String(400, "json合并错误：%v", err)
		} else if err := json.Unmarshal(str, &m); err != nil {
			ctx.String(400, "json解析错误：%v", err)
		} else {
			ctx.JSON(m, nil)
		}
	}
}

func (ctl *ModelsController) HandleDelete(ctx *bm.Context) {
	if mid, exists := ctx.Get("mid"); exists {
		ctx.JSON(mid, nil)
	}
}
