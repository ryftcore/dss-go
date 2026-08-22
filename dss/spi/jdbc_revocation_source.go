// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/JdbcRevocationSource.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.client.jdbc's JdbcCacheConnector/SqlQuery/SqlSelectQuery/SqlRecord are
// ported in the sibling package dss/spi/client/jdbc, but with concrete struct shapes that do not
// match what this class needs (it is typed against them in every abstract method and half the
// concrete ones). The minimal interfaces are therefore declared below as structural stand-ins.
//
// TODO: make dss/spi/client/jdbc's concrete types satisfy these interfaces, or retype this file
// against them, and drop the stand-ins.
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SqlQuery is a structural stand-in for
// eu.europa.esig.dss.spi.client.jdbc.query.SqlQuery; see the file header.
type SqlQuery interface {
}

// SqlSelectQuery is a structural stand-in for
// eu.europa.esig.dss.spi.client.jdbc.query.SqlSelectQuery; see the file header.
type SqlSelectQuery interface {
	SqlQuery
}

// SqlRecord is a structural stand-in for
// eu.europa.esig.dss.spi.client.jdbc.record.SqlRecord, the row type
// buildRevocationTokenFromResult reads from; see the file header.
type SqlRecord interface {
}

// JdbcCacheConnector is a structural stand-in for
// eu.europa.esig.dss.spi.client.jdbc.JdbcCacheConnector, connecting to an SQL database and
// performing queries; see the file header. Method names/shapes are inferred from the four call
// sites in JdbcRevocationSource.java (select, execute, executeThrowable, tableQuery).
type JdbcCacheConnector interface {
	// Select executes query with args and returns the resulting records. Port of
	// select(SqlSelectQuery, Object...).
	Select(query SqlSelectQuery, args ...any) ([]SqlRecord, error)
	// Execute executes query with args. Port of execute(SqlQuery, Object...).
	Execute(query SqlQuery, args ...any) error
	// ExecuteThrowable executes query, surfacing a connection-level failure as an error. Port
	// of executeThrowable(SqlQuery), which declares `throws SQLException`.
	ExecuteThrowable(query SqlQuery) error
	// TableQuery runs query and reports whether it found the table. Port of tableQuery(SqlQuery).
	TableQuery(query SqlQuery) bool
}

// JdbcRevocationSourceOverrides captures what JdbcRevocationSourceBase needs to reach through
// virtual dispatch: the operations Java declares abstract on JdbcRevocationSource<R>.
type JdbcRevocationSourceOverrides[R revocation.Revocation] interface {
	// CreateTableQuery returns the CREATE_TABLE sql query. Port of the abstract
	// getCreateTableQuery().
	CreateTableQuery() SqlQuery
	// TableExistenceQuery returns an sql query to check table existence. Port of the abstract
	// getTableExistenceQuery().
	TableExistenceQuery() SqlQuery
	// DeleteTableQuery returns an sql query to remove a table from DB. Port of the abstract
	// getDeleteTableQuery().
	DeleteTableQuery() SqlQuery
	// InsertRevocationTokenEntryQuery returns an SQL query to insert a new revocation token
	// into a table. Port of the abstract getInsertRevocationTokenEntryQuery().
	InsertRevocationTokenEntryQuery() SqlQuery
	// UpdateRevocationTokenEntryQuery returns an SQL query to update a revocation token in a
	// table. Port of the abstract getUpdateRevocationTokenEntryQuery().
	UpdateRevocationTokenEntryQuery() SqlQuery
	// RemoveRevocationTokenEntryQuery returns an sql query to remove a revocation token from
	// DB. Port of the abstract getRemoveRevocationTokenEntryQuery().
	RemoveRevocationTokenEntryQuery() SqlQuery
	// BuildRevocationTokenFromResult builds a RevocationToken from the obtained record. Port of
	// the abstract buildRevocationTokenFromResult(SqlRecord, CertificateToken, CertificateToken),
	// which declares `throws DSSExternalResourceException` (spi/exception.DSSExternalResourceException,
	// already ported); the Go port keeps the return type the plain error interface so any
	// concrete override can wrap it or another error type without this file importing that
	// package for a type it never itself constructs.
	BuildRevocationTokenFromResult(response SqlRecord, certificateToken, issuerCertificateToken *model.CertificateToken) (RevocationToken[R], error)
	// RevocationDataExtractQuery returns a request to find revocation data. Port of the
	// abstract getRevocationDataExtractQuery().
	RevocationDataExtractQuery() SqlSelectQuery
}

