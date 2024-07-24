package middlewares

import (
	"time"

	"github.com/AndreanDjabbar/CaysAPIHub/internal/utils"
	"github.com/gin-gonic/gin"
)

var cookieKey string = "cookie"

func SetCookies(c *gin.Context, username string) {
	token, err := utils.GenerateSecureToken(username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate secure token"})
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
		return ""
	} 
	resultCookie, err := utils.ExtractUsername(cookie)
	if err != nil {
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