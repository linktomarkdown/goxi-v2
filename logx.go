package goxi_v2

import (
	"github.com/go-resty/resty/v2"
	"log"
)

type LogXLogic struct {
	Endpoint  string
	AppId     string
	AccessKey string
	SecretKey string
}

func NewLogXLogic(endpoint string, appId string, accessKey string, secretKey string) *LogXLogic {
	return &LogXLogic{
		Endpoint:  endpoint,
		AppId:     appId,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

type LogXCreateSchema struct {
	Body     map[string]interface{} `json:"body"`
	Content  string                 `json:"content"`
	Device   map[string]interface{} `json:"device"`
	From     string                 `json:"from"`
	Header   map[string]interface{} `json:"header"`
	Network  map[string]interface{} `json:"network"`
	Response map[string]interface{} `json:"response"`
	Type     string                 `json:"type"`
}

// Up 上报日志,非线程阻塞，不会等待返回结果
func (l *LogXLogic) Up(schema LogXCreateSchema) {
	go func() {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("AccessKey", l.AccessKey).
			SetHeader("SecretKey", l.SecretKey).
			SetBody(schema).
			Post(l.Endpoint + "/logx/log")
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		log.Println(resp.String())
	}()
}

// UpMin 小日志上报
func (l *LogXLogic) UpMin(t string, msg string) {
	var param = LogXCreateSchema{
		Body:     map[string]interface{}{},
		Content:  msg,
		Device:   map[string]interface{}{},
		From:     l.AppId,
		Header:   map[string]interface{}{},
		Network:  map[string]interface{}{},
		Response: map[string]interface{}{},
		Type:     t,
	}
	l.Up(param)
}
