package pgunit_test

import (
	"os"
	"testing"

	"github.com/JordanOcokoljic/pgkit"
	"github.com/JordanOcokoljic/pgkit/pgunit"
)

func TestTransactionedTestCase(t *testing.T) {
	db, err := pgkit.Open(os.Getenv("PGKIT_TEST_URL"))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer db.Close()

	preTableNames := pgunit.GetSchemaTableNames(t, db, "public")

	pgunit.TransactionedTestCase(
		t, db,
		func(s *testing.T, dp pgkit.DataProvider) {
			_, err := dp.Exec("CREATE TABLE testing (pk INT)")
			if err != nil {
				s.Fatalf(err.Error())
			}
		},
	)

	postTableNames := pgunit.GetSchemaTableNames(t, db, "public")

	if len(postTableNames) != len(preTableNames) {
		t.Fatalf("table name slice lengths did not match")
	}

	for i := 0; i < len(postTableNames); i++ {
		if preTableNames[i] != postTableNames[i] {
			t.Fatalf("database was modified")
		}
	}
}

func TestTemporaryDatabaseTestCase(t *testing.T) {
	db, err := pgkit.Open(os.Getenv("PGKIT_TEST_URL"))
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer db.Close()

	var dbConn pgkit.ConnectionDetail

	pgunit.TemporaryDatabaseTestCase(t, db, func(s *testing.T, sdb pgkit.DB) {
		_, err := sdb.Exec("CREATE TABLE pgkit (pk INT)")
		if err != nil {
			s.Fatalf(err.Error())
		}

		_, err = sdb.Exec("DROP TABLE pgkit")
		if err != nil {
			s.Fatalf(err.Error())
		}

		dbConn = sdb.Connection.Copy()
	})

	ndb, err := pgkit.Open(dbConn.String())
	if err == nil {
		ndb.Close()
		t.Fatal("database was accessible after subtest")
	}
}

func TestDatabaseGetsTornDownIfFatalOccurs(t *testing.T) {
	db, err := pgkit.Open(os.Getenv("PGKIT_TEST_URL"))
	if err != nil {
		t.Fatal(err.Error())
	}

	var dbConn pgkit.ConnectionDetail
	defer func() {
		db.Close()

		if r := recover(); r == nil {
			t.Error("did not recover")
		}

		ndb, err := pgkit.Open(dbConn.String())
		if err == nil {
			ndb.Close()
			t.Error("database was accessible after teardown")
		}
	}()

	pgunit.TemporaryDatabaseTestCase(t, db, func(s *testing.T, sdb pgkit.DB) {
		dbConn = sdb.Connection.Copy()
		panic("panic")
	})
}
