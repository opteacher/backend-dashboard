package globalmw

import (
	"strings"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type corsSetup struct {
}

func (s *corsSetup) ServeHTTP(ctx *bm.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization,Accept,X-Requested-With")
	if strings.ToUpper(ctx.Request.Method) == "OPTIONS" {
		ctx.Status(200)
	} else {
		ctx.Next()
	}
}

func SetupCORS() (p *corsSetup) {
	return &corsSetup{}
}
