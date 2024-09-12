package config

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func SetUpSessionStore() sessions.Store {
	mainStore := cookie.NewStore([]byte("main"))
	mainStore.Options(sessions.Options{
		MaxAge:   0,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
	})
	return mainStore
}