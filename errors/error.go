package errors

import "github.com/pkg/errors"

var (
	// ErrSystem System error unknown
	ErrSystem = New("System error unknown")
	// ErrToken Token is expired
	ErrToken = New("Token is expired")
	// ErrAccountForbidden account forbidden
	ErrAccountForbidden = New("Account forbidden")
	// ErrDB db error
	ErrDB = New("DB error")
	// ErrDelUsed item is used can not delete
	ErrDelUsed = New("Can not delete used item")
	// ErrParams Params error
	ErrParams = New("Params error")
	// ErrPermission Permission denied
	ErrPermission = New("Permission denied")
	// ErrTypeMismatch  mismatch Type
	ErrTypeMismatch = New("Type mismatch")
	// ErrInvalid Invalid
	ErrInvalid = New("Invalid")
	// ErrRepeat Repeat
	ErrRepeat = New("Repeat")
	// ErrSystemBusy system busy
	ErrSystemBusy = New("System busy")
	// ErrLimitExceed 超出限制
	ErrLimitExceed = New("LimitExceed")
	// ErrNotFount not fount
	ErrNotFount = New("not fount")
)
var errCodeMap = map[error]int{
	ErrToken:            401,
	ErrAccountForbidden: 403,
	ErrSystem:           500,
	ErrDB:               511,
	ErrDelUsed:          512,
	ErrParams:           513,
	ErrPermission:       514,
	ErrTypeMismatch:     515,
	ErrInvalid:          516,
	ErrRepeat:           517,
	ErrSystemBusy:       518,
	ErrLimitExceed:      519,
	ErrNotFount:         520,
}

// Register 注册新的类型错误
func Register(m map[error]int) {
	for err, code := range m {
		errCodeMap[err] = code
	}
}

// ErrCode get error code
func ErrCode(err error) int {
	code := errCodeMap[err]
	if code == 0 {
		code = -1
	}
	return code
}

// ErrCodeMessage get error code and message
func ErrCodeMessage(err error) (string, int) {
	if errCodeMap[err] == 0 {
		err = ErrSystem
	}
	return err.Error(), errCodeMap[err]
}

// New error new
func New(message string) error {
	return errors.New(message)
}

// Wrap error wrap
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf error wrapf
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Cause error cause
func Cause(err error) error {
	return errors.Cause(err)
}

//Errorf error errorf
func Errorf(message string, args ...interface{}) error {
	return errors.Errorf(message, args...)
}

//WithStack error with stack
func WithStack(err error) error {
	return errors.WithStack(err)
}
