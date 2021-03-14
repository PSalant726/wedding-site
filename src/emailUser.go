package main

import (
	"fmt"
	"log"
	"strings"

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
				Logo:        "https://raw.githubusercontent.com/PSalant726/wedding-site/master/assets/images/logo.png",
				Copyright:   "Copyright © 2020 Phil Salant. All rights reserved.",
				TroubleText: "Can't {ACTION}? Copy and paste this URL into your web browser instead:",
			},
			TextDirection: hermes.TDLeftToRight,
		},
	}
}

func (eu *EmailUser) SendSubscriberNotification(user string, isSubscribing bool) error {
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
		return fmt.Errorf("failed to send subscriber notification: %w", err)
	}

	return nil
}

func (eu *EmailUser) SendQuestionNotification(userName, userEmail, question string) error {
	var (
		m    = gomail.NewMessage()
		body = fmt.Sprintf("%s (%s) has asked the following question:\n\n%s", userName, userEmail, question)
	)

	m.SetBody("text/plain", body)
	m.SetHeaders(map[string][]string{
		"From":     {m.FormatAddress(eu.Username, "Wedding Guest Questions")},
		"To":       {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"Reply-To": {m.FormatAddress(userEmail, userName)},
		"Subject":  {fmt.Sprintf("[Guest Question] %s has asked a question", userName)},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send question notification: %w", err)
	}

	return nil
}

func (eu *EmailUser) SendRSVPNotification(rsvp RSVP) error {
	var (
		m    = gomail.NewMessage()
		body = fmt.Sprintf(
			"%s has RSVP-ed!\n\nEmail:\n%s\n\nGuests:\n%s\n\nAttending:\n%t\n\nMessage:\n%s",
			rsvp.Name,
			rsvp.Email,
			strings.Join(rsvp.Guests, ", "),
			rsvp.Attending,
			rsvp.Message,
		)
	)

	m.SetBody("text/plain", body)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "Wedding RSVP's")},
		"To":      {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"Subject": {fmt.Sprintf("[RSVP] %s has RSVP-ed!", rsvp.Name)},
	})

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send RSVP notification: %w", err)
	}

	return nil
}

func (eu *EmailUser) GetGomailMessage(message Message) (*gomail.Message, error) {
	plainText, _ := eu.Hermes.GeneratePlainText(message.Body)

	html, err := eu.Hermes.GenerateHTML(message.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML for message: %w", err)
	}

	m := gomail.NewMessage()
	m.SetBody("text/plain", plainText)
	m.AddAlternative("text/html", html)
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(eu.Username, "RhiPhil Wedding")},
		"To":      {m.FormatAddress(message.Recipient, "")},
		"Subject": {message.Subject},
	})

	return m, nil
}

func (eu *EmailUser) SendHermesMessage(message Message) error {
	m, err := eu.GetGomailMessage(message)
	if err != nil {
		return err
	}

	if err := eu.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (eu *EmailUser) SendSubscriberCommunication(subscriberList map[string]string, communication string) error {
	sender, err := eu.Dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to authenticate with email server: %w", err)
	}
	defer sender.Close()

	var successfulSends int

	for emailAddress, name := range subscriberList {
		message, err := eu.GetGomailMessage(*NewSubscriberCommunicationMessage(name, emailAddress, communication))
		if err != nil {
			log.Printf("Failed to create message for %s: %v", emailAddress, err)
			continue
		}

		if err := gomail.Send(sender, message); err != nil {
			log.Printf("Failed to send message to %s: %v", emailAddress, err)
			continue
		}

		successfulSends++
	}

	log.Printf("Successfully sent messages to %d of %d subscribers", successfulSends, len(subscriberList))

	return nil
}
