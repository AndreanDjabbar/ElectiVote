package middlewares

import (
	"log/slog"
	"os"
	"time"

	"github.com/AndreanDjabbar/ElectiVote/config"
	"github.com/AndreanDjabbar/ElectiVote/internal/utils"
	"github.com/gin-gonic/gin"
)

var cookieKey string = os.Getenv("COOKIE_KEY")
var logger *slog.Logger = config.SetUpLogger()

func SetCookies(c *gin.Context, username string) {
	logger.Info(
		"SetCookies - setting cookies",
	)
	token, err := utils.GenerateSecureToken(username)
	if err != nil {
		logger.Error(
			"SetCookies - error generating secure token",
			"error", err,
		)
		return
	}
	expiration := time.Now().Add(7 * 24 * time.Hour)
	c.SetCookie(
		cookieKey,
		token,
		int(expiration.Unix()),
		"/",
		"localhost",
		false,
		true,
	)
}

func GetCookies(c *gin.Context) string {
	logger.Info(
		"GetCookies - getting cookies",
	)
	cookie, err := c.Cookie(cookieKey)
	if err != nil {
		logger.Error(
			"GetCookies - error getting cookies",
			"error", err,
		)
		return ""
	} 
	resultCookie, err := utils.ExtractUsername(cookie)
	if err != nil {
		logger.Error(
			"GetCookies - error extracting username from cookies",
			"error", err,
		)
		return ""
	}
	return resultCookie
}

func IsLogged(c *gin.Context) bool {
	logger.Info(
		"IsLogged - checking if user is logged in",
	)
	if GetSession(c) != "" {
		return true
	} else if GetCookies(c) != "" {
		return true
	}
	return false
}

func DeleteCookie(c *gin.Context) {
	logger.Info(
		"DeleteCookie - deleting cookies",
	)
	c.SetCookie(
		cookieKey,
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)
}

func GetUserData(c *gin.Context) string {
	logger.Info(
		"GetUserData - getting user session or cookies data",
	)
	if GetSession(c) != "" {
		return GetSession(c)
	} else if GetCookies(c) != "" {
		return GetCookies(c)
	}
	return ""
}