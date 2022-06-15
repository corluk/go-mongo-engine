package engine

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoEngine struct {
	Uri            string
	DBName         string
	CollectionName string
	Connection     *mongo.Client
	Collection     *mongo.Collection
	Context        func() (context.Context, context.CancelFunc)
}
type Context interface {
	GetContext() context.Context
}

func New(uri string, dbName string, colName string) MongoEngine {

	engine := MongoEngine{
		Uri:            uri,
		DBName:         dbName,
		CollectionName: colName,
	}
	engine.Context = func() (context.Context, context.CancelFunc) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		return ctx, cancel
	}

	return engine

}
func (mongoEngine *MongoEngine) AddIndex(model mongo.IndexModel, opts *options.CreateIndexesOptions) error {

	return mongoEngine.Exec(func(col *mongo.Collection) error {

		_, err := col.Indexes().CreateOne(context.TODO(), model, opts)
		return err
	})

}
func (mongoEngine *MongoEngine) AddIndexes(model []mongo.IndexModel, opts *options.CreateIndexesOptions) error {

	return mongoEngine.Exec(func(col *mongo.Collection) error {

		_, err := col.Indexes().CreateMany(context.TODO(), model, opts)
		return err
	})

}

func (mongoEngine *MongoEngine) Connect() error {

	ctx, _ := mongoEngine.Context()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEngine.Uri))
	if err != nil {
		return err
	}
	mongoEngine.Connection = client
	mongoEngine.Collection = mongoEngine.Connection.Database(mongoEngine.DBName).Collection(mongoEngine.CollectionName)
	return nil
}

func (mongoEngine *MongoEngine) Disconnect() error {
	ctx, _ := mongoEngine.Context()
	err := mongoEngine.Connection.Disconnect(ctx)
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
func (mongoEngine *MongoEngine) DropCollection() {

	mongoEngine.Exec(func(col *mongo.Collection) error {

		return col.Drop(context.TODO())
	})

}
func (mongoEngine *MongoEngine) FindOne(filter interface{}, onCursor func(cursor *mongo.SingleResult, ctx *context.Context) error, opts *options.FindOneOptions) {

	mongoEngine.Exec(func(collection *mongo.Collection) error {

		cursor := collection.FindOne(context.TODO(), filter, opts)
		if cursor.Err() != nil {
			return cursor.Err()
		}
		ctx, _ := mongoEngine.Context()
		return onCursor(cursor, &ctx)

	})

}

func (mongoEngine *MongoEngine) Count(filter interface{}, opts *options.CountOptions) (int64, error) {
	var size int64
	err := mongoEngine.Exec(func(collection *mongo.Collection) error {
		_size, err := collection.CountDocuments(context.TODO(), filter, opts)
		if err != nil {
			return err
		}
		size = _size
		return nil
	})

	return size, err

}
func (mongoEngine *MongoEngine) Find(filter interface{}, onCursor func(cursor *mongo.Cursor) error, opts *options.FindOptions) error {

	return mongoEngine.Exec(func(collection *mongo.Collection) error {

		cursor, err := collection.Find(context.TODO(), filter, opts)

		if err != nil {
			return err
		}

		return onCursor(cursor)

	})

}

func (mongoEngine *MongoEngine) Save(doc interface{}, filter interface{}) error {

	return mongoEngine.Exec(func(col *mongo.Collection) error {

		opts := &options.ReplaceOptions{}
		opts.SetUpsert(true)
		_, err := col.ReplaceOne(context.TODO(), filter, doc, opts)

		return err

	})

}

func (mongoEngine *MongoEngine) SearchByText(q string, onCursor func(cursor *mongo.Cursor) error, opts *options.FindOptions) error {

	return mongoEngine.Exec(func(col *mongo.Collection) error {

		q := bson.D{primitive.E{Key: "$text", Value: bson.D{primitive.E{Key: "$search", Value: q}}}}
		findCtx, _ := mongoEngine.Context()
		cursor, err := col.Find(findCtx, q, opts)

		if err != nil {
			return err
		}

		return onCursor(cursor)

	})

}
