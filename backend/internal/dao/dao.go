package dao

import (
	"context"
	gsql "database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"os/exec"

	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
	"github.com/bilibili/kratos/pkg/log"

	"backend/internal/utils"
)

// Dao dao.
type Dao struct {
	dc *sql.Config
	db *sql.DB
	tx *sql.Tx
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
	)
	checkErr(paladin.Get("mysql.toml").UnmarshalTOML(&dc))
	dao = &Dao{
		// mysql
		dc: dc.Demo,
		db: sql.NewMySQL(dc.Demo),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	return d.db.Ping(ctx)
}

func (d *Dao) BeginTx(ctx context.Context) (*sql.Tx, error) {
	log.Info("database: Transaction(START)")
	return d.db.Begin(ctx)
}

func (d *Dao) CommitTx(tx *sql.Tx) (err error) {
	if err = tx.Commit(); err != nil {
		d.RollbackTx(tx)
	} else {
		log.Info("database: Transaction(COMMIT)")
	}
	return err
}

func (d *Dao) Commit() (err error) {
	if err = d.tx.Commit(); err != nil {
		d.RollbackTx(d.tx)
	} else {
		log.Info("database: Transaction(COMMIT)")
	}
	return err
}

func (d *Dao) RollbackTx(tx *sql.Tx) error {
	log.Info("database: Transaction(ROLLBACK)")
	return tx.Rollback()
}

func (d *Dao) Rollback() error {
	log.Info("database: Transaction(ROLLBACK)")
	return d.tx.Rollback()
}

var typeMap = map[reflect.Kind]string{
	reflect.String: "VARCHAR(255)",
	reflect.Int32:  "INT",
	reflect.Int64:  "INT(11)",
	reflect.Int:    "INT",
	reflect.Uint8:  "INT",
	reflect.Bool:   "BOOLEAN",
}

var ModelMap = make(map[string]reflect.Type)

