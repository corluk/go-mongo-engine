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

		err := mongoEngine.Exec(func(col *mongo.Collection) error {
			_, err := col.InsertMany(context.TODO(), ui, nil)
			return err
		})

		assert.Nil(t, err)
		var itemsB []Item
		err = mongoEngine.Find(bson.M{}, func(cursor *mongo.Cursor) error {
			return cursor.All(context.TODO(), &itemsB)
		}, nil)
		assert.Nil(t, err)
		paginate := Paginate{
			PerPage: 10,
			Page:    1,
			Engine:  &mongoEngine,
		}
		assert.NotNil(t, paginate)
		var items []Item
		opts := options.FindOptions{}
		count, err := mongoEngine.Count(bson.M{}, nil)

		assert.Nil(t, err)
		if err != nil {
			t.Logf("err : %s", err.Error())
		}

		assert.Greater(t, count, int64(1))

		err = paginate.Find(bson.M{}, func(cursor *mongo.Cursor) error {
			return cursor.All(context.TODO(), &items)
		}, &opts)

		assert.Nil(t, err)
		assert.Equal(t, 10, len(items))

	})

}
