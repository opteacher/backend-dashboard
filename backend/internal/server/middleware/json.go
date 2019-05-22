package globalmw

import (
	"encoding/json"
	"io/ioutil"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type jsonParser struct {
}

func (p *jsonParser) ServeHTTP(ctx *bm.Context) {
	// 读取参数
	var body []byte
	var err error
	if body, err = ioutil.ReadAll(ctx.Request.Body); err != nil {
		ctx.String(400, "读取请求体错误：%v", err)
		return
	}
	defer ctx.Request.Body.Close()

	// 解析参数
	var value interface{}
	if err = json.Unmarshal(body, &value); err != nil {
		ctx.String(400, "解析请求体错误：%v", err)
		return
	}
	ctx.Set("body", value)
	ctx.Next()
}

func ParseJSON() (p *jsonParser) {
	return &jsonParser{}
}
