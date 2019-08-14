package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	pb "backend/api"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/server"
	"backend/internal/utils"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
)

// Service service.
type Service struct {
	ac *paladin.Map
	cc struct {
		Qiniu *utils.StorageConfig
	}
	dao      *dao.Dao
	gbuilder *ProjGenBuilder
}

// New new a service and return.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		ac:  ac,
		dao: dao.New(),
	}
	if err := paladin.Get("cdn.toml").UnmarshalTOML(&s.cc); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		panic(err)
	} else if err := s.dao.CreateTx(tx, model.MODELS_TABLE, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
		panic(err)
	} else if _, err := s.dao.DeleteTx(tx, model.MODELS_TABLE, "", []interface{}{}); err != nil {
		panic(err)
	} else if mp, err := utils.ToMap(pb.Model{
		Name:  "Nil",
		Type:  "struct",
		Props: []*pb.Prop{},
	}); err != nil {
		panic(err)
	} else if _, err := s.dao.InsertTx(tx, model.MODELS_TABLE, mp); err != nil {
		panic(err)
	} else if mp, err := utils.ToMap(pb.Model{
		Name: "IdenReqs",
		Type: "struct",
		Props: []*pb.Prop{
			&pb.Prop{
				Name: "id",
				Type: "int32",
			},
		},
	}); err != nil {
		panic(err)
	} else if _, err := s.dao.InsertTx(tx, model.MODELS_TABLE, mp); err != nil {
		panic(err)
	} else if mp, err := utils.ToMap(pb.Model{
		Name: "NameReqs",
		Type: "struct",
		Props: []*pb.Prop{
			&pb.Prop{
				Name: "name",
				Type: "string",
			},
		},
	}); err != nil {
		panic(err)
	} else if _, err := s.dao.InsertTx(tx, model.MODELS_TABLE, mp); err != nil {
		panic(err)
	} else if err := s.dao.CreateTx(tx, model.LINKS_TABLE, reflect.TypeOf((*pb.Link)(nil)).Elem()); err != nil {
		panic(err)
	} else if s.gbuilder, err = NewProjGenBuilder(s.dao, tx); err != nil {
		panic(err)
	} else if err := s.setupDAO(tx); err != nil {
		panic(err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		panic(err)
	}
	return s
}

