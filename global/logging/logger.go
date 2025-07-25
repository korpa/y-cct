package logging

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

func GetLogger() *slog.Logger {
	logHandler := log.NewWithOptions(os.Stderr, log.Options{
		Prefix:          "",
		ReportTimestamp: true,
		ReportCaller:    false,
		TimeFormat:      "15:04:05", //  time.RFC3339,
		Formatter:       log.TextFormatter,
	})

	// logHandler.SetStyles(styles)
	logger := slog.New(logHandler)
	return logger
}
