package deps

import (
	"context"
	"micro_service/app/constant"
	"micro_service/app/global"
	"micro_service/app/library/utils"
	"micro_service/pkg/microerr"
	"micro_service/pkg/mysql"
	"micro_service/pkg/mysql/sqlt"
)

var (
	MysqlDB *sqlt.DB
)

func InitMysql() {
	MysqlDB = mysql.New(&global.AppConfig.Mysql)
}

func StopMysql() {
	var err error
	ctx := utils.CreateTracingToContext(context.Background(), global.AppConfig.Common.ServiceName)

	err = MysqlDB.Close()
	if err != nil {
		sysErr := microerr.ServerCommonError.Wrap(err, "DB mysql Close fali")
		Logger.ErrorRaw(sysErr, utils.GetContextTraceId(ctx), "NA")
	}

}
