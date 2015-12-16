package lib

import (
	"../../util"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
)

// ログ
var errorLog = logrus.New()
var accessLog = logrus.New()

func InitLog() {
	// ディレクトリの作成
	err := os.MkdirAll("log", 0777)
	if err != nil {
		util.Perror(err)
	}

	accessLogFile, err := os.OpenFile("log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	errorLogFile, err := os.OpenFile("log/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("[Error]: %s", err))
	}

	access_out := io.MultiWriter(os.Stdout, accessLogFile)
	accessLog.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: true}
	accessLog.Out = access_out

	error_out := io.MultiWriter(os.Stdout, errorLogFile)
	errorLog.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: true}
	errorLog.Out = error_out
}

func Info(v ...interface{}) {
	accessLog.Info(fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	accessLog.Debug(fmt.Sprint(v...))
}

func Error(v ...interface{}) {
	errorLog.Info(fmt.Sprint(v...))
}

func Fatal(v ...interface{}) {
	errorLog.Error(fmt.Sprint(v...))
}

func Fatalf(format string, v ...interface{}) {
	errorLog.Error(fmt.Sprintf(format, v...))
}
