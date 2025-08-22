package k8sutil

import (
	"log/slog"

	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/go-logr/logr"
)

func NewControllerRuntimeLogger(level slog.Level) logr.Logger {
	return logger.NewLogr(logger.WithTimeFormat("[15:04:05]"), logger.WithLogLevel(level))
}
