package utils

import (
	"douban/utils/logs"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(sub, format string, a ...interface{}) {
	emailConf, err := GetEmailConf()
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("sign <%s>", emailConf.Username)
	e.To = []string{emailConf.Username}
	e.Subject = sub
	msg := fmt.Sprintf(format, a...)
	e.Text = []byte(msg)

	err = e.Send("smtp.qq.com:587", smtp.PlainAuth("", emailConf.Username, emailConf.Password, "smtp.qq.com"))
	if err != nil {
		logs.Logger.Error("%s", err)
	}
}
