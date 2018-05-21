package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"unicode/utf8"
	"github.com/kanyways/wechat-robot/utils"
)

var jsonData map[string]interface{}

func initJSON() {

	bytes, err := ioutil.ReadFile(utils.GetFilePath("config.json"))
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		fmt.Println("invalid config: ", err.Error())
		os.Exit(-1)
	}
}

type dBConfig struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

// DBConfig 数据库相关配置
var DBConfig dBConfig

func initDB() {
	utils.SetStructByJSON(&DBConfig, jsonData["database"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
	DBConfig.URL = url
}

type redisConfig struct {
	Host      string
	Port      int
	Password  string
	Database  int
	URL       string
	MaxIdle   int
	MaxActive int
}

// RedisConfig redis相关配置
var RedisConfig redisConfig

func initRedis() {
	utils.SetStructByJSON(&RedisConfig, jsonData["redis"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%d", RedisConfig.Host, RedisConfig.Port)
	RedisConfig.URL = url
}

type serverConfig struct {
	APIPoweredBy string
	SiteName     string
	Env          string
	LogDir       string
	LogFile      string
	APIPrefix    string
	Port         int
	ApiKey       string
	ApiSecret    string
}

// ServerConfig 服务器相关配置
var ServerConfig serverConfig

func initServer() {
	utils.SetStructByJSON(&ServerConfig, jsonData["go"].(map[string]interface{}))
	sep := string(os.PathSeparator)

	ymdStr := utils.GetTodayYMD("-")

	if ServerConfig.LogDir == "" {
		ServerConfig.LogDir = utils.GetFilePath(ServerConfig.LogDir + sep + "logs")
	}
	length := utf8.RuneCountInString(ServerConfig.LogDir)
	lastChar := ServerConfig.LogDir[length-1:]
	if lastChar != sep {
		ServerConfig.LogDir = ServerConfig.LogDir + sep
	}
	ServerConfig.LogFile = ServerConfig.LogDir + ymdStr + ".log"
}

func init() {
	initJSON()
	initDB()
	initRedis()
	initServer()
}
