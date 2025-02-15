package common

import (
	"fmt"
	"log/slog"
)

type Slogger struct {
	*slog.Logger
}

func (l Slogger) Debug(v ...interface{})                   {}
func (l Slogger) Debugf(format string, v ...interface{})   { l.Logger.Debug(fmt.Sprintf(format, v...)) }
func (l Slogger) Error(v ...interface{})                   {}
func (l Slogger) Errorf(format string, v ...interface{})   { l.Logger.Error(fmt.Sprintf(format, v...)) }
func (l Slogger) Info(v ...interface{})                    {}
func (l Slogger) Infof(format string, v ...interface{})    { l.Logger.Info(fmt.Sprintf(format, v...)) }
func (l Slogger) Warning(v ...interface{})                 {}
func (l Slogger) Warningf(format string, v ...interface{}) { l.Logger.Warn(fmt.Sprintf(format, v...)) }
func (l Slogger) Fatal(v ...interface{})                   {}
func (l Slogger) Fatalf(format string, v ...interface{})   { l.Logger.Error(fmt.Sprintf(format, v...)) }
func (l Slogger) Panic(v ...interface{})                   {}
func (l Slogger) Panicf(format string, v ...interface{})   { l.Logger.Error(fmt.Sprintf(format, v...)) }
