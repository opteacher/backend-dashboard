package test

import (
	"backend/internal/dao"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestPing(t *testing.T) {
	d := dao.New()
	assert.NoError(t, d.Ping(context.Background()))
}

func TestCreate(t *testing.T) {
	d := dao.New()
	assert.NoError(t, d.Create(context.Background(), "abcd"))
}

func TestDrop(t *testing.T) {
	d := dao.New()
	assert.NoError(t, d.Drop(context.Background(), "abcd"))
}

func TestInsert(t *testing.T) {
	d := dao.New()
	res, err := d.Insert(context.Background(), "abcd", map[string]interface{}{
		"name": "cctv",
		"age": 15,
	})
	assert.NoError(t, err)
	println(res)
}

func TestQueryOne(t *testing.T) {
	d := dao.New()
	rowid, err := primitive.ObjectIDFromHex("5dd4f845c4f5b19c28963c16")
	assert.NoError(t, err)
	res, err := d.QueryOne(context.Background(), "abcd", bson.D{{
		"_id", rowid,
	}})
	sres, err := json.Marshal(res)
	assert.NoError(t, err)
	println(string(sres))
	assert.NoError(t, err)
}

func TestQuery(t *testing.T) {
	d := dao.New()
	ress, err := d.Query(context.Background(), "abcd", bson.D{{
		"name", "opower",
	}})
	assert.NoError(t, err)
	for _, res := range ress {
		bytes, err := json.Marshal(res)
		assert.NoError(t, err)
		println(string(bytes))
	}
}

func TestDelete(t *testing.T) {
	d := dao.New()
	res, err := d.Delete(context.Background(), "abcd", bson.D{{
		"name", "opower",
	}})
	assert.NoError(t, err)
	println(res)
}

func TestUpdate(t *testing.T) {
	d := dao.New()
	res, err := d.Update(context.Background(), "abcd", bson.D{{
		"name", "opower",
	}}, map[string]interface{}{
		"age": 27,
	})
	assert.NoError(t, err)
	print(res)
}

func TestSave(t *testing.T) {
	d := dao.New()
	rowid, err := primitive.ObjectIDFromHex("5dd4f845c4f5b19c28963c16")
	assert.NoError(t, err)
	res, err := d.Save(context.Background(), "abcd", bson.D{{
		"_id", rowid,
	}}, bson.D{{
		"$set", bson.D{{"age", 5}},
	}})
	assert.NoError(t, err)
	bytes, err := json.Marshal(res)
	assert.NoError(t, err)
	println(string(bytes))
}

func TestSource(t *testing.T) {
	d := dao.New()
	assert.NoError(t, d.Source(context.Background(), "/Users/zhaojiachen/Projects/backend-dashboard/backend/datas/abcd.json"))
}