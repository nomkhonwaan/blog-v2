//go:generate mockgen -destination=./mock/single_result_mock.go github.com/nomkhonwaan/myblog/pkg/mongo SingleResult

package mongo

// SingleResult is a wrapped interface to the original mongo.SingleResult for testing benefit
type SingleResult interface {
	Decode(interface{}) error
}
