package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"ip2country/internal/config"
)

func InitLogger(cnf *config.Config) {
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, options)
	if cnf.IsDebug {
		options.Level = slog.LevelDebug
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      options.Level,
			TimeFormat: time.Kitchen,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if err, ok := a.Value.Any().(error); ok {
					aErr := tint.Err(err)
					aErr.Key = a.Key
					return aErr
				}
				return a
			},
		})
	}
	logger := slog.New(handler)
	if !cnf.IsDebug {
		logger = logger.With(
			slog.String("service_name", cnf.Logger.ServiceName),
			slog.String("service_version", cnf.Logger.ServiceVersion),
		)
	}

	slog.SetDefault(logger)
}
