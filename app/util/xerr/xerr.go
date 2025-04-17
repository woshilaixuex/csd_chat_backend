package xerr

import (
	"fmt"
	"log"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 自定义错误 background err 是默认错误，如果要返回记得调用push，不然返回的是background err
 * @Date: 2025-02-18 23:02
 */

type AppError struct {
	Code    uint
	Message string
	Err     error
}

func (xerr *AppError) Error() string {
	if xerr.Err != nil {
		return fmt.Sprintf("Code: %v, Message: %s, Inner Error: | %s", xerr.Code, xerr.Message, xerr.Err.Error())
	}
	return fmt.Sprintf("Code: %v, Message: %s", xerr.Code, xerr.Message)
}

// 对比错误
func (xerr *AppError) Equal(apperr *AppError) bool {
	return xerr.Code == apperr.Code
}

// 是否存在该错误
func (xerr *AppError) HaveErr(apperr *AppError) bool {
	if xerr == nil || apperr == nil {
		return false
	}
	codeMatch := func(err1, err2 *AppError) bool {
		return err1.Equal(err2)
	}
	currentErr := xerr

	for currentErr != nil {
		if codeMatch(currentErr, apperr) {
			return true
		}
		nextErr, ok := IsDefined(currentErr.Unwrap())
		if !ok {
			log.Println("unregulated use AppError")
			return false
		}
		currentErr = nextErr
	}

	return codeMatch(currentErr, apperr)
}

// 是否是自己定义的工具err
func IsDefined(err error) (*AppError, bool) {
	xerr, ok := err.(*AppError)
	return xerr, ok
}

// 推出错误栈
func (xerr *AppError) Unwrap() error {
	return xerr.Err
}

// 堆叠错误栈
func (xerr *AppError) Wrap(upe error) *AppError {
	if upe == nil {
		return nil
	}
	xerr, ok := IsDefined(upe)
	if !ok {
		return xerr.otherErr(upe)
	}
	return &AppError{
		Code:    xerr.Code,
		Message: xerr.Message,
		Err:     xerr,
	}
}

func (xerr *AppError) Submit() error {
	if xerr.Err != nil {
		return xerr
	}
	return nil
}

// 其他工具提供的错误
func (xerr *AppError) otherErr(err error) *AppError {
	return &AppError{
		Code:    OtherError.Code,
		Message: err.Error(),
	}
}

// 定义全局错误
func NewAppError(code uint, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     nil,
	}
}
