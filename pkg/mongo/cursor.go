package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// ScanAll performs scanning all documents and decoding.
// This function requires address of docs variable to slice of struct,
// if non-slice variable has been sent, reject with error "non-slice variable given".
func ScanAll(ctx context.Context, cur *mongo.Cursor, docs interface{}) error {
	v := reflect.ValueOf(docs)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("docs argument must be a slice address")
	}

	v = v.Elem()
	elemType := reflect.TypeOf(docs).Elem().Elem()

	for cur.Next(ctx) {
		elem := reflect.New(elemType).Interface()
		if err := cur.Decode(elem); err != nil {
			return err
		}

		v.Set(reflect.Append(v, reflect.ValueOf(elem).Elem()))
	}

	return nil
}
