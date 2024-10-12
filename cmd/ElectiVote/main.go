package main

import (
	"fmt"
	"os"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/routes"
	"github.com/gin-contrib/sessions"
	"github.com/joho/godotenv" // Pastikan package ini terinstal
)

func init() {
	db.ConnectToDatabase()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	logger := config.SetUpLogger()
	logger.Info("Start setting up server")

	router := config.SetUpRouter()
	router.Use(sessions.Sessions("mainSession", config.SetUpSessionStore()))
	routes.SetUpRoutes(router)

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server is run on", "host", host, "port", port)
	err = router.Run(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Server failed to start", "error", err)
		panic(err)
	}
}
