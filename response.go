package goxi_v2

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Code = -1
		body.Msg = err.Error()
	} else {
		body.Code = http.StatusOK
		body.Msg = "OK"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

func Fail(w http.ResponseWriter, Code int, Msg string) {
	var body Body
	body.Code = Code
	body.Msg = Msg
	httpx.OkJson(w, body)
}

func failOnError(err error, msg string) {
	if err != nil {
		logx.Error("%s: %s", msg, err)
	}
}
