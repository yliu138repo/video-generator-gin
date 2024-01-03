server:
	go run main.go

dev-server:
	gin run main.go

install:
	go mod download

backup-env:
	cp pkg/common/envs/.env pkg/common/envs/.env.sample