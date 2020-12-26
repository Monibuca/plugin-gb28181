package transaction

import (
	"errors"
)

//transaction 的错误定义
var (
	ErrorSyntax  = errors.New("message syntax error")
	ErrorCheck   = errors.New("message check failed")
	ErrorParse   = errors.New("message parse failed")
	ErrorUnknown = errors.New("message unknown")
)
