# sqlmon

[![CircleCI](https://circleci.com/gh/tetafro/sqlmon.svg?style=shield)](https://circleci.com/gh/tetafro/sqlmon)
[![Codecov](https://codecov.io/gh/tetafro/sqlmon/branch/master/graph/badge.svg)](https://codecov.io/gh/tetafro/sqlmon)
[![Go Report](https://goreportcard.com/badge/github.com/tetafro/sqlmon)](https://goreportcard.com/report/github.com/tetafro/sqlmon)
[![Release](https://img.shields.io/github/tag/tetafro/sqlmon.svg)](https://github.com/tetafro/sqlmon/releases)

Wrapper for `database/sql` with callbacks for driver operations.

## Usage

Here is an example of using this wrapper with `lib/pq` driver
```go
import (
    "fmt"
    "database/sql"

    "github.com/lib/pq"
    "github.com/tetafro/sqlmon"
)

func main() {
    drv := sqlmon.Wrap(&pq.Driver{})
    drv.RegisterCallback(OnStmtQuery, func(op string, dur time.Duration, err error) {
        fmt.Println(op, dur, err)
    })
    sql.Register("postgres-wrapped", drv)

    db, err := sql.Open("postgres-wrapped", "host=locahost dbname=testdb")
}
```
