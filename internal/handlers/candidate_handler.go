package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-gonic/gin"
)

var logger *slog.Logger = config.SetUpLogger()

func ViewAddCandidatePage(c *gin.Context) {
	logger.Info(
		"ViewAddCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewAddCandidatePage - User not logged in",
			"method", c.Request.Method,
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
			"ViewAddCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	logger.Info(
		"ViewAddCandidatePage - Rendering Page",
		"method", c.Request.Method,
	)
	context := gin.H {
		"title": "Add Candidate",
		"voteID": voteID,
	}
	c.HTML(
		http.StatusOK,
		"addCandidate.html",
		context,
	)
}

func AddCandidatePage(c *gin.Context) {
	logger.Info(
		"AddCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"AddCandidatePage - User not logged in",
			"method", c.Request.Method,
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
			"AddCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	candidateNameErr := ""
	candidateName := c.PostForm("candidateName")
	candidateDescription := c.PostForm("candidateDesc")
	candidatePicture, candidatePictureErr := c.FormFile("candidatePicture")

	if len(candidateName) < 3 {
		candidateNameErr = "Candidate name must be at least 3 characters"
	}

	if candidateNameErr != "" {
		logger.Warn(
			"AddCandidatePage - Invalid Input",
			"method", c.Request.Method,
			"action", "rendering page with error message",
		)
		context := gin.H {
			"title": "Add Candidate",
			"voteID": voteID,
			"candidateNameErr": candidateNameErr,
			"candidateName": candidateName,
			"candidateDescription": candidateDescription,
		}
		c.HTML(
			http.StatusBadRequest,
			"addCandidate.html",
			context,
		)
		return
	}
	newCandidate := factories.CandidateFactory(candidateName, candidateDescription, "default.png", uint(voteID))

	if candidatePicture != nil {
		if candidatePictureErr != nil {
			logger.Warn(
				"AddCandidatePage - Error Uploading Picture",
				"method", c.Request.Method,
				"action", "rendering page with error message",
			)
			utils.RenderError(
				c,
				http.StatusBadRequest,
				candidatePictureErr.Error(),
				"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
			)
		}
		err := c.SaveUploadedFile(
			candidatePicture,
			"internal/assets/images/"+candidatePicture.Filename,
		)
		if err != nil {
			logger.Warn(
				"AddCandidatePage - Error Saving Picture",
				"method", c.Request.Method,
				"action", "rendering page with error message",
			)
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
			)
		}
		logger.Info(
			"AddCandidatePage - Picture Uploaded",
			"method", c.Request.Method,
		)
		newCandidate.CandidatePicture = candidatePicture.Filename
	}
	_, err := repositories.AddCandidate(newCandidate)
	if err != nil {
		logger.Warn(
			"AddCandidatePage - Error Adding Candidate",
			"method", c.Request.Method,
			"action", "rendering page with error message",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
		)
	}
	logger.Info(
		"AddCandidatePage - Candidate Added",
		"method", c.Request.Method,
		"action", "redirecting to manage vote page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}

func ViewManageCandidatePage(c *gin.Context) {
	logger.Info(
		"ViewManageCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewManageCandidatePage - User not logged in",
			"method", c.Request.Method,
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
	candidateID, _ := strconv.Atoi(c.Param("candidateID"))

	if !repositories.IsValidVoteModerator(username, uint(voteID)) || !repositories.IsValidCandidateModerator(username, uint(candidateID)) {
		logger.Warn(
			"ViewManageCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
		logger.Warn(
			"ViewManageCandidatePage - Error Getting Candidate",
			"method", c.Request.Method,
			"action", "rendering error page",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}

	logger.Info(
		"ViewManageCandidatePage - Rendering Page",
		"method", c.Request.Method,
	)
	context := gin.H {
		"title": "Manage Candidate",
		"voteID": voteID,
		"candidateID": candidateID,
		"candidateData": candidateData,
	}
	c.HTML(
		http.StatusOK,
		"manageCandidate.html",
		context,
	)
}

func ManageCandidatePage(c *gin.Context) {
	logger.Info(
		"ManageCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ManageCandidatePage - User not logged in",
			"method", c.Request.Method,
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
	candidateID, _ := strconv.Atoi(c.Param("candidateID"))

	if !repositories.IsValidVoteModerator(username, uint(voteID)) || !repositories.IsValidCandidateModerator(username, uint(candidateID)) {
		logger.Warn(
			"ManageCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
		logger.Warn(
			"ManageCandidatePage - Error Getting Candidate",
			"method", c.Request.Method,
			"action", "rendering error page",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}

	candidateNameErr := ""
	candidateName := c.PostForm("candidateName")
	candidateDescription := c.PostForm("candidateDesc")
	candidatePicture, candidatePictureErr := c.FormFile("candidatePicture")

	if len(candidateName) < 3 {
		candidateNameErr = "Candidate name must be at least 3 characters"
	}

	if candidateNameErr != "" {
		logger.Warn(
			"ManageCandidatePage - Invalid Input",
			"method", c.Request.Method,
			"action", "rendering page with error message",
		)
		context := gin.H {
			"title": "Manage Candidate",
			"voteID": voteID,
			"candidateID": candidateID,
			"candidateNameErr": candidateNameErr,
			"candidateName": candidateName,
			"candidateDescription": candidateDescription,
			"candidateData": candidateData,
		}
		c.HTML(
			http.StatusBadRequest,
			"manageCandidate.html",
			context,
		)
		return
	}

	updatedCandidate := factories.CandidateFactory(candidateName, candidateDescription, candidateData.CandidatePicture, uint(voteID))

	if candidatePicture != nil {
		if candidatePictureErr != nil {
			logger.Warn(
				"ManageCandidatePage - Error Uploading Picture",
				"method", c.Request.Method,
				"action", "rendering page with error message",
			)
			utils.RenderError(
				c,
				http.StatusBadRequest,
				candidatePictureErr.Error(),
				"/electivote/manage-candidate-page/"+strconv.Itoa(voteID)+"/"+strconv.Itoa(candidateID),
			)
		}
		err := c.SaveUploadedFile(
			candidatePicture,
			"internal/assets/images/"+candidatePicture.Filename,
		)
		if err != nil {
			logger.Warn(
				"ManageCandidatePage - Error Saving Picture",
				"method", c.Request.Method,
				"action", "rendering page with error message",
			)
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/manage-candidate-page/"+strconv.Itoa(voteID)+"/"+strconv.Itoa(candidateID),
			)
		}
		updatedCandidate.CandidatePicture = candidatePicture.Filename
	}

	_, err = repositories.UpdateCandidate(uint(candidateID), updatedCandidate)
	if err != nil {	
		logger.Warn(
			"ManageCandidatePage - Error Updating Candidate",
			"method", c.Request.Method,
			"action", "rendering page with error message",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-candidate-page/"+strconv.Itoa(voteID)+"/"+strconv.Itoa(candidateID),
		)
	}
	logger.Info(
		"ManageCandidatePage - Candidate Updated",
		"method", c.Request.Method,
		"action", "redirecting to manage vote page",
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}

func ViewDeleteCandidatePage(c *gin.Context) {
	logger.Info(
		"ViewDeleteCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"ViewDeleteCandidatePage - User not logged in",
			"method", c.Request.Method,
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
	candidateID, _ := strconv.Atoi(c.Param("candidateID"))

	if !repositories.IsValidVoteModerator(username, uint(voteID)) || !repositories.IsValidCandidateModerator(username, uint(candidateID)) {
		logger.Warn(
			"ViewDeleteCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
		logger.Warn(
			"ViewDeleteCandidatePage - Error Getting Candidate",
			"method", c.Request.Method,
			"action", "rendering error page",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}
	logger.Info(
		"ViewDeleteCandidatePage - Rendering Page",
		"method", c.Request.Method,
	)
	context := gin.H {
		"title": "Delete Candidate",
		"voteID": voteID,
		"candidateID": candidateID,
		"candidateData": candidateData,
	}
	c.HTML(
		http.StatusOK,
		"deleteCandidate.html",
		context,
	)
}

func DeleteCandidatePage(c *gin.Context) {
	logger.Info(
		"DeleteCandidatePage - Page Accessed",
		"method", c.Request.Method,
	)
	if !middlewares.IsLogged(c) {
		logger.Warn(
			"DeleteCandidatePage - User not logged in",
			"method", c.Request.Method,
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
	candidateID, _ := strconv.Atoi(c.Param("candidateID"))

	if !repositories.IsValidVoteModerator(username, uint(voteID)) || !repositories.IsValidCandidateModerator(username, uint(candidateID)) {
		logger.Warn(
			"DeleteCandidatePage - User not authorized",
			"method", c.Request.Method,
			"action", "redirecting to home page",
		)
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	err := repositories.DeleteCandidate(uint(candidateID))
	if err != nil {
		logger.Warn(
			"DeleteCandidatePage - Error Deleting Candidate",
			"method", c.Request.Method,
			"action", "rendering page with error message",
		)
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}
	logger.Info(
		"DeleteCandidatePage - Candidate Deleted",
		"method", c.Request.Method,
	)
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}