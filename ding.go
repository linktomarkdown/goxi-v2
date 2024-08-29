package goxi_v2

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DingLogic struct {
	AccessToken string
	Secret      string
	EnableAt    bool
	AtAll       bool
}

func NewDingLogic(params *DingLogic) *DingLogic {
	return &DingLogic{
		AccessToken: params.AccessToken,
		Secret:      params.Secret,
		EnableAt:    params.EnableAt,
		AtAll:       params.AtAll,
	}
}

// Send 发送钉钉机器人webhook消息
func (ding *DingLogic) Send(s string, at ...string) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": s,
		},
	}
	if ding.EnableAt {
		if ding.AtAll {
			if len(at) > 0 {
				return errors.New("the parameter \"AtAll\" is \"true\", but the \"at\" parameter of SendMessage is not empty")
			}
			msg["at"] = map[string]interface{}{
				"isAtAll": ding.AtAll,
			}
		} else {
			msg["at"] = map[string]interface{}{
				"atMobiles": at,
				"isAtAll":   ding.AtAll,
			}
		}
	} else {
		if len(at) > 0 {
			return errors.New("the parameter \"EnableAt\" is \"false\", but the \"at\" parameter of SendMessage is not empty")
		}
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(ding.getURL(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (ding *DingLogic) hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (ding *DingLogic) getURL() string {
	wh := "https://oapi.dingtalk.com/robot/send?access_token=" + ding.AccessToken
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, ding.Secret)
	sign := ding.hmacSha256(stringToSign, ding.Secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", wh, timestamp, sign)
	return url
}
