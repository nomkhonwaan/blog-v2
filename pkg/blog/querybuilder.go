package blog

// QueryBuilder is query builder for supporting data store querying syntax
type QueryBuilder interface {
	// Returns offset index
	Offset() int64

	// Allows to override default offset index
	WithOffset(offset int64) QueryBuilder

	// Returns number of the maximum items to-be returned
	Limit() int64

	// Allows to override default maximum items
	WithLimit(limit int64) QueryBuilder
}

type mongoQueryBuilder struct {
	// Initial index of the returned records
	offset int64

	// Maximum items to-be returned from executing
	limit int64
}

// NewQueryBuilder returns new query builder object with default values
func NewQueryBuilder() QueryBuilder {
	return &mongoQueryBuilder{limit: 5}
}

func (qb *mongoQueryBuilder) Offset() int64 {
	return qb.offset
}

func (qb *mongoQueryBuilder) WithOffset(offset int64) QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *mongoQueryBuilder) Limit() int64 {
	return qb.limit
}

func (qb *mongoQueryBuilder) WithLimit(limit int64) QueryBuilder {
	qb.limit = limit
	return qb
}
