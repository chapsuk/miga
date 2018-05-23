package migrate

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Logger interface for migrator
type Logger interface {
	Infof(format string, args ...interface{})
}

// Options for migrator
type Options struct {
	// DB connection
	DB DB
	// Path to migrations files
	Path string
	// Logger
	Logger Logger
}

// Migrator interface
type Migrator interface {
	Up(steps int) error
	Down(steps int) error
	List() (Migrations, error)
	Plan() (Migrations, error)
	Version() (int64, error)
	VersionName() (string, error)
}

// DB interface
type DB interface {
	RunInTransaction(fn func(*pg.Tx) error) error
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecOne(query interface{}, params ...interface{}) (orm.Result, error)
	Query(model, query interface{}, params ...interface{}) (orm.Result, error)
	QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error)
}

// migrate is implementation of Migrator
type migrate struct {
	Options
	Migrations
}

// Migration item
type Migration struct {
	Version   int64
	Name      string
	CreatedAt time.Time
	Up        func(DB) error
	Down      func(DB) error
}

// RealName return formatted filename
func (m Migration) RealName() string {
	return strconv.FormatInt(m.Version, 10) + "_" + m.Name
}

// Migrations slice
type Migrations []*Migration

// New creates new Migrator
func New(opts Options) (Migrator, error) {
	var err error

	if opts.DB == nil {
		return nil, ErrNoDB
	}

	if err = createTables(opts.DB); err != nil {
		return nil, err
	}

	return &migrate{Options: opts}, nil
}

func prepareMigrations(migrate *migrate) (err error) {
	var (
		files []os.FileInfo
		opts  = migrate.Options
	)

	if files, err = findMigrations(opts.Path); err != nil {
		return
	}

	migrate.Migrations, err = extractMigrations(opts.Logger, opts.Path, files)

	return
}

// createTables for migrations
func createTables(db DB) error {
	var err error
	if len(schemaName) > 0 {
		if _, err = db.Exec(
			sqlCreateSchema,
			pg.Q(schemaName),
		); err != nil {
			return err
		}
	}

	_, err = db.Exec(sqlCreateTable, pg.Q(tableName))

	return err
}

// Up, roll up multiple migrations
func (m *migrate) Up(steps int) error {
	var (
		err     error
		version int64
		count   int
	)

	if err = prepareMigrations(m); err != nil {
		return err
	}

	count = len(m.Migrations)

	if steps < 0 {
		return ErrPositiveSteps
	}

	if steps == 0 {
		steps = count
	}

	if version, err = m.Version(); err != nil {
		return err
	}

	items := make(Migrations, count)

	copy(items, m.Migrations)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Version < items[j].Version
	})

	for i, item := range items {
		if steps <= 0 {
			break
		}

		if item.Version <= version {
			continue
		}

		m.Logger.Infof("(%d) migrate up to: %d_%s", i+1, item.Version, item.Name)
		if err = item.Up(m.DB); err != nil {
			return err
		}

		steps--
	}

	return nil
}

// Down rollback some migrations
func (m *migrate) Down(steps int) error {
	var (
		err     error
		version int64
		count   int
	)

	if err = prepareMigrations(m); err != nil {
		return err
	}

	count = len(m.Migrations)

	if steps < 0 {
		return ErrPositiveSteps
	}

	if steps > count || steps == 0 {
		steps = count
	}

	if version, err = m.Version(); err != nil {
		return err
	}

	items := make(Migrations, count)

	copy(items, m.Migrations)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Version > items[j].Version
	})

	for _, item := range items {
		if steps <= 0 {
			break
		}

		if item.Version > version {
			continue
		}

		m.Logger.Infof("(%d) migrate down to: %d_%s", steps, item.Version, item.Name)
		if err = item.Down(m.DB); err != nil {
			return err
		}

		steps--
	}

	return nil
}

func (m *migrate) List() (Migrations, error) {
	var v []struct {
		Version   int64
		Name      string
		CreatedAt time.Time
	}

	if _, err := m.DB.Query(&v, sqlSelectVersion, getTableName()); err != nil {
		return nil, err
	}

	result := make(Migrations, 0, len(v))

	for _, item := range v {
		name := strings.Replace(item.Name, strconv.FormatInt(item.Version, 10)+"_", "", -1)
		result = append(result, &Migration{
			Version:   item.Version,
			Name:      name,
			CreatedAt: item.CreatedAt,
		})
	}

	return result, nil
}

func (m *migrate) Plan() (Migrations, error) {
	var v, err = m.Version()
	if err != nil {
		return nil, err
	}

	if err = prepareMigrations(m); err != nil {
		return nil, err
	}

	var result Migrations

	for _, mig := range m.Migrations {
		if mig.Version > v {
			result = append(result, mig)
		}
	}

	return result, nil
}

// Version fetching from database
func (m *migrate) Version() (version int64, err error) {
	version = -1

	if err = createTables(m.DB); err != nil {
		return
	}

	if _, err = m.DB.QueryOne(
		pg.Scan(&version),
		sqlGetVersion,
		getTableName(),
	); err != nil && err == pg.ErrNoRows {
		err = nil
		version = 0
	}

	return
}

func (m *migrate) VersionName() (version string, err error) {
	if err = createTables(m.DB); err != nil {
		return
	}

	if _, err = m.DB.QueryOne(
		pg.Scan(&version),
		sqlGetName,
		getTableName(),
	); err != nil && err == pg.ErrNoRows {
		err = nil
	}

	return
}
