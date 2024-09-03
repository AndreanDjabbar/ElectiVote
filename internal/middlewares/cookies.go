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
		"SetCookies - Setting Cookies",
	)
	token, err := utils.GenerateSecureToken(username)
	if err != nil {
		logger.Error(
			"SetCookies - Failed to generate secure token",
			"error", err.Error(),
		)
		return
	}
	logger.Info(
		"SetCookies - setting cookies finished",
	)
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
		"GetCookies - Getting Cookies",
	)
	cookie, err := c.Cookie(cookieKey)
	if err != nil {
		logger.Warn(
			"GetCookies - No Cookies Found",
		)
		return ""
	} 
	resultCookie, err := utils.ExtractUsername(cookie)
	if err != nil {
		logger.Error(
			"GetCookies - Failed to extract username from cookies",
		)
		return ""
	}
	logger.Info(
		"GetCookies - Cookies Found",
	)
	return resultCookie
}

func IsLogged(c *gin.Context) bool {
	logger.Info(
		"IsLogged - Checking User Login Status",
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
		"DeleteCookie - Deleting Cookies",
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
		"GetUserData - Getting User Cookies Data",
	)
	if GetSession(c) != "" {
		return GetSession(c)
	} else if GetCookies(c) != "" {
		return GetCookies(c)
	}
	return ""
}