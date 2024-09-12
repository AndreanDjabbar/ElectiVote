package config

import (
	"os"
	"github.com/gin-gonic/gin"
)

func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		return "localhost"
	}
	return host
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("internal/views/html/*.html")
	router.Static("/images", "internal/assets/images")
	router.MaxMultipartMemory = 8 << 20
	return router
}