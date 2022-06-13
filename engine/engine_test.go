package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestItem struct {
	Name   string    `json:"name"`
	Value  string    `json:"value,omitempty" bson:"value,omitempty"`
	Time   time.Time `json:"time"`
	Unique string    `json:"unique,omitempty" bson:"unique,omitempty"`
}

func TestInsertOne(t *testing.T) {

	//var testItem TestItem
	doc := TestItem{

		Name: "test item val:",

		Time: time.Now(),
	}
	mongoEngine := New("mongodb://localhost:27017", "test", "test")
	mongoEngine.DropCollection()
	model := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "name", Value: "text"}, primitive.E{Key: "value", Value: "text"}},
		Options: options.Index().SetDefaultLanguage("turkish"),
	}
	opts := options.Index()
	opts.SetUnique(true)
	opts.SetSparse(true)
	models := mongo.IndexModel{

		Keys:    bson.D{primitive.E{Key: "unique1", Value: 1}, primitive.E{Key: "unique2", Value: 1}},
		Options: opts,

		/*, {
			Keys:    bson.D{primitive.E{Key: "unique3", Value: 1}},
			Options: options.Index().SetUnique(true),
		},*/
	}

	err := mongoEngine.AddIndex(models, nil)
	if err != nil {
		t.Logf("err %s", err.Error())
		t.FailNow()
	}
	err = mongoEngine.AddIndex(model, nil)
	if err != nil {
		t.Logf("err %s", err.Error())
		t.FailNow()
	}

	for i := 1; i <= 10; i++ {
		doc.Name = fmt.Sprintf("test value %d", i)

		filter := bson.D{primitive.E{Key: "name", Value: doc.Name}}
		err = mongoEngine.Save(doc, filter)

		if err != nil {
			t.Logf("err %s", err.Error())
			t.FailNow()
		}
	}

	var items []TestItem
	err = mongoEngine.SearchByText("test", func(cursor *mongo.Cursor) error {

		return cursor.All(context.TODO(), &items)
	}, nil)
	if err != nil {
		t.Logf("err %s", err.Error())
		t.FailNow()
	}

	if len(items) < 1 {
		t.Logf("err %s", err.Error())
		t.FailNow()
	}
	filter := bson.D{}
	err = mongoEngine.Find(items, filter, nil)

	if err != nil {
		t.Logf("err %s", err.Error())
		t.FailNow()
	}

}
