package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
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
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	var wg sync.WaitGroup
	var voteData models.Vote
	var candidates []models.Candidate
	var voteDataErr error
	var candidatesErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		voteData, voteDataErr = repositories.GetVoteDataByVoteID(uint(voteID))
	}()

	go func() {
		defer wg.Done()
		candidates, candidatesErr = repositories.GetCandidatesByVoteID(uint(voteID))
	}()

	wg.Wait()

	if voteDataErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			voteDataErr.Error(),
			"/electivote/manage-vote-page/",
		)
		return
	}

	if candidatesErr != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			candidatesErr.Error(),
			"/electivote/manage-vote-page/",
		)
		return
	}

	context := gin.H{
		"title":      "Manage Vote",
		"voteData":   voteData,
		"candidates": candidates,
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
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
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
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
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
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
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

func ViewJoinVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	context := gin.H {
		"title": "Join Vote",
	}
	c.HTML(
		http.StatusOK,
		"joinVote.html",
		context,
	)
}

func JoinVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	voteCodeErr := ""
	voteCode := c.PostForm("voteCode")
	username := middlewares.GetUserData(c)
	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/join-vote-page/",
		)
		return
	}

	if repositories.IsVoted(uint(userID), voteCode) {
		voteCodeErr = "You already voted in this vote"
	}

	if len(voteCode) != 6 {
		voteCodeErr = "Vote code must be 6 characters"
	}

	_, err = repositories.GetVoteByVoteCode(voteCode)
	if len(voteCode) == 6 && err != nil {
		voteCodeErr = "Vote code not found"
	}

	if voteCodeErr != "" {
		context := gin.H {
			"title": "Join Vote",
			"voteCode": voteCode,
			"voteCodeErr": voteCodeErr,
		}
		c.HTML(
			http.StatusOK,
			"joinVote.html",
			context,
		)
		return
	}
	c.Redirect(
		http.StatusFound,
		"/electivote/vote-page/" + voteCode,
	)
}

func ViewVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	voteCode := c.Param("voteCode")
	voteID, err := repositories.GetVoteIDByVoteCode(voteCode)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	VoteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}
	
	context := gin.H {
		"title": "Vote",
		"candidates": candidates,
		"voteTitle": VoteData.VoteTitle,
		"voteDescription": VoteData.VoteDescription,
		"voteCode": voteCode,
	}
	c.HTML(
		http.StatusOK,
		"vote.html",
		context,
	)
}	

func VotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	votedErr := ""
	voted := c.PostForm("voted")
	voteCode := c.Param("voteCode")
	voteID, err := repositories.GetVoteIDByVoteCode(voteCode)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	VoteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	username := middlewares.GetUserData(c)
	userID, err := repositories.GetUserIdByUsername(username)

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	if voted == "" {
		votedErr = "Please select a candidate"
	}

	if votedErr != "" {
		context := gin.H {
			"title": "Vote",
			"votedErr": votedErr,
			"voteCode": voteCode,
			"voted":voted,
			"voteTitle": VoteData.VoteTitle,
			"voteDescription": VoteData.VoteDescription,
			"candidates": candidates,
		}
		c.HTML(
			http.StatusOK,
			"vote.html",
			context,
		)
		return
	}

	votedInt, _ := strconv.Atoi(voted)
	votedTime := models.CustomTime{Time: time.Now()}

	votedRecord := factories.VoteRecordFactory(uint(voteID), uint(userID), uint(votedInt), votedTime)
	_, err = repositories.CreateVoteRecord(votedRecord)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
		return
	}

	err = repositories.IncrementCandidateVote(uint(votedInt))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
		return
	}

	c.Redirect(
		http.StatusFound,
		"/electivote/home-page/",
	)
}

func ViewVoteResultPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	username := middlewares.GetUserData(c)
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
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

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	candidatesJson, err := json.Marshal(candidates)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	isExist := true

	if len(candidates) == 0 {
		isExist = false
	}

	context := gin.H {
		"title": "Vote Result",
		"voteID": voteID,
		"voteData": voteData,
		"candidates": candidates,
		"candidatesJson": string(candidatesJson),
		"isExist": isExist,
		"voteTitle": voteData.VoteTitle,
	}
	c.HTML(
		http.StatusOK,
		"voteResult.html",
		context,
	)
}