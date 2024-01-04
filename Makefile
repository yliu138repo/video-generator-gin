server: doc
	go run main.go

dev-server:
	gin run main.go

install:
	go mod download

doc:
	swag init

backup-env:
	cp pkg/common/envs/.env pkg/common/envs/.env.sample