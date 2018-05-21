package model

import (
	"net/url"
	"github.com/kanyways/wechat-robot/utils"
	"log"
	"encoding/json"
	"strings"
	"github.com/silenceper/wechat/message"
)

// 青云客的数据结构
type QyRobot struct {
	code    int
	Content string `json:"content"`
}

func qykSendMessage(content string) string {
	var info = url.QueryEscape(content);
	requestUrl := "http://api.qingyunke.com/api.php?key=free&appid=0&msg=" + info

	r, err := utils.Get(requestUrl, nil)
	if err != nil {
		log.Println("请求失败")
		return content
	}
	reply := new(QyRobot)
	json.Unmarshal(r, reply)

	reply.Content = strings.Replace(reply.Content, "qingyunke.com", "kany.me", -1)
	reply.Content = strings.Replace(reply.Content, "{br}", "\r\n", -1)

	return reply.Content
}

func QykRecevie(msg *message.MixMessage) *message.Reply {
	var text *message.Text
	switch msg.MsgType {
	case message.MsgTypeVoice:
		//将语音使用文本处理
		if len(msg.Recognition) > 0 {
			msg.Content = msg.Recognition
			content := qykSendMessage(msg.Content)
			text = message.NewText(content)
		} else {
			text = message.NewText("是风太大还是我没有听清，要不然再说一下？")
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	case message.MsgTypeImage:
		//将图片发送回去
		image := message.NewImage(msg.MediaID)
		return &message.Reply{MsgType: message.MsgTypeImage, MsgData: image}
	case message.MsgTypeVideo:
		//将视频发送回去
		text = message.NewText("是风太大还是我没有听清，要不然再说一下？")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	case message.MsgTypeShortVideo:
		//将视频发送回去
		text = message.NewText("是风太大还是我没有听清，要不然再说一下？")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	default: //默认使用文本信息处理
		if strings.Contains(msg.Content, "测试的") {
			text = message.NewText("你想要干嘛？")
		} else {
			content := qykSendMessage(msg.Content)
			text = message.NewText(content)
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
