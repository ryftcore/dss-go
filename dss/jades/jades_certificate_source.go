// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESCertificateSource.java
// (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES (sibling chunks of the same package, not in this manifest - see
// S6_BRIEF.md's VAL chunk split):
//
//	JAdESEtsiUHeader with IsExist() bool and Attributes() []*EtsiUComponent (already assumed,
//	  identically, by dss_json_utils.go's DSSJsonUtilsUnsignedPropertiesWithHeaderName)
//	EtsiUComponent, embedding JAdESAttribute (see jades_attribute.go's header), so
//	  HeaderName()/Value() reach it through promotion
//	JWSHeaderParameterNames* additions this file needs beyond what dss_json_utils.go /
//	  jades_level_baseline_*.go already established: JAdESHeaderParameterNamesAxVals ("axVals"),
//	  JAdESHeaderParameterNamesXRefs ("xRefs"), JAdESHeaderParameterNamesAxRefs ("axRefs"),
//	  JAdESHeaderParameterNamesOtherCert ("otherCert") - following the same "capitalize the JSON
//	  value's first letter" convention already visible in the landed constants (XVals, RVals,
//	  TstVD, Encoding, ...).
//
// VIRTUAL-DISPATCH GAP (frozen package, flagged for integrator - see PORTING.md "no edits to
// frozen packages"): Java's TokenCertificateSource.findTokensFromRefs(refs) loops calling
// `this.findTokensFromCertRef(ref)`, reaching JAdESCertificateSource's override via ordinary
// Java virtual dispatch. spi.TokenCertificateSource.FindTokensFromRefs (already landed, frozen)
// has no overrides-interface hook for FindTokensFromCertRef, so calling it here would statically
// resolve to CommonCertificateSource's base implementation and silently drop the kid/x5u
// resolution FindTokensFromCertRef adds below. KeyIdentifierCertificates therefore reimplements
// the same loop locally (jadesCertificateSourceFindTokensFromRefs), calling this type's own
// FindTokensFromCertRef directly - reproducing upstream's virtual-dispatch outcome without
// touching the frozen spi package.
//
// DEVIATION: extractPublicKey's PublicJsonWebKey#getPublicKey() has no encoded-bytes
// counterpart in internal/jose's JwkHeader (a JWK carries no DER SubjectPublicKeyInfo to
// preserve verbatim, unlike every other model.PublicKey constructed in this port from a
// certificate's own encoding - see model/public_key.go's header). jadesCertificateSourcePublicKey
// re-derives the DER via x509.MarshalPKIXPublicKey, the same synthesis a JCE provider performs
// for PublicJsonWebKey#getPublicKey().getEncoded() upstream.
package jades

