SOURCE:=client server store types util ulogme loggers
# TODO
all: gologme_server gologme_client

gologme_server: bin/gologme_server $(SOURCE) server/bindata.go
	go build github.com/erasche/gologme/bin/gologme_server/

gologme_client: bin/gologme_client $(SOURCE) server/bindata.go
	go build github.com/erasche/gologme/bin/gologme_client/

deps:
	go get github.com/Masterminds/glide/...
	go install github.com/Masterminds/glide/...
	glide install

gofmt:
	goimports -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

qc_deps:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update

qc:
	gometalinter --cyclo-over=10 --deadline=30s --vendor --json ./... > report.json

test: deps gofmt
	go test -v $$(glide novendor)

server/bindata.go: server/static/*
	go-bindata-assetfs -pkg server -debug server/static/*
	mv bindata_assetfs.go server/bindata.go
