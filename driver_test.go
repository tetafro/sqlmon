package sqlmon

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestDriverCallbacks(t *testing.T) {
	drv := Wrap(&driverFake{})

	// Register all types of callbacks
	cb := testCallbacker{counter: make(map[string]int)}
	drv.RegisterCallback(OnDriverOpen, cb.callback)
	drv.RegisterCallback(OnConnBegin, cb.callback)
	drv.RegisterCallback(OnConnClose, cb.callback)
	drv.RegisterCallback(OnConnPrepare, cb.callback)
	drv.RegisterCallback(OnTxCommit, cb.callback)
	drv.RegisterCallback(OnTxRollback, cb.callback)
	drv.RegisterCallback(OnStmtClose, cb.callback)
	drv.RegisterCallback(OnStmtExec, cb.callback)
	drv.RegisterCallback(OnStmtQuery, cb.callback)

	sql.Register("wrapped", drv)

	t.Run("open database", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		err = db.Ping()
		if err != nil {
			t.Fatalf("ping error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen: 1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("plain query", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		_, err = db.Query("fake query")
		if err != nil {
			t.Fatalf("query error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnPrepare: 1,
			OnStmtQuery:   1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("plain exec", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		_, err = db.Exec("fake query")
		if err != nil {
			t.Fatalf("exec query error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnBegin:   1,
			OnConnPrepare: 1,
			OnStmtExec:    1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("prepared statement query", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		stmt, err := db.Prepare("fake query")
		if err != nil {
			t.Fatalf("prepare query error: %v", err)
		}
		_, err = stmt.Query()
		if err != nil {
			t.Fatalf("statement query error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnPrepare: 1,
			OnStmtQuery:   1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("prepared statement exec", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		stmt, err := db.Prepare("fake query")
		if err != nil {
			t.Fatalf("prepare query error: %v", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			t.Fatalf("statement exec query error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnPrepare: 1,
			OnStmtExec:    1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("commit prepared statement exec in transaction", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("transaction begin error: %v", err)
		}
		stmt, err := tx.Prepare("fake query")
		if err != nil {
			t.Fatalf("transaction prepare query error: %v", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			t.Fatalf("transaction statement query exec error: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			t.Fatalf("transaction commit error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnBegin:   2,
			OnConnPrepare: 1,
			OnStmtExec:    1,
			OnTxCommit:    1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
		cb.counter = make(map[string]int)
	})

	t.Run("rollback prepared statement exec in transaction", func(t *testing.T) {
		db, err := sql.Open("wrapped", "")
		if err != nil {
			t.Fatalf("connection error: %v", err)
		}
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("transaction begin error: %v", err)
		}
		stmt, err := tx.Prepare("fake query")
		if err != nil {
			t.Fatalf("transaction prepare query error: %v", err)
		}
		_, err = stmt.Exec()
		if err != nil {
			t.Fatalf("transaction statement query exec error: %v", err)
		}
		err = tx.Rollback()
		if err != nil {
			t.Fatalf("transaction rollback error: %v", err)
		}
		err = checkCounters(cb.counter, map[string]int{
			OnDriverOpen:  1,
			OnConnBegin:   2,
			OnConnPrepare: 1,
			OnStmtExec:    1,
			OnTxRollback:  1,
		})
		if err != nil {
			t.Fatalf("wrong callback counters: %v", err)
		}
	})
}

func checkCounters(got, want map[string]int) error {
	for k, v := range want {
		gv, ok := got[k]
		if !ok {
			return fmt.Errorf("missing key: %s", k)
		}
		if gv != v {
			return fmt.Errorf("wrong value for key '%s', got %d, want %d", k, gv, v)
		}
	}
	for k := range got {
		if _, ok := want[k]; !ok {
			return fmt.Errorf("got unexpected key: %s", k)
		}
	}
	return nil
}

type testCallbacker struct {
	counter map[string]int
}

func (cb *testCallbacker) callback(op string, dur time.Duration, err error) {
	cb.counter[op]++
}
