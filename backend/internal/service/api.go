package service

import (
	"backend/internal/dao"
	"backend/internal/model"
	"backend/utils"
	"context"
	"reflect"
	"sync"
	"time"

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
			// NOTE: 发现一个现象：新启动App后，如果数据源的表已存在，删除之后再Create，
			// 		 会发生Context超Deadline的错误。但同一次运行过程，Create几次都不会报错
			c, cancel := context.WithDeadline(ctx, time.Now().Add(60*time.Second))
			defer cancel()
			switch method {
			case CREATE:
				if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "事务开启失败：%v", err)
				} else if err := s.dao.Create(tx, model.MODELS_NAME, reflect.TypeOf((*model.Model)(nil)).Elem()); err != nil {
					s.dao.RollbackTx(tx)
					ctx.String(400, "创建表%s失败：%v", model.MODELS_NAME, err)
				} else if err := s.dao.CommitTx(tx); err != nil {
					s.dao.RollbackTx(tx)
					ctx.String(400, "提交创建的表集合失败：%v", err)
				} else {
					ctx.String(200, "创建表成功")
				}
			case INSERT:
				pamlst := params.([]interface{})
				if len(pamlst) < 1 {
					ctx.String(400, "需要指定要插入的元组")
				} else if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "开启事务失败：%v", err)
				} else {
					var respData []interface{}
					for _, obj := range pamlst {
						if !reflect.TypeOf(obj).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
							s.dao.RollbackTx(tx)
							ctx.String(400, "参数为元组，必须指定为object")
							return
						} else {
							objmap := obj.(map[string]interface{})
							if id, err := s.dao.InsertTx(tx, mname, objmap); err != nil {
								ctx.String(400, "插入数据源错误：%v", err)
								return
							} else {
								objmap["id"] = id
								respData = append(respData, objmap)
							}
						}
					}
					if err := s.dao.CommitTx(tx); err != nil {
						s.dao.RollbackTx(tx)
						ctx.String(400, "提交数据源失败：%v", err)
					} else {
						ctx.JSON(respData, nil)
					}
				}
			case DELETE:
				pamlst := params.([]interface{})
				if len(pamlst) < 1 {
					ctx.String(400, "需要指定要删除的元组")
				} else if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "开启事务失败：%v", err)
				} else {
					for _, obj := range pamlst {
						if !reflect.TypeOf(obj).ConvertibleTo(reflect.TypeOf((*int64)(nil)).Elem()) {
							s.dao.RollbackTx(tx)
							ctx.String(400, "参数为id，必须指定为int")
							return
						} else if _, err := s.dao.DeleteTxByID(tx, mname, int64(obj.(float64))); err != nil {
							ctx.String(400, "删除数据源错误：%v", err)
							return
						}
					}
					if err := s.dao.CommitTx(tx); err != nil {
						s.dao.RollbackTx(tx)
						ctx.String(400, "提交数据源失败：%v", err)
					} else {
						ctx.JSON(len(pamlst), nil)
					}
				}
			case UPDATE:
				pamlst := params.([]interface{})
				if len(pamlst) < 1 {
					ctx.String(400, "需要指定要更新的元组")
				} else if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "开启事务失败：%v", err)
				} else {
					var resp []map[string]interface{}
					for _, obj := range pamlst {
						if !reflect.TypeOf(obj).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
							s.dao.RollbackTx(tx)
							ctx.String(400, "参数为元组，必须指定为object")
							return
						} else if res, err := s.dao.UpdateTxByID(tx, mname, obj.(map[string]interface{})); err != nil {
							ctx.String(400, "更新数据源错误：%v", err)
							return
						} else {
							resp = append(resp, res)
						}
					}
					if err := s.dao.CommitTx(tx); err != nil {
						s.dao.RollbackTx(tx)
						ctx.String(400, "提交数据源失败：%v", err)
					} else {
						ctx.JSON(resp, nil)
					}
				}
			default:
				ctx.String(400, "未知Method：%s", method)
			}
		}
	})
}
