package xerr

import "fmt"

type CodeError struct {
	errCode  int32
	errMsg   string
	httpCode int
}

func (e *CodeError) ErrCode() int32 {
	return e.errCode
}

func (e *CodeError) ErrMsg() string {
	return e.errMsg
}

func (e *CodeError) HTTPCode() int {
	return e.httpCode
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func New(errCode int32, errMsg string, httpCode int) *CodeError {
	return &CodeError{errCode, errMsg, httpCode}
}
