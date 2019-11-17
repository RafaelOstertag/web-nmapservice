package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"guengel.ch/nmapservice/service"

	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"

	gnms "github.com/RafaelOstertag/grpcnmapservice"
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
		MaxSize:    1, // megabytes
		MaxBackups: 10,
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

func getServiceCoordinates() (string, int) {
	myListenAddress := getListenAddress()
	components := strings.Split(myListenAddress, ":")

	var myAddress string
	if components[0] == "" {
		myAddress = service.GetOutboundIP()
	} else {
		myAddress = components[0]
	}

	myPort, err := strconv.Atoi(components[1])
	if err != nil {
		log.Panicf("Error getting listening port: %v", err)
	}

	return myAddress, myPort
}

func main() {
	setUpLogging()

	host, port := getServiceCoordinates()

	err := service.Register(host, port)
	if err != nil {
		log.Printf("Error during service registration: %v", err)
	}

	var listenAddress = getListenAddress()
	log.Printf("Starting server on %s", listenAddress)

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	gnms.RegisterNmapServer(grpcServer, &service.NmapService{})
	gnms.RegisterHealthServer(grpcServer, &service.HealthService{Health: make(chan gnms.HealthCheckResponse_ServingStatus)})
	grpcServer.Serve(lis)
}
