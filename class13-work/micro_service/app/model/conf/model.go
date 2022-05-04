package conf

import (
	"micro_service/pkg/mysql"
	"micro_service/pkg/redis"
	"time"
)

type AppConfigs struct {
	MysqlDB                 mysql.Conf
	RedisDB                       redis.Conf
	Common                      Common
	RpcServer                   RpcServer
}

type Common struct {
	GrpcAddress         string
	HttpAddress         string
	ServiceName         string
}