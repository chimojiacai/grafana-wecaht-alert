package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type Hook struct {
	DashboardId string `json:"dashboardId"`
	EvalMatches string `json:"evalMatches"`
	ImageUrl    string `json:"imageUrl"`
	Message     string `json:"message"`
	OrgId       string `json:"orgId"`
	PanelId     string `json:"panelId"`
	RuleId      string `json:"ruleId"`
	RuleName    string `json:"ruleName"`
	RuleUrl     string `json:"ruleUrl"`
	State       string `json:"state"`
	Tags        string `json:"tags"`
	Title       string `json:"title"`
}

var sentCount = 0

const (
	Url         = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	OKMsg       = "告警恢复"
	AlertingMsg = "触发告警"
	OK          = "OK"
	Alerting    = "Alerting"
	ColorGreen  = "info"
	ColorGray   = "comment"
	ColorRed    = "warning"
)

// 记录发送次数
func GetSendCount(c *gin.Context) {
	_, _ = c.Writer.WriteString("G2WW Server created by Nova Kwok is running! \nParsed & forwarded \" + strconv.Itoa(sentCount) + \" messages to WeChat Work!")
	return
}

// 发送消息
func SendMsg(c *gin.Context) {
	h := &Hook{}
	if err := c.BindJSON(&h); err != nil {
		fmt.Println(err)
		_, _ = c.Writer.WriteString("Error on JSON format")
		return
	}

	marshal, _ := json.Marshal(h)
	fmt.Println("接受参数数据：", string(marshal))
	// 字符串替换
	h.RuleUrl = strings.ReplaceAll(h.RuleUrl, ":3000", "")
	color := ColorGreen
	if strings.Contains(h.Title, OK) {
		h.Title = strings.ReplaceAll(h.Title, OK, OKMsg)
	} else {
		h.Title = strings.ReplaceAll(h.Title, Alerting, AlertingMsg)
		color = ColorRed
	}

	// Send to WeChat Work
	url := Url + c.Query("key")
	// 处理数据格式
	msgStr := MsgMarkdown(h, color)
	if c.Query("type") == "news" {
		msgStr = MsgNews(h)
	}

	fmt.Println("发送的消息是：", msgStr)

	jsonStr := []byte(msgStr)
	// 发送http请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = c.Writer.WriteString("Error sending to WeChat Work API")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("shuju:", string(body))

	_, _ = c.Writer.Write(body)
	sentCount++

	return
}

// 发送消息类型 news
func MsgNews(h *Hook) string {
	return fmt.Sprintf(`
		{
			"msgtype": "news",
			"news": {
			  "articles": [
				{
				  "title": "%s",
				  "description": "%s",
				  "url": "%s",
				  "picurl": "%s"
				}
			  ]
			}
		  }
		`, h.Title, h.Message, h.RuleUrl, h.ImageUrl)
}

// 发送消息类型
func MsgMarkdown(h *Hook, color string) string {
	return fmt.Sprintf(`
	{
       "msgtype": "markdown",
       "markdown": {
           "content": "<font color=\"%s\">%s</font>\r\n<font color=\"comment\">%s\r\n[点击查看详情](%s)![](%s)</font>",
            "mentioned_list":["龙","@all"],
            "mentioned_mobile_list":["13701390306","@all"]
       }
  }`, color, h.Title, h.Message, h.RuleUrl, h.ImageUrl)
}
