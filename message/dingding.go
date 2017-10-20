package message

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"udp_server/logs"
)

var APIURL = "https://oapi.dingtalk.com/robot/send?access_token=%s"

const (
	TOKEN = "e4fa2ac5d689ec738e46800d4f57df1198fff5df6932e16de96ed0297917a117"
)

type DingDingMsg struct {
	MsgType  string  `json:"msgtype"`
	Markdown Content `json:"markdown"`
	At       At      `json:"at"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func SendDD(title, content string, phones []string, token string) {
	dm := DingDingMsg{
		MsgType: "markdown",
		Markdown: Content{
			Title: "👻" + title,
			Text:  content,
		},
		At: At{
			AtMobiles: phones,
			IsAtAll:   false,
		},
	}
	if token == "" {
		token = TOKEN
	}
	client := &http.Client{}
	dmm, err := json.Marshal(dm)
	if err != nil {
		logs.Log(nil).Errorf("SendDD err:%v", err)
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf(APIURL, token), nil)
	req.Header.Set("Content-Type", "application/json")

	req.Body = ioutil.NopCloser(strings.NewReader(string(dmm)))
	resp, _ := client.Do(req)
	logs.Log(nil).Infof("dmm=%+v, resp=%+v", string(dmm), resp)
	if resp != nil {
		defer resp.Body.Close()
	}
}
