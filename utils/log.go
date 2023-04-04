package utils

import (
	"encoding/json"

	"github.com/ghettovoice/gosip/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	m7slog "m7s.live/engine/v4/log"
)

type ZapLogger struct {
	log     *m7slog.Logger
	prefix  string
	fields  log.Fields
	sugared *zap.SugaredLogger
	level   log.Level
}

func NewZapLogger(log *m7slog.Logger, prefix string, fields log.Fields) (z *ZapLogger) {
	z = &ZapLogger{
		log:    log,
		prefix: prefix,
		fields: fields,
	}
	z.sugared = z.prepareEntry()
	return
}

func (l *ZapLogger) Print(args ...interface{}) {
	if l.level >= log.InfoLevel {
		l.sugared.Info(args...)
	}
}

func (l *ZapLogger) Printf(format string, args ...interface{}) {
	if l.level >= log.InfoLevel {
		l.sugared.Infof(format, args...)
	}
}

func (l *ZapLogger) Trace(args ...interface{}) {
	if l.level >= log.TraceLevel {
		l.sugared.Debug(args...)
	}
}

func (l *ZapLogger) Tracef(format string, args ...interface{}) {
	if l.level >= log.TraceLevel {
		l.sugared.Debugf(format, args...)
	}
}

func (l *ZapLogger) Debug(args ...interface{}) {
	if l.level >= log.DebugLevel {
		l.sugared.Debug(args...)
	}
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	if l.level >= log.DebugLevel {
		l.sugared.Debugf(format, args...)
	}
}

func (l *ZapLogger) Info(args ...interface{}) {
	if l.level >= log.InfoLevel {
		l.sugared.Info(args...)
	}
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	if l.level >= log.InfoLevel {
		l.sugared.Infof(format, args...)
	}
}

func (l *ZapLogger) Warn(args ...interface{}) {
	if l.level >= log.WarnLevel {
		l.sugared.Warn(args...)
	}
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	if l.level >= log.WarnLevel {
		l.sugared.Warnf(format, args...)
	}
}

func (l *ZapLogger) Error(args ...interface{}) {
	if l.level >= log.ErrorLevel {
		l.sugared.Error(args...)
	}
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	if l.level >= log.ErrorLevel {
		l.sugared.Errorf(format, args...)
	}
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	if l.level >= log.FatalLevel {
		l.sugared.Fatal(args...)
	}
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	if l.level >= log.FatalLevel {
		l.sugared.Fatalf(format, args...)
	}
}

func (l *ZapLogger) Panic(args ...interface{}) {
	if l.level >= log.PanicLevel {
		l.sugared.Panic(args...)
	}
}

func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	if l.level >= log.PanicLevel {
		l.sugared.Panicf(format, args...)
	}
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
	l.level = level
}
