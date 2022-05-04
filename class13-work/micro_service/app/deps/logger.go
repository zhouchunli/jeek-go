package deps

import (
	"micro_service/pkg/logger"
)

var (
	Logger *logger.Logger
)

const svcName = "micro_service"

func InitLogger() {
	Logger = logger.New(logger.Config{
		Glv:        logger.Info,
		Svc:        svcName,
		StackDepth: 3,
	})
}