func (s *Service) setupDAO(tx *sql.Tx) error {
	DaoCategoryType := reflect.TypeOf((*pb.DaoCategory)(nil)).Elem()
	DaoInterfaceType := reflect.TypeOf((*pb.DaoInterface)(nil)).Elem()
	DaoGroupType := reflect.TypeOf((*pb.DaoGroup)(nil)).Elem()
	DaoConfigType := reflect.TypeOf((*pb.DaoConfig)(nil)).Elem()

	if err := s.dao.DropTxByType(tx, model.DAO_INTERFACES_TABLE, DaoInterfaceType); err != nil {
		return err
	} else if err := s.dao.DropTxByType(tx, model.DAO_GROUPS_TABLE, DaoGroupType); err != nil {
		return err
	} else if err := s.dao.DropTxByType(tx, model.DAO_CATEGORIES_TABLE, DaoCategoryType); err != nil {
		return err
	} else if err := s.dao.DropTxByType(tx, model.DAO_CONFIGS_TABLE, DaoConfigType); err != nil {
		return err
	} else if err := s.dao.CreateTx(tx, model.DAO_CATEGORIES_TABLE, DaoCategoryType); err != nil {
		return err
	} else if err := s.dao.CreateTx(tx, model.DAO_INTERFACES_TABLE, DaoInterfaceType); err != nil {
		return err
	} else if err := s.dao.CreateTx(tx, model.DAO_GROUPS_TABLE, DaoGroupType); err != nil {
		return err
	} else if err := s.dao.CreateTx(tx, model.DAO_CONFIGS_TABLE, DaoConfigType); err != nil {
		return err
	} else if err := s.dao.SourceTx(tx, "/Users/zhaojiachen/Projects/backend-dashboard/backend/datas/dao_categories.sql"); err != nil {
		return err
	}
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_CATEGORIES_TABLE, map[string]interface{}{
	// 	"name": "database_notx", "desc": "数据库（无事务）", "lang": "golang",
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_CATEGORIES_TABLE, map[string]interface{}{
	// 	"name": "database_tx", "desc": "数据库（带事务）", "lang": "golang",
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "Close", "category": "database_tx", "desc": "关闭数据库",
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "Ping", "category": "database_tx", "desc": "测试数据库通讯状态",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "ctx", Type: "context.Context"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"context"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "BeginTx", "category": "database_tx", "desc": "开启事务",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "ctx", Type: "context.Context"},
	// 	},
	// 	"returns": []interface{}{"*sql.Tx", "error"},
	// 	"requires": []interface{}{"context", "github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "CommitTx", "category": "database_tx", "desc": "提交事务",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "RollbackTx", "category": "database_tx", "desc": "回滚事务",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "CreateTx", "category": "database_tx", "desc": "创建表",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "typ", Type: "reflect.Type"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"reflect", "github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "DropTx", "category": "database_tx", "desc": "删除表",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "typ", Type: "reflect.Type"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"reflect", "github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "ExecTx", "category": "database_tx", "desc": "执行SQL",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "sql", Type: "string"},
	// 		&pb.Prop{Name: "args", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "QueryTx", "category": "database_tx", "desc": "查询",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"[]map[string]interface{}", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "QueryOneTx", "category": "database_tx", "desc": "查询（单个）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"map[string]interface{}", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "QueryTxBySQL", "category": "database_tx", "desc": "查询（用SQL语句）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "ssql", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"[]map[string]interface{}", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "QueryTxOfOption", "category": "database_tx", "desc": "查询（带选项：GROUP BY、LIMIT、ORDER BY等）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 		&pb.Prop{Name: "options", Type: "[]string"},
	// 	},
	// 	"returns": []interface{}{"[]map[string]interface{}", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "QueryTxIdenCol", "category": "database_tx", "desc": "查询（指定列）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cols", Type: "[]string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"[]map[string]interface{}", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "InsertTx", "category": "database_tx", "desc": "插入（不检测存在与否）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "entry", Type: "map[string]interface{}"},
	// 	},
	// 	"returns": []interface{}{"int64", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "SaveTx", "category": "database_tx", "desc": "插入/保存（检测存在与否）",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 		&pb.Prop{Name: "entry", Type: "map[string]interface{}"},
	// 		&pb.Prop{Name: "commit", Type: "bool"},
	// 	},
	// 	"returns": []interface{}{"int64", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_INTERFACES_TABLE, map[string]interface{}{
	// 	"name": "DeleteTx", "category": "database_tx", "desc": "删除",
	// 	"params": []interface{}{
	// 		&pb.Prop{Name: "tx", Type: "*sql.Tx"},
	// 		&pb.Prop{Name: "table", Type: "string"},
	// 		&pb.Prop{Name: "cstr", Type: "string"},
	// 		&pb.Prop{Name: "carg", Type: "[]interface{}"},
	// 	},
	// 	"returns": []interface{}{"int64", "error"},
	// 	"requires": []interface{}{"github.com/bilibili/kratos/pkg/database/sql"},
	// }); err != nil {
	// 	return err
	// } else if _, err := s.dao.InsertTx(tx, model.DAO_CATEGORIES_TABLE, map[string]interface{}{
	// 	"name": "cache", "desc": "缓存", "lang": "golang",
	// }); err != nil {
	// 	return err
	// }
	return nil
}

func (s *Service) AppID() string {
	appID, _ := s.ac.Get("appID").String()
	return appID
}

func (s *Service) SwaggerFile() string {
	swagger, _ := s.ac.Get("swaggerFile").String()
	return swagger
}

func (s *Service) ModelsInsert(ctx context.Context, req *pb.Model) (*pb.Model, error) {
	mm := make(map[string]interface{})
	if bm, err := json.Marshal(req); err != nil {
		return nil, fmt.Errorf("转JSON失败：%v", err)
	} else if err := json.Unmarshal(bm, &mm); err != nil {
		return nil, fmt.Errorf("提交的请求参数格式有误：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if id, err := s.dao.InsertTx(tx, model.MODELS_TABLE, mm); err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	} else if mm, err = s.dao.QueryTxByID(tx, model.MODELS_TABLE, id); err != nil {
		return nil, fmt.Errorf("查询记录：%d失败：%v", id, err)
	} else if _, err := GenModelApiInfo(s.dao, tx, mm, nil); err != nil {
		return nil, fmt.Errorf("生成模块接口失败：%v", err)
	} else {
		if len(req.Type) == 0 || req.Type == "model" {
			if aryMap, err := utils.ToMap(pb.Model{
				Name:  req.Name + "Array",
				Type:  "struct",
				Model: req.Name,
				Props: []*pb.Prop{
					&pb.Prop{
						Name: utils.ToPlural(strings.ToLower(req.Name)),
						Type: "repeated " + req.Name,
					},
				},
			}); err != nil {
				return nil, fmt.Errorf("生成集合结构失败：%v", err)
			} else if _, err := s.dao.InsertTx(tx, model.MODELS_TABLE, aryMap); err != nil {
				return nil, fmt.Errorf("插入集合结构失败：%v", err)
			}
		}
		if err := s.dao.CommitTx(tx); err != nil {
			return nil, fmt.Errorf("提交插入事务失败：%v", err)
		}
		req.Id = id
		return req, nil
	}
}

func (s *Service) ModelsDelete(ctx context.Context, req *pb.NameID) (*pb.Model, error) {
	conds := "`name`=?"
	argus := []interface{}{req.Name}
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if res, err := s.dao.QueryTx(tx, model.MODELS_TABLE, conds, argus); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else if _, err := s.dao.DeleteTx(tx, model.MODELS_TABLE, conds, argus); err != nil {
		return nil, fmt.Errorf("删除数据库记录失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("数据库提交失败：%v", err)
	} else if bd, err := json.Marshal(res[0]); err != nil {
		return nil, fmt.Errorf("转JSON字节码失败：%v", err)
	} else if resp, err := utils.UnmarshalJSON(bd, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
		return nil, fmt.Errorf("转Model对象失败：%v", err)
	} else {
		return resp.(*pb.Model), nil
	}
}

func (s *Service) ModelsUpdate(ctx context.Context, req *pb.Model) (*pb.Empty, error) {
	if req.Id == 0 {
		return nil, errors.New("需要给出要更新的模型ID")
	} else if bytes, err := json.Marshal(req); err != nil {
		return nil, fmt.Errorf("转JSON字节码失败：%v", err)
	} else if pmp, err := utils.UnmarshalJSON(bytes, reflect.TypeOf((*map[string]interface{})(nil)).Elem()); err != nil {
		return nil, fmt.Errorf("转Map失败：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if _, err := s.dao.UpdateTxByID(tx, model.MODELS_TABLE, *(pmp.(*map[string]interface{}))); err != nil {
		return nil, fmt.Errorf("更新模型失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("数据库提交失败：%v", err)
	} else {
		return &pb.Empty{}, nil
	}
}

func (s *Service) ModelsSelectAll(ctx context.Context, req *pb.TypeIden) (*pb.ModelArray, error) {
	cstr := ""
	carg := make([]interface{}, 0)
	if req.Type != "" {
		cstr = "`type`=?"
		carg = append(carg, req.Type)
	}
	if res, err := s.dao.Query(ctx, model.MODELS_TABLE, cstr, carg); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.ModelArray)
		for _, entry := range res {
			if bdata, err := json.Marshal(entry); err != nil {
				return nil, fmt.Errorf("转JSON字节码失败：%v", err)
			} else if mdl, err := utils.UnmarshalJSON(bdata, reflect.TypeOf((*pb.Model)(nil)).Elem()); err != nil {
				return nil, fmt.Errorf("转Model对象失败：%v", err)
			} else {
				resp.Models = append(resp.Models, mdl.(*pb.Model))
			}
		}
		return resp, nil
	}
}

func (s *Service) ModelsSelectByName(context.Context, *pb.NameID) (*pb.Model, error) {
	// TODO:
	return nil, nil
}

func (s *Service) StructsSelectAllBases(context.Context, *pb.Empty) (*pb.NameArray, error) {
	return &pb.NameArray{
		Names: []string{"Nil", "IdenReqs", "NameReqs"},
	}, nil
}

func (s *Service) LinksInsert(ctx context.Context, req *pb.Link) (*pb.Link, error) {
	if bm, err := json.Marshal(req); err != nil {
		return nil, fmt.Errorf("转JSON失败：%v", err)
	} else if mm, err := utils.UnmarshalJSON(bm, reflect.TypeOf((*map[string]interface{})(nil)).Elem()); err != nil {
		return nil, fmt.Errorf("提交的请求参数格式有误：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if id, err := s.dao.InsertTx(tx, model.LINKS_TABLE, *(mm.(*map[string]interface{}))); err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交插入事务失败：%v", err)
	} else {
		req.Id = id
		return req, nil
	}
}

func (s *Service) LinksSelectAll(ctx context.Context, req *pb.Empty) (*pb.LinkArray, error) {
	if res, err := s.dao.Query(ctx, model.LINKS_TABLE, "", []interface{}{}); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.LinkArray)
		for _, entry := range res {
			if bdata, err := json.Marshal(entry); err != nil {
				return nil, fmt.Errorf("转JSON字节码失败：%v", err)
			} else if mdl, err := utils.UnmarshalJSON(bdata, reflect.TypeOf((*pb.Link)(nil)).Elem()); err != nil {
				return nil, fmt.Errorf("转Link对象失败：%v", err)
			} else {
				resp.Links = append(resp.Links, mdl.(*pb.Link))
			}
		}
		return resp, nil
	}
}

func (s *Service) LinksDeleteBySymbol(context.Context, *pb.SymbolID) (*pb.Link, error) {
	// TODO:
	return nil, nil
}

func (s *Service) ApisSelectByName(ctx context.Context, req *pb.NameID) (*pb.ApiInfo, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if api, err := ApiInfoFmDbByTx(s.dao, tx, req.Name); err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交事务失败：%v", err)
	} else {
		return api, nil
	}
}

