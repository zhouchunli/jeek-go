/*
 * 参考：https://github.com/henrmota/errors-handling-example/blob/master/errors.go
 */

package microerr

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	//预先定义的统一的错误返回值
	SuccessErr = Success.New("")
)

type ErrType int

type SysErr struct {
	code ErrType
	err  error
}

func (err *SysErr) Error() string {
	if err.err == nil {
		return ""
	}
	return err.err.Error()
}

func (err *SysErr) GetType() ErrType {
	return err.code
}

func (err *SysErr) GetCode() int {
	return int(err.GetType())
}

func (err *SysErr) GetError() error {
	return err.err
}

// 根据错误码创建一个新错误
func (code ErrType) New(msg string) error {
	return &SysErr{code: code, err: errors.New(msg)}
}

// 用格式化字符串创建一个错误
func (code ErrType) Newf(format string, args ...interface{}) error {
	return &SysErr{code: code, err: errors.Errorf(format, args...)}
}

// 包裹一个错误信息，增加上下文，并附加堆栈信息
func (code ErrType) Wrap(err error, msg string) error {
	return code.Wrapf(err, msg)
}

// 包裹一个错误信息，增加格式化字符串的上下文，并附加堆栈信息
func (code ErrType) Wrapf(err error, format string, args ...interface{}) error {
	return &SysErr{code: code, err: errors.Wrapf(err, format, args...)}
}

// 增加一个无类型的错误，添加堆栈信息
func New(msg string) error {
	return &SysErr{code: NoTypeErr, err: errors.New(msg)}
}

// 增加一个无类型的错误，添加堆栈信息
func Newf(format string, args ...interface{}) error {
	return &SysErr{code: NoTypeErr, err: errors.Errorf(format, args...)}
}

// 将一个错误进行包装，添加上下文并增加堆栈信息
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// 将一个错误进行包装，添加格式化上下文并增加堆栈信息
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if sysErr, ok := err.(*SysErr); ok {
		return &SysErr{
			code: sysErr.code,
			err:  wrappedError,
		}
	}

	return &SysErr{code: NoTypeErr, err: wrappedError}
}

// 给一个错误加上上下文消息
func AddContext(err error, context string) error {
	wrappedError := errors.WithMessage(err, context)
	if sysErr, ok := err.(*SysErr); ok {
		return &SysErr{
			code: sysErr.code,
			err:  wrappedError,
		}
	}
	return &SysErr{code: NoTypeErr, err: wrappedError}
}

// 给一个错误加上格式化的上下文消息
func AddContextf(err error, format string, args ...interface{}) error {
	wrappedError := errors.WithMessagef(err, format, args...)
	if sysErr, ok := err.(*SysErr); ok {
		return &SysErr{
			code: sysErr.code,
			err:  wrappedError,
		}
	}
	return &SysErr{code: NoTypeErr, err: wrappedError}
}

// 获取一个错误的错误类型
func GetErrType(err error) ErrType {
	if sysErr, ok := err.(*SysErr); ok {
		return sysErr.code
	}
	return NoTypeErr
}

// 获取一个错误的错误码
func GetErrCode(err error) int {
	return int(GetErrType(err))
}

// 返回最初始的SysErr类型错误
// 判断方法是：err字段没有被errors包进行包装 或者 err字段被包装的不是SysErr类型
func GetCause(err *SysErr) *SysErr {
	if err == nil {
		return nil
	}
	originErr := err
	for {
		causeErr := errors.Cause(originErr.err)
		// 如果err字段没有被包装，则返回
		if originErr.err == causeErr {
			return originErr
		}
		tmp, ok := causeErr.(*SysErr)
		// 如果被包裹的Err不是SysErr类型，则返回
		if !ok {
			return originErr
		}
		originErr = tmp
	}
}

// error如果有堆栈信息，获取堆栈信息
// depth 表示获取从顶部数的多少个栈帧，0表示所有栈帧
// 注意：堆栈信息只应该在创建错误时加入，不应该在添加上下文的时候重复加入
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func GetStackTrace(err *SysErr, depth int) string {
	err = GetCause(err)
	if errWithStack, ok := err.err.(stackTracer); ok {
		stacks := errWithStack.StackTrace()
		stacksDepth := len(stacks)
		if depth == 0 {
			depth = stacksDepth
		} else if depth > stacksDepth-1 {
			depth = stacksDepth
		} else {
			depth += 1
		}
		//因为把errors.New等封装在一个函数中，最顶城堆栈其实不关心，这里其删掉
		return fmt.Sprintf("%+v", stacks[1:depth])
	}
	return ""
}
