[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/golang-migrate/migrate/ci.yaml?branch=master)](https://github.com/golang-migrate/migrate/actions/workflows/ci.yaml?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/golang-migrate/migrate)](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
[![Coverage Status](https://img.shields.io/coveralls/github/golang-migrate/migrate/master.svg)](https://coveralls.io/github/golang-migrate/migrate?branch=master)
[![packagecloud.io](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/golang-migrate/migrate?filter=debs)
[![Docker Pulls](https://img.shields.io/docker/pulls/migrate/migrate.svg)](https://hub.docker.com/r/migrate/migrate/)
![Supported Go Versions](https://img.shields.io/badge/Go-1.21%2C%201.22-lightgrey.svg)
[![GitHub Release](https://img.shields.io/github/release/golang-migrate/migrate.svg)](https://github.com/golang-migrate/migrate/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-migrate/migrate/v4)](https://goreportcard.com/report/github.com/golang-migrate/migrate/v4)

# [Backend #16] How to handle DB errors in GOlang correctly

- Use sqlc to generate go code for users table
- other tests failed since users table introduce a new foreign key constraint, need to fix that in `account_test.go`
- api tests fail, `*mockdb.MockStore does not implement db.Store`: since we create two new go functions for users table, `Querier` interface requires the mockdb to implement those two new functions. Need to run `make mock` to regenerate the code for the `MockStore`
- then the HTTP api requests also fail because of the new introduced foreign key constraint,

# [Backend #15] Add users table with unique & foreign key constraints in PostgreSQL

- Add users table to dbdiagram
- instead of overwriting migration up sql, we should create another migration version: `migrate create -ext sql -dir db/migration -seq add_users`
- add migrateup1 and migratedown1 to Makefile

# [Backend #14] Implement transfer money API with a custom params validator

- add `api/transfer.go`
- add `api/validattor.go` for custom currency validator, then register the validator in `server.go`

# [Backend #13] Mock DB for testing HTTP API in Go and achieve 100% coverage

- Mocking db vs using a real db: isolate tests data to avoid conflicts, faster tests, 100% coverage
- How to mock?
  - use fake db: memory: implement a fake version of db: store data in memory
  - use db stubs: GOMOCK: generate and build stubs that returns hard-coded values
- At the moment, `NewStore(store *db.Store)` is accepting a db store object, which always connects to a real db; so in order to use a mock db in the api server tests, we replace `Store` struct object with an interface. Create a `Store` interface, and rename the old `Store` struct to `SQLStore` to be a real implementation of the `Store` interface. Look at [Backend #6] for details.
- create a `/mock` folder
- `mockgen -package mockdb -destination db/mock/store.go github.com/YimingLi-Billy/simplebank/db/sqlc Store` generates `/mock/store.go`
  - `MockStore`: a struct that implements all required functions of the Store interface
  - `MockStoreMockRecorder`: also includes the functions of the same name and the same number of arguments, but the types of the arguments are just general interface types. The idea is we can specify how many times the function should be called, and with what values of the arguments.
- This GoMock code checks the inputs (arguments) and the number of calls, but it doesn't check what the method returns, instead it specifies what the mock should return when the method is called with the expected inputs.
  ```go
    store.EXPECT().
  		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
  		Times(1).
  		Return(account, nil)
  ```

# [Backend #12] Load config from file & environment variables in Golang with Viper

- Install go viper package
- `app.env`: declare environment variables
- `config.go`: declare a struct for loading config values, and declare `LoadConfig()`

# [Backend #11] Implement RESTful HTTP API in Go using Gin

- `Server`: `api/server.go` contains `db.Store` and router `gin.Engine`, the router will help us send each API request to the correct handler for processing
- `account.go` contains all router functions regarding account
- gin validator

# [Backend #10] Setup GithubActions for Golang + Postgres to run automated tests

Skipped

# [Backend #9] Understand isolation levels & read phenomena in MySQL & PostgreSQL via examples

Skipped

# [Backend #8] How to avoid deadlock in DB transaction? Queries order matters!

Skipped since I don't work with DB that much

# [Backend #7] DB transaction lock & How to handle deadlock in Golang

- TDD: test driven development. Write tests first, then fix the code to make the tests pass
- Add money: 1. get account balance; 2. add amount to the current balance.
  - often done incorrectly without proper locking mechanism: GET balannce doesn't block the transaction so other coroutines can read the same number before the updates
  - this lecture demonstrate this behavior with two terminal sessions of postgres, using `SELECT FOR UPDATE` to block GET balance. After using `SELECT FOR UPDATE`, the deadlock occurs.
  - Debug deadlock: print out logs to see which transaction is calling which query and in which order, for that we have to assign a name for each transaction and pass it into the `TransferTx()` via the context argument.
  - Once print out the actions in each transaction, we can try the same actions using two postgres terminals. We find out that terminal 2 `INSERT` entry is blocking terminal 1 `SELECT FOR UPDATE`.
  - google postgres lock for a query to look for blocked queries and what is blocking them. Copy and run that in TablePlus
  - deadlock caused by foreign key constraint

# [Backend #6] A clean way to implement database transaction in Golang

- ACID:
  - Atomicity: Either all operations complete successfully or the transaction fails and the db is unchanged
  - Consistency: The db state must be valid after the transaction. All constraints must be satisfied
  - Isolation: Concurrent transactions must not affect each other
  - Durability: Data written by a successful transaction must be recorded in persistent storage
- `store.go` provides all functions to execute db queries and transactions
- `db`, `Quries`, `SQLStore`, `Querier`, and `Store`:
  - `db`: database
  - `Queries`: a wrapper of `db`, and a bunch of go queries function generated by sqlc
  - `SQLStore`: a wrapper of both `Queries` and `db`, note that `db` and `Queries`'s db are the same `db` (weird)
  - `Querier`: an interface generated by sqlc that implements all `Queries` functions
  - `Store`: a wrapper of `Querier` interface, plus `TransferTx` function, bc `TransferTx` is a transaction, and `Querier` functions are regular database CRUD.
- `execTx` function uses postgres `db.BeginTx` and `tx.Rollback` to execute a transaction (a callback taken in as an input)
- Implement `TransferTx` transaction function
- Implement `store_test.go` unit tests, use Go routine to test transaction
  - channel: connect concurrent Go routines, allow them to safely share data with each other without explicit locking

# [Backend #5] Write Golang unit tests for database CRUD with random data

- Create `main_test.go` at the same place as the go functions
- `TestMain` function in Go will run every test that starts with `Test` in the same package when you run the tests using `go test`. In Go, any function that starts with `Test` and takes a parameter of type `*testing.T` is considered a test function. The Go testing framework looks for functions that follow this naming convention and runs them when `go test` is executed.
- Create `account_test.go` to contain tests for the target functions.
- `make test`

# [Backend #4] Generate CRUD Golang code from SQL | Compare db/sql, gorm, sqlx & sqlc

- sqlc is a tool that gernates type-safe Go code from SQL queries. It helps devs work with databases in a more efficient way by allowing them to write SQL queries directly in their Go code and then automatically generating the necessary Go types and functions that correspond to those queries.
- `sqlc init` creates a `sqlc.yaml`
- Go to sqlc github page and refer to the settings for updating `sqlc.yaml`
- Add queries in db/query
- `sqlc generate` generates `models.go`, `db.go`, and `account.sql.go`:
  - `models.go`: contains the struct definition of 3 models
  - `db.go`: contains DBTX interface, this allows us to freely use either a db or transaction to execute a query, depend on whether we want to execute just 1 sql query or a set of multiple queries within a transaction.
  - `account.sql.go`: contains generated go code.
  - Do not change the generated files, since they will be overwritten when run `sqlc generate`

# [Backend #3] How to write & run database migration in Golang

- Install golang-migrate: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
- `migrate --version`
- Create migration files: `migrate create -ext sql -dir db/migration -seq init_schema`
- Copy .sql file generated by dbdiagram.io to `init_schema.up.sql`
- Add the following to `init_schema.down.sql`:
  - ```sql
      DROP TABLE IF EXISTS entries;
      DROP TABLE IF EXISTS transfers;
      DROP TABLE IF EXISTS accounts;
    ```
- check if postgres12 container is running
- Create Makefile for createdb and dropdb command

# [Backend #2] Install & use Docker + Postgres + TablePlus to create DB schema

Docker:

- `docker pull postgres:12-alpine`
- `docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine`

- `docker start postgres12`
- Access Postgres console with psql: `docker exec -it postgres12 psql -U root`
- docker logs postgres12
- docker stop postgres12

TablePlus:

- connect to the postgres database via port

# migrate

**Database migrations written in Go. Use as [CLI](#cli-usage) or import as [library](#use-in-your-go-project).**

- Migrate reads migrations from [sources](#migration-sources)
  and applies them in correct order to a [database](#databases).
- Drivers are "dumb", migrate glues everything together and makes sure the logic is bulletproof.
  (Keeps the drivers lightweight, too.)
- Database drivers don't assume things or try to correct user input. When in doubt, fail.

Forked from [mattes/migrate](https://github.com/mattes/migrate)

## Databases

Database drivers run migrations. [Add a new database?](database/driver.go)

- [PostgreSQL](database/postgres)
- [PGX v4](database/pgx)
- [PGX v5](database/pgx/v5)
- [Redshift](database/redshift)
- [Ql](database/ql)
- [Cassandra / ScyllaDB](database/cassandra)
- [SQLite](database/sqlite)
- [SQLite3](database/sqlite3) ([todo #165](https://github.com/mattes/migrate/issues/165))
- [SQLCipher](database/sqlcipher)
- [MySQL / MariaDB](database/mysql)
- [Neo4j](database/neo4j)
- [MongoDB](database/mongodb)
- [CrateDB](database/crate) ([todo #170](https://github.com/mattes/migrate/issues/170))
- [Shell](database/shell) ([todo #171](https://github.com/mattes/migrate/issues/171))
- [Google Cloud Spanner](database/spanner)
- [CockroachDB](database/cockroachdb)
- [YugabyteDB](database/yugabytedb)
- [ClickHouse](database/clickhouse)
- [Firebird](database/firebird)
- [MS SQL Server](database/sqlserver)
- [rqlite](database/rqlite)

### Database URLs

Database connection strings are specified via URLs. The URL format is driver dependent but generally has the form: `dbdriver://username:password@host:port/dbname?param1=true&param2=false`

Any [reserved URL characters](https://en.wikipedia.org/wiki/Percent-encoding#Percent-encoding_reserved_characters) need to be escaped. Note, the `%` character also [needs to be escaped](https://en.wikipedia.org/wiki/Percent-encoding#Percent-encoding_the_percent_character)

Explicitly, the following characters need to be escaped:
`!`, `#`, `$`, `%`, `&`, `'`, `(`, `)`, `*`, `+`, `,`, `/`, `:`, `;`, `=`, `?`, `@`, `[`, `]`

It's easiest to always run the URL parts of your DB connection URL (e.g. username, password, etc) through an URL encoder. See the example Python snippets below:

```bash
$ python3 -c 'import urllib.parse; print(urllib.parse.quote(input("String to encode: "), ""))'
String to encode: FAKEpassword!#$%&'()*+,/:;=?@[]
FAKEpassword%21%23%24%25%26%27%28%29%2A%2B%2C%2F%3A%3B%3D%3F%40%5B%5D
$ python2 -c 'import urllib; print urllib.quote(raw_input("String to encode: "), "")'
String to encode: FAKEpassword!#$%&'()*+,/:;=?@[]
FAKEpassword%21%23%24%25%26%27%28%29%2A%2B%2C%2F%3A%3B%3D%3F%40%5B%5D
$
```

## Migration Sources

Source drivers read migrations from local or remote sources. [Add a new source?](source/driver.go)

- [Filesystem](source/file) - read from filesystem
- [io/fs](source/iofs) - read from a Go [io/fs](https://pkg.go.dev/io/fs#FS)
- [Go-Bindata](source/go_bindata) - read from embedded binary data ([jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata))
- [pkger](source/pkger) - read from embedded binary data ([markbates/pkger](https://github.com/markbates/pkger))
- [GitHub](source/github) - read from remote GitHub repositories
- [GitHub Enterprise](source/github_ee) - read from remote GitHub Enterprise repositories
- [Bitbucket](source/bitbucket) - read from remote Bitbucket repositories
- [Gitlab](source/gitlab) - read from remote Gitlab repositories
- [AWS S3](source/aws_s3) - read from Amazon Web Services S3
- [Google Cloud Storage](source/google_cloud_storage) - read from Google Cloud Platform Storage

## CLI usage

- Simple wrapper around this library.
- Handles ctrl+c (SIGINT) gracefully.
- No config search paths, no config files, no magic ENV var injections.

**[CLI Documentation](cmd/migrate)**

### Basic usage

```bash
$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 2
```

### Docker usage

```bash
$ docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://localhost:5432/database up 2
```

## Use in your Go project

- API is stable and frozen for this release (v3 & v4).
- Uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies.
- To help prevent database corruptions, it supports graceful stops via `GracefulStop chan bool`.
- Bring your own logger.
- Uses `io.Reader` streams internally for low memory overhead.
- Thread-safe and no goroutine leaks.

**[Go Documentation](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)**

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {
    m, err := migrate.New(
        "github://mattes:personal-access-token@mattes/migrate_test",
        "postgres://localhost:5432/database?sslmode=enable")
    m.Steps(2)
}
```

Want to use an existing database client?

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "postgres", driver)
    m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}
```

## Getting started

Go to [getting started](GETTING_STARTED.md)

## Tutorials

- [CockroachDB](database/cockroachdb/TUTORIAL.md)
- [PostgreSQL](database/postgres/TUTORIAL.md)

(more tutorials to come)

## Migration files

Each migration has an up and down migration. [Why?](FAQ.md#why-two-separate-files-up-and-down-for-a-migration)

```bash
1481574547_create_users_table.up.sql
1481574547_create_users_table.down.sql
```

[Best practices: How to write migrations.](MIGRATIONS.md)

## Coming from another db migration tool?

Check out [migradaptor](https://github.com/musinit/migradaptor/).
_Note: migradaptor is not affiliated or supported by this project_

## Versions

| Version    | Supported?         | Import                                                                                                                                 | Notes                                        |
| ---------- | ------------------ | -------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| **master** | :white_check_mark: | `import "github.com/golang-migrate/migrate/v4"`                                                                                        | New features and bug fixes arrive here first |
| **v4**     | :white_check_mark: | `import "github.com/golang-migrate/migrate/v4"`                                                                                        | Used for stable releases                     |
| **v3**     | :x:                | `import "github.com/golang-migrate/migrate"` (with package manager) or `import "gopkg.in/golang-migrate/migrate.v3"` (not recommended) | **DO NOT USE** - No longer supported         |

## Development and Contributing

Yes, please! [`Makefile`](Makefile) is your friend,
read the [development guide](CONTRIBUTING.md).

Also have a look at the [FAQ](FAQ.md).

---

Looking for alternatives? [https://awesome-go.com/#database](https://awesome-go.com/#database).
