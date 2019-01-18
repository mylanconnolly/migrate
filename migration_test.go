package migrate

import (
	"path/filepath"
	"reflect"
	"testing"
)

var testDown1 = `DROP TABLE test_table;
`

var testUp1 = `CREATE TABLE test_table (
  id BIGSERIAL PRIMARY KEY,
  inserted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX test_table_idx ON test_table (inserted_at);
`

func TestMigrationFromFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    Migration
		wantErr bool
	}{
		{"missing up", filepath.Join("testdata", "missing_up", "1_create_test_table.up.sql"), Migration{}, true},
		{"missing down", filepath.Join("testdata", "missing_down", "1_create_test_table.up.sql"), Migration{}, true},
		{"invalid filenames", filepath.Join("testdata", "invalid_filenames", "create_test_table.up.sql"), Migration{}, true},
		{
			"valid", filepath.Join("testdata", "single_valid", "1_create_test_table.up.sql"),
			Migration{up: testUp1, down: testDown1, version: 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := migrationFromFile(tt.path)

			if tt.wantErr && err == nil {
				t.Fatalf("Wanted an error, got nil")
			}
			if err != nil && !tt.wantErr {
				t.Fatalf("Got error: %#v, wanted nil", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Got: %#v, wanted: %#v", got, tt.want)
			}
		})
	}
}
