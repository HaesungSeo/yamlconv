package yamlconv

import "errors"

var (
	ErrNotFoundError         = errors.New("not found")
	ErrInvalidIndexError     = errors.New("invalid index")
	ErrIndexOutOfRangeError  = errors.New("index out of range")
	ErrSearchKeyTooLongError = errors.New("too many keys")
)

type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

func (e *NotFoundError) Unwrap() error { return e.Err }

type InvalidIndexError struct {
	Err error
}

func (e *InvalidIndexError) Error() string {
	return e.Err.Error()
}

func (e *InvalidIndexError) Unwrap() error { return e.Err }

type IndexOutOfRangeError struct {
	Err error
}

func (e *IndexOutOfRangeError) Error() string {
	return e.Err.Error()
}

func (e *IndexOutOfRangeError) Unwrap() error { return e.Err }

type SearchKeyTooLongError struct {
	Err error
}

func (e *SearchKeyTooLongError) Error() string {
	return e.Err.Error()
}

func (e *SearchKeyTooLongError) Unwrap() error { return e.Err }