func (s *Service) ApisSelectAll(ctx context.Context, req *pb.Empty) (*pb.ApiInfoArray, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if res, err := s.dao.QueryTx(tx, model.API_INFO_TABLE, "", []interface{}{}); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.ApiInfoArray)
		for _, entry := range res {
			if api, err := CvtApiInfoFmMap(s.dao, tx, entry); err != nil {
				return nil, fmt.Errorf("API map转为ApiInfo失败：%v", err)
			} else {
				resp.Infos = append(resp.Infos, api)
			}
		}
		if err := s.dao.CommitTx(tx); err != nil {
			return nil, fmt.Errorf("提交事务失败：%v", err)
		} else {
			return resp, nil
		}
	}
}

func (s *Service) ApisInsert(ctx context.Context, req *pb.ApiInfo) (*pb.ApiInfo, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if _, err := ApiInfoToDbByTx(s.dao, tx, req); err != nil {
		return nil, fmt.Errorf("插入数据库失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交插入事务失败：%v", err)
	} else {
		return req, nil
	}
}

func (s *Service) ApisDeleteByName(ctx context.Context, req *pb.NameID) (*pb.ApiInfo, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if api, err := ApiInfoFmDbByTx(s.dao, tx, req.Name); err != nil {
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	} else if _, err := s.dao.DeleteTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{req.Name}); err != nil {
		return nil, fmt.Errorf("删除接口失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交事务失败：%v", err)
	} else {
		return api, nil
	}
}

