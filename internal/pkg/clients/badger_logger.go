package clients

import "go.uber.org/zap"

type BadgerZapLogger struct {
	Log *zap.SugaredLogger
}

func (l *BadgerZapLogger) Errorf(s string, args ...interface{}) {
	l.Log.Errorf(s, args...)
}

func (l *BadgerZapLogger) Warningf(s string, args ...interface{}) {
	l.Log.Warnf(s, args...)
}

func (l *BadgerZapLogger) Infof(s string, args ...interface{}) {
	l.Log.Infof(s, args...)
}

func (l *BadgerZapLogger) Debugf(s string, args ...interface{}) {
	l.Log.Debugf(s, args...)
}
