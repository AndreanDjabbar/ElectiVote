package routes

import (
	"net/http"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)	
}

func MainRootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}

func SetUpRoutes(router *gin.Engine) {
	mainRouter := router.Group("/electivote")
	{
		router.GET("/", RootHandler)
		mainRouter.GET("/", MainRootHandler)
	}
	{
		mainRouter.GET("login-page/", handlers.ViewLoginPage)
		mainRouter.POST("login-page/", handlers.LoginPage)
	}
}