package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"

	"github.com/astaxie/beego/logs"
)

type Mailer struct {
	SmtpServer string
	SmtpPort   string
	From       string
	To         []string
	Body       string
	Subject    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewMailer(smtpServer, smtpPort, from, subject string, to []string) *Mailer {
	return &Mailer{
		SmtpServer: smtpServer,
		SmtpPort:   smtpPort,
		From:       from,
		To:         to,
		Subject:    subject,
	}
}

func (m *Mailer) Send(templateName string, items interface{}) {
	err := m.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := m.sendMail(); ok {
		log.Printf("Email has been sent to %s\n", m.To)
	} else {
		log.Printf("Failed to send the email to %s\n", m.To)
	}
}

func (m *Mailer) parseTemplate(filename string, data interface{}) error {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}

	m.Body = buffer.String()

	return nil
}

func (m *Mailer) sendMail() bool {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")
	c, err := smtp.Dial(fmt.Sprintf("%s:%s", m.SmtpServer, m.SmtpPort))
	if err != nil {
		logs.Error("can not send mail to %v, err: %v", m.To, err)
		return false
	}
	defer c.Close()

	if err = c.Mail(r.Replace(m.From)); err != nil {
		logs.Error("can not set source mail , err: %v", m.To, err)
		return false
	}

	var rcptSuccess int = 0
	for i := range m.To {
		if err = c.Rcpt(r.Replace(m.To[i])); err != nil {
			logs.Error("can not set destination mail [%v], err : %v", m.To[i], err)
			continue
		}
		rcptSuccess = rcptSuccess + 1
	}

	if rcptSuccess == 0 {
		logs.Error("can not set destination mail")
		return false
	}

	w, err := c.Data()
	if err != nil {
		logs.Error("create writecloser fail, err =%v", err)
		return false
	}

	body := "To: " + strings.Join(m.To, ",") + "\r\n" +
		"From: " + m.From + "\r\n" +
		"Subject: " + m.Subject + "\r\n" +
		MIME + "\r\n" +
		m.Body

	//logs.Warning("%#v", body)
	_, err = w.Write([]byte(body))
	if err != nil {
		logs.Error("can not write body message, err = %v", err)
		return false
	}

	err = w.Close()
	if err != nil {
		logs.Error("can not close writeCloser, err = %v", err)
		return false
	}

	err = c.Quit()
	if err != nil {
		logs.Error("can not quit client email, err = %v", err)
		return false
	}

	return true
}
