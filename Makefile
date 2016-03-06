GOOS ?= linux
GOARCH ?= amd64
SRC := $(wildcard *.go)

# TODO
all:
	go build github.com/erasche/gologme/bin/gologme_client/
	go build github.com/erasche/gologme/bin/gologme_server/

deps:
	go get github.com/Masterminds/glide/...
	go install github.com/Masterminds/glide/...
	glide install

gofmt:
	find $(glide novendor) -name '*.go' -exec gofmt -w '{}' \;

qc_deps:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update

qc: qc_deps
	gometalinter --cyclo-over=10 $(glide novendor)

test: $(SRC) deps gofmt
	go test -v $(glide novendor)
