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
	if !repositories.IsValidModerator(username, uint(voteID)) {
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
	if !repositories.IsValidModerator(username, uint(voteID)) {
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
		newCandidate.CandidatePicture = candidatePicture.Filename
		if err != nil {
			utils.RenderError(
				c,
				http.StatusInternalServerError,
				err.Error(),
				"/electivote/add-candidate-page/"+strconv.Itoa(voteID),
			)
		}
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