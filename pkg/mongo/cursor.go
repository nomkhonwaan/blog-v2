//go:generate mockgen -destination=./mock/cursor_mock.go github.com/nomkhonwaan/myblog/pkg/mongo Cursor

package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// Cursor is a wrapped interface to the original mongo.Cursor for testing benefit
type Cursor interface {
	Close(context.Context) error
	Decode(interface{}) error
	Next(context.Context) bool
}

type cursor struct {
	context.Context
	*mongo.Cursor
}

// Decode extends from the original Cursor.Decode function which accepts slice decoding
func (cur cursor) Decode(val interface{}) error {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		return errors.New("val argument must be a pointer")
	}

	if v.Elem().Kind() != reflect.Slice {
		return cur.Cursor.Decode(val)
	}

	v = v.Elem()
	elemType := reflect.TypeOf(val).Elem().Elem()

	for cur.Cursor.Next(cur.Context) {
		elem := reflect.New(elemType).Interface()
		if err := cur.Cursor.Decode(elem); err != nil {
			return err
		}

		v.Set(reflect.Append(v, reflect.ValueOf(elem).Elem()))
	}

	return nil
}
