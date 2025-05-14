package main

import (
	"introme-api/pkg/admin-service/entities"
	"introme-api/pkg/authentication"
	openaiservice "introme-api/pkg/openai-service"
	"introme-api/pkg/shared/helper"
	"introme-api/server"
	"log"

	"github.com/joho/godotenv"
)

var OrgID = []string{"amsort"}

func main() {

	// Load environment variables from the .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Server initialization
	app := server.Create()

	// By Default try to connect shared db
	// database.Init()
	helper.InitSuggestionCache()
	// Set up authentication routes for routes that do not require a token.
	authentication.SetupRoutes(app)

	// Set up all routes for the application.
	entities.SetupAllRoutes(app)
	openaiservice.SetupAiRoutes(app)
	// Initialize custom validators for data validation.gvg///]oo00-=}}]]]
	helper.InitCustomValidator() //testing

	// Create a context for the background service
	// _, cancel := context.WithCancel(context.Background())
	// var wg sync.WaitGroup

	// Start the background service with a 2-minute delay
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	time.Sleep(10 * time.Second) // Wait for 2 minutes
	//server start after 2 mins removed the data from temp user collection
	// 	helper.ServerInitstruct(OrgID)
	// 	helper.DeleteTempDocumentsBasedexpireTime()

	// Once the server is loaded, create and load structs based on OrgId and add them to the TypeMap.
	// }()

	// Start the server immediately
	// go func() {
	go func() {
		helper.ServerInitstruct(OrgID)
	}()
	if err := server.Listen(app); err != nil {
		log.Panic(err)
	}
	// }()

	// Wait for a signal to gracefully shut down the background service
	// time.Sleep(time.Hour * 24)

	// Signal the background service to shut down gracefully
	// cancel()

	// Wait for the background service to finish
	// wg.Wait()

}
