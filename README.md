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

PackageName | Version | Postgres            | MySQL    | Clickhouse | Vertica
----------- | ------- | ------------------- | -------- | ---------- | ----
[goose](https://github.com/pressly/goose)       |  3.9.1 ([patch](https://github.com/chapsuk/goose/commit/d8dae35e216b5b70d3db4e986884f715b5a280cc)  |  :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark:
[migrate](https://github.com/golang-migrate/migrate)     |  4.2.5  |  :heavy_check_mark: | :heavy_check_mark: | |
[impg](https://github.com/im-kulikov/migrate)        |   0.1   |  :heavy_check_mark: | | |


## Features

* One from several migration packages of your choice;
* Converting one migration format to another;
* Testing;

## Usage

Miga CLI inherit from goose CLI and may not be familiar to users of other utilities.
See commands description before usage

```text
≻ ./bin/miga
NAME:
   miga - Migrations root command

USAGE:
   miga command [command options] [arguments...]

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

The config file after v1.0.0 required

### Flags

```bash
--config value   Config file name (default: "miga.yml")
```

### Config file

```yml
miga:
  driver: goose
  path: ./migrations
  table: db_version

db:
  dsn: "postgres://user:password@127.0.0.1:5432/miga?sslmode=disable"
  dialect: postgres

logger:
  level: info
  format: console
```

### Using without config

```bash
> miga -c miga.yml up
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
