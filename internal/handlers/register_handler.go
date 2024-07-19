package handlers

import (
	"net/http"
	"sync"

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
	password2 := c.PostForm("password2")
	email := c.PostForm("email")

	usernameErr, passwordErr, password2Err, emailErr := utils.ValidateRegisterInput(username, password, password2, email)

	if usernameErr == "" && passwordErr == "" && password2Err == "" && emailErr == "" {
		var wg sync.WaitGroup
		var hashedPassword string
		var hashErr, err error
		var regMu sync.Mutex

		hashedPassword, hashErr = utils.HashPassword(password)
		
		if hashErr != nil {
			context := gin.H {
				"title": "Error",
				"error": hashErr.Error(),
				"source": "/electivote/register-page/",
			}
			c.HTML(
				http.StatusInternalServerError,
				"error.html",
				context,
			)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			regMu.Lock()
			newUser := models.User{
				Username: username,
				Password: hashedPassword,
				Email: email,
				Role: "user",
			}
			_, err = repositories.RegisterUser(newUser)
			regMu.Unlock()
		}()
		wg.Wait()

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
		"password2Err": password2Err,
		"emailErr": emailErr,
		"username": username,
		"password": password,
		"password2": password2,
		"email": email,
	}
	c.HTML(
		http.StatusOK,
		"register.html",
		context,
	)
}