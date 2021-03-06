package main

import (
	"fmt"
	"log"
	"os"
)

func SetLogOutputToFile() {
	logfile, err = os.OpenFile("logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(logfile)
	log.Println("----------------Log start----------------")
}

func CloseLogFile() {
	log.Println("----------------Log end------------------")
	err := logfile.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func SimpleFatalErrorHandler(err error) {
	if err != nil {
		FatalError(err)
	}
}

func FatalError(v ...interface{}) {
	log.Println(v...)
	CloseLogFile()
	os.Exit(-1)
}

func EndProgramWithMessage(message string) {
	fmt.Println(message)
	fmt.Println("Please configure config.cfg")
	fmt.Println("Press the Enter Key to terminate the program")
	_, _ = fmt.Scanln() // wait for Enter Key
	FatalError()
}
