package controllers

import (
	"fmt"
	"log"

	"github.com/ttacon/libphonenumber"
)

func RemoveDialingCode(number string) (string, error) {
	parsedNumber, err := libphonenumber.Parse(number, "US")
	if err != nil {
		log.Println("parseTwilioSmsDetails: [ERROR] Failed to parse phone number:", err)
		return "", err
	}

	// Extract country code and national number from the parsed phone number
	//countryCodeInt := parsedNumber.GetCountryCode()
	phoneInt := parsedNumber.GetNationalNumber()

	// Convert the national number to string format
	phoneNumber := fmt.Sprintf("%v", phoneInt)
	return phoneNumber, nil
}
