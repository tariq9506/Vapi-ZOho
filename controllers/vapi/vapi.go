package vapicontroller

import (
	"io"
	"log"
	"net/http"
	"tutree-vapi/service"

	"github.com/gin-gonic/gin"
)

func GetCallDetails(c *gin.Context) {
	callerID := c.Query("caller-id")
	if len(callerID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone_number is required"})
		return
	}

	userResponseInfo, err := service.GetCallDetails(callerID)
	if err != nil {
		log.Println("Error getting call details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	log.Println("GetCallDetails LOG", userResponseInfo)
	phone := userResponseInfo["phoneNumber"]
	log.Println("GetCallDetails LOG", phone)
	if userResponseInfo["offer_response"] == "Interested" {

		err = service.SendMessage(phone, "Hello, this is a test message")
		if err != nil {
			log.Println("Error sending message:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": userResponseInfo})
}
func HandleTheResponnceOfVapi(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
	}
	log.Println("HandleTheResponnceOfVapi LOG", string(bodyBytes))
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
func RouteHit(c *gin.Context) {
	log.Println("RouteHit")
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
