package main

import (
	"flag"

	"github.com/mlmhl/compiler/gdync/interpreter"
)

func main() {
	var fileName = flag.String("fileName", "test", "file name")
	flag.Parse()

	inter := interpreter.NewInterpreter()
	inter.Interpret(*fileName)
}