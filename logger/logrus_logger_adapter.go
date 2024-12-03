package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type LogrusLoggerAdapter struct {
	log *logrus.Logger
}

func (l *LogrusLoggerAdapter) Println(i ...interface{}) {
	l.log.Println(i...)
}

func (l *LogrusLoggerAdapter) Debug(i ...interface{}) {
	l.log.Debug(i...)
}

func (l *LogrusLoggerAdapter) Debugf(format string, args ...interface{}) {
	l.log.Debugf(fmt.Sprintf("[DEBUG] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Info(i ...interface{}) {
	l.log.Info(i...)
}

func (l *LogrusLoggerAdapter) Infof(format string, args ...interface{}) {
	l.log.Infof(fmt.Sprintf("[INFO] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Warn(i ...interface{}) {
	l.log.Warn(i...)
}

func (l *LogrusLoggerAdapter) Warnf(format string, args ...interface{}) {
	l.log.Warnf(fmt.Sprintf("[WARN] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Error(i ...interface{}) {
	l.log.Error(i...)
}

func (l *LogrusLoggerAdapter) Errorf(format string, args ...interface{}) {
	l.log.Errorf(fmt.Sprintf("[ERROR] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Fatal(i ...interface{}) {
	l.log.Fatal(i...)
}

func (l *LogrusLoggerAdapter) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(fmt.Sprintf("[FATAL] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Panic(i ...interface{}) {
	l.log.Panic(i...)
}

func (l *LogrusLoggerAdapter) Panicf(format string, args ...interface{}) {
	l.log.Panicf(fmt.Sprintf("[PANIC] %s", format), args...)
}

func (l *LogrusLoggerAdapter) Writer() io.Writer {
	return l.log.Writer()
}

func NewLogrusLoggerAdapter() *LogrusLoggerAdapter {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{}) // Choose any format (default: &logrus.TextFormatter{})
	log.SetOutput(os.Stdout)

	// If the APP_ENV variable isn't recognized, try to load the .env file first
	err := godotenv.Load("/usr/src/app/.env")
	if err != nil {
		log.Errorln("Error: .env file couldn't be loaded.")
	}

	// Configure log levels according environment variables
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.SetLevel(logrus.ErrorLevel)
		log.Println("Production environment detected.")
	} else if env == "deveploment" {
		log.SetLevel(logrus.DebugLevel)
		log.Println("Development environment detected.")
	} else {
		log.SetLevel(logrus.DebugLevel) // If environment var "APP_ENV" was not provided --> Default level = InfoLevel
		//log.SetLevel(logrus.InfoLevel) // If environment var "APP_ENV" was not provided --> Default level = InfoLevel
		log.Println("APP_ENV variable wasn't set.")
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
