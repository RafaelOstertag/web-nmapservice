GO = go
GOCOV = $${HOME}/go/bin/gocov

all: nmapservice tests coverage

tests:
	$(GO) test -v guengel.ch/nmapservice/service

$(GOCOV):
	$(GO) get github.com/axw/gocov/gocov

coverage: $(GOCOV)
	$(GOCOV) test guengel.ch/nmapservice/service | $(GOCOV) report

clean:
	 rm -f nmapservice

nmapservice:
	$(GO) build -v guengel.ch/nmapservice