// JdbcRevocationSourceBase retrieves revocation tokens from a JDBC datasource. A concrete
// source embeds it and registers itself with InitJdbcRevocationSource; the outstanding
// RepositoryRevocationSourceOverrides methods (InitRevocationTokenKeys, RevocationAccessURLs,
// RevocationTokenKey) still need to come from that concrete source, since JdbcRevocationSource
// itself is abstract on those in Java too.
type JdbcRevocationSourceBase[R revocation.Revocation] struct {
	RepositoryRevocationSourceBase[R]

	// overrides points back at the concrete source; see InitJdbcRevocationSource.
	overrides JdbcRevocationSourceOverrides[R]

	// jdbcCacheConnector connects to the SQL database and performs queries. Java marks the
	// field transient; the Go port carries no such marker since Go values are not
	// java.io.Serializable.
	jdbcCacheConnector JdbcCacheConnector
}

// NewJdbcRevocationSourceBase builds an empty JDBC revocation source. Port of the protected
// default constructor.
func NewJdbcRevocationSourceBase[R revocation.Revocation]() JdbcRevocationSourceBase[R] {
	return JdbcRevocationSourceBase[R]{RepositoryRevocationSourceBase: NewRepositoryRevocationSourceBase[R]()}
}

// InitJdbcRevocationSource registers the concrete source with its base, and forwards the
// registration to the embedded RepositoryRevocationSourceBase (overrides must also implement
// RepositoryRevocationSourceOverrides - the concrete source's own InitRevocationTokenKeys,
// RevocationAccessURLs and RevocationTokenKey combine with FindRevocations/RemoveRevocation
// promoted from this struct; InsertRevocation/UpdateRevocation stay unimplemented here, exactly
// as in Java, which leaves the cache write path to be added by whatever populates the table).
// It must be called by the outermost concrete source's constructor before the source is used.
func (s *JdbcRevocationSourceBase[R]) InitJdbcRevocationSource(overrides JdbcRevocationSourceOverrides[R]) {
	s.overrides = overrides
	repositoryOverrides, ok := overrides.(RepositoryRevocationSourceOverrides[R])
	if !ok {
		panic("JdbcRevocationSource was not initialised: the concrete source must also implement " +
			"RepositoryRevocationSourceOverrides (InitRevocationTokenKeys, RevocationAccessURLs, RevocationTokenKey)")
	}
	s.RepositoryRevocationSourceBase.InitRepositoryRevocationSource(repositoryOverrides)
}

// jdbcRevocationSourceBaseOverrides returns the registered overrides, panicking when the
// concrete source forgot to call InitJdbcRevocationSource.
func (s *JdbcRevocationSourceBase[R]) jdbcRevocationSourceBaseOverrides() JdbcRevocationSourceOverrides[R] {
	if s.overrides == nil {
		panic("JdbcRevocationSource was not initialised: the concrete source must call InitJdbcRevocationSource in its constructor")
	}
	return s.overrides
}

// JdbcCacheConnector gets the SQL connection DataSource. Port of the protected
// getJdbcCacheConnector().
//
// Panics with the Java message when it has not been set (Objects.requireNonNull).
func (s *JdbcRevocationSourceBase[R]) JdbcCacheConnector() JdbcCacheConnector {
	if s.jdbcCacheConnector == nil {
		panic("JdbcCacheConnector shall be provided! Use SetJdbcCacheConnector(jdbcCacheConnector) method.")
	}
	return s.jdbcCacheConnector
}

// SetJdbcCacheConnector sets the SQL connection DataSource. Port of
// setJdbcCacheConnector(JdbcCacheConnector).
func (s *JdbcRevocationSourceBase[R]) SetJdbcCacheConnector(jdbcCacheConnector JdbcCacheConnector) {
	s.jdbcCacheConnector = jdbcCacheConnector
}

// FindRevocations queries the revocation data extract query for key and builds a RevocationToken
// per returned record. Port of the findRevocations(String, CertificateToken, CertificateToken)
// override.
//
// Java declares `throws DSSExternalResourceException` here implicitly through
// buildRevocationTokenFromResult but does not catch it, so it propagates unchecked out of
// findRevocations; the Go port panics with the same error for the same reason - RepositoryRevocationSourceOverrides.FindRevocations
// has no error channel to return it through.
func (s *JdbcRevocationSourceBase[R]) FindRevocations(key string, certificateToken, issuerCertificateToken *model.CertificateToken) []RevocationToken[R] {
	overrides := s.jdbcRevocationSourceBaseOverrides()
	responses, err := s.JdbcCacheConnector().Select(overrides.RevocationDataExtractQuery(), key)
	if err != nil {
		panic(err)
	}
	// Upstream logs "Record obtained : {}"
	if utils.IsCollectionNotEmpty(responses) {
		return s.revocationDataFromRecords(responses, certificateToken, issuerCertificateToken)
	}
	return nil
}

