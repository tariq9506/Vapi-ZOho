package main

import (
	"log"
	"os"
	route "tutree-vapi/router"

	"github.com/gin-contrib/pprof"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {

		log.Fatal("Error loading .env file -> ", err)
	}

	//setup routes
	r := route.SetupRouter()
	pprof.Register(r)
	// running
	r.Run(":" + os.Getenv("PORT"))
}
