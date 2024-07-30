package handlers

import (
	"net/http"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-gonic/gin"
)

func ViewCreateVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	context := gin.H {
		"title": "Vote",
	}
	c.HTML(
		http.StatusOK,
		"createVote.html",
		context,
	)
}

func CreateVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	username := middlewares.GetUserData(c)
	voteTitle := c.PostForm("voteTitle")
	voteDesc := c.PostForm("voteDesc")
	voteCode := utils.GenerateVoteCode()
	start := time.Now()
	formattedStart := start.Format("2006-01-02 15:04:05")
	parsedTime, err := time.Parse("2006-01-02 15:04:05", formattedStart)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/create-vote-page/",
		)
	}
	moderatorID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/create-vote-page/",
		)
	}

	newVote := factories.StartVoteFactory(voteTitle, voteDesc, voteCode, uint(moderatorID), parsedTime)

	_, err = repositories.CreateVote(newVote)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/create-vote-page/",
		)
	}
}