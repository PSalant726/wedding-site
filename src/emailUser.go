package main

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"gopkg.in/gomail.v2"
)

type EmailUser struct {
	Username string
	Password string
	Dialer   *gomail.Dialer
	Hermes   hermes.Hermes
}

func NewGmailUser(username, password string) *EmailUser {
	return &EmailUser{
		Username: username,
		Password: password,
		Dialer:   gomail.NewDialer("smtp.gmail.com", 587, username, password),
		Hermes: hermes.Hermes{
			DisableCSSInlining: false,
			Product: hermes.Product{
				Name:        "Phil & Rhiannon",
				Link:        "https://www.rhiphilwedding.com",
				Logo:        "", // TODO: Link to the logo image
				Copyright:   "Copyright © 2020 Phil Salant. All rights reserved.",
				TroubleText: "Can't '{ACTION}'? Copy and paste this URL into your web browser instead:",
			},
			TextDirection: hermes.TDLeftToRight,
		},
	}
}

func (eu *EmailUser) SendNotification(user string, isSubscribing bool) error {
	var (
		m       = gomail.NewMessage()
		subject string
		body    string
	)

	if isSubscribing {
		subject = fmt.Sprintf("[Wedding Details] %s has subscribed", user)
		body = fmt.Sprintf("%s has subscribed to receive updates about wedding details.", user)
	} else {
		subject = fmt.Sprintf("[Wedding Details][Unsubscribe] %s has unsubscribed", user)
		body = fmt.Sprintf("%s has unsubscribed from updates about wedding details.", user)
	}

	m.SetBody("text/plain", body)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"To":      {m.FormatAddress(eu.Username, "Subscriber Notification")},
		"Subject": {subject},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (eu *EmailUser) SendHermesMessage(recipient, subject string, message hermes.Email) error {
	var (
		m            = gomail.NewMessage()
		html, _      = eu.Hermes.GenerateHTML(message)
		plainText, _ = eu.Hermes.GeneratePlainText(message)
	)

	m.SetBody("text/plain", plainText)
	m.AddAlternative("text/html", html)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"To":      {m.FormatAddress(recipient, "")},
		"Subject": {subject},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
