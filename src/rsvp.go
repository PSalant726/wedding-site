package main

import (
	"fmt"
	"html/template"
	"log"
	"net/mail"
	"net/url"
	"os"
	"strings"

	"github.com/fabioberger/airtable-go"
)

const shouldRetryIfRateLimited = true

var (
	airtableAPIKey = os.Getenv("AIRTABLE_API_KEY")
	airtableBaseID = os.Getenv("AIRTABLE_BASE_ID")
	airtableClient = airtable.New(airtableAPIKey, airtableBaseID, shouldRetryIfRateLimited)
)

type RSVP struct {
	Name      string
	Email     string
	Zip       string
	Guests    []string
	Message   string
	Attending bool
	Rehearsal bool
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

func NewRSVP(response url.Values, isRehearsal bool) RSVP {
	rsvp := &RSVP{
		Name:      strings.TrimSpace(strings.Title(template.HTMLEscapeString(response.Get("name")))),
		Zip:       strings.TrimSpace(template.HTMLEscapeString(response.Get("zip"))),
		Guests:    parseGuests(template.HTMLEscapeString(response.Get("guests"))),
		Message:   strings.TrimSpace(template.HTMLEscapeString(response.Get("message"))),
		Attending: response.Get("response") == "1",
		Rehearsal: isRehearsal,
	}

	email, err := mail.ParseAddress(template.HTMLEscapeString(response.Get("email")))
	if err == nil {
		rsvp.Email = email.Address
	}

	return *rsvp
}

func (r RSVP) Validate() error {
	rehearsalFilter := ")"
	if r.Rehearsal {
		rehearsalFilter = ", NOT({Rehearsal Dinner RSVP} = 'Not Invited'))"
	}

	var (
		responders = make([]guest, 0)
		listParams = airtable.ListParameters{
			Fields: []string{"Guest", "ZIP"},
			FilterByFormula: fmt.Sprintf(
				"AND(LOWER({Guest}) = '%s', {ZIP} = '%s'%s",
				strings.ToLower(r.Name),
				r.Zip,
				rehearsalFilter,
			),
			MaxRecords: 5,
		}
	)

	if err := airtableClient.ListRecords("Guests", &responders, listParams); err != nil {
		log.Printf("Failed to query Airtable for RSVP: %+v. Error: %s", r, err)
		return fmt.Errorf("An internal error occurred. Please try again.")
	}

	if len(responders) == 0 {
		return fmt.Errorf("Please verify that your name and zip code match your invitation's envelope.")
	}

	return nil
}

// // This function would be useful if changing from the current combined
// //  name & zip validation to two-step name *then* zip validation.

// func (r RSVP) validateZip(responders []guest) error {
// 	var zipMatch bool

// 	for _, responder := range responders {
// 		zip, err := strconv.Atoi(responder.Fields.ZIP)
// 		if err != nil {
// 			return fmt.Errorf("Failed to parse ZIP code: %w", err)
// 		} else if r.Zip == zip {
// 			zipMatch = true
// 		}
// 	}

// 	if !zipMatch {
// 		return fmt.Errorf("%s's invitation was send to a different ZIP code. Please try again.", r.Name)
// 	}

// 	return nil
// }

func parseGuests(raw string) []string {
	guests := make([]string, 0)

	for _, guest := range strings.Split(raw, ",") {
		guest = strings.TrimSpace(guest)
		guest = strings.TrimPrefix(guest, "and")
		guest = strings.TrimSpace(guest)

		guests = append(guests, strings.Title(guest))
	}

	return guests
}