func (s *Service) StepsInsert(ctx context.Context, req *pb.StepReqs) (*pb.Empty, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if id, err := OperStepToDbByTx(s.dao, tx, req.OperStep); err != nil {
		// 添加OperStep
		return nil, fmt.Errorf("插入步骤失败：%v", err)
	} else if mapi, err := s.dao.QueryOneTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{
		req.OperStep.ApiName,
	}); err != nil {
		// 获取当前API的所有步骤
		return nil, fmt.Errorf("查询接口信息失败：%v", err)
	} else if steps := strings.Split(mapi["steps"].(string), ","); false {
		// 修改Api的步骤顺序
		return nil, nil
	} else if rear := append([]string{}, steps[req.Index:]...); false {
		return nil, nil
	} else if steps = append(steps[:req.Index], strconv.Itoa(int(id))); false {
		return nil, nil
	} else if steps = append(steps, rear...); false {
		return nil, nil
	} else if _, err := s.dao.SaveTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{
		req.OperStep.ApiName,
	}, map[string]interface{}{
		"steps": strings.Trim(strings.Join(steps, ","), ","),
	}, false); err != nil {
		return nil, fmt.Errorf("保存接口信息失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交保存事务失败：%v", err)
	} else {
		return &pb.Empty{}, nil
	}
}

