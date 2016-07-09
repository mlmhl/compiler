package clog

import (
	"io"
	"os"

	gio "github.com/mlmhl/goutil/io"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"fmt"
)

type Logger struct {
	output io.Writer
}

func NewLogger(outputFile string) (*Logger, error) {
	var err error
	var output io.Writer

	if outputFile == "stdout" {
		output = os.Stdout
	} else if outputFile == "stderr" {
		output = os.Stderr
	} else {
		if gio.IsExist(outputFile) {
			os.Remove(outputFile)
		}
		if output, err = os.Create(outputFile); err != nil {
			return nil, err
		}
	}

	return &Logger{output}, nil
}

// log a internal error
func (logger *Logger) InternalError(err gerror.Error) {
	logger.output.Write([]byte(err.GetMessage() + "\n"))
	os.Exit(1)
}

// log a compile error
func (logger *Logger) CompileError(err gerror.Error) {
	logger.logError(err)
}

// log a runtime error
func (logger *Logger) RuntimeError(err gerror.Error) {
	logger.logError(err)
}

func (logger *Logger) logError(err gerror.Error) {
	location := err.GetLocation()
	logger.output.Write([]byte(fmt.Sprintf("%s,%d,%d: %s\n", location.GetFileName(),
		location.GetLine(), location.GetPosition(), err.GetMessage())))
	os.Exit(1)
}