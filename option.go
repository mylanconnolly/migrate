package migrate

import "github.com/jackc/pgx"

// Option is used to define an option type so that users can supply a varying
// amount of options to override default connection behavior.
type Option interface {
	ApplyConfig(*pgx.ConnConfig)
}

// UsernameOption is used to define the username for a connection.
type UsernameOption string

// ApplyConfig is used to implement Option.
func (u UsernameOption) ApplyConfig(cfg *pgx.ConnConfig) { cfg.User = string(u) }

// PasswordOption is used to define the password for a connection.
type PasswordOption string

// ApplyConfig is used to implement Option.
func (p PasswordOption) ApplyConfig(cfg *pgx.ConnConfig) { cfg.Password = string(p) }

// DatabaseOption is used to define the database name used for a connection.
type DatabaseOption string

// ApplyConfig is used to implement Option.
func (d DatabaseOption) ApplyConfig(cfg *pgx.ConnConfig) { cfg.Database = string(d) }

// HostnameOption is used to define the database server hostname used for a
// connection.
type HostnameOption string

// ApplyConfig is used to implement Option.
func (h HostnameOption) ApplyConfig(cfg *pgx.ConnConfig) { cfg.Host = string(h) }

// PortOption is used to define the port number for a connection.
type PortOption uint16

// ApplyConfig is used to implement Option.
func (p PortOption) ApplyConfig(cfg *pgx.ConnConfig) { cfg.Port = uint16(p) }
