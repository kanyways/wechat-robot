package model

import (
	"net/url"
	"github.com/kanyways/wechat-robot/utils"
	"log"
	"encoding/json"
	"github.com/kanyways/wechat-robot/config"
	"github.com/silenceper/wechat/message"
	"strings"
)

// 机器人回复的数据的json绑定对象
type Robot struct {
	code int
	Text string `json:"text"`
}

func sendMessage(fromUserName string, content string) string {
	var info = url.QueryEscape(content);
	var requestUrl = "http://www.tuling123.com/openapi/api?key=" + config.ServerConfig.ApiKey + "&info=" + info + "&userid=" + fromUserName
	r, err := utils.Get(requestUrl, nil)
	if err != nil {
		log.Println("请求失败")
		return content
	}
	reply := new(Robot)
	json.Unmarshal(r, reply)
	return reply.Text
}

func SendMessageV2(fromUserName string, content string) string {
	var info = url.QueryEscape(content);
	var requestUrl = "http://www.tuling123.com/openapi/api?key=" + config.ServerConfig.ApiKey + "&info=" + info + "&userid=" + fromUserName
	r, err := utils.Get(requestUrl, nil)
	if err != nil {
		log.Println("请求失败")
		return content
	}
	reply := new(Robot)
	json.Unmarshal(r, reply)
	return reply.Text
}

func Recevie(msg *message.MixMessage) *message.Reply {
	var text *message.Text
	switch msg.MsgType {
	case message.MsgTypeVoice:
		//将语音使用文本处理
		if len(msg.Recognition) > 0 {
			msg.Content = msg.Recognition
			content := sendMessage(string(msg.FromUserName), msg.Content)
			text = message.NewText(content)
		} else {
			text = message.NewText("是风太大还是我没有听清，要不然再说一下？")
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	case message.MsgTypeImage:
		//将图片发送回去
		image := message.NewImage(msg.MediaID)
		return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
	default: //默认使用文本信息处理
		if strings.Contains(msg.Content, "测试的") {
			text = message.NewText("你想要干嘛？")
		} else {
			content := sendMessage(string(msg.FromUserName), msg.Content)
			text = message.NewText(content)
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
