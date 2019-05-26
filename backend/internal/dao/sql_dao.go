package dao

import (
	"backend/utils"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"
	gsql "database/sql"
)

type MySqlDao struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlDao() (dao *MySqlDao) {
	var config struct {
		Demo *sql.Config
	}
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&config))
	dao = &MySqlDao{
		db: sql.NewMySQL(config.Demo),
	}
	return
}

func (d *MySqlDao) Close() {
	d.db.Close()
}

func (d *MySqlDao) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.db.Begin(ctx)
}

func (d *MySqlDao) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (d *MySqlDao) Commit() error {
	return d.tx.Commit()
}

func (d *MySqlDao) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (d *MySqlDao) Rollback() error {
	return d.tx.Rollback()
}

var typeMap = map[string]string {
	"string": "VARCHAR(255)",
	"int64":  "INT(11)",
}

func (d *MySqlDao) Create(tx *sql.Tx, table string, typ reflect.Type) (err error) {	
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (`id` INT(11) AUTO_INCREMENT,", table)
	var pkeys = []string{"id"}
	var indexes []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// 提取字段名
		fname := strings.ToLower(field.Name)
		fattr := "" // 字段属性
		if field.Tag.Get("orm") != "" {
			// 用逗号分隔字段名和字段修饰说明
			strs := strings.Split(field.Tag.Get("orm"), ",")
			switch len(strs) {
			case 1:
				fname = strs[0]
			case 2:
				if len(strs[0]) != 0 {
					fname = strs[0]
				}
				strs = strings.Split(strs[1], "|")
				for _, str := range strs {
					switch str {
					case "PRIMARY_KEY":
						pkeys = append(pkeys, fname)
					case "INDEX":
						indexes = append(indexes, fname)
					case "NOT_NULL":
						fattr += "NOT NULL "
					}
				}
			default:
				return errors.New("错误的ORM字段说明")
			}
		}
		// 提取类型
		ftype := field.Type
		tname := fmt.Sprintf("%s", ftype)
		pidField := reflect.StructField{
			Name: "ParentID",
			Type: reflect.TypeOf(int64(0)),
			Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s_id,NOT_NULL|INDEX\"", utils.ToSingular(table))),
		}
		if tn, exists := typeMap[tname]; exists {
			// 存在于类型表的直接拿来用
			sql += fmt.Sprintf("`%s` %s %s,", fname, tn, fattr)
		} else if tname[:2] == "[]" {
			// 判断是否是普通类型数组
			stable := fmt.Sprintf("%s_%s_mapper", table, fname)
			if tn, exists = typeMap[tname[2:]]; exists {
				// 只有一条字段的普通集合，建表与值关联
				if err := d.Create(tx, stable, reflect.StructOf([]reflect.StructField{
					pidField,
					{
						Name: "Value",
						Type: ftype.Elem(),
						Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s\"", utils.ToSingular(fname))),
					},
				})); err != nil {
					return err
				}
				// 自定义类型数组，递归建表
			} else if err := d.Create(tx, stable, reflect.StructOf(append(utils.GetAllFields(ftype.Elem()), pidField))); err != nil {
				return err
			}
			// 自定义类型，递归建表
		} else if err := d.Create(tx, fname, reflect.StructOf(append(utils.GetAllFields(ftype), pidField))); err != nil {
			return err
		}
	}
	sql = strings.TrimRight(sql, ",")
	if len(pkeys) != 0 {
		sql += ",PRIMARY KEY("
		for _, pk := range pkeys {
			sql += fmt.Sprintf("`%s`,", pk)
		}
		sql = strings.TrimRight(sql, ",")
		sql += ")"
	}
	if len(indexes) != 0 {
		sql += ",INDEX "
		for _, idx := range indexes {
			sql += fmt.Sprintf("`%s`(`%s`),", idx, idx)
		}
		sql = strings.TrimRight(sql, ",")
	}
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	log.Info("database: SQL(%s)", sql)
	_, err = tx.Exec(sql)
	return err
}

func (d *MySqlDao) QueryTx(tx *sql.Tx, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condStr)
	log.Info("database: SQL(%s)", sql)
	if rows, err := tx.Query(sql, condArgs...); err != nil {
		return nil, err
	} else {
		return d.query(rows, table, condStr, condArgs)
	}
}

func (d *MySqlDao) Query(ctx context.Context, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, condStr)
	log.Info("database: SQL(%s)", sql)
	if rows, err := d.db.Query(ctx, sql, condArgs...); err != nil {
		return nil, err
	} else {
		return d.query(rows, table, condStr, condArgs)
	}
}

