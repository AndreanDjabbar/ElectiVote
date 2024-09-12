package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
)

func ViewHomePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewHomePage - User is not logged in",
			"Client IP", c.ClientIP(),
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)

	logger.Info(
		"ViewHomePage - rendering home page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
	context := gin.H {
		"title": "Home",
	}
	c.HTML(
		http.StatusOK,
		"home.html",
		context,
	)
}