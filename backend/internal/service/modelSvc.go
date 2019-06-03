package service

import (
	"backend/internal/dao"
	"backend/utils"
	"context"
	"sync"
	"time"
	"backend/api"
	"strings"
	"encoding/json"
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
	if req == nil {
		return
	} else if len(req.Method) == 0 {
		return utils.WrapError("必须指定method"), nil
	} else if len(req.Params) == 0 {
		// NOTE: 参数列表对应增删改查四个操作，第一个参数必须是操作的表名。之后，对于删改查：
		//       其前两个参数为条件限定参数第二个参数为条件限定字符串；第三个参数为条件参数列
		//       表，之后再跟实例。增操作除外，直接是实例列表
		return utils.WrapError("必须在params中指定操作的表名"), nil
	} else {
		mname := string(req.Params[0])
		// NOTE: 发现一个现象：新启动App后，如果数据源的表已存在，删除之后再Create，
		// 		 会发生Context超Deadline的错误。但同一次运行过程，Create几次都不会报错
		c, cancel := context.WithDeadline(ctx, time.Now().Add(60*time.Second))
		defer cancel()
		switch strings.ToLower(req.Method) {
		case CREATE:
			if tx, err := s.dao.BeginTx(c); err != nil {
				return utils.WrapError("事务开启失败：%v", err), nil
			} else if err := s.dao.Create(tx, mname, dao.ModelMap[mname]); err != nil {
				return utils.WrapError("创建表%s失败：%v", mname, err), nil
			} else if err := s.dao.CommitTx(tx); err != nil {
				return utils.WrapError("提交创建的表集合失败：%v", err), nil
			} else {
				return &api.GrpcResp{
					Status: 200,
					Message: "创建表成功",
				}, nil
			}
		case INSERT:
			if len(req.Params) < 2 {
				return utils.WrapError("需要指定要插入的元组"), nil
			} else if tx, err := s.dao.BeginTx(c); err != nil {
				return utils.WrapError("开启事务失败：%v", err), nil
			} else {
				var respData []interface{}
				for i := 1;  i < len(req.Params); i++ {
					objmap := make(map[string]interface{})
					if err := json.Unmarshal(req.Params[i], &objmap); err != nil {
						s.dao.RollbackTx(tx)
						return utils.WrapError("解析待插入的元祖失败：%v", err), nil
					} else if id, err := s.dao.InsertTx(tx, mname, objmap); err != nil {
						return utils.WrapError("插入数据源错误：%v", err), nil
					} else {
						objmap["id"] = id
						respData = append(respData, objmap)
					}
				}
				if err := s.dao.CommitTx(tx); err != nil {
					return utils.WrapError("提交数据源失败：%v", err), nil
				} else if jsonResp, err := json.Marshal(respData); err != nil {
					return utils.WrapError("转向JSON失败：%v", err), nil
				} else {
					return &api.GrpcResp{
						Status: 200,
						Message: "插入成功",
						Data: jsonResp,
					}, nil
				}
			}
		//case DELETE:
		//	if len(req.Params) < 3 {
		//		return utils.WrapError("需要指定要删除的条件字符串和参数列表"), nil
		//	} else if condStr := string(req.Params[1]); len(condStr) == 0 {
		//		return utils.WrapError("第二个参数（条件字符串）为空"), nil
		//	} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
		//		ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
		//	} else if tx, err := s.dao.BeginTx(c); err != nil {
		//		ctx.String(400, "开启事务失败：%v", err)
		//	} else if n, err := s.dao.DeleteTx(tx, mname, pamlst[0].(string), pamlst[1].([]interface{})); err != nil {
		//		ctx.String(400, "删除数据源错误：%v", err)
		//	} else if err := s.dao.CommitTx(tx); err != nil {
		//		ctx.String(400, "提交数据源失败：%v", err)
		//	} else {
		//		ctx.JSON(n, nil)
		//	}
		//case UPDATE:
		//	pamlst := params.([]interface{})
		//	if len(pamlst) < 3 {
		//		ctx.String(400, "需要指定要更新的条件字符串和参数列表")
		//	} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
		//		ctx.String(400, "第一个参数必须是string类型（条件字符串），但收到的是：%s", reflect.TypeOf(pamlst[0]))
		//	} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
		//		ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
		//	} else if tx, err := s.dao.BeginTx(c); err != nil {
		//		ctx.String(400, "开启事务失败：%v", err)
		//	} else {
		//		condStr := pamlst[0].(string)
		//		condArgs := pamlst[1].([]interface{})
		//		var updNum int64
		//		for i := 2; i < len(pamlst); i++ {
		//			if obj := pamlst[i]; !reflect.TypeOf(obj).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
		//				s.dao.RollbackTx(tx)
		//				ctx.String(400, "参数为元组，必须指定为object")
		//				return
		//			} else if n, err := s.dao.SaveTx(tx, mname, condStr, condArgs, obj.(map[string]interface{}), false); err != nil {
		//				ctx.String(400, "更新数据源错误：%v", err)
		//				return
		//			} else {
		//				updNum += n
		//			}
		//		}
		//		if err := s.dao.CommitTx(tx); err != nil {
		//			ctx.String(400, "提交数据源失败：%v", err)
		//		} else {
		//			ctx.JSON(updNum, nil)
		//		}
		//	}
		//case SELECT:
		//	pamlst := params.([]interface{})
		//	if len(pamlst) < 2 {
		//		ctx.String(400, "需要指定要查询的条件字符串和参数列表")
		//	} else if !reflect.TypeOf(pamlst[0]).ConvertibleTo(reflect.TypeOf((*string)(nil)).Elem()) {
		//		ctx.String(400, "第一个参数必须是string类型（条件字符串），但收到的是：%s", reflect.TypeOf(pamlst[0]))
		//	} else if !reflect.TypeOf(pamlst[1]).ConvertibleTo(reflect.TypeOf((*[]interface{})(nil)).Elem()) {
		//		ctx.String(400, "第二个参数必须是[]interface{}（条件参数），但收到的是：%s", reflect.TypeOf(pamlst[1]))
		//	} else if res, err := s.dao.Query(c, mname, pamlst[0].(string), pamlst[1].([]interface{})); err != nil {
		//		ctx.String(400, "查询数据源错误：%v", err)
		//	} else {
		//		ctx.JSON(res, nil)
		//	}
		default:
			return utils.WrapError("未知Method：%s", req.Method), nil
		}
	}
}