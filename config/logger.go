package config

import (
    "log/slog"
    "os"
    "time"
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

    // Load Asia/Jakarta timezone
    location, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        panic(err)
    }

    // Custom time format function
    timeFormat := func(t time.Time) string {
        return t.In(location).Format("2006-01-02 15:04:05")
    }

    handler := slog.NewTextHandler(file, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        TimeFormat: timeFormat,
    })

    logger := slog.New(handler)
    return logger
}