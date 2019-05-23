package dao

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/bilibili/kratos/pkg/cache/memcache"
	"github.com/bilibili/kratos/pkg/cache/redis"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	xtime "github.com/bilibili/kratos/pkg/time"
)

// Dao dao.
type Dao struct {
	db          *sql.DB
	redis       *redis.Pool
	redisExpire int32
	mc          *memcache.Memcache
	mcExpire    int32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a dao and return.
func New() (dao *Dao) {
	var (
		dc struct {
			Demo *sql.Config
		}
		rc struct {
			Demo       *redis.Config
			DemoExpire xtime.Duration
		}
		mc struct {
			Demo       *memcache.Config
			DemoExpire xtime.Duration
		}
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	checkErr(paladin.Get("redis.toml").UnmarshalTOML(&rc))
	checkErr(paladin.Get("memcache.toml").UnmarshalTOML(&mc))
	dao = &Dao{
		// mysql
		db: sql.NewMySQL(dc.Demo),
		// redis
		redis:       redis.NewPool(rc.Demo),
		redisExpire: int32(time.Duration(rc.DemoExpire) / time.Second),
		// memcache
		mc:       memcache.New(mc.Demo),
		mcExpire: int32(time.Duration(mc.DemoExpire) / time.Second),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	if err = d.pingMC(ctx); err != nil {
		return
	}
	if err = d.pingRedis(ctx); err != nil {
		return
	}
	return d.db.Ping(ctx)
}

func (d *Dao) pingMC(ctx context.Context) (err error) {
	if err = d.mc.Set(ctx, &memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *Dao) pingRedis(ctx context.Context) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", "ping", "pong"); err != nil {
		log.Error("conn.Set(PING) error(%v)", err)
	}
	return
}

func (d *Dao) QueryTx(tx *sql.Tx, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	if rows, err := tx.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condStr), condArgs...); err != nil {
		return nil, err
	} else {
		return d.query(rows, table, condStr, condArgs)
	}
}

func (d *Dao) Query(ctx context.Context, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	if rows, err := d.db.Query(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condStr), condArgs...); err != nil {
		return nil, err
	} else {
		return d.query(rows, table, condStr, condArgs)
	}
}

func (d *Dao) query(rows *sql.Rows, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	if ctypes, err := rows.ColumnTypes(); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		var res []map[string]interface{}
		for rows.Next() {
			row := make(map[string]interface{})
			var value []interface{}
			for _, ctype := range ctypes {
				val := reflect.New(ctype.ScanType()).Interface()
				row[ctype.Name()] = val
				value = append(value, val)
			}
			if err := rows.Scan(value...); err != nil {
				continue
			}
			res = append(res, row)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (d *Dao) Insert(ctx context.Context, table string, entry map[string]interface{}) (int64, error) {
	// NOTE: 理论上不存在负数的ID
	return d.Save(ctx, table, "`id`=?", []interface{}{-99}, entry)
}

func (d *Dao) Save(ctx context.Context, table string, condStr string, condArgs []interface{}, entry map[string]interface{}) (int64, error) {
	if tx, err := d.db.Begin(ctx); err != nil {
		return 0, err
	} else if items, err := d.QueryTx(tx, table, condStr, condArgs); err != nil {
		tx.Rollback()
		return 0, err
	} else if len(items) == 0 {
		// 新增
		ks, vs := splitKeyAndVal(entry)
		kstr, vstr := combineInsert(ks)
		if res, err := tx.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, kstr, vstr), vs...); err != nil {
			tx.Rollback()
			return 0, err
		} else {
			tx.Commit()
			return res.RowsAffected()
		}
	} else {
		// 更新
		ks, vs := splitKeyAndVal(entry)
		kstr := combineUpdate(ks)
		if res, err := tx.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, kstr, condStr), append(vs, condArgs)...); err != nil {
			tx.Rollback()
			return 0, err
		} else {
			tx.Commit()
			return res.RowsAffected()
		}
	}
}

func splitKeyAndVal(entry map[string]interface{}) ([]string, []interface{}) {
	var keys []string
	var vals []interface{}
	for k, v := range entry {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	return keys, vals
}

func combineInsert(keys []string) (kstr string, vstr string) {
	for _, key := range keys {
		kstr += key + ","
		vstr += "?,"
	}
	kstr = strings.TrimRight(kstr, ",")
	vstr = strings.TrimRight(vstr, ",")
	return kstr, vstr
}

func combineUpdate(keys []string) string {
	str := ""
	for _, key := range keys {
		str += key + "=?,"
	}
	str = strings.TrimRight(str, ",")
	return str
}
