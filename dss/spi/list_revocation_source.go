// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ListRevocationSource.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// OfflineRevocationSource is the Go-only interface counterpart of the Java abstract class
// OfflineRevocationSource<R> (ported as the OfflineRevocationSourceBase[R] struct in
// offline_revocation_source.go). ListRevocationSource holds a heterogeneous list of concrete
// offline sources (e.g. a future CRL-backed one and a future OCSP-backed one both instantiated
// for the same R) polymorphically, which in Go requires an interface rather than the embeddable
// struct; every method here is implemented on OfflineRevocationSourceBase[R] and promoted to
// whatever concrete source embeds it.
type OfflineRevocationSource[R revocation.Revocation] interface {
	RevocationSource[R]
	MultipleRevocationSource[R]

	// IsEmpty checks if the current source is empty. Port of isEmpty().
	IsEmpty() bool
	// AllRevocationBinaries retrieves all found revocation binaries. Port of
	// getAllRevocationBinaries().
	AllRevocationBinaries() []EncapsulatedRevocationTokenIdentifier[R]
	// FindBinaryForReference returns the incorporated EncapsulatedRevocationTokenIdentifier
	// corresponding to the provided reference, nil when none matches. Port of
	// findBinaryForReference(RevocationRef).
	FindBinaryForReference(reference RevocationRef[R]) EncapsulatedRevocationTokenIdentifier[R]
	// IsOrphan checks if the source does not contain revocation identifiers matching the
	// reference. Port of isOrphan(RevocationRef).
	IsOrphan(reference RevocationRef[R]) bool
}

// ListRevocationSource allows handling a list of OfflineRevocationSources.
type ListRevocationSource[R revocation.Revocation] struct {
	// sources is the list of revocation sources.
	sources []OfflineRevocationSource[R]
}

// NewListRevocationSource builds an empty list. Port of the default constructor.
func NewListRevocationSource[R revocation.Revocation]() *ListRevocationSource[R] {
	return &ListRevocationSource[R]{}
}

// NewListRevocationSourceFrom initializes the list with an OfflineRevocationSource. Port of the
// ListRevocationSource(OfflineRevocationSource) constructor.
func NewListRevocationSourceFrom[R revocation.Revocation](revocationSource OfflineRevocationSource[R]) *ListRevocationSource[R] {
	list := NewListRevocationSource[R]()
	list.Add(revocationSource)
	return list
}

// Add adds revocationSource to the list, keeping old values. Returns whether it was added
// (Java's contains() check uses OfflineRevocationSource's inherited Object.equals(), i.e.
// reference identity, since the Java class does not override equals(); the interface value
// comparison below is Go's equivalent). Port of add(OfflineRevocationSource).
func (l *ListRevocationSource[R]) Add(revocationSource OfflineRevocationSource[R]) bool {
	if revocationSource == nil {
		return false
	}
	for _, existing := range l.sources {
		if existing == revocationSource {
			return false
		}
	}
	l.sources = append(l.sources, revocationSource)
	return true
}

// AddAllFrom adds all sources from listRevocationSource to the list, keeping old values. Port
// of addAll(ListRevocationSource).
func (l *ListRevocationSource[R]) AddAllFrom(listRevocationSource *ListRevocationSource[R]) {
	l.AddAll(listRevocationSource.Sources())
}

// AddAll adds all revocationSources to the list, keeping old values. Port of
// addAll(List<OfflineRevocationSource>).
//
// Java's addAll(Collection) appends unconditionally, without the add(single) method's dedup
// check; the Go port mirrors that difference faithfully.
func (l *ListRevocationSource[R]) AddAll(revocationSources []OfflineRevocationSource[R]) {
	l.sources = append(l.sources, revocationSources...)
}

// Sources gets a list of all embedded sources. Port of getSources().
func (l *ListRevocationSource[R]) Sources() []OfflineRevocationSource[R] {
	return l.sources
}

// IsEmpty checks if the current ListRevocationSource and its children are empty. Port of
// isEmpty().
func (l *ListRevocationSource[R]) IsEmpty() bool {
	for _, revocationSource := range l.sources {
		if !revocationSource.IsEmpty() {
			return false
		}
	}
	return true
}

// RevocationTokens retrieves the union of the revocation tokens from every embedded source, for
// the given certificate/issuer couple. Port of getRevocationTokens(CertificateToken,
// CertificateToken).
//
// Java collects into a HashSet<RevocationToken<R>>, deduplicating with RevocationToken#equals();
// the Go port dedups with RevocationTokenBase.Equals through the same linear-scan approach used
// throughout this port for types with no comparable Go representation.
func (l *ListRevocationSource[R]) RevocationTokens(certificateToken, issuerCertificateToken *model.CertificateToken) ([]RevocationToken[R], error) {
	var allTokens []RevocationToken[R]
	for _, revocationSource := range l.sources {
		tokens, err := revocationSource.RevocationTokens(certificateToken, issuerCertificateToken)
		if err != nil {
			return nil, err
		}
		for _, token := range tokens {
			if !listRevocationSourceContainsToken(allTokens, token) {
				allTokens = append(allTokens, token)
			}
		}
	}
	return allTokens, nil
}

// listRevocationSourceContainsToken reports whether tokens already contains an equal token.
func listRevocationSourceContainsToken[R revocation.Revocation](tokens []RevocationToken[R], token RevocationToken[R]) bool {
	for _, existing := range tokens {
		if existing.Equals(token) {
			return true
		}
	}
	return false
}

// AllRevocationBinaries gets all revocation token binaries from all incorporated sources. Port
// of getAllRevocationBinaries().
//
// Java collects into a HashSet<EncapsulatedRevocationTokenIdentifier<R>>, deduplicating with
// IdentifierBase#equals() (digest + Java class); the Go port dedups on AsXmlID(), which encodes
// exactly that (prefix + hex digest, with the class name fixed per binary implementation).
func (l *ListRevocationSource[R]) AllRevocationBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	var allBinaries []EncapsulatedRevocationTokenIdentifier[R]
	seen := make(map[string]struct{})
	for _, revocationSource := range l.sources {
		for _, binary := range revocationSource.AllRevocationBinaries() {
			id := binary.AsXmlID()
			if _, found := seen[id]; !found {
				seen[id] = struct{}{}
				allBinaries = append(allBinaries, binary)
			}
		}
	}
	return allBinaries
}

// FindBinaryForReference gets the incorporated EncapsulatedRevocationTokenIdentifier
// corresponding to the provided reference. Port of findBinaryForReference(RevocationRef).
func (l *ListRevocationSource[R]) FindBinaryForReference(reference RevocationRef[R]) EncapsulatedRevocationTokenIdentifier[R] {
	for _, revocationSource := range l.sources {
		if tokenIdentifier := revocationSource.FindBinaryForReference(reference); tokenIdentifier != nil {
			return tokenIdentifier
		}
	}
	return nil
}

// IsOrphan checks if the source does not contain revocation identifiers matching the reference.
// Port of isOrphan(RevocationRef).
func (l *ListRevocationSource[R]) IsOrphan(reference RevocationRef[R]) bool {
	for _, revocationSource := range l.sources {
		if !revocationSource.IsOrphan(reference) {
			return false
		}
	}
	return true
}

// NumberOfSources returns the number of set RevocationSources. Port of getNumberOfSources().
func (l *ListRevocationSource[R]) NumberOfSources() int {
	return len(l.sources)
}

// compile-time assertion: a ListRevocationSource is a MultipleRevocationSource.
var _ MultipleRevocationSource[revocation.CRL] = (*ListRevocationSource[revocation.CRL])(nil)
