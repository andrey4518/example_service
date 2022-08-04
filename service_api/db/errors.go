package db

type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

type QueryConditionError struct {
	Message string
}

func (e *QueryConditionError) Error() string {
	return e.Message
}

var intErr *InternalError
var qCondErr *QueryConditionError
