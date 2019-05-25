package service

import (
	"backend/internal/dao"
	"backend/internal/model"
	"backend/utils"
	"context"
	"sync"
	"reflect"

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
	return nil
}

const CREATE = "create"
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
		} else if params, exs := mbody["params"]; !exs {
			ctx.String(400, "必须指定params")
		} else {
			switch method {
			case CREATE:
				if err := s.dao.Create(ctx, model.MODELS_NAME, reflect.TypeOf((*model.Model)(nil)).Elem()); err != nil {
					s.dao.Rollback()
					ctx.String(400, "创建表%s失败：%v", model.MODELS_NAME, err)
				} else if err := s.dao.Commit(); err != nil {
					s.dao.Rollback()
					ctx.String(400, "提交创建的表集合失败：%v", err)
				} else {
					ctx.String(200, "创建表成功")
				}
			case INSERT:
				pamlst := params.([]interface{})
				if len(pamlst) < 2 {
					ctx.String(400, "需要指定域名和要插入的元组")
				} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
					ctx.String(400, "第一个参数为域名，必须指定为string")
				} else {
					table := pamlst[0].(string)
					if tx, err := s.dao.BeginTx(ctx); err != nil {
						ctx.String(400, "开启事务失败：%v", err)
					} else {
						for i := 1; i < len(pamlst); i++ {
							if !reflect.TypeOf(pamlst[i]).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
								s.dao.RollbackTx(tx)
								ctx.String(400, "第二个参数为元组，必须指定为object")
							} else {
								objmap := pamlst[i].(map[string]interface{})
								if _, err := s.dao.InsertTx(tx, table, objmap); err != nil {
									s.dao.RollbackTx(tx)
									ctx.String(400, "插入数据源错误：%v", err)
								}
							}
						}
						if err := s.dao.CommitTx(tx); err != nil {
							ctx.String(400, "提交数据源失败：%v", err)
						} else {
							ctx.String(200, "插入数据源成功")
						}
					}
				}
			default:
				ctx.String(400, "未知Method：%s", method)
			}
		}
	})
}
