package main

import (
	"fmt"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/routes"
	"github.com/gin-contrib/sessions"
)
func init() {
	db.ConnectToDatabase()
}

func main() {
	logger := config.SetUpLogger()
	logger.Info("Start setting up server")

	router := config.SetUpRouter()
	router.Use(sessions.Sessions("mainSession", config.SetUpSessionStore()))
	routes.SetUpRoutes(router)

	host := config.GetHost()
	port := config.GetPort()

	logger.Info("Server is run on", "host", host, "port", port)
	err := router.Run(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Server failed to start", "error", err.Error())
		return 
	}
}