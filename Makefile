.PHONY : build local stop-local prod down test network prune clean

default: build

build:
	@echo '### BUILDING GO BINARY'
	@go build -o api cmd/api/main.go

local:
	@echo '### STARTING LOCAL DOCKER-COMPOSE'
	@cp .env.local .env;
	@docker-compose -f docker-compose.yml up -d

stop-local:
	@echo '### STOPPING LOCAL DOCKER-COMPOSE'
	@docker-compose -f docker-compose.yml down

prod:
	@echo '### LOADING PRODUCTION ENV'
	@cp .env.prod .env;

down:
	@echo '### DOCKER COMPOSE DOWN'
	@docker-compose down --remove-orphans

test:
	@go test -v ./...

network:
	@docker network create gym-net

prune:
	@echo '### PRUNING DOCKER SYSTEM'
	@docker system prune -a

clean:
	@docker rm -vf $$(docker ps -aq)
	@docker rmi -f $$(docker images -a -q)
