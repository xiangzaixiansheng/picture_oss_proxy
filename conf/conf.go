package conf

import (
	"fmt"
	"picture-oss-proxy/cache"
	util "picture-oss-proxy/pkg/utils"

	"os"

	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	RunMode  string
	HttpPort string

	ENV string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func Init() {
	//env
	if _env := os.Getenv("ENV"); _env != "" {
		ENV = _env
	} else {
		ENV = "dev"
	}
	fmt.Println("环境变量是", ENV)
	configFilePath := fmt.Sprintf("./conf/app.%s.ini", ENV)

	file, err := ini.Load(configFilePath)
	if err != nil {
		util.LogrusObj.Errorln("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadRedisData(file)
	if err := LoadLocales("conf/locales/zh-cn.yaml"); err != nil {
		logging.Info(err) //日志内容
		panic(err)
	}
	//redis
	util.LogrusObj.Infoln("[redis]init RedisAddr", RedisAddr, "RedisDbName", RedisDbName)
	cache.NewRedis(RedisAddr, RedisDbName, "")
}

func LoadServer(file *ini.File) {
	RunMode = file.Section("server").Key("RunMode").String()
	HttpPort = file.Section("server").Key("HttpPort").String()
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
