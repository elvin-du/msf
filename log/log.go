package log

import (
	"github.com/sirupsen/logrus"
)

var LogrusEntry *logrus.Entry

func init() {
	LogrusEntry = logrus.NewEntry(logrus.StandardLogger())
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	LogrusEntry.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	LogrusEntry.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	LogrusEntry.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	LogrusEntry.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	LogrusEntry.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	LogrusEntry.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	LogrusEntry.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	LogrusEntry.Fatal(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	LogrusEntry.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	LogrusEntry.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	LogrusEntry.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	LogrusEntry.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	LogrusEntry.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	LogrusEntry.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	LogrusEntry.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	LogrusEntry.Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	LogrusEntry.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	LogrusEntry.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	LogrusEntry.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	LogrusEntry.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	LogrusEntry.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	LogrusEntry.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	LogrusEntry.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	LogrusEntry.Fatalln(args...)
}
