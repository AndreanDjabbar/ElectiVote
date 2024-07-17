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