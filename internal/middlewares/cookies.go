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
	cookie, err := c.Cookie(cookieKey)
	if err != nil {
		logger.Error(
			"GetCookies - error getting cookies",
			"error", err,
			"Client IP", c.ClientIP(),
		)
		return ""
	} 
	resultCookie, err := utils.ExtractUsername(cookie)
	if err != nil {
		logger.Error(
			"GetCookies - error extracting username from cookies",
			"error", err,
			"Client IP", c.ClientIP(),
		)
		return ""
	}
	return resultCookie
}

func IsLogged(c *gin.Context) bool {
	if GetSession(c) != "" {
		return true
	} else if GetCookies(c) != "" {
		return true
	}
	return false
}

func DeleteCookie(c *gin.Context) {
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
	if GetSession(c) != "" {
		return GetSession(c)
	} else if GetCookies(c) != "" {
		return GetCookies(c)
	}
	return ""
}