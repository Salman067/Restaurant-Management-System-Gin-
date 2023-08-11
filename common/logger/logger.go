package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var logger *logrus.Logger
var loggerEntry *logrus.Entry

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
}

type LoggerInterface interface {
	Error(args ...interface{})
}

func NewLogger(client *raven.Client) LoggerInterface {
	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	buildTag := "develop"

	if viper.GetString("JSON_LOG") == "true" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	if viper.GetString("GIN_MODE") == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}

	loggerEntry = logrus.NewEntry(logger)

	if viper.GetString("BUILD_TAG") != "" {
		buildTag = viper.GetString("BUILD_TAG")
	}

	loggerEntry = loggerEntry.WithField("build", buildTag)

	if client != nil {
		hook, err := logrus_sentry.NewWithClientSentryHook(client, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		timeout := viper.GetInt("SENTRY_TIMEOUT")
		hook.Timeout = time.Duration(timeout) * time.Second
		hook.StacktraceConfiguration.Enable = true

		if err == nil {
			logger.Hooks.Add(hook)
		}
	}
	return logger
}

func NewRavenClient() *raven.Client {
	dsn := viper.Get("SENTRY_DSN")
	if dsn == nil {
		return nil
	}

	client, err := raven.New(dsn.(string))
	if err != nil {
		fmt.Println("Fatal")
		fmt.Println(err)
	}
	return client
}

func LogInfo(args ...interface{}) {
	// loggerEntry.Data["file"] = fileInfo(2)
	loggerEntry.WithField("file", fileInfo(2)).Info(args...)
}

func LogError(args ...interface{}) {
	// loggerEntry.Data["file"] = fileInfo(2)
	loggerEntry.WithField("file", fileInfo(2)).Error(args...)
}

func LogDebug(args ...interface{}) {
	// loggerEntry.Data["file"] = fileInfo(2)
	loggerEntry.WithField("file", fileInfo(2)).Debug(args...)
}

func WithField(key string, value interface{}) {
	loggerEntry = loggerEntry.WithField(key, value)
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}
	// else {
	// 	slash := strings.LastIndex(file, "/")
	//
	// 	if slash >= 0 {
	// 		file = file[slash+1:]
	// 	}
	// }
	return fmt.Sprintf("%s:%d", file, line)
}
