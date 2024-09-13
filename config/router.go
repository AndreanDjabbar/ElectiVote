package config

import (
	"os"
	"text/template"
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
	router.SetFuncMap(template.FuncMap{
		"AddOne": func(i int) int {
			return i + 1
		},
	})
	router.LoadHTMLGlob("internal/views/html/*.html")
	router.Static("/images", "internal/assets/images")
	router.MaxMultipartMemory = 8 << 20
	return router
}