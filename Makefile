server:
	go run cmd/main.go

dev-server:
	gin run cmd/main.go

install:
	go mod download

backup-env:
	cp pkg/common/envs/.env pkg/common/envs/.env.sample