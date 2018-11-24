// Package sqlmon provides a wrapper for standard SQL driver
// that adds ability to create callbacks for driver operations.
package sqlmon

import (
	"database/sql/driver"
	"time"
)

// Driver is a wrapper for driver.Driver, that has a set of callbacks
// for different operations.
type Driver struct {
	origin    driver.Driver
	callbacks map[string]Callback
}

// Wrap wraps given database driver and returns a driver, that can
// uses callbacks for operations.
func Wrap(drv driver.Driver) *Driver {
	return &Driver{origin: drv, callbacks: make(map[string]Callback)}
}

// Callback is a function that is executed at a certain moment.
type Callback func(op string, dur time.Duration, err error)

// Set of available callback types.
const (
	OnDriverOpen  = "driver.Open"
	OnConnBegin   = "conn.Begin"
	OnConnClose   = "conn.Close"
	OnConnPrepare = "conn.Prepare"
	OnTxCommit    = "tx.Commit"
	OnTxRollback  = "tx.Rollback"
	OnStmtClose   = "stmt.Close"
	OnStmtExec    = "stmt.Exec"
	OnStmtQuery   = "stmt.Query"
)

// RegisterCallback adds callback function to driver.
func (d *Driver) RegisterCallback(typ string, cb Callback) {
	d.callbacks[typ] = cb
}

// Open returns a new connection to the database.
func (d *Driver) Open(name string) (driver.Conn, error) {
	t := time.Now()
	cn, err := d.origin.Open(name)

	if cb, ok := d.callbacks[OnDriverOpen]; ok {
		cb(OnDriverOpen, time.Since(t), err)
	}

	return &Conn{driver: d, origin: cn}, nil
}

// Conn is a wrapper for driver.Conn, that can execute callbacks
// for different operations.
type Conn struct {
	driver *Driver
	origin driver.Conn
}

// Begin implements driver.Conn interface.
func (cn *Conn) Begin() (driver.Tx, error) {
	t := time.Now()
	tx, err := cn.origin.Begin()

	if cb, ok := cn.driver.callbacks[OnConnBegin]; ok {
		cb(OnConnBegin, time.Since(t), err)
	}

	if err != nil {
		return tx, err
	}
	return &Tx{driver: cn.driver, origin: tx}, nil
}

// Close implements driver.Conn interface.
func (cn *Conn) Close() error {
	t := time.Now()
	err := cn.origin.Close()

	if cb, ok := cn.driver.callbacks[OnConnClose]; ok {
		cb(OnConnClose, time.Since(t), err)
	}

	return err
}

// Prepare implements driver.Conn interface.
func (cn *Conn) Prepare(q string) (driver.Stmt, error) {
	t := time.Now()
	stmt, err := cn.origin.Prepare(q)

	if cb, ok := cn.driver.callbacks[OnConnPrepare]; ok {
		cb(OnConnPrepare, time.Since(t), err)
	}

	return &Stmt{driver: cn.driver, origin: stmt}, nil
}

// Tx is a wrapper for driver.Tx, that can execute callbacks
// for different operations.
type Tx struct {
	driver *Driver
	origin driver.Tx
}

// Commit implements driver.Tx interface.
func (tx *Tx) Commit() error {
	t := time.Now()
	err := tx.origin.Commit()

	if cb, ok := tx.driver.callbacks[OnTxCommit]; ok {
		cb(OnTxCommit, time.Since(t), err)
	}

	return err
}

// Rollback implements driver.Tx interface.
func (tx *Tx) Rollback() error {
	t := time.Now()
	err := tx.origin.Rollback()

	if cb, ok := tx.driver.callbacks[OnTxRollback]; ok {
		cb(OnTxRollback, time.Since(t), err)
	}

	return err
}

// Stmt is a wrapper for driver.Stmt, that can execute callbacks
// for different operations.
type Stmt struct {
	driver *Driver
	origin driver.Stmt
}

// Close implements driver.Stmt interface.
func (stmt *Stmt) Close() error {
	t := time.Now()
	err := stmt.origin.Close()

	if cb, ok := stmt.driver.callbacks[OnConnBegin]; ok {
		cb(OnConnBegin, time.Since(t), err)
	}

	return err
}

// NumInput implements driver.Stmt interface.
func (stmt *Stmt) NumInput() int {
	return stmt.origin.NumInput()
}

// Exec implements driver.Stmt interface.
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	t := time.Now()
	res, err := stmt.origin.Exec(args)

	if cb, ok := stmt.driver.callbacks[OnStmtExec]; ok {
		cb(OnStmtExec, time.Since(t), err)
	}

	return res, err
}

// Query implements driver.Stmt interface.
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	t := time.Now()
	rows, err := stmt.origin.Query(args)

	if cb, ok := stmt.driver.callbacks[OnStmtQuery]; ok {
		cb(OnStmtQuery, time.Since(t), err)
	}

	return rows, err
}
