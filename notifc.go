package main

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"
)

func push_notifcation(file string, files_has_change []string) {
	current_time := time.Now()
	format_time := fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d]",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second(),
	)

	slack_header := make(map[string]string)

	slack_body := []byte(fmt.Sprintf(`{"text":"*There are changes in this urls* [%s] %s\n\t%s"}`, file, format_time, strings.Join(files_has_change, "\n")))

	slack_header["Content-Type"] = "application/json"

	if data.Config.Channel_use == "slack" {
		send_request(
			data.Config.Channel.Slack, "POST", slack_header, slack_body,
		)

	} else if data.Config.Channel_use == "mail" {

		email_body := fmt.Sprintf(`<h3>There are changes in this urls [%s]%s</h3> %s`, file, format_time, strings.Join(files_has_change, "<br>"))
		email_title := fmt.Sprintf(`Change file %s`, format_time)
		m := gomail.NewMessage()

		m.SetHeader("From", data.Config.Channel.Mail.From)
		m.SetHeader("To", data.Config.Channel.Mail.To)
		m.SetHeader("Subject", email_title)
		m.SetBody("text/html", email_body)

		d := gomail.NewDialer(
			data.Config.Channel.Mail.Host,
			data.Config.Channel.Mail.Port,
			data.Config.Channel.Mail.Email,
			data.Config.Channel.Mail.Password,
		)

		d.TLSConfig = &tls.Config{InsecureSkipVerify: data.Config.Channel.Mail.Tls}

		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
		}
	}
}
