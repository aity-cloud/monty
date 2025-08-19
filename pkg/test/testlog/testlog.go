package testlog

import (
	"log/slog"

	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/test/testruntime"
	"github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
)

var Log = lo.TernaryF(testruntime.IsTesting, func() *slog.Logger {
	return logger.New(logger.WithLogLevel(logger.DefaultLogLevel.Level()), logger.WithWriter(ginkgo.GinkgoWriter)).WithGroup("test")
}, func() *slog.Logger {
	return logger.New(logger.WithLogLevel(logger.DefaultLogLevel.Level())).WithGroup("test")
})
