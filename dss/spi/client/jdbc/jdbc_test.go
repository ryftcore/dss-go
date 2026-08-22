// Tests for the jdbc package. There is no upstream JUnit for
// CacheConnector/SqlQuery/SqlSelectQuery to port test vectors from
// (see dss-spi client/jdbc upstream), so behavior is
// verified against a lightweight in-process database/sql/driver fake
// registered below, exercising the same success/rollback contracts the
// Java implementation documents (see jdbc_cache_connector.go).
package jdbc

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync"
	"testing"
)

// ---- fake database/sql driver -------------------------------------------

type fakeBehavior struct {
	beginErr    error
	execFunc    func(query string, args []driver.Value) (driver.Result, error)
	queryFunc   func(query string, args []driver.Value) (driver.Rows, error)
	commitErr   error
	rollbackErr error
}

var (
	fakeRegistryMu sync.Mutex
	fakeRegistry   = map[string]*fakeBehavior{}
	fakeRegistered bool
)

func registerFakeDriverOnce() {
	fakeRegistryMu.Lock()
	defer fakeRegistryMu.Unlock()
	if !fakeRegistered {
		sql.Register("fakejdbc", fakeDriver{})
		fakeRegistered = true
	}
}

// openFakeDB registers b under a fresh unique DSN and returns a *sql.DB
// backed by it.
func openFakeDB(t *testing.T, b *fakeBehavior) *sql.DB {
	t.Helper()
	registerFakeDriverOnce()
	dsn := t.Name()
	fakeRegistryMu.Lock()
	fakeRegistry[dsn] = b
	fakeRegistryMu.Unlock()

	db, err := sql.Open("fakejdbc", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

type fakeDriver struct{}

func (fakeDriver) Open(name string) (driver.Conn, error) {
	fakeRegistryMu.Lock()
	b, ok := fakeRegistry[name]
	fakeRegistryMu.Unlock()
	if !ok {
		return nil, errors.New("fakejdbc: unknown dsn " + name)
	}
	return &fakeConn{behavior: b}, nil
}

type fakeConn struct {
	behavior *fakeBehavior
}

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeStmt{conn: c, query: query}, nil
}

func (c *fakeConn) Close() error { return nil }

func (c *fakeConn) Begin() (driver.Tx, error) {
	if c.behavior.beginErr != nil {
		return nil, c.behavior.beginErr
	}
	return &fakeTx{behavior: c.behavior}, nil
}

type fakeTx struct {
	behavior *fakeBehavior
}

func (tx *fakeTx) Commit() error   { return tx.behavior.commitErr }
func (tx *fakeTx) Rollback() error { return tx.behavior.rollbackErr }

type fakeStmt struct {
	conn  *fakeConn
	query string
}

func (s *fakeStmt) Close() error  { return nil }
func (s *fakeStmt) NumInput() int { return -1 }

func (s *fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	if s.conn.behavior.execFunc == nil {
		return fakeResult{}, nil
	}
	return s.conn.behavior.execFunc(s.query, args)
}

func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	if s.conn.behavior.queryFunc == nil {
		return &fakeRows{}, nil
	}
	return s.conn.behavior.queryFunc(s.query, args)
}

type fakeResult struct {
	lastInsertID int64
	rowsAffected int64
	rowsErr      error
}

func (r fakeResult) LastInsertId() (int64, error) { return r.lastInsertID, nil }
func (r fakeResult) RowsAffected() (int64, error) {
	if r.rowsErr != nil {
		return 0, r.rowsErr
	}
	return r.rowsAffected, nil
}

type fakeRows struct {
	cols []string
	data [][]driver.Value
	idx  int
}

func (r *fakeRows) Columns() []string { return r.cols }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.data) {
		return io.EOF
	}
	copy(dest, r.data[r.idx])
	r.idx++
	return nil
}

// ---- SqlQuery -------------------------------------------------------------

