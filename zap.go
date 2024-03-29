package zap // import "go.unistack.org/micro-logger-zap/v3"

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.unistack.org/micro/v3/logger"
)

type zaplog struct {
	zap  *zap.Logger
	opts logger.Options
	sync.RWMutex
}

func (l *zaplog) Level(lvl logger.Level) {
}

func (l *zaplog) Clone(opts ...logger.Option) logger.Logger {
	nl := &zaplog{zap: l.zap, opts: l.opts}
	for _, o := range opts {
		o(&nl.opts)
	}
	return nl
}

func (l *zaplog) Init(opts ...logger.Option) error {
	var err error

	for _, o := range opts {
		o(&l.opts)
	}

	if zlog, ok := l.opts.Context.Value(loggerKey{}).(*zap.Logger); ok {
		l.zap = zlog
		return nil
	}

	zapConfig := zap.NewProductionConfig()
	if zconfig, ok := l.opts.Context.Value(configKey{}).(zap.Config); ok {
		zapConfig = zconfig
	}

	if zcconfig, ok := l.opts.Context.Value(encoderConfigKey{}).(zapcore.EncoderConfig); ok {
		zapConfig.EncoderConfig = zcconfig
	}

	skip, ok := l.opts.Context.Value(callerSkipKey{}).(int)
	if !ok || skip < 1 {
		skip = 1
	}

	// Set log Level if not default
	zapConfig.Level = zap.NewAtomicLevel()
	if l.opts.Level != logger.InfoLevel {
		zapConfig.Level.SetLevel(loggerToZapLevel(l.opts.Level))
	}

	log, err := zapConfig.Build(zap.AddCallerSkip(skip))
	if err != nil {
		return err
	}

	log = log.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
			zapcore.Lock(zapcore.AddSync(l.opts.Out)),
			loggerToZapLevel(l.opts.Level),
		)
	}))

	// Adding seed fields if exist
	if l.opts.Fields != nil {
		data := make([]zap.Field, 0, len(l.opts.Fields)/2)
		for i := 0; i < len(l.opts.Fields); i += 2 {
			data = append(data, zap.Any(l.opts.Fields[i].(string), l.opts.Fields[i+1]))
		}
		log = log.With(data...)
	}

	// Adding namespace
	if namespace, ok := l.opts.Context.Value(namespaceKey{}).(string); ok {
		log = log.With(zap.Namespace(namespace))
	}

	// defer log.Sync() ??

	l.zap = log

	return nil
}

func (l *zaplog) Fields(fields ...interface{}) logger.Logger {
	data := make([]zap.Field, 0, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		data = append(data, zap.Any(fields[i].(string), fields[i+1]))
	}

	zl := &zaplog{
		zap:  l.zap.With(data...),
		opts: l.opts,
	}

	return zl
}

func (l *zaplog) Errorf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.ErrorLevel, msg, args...)
}

func (l *zaplog) Debugf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.DebugLevel, msg, args...)
}

func (l *zaplog) Infof(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.InfoLevel, msg, args...)
}

func (l *zaplog) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.FatalLevel, msg, args...)
	os.Exit(1)
}

func (l *zaplog) Tracef(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.TraceLevel, msg, args...)
}

func (l *zaplog) Warnf(ctx context.Context, msg string, args ...interface{}) {
	l.Logf(ctx, logger.WarnLevel, msg, args...)
}

func (l *zaplog) Error(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.ErrorLevel, args...)
}

func (l *zaplog) Debug(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.DebugLevel, args...)
}

func (l *zaplog) Info(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.InfoLevel, args...)
}

func (l *zaplog) Fatal(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.FatalLevel, args...)
	os.Exit(1)
}

func (l *zaplog) Trace(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.TraceLevel, args...)
}

func (l *zaplog) Warn(ctx context.Context, args ...interface{}) {
	l.Log(ctx, logger.WarnLevel, args...)
}

func (l *zaplog) Log(ctx context.Context, level logger.Level, args ...interface{}) {
	if !l.V(level) {
		return
	}

	msg := fmt.Sprint(args...)
	switch loggerToZapLevel(level) {
	case zap.DebugLevel:
		l.zap.Debug(msg)
	case zap.InfoLevel:
		l.zap.Info(msg)
	case zap.WarnLevel:
		l.zap.Warn(msg)
	case zap.ErrorLevel:
		l.zap.Error(msg)
	case zap.FatalLevel:
		l.zap.Fatal(msg)
	}
}

func (l *zaplog) Logf(ctx context.Context, level logger.Level, format string, args ...interface{}) {
	if !l.V(level) {
		return
	}

	msg := fmt.Sprintf(format, args...)
	switch loggerToZapLevel(level) {
	case zap.DebugLevel:
		l.zap.Debug(msg)
	case zap.InfoLevel:
		l.zap.Info(msg)
	case zap.WarnLevel:
		l.zap.Warn(msg)
	case zap.ErrorLevel:
		l.zap.Error(msg)
	case zap.FatalLevel:
		l.zap.Fatal(msg)
	}
}

func (l *zaplog) V(level logger.Level) bool {
	return l.zap.Core().Enabled(loggerToZapLevel(level))
}

func (l *zaplog) String() string {
	return "zap"
}

func (l *zaplog) Options() logger.Options {
	return l.opts
}

// New builds a new logger based on options
func NewLogger(opts ...logger.Option) logger.Logger {
	l := &zaplog{opts: logger.NewOptions(opts...)}
	return l
}

func loggerToZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.TraceLevel, logger.DebugLevel:
		return zap.DebugLevel
	case logger.InfoLevel:
		return zap.InfoLevel
	case logger.WarnLevel:
		return zap.WarnLevel
	case logger.ErrorLevel:
		return zap.ErrorLevel
	case logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

/*
func zapToLoggerLevel(level zapcore.Level) logger.Level {
	switch level {
	case zap.DebugLevel:
		return logger.DebugLevel
	case zap.InfoLevel:
		return logger.InfoLevel
	case zap.WarnLevel:
		return logger.WarnLevel
	case zap.ErrorLevel:
		return logger.ErrorLevel
	case zap.FatalLevel:
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}
*/
