package main

import (
	"fmt"
	"os"

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
	host := os.Getenv("HOST")
	routes.SetUpRoutes(router)
	err := router.Run(fmt.Sprintf("%s:%s",host, port))
	if err != nil {
		panic(err)
	}
}