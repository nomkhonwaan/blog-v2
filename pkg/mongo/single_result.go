// go:generate mockgen -destination=./mock/single_result_mock.go github.com/nomkhonwaan/myblog/mongo SingleResult

package mongo

import "go.mongodb.org/mongo-driver/mongo"

// SingleResult is a wrapped interface to the `mongo.SingleResult` for testing purpose
type SingleResult interface {
	// Decode the current document into given variable
	Decode(interface{}) error
}

// IsErrorRecordNotFound returns boolean based on given error object is equal to no documents error or not
func IsErrorRecordNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}
