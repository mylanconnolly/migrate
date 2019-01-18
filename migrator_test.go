// +build integration

package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx"
)

var (
	testUser            = "postgres"
	testHostname        = "localhost"
	testPort     uint16 = 5432
	testDB       string
	testPassword string
)

func newConn() (*pgx.ConnPool, error) {
	return pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     testUser,
			Password: testPassword,
			Host:     testHostname,
			Port:     testPort,
		},
	})
}

func setup() error {
	conn, err := newConn()

	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, testDB))
	return err
}

func teardown() error {
	conn, err := newConn()

	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, testDB))
	return nil
}

func TestMain(m *testing.M) {
	testDB = fmt.Sprintf("migration_test_%s", time.Now().Format("2006010215040506"))

	if err := setup(); err != nil {
		fmt.Println("Could not perform setup tasks:", err.Error())
		os.Exit(1)
	}

	retcode := m.Run()

	if err := teardown(); err != nil {
		fmt.Println("Could not perform teardown tasks:", err.Error())
		os.Exit(1)
	}
	os.Exit(retcode)
}

func TestMigrations(t *testing.T) {
	tests := []struct {
		name        string
		dir         string
		wantTables  []string
		wantIndexes []string
		wantErrUp   bool
		wantErrDown bool
	}{
		{"missing up", filepath.Join("testdata", "missing_up"), []string{}, []string{}, false, false},
		{"missing down", filepath.Join("testdata", "missing_down"), []string{}, []string{}, true, false},
		{"invalid filenames", filepath.Join("testdata", "invalid_filenames"), []string{}, []string{}, true, false},
		{"single valid", filepath.Join("testdata", "single_valid"), []string{"test_table"}, []string{"test_table_idx"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := New(tt.dir, DatabaseOption(testDB), HostnameOption(testHostname), UsernameOption(testUser), PasswordOption(testPassword), PortOption(testPort))

			if err != nil {
				t.Fatalf("Got error connecting to database: %#v", err)
			}
			errUp := m.Up()

			if tt.wantErrUp {
				if errUp == nil {
					t.Fatalf("Wanted error migrating up, got nil")
				}
				return
			}
			if errUp != nil {
				t.Fatalf("Got error migrationg up: %#v, wanted nil", errUp)
			}
			errDown := m.Down()

			if tt.wantErrDown {
				if errDown == nil {
					t.Fatalf("Wanted error migrating down, got nil")
				}
				return
			}
			if errDown != nil {
				t.Fatalf("Got error migrationg down: %#v, wanted nil", errUp)
			}
		})
	}
}
