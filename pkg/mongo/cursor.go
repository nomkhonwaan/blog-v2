package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// Cursor is a wrapped interface to the `mongo.Cursor` for testing purpose
type Cursor interface {
	// Close the cursor
	Close(context.Context) error

	// Decode the current document into given variable
	Decode(interface{}) error

	// Get the next result and return true if there were no errors and the next result is available
	Next(context.Context) bool
}

// CustomCursor provides customized cursor methods on top of the original `mongo.Cursor`
type CustomCursor struct {
	context.Context
	*mongo.Cursor
}

// ScanAll performs scanning all documents and decoding to the given docs variable.
func (cur CustomCursor) Decode(val interface{}) error {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		return errors.New("val argument must be a pointer")
	}

	// call the original decoding function for non-slice variable
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
