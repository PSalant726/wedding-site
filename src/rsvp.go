package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/mail"
	"net/url"
	"os"
	"strings"

	"github.com/fabioberger/airtable-go"
	"github.com/matcornic/hermes/v2"
)

var (
	airtableAPIKey    = os.Getenv("AIRTABLE_API_KEY")
	airtableBaseID    = os.Getenv("AIRTABLE_BASE_ID")
	airtableClient, _ = airtable.New(airtableAPIKey, airtableBaseID)
)

type RSVP struct {
	Name      string
	Email     string
	Zip       string
	Guests    []string
	Message   string
	Attending bool
	Rehearsal bool
	Event     string
}

type guest struct {
	AirtableID string
	Fields     struct {
		Guest     string
		Email     string
		ZIP       string
		RSVP      string
		PlusOne   []string
		PlusOneOf []string
	}
}

func NewRSVP(response url.Values, rehearsal string) RSVP {
	rsvp := &RSVP{
		Name:      strings.TrimSpace(strings.Title(template.HTMLEscapeString(response.Get("name")))),
		Zip:       strings.TrimSpace(template.HTMLEscapeString(response.Get("zip"))),
		Guests:    parseGuests(response.Get("guests")),
		Message:   strings.TrimSpace(template.HTMLEscapeString(response.Get("message"))),
		Attending: response.Get("response") == "1",
		Rehearsal: rehearsal == "true",
		Event:     "ceremony & reception",
	}

	email, err := mail.ParseAddress(template.HTMLEscapeString(response.Get("email")))
	if err == nil {
		rsvp.Email = email.Address
	}

	if rsvp.Rehearsal {
		rsvp.Event = "rehearsal dinner"
	}

	return *rsvp
}

func (r RSVP) AttendingText() string {
	if r.Attending {
		return "Yes!"
	}

	return "No"
}

func (r RSVP) ConfirmationActions() []hermes.Action {
	var (
		travelInfo     hermes.Action
		changeResponse = hermes.Action{
			Instructions: "Need to change your response? You can do that any time through April 30th by clicking here:",
			Button: hermes.Button{
				Color: "#83D3C9",
				Link:  "https://www.rhiphilwedding.com/rsvp",
				Text:  "Change your RSVP",
			},
		}
	)

	if r.Rehearsal {
		changeResponse.Button.Link = "https://www.rhiphilwedding.com/rehearsal"
	} else if r.Attending {
		travelInfo = hermes.Action{
			Instructions: "Need to book a place to stay? " +
				"Get more information about travel and accommodations here:",
			Button: hermes.Button{
				Color: "#331929",
				Link:  "https://www.rhiphilwedding.com/travel",
				Text:  "Get Traveler Information",
			},
		}
	}

	return []hermes.Action{changeResponse, travelInfo}
}

func (r RSVP) ConfirmationIntros() []string {
	intros := []string{
		fmt.Sprintf("This email confirms that we've received your RSVP to our %s, thanks!", r.Event),
	}

	detail := "We're sorry to hear you won't be able to join us, " +
		"but we hope to be able to catch up soon. If things change, " +
		"and it turns out you can be there, you can change your " +
		"response any time through April 30th using the link below."

	if r.Attending {
		date := 5
		if r.Rehearsal {
			date = 4
		}

		detail = fmt.Sprintf("We're glad to hear you can make it! "+
			"There's nothing more you need to do. "+
			"We can't wait to celebrate with you on June %dth!", date)
	}

	return append(intros, detail, "Here's how you replied:")
}

func (r RSVP) ConfirmationTable() hermes.Table {
	guests := strings.Join(r.Guests, ", ")
	if guests == "" {
		guests = "None"
	}

	message := r.DecodedMessage()
	if r.Message == "" {
		message = "None"
	}

	return hermes.Table{
		Data: [][]hermes.Entry{
			{
				{Key: "Field", Value: "Name"},
				{Key: "Your Reply:", Value: r.Name},
			},
			{
				{Key: "Field", Value: "Attending"},
				{Key: "Your Reply:", Value: r.AttendingText()},
			},
			{
				{Key: "Field", Value: "Guests"},
				{Key: "Your Reply:", Value: guests},
			},
			{
				{Key: "Field", Value: "Message"},
				{Key: "Your Reply:", Value: message},
			},
		},
		Columns: hermes.Columns{
			CustomWidth: map[string]string{"Field": "25%"},
		},
	}
}

func (r RSVP) DecodedMessage() string {
	return strings.TrimSpace(html.UnescapeString(r.Message))
}

func (r RSVP) Validate() error {
	var (
		responders = make([]guest, 0)
		listParams = airtable.ListParameters{
			Fields:          []string{"Guest", "ZIP"},
			FilterByFormula: r.filterByFormula(),
			MaxRecords:      5,
		}
	)

	airtableClient.ShouldRetryIfRateLimited = true
	if err := airtableClient.ListRecords("Guests", &responders, listParams); err != nil {
		log.Printf("Failed to query Airtable for RSVP: %+v. Error: %s", r, err)
		return fmt.Errorf("An internal error occurred. Please try again.")
	}

	if len(responders) == 0 {
		return fmt.Errorf("Please verify that your name and zip code match your invitation's envelope.")
	}

	return nil
}

func (r RSVP) filterByFormula() string {
	rehearsalFilter := ")"
	if r.Rehearsal {
		rehearsalFilter = ", NOT({Rehearsal Dinner RSVP} = 'Not Invited'))"
	}

	return fmt.Sprintf(
		"AND(LOWER({Guest}) = '%s', {ZIP} = '%s'%s",
		strings.ToLower(r.Name),
		r.Zip,
		rehearsalFilter,
	)
}

// // This function would be useful if changing from the current combined
// //  name & zip validation to two-step name *then* zip validation.
// func (r RSVP) validateZip(responders []guest) error {
// 	var zipMatch bool
//
// 	for _, responder := range responders {
// 		zip, err := strconv.Atoi(responder.Fields.ZIP)
// 		if err != nil {
// 			return fmt.Errorf("Failed to parse ZIP code: %w", err)
// 		} else if r.Zip == zip {
// 			zipMatch = true
// 		}
// 	}
//
// 	if !zipMatch {
// 		return fmt.Errorf("%s's invitation was send to a different ZIP code. Please try again.", r.Name)
// 	}
//
// 	return nil
// }

func parseGuests(raw string) []string {
	raw = template.HTMLEscapeString(raw)
	guests := make([]string, 0)

	for _, guest := range strings.Split(raw, ",") {
		guest = strings.TrimSpace(guest)
		guest = strings.TrimPrefix(guest, "and")
		guest = strings.TrimSpace(guest)

		guests = append(guests, strings.Title(guest))
	}

	return guests
}
