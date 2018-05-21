package model

type Wechat struct {
	ID          uint32 `json:"id"`
	UserID      uint32 `json:"userId"`
	AppID       string `json:"appId"`
	AppSecret   string `json:"appSecret"`
	Token       string `json:"token"`
	EncodingKey string `json:"encodingKey"`
	Remark      string `json:"remark"`
	Robot       uint32 `json:"robot"`
}
