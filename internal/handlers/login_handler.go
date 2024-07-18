package handlers

import (
	"net/http"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/repositories"
	"github.com/gin-gonic/gin"
)

func ViewLoginPage(c *gin.Context) {
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

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
	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.PostForm("remember")

	usernameErr, passwordErr := "", ""

	if username == "" {
		usernameErr = "Username Must be Filled"
	}

	if password == "" {
		passwordErr = "Password Must be Filled"
	}

	if len(username) != 0 && (len(username) < 5 || len(username) > 255) {
		usernameErr = "Username must be between 5 and 255 characters"
	}

	if len(password) != 0 && (len(password) < 5 || len(password) > 255) {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	_, err := repositories.GetUserByUsername(username)
	if err != nil && username != "" && (len(username) >= 5 && len(username) <= 255) {
		usernameErr = "Invalid Username"
	}

	_, err = repositories.CheckPasswordByUSername(username, password)
	if err != nil && password != "" && (len(password) >= 5 && len(password) <= 255) {
		passwordErr = "Invalid Password"
	}
	
	if usernameErr == "" && passwordErr == "" {
		if remember == "on" {
			middlewares.SetCookies(c, username)
		} else {
			middlewares.SetSession(c, username)
		}
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Login",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"username": username,
		"password": password,
	}
	c.HTML(
		http.StatusOK,
		"login.html",
		context,
	)
}

func LogoutPage(c *gin.Context) {
	middlewares.DeleteSession(c)
	middlewares.DeleteCookie(c)
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}