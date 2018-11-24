package sqlmon

import "database/sql/driver"

type driverFake struct{}

func (d *driverFake) Open(name string) (driver.Conn, error) {
	return &connFake{}, nil
}

type connFake struct{}

func (cn *connFake) Begin() (driver.Tx, error) {
	return &txFake{}, nil
}

func (cn *connFake) Close() error {
	return nil
}

func (cn *connFake) Prepare(q string) (driver.Stmt, error) {
	return &stmtFake{}, nil
}

type txFake struct{}

func (tx *txFake) Commit() error {
	return nil
}

func (tx *txFake) Rollback() error {
	return nil
}

type stmtFake struct{}

func (stmt *stmtFake) Close() error {
	return nil
}

func (stmt *stmtFake) NumInput() int {
	return 0
}

func (stmt *stmtFake) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (stmt *stmtFake) Query(args []driver.Value) (driver.Rows, error) {
	return &rowsFake{}, nil
}

type rowsFake struct{}

func (r *rowsFake) Columns() []string {
	return []string{}
}

func (r *rowsFake) Close() error {
	return nil
}

func (r *rowsFake) Next(dest []driver.Value) error {
	return nil
}
