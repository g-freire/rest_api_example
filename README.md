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

### Installing
#### Using GOMODULE

```bash
go mod download
```

Compose the local environment (cache + db)
```bash
make local
make stop-local
```

Compose the dev environment (app + cache + db)
```bash
make dev
make stop-dev
```

Run the app locally
```bash
go run cmd/api/main.go
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


