package logger

import "log"

func Debug(msg string) {
	log.Println("Debug", msg)
}

func Debugf(format string, v ...interface{}) {
	log.Printf("Debug "+format, v)
}
func Info(msg string) {
	log.Println("Info", msg)
}
func Warn(msg string) {
	log.Println("Warn", msg)
}
func Error(msg string) {
	log.Println("Error", msg)
}
