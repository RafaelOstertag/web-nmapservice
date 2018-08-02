package router

import (
	"log"
	"net/http"

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
		if _, ok := err.(nmap.HostSpecError); ok == true {
			return routing.NewHTTPError(400, err.Error())
		}

		if _, ok := err.(nmap.PortSpecError); ok == true {
			return routing.NewHTTPError(400, err.Error())
		}

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

// ApplicationRouting sets up the router for the application
func ApplicationRouting() *routing.Router {
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

	return router
}
