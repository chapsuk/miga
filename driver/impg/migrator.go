package impg

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chapsuk/miga/logger"
	"github.com/chapsuk/miga/utils"
	"github.com/go-pg/pg"
	orig "github.com/im-kulikov/migrate"
)

type Migrator struct {
	stump orig.Migrator
	db    *pg.DB
	dir   string
}

func New(dialect, dsn, tableName, dir string) (*Migrator, error) {
	if dialect != "postgres" {
		return nil, fmt.Errorf("unsupported dialect %s for STUMP driver", strings.ToUpper(dialect))
	}

	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	s, err := orig.New(orig.Options{
		DB:     db,
		Logger: logger.G(),
		Path:   dir,
	})

	if err != nil {
		return nil, err
	}

	// orig.SetTableName(tableName)

	return &Migrator{
		stump: s,
		db:    db,
		dir:   dir,
	}, nil
}

func (m Migrator) Create(name, ext string) error {
	_, _, err := utils.CreateMigrationsFiles(time.Now().Unix(), m.dir, name, ext)
	return err
}

func (m Migrator) Close() error {
	return m.db.Close()
}

func (m Migrator) Down() error {
	return m.stump.Down(1)
}

func (m Migrator) DownTo(v string) error {
	version, err := versionToInt64(v)
	if err != nil {
		return err
	}

	current, err := m.stump.Version()
	if err != nil {
		return err
	}

	diff := current - version
	if diff <= 0 {
		return fmt.Errorf("nothing to update, current: %d dest: %d", current, version)
	}

	return m.stump.Down(int(diff))
}

func (m Migrator) Redo() error {
	err := m.stump.Down(1)
	if err != nil {
		logger.G().Warnf("down cmd on redo failed: %s", err)
	}
	return m.stump.Up(1)
}

func (m Migrator) Reset() error {
	return m.stump.Down(0)
}

func (m Migrator) Status() error {
	return m.Version()
}

func (m Migrator) Up() error {
	return m.stump.Up(0)
}

func (m Migrator) UpTo(v string) error {
	version, err := versionToInt64(v)
	if err != nil {
		return err
	}

	current, err := m.stump.Version()
	if err != nil {
		return err
	}

	diff := version - current
	if diff <= 0 {
		return fmt.Errorf("nothing to update, current: %d dest: %d", current, version)
	}

	return m.stump.Up(int(diff))
}

func (m Migrator) Version() error {
	v, err := m.stump.Version()
	if err != nil {
		return err
	}
	n, err := m.stump.VersionName()
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