func delStepFmApiSteps(steps string, delId string) string {
	ary := strings.Split(steps, ",")
	for i := 0; i < len(ary); i++ {
		if ary[i] == delId {
			return strings.Join(append(ary[:i], ary[i+1:]...), ",")
		}
	}
	return steps
}

func (s *Service) StepsDelete(ctx context.Context, req *pb.DelStepReqs) (*pb.Empty, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if _, err := s.dao.DeleteTxByID(tx, model.OPER_STEP_TABLE, req.StepId); err != nil {
		return nil, fmt.Errorf("删除步骤失败：%v", err)
	} else if mps, err := s.dao.QueryTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{req.ApiName}); err != nil {
		return nil, fmt.Errorf("查询接口失败：%v", err)
	} else if len(mps) == 0 {
		return &pb.Empty{}, nil
	} else if mp := mps[0]; false {
		return nil, nil
	} else if _, err := s.dao.SaveTx(tx, model.API_INFO_TABLE, "`name`=?", []interface{}{req.ApiName}, map[string]interface{}{
		"steps": delStepFmApiSteps(mp["steps"].(string), strconv.Itoa(int(req.StepId))),
	}, true); err != nil {
		return nil, fmt.Errorf("更新接口的流程失败：%v", err)
	} else {
		return &pb.Empty{}, nil
	}
}

// 这是添加步骤模板，可以通过设置apiName来指定要插入的接口，但只能追加到api流程的最后
// 如果需要插入到流程中间，则需要使用StepsInsert
func (s *Service) OperStepsInsert(context.Context, *pb.OperStep) (*pb.OperStep, error) {
	return nil, nil
}

func (s *Service) OperStepsSelectTemp(ctx context.Context, req *pb.Empty) (*pb.OperStepArray, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if res, err := s.dao.QueryTx(tx, model.OPER_STEP_TABLE, "`api_name` IS NULL", nil); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.OperStepArray)
		for _, entry := range res {
			resp.Steps = append(resp.Steps, CvtOperStepFmMap(entry))
		}
		if err := s.dao.CommitTx(tx); err != nil {
			return nil, fmt.Errorf("提交事务失败：%v", err)
		} else {
			return resp, nil
		}
	}
	return nil, nil
}

