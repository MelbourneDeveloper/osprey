// Package logging provides process-wide structured logging for the compiler.
//
// Logging is disabled unless OSPREY_LOG_LEVEL is set to trace, debug, info,
// warn, or error. Records are written to OSPREY_LOG_FILE, or ./osprey.log when
// no file is supplied. JSON is the default format; set OSPREY_LOG_FORMAT=text
// for text output. OSPREY_LOG_SOURCE=true adds source locations.
package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	envFile   = "OSPREY_LOG_FILE"
	envFormat = "OSPREY_LOG_FORMAT"
	envLevel  = "OSPREY_LOG_LEVEL"
	envSource = "OSPREY_LOG_SOURCE"

	defaultLogFile = "osprey.log"
	logDirMode     = 0o755
	logFileMode    = 0o600
)

const (
	// LevelTrace is more verbose than slog.LevelDebug and is intended for
	// diagnostic compiler traces that are too noisy for normal debugging.
	LevelTrace slog.Level = slog.LevelDebug - 4
	levelOff   slog.Level = slog.LevelError + 100
)

//nolint:gochecknoglobals // The compiler logger is process-wide so all packages share one sink and level.
var (
	baseLogger *slog.Logger
	levelVar   slog.LevelVar
	once       sync.Once
)

// Logger returns a structured logger tagged with the supplied component.
func Logger(component string) *slog.Logger {
	configure()

	return baseLogger.With("component", component)
}

func configure() {
	once.Do(func() {
		level := parseLevel(os.Getenv(envLevel))
		levelVar.Set(level)

		writer := io.Discard
		if level != levelOff {
			file, err := openLogFile()
			if err == nil {
				writer = file
			}
		}

		options := &slog.HandlerOptions{
			Level:       &levelVar,
			AddSource:   parseBool(os.Getenv(envSource)),
			ReplaceAttr: replaceLevelAttr,
		}

		if strings.EqualFold(strings.TrimSpace(os.Getenv(envFormat)), "text") {
			baseLogger = slog.New(slog.NewTextHandler(writer, options))
			return
		}

		baseLogger = slog.New(slog.NewJSONHandler(writer, options))
	})
}

func openLogFile() (*os.File, error) {
	path := strings.TrimSpace(os.Getenv(envFile))
	if path == "" {
		path = defaultLogFile
	}

	if dir := filepath.Dir(path); dir != "." && dir != "" {
		err := os.MkdirAll(dir, logDirMode)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFileMode)
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "trace":
		return LevelTrace
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return levelOff
	}
}

func parseBool(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func replaceLevelAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Key != slog.LevelKey {
		return attr
	}

	if level, ok := attr.Value.Any().(slog.Level); ok && level == LevelTrace {
		attr.Value = slog.StringValue("TRACE")
	}

	return attr
}
