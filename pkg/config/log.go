package config

import (
	"io"
	"log/slog"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"
)

func setupLog(conf *logConf) *slog.Logger {
	var level slog.Leveler = slog.LevelDebug
	for _, l := range []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError} {
		if l.String() == conf.Level {
			level = l
			break
		}
	}

	if conf.MaxSize == 0 {
		conf.MaxSize = 256
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = 7
	}

	var writer io.Writer
	if !conf.OnConsole {
		writer = &lumberjack.Logger{
			Filename:  path.Join("log", "app.log"),
			MaxSize:   conf.MaxSize,
			MaxAge:    conf.MaxAge,
			Compress:  true,
			LocalTime: true,
		}
	} else {
		writer = os.Stdout
	}

	var handler slog.Handler
	if conf.IsJson {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level, AddSource: conf.Trace})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level, AddSource: conf.Trace})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
