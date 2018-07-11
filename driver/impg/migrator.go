package impg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chapsuk/miga/logger"
	"github.com/chapsuk/miga/utils"
	"github.com/go-pg/pg"
	orig "github.com/im-kulikov/migrate"
)

var ErrMissingDBConfig = errors.New("missing DB config for impg driver")

type Migrator struct {
	impg orig.Migrator
	db   *pg.DB
	dir  string
}

func New(dialect, dsn, tableName, dir string) (*Migrator, error) {
	if dialect != "postgres" {
		return nil, fmt.Errorf("unsupported dialect %s for impg driver", strings.ToUpper(dialect))
	}

	if dsn == "" {
		// may be used for create command, who did not need db connection
		return &Migrator{dir: dir}, nil
	}

	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	impg, err := orig.New(orig.Options{
		DB:     db,
		Logger: logger.G(),
		Path:   dir,
	})

	if err != nil {
		return nil, err
	}

	orig.SetTableName(tableName)

	return &Migrator{
		impg: impg,
		db:   db,
		dir:  dir,
	}, nil
}

func (m Migrator) Create(name, ext string) error {
	_, _, err := utils.CreateMigrationsFiles(time.Now().Unix(), m.dir, name, ext)
	return err
}

func (m Migrator) Close() error {
	if m.db == nil {
		return nil
	}
	return m.db.Close()
}

func (m Migrator) Down() error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}
	return m.impg.Down(1)
}

func (m Migrator) DownTo(v string) error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}

	version, err := versionToInt64(v)
	if err != nil {
		return err
	}

	current, err := m.impg.Version()
	if err != nil {
		return err
	}

	diff := current - version
	if diff <= 0 {
		return fmt.Errorf("nothing to update, current: %d dest: %d", current, version)
	}

	return m.impg.Down(int(diff))
}

func (m Migrator) Redo() error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}

	err := m.impg.Down(1)
	if err != nil {
		logger.G().Warnf("down cmd on redo failed: %s", err)
	}
	return m.impg.Up(1)
}

func (m Migrator) Reset() error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}
	return m.impg.Down(0)
}

func (m Migrator) Status() error {
	return m.Version()
}

func (m Migrator) Up() error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}
	return m.impg.Up(0)
}

func (m Migrator) UpTo(v string) error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}
	version, err := versionToInt64(v)
	if err != nil {
		return err
	}

	current, err := m.impg.Version()
	if err != nil {
		return err
	}

	diff := version - current
	if diff <= 0 {
		return fmt.Errorf("nothing to update, current: %d dest: %d", current, version)
	}

	return m.impg.Up(int(diff))
}

func (m Migrator) Version() error {
	if m.impg == nil {
		return ErrMissingDBConfig
	}
	v, err := m.impg.Version()
	if err != nil {
		return err
	}
	n, err := m.impg.VersionName()
	if err != nil {
		return err
	}

	logger.G().Infof("Current version %d (%s)", v, n)
	return nil
}

func versionToInt64(version string) (int64, error) {
	v, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
