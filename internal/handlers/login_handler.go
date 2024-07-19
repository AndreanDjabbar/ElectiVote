package handlers

import (
	"net/http"
	"sync"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/repositories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/utils"
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

	usernameErr, passwordErr := utils.ValidateLoginInput(username, password)
	var usernameCheckErr, passwordCheckErr error
	var wg sync.WaitGroup

	if usernameErr == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repositories.GetUserByUsername(username); err != nil {
				usernameCheckErr = err
			}
		}() 
	}

	if passwordErr == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := repositories.CheckPasswordByUSername(username, password); err != nil {
				passwordCheckErr = err
			}
		}()
	
	}
	wg.Wait()

	if usernameCheckErr != nil {
		usernameErr = "Username not found"
	}

	if passwordCheckErr != nil {
		passwordErr = "Password is incorrect"
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