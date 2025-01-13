package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type CallRequest struct {
	To           string `json:"phoneNumber"`
	Text         string `json:"text"`
	From         string `json:"from"`
	AssistanceID string `json:"your_assistant_id"`
}

// MakeACallOfZohoColdLeads makes a call to a phone number using the vapi.ai API.
// Parameters:
// - phoneNumber: The phone number to call.
// Returns:
// - int: The status code of the API request.
// - error: Any error encountered during the process.

func MakeACallOfZohoColdLeads(phoneNumber, districtName string) (int, string, error) {
	// vapi.ai API URL
	url := os.Getenv("VAPI_AI_CALL_PHONE_URL")
	method := "POST"
	// Payload for the API request
	// Corrected the payload to match the API documentation
	assistantID := os.Getenv("YOUR_ASSISTANT_ID")
	if len(assistantID) == 0 {
		log.Println("Assistant ID not found in the environment variables")
		return 0, "", fmt.Errorf("assistant ID not found in the environment variables")
	}
	phoneNumberID := os.Getenv("YOUR_PHONE_NUMBER_ID")
	if len(phoneNumberID) == 0 {
		log.Println("Phone number ID not found in the environment variables")

		return 0, "", fmt.Errorf("phone number ID not found in the environment variables")
	}
	payload := fmt.Sprintf(`{
	"assistantId": "%s",
	"assistantOverrides": {
		"variableValues": {
			"school_district_name": "%s"
			
		}
	},
	"customer": {
		"number": "%s"
	},
	"phoneNumberId": "%s"
}
`, assistantID, districtName, phoneNumber, phoneNumberID)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	apikey := os.Getenv("VAPI_AI_API_KEY")
	if len(apikey) == 0 {
		log.Println("API key not found in the environment variables")
		return 0, "", fmt.Errorf("API key not found in the environment variables")
	}
	// Add headers
	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, "", err
	}
	defer resp.Body.Close()
	// Check the response status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success: Log and return the status code
		fmt.Println("Call initiated successfully!")
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Client error (e.g., bad request)
		fmt.Printf("Client error: %d - Please check your request.\n", resp.StatusCode)
	} else if resp.StatusCode >= 500 && resp.StatusCode < 600 {
		// Server error (e.g., API issue)
		fmt.Printf("Server error: %d - There was an issue with the API server.\n", resp.StatusCode)
	} else {
		// Other status codes
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

	// Read and parse the response body for more info (optional)
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		fmt.Println("Error decoding response:", err)
	} else {
		//Log or process the responseBody for further insights
		// fmt.Printf("Response body: %+v\n", responseBody)
	}

	// Step 2: Fetch the phoneCallProviderId directly from the decoded map
	callerID, ok := responseBody["id"].(string)
	if !ok {
		log.Printf("Error: phoneCallProviderId not found or not a string in response")
	}
	return resp.StatusCode, callerID, nil
}
func GetCallDetails(callID string) (map[string]string, error) {
	// Replace with the actual call ID
	apiKey := os.Getenv("VAPI_AI_API_KEY") // Replace with your actual API key
	vapiURL := os.Getenv("VAPI_AI_CALL_DETAILS_URL")

	// Construct the URL to fetch the call details
	url := fmt.Sprintf("%s%s", vapiURL, callID)
	log.Println("URL:", url)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers for the request
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-OK status codes
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: received status code", resp.StatusCode)
		return nil, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	// Parse the JSON response body into a map
	var responseBody map[string]interface{}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Println("Error decoding JSON:", err)
		return nil, err
	}

	// Log the response body for debugging
	//log.Printf("Response body: %+v\n", responseBody)
	userResponseInfo := make(map[string]string)
	// Fetch the call status and ended reason, checking if they exist
	status, ok := responseBody["status"].(string)
	if !ok {
		log.Println("Error: Unable to fetch 'status' from the response")
		// return nil, fmt.Errorf("unable to fetch 'status' from the response")
	}
	userResponseInfo["status"] = status

	endedReason, ok := responseBody["endedReason"].(string)
	if !ok {
		log.Println("Error: Unable to fetch 'endedReason' from the response")
		// return nil, fmt.Errorf("unable to fetch 'endedReason' from the response")
	}
	userResponseInfo["endedReason"] = endedReason
	// Navigate to analysis -> structuredData -> offer_response
	analysis, ok := responseBody["analysis"].(map[string]interface{})
	if !ok {
		log.Println("Error: Unable to fetch 'analysis' from the response")
		// return userResponseInfo, nil
	}
	// Fetch the phone number
	customer, ok := responseBody["customer"].(map[string]interface{})
	if !ok {
		log.Println("Error: Unable to fetch 'customer' from the response")
		// return userResponseInfo, nil
	}

	phoneNumber, ok := customer["number"].(string)
	if !ok {
		log.Println("Error: Unable to fetch 'number' from the 'customer' field")
		// return userResponseInfo, nil
	}
	userResponseInfo["phoneNumber"] = phoneNumber

	structuredData, ok := analysis["structuredData"].(map[string]interface{})
	if !ok {
		log.Println("Error: Unable to fetch 'structuredData' from the response")
		// return userResponseInfo, nil
	}

	offerResponse, ok := structuredData["offer_response"].(string)
	if !ok {
		log.Println("Error: Unable to fetch 'offer_response' from the response")
		// return userResponseInfo, nil
	}
	userResponseInfo["offer_response"] = offerResponse
	// Return the status and endedReason
	return userResponseInfo, nil
}
