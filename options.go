package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.unistack.org/micro/v3/logger"
)

type callerSkipKey struct{}

func WithCallerSkip(i int) logger.Option {
	return logger.SetOption(callerSkipKey{}, i)
}

type configKey struct{}

// WithConfig pass zap.Config to logger
func WithConfig(c zap.Config) logger.Option {
	return logger.SetOption(configKey{}, c)
}

type loggerKey struct{}

// WithLogger pass *zap.Logger to logger
func WithLogger(l *zap.Logger) logger.Option {
	return logger.SetOption(loggerKey{}, l)
}

type encoderConfigKey struct{}

// WithEncoderConfig pass zapcore.EncoderConfig to logger
func WithEncoderConfig(c zapcore.EncoderConfig) logger.Option {
	return logger.SetOption(encoderConfigKey{}, c)
}

type namespaceKey struct{}

func WithNamespace(namespace string) logger.Option {
	return logger.SetOption(namespaceKey{}, namespace)
}
