package config

import (
	"log/slog"
	"os"
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