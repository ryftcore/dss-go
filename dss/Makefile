.PHONY: build vet test race check

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
