package zap

import (
	"bytes"
	"context"
	"testing"

	"go.uber.org/zap"
	"go.unistack.org/micro/v3/logger"
)

func TestFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithLevel(logger.TraceLevel), logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Fields("key", "val").Info(ctx, "message")
	if !bytes.Contains(buf.Bytes(), []byte(`"key":"val"`)) {
		t.Fatalf("logger fields not works, buf contains: %s", buf.Bytes())
	}
}

func TestOutput(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewLogger(logger.WithOutput(buf))
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Infof(context.TODO(), "test logger name: %s", "name")
	if !bytes.Contains(buf.Bytes(), []byte(`test logger name`)) {
		t.Fatalf("log not redirected: %s", buf.Bytes())
	}
}

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
	logger.Infof(context.TODO(), "test logf: %s", "name")
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
	l.Debugf(context.TODO(), "test show debug: %s", "debug msg")

	if err := logger.Init(logger.WithLevel(logger.InfoLevel)); err != nil {
		t.Fatal(err)
	}
	l.Debugf(context.TODO(), "test non-show debug: %s", "debug msg")
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

	logger.Infof(context.TODO(), "test logf: %s", "name")
}
