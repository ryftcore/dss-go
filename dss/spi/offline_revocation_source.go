// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/OfflineRevocationSource.java (DSS 6.5.RC1).
//
// OfflineCRLSourceBase and OfflineOCSPSourceBase (chunk CRLOCSP, a sibling of this phase 2a
// chunk) already embed OfflineRevocationSourceBase[R] (built with
// NewOfflineRevocationSourceBase[R](tokenRefMatcher)) and call AllRevocationBinaries() and
// AddRevocationWithBinary(token, binary); this file defines OfflineRevocationSourceBase to
// match that shape. Their doc comments describe a concrete source registering itself with
// InitOfflineRevocationSource "directly when it overrides RevocationTokens, with the base value
// otherwise" - neither constructor actually calls it, which is consistent with OfflineCRLSource/
// OfflineOCSPSource still being abstract-ish bases themselves in this phase: the eventual leaf
// source (a future phase's SignatureCRLSource et al., embedding OfflineCRLSourceBase) is the one
// expected to call InitOfflineRevocationSource(leafSource) in its own constructor, exactly the
// way InitToken/InitRevocationToken are always called by the outermost concrete type.
package spi

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/utils"
)

// OfflineRevocationSourceOverrides captures what OfflineRevocationSourceBase needs to reach
// through virtual dispatch: the plural RevocationTokens a concrete source (OfflineCRLSourceBase,
// OfflineOCSPSourceBase, ...) must provide, so that RevocationToken (singular) - which Java's
// getRevocationToken(CertificateToken, CertificateToken) implements concretely in terms of it -
// reaches the right implementation.
type OfflineRevocationSourceOverrides[R revocation.Revocation] interface {
	RevocationTokens(certificateToken, issuerCertificateToken *model.CertificateToken) ([]RevocationToken[R], error)
}

// offlineRevocationSourceBinaryEntry pairs a revocation data binary with the origins it has
// been added with.
type offlineRevocationSourceBinaryEntry[R revocation.Revocation] struct {
	binary  EncapsulatedRevocationTokenIdentifier[R]
	origins []enumerations.RevocationOrigin
}

// offlineRevocationSourceTokenEntry pairs a computed RevocationToken with the origins it has
// been added with.
type offlineRevocationSourceTokenEntry[R revocation.Revocation] struct {
	token   RevocationToken[R]
	origins []enumerations.RevocationOrigin
}

// offlineRevocationSourceRefEntry pairs a revocation reference with the origins it has been
// added with.
type offlineRevocationSourceRefEntry[R revocation.Revocation] struct {
	reference RevocationRef[R]
	origins   []enumerations.RevocationRefOrigin
}

// OfflineRevocationSourceBase carries the state and the concrete behaviour of the Java abstract
// class OfflineRevocationSource<R>. A concrete source (OfflineCRLSourceBase, OfflineOCSPSourceBase)
// embeds it and registers itself with InitOfflineRevocationSource.
//
// Java's three java.util.HashMaps (EncapsulatedRevocationTokenIdentifier<R>/RevocationToken<R>/
// RevocationRef<R> -> Set<Origin>) each key on a type with no comparable Go representation
// usable directly as a map key: RevocationToken<R>'s equals() also compares the related
// certificate (see RevocationTokenBase.Equals), and EncapsulatedRevocationTokenIdentifier<R>/
// RevocationRef<R> compare by digest, not by pointer identity. They are ported as
// insertion-ordered slices of pairs, searched with the appropriate equality, mirroring the
// convention already established for TokenCertificateSource's LinkedHashMap fields. DEVIATION,
// flagged for integrator reconciliation: RevocationRef<R>'s fuller equals() (CRLRef also
// compares crlIssuer/crlIssueTime/crlNumber; OCSPRef also compares producedAt/responderId) is
// defined on their own concrete *CRLRef/*OCSPRef receivers, not reachable from this R-agnostic
// code, so referenceOrigins dedups on the reference's DSS identifier (digest) alone; two
// references sharing a digest but disagreeing on those auxiliary fields - not expected in
// practice for the same revocation data - collapse into one entry here.
type OfflineRevocationSourceBase[R revocation.Revocation] struct {
	// overrides points back at the concrete source; see InitOfflineRevocationSource.
	overrides OfflineRevocationSourceOverrides[R]

	// binaryOrigins pairs revocation token identifiers with their corresponding origins.
	binaryOrigins []offlineRevocationSourceBinaryEntry[R]
	// tokenOrigins pairs computed RevocationTokens with their origins.
	tokenOrigins []offlineRevocationSourceTokenEntry[R]
	// referenceOrigins pairs revocation references with their origins.
	referenceOrigins []offlineRevocationSourceRefEntry[R]

	// tokenRefMatcher is the used RevocationTokenRefMatcher.
	tokenRefMatcher RevocationTokenRefMatcher[R]
}