func genCreateSQL(table string, typ reflect.Type) (sqls []string, err error) {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (`id` INT(11) AUTO_INCREMENT,", table)
	var pkeys = []string{"id"}
	var fkeys = make(map[string]string)
	var ukeys = []string{}
	var indexes []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// 跳过开头为XXX_的字段
		if len(field.Name) >= 4 && field.Name[0:4] == "XXX_" {
			continue
		}
		// 提取字段名
		fname := utils.CamelToPascal(field.Name)
		if fname == "id" {
			continue
		}
		fattr := "" // 字段属性
		if field.Tag.Get("orm") != "" {
			// 用逗号分隔字段名和字段修饰说明
			strs := strings.Split(field.Tag.Get("orm"), ",")
			switch len(strs) {
			case 1:
				if strs[0] == "*" {
					continue
				}
				fname = strs[0]
			case 2:
				if len(strs[0]) != 0 {
					fname = strs[0]
				}
				strs = strings.Split(strs[1], "|")
				for _, str := range strs {
					switch {
					case str == "PRIMARY_KEY":
						pkeys = append(pkeys, fname)
					case str == "UNIQUE_KEY":
						ukeys = append(ukeys, fname)
						//fattr += "UNIQUE "
					case str == "INDEX":
						indexes = append(indexes, fname)
					case str == "NOT_NULL":
						fattr += "NOT NULL "
					case len(str) > 7 && str[:7] == "DEFAULT":
						pattern := regexp.MustCompile(`\((.|\'|\")+\)`)
						fattr += fmt.Sprintf("DEFAULT %s ", strings.Trim(pattern.FindString(str[7:]), "()"))
					case len(str) > 11 && str[:11] == "FOREIGN_KEY":
						pattern := regexp.MustCompile(`(\w+):`)
						cname := pattern.FindStringSubmatch(str[11:])
						pattern = regexp.MustCompile(`(\w+)\.(\w+)`)
						fkpair := pattern.FindStringSubmatch(str[11:])
						fkey := fname
						if len(cname) >= 2 {
							fkey += ":" + cname[1]
						}
						fkeys[fkey] = fmt.Sprintf("`%s`(`%s`)", fkpair[1], fkpair[2])
					case len(str) > 7 && str[:7] == "COMMENT":
						pattern := regexp.MustCompile(`\(.+\)`)
						fattr += fmt.Sprintf("COMMENT '%s' ", strings.Trim(pattern.FindString(str[7:]), "()"))
					}
				}
			default:
				return nil, errors.New("错误的ORM字段说明")
			}
		}
		// 提取类型
		ftype := field.Type
		fkind := ftype.Kind()
		if fkind == reflect.Ptr {
			ftype = ftype.Elem()
			fkind = ftype.Kind()
		}
		pidField := reflect.StructField{
			Name: "ParentID",
			Type: reflect.TypeOf(int64(0)),
			Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s_id,NOT_NULL|FOREIGN_KEY(%s.id)\"", utils.ToSingular(table), table)),
		}
		stable := fmt.Sprintf("%s_%s_mapper", utils.ToSingular(table), fname)
		if tn, exists := typeMap[fkind]; exists {
			// 存在于类型表的直接拿来用
			sql += fmt.Sprintf("`%s` %s %s,", fname, tn, fattr)
		} else if fkind == reflect.Slice || fkind == reflect.Array {
			// 判断是否是普通类型数组
			eleTyp := ftype.Elem()
			if eleTyp.Kind() == reflect.Ptr {
				eleTyp = eleTyp.Elem()
			}
			if tn, exists = typeMap[eleTyp.Kind()]; exists {
				// 只有一条字段的普通集合，建表与值关联
				if subsqls, err := genCreateSQL(stable, reflect.StructOf([]reflect.StructField{
					pidField,
					{
						Name: "Value",
						Type: eleTyp,
						Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s\"", utils.ToSingular(fname))),
					},
				})); err != nil {
					return nil, err
				} else {
					sqls = append(sqls, subsqls...)
				}
				// 自定义类型数组，递归建表
			} else if subsqls, err := genCreateSQL(stable, reflect.StructOf(append(utils.GetAllFields(eleTyp), pidField))); err != nil {
				return nil, err
			} else {
				sqls = append(sqls, subsqls...)
			}
		} else if fkind == reflect.Map {
			if subsqls, err := genCreateSQL(stable, reflect.StructOf([]reflect.StructField{
				pidField,
				{
					Name: "Key",
					Type: ftype.Key(),
					Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s_key,PRIMARY_KEY\"", fname)),
				},
				{
					Name: "Value",
					Type: ftype.Elem(),
					Tag:  reflect.StructTag(fmt.Sprintf("orm:\"%s_value\"", fname)),
				},
			})); err != nil {
				return nil, err
			} else {
				sqls = append(sqls, subsqls...)
			}
			// 自定义类型，递归建表
		} else if subsqls, err := genCreateSQL(stable, reflect.StructOf(append(utils.GetAllFields(ftype), pidField))); err != nil {
			return nil, err
		} else {
			sqls = append(sqls, subsqls...)
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
	if len(ukeys) != 0 {
		sql += ",UNIQUE("
		for _, uk := range ukeys {
			sql += fmt.Sprintf("`%s`,", uk)
		}
		sql = strings.TrimRight(sql, ",") + ")"
	}
	if len(indexes) != 0 {
		sql += ",INDEX "
		for _, idx := range indexes {
			sql += fmt.Sprintf("`%s`(`%s`),", idx, idx)
		}
		sql = strings.TrimRight(sql, ",")
	}
	if len(fkeys) != 0 {
		for pn, ref := range fkeys {
			// 关联父表的删除操作，无视更新操作
			cname := strings.TrimRight(table, "_mapper")
			splitIdx := strings.Index(pn, ":")
			if splitIdx != -1 {
				cname = pn[splitIdx + 1:]
				pn = pn[0:splitIdx]
			}
			sql += fmt.Sprintf(",CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES %s ON DELETE CASCADE ON UPDATE NO ACTION", cname, pn, ref)
		}
	}
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	sqls = append(sqls, sql)
	return sqls, nil
}

