package engine

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEngine struct {
	Uri            string
	DBName         string
	CollectionName string
	Connection     *mongo.Client
	Collection     *mongo.Collection
}

func (mongoEngine *MongoEngine) Connect() error {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEngine.Uri))
	if err != nil {
		return err
	}
	mongoEngine.Connection = client
	mongoEngine.Collection = mongoEngine.Connection.Database(mongoEngine.DBName).Collection(mongoEngine.CollectionName)
	return nil
}

func (mongoEngine *MongoEngine) Disconnect() error {

	err := mongoEngine.Connection.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	mongoEngine.Connection = nil
	mongoEngine.Collection = nil
	return nil
}

func (mongoEngine *MongoEngine) SetCollection(db string, name string) (*mongo.Collection, error) {

	mongoEngine.Collection = mongoEngine.Connection.Database(db).Collection(name)
	return mongoEngine.Collection, nil
}
