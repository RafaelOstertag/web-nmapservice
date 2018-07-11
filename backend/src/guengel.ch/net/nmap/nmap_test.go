package nmap

import (
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
		{"empty spec", args{""}, true},
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
		{"invalid host", args{"www.example"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("validateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
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
