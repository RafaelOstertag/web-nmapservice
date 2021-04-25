package service

import (
	"encoding/xml"
	"log"
	"strconv"
)

type xmlResult struct {
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

// Result holds the result of the Nmap scan in a form suitable for JSON output.
type Result struct {
	State     string
	Addresses []string
	Hostnames []string
	Ports     []Port
}

// ParseResult reads the XML result from Nmap
func (r *xmlResult) ParseResult(nmapResult []byte) (e error) {
	e = xml.Unmarshal(nmapResult, r)
	return
}

// ToResult returns the XMLResult as Result
func (r *xmlResult) ToResult() *Result {
	result := new(Result)

	result.State = r.Status.State
	result.Hostnames = make([]string, len(r.Hostnames))
	for i, hostname := range r.Hostnames {
		result.Hostnames[i] = hostname.Hostname
	}

	result.Addresses = make([]string, len(r.Addresses))
	for i, address := range r.Addresses {
		result.Addresses[i] = address.Address
	}

	result.Ports = make([]Port, len(r.Ports))
	for i, port := range r.Ports {
		result.Ports[i].Name = port.Service.Name

		var err error
		if result.Ports[i].Number, err = strconv.Atoi(port.Portid); err != nil {
			result.Ports[i].Number = -1
			log.Printf("Error converting port number: %v", err)
		}

		result.Ports[i].State = port.State.State
	}

	return result
}
