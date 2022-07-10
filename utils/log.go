package utils

import (
	"encoding/json"

	"github.com/ghettovoice/gosip/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log    *zap.Logger
	prefix string
	fields log.Fields
}

func NewZapLogger(log *zap.Logger, prefix string, fields log.Fields) *ZapLogger {
	return &ZapLogger{
		log:    log,
		prefix: prefix,
		fields: fields,
	}
}

func (l *ZapLogger) Print(args ...interface{}) {
	l.prepareEntry().Debug(args...)
}

func (l *ZapLogger) Printf(format string, args ...interface{}) {
	l.prepareEntry().Debugf(format, args...)
}

func (l *ZapLogger) Trace(args ...interface{}) {
	l.prepareEntry().Warn(args...)
}

func (l *ZapLogger) Tracef(format string, args ...interface{}) {
	l.prepareEntry().Warnf(format, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.prepareEntry().Debug(args...)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.prepareEntry().Debugf(format, args...)
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.prepareEntry().Info(args...)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.prepareEntry().Infof(format, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.prepareEntry().Warn(args...)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.prepareEntry().Warnf(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.prepareEntry().Error(args...)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.prepareEntry().Errorf(format, args...)
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.prepareEntry().Fatal(args...)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.prepareEntry().Fatalf(format, args...)
}

func (l *ZapLogger) Panic(args ...interface{}) {
	l.prepareEntry().Panic(args...)
}

func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.prepareEntry().Panicf(format, args...)
}

func (l *ZapLogger) WithPrefix(prefix string) log.Logger {
	return NewZapLogger(l.log, prefix, l.Fields())
}

func (l *ZapLogger) Prefix() string {
	return l.prefix
}

func (l *ZapLogger) WithFields(fields log.Fields) log.Logger {
	return NewZapLogger(l.log, l.Prefix(), l.Fields().WithFields(fields))
}

func (l *ZapLogger) Fields() log.Fields {
	return l.fields
}

func (l *ZapLogger) prepareEntry() *zap.SugaredLogger {

	newlog := l.log.With(zap.String("prefix", l.Prefix()))
	if l.fields != nil {
		fields := make([]zapcore.Field, len(l.fields))
		idx := 0
		for k, v := range l.fields {
			s, _ := json.Marshal(v)
			fields[idx] = zap.String(k, string(s))
			idx++
		}
		newlog = newlog.With(fields...)
	}
	return newlog.Sugar()
}

func (l *ZapLogger) SetLevel(level log.Level) {
	// zapcore.Level
}
