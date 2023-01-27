package goose

import (
	"database/sql"
	"fmt"

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
	g.clickhouseHackEnsureTable()
	return orig.Run("down", g.db, g.dir)
}

func (g Goose) DownTo(version string) error {
	g.clickhouseHackEnsureTable()
	return orig.Run("down-to", g.db, g.dir, version)
}

func (g Goose) Redo() error {
	g.clickhouseHackEnsureTable()
	return orig.Run("redo", g.db, g.dir)
}

func (g Goose) Reset() error {
	g.clickhouseHackEnsureTable()
	return orig.Run("reset", g.db, g.dir)
}

func (g Goose) Status() error {
	g.clickhouseHackEnsureTable()
	return orig.Run("status", g.db, g.dir)
}

func (g Goose) Up() error {
	g.clickhouseHackEnsureTable()
	return orig.Run("up", g.db, g.dir)
}

func (g Goose) UpTo(version string) error {
	g.clickhouseHackEnsureTable()
	return orig.Run("up-to", g.db, g.dir, version)
}

func (g Goose) Version() error {
	g.clickhouseHackEnsureTable()
	return orig.Run("version", g.db, g.dir)
}

func (g Goose) clickhouseHackEnsureTable() {
	if !g.isClickhouse() || len(g.clickhouseClusterName) == 0 {
		return
	}

	var (
		queries     = []string{}
		schemaTable = fmt.Sprintf("%s.%s", g.clickhouseSchema, g.versionTableName)
	)

	if g.clickhouseSharded {
		queries = append(queries, fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s_shard ON CLUSTER '%s' (
				version_id Int64,
				is_applied UInt8,
				date Date default now(),
				tstamp DateTime default now()
			) Engine = %s
			ORDER BY tstamp
		`, schemaTable, g.clickhouseClusterName, g.clickhouseEngine))

		queries = append(queries, fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s ON CLUSTER '%s' AS %s_shard
			ENGINE = Distributed('%s', %s, %s_shard, rand())
		`, schemaTable, g.clickhouseClusterName, schemaTable, g.clickhouseClusterName, g.clickhouseSchema, g.versionTableName))
	} else {
		queries = append(queries, fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s ON CLUSTER '%s' (
				version_id Int64,
				is_applied UInt8,
				date Date default now(),
				tstamp DateTime default now()
			) Engine = %s
			ORDER BY tstamp
		`, schemaTable, g.clickhouseClusterName, g.clickhouseEngine))
	}

	for _, q := range queries {
		if _, err := g.db.Exec(q); err != nil {
			panic("Failed applly clickhouse dirty hack: " + err.Error())
		}
	}

	var total int
	err := g.db.
		QueryRow(fmt.Sprintf("select count(*) from %s", schemaTable)).
		Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		panic("Failed get last applied version: " + err.Error())
	} else if total > 0 {
		return
	}

	tx, err := g.db.Begin()
	if err != nil {
		panic("Failed begin tx: " + err.Error())
	}
	defer tx.Rollback() // nolint: errcheck

	if _, err := tx.Exec(fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES (?, ?)", schemaTable), 0, 1); err != nil {
		panic("Failed insert initial version: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		panic("Failed commit initial version: " + err.Error())
	}
}
