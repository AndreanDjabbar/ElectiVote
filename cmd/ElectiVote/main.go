package main

import (
	"fmt"
	"os"

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
		router.LoadHTMLGlob("internal/views/html/*.html")
		router.Static("/images", "internal/assets/images")
		router.MaxMultipartMemory = 8 << 20
		router.Use(sessions.Sessions(SessionKey, store))
		store.Options(sessions.Options{
			MaxAge: 0,
			HttpOnly: true,
			Path: "/",
			Secure: false,
		})
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	routes.SetUpRoutes(router)
	err := router.Run(fmt.Sprintf("%s:%s",host, port))
	if err != nil {
		panic(err)
	}
}