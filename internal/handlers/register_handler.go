package handlers

import (
	"net/http"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/repositories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/utils"
	"github.com/gin-gonic/gin"
)

func ViewRegisterPage(c *gin.Context) {

	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	context := gin.H {
		"title": "Register",
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}

func RegisterPage(c *gin.Context) {

	if middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	usernameErr, passwordErr, emailErr := "", "", ""

	if username == "" {
		usernameErr = "Username must be filled"
	}

	if password == "" {
		passwordErr = "Password must be filled"
	}

	if email == "" {
		emailErr = "Email must be filled"
	}

	if len(username) != 0 && (len(username) < 5 || len(username) > 255) {
		usernameErr = "Username must be between 5 and 255 characters"
	}

	if len(password) != 0 && (len(password) < 5 || len(password) > 255) {
		passwordErr = "Password must be between 5 and 255 characters"
	}

	if email != "" && !utils.IsValidEmail(email) {
		emailErr = "Email must contain @ and end with .com or .co.id"
	}

	if usernameErr == "" && passwordErr == "" && emailErr == "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			context := gin.H {
				"title": "Error",
				"error": err.Error(),
				"source": "/electivote/register-page/",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
		}
		newUser := models.User{
			Username: username,
			Password: hashedPassword,
			Email: email,
			Role: "user",
		}
		_, err = repositories.RegisterUser(newUser)
		if err != nil {
			context := gin.H {
				"title": "Error",
				"error": err.Error(),
				"source": "/electivote/register-page/",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
		}
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	context := gin.H {
		"title": "Register",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"emailErr": emailErr,
		"username": username,
		"password": password,
		"email": email,
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}