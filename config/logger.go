package config

import (
	"log/slog"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetUpLogger() *slog.Logger {
	file, err := os.OpenFile(
		"logs/ElectiVote.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err) 
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	logger := slog.New(handler)
	return logger
}

func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		return "localhost"
	}
	return host
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

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


func SetUpRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("internal/views/html/*.html")
	router.Static("/images", "internal/assets/images")
	router.MaxMultipartMemory = 8 << 20
	return router
}