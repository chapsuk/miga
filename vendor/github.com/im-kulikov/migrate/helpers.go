package migrate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/types"
)

// getTableName for quote table name
func getTableName() types.ValueAppender {
	return pg.Q(tableName)
}

// extractAttributes, such as  version, name and migration-type
func extractAttributes(filename string) (version int64, name, mType string, err error) {
	parts := strings.SplitN(filename, "_", 2)

	if len(parts) != 2 {
		err = fmt.Errorf(errFileNamingTpl, filename)
		return
	}

	if version, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return
	} else if version <= 0 {
		err = fmt.Errorf(errFileVersionTpl, parts[0])
		return
	}

	parts = strings.SplitN(parts[1], ".", 3)

	if len(parts) != 3 || parts[1] != "down" && parts[1] != "up" {
		err = fmt.Errorf(errFileNamingTpl, filename)
		return
	}

	name, mType = parts[0], parts[1]

	return
}

// findMigrations in specified folder (path)
func findMigrations(path string) ([]os.FileInfo, error) {
	var (
		err error
		dir *os.File
	)

	if dir, err = os.Open(path); os.IsNotExist(err) {
		return nil, ErrDirNotExists
	}

	defer dir.Close()

	return dir.Readdir(0)
}

// updateVersion abstraction
type updateVersion func(tx *pg.Tx, version int64, name string) error

// remVersion migration from database
func remVersion(tx *pg.Tx, version int64, name string) error {
	_, err := tx.Exec(sqlRemVersion, getTableName(), version, name)
	return err
}

// addVersion migration to database
func addVersion(tx *pg.Tx, version int64, name string) error {
	_, err := tx.Exec(sqlNewVersion, getTableName(), version, name)
	return err
}

// doMigrate closure
func doMigrate(version int64, name, sql string, fn updateVersion) func(db DB) error {
	return func(db DB) error {
		return db.RunInTransaction(func(tx *pg.Tx) error {
			if _, errQuery := tx.Exec(sql); errQuery != nil {
				return errQuery
			}

			if errVersion := fn(tx, version, name); errVersion != nil {
				return errVersion
			}

			return nil
		})
	}
}

// extractMigrations, find files in migration folder and convert to Migration-item
func extractMigrations(log Logger, folder string, files []os.FileInfo) (Migrations, error) {
	var (
		err          error
		data         []byte
		migrateParts = make(map[string]*Migration)
		items        Migrations
	)

	for _, file := range files {
		// Ignore non sql files:
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		log.Infof("Prepare migration file: %s", file.Name())

		if data, err = ioutil.ReadFile(path.Join(folder, file.Name())); err != nil {
			return nil, err
		}

		ver, name, mType, err := extractAttributes(file.Name())
		if err != nil {
			return nil, err
		}

		m, ok := migrateParts[name]
		if !ok {
			m = &Migration{
				Version: ver,
				Name:    name,
			}
		} else if m.Version != ver {
			return nil, fmt.Errorf(errVersionNotEqualTpl, m.Version, ver)
		}

		switch mType {
		case "up":
			m.Up = doMigrate(m.Version, m.RealName(), string(data), addVersion)
		case "down":
			m.Down = doMigrate(m.Version, m.RealName(), string(data), remVersion)
		}

		migrateParts[name] = m
	}

	items = make(Migrations, 0, len(migrateParts))
	for name, m := range migrateParts {
		log.Infof("Prepare migration: %s", name)

		if m.Down == nil || m.Up == nil {
			return nil, ErrBothMigrateTypes
		}

		items = append(items, m)
	}

	return items, nil
}

// CreateMigration files
func CreateMigration(folder, name string) error {
	var (
		err   error
		f     *os.File
		dt    = time.Now().Unix()
		items = []string{"up", "down"}
	)

	for _, item := range items {
		filename := path.Join(folder, fmt.Sprintf(fileNameTpl, dt, name, item))
		if f, err = os.Create(filename); err != nil {
			return err
		}

		f.Close()
	}

	return nil
}
