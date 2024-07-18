package main

import (
	"github.com/AndreanDjabbar/CaysAPIHub/internal/db"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)
func init() {
	db.ConnectToDatabase()
}

var SessionKey string = "session"

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	{
		router.LoadHTMLGlob("../../internal/views/html/*.html")
		router.Use(sessions.Sessions(SessionKey, store))
		store.Options(sessions.Options{
			MaxAge: 0,
			HttpOnly: true,
			Path: "/",
			Secure: false,
		})
	}
	
	routes.SetUpRoutes(router)
	err := router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}