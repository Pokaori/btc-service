package utils

import (
	"fmt"
	"log"

	gomail "gopkg.in/gomail.v2"
)

type EmailNotifier interface {
	SendEmails(emails []string)
}

type EmailBTCtoUAHNotifier struct {
	Host     string
	Port     int
	From     string
	Password string
	Rate     float64
}

func (notifier *EmailBTCtoUAHNotifier) SendEmails(emails []string) {
	subject := "Changes in BTC to UAH currency rate"
	body := "The rate is " + fmt.Sprint(notifier.Rate)

	c := make(chan string)
	for _, email := range emails {
		go notifier.sendEmail(c, email, subject, body)
	}
	for i := 0; i < len(emails); i++ {
		log.Println(<-c)
	}
	log.Println("Emails Sent!")
}

func (notifier *EmailBTCtoUAHNotifier) sendEmail(c chan string, email string, subject string, body string) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", notifier.From)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	n := gomail.NewDialer("smtp.gmail.com", 587, notifier.From, notifier.Password)

	if err := n.DialAndSend(msg); err != nil {
		c <- "Error for email " + email + ":\n" + err.Error()
	} else {
		c <- "Successed for email " + email
	}

}
