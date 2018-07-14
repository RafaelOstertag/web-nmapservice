package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
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

func getListenAddress() string {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		return ":8081"
	}
	return listen
}

func main() {
	router := routing.New()

	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	api := router.Group("/v1")
	api.Use(
		content.TypeNegotiator(content.JSON),
	)

	api.Get("/scan/<host:[a-zA-Z0-9.-]+>/<portSpec:[\\d,-]+>", handleScanRequest)
	api.Get("/scan/<host:[a-zA-Z0-9.-]+>", handleScanRequest)

	http.Handle("/", router)
	log.Print("Starting server")
	log.Fatal(http.ListenAndServe(getListenAddress(), nil))
}
