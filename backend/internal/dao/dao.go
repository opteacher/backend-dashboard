package dao

import (
	"backend/internal/utils"
	"bytes"
	"context"
	"errors"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os/exec"
	"path/filepath"
	"strings"
)

// Dao dao.
type Dao struct {
	cliOpns *options.ClientOptions
	dbName string
	mod string
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
			Demo *struct{
				Url string
				Db string
				Mod string
			}
		}
	)
	checkErr(paladin.Get("mongo.toml").UnmarshalTOML(&dc))
	dao = &Dao{
		cliOpns: options.Client().ApplyURI(dc.Demo.Url),
		dbName: dc.Demo.Db,
		mod: dc.Demo.Mod,
	}
	return dao
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) error {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	return err
}

func (d *Dao) Create(ctx context.Context, colcName string) error {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	if d.mod == "DROP_CREATE" {
		return colc.Drop(ctx)
	} else if d.mod == "UPDATE" {
		if _, err := colc.DeleteMany(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dao) Drop(ctx context.Context, colcName string) error {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	return colc.Drop(ctx)
}

func (d *Dao) Source(ctx context.Context, file string) error {
	fname := filepath.Base(file)
	cname := strings.SplitN(fname, ".", 2)[0]
	cmd := exec.CommandContext(ctx, "mongoimport", []string{
		"--drop",
		"--jsonArray",
		"--db", d.dbName,
		"--collection", cname,
		"--file", file,
	}...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(err.Error() + "\n" + stderr.String())
	}
	return nil
}

func (d *Dao) QueryOne(ctx context.Context, colcName string, conds bson.D) (map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	var res bson.M
	err = colc.FindOne(ctx, conds).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Dao) Query(ctx context.Context, colcName string, conds bson.D) ([]map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	cursor, err := colc.Find(ctx, conds)
	if err != nil {
		return nil, err
	}
	var ress []bson.M
	if err := cursor.All(ctx, &ress); err != nil {
		return nil, err
	}
	var mress []map[string]interface{}
	for _, res := range ress {
		mres, err := utils.ToMap(res)
		if err != nil {
			return nil, err
		}
		mress = append(mress, mres)
	}
	return mress, nil
}

func (d *Dao) Insert(ctx context.Context, colcName string, entry interface{}) (string, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return "", err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	res, err := colc.InsertOne(ctx, entry)
	if err != nil {
		return "", err
	} else {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
}

func (d *Dao) Save(ctx context.Context, colcName string, conds bson.D, entry interface{}) (map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var res bson.M
	err = colc.FindOneAndUpdate(ctx, conds, entry, opts).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Dao) Delete(ctx context.Context, colcName string, conds bson.D) (int64, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return -1, err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	res, err := colc.DeleteMany(ctx, conds)
	if err != nil {
		return -1, err
	} else {
		return res.DeletedCount, nil
	}
}

func (d *Dao) Update(ctx context.Context, colcName string, conds bson.D, entry interface{}) (int64, error) {
	cli, err := mongo.Connect(ctx, d.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return -1, err
	}
	db := cli.Database(d.dbName)
	colc := db.Collection(colcName)
	res, err := colc.UpdateMany(ctx, conds, entry)
	if err != nil {
		return -1, err
	} else {
		return res.UpsertedCount, nil
	}
}