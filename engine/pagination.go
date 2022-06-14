package engine

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Paginate struct {
	PerPage int64
	Page    int64
	Engine  *MongoEngine
}

func (paginate *Paginate) setPaginate(opts *options.FindOptions) {

	if opts == nil {
		opts = &options.FindOptions{}
	}
	skip := paginate.PerPage * (paginate.Page - 1)
	limit := paginate.PerPage

	opts.SetLimit(limit)
	opts.SetSkip(skip)

}
func (paginate *Paginate) SearchByText(q string, onCursor func(cursor *mongo.Cursor) error, opts *options.FindOptions) error {

	paginate.setPaginate(opts)

	return paginate.Engine.SearchByText(q, onCursor, opts)
}

func (paginate *Paginate) Find(filter interface{}, onCursor func(cursor *mongo.Cursor) error, opts *options.FindOptions) error {
	paginate.setPaginate(opts)
	return paginate.Engine.Find(filter, onCursor, opts)
}
