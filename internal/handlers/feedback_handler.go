package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-gonic/gin"
)

func ViewGiveFeedbackPage(c *gin.Context) {
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
		"ViewGiveFeedbackPage - rendering give feedback page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	context := gin.H{
		"title":   "Feedback",
		"siteKey": siteKey,
	}
	c.HTML(
		http.StatusOK,
		"feedback.html",
		context,
	)
}

func GiveFeedbackPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"GiveFeedbackPage - User is not logged in",
			"Client IP", c.ClientIP(),
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	siteKey := os.Getenv("RECAPTCHA_SITE_KEY")
	username := middlewares.GetUserData(c)
	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"GiveFeedbackPage - failed to get user id",
			"Username", username,
			"Client IP", c.ClientIP(),
			"error", err,
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	feedbackMessage := c.PostForm("feedback")
	feedbackRateStr := c.PostForm("rating")
	feedbackRate, err := strconv.Atoi(feedbackRateStr)
	if err != nil {
		logger.Error(
			"GiveFeedbackPage - failed to convert feedback rate to int",
			"Username", username,
			"Client IP", c.ClientIP(),
		)
		utils.RenderError(
			c,
			http.StatusBadRequest,
			"Failed to convert feedback rate to int",
			"home-page",
		)
		return
	}
	feedbackMessageErr, feedbackRateErr, captchaErr := utils.ValidateFeedbackInput(feedbackMessage, uint(feedbackRate), c)
	if feedbackMessageErr != "" || feedbackRateErr != "" || captchaErr != "" {
		logger.Warn(
			"GiveFeedbackPage - invalid input",
			"Username", username,
			"Client IP", c.ClientIP(),
			"feedbackMessageErr", feedbackMessageErr,
			"feedbackRateErr", feedbackRateErr,
			"captchaErr", captchaErr,
		)
		context := gin.H{
			"title":              "Feedback",
			"siteKey":            siteKey,
			"feedbackMessage":    feedbackMessage,
			"feedbackRate":       feedbackRateStr,
			"feedbackMessageErr": feedbackMessageErr,
			"feedbackRateErr":    feedbackRateErr,
			"captchaErr":         captchaErr,
		}
		c.HTML(
			http.StatusBadRequest,
			"feedback.html",
			context,
		)
		return
	}
	timeCreated := models.CustomTime{
		Time: time.Now(),
	}
	feedback := factories.FeedbackFactory(uint(userID), feedbackMessage, uint(feedbackRate), timeCreated)
	_, err = repositories.CreateFeedback(feedback)
	if err != nil {
		logger.Error(
			"GiveFeedbackPage - failed to create feedback",
			"Username", username,
			"Client IP", c.ClientIP(),
			"error", err,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			"Failed to create feedback",
			"home-page",
		)
		return
	}
	logger.Info(
		"GiveFeedbackPage - feedback created",
		"Username", username,
		"Client IP", c.ClientIP(),
		"action", "redirecting to about us page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/about-us-page/",
	)
}