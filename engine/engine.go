package engine

import (
	"context"
	"time"

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
		ctx := context.Background()
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
func (mongoEngine *MongoEngine) Find(docs []interface{}, filter interface{}, opts *options.FindOptions) {

	mongoEngine.Exec(func(collection *mongo.Collection, ctx *context.Context) error {

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

/*
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
func (mongoEngine *MongoEngine) InsertMany(docs []interface{}, opts *options.InsertOneOptions) (*mongo.InsertManyResult, error) {
	err := mongoEngine.Connect()
	if err != nil {
		return nil, err
	}
	defer mongoEngine.Disconnect()

	return mongoEngine.Collection.InsertMany(context.TODO(), docs, opts)
}
func (mongoEngine *MongoEngine) InsertOne(doc interface{}, opts *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	err := mongoEngine.Connect()
	if err != nil {
		return nil, err
	}
	defer mongoEngine.Disconnect()

	return mongoEngine.Collection.InsertOne(context.TODO(), doc, opts)
}


/*

func (mongoEngine *MongoEngine) FindOne(item interface{}, filter interface{}, opts *options.FindOneOptions) error {
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



func (mongoEngine *MongoEngine) ReplaceOne(doc interface{}, filter interface{}, opts *options.ReplaceOptions) (*mongo.UpdateResult, error) {
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
		_, err = mongoEngine.ReplaceOne(doc, filter, replaceOpts)
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
*/
