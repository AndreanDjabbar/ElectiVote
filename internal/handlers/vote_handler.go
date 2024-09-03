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
	logger.Info("ViewCreateVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewCreateVotePage - User is not logged in",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	logger.Info(
		"ViewCreateVotePage - rendering create vote page",
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
	logger.Info("CreateVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"CreateVotePage - User is not logged in",
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
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	logger.Info(
		"CreateVotePage - rendering create vote page with error messages",
	)
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
	logger.Info("ViewManageVotesPage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewManageVotesPage - User is not logged in",
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
	logger.Info("ViewManageVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewManageVotePage - User is not logged in",
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
	logger.Info("ManageVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ManageVotePage - User is not logged in",
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
			"action", "redirecting to manage vote page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/manage-vote-page/",
		)
		return
	}

	logger.Info(
		"ManageVotePage - rendering manage vote page with error messages",
	)
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
	logger.Info("ViewDeleteVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewDeleteVotePage - User is not logged in",
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
	logger.Info("DeleteVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"DeleteVotePage - User is not logged in",
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
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	err := repositories.DeleteVote(uint(voteID))
	if err != nil {
		logger.Error(
			"DeleteVotePage - failed to delete vote",
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
		"DeleteVotePage - vote deleted",
		"action", "redirecting to manage votes page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/",
	)
}

func ViewJoinVotePage(c *gin.Context) {
	logger.Info("ViewJoinVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewJoinVotePage - User is not logged in",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}

	logger.Info(
		"ViewJoinVotePage - rendering join vote page",
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
	logger.Info("JoinVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"JoinVotePage - User is not logged in",
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
		)
		voteCodeErr = "You already voted in this vote"
	}

	if len(voteCode) != 6 {
		logger.Warn(
			"JoinVotePage - vote code must be 6 characters",
		)
		voteCodeErr = "Vote code must be 6 characters"
	}

	_, err = repositories.GetVoteByVoteCode(voteCode)
	if len(voteCode) == 6 && err != nil {
		logger.Warn(
			"JoinVotePage - vote code not found",
		)
		voteCodeErr = "Vote code not found"
	}

	if voteCodeErr != "" {
		logger.Info(
			"JoinVotePage - rendering join vote page with error messages",
		)
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
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/vote-page/" + voteCode,
	)
}

func ViewVotePage(c *gin.Context) {
	logger.Info("ViewVotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVotePage - User is not logged in",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	voteCode := c.Param("voteCode")
	voteID, err := repositories.GetVoteIDByVoteCode(voteCode)
	if err != nil {
		logger.Error(
			"ViewVotePage - failed to get vote ID by vote code",
			"error", err.Error(),
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
	logger.Info("VotePage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"VotePage - User is not logged in",
			"action", "redirecting to login page",
		)
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
		logger.Error(
			"VotePage - failed to get vote ID by vote code",
			"error", err.Error(),
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
		)
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
		logger.Error(
			"VotePage - failed to get candidates by vote ID",
			"error", err.Error(),
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
		)
		votedErr = "Please select a candidate"
	}

	if votedErr != "" {
		logger.Info(
			"VotePage - rendering vote page with error messages",
		)
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
		"action", "redirecting to home page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/home-page/",
	)
}

func ViewVoteResultPage(c *gin.Context) {
	logger.Info("ViewVoteResultPage - page accessed")
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVoteResultPage - User is not logged in",
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
		)
		isExist = false
	}

	logger.Info(
		"ViewVoteResultPage - rendering vote result page",
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