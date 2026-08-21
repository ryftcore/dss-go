// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/identifier/UserFriendlyIdentifierProvider.java (DSS 6.5.RC1).
//
// Creates an identifier for a given token by the template:
//
//	TOKEN-CommonCertName-CreationDate-id(optional)
//
// Examples:
//
//	SIGNATURE-JohnConner-20201015-2045
//	CERTIFICATE-CryptoSign-20151014-1425
//
// Java's single `getIdAsString(IdentifierBasedObject)` dispatches on the object's runtime class
// through a chain of `instanceof` checks (AdvancedSignature, Token, SignatureScope,
// CertificateRef, RevocationRef<?>, EncapsulatedRevocationTokenIdentifier<?>, EvidenceRecord,
// EAA, TLInfo, LoTEInfo - in that order, since e.g. a `TLInfo` also matches `instanceof TLInfo`
// whether its runtime class is TLInfo, LOTLInfo or PivotInfo). Go's type switch reproduces this:
// a `case` naming an interface type matches any concrete value whose method set satisfies it,
// including one reached only through embedding (LOTLInfo/PivotInfo promote TLInfo's methods,
// EAARevocationToken's Token-shaped siblings promote model.TokenBase's, etc.), so the same
// same-order chain below preserves the upstream dispatch precisely without needing a
// TLInfo/LOTLInfo/PivotInfo (or LoTEInfo/LoLoTEInfo) common base type to exist in Go.
package identifier

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/lote"
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/model/tsl"
	"github.com/utain/esig/dss/spi"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// String is used to separate different parts of the identifier.
const stringDelimiter = "_"

// String to be used to replace non-alphanumeric characters in a certificate's common name.
const nameReplacement = "-"

// String to be used to replace invalid XML characters.
const unsupportedCharacter = "?"

// To be used to define an issuer name.
const issuerLabel = "ISSUER-"

// To be used to define a serial number.
const serialLabel = "SERIAL-"

// String used when token's signing certificate is not identified.
const unknownSigner = "UNKNOWN-SIGNER"

// String used when token's signing certificate does not have a human-readable name.
const unnamedSigner = "UNNAMED-SIGNER"

// revocationRefLike is satisfied by *spi.CRLRef and *spi.OCSPRef (both embed
// spi.RevocationRefBase[R] for their own R, promoting these two methods with an R-independent
// signature). It stands in for Java's use of the raw, wildcarded `RevocationRef<?>` type.
type revocationRefLike interface {
	Digest() model.Digest
	DSSIDAsString() string
}

// revocationTokenIdentifierLike is satisfied by *crlparser.CRLBinary and
// *spi.OCSPResponseBinary (both embed model.EncapsulatedRevocationTokenIdentifier[R] for their
// own R). It stands in for Java's use of the raw, wildcarded
// `EncapsulatedRevocationTokenIdentifier<?>` type.
type revocationTokenIdentifierLike interface {
	model.Identifier
	DigestValue(digestAlgorithm enumerations.DigestAlgorithm) ([]byte, error)
}

// tlInfoLike is satisfied by *tsl.TLInfo, *tsl.LOTLInfo and *tsl.PivotInfo, matching Java's
// `object instanceof TLInfo` catching all three (LOTLInfo/PivotInfo extend TLInfo upstream).
type tlInfoLike interface {
	TLParsingCacheInfo() (tsl.TLParsingInfoRecord, bool)
	DSSIDAsString() string
}

// loteInfoLike is satisfied by *lote.LoTEInfo and *lote.LoLoTEInfo, matching Java's
// `object instanceof LoTEInfo` catching both (LoLoTEInfo extends LoTEInfo upstream).
type loteInfoLike interface {
	ParsingCacheInfo() lote.LoTEParsingInfoRecord
	DSSIDAsString() string
}

// UserFriendlyIdentifierProvider creates a human-friendly String identifier for tokens,
// signature scopes, references, evidence records, EAAs and trusted-list summaries.
type UserFriendlyIdentifierProvider struct {
	// uniqueTokenIdsMap represents a cached values map for processed tokens between the
	// original DSS hash-based Id and the one computed by the class (including preservation
	// from duplicates).
	uniqueTokenIdsMap map[string]string

	// generatedTokenIdsMap maps DSS identifiers of processed tokens to the computed ids by
	// the provider. The map can contain duplicates. Used to identify duplicates.
	generatedTokenIdsMap map[string]string

	// signaturePrefix is the prefix to be used for a signature identifier creation.
	signaturePrefix string
	// counterSignaturePrefix is the prefix to be used for a counter-signature identifier
	// creation.
	counterSignaturePrefix string
	// timestampPrefix is the prefix to be used for a timestamp identifier creation.
	timestampPrefix string
	// certificatePrefix is the prefix to be used for a certificate identifier creation.
	certificatePrefix string
	// crlPrefix is the prefix to be used for an CRL identifier creation.
	crlPrefix string
	// ocspPrefix is the prefix to be used for an OCSP identifier creation.
	ocspPrefix string
	// signedDataPrefix is the prefix to be used for an original document identifier creation.
	signedDataPrefix string
	// evidenceRecordPrefix is the prefix to be used for an evidence record identifier
	// creation.
	evidenceRecordPrefix string
	// eaaPrefix is the prefix to be used for an EAA identifier creation.
	eaaPrefix string
	// eaaStatusTokenPrefix is the prefix to be used for an EAA revocation token identifier
	// creation.
	eaaStatusTokenPrefix string
	// lotlPrefix is the prefix to be used for a List of Trusted Lists identifier creation.
	lotlPrefix string
	// tlPrefix is the prefix to be used for a Trusted List identifier creation.
	tlPrefix string
	// pivotPrefix is the prefix to be used for a pivot identifier creation.
	pivotPrefix string
	// lolotePrefix is the prefix to be used for a List of Lists of Trusted Entities identifier
	// creation.
	//
	// NOTE: dead field, faithfully preserved from upstream: Java declares this field, but no
	// setter exists for it and getLoTEPrefix (see loTEPrefix below) never reads it either -
	// see that method's own NOTE.
	lolotePrefix string
	// lotePrefix is the prefix to be used for a List of Trusted Entities identifier creation.
	//
	// NOTE: dead field; see lolotePrefix above and loTEPrefix below.
	lotePrefix string
	// dateFormat is the date format to be used for a token identifier creation, a
	// java.text.SimpleDateFormat pattern.
	dateFormat string
}

// compile-time interface assertion.
var _ model.TokenIdentifierProvider = (*UserFriendlyIdentifierProvider)(nil)

// NewUserFriendlyIdentifierProvider is the default constructor, instantiating empty maps of
// processed tokens and every prefix/date format at its documented default.
func NewUserFriendlyIdentifierProvider() *UserFriendlyIdentifierProvider {
	return &UserFriendlyIdentifierProvider{
		uniqueTokenIdsMap:      make(map[string]string),
		generatedTokenIdsMap:   make(map[string]string),
		signaturePrefix:        "SIGNATURE",
		counterSignaturePrefix: "COUNTER-SIGNATURE",
		timestampPrefix:        "TIMESTAMP",
		certificatePrefix:      "CERTIFICATE",
		crlPrefix:              "CRL",
		ocspPrefix:             "OCSP",
		signedDataPrefix:       "DOCUMENT",
		evidenceRecordPrefix:   "EVIDENCE-RECORD",
		eaaPrefix:              "EAA",
		eaaStatusTokenPrefix:   "EAA-STATUS",
		lotlPrefix:             "LOTL",
		tlPrefix:               "TL",
		pivotPrefix:            "PIVOT",
		lolotePrefix:           "LOLOTE",
		lotePrefix:             "LOTE",
		dateFormat:             "yyyyMMdd-HHmm",
	}
}

// SetSignaturePrefix sets the prefix to be used for signature identifiers. Default =
// "SIGNATURE".
func (p *UserFriendlyIdentifierProvider) SetSignaturePrefix(signaturePrefix string) {
	assertNotBlank(signaturePrefix)
	p.signaturePrefix = signaturePrefix
}

// SetCounterSignaturePrefix sets the prefix to be used for counter-signature identifiers.
// Default = "COUNTER-SIGNATURE".
func (p *UserFriendlyIdentifierProvider) SetCounterSignaturePrefix(counterSignaturePrefix string) {
	assertNotBlank(counterSignaturePrefix)
	p.counterSignaturePrefix = counterSignaturePrefix
}

// SetTimestampPrefix sets the prefix to be used for timestamp identifiers. Default =
// "TIMESTAMP".
func (p *UserFriendlyIdentifierProvider) SetTimestampPrefix(timestampPrefix string) {
	assertNotBlank(timestampPrefix)
	p.timestampPrefix = timestampPrefix
}

// SetCertificatePrefix sets the prefix to be used for certificate identifiers. Default =
// "CERTIFICATE".
func (p *UserFriendlyIdentifierProvider) SetCertificatePrefix(certificatePrefix string) {
	assertNotBlank(certificatePrefix)
	p.certificatePrefix = certificatePrefix
}

// SetCrlPrefix sets the prefix to be used for CRL identifiers. Default = "CRL".
func (p *UserFriendlyIdentifierProvider) SetCrlPrefix(crlPrefix string) {
	assertNotBlank(crlPrefix)
	p.crlPrefix = crlPrefix
}

// SetOcspPrefix sets the prefix to be used for OCSP identifiers. Default = "OCSP".
func (p *UserFriendlyIdentifierProvider) SetOcspPrefix(ocspPrefix string) {
	assertNotBlank(ocspPrefix)
	p.ocspPrefix = ocspPrefix
}

// SetSignedDataPrefix sets the prefix to be used for original document identifiers. Default =
// "DOCUMENT".
func (p *UserFriendlyIdentifierProvider) SetSignedDataPrefix(signedDataPrefix string) {
	assertNotBlank(signedDataPrefix)
	p.signedDataPrefix = signedDataPrefix
}

// SetEvidenceRecordPrefix sets the prefix to be used for evidence record identifiers. Default
// = "EVIDENCE-RECORD".
func (p *UserFriendlyIdentifierProvider) SetEvidenceRecordPrefix(evidenceRecordPrefix string) {
	assertNotBlank(evidenceRecordPrefix)
	p.evidenceRecordPrefix = evidenceRecordPrefix
}

// SetLOTLPrefix sets the prefix to be used for a LOTL identifier. Default = "LOTL".
func (p *UserFriendlyIdentifierProvider) SetLOTLPrefix(lotlPrefix string) {
	assertNotBlank(lotlPrefix)
	p.lotlPrefix = lotlPrefix
}

// SetTLPrefix sets the prefix to be used for TL identifiers. Default = "TL".
func (p *UserFriendlyIdentifierProvider) SetTLPrefix(tlPrefix string) {
	assertNotBlank(tlPrefix)
	p.tlPrefix = tlPrefix
}

// SetPivotPrefix sets the prefix to be used for pivot identifiers. Default = "PIVOT".
func (p *UserFriendlyIdentifierProvider) SetPivotPrefix(pivotPrefix string) {
	assertNotBlank(pivotPrefix)
	p.pivotPrefix = pivotPrefix
}

// SetEAAPrefix sets the prefix to be used for EAA identifiers. Default = "EAA".
//
// NOTE: Java's setEAAPrefix(String) calls `assertNotBlank(pivotPrefix)`, validating the
// current pivotPrefix FIELD rather than the eaaPrefix parameter being set - every other setter
// in this class validates its own parameter, so this looks like a copy-paste bug upstream. It
// is reproduced verbatim: since pivotPrefix defaults to "PIVOT" and can only be replaced by
// SetPivotPrefix (itself correctly blank-guarded), this check can never actually reject a
// blank eaaPrefix argument in practice.
func (p *UserFriendlyIdentifierProvider) SetEAAPrefix(eaaPrefix string) {
	assertNotBlank(p.pivotPrefix)
	p.eaaPrefix = eaaPrefix
}

// SetEAAStatusTokenPrefix sets the prefix to be used for EAA Status Token identifiers. Default
// = "EAA-STATUS".
func (p *UserFriendlyIdentifierProvider) SetEAAStatusTokenPrefix(eaaStatusTokenPrefix string) {
	assertNotBlank(eaaStatusTokenPrefix)
	p.eaaStatusTokenPrefix = eaaStatusTokenPrefix
}

// SetDateFormat sets the dataFormat to be used for identifiers creation. Default =
// "yyyyMMdd-HHmm".
//
// Panics with the Java message when dateFormat is empty (Objects.requireNonNull("The
// dataFormat cannot be null!") - note the "dataFormat"/"dateFormat" mismatch is upstream's own
// typo, reproduced verbatim).
func (p *UserFriendlyIdentifierProvider) SetDateFormat(dateFormat string) {
	if dateFormat == "" {
		panic("The dataFormat cannot be null!")
	}
	p.dateFormat = dateFormat
}

// IDAsString returns a String identifier for the given object. Port of getIdAsString().
//
// Panics when object is nil (Java's Objects.requireNonNull("The object cannot be null!")).
func (p *UserFriendlyIdentifierProvider) IDAsString(object model.IdentifierBasedObject) string {
	if object == nil {
		panic("The object cannot be null!")
	}

	if cachedIdentifier := p.cachedIdentifier(object); cachedIdentifier != "" {
		return cachedIdentifier
	}

	switch v := object.(type) {
	case spivalidation.AdvancedSignature:
		return p.idAsStringForSignature(v)
	case model.Token:
		return p.idAsStringForToken(v)
	case scope.SignatureScope:
		return p.idAsStringForSignatureScope(v)
	case *spi.CertificateRef:
		return p.idAsStringForCertRef(v)
	case revocationRefLike:
		return p.idAsStringForRevRef(v)
	case revocationTokenIdentifierLike:
		return p.idAsStringForRevTokenIdentifier(v)
	case spivalidation.EvidenceRecord:
		return p.idAsStringForEvidenceRecordIdentifier(v)
	case spivalidation.EAA:
		return p.idAsStringForEAAIdentifier(v)
	case tlInfoLike:
		return p.idAsStringForTL(v)
	case loteInfoLike:
		return p.idAsStringForLoTE(v)
	}
	// LOG.warn("The class '{}' is not supported! ...") dropped: not load-bearing (PORTING.md).
	return object.DSSID().AsXmlID()
}

// cachedIdentifier looks the object's original DSS identifier up in uniqueTokenIdsMap,
// returning "" (Java's null) on a cache miss. Port of the private getCachedIdentifier(...).
//
// Panics when the object's DSSID() is nil (Java's thrown IllegalArgumentException("The
// returned Identifier cannot be null for the object of class '%s'!")).
func (p *UserFriendlyIdentifierProvider) cachedIdentifier(object model.IdentifierBasedObject) string {
	identifier := object.DSSID()
	if identifier == nil {
		panic(fmt.Sprintf("The returned Identifier cannot be null for the object of class '%T'!", object))
	}
	originalIdentifier := identifier.AsXmlID()
	return p.uniqueTokenIdsMap[originalIdentifier]
}

// idAsStringForSignature gets a String identifier for a given AdvancedSignature. Port of the
// protected getIdAsStringForSignature(AdvancedSignature).
func (p *UserFriendlyIdentifierProvider) idAsStringForSignature(signature spivalidation.AdvancedSignature) string {
	var subject *model.X500PrincipalHelper
	if signingCertificate := signature.SigningCertificateToken(); signingCertificate != nil {
		subject = signingCertificate.Subject()
	}
	prefix := p.signaturePrefix
	if signature.IsCounterSignature() {
		prefix = p.counterSignaturePrefix
	}
	var signingTime time.Time
	if t := signature.SigningTime(); t != nil {
		signingTime = *t
	}
	return p.createIDString(prefix, subject, signingTime, signature.ID())
}

// idAsStringForToken gets a String identifier for a given Token. Port of the protected
// getIdAsStringForToken(Token).
func (p *UserFriendlyIdentifierProvider) idAsStringForToken(token model.Token) string {
	return p.createIDString(p.tokenPrefix(token), p.tokenSubject(token), token.CreationDate(), token.DSSIDAsString())
}

// idAsStringForSignatureScope gets a String identifier for a given SignatureScope. Port of the
// protected getIdAsStringForSignatureScope(SignatureScope).
func (p *UserFriendlyIdentifierProvider) idAsStringForSignatureScope(signatureScope scope.SignatureScope) string {
	var sb strings.Builder
	sb.WriteString(p.signedDataPrefix)
	sb.WriteString(stringDelimiter)
	documentName := signatureScope.Name(p)
	if utils.IsStringNotBlank(documentName) {
		sb.WriteString(p.userFriendlyString(documentName))
	} else {
		sb.WriteString(string(signatureScope.Type()))
	}
	return p.generateID(&sb, signatureScope.DSSIDAsString())
}

// idAsStringForTL gets a String identifier for a given TLInfo (or a LOTLInfo/PivotInfo, both
// of which are-a TLInfo upstream). Port of the protected getIdAsStringForTL(TLInfo).
func (p *UserFriendlyIdentifierProvider) idAsStringForTL(tlInfo tlInfoLike) string {
	var sb strings.Builder
	sb.WriteString(p.tlInfoPrefix(tlInfo))
	if parsingCacheInfo, ok := tlInfo.TLParsingCacheInfo(); ok {
		if utils.IsStringNotBlank(parsingCacheInfo.Territory()) {
			sb.WriteString(stringDelimiter)
			sb.WriteString(p.userFriendlyString(parsingCacheInfo.Territory()))
		}
		if issueDate := parsingCacheInfo.IssueDate(); !issueDate.IsZero() {
			sb.WriteString(stringDelimiter)
			sb.WriteString(spi.DSSUtilsFormatDateWithCustomFormat(issueDate, p.dateFormat))
		}
	}
	return p.generateID(&sb, tlInfo.DSSIDAsString())
}

// idAsStringForLoTE gets a String identifier for a given LoTEInfo (or a LoLoTEInfo, which
// is-a LoTEInfo upstream). Port of the protected getIdAsStringForLoTE(LoTEInfo).
func (p *UserFriendlyIdentifierProvider) idAsStringForLoTE(listInfo loteInfoLike) string {
	var sb strings.Builder
	sb.WriteString(p.loTEInfoPrefix(listInfo))
	if parsingCacheInfo := listInfo.ParsingCacheInfo(); parsingCacheInfo != nil {
		if utils.IsStringNotBlank(parsingCacheInfo.Territory()) {
			sb.WriteString(stringDelimiter)
			sb.WriteString(p.userFriendlyString(parsingCacheInfo.Territory()))
		}
		if issueDate := parsingCacheInfo.IssueDate(); !issueDate.IsZero() {
			sb.WriteString(stringDelimiter)
			sb.WriteString(spi.DSSUtilsFormatDateWithCustomFormat(issueDate, p.dateFormat))
		}
	}
	return p.generateID(&sb, listInfo.DSSIDAsString())
}

// idAsStringForCertRef gets a String identifier for a given CertificateRef. Port of the
// protected getIdAsStringForCertRef(CertificateRef).
func (p *UserFriendlyIdentifierProvider) idAsStringForCertRef(certificateRef *spi.CertificateRef) string {
	var sb strings.Builder
	sb.WriteString(p.certificatePrefix)
	if certificateRef.ResponderId() != nil && certificateRef.ResponderId().X500Principal() != nil {
		sb.WriteString(stringDelimiter)
		x500PrincipalHelper := model.NewX500PrincipalHelper(certificateRef.ResponderId().X500Principal())
		sb.WriteString(p.humanReadableName(x500PrincipalHelper))

	} else if certificateRef.CertificateIdentifier() != nil {
		if certificateRef.CertificateIdentifier().IssuerName() != nil {
			sb.WriteString(stringDelimiter)
			sb.WriteString(issuerLabel)
			x500PrincipalHelper := model.NewX500PrincipalHelper(certificateRef.CertificateIdentifier().IssuerName())
			sb.WriteString(p.humanReadableName(x500PrincipalHelper))
		}
		if certificateRef.CertificateIdentifier().SerialNumber() != nil {
			sb.WriteString(stringDelimiter)
			sb.WriteString(serialLabel)
			sb.WriteString(certificateRef.CertificateIdentifier().SerialNumber().String())
		}

	} else if certDigest := certificateRef.CertDigest(); certDigest.Value() != nil {
		sb.WriteString(stringDelimiter)
		sb.WriteString(certDigest.HexValue())
	}
	return p.generateID(&sb, certificateRef.DSSIDAsString())
}

// idAsStringForRevRef gets a String identifier for a given RevocationRef. Port of the
// protected getIdAsStringForRevRef(RevocationRef<?>).
func (p *UserFriendlyIdentifierProvider) idAsStringForRevRef(revocationRef revocationRefLike) string {
	var sb strings.Builder
	sb.WriteString(p.revocationRefPrefix(revocationRef))
	sb.WriteString(stringDelimiter)
	sb.WriteString(revocationRef.Digest().HexValue())
	return p.generateID(&sb, revocationRef.DSSIDAsString())
}

// idAsStringForRevTokenIdentifier gets a String identifier for a given
// EncapsulatedRevocationTokenIdentifier. Port of the protected
// getIdAsStringForRevTokenIdentifier(EncapsulatedRevocationTokenIdentifier<?>).
func (p *UserFriendlyIdentifierProvider) idAsStringForRevTokenIdentifier(revocationIdentifier revocationTokenIdentifierLike) string {
	var sb strings.Builder
	sb.WriteString(p.revocationIdentifierPrefix(revocationIdentifier))
	sb.WriteString(stringDelimiter)
	digestValue, err := revocationIdentifier.DigestValue(enumerations.DigestAlgorithm_SHA256)
	if err != nil {
		// Java's DSSException for an unavailable algorithm is unchecked and propagates
		// straight out of getIdAsStringForRevTokenIdentifier (and out of getIdAsString
		// itself, which declares no checked exception either).
		panic(err.Error())
	}
	sb.WriteString(utils.ToHex(digestValue))
	return p.generateID(&sb, revocationIdentifier.AsXmlID())
}

// idAsStringForEvidenceRecordIdentifier gets a String identifier for a given EvidenceRecord.
// Port of the protected getIdAsStringForEvidenceRecordIdentifier(EvidenceRecord).
func (p *UserFriendlyIdentifierProvider) idAsStringForEvidenceRecordIdentifier(evidenceRecord spivalidation.EvidenceRecord) string {
	var sb strings.Builder
	sb.WriteString(p.evidenceRecordPrefix)
	sb.WriteString(stringDelimiter)
	if timestamps := evidenceRecord.Timestamps(); utils.IsCollectionNotEmpty(timestamps) {
		sb.WriteString(p.deterministicIDPartForToken(timestamps[0]))
	} else {
		// Java: `(MultipleDigestIdentifier) evidenceRecord.getDSSId()`, an unconditional cast
		// that throws ClassCastException (an unchecked exception) for any other concrete
		// Identifier - reproduced as the panicking one-result type assertion form.
		multipleDigestIdentifier := evidenceRecord.DSSID().(*model.MultipleDigestIdentifier)
		digestValue, err := multipleDigestIdentifier.DigestValue(enumerations.DigestAlgorithm_SHA256)
		if err != nil {
			panic(err.Error())
		}
		sb.WriteString(utils.ToHex(digestValue))
	}
	return p.generateID(&sb, evidenceRecord.Id())
}

// idAsStringForEAAIdentifier gets an identifier for the EAA. Port of the protected
// getIdAsStringForEAAIdentifier(EAA).
func (p *UserFriendlyIdentifierProvider) idAsStringForEAAIdentifier(eaa spivalidation.EAA) string {
	var sb strings.Builder
	sb.WriteString(p.eaaPrefix)

	payload := eaa.Payload()
	if subject := payload.Subject(); !subject.IsNullOrEmpty() {
		sb.WriteString(stringDelimiter)
		sb.WriteString(subject.StringValue())
	}
	if docType := payload.DocType(); !docType.IsNullOrEmpty() {
		sb.WriteString(stringDelimiter)
		sb.WriteString(docType.StringValue())
	}
	if issuedAtTime := payload.IssuedAtTime(); !issuedAtTime.IsNullOrEmpty() {
		sb.WriteString(stringDelimiter)
		sb.WriteString(spi.DSSUtilsFormatDateWithCustomFormat(*issuedAtTime.DateValue(), p.dateFormat))
	}

	return p.generateID(&sb, eaa.ID())
}

// createIDString builds "<prefix>_<deterministic-id-part>", then routes through generateID for
// duplicate detection. Port of the private createIdString(String, X500PrincipalHelper, Date,
// String).
func (p *UserFriendlyIdentifierProvider) createIDString(prefix string, subject *model.X500PrincipalHelper, creationDate time.Time, dssID string) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(stringDelimiter)
	sb.WriteString(p.deterministicIDPart(subject, creationDate))
	return p.generateID(&sb, dssID)
}

// deterministicIDPartForToken ports the private getDeterministicIdPart(Token) overload.
func (p *UserFriendlyIdentifierProvider) deterministicIDPartForToken(token model.Token) string {
	return p.deterministicIDPart(p.tokenSubject(token), token.CreationDate())
}

// deterministicIDPart ports the private getDeterministicIdPart(X500PrincipalHelper, Date)
// overload: "<subject-or-UNKNOWN-SIGNER>[_<formatted-creation-date>]".
func (p *UserFriendlyIdentifierProvider) deterministicIDPart(subject *model.X500PrincipalHelper, creationDate time.Time) string {
	var sb strings.Builder
	if subject != nil {
		sb.WriteString(p.humanReadableName(subject))
	} else {
		sb.WriteString(unknownSigner)
	}
	if !creationDate.IsZero() {
		sb.WriteString(stringDelimiter)
		sb.WriteString(spi.DSSUtilsFormatDateWithCustomFormat(creationDate, p.dateFormat))
	}
	return sb.String()
}

// humanReadableName ports the private getHumanReadableName(X500PrincipalHelper).
func (p *UserFriendlyIdentifierProvider) humanReadableName(subject *model.X500PrincipalHelper) string {
	name := spi.DSSASN1UtilsHumanReadableNameOfPrincipal(subject)
	if utils.IsStringNotEmpty(name) {
		return p.userFriendlyString(name)
	}
	return unnamedSigner
}

// generateID ports the private generateId(StringBuilder, String): appends a "_<n>" duplicate
// suffix to sb when the id built so far collides with an id already generated for a different
// dssID, then records both the pre-suffix generated id (for future duplicate detection) and
// the final, possibly-suffixed unique id.
//
// Java re-reads stringBuilder.toString() a second time, after the possible append, for the
// value stored in uniqueTokenIdsMap and returned; generatedTokenIdsMap however keeps the
// pre-append value captured at the top of the method. Both reads are reproduced precisely
// (generatedID vs sb.String() below), since the whole duplicate-counting scheme depends on
// comparing against the pre-suffix form on every subsequent call.
func (p *UserFriendlyIdentifierProvider) generateID(sb *strings.Builder, dssID string) string {
	generatedID := sb.String()
	duplicatesNumber := p.duplicatesNumber(generatedID, dssID)
	if duplicatesNumber != 0 {
		duplicatesNumber++
		sb.WriteString(stringDelimiter)
		sb.WriteString(strconv.FormatInt(duplicatesNumber, 10))
	}
	p.generatedTokenIdsMap[dssID] = generatedID

	uniqueID := sb.String()
	p.uniqueTokenIdsMap[dssID] = uniqueID
	return uniqueID
}

// duplicatesNumber ports the private getDuplicatesNumber(String, String): counts how many
// other dssIDs already generated exactly builtID. Map iteration order does not influence the
// result (a count), so no determinism concern applies here.
func (p *UserFriendlyIdentifierProvider) duplicatesNumber(builtID, dssID string) int64 {
	var count int64
	for otherDSSID, generatedValue := range p.generatedTokenIdsMap {
		if otherDSSID != dssID && generatedValue == builtID {
			count++
		}
	}
	return count
}

// tokenPrefix ports the private getTokenPrefix(Token).
//
// Panics for an unsupported token class (Java's thrown IllegalArgumentException).
func (p *UserFriendlyIdentifierProvider) tokenPrefix(token model.Token) string {
	switch token.(type) {
	case *model.CertificateToken:
		return p.certificatePrefix
	case *spi.CRLToken:
		return p.crlPrefix
	case *spi.OCSPToken:
		return p.ocspPrefix
	case *spivalidation.TimestampToken:
		return p.timestampPrefix
	case *spivalidation.EAARevocationToken:
		return p.eaaStatusTokenPrefix
	default:
		panic(fmt.Sprintf("Unsupported token of class '%T' has been reached!", token))
	}
}

// tokenSubject ports the private getTokenSubject(Token).
func (p *UserFriendlyIdentifierProvider) tokenSubject(token model.Token) *model.X500PrincipalHelper {
	if certificateToken, ok := token.(*model.CertificateToken); ok {
		return certificateToken.Subject()
	}
	if issuer := token.IssuerX500Principal(); issuer != nil {
		return model.NewX500PrincipalHelper(issuer)
	}
	return nil
}

// tlInfoPrefix ports the private getTlPrefix(TLInfo).
func (p *UserFriendlyIdentifierProvider) tlInfoPrefix(tlInfo tlInfoLike) string {
	switch tlInfo.(type) {
	case *tsl.PivotInfo:
		return p.pivotPrefix
	case *tsl.LOTLInfo:
		return p.lotlPrefix
	default:
		return p.tlPrefix
	}
}

// loTEInfoPrefix ports the private getLoTEPrefix(LoTEInfo).
//
// NOTE: reproduced verbatim from an upstream bug: for a LoLoTEInfo this returns lotlPrefix
// ("LOTL") rather than lolotePrefix ("LOLOTE"), and for a plain LoTEInfo it returns tlPrefix
// ("TL") rather than lotePrefix ("LOTE") - upstream's getLoTEPrefix body is a byte-for-byte
// copy of getTlPrefix with the type names changed but not the returned prefix fields, so
// lolotePrefix/lotePrefix are unreachable dead fields (see their declarations above).
func (p *UserFriendlyIdentifierProvider) loTEInfoPrefix(listInfo loteInfoLike) string {
	switch listInfo.(type) {
	case *lote.LoLoTEInfo:
		return p.lotlPrefix
	default:
		return p.tlPrefix
	}
}

// revocationRefPrefix ports the private getRevocationRefPrefix(RevocationRef<?>).
//
// Panics for an unsupported reference class (Java's thrown IllegalArgumentException).
func (p *UserFriendlyIdentifierProvider) revocationRefPrefix(revocationRef revocationRefLike) string {
	switch revocationRef.(type) {
	case *spi.CRLRef:
		return p.crlPrefix
	case *spi.OCSPRef:
		return p.ocspPrefix
	default:
		panic(fmt.Sprintf("Unsupported RevocationRef of class '%T' has been reached!", revocationRef))
	}
}

// revocationIdentifierPrefix ports the private
// getRevocationIdentifierPrefix(EncapsulatedRevocationTokenIdentifier<?>).
//
// Panics for an unsupported identifier class (Java's thrown IllegalArgumentException).
func (p *UserFriendlyIdentifierProvider) revocationIdentifierPrefix(revocationIdentifier revocationTokenIdentifierLike) string {
	switch revocationIdentifier.(type) {
	case *crlparser.CRLBinary:
		return p.crlPrefix
	case *spi.OCSPResponseBinary:
		return p.ocspPrefix
	default:
		panic(fmt.Sprintf("Unsupported RevocationTokenIdentifier of class '%T' has been reached!", revocationIdentifier))
	}
}

// userFriendlyString ports the private getUserFriendlyString(String).
func (p *UserFriendlyIdentifierProvider) userFriendlyString(str string) string {
	str = spi.DSSUtilsRemoveControlCharacters(str)
	str = spi.DSSUtilsReplaceInvalidXmlCharacters(str, unsupportedCharacter)
	str = spi.DSSUtilsReplaceAllNonAlphanumericCharacters(str, nameReplacement)
	return trimRepeated(str, nameReplacement)
}

// trimRepeated ports the private trim(String, String): strips every leading and trailing
// occurrence of trimmedStr from str (not just one), stopping once removing another occurrence
// would leave nothing (str no longer strictly longer than trimmedStr).
func trimRepeated(str, trimmedStr string) string {
	for len(str) > len(trimmedStr) && strings.HasPrefix(str, trimmedStr) {
		str = utils.SubstringAfter(str, trimmedStr)
	}
	for len(str) > len(trimmedStr) && strings.HasSuffix(str, trimmedStr) {
		str = str[:len(str)-len(trimmedStr)]
	}
	return str
}

// assertNotBlank ports the private assertNotBlank(String).
//
// Panics with the Java message when str is blank (Java's thrown IllegalArgumentException("The
// prefix cannot be null or blank!")).
func assertNotBlank(str string) {
	if utils.IsStringBlank(str) {
		panic("The prefix cannot be null or blank!")
	}
}
