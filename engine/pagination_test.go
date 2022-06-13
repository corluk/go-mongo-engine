package engine

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Name string `json:"name"`
}

func TestPagination(t *testing.T) {

	t.Run("Should pagination page 1 returns 10 item ", func(t *testing.T) {

		mongoEngine := New("mongodb://localhost:27017", "test", "testpaginaton")
		mongoEngine.DropCollection()
		var itemA []Item
		for i := 0; i < 100; i++ {

			item := Item{
				Name: "name value " + strconv.Itoa(i),
			}

			itemA = append(itemA, item)

		}
		var ui []interface{}
		for _, t := range itemA {
			ui = append(ui, t)
		}
		var itemsB *[]Item
		err := mongoEngine.Exec(func(col *mongo.Collection, ctx *context.Context) error {
			_, err := col.InsertMany(*ctx, ui, nil)
			return err
		})

		assert.Nil(t, err)

		err = mongoEngine.Find(itemsB, bson.M{}, nil)
		assert.NotNil(t, err)
		paginate := Paginate{
			PerPage: 10,
			Page:    1,
			Engine:  &mongoEngine,
		}
		assert.NotNil(t, paginate)
		var items []Item
		opts := options.FindOptions{}
		count, err := mongoEngine.Count(bson.M{}, nil)

		assert.Greater(t, 1, count)
		err = mongoEngine.Find(&items, bson.M{}, &opts)

		assert.Nil(t, err)

		assert.Equal(t, 10, len(items))

	})

}
