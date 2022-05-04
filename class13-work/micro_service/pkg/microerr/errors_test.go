package microerr

import (
	goerrors "errors"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testErr struct {
	msg string
}

func (c testErr) Error() string {
	return c.msg
}

func TestSysErrNew(t *testing.T) {
	errMsg := "[1016] Db access microerr"
	e := DbError.New(errMsg)

	got := e.Error()
	assert.Equal(t, errMsg, got)
}

func TestSysErrNewf(t *testing.T) {
	errMsg := "[1016] Db access microerr"
	e := DbError.Newf("[%d] %s", 1016, "Db access microerr")

	got := e.Error()
	assert.Equal(t, errMsg, got)
}

func TestSysErrWrap(t *testing.T) {
	msg := "reading user from db"
	origin := errors.New("[1016] Db access microerr")
	wrapErr := DbError.Wrap(origin, msg)

	got := wrapErr.Error()
	want := fmt.Sprintf("%s: %s", msg, origin.Error())
	assert.Equal(t, want, got)
}

func TestSysErrWrapf(t *testing.T) {
	msg := "reading user from db"
	origin := errors.New("[1016] Db access microerr")
	wrapErr := DbError.Wrapf(origin, "wrapped: %s", msg)

	got := wrapErr.Error()
	want := fmt.Sprintf("wrapped: %s: %s", msg, origin.Error())
	assert.Equal(t, want, got)
}

func TestNew(t *testing.T) {
	errMsg := "[1016] Db access microerr"
	e := New(errMsg)

	got := e.Error()
	assert.Equal(t, errMsg, got)
}

func TestNewf(t *testing.T) {
	errMsg := "[1016] Db access microerr"
	e := Newf("[%d] %s", 1016, "Db access microerr")

	got := e.Error()
	assert.Equal(t, errMsg, got)
}

func TestAddContext(t *testing.T) {
	msg1 := "wrapped message 1"
	msg2 := "wrapped message 2"
	msg3 := "wrapped message 3"
	msgOrigin := "Mysql has gone away"
	err := DbError.New(msgOrigin)

	err = AddContext(err, msg1)
	err = AddContext(err, msg2)
	err = AddContext(err, msg3)

	got := err.Error()
	want := fmt.Sprintf("%s: %s: %s: %s", msg3, msg2, msg1, msgOrigin)
	assert.Equal(t, want, got)
}

func TestAddContextf(t *testing.T) {
	msg1 := "wrapped message %d"
	msg2 := "wrapped message %d"
	msg3 := "wrapped message %d"
	msgOrigin := "Mysql has gone away"
	err := DbError.New(msgOrigin)

	err = AddContextf(err, msg1, 1)
	err = AddContextf(err, msg2, 2)
	err = AddContextf(err, msg3, 3)

	got := err.Error()
	want := fmt.Sprintf("wrapped message %d: wrapped message %d: wrapped message %d: %s", 3, 2, 1, msgOrigin)
	assert.Equal(t, want, got)
}

func TestGetErrType(t *testing.T) {
	err := DbError.New("Origin err")
	err = Wrap(err, "wrapped message")
	err = AddContext(err, "context message")

	got := GetErrType(err)
	assert.Equal(t, DbError, got)
}

func TestGetStackTrace(t *testing.T) {
	// 无包裹的错误
	err1 := func() error {
		err := DbError.New("microerr message")
		return err
	}()
	trace11 := GetStackTrace(err1.(*SysErr), 0)
	trace12 := GetStackTrace(err1.(*SysErr), 100)
	assert.NotEmpty(t, trace11)
	assert.NotEmpty(t, trace12)
	assert.Equal(t, trace11, trace12)
	trace13 := GetStackTrace(err1.(*SysErr), 1)
	assert.NotEmpty(t, trace13)

	// 有包裹的错误1
	err2 := func() error {
		err := DbError.New("microerr message")
		err = Wrap(err, "message 1")
		err = Wrap(err, "message 2")
		return err
	}()
	trace21 := GetStackTrace(err2.(*SysErr), 0)
	trace22 := GetStackTrace(err2.(*SysErr), 100)
	assert.NotEmpty(t, trace21)
	assert.NotEmpty(t, trace22)
	assert.Equal(t, trace21, trace22)
	trace23 := GetStackTrace(err2.(*SysErr), 1)
	assert.NotEmpty(t, trace23)

	// 有包裹的错误2
	err3 := func() error {
		err := DbError.New("microerr message")
		err = AddContext(err, "message 1")
		err = AddContext(err, "message 2")
		return err
	}()
	trace31 := GetStackTrace(err3.(*SysErr), 0)
	trace32 := GetStackTrace(err3.(*SysErr), 100)
	assert.NotEmpty(t, trace31)
	assert.NotEmpty(t, trace32)
	assert.Equal(t, trace31, trace32)
	trace33 := GetStackTrace(err3.(*SysErr), 1)
	assert.NotEmpty(t, trace33)

	// 自定义错误1
	err4 := func() error {
		err := AddContext(testErr{msg: "custom microerr message"}, "message 1")
		err = AddContext(err, "message 2")
		return err
	}()
	trace41 := GetStackTrace(err4.(*SysErr), 0)
	trace42 := GetStackTrace(err4.(*SysErr), 100)
	trace43 := GetStackTrace(err4.(*SysErr), 1)
	assert.Empty(t, trace41)
	assert.Empty(t, trace42)
	assert.Empty(t, trace43)

	// 自定义错误2
	err5 := func() error {
		err := Wrap(testErr{msg: "custom microerr message"}, "message 1")
		err = Wrap(err, "message 2")
		err = Wrap(err, "message 3")
		return err
	}()
	trace51 := GetStackTrace(err5.(*SysErr), 0)
	trace52 := GetStackTrace(err5.(*SysErr), 100)
	assert.NotEmpty(t, trace51)
	assert.NotEmpty(t, trace52)
	assert.Equal(t, trace51, trace52)
	trace53 := GetStackTrace(err5.(*SysErr), 1)
	assert.NotEmpty(t, trace53)

}

func TestGetSysErrCause(t *testing.T) {
	msg1 := "wrapped message 1"
	msg2 := "wrapped message 2"
	msg3 := "wrapped message 3"
	originMsg := "Mysql has gone away"
	origin := DbError.New(originMsg)

	// 测试为包装的SysErr
	causeErr0 := GetCause(origin.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr0)

	// 测试ErrType.New + AddContext
	err1 := AddContext(origin, msg1)
	err1 = AddContext(err1, msg2)
	err1 = AddContext(err1, msg3)
	causeErr1 := GetCause(err1.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr1)

	// 测试ErrType.New + Wrap
	err2 := Wrap(origin, msg1)
	err2 = Wrap(origin, msg2)
	err2 = Wrap(origin, msg3)
	causeErr2 := GetCause(err2.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr2)

	// 原生error类型 + Wrap
	err3 := goerrors.New("Raw Error Message")
	origin = Wrap(err3, msg1)
	err3 = Wrap(origin, msg2)
	causeErr3 := GetCause(err3.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr3)

	// 原生error类型 + AddContext
	err4 := goerrors.New("Raw Error Message")
	origin = AddContext(err4, msg1)
	err4 = AddContext(origin, msg2)
	causeErr4 := GetCause(err4.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr4)

	// 自定义错误类型 + wrap
	err5 := testErr{msg: "Custom Error Message"}
	origin = Wrap(err5, msg1)
	err55 := Wrap(origin, msg2)
	err55 = Wrap(err55, msg3)
	causeErr5 := GetCause(err55.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr5)

	// 自定义错误类型 + wrap
	err6 := testErr{msg: "Custom Error Message"}
	origin = AddContext(err6, msg1)
	err66 := Wrap(origin, msg2)
	err66 = Wrap(err66, msg3)
	causeErr6 := GetCause(err66.(*SysErr))
	assert.Equal(t, origin.(*SysErr), causeErr6)

}
