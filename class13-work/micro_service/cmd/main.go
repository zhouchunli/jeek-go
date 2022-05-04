package main

import (
	"micro_service/app/deps"
)

func main() {
	// 加载项目配置
	deps.LoadConfig()
	// 启动依赖
	deps.Start()

}
