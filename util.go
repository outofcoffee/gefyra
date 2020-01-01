package main

import "log"

func fatalIfError(err interface{}, message string) {
	if err != nil {
		log.Fatal(message, " - ", err)
	}
}

func signalFailure(control chan int) {
	control <- 1
}
