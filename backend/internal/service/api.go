package service

import (
	"backend/internal/dao"
	"backend/internal/model"
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
			ctx.String(400, "必须指定method")
		} else {
			params := mbody["params"]
			switch method {
			case INSERT:
				if params == nil {
					ctx.String(400, "做新增操作必须提交对象作为参数params")
					return
				}
				if tx, err := s.dao.BeginTx(ctx); err != nil {
					ctx.String(400, "事务开启失败：%v", err)
				} else {
					for _, m := range params.([]interface{}) {
						mdl := m.(map[string]interface{})
						if _, err := s.dao.InsertTx(tx, model.MODELS_NAME, mdl); err != nil {
							s.dao.Rollback(tx)
							ctx.String(400, "插入数据源失败：%v", err)
						}
					}
					if err := s.dao.Commit(tx); err != nil {
						ctx.String(400, "事务提交失败：%v", err)
					}
				}
			default:
				ctx.String(400, "未知Method：%s", method)
			}
		}
	})
}
