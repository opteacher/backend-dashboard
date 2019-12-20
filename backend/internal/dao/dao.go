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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Dao dao.
type MongoDao struct {
	cliOpns *options.ClientOptions
	dbName string
	mod string
}

// New new a dao and return.
func NewMongo() *MongoDao {
	var (
		dc struct {
			Demo *struct{
				Url string
				Db string
				Mod string
				InitData string
			}
		}
	)
	utils.CheckErr(paladin.Get("mongo.toml").UnmarshalTOML(&dc))
	mongoDao := &MongoDao{
		cliOpns: options.Client().ApplyURI(dc.Demo.Url),
		dbName: dc.Demo.Db,
		mod: dc.Demo.Mod,
	}
	if len(dc.Demo.InitData) != 0 {
		var ac struct {
			ProjPath string
		}
		err := paladin.Get("application.toml").UnmarshalTOML(&ac)
		if err != nil {
			panic(err)
		}
		initDataPath := filepath.Join(ac.ProjPath, dc.Demo.InitData)
		jfnames, err := utils.ScanAllFilesByFolder(initDataPath, "json")
		if err != nil {
			panic(err)
		}
		ctx := context.TODO()
		for _, jfname := range jfnames {
			if err = mongoDao.Source(ctx, jfname, true); err != nil {
				panic(err)
			}
		}
	}
	return mongoDao
}

// Close close the resource.
func (md *MongoDao) Close() {
}

// Ping ping the resource.
func (md *MongoDao) Ping(ctx context.Context) error {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	return err
}

func (md *MongoDao) Create(ctx context.Context, colcName string) error {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	if md.mod == "DROP_CREATE" {
		return colc.Drop(ctx)
	} else if md.mod == "UPDATE" {
		if _, err := colc.DeleteMany(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}

func (md *MongoDao) Drop(ctx context.Context, colcName string) error {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	return colc.Drop(ctx)
}

func (md *MongoDao) Source(ctx context.Context, file string, create bool) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}
	fname := filepath.Base(file)
	cname := strings.SplitN(fname, ".", 2)[0]
	cmd := exec.CommandContext(ctx, "mongoimport", []string{
		"--drop",
		"--jsonArray",
		"--db", md.dbName,
		"--collection", cname,
		"--file", file,
	}...)
	if create {
		if err := md.Create(ctx, cname); err != nil {
			return err
		}
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(err.Error() + "\n" + stderr.String())
	}
	return nil
}

func (md *MongoDao) QueryOne(ctx context.Context, colcName string, conds bson.D) (map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	var res bson.M
	err = colc.FindOne(ctx, conds).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (md *MongoDao) Query(ctx context.Context, colcName string, conds bson.D) ([]map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(md.dbName)
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

func (md *MongoDao) Insert(ctx context.Context, colcName string, entry interface{}) (string, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return "", err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	res, err := colc.InsertOne(ctx, entry)
	if err != nil {
		return "", err
	} else {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
}

func (md *MongoDao) InsertMany(ctx context.Context, colcName string, entries []interface{}) (int64, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return -1, err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	res, err := colc.InsertMany(ctx, entries)
	if err != nil {
		return -1, err
	} else {
		return int64(len(res.InsertedIDs)), nil
	}
}

func (md *MongoDao) Save(ctx context.Context, colcName string, conds bson.D, entry interface{}) (map[string]interface{}, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var res bson.M
	err = colc.FindOneAndUpdate(ctx, conds, entry, opts).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (md *MongoDao) Delete(ctx context.Context, colcName string, conds bson.D) (int64, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return -1, err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	res, err := colc.DeleteMany(ctx, conds)
	if err != nil {
		return -1, err
	} else {
		return res.DeletedCount, nil
	}
}

func (md *MongoDao) Update(ctx context.Context, colcName string, conds bson.D, entry interface{}) (int64, error) {
	cli, err := mongo.Connect(ctx, md.cliOpns)
	defer cli.Disconnect(ctx)
	if err != nil {
		return -1, err
	}
	db := cli.Database(md.dbName)
	colc := db.Collection(colcName)
	res, err := colc.UpdateMany(ctx, conds, entry)
	if err != nil {
		return -1, err
	} else {
		return res.UpsertedCount, nil
	}
}