import (
	"crypto"
	"crypto/x509"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// JAdESCertificateSource extracts and stores certificates from a JAdES signature. Port of the
// class JAdESCertificateSource, extending spi.SignatureCertificateSource.
type JAdESCertificateSource struct {
	spi.SignatureCertificateSource

	// jws is the JWS Signature to extract certificates from. Port of the private transient
	// final JWS jws field.
	jws *JWS

	// etsiUHeader represents the unsigned 'etsiU' header. Port of the private transient final
	// JAdESEtsiUHeader etsiUHeader field.
	etsiUHeader *JAdESEtsiUHeader

	// kidMap holds 'kid' certificates, when present. Port of the private final Map<String,
	// CertificateToken> kidMap field.
	kidMap map[string]*model.CertificateToken

	// x509UrlMap holds 'x5u' certificates, when present. Port of the private final Map<String,
	// Collection<CertificateToken>> x509UrlMap field.
	x509UrlMap map[string][]*model.CertificateToken
}

// NewJAdESCertificateSource is the default constructor. Port of the public
// JAdESCertificateSource(JWS, JAdESEtsiUHeader) constructor. All certificates are extracted
// during instantiation.
//
// Panics with the Java messages when jws or etsiUHeader is missing (Objects.requireNonNull).
func NewJAdESCertificateSource(jws *JWS, etsiUHeader *JAdESEtsiUHeader) *JAdESCertificateSource {
	if jws == nil {
		panic("JSON Web signature cannot be null")
	}
	if etsiUHeader == nil {
		panic("etsiUHeader cannot be null")
	}

	s := &JAdESCertificateSource{
		jws:         jws,
		etsiUHeader: etsiUHeader,
		kidMap:      make(map[string]*model.CertificateToken),
		x509UrlMap:  make(map[string][]*model.CertificateToken),
	}
	s.InitSignatureCertificateSource(s)

	// signing certificate
	s.extractX5T()
	s.extractX5TS256()
	s.extractX5TO()
	s.extractSigX5Ts()
	s.extractKid()
	s.extractX509Url()

	// certificate chain
	s.extractX5C()

	// unsigned properties
	s.extractEtsiU()

	return s
}

// KeyIdentifierCertificateRefs retrieves the list of CertificateRefs referenced within a 'kid'
// (key identifier) header. Port of getKeyIdentifierCertificateRefs().
func (s *JAdESCertificateSource) KeyIdentifierCertificateRefs() []*spi.CertificateRef {
	return s.CertificateRefsByOrigin(enumerations.CertificateRefOriginKeyIdentifier)
}

// KeyIdentifierCertificates retrieves the CertificateTokens according to a reference present
// within a 'kid' (key identifier) header. Port of getKeyIdentifierCertificates().
func (s *JAdESCertificateSource) KeyIdentifierCertificates() map[string]*model.CertificateToken {
	return s.jadesCertificateSourceFindTokensFromRefs(s.KeyIdentifierCertificateRefs())
}

// jadesCertificateSourceFindTokensFromRefs is the local, virtual-dispatch-correct counterpart of
// spi.TokenCertificateSource.FindTokensFromRefs; see the file header.
func (s *JAdESCertificateSource) jadesCertificateSourceFindTokensFromRefs(certificateRefs []*spi.CertificateRef) map[string]*model.CertificateToken {
	result := make(map[string]*model.CertificateToken)
	for _, certificateRef := range certificateRefs {
		for key, token := range s.FindTokensFromCertRef(certificateRef) {
			result[key] = token
		}
	}
	return result
}

func (s *JAdESCertificateSource) extractX5T() {
	base64UrlSHA1Certificate := s.jws.ProtectedHeaderValueAsString(jose.HeaderX509CertificateThumbprint)
	if utils.IsStringNotEmpty(base64UrlSHA1Certificate) { //nolint:staticcheck // mirrors upstream JAdESCertificateSource#extractX5T: the guard is kept because Java's body is `LOG.warn("Found {} with value {} but not supported by the JAdES standard", ...)` only.
		// Upstream builds a Digest(SHA1, ...) purely to log "Found {} with value {} but not
		// supported by the JAdES standard"; with slf4j dropped per PORTING.md there is nothing
		// left for this branch to do.
	}
}

func (s *JAdESCertificateSource) extractX5TS256() {
	base64UrlSHA256Certificate := s.jws.ProtectedHeaderValueAsString(jose.HeaderX509CertificateSHA256Thumbprint)
	if utils.IsStringNotEmpty(base64UrlSHA256Certificate) {
		certRef := spi.NewCertificateRef()
		certRef.SetCertDigest(model.NewDigest(enumerations.DigestAlgorithmSHA256, DSSJsonUtilsFromBase64Url(base64UrlSHA256Certificate)))
		s.AddCertificateRef(certRef, enumerations.CertificateRefOriginSigningCertificate)
	}
}

func (s *JAdESCertificateSource) extractX5TO() {
	s.extractX5TOFromMap(s.jws.ProtectedHeaderValueAsMap(JAdESHeaderParameterNamesX5tO))
}

func (s *JAdESCertificateSource) extractX5TOFromMap(x5TO *jose.Object) {
	if x5TO.Size() != 0 {
		digest, ok := DSSJsonUtilsDigest(x5TO)
		if ok {
			certRef := spi.NewCertificateRef()
			certRef.SetCertDigest(digest)
			s.AddCertificateRef(certRef, enumerations.CertificateRefOriginSigningCertificate)
		}
	}
}

func (s *JAdESCertificateSource) extractSigX5Ts() {
	sigX5tsList := s.jws.ProtectedHeaderValueAsList(JAdESHeaderParameterNamesSigX5tS)
	for _, item := range sigX5tsList {
		s.extractX5TOFromMap(DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesX5tO))
	}
}

func (s *JAdESCertificateSource) extractKid() {
	kid := s.jws.KeyIDHeaderValue()
	if kid != "" {
		certificateRef := spi.NewCertificateRef()
		issuerSerial := DSSJsonUtilsIssuerSerial(kid)
		if issuerSerial != nil {
			certificateRef.SetCertificateIdentifier(spi.DSSASN1UtilsToSignerIdentifierFromIssuerSerial(issuerSerial))
		} else {
			certificateRef.SetKid(kid)
		}
		s.AddCertificateRef(certificateRef, enumerations.CertificateRefOriginKeyIdentifier)
	}
}

