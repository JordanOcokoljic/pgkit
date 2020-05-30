# pgKit

pgKit simplifies using Postgres in your Go applications. It does this by acting
as a wrapper around your connection to the database. It uses `lib/pq` under the
hood. Additionally, it provides some convience utilities for handling testing
with your databases - such as providing a way to setup transactioned test cases
and temporary databases with your tests.

## Documentation
Documentation can be found [here](https://pkg.go.dev/github.com/JordanOcokoljic/pgkit).

## Testing
pgKit has unit tests, they can be run with:
``` bash
go test
```