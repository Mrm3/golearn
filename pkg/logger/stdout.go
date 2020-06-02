package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"immortality/service/db_model"
)

type StdoutLogger struct {
	w           io.WriteCloser

	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var Log StdoutLogger

func NewStdoutLogger(w io.WriteCloser, prefix string) StdoutLogger {
	return StdoutLogger{
		w: w,
		infoLogger: log.New(w, prefix+" [INFO] ",
			log.Ldate|log.Ltime|log.Lmicroseconds),
		warnLogger: log.New(w, prefix+" [WARN] ",
			log.Ldate|log.Ltime|log.Lmicroseconds),
		errorLogger: log.New(w, prefix+" [ERROR] ",
			log.Ldate|log.Ltime|log.Lmicroseconds),
		debugLogger: log.New(w, prefix+" [DEBUG] ",
			log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

func getCaller(skip int) string {
	_, fullPath, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	fileParts := strings.Split(fullPath, "/")
	file := fileParts[len(fileParts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func (l StdoutLogger) Info(args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller}, args...)
	l.infoLogger.Println(args...)
}

func (l StdoutLogger) Debug(args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller}, args...)
	l.debugLogger.Println(args...)
}

func (l StdoutLogger) Debugf(requestID, format string, args ...interface{}) {
	caller := getCaller(2)
	l.debugLogger.Printf("%s [%s] %s",
		caller, requestID, fmt.Sprintf(format, args...))
}

func (l StdoutLogger) Infof(requestID, format string, args ...interface{}) {
	caller := getCaller(2)
	l.infoLogger.Printf("%s [%s] %s",
		caller, requestID, fmt.Sprintf(format, args...))
}

func (l StdoutLogger) Warn(args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller}, args...)
	l.warnLogger.Println(args...)
}

func (l StdoutLogger) Warnf(requestID, format string, args ...interface{}) {
	caller := getCaller(2)
	l.warnLogger.Printf("%s [%s] %s",
		caller, requestID, fmt.Sprintf(format, args...))
}

func (l StdoutLogger) Error(args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller}, args...)
	l.errorLogger.Println(args...)
}

func (l StdoutLogger) Errorf(requestID, format string, args ...interface{}) {
	caller := getCaller(2)
	l.errorLogger.Printf("%s [%s] %s",
		caller, requestID, fmt.Sprintf(format, args...))
}

func (l StdoutLogger) Info1(requestId string, args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller, "["+requestId+"]"}, args...)
	l.infoLogger.Println(args...)
}
func (l StdoutLogger) Error1(requestId string, args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller, "["+requestId+"]"}, args...)
	l.errorLogger.Println(args...)
}
func (l StdoutLogger) Warn1(requestId string, args ...interface{}) {
	caller := getCaller(2)
	args = append([]interface{}{caller, "["+requestId+"]"}, args...)
	l.warnLogger.Println(args...)
}

func (l StdoutLogger) XOLog(sql string, args ...interface{}) {
	caller := getCaller(3)
	args = append([]interface{}{caller, sql + ";", "arguments:"}, args...)
	l.infoLogger.Println(args...)
}

func (l StdoutLogger) Close() {
	_ = l.w.Close()
}

func init() {
	err := os.MkdirAll("/var/log/immortality", os.ModePerm)
	if err != nil {
		panic("failed to create log dir " + "/var/log/immortality")
	}
	w, err := os.OpenFile("/var/log/immortality/immortality.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file " + "/var/log/immortality/immortality.log")
	}

	Log = NewStdoutLogger(w, "")

	db_model.XOLog = Log.XOLog
}
