_svr_name = mcprouter
_build_time = $(shell date +%y%m%d%H%M%S) 
_image_name = ${_svr_name}:v${_build_time}

build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${_svr_name}

docker-build:build-linux-amd64
	docker build -f Dockerfile -t ${_svr_name} .

dev:
	air -c .air.toml server

tidy:
	go mod tidy

