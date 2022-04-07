.PHONY : build local stop-local prod down test network prune clean

default: build

build:
	@echo '### BUILDING GO BINARY'
	@go build -o api cmd/api/main.go

local:
	@echo '### STARTING LOCAL DOCKER-COMPOSE'
	@cp .env.local .env;
	@docker-compose -f docker-compose-local.yml up -d --build

stop-local:
	@echo '### STOPPING LOCAL DOCKER-COMPOSE'
	@docker-compose -f docker-compose-local.yml down --remove-orphans

dev:
	@echo '### LOADING DEV ENV'
	@cp .env.dev .env;
	@docker-compose -f docker-compose-dev.yml up -d --build

stop-dev:
	@echo '### STOPPING DEV DOCKER-COMPOSE'
	@docker-compose -f docker-compose-dev.yml down --remove-orphans

test:
	@go test -v ./...

api-test:
	@go test -v ./cmd/api

network:
	@docker network create gym-net

prune:
	@echo '### PRUNING DOCKER SYSTEM'
	@docker system prune -a

clean:
	@docker rm -vf $$(docker ps -aq)
	@docker rmi -f $$(docker images -a -q)