// NewOfflineRevocationSourceBase builds the base state of an offline revocation source, whose
// references are matched with tokenRefMatcher.
//
// Panics with the Java message when tokenRefMatcher is missing (Objects.requireNonNull). Port
// of the protected OfflineRevocationSource(RevocationTokenRefMatcher<R>) constructor.
func NewOfflineRevocationSourceBase[R revocation.Revocation](tokenRefMatcher RevocationTokenRefMatcher[R]) OfflineRevocationSourceBase[R] {
	if tokenRefMatcher == nil {
		panic("tokenRefMatcher cannot be null")
	}
	return OfflineRevocationSourceBase[R]{tokenRefMatcher: tokenRefMatcher}
}

// InitOfflineRevocationSource registers the concrete source with its base so that the base can
// dispatch RevocationToken (singular) to the operations Java would reach through virtual
// dispatch on getRevocationTokens (plural). It must be called by the outermost concrete source's
// constructor before RevocationToken (singular) is used.
func (s *OfflineRevocationSourceBase[R]) InitOfflineRevocationSource(overrides OfflineRevocationSourceOverrides[R]) {
	s.overrides = overrides
}

// offlineRevocationSourceBaseOverrides returns the registered overrides, panicking when the
// concrete source forgot to call InitOfflineRevocationSource.
func (s *OfflineRevocationSourceBase[R]) offlineRevocationSourceBaseOverrides() OfflineRevocationSourceOverrides[R] {
	if s.overrides == nil {
		panic("OfflineRevocationSource was not initialised: the concrete source must call InitOfflineRevocationSource in its constructor")
	}
	return s.overrides
}

// AddBinary adds a token binary with its origin. Port of addBinary(EncapsulatedRevocationTokenIdentifier, RevocationOrigin).
//
// Panics with the Java messages when binary or origin is missing (Objects.requireNonNull).
func (s *OfflineRevocationSourceBase[R]) AddBinary(binary EncapsulatedRevocationTokenIdentifier[R], origin enumerations.RevocationOrigin) {
	if binary == nil {
		panic("The binary is null")
	}
	if origin == "" {
		panic("The origin is null")
	}
	for i := range s.binaryOrigins {
		entry := &s.binaryOrigins[i]
		if entry.binary.AsXmlID() == binary.AsXmlID() {
			entry.origins = offlineRevocationSourceAddOrigin(entry.origins, origin)
			return
		}
	}
	s.binaryOrigins = append(s.binaryOrigins, offlineRevocationSourceBinaryEntry[R]{
		binary: binary, origins: []enumerations.RevocationOrigin{origin},
	})
}

// AddRevocation adds a revocation token with its origin. Port of addRevocation(RevocationToken, RevocationOrigin).
//
// Panics with the Java messages when token or origin is missing (Objects.requireNonNull).
func (s *OfflineRevocationSourceBase[R]) AddRevocation(token RevocationToken[R], origin enumerations.RevocationOrigin) {
	if token == nil {
		panic("The token is null")
	}
	if origin == "" {
		panic("The origin is null")
	}
	for i := range s.tokenOrigins {
		entry := &s.tokenOrigins[i]
		if entry.token.Equals(token) {
			entry.origins = offlineRevocationSourceAddOrigin(entry.origins, origin)
			return
		}
	}
	s.tokenOrigins = append(s.tokenOrigins, offlineRevocationSourceTokenEntry[R]{
		token: token, origins: []enumerations.RevocationOrigin{origin},
	})
}