func (s *JAdESCertificateSource) extractX509Url() {
	x5u := s.jws.ProtectedHeaderValueAsString(jose.HeaderX509URL)
	if utils.IsStringNotEmpty(x5u) {
		certificateRef := spi.NewCertificateRef()
		certificateRef.SetX509Url(x5u)
		s.AddCertificateRef(certificateRef, enumerations.CertificateRefOriginX509URL)
	}
}

func (s *JAdESCertificateSource) extractX5C() {
	x509CertChain := s.jws.ProtectedHeaderValueAsList(jose.HeaderX509CertificateChain)
	for _, item := range x509CertChain {
		certificateBase64 := DSSJsonUtilsToString(item)
		if utils.IsStringNotEmpty(certificateBase64) {
			certificate, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(certificateBase64)
			if err != nil {
				// Upstream logs "Unable to decode a certificate from '{}'! Reason : {}".
				continue
			}
			s.AddCertificateWithOrigin(certificate, enumerations.CertificateOriginKeyInfo)
		}
	}
}

func (s *JAdESCertificateSource) extractEtsiU() {
	if !s.etsiUHeader.IsExist() {
		return
	}

	for _, attribute := range s.etsiUHeader.Attributes() {
		s.extractCertificateValues(attribute)
		s.extractAttrAuthoritiesCertValues(attribute)
		s.extractTimestampValidationData(attribute)
		s.extractAnyValidationData(attribute)

		s.extractCompleteCertificateRefs(attribute)
		s.extractAttributeCertificateRefs(attribute)
	}
}

func (s *JAdESCertificateSource) extractCertificateValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesXVals == attribute.HeaderName() {
		s.extractCertificateValuesFromList(DSSJsonUtilsToList(attribute.Value(), JAdESHeaderParameterNamesXVals),
			enumerations.CertificateOriginCertificateValues)
	}
}

func (s *JAdESCertificateSource) extractAttrAuthoritiesCertValues(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesAxVals == attribute.HeaderName() {
		s.extractCertificateValuesFromList(DSSJsonUtilsToList(attribute.Value(), JAdESHeaderParameterNamesAxVals),
			enumerations.CertificateOriginAttrAuthoritiesCertValues)
	}
}

func (s *JAdESCertificateSource) extractTimestampValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesTstVD, enumerations.CertificateOriginTimestampValidationData)
}

func (s *JAdESCertificateSource) extractAnyValidationData(attribute *EtsiUComponent) {
	s.extractValidationData(attribute, JAdESHeaderParameterNamesAnyValData, enumerations.CertificateOriginAnyValidationData)
}

func (s *JAdESCertificateSource) extractValidationData(attribute *EtsiUComponent, headerName string, origin enumerations.CertificateOrigin) {
	if headerName == attribute.HeaderName() {
		tstVd := DSSJsonUtilsToMap(attribute.Value(), headerName)
		xVals := DSSJsonUtilsGetAsList(tstVd, JAdESHeaderParameterNamesXVals)
		if utils.IsCollectionNotEmpty(xVals) {
			s.extractCertificateValuesFromList(xVals, origin)
		}
	}
}

func (s *JAdESCertificateSource) extractCompleteCertificateRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesXRefs == attribute.HeaderName() {
		s.extractCertificateRefsFromList(DSSJsonUtilsToList(attribute.Value(), JAdESHeaderParameterNamesXRefs),
			enumerations.CertificateRefOriginCompleteCertificateRefs)
	}
}

func (s *JAdESCertificateSource) extractAttributeCertificateRefs(attribute *EtsiUComponent) {
	if JAdESHeaderParameterNamesAxRefs == attribute.HeaderName() {
		s.extractCertificateRefsFromList(DSSJsonUtilsToList(attribute.Value(), JAdESHeaderParameterNamesAxRefs),
			enumerations.CertificateRefOriginAttributeCertificateRefs)
	}
}

func (s *JAdESCertificateSource) extractCertificateValuesFromList(xVals []any, origin enumerations.CertificateOrigin) {
	for _, item := range xVals {
		xVal := DSSJsonUtilsToMapValue(item)
		x509Cert := DSSJsonUtilsGetAsMap(xVal, JAdESHeaderParameterNamesX509Cert)
		otherCert := DSSJsonUtilsGetAsMap(xVal, JAdESHeaderParameterNamesOtherCert)
		if x509Cert.Size() != 0 {
			s.extractX509Cert(x509Cert, origin)
		} else if otherCert.Size() != 0 {
			// Upstream logs "The header '{}' is not supported! The entry is skipped.".
		}
	}
}

func (s *JAdESCertificateSource) extractCertificateRefsFromList(xRefs []any, origin enumerations.CertificateRefOrigin) {
	for _, item := range xRefs {
		xref := DSSJsonUtilsToMapValue(item)
		certificateRef := JAdESCertificateRefExtractionUtilsCreateCertificateRef(xref)
		if certificateRef != nil {
			s.AddCertificateRef(certificateRef, origin)
		}
	}
}

func (s *JAdESCertificateSource) extractX509Cert(x509Cert *jose.Object, origin enumerations.CertificateOrigin) {
	encoding := DSSJsonUtilsGetAsString(x509Cert, JAdESHeaderParameterNamesEncoding)
	if utils.IsStringEmpty(encoding) || utils.AreStringsEqual(enumerations.PKIEncodingDER.URI(), encoding) {
		val := DSSJsonUtilsGetAsString(x509Cert, JAdESHeaderParameterNamesVal)
		if utils.IsStringNotEmpty(val) {
			certificate, err := spi.DSSUtilsLoadCertificateFromBase64EncodedString(val)
			if err != nil {
				// Upstream logs "Unable to decode a certificate from '{}'! Reason : {}".
				return
			}
			s.AddCertificateWithOrigin(certificate, origin)
		}
	} else {
		// Upstream logs "Unsupported encoding header value : '{}'".
	}
}

// ExtractCandidatesForSigningCertificate implements spi.SignatureCertificateSourceOverrides.
// Port of the protected extractCandidatesForSigningCertificate(CertificateSource) override.
func (s *JAdESCertificateSource) ExtractCandidatesForSigningCertificate(signingCertificateSource spi.CertificateSource) *spi.CandidatesForSigningCertificate {
	candidatesForSigningCertificate := s.InitCandidatesList(signingCertificateSource)
	if !candidatesForSigningCertificate.IsEmpty() {
		return candidatesForSigningCertificate
	}

	for _, certificateToken := range s.KeyInfoCertificates() {
		candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
	}

	// if x5u does not contain certificates, check other certificates embedded into the signature
	if candidatesForSigningCertificate.IsEmpty() {

		// From JWK (not JAdES)
		publicKey := s.extractPublicKey()
		if publicKey != nil {
			candidatesForSigningCertificate.Add(spi.NewCertificateValidityFromPublicKey(publicKey))

		} else {
			// Add all found certificates
			for _, certificateToken := range s.Certificates() {
				candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
			}
		}
	}

	if signingCertificateSource != nil {
		s.resolveFromSource(signingCertificateSource, candidatesForSigningCertificate)
	}

	s.checkSigningCertificateRef(candidatesForSigningCertificate)

	return candidatesForSigningCertificate
}

func (s *JAdESCertificateSource) resolveFromSource(signingCertificateSource spi.CertificateSource, candidatesForSigningCertificate *spi.CandidatesForSigningCertificate) {
	kidCandidate := s.resolveByKid(signingCertificateSource)
	if kidCandidate != nil {
		// Upstream logs "Resolved certificate by kid".
		s.SignatureCertificateSource.AddCertificate(kidCandidate)
		candidatesForSigningCertificate.Add(spi.NewCertificateValidity(kidCandidate))
		return
	}

	uriCandidates := s.resolveByUri(signingCertificateSource)
	if utils.IsCollectionNotEmpty(uriCandidates) {
		// Upstream logs "Resolved certificates by x5u".
		for _, externalCandidate := range uriCandidates {
			s.SignatureCertificateSource.AddCertificate(externalCandidate)
			candidatesForSigningCertificate.Add(spi.NewCertificateValidity(externalCandidate))
		}
		return
	}

	certificateDigest := s.signingCertificateDigest()
	if certificateDigest != nil {
		certificatesByDigest := signingCertificateSource.ByCertificateDigest(*certificateDigest)
		if utils.IsMapNotEmpty(certificatesByDigest) {
			// Upstream logs "Resolved certificate by digest".
			for _, certificateToken := range certificatesByDigest {
				candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
			}
		}

	} else {
		certificates := signingCertificateSource.Certificates()
		// Upstream logs "No signing certificate reference found. Resolve all {} certificates
		// from the provided certificate source as signing candidates.".
		for _, certCandidate := range certificates {
			candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certCandidate))
		}
	}
}

func (s *JAdESCertificateSource) resolveByKid(signingCertificateSource spi.CertificateSource) *model.CertificateToken {
	kidHeader := s.jws.KeyIDHeaderValue()
	if utils.IsStringNotEmpty(kidHeader) {
		if kidCertificateSource, ok := signingCertificateSource.(*spi.KidCertificateSource); ok {
			certificateByKid := kidCertificateSource.CertificateByKid(kidHeader)
			if certificateByKid != nil {
				s.kidMap[kidHeader] = certificateByKid
			}
			return certificateByKid
		}
		// Upstream logs "JWS/JAdES contains a 'kid' header (provide a KidCertificateSource to
		// resolve it)".
	}
	return nil
}

func (s *JAdESCertificateSource) resolveByUri(signingCertificateSource spi.CertificateSource) []*model.CertificateToken {
	x5uHeader := s.jws.ProtectedHeaderValueAsString(jose.HeaderX509URL)
	if utils.IsStringNotEmpty(x5uHeader) {
		if x509URLCertificateSource, ok := signingCertificateSource.(spi.X509URLCertificateSource); ok {
			certificatesByUri := x509URLCertificateSource.CertificatesByUrl(x5uHeader)
			if utils.IsCollectionNotEmpty(certificatesByUri) {
				s.x509UrlMap[x5uHeader] = certificatesByUri
			}
			return certificatesByUri
		}
		// Upstream logs "JWS/JAdES contains a 'x5u' header (provide a X509URLCertificateSource
		// to resolve it)".
	}
	return nil
}

// extractPublicKey ports the private extractPublicKey(); see the file header DEVIATION note.
func (s *JAdESCertificateSource) extractPublicKey() *model.PublicKey {
	key, err := s.jws.JwkHeader()
	if err != nil {
		// Upstream logs "Unable to extract the public key".
		return nil
	}
	return jadesCertificateSourcePublicKey(key)
}

// jadesCertificateSourcePublicKey wraps a crypto.PublicKey parsed from a 'jwk' header into a
// model.PublicKey; see the file header DEVIATION note.
func jadesCertificateSourcePublicKey(key crypto.PublicKey) *model.PublicKey {
	if key == nil {
		return nil
	}
	encoded, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil
	}
	return model.NewPublicKeyFromEncoded(encoded, key)
}

func (s *JAdESCertificateSource) checkSigningCertificateRef(candidates *spi.CandidatesForSigningCertificate) {
	var signingCertRef *spi.CertificateRef
	potentialSigningCertificates := s.SigningCertificateRefs()
	if utils.IsCollectionNotEmpty(potentialSigningCertificates) {
		// first reference shall be a reference to a signing certificate
		signingCertRef = potentialSigningCertificates[0]
	}

	var kidCertRef *spi.CertificateRef
	keyIdentifierCertificateRefs := s.KeyIdentifierCertificateRefs()
	if utils.IsCollectionNotEmpty(keyIdentifierCertificateRefs) {
		kidCertRef = keyIdentifierCertificateRefs[0]
	}

	if signingCertRef != nil {
		var bestCertificateValidity *spi.CertificateValidity
		// check all certificates against the signingCert ref and find the best one
		certificateValidityList := candidates.CertificateValidityList()
		for _, certificateValidity := range certificateValidityList {
			if s.isValid(certificateValidity, signingCertRef, kidCertRef) {
				bestCertificateValidity = certificateValidity
			}
		}
		if bestCertificateValidity != nil {
			_ = candidates.SetTheCertificateValidity(bestCertificateValidity)
		}
	}
}

func (s *JAdESCertificateSource) isValid(certificateValidity *spi.CertificateValidity, signingCertRef, kidCertRef *spi.CertificateRef) bool {
	certificateValidity.SetDigestPresent(signingCertRef != nil && !signingCertRef.CertDigest().IsEmpty())
	certificateValidity.SetIssuerSerialPresent(kidCertRef != nil && kidCertRef.CertificateIdentifier() != nil)

	certificateToken := certificateValidity.CertificateToken()
	if certificateToken != nil {
		certificateMatcher := s.CertificateMatcher()
		if signingCertRef != nil {
			certificateValidity.SetDigestEqual(certificateMatcher.MatchByDigest(certificateToken, signingCertRef))
		}
		if kidCertRef != nil {
			certificateValidity.SetSerialNumberEqual(certificateMatcher.MatchBySerialNumber(certificateToken, kidCertRef))
			certificateValidity.SetDistinguishedNameEqual(certificateMatcher.MatchByIssuerName(certificateToken, kidCertRef))
		}
	}
	return certificateValidity.IsValid()
}

func (s *JAdESCertificateSource) signingCertificateDigest() *model.Digest {
	signingCertificateRefs := s.SigningCertificateRefs()
	if utils.IsCollectionNotEmpty(signingCertificateRefs) {
		// must contain only one reference
		signingCert := signingCertificateRefs[0]
		digest := signingCert.CertDigest()
		return &digest
	}
	return nil
}

// OrphanCertificateRefs returns the list of CertificateRefs left without a matched
// CertificateToken, plus any 'x5u' hint refs not carried by the base implementation.
// Port of the getOrphanCertificateRefs() override.
func (s *JAdESCertificateSource) OrphanCertificateRefs() []*spi.CertificateRef {
	certRefs := s.SignatureCertificateSource.OrphanCertificateRefs()
	x509CertUriRefs := s.CertificateRefsByOrigin(enumerations.CertificateRefOriginX509URL)
	for _, certificateRef := range x509CertUriRefs {
		if !jadesCertificateSourceContainsRef(certRefs, certificateRef) {
			certRefs = append(certRefs, certificateRef)
		}
	}
	return certRefs
}

func jadesCertificateSourceContainsRef(refs []*spi.CertificateRef, ref *spi.CertificateRef) bool {
	for _, existing := range refs {
		if existing.Equals(ref) {
			return true
		}
	}
	return false
}

// ReferencesForCertificateToken returns the CertificateRefs matching certificateToken, adding
// those resolved through the 'kid'/'x5u' maps. Port of the getReferencesForCertificateToken(
// CertificateToken) override.
//
// kidMap/x509UrlMap are plain java.util.HashMaps upstream (not LinkedHashMap), so their Java
// iteration order is already unspecified; iterating the Go maps below in Go's own (also
// unspecified) order reproduces that same absence of an ordering guarantee rather than a
// deviation from it.
func (s *JAdESCertificateSource) ReferencesForCertificateToken(certificateToken *model.CertificateToken) []*spi.CertificateRef {
	result := s.SignatureCertificateSource.ReferencesForCertificateToken(certificateToken)
	for kid, token := range s.kidMap {
		if token.Equals(certificateToken) {
			for _, certificateRef := range s.CertificateRefsByOrigin(enumerations.CertificateRefOriginKeyIdentifier) {
				if kid == certificateRef.Kid() {
					result = append(result, certificateRef)
				}
			}
		}
	}
	for x5u, tokens := range s.x509UrlMap {
		if jadesCertificateSourceContainsToken(tokens, certificateToken) {
			for _, certificateRef := range s.CertificateRefsByOrigin(enumerations.CertificateRefOriginX509URL) {
				if x5u == certificateRef.X509Url() {
					result = append(result, certificateRef)
				}
			}
		}
	}
	return result
}

func jadesCertificateSourceContainsToken(tokens []*model.CertificateToken, certificateToken *model.CertificateToken) bool {
	for _, token := range tokens {
		if token.Equals(certificateToken) {
			return true
		}
	}
	return false
}

// FindTokensFromCertRef returns the certificate tokens for the provided CertificateRef, keyed by
// DSSIDAsString(), adding the 'kid'/'x5u' matches on top of the base lookup. Port of the
// findTokensFromCertRef(CertificateRef) override.
func (s *JAdESCertificateSource) FindTokensFromCertRef(certificateRef *spi.CertificateRef) map[string]*model.CertificateToken {
	certificates := s.SignatureCertificateSource.FindTokensFromCertRef(certificateRef)
	if certificates == nil {
		certificates = make(map[string]*model.CertificateToken)
	}
	if utils.IsStringNotEmpty(certificateRef.Kid()) {
		if certificateTokenByKid, ok := s.kidMap[certificateRef.Kid()]; ok && certificateTokenByKid != nil {
			certificates[certificateTokenByKid.DSSIDAsString()] = certificateTokenByKid
		}
	}
	if utils.IsStringNotEmpty(certificateRef.X509Url()) {
		for _, certificateToken := range s.x509UrlMap[certificateRef.X509Url()] {
			certificates[certificateToken.DSSIDAsString()] = certificateToken
		}
	}
	return certificates
}

// compile-time assertion: a JAdESCertificateSource satisfies its own overrides contract.
var _ spi.SignatureCertificateSourceOverrides = (*JAdESCertificateSource)(nil)
