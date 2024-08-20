package dialectquery

import "fmt"

type Starrocks struct{}

var _ Querier = (*Starrocks)(nil)

func (m *Starrocks) CreateTable(tableName string) string {
	q := `CREATE TABLE IF NOT EXISTS %s (
		version_id bigint NOT NULL,
		id bigint NOT NULL AUTO_INCREMENT,
		is_applied boolean NOT NULL,
		tstamp datetime NULL default CURRENT_TIMESTAMP
	)
	PRIMARY KEY (version_id,id)
	DISTRIBUTED BY HASH (id)
	ORDER BY (version_id)`
	return fmt.Sprintf(q, tableName)
}

func (m *Starrocks) InsertVersion(tableName string) string {
	q := `INSERT INTO %s (version_id, is_applied) VALUES (?, ?)`
	return fmt.Sprintf(q, tableName)
}

func (m *Starrocks) DeleteVersion(tableName string) string {
	q := `DELETE FROM %s WHERE version_id=?`
	return fmt.Sprintf(q, tableName)
}

func (m *Starrocks) GetMigrationByVersion(tableName string) string {
	q := `SELECT tstamp, is_applied FROM %s WHERE version_id=? ORDER BY tstamp DESC LIMIT 1`
	return fmt.Sprintf(q, tableName)
}

func (m *Starrocks) ListMigrations(tableName string) string {
	q := `SELECT version_id, is_applied from %s ORDER BY version_id DESC`
	return fmt.Sprintf(q, tableName)
}

func (m *Starrocks) GetLatestVersion(tableName string) string {
	q := `SELECT MAX(version_id) FROM %s`
	return fmt.Sprintf(q, tableName)
}
