package mongo

// SingleResult is a wrapped interface to the `mongo.SingleResult` for testing purpose
type SingleResult interface {
	// Decode the current document into given variable
	Decode(interface{}) error
}
