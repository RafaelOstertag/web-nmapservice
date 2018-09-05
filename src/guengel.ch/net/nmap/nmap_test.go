package nmap

import (
	"os"
	"reflect"
	"testing"
)

func Test_validatePortSpec(t *testing.T) {
	type args struct {
		portSpec string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"single port", args{"12"}, false},
		{"port range", args{"12-2"}, false},
		{"port list", args{"12,13"}, false},
		{"mixed spec", args{"12,13-25,2,4"}, false},
		{"mixed spec w/ starting range", args{"13,25,12,2"}, false},
		{"invalid spec", args{"ssh,telnet"}, true},
		{"empty spec", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePortSpec(tt.args.portSpec); (err != nil) != tt.wantErr {
				t.Errorf("validatePortSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"host name", args{"www.example.org"}, false},
		{"ip v4 address", args{"10.0.0.1"}, false},
		{"ip v6 address", args{"::1"}, true},
		{"invalid ip", args{"10.0"}, true},
		{"short host", args{"www.example"}, false},
		{"invalid host", args{"example"}, true},
		{"invalid host", args{"ch."}, true},
		{"invalid host", args{".ch"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("validateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRun requires to be run with mocks/nmap.sh
func TestRun(t *testing.T) {
	goPath := os.Getenv("GOPATH")
	os.Setenv("NMAP_CMD", goPath+"/mocks/nmap.sh")

	type args struct {
		host     string
		portspec string
	}
	tests := []struct {
		name    string
		args    args
		want    *Result
		wantErr bool
	}{
		{
			"run nmap",
			args{"gizmo.kruemel.home", "22"},
			&Result{"up", []string{"192.168.100.1"}, []string{"gizmo.kruemel.home"}, []Port{{Name: "ssh", Number: 22, State: "open"}}},
			false,
		},
		{
			"run against mock with invalid parameters",
			args{"gizmo.kruemel.home", ""},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.args.host, tt.args.portspec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