func (d *Dao) Create(ctx context.Context, table string, typ reflect.Type) error {
	if tx, err := d.BeginTx(ctx); err != nil {
		return err
	} else if err := d.CreateTx(tx, table, typ); err != nil {
		return err
	} else if err := d.CommitTx(tx); err != nil {
		return err
	} else {
		return nil
	}
}

func (d *Dao) CreateTx(tx *sql.Tx, table string, typ reflect.Type) error {
	if sqls, err := genCreateSQL(table, typ); err != nil {
		d.RollbackTx(tx)
		return err
	} else {
		for i := len(sqls) - 1; i >= 0; i-- {
			log.Info("database: SQL(%s)", sqls[i])
			if _, err := tx.Exec(sqls[i]); err != nil {
				d.RollbackTx(tx)
				return err
			}
		}
		ModelMap[table] = typ
		return nil
	}
}

func (d *Dao) DropTx(tx *sql.Tx, table string) error {
	if _, err := tx.Exec("DROP TABLE IF EXISTS " + table, []interface{}{}...); err != nil {
		d.RollbackTx(tx)
		return err
	} else {
		return nil
	}
}

func (d *Dao) DropTxByType(tx *sql.Tx, table string, typ reflect.Type) error {
	stable := utils.ToSingular(table)
	for i := 0; i < typ.NumField(); i++ {
		fname := utils.CamelToPascal(typ.Field(i).Name)
		subTable := fmt.Sprintf("%s_%s_mapper", stable, fname)
		if _, err := tx.Exec("DROP TABLE IF EXISTS " + subTable, []interface{}{}...); err != nil {
			d.RollbackTx(tx)
			return err
		}
	}
	return d.DropTx(tx, table)
}

func (d *Dao) ExecTx(tx *sql.Tx, sql string, args []interface{}) (gsql.Result, error) {
	if res, err := tx.Exec(sql, args...); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else {
		return res, nil
	}
}

func (d *Dao) SourceTx(tx *sql.Tx, file string) error {
	pattern := regexp.MustCompile(`(\w+):(\w+)@tcp\(([\w|\.]+):(\d+)\)/([\w|-]+)?`)
	strs := pattern.FindStringSubmatch(d.dc.DSN)
	if len(strs) != 6 {
		return fmt.Errorf("数据库DSN格式有变：%s", d.dc.DSN)
	}

	command := "mysql -P {port} -h {host} -u{username} -p{password} {database} < {source}"
	command = strings.Replace(command, "{username}", strs[1], 1)
	command = strings.Replace(command, "{password}", strs[2], 1)
	command = strings.Replace(command, "{host}", strs[3], 1)
	command = strings.Replace(command, "{port}", strs[4], 1)
	command = strings.Replace(command, "{database}", strs[5], 1)
	command = strings.Replace(command, "{source}", file, 1)
	log.Info("database: Exec SQL file(%s)", command)
	
	cmd := exec.Command("/bin/sh", "-C", command)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行SQL文件（%s）失败：%v", file, err)
	} else {
		log.Info("database: Exec SQL file(%s)", string(out))
	}
	return nil
}

func (d *Dao) QueryOneTx(tx *sql.Tx, table string, condStr string, condArgs []interface{}) (map[string]interface{}, error) {
	if ary, err := d.QueryTx(tx, table, condStr, condArgs); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else if len(ary) == 0 {
		d.RollbackTx(tx)
		return nil, errors.New("Queried data empty")
	} else {
		return ary[0], nil
	}
}

func genQuerySQL(table string, columns []string, cstr string, options []string) string {
	str := "SELECT "
	if columns != nil && len(columns) != 0 {
		str += strings.Join(columns, ", ")
	} else {
		str += "*"
	}
	str += fmt.Sprintf(" FROM `%s`", table)
	if len(cstr) != 0 {
		str += fmt.Sprintf(" WHERE %s", cstr)
	}
	if options != nil && len(options) != 0 {
		str += " " + strings.Join(options, " ")
	}
	return str
}

