package handlers

import (
	"net/http"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/factories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
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
	password2 := c.PostForm("password2")
	email := c.PostForm("email")

	usernameErr, passwordErr, password2Err, emailErr := utils.ValidateRegisterInput(username, password, password2, email)

	if usernameErr == "" && passwordErr == "" && password2Err == "" && emailErr == "" {
		hashedPassword, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			utils.RenderError(c, http.StatusInternalServerError, hashErr.Error(), "/electivote/register-page/")
			return
		}

		newUser := factories.CreateUser(username, hashedPassword, email, "user")

		_, err := repositories.RegisterUser(newUser)
		if err != nil {
			utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
			return
		}

		userID, err := repositories.GetUserIdByUsername(username)
		if err != nil {
			utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
			return
		}

		newProfile := factories.CreateFirstProfile(userID)

		_, err = repositories.CreateProfile(newProfile)
		if err != nil {
			utils.RenderError(c, http.StatusInternalServerError, err.Error(), "/electivote/register-page/")
			return
		}

		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	context := gin.H{
		"title":       "Register",
		"usernameErr": usernameErr,
		"passwordErr": passwordErr,
		"password2Err": password2Err,
		"emailErr":    emailErr,
		"username":    username,
		"password":    password,
		"password2":   password2,
		"email":       email,
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}
