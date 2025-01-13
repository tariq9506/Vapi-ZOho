package vapicontroller

import (
	"log"
	"net/http"
	callschedule "tutree-vapi/controllers/call_schedule"
	"tutree-vapi/models"

	"github.com/gin-gonic/gin"
)

func GetPhoneOfColdLeads(c *gin.Context) {
	leads, err := models.GetZoheColdLeads()
	if err != nil {
		log.Println("Error getting cold leads:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Step 2: Iterate over each lead and initiate the call
	var statusCode int

	// Make the call to the current lead's phone number
	err = callschedule.ScheduleCalls(leads)
	if err != nil {
		log.Println("Error scheduling calls:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"leads":      leads,
		"statusCode": statusCode,
	})
}
