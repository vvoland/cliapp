package cliapp

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/lmittmann/tint"
)

var (
	LogLevel        = slog.LevelDebug
	ShutdownTimeout = 1 * time.Second
)

func Init() context.Context {
	var handler slog.Handler

	if isTerminal(os.Stderr) {
		handler = niceHandler()
	} else {
		handler = plainHandler()
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx, done := context.WithCancel(context.Background())

	signalCtx, cancelSignal := signal.NotifyContext(ctx, os.Interrupt)
	context.AfterFunc(signalCtx, func() {
		slog.Warn("Received interrupt signal, shutting down gracefully...")
		done()
		time.AfterFunc(ShutdownTimeout, func() {
			slog.Warn("Shutdown timeout reached, next signal will terminate immediately")
			cancelSignal()
		})
	})

	return ctx
}

func plainHandler() slog.Handler {
	return slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey && len(groups) == 0 {
				level := a.Value.Any().(slog.Level)
				switch level {
				case slog.LevelDebug:
					return slog.String(a.Key, "DBG")
				case slog.LevelInfo:
					return slog.String(a.Key, "INF")
				case slog.LevelWarn:
					return slog.String(a.Key, "WRN")
				case slog.LevelError:
					return slog.String(a.Key, "ERR")
				}
			}
			return a
		},
	})
}

func niceHandler() slog.Handler {
	return tint.NewHandler(os.Stderr, &tint.Options{
		TimeFormat: "06-01-02 3:04PM",
		Level:      LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey && len(groups) == 0 {
				level := a.Value.Any().(slog.Level)
				switch level {
				case slog.LevelDebug:
					return tint.Attr(8+6, slog.String(a.Key, "DBG"))
				case slog.LevelInfo:
					return tint.Attr(8+2, slog.String(a.Key, "INF"))
				case slog.LevelWarn:
					return tint.Attr(0+3, slog.String(a.Key, "WRN"))
				case slog.LevelError:
					return tint.Attr(8+1, slog.String(a.Key, "ERR"))
				}
			}
			return a
		},
	})
}

func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
