package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
)

func ViewHomePage(c *gin.Context) {
	logger.Info(
		"ViewHomePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewHomePage - User Not Logged In",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	dataUser := middlewares.GetUserData(c)

	logger.Info(
		"ViewHomePage - User Logged In",
	)
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