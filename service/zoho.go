package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"tutree-vapi/models"
	"tutree-vapi/utility"
)

// CreateZohoCRMAccessToken generates an access token for Zoho CRM using the provided refresh token, client ID, and client secret.
//
// This function performs the following steps:
// 1. Retrieves the refresh token, client ID, client secret, and Zoho website URL from environment variables.
// 2. Constructs a POST request to the Zoho API to refresh the access token.
// 3. Sends the request to the Zoho server and handles the response.
// 4. Parses the response JSON to extract the access token.
// 5. Returns the generated access token if successful, or an error if any issues occur during the process.
//
// Parameters:
// - None (all necessary information is fetched from environment variables)
//
// Returns:
// - The generated access token as a string if successful, or an error if any issue occurs.
func CreateZohoCRMAccessToken() (string, error) {
	var accessToken string
	data := url.Values{}
	// fetch refresh token from env
	refreshToken := utility.GetZohoCRMRefreshToken()
	if len(refreshToken) == 0 {
		log.Println("[NOT FOUND] CreateZohoCRMAccessToken : Failed to get refresh token from ENV.")
		return accessToken, errors.New("refresh token not found")
	}
	data.Set("refresh_token", refreshToken)
	fmt.Println("REFERESH TOKEN=================", refreshToken)
	// fetch client id from env
	clientID := utility.GetZohoCRMClientID()
	if len(clientID) == 0 {
		log.Println("[NOT FOUND] CreateZohoCRMAccessToken : Failed to get client id from ENV.")
		return accessToken, errors.New("client id not found")
	}
	data.Set("client_id", clientID)
	// fetch client secret from env
	clientSecret := utility.GetZohoCRMClientSecret()
	if len(clientSecret) == 0 {
		log.Println("[NOT FOUND] CreateZohoCRMAccessToken : Failed to get client secret from ENV.")
		return accessToken, errors.New("client secret not found")
	}
	data.Set("client_secret", clientSecret)

	data.Set("grant_type", "refresh_token")
	zohoWebsite := utility.GetZohoCRMWebsiteLinkToGenerateAccessToken()
	if len(zohoWebsite) == 0 {
		log.Println("[NOT FOUND] CreateZohoCRMAccessToken : Failed to get client secret from ENV.")
		return accessToken, errors.New("zoho website not found")
	}

	body := bytes.NewBufferString(data.Encode())
	// Create the HTTP request
	req, err := http.NewRequest("POST", zohoWebsite, body)
	if err != nil {
		log.Println("[ERROR] CreateZohoCRMAccessToken : Failed to make a request to generate access token with error :", err)
		return accessToken, err
	}
	// Set the Content-Type header to x-www-form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR] CreateZohoCRMAccessToken : Failed to send request to generate access token with error :", err)
		return accessToken, err
	}
	defer resp.Body.Close()
	// Read and print the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] CreateZohoCRMAccessToken : Failed to read response to generate access token with error :", err)
		return accessToken, err
	}
	// Parse the JSON response to extract access_token
	var requestResult map[string]interface{}
	err = json.Unmarshal(respBody, &requestResult)
	if err != nil {
		fmt.Println("[ERROR] CreateZohoCRMAccessToken : Failed to parse response into json to generate access token with error :", err)
		return accessToken, err
	}
	// Extract and print the access_token
	if token, ok := requestResult["access_token"].(string); ok {
		accessToken = token
		fmt.Println("Access Token:", accessToken)
	} else {
		fmt.Println("Access Token not found in the response")
	}
	return accessToken, nil
}

type Lead struct {
	Phone      string `json:"Phone"`
	LeadSource string `json:"Lead_Source"`
	District   string `json:"district"`
}

type Response struct {
	Data []Lead `json:"data"`
}

func GetPhoneOfColdLeads(accessToken string) ([]Lead, error) {
	// ZOHO API URL
	url := "https://www.zohoapis.in/crm/v2/Leads?fields=Phone,Lead_Source,district"

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []Lead{}, err
	}

	// Set headers
	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return []Lead{}, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return []Lead{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return []Lead{}, err
	}
	leads := []Lead{}
	// Extract and print phone numbers where Lead_Source is "Cold_data"
	for _, lead := range response.Data {
		if lead.LeadSource == "Partner" {
			lead.Phone = fmt.Sprintf("+1%s", lead.Phone)
			leads = append(leads, lead)
		}
	}
	log.Println("Leads --------------->", leads)
	return leads, nil
}
func UpdateLeadsDataOnCRM(accessToken string, leadID string, data map[string]interface{}) error {
	// Prepare payload
	payload := map[string]interface{}{
		"data": []map[string]interface{}{data},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("https://www.zohoapis.in/crm/v2/Leads/%s", leadID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the headers
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	// responseBody, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalf("Error reading response: %v", err)
	// }

	// Print the response
	fmt.Printf("Response Status: %s\n", resp.Status)
	//fmt.Printf("Response Body: %s\n", string(responseBody))
	return nil
}

func GetZohoLeads(accessToken string, page int) ([]models.Lead, bool, error) {
	apiURL := "https://www.zohoapis.in/crm/v2/Leads/search"
	userID := os.Getenv("ZOHO_USER_ID")
	// filter out the data for cold leads
	criteria := fmt.Sprintf("((Created_By.id:equals:%s) and (Lead_Source:equals:Cold_data))", userID)
	encodedCriteria := url.QueryEscape(criteria)
	// Max records per page
	perPage := 200
	// NOTE : zoho allow only 200 data per page per api hit
	queryParams := fmt.Sprintf("?criteria=%s&fields=id,district,created_at,gender,Full_Name,Lead_Status,Phone,query_param,session_type&per_page=%d&page=%d", encodedCriteria, perPage, page)
	fullURL := apiURL + queryParams

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, true, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, true, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, true, fmt.Errorf("failed to fetch data: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, fmt.Errorf("error reading response body: %v", err)
	}

	var response struct {
		Data []map[string]interface{} `json:"data"`
		Info struct {
			MoreRecords bool `json:"more_records"`
		} `json:"info"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, true, fmt.Errorf("error parsing JSON response: %v", err)
	}
	// Process the current batch of leads
	leads := ProcessLeadsBatch(response.Data)

	return leads, response.Info.MoreRecords, nil
}
func ProcessLeadsBatch(rawLeads []map[string]interface{}) []models.Lead {
	var leads []models.Lead

	for _, rawLead := range rawLeads {
		lead := models.Lead{
			ID:       getString(rawLead["id"]),
			FullName: getString(rawLead["Full_Name"]),
			Phone:    getString(rawLead["Phone"]),
			District: getString(rawLead["district"]),
			Gender:   getString(rawLead["gender"]),
		}
		leads = append(leads, lead)
	}
	return leads
}
func getString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}
