package logger

import (
	"log"
	"os"
)

var InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

var ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

var WarnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime)
