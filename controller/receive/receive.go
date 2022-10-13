package receive

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/model"
	"github.com/kanyways/wechat-robot/utils"
	"github.com/silenceper/wechat"
	"github.com/silenceper/wechat/cache"
	"github.com/silenceper/wechat/message"
	"log"
	"net/http"
	"strconv"
	"time"
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
	weChatId, err := strconv.Atoi(ctx.Param("id"))
	RedisConn := model.RedisPool.Get()
	defer RedisConn.Close()

	//Redis存放的Key
	weChatRedisKey := fmt.Sprintf("wx_app_data:id:%v", weChatId)
	var wx model.Wechat
	//先判断是否存在，不存在就不走下面
	isexit, _ := redis.Bool(RedisConn.Do("EXISTS", weChatRedisKey))
	if isexit {
		wechatBytes, _ := redis.Bytes(RedisConn.Do("GET", weChatRedisKey))
		if len(wechatBytes) > 0 {
			err = json.Unmarshal(wechatBytes, &wx)
			if err != nil {
				log.Println("错误的微信用户ID")
				return
			}
		}
	}

	if wx.ID < 1 {
		if err := model.DB.First(&wx, weChatId).Error; err != nil {
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
		RedisConn.Do("SET", weChatRedisKey, wechatBytes)
	}

	// 由于使用的插件会去校验，所以此处的校验取消。
	//// 检查请求来源
	//if !utils.ValidateGetUrl(ctx, wx.Token) {
	//	log.Println("Wechat Service: this http request is not from Wechat platform!")
	//	ctx.JSON(http.StatusForbidden, gin.H{
	//		"message": 403,
	//	})
	//	return
	//}

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
		//data, err := json.Marshal(msg)
		//if err != nil {
		//	log.Fatalf("Json marshaling failed：%s", err)
		//}
		//fmt.Printf("%s\n", data)
		switch msg.MsgType {
		case message.MsgTypeEvent:
			return model.EventReceived(&msg)
		default:
			text := checkHandler(&msg)
			if text != nil {
				return text
			}
			switch wx.Robot {
			case 1:
				return model.TxAiRecevie(&msg)
			default:
				return model.TxAiRecevie(&msg)
			}

		}

		//回复消息：演示回复用户发送的消息
		//var text *message.Text
		//Golang 判断字符串中是否包含 chars 中的任何一个字符
		//fmt.Println(strings.ContainsAny("widuu", "wi")) //true

		//Golang 判断字符串中是否包含其他某字符
		//fmt.Println(strings.Contains("widuu", "wi")) //true

		// see https://studygolang.com/articles/3447
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

	// 读取所有的Body数据，将xml转化成Map
	//defer c.Request.Body.Close()
	//var body = make([]byte, 0)
	//var len = 0
	//for {
	//	var buffer = make([]byte, 1024)
	//	n, err := c.Request.Body.Read(buffer)
	//	if err != nil && err != io.EOF {
	//		c.JSON(http.StatusInternalServerError, gin.H{
	//			"status": http.StatusInternalServerError,
	//			"errMsg": http.StatusText(http.StatusInternalServerError),
	//		})
	//		return
	//	}
	//	body = append(body, buffer[:n]...)
	//	len = len + n
	//	if err == io.EOF {
	//		break
	//	}
	//}
	//var rm map[string]string
	//xml.Unmarshal(body, (*model.Map)(&rm))
	//fmt.Println(rm["name"])
}

func checkHandler(msg *message.MixMessage) *message.Reply {
	RedisConn := model.RedisPool.Get()
	defer RedisConn.Close()
	currentUser := string(msg.ToUserName)
	userId := string(msg.FromUserName)
	unSubscribeKey := fmt.Sprintf("wx:unsubscribe:%s:user:%s", currentUser, userId)

	isExist, _ := redis.Bool(RedisConn.Do("EXISTS", unSubscribeKey))
	if isExist {
		currentTime := time.Now().Unix()
		ttlTime, _ := redis.Int64(RedisConn.Do("TTL", unSubscribeKey))
		text := fmt.Sprintf("由于您之前的操作过于高端，\n特决定您在%s之后才能正常使用功能。", utils.FormatS(time.Unix(currentTime+ttlTime, 0), "2006-01-02 15:04:05"))
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(text)}
	}
	return nil
}
