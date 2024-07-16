package main

import "github.com/gin-gonic/gin"

func init() {

}

func main() {
	router := gin.Default()
	{
		router.LoadHTMLGlob("/internal/static/html/*.html")
	}
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}