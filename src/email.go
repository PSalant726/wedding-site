package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailUser struct {
	Username string
	Password string
	Dialer   *gomail.Dialer
}

func NewGmailUser(username, password string) *EmailUser {
	return &EmailUser{
		Username: username,
		Password: password,
		Dialer:   gomail.NewDialer("smtp.gmail.com", 587, username, password),
	}
}

func (eu *EmailUser) SendSubscribeThankYouMessage(recipient string) error {
	m := gomail.NewMessage()
	m.SetBody("text/plain", "Thanks for subscribing to updates about our wedding! We've added you to our mailing list, and we'll let you know when more details become available.")
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"To":      {m.FormatAddress(recipient, "")},
		"Subject": {"Thanks for subscribing to updates about our wedding!"},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (eu *EmailUser) SendNewSubscriberNotification(subscriber string) error {
	m := gomail.NewMessage()
	m.SetBody("text/plain", fmt.Sprintf("%s has subscribed to receive updates about wedding details.", subscriber))
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"To":      {m.FormatAddress(eu.Username, "Subscriber Notification")},
		"Subject": {"[Wedding Details] " + subscriber + " has subscribed"},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
