package handlers

import (
	"net/http"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/middlewares"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/models"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/repositories"
	"github.com/AndreanDjabbar/CaysAPIHub/internal/utils"
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

	var userProfile models.Profile
	var err error
	username := middlewares.GetUserData(c)

	userProfile, err = repositories.GetProfilesByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	}

	userEmail, err := repositories.GetUserEmailByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	
	}
	context := gin.H{
		"title": "Profile",
		"username": username,
		"userProfile": userProfile,
		"userEmail": userEmail,
	}
	c.HTML(
		http.StatusOK,
		"profile.html",
		context,
	)
}

func ViewEditProfilePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	
	var userProfile models.Profile
	var err error
	username := middlewares.GetUserData(c)

	userProfile, err = repositories.GetProfilesByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	}

	userEmail, err := repositories.GetUserEmailByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/profile-page/",
		)
	
	}
	context := gin.H{
		"title": "Profile",
		"username": username,
		"userProfile": userProfile,
		"userEmail": userEmail,
	}
	c.HTML(
		http.StatusOK,
		"editProfile.html",
		context,
	)
}