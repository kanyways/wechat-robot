/**
 * Project Name:wechat-robot
 * File Name:txairobot.go
 * Package Name:model
 * Date:2019年07月16日 14:19
 * Function:
 * Copyright (c) 2019, Jason.Wang All Rights Reserved.
 */
package model

import (
	"encoding/json"
	"github.com/kanyways/wechat-robot/config"
	"github.com/kanyways/wechat-robot/utils"
	"github.com/silenceper/wechat/message"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type AiRequest struct {
	AppId     float64 `json:"app_id"`
	TimeStamp int64   `json:"time_stamp"`
	NonceStr  string  `json:"nonce_str"`
	Sign      string  `json:"sign"`
	Session   string  `json:"session"`
	Question  string  `json:"question"`
}

type AiResponseData struct {
	Session string `json:"session"`
	Answer  string `json:"answer"`
}

type AiResponse struct {
	Ret  int64          `json:"ret"`
	Msg  string         `json:"msg"`
	Data AiResponseData `json:"data"`
}

//发送数据
func SendTxAiMessage(fromUserName string, content string) (data string) {
	aiRequest := AiRequest{}
	aiRequest.AppId = config.ServerConfig.TxApiKey
	aiRequest.TimeStamp = time.Now().Unix()
	aiRequest.NonceStr = strconv.FormatInt(time.Now().Unix()+rand.Int63(), 10)
	aiRequest.Session = fromUserName
	aiRequest.Question = strings.TrimSpace(content)
	aiRequest.Sign = utils.GetSignEncodeToUpper(aiRequest, "&app_key="+config.ServerConfig.TxApiSecret)

	r, err := utils.Post("https://api.ai.qq.com/fcgi-bin/nlp/nlp_textchat", utils.Struct2Map(aiRequest))
	if err != nil {
		log.Println("请求失败")
		return
	}
	log.Println(string(r))

	reply := AiResponse{}
	json.Unmarshal(r, &reply)
	if reply.Ret == 0 {
		data = reply.Data.Answer
	}
	return
}

func TxAiRecevie(msg *message.MixMessage) *message.Reply {
	var text *message.Text
	switch msg.MsgType {
	case message.MsgTypeVoice:
		//将语音使用文本处理
		if len(msg.Recognition) > 0 {
			msg.Content = msg.Recognition
			content := SendTxAiMessage(string(msg.FromUserName), msg.Content)
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
		if strings.Contains(msg.Content, "测试") {
			text = message.NewText("你想要干嘛？")
		} else {
			content := SendTxAiMessage(string(msg.FromUserName), msg.Content)
			text = message.NewText(content)
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
