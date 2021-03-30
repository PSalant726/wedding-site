package main

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"

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
				Title: "Thanks for subscribing to updates about our wedding!",
				Intros: []string{
					"We're excited to continue planning, and we'll let you know as soon as we have more details to share. " +
						"In the meantime, please add no-reply@rhiphilwedding.com to your address book, " +
						"and join the conversation on social media with #rhiphil !",
				},
				Actions: []hermes.Action{
					{
						Instructions: "To unsubscribe, please click here:",
						Button: hermes.Button{
							Color: "#331929",
							Link:  "https://www.rhiphilwedding.com/unsubscribe?address=" + url.QueryEscape(recipient),
							Text:  "Unsubscribe",
						},
					},
				},
				Outros:    []string{"Looking forward to seeing you there!"},
				Signature: "Sincerely",
			},
		},
	}
}

func NewSubscriberCommunicationMessage(name, recipient, messageText string) *Message {
	return &Message{
		Recipient: recipient,
		Subject:   "Important information about Rhiannon & Phil's wedding!",
		Body: hermes.Email{
			Body: hermes.Body{
				Name: name,
				Intros: append(
					[]string{"We have an update about our wedding, and we wanted you to be the first to know:"},
					messageText,
				),
				Actions: []hermes.Action{
					{
						Button: hermes.Button{
							Color: "#83D3C9",
							Link:  "https://www.rhiphilwedding.com",
							Text:  "Visit RhiPhilWedding.com",
						},
					},
					{
						Instructions: "You are receiving this message because you previously subscribed to receive updates about our wedding. " +
							"To ensure continued delivery, please add no-reply@rhiphilwedding.com to your address book. " +
							"To unsubscribe, click here:",
						Button: hermes.Button{
							Color: "#331929",
							Link:  "https://www.rhiphilwedding.com/unsubscribe?address=" + url.QueryEscape(recipient),
							Text:  "Unsubscribe",
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
							Color: "#83D3C9",
							Link:  "https://www.rhiphilwedding.com/subscribe?address=" + url.QueryEscape(address),
							Text:  "Subscribe",
						},
					},
				},
				Signature: "Sincerely",
			},
		},
	}
}

func NewQuestionReceivedMessage(sender *mail.Address, question string) *Message {
	return &Message{
		Recipient: sender.Address,
		Subject:   "We've received your question.",
		Body: hermes.Email{
			Body: hermes.Body{
				Name: sender.Name,
				Intros: []string{
					"Thanks for sending us your question!",
					"We probably haven't seen it yet, but we're sure it was a good one. " +
						"Keep an eye on your inbox for an answer; we'll get back to you as soon as we can. " +
						"To ensure delivery of your answer, please add no-reply@rhiphilwedding.com to your address book.",
				},
				Dictionary: []hermes.Entry{
					{Key: "You Asked", Value: question},
					{Key: "We'll send our answer to", Value: sender.Address},
				},
				Signature: "Sincerely",
			},
		},
	}
}

func NewRSVPConfirmationMessage(rsvp RSVP) *Message {
	msg := &Message{
		Recipient: rsvp.Email,
		Subject:   fmt.Sprintf("Thanks for your RSVP to our wedding %s!", rsvp.Event),
		Body: hermes.Email{
			Body: hermes.Body{
				Name:      strings.Split(rsvp.Name, " ")[0],
				Intros:    rsvp.ConfirmationIntros(),
				Table:     rsvp.ConfirmationTable(),
				Actions:   rsvp.ConfirmationActions(),
				Signature: "Sincerely",
			},
		},
	}

	if rsvp.Attending {
		msg.Body.Body.Outros = []string{"We're looking forward to seeing you there!"}
	}

	return msg
}
