package migrate

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"miga/logger"
	"miga/utils"

	orig "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ErrMissingDBConfig = errors.New("missing db config for migrate driver")

type Migrator struct {
	backend *orig.Migrate
	dir     string
	db      *sql.DB
	hack    source.Driver
}

func New(dialect, dsn, tableName, dir string) (*Migrator, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	if dsn == "" {
		return &Migrator{dir: dir}, nil
	}

	var driver database.Driver
	switch dialect {
	case "postgres":
		driver, err = postgres.WithInstance(db, &postgres.Config{
			MigrationsTable: tableName,
		})
	case "mysql":
		driver, err = mysql.WithInstance(db, &mysql.Config{
			MigrationsTable: tableName,
		})
	default:
		return nil, errors.New("Unsupported dialect")
	}

	if err != nil {
		return nil, err
	}

	// hack for get actual db version and fix dirty state
	source := "file://" + dir
	f := &file.File{}
	h, err := f.Open(source)
	if err != nil {
		return nil, err
	}

	m, err := orig.NewWithDatabaseInstance(source, dialect, driver)
	if err != nil {
		return nil, err
	}
	m.Log = &utils.StdLogger{}

	return &Migrator{
		backend: m,
		dir:     dir,
		db:      db,
		hack:    h,
	}, nil
}

func (m Migrator) Close() error {
	if m.db == nil {
		return nil
	}
	return m.db.Close()
}

func (m Migrator) Create(name, ext string) error {
	_, _, err := utils.CreateMigrationsFiles(time.Now().Unix(), m.dir, name, ext)
	return err
}

func (m Migrator) Down() error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	return m.fixDirtyState(m.backend.Steps(-1))
}

func (m Migrator) DownTo(version string) error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	v, err := versionToUint(version)
	if err != nil {
		return err
	}

	current, _, err := m.backend.Version()
	if v >= current {
		return fmt.Errorf("Nothing to update, current version: %d", current)
	}

	return m.fixDirtyState(m.backend.Migrate(uint(v)))
}

func (m Migrator) Redo() error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	err := m.backend.Steps(-1)
	if err != nil {
		return m.fixDirtyState(err)
	}

	return m.fixDirtyState(m.backend.Steps(1))
}

func (m Migrator) Reset() error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	return m.fixDirtyState(m.backend.Down())
}

func (m Migrator) Status() error {
	return m.Version()
}

func (m Migrator) Up() error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	return m.fixDirtyState(m.backend.Up())
}

func (m Migrator) UpTo(version string) error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	v, err := versionToUint(version)
	if err != nil {
		return err
	}

	current, _, err := m.backend.Version()
	if v <= current {
		return fmt.Errorf("Nothing to update, current version: %d", current)
	}

	return m.fixDirtyState(m.backend.Migrate(uint(v)))
}

func (m Migrator) Version() error {
	if m.db == nil {
		return ErrMissingDBConfig
	}
	version, dirty, err := m.backend.Version()
	if err != nil {
		return err
	}
	logger.G().Infof("Current version: %d dirty: %t", version, dirty)
	return nil
}

func (m Migrator) fixDirtyState(err error) error {
	if err == nil {
		return nil
	}

	current, dirty, verr := m.backend.Version()
	if verr != nil {
		return errors.Wrapf(err, "get current version for fix dirty state error: %s", verr)
	}

	if !dirty {
		return err
	}

	actual, aerr := m.hack.Prev(current)
	if aerr != nil {
		logger.G().Warnf("(skip) get prev version for fix dirty state error: %s", aerr)
	}

	ferr := m.backend.Force(int(actual))
	if ferr != nil {
		return errors.Wrapf(err, "set force version error: %s", ferr)
	}

	return err
}

func versionToUint(version string) (uint, error) {
	v, err := strconv.Atoi(version)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
