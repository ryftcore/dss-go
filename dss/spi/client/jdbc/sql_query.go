// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/jdbc/query/SqlQuery.java (DSS 6.5.RC1).
package jdbc

import "fmt"

// SqlQuery represents a stateless query to be made to an SQL database.
type SqlQuery struct {
	// queryString is the executable SQL query string.
	queryString string
}

// NewSqlQuery creates a SqlQuery with the given query string. Ports both
// the protected constructor and the public static factory
// SqlQuery#createQuery(String), which upstream exposes as the sole way to
// build an instance.
func NewSqlQuery(queryString string) *SqlQuery {
	return &SqlQuery{queryString: queryString}
}

// QueryString returns the query string. Ports #getQueryString.
func (q *SqlQuery) QueryString() string {
	return q.queryString
}

// String ports #toString.
func (q *SqlQuery) String() string {
	return fmt.Sprintf("JdbcQuery[queryString='%s']", q.queryString)
}
