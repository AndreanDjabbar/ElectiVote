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