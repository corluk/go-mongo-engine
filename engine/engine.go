package engine

import (
	"context"
	"fmt"
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

	return mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		_, err := col.Indexes().CreateOne(*ctx, model, opts)
		return err
	})

}
func (mongoEngine *MongoEngine) AddIndexes(model []mongo.IndexModel, opts *options.CreateIndexesOptions) error {

	return mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		_, err := col.Indexes().CreateMany(*ctx, model, opts)
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

func (mongoEngine *MongoEngine) Exec(onExec func(collection *mongo.Collection, context *context.Context) error) error {
	err := mongoEngine.Connect()
	if err != nil {
		return err
	}
	defer mongoEngine.Disconnect()
	ctx, _ := mongoEngine.Context()
	return onExec(mongoEngine.Collection, &ctx)

}
func (mongoEngine *MongoEngine) DropCollection() {

	mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		return col.Drop(*ctx)
	})

}
func (mongoEngine *MongoEngine) FindOne(doc interface{}, filter interface{}, opts *options.FindOneOptions) {

	mongoEngine.Exec(func(collection *mongo.Collection, context *context.Context) error {

		singleResult := collection.FindOne(*context, filter, opts)
		if singleResult.Err() != nil {
			return singleResult.Err()
		}

		singleResult.Decode(&doc)
		return nil
	})

}

/*
func (mongoEngine *MongoEngine) Save(doc interface{}, find *interface{}, opts *options.FindOneAndReplaceOptions) error {

	if opts == nil {
		opts = options.FindOneAndReplace()

	}

	opts.SetUpsert(true)
	return mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		if find != nil {
			result := col.FindOneAndReplace(*ctx, find, doc, opts)
			return result.Err()
		}

		return nil

	})

}
*/
func (mongoEngine *MongoEngine) Count(filter interface{}, opts *options.CountOptions) (int64, error) {
	var size int64
	err := mongoEngine.Exec(func(collection *mongo.Collection, context *context.Context) error {
		_size, err := collection.CountDocuments(*context, filter, opts)
		if err != nil {
			return err
		}
		size = _size
		return nil
	})

	return size, err

}
func (mongoEngine *MongoEngine) Find(docs interface{}, filter interface{}, opts *options.FindOptions) error {

	return mongoEngine.Exec(func(collection *mongo.Collection, ctx *context.Context) error {

		cursor, err := collection.Find(*ctx, filter, opts)

		if err != nil {
			return err
		}

		if cursor.Err() != nil {
			return cursor.Err()
		}

		cursor.All(*ctx, &docs)
		return nil
	})

}

func (mongoEngine *MongoEngine) Save(doc interface{}, filter interface{}) error {

	err := mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		opts := &options.ReplaceOptions{}
		opts.SetUpsert(true)
		_, err := col.ReplaceOne(*ctx, filter, doc, opts)

		return err

	})
	return err

	// })

}
func (mongoEngine *MongoEngine) SearchByText(q string, onCursor func(cursor *mongo.Cursor) error, opts *options.FindOptions) error {

	return mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {

		q := bson.D{primitive.E{Key: "$text", Value: bson.D{primitive.E{Key: "$search", Value: q}}}}

		temporaryBytes, err := bson.MarshalExtJSON(q, true, true)
		if err != nil {

			return err
		}
		str := string(temporaryBytes)
		fmt.Printf(" json %s ", str)
		cursor, err := col.Find(*ctx, q, opts)
		if err != nil {
			return err
		}
		if cursor.Err() != nil {
			return cursor.Err()
		}
		return onCursor(cursor)

	})

}
