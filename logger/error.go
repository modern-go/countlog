package logger

import (
	"log"
	"os"
)

var ErrorLogger = log.New(os.Stderr, "", 0)