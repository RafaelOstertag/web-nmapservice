GOENV = env GOPATH=`pwd`
GOCOV = bin/gocov

DEPENDENCIES = 

all: tests coverage nmapservice

tests:
	$(GOENV) go test -v guengel.ch/nmapservice

$(GOCOV):
	$(GOENV) go get github.com/axw/gocov/gocov


coverage: $(GOCOV)
	$(GOENV) $(GOCOV) test guengel.ch/nmapservice/service | $(GOCOV) report


clean:
	rm -rf bin pkg
	rm -rf src/github.com src/golang.org

nmapservice:
	$(GOENV) go get -v guengel.ch/nmapservice
