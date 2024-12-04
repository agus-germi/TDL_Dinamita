package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/agus-germi/TDL_Dinamita/utils"
	"github.com/sirupsen/logrus"
)

type LogrusLoggerAdapter struct {
	log *logrus.Logger
}

func (l *LogrusLoggerAdapter) Println(i ...interface{}) {
	l.log.Println(i...)
}

func (l *LogrusLoggerAdapter) Debug(i ...interface{}) {
	l.log.Debug(fmt.Sprintf("[DEBUG] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.log.Debugf(fmt.Sprintf("[DEBUG] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Info(i ...interface{}) {
	l.log.Info(fmt.Sprintf("[INFO] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Infof(format string, args ...interface{}) {
	l.log.Infof(fmt.Sprintf("[INFO] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Warn(i ...interface{}) {
	l.log.Warn(fmt.Sprintf("[WARN] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Warnf(format string, args ...interface{}) {
	l.log.Warnf(fmt.Sprintf("[WARN] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Error(i ...interface{}) {
	l.log.Error(fmt.Sprintf("[ERROR] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.log.Errorf(fmt.Sprintf("[ERROR] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Fatal(i ...interface{}) {
	l.log.Fatal(fmt.Sprintf("[FATAL] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(fmt.Sprintf("[FATAL] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Panic(i ...interface{}) {
	l.log.Panic(fmt.Sprintf("[PANIC] %s", fmt.Sprint(i...)))
}

func (l *LogrusLoggerAdapter) Panicf(format string, args ...interface{}) {
	l.log.Panicf(fmt.Sprintf("[PANIC] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Writer() io.Writer {
	return l.log.Writer()
}

func NewLogrusLoggerAdapter() *LogrusLoggerAdapter {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true, // Enable colors
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    false, // Make sure colors are enabled (no one can disable them)
		QuoteEmptyFields: true,
	})

	log.SetOutput(os.Stdout) // Choose any output (default: os.Stderr) --> Maybe we can use a file instead of stdout

	// Configure log levels according environment variables
	env, err := utils.GetEnv("APP_ENV")
	if err != nil {
		log.Errorf("[ERROR] Couldn't get 'APP_ENV' variable: %v", err)
	}

	if env == "prod" {
		log.SetLevel(logrus.ErrorLevel)
		log.Info("[INFO] Production environment detected.")
	} else if env == "dev" {
		log.SetLevel(logrus.DebugLevel)
		log.Info("[INFO] Development environment detected.")
	} else {
		log.SetLevel(logrus.InfoLevel) // If environment var "APP_ENV" was not provided --> Default level = InfoLevel
		log.Info("[INFO] APP_ENV variable wasn't set.")
	}
	return &LogrusLoggerAdapter{log: log}
}

/*
-------------- Log levels of logrus --------------
TraceLevel: El nivel más bajo, ideal para registros extremadamente detallados.
DebugLevel: Mensajes de depuración, útiles durante el desarrollo.
InfoLevel: Información general sobre el flujo de la aplicación, como eventos esperados.
WarnLevel: Advertencias sobre posibles problemas o situaciones inusuales.
ErrorLevel: Errores que ocurren, pero que no impiden la ejecución.
FatalLevel: Errores graves que causan que la aplicación termine.
PanicLevel: Errores graves que causan que la aplicación termine de inmediato, pero antes de hacerlo realiza un "panic".
*/