func (d *Dao) QueryTxBySQL(tx *sql.Tx, table string, sql string, carg []interface{}) ([]map[string]interface{}, error) {
	log.Info("database: SQL(%s); ARGS(%v)", sql, carg)
	if rows, err := tx.Query(sql, carg...); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else if entries, err := d.ProcsResult(rows); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else if ftmapper := chkExistsMapper(table); len(ftmapper) == 0 {
		// 检查是否有外联属性，没有正常返回
		for idx, entry := range entries {
			entries[idx] = FixQueryResult(entry)
		}
		return entries, nil
	} else {
		// 有外联属性，逐个记录分解
		for idx, entry := range entries {
			id := entry["id"]
			// 逐个外联属性处理
			for fname, stable := range ftmapper {
				fkey := fmt.Sprintf("%s_id", utils.ToSingular(table))
				// 通过映射表和外联ID查找所有符合条件的项目
				if prop, err := d.QueryTx(tx, stable, fmt.Sprintf("`%s`=?", fkey), []interface{}{id}); err != nil {
					d.RollbackTx(tx)
					return nil, err
				} else {
					// TODO: 需要判断是否是集合类型的属性
					entries[idx][fname] = prop
				}
			}
		}
		// 修正结果，集合子字段中的id和父id需要剔除；如果集合子字段有字段和子字段名相同，则直接替换
		for idx, entry := range entries {
			entries[idx] = FixQueryResult(entry)
		}
		return entries, nil
	}
}

func (d *Dao) QueryTxOfOption(tx *sql.Tx, table string, cstr string, carg []interface{}, options []string) ([]map[string]interface{}, error) {
	return d.QueryTxBySQL(tx, table, genQuerySQL(table, nil, cstr, options), carg)
}

func (d *Dao) QueryTxIdenCol(tx *sql.Tx, table string, cols []string, cstr string, carg []interface{}) ([]map[string]interface{}, error) {
	return d.QueryTxBySQL(tx, table, genQuerySQL(table, cols, cstr, nil), carg)
}

func (d *Dao) QueryTx(tx *sql.Tx, table string, cstr string, carg []interface{}) ([]map[string]interface{}, error) {
	return d.QueryTxOfOption(tx, table, cstr, carg, nil)
}

func (d *Dao) QueryTxByID(tx *sql.Tx, table string, id interface{}) (map[string]interface{}, error) {
	if ress, err := d.QueryTx(tx, table, "`id`=?", []interface{}{id}); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else if len(ress) == 0 {
		return nil, errors.New("没有找到指定记录")
	} else {
		return ress[0], nil
	}
}

func (d *Dao) Query(ctx context.Context, table string, condStr string, condArgs []interface{}) ([]map[string]interface{}, error) {
	if tx, err := d.BeginTx(ctx); err != nil {
		return nil, err
	} else if entries, err := d.QueryTx(tx, table, condStr, condArgs); err != nil {
		return nil, err
	} else if err := d.CommitTx(tx); err != nil {
		return nil, err
	} else {
		return entries, nil
	}
}

func combineWhereIn(which string, size int) (sql string) {
	sql = fmt.Sprintf("%s IN (", which)
	for i := 0; i < size; i++ {
		sql += "?,"
	}
	sql = strings.TrimRight(sql, ",")
	sql += ")"
	return
}

func chkExistsMapper(table string) (ftmapper map[string]string) {
	ftmapper = make(map[string]string)
	var rowType reflect.Type
	var exists bool
	if rowType, exists = ModelMap[table]; !exists {
		return
	}
	for i := 0; i < rowType.NumField(); i++ {
		field := rowType.Field(i)
		fkind := field.Type.Kind()
		if _, exs := typeMap[fkind]; exs {
			continue
		}
		// 跳过开头为XXX_的字段
		if len(field.Name) >= 4 && field.Name[0:4] == "XXX_" {
			continue
		}
		fname := utils.CamelToPascal(field.Name)
		if fkind == reflect.Slice || fkind == reflect.Array {
			table = utils.ToSingular(table)
		}
		stable := fmt.Sprintf("%s_%s_mapper", table, fname)
		ftmapper[fname] = stable
	}
	return
}

