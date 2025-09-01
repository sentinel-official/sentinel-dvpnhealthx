.PHONY: install
install:
	go build -o "${GOPATH}/bin/sentinel-dvpnhealthx" main.go
