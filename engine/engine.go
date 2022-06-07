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

func (mongoEngine *MongoEngine) Exec(onExec func(collection *mongo.Collection) error) error {
	err := mongoEngine.Connect()
	if err != nil {
		return err
	}
	defer mongoEngine.Disconnect()
	return onExec(mongoEngine.Collection)

}

func (mongoEngine *MongoEngine) Find(filter interface{}, opts *options.FindOptions, items []interface{}) error {

	err := mongoEngine.Connect()
	if err != nil {
		return err
	}
	defer mongoEngine.Disconnect()
	cursor, err := mongoEngine.Collection.Find(context.TODO(), filter, opts)

	if err != nil {
		return err
	}

	err = cursor.All(context.TODO(), items)
	return err

}

func (mongoEngine *MongoEngine) InsertOne(doc interface{}, opts *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	err := mongoEngine.Connect()
	if err != nil {
		return nil, err
	}
	defer mongoEngine.Disconnect()

	return mongoEngine.Collection.InsertOne(context.TODO(), doc, opts)
}

func (mongoEngine *MongoEngine) FindOne(item interface{},filter interface{}, opts *options.FindOneOptions) error {
	err := mongoEngine.Connect()
	if err != nil {
		return err
	}
	defer mongoEngine.Disconnect()

	cursor := mongoEngine.Collection.FindOne(context.TODO(), filter, opts)
	err = cursor.Decode(item)
	return err

}
func (mongoEngine *MongoEngine) Exists(filter interface{}, opts *options.FindOneOptions) (bool, error) {

	err := mongoEngine.Connect()
	if err != nil {
		return true, err
	}
	defer mongoEngine.Disconnect()
	cursor := mongoEngine.Collection.FindOne(context.TODO(), filter, opts)
	return cursor.Err().Error() == "ErrNoDocuments", nil

}
func (mongoEngine *MongoEngine) DropCollection() error {
	err := mongoEngine.Connect()
	if err != nil {
		return err
	}
	defer mongoEngine.Disconnect()

	return mongoEngine.Collection.Drop(context.TODO())

}
func (mongoEngine *MongoEngine) Replace(doc interface{}, filter interface{}, opts *options.ReplaceOptions) (*mongo.UpdateResult, error) {
	err := mongoEngine.Connect()
	if err != nil {
		return nil, err
	}
	defer mongoEngine.Disconnect()

	return mongoEngine.Collection.ReplaceOne(context.TODO(), filter, doc, opts)

}

func (mongoEngine *MongoEngine) Save(doc interface{}, filter interface{}, findOpts *options.FindOneOptions, replaceOpts *options.ReplaceOptions, insertOpts *options.InsertOneOptions) error {

	exists, err := mongoEngine.Exists(filter, findOpts)
	if err != nil {
		return err

	}
	if exists {
		_, err = mongoEngine.Replace(doc, filter, replaceOpts)
		if err != nil {
			return err
		}
	} else {
		_, err := mongoEngine.InsertOne(doc, insertOpts)
		if err != nil {
			return err
		}
	}
	return nil
}
