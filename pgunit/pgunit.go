package pgunit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/JordanOcokoljic/pgkit"
)

// TransactionedTestCase will execute the provided function with a transaction
// and rollback once the test is complete. It passes the given *testing.T to
// the function, so failing in the test function triggers the failure of the
// outer test.
func TransactionedTestCase(
	t *testing.T,
	db pgkit.DB,
	fn func(*testing.T, pgkit.DataProvider),
) {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer tx.Rollback()

	fn(t, tx)
}

// letterBytes is used to generate random names for databases.
const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// generateRandomName returns a random string of characters to provide names
// for temporary objects.
func generateRandomName() string {
	b := make([]byte, 12)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

// TemporaryDatabaseTestCase creates a new database and executes the provided
// function passing it a reference to the newly created database. When the test
// is complete, the temporary database is destroyed.
func TemporaryDatabaseTestCase(
	t *testing.T,
	db pgkit.DB,
	fn func(*testing.T, pgkit.DB),
) {
	t.Helper()

	name := generateRandomName()
	cmd := fmt.Sprintf("CREATE DATABASE %s", name)
	_, err := db.Exec(cmd)
	if err != nil {
		t.Fatalf(err.Error())
	}

	connDetail := db.Connection.Copy()
	connDetail.Database = name

	ndb, err := pgkit.Open(connDetail.String())
	if err != nil {
		t.Fatalf(err.Error())
	}

	fn(t, ndb)

	ndb.Close()

	cmd = fmt.Sprintf("DROP DATABASE %s", name)
	_, err = db.Exec(cmd)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// GetSchemaTableNames will return the names of all the tables in the public
// schema.
func GetSchemaTableNames(t *testing.T, db pgkit.DB, schema string) []string {
	t.Helper()

	rows, err := db.Query("SELECT * FROM information_schema.tables WHERE table_schema = $1", schema)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			t.Fatalf(err.Error())
		}

		tables = append(tables, table)
	}

	return tables
}
