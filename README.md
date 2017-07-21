# Scheduler

## Migrations
- get migrations cli:
    `go get -u -d github.com/mattes/migrate/cli`
    `go build -tags 'postgres' -o $GOPATH/bin/migrate github.com/mattes/migrate/cli`
- run migraions:
    `migrate -source file:///$GOPATH/src/github.com/bolsunovskyi/scheduler/_migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up`
