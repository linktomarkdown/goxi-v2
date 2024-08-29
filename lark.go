package goxi_v2

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

type LarkMessage struct {
	Config struct {
		WideScreenMode bool `json:"wide_screen_mode"`
	} `json:"config"`
	Elements []struct {
		Tag     string `json:"tag"`
		Content string `json:"content,omitempty"`
		Actions []struct {
			Tag  string `json:"tag"`
			Text struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"text"`
			Type     string `json:"type"`
			MultiUrl struct {
				Url        string `json:"url"`
				AndroidUrl string `json:"android_url"`
				IosUrl     string `json:"ios_url"`
				PcUrl      string `json:"pc_url"`
			} `json:"multi_url"`
		} `json:"actions,omitempty"`
	} `json:"elements"`
	Header struct {
		Template string `json:"template"`
		Title    struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"title"`
	} `json:"header"`
}

type LarkLogic struct {
	LarkWebHook string
	Struct2Json func(v interface{}) string
}

func NewLarkLogic(larkURI string) *LarkLogic {
	return &LarkLogic{
		LarkWebHook: larkURI,
		Struct2Json: Struct2Json,
	}
}

// Send 有交互按钮
func (b *LarkLogic) Send(title string, content string) {
	// 组装消息构造体
	msg := LarkMessage{
		Config: struct {
			WideScreenMode bool `json:"wide_screen_mode"`
		}{
			WideScreenMode: true,
		},
		Elements: []struct {
			Tag     string `json:"tag"`
			Content string `json:"content,omitempty"`
			Actions []struct {
				Tag  string `json:"tag"`
				Text struct {
					Tag     string `json:"tag"`
					Content string `json:"content"`
				} `json:"text"`
				Type     string `json:"type"`
				MultiUrl struct {
					Url        string `json:"url"`
					AndroidUrl string `json:"android_url"`
					IosUrl     string `json:"ios_url"`
					PcUrl      string `json:"pc_url"`
				} `json:"multi_url"`
			} `json:"actions,omitempty"`
		}{
			{
				Tag:     "markdown",
				Content: content,
			},
			{
				Tag: "action",
				Actions: []struct {
					Tag  string `json:"tag"`
					Text struct {
						Tag     string `json:"tag"`
						Content string `json:"content"`
					} `json:"text"`
					Type     string `json:"type"`
					MultiUrl struct {
						Url        string `json:"url"`
						AndroidUrl string `json:"android_url"`
						IosUrl     string `json:"ios_url"`
						PcUrl      string `json:"pc_url"`
					} `json:"multi_url"`
				}{
					{
						Tag: "button",
						Text: struct {
							Tag     string `json:"tag"`
							Content string `json:"content"`
						}{
							Tag:     "plain_text",
							Content: "立即前往",
						},
						Type: "primary",
						MultiUrl: struct {
							Url        string `json:"url"`
							AndroidUrl string `json:"android_url"`
							IosUrl     string `json:"ios_url"`
							PcUrl      string `json:"pc_url"`
						}{
							Url:        content,
							AndroidUrl: "",
							IosUrl:     "",
							PcUrl:      "",
						},
					},
				},
			},
		},
		Header: struct {
			Template string `json:"template"`
			Title    struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"title"`
		}{
			Template: "blue",
			Title: struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			}{
				Content: title,
				Tag:     "plain_text",
			},
		},
	}
	client := resty.New()
	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"msg_type": "interactive",
			"card":     b.Struct2Json(msg),
		}).
		Post(b.LarkWebHook)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if post.StatusCode() != 200 {
		log.Fatal("发送飞书消息失败")
		return
	}
}

// SendNoAction 没有交互按钮
func (b *LarkLogic) SendNoAction(title string, content string) {
	// 组装消息构造体
	msg := LarkMessage{
		Config: struct {
			WideScreenMode bool `json:"wide_screen_mode"`
		}{
			WideScreenMode: true,
		},
		Elements: []struct {
			Tag     string `json:"tag"`
			Content string `json:"content,omitempty"`
			Actions []struct {
				Tag  string `json:"tag"`
				Text struct {
					Tag     string `json:"tag"`
					Content string `json:"content"`
				} `json:"text"`
				Type     string `json:"type"`
				MultiUrl struct {
					Url        string `json:"url"`
					AndroidUrl string `json:"android_url"`
					IosUrl     string `json:"ios_url"`
					PcUrl      string `json:"pc_url"`
				} `json:"multi_url"`
			} `json:"actions,omitempty"`
		}{
			{
				Tag:     "markdown",
				Content: content,
			},
		},
		Header: struct {
			Template string `json:"template"`
			Title    struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"title"`
		}{
			Template: "blue",
			Title: struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			}{
				Content: title,
				Tag:     "plain_text",
			},
		},
	}
	client := resty.New()
	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"msg_type": "interactive",
			"card":     b.Struct2Json(msg),
		}).
		Post(b.LarkWebHook)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if post.StatusCode() != 200 {
		log.Fatal("发送飞书消息失败")
		return
	}
}

func Struct2Json(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
