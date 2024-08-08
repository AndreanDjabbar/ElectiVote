package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndreanDjabbar/ElectiVote/internal/factories"
	"github.com/AndreanDjabbar/ElectiVote/internal/middlewares"
	"github.com/AndreanDjabbar/ElectiVote/internal/repositories"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-gonic/gin"
)

func ViewAddCandidatePage(c *gin.Context) {
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
	candidateNameErr := ""
	candidateName := c.PostForm("candidateName")
	candidateDescription := c.PostForm("candidateDesc")
	candidatePicture, candidatePictureErr := c.FormFile("candidatePicture")

	if len(candidateName) < 3 {
		candidateNameErr = "Candidate name must be at least 3 characters"
	}

	if candidateNameErr != "" {
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
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
			)
		}
		newCandidate.CandidatePicture = candidatePicture.Filename
	}
	_, err := repositories.AddCandidate(newCandidate)
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
		)
	}
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}

func ViewManageCandidatePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
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
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}

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
	if !middlewares.IsLogged(c) {
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
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
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
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-candidate-page/"+strconv.Itoa(voteID)+"/"+strconv.Itoa(candidateID),
		)
	}

	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}

func ViewDeleteCandidatePage(c *gin.Context) {
	if !middlewares.IsLogged(c) {
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
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}

	candidateData, err := repositories.GetCandidateByCandidateID(uint(candidateID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}
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
	if !middlewares.IsLogged(c) {
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
		c.Redirect(
			http.StatusFound,
			"/electivote/home-page/",
		)
		return
	}
	err := repositories.DeleteCandidate(uint(candidateID))
	if err != nil {
		utils.RenderError(
			c,
			http.StatusInternalServerError,
			err.Error(),
			"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
		)
	}
	c.Redirect(
		http.StatusFound,
		"/electivote/manage-vote-page/"+strconv.Itoa(voteID),
	)
}