package handlers

import (
	"net/http"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/gin-gonic/gin"
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

func ViewProfilePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	
	dataUser := middlewares.GetUserData(c)

	context := gin.H {
		"title": "Profile",
		"dataUser": dataUser,
	}
	c.HTML(
		http.StatusOK,
		"profile.html",
		context,
	)
}