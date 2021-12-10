package email

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"redisData/pkg/email/google"
	"redisData/pkg/email/qq"
	"redisData/pkg/logger"
	"strings"
)

func SendEmail(title string, text string, toId string) error {
	var err error
	emailType := viper.GetString("email.email_type")
	if emailType == "QQ" {
		err = qq.SendEmail(title, viper.GetString("email.qq_email_user"), toId, text)
		if err != nil {
			return err
		}
	}
	if emailType == "GOOGLE" {
		err = google.New().Send(title, text, toId)
		if err != nil {
			return err
		}
	}
	return nil
}


func SendDingMsg(name ,msg string) {
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=8dc61e631bf181a624549409bee798031b83d845bd1d27791f8dd4f66a93a4d4`
	content := `{"msgtype": "text",
		"text": {"content": "` + name+ ":" + msg + `"}
	}`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		logger.Error(err)
		return
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		// handle error
		fmt.Println(err)
	}
}
