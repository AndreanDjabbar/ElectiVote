package handlers

import (
	"net/http"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func ViewLoginPage(c *gin.Context) {
	context := gin.H {
		"title": "Login",
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LoginPage(c *gin.Context) {
	username := c.PostForm("username")
	remember := c.PostForm("remember")
	if remember == "on" {
		middlewares.SetCookies(c, username)
	} else {
		middlewares.SetSession(c, username)
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		nil,
	)
}