package deps

import (
	"log"
	"micro_service/app/global"
	"micro_service/pkg/config"
	"os"
	"strings"
)

const (
	ProjectName = "micro_service"
	configPath  = "app" + string(os.PathSeparator) + "config"
	configName  = "config"
)

// 找到配置文件所在目录
func getConfPath() string {
	pwd, _ := os.Getwd()
	index := strings.LastIndex(pwd, ProjectName)
	newPwd := pwd[:index+len(ProjectName)]
	ps := string(os.PathSeparator)
	var configPath = newPwd + ps + configPath + ps + configName + ".yaml"
	log.Println("", configPath)
	return configPath
}

func LoadConfig() {
	config.Load(getConfPath(), &global.AppConfig)
}
