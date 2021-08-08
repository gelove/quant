package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"quant/pkg/utils/date"
	"runtime"
	"time"
)

type Level int

var (
	F      *os.File
	logger *log.Logger

	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func MustLoad(logPath, logFileExt string) {
	var err error
	filePath := getLogFilePath(logPath)
	fileName := getLogFileName(logFileExt)
	F, err = MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
}

func DebugF(format string, v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(format+"\n", v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
}

func InfoF(format string, v ...interface{}) {
	setPrefix(INFO)
	logger.Printf(format+"\n", v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
}

func WarnF(format string, v ...interface{}) {
	setPrefix(WARNING)
	logger.Printf(format+"\n", v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
}

func ErrorF(format string, v ...interface{}) {
	setPrefix(ERROR)
	logger.Printf(format+"\n", v...)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Fatalln(v...)
}

func FatalF(format string, v ...interface{}) {
	setPrefix(FATAL)
	logger.Printf(format+"\n", v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}

func getCurrentDay() string {
	return time.Now().Format(date.YYYY_MM_DD)
}

func getLogFilePath(logPath string) string {
	return fmt.Sprintf("%s", logPath)
}

func getLogFileName(logFileExt string) string {
	return fmt.Sprintf("%s%s", getCurrentDay(), logFileExt)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
