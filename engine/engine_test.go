package engine

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestItem struct {
	ID    primitive.ObjectID `json:"_id"`
	Name  string             `json:"name"`
	Value string             `json:"value"`
	Time  time.Time          `json:"time"`
}

func TestInsertOne(t *testing.T) {

	//var testItem TestItem
	doc := TestItem{
		ID:   primitive.NewObjectID(),
		Name: "test item 12",

		Time: time.Now(),
	}
	mongoEngine := New("mongodb://localhost:27017", "test", "test")

	model := mongo.IndexModel{Keys: bson.D{{"title", "text"}, {"value", "text"}}}

	err := mongoEngine.AddIndex(model, nil)
	if err != nil {

		t.FailNow()
	}

	err = mongoEngine.Save(doc, nil, nil)

	if err != nil {

		t.FailNow()
	}
	/*
		err := mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

			_, err := col.InsertOne(*ctx, doc)

			if err != nil {
				errstr := err.Error()
				t.Log(errstr)
				return err
			}

			return nil

		})

		if err != nil {
			t.FailNow()
		}
		/*
			mongoEngine.Exec(func (col *mongo.Collection) error{
				context , cancel := mongoEngine.Context()
				cursor , err := col.Find(context,bson.M{})


			})
			err := mongoEngine.DropCollection()

			if err != nil {
				t.FailNow()
			}
			doc := TestItem{
				Name: "test item x",
				Time: time.Now(),
			}
			_, err = mongoEngine.InsertOne(&doc, nil)
			if err != nil {
				t.FailNow()
			}
			var doc2 TestItem
			err = mongoEngine.FindOne(&doc2, bson.M{"name": "test item x"}, nil)
			if err != nil {
				t.FailNow()

			}
			if doc2.Name != "test item x" {
				t.FailNow()
			}
			/*mongoEngine.Exec(func(collection *mongo.Collection) error {

				result, err := collection.InsertOne(context.TODO(), doc, nil)

				if err != nil {
					return err
				}
				fmt.Printf("result.InsertedID: %v\n", result.InsertedID)
				return nil
			})*/

}