func (d *Dao) ProcsResult(rows *sql.Rows) (res []map[string]interface{}, err error) {
	var ctypes []*gsql.ColumnType
	if ctypes, err = rows.ColumnTypes(); err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		row := make(map[string]interface{})
		var value []interface{}
		for _, ctype := range ctypes {
			row[ctype.Name()] = reflect.New(ctype.ScanType()).Interface()
			value = append(value, row[ctype.Name()])
		}
		if err = rows.Scan(value...); err != nil {
			continue
		}
		// Golang对数据库查询的类型做了包装，需要转一下
		for cname, col := range row {
			tcol := reflect.TypeOf(col)
			switch {
			case tcol.ConvertibleTo(reflect.TypeOf((*int32)(nil))):
				row[cname] = *col.(*int32)
			case tcol.ConvertibleTo(reflect.TypeOf((*gsql.RawBytes)(nil))):
				row[cname] = string(*col.(*gsql.RawBytes))
			case tcol.ConvertibleTo(reflect.TypeOf((*gsql.NullInt64)(nil))):
				row[cname] = (col.(*gsql.NullInt64)).Int64
			case tcol.ConvertibleTo(reflect.TypeOf((*gsql.NullFloat64)(nil))):
				row[cname] = (col.(*gsql.NullFloat64)).Float64
			case tcol.ConvertibleTo(reflect.TypeOf((*gsql.NullString)(nil))):
				row[cname] = (col.(*gsql.NullString)).String
			case tcol.ConvertibleTo(reflect.TypeOf((*gsql.NullBool)(nil))):
				row[cname] = (col.(*gsql.NullBool)).Bool
			}
		}
		res = append(res, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Dao) Insert(ctx context.Context, table string, entry map[string]interface{}) (int64, error) {
	// NOTE: 理论上不存在负数的ID
	return d.Save(ctx, table, "`id`=?", []interface{}{-99}, entry)
}

func (d *Dao) Save(ctx context.Context, table string, condStr string, condArgs []interface{}, entry map[string]interface{}) (int64, error) {
	if tx, err := d.db.Begin(ctx); err != nil {
		return 0, err
	} else {
		return d.SaveTx(tx, table, condStr, condArgs, entry, true)
	}
}

func (d *Dao) InsertTx(tx *sql.Tx, table string, entry map[string]interface{}) (int64, error) {
	return d.SaveTx(tx, table, "`id`=?", []interface{}{-99}, entry, false)
}

func (d *Dao) SaveTx(tx *sql.Tx, table string, condStr string, condArgs []interface{}, entry map[string]interface{}, commit bool) (int64, error) {
	// 提取集合类型和自定义类型
	cmpProps := make(map[string]interface{})
	aryProps := make(map[string][]interface{})
	for pname, prop := range entry {
		tkind := reflect.TypeOf(prop).Kind()
		// 只要是数字，传过来的都是float64……
		if tkind == reflect.Float64 {
			continue
		} else if _, exs := typeMap[tkind]; !exs {
			if tkind == reflect.Slice || tkind == reflect.Array {
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
		d.RollbackTx(tx)
		return 0, err
	} else if len(items) == 0 {
		// 新增
		ks, vs := splitKeyAndVal(entry)
		kstr, vstr := combineInsert(ks)
		sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, kstr, vstr)
		log.Info("database: SQL(%s); ARGS(%v)", sql, vs)
		if res, err = tx.Exec(sql, vs...); err != nil {
			d.RollbackTx(tx)
			return 0, err
		}
	} else {
		// 更新
		ks, vs := splitKeyAndVal(entry)
		kstr := combineUpdate(ks)
		sql := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, kstr, condStr)
		args := append(vs, condArgs...)
		log.Info("database: SQL(%s); ARGS(%v)", sql, args)
		if res, err = tx.Exec(sql, args...); err != nil {
			d.RollbackTx(tx)
			return 0, err
		}
	}
	// 处理复合类型和集合类型
	if id, err := res.LastInsertId(); err != nil {
		d.RollbackTx(tx)
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
			if err := d.CommitTx(tx); err != nil {
				d.RollbackTx(tx)
				return 0, err
			}
		}
		return id, nil
	}
}

