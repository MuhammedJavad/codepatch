package rabbitmq

import "fmt"

type rejectError struct{
	err error
}

func (err rejectError) Error() string {
	return fmt.Sprintf("reject error: %v", err.err)
}

func (err rejectError) Unwrap() error {
	return err.err
}

func RejectError(err error) error {
	return rejectError{err: err}
}

func isRejectError(err error) bool {
	_, ok := err.(rejectError)
	return ok
}
