//go:generate mockgen -destination=./mock/single_result_mock.go github.com/nomkhonwaan/myblog/pkg/mongo SingleResult

package mongo

import "go.mongodb.org/mongo-driver/mongo"

// SingleResult is a wrapped interface to the original mongo.SingleResult for testing benefit
type SingleResult interface {
	Decode(interface{}) error
}

// IsErrorRecordNotFound validates an error object is an mongo.ErrNoDocuments or not
func IsErrorRecordNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}