func TestNewSqlQuery(t *testing.T) {
	q := NewSqlQuery("SELECT 1")
	if got := q.QueryString(); got != "SELECT 1" {
		t.Fatalf("QueryString() = %q, want %q", got, "SELECT 1")
	}
	want := "JdbcQuery[queryString='SELECT 1']"
	if got := q.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

// ---- SqlSelectQuery / GetRecords -------------------------------------------

type fakeRecord struct{ v int }

type fakeSelectQuery struct {
	*SqlQuery
	getRecordErr error
}

func (q *fakeSelectQuery) GetRecord(rs *sql.Rows) (SqlRecord, error) {
	if q.getRecordErr != nil {
		return nil, q.getRecordErr
	}
	var v int
	if err := rs.Scan(&v); err != nil {
		return nil, err
	}
	return fakeRecord{v: v}, nil
}

func newFakeSelectQuery(queryString string) *fakeSelectQuery {
	return &fakeSelectQuery{SqlQuery: NewSqlQuery(queryString)}
}

func TestGetRecords(t *testing.T) {
	b := &fakeBehavior{
		queryFunc: func(query string, args []driver.Value) (driver.Rows, error) {
			return &fakeRows{
				cols: []string{"v"},
				data: [][]driver.Value{{int64(1)}, {int64(2)}, {int64(3)}},
			}, nil
		},
	}
	db := openFakeDB(t, b)

	rows, err := db.Query("SELECT v")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	q := newFakeSelectQuery("SELECT v")
	records, err := GetRecords(q, rows)
	if err != nil {
		t.Fatalf("GetRecords: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("len(records) = %d, want 3", len(records))
	}
	for i, want := range []int{1, 2, 3} {
		got, ok := records[i].(fakeRecord)
		if !ok || got.v != want {
			t.Fatalf("records[%d] = %v, want fakeRecord{%d}", i, records[i], want)
		}
	}
}

func TestGetRecordsPropagatesGetRecordError(t *testing.T) {
	wantErr := errors.New("boom")
	b := &fakeBehavior{
		queryFunc: func(query string, args []driver.Value) (driver.Rows, error) {
			return &fakeRows{cols: []string{"v"}, data: [][]driver.Value{{int64(1)}}}, nil
		},
	}
	db := openFakeDB(t, b)
	rows, err := db.Query("SELECT v")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	q := newFakeSelectQuery("SELECT v")
	q.getRecordErr = wantErr
	if _, err := GetRecords(q, rows); !errors.Is(err, wantErr) {
		t.Fatalf("GetRecords error = %v, want %v", err, wantErr)
	}
}

// ---- CacheConnector.Execute --------------------------------------------

func TestJdbcCacheConnectorExecuteNilQueryPanics(t *testing.T) {
	c := NewCacheConnector(openFakeDB(t, &fakeBehavior{}))
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil query")
		}
	}()
	c.Execute(nil)
}

func TestJdbcCacheConnectorExecuteSuccess(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			if query != "UPDATE t SET a=1" {
				t.Fatalf("unexpected query %q", query)
			}
			return fakeResult{rowsAffected: 5}, nil
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	got := c.Execute(NewSqlQuery("UPDATE t SET a=1"))
	if got != 5 {
		t.Fatalf("Execute() = %d, want 5", got)
	}
}

func TestJdbcCacheConnectorExecuteBeginError(t *testing.T) {
	b := &fakeBehavior{beginErr: errors.New("no conn")}
	c := NewCacheConnector(openFakeDB(t, b))
	if got := c.Execute(NewSqlQuery("UPDATE t")); got != 0 {
		t.Fatalf("Execute() = %d, want 0", got)
	}
}

func TestJdbcCacheConnectorExecuteExecErrorRollsBack(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return nil, errors.New("exec failed")
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if got := c.Execute(NewSqlQuery("UPDATE t")); got != 0 {
		t.Fatalf("Execute() = %d, want 0", got)
	}
}

func TestJdbcCacheConnectorExecuteCommitError(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return fakeResult{rowsAffected: 1}, nil
		},
		commitErr: errors.New("commit failed"),
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if got := c.Execute(NewSqlQuery("UPDATE t")); got != 0 {
		t.Fatalf("Execute() = %d, want 0", got)
	}
}

// ---- CacheConnector.Select --------------------------------------------

