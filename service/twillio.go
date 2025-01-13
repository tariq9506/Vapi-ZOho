package service

import (
	"log"
	"os"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendMessage sends an SMS message to a specified phone number using the Twilio API.
// Parameters:
// - phone: The recipient's phone number (without the country code).
// - message: The content of the SMS message to be sent.
// - dialingCode: The international dialing code for the recipient's country (e.g., "+1" for the US).
// Returns:
// - error: Returns an error if the message sending fails, otherwise returns nil.
func SendMessage(phone string, message string) error {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_ACCOUNT_AUTH_TOKEN")
	twilioFrom := os.Getenv("TWILIO_FROM_NUMBER")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}

	params.SetTo(phone)
	params.SetFrom(twilioFrom)
	params.SetBody(message)
	log.Println("SendMessage : sending message to the tutree user with phone number:", phone)

	_, err := client.Api.CreateMessage(params)

	if err != nil {
		log.Println("SendMessage : failed while sending message to the tutree user, with error:", err)
		return err
	}

	return nil
}
