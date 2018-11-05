package main

import (
	"SendingMail/utils"
	"time"
)

func main() {
	subject := "Information from System"
	from := "noreply@mail.com"
	to := []string{"abc@gmail.com"}
	smtpServer := "smtp.mail.com"
	smtpPort := "587"
	mailer := utils.NewMailer(smtpServer, smtpPort, from, subject, to)
	items := map[string]string{
		"provider":   "PROVIDER",
		"filename":   "Test file",
		"uploadDay":  time.Now().Format("02/01/2006"),
		"uploadTime": time.Now().Format("02/01/2006 15:04:05"),
	}
	mailer.Send("templates/mail-template.html", items)
}
