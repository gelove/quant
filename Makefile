# note: 注意文件名大小写 Makefile
.PHONY: all start build wire

NOW = $(shell date '+%FT%T')

MOMENT_MAC = $(shell date -v-7d  +%FT%T)

MOMENT_LINUX = $(shell date -d '-7 days' +%FT%T)

COUNT = 20

APP = quant

all: start

build:
	@go build -a -v -ldflags "-s -w" -o $(APP) . && upx $(APP)

build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -ldflags "-s -w" -o $(APP) . && upx $(APP)

start:
	go run -race main.go start

bounce:
	go run main.go bounce $(MOMENT_MAC)

swing:
	go run main.go swing -n $(COUNT)

swagger:
	swag init --generalInfo ./internal/app/swagger.go --output ./internal/app/swagger

wire:
	wire gen ./internal/app/initial

test:
	@go test -v ./... -coverprofile coverage.txt

test-ci:
	@CI=true go test -v ./...

clean:
	go clean
	rm -rf ./$(APP)
