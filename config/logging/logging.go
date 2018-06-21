package logging

import (
	"io"
	"log"
)

var (
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
)

func Debug(v ...interface{}) {
	debug.Println(v)
}

func Error(v ...interface{}) {
	error.Println(v)
}

func Info(v ...interface{}) {
	info.Println(v)
}

func InitLogging(
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
