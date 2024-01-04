server: doc
	go run main.go

build: doc
	mkdir -p build
	rm -rf build/*
	go build -o build/video-generator-api

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