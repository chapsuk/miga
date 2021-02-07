# miga

[![Build Status](https://travis-ci.com/chapsuk/miga.svg?token=m33r59zSHRPMSbqfFKFk&branch=master)](https://travis-ci.com/chapsuk/miga)

Miga is a command line utility around several migration packages with single interface.
Aimed to add extra features and hide some limitations of existing golang migration CLI`s.

```command
> go get -u github.com/chapsuk/miga
```

```command
> docker run -it chapsuk/miga --help
```

## Supporting

PackageName | Version | Postgres            | MySQL    | Clickhouse
----------- | ------- | ------------------- | -------- | --------
[goose](https://github.com/pressly/goose)       |  2.7.0 (+ [patches](https://github.com/pressly/goose/compare/v2.7.0...chapsuk:clickhouse?expand=1))  |  :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:
[migrate](https://github.com/golang-migrate/migrate)     |  4.2.5  |  :heavy_check_mark: | :heavy_check_mark: |
[impg](https://github.com/im-kulikov/migrate)        |   0.1   |  :heavy_check_mark: | |


## Features

* Configuration by file or environment variables;
* Seeding db by migration way;
* One from several migration packages of your choice;
* Converting one migration format to another;
* Testing;

## Usage

Miga CLI inherit from goose CLI and may not be familiar to users of other utilities.
See commands description before usage

```text
≻ ./bin/miga migrate
NAME:
   miga migrate - Migrations root command

USAGE:
   miga migrate command [command options] [arguments...]

COMMANDS:
     convert  Converting migrations FROM_FORMAT to TO_FORMAT and store to DESTENITION_PATH
     create   Creates new migration sql file
     down     Roll back the version by 1
     down-to  Roll back to a specific VERSION
     redo     Re-run the latest migration
     reset    Roll back all migrations
     status   Dump the migration status for the current DB
     up       Migrate the DB to the most recent version available
     up-to    Migrate the DB to a specific VERSION
     version  Print the current version of the database
     help, h  Shows a list of commands or help for one command
```

## Configuration

Miga has 3 level configuration options by priority:

### Flags

```bash
--config value      Config file name (default: "")
--driver value      Migration driver name: goose, migrate, stump (default: "goose")
--log.level value   Logger level [debug|info|...] (default: "debug")
--log.format value  Logger output format console|json (default: "console")
```

### Environment variables

Name                        | SettingDefault | Description
--------------------------- | -------------- | -----------------------
MIGA_CONFIG                 | miga.yml       | config file
MIGA_DRIVER                 | goose          | one from [list](#supporting)
MIGA_POSTGRES_DSN           |                | postgres DSN string
MIGA_MYSQL_DSN              |                | mysql DSN string
MIGA_MIGRATE_PATH           | ./migrations   | migrations dir
MIGA_MIGRATE_TABLE_NAME     | db_version     | migrations db version table name
MIGA_SEED_PATH              | ./seeds        | seeds dir
MIGA_SEED_TABLE_NAME        | seed_version   | seeds version table name
MIGA_LOG_LEVEL              | info           | logging level
MIGA_LOG_FORMAT             | console        | logging format (console or json)

*prefix `MIGA` may be changed by build flag `-ldflags "-X main.Name=<NAME>"`

### Config file

```yml
driver: goose
postgres:
  dsn: "postgres://user:password@127.0.0.1:5432/miga?sslmode=disable"
# mysql:
#   dsn: "user:password@tcp(127.0.0.1:3306)/miga"
migrate:
  path: ./migrations
  table_name: db_version
seed:
  path: ./seeds
  table_name: seed_version
```

### Using without config

```bash
> MIGA_POSTGRES_DSN="postgres://user:password@127.0.0.1:5432/miga?sslmode=disable" \
  MIGA_MIGRATE_PATH=./tests/migrations/goose/ \
  MIGA_SEED_PATH=./tests/seeds/goose \
  miga --driver goose migrate up
```

## Tests

```text
≻ make db_up
≻ GOCACHE=off go test -v ./tests/
=== RUN   TestConvert

  Convert from MIGRATE to GOOSE source: ./migrations/migrate dest: ./migrations/tmp/goose ✔✔
    Given migrations GOOSE driver with MYSQL dialect ✔✔
      #1 up to first migartions (create users table) ✔✔✔✔✔✔✔✔✔
      #2 up to second (create wallets table) ✔✔✔✔✔✔✔✔✔
      #3 up to third (alter tables) ✔✔✔✔✔✔✔✔✔
      #4 retry last migration ✔✔✔✔✔✔✔✔✔
      #5 query with `;` ✔✔✔✔✔✔✔✔✔
      #101 incorrect migration (duplicate of 3 migration) ✔✔✔✔✔✔✔✔✔✔
      #102 try jump over failed migration ✔✔✔✔✔✔✔
      #103 down to 2 migartion ✔✔✔✔✔✔✔
      #104 reset all ✔✔✔✔✔✔✔✔
      #105 up to latest, but stop on failed ✔✔✔✔✔✔✔✔✔✔
      #106 down to first ✔✔✔✔✔✔✔✔
      #107 reset all ✔✔✔✔✔
...

  Given migrations GOOSE driver with POSTGRES dialect ✔✔
    #1 up to first migartions (create users table) ✔✔✔✔✔✔
    #2 up to second (create wallets table) ✔✔✔✔✔✔
    #3 up to third (alter tables) ✔✔✔✔✔✔
    #4 retry last migration ✔✔✔✔✔✔
    #5 query with `;` ✔✔✔✔✔✔
    #6 plpsql statement, should create histories table and func for create inheritans ✔✔✔✔✔✔✔✔✔✔
    #101 incorrect migration (duplicate of 3 migration) ✔✔✔✔✔✔✔
    #102 try jump over failed migration ✔✔✔✔
    #103 down to 2 migartion ✔✔✔✔
    #104 reset all ✔✔✔✔✔
    #105 up to latest, but stop on failed ✔✔✔✔✔✔✔
    #106 down to first ✔✔✔✔✔
    #107 reset all ✔✔✔✔


1488 total assertions

--- PASS: TestMigrations (3.35s)
PASS
ok    miga/tests	10.521s
```
