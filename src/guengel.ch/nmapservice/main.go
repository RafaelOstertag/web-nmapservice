package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"guengel.ch/nmapservice/service"

	"google.golang.org/grpc"
	// this is currently having difficulties to be resolved: "gopkg.in/natefinch/lumberjack.v2"
	"github.com/natefinch/lumberjack"

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

func main() {
	setUpLogging()

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
