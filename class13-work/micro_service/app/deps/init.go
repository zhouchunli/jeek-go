package deps

// 初始化整个项目用到的连接库
func Start() {
	// 加载项目配置
	LoadConfig()
	// redis
	InitRedis()
	// mysql
	InitMysql()
	// logger
	InitLogger()
	//rpc
	//InitRpcClient()
	// dynamo
	//InitDynamoDB()

	//InitEs()
}

func Stop() {
	StopRedis()
	StopMysql()
}
