package main

import (
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
						"In the meantime, join the conversation on social media with #rhiphil !",
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
						Instructions: "You are receiving this message because you previously subscribed to receive updates " +
							"about our wedding. To unsubscribe, please click here:",
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

func NewQuestionReceivedMessage(userName, userEmail, question string) *Message {
	return &Message{
		Recipient: userEmail,
		Subject:   "We've received your question.",
		Body: hermes.Email{
			Body: hermes.Body{
				Name: userName,
				Intros: []string{
					"Thanks for sending us your question!",
					"We probably haven't seen it yet, but we're sure it was a good one. " +
						"Keep an eye on your inbox for an answer; we'll get back to you as soon as we can.",
				},
				Dictionary: []hermes.Entry{
					{Key: "You Asked", Value: question},
					{Key: "We'll send our answer to", Value: userEmail},
				},
				Signature: "Sincerely",
			},
		},
	}
}

func NewRSVPConfirmationMessage(rsvp RSVP) *Message {
	guests := strings.Join(rsvp.Guests, ", ")
	if guests == "" {
		guests = "None"
	}

	message := rsvp.Message
	if message == "" {
		message = "None"
	}

	if rsvp.Rehearsal {
		return NewRehearsalRSVPConfirmationMessage(rsvp, guests, message)
	}

	msg := &Message{
		Recipient: rsvp.Email,
		Subject:   "Thanks for your RSVP!",
		Body: hermes.Email{
			Body: hermes.Body{
				Name:   strings.Split(rsvp.Name, " ")[0],
				Intros: []string{"This email confirms that your RSVP has been received, thanks!"},
				Table: hermes.Table{
					Data: [][]hermes.Entry{
						{
							{Key: "Field", Value: "Name"},
							{Key: "Your Reply:", Value: rsvp.Name},
						},
						{
							{Key: "Field", Value: "Guests"},
							{Key: "Your Reply:", Value: guests},
						},
						{
							{Key: "Field", Value: "Message"},
							{Key: "Your Reply:", Value: message},
						},
						{{Key: "Field", Value: "Attending"}},
					},
					Columns: hermes.Columns{
						CustomWidth: map[string]string{"Field": "25%"},
					},
				},
				Actions: []hermes.Action{
					{
						Instructions: "Need to change your response? You can do that any time through April 30th by clicking here:",
						Button: hermes.Button{
							Color: "#83D3C9",
							Link:  "https://www.rhiphilwedding.com/rsvp",
							Text:  "Change your RSVP",
						},
					},
				},
				Signature: "Sincerely",
			},
		},
	}

	detail := "We're sorry to hear you won't be able to join us, " +
		"but we hope to be able to catch up soon. If things change, " +
		"and it turns out you can be there, you can change your " +
		"response any time through April 30th using the link below."
	isAttending := hermes.Entry{Key: "Your Reply:", Value: "No"}

	if rsvp.Attending {
		detail = "We're glad to hear you can make it! " +
			"There's nothing more you need to do. " +
			"We can't wait to celebrate with you on June 5th!"

		isAttending = hermes.Entry{Key: "Your Reply:", Value: "Yes!"}
		msg.Body.Body.Outros = []string{"We're looking forward to seeing you there!"}
		msg.Body.Body.Actions = append(
			msg.Body.Body.Actions,
			hermes.Action{
				Instructions: "Need to book a place to stay? " +
					"Get more information about travel and accommodations here:",
				Button: hermes.Button{
					Color: "#331929",
					Link:  "https://www.rhiphilwedding.com/travel",
					Text:  "Get Traveler Information",
				},
			},
		)
	}

	msg.Body.Body.Intros = append(msg.Body.Body.Intros, detail, "Here's how you replied:")
	msg.Body.Body.Table.Data[len(msg.Body.Body.Table.Data)-1] = append(
		msg.Body.Body.Table.Data[len(msg.Body.Body.Table.Data)-1], isAttending,
	)

	return msg
}

func NewRehearsalRSVPConfirmationMessage(rsvp RSVP, guests, message string) *Message {
	msg := &Message{
		Recipient: rsvp.Email,
		Subject:   "Thanks for your RSVP to our Rehearsal Dinner!",
		Body: hermes.Email{
			Body: hermes.Body{
				Name:   strings.Split(rsvp.Name, " ")[0],
				Intros: []string{"This email confirms that we've received your RSVP to our rehearsal dinner, thanks!"},
				Table: hermes.Table{
					Data: [][]hermes.Entry{
						{
							{Key: "Field", Value: "Name"},
							{Key: "Your Reply:", Value: rsvp.Name},
						},
						{
							{Key: "Field", Value: "Guests"},
							{Key: "Your Reply:", Value: guests},
						},
						{
							{Key: "Field", Value: "Message"},
							{Key: "Your Reply:", Value: message},
						},
						{{Key: "Field", Value: "Attending"}},
					},
					Columns: hermes.Columns{
						CustomWidth: map[string]string{"Field": "25%"},
					},
				},
				Actions: []hermes.Action{
					{
						Instructions: "Need to change your response? You can do that any time through April 30th by clicking here:",
						Button: hermes.Button{
							Color: "#83D3C9",
							Link:  "https://www.rhiphilwedding.com/rehearsal",
							Text:  "Change your RSVP",
						},
					},
				},
				Signature: "Sincerely",
			},
		},
	}

	detail := "We're sorry to hear you won't be able to make it, " +
		"but we hope you'll still be able to attend the wedding. " +
		"If things change, and it turns out you can be there, " +
		"you can change your response any time through April 30th using the link below."
	isAttending := hermes.Entry{Key: "Your Reply:", Value: "No"}

	if rsvp.Attending {
		detail = "We're glad to hear you can make it! " +
			"If you haven't also RSVP-ed to the ceremony and reception, " +
			"please be sure to do so by April 30th. Otherwise, " +
			"there's nothing more you need to do. " +
			"We appreciate you taking the time to make sure " +
			"everything goes well on our special day."

		isAttending = hermes.Entry{Key: "Your Reply:", Value: "Yes!"}
		msg.Body.Body.Outros = []string{"We're looking forward to seeing you there!"}
	}

	msg.Body.Body.Intros = append(msg.Body.Body.Intros, detail, "Here's how you replied:")
	msg.Body.Body.Table.Data[len(msg.Body.Body.Table.Data)-1] = append(
		msg.Body.Body.Table.Data[len(msg.Body.Body.Table.Data)-1], isAttending,
	)

	return msg
}
