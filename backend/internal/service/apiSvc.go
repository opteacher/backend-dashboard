package service

import (
	"backend/internal/dao"
	"backend/utils"
	"context"
	"reflect"
	"sync"
	"time"
	"fmt"

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

func (s *ApiService) AddModelAPI(g *bm.RouterGroup, mname string, methods []string) {
	g.POST("/"+mname, func(ctx *bm.Context) {
		if tm, ok := ctx.Deadline(); ok {
			fmt.Println(tm.Sub(time.Now()))
		}
		mbody := utils.GetReqBody(ctx)
		if mbody == nil {
			return
		} else if method, exs := mbody["method"]; !exs {
			ctx.String(400, "必须指定method")
		} else if params, exs := mbody["params"]; !exs {
			// NOTE: 参数列表对应增删改查四个操作，对于删改查：其前两个参数为条件限定参数
			//       第一个参数为条件限定字符串；第二个参数为条件参数列表，之后再跟实例。
			//       增操作除外，直接是实例列表
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
				} else if err := s.dao.Create(tx, mname, dao.ModelMap[mname]); err != nil {
					ctx.String(400, "创建表%s失败：%v", mname, err)
				} else if err := s.dao.CommitTx(tx); err != nil {
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
						ctx.String(400, "提交数据源失败：%v", err)
					} else {
						ctx.JSON(respData, nil)
					}
				}
			case DELETE:
				pamlst := params.([]interface{})
				if len(pamlst) < 2 {
					ctx.String(400, "需要指定要删除的条件字符串和参数列表")
				} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
					ctx.String(400, "第一个参数必须是string类型（条件字符串），但收到的是：%s", reflect.TypeOf(pamlst[0]))
				} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
					ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
				} else if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "开启事务失败：%v", err)
				} else if n, err := s.dao.DeleteTx(tx, mname, pamlst[0].(string), pamlst[1].([]interface{})); err != nil {
					ctx.String(400, "删除数据源错误：%v", err)
				} else if err := s.dao.CommitTx(tx); err != nil {
					ctx.String(400, "提交数据源失败：%v", err)
				} else {
					ctx.JSON(n, nil)
				}
			case UPDATE:
				pamlst := params.([]interface{})
				if len(pamlst) < 3 {
					ctx.String(400, "需要指定要更新的条件字符串和参数列表")
				} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
					ctx.String(400, "第一个参数必须是string类型（条件字符串），但收到的是：%s", reflect.TypeOf(pamlst[0]))
				} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
					ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
				} else if tx, err := s.dao.BeginTx(c); err != nil {
					ctx.String(400, "开启事务失败：%v", err)
				} else {
					condStr := pamlst[0].(string)
					condArgs := pamlst[1].([]interface{})
					var updNum int64
					for i := 2; i < len(pamlst); i++ {
						if obj := pamlst[i]; !reflect.TypeOf(obj).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
							s.dao.RollbackTx(tx)
							ctx.String(400, "参数为元组，必须指定为object")
							return
						} else if n, err := s.dao.SaveTx(tx, mname, condStr, condArgs, obj.(map[string]interface{}), false); err != nil {
							ctx.String(400, "更新数据源错误：%v", err)
							return
						} else {
							updNum += n
						}
					}
					if err := s.dao.CommitTx(tx); err != nil {
						ctx.String(400, "提交数据源失败：%v", err)
					} else {
						ctx.JSON(updNum, nil)
					}
				}
			case SELECT:
				pamlst := params.([]interface{})
				if len(pamlst) < 2 {
					ctx.String(400, "需要指定要查询的条件字符串和参数列表")
				} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
					ctx.String(400, "第一个参数必须是string类型（条件字符串），但收到的是：%s", reflect.TypeOf(pamlst[0]))
				} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
					ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
				} else if res, err := s.dao.Query(c, mname, pamlst[0].(string), pamlst[1].([]interface{})); err != nil {
					ctx.String(400, "查询数据源错误：%v", err)
				} else {
					ctx.JSON(res, nil)
				}
			default:
				ctx.String(400, "未知Method：%s", method)
			}
		}
	})
}