// AddRevocationWithBinary adds a RevocationToken built from binary, under every origin binary
// was itself added with. Port of addRevocation(RevocationToken, EncapsulatedRevocationTokenIdentifier).
//
// Panics with the Java messages: token nil (Objects.requireNonNull("The token is null")), binary
// nil (Objects.requireNonNull("The origin is null") - copied verbatim from upstream, which
// reuses that message for this check too), and binary not found among the ones already added
// with AddBinary (IllegalStateException("Unable to find the binary '%s'")).
func (s *OfflineRevocationSourceBase[R]) AddRevocationWithBinary(token RevocationToken[R], binary EncapsulatedRevocationTokenIdentifier[R]) {
	if token == nil {
		panic("The token is null")
	}
	if binary == nil {
		panic("The origin is null")
	}
	origins := s.originsForBinary(binary)
	if origins == nil {
		panic(fmt.Sprintf("Unable to find the binary '%s'", binary.AsXmlID()))
	}
	for _, origin := range origins {
		s.AddRevocation(token, origin)
	}
}

// originsForBinary returns the origins binary was added with, nil when it was never added.
func (s *OfflineRevocationSourceBase[R]) originsForBinary(binary EncapsulatedRevocationTokenIdentifier[R]) []enumerations.RevocationOrigin {
	for _, entry := range s.binaryOrigins {
		if entry.binary.AsXmlID() == binary.AsXmlID() {
			return entry.origins
		}
	}
	return nil
}

// AddRevocationReference adds a revocation reference with its origin. Port of
// addRevocationReference(RevocationRef, RevocationRefOrigin).
//
// Panics with the Java messages when reference or origin is missing (Objects.requireNonNull).
func (s *OfflineRevocationSourceBase[R]) AddRevocationReference(reference RevocationRef[R], origin enumerations.RevocationRefOrigin) {
	if reference == nil {
		panic("The reference is null")
	}
	if origin == "" {
		panic("The origin is null")
	}
	for i := range s.referenceOrigins {
		entry := &s.referenceOrigins[i]
		if entry.reference.DSSIDAsString() == reference.DSSIDAsString() {
			entry.origins = offlineRevocationSourceAddRefOrigin(entry.origins, origin)
			return
		}
	}
	s.referenceOrigins = append(s.referenceOrigins, offlineRevocationSourceRefEntry[R]{
		reference: reference, origins: []enumerations.RevocationRefOrigin{origin},
	})
}

// AllRevocationBinaries retrieves all found revocation binaries. Port of getAllRevocationBinaries().
func (s *OfflineRevocationSourceBase[R]) AllRevocationBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	result := make([]EncapsulatedRevocationTokenIdentifier[R], 0, len(s.binaryOrigins))
	for _, entry := range s.binaryOrigins {
		result = append(result, entry.binary)
	}
	return result
}

// AllRevocationTokens retrieves a slice of all found RevocationTokens. Port of
// getAllRevocationTokens() (a Set<RevocationToken<R>> in Java).
func (s *OfflineRevocationSourceBase[R]) AllRevocationTokens() []RevocationToken[R] {
	result := make([]RevocationToken[R], 0, len(s.tokenOrigins))
	for _, entry := range s.tokenOrigins {
		result = append(result, entry.token)
	}
	return result
}

// UniqueRevocationTokensWithOrigins returns the unique RevocationTokens - based on the DSS
// identifier alone, since a single binary can cover several certificates - together with their
// origins. Port of getUniqueRevocationTokensWithOrigins().
func (s *OfflineRevocationSourceBase[R]) UniqueRevocationTokensWithOrigins() []RevocationTokenOriginsEntry[R] {
	var result []RevocationTokenOriginsEntry[R]
	var knownIDs []string
	for _, entry := range s.tokenOrigins {
		currentID := entry.token.DSSID().AsXmlID()
		if !offlineRevocationSourceContainsString(knownIDs, currentID) {
			result = append(result, RevocationTokenOriginsEntry[R]{Token: entry.token, Origins: entry.origins})
			knownIDs = append(knownIDs, currentID)
		}
	}
	return result
}

// AllRevocationReferences retrieves all found RevocationRefs. Port of getAllRevocationReferences().
func (s *OfflineRevocationSourceBase[R]) AllRevocationReferences() []RevocationRef[R] {
	result := make([]RevocationRef[R], 0, len(s.referenceOrigins))
	for _, entry := range s.referenceOrigins {
		result = append(result, entry.reference)
	}
	return result
}

// RevocationToken returns the latest issued revocation token among the ones found for the given
// certificateToken, nil if none is found. Port of the getRevocationToken(CertificateToken,
// CertificateToken) override.
//
// Java's getRevocationTokens has no throws clause; the error the Go MultipleRevocationSource
// override can return is therefore repanicked here, since RevocationSource has no error channel
// either. Java's `latestRevocationToken.getThisUpdate().before(...)` can NullPointerException
// when the first-seen token's thisUpdate is null and a later one's isn't; time.Time's zero value
// compares safely instead, so the Go port simply never hits that edge case.
func (s *OfflineRevocationSourceBase[R]) RevocationToken(certificateToken, issuerCertificateToken *model.CertificateToken) RevocationToken[R] {
	revocationTokens, err := s.offlineRevocationSourceBaseOverrides().RevocationTokens(certificateToken, issuerCertificateToken)
	if err != nil {
		panic(err)
	}

	var latestRevocationToken RevocationToken[R]
	if utils.IsCollectionNotEmpty(revocationTokens) {
		for _, revocationToken := range revocationTokens {
			if latestRevocationToken == nil ||
				(!revocationToken.ThisUpdate().IsZero() && latestRevocationToken.ThisUpdate().Before(revocationToken.ThisUpdate())) {
				latestRevocationToken = revocationToken
			}
		}
	}
	return latestRevocationToken
}

// CMSSignedDataRevocationBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the CMS SignedData. NOTE: applicable only for CAdES revocation sources. Port of
// getCMSSignedDataRevocationBinaries().
func (s *OfflineRevocationSourceBase[R]) CMSSignedDataRevocationBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_CMS_SIGNED_DATA)
}

// CMSSignedDataRevocationTokens retrieves the list of all RevocationTokens present in the CMS
// SignedData. NOTE: applicable only for CAdES revocation sources. Port of
// getCMSSignedDataRevocationTokens().
func (s *OfflineRevocationSourceBase[R]) CMSSignedDataRevocationTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_CMS_SIGNED_DATA)
}

// RevocationValuesBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the 'RevocationValues' element. Port of getRevocationValuesBinaries().
func (s *OfflineRevocationSourceBase[R]) RevocationValuesBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_REVOCATION_VALUES)
}

// RevocationValuesTokens retrieves the list of all RevocationTokens present in the
// 'RevocationValues' element. Port of getRevocationValuesTokens().
func (s *OfflineRevocationSourceBase[R]) RevocationValuesTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_REVOCATION_VALUES)
}

// AttributeRevocationValuesBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the 'AttributeRevocationValues' element. Port of getAttributeRevocationValuesBinaries().
func (s *OfflineRevocationSourceBase[R]) AttributeRevocationValuesBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES)
}

// AttributeRevocationValuesTokens retrieves the list of all RevocationTokens present in the
// 'AttributeRevocationValues' element. Port of getAttributeRevocationValuesTokens().
func (s *OfflineRevocationSourceBase[R]) AttributeRevocationValuesTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_ATTRIBUTE_REVOCATION_VALUES)
}

// TimestampValidationDataBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the 'TimestampValidationData' element. Port of getTimestampValidationDataBinaries().
func (s *OfflineRevocationSourceBase[R]) TimestampValidationDataBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA)
}

// TimestampValidationDataTokens retrieves the list of all RevocationTokens present in the
// 'TimestampValidationData' element. Port of getTimestampValidationDataTokens().
func (s *OfflineRevocationSourceBase[R]) TimestampValidationDataTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_TIMESTAMP_VALIDATION_DATA)
}

// AnyValidationDataBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the 'AnyValidationData' element. Port of getAnyValidationDataBinaries().
func (s *OfflineRevocationSourceBase[R]) AnyValidationDataBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_ANY_VALIDATION_DATA)
}

// AnyValidationDataTokens retrieves the list of all RevocationTokens present in the
// 'AnyValidationData' element. Port of getAnyValidationDataTokens().
func (s *OfflineRevocationSourceBase[R]) AnyValidationDataTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_ANY_VALIDATION_DATA)
}

// DSSDictionaryBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers present
// in 'DSS' dictionary. NOTE: applicable only for PAdES revocation source. Port of
// getDSSDictionaryBinaries().
func (s *OfflineRevocationSourceBase[R]) DSSDictionaryBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_DSS_DICTIONARY)
}

// DSSDictionaryTokens retrieves the list of all RevocationTokens present in 'DSS' dictionary.
// NOTE: applicable only for PAdES revocation source. Port of getDSSDictionaryTokens().
func (s *OfflineRevocationSourceBase[R]) DSSDictionaryTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_DSS_DICTIONARY)
}

// VRIDictionaryBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers present
// in 'VRI' dictionary. NOTE: applicable only for PAdES revocation source. Port of
// getVRIDictionaryBinaries().
func (s *OfflineRevocationSourceBase[R]) VRIDictionaryBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_VRI_DICTIONARY)
}

// VRIDictionaryTokens retrieves the list of all RevocationTokens present in 'VRI' dictionary.
// NOTE: applicable only for PAdES revocation source. Port of getVRIDictionaryTokens().
func (s *OfflineRevocationSourceBase[R]) VRIDictionaryTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_VRI_DICTIONARY)
}

// ADBERevocationValuesBinaries retrieves the list of all EncapsulatedRevocationTokenIdentifiers
// present in the ADBE signed attribute. Port of getADBERevocationValuesBinaries().
func (s *OfflineRevocationSourceBase[R]) ADBERevocationValuesBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	return s.binariesByOrigin(enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL)
}

// ADBERevocationValuesTokens retrieves the list of all RevocationTokens present in the ADBE
// signed attribute. NOTE: applicable only for PAdES revocation source. Port of
// getADBERevocationValuesTokens().
func (s *OfflineRevocationSourceBase[R]) ADBERevocationValuesTokens() []RevocationToken[R] {
	return s.tokensByOrigin(enumerations.RevocationOrigin_ADBE_REVOCATION_INFO_ARCHIVAL)
}

// CompleteRevocationRefs retrieves the list of all RevocationRefs present in the signature
// 'complete-revocation-references' attribute (used in CAdES and XAdES). Port of
// getCompleteRevocationRefs().
func (s *OfflineRevocationSourceBase[R]) CompleteRevocationRefs() []RevocationRef[R] {
	return s.referencesByOrigin(enumerations.RevocationRefOrigin_COMPLETE_REVOCATION_REFS)
}

// AttributeRevocationRefs retrieves the list of all RevocationRefs present in the signature
// 'attribute-revocation-references' attribute (used in CAdES and XAdES). Port of
// getAttributeRevocationRefs().
func (s *OfflineRevocationSourceBase[R]) AttributeRevocationRefs() []RevocationRef[R] {
	return s.referencesByOrigin(enumerations.RevocationRefOrigin_ATTRIBUTE_REVOCATION_REFS)
}

// FindRefsAndOriginsForRevocationToken retrieves the RevocationRefs with their origins found for
// the given RevocationToken. Port of findRefsAndOriginsForRevocationToken(RevocationToken).
//
// Java's tokenRefMatcher.match has no throws clause; the Go RevocationTokenRefMatcher.Match can
// return an error (an unset digest on the reference), which - having no channel to return it
// through here either - is repanicked.
func (s *OfflineRevocationSourceBase[R]) FindRefsAndOriginsForRevocationToken(revocationToken RevocationToken[R]) []RevocationRefOriginsEntry[R] {
	var result []RevocationRefOriginsEntry[R]
	for _, entry := range s.referenceOrigins {
		matched, err := s.tokenRefMatcher.Match(revocationToken, entry.reference)
		if err != nil {
			panic(err)
		}
		if matched {
			result = append(result, RevocationRefOriginsEntry[R]{Reference: entry.reference, Origins: entry.origins})
		}
	}
	return result
}

// FindRefsAndOriginsForBinary retrieves the orphan RevocationRefs with their RevocationRefOrigins
// for a given EncapsulatedRevocationTokenIdentifier. Port of
// findRefsAndOriginsForBinary(EncapsulatedRevocationTokenIdentifier).
func (s *OfflineRevocationSourceBase[R]) FindRefsAndOriginsForBinary(identifier EncapsulatedRevocationTokenIdentifier[R]) []RevocationRefOriginsEntry[R] {
	var result []RevocationRefOriginsEntry[R]
	for _, entry := range s.referenceOrigins {
		matched, err := s.tokenRefMatcher.MatchBinary(identifier, entry.reference)
		if err != nil {
			panic(err)
		}
		if matched {
			result = append(result, RevocationRefOriginsEntry[R]{Reference: entry.reference, Origins: entry.origins})
		}
	}
	return result
}

// FindBinaryForReference returns the linked EncapsulatedRevocationTokenIdentifier for a given
// RevocationRef, nil when none matches. Port of findBinaryForReference(RevocationRef).
func (s *OfflineRevocationSourceBase[R]) FindBinaryForReference(ref RevocationRef[R]) EncapsulatedRevocationTokenIdentifier[R] {
	for _, entry := range s.binaryOrigins {
		matched, err := s.tokenRefMatcher.MatchBinary(entry.binary, ref)
		if err != nil {
			panic(err)
		}
		if matched {
			return entry.binary
		}
	}
	return nil
}

// OrphanRevocationReferencesWithOrigins retrieves the orphan RevocationRefs with their origins.
// Port of getOrphanRevocationReferencesWithOrigins().
func (s *OfflineRevocationSourceBase[R]) OrphanRevocationReferencesWithOrigins() []RevocationRefOriginsEntry[R] {
	var result []RevocationRefOriginsEntry[R]
	for _, entry := range s.referenceOrigins {
		if s.IsOrphan(entry.reference) {
			result = append(result, RevocationRefOriginsEntry[R]{Reference: entry.reference, Origins: entry.origins})
		}
	}
	return result
}

// IsOrphan verifies whether the given RevocationRef is an orphan (not linked to a complete
// RevocationToken). Port of isOrphan(RevocationRef).
func (s *OfflineRevocationSourceBase[R]) IsOrphan(reference RevocationRef[R]) bool {
	for _, entry := range s.tokenOrigins {
		matched, err := s.tokenRefMatcher.Match(entry.token, reference)
		if err != nil {
			panic(err)
		}
		if matched {
			return false
		}
	}
	return true
}

// AllReferencedRevocationBinaries retrieves the set of EncapsulatedRevocationTokenIdentifiers
// which have a reference. Port of getAllReferencedRevocationBinaries().
func (s *OfflineRevocationSourceBase[R]) AllReferencedRevocationBinaries() []EncapsulatedRevocationTokenIdentifier[R] {
	var result []EncapsulatedRevocationTokenIdentifier[R]
	seen := make(map[string]struct{})
	for _, refEntry := range s.referenceOrigins {
		for _, binEntry := range s.binaryOrigins {
			matched, err := s.tokenRefMatcher.MatchBinary(binEntry.binary, refEntry.reference)
			if err != nil {
				panic(err)
			}
			if matched {
				id := binEntry.binary.AsXmlID()
				if _, found := seen[id]; !found {
					seen[id] = struct{}{}
					result = append(result, binEntry.binary)
				}
			}
		}
	}
	return result
}

// IsEmpty checks if the revocation source is empty. Port of isEmpty().
func (s *OfflineRevocationSourceBase[R]) IsEmpty() bool {
	return len(s.binaryOrigins) == 0 && len(s.tokenOrigins) == 0 && len(s.referenceOrigins) == 0
}

// binariesByOrigin retrieves a slice of EncapsulatedRevocationTokenIdentifiers for a given
// RevocationOrigin. Port of the private getBinariesByOrigin(RevocationOrigin).
func (s *OfflineRevocationSourceBase[R]) binariesByOrigin(origin enumerations.RevocationOrigin) []EncapsulatedRevocationTokenIdentifier[R] {
	var result []EncapsulatedRevocationTokenIdentifier[R]
	for _, entry := range s.binaryOrigins {
		if offlineRevocationSourceContainsOrigin(entry.origins, origin) {
			result = append(result, entry.binary)
		}
	}
	return result
}

// tokensByOrigin retrieves a slice of RevocationTokens for a given RevocationOrigin. Port of the
// private getTokensByOrigin(RevocationOrigin).
func (s *OfflineRevocationSourceBase[R]) tokensByOrigin(origin enumerations.RevocationOrigin) []RevocationToken[R] {
	var result []RevocationToken[R]
	for _, entry := range s.tokenOrigins {
		if offlineRevocationSourceContainsOrigin(entry.origins, origin) {
			result = append(result, entry.token)
		}
	}
	return result
}

// referencesByOrigin retrieves a slice of RevocationRefs for a given RevocationRefOrigin. Port
// of the private getReferencesByOrigin(RevocationRefOrigin).
func (s *OfflineRevocationSourceBase[R]) referencesByOrigin(origin enumerations.RevocationRefOrigin) []RevocationRef[R] {
	var result []RevocationRef[R]
	for _, entry := range s.referenceOrigins {
		if offlineRevocationSourceContainsRefOrigin(entry.origins, origin) {
			result = append(result, entry.reference)
		}
	}
	return result
}

// RevocationTokenOriginsEntry pairs a RevocationToken with the origins it has been found with,
// standing in for the Java Map<RevocationToken<R>, Set<RevocationOrigin>> return values.
type RevocationTokenOriginsEntry[R revocation.Revocation] struct {
	Token   RevocationToken[R]
	Origins []enumerations.RevocationOrigin
}

// RevocationRefOriginsEntry pairs a RevocationRef with the origins it has been found with,
// standing in for the Java Map<RevocationRef<R>, Set<RevocationRefOrigin>> return values.
type RevocationRefOriginsEntry[R revocation.Revocation] struct {
	Reference RevocationRef[R]
	Origins   []enumerations.RevocationRefOrigin
}

// offlineRevocationSourceAddOrigin appends origin to origins unless it is already present,
// mirroring the Set<RevocationOrigin> semantics of the corresponding Java field.
func offlineRevocationSourceAddOrigin(origins []enumerations.RevocationOrigin, origin enumerations.RevocationOrigin) []enumerations.RevocationOrigin {
	if offlineRevocationSourceContainsOrigin(origins, origin) {
		return origins
	}
	return append(origins, origin)
}

// offlineRevocationSourceContainsOrigin reports whether origins contains origin.
func offlineRevocationSourceContainsOrigin(origins []enumerations.RevocationOrigin, origin enumerations.RevocationOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}

// offlineRevocationSourceAddRefOrigin appends origin to origins unless it is already present,
// mirroring the Set<RevocationRefOrigin> semantics of the corresponding Java field.
func offlineRevocationSourceAddRefOrigin(origins []enumerations.RevocationRefOrigin, origin enumerations.RevocationRefOrigin) []enumerations.RevocationRefOrigin {
	if offlineRevocationSourceContainsRefOrigin(origins, origin) {
		return origins
	}
	return append(origins, origin)
}

// offlineRevocationSourceContainsRefOrigin reports whether origins contains origin.
func offlineRevocationSourceContainsRefOrigin(origins []enumerations.RevocationRefOrigin, origin enumerations.RevocationRefOrigin) bool {
	for _, o := range origins {
		if o == origin {
			return true
		}
	}
	return false
}

// offlineRevocationSourceContainsString reports whether values contains value.
func offlineRevocationSourceContainsString(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

// compile-time assertion: an OfflineRevocationSourceBase satisfies RevocationSource on its own
// (RevocationToken is concrete); RevocationTokens (plural) stays abstract, dispatched through
// OfflineRevocationSourceOverrides, so MultipleRevocationSource is only satisfied once a
// concrete source (e.g. OfflineCRLSourceBase) provides its own RevocationTokens override.
var _ RevocationSource[revocation.CRL] = (*OfflineRevocationSourceBase[revocation.CRL])(nil)
