package zoho

import (
	"log"
	"net/http"
	"strings"
	"tutree-vapi/models"
	"tutree-vapi/service"

	"github.com/gin-gonic/gin"
)

func UpdateLeadsOnZohoCRM(leadsData models.ZohoCRMLeads) error {
	// Step 1: Generate Zoho CRM access token
	acessToken, err := service.CreateZohoCRMAccessToken()
	if err != nil {
		log.Println("[ERROR] UpdateLeadsOnZohoCRM : Failed to get Zoho access token with error :", err)
		return err
	}
	// Step 3: Convert user details to map format for Zoho CRM
	userMap := convertLeadDataStructToMap(leadsData)
	// var zohoLeadId int
	// Step 5: Update user lead data on Zoho CRM
	err = service.UpdateLeadsDataOnCRM(acessToken, leadsData.ZohoLeadId, userMap)
	if err != nil {
		log.Println("UpdateLeadsOnZohoCRM: Failed to post the user's data on zohoCRM with", err)
		return err
	}
	return nil
}

// ConvertStudentStructToMap converts a ZohoCRMLeads struct into a map[string]interface{}.
//
// This function dynamically processes the input struct and only inserts non-empty fields into the resulting map.
// Key functionalities include:
// This function ensures clean and efficient data transformation for CRM-related operations.

func convertLeadDataStructToMap(leads models.ZohoCRMLeads) map[string]interface{} {
	studentMap := make(map[string]interface{})

	// Only insert the non-empty fields into the map
	if strings.TrimSpace(leads.CallEndedReason) != "" {
		studentMap["Lead_Status"] = leads.CallEndedReason
	}

	return studentMap
}
func GetZohoLeads(c *gin.Context) {
	page := 1
	var (
		responce      []models.Lead
		err           error
		isDataPresent bool
	)
	accessToken, err := service.CreateZohoCRMAccessToken()
	if err != nil || accessToken == "" {
		return
	}
	for {
		responce, isDataPresent, err = service.GetZohoLeads(accessToken, page)
		if err != nil {
			log.Println("[ERROR] GetZohoLeads -> service.GetZohoLeads Failed to get data from zoho with error :", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Failed",
				"message": "Failed to get data from zoho leads!",
				"error":   err.Error(),
			})
			return
		}
		err = models.InsertAllLeads(responce)
		if err != nil {
			log.Println("[ERROR] GetZohoLeads -> models.InsertAllLeads Failed to save zoho leads data in database with error :", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Failed",
				"message": "Failed to save zoho leads data in table!",
				"error":   err.Error(),
			})
			return
		}
		if !isDataPresent {
			log.Println("[SUCCESS] GetZohoLeads successfully get zoho leads from zoho leads dashboard.")
			c.JSON(http.StatusOK, gin.H{
				"status":  "Success",
				"message": "Successfully get zoho leads from zoho dashboard !!!",
			})
			break
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"response": responce,
		"status":   "Success",
		"message":  "Successfully get zoho leads from zoho dashboard !!!",
	})
}
