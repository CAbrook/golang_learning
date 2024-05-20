package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	LevelError = iota
	LevelWarning
	LevelInfo
)

var (
	errorLogger   = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	infoLogger    = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	currentLevel  = LevelInfo
)

func SetLevel(level int) {
	currentLevel = level
}

func Error(v ...interface{}) {
	if currentLevel >= LevelError {
		errorLogger.Println(v...)
	}
}

func Warning(v ...interface{}) {
	if currentLevel >= LevelWarning {
		warningLogger.Println(v...)
	}
}

func Info(v ...interface{}) {
	if currentLevel >= LevelInfo {
		infoLogger.Println(v...)
	}
}

func LogError[T any](err error, context T) bool {
	if err != nil {
		funcName, file, line := getCallerInfo()
		Error(fmt.Sprintf("Error: %v, Function: %s, File: %s, Line: %d, Context: %v", err, funcName, file, line, context))
		return true
	}
	return false
}

func getCallerInfo() (string, string, int) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", "unknown", 0
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", file, line
	}
	return fn.Name(), file, line
}
