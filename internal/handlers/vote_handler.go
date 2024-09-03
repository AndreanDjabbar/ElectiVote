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
	logger.Info(
		"ViewCreateVotePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewCreateVotePage - User Not Logged In",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	logger.Info(
		"ViewCreateVotePage - Rendering Page",
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
	logger.Info(
		"CreateVotePage - Creating Vote",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"CreateVotePage - User Not Logged In",
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
			"CreateVotePage - Error Getting Moderator ID",
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
			"CreateVotePage - Vote is less than 5 characters",
		)
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.StartVoteFactory(voteTitle, voteDesc, voteCode, uint(moderatorID), start)
	
		_, err = repositories.CreateVote(newVote)
		if err != nil {
			logger.Error(
				"CreateVotePage - Error Creating Vote",
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
			"CreateVotePage - Vote Created",
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	logger.Info(
		"CreateVotePage - Rendering Page",
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
		logger.Error(
			"ViewManageVotesPage - Error Getting Votes Data",
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
		"ViewManageVotesPage - Rendering Page",
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
	logger.Info(
		"ViewManageVotePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewManageVotePage - User Not Logged In",
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
			"ViewManageVotePage - User Not Moderator",
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
			"ViewManageVotePage - Error Getting Vote Data",
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
			"ViewManageVotePage - Error Getting Candidates",
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
		"ViewManageVotePage - Rendering Page",
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
	logger.Info(
		"ManageVotePage - Managing Vote",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ManageVotePage - User Not Logged In",
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
			"ManageVotePage - User Not Moderator",
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
			"ManageVotePage - Error Getting Vote Data",
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
			"ManageVotePage - Vote is less than 5 characters",
		)
		voteTitleErr = "Vote title must be at least 5 characters"
	}

	if voteTitleErr == "" {
		newVote := factories.UpdateVoteFactory(voteTitle, voteDesc)
		_, err := repositories.UpdateVote(uint(voteID), newVote)
		if err != nil {
			logger.Error(
				"ManageVotePage - Error Updating Vote",
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
			"ManageVotePage - Vote Updated",
			"action", "redirecting to manage vote page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/manage-vote-page/",
		)
		return
	}
	logger.Info(
		"ManageVotePage - Rendering Page",
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
	logger.Info(
		"ViewDeleteVotePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewDeleteVotePage - User Not Logged In",
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
			"ViewDeleteVotePage - User Not Moderator",
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
			"ViewDeleteVotePage - Error Getting Vote Data",
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
		"ViewDeleteVotePage - Rendering Page",
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
	logger.Info(
		"DeleteVotePage - Deleting Vote",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"DeleteVotePage - User Not Logged In",
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
			"DeleteVotePage - User Not Moderator",
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
			"DeleteVotePage - Error Deleting Vote",
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
		"DeleteVotePage - Vote Deleted",
		"action", "redirecting to manage votes page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/",
	)
}

func ViewJoinVotePage(c *gin.Context) {
	logger.Info(
		"ViewJoinVotePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewJoinVotePage - User Not Logged In",
			"action", "redirecting to login page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/login-page/",
		)
		return
	}
	logger.Info(
		"ViewJoinVotePage - Rendering Page",
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
	logger.Info(
		"JoinVotePage - Joining Vote",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"JoinVotePage - User Not Logged In",
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
			"JoinVotePage - Error Getting User ID",
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
			"JoinVotePage - User Already Voted",
		)
		voteCodeErr = "You already voted in this vote"
	}

	if len(voteCode) != 6 {
		logger.Warn(
			"JoinVotePage - Vote Code is not 6 characters",
		)
		voteCodeErr = "Vote code must be 6 characters"
	}

	_, err = repositories.GetVoteByVoteCode(voteCode)
	if len(voteCode) == 6 && err != nil {
		logger.Warn(
			"JoinVotePage - Vote Code Not Found",
		)
		voteCodeErr = "Vote code not found"
	}
	logger.Info(
		"JoinVotePage - Rendering Page",
	)
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
		"JoinVotePage - Redirecting to Vote Page",
		"action", "redirecting to vote page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/vote-page/" + voteCode,
	)
}

func ViewVotePage(c *gin.Context) {
	logger.Info(
		"ViewVotePage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVotePage - User Not Logged In",
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
			"ViewVotePage - Error Getting Vote ID",
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
			"ViewVotePage - Error Getting Vote Data",
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
			"ViewVotePage - Error Getting Candidates",
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
		"ViewVotePage - Rendering Page",
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
	logger.Info(
		"VotePage - Voting Candidate",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"VotePage - User Not Logged In",
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
			"VotePage - Error Getting Vote ID",
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
			"VotePage - Error Getting Vote Data",
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
	if err != nil {
		logger.Error(
			"VotePage - Error Getting User ID",
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
			"VotePage - Error Getting Candidates",
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
			"VotePage - No Candidate Selected",
		)
		votedErr = "Please select a candidate"
	}

	if votedErr != "" {
		logger.Info(
			"VotePage - Rendering Page",
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
			"VotePage - Error Creating Vote Record",
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
			"VotePage - Error Incrementing Candidate Vote",
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
		"VotePage - Voting Success",
		"action", "redirecting to home page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/home-page/",
	)
}

func ViewVoteResultPage(c *gin.Context) {
	logger.Info(
		"ViewVoteResultPage - Page Accessed",
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewVoteResultPage - User Not Logged In",
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
			"ViewVoteResultPage - User Not Moderator",
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
			"ViewVoteResultPage - Error Getting Vote Data",
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
			"ViewVoteResultPage - Error Getting Candidates",
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
			"ViewVoteResultPage - Error Marshalling Candidates",
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
		isExist = false
	}

	logger.Info(
		"ViewVoteResultPage - Rendering Page",
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