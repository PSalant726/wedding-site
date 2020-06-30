package main

import (
	"net/url"

	"github.com/matcornic/hermes/v2"
)

func NewSubscriberThankYouMessage(recipient string) (hermes.Email, string) {
	return hermes.Email{
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
	}, "Thanks for subscribing to updates about our wedding!"
}

func NewUnsubscribeConfirmationMessage(address string) (hermes.Email, string) {
	return hermes.Email{
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
			Outros:    []string{"Looking forward to seeing you there!"},
			Signature: "Sincerely",
		},
	}, "You have unsubscribed from RhiPhil wedding updates."
}
