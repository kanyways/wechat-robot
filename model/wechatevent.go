/**
 * Project Name:wechat-robot
 * File Name:event.go
 * Package Name:wechat
 * Date:2020年06月12日 11:22
 * Function:
 * Copyright (c) 2020, Jason.Wang All Rights Reserved.
 */

package model

import (
	"fmt"
	"github.com/silenceper/wechat/message"
)

func EventReceived(msg *message.MixMessage) *message.Reply {
	var text *message.Text
	switch msg.Event {
	// 取消订阅
	case message.EventUnsubscribe:
		currentUser := string(msg.ToUserName)
		userId := string(msg.FromUserName)
		RedisConn := RedisPool.Get()
		defer RedisConn.Close()
		unSubscribeKey := fmt.Sprintf("wx:unsubscribe:%s:user:%s", currentUser, userId)
		RedisConn.Do("SETEX", unSubscribeKey, 24*60*60, 1)
		return nil
	default:
		// 默认使用文本信息处理
		// message.EventSubscribe
		text = message.NewText("你来了？\n欢迎关注，o(∩_∩)o 哈哈")
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	}
}