func TestJdbcCacheConnectorSelectSuccess(t *testing.T) {
	b := &fakeBehavior{
		queryFunc: func(query string, args []driver.Value) (driver.Rows, error) {
			return &fakeRows{cols: []string{"v"}, data: [][]driver.Value{{int64(7)}, {int64(8)}}}, nil
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	records := c.Select(newFakeSelectQuery("SELECT v"))
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}
}

func TestJdbcCacheConnectorSelectBeginError(t *testing.T) {
	b := &fakeBehavior{beginErr: errors.New("no conn")}
	c := NewCacheConnector(openFakeDB(t, b))
	if records := c.Select(newFakeSelectQuery("SELECT v")); records != nil {
		t.Fatalf("Select() = %v, want nil", records)
	}
}

func TestJdbcCacheConnectorSelectQueryErrorRollsBack(t *testing.T) {
	b := &fakeBehavior{
		queryFunc: func(query string, args []driver.Value) (driver.Rows, error) {
			return nil, errors.New("query failed")
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if records := c.Select(newFakeSelectQuery("SELECT v")); records != nil {
		t.Fatalf("Select() = %v, want nil", records)
	}
}

func TestJdbcCacheConnectorSelectCommitError(t *testing.T) {
	b := &fakeBehavior{
		queryFunc: func(query string, args []driver.Value) (driver.Rows, error) {
			return &fakeRows{cols: []string{"v"}, data: [][]driver.Value{{int64(1)}}}, nil
		},
		commitErr: errors.New("commit failed"),
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if records := c.Select(newFakeSelectQuery("SELECT v")); records != nil {
		t.Fatalf("Select() = %v, want nil", records)
	}
}

// ---- CacheConnector.TableQuery -----------------------------------------

func TestJdbcCacheConnectorTableQuerySuccess(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return fakeResult{}, nil
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if !c.TableQuery(NewSqlQuery("CREATE TABLE t (a int)")) {
		t.Fatal("TableQuery() = false, want true")
	}
}

func TestJdbcCacheConnectorTableQueryBeginError(t *testing.T) {
	b := &fakeBehavior{beginErr: errors.New("no conn")}
	c := NewCacheConnector(openFakeDB(t, b))
	if c.TableQuery(NewSqlQuery("CREATE TABLE t (a int)")) {
		t.Fatal("TableQuery() = true, want false")
	}
}

func TestJdbcCacheConnectorTableQueryExecError(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return nil, errors.New("exec failed")
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if c.TableQuery(NewSqlQuery("DROP TABLE t")) {
		t.Fatal("TableQuery() = true, want false")
	}
}

func TestJdbcCacheConnectorTableQueryCommitError(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return fakeResult{}, nil
		},
		commitErr: errors.New("commit failed"),
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if c.TableQuery(NewSqlQuery("DROP TABLE t")) {
		t.Fatal("TableQuery() = true, want false")
	}
}

// ---- CacheConnector.ExecuteThrowable -----------------------------------

func TestJdbcCacheConnectorExecuteThrowableSuccess(t *testing.T) {
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return fakeResult{rowsAffected: 3}, nil
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	n, err := c.ExecuteThrowable(NewSqlQuery("UPDATE t"))
	if err != nil {
		t.Fatalf("ExecuteThrowable() error = %v", err)
	}
	if n != 3 {
		t.Fatalf("ExecuteThrowable() = %d, want 3", n)
	}
}

func TestJdbcCacheConnectorExecuteThrowableBeginError(t *testing.T) {
	wantErr := errors.New("no conn")
	b := &fakeBehavior{beginErr: wantErr}
	c := NewCacheConnector(openFakeDB(t, b))
	if _, err := c.ExecuteThrowable(NewSqlQuery("UPDATE t")); err == nil {
		t.Fatal("expected error")
	}
}

func TestJdbcCacheConnectorExecuteThrowableExecErrorRollsBack(t *testing.T) {
	wantErr := errors.New("exec failed")
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return nil, wantErr
		},
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if _, err := c.ExecuteThrowable(NewSqlQuery("UPDATE t")); err == nil {
		t.Fatal("expected error")
	}
}

func TestJdbcCacheConnectorExecuteThrowableCommitError(t *testing.T) {
	wantErr := errors.New("commit failed")
	b := &fakeBehavior{
		execFunc: func(query string, args []driver.Value) (driver.Result, error) {
			return fakeResult{rowsAffected: 1}, nil
		},
		commitErr: wantErr,
	}
	c := NewCacheConnector(openFakeDB(t, b))
	if _, err := c.ExecuteThrowable(NewSqlQuery("UPDATE t")); err == nil {
		t.Fatal("expected error")
	}
}
