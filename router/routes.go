package route

import (
	"net/http"
	vapicontroller "tutree-vapi/controllers/vapi"
	"tutree-vapi/controllers/zoho"
	"tutree-vapi/webhook"

	"github.com/gin-gonic/gin"
)

// AddRoutes is responsible for adding all the routes so the server can handle
// new routes. this means that we can reuse this function for multiple prefixes.
// prefixes like job_portal are necessary for legacy url handling.
func AddRoutes(router *gin.RouterGroup) {

	// NOTE :- all api must be in this group and for every particular feature apis must be create new group.
	router.GET("/", vapicontroller.RouteHit)
	router.GET("/zoho-leads", vapicontroller.GetPhoneOfColdLeads)
	router.GET("/call-details", vapicontroller.GetCallDetails)
	router.POST("/server-messages", webhook.WebhookHandler)
	router.OPTIONS("/server-messages", optionsHandler)
	router.POST("/websocket/send-link", webhook.HandleToSendLinkViaSMS)
	router.OPTIONS("/websocket/send-link", optionsHandler)
	router.GET("/leads", zoho.GetZohoLeads)
}

// SetupRouter sets up routes
func SetupRouter() *gin.Engine {

	router := gin.Default()
	// Add all current URls
	AddRoutes(&router.RouterGroup)
	router.Use(corsMiddleware())
	return router
}

// corsMiddleware sets the necessary headers for CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization,token")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// optionsHandler handles preflight OPTIONS requests
func optionsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
	c.Status(http.StatusOK)
}
