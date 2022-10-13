package config

import (
	"fmt"
	"github.com/kanyways/configs"
	"github.com/kanyways/wechat-robot/utils"
	"os"
	"unicode/utf8"
)

// 系统配置
type serverConfig struct {
	ApiPoweredBy string  `yaml:"apiPoweredBy"`
	SiteName     string  `yaml:"siteName"`
	Env          string  `yaml:"env"`
	ApiPrefix    string  `yaml:"apiPrefix"`
	Port         int     `yaml:"port"`
	LogDir       string  `yaml:"logDir"`
	LogFile      string  `yaml:"logFile"`
	ApiKey       string  `yaml:"apiKey"`
	ApiSecret    string  `yaml:"apiSecret"`
	TxApiKey     float64 `yaml:"txApiKey"`
	TxApiSecret  string  `yaml:"txApiSecret"`
}

// 数据库配置
type dBConfig struct {
	Dialect      string `yaml:"dialect"`
	Database     string `yaml:"database"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Charset      string `yaml:"charset"`
	URL          string
	MaxIdleConns int `yaml:"maxIdleConns"`
	MaxOpenConns int `yaml:"maxOpenConns"`
}

// redis的配置
type redisConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Password  string `yaml:"password"`
	Database  int    `yaml:"database"`
	MaxIdle   int    `yaml:"maxIdle"`
	MaxActive int    `yaml:"maxActive"`
	URL       string
}

// 整体配置文件
type Config struct {
	Server   serverConfig `yaml:"server"`
	Database dBConfig     `yaml:"database"`
	Redis    redisConfig  `yaml:"redis"`
}

// 系统整体配置
var conf Config

// DBConfig 数据库相关配置
var DBConfig dBConfig

// RedisConfig redis相关配置
var RedisConfig redisConfig

// ServerConfig 服务器相关配置
var ServerConfig serverConfig

func initConfig() {
	configs.Parse(&conf, configs.GetConfigAbsolutePath("config.yml"))
}

func initDB() {
	DBConfig = conf.Database
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
	DBConfig.URL = url
}

func initRedis() {
	RedisConfig = conf.Redis
	url := fmt.Sprintf("%s:%d", RedisConfig.Host, RedisConfig.Port)
	RedisConfig.URL = url
}

func initServer() {
	ServerConfig = conf.Server
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
	initConfig()
	initDB()
	initRedis()
	initServer()
}
