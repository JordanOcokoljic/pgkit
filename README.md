# pgKit

pgKit simplifies using Postgres in your Go applications. It does this by acting
as a wrapper around your connection to the database. It uses `lib/pq` under the
hood. Additionally, it provides some convience utilities for handling testing
with your databases - such as providing a way to setup transactioned test cases
and temporary databases with your tests.

## Documentation
Documentation can be found [here](https://pkg.go.dev/github.com/JordanOcokoljic/pgkit).

## Testing
pgKit has unit tests, they require a user and a database to run, to get these
setup:

1. Create a new user with the `CREATEDB` permission.
2. Create a new database.
3. Provide the connection url to the test via the PGKIT_TEST_URL environment
variable.

For example, given the user `pgkit_test` with the password `pgkit` and a the
database `pgkit_test` running on localhost with default port and SSL options,
the provided URL would be:
`postgresql://pgkit_test:pgkit@localhost:5432/pgkit_test?sslmode=disable`

To run the tests then:
``` bash
PGKIT_TEST_URL="postgresql://pgkit_test:pgkit@localhost:5432/pgkit_test?sslmode=disable" go test ./...
```