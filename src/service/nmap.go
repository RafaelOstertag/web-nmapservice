package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

const (
	defaultNmapCommand   = "nmap"
	portSpecRegexpString = "^(\\d+|\\d+-\\d+|,)*$"
	hostRegexString      = "^([\\da-zA-Z-]+\\.)+([a-zA-Z]+)$|^([\\d]{1,3}\\.){3}[\\d]+$"
)

var (
	portSpecRegexp = regexp.MustCompile(portSpecRegexpString)
	hostRegex      = regexp.MustCompile(hostRegexString)
)

// PortSpecError indicates an invalid port specification
type PortSpecError string

func (t PortSpecError) Error() string {
	return string(t)
}

// HostSpecError indicates an invalid host specification
type HostSpecError string

func (t HostSpecError) Error() string {
	return string(t)
}

// ScannerError indicates a scanner issue
type ScannerError string

func (t ScannerError) Error() string {
	return string(t)
}

// Run nmap against host using portspec. Portspec may only contain digits, `-`, and `,` or be empty.
func Run(host, portspec string) (*Result, error) {
	nmapCommand := getNmapCommand()

	if err := validatePortSpec(portspec); err != nil {
		log.Printf("Portspec invalid: %v", err)
		return nil, PortSpecError("Invalid portspec")
	}

	if err := validateHost(host); err != nil {
		log.Printf("Host invalid: %v", err)
		return nil, HostSpecError("Invalid host")
	}

	var cmd *exec.Cmd
	if portspec == "" {
		cmd = exec.Command(nmapCommand, "-oX", "-", host)
	} else {
		cmd = exec.Command(nmapCommand, "-oX", "-", "-p", portspec, host)
	}

	var output []byte
	var err error
	if output, err = cmd.Output(); err != nil {
		log.Printf("Error running nmap: %v", err)
		if _, ok := err.(*exec.ExitError); ok == true {
			return nil, ScannerError("Error invoking scanner")
		}
		return nil, ScannerError(err.Error())

	}

	xmlResult := new(xmlResult)

	if err = xmlResult.ParseResult(output); err != nil {
		log.Printf("Error digesting nmap output: %v", err)
		return nil, err
	}

	return xmlResult.ToResult(), nil
}

func getNmapCommand() string {
	nmapCommandFromEnv := os.Getenv("NMAP_CMD")
	if nmapCommandFromEnv == "" {
		return defaultNmapCommand
	}

	return nmapCommandFromEnv
}

func validatePortSpec(portSpec string) error {
	if portSpecRegexp.MatchString(portSpec) {
		return nil
	}

	errorMsg := fmt.Sprintf("'%s' is an invalid portspec. It does not match the regexp '%s'", portSpec, portSpecRegexpString)
	return errors.New(errorMsg)
}

func validateHost(host string) error {
	if hostRegex.MatchString(host) {
		return nil
	}

	errorMsg := fmt.Sprintf("'%s' is an invalid host. It does not match the regexp '%s'", host, hostRegexString)
	return errors.New(errorMsg)
}
