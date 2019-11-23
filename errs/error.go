package errs

import (
	"fmt"
	"github.com/pkg/errors"
)

type Errno struct {
	Code    int
	Message string
	Hint    string
}

func (err Errno) Error() string {
	return err.Message
}

// Err represents an error
type Err struct {
	Errno
	Err error
}

func New(errno *Errno, err error) *Err {
	return &Err{Errno: *errno, Err: err}
}

func Newf(errno *Errno, format string, args ...interface{}) *Err {
	return &Err{Errno: *errno, Err: errors.Errorf(format, args...)}
}

func (err *Err) Add(message string) *Err {
	err.Message += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) *Err {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) SetHit(format string, args ...interface{}) *Err {
	err.Hint = fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

//Cause return code, message, hint, err
func Cause(err error) (code int, msg string, hint string, errMsg string) {
	if err == nil {
		return OK.Code, OK.Message, OK.Hint, ""
	}

	switch typed := errors.Cause(err).(type) {
	case *Err:
		return typed.Code, typed.Message, typed.Hint, err.Error()
	case *Errno:
		return typed.Code, typed.Message, typed.Hint, err.Error()
	default:
	}

	return InternalServerError.Code, InternalServerError.Message, InternalServerError.Hint, err.Error()
}
