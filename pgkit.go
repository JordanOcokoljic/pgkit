package pgkit

import (
	"database/sql"

	// Used to establish a connections to postgres databases.
	_ "github.com/lib/pq"
)

// DataProvider provides methods for interating with a SQL database, and is
// fulfilled by both *sql.DB and *sql.TX which allow for transactioned unit
// tests or mock structs.
type DataProvider interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}

// ConnectionDetail represents the information used to connect to the database.
type ConnectionDetail struct {
	User     string
	Password string
	Location string
	Port     string
	Database string
	Options  map[string]string
}

// NewConnectionDetail will return a ConnectionDetail with the Options map
// correctly initialized.
func NewConnectionDetail() ConnectionDetail {
	cd := ConnectionDetail{}
	cd.Options = make(map[string]string)
	return cd
}

// ParseDetails extracts the connection details out of the connection URI.
func ParseDetails(connection string) (ConnectionDetail, error) {
	cd := NewConnectionDetail()
	return cd, nil
}

// DB wraps around *sql.DB so that additional information, such as the
// name of the server and current database name can be recorded.
type DB struct {
	*sql.DB
	ConnectionDetail
}

// Open will attempt to open a connection to a Postgres Database as specified
// by the connection string provided, it will then ping the database to see if
// the connection is valid.
func Open(connection string) (*DB, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	cd, err := ParseDetails(connection)
	if err != nil {
		return nil, err
	}

	return &DB{db, cd}, nil
}