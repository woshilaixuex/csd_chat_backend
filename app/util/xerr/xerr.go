package xerr

import "fmt"

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 自定义错误
 * @Date: 2025-02-18 23:02
 */

type AppError struct {
	Code    uint
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Code: %d, Message: %s, Inner Error: %s|", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func (e *AppError) Equal(apperr *AppError) bool {
	return e.Code == apperr.Code
}

func IsDefined(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func (e *AppError) Unwrap() error {
	return e.Err
}
func (e *AppError) Wrap(upe *AppError) error {
	if upe.Err != nil {
		return e
	}
	return &AppError{
		Code:    upe.Code,
		Message: upe.Message,
		Err:     e,
	}
}
func NewAppError(code uint, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
