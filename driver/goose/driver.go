package goose

import (
	"database/sql"

	"miga/utils"

	orig "github.com/pressly/goose/v3"
)

type Goose struct {
	db                    *sql.DB
	dir                   string
	dialect               string
	versionTableName      string
	clickhouseSchema      string
	clickhouseClusterName string
	clickhouseEngine      string
	clickhouseSharded     bool
}

func New(
	dialect, dsn, tableName, dir string,
	clickhouseSchema, clickhouseClusterName, clickhouseEngine string,
	clickhouseSharded bool,
) (*Goose, error) {
	err := orig.SetDialect(dialect)
	if err != nil {
		return nil, err
	}

	orig.SetTableName(tableName)
	orig.SetLogger(&utils.StdLogger{})

	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	return &Goose{
		db:                    db,
		dir:                   dir,
		dialect:               dialect,
		versionTableName:      tableName,
		clickhouseSchema:      clickhouseSchema,
		clickhouseClusterName: clickhouseClusterName,
		clickhouseEngine:      clickhouseEngine,
		clickhouseSharded:     clickhouseSharded,
	}, nil
}

func (g Goose) isClickhouse() bool {
	return g.dialect == "clickhouse"
}

func (g Goose) Close() error {
	return g.db.Close()
}

func (g Goose) Create(name, ext string) error {
	return orig.Run("create", g.db, g.dir, name, ext)
}

func (g Goose) Down() error {
	return orig.Run("down", g.db, g.dir)
}

func (g Goose) DownTo(version string) error {
	return orig.Run("down-to", g.db, g.dir, version)
}

func (g Goose) Redo() error {
	return orig.Run("redo", g.db, g.dir)
}

func (g Goose) Reset() error {
	return orig.Run("reset", g.db, g.dir)
}

func (g Goose) Status() error {
	return orig.Run("status", g.db, g.dir)
}

func (g Goose) Up() error {
	return orig.Run("up", g.db, g.dir)
}

func (g Goose) UpTo(version string) error {
	return orig.Run("up-to", g.db, g.dir, version)
}

func (g Goose) Version() error {
	return orig.Run("version", g.db, g.dir)
}
