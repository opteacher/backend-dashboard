package service

import (
	"backend/internal/dao"
	"backend/utils"
	"context"
	"sync"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type ApiService struct {
	dao      *dao.MySqlDao
	handlers map[string]func(*bm.Context)
}

var newOnce sync.Once
var svc *ApiService

func NewApiService() *ApiService {
	newOnce.Do(func() {
		svc = &ApiService{
			dao:      dao.NewSqlDao(),
			handlers: make(map[string]func(*bm.Context)),
		}
	})
	return svc
}

func InsApiService() *ApiService {
	return svc
}

func (s *ApiService) Ping(ctx context.Context) error {
	return s.dao.Ping(ctx)
}

const INSERT = "insert"
const DELETE = "delete"
const UPDATE = "update"
const SELECT = "select"

func (s *ApiService) AddModelAPI(g *bm.RouterGroup, mname string, methods []string) {
	g.POST("/"+mname, func(ctx *bm.Context) {
		mbody := utils.GetReqBody(ctx)
		if mbody == nil {
			return
		} else if method, exs := mbody["method"]; !exs {
			ctx.String(400, "必须指定Method")
		} else {
			params := mbody["params"]
			switch method {
			case INSERT:
			default:
				ctx.String(400, "未知Method：%s", method)
			}
		}
	})
}