// revocationDataFromRecords builds a RevocationToken per record, skipping the ones that yield
// none. Port of the private getRevocationDataFromRecords(Collection<SqlRecord>, CertificateToken,
// CertificateToken).
func (s *JdbcRevocationSourceBase[R]) revocationDataFromRecords(records []SqlRecord, certificateToken, issuerCertificateToken *model.CertificateToken) []RevocationToken[R] {
	overrides := s.jdbcRevocationSourceBaseOverrides()
	var revocationTokens []RevocationToken[R]
	for _, sqlRecord := range records {
		revocationToken, err := overrides.BuildRevocationTokenFromResult(sqlRecord, certificateToken, issuerCertificateToken)
		if err != nil {
			panic(err)
		}
		if revocationToken != nil {
			revocationTokens = append(revocationTokens, revocationToken)
		}
	}
	return revocationTokens
}

// RemoveRevocation removes the entry for revocationTokenKey from the database. Port of the
// removeRevocation(String) override.
func (s *JdbcRevocationSourceBase[R]) RemoveRevocation(revocationTokenKey string) {
	overrides := s.jdbcRevocationSourceBaseOverrides()
	if err := s.JdbcCacheConnector().Execute(overrides.RemoveRevocationTokenEntryQuery(), revocationTokenKey); err != nil {
		panic(err)
	}
}

// InitTable initializes the revocation token table by creating it if it does not exist. Port of
// initTable(), which declares `throws SQLException`.
func (s *JdbcRevocationSourceBase[R]) InitTable() error {
	tableExists, err := s.IsTableExists()
	if err != nil {
		return err
	}
	if !tableExists {
		// Upstream logs "Table does not exist. Creating a new table..."
		if err := s.createTable(); err != nil {
			return err
		}
		// Upstream logs "Table was created."
	}
	// Upstream logs "Table already exists." when it does.
	return nil
}

// createTable creates the revocation token table. Port of the private createTable(), which
// declares `throws SQLException`.
func (s *JdbcRevocationSourceBase[R]) createTable() error {
	return s.JdbcCacheConnector().ExecuteThrowable(s.jdbcRevocationSourceBaseOverrides().CreateTableQuery())
}

// IsTableExists verifies if the table exists. Port of isTableExists().
//
// Java's tableQuery(SqlQuery) declares no throws clause; the Go JdbcCacheConnector.TableQuery
// mirrors that (no error return), so this always returns a nil error - the signature only
// exists for symmetry with InitTable/DestroyTable, whose underlying calls do declare
// `throws SQLException`.
func (s *JdbcRevocationSourceBase[R]) IsTableExists() (bool, error) {
	return s.JdbcCacheConnector().TableQuery(s.jdbcRevocationSourceBaseOverrides().TableExistenceQuery()), nil
}

// DestroyTable removes the table from the database. Port of destroyTable(), which declares
// `throws SQLException`.
func (s *JdbcRevocationSourceBase[R]) DestroyTable() error {
	tableExists, err := s.IsTableExists()
	if err != nil {
		return err
	}
	if tableExists {
		// Upstream logs "Table exists. Removing the table..."
		if err := s.dropTable(); err != nil {
			return err
		}
		// Upstream logs "Table was destroyed."
		return nil
	}
	// Upstream logs "Cannot drop the table. Table does not exist."
	return nil
}

// dropTable drops the revocation token table. Port of the private dropTable(), which declares
// `throws SQLException`.
func (s *JdbcRevocationSourceBase[R]) dropTable() error {
	return s.JdbcCacheConnector().ExecuteThrowable(s.jdbcRevocationSourceBaseOverrides().DeleteTableQuery())
}

// compile-time assertion: a JdbcRevocationSourceBase is a RevocationSource and a
// MultipleRevocationSource (promoted from RepositoryRevocationSourceBase).
var (
	_ RevocationSource[revocation.CRL]         = (*JdbcRevocationSourceBase[revocation.CRL])(nil)
	_ MultipleRevocationSource[revocation.CRL] = (*JdbcRevocationSourceBase[revocation.CRL])(nil)
)
