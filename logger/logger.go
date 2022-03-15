package logger

import "log"

var level = INFO

const (
	DEBUG = 0
	INFO  = 1
	WARN  = 2
	ERROR = 3
)

func SetLevel(lv int) {
	level = lv
}

func Debug(msg string) {
	if DEBUG < level {
		return
	}
	log.Println("Debug", msg)
}

func Debugf(format string, v ...interface{}) {
	if DEBUG < level {
		return
	}
	log.Printf("Debug "+format, v...)
}

func Info(msg string) {
	if INFO < level {
		return
	}
	log.Println("Info", msg)
}

func Infof(format string, v ...interface{}) {
	if INFO < level {
		return
	}
	log.Printf("Info "+format, v...)
}

func Warn(msg string) {
	if WARN < level {
		return
	}
	log.Println("Warn", msg)
}

func Warnf(format string, v ...interface{}) {
	if WARN < level {
		return
	}
	log.Printf("Warn "+format, v...)
}

func Error(msg string) {
	log.Println("Error", msg)
}
func Errorf(format string, v ...interface{}) {
	log.Printf("Error "+format, v...)
}
