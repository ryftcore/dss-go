.PHONY: build vet test race check cli-build

build:
	go build ./...

vet:
	go vet ./...

test:
	go test ./... -timeout 40m

# dss/asic's merge cross-validation exceeds the default 10m under -race.
race:
	go test ./... -race -timeout 60m

check: build vet test

# cli-build builds the esig binary into bin/, named for the host OS/arch so
# repeat builds on different platforms do not clobber each other.
cli-build:
	mkdir -p bin
	go build -o bin/esig$(shell go env GOEXE) ./cmd/esig
