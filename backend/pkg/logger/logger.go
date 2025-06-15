package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var logLevel = DEBUG

func init() {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		switch level {
		case "DEBUG":
			logLevel = DEBUG
		case "INFO":
			logLevel = INFO
		case "WARNING":
			logLevel = WARNING
		case "ERROR":
			logLevel = ERROR
		case "FATAL":
			logLevel = FATAL
		default:
			fmt.Printf("Unknown log level: %s\n", level)
			logLevel = DEBUG
		}
	}
	fmt.Printf("Log level set to: %d\n", logLevel)
}

// Debug level log
func Debugf(message string, args ...interface{}) {
	if logLevel > DEBUG{
		return
	}
	formatted := fmt.Sprintf(message, args...)
	fmt.Printf("[DEBUG] %s: %s\n", time.Now().Format(time.DateTime), formatted)
}

// Info level log
func Infof(message string, args ...interface{}) {
	if logLevel > INFO {
		return
	}
	formatted := fmt.Sprintf(message, args...)
	fmt.Printf("[INFO] %s: %s\n", time.Now().Format(time.DateTime), formatted)
}

// Warning level log
func Warningf(message string, args ...interface{}) {
	if logLevel > WARNING {
		return
	}
	formatted := fmt.Sprintf(message, args...)
	fmt.Printf("[WARNING] %s: %s\n", time.Now().Format(time.DateTime), formatted)
}

// Error level log
func Errorf(message string, trace bool, args ...interface{}) {
	if logLevel > ERROR {
		return
	}
	formatted := fmt.Sprintf(message, args...)
	if !trace {
		fmt.Printf("[ERROR] %s: %s\n", time.Now().Format(time.DateTime), formatted)
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("[ERROR] %s: %s\n", time.Now().Format(time.DateTime), formatted)
		return
	}

	_, filename := filepath.Split(file)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("[ERROR] %s (%s:%d %s): %s\n",
		time.Now().Format(time.DateTime),
		filename, line, funcName, formatted)
}

// Fatal level log
func Fatalf(message string, exit bool, args ...interface{}) {
	if logLevel > FATAL {
		return
	}
	formatted := fmt.Sprintf(message, args...)
	fmt.Printf("[FATAL] %s: %s\n", time.Now().Format(time.DateTime), formatted)
	debug.PrintStack()
	if exit {
		os.Exit(1)
	}
}