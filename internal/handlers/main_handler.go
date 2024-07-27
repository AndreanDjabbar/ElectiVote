package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
)

func ViewHomePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	dataUser := middlewares.GetUserData(c)

	context := gin.H {
		"title": "Home",
		"dataUser": dataUser,
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}