package utility

import "os"

// GetZohoCRMRefreshToken retrieves the Zoho CRM refresh token from the environment variables.
// This token is used to generate a new access token for Zoho CRM API calls.
func GetZohoCRMRefreshToken() string {
	return os.Getenv("ZOHO_CRM_REFRESH_TOKEN")
}

// GetZohoCRMClientID retrieves the Zoho CRM Client ID from the environment variables.
// The Client ID is required for authentication with Zoho's OAuth2 API.
func GetZohoCRMClientID() string {
	return os.Getenv("ZOHO_CLIENT_ID")
}

// GetZohoCRMClientSecret retrieves the Zoho CRM Client Secret from the environment variables.
// The Client Secret, along with the Client ID, is used to authenticate with Zoho's OAuth2 API.
func GetZohoCRMClientSecret() string {
	return os.Getenv("ZOHO_CLIENT_SECRET")
}

// GetZohoCRMWebsiteLinkToGenerateAccessToken retrieves the website link to generate the Zoho CRM access token.
// This link is used during the OAuth2 flow to obtain a new access token.
func GetZohoCRMWebsiteLinkToGenerateAccessToken() string {
	return os.Getenv("ZOHO_WEBSITE_LINK_TO_GENERATE_ACCESS_TOKEN")
}

// GetZohoApiToCreateLeadOnZohoCRMLeadBoard retrieves the API endpoint URL for creating leads on Zoho CRM.
// This endpoint allows posting new lead data to Zoho's CRM system.
func GetZohoApiToCreateLeadOnZohoCRMLeadBoard() string {
	return os.Getenv("ZOHO_CRM_WEBSITE_TO_SEND_LEADS")
}

// GetZohoApiToUpdateLeadOnZohoCRMLeadBoard retrieves the API endpoint URL for updating leads on Zoho CRM.
// This endpoint allows updating existing lead data on Zoho's CRM system.
func GetZohoApiToUpdateLeadOnZohoCRMLeadBoard() string {
	return os.Getenv("ZOHO_CRM_WEBSITE_TO_UPDATE_LEADS")
}
