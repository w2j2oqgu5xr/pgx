// Package pgx is a PostgreSQL driver and toolkit for Go.
//
// The pgx driver is a low-level, high performance interface that exposes
// PostgreSQL-specific features such as LISTEN/NOTIFY and COPY. It also includes
// a standard database/sql driver.
//
// Establishing a connection:
//
//	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer conn.Close(context.Background())
//
// Querying:
//
//	var name string
//	var weight int64
//	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
//		os.Exit(1)
//	}
package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

const (
	// TextFormatCode is the format code for text format.
	TextFormatCode = pgproto3.TextFormat
	// BinaryFormatCode is the format code for binary format.
	BinaryFormatCode = pgproto3.BinaryFormat
)

// ErrNoRows is returned by Scan when QueryRow doesn't return a row.
// In such a case, QueryRow returns a placeholder *Row value that defers
// this error until a Scan.
var ErrNoRows = fmt.Errorf("no rows in result set")

// ErrTxClosed is returned when a transaction is already closed.
var ErrTxClosed = fmt.Errorf("tx is closed")

// ErrTxCommitRollback is returned when a commit is attempted on a transaction
// that has already been rolled back.
var ErrTxCommitRollback = fmt.Errorf("commit unexpectedly resulted in rollback")

// Connect establishes a connection with a PostgreSQL server with a connection
// string. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
// for details.
func Connect(ctx context.Context, connString string) (*Conn, error) {
	return ConnectConfig(ctx, mustParseConfig(connString))
}

// ConnectConfig establishes a connection with a PostgreSQL server with a
// configuration struct. connConfig must have been created by ParseConfig.
func ConnectConfig(ctx context.Context, connConfig *Config) (*Conn, error) {
	return connect(ctx, connConfig)
}

// mustParseConfig panics on error — for internal use only.
func mustParseConfig(connString string) *Config {
	cfg, err := ParseConfig(connString)
	if err != nil {
		panic(fmt.Sprintf("pgx: unable to parse config: %v", err))
	}
	return cfg
}

// Identifier a PostgreSQL identifier or name. When used with QueryExecModeSimpleProtocol
// Identifiers are sanitized with (pgx.Identifier).Sanitize(). If an Identifier
// will be used with the simple protocol, it should be created as a pgx.Identifier
// rather than a string.
type Identifier []string

// Sanitize returns a sanitized string safe for SQL interpolation.
func (ident Identifier) Sanitize() string {
	var sb string
	for i, part := range ident {
		if i > 0 {
			sb += "."
		}
		sb += `"` + sanitizeIdentifierPart(part) + `"`
	}
	return sb
}

// sanitizeIdentifierPart escapes double quotes within an identifier part.
func sanitizeIdentifierPart(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			result = append(result, '"', '"')
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// CommandTag is the result of an Exec function.
type CommandTag = pgconn.CommandTag