func (d *MySqlDao) query(rows *sql.Rows, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
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

func (d *MySqlDao) Insert(ctx context.Context, table string, entry map[string]interface{}) (int64, error) {
	// NOTE: 理论上不存在负数的ID
	return d.Save(ctx, table, "`id`=?", []interface{}{-99}, entry)
}

func (d *MySqlDao) Save(ctx context.Context, table string, condStr string, condArgs []interface{}, entry map[string]interface{}) (int64, error) {
	if tx, err := d.db.Begin(ctx); err != nil {
		return 0, err
	} else {
		return d.SaveTx(tx, table, condStr, condArgs, entry, true)
	}
}

func (d *MySqlDao) InsertTx(tx *sql.Tx, table string, entry map[string]interface{}) (int64, error) {
	return d.SaveTx(tx, table, "`id`=?", []interface{}{-99}, entry, false)
}

func (d *MySqlDao) SaveTx(tx *sql.Tx, table string, condStr string, condArgs []interface{}, entry map[string]interface{}, commit bool) (int64, error) {
	// 提取集合类型和自定义类型
	cmpProps := make(map[string]interface{})
	aryProps := make(map[string][]interface{})
	for pname, prop := range entry {
		tname := fmt.Sprintf("%s", reflect.TypeOf(prop))
		if _, exs := typeMap[tname]; !exs {
			if tname[0:2] == "[]" {
				aryProps[pname] = prop.([]interface{})
			} else {
				cmpProps[pname] = prop
			}
			delete(entry, pname)
		}
	}
	// 检查要插入的对象是否存在
	var res gsql.Result
	if items, err := d.QueryTx(tx, table, condStr, condArgs); err != nil {
		tx.Rollback()
		return 0, err
	} else if len(items) == 0 {
		// 新增
		ks, vs := splitKeyAndVal(entry)
		kstr, vstr := combineInsert(ks)
		sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, kstr, vstr)
		log.Info("database: SQL(%s)", sql)
		if res, err = tx.Exec(sql, vs...); err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		// 更新
		ks, vs := splitKeyAndVal(entry)
		kstr := combineUpdate(ks)
		sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, kstr, condStr)
		log.Info("database: SQL(%s)", sql)
		if res, err = tx.Exec(sql, append(vs, condArgs)...); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// 处理复合类型和集合类型
	if id, err := res.LastInsertId(); err != nil {
		tx.Rollback()
		return 0, err
	} else {
		for pname, prop := range cmpProps {
			if _, err := d.SaveCompPropTx(tx, pname, prop, table, id); err != nil {
				return 0, err
			}
		}
		for pname, prop := range aryProps {
			if _, err := d.SaveArrayPropTx(tx, pname, prop, table, id); err != nil {
				return 0, err
			}
		}
		if commit {
			tx.Commit()
		}
	}
	return res.RowsAffected()
}

func (d *MySqlDao) SaveArrayPropTx(tx *sql.Tx, pname string, prop []interface{}, parent string, pid int64) (int64, error) {
	table := fmt.Sprintf("%s_%s_mapper", strings.ToLower(parent), strings.ToLower(pname))
	var num int64
	for _, p := range prop {
		if _, exs := typeMap[reflect.TypeOf(p).Name()]; exs {
			// 普通类型
			if n, err := d.InsertTx(tx, table, map[string]interface{}{
				fmt.Sprintf("%s_id", utils.ToSingular(parent)): pid,
				pname: p,
			}); err != nil {
				return 0, err
			} else {
				num += n
			}
		} else if fmt.Sprintf("%s", p)[0:2] == "[]" {
			// 集合中的集合类型
			// TODO: 应该不常用
		} else {
			// 自定义类型
			if n, err := d.SaveCompPropTx(tx, pname, p, parent, pid); err != nil {
				return 0, err
			} else {
				num += n
			}
		}
	}
	return num, nil
}

func (d *MySqlDao) SaveCompPropTx(tx *sql.Tx, pname string, prop interface{}, parent string, pid int64) (int64, error) {
	table := fmt.Sprintf("%s_%s_mapper", strings.ToLower(parent), strings.ToLower(pname))
	if !reflect.TypeOf(prop).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
		tx.Rollback()
		return 0, fmt.Errorf("非对象键值对类型：%s", reflect.TypeOf(prop).Name())
	}
	entry := prop.(map[string]interface{})
	entry[fmt.Sprintf("%s_id", utils.ToSingular(parent))] = pid
	return d.InsertTx(tx, table, entry)
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
		kstr += fmt.Sprintf("`%s`,", key)
		vstr += "?,"
	}
	kstr = strings.TrimRight(kstr, ",")
	vstr = strings.TrimRight(vstr, ",")
	return kstr, vstr
}

func combineUpdate(keys []string) string {
	str := ""
	for _, key := range keys {
		str += fmt.Sprintf("`%s`=?,", key)
	}
	str = strings.TrimRight(str, ",")
	return str
}
