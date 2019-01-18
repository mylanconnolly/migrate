package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx"
)

// Migrator is the type that performs migrations.
type Migrator struct {
	conn *pgx.ConnPool
	dir  string
}

func (m Migrator) ensureTableExists(tx *pgx.Tx) error {
	found := false
	findTableSQL := "SELECT EXISTS(SELECT * FROM information_schema.tables WHERE table_name = 'schema_migrations')"

	if err := tx.QueryRow(findTableSQL).Scan(&found); err != nil {
		return err
	}
	if found {
		return nil
	}
	makeTableSQL := "CREATE TABLE schema_migrations (version BIGINT PRIMARY KEY, migrated_at TIMESTAMPTZ)"
	_, err := tx.Exec(makeTableSQL)
	return err
}

func (m Migrator) getVersion(tx *pgx.Tx) (int64, error) {
	var version int64

	if err := m.ensureTableExists(tx); err != nil {
		return 0, err
	}
	// Although we're only maintaining one row in the schema_migrations table,
	// you never know what might happen with the end user, and the table may end
	// up in an inconsistent state. We will make a best-effort attempt to get the
	// version of the database.
	getVersionSQL := "SELECT version FROM schema_migrations ORDER BY migrated_at DESC LIMIT 1"
	err := tx.QueryRow(getVersionSQL).Scan(&version)

	switch err {
	case pgx.ErrNoRows:
		return 0, nil
	case nil:
		return version, nil
	default:
		return 0, err
	}
}

func (m Migrator) getMigrations() ([]Migration, error) {
	var migrations []Migration
	stat, err := os.Stat(m.dir)

	if err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("could not open directory '%s'", m.dir)
	}
	files, err := filepath.Glob(filepath.Join(m.dir, "*.up.sql"))

	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	for _, filename := range files {
		m, err := migrationFromFile(filename)

		if err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}
	return migrations, nil
}

// Up is used to migrate the database up to the newest version; that is, the
// state after running all of the "up" migrations.
func (m Migrator) Up() error {
	tx, err := m.conn.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	version, err := m.getVersion(tx)

	if err != nil {
		return err
	}
	migrations, err := m.getMigrations()

	if err != nil {
		return err
	}
	for _, migration := range migrations {
		// The migration appears to have already been applied so we can skip it.
		if migration.version <= version {
			continue
		}
		if _, err = m.conn.Exec(migration.up); err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// Down is used to migrate the "clean" state; that is, the state after running
// all of the "down" migrations.
func (m Migrator) Down() error {
	tx, err := m.conn.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	version, err := m.getVersion(tx)

	if err != nil {
		return err
	}
	migrations, err := m.getMigrations()

	if err != nil {
		return err
	}
	for i := len(migrations) - 1; i >= 0; i-- {
		// The migration appears to not be applied yet so we can skip it.
		if migrations[i].version >= version {
			continue
		}
		if _, err = m.conn.Exec(migrations[i].down); err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// New is used to create a new migrator with the given options. If options are
// not provided, default values are used instead. This allows for simplified
// development, where you may be using default connection information, anyway.
func New(dir string, options ...Option) (*Migrator, error) {
	config := pgx.ConnConfig{}

	for _, opt := range options {
		opt.ApplyConfig(&config)
	}
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: config,
	})
	if err != nil {
		return nil, err
	}
	return &Migrator{conn: conn, dir: dir}, nil
}

// NewConn is an alternative way to create a new migrator by supplying a conn
// pool directly, rather than letting us create it.
func NewConn(dir string, conn *pgx.ConnPool) *Migrator {
	return &Migrator{conn: conn, dir: dir}
}
