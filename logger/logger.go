package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
	initErr  error // Package-level error tracker for sync.Once
)

func GetLogger() (*slog.Logger, error) {
	once.Do(func() {
		file, err := os.OpenFile(
			"logs.log",
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			0644,
		)
		if err != nil {
			initErr = err
			return
		}

		writer := io.MultiWriter(os.Stdout, file)

		instance = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}
