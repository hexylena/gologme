GOOS ?= linux
GOARCH ?= amd64
SRC := $(wildcard *.go)

# TODO
all: gologme_server gologme_client

gologme_server: bin client server store types util ulogme
	go build github.com/erasche/gologme/bin/gologme_server/

gologme_client: bin client server store types util ulogme
	go build github.com/erasche/gologme/bin/gologme_client/

deps:
	go get github.com/Masterminds/glide/...
	go install github.com/Masterminds/glide/...
	glide install

gofmt:
	goimports $$(find . -type f -name '*.go' -not -path "./vendor/*")
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

qc_deps:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update

qc:
	gometalinter --cyclo-over=10 --deadline=30s --vendor --json ./... > report.json

test: $(SRC) deps gofmt
	go test -v $$(glide novendor)
