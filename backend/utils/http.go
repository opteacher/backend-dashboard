package utils

import (
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"backend/api"
	"fmt"
)

func GetReqBody(ctx *bm.Context) map[string]interface{} {
	if body, exists := ctx.Get("body"); !exists {
		ctx.String(400, "未给出参数")
	} else if mbody, succeed := body.(map[string]interface{}); !succeed {
		ctx.String(400, "参数无法转成map，检查是否是合法的JSON")
	} else {
		return mbody
	}
	return nil
}

func WrapError(msg string, params ...interface{}) *api.GrpcResp {
	return &api.GrpcResp {
		Status: 400,
		Message: fmt.Sprintf(msg, params...),
	}
}