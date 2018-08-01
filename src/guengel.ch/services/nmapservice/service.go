package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"gopkg.in/natefinch/lumberjack.v2"
	"guengel.ch/net/nmap"
)

func handleScanRequest(c *routing.Context) error {
	var err error
	var result *nmap.Result

	if result, err = nmap.Run(c.Param("host"), c.Param("portSpec")); err != nil {
		return err
	}

	return c.Write(result)
}

func handleHealthRequest(c *routing.Context) error {
	type HealthStruct struct {
		Status string `json:"status"`
	}

	health := HealthStruct{"ok"}

	return c.Write(health)
}

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

	router := routing.New()

	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		content.TypeNegotiator(content.JSON),
	)

	router.Get("/health", handleHealthRequest)

	api := router.Group("/v1")

	api.Get("/scan/<host:[a-zA-Z0-9.-]+>/<portSpec:[\\d,-]+>", handleScanRequest)
	api.Get("/scan/<host:[a-zA-Z0-9.-]+>", handleScanRequest)

	http.Handle("/", router)
	var listenAddress = getListenAddress()
	log.Printf("Starting server on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
