// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/timestamp/AbstractTimestampSource.java (DSS 6.5.RC1).
//
// # Package layout
//
// eu.europa.esig.dss.spi.validation.timestamp is a 1:1 sibling package of the flattened
// dss/spi/validation package: AdvancedSignature, SignatureAttribute,
// SignatureProperties, TimestampToken, TimestampedReference, TimestampSource and EvidenceRecord
// all live there and are referenced here as validation.X.
//
// # Methods became package-level functions
//
// Java's AbstractTimestampSource declares no fields at all: every one of its protected methods
// is pure with respect to `this`, existing on the class only so CAdES/XAdES/JAdES timestamp
// sources (embedding it, in later phases, from other Go packages) inherit it. Go has no
// implementation inheritance; the idiomatic port of a stateless "protected helper for
// subclasses" is a set of exported package-level functions, not methods on an empty receiver -
// so that is what this file provides. AbstractTimestampSource itself is kept as an empty marker
// struct purely so DetachedTimestampSource and SignatureTimestampSource can embed it, mirroring
// Java's "extends AbstractTimestampSource" for structural/diff fidelity; no method is promoted
// through the embedding since none of the functions below are methods.
//
// # EvidenceRecord's revocation-source types
//
// This file, detached_timestamp_source.go and signature_timestamp_source.go call a subset of
// validation.EvidenceRecord's methods: Id(), TimestampedReferences(),
// SetTimestampedReferences(...), Timestamps(), CertificateSource(), CRLSource(), OCSPSource(),
// DetachedEvidenceRecords(), ManifestFile(). Two judgment calls about that subset, flagged for
// integrator awareness:
//
//   - CRLSource()/OCSPSource() are typed as evidenceRecordRevocationSource[R] (defined below in
//     this file), not the spi.OfflineRevocationSource[R] interface Java's getCRLSource()/
//     getOCSPSource() return: this file both merges the result into this package's own
//     *spi.ListRevocationSource[R] (needing the full spi.OfflineRevocationSource[R] contract)
//     and reads AllRevocationReferences() off it directly (which spi.OfflineRevocationSource[R]
//     does not expose - see revocationBinaryLookup's doc comment below), so it needs both at
//     once. TimestampToken.CRLSource()/OCSPSource() (dss/spi/validation/timestamp_crl_source.go)
//     hit a related problem and resolved it by returning a concrete type instead of the
//     interface; evidenceRecordRevocationSource generalizes that to an interface combining both
//     needs, since a concrete EvidenceRecord's revocation source type is not known here.
//   - SetTimestampedReferences has no Java counterpart: Java mutates the List<TimestampedReference>
//     getTimestampedReferences() returns in place (List reference semantics). Go slices returned
//     by value do not alias the field they came from, so TimestampToken.SetTimestampedReferences
//     (dss/spi/validation) exists purely to make such a mutation observable to later callers.
package timestamp

