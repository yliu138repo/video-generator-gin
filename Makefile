server: doc
	go run main.go

build: doc cp-env-sample
	mkdir -p build
	rm -rf build/*
	go build -o build/video-generator-api

build-legacy: doc cp-env-sample
	mkdir -p build
	rm -rf build/*
	CGO_ENABLED=0 go build -o build/video-generator-api

dev-server:
	gin run main.go

install:
	go mod download

doc:
	swag init

backup-env:
	cp .env .env.sample

cp-env-sample:
	cp .env.sample .env