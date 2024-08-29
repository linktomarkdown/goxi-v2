package goxi_v2

import (
	"github.com/go-resty/resty/v2"
	"log"
)

type BarkLogic struct {
	BarkURI  string
	BarkIcon string
}

func NewBarkLogic(urlURI string, iconURI string) *BarkLogic {
	return &BarkLogic{
		BarkURI:  urlURI,
		BarkIcon: iconURI,
	}
}

func (b *BarkLogic) Send(key string, groupName string, title string, content string) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"title":      title,
			"body":       content,
			"sound":      "healthnotification",
			"icon":       b.BarkIcon,
			"group":      groupName,
			"isArchive":  1,
			"device_key": key,
		}).
		Post(b.BarkURI)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println(resp.String())
}

// SendMulti 发送多条
func (b *BarkLogic) SendMulti(key []string, groupName string, title string, content string) {
	for _, v := range key {
		b.Send(v, groupName, title, content)
	}
}
