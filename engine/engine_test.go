package engine

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func getEngine() MongoEngine {

	mongoEngine := MongoEngine{
		Uri:            "mongodb://localhost:27017",
		DBName:         "test",
		CollectionName: "test",
	}

	return mongoEngine
}

type TestItem struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

func TestInsertOne(t *testing.T) {

	mongoEngine := getEngine()
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
