# Gym Backend Core

Microservice responsible for main backend services.

![Go](https://img.shields.io/badge/Golang-1.18-blue.svg?logo=go&longCache=true&style=flat)
![Postgres](https://img.shields.io/badge/Postgres-14.2-lightblue.svg?logo=postgresql&longCache=true&style=flat)

## Getting Started

This project uses the **Go** programming language (Golang) and PostgreSQL relational database.

## Steps to run

Clone the repository and enter the folder
```bash
https://gitlab.com/gym-global/backend.git
cd backend
```

### Installing
#### Using GOMODULE

```bash
go mod download
```

Compose the local environment
```bash
make local
```

## Running the tests

```bash
go test ./...
```
#### Creating a migration
```bash
migrate create -ext sql -dir migration -seq some_migration_name
```


