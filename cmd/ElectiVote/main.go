package main

import (
	"fmt"
	"os"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/db"
	"github.com/AndreanDjabbar/ElectiVote/internal/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)
func init() {
	db.ConnectToDatabase()
}

func main() {
	logger := config.SetUpLogger()
	logger.Info("Start setting up server")
	router := gin.Default()
	router.LoadHTMLGlob("internal/views/html/*.html")
	router.Static("/images", "internal/assets/images")
	router.MaxMultipartMemory = 8 << 20

	mainStore := cookie.NewStore([]byte("main"))
	mainStore.Options(sessions.Options{
		MaxAge: 0,
		HttpOnly: true,
		Path: "/",
		Secure: false,
	})

	router.Use(sessions.Sessions("mainSession", mainStore))
	port := os.Getenv("PORT")
	if port == "" {
		logger.Warn("PORT is not set, using default port 8080")
		port = "8080" 
	}
	host := os.Getenv("HOST")
	if host == "" {
		logger.Warn("HOST is not set, using default host localhost")
		host = "localhost"
	}
	routes.SetUpRoutes(router)
	logger.Info("Server is run on", "host", host, "port", port)
	err := router.Run(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Server failed to start", "error", err)
		panic(err)
	}
}