package syslog

import (
	"fmt"
	"log"
	"os"
)

// log file path
var fileName string = "gdync.log"
var logger *log.Logger

const (
	logLevel = log.LstdFlags | log.Lshortfile | log.Lmicroseconds
)

// log level tag
const (
	dEBUG = "DEBUG"
	iNFO  = "INFO"
	wARN  = "WARN"
	eRROR = "ERROR"
	fATAL = "FATAL"
)

func init() {
	var file *os.File
	var err error
	if isExist(fileName) {
		file, err = os.Open(fileName)
	} else {
		file, err = os.Create(fileName)
	}

	if err != nil {
		logger = log.New(os.Stdout, "", logLevel)
	} else {
		logger = log.New(file, "", logLevel)
	}
}

// Check whether the file is exist or not
func isExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func output(tag string, v ...interface{}) {
	logger.Output(3, "["+tag+"] "+fmt.Sprintln(v...))
}

func outputf(tag, format string, v ...interface{}) {
	logger.Output(3, "["+tag+"] "+fmt.Sprintf(format, v))
}

func Debug(v ...interface{}) {
	output(dEBUG, v)
}

func Debugf(format string, v ...interface{}) {
	outputf(dEBUG, format, v)
}

func Info(v ...interface{}) {
	output(iNFO, v)
}

func Infof(format string, v ...interface{}) {
	outputf(iNFO, format, v)
}

func Warn(v ...interface{}) {
	output(wARN, v)
}

func Warnf(format string, v ...interface{}) {
	outputf(wARN, format, v)
}

func Error(v ...interface{}) {
	output(eRROR, v)
}

func Errorf(format string, v ...interface{}) {
	outputf(eRROR, format, v)
}

// Fatal causes the current program to exit with status code 1.
func Fatal(v ...interface{}) {
	output(fATAL, v)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	outputf(fATAL, format, v)
	os.Exit(1)
}