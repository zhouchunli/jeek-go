package logger

import (
	"micro_service/pkg/microerr"
	"testing"

	"github.com/pkg/errors"
)

func TestLogger_ErrorRaw(t *testing.T) {
	origin := errors.New("original error message")
	wrapped := microerr.DbError.Wrap(origin, "wrap message")
	logger := New(Config{
		Glv:        Debug,
		Svc:        "test-service",
		StackDepth: 3,
	})
	logger.ErrorRaw(wrapped, "xxxxxxxx", "mkey")
	logger.Error(wrapped.(*microerr.SysErr).GetCode(), wrapped.Error(), microerr.GetStackTrace(wrapped.(*microerr.SysErr), 10), "yyyyyyy", "mkey")
}
