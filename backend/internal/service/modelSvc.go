package service

import (
	"backend/internal/dao"
	"context"
	"sync"
	"backend/api"
	"strings"
	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
	"unsafe"
	"fmt"
	"github.com/pkg/errors"
)

type ModelService struct {
	svc *Service
	dao *dao.MySqlDao
}

var modelSvcNewOnce sync.Once
var mdlSvc *ModelService

func NewModelService() *ModelService {
	modelSvcNewOnce.Do(func() {
		mdlSvc = &ModelService{
			svc: New(),
			dao: dao.NewSqlDao(),
		}
	})
	return mdlSvc
}

func InsModelService() *ModelService {
	return mdlSvc
}

func (s *ModelService) Ping(ctx context.Context) error {
	return s.svc.Ping(ctx)
}

func (s *ModelService) Close() {
	s.svc.Close()
}

const CREATE = "create"
const INSERT = "insert"
const DELETE = "delete"
const UPDATE = "update"
const SELECT = "select"

func (s *ModelService) HdlModelURL(ctx context.Context, req *api.GrpcReqs) (resp *api.GrpcResp, err error) {
	c := (*bm.Context)(unsafe.Pointer(&ctx))
	if req == nil {
		return
	} else if len(req.Method) == 0 {
		return nil, errors.New("必须指定method")
	} else if len(req.Table) == 0 {
		return nil, errors.New("必须在指定table")
	} else {
		// NOTE: 发现一个现象：新启动App后，如果数据源的表已存在，删除之后再Create，
		// 		 会发生Context超Deadline的错误。但同一次运行过程，Create几次都不会报错
		switch strings.ToLower(req.Method) {
		case CREATE:
			if tx, err := s.dao.BeginTx(c); err != nil {
				return nil, fmt.Errorf("事务开启失败：%v", err)
			} else if err := s.dao.Create(tx, req.Table, dao.ModelMap[req.Table]); err != nil {
				return nil, fmt.Errorf("创建表%s失败：%v", req.Table, err)
			} else if err := s.dao.CommitTx(tx); err != nil {
				return nil, fmt.Errorf("提交创建的表集合失败：%v", err)
			} else {
				return nil, nil
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
		default:
			return nil, fmt.Errorf("未知Method：%s", req.Method)
		}
	}
}