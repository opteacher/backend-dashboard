package service

import (
	"backend/internal/dao"
	"fmt"
	"sync"

	bm "github.com/bilibili/kratos/pkg/net/http/blademaster"
)

type ApiService struct {
	dao *dao.Dao
}

var newOnce sync.Once
var svc *ApiService

func NewApiService() *ApiService {
	newOnce.Do(func() {
		svc = &ApiService{
			dao: dao.New(),
		}
	})
	return svc
}

func InsApiService() *ApiService {
	return svc
}

const INSERT = "insert"
const DELETE = "delete"
const UPDATE = "update"
const SELECT = "select"

func (s *ApiService) AddModelAPI(g *bm.RouterGroup, mname string, methods []string) ([]string, error) {
	var ms []string

	for _, method := range methods {
		path := fmt.Sprintf("/%s/%s", method, mname)
		ms = append(ms, path)
		var handler func(*bm.Context)
		switch method {
		case INSERT:
			handler = func(ctx *bm.Context) {
				if body, exists := ctx.Get("body"); !exists {
					ctx.String(400, "未给出参数")
				} else if mbody, succeed := body.(map[string]interface{}); !succeed {
					ctx.String(400, "参数无法转成map，检查是否是合法的JSON")
				} else {
					if id, exists := mbody["id"]; !exists {
						s.dao.Insert(ctx, mname, mbody)
					} else {
						s.dao.Save(ctx, mname, "`id`=?", []interface{}{id}, mbody)
					}
					ctx.JSON(mbody, nil)
				}
			}
		}
		g.POST(path, handler)
	}
	return ms, nil
}
