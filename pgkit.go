package pgkit

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

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

// Copy returns a deep copy of the ConnectionDetail so it can be manipulated
// without impacting the original.
func (cd *ConnectionDetail) Copy() ConnectionDetail {
	ncd := NewConnectionDetail()
	ncd.User = cd.User
	ncd.Password = cd.Password
	ncd.Location = cd.Location
	ncd.Port = cd.Port
	ncd.Database = cd.Database

	for key, value := range cd.Options {
		ncd.Options[key] = value
	}

	return ncd
}

// IsValid indicates if the information in the ConnectionDetail can be used to
// construct a valid connection.
func (cd *ConnectionDetail) IsValid() bool {
	// Password is set without user
	if cd.Password != "" && cd.User == "" {
		return false
	}

	// Port is set without location
	if cd.Port != "" && cd.Location == "" {
		return false
	}

	// Option without a key
	if _, ok := cd.Options[""]; ok {
		return false
	}

	// Port is non numeric
	if _, err := strconv.Atoi(cd.Port); err != nil {
		return false
	}

	return true
}

// Returns a URL with the same values as the ConnectionDetail represents. It
// will not check if the returned URL is valid, that is left up to the user.
func (cd *ConnectionDetail) String() string {
	var str strings.Builder
	str.WriteString("postgresql://")

	if cd.User != "" {
		str.WriteString(cd.User)

		if cd.Password != "" {
			str.WriteString(fmt.Sprintf(":%s", cd.Password))
		}

		str.WriteString("@")
	}

	if cd.Location != "" {
		str.WriteString(cd.Location)

		if cd.Port != "" {
			str.WriteString(fmt.Sprintf(":%s", cd.Port))
		}
	}

	if cd.Database != "" {
		str.WriteString(fmt.Sprintf("/%s", cd.Database))
	}

	if len(cd.Options) > 0 {
		str.WriteString("?")

		var sub strings.Builder
		for key, value := range cd.Options {
			sub.WriteString(fmt.Sprintf("%s=%s&", key, value))
		}

		options := sub.String()
		str.WriteString(options[:len(options)-1])
	}

	return str.String()
}

// ParseDetails extracts the connection details out of the connection URI.
func ParseDetails(connection string) (ConnectionDetail, error) {
	cd := NewConnectionDetail()

	u, err := url.Parse(connection)
	if err != nil {
		return cd, err
	}

	cd.User = u.User.Username()
	cd.Location = u.Hostname()
	cd.Port = u.Port()

	if path := u.Path; path != "" {
		cd.Database = u.Path[1:]
	}

	if password, hasPassword := u.User.Password(); hasPassword {
		cd.Password = password
	}

	for key, value := range u.Query() {
		cd.Options[key] = value[0]
	}

	return cd, nil
}

// DB wraps around *sql.DB so that additional information, such as the
// name of the server and current database name can be recorded.
type DB struct {
	*sql.DB
	Connection ConnectionDetail
}

// Open will attempt to open a connection to a Postgres Database as specified
// by the connection string provided, it will then ping the database to see if
// the connection is valid.
func Open(connection string) (DB, error) {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return DB{}, err
	}

	err = db.Ping()
	if err != nil {
		return DB{}, err
	}

	cd, err := ParseDetails(connection)
	if err != nil {
		return DB{}, err
	}

	return DB{db, cd}, nil
}
