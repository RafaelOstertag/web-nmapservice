package nmap

import (
	"encoding/xml"
	"log"
	"strconv"
)

// XMLResult holds the result of the scan.
type XMLResult struct {
	Status struct {
		State string `xml:"state,attr"`
	} `xml:"host>status"`
	Addresses []struct {
		Address string `xml:"addr,attr"`
	} `xml:"host>address"`
	Hostnames []struct {
		Hostname string `xml:"name,attr"`
	} `xml:"host>hostnames>hostname"`
	Ports []struct {
		Portid string `xml:"portid,attr"`
		State  struct {
			State string `xml:"state,attr"`
		} `xml:"state"`
		Service struct {
			Name string `xml:"name,attr"`
		} `xml:"service"`
	} `xml:"host>ports>port"`
}

// Port holds the port state
type Port struct {
	Number int
	State  string
	Name   string
}

// JSONResult holds the result of the Nmap scan in a form suitable for JSON output.
type JSONResult struct {
	State     string
	Addresses []string
	Hostnames []string
	Ports     []Port
}

// ParseResult reads the XML result from Nmap
func (r *XMLResult) ParseResult(nmapResult []byte) (e error) {
	e = xml.Unmarshal(nmapResult, r)
	return
}

// ToJSONResult returns the XMLResult as JSONResult
func (r *XMLResult) ToJSONResult() *JSONResult {
	jsonResult := new(JSONResult)

	jsonResult.State = r.Status.State
	jsonResult.Hostnames = make([]string, len(r.Hostnames))
	for i, hostname := range r.Hostnames {
		jsonResult.Hostnames[i] = hostname.Hostname
	}

	jsonResult.Addresses = make([]string, len(r.Addresses))
	for i, address := range r.Addresses {
		jsonResult.Addresses[i] = address.Address
	}

	jsonResult.Ports = make([]Port, len(r.Ports))
	for i, port := range r.Ports {
		jsonResult.Ports[i].Name = port.Service.Name

		var err error
		if jsonResult.Ports[i].Number, err = strconv.Atoi(port.Portid); err != nil {
			jsonResult.Ports[i].Number = -1
			log.Printf("Error converting port number: %v", err)
		}

		jsonResult.Ports[i].State = port.State.State
	}

	return jsonResult
}
