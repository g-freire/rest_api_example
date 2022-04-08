# Gym Backend Core

Microservice responsible for main backend services.

![Go](https://img.shields.io/badge/go-1.18-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-14.2-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-6.2-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)

## Getting Started

This project uses the **Go** programming language (Golang), **PostgreSQL** as the relational database and **Redis** as the cache.

## Steps to run

Clone the repository and enter the folder
```bash
git clone https://gitlab.com/gym-global/backend.git && cd backend
```
##### At docker compose network
Compose the dev environment (app + cache + db)
```bash
make dev
make stop-dev
```
Test the app
```bash
curl --location --request GET 'http://localhost:6000/health'
curl --location --request GET 'http://localhost:6000/'
```

##### Locally - debug mode
Compose the local environment (cache + db)
```bash
make local
make stop-local
```
Install, Build and Run Go binary
```bash
go run cmd/api/main.go or
go mod download && go build -o gym cmd/api/main.go && ./gym
```
Test the app
```bash
curl --location --request GET 'http://localhost:8080/health'
curl --location --request GET 'http://localhost:8080/'
```

### Running the tests

```bash
go test ./...
or 
make test
make api-test
```
#### Creating a migration
```bash
migrate create -ext sql -dir migration -seq some_migration_name
```


