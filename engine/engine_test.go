package engine

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestItem struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `json:"name"`
	Value string             `json:"value"`
	Time  time.Time          `json:"time"`
}

func TestInsertOne(t *testing.T) {

	//var testItem TestItem
	doc := TestItem{
		ID:   primitive.NewObjectID(),
		Name: "test item 122s2",

		Time: time.Now(),
	}
	mongoEngine := New("mongodb://localhost:27017", "test", "test")
	mongoEngine.DropCollection()
	model := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "name", Value: "text"}, primitive.E{Key: "value", Value: "text"}},
		Options: options.Index().SetDefaultLanguage("turkish"),
	}
	models := mongo.IndexModel{

		Keys:    bson.D{primitive.E{Key: "unique1", Value: 1}, primitive.E{Key: "unique2", Value: 1}},
		Options: options.Index().SetUnique(true),

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

		t.FailNow()
	}
	filter := bson.D{primitive.E{Key: "name", Value: doc.Name}}
	err = mongoEngine.Save(doc, filter)

	if err != nil {

		t.FailNow()
	}
	var items []TestItem
	err = mongoEngine.SearchByText("test", func(cursor *mongo.Cursor) error {

		return cursor.All(context.TODO(), &items)
	}, nil)
	if err != nil {

		t.FailNow()
	}

	if len(items) < 1 {
		t.FailNow()
	}

}
