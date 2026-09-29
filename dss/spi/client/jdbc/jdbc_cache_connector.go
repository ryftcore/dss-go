// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/jdbc/JdbcCacheConnector.java (DSS 6.5.RC1).
package jdbc

import "database/sql"

// JdbcCacheConnector executes calls to a *sql.DB (Java's
// javax.sql.DataSource; the caller passes an already-configured *sql.DB).
type CacheConnector struct {
	// dataSource is the SQL data source to create connections with.
	dataSource *sql.DB
}

// NewCacheConnector is the default constructor.
func NewCacheConnector(dataSource *sql.DB) *CacheConnector {
	return &CacheConnector{dataSource: dataSource}
}

// Execute executes a query with a custom set of arguments, such as
// UPDATE or DELETE, by handling the error. On failure the transaction is
// rolled back and 0 is returned, mirroring
// #execute(SqlQuery, Object...), which never propagates the SQLException.
func (c *CacheConnector) Execute(query *SqlQuery, arguments ...any) int {
	if query == nil {
		panic("Query cannot be null!")
	}

	tx, err := c.dataSource.Begin()
	if err != nil {
		return 0
	}

	result, err := tx.Exec(query.QueryString(), arguments...)
	if err != nil {
		_ = tx.Rollback()
		return 0
	}

	if err := tx.Commit(); err != nil {
		return 0
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0
	}
	return int(n)
}

// Select executes the query and returns a collection of selected
// objects. Ports #select(SqlSelectQuery, Object...): on any failure it
// rolls back and returns an empty (nil) collection rather than
// propagating the error, matching the Java method's
// Collections.emptySet() fallback.
func (c *CacheConnector) Select(selectQuery SqlSelectQuery, arguments ...any) []SqlRecord {
	tx, err := c.dataSource.Begin()
	if err != nil {
		return nil
	}

	rows, err := tx.Query(selectQuery.QueryString(), arguments...)
	if err != nil {
		_ = tx.Rollback()
		return nil
	}

	records, err := GetRecords(selectQuery, rows)
	_ = rows.Close()
	if err != nil {
		_ = tx.Rollback()
		return nil
	}

	if err := tx.Commit(); err != nil {
		return nil
	}
	return records
}

// TableQuery allows table creation, removal and existence check. Returns
// true if the query executed successfully, false otherwise. Ports
// #tableQuery(SqlQuery).
//
// DEVIATION: Java returns Statement#execute's boolean, which distinguishes
// "first result is a ResultSet" from "first result is an update count" -
// a distinction database/sql's Exec does not expose. Callers of
// tableQuery only rely on the success/failure signal (DDL and existence
// probes), so this returns whether the statement executed without error.
func (c *CacheConnector) TableQuery(query *SqlQuery) bool {
	tx, err := c.dataSource.Begin()
	if err != nil {
		return false
	}

	if _, err := tx.Exec(query.QueryString()); err != nil {
		// Java's tableQuery catches SQLException and returns false without an
		// explicit rollback (closeQuietly closes the connection, which discards
		// the uncommitted transaction). A database/sql transaction must be
		// ended, or it keeps its connection checked out, so it is rolled back
		// here.
		_ = tx.Rollback()
		return false
	}

	if err := tx.Commit(); err != nil {
		return false
	}
	return true
}

// ExecuteThrowable allows executing INSERT, UPDATE or DELETE queries,
// returning an error in case of failure. Ports #executeThrowable(SqlQuery).
func (c *CacheConnector) ExecuteThrowable(query *SqlQuery) (int, error) {
	tx, err := c.dataSource.Begin()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query.QueryString())
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
