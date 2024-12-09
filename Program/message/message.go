package message

import (
	"fmt"
	"github.com/go-mail/mail/v2"
)

// Отправка сообщений о новом ip адресе
func SendMessage() {
	m := mail.NewMessage()
	m.SetHeader("From", "email@mail.ru")
	m.SetHeader("To", "user_email@gmail.com")
	m.SetHeader("Subject", "New ip address")
	m.SetBody("text/plain", "Work has been detected on your account from a different IP address.")

	fmt.Printf("From: %s,\nTo: %s,\nSubject: %s\n", m.GetHeader("From"), m.GetHeader("To"), m.GetHeader("Subject"))

}
