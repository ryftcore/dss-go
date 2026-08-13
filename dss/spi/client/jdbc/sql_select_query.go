// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/jdbc/query/SqlSelectQuery.java (DSS 6.5.RC1).
package jdbc

import "database/sql"

// SqlSelectQuery is a query containing logic to extract records from a
// *sql.Rows result set (Java ResultSet). Java models this as an abstract
// class extending SqlQuery with an abstract getRecord(ResultSet); Go lacks
// abstract methods, so it is ported as an interface. Implementations
// typically embed *SqlQuery for QueryString and implement GetRecord; the
// shared #getRecords(ResultSet) loop is ported as the package-level
// GetRecords function below rather than duplicated per implementation.
type SqlSelectQuery interface {
	// QueryString returns the executable SQL query string. Ports the
	// inherited SqlQuery#getQueryString.
	QueryString() string

	// GetRecord returns a response extracted from the given rs at its
	// current row position. Ports the abstract #getRecord(ResultSet).
	GetRecord(rs *sql.Rows) (SqlRecord, error)
}

// GetRecords extracts a collection of SqlRecords from rs by calling
// q.GetRecord for each row. Ports
// SqlSelectQuery#getRecords(ResultSet).
func GetRecords(q SqlSelectQuery, rs *sql.Rows) ([]SqlRecord, error) {
	var records []SqlRecord
	for rs.Next() {
		record, err := q.GetRecord(rs)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
