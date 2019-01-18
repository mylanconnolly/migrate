package migrate

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx"
)

func TestUsernameApplyConfig(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		config pgx.ConnConfig
		want   pgx.ConnConfig
	}{
		{"empty config", "foo", pgx.ConnConfig{}, pgx.ConnConfig{User: "foo"}},
		{"replace value", "foo", pgx.ConnConfig{User: "bar"}, pgx.ConnConfig{User: "foo"}},
		{"respect other values", "foo", pgx.ConnConfig{Password: "xxx"}, pgx.ConnConfig{User: "foo", Password: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UsernameOption(tt.val).ApplyConfig(&tt.config)

			if !reflect.DeepEqual(tt.want, tt.config) {
				t.Fatalf("Got: %#v, want: %#v", tt.config, tt.want)
			}
		})
	}
}

func TestPasswordApplyConfig(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		config pgx.ConnConfig
		want   pgx.ConnConfig
	}{
		{"empty config", "foo", pgx.ConnConfig{}, pgx.ConnConfig{Password: "foo"}},
		{"replace value", "foo", pgx.ConnConfig{Password: "bar"}, pgx.ConnConfig{Password: "foo"}},
		{"respect other values", "foo", pgx.ConnConfig{User: "xxx"}, pgx.ConnConfig{Password: "foo", User: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PasswordOption(tt.val).ApplyConfig(&tt.config)

			if !reflect.DeepEqual(tt.want, tt.config) {
				t.Fatalf("Got: %#v, want: %#v", tt.config, tt.want)
			}
		})
	}
}

func TestDatabaseApplyConfig(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		config pgx.ConnConfig
		want   pgx.ConnConfig
	}{
		{"empty config", "foo", pgx.ConnConfig{}, pgx.ConnConfig{Database: "foo"}},
		{"replace value", "foo", pgx.ConnConfig{Database: "bar"}, pgx.ConnConfig{Database: "foo"}},
		{"respect other values", "foo", pgx.ConnConfig{Password: "xxx"}, pgx.ConnConfig{Database: "foo", Password: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DatabaseOption(tt.val).ApplyConfig(&tt.config)

			if !reflect.DeepEqual(tt.want, tt.config) {
				t.Fatalf("Got: %#v, want: %#v", tt.config, tt.want)
			}
		})
	}
}

func TestHostnameApplyConfig(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		config pgx.ConnConfig
		want   pgx.ConnConfig
	}{
		{"empty config", "foo", pgx.ConnConfig{}, pgx.ConnConfig{Host: "foo"}},
		{"replace value", "foo", pgx.ConnConfig{Host: "bar"}, pgx.ConnConfig{Host: "foo"}},
		{"respect other values", "foo", pgx.ConnConfig{Password: "xxx"}, pgx.ConnConfig{Host: "foo", Password: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HostnameOption(tt.val).ApplyConfig(&tt.config)

			if !reflect.DeepEqual(tt.want, tt.config) {
				t.Fatalf("Got: %#v, want: %#v", tt.config, tt.want)
			}
		})
	}
}

func TestPortApplyConfig(t *testing.T) {
	tests := []struct {
		name   string
		val    uint16
		config pgx.ConnConfig
		want   pgx.ConnConfig
	}{
		{"empty config", 5555, pgx.ConnConfig{}, pgx.ConnConfig{Port: 5555}},
		{"replace value", 5555, pgx.ConnConfig{Port: 5432}, pgx.ConnConfig{Port: 5555}},
		{"respect other values", 5555, pgx.ConnConfig{Password: "xxx"}, pgx.ConnConfig{Port: 5555, Password: "xxx"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PortOption(tt.val).ApplyConfig(&tt.config)

			if !reflect.DeepEqual(tt.want, tt.config) {
				t.Fatalf("Got: %#v, want: %#v", tt.config, tt.want)
			}
		})
	}
}
