package testing

import (
	"crypto/rand"
	"fmt"
	"testing"

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

// generateRandomName returns a random string of characters to provide names
// for temporary objects.
func generateRandomName(t *testing.T) string {
	t.Helper()
	name := make([]byte, 24)
	if _, err := rand.Read(name); err != nil {
		t.Fatalf(err.Error())
	}

	return string(name)
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
	name := generateRandomName(t)
	cmd := fmt.Sprintf("CREATE DATABASE %s;", name)
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

	defer ndb.Close()

	fn(t, ndb)

	cmd = fmt.Sprintf("DROP DATABASE %sl;", name)
	_, err = db.Exec(cmd)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
