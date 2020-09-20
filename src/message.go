package main

import (
	"net/url"

	"github.com/matcornic/hermes/v2"
)

type Message struct {
	Recipient string
	Subject   string
	Body      hermes.Email
}

func NewSubscriberThankYouMessage(recipient string) *Message {
	return &Message{
		Recipient: recipient,
		Subject:   "Thanks for subscribing to updates about our wedding!",
		Body: hermes.Email{
			Body: hermes.Body{
				Title:  "Thanks for subscribing to updates about our wedding!",
				Intros: []string{"We're excited to start planning, and we'll let you know as soon as we have more details to share. In the meantime, join the conversation on social media with #rhiphil !"},
				Actions: []hermes.Action{
					{
						Instructions: "To unsubscribe, please click here:",
						Button: hermes.Button{
							Text: "Unsubscribe",
							Link: "https://www.rhiphilwedding.com/unsubscribe?address=" + url.QueryEscape(recipient),
						},
					},
				},
				Outros:    []string{"Looking forward to seeing you there!"},
				Signature: "Sincerely",
			},
		},
	}
}

func NewUnsubscribeConfirmationMessage(address string) *Message {
	return &Message{
		Recipient: address,
		Subject:   "You have unsubscribed from RhiPhil wedding updates.",
		Body: hermes.Email{
			Body: hermes.Body{
				Title:  "You have successfully unsubscribed.",
				Intros: []string{"We're sorry to see you go, but you won't hear from us again."},
				Actions: []hermes.Action{
					{
						Instructions: "To re-subscribe, please click here:",
						Button: hermes.Button{
							Text: "Subscribe",
							Link: "https://www.rhiphilwedding.com/subscribe?address=" + url.QueryEscape(address),
						},
					},
				},
				Signature: "Sincerely",
			},
		},
	}
}

func NewQuestionReceivedMessage(userName, userEmail, question string) *Message {
	return &Message{
		Recipient: userEmail,
		Subject:   "We've received your question",
		Body: hermes.Email{
			Body: hermes.Body{
				Name: userName,
				Intros: []string{
					"Thanks for sending us your question!",
					"We probably haven't seen it yet, but we're sure it was a good one. " +
						"Keep an eye on your inbox for an answer; we'll get back to you as soon as we can.",
				},
				Dictionary: []hermes.Entry{
					{
						Key:   "You Asked",
						Value: question,
					},
					{
						Key:   "We'll send our answer to",
						Value: userEmail,
					},
				},
				Signature: "Sincerely",
			},
		},
	}
}