func (s *Service) DaoGroupsSelectAll(ctx context.Context, req *pb.Empty) (*pb.DaoGroupArray, error) {
	if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if res, err := s.dao.QueryTx(tx, model.DAO_GROUPS_TABLE, "", nil); err != nil {
		return nil, fmt.Errorf("查询数据库失败：%v", err)
	} else {
		resp := new(pb.DaoGroupArray)
		for _, entry := range res {
			if grp, err := utils.ToObj(entry, reflect.TypeOf((*pb.DaoGroup)(nil)).Elem()); err != nil {
				return nil, fmt.Errorf("JSON转成DAO组失败：%v", err)
			} else {
				resp.Groups = append(resp.Groups, grp.(*pb.DaoGroup))
			}
		}
		if err := s.dao.CommitTx(tx); err != nil {
			return nil, fmt.Errorf("提交事务失败：%v", err)
		} else {
			return resp, nil
		}
	}
	return nil, nil
}

func (s *Service) DaoGroupsInsert(ctx context.Context, req *pb.DaoGroup) (*pb.DaoGroup, error) {
	if mp, err := utils.ToMap(req); err != nil {
		return nil, fmt.Errorf("转成JSON失败：%v", err)
	} else if tx, err := s.dao.BeginTx(ctx); err != nil {
		return nil, fmt.Errorf("开启事务失败：%v", err)
	} else if _, err := s.dao.InsertTx(tx, model.DAO_GROUPS_TABLE, mp); err != nil {
		return nil, fmt.Errorf("DAO组插入数据库失败：%v", err)
	} else if err := s.dao.CommitTx(tx); err != nil {
		return nil, fmt.Errorf("提交事务失败：%v", err)
	} else {
		return req, nil
	}
}

func (s *Service) Export(ctx context.Context, req *pb.ExpOptions) (*pb.UrlResp, error) {
	if wsPath, err := s.ac.Get("workspace").String(); err != nil {
		return nil, fmt.Errorf("配置文件中未定义工作区目录：%v", err)
	} else if bin, err := time.Now().MarshalBinary(); err != nil {
		return nil, fmt.Errorf("生成临时文件夹名失败：%v", err)
	} else if cchName := fmt.Sprintf("%x", md5.Sum(bin)); false {
		return nil, nil
	} else if cchPath := path.Join(wsPath, "cache", cchName); false {
		return nil, nil
	} else if err := os.MkdirAll(cchPath, 0755); err != nil {
		return nil, fmt.Errorf("创建临时文件夹：%s失败：%v", cchPath, err)
	} else if tmpPath := path.Join(wsPath, "template", req.Type); false {
		return nil, nil
	} else if pathName := path.Join(cchPath, req.Type); false {
		return nil, nil
	} else if utils.CopyFolder(tmpPath, pathName); false {
		return nil, nil
	} else if err := s.gbuilder.Build(ctx, req, pathName).Adjust(ctx); err != nil {
		return nil, fmt.Errorf("编辑项目失败：%v", err)
	} else if wsFile, err := os.Open(pathName); err != nil {
		return nil, fmt.Errorf("工作区目录有误，打开失败：%v", err)
	} else if zipPath := path.Join(cchPath, req.Name); false {
		return nil, nil
	} else if err := utils.Compress([]*os.File{wsFile}, zipPath); err != nil {
		wsFile.Close()
		return nil, fmt.Errorf("压缩项目失败：%v", err)
	} else if url, err := utils.Upload(zipPath, *s.cc.Qiniu); err != nil {
		wsFile.Close()
		return nil, fmt.Errorf("上传项目包失败：%v", err)
		//} else if err := os.RemoveAll(cchPath); err != nil {
		//	wsFile.Close()
		//	return nil, fmt.Errorf("删除临时文件夹失败：%v", err)
	} else {
		wsFile.Close()
		return &pb.UrlResp{Url: url}, nil
	}
}

func (s *Service) SpecialSymbols(context.Context, *pb.Empty) (*pb.SymbolsResp, error) {
	return &pb.SymbolsResp{
		Values: pb.SpcSymbol_value,
		Names:  pb.SpcSymbol_name,
	}, nil
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
	// 注销服务
	if cli, err := server.RegisterService(); err != nil {
		panic(err)
	} else if _, err := cli.Cancel(context.Background(), &pb.IdenSvcReqs{AppID: s.AppID()}); err != nil {
		panic(err)
	}
}
