package zap

import (
	"testing"

	"github.com/unistack-org/micro/v3/logger"
	"go.uber.org/zap"
)

func TestName(t *testing.T) {
	l := NewLogger()
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l := NewLogger()
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l
	logger.Logf(logger.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l := NewLogger()
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	logger.DefaultLogger = l

	if err := logger.Init(logger.WithLevel(logger.DebugLevel)); err != nil {
		t.Fatal(err)
	}
	l.Logf(logger.DebugLevel, "test show debug: %s", "debug msg")

	if err := logger.Init(logger.WithLevel(logger.InfoLevel)); err != nil {
		t.Fatal(err)
	}
	l.Logf(logger.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWrapper(t *testing.T) {
	z, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	zl := NewLogger(WithLogger(z))
	if err := zl.Init(); err != nil {
		t.Fatal(err)
	}
	logger.DefaultLogger = zl

	logger.Logf(logger.InfoLevel, "test logf: %s", "name")
}
