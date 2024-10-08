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
		logger.Warn(
			"ViewCreateVotePage - User is not logged in",
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
		"ViewCreateVotePage - rendering create vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"CreateVotePage - User is not logged in",
			"Client IP", c.ClientIP(),
			"action", "redirecting to login page",
		)
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
		logger.Error(
			"CreateVotePage - failed to get user ID by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/create-vote-page/",
		)
	}

	if len(voteTitle) < 5 {
		logger.Warn(
			"CreateVotePage - vote title must be at least 5 characters",
			"Vote Title Inputted", voteTitle,
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.StartVoteFactory(voteTitle, voteDesc, voteCode, uint(moderatorID), start)
	
		_, err = repositories.CreateVote(newVote)
		if err != nil {
			logger.Error(
				"CreateVotePage - failed to create vote",
				"error", err.Error(),
				"Client IP", c.ClientIP(),
				"Username", username,
			)
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/create-vote-page/",
			)
		}

		logger.Info(
			"CreateVotePage - vote created",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
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
		logger.Warn(
			"ViewManageVotesPage - User is not logged in",
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
	votesData, err := repositories.GetVotesDataByUsername(username)
	if err != nil {
		logger.Error(
			"ViewManageVotesPage - failed to get votes data by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	logger.Info(
		"ViewManageVotesPage - rendering manage votes page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"ViewManageVotePage - User is not logged in",
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
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		logger.Warn(
			"ViewManageVotePage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
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
		logger.Error(
			"ViewManageVotePage - failed to get vote data by vote ID",
			"error", voteDataErr.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			voteDataErr.Error(),
			"/electivote/manage-vote-page/",
		)
		return
	}

	if candidatesErr != nil {
		logger.Error(
			"ViewManageVotePage - failed to get candidates by vote ID",
			"error", candidatesErr.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			candidatesErr.Error(),
			"/electivote/manage-vote-page/",
		)
		return
	}

	logger.Info(
		"ViewManageVotePage - rendering manage vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"ManageVotePage - User is not logged in",
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
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		logger.Warn(
			"ManageVotePage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ManageVotePage - failed to get vote data by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
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
		logger.Warn(
			"ManageVotePage - vote title must be at least 5 characters",
			"Vote Title Inputted", voteTitle,
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.UpdateVoteFactory(voteTitle, voteDesc)
		_, err := repositories.UpdateVote(uint(voteID), newVote)
		if err != nil {
			logger.Error(
				"ManageVotePage - failed to update vote",
				"error", err.Error(),
				"Client IP", c.ClientIP(),
				"Username", username,
			)
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/manage-vote-page/",
			)
		}

		logger.Info(
			"ManageVotePage - vote updated",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to manage vote page",
		)
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
		logger.Warn(
			"ViewDeleteVotePage - User is not logged in",
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
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		logger.Warn(
			"ViewDeleteVotePage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ViewDeleteVotePage - failed to get vote data by vote ID",
			"Client IP", c.ClientIP(),
			"Username", username,
			"error", err.Error(),
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	logger.Info(
		"ViewDeleteVotePage - rendering delete vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"DeleteVotePage - User is not logged in",
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
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		logger.Warn(
			"DeleteVotePage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to get vote data by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	moderatorName, err := repositories.GetModeratorNameByModeratorID(voteData.ModeratorID)
	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to get moderator name by moderator ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	candidateWinner, err := repositories.GetCandidateWinner(uint(voteID))
	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to get candidate winner",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
	}

	voteTitle := voteData.VoteTitle
	voteDescription := voteData.VoteDescription
	start := voteData.Start
	candidateWinnerName := candidateWinner.CandidateName
	candidateWinnerVotes := candidateWinner.TotalVotes
	candidateWinnerPicture := candidateWinner.CandidatePicture
	end := models.CustomTime{Time: time.Now()}
	voteHistory := factories.VoteHistoryFactory(voteData.ModeratorID, candidateWinnerVotes, moderatorName, voteTitle, voteDescription, candidateWinnerName, candidateWinnerPicture, start, end)
	err = repositories.CreateVoteHistory(voteHistory)

	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to create vote history",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	err = repositories.DeleteVote(uint(voteID))
	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to delete vote",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	logger.Info(
		"DeleteVotePage - vote deleted",
		"Client IP", c.ClientIP(),
		"Username", username,
		"action", "redirecting to manage votes page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/",
	)
}

func ViewJoinVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewJoinVotePage - User is not logged in",
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
		"ViewJoinVotePage - rendering join vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"JoinVotePage - User is not logged in",
			"Client IP", c.ClientIP(),
			"action", "redirecting to login page",
		)
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
		logger.Error(
			"JoinVotePage - failed to get user ID by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/join-vote-page/",
		)
		return
	}

	if repositories.IsVoted(uint(userID), voteCode) {
		logger.Warn(
			"JoinVotePage - user already voted in this vote",
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		voteCodeErr = "You already voted in this vote"
	}

	if len(voteCode) != 6 {
		logger.Warn(
			"JoinVotePage - vote code must be 6 characters",
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		voteCodeErr = "Vote code must be 6 characters"
	}

	_, err = repositories.GetVoteByVoteCode(voteCode)
	if len(voteCode) == 6 && err != nil {
		logger.Warn(
			"JoinVotePage - vote code not found",
			"Client IP", c.ClientIP(),
			"Username", username,
		)
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

	logger.Info(
		"JoinVotePage - redirecting to vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/vote-page/" + voteCode,
	)
}

func ViewVotePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVotePage - User is not logged in",
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
	voteCode := c.Param("voteCode")
	voteID, err := repositories.GetVoteIDByVoteCode(voteCode)
	if err != nil {
		logger.Error(
			"ViewVotePage - failed to get vote ID by vote code",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	VoteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ViewVotePage - failed to get vote data by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ViewVotePage - failed to get candidates by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}
	
	logger.Info(
		"ViewVotePage - rendering vote page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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
		logger.Warn(
			"VotePage - User is not logged in",
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
	votedErr := ""
	voted := c.PostForm("voted")
	voteCode := c.Param("voteCode")
	voteID, err := repositories.GetVoteIDByVoteCode(voteCode)
	if err != nil {
		logger.Error(
			"VotePage - failed to get vote ID by vote code",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	VoteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"VotePage - failed to get vote data by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"VotePage - failed to get user ID by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"VotePage - failed to get candidates by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}

	if voted == "" {
		logger.Warn(
			"VotePage - please select a candidate",
			"Client IP", c.ClientIP(),
			"Username", username,
		)
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
		logger.Error(
			"VotePage - failed to create vote record",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
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
		logger.Error(
			"VotePage - failed to increment candidate vote",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
		return
	}

	logger.Info(
		"VotePage - vote recorded",
		"Client IP", c.ClientIP(),
		"action", "redirecting to home page",
		"Username", username,
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/home-page/",
	)
}

func ViewVoteResultPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVoteResultPage - User is not logged in",
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
	voteID, _ := strconv.Atoi(c.Param("voteID"))
	if !repositories.IsValidVoteModerator(username, uint(voteID)) {
		logger.Warn(
			"ViewVoteResultPage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	voteData, err := repositories.GetVoteDataByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ViewVoteResultPage - failed to get vote data by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	candidates, err := repositories.GetCandidatesByVoteID(uint(voteID))
	if err != nil {
		logger.Error(
			"ViewVoteResultPage - failed to get candidates by vote ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	candidatesJson, err := json.Marshal(candidates)
	if err != nil {
		logger.Error(
			"ViewVoteResultPage - failed to marshal candidates",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/",
		)
	}

	isExist := true

	if len(candidates) == 0 {
		logger.Warn(
			"ViewVoteResultPage - no candidates found",
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		isExist = false
	}

	logger.Info(
		"ViewVoteResultPage - rendering vote result page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
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

func ViewVoteHistoryPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVoteHistoryPage - User is not logged in",
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
	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"ViewVoteHistoryPage - failed to get user ID by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
	}
	voteHistories, err := repositories.GetVoteHistoriesByUserID(uint(userID))
	if err != nil {
		logger.Error(
			"ViewVoteHistoryPage - failed to get vote records by user ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/home-page/",
		)
		return
	}

	logger.Info(
		"ViewVoteHistoryPage - rendering vote history page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
	context := gin.H {
		"title": "Vote History",
		"isExist": len(voteHistories) > 0,
		"voteHistories": voteHistories,
	}
	c.HTML(
		http.StatusOK,
		"voteHistory.html",
		context,
	)
}

func ViewVoteHistoryDetailPage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVoteHistoryDetailPage - User is not logged in",
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
	userID, err := repositories.GetUserIdByUsername(username)
	if err != nil {
		logger.Error(
			"ViewVoteHistoryDetailPage - failed to get user ID by username",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/vote-history-page/",
		)
		return
	}
	voteHistoryID, _ := strconv.Atoi(c.Param("voteHistoryID"))
	voteHistory, err := repositories.GetVoteHistoryByVoteHistoryID(uint(voteHistoryID))
	if err != nil {
		logger.Error(
			"ViewVoteHistoryDetailPage - failed to get vote history by vote history ID",
			"error", err.Error(),
			"Client IP", c.ClientIP(),
			"Username", username,
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/vote-history-page/",
		)
	}

	if uint(userID) != voteHistory.ModeratorID {
		logger.Warn(
			"ViewVoteHistoryDetailPage - User is not a valid vote moderator",
			"Client IP", c.ClientIP(),
			"Username", username,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/vote-history/page/",
		)
		return
	}
	isWinnerExist := true
	if voteHistory.CandidateWinnerName == "None" {
		isWinnerExist = false
	}

	logger.Info(
		"ViewVoteHistoryDetailPage - rendering vote history detail page",
		"Client IP", c.ClientIP(),
		"Username", username,
	)
	context := gin.H {
		"title": "Vote History Detail",
		"voteHistory": voteHistory,
		"isWinnerExist": isWinnerExist,
	}
	c.HTML(
		http.StatusOK,
		"voteHistoryDetail.html",
		context,
	)
}