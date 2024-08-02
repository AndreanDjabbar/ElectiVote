package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/models"
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
		"title": "Create Vote",
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
	voteTitleErr := ""
	username := middlewares.GetUserData(c)
	voteTitle := c.PostForm("voteTitle")
	voteDesc := c.PostForm("voteDesc")
	voteCode := utils.GenerateVoteCode()
	start := models.CustomTime{Time: time.Now()}

	moderatorID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/create-vote-page/",
		)
	}

	if len(voteTitle) < 5 {
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.StartVoteFactory(voteTitle, voteDesc, voteCode, uint(moderatorID), start)
	
		_, err = repositories.CreateVote(newVote)
		if err != nil {
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/create-vote-page/",
			)
		}
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	context := gin.H {
		"title": "Create Vote",
		"voteTitleErr": voteTitleErr,
		"voteTitle": voteTitle,
		"voteDesc": voteDesc,
	}
	c.HTML(
		http.StatusOK,
		"createVote.html",
		context,
	)
}

func ViewManageVotesPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	username := middlewares.GetUserData(c)
	votesData, err := repositories.GetVotesDataByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}
	context := gin.H {
		"title": "Manage Votes",
		"votes": votesData,
	}
	c.HTML(
		http.StatusOK,
		"manageVotes.html",
		context,
	)
}

func ViewManageVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidModerator(username, uint(voteID)) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}
	context := gin.H {
		"title": "Manage Vote",
		"voteData": voteData,
	}
	c.HTML(
		http.StatusOK,
		"manageVote.html",
		context,
	)
}

func ManageVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidModerator(username, uint(voteID)) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	voteTitleErr := ""
	voteTitle := c.PostForm("voteTitle")
	voteDesc := c.PostForm("voteDesc")

	if len(voteTitle) < 5 {
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.UpdateVoteFactory(voteTitle, voteDesc)
		_, err := repositories.UpdateVote(uint(voteID), newVote)
		if err != nil {
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/manage-vote-page/",
			)
		}
		c.Redirect(
			http.StatusFound,
			"/electivote/manage-vote-page/",
		)
		return
	}
	context := gin.H {
		"title": "Manage Vote",
		"voteData": voteData,
		"voteTitleErr": voteTitleErr,
		"voteTitle": voteTitle,
		"voteDesc": voteDesc,
	}
	c.HTML(
		http.StatusOK,
		"manageVote.html",
		context,
	)
}

func ViewDeleteVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidModerator(username, uint(voteID)) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	context := gin.H {
		"title": "Delete Vote",
		"voteData": voteData,
	}
	c.HTML(
		http.StatusOK,
		"deleteVote.html",
		context,
	)
}

func DeleteVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidModerator(username, uint(voteID)) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	err := repositories.DeleteVote(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/",
	)
}