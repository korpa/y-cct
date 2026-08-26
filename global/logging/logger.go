package logging

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

func GetLogger() *slog.Logger {
	if !isTerminal(os.Stdout) {
		// charmbracelet/log's SetOutput() unconditionally does
		// termenv.WithColorCache(true), which eagerly sends an OSC10 fg-color
		// query and a CSI6n cursor-position query to the terminal and blocks
		// on the reply. When stdout is redirected (e.g. process substitution
		// like `source <(y completion bash)`) but stderr is still a tty, that
		// query still fires and can race with anything else writing to stdin
		// at the same time. Skip the styled renderer entirely in that case.
		return slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

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

func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
