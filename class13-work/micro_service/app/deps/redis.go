package deps

import (
	"context"
	"micro_service/app/constant"
	"micro_service/app/global"
	"micro_service/app/library/utils"
	"micro_service/pkg/microerr"
	lib "micro_service/pkg/redis"

	"github.com/go-redis/redis/v8"
)

var (
	RedisDB *redis.Client
)

func InitRedis() {
	ctx := context.Background()
	RedisDB = lib.New(ctx, &global.AppConfig.Redis)
}

func StopRedis() {
	var err error
	ctx := utils.CreateTracingToContext(context.Background(), global.AppConfig.Common.ServiceName)
	err = RedisDB.Close()
	if err != nil {
		sysErr := microerr.ServerCommonError.Wrap(err, "RedisDB Close fali")
		Logger.ErrorRaw(sysErr, utils.GetContextTraceId(ctx), "NA")
	}
}
