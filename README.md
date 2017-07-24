# Scheduler

## Migrations
- get migrations cli:
    `go get -u -d github.com/mattes/migrate/cli`
    `go build -tags 'postgres' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli`
- run migraions:
    `migrate -source file:///$GOPATH/src/github.com/bolsunovskyi/scheduler/_migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up`

## Postgres
run postgres:
    `docker run -p 5434:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres postgres`
    
## Build plugin
    `go build -buildmode=plugin`