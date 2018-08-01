package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"guengel.ch/services/nmapservice/router" 

	"gopkg.in/natefinch/lumberjack.v2"
)

func getListenAddress() string {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		return ":8081"
	}
	return listen
}

func setUpLogging() {
	logpath, isLogPathSet := os.LookupEnv("LOGPATH")
	if isLogPathSet == false {
		log.Print("Logging to stdout")
		return
	}

	lj := lumberjack.Logger{
		Filename:   logpath + "/nmapservice.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	log.SetOutput(&lj)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			lj.Rotate()
			log.Print("Rotated log")
		}
	}()
}

func main() {
	setUpLogging()

	http.Handle("/", router.ApplicationRouting())
	var listenAddress = getListenAddress()
	log.Printf("Starting server on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
