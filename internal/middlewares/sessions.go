package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var sessionKey string = "session"

func SetSession(c *gin.Context, value string) {
	session := sessions.Default(c)
	session.Set(sessionKey, value)
	session.Save()
}

func GetSession(c *gin.Context) string {
	session := sessions.Default(c)
	value := session.Get(sessionKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

func DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(sessionKey)
	session.Save()
}