import (
	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// must panics with a model.DSSError wrapping err when err is non-nil, otherwise returning value
// unchanged. It is the boundary this package uses wherever a method must satisfy a frozen,
// error-free method signature (SignatureTimestampSource implements validation.TimestampSource,
// whose methods - like every Java method in this manifest - declare no throws clause) but
// internally calls a helper this port made fallible (e.g. CreateReferencesForOCSPBinary, whose
// spi.NewOCSPCertificateSource call can fail); it reproduces Java's unchecked-DSSException
// propagation through a method that does not declare it.
func must[T any](value T, err error) T {
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return value
}

// AbstractTimestampSource contains a set of TimestampTokens found in a
// validation.AdvancedSignature object.
//
// It carries no state (see the package doc comment); DetachedTimestampSource and
// SignatureTimestampSource embed it for structural fidelity with the Java class hierarchy.
//
// java.io.Serializable is dropped (no Go counterpart).
type AbstractTimestampSource struct{}

// revocationBinaryLookup is the minimal contract CreateReferencesForCRLRefs and
// CreateReferencesForOCSPRefs need from their currentXXXSource parameter: looking up the
// encapsulated binary matching a reference. Java's OfflineRevocationSource<R> parameter type
// exposes much more (all of RevocationSource, MultipleRevocationSource, IsEmpty, ...), none of
// which either function uses.
//
// Narrowing to this shape (rather than requiring the full spi.OfflineRevocationSource[R]
// interface) is also what makes both call sites of these two functions compile: the merged
// ListRevocationSource-backed sources (*validation.TimestampCRLSource et al.) satisfy the full
// interface, but EvidenceRecord.CRLSource()/OCSPSource() (see this file's package doc comment)
// is typed as the concrete *spi.OfflineRevocationSourceBase[R],
// which does not - it only ever gets RevocationTokens (plural) once a concrete leaf source
// embeds it and provides an override, which is not the case for a bare base value.
type revocationBinaryLookup[R revocation.Revocation] interface {
	FindBinaryForReference(reference spi.RevocationRef[R]) spi.EncapsulatedRevocationTokenIdentifier[R]
}

// evidenceRecordRevocationSource is what EvidenceRecord.CRLSource()/OCSPSource() returns: the
// full spi.OfflineRevocationSource[R] contract, so it can be folded into this package's merged
// *spi.ListRevocationSource[R] via .Add(), plus AllRevocationReferences() (which
// spi.OfflineRevocationSource[R] itself does not expose - see revocationBinaryLookup above - but
// which every concrete offline revocation source built on spi.OfflineRevocationSourceBase[R],
// e.g. TimestampCRLSource, provides).
type evidenceRecordRevocationSource[R revocation.Revocation] interface {
	spi.OfflineRevocationSource[R]
	AllRevocationReferences() []spi.RevocationRef[R]
}

// xmlIdentifiable is the minimal contract createReferenceForIdentifier/createReferencesForIdentifiers
// and their CRL/OCSP-binary callers need: an object exposing its XML-conformant Id string. Every
// model.Identifier (certificate tokens, signature scopes, ...) and every
// spi.EncapsulatedRevocationTokenIdentifier[R]/*crlparser.CRLBinary/*spi.OCSPResponseBinary
// (CRL/OCSP binaries, whether reached through the merged ListRevocationSource shape or a
// format-specific concrete shape) already satisfies it.
type xmlIdentifiable interface {
	AsXmlID() string
}

// addReference adds referenceToAdd to *referenceList without duplicates.
// Port of addReference(List, TimestampedReference).
func addReference(referenceList *[]*validation.TimestampedReference, referenceToAdd *validation.TimestampedReference) {
	addReferences(referenceList, []*validation.TimestampedReference{referenceToAdd})
}

// addReferenceForIdentifier adds a reference for the given identifier and category to
// *referenceList without duplicates. Port of the addReference(List, Identifier,
// TimestampedObjectType) overload; Go has no overloading, so it is named distinctly.
func addReferenceForIdentifier(referenceList *[]*validation.TimestampedReference, identifier xmlIdentifiable,
	category enumerations.TimestampedObjectType) {
	addReference(referenceList, validation.NewTimestampedReference(identifier.AsXmlID(), category))
}

// addReferences adds referencesToAdd to *referenceList without duplicates.
// Port of addReferences(List, List) which delegates to DSSUtils.enrichCollection(Collection,
// Collection).
//
// DEVIATION: spi.DSSUtilsEnrichCollection (the already-ported, generic port of enrichCollection)
// is deliberately not reused here: it dedups on Go's == over a comparable type parameter, which
// for a pointer type is identity, not TimestampedReference's value-based equals(); see its own
// doc comment. Deduplicating by *validation.TimestampedReference.Equals (category + object id)
// below reproduces Java's List#contains(Object) semantics exactly.
func addReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
	for _, candidate := range referencesToAdd {
		found := false
		for _, existing := range *referenceList {
			if existing.Equals(candidate) {
				found = true
				break
			}
		}
		if !found {
			*referenceList = append(*referenceList, candidate)
		}
	}
}

// ReferencesFromTimestamp incorporates all references from the given timestampToken.
// Port of getReferencesFromTimestamp(TimestampToken, ListCertificateSource, ListRevocationSource,
// ListRevocationSource).
func ReferencesFromTimestamp(timestampToken *validation.TimestampToken, certificateSource *spi.ListCertificateSource,
	crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) (
	[]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}
	addReference(&references, validation.NewTimestampedReference(timestampToken.DSSIDAsString(), enumerations.TimestampedObjectTypeTimestamp))
	addReferences(&references, timestampToken.TimestampedReferences())
	encapsulated, err := EncapsulatedValuesFromTimestamp(timestampToken, certificateSource, crlSource, ocspSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, encapsulated)
	return references, nil
}

// EncapsulatedValuesFromTimestamp gets a list of all validation data embedded to the
// timestampedTimestamp. Port of getEncapsulatedValuesFromTimestamp(TimestampToken,
// ListCertificateSource, ListRevocationSource, ListRevocationSource).
func EncapsulatedValuesFromTimestamp(timestampedTimestamp *validation.TimestampToken, certificateSource *spi.ListCertificateSource,
	crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) (
	[]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}

	timestampCertificateSource := timestampedTimestamp.CertificateSource()
	addReferences(&references, CreateReferencesForCertificates(timestampCertificateSource.Certificates()))
	addReferences(&references, CreateReferencesForCertificateRefs(timestampCertificateSource.AllCertificateRefs(),
		timestampCertificateSource, certificateSource))

	timestampCRLSource := timestampedTimestamp.CRLSource()
	addReferences(&references, CreateReferencesForCRLBinaries(timestampCRLSource.AllRevocationBinaries()))
	addReferences(&references, CreateReferencesForCRLRefs(timestampCRLSource.AllRevocationReferences(),
		timestampCRLSource, crlSource))

	timestampOCSPSource := timestampedTimestamp.OCSPSource()
	ocspBinaryReferences, err := CreateReferencesForOCSPBinaries(timestampOCSPSource.AllRevocationBinaries(), certificateSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, ocspBinaryReferences)
	ocspRefReferences, err := CreateReferencesForOCSPRefs(timestampOCSPSource.AllRevocationReferences(),
		timestampOCSPSource, certificateSource, ocspSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, ocspRefReferences)

	return references, nil
}

// ReferencesFromEvidenceRecord incorporates all references from the given evidenceRecord.
// Port of getReferencesFromEvidenceRecord(EvidenceRecord, ListCertificateSource,
// ListRevocationSource, ListRevocationSource).
func ReferencesFromEvidenceRecord(evidenceRecord validation.EvidenceRecord, certificateSource *spi.ListCertificateSource,
	crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) (
	[]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}
	addReference(&references, validation.NewTimestampedReference(evidenceRecord.Id(), enumerations.TimestampedObjectTypeEvidenceRecord))
	addReferences(&references, evidenceRecord.TimestampedReferences())
	encapsulated, err := EncapsulatedValuesFromEvidenceRecord(evidenceRecord, certificateSource, crlSource, ocspSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, encapsulated)
	return references, nil
}

// EncapsulatedValuesFromEvidenceRecord gets a list of all validation data embedded to the
// evidenceRecord. Port of getEncapsulatedValuesFromEvidenceRecord(EvidenceRecord,
// ListCertificateSource, ListRevocationSource, ListRevocationSource).
func EncapsulatedValuesFromEvidenceRecord(evidenceRecord validation.EvidenceRecord, certificateSource *spi.ListCertificateSource,
	crlSource *spi.ListRevocationSource[revocation.CRL], ocspSource *spi.ListRevocationSource[revocation.OCSP]) (
	[]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}

	for _, timestampToken := range evidenceRecord.Timestamps() {
		fromTimestamp, err := ReferencesFromTimestamp(timestampToken, certificateSource, crlSource, ocspSource)
		if err != nil {
			return nil, err
		}
		addReferences(&references, fromTimestamp)
	}

	erCertificateSource := evidenceRecord.CertificateSource()
	addReferences(&references, CreateReferencesForCertificates(erCertificateSource.Certificates()))
	addReferences(&references, CreateReferencesForCertificateRefs(erCertificateSource.AllCertificateRefs(),
		erCertificateSource, certificateSource))

	erCRLSource := evidenceRecord.CRLSource()
	addReferences(&references, CreateReferencesForCRLBinaries(erCRLSource.AllRevocationBinaries()))
	addReferences(&references, CreateReferencesForCRLRefs(erCRLSource.AllRevocationReferences(),
		erCRLSource, crlSource))

	erOCSPSource := evidenceRecord.OCSPSource()
	ocspBinaryReferences, err := CreateReferencesForOCSPBinaries(erOCSPSource.AllRevocationBinaries(), certificateSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, ocspBinaryReferences)
	ocspRefReferences, err := CreateReferencesForOCSPRefs(erOCSPSource.AllRevocationReferences(),
		erOCSPSource, certificateSource, ocspSource)
	if err != nil {
		return nil, err
	}
	addReferences(&references, ocspRefReferences)

	return references, nil
}

// SignerDataTimestampedReferences creates a list of TimestampedReferences from a given list of
// SignatureScopes. Port of getSignerDataTimestampedReferences(List).
func SignerDataTimestampedReferences(signatureScopes []scope.SignatureScope) []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	for _, signatureScope := range signatureScopes {
		addReference(&references, validation.NewTimestampedReference(signatureScope.DSSIDAsString(), enumerations.TimestampedObjectTypeSignedData))
		if children := signatureScope.Children(); len(children) > 0 {
			addReferences(&references, SignerDataTimestampedReferences(children))
		}
	}
	return references
}

// CreateReferencesForCertificates creates a list of TimestampedReferences for the provided list
// of certificates. Port of createReferencesForCertificates(Collection).
func CreateReferencesForCertificates(certificates []*model.CertificateToken) []*validation.TimestampedReference {
	references := []*validation.TimestampedReference{}
	for _, certificateToken := range certificates {
		addReference(&references, CreateReferenceForCertificate(certificateToken))
	}
	return references
}

// CreateReferenceForCertificate creates a TimestampedReference for the provided CertificateToken.
// Port of createReferenceForCertificate(CertificateToken).
func CreateReferenceForCertificate(certificateToken *model.CertificateToken) *validation.TimestampedReference {
	return CreateReferenceForIdentifier(certificateToken.DSSID(), enumerations.TimestampedObjectTypeCertificate)
}

// CreateReferencesForIdentifiers creates a list of TimestampedReferences from the identifiers of
// a given type. Port of createReferencesForIdentifiers(Collection, TimestampedObjectType);
// T mirrors Java's Collection<? extends Identifier>.
func CreateReferencesForIdentifiers[T xmlIdentifiable](identifiers []T,
	timestampedObjectType enumerations.TimestampedObjectType) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	for _, identifier := range identifiers {
		timestampedReferences = append(timestampedReferences, CreateReferenceForIdentifier(identifier, timestampedObjectType))
	}
	return timestampedReferences
}

// CreateReferenceForIdentifier creates a TimestampedReference for the given identifier.
// Port of createReferenceForIdentifier(Identifier, TimestampedObjectType).
func CreateReferenceForIdentifier(identifier xmlIdentifiable, timestampedObjectType enumerations.TimestampedObjectType) *validation.TimestampedReference {
	return validation.NewTimestampedReference(identifier.AsXmlID(), timestampedObjectType)
}

// CreateReferencesForCRLBinaries creates a list of TimestampedReferences from a collection of
// CRL binaries. Port of createReferencesForCRLBinaries(Collection).
func CreateReferencesForCRLBinaries[T xmlIdentifiable](crlBinaryIdentifiers []T) []*validation.TimestampedReference {
	return CreateReferencesForIdentifiers(crlBinaryIdentifiers, enumerations.TimestampedObjectTypeRevocation)
}

// CreateReferencesForOCSPBinaries creates a list of TimestampedReferences from a collection of
// OCSP response binaries, expanding each into references for its embedded certificates too.
// Port of createReferencesForOCSPBinaries(Collection, ListCertificateSource).
//
// Java filters the collection with `instanceof OCSPResponseBinary`; the type assertion below,
// which simply skips a non-matching element instead of casting unconditionally, is the direct
// port of that defensive instanceof-guarded cast (contrast with
// createReferencesForOCSPBinary's caller in createReferencesForOCSPRefs, which ports an
// unconditional cast and panics instead).
func CreateReferencesForOCSPBinaries[T xmlIdentifiable](ocspBinaryIdentifiers []T,
	certificateSource *spi.ListCertificateSource) ([]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}
	for _, ocspIdentifier := range ocspBinaryIdentifiers {
		ocspResponseBinary, ok := any(ocspIdentifier).(*spi.OCSPResponseBinary)
		if !ok {
			continue
		}
		binaryReferences, err := CreateReferencesForOCSPBinary(ocspResponseBinary, certificateSource)
		if err != nil {
			return nil, err
		}
		addReferences(&references, binaryReferences)
	}
	return references, nil
}

// CreateReferencesForOCSPBinary creates a list of TimestampedReferences for a OCSPResponseBinary.
// Port of createReferencesForOCSPBinary(OCSPResponseBinary, ListCertificateSource).
//
// Java's `new OCSPCertificateSource(...)` declares no checked exception; the Go port
// (spi.NewOCSPCertificateSource) surfaces its certificate-extraction failures as an error
// instead, which this function - and everything that calls it - therefore propagates.
func CreateReferencesForOCSPBinary(ocspResponseBinary *spi.OCSPResponseBinary, certificateSource *spi.ListCertificateSource) (
	[]*validation.TimestampedReference, error) {
	references := []*validation.TimestampedReference{}

	addReference(&references, CreateReferenceForIdentifier(ocspResponseBinary, enumerations.TimestampedObjectTypeRevocation))

	ocspCertificateSource, err := spi.NewOCSPCertificateSource(ocspResponseBinary.BasicOCSPResp())
	if err != nil {
		return nil, err
	}
	addReferences(&references, CreateReferencesForCertificates(ocspCertificateSource.Certificates()))
	addReferences(&references, CreateReferencesForCertificateRefs(ocspCertificateSource.AllCertificateRefs(),
		ocspCertificateSource, certificateSource))

	return references, nil
}

// CreateReferencesForCertificateRefs returns a list of timestamped references from the given
// collection of certificateRefs. Port of createReferencesForCertificateRefs(Collection,
// CertificateSource, ListCertificateSource).
func CreateReferencesForCertificateRefs(certificateRefs []*spi.CertificateRef, currentCertificateSource spi.CertificateSource,
	listCertificateSource *spi.ListCertificateSource) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	for _, certRef := range certificateRefs {
		certificateTokens := currentCertificateSource.FindTokensFromCertRef(certRef)
		if len(certificateTokens) == 0 {
			certificateTokens = listCertificateSource.FindTokensFromCertRef(certRef)
		}
		if len(certificateTokens) > 0 {
			addReferences(&timestampedReferences, CreateReferencesForCertificates(certificateTokenMapValues(certificateTokens)))
		} else {
			addReference(&timestampedReferences, validation.NewTimestampedReference(certRef.DSSIDAsString(), enumerations.TimestampedObjectTypeCertificate))
		}
	}
	return timestampedReferences
}

// certificateTokenMapValues collects the values of a Set<CertificateToken>-equivalent map (this
// port's convention for Java's Set<CertificateToken>, see spi.ListCertificateSource) into a
// slice, for handing to CreateReferencesForCertificates.
func certificateTokenMapValues(tokens map[string]*model.CertificateToken) []*model.CertificateToken {
	result := make([]*model.CertificateToken, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, token)
	}
	return result
}

// CreateReferencesForCRLRefs returns a list of timestamped references from the given collection
// of crlRefs. Port of createReferencesForCRLRefs(Collection, OfflineRevocationSource,
// ListRevocationSource).
//
// T mirrors Java's Collection<? extends RevocationRef<CRL>>: one call site in this file passes
// the merged-source shape spi.RevocationRef[revocation.CRL] (OfflineRevocationSourceBase.
// AllRevocationReferences()'s element type), the other (signature_timestamp_source.go, via the
// abstract GetCRLRefs) the concrete []*spi.CRLRef.
func CreateReferencesForCRLRefs[T spi.RevocationRef[revocation.CRL]](crlRefs []T, currentCRLSource revocationBinaryLookup[revocation.CRL],
	listCRLSource *spi.ListRevocationSource[revocation.CRL]) []*validation.TimestampedReference {
	timestampedReferences := []*validation.TimestampedReference{}
	for _, crlRef := range crlRefs {
		token := currentCRLSource.FindBinaryForReference(crlRef)
		if token == nil {
			token = listCRLSource.FindBinaryForReference(crlRef)
		}
		if token != nil {
			addReference(&timestampedReferences, validation.NewTimestampedReference(token.AsXmlID(), enumerations.TimestampedObjectTypeRevocation))
		} else {
			addReference(&timestampedReferences, validation.NewTimestampedReference(crlRef.DSSIDAsString(), enumerations.TimestampedObjectTypeRevocation))
		}
	}
	return timestampedReferences
}

// CreateReferencesForOCSPRefs returns a list of timestamped references from the given collection
// of ocspRefs. Port of createReferencesForOCSPRefs(Collection, OfflineRevocationSource,
// ListCertificateSource, ListRevocationSource).
//
// T mirrors Java's Collection<? extends RevocationRef<OCSP>>; see CreateReferencesForCRLRefs.
func CreateReferencesForOCSPRefs[T spi.RevocationRef[revocation.OCSP]](ocspRefs []T, currentOCSPSource revocationBinaryLookup[revocation.OCSP],
	listCertificateSource *spi.ListCertificateSource, listOCSPSource *spi.ListRevocationSource[revocation.OCSP]) (
	[]*validation.TimestampedReference, error) {
	timestampedReferences := []*validation.TimestampedReference{}
	for _, ocspRef := range ocspRefs {
		token := currentOCSPSource.FindBinaryForReference(ocspRef)
		if token == nil {
			token = listOCSPSource.FindBinaryForReference(ocspRef)
		}
		if token != nil {
			// Java casts unconditionally: `(OCSPResponseBinary) token`. Every concrete
			// EncapsulatedRevocationTokenIdentifier[OCSP] in this port is a *spi.OCSPResponseBinary,
			// so this assertion always succeeds in practice; a genuine mismatch panics, exactly
			// as Java's ClassCastException would.
			ocspResponseBinary := token.(*spi.OCSPResponseBinary)
			binaryReferences, err := CreateReferencesForOCSPBinary(ocspResponseBinary, listCertificateSource)
			if err != nil {
				return nil, err
			}
			addReferences(&timestampedReferences, binaryReferences)
		} else {
			addReference(&timestampedReferences, validation.NewTimestampedReference(ocspRef.DSSIDAsString(), enumerations.TimestampedObjectTypeRevocation))
		}
	}
	return timestampedReferences, nil
}

// ProcessEvidenceRecordTimestamps enriches embedded time-stamp tokens with evidence record
// references. Port of processEvidenceRecordTimestamps(EvidenceRecord).
func ProcessEvidenceRecordTimestamps(evidenceRecord validation.EvidenceRecord) {
	for _, timestampToken := range evidenceRecord.Timestamps() {
		ensureOnlyDataTimestampReferencesPresent(timestampToken, evidenceRecord.TimestampedReferences())
		timestampAddReferences(timestampToken, evidenceRecord.TimestampedReferences())
	}
}

// ProcessEmbeddedEvidenceRecords enriches embedded evidence records with the covered references.
// Port of processEmbeddedEvidenceRecords(EvidenceRecord).
func ProcessEmbeddedEvidenceRecords(evidenceRecord validation.EvidenceRecord) {
	for _, embeddedEvidenceRecord := range evidenceRecord.DetachedEvidenceRecords() {
		embeddedEvidenceRecord.SetTimestampedReferences(mergeReferences(embeddedEvidenceRecord.TimestampedReferences(), evidenceRecord.TimestampedReferences()))
		ProcessEvidenceRecordTimestamps(embeddedEvidenceRecord)
	}
}

// mergeReferences returns a copy of base with every reference of additional not already present
// (by TimestampedReference.Equals) appended, i.e. the value-returning shape addReferences needs
// when the list being grown is reached through a getter that returns by value (EvidenceRecord.
// TimestampedReferences(), TimestampToken.TimestampedReferences()) rather than through a field
// this package can take the address of.
func mergeReferences(base []*validation.TimestampedReference, additional []*validation.TimestampedReference) []*validation.TimestampedReference {
	merged := append([]*validation.TimestampedReference(nil), base...)
	addReferences(&merged, additional)
	return merged
}

// timestampAddReferences enriches timestampToken's TimestampedReferences with referencesToAdd,
// without duplicates. Port of the `addReferences(timestampToken.getTimestampedReferences(), ...)`
// call sites, which upstream rely on Java's List reference semantics to mutate the
// TimestampToken's own list in place. TimestampToken.TimestampedReferences() returns its slice
// by value, so the mutation is made observable to later TimestampedReferences() calls through
// TimestampToken.SetTimestampedReferences instead.
func timestampAddReferences(timestampToken *validation.TimestampToken, referencesToAdd []*validation.TimestampedReference) {
	timestampToken.SetTimestampedReferences(mergeReferences(timestampToken.TimestampedReferences(), referencesToAdd))
}

// ensureOnlyDataTimestampReferencesPresent is a workaround to ensure time-stamps from an
// evidence record do not refer to signature or time-stamp files in addition to token references.
// Port of the private ensureOnlyDataTimestampReferencesPresent(List, List), adapted to operate
// directly on the TimestampToken (see timestampAddReferences) instead of on a caller-owned list
// reference Go cannot express here.
func ensureOnlyDataTimestampReferencesPresent(timestampToken *validation.TimestampToken, referencesToCheck []*validation.TimestampedReference) {
	current := timestampToken.TimestampedReferences()
	filtered := make([]*validation.TimestampedReference, 0, len(current))
	for _, timestampedReference := range current {
		remove := timestampedReference.Category() == enumerations.TimestampedObjectTypeSignedData &&
			!containsEqualReference(referencesToCheck, timestampedReference)
		if !remove {
			filtered = append(filtered, timestampedReference)
		}
	}
	timestampToken.SetTimestampedReferences(filtered)
}

// containsEqualReference reports whether references contains an entry equal (per
// TimestampedReference.Equals) to candidate. Port of the `referencesToCheck.stream().noneMatch(...)`
// check embedded in ensureOnlyDataTimestampReferencesPresent's removeIf predicate.
func containsEqualReference(references []*validation.TimestampedReference, candidate *validation.TimestampedReference) bool {
	for _, reference := range references {
		if reference.Equals(candidate) {
			return true
		}
	}
	return false
}

// compile-time assertion: CRLBinary/OCSPResponseBinary/model.Identifier all satisfy
// xmlIdentifiable, matching every T this file's generic functions are actually instantiated
// with across the manifest.
var (
	_ xmlIdentifiable = (*crlparser.CRLBinary)(nil)
	_ xmlIdentifiable = (*spi.OCSPResponseBinary)(nil)
	_ xmlIdentifiable = (model.Identifier)(nil)
	_ xmlIdentifiable = (spi.EncapsulatedRevocationTokenIdentifier[revocation.CRL])(nil)
)