func (d *Dao) SaveArrayPropTx(tx *sql.Tx, pname string, prop []interface{}, parent string, pid int64) (int64, error) {
	table := fmt.Sprintf("%s_%s_mapper", utils.ToSingular(parent), utils.CamelToPascal(pname))
	var num int64
	for _, p := range prop {
		if _, exs := typeMap[reflect.TypeOf(p).Kind()]; exs {
			// 普通类型
			if n, err := d.InsertTx(tx, table, map[string]interface{}{
				fmt.Sprintf("%s_id", utils.ToSingular(parent)): pid,
				utils.ToSingular(pname):                        p,
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

func (d *Dao) SaveCompPropTx(tx *sql.Tx, pname string, prop interface{}, parent string, pid int64) (int64, error) {
	table := fmt.Sprintf("%s_%s_mapper", utils.ToSingular(parent), utils.CamelToPascal(pname))
	if !reflect.TypeOf(prop).ConvertibleTo(reflect.TypeOf((*map[string]interface{})(nil)).Elem()) {
		d.RollbackTx(tx)
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
		kstr += fmt.Sprintf("`%s`,", utils.CamelToPascal(key))
		vstr += "?,"
	}
	kstr = strings.TrimRight(kstr, ",")
	vstr = strings.TrimRight(vstr, ",")
	return kstr, vstr
}

func combineUpdate(keys []string) string {
	str := ""
	for _, key := range keys {
		str += fmt.Sprintf("`%s`=?,", utils.CamelToPascal(key))
	}
	str = strings.TrimRight(str, ",")
	return str
}

func (d *Dao) DeleteTx(tx *sql.Tx, table string, cstr string, carg []interface{}) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM `%s`", table)
	if len(cstr) != 0 {
		sql += fmt.Sprintf(" WHERE %s", cstr)
	}
	log.Info("database: SQL(%s); ARGS(%v)", sql, carg)
	if res, err := tx.Exec(sql, carg...); err != nil {
		d.RollbackTx(tx)
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

func (d *Dao) DeleteTxByID(tx *sql.Tx, table string, id interface{}) (int64, error) {
	return d.DeleteTx(tx, table, "`id`=?", []interface{}{id})
}

func (d *Dao) UpdateTxByID(tx *sql.Tx, table string, entry map[string]interface{}) (map[string]interface{}, error) {
	if id, exs := entry["id"]; !exs {
		d.RollbackTx(tx)
		return nil, errors.New("使用ID指定更新需要给出ID")
	} else if delete(entry, "id"); false {
		return nil, nil
	} else if _, err := d.SaveTx(tx, table, "`id`=?", []interface{}{id}, entry, false); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else if res, err := d.QueryTxByID(tx, table, id); err != nil {
		d.RollbackTx(tx)
		return nil, err
	} else {
		return res, nil
	}
}

func FixQueryResult(result map[string]interface{}) map[string]interface{} {
	// 子集合字段中，集合项和子字段名重合（无视单复数）；则直接用该同名集合项字段代替子集合项
	// e.g: "methods":[{"id":13,"method":"PUT","model_id":8},...  => "methods":["PUT",...
	for fname, fvalue := range result {
		// 修改查询结果的字段名字，下划线分隔改为驼峰命名
		// NOTE: 放到前端处理
		// dname := fname
		// fname = utils.PascalToCamel(fname)
		// if dname != fname {
		// 	result[fname] = fvalue
		// 	delete(result, dname)
		// }
		if reflect.TypeOf(fvalue).Kind() != reflect.Slice {
			continue
		}
		farray := make([]interface{}, 0)
		field := fvalue.([]map[string]interface{})
		for idx, element := range field {
			for sfname, sfvalue := range element {
				// NOTE: 结果会删除id和以_id结尾的集合子字段，注意命名！
				if len(sfname) >= 3 && sfname[len(sfname)-3:] == "_id" {
					delete(element, sfname)
				} else if utils.ToPlural(sfname) == fname {
					farray = append(farray, sfvalue)
				}
			}
			field[idx] = element
		}
		if len(farray) != 0 {
			result[fname] = farray
		} else {
			result[fname] = field
		}
	}
	return result
}