package routes

import (
	"net/http"

	"github.com/AndreanDjabbar/ElectiVote/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)	
}

func MainRootHandler(c *gin.Context) {
	c.Redirect(
		http.StatusFound,
		"/electivote/login-page/",
	)
}

func SetUpRoutes(router *gin.Engine) {
	mainRouter := router.Group("/electivote")
	{
		router.GET("/", RootHandler)
		mainRouter.GET("/", MainRootHandler)
	}
	{
		mainRouter.GET("login-page/", handlers.ViewLoginPage)
		mainRouter.POST("login-page/", handlers.LoginPage)
		mainRouter.GET("register-page/", handlers.ViewRegisterPage)
		mainRouter.POST("register-page/", handlers.RegisterPage)
		mainRouter.GET("logout/", handlers.LogoutPage)
		mainRouter.GET("home-page/", handlers.ViewHomePage)
	}
	{
		mainRouter.GET("profile-page/", handlers.ViewProfilePage)
		mainRouter.GET("edit-profile-page/", handlers.ViewEditProfilePage)
		mainRouter.POST("edit-profile-page/", handlers.EditProfilePage)
	}
	{
		mainRouter.GET("create-vote-page/", handlers.ViewCreateVotePage)
		mainRouter.POST("create-vote-page/", handlers.CreateVotePage)
		mainRouter.GET("manage-vote-page/", handlers.ViewManageVotesPage)
		mainRouter.GET("manage-vote-page/:voteID/", handlers.ViewManageVotePage)
		mainRouter.POST("manage-vote-page/:voteID/", handlers.ManageVotePage)
		mainRouter.GET("delete-vote-page/:voteID/", handlers.ViewDeleteVotePage)
		mainRouter.GET("delete-vote/:voteID/", handlers.DeleteVotePage)
	}
	{
		mainRouter.GET("add-candidate-page/:voteID/", handlers.ViewAddCandidatePage)
		mainRouter.POST("add-candidate-page/:voteID/", handlers.AddCandidatePage)
		mainRouter.GET("manage-candidate-page/:voteID/:candidateID/", handlers.ViewManageCandidatePage)
		mainRouter.POST("manage-candidate-page/:voteID/:candidateID/", handlers.ManageCandidatePage)
		mainRouter.GET("delete-candidate-page/:voteID/:candidateID/", handlers.ViewDeleteCandidatePage)
		mainRouter.GET("delete-candidate/:voteID/:candidateID/", handlers.DeleteCandidatePage)
	}
}