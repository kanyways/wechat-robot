package receive

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/utils"
	"log"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/message"
	"strconv"
	"github.com/kanyways/wechat-robot/model"
	"net/http"
	"github.com/silenceper/wechat/cache"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func GetMessage(ctx *gin.Context) {

	// 获取当前的用户的微信的配置
	wechatId, err := strconv.Atoi(ctx.Param("id"))

	var wx model.Wechat

	if err := model.DB.First(&wx, wechatId).Error; err != nil {
		fmt.Printf(err.Error())
		return
	}

	if err != nil {
		log.Println("错误的微信用户ID")
		return
	}
	// 检查请求来源
	if !utils.ValidateGetUrl(ctx, wx.Token) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": 403,
		})
		return
	} else {
		ctx.String(http.StatusOK, ctx.Query("echostr"))
		return
	}
}

func PostMessage(ctx *gin.Context) {

	// 获取当前的用户的微信的配置
	wechatId, err := strconv.Atoi(ctx.Param("id"))
	RedisConn := model.RedisPool.Get()
	defer RedisConn.Close()

	//Redis存放的Key
	wechatRedisKey := "wx_app_data:id:" + strconv.Itoa(wechatId)
	var wx model.Wechat
	//先判断是否存在，不存在就不走下面
	isexit, _ := redis.Bool(RedisConn.Do("EXISTS", wechatRedisKey));
	if isexit {
		wechatBytes, _ := redis.Bytes(RedisConn.Do("GET", wechatRedisKey))
		if len(wechatBytes) > 0 {
			err = json.Unmarshal(wechatBytes, &wx)
			if err != nil {
				log.Println("错误的微信用户ID")
				return
			}
		}
	}

	if wx.ID < 1 {
		if err := model.DB.First(&wx, wechatId).Error; err != nil {
			fmt.Printf(err.Error())
			return
		}
		if err != nil {
			log.Println("错误的微信用户ID")
			return
		}
		wechatBytes, err := json.Marshal(wx)
		if err != nil {
			log.Println("链接Redis出错了")
		}
		RedisConn.Do("SET", wechatRedisKey, wechatBytes)
	}

	redisOpts := &cache.RedisOpts{
		Host:      config.RedisConfig.Host,
		Password:  config.RedisConfig.Password,
		Database:  config.RedisConfig.Database,
		MaxIdle:   config.RedisConfig.MaxIdle,
		MaxActive: config.RedisConfig.MaxActive,
	}

	//使用memcache保存access_token，也可选择redis或自定义cache
	redisCache := cache.NewRedis(redisOpts)

	//配置微信参数
	config := &wechat.Config{
		AppID:          wx.AppID,
		AppSecret:      wx.AppSecret,
		Token:          wx.Token,
		EncodingAESKey: wx.EncodingKey,
		Cache:          redisCache,
	}

	wc := wechat.NewWechat(config)

	// 传入request和responseWriter
	server := wc.GetServer(ctx.Request, ctx.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		switch wx.Robot {
		case 1:
			return model.Recevie(&msg)
		default:
			return model.QykRecevie(&msg)
		}
	})

	//处理消息接收以及回复
	err = server.Serve()
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": 403,
		})
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Println(err)
	}
}
