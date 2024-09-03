package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var sessionKey string = os.Getenv("SESSION_KEY")

func SetSession(c *gin.Context, value string) {
	logger.Info(
		"SetSession - setting session",
	)
	session := sessions.Default(c)
	session.Set(sessionKey, value)
	session.Save()
}

func GetSession(c *gin.Context) string {
	logger.Info(
		"GetSession - getting session",
	)
	session := sessions.Default(c)
	value := session.Get(sessionKey)
	if value == nil {
		logger.Warn(
			"GetSession - session is empty",
		)
		return ""
	}
	return value.(string)
}

func DeleteSession(c *gin.Context) {
	logger.Info(
		"DeleteSession - deleting session",
	)
	session := sessions.Default(c)
	session.Delete(sessionKey)
	session.Save()
}

func SetRegisterSession(c *gin.Context, username, email, password, otp string) {
	logger.Info(
		"SetRegisterSession - setting register session",
	)
	session := sessions.Default(c)
	session.Set("username", username)
	session.Set("email", email)
	session.Set("password", password)
	session.Set("otp", otp)
	session.Set("created_at", time.Now().Unix())
	if err := session.Save(); err != nil {
		logger.Error(
			"SetRegisterSession - error saving session",
		)
		fmt.Println("Error saving session: ", err)
	}
}

func DeleteRegisterSession(c *gin.Context) {
	logger.Info(
		"DeleteRegisterSession - deleting register session",
	)
	session := sessions.Default(c)
	session.Delete("username")
	session.Delete("email")
	session.Delete("password")
	session.Delete("otp")
	session.Save()
}