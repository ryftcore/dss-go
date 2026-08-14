// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/RevocationDataVerifier.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCIES (flagged per S2B_BRIEF.md): this file is one of the contracts
// signature_validation_context.go (sibling chunk) already documents its expectations of; the
// public surface below matches that header comment's "RevocationDataVerifier (VAL-B)" entry
// verbatim (constructor names, TrustAnchorVerifier()/SetTrustAnchorVerifier,
// SetValidationContext, IsAcceptable, CheckCertificateNotRevoked, IsRevocationDataSkip,
// IsRevocationDataFresh, IsAcceptableForChain, IsAfterThisUpdateAndBeforeNextUpdate). Two types
// referenced here are owned by sibling chunks and used opaquely / by their documented shape:
//
//   - TrustAnchorVerifier (VAL-C): NewDefaultTrustAnchorVerifier() *TrustAnchorVerifier,
//     IsTrustedAtTime(cert *model.CertificateToken, controlTime time.Time, context enumerations.Context) bool
//     (see signature_validation_context.go's own header comment).
//   - ValidationContext (VAL-C): the interface SignatureValidationContext implements. This file
//     only needs two of its members, inferred from every ValidationContext call the Java source
//     makes (addCertificateTokenForVerification(CertificateToken), getRevocationData(CertificateToken)):
//     AddCertificateTokenForVerification(certificateToken *model.CertificateToken)
//     GetRevocationData(certificateToken *model.CertificateToken) []AnyRevocationToken
//     (integration note: originally written as RevocationData per PORTING.md's get-dropped
//     convention; renamed to GetRevocationData to match ValidationContext's landed interface,
//     which keeps the get-prefix for consistency with SignatureValidationContext's other
//     getters - see validation_context.go's header for the full cross-chunk conflict writeup.)
//
// AnyRevocationToken (the Go analogue of Java's raw RevocationToken<?>) is already defined in
// crl_first_revocation_data_loading_strategy.go, landed earlier in this same chunk.
package validation

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
)

// revocationDataVerifierDefaultDigestAlgorithms is the default collection of Digest Algorithms
// to accept from CRL/OCSP responders, synchronized with ETSI 119 312 V1.4.2.
var revocationDataVerifierDefaultDigestAlgorithms = []enumerations.DigestAlgorithm{
	enumerations.DigestAlgorithm_SHA224, enumerations.DigestAlgorithm_SHA256,
	enumerations.DigestAlgorithm_SHA384, enumerations.DigestAlgorithm_SHA512,
	enumerations.DigestAlgorithm_SHA3_256, enumerations.DigestAlgorithm_SHA3_384,
	enumerations.DigestAlgorithm_SHA3_512,
}

// revocationDataVerifierDefaultEncryptionAlgorithmsKeyLength is the default map of acceptable
// Encryption Algorithms and their corresponding minimal key length.
var revocationDataVerifierDefaultEncryptionAlgorithmsKeyLength = map[enumerations.EncryptionAlgorithm]int{
	enumerations.EncryptionAlgorithm_DSA:         2048,
	enumerations.EncryptionAlgorithm_RSA:         1900,
	enumerations.EncryptionAlgorithm_RSASSA_PSS:  1900,
	enumerations.EncryptionAlgorithm_ECDSA:       256,
	enumerations.EncryptionAlgorithm_PLAIN_ECDSA: 256,
}

// revocationDataVerifierDefaultMaximumRevocationFreshness is the default maximum revocation
// freshness (0 : revocation must be issued after the best-signature-time / lowest POE).
const revocationDataVerifierDefaultMaximumRevocationFreshness = int64(0)

// revocationDataVerifierDefaultRevocationSkipCertificateExtensions is the default collection of
// certificate extension OIDs indicating the revocation check is not required for those
// certificates: valassured-ST-certs, id_pkix_ocsp_nocheck and noRevAvail.
var revocationDataVerifierDefaultRevocationSkipCertificateExtensions = []string{
	spi.OID_id_etsi_ext_valassured_ST_certs.String(),
	enumerations.CertificateExtensionEnum_OCSP_NOCHECK.OID(),
	enumerations.CertificateExtensionEnum_NO_REVOCATION_AVAILABLE.OID(),
}

// RevocationDataVerifier is used to verify acceptance of a revocation data for the following
// validation process, whether the revocation data has been extracted from a document or
// obtained from an online source. The class verifies the consistency of the given revocation
// information and applicability of the used cryptographic constraints used to create this
// token.
//
// NOTE: It is not recommended to use a single instance of RevocationDataVerifier within
// different CertificateVerifiers, as it may lead to concurrency issues during the execution in
// multi-threaded environments (Go: this port is not goroutine-safe, matching the rest of the
// value objects in this package). Please use a new RevocationDataVerifier per each
// CertificateVerifier.
type RevocationDataVerifier struct {
	// acceptableDigestAlgorithms is a collection of Digest Algorithms to accept from CRL/OCSP
	// responders.
	acceptableDigestAlgorithms []enumerations.DigestAlgorithm

	// acceptableEncryptionAlgorithmKeyLength maps acceptable Encryption Algorithms to their
	// corresponding minimal acceptable key length.
	acceptableEncryptionAlgorithmKeyLength map[enumerations.EncryptionAlgorithm]int

	// acceptableSignatureAlgorithmKeyLength maps acceptable Signature Algorithms to their
	// corresponding minimal acceptable key length.
	acceptableSignatureAlgorithmKeyLength map[enumerations.SignatureAlgorithm]int

	// revocationSkipCertificateExtensions is a collection of certificate extension identifiers
	// indicating the revocation check is not required for those certificates.
	revocationSkipCertificateExtensions []string

	// revocationSkipCertificatePolicies is a collection of certificate policy identifiers
	// indicating the revocation check is not required for those certificates.
	revocationSkipCertificatePolicies []string

	// signatureMaximumRevocationFreshness defines the maximum allowed revocation freshness (in
	// milliseconds) for a signature's certificate chain. nil disables the check.
	signatureMaximumRevocationFreshness *int64

	// timestampMaximumRevocationFreshness defines the maximum allowed revocation freshness (in
	// milliseconds) for a timestamp's certificate chain. nil disables the check.
	timestampMaximumRevocationFreshness *int64

	// revocationMaximumRevocationFreshness defines the maximum allowed revocation freshness (in
	// milliseconds) for a revocation's certificate chain. nil disables the check.
	revocationMaximumRevocationFreshness *int64

	// checkRevocationFreshnessNextUpdate, when true and no revocation maximum freshness is
	// defined for the given context, enforces revocation freshness check using the difference
	// between revocation's nextUpdate and thisUpdate as the maximum acceptable freshness.
	checkRevocationFreshnessNextUpdate bool

	// acceptTimestampCertificatesWithoutRevocation indicates whether timestamp certificates
	// without revocation data should be accepted.
	acceptTimestampCertificatesWithoutRevocation bool

	// acceptRevocationCertificatesWithoutRevocation indicates whether revocation certificates
	// without revocation data should be accepted.
	acceptRevocationCertificatesWithoutRevocation bool

	// trustAnchorVerifier verifies whether a given certificate token is a trust anchor at the
	// control time.
	trustAnchorVerifier *TrustAnchorVerifier

	// validationContext is the signature validation context.
	validationContext ValidationContext
}

// NewEmptyRevocationDataVerifier creates an empty instance of RevocationDataVerifier. All
// constraints should be configured manually. Port of createEmptyRevocationDataVerifier().
func NewEmptyRevocationDataVerifier() *RevocationDataVerifier {
	return &RevocationDataVerifier{}
}

// NewDefaultRevocationDataVerifier instantiates a new RevocationDataVerifier, using the
// default validation constraints (synchronized with default validation policy). Port of
// createDefaultRevocationDataVerifier().
func NewDefaultRevocationDataVerifier() *RevocationDataVerifier {
	v := &RevocationDataVerifier{}
	v.SetAcceptableDigestAlgorithms(revocationDataVerifierDefaultDigestAlgorithms)
	v.SetAcceptableEncryptionAlgorithmKeyLength(revocationDataVerifierDefaultEncryptionAlgorithmsKeyLength)
	v.SetRevocationSkipCertificateExtensions(revocationDataVerifierDefaultRevocationSkipCertificateExtensions)
	// #revocationSkipCertificatePolicies are empty
	defaultFreshness := revocationDataVerifierDefaultMaximumRevocationFreshness
	v.SetSignatureMaximumRevocationFreshness(&defaultFreshness)
	v.SetTimestampMaximumRevocationFreshness(&defaultFreshness)
	v.SetRevocationMaximumRevocationFreshness(&defaultFreshness)
	// #checkRevocationFreshnessNextUpdate is false
	// #acceptRevocationIssuersWithoutRevocation is false
	return v
}

// SetAcceptableDigestAlgorithms sets a collection of Digest Algorithms for acceptance. Port of
// setAcceptableDigestAlgorithms(Collection); Java's Objects.requireNonNull becomes a panic.
func (v *RevocationDataVerifier) SetAcceptableDigestAlgorithms(acceptableDigestAlgorithms []enumerations.DigestAlgorithm) {
	if acceptableDigestAlgorithms == nil {
		panic("Collection of DigestAlgorithms for acceptance cannot be null!")
	}
	v.acceptableDigestAlgorithms = acceptableDigestAlgorithms
}

// SetAcceptableEncryptionAlgorithmKeyLength sets a map of acceptable Encryption Algorithms and
// their corresponding minimal key length values. Port of
// setAcceptableEncryptionAlgorithmKeyLength(Map); Java's Objects.requireNonNull becomes a panic.
func (v *RevocationDataVerifier) SetAcceptableEncryptionAlgorithmKeyLength(acceptableEncryptionAlgorithmKeyLength map[enumerations.EncryptionAlgorithm]int) {
	if acceptableEncryptionAlgorithmKeyLength == nil {
		panic("Map of Encryption Algorithms for acceptance cannot be null!")
	}
	v.acceptableEncryptionAlgorithmKeyLength = acceptableEncryptionAlgorithmKeyLength
}

// SetAcceptableSignatureAlgorithmKeyLength sets a map of acceptable Signature Algorithms and
// their corresponding minimal key length values. Port of
// setAcceptableSignatureAlgorithmKeyLength(Map); Java's Objects.requireNonNull becomes a panic.
func (v *RevocationDataVerifier) SetAcceptableSignatureAlgorithmKeyLength(acceptableSignatureAlgorithmKeyLength map[enumerations.SignatureAlgorithm]int) {
	if acceptableSignatureAlgorithmKeyLength == nil {
		panic("Map of Signature Algorithms for acceptance cannot be null!")
	}
	v.acceptableSignatureAlgorithmKeyLength = acceptableSignatureAlgorithmKeyLength
}

// SetRevocationSkipCertificateExtensions sets a collection of certificate extension OIDs
// indicating the revocation check shall be skipped for the given certificate. Port of
// setRevocationSkipCertificateExtensions(Collection).
func (v *RevocationDataVerifier) SetRevocationSkipCertificateExtensions(revocationSkipCertificateExtensions []string) {
	v.revocationSkipCertificateExtensions = revocationSkipCertificateExtensions
}

// SetRevocationSkipCertificatePolicies sets a collection of certificate policy OIDs indicating
// the revocation check shall be skipped for the given certificate. Port of
// setRevocationSkipCertificatePolicies(Collection).
func (v *RevocationDataVerifier) SetRevocationSkipCertificatePolicies(revocationSkipCertificatePolicies []string) {
	v.revocationSkipCertificatePolicies = revocationSkipCertificatePolicies
}

// SetSignatureMaximumRevocationFreshness sets the maximum accepted freshness (in milliseconds)
// for revocation data issued for a signature's certificate chain certificates. nil disables
// the check. Port of setSignatureMaximumRevocationFreshness(Long).
func (v *RevocationDataVerifier) SetSignatureMaximumRevocationFreshness(signatureMaximumRevocationFreshness *int64) {
	v.signatureMaximumRevocationFreshness = signatureMaximumRevocationFreshness
}

// SetTimestampMaximumRevocationFreshness sets the maximum accepted freshness (in milliseconds)
// for revocation data issued for a timestamp's certificate chain certificates. nil disables the
// check. Port of setTimestampMaximumRevocationFreshness(Long).
func (v *RevocationDataVerifier) SetTimestampMaximumRevocationFreshness(timestampMaximumRevocationFreshness *int64) {
	v.timestampMaximumRevocationFreshness = timestampMaximumRevocationFreshness
}

// SetRevocationMaximumRevocationFreshness sets the maximum accepted freshness (in milliseconds)
// for revocation data issued for a revocation data's certificate chain certificates (CRL or
// OCSP). nil disables the check. Port of setRevocationMaximumRevocationFreshness(Long).
func (v *RevocationDataVerifier) SetRevocationMaximumRevocationFreshness(revocationMaximumRevocationFreshness *int64) {
	v.revocationMaximumRevocationFreshness = revocationMaximumRevocationFreshness
}

// SetCheckRevocationFreshnessNextUpdate sets whether the difference between revocation's
// nextUpdate and thisUpdate fields shall be taken as a maximum acceptable revocation freshness
// in case no maximum revocation freshness constraint is defined for the given context. Port of
// setCheckRevocationFreshnessNextUpdate(boolean).
func (v *RevocationDataVerifier) SetCheckRevocationFreshnessNextUpdate(checkRevocationFreshnessNextUpdate bool) {
	v.checkRevocationFreshnessNextUpdate = checkRevocationFreshnessNextUpdate
}

// SetAcceptTimestampCertificatesWithoutRevocation sets whether a timestamp certificate without
// a valid revocation data should be accepted by the verifier. Port of
// setAcceptTimestampCertificatesWithoutRevocation(boolean).
func (v *RevocationDataVerifier) SetAcceptTimestampCertificatesWithoutRevocation(acceptTimestampCertificatesWithoutRevocation bool) {
	v.acceptTimestampCertificatesWithoutRevocation = acceptTimestampCertificatesWithoutRevocation
}

// SetAcceptRevocationCertificatesWithoutRevocation sets whether a revocation certificate
// without a valid revocation data should be accepted by the verifier. Port of
// setAcceptRevocationCertificatesWithoutRevocation(boolean).
func (v *RevocationDataVerifier) SetAcceptRevocationCertificatesWithoutRevocation(acceptRevocationCertificatesWithoutRevocation bool) {
	v.acceptRevocationCertificatesWithoutRevocation = acceptRevocationCertificatesWithoutRevocation
}

// TrustAnchorVerifier gets the trust anchor verifier. This method is used internally within
// SignatureValidationContext to identify whether the configuration is already present and a
// trustAnchorVerifier should be set. Port of getTrustAnchorVerifier().
func (v *RevocationDataVerifier) TrustAnchorVerifier() *TrustAnchorVerifier {
	return v.trustAnchorVerifier
}

// SetTrustAnchorVerifier sets whether a certificate token can be considered as a trust anchor
// at the given control time. Port of setTrustAnchorVerifier(TrustAnchorVerifier).
func (v *RevocationDataVerifier) SetTrustAnchorVerifier(trustAnchorVerifier *TrustAnchorVerifier) {
	v.trustAnchorVerifier = trustAnchorVerifier
}

// SetValidationContext sets the validation context for certificates validation. Port of the
// protected setValidationContext(ValidationContext).
func (v *RevocationDataVerifier) SetValidationContext(validationContext ValidationContext) {
	v.validationContext = validationContext
}

// IsAcceptableNow verifies the validity of the given RevocationToken using the embedded issuer
// certificate token at the current time. Port of the isAcceptable(RevocationToken) overload.
func (v *RevocationDataVerifier) IsAcceptableNow(revocationToken AnyRevocationToken) bool {
	return v.IsAcceptable(revocationToken, time.Now())
}

// IsAcceptable verifies the validity of the given RevocationToken at the given controlTime
// using the embedded issuer certificate token. Port of the isAcceptable(RevocationToken, Date)
// overload; exported under this exact two-argument name because
// signature_validation_context.go's "RevocationDataVerifier (VAL-B)" header comment already
// depends on it: "IsAcceptable(token AnyRevocationToken, controlTime time.Time) bool".
func (v *RevocationDataVerifier) IsAcceptable(revocationToken AnyRevocationToken, controlTime time.Time) bool {
	return v.isAcceptableWithIssuer(revocationToken, revocationToken.IssuerCertificateToken(), controlTime)
}

// isAcceptableWithIssuer verifies the validity of the given RevocationToken at controlTime.
// Port of the isAcceptable(RevocationToken, CertificateToken, Date) overload.
func (v *RevocationDataVerifier) isAcceptableWithIssuer(revocationToken AnyRevocationToken, issuerCertificateToken *model.CertificateToken, controlTime time.Time) bool {
	return v.IsAcceptableForChain(revocationToken, issuerCertificateToken, nil, controlTime)
}

// IsAcceptableForChain verifies the validity of the given RevocationToken at controlTime
// against the given issuer certificate chain. Port of the
// isAcceptable(RevocationToken, CertificateToken, List, Date) overload.
func (v *RevocationDataVerifier) IsAcceptableForChain(revocationToken AnyRevocationToken, issuerCertificateToken *model.CertificateToken,
	certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	return v.isRevocationTokenValid(revocationToken) && v.isRevocationDataComplete(revocationToken) &&
		v.isGoodIssuer(revocationToken, issuerCertificateToken, controlTime) &&
		v.IsCertificateChainValid(certificateChain, controlTime, enumerations.Context_REVOCATION) && v.isConsistent(revocationToken) &&
		v.isAcceptableSignatureAlgorithm(revocationToken, issuerCertificateToken)
}

// isRevocationTokenValid verifies whether the revocation token is cryptographically valid.
// Port of isRevocationTokenValid(RevocationToken).
func (v *RevocationDataVerifier) isRevocationTokenValid(revocationToken AnyRevocationToken) bool {
	// upstream logs the invalidity reason here, dropped per PORTING.md's slf4j rule.
	return revocationToken.IsValid()
}

// isRevocationDataComplete verifies whether the revocation token contains all required data.
// Port of isRevocationDataComplete(RevocationToken).
func (v *RevocationDataVerifier) isRevocationDataComplete(revocationToken AnyRevocationToken) bool {
	if revocationToken.RelatedCertificate() == nil {
		return false
	}
	if revocationToken.Status() == "" {
		return false
	}
	if revocationToken.ThisUpdate().IsZero() {
		return false
	}
	return true
}

// isGoodIssuer verifies validity of the issuerCertificateToken of revocationToken. Port of
// isGoodIssuer(RevocationToken, CertificateToken, Date).
func (v *RevocationDataVerifier) isGoodIssuer(revocationToken AnyRevocationToken, issuerCertificateToken *model.CertificateToken, controlTime time.Time) bool {
	if issuerCertificateToken == nil {
		return false
	}
	if enumerations.RevocationType_OCSP == revocationToken.RevocationType() &&
		!spi.DSSRevocationUtilsCheckIssuerValidAtRevocationProductionTime(revocationToken, issuerCertificateToken) {
		return false
	}
	if !v.isCertificateValid(issuerCertificateToken, revocationToken.Certificates(), controlTime) {
		return false
	}
	return true
}

// isConsistent verifies whether the revocation token is consistent. Port of
// isConsistent(RevocationToken).
func (v *RevocationDataVerifier) isConsistent(revocation AnyRevocationToken) bool {
	certToken := revocation.RelatedCertificate()

	if !isRevocationIssuedAfterCertificateNotBefore(revocation, certToken) {
		return false
	}
	if !doesRevocationKnowCertificate(revocation, certToken) {
		return false
	}
	return true
}

// isRevocationIssuedAfterCertificateNotBefore is the port of the private static
// isRevocationIssuedAfterCertificateNotBefore(RevocationToken, CertificateToken).
func isRevocationIssuedAfterCertificateNotBefore(revocationToken AnyRevocationToken, certificateToken *model.CertificateToken) bool {
	return !certificateToken.NotBefore().After(revocationToken.ThisUpdate())
}

// doesRevocationKnowCertificate is the port of the private static
// doesRevocationKnowCertificate(RevocationToken, CertificateToken).
func doesRevocationKnowCertificate(revocationToken AnyRevocationToken, certificateToken *model.CertificateToken) bool {
	return revocationInformationAssured(revocationToken, certificateToken) || certHashMatch(revocationToken)
}

// revocationInformationAssured is the port of the private static
// revocationInformationAssured(RevocationToken, CertificateToken).
func revocationInformationAssured(revocationToken AnyRevocationToken, certificateToken *model.CertificateToken) bool {
	notAfterRevoc := revocationToken.ThisUpdate()
	certNotAfter := certificateToken.NotAfter()

	if expiredCertsOnCRL := revocationToken.ExpiredCertsOnCRL(); !expiredCertsOnCRL.IsZero() {
		notAfterRevoc = expiredCertsOnCRL
	}
	if archiveCutOff := revocationToken.ArchiveCutOff(); !archiveCutOff.IsZero() {
		notAfterRevoc = archiveCutOff
	}

	return !certNotAfter.Before(notAfterRevoc)
}

// certHashMatch is the port of the private static certHashMatch(RevocationToken).
func certHashMatch(revocationToken AnyRevocationToken) bool {
	return revocationToken.CertHashPresent() && revocationToken.CertHashMatch()
}

// isAcceptableSignatureAlgorithm verifies validity of the used signature algorithm on
// revocation data creation against the specified cryptographic constraints. Port of
// isAcceptableSignatureAlgorithm(RevocationToken, CertificateToken).
func (v *RevocationDataVerifier) isAcceptableSignatureAlgorithm(revocationToken AnyRevocationToken, issuerCertificateToken *model.CertificateToken) bool {
	// Code supports validation against both, explicit signature algorithm definition and
	// a digest + encryption algorithm pairs
	signatureAlgorithm := revocationToken.SignatureAlgorithm()
	if signatureAlgorithm == "" {
		return false
	}
	if issuerCertificateToken == nil {
		return false
	}

	if utils.IsMapNotEmpty(v.acceptableSignatureAlgorithmKeyLength) {
		if signatureAlgorithmMinKeySize, ok := v.acceptableSignatureAlgorithmKeyLength[signatureAlgorithm]; ok {
			if v.isPublicKeySizeSupported(signatureAlgorithmMinKeySize, issuerCertificateToken) {
				return true
			}
		}
		// continue with Digest+Encryption algorithm pairs
	}

	if utils.IsCollectionEmpty(v.acceptableDigestAlgorithms) {
		return false
	} else if utils.IsMapEmpty(v.acceptableEncryptionAlgorithmKeyLength) {
		return false
	}

	if !containsDigestAlgorithm(v.acceptableDigestAlgorithms, signatureAlgorithm.DigestAlgorithm()) {
		return false
	}
	encryptionAlgorithmMinKeySize, ok := v.acceptableEncryptionAlgorithmKeyLength[signatureAlgorithm.EncryptionAlgorithm()]
	if !ok {
		return false
	}
	return v.isPublicKeySizeSupported(encryptionAlgorithmMinKeySize, issuerCertificateToken)
}

// containsDigestAlgorithm reports whether digestAlgorithms contains algorithm. Java relies on
// Collection#contains; Go has no generic Collection helper for a plain equality scan.
func containsDigestAlgorithm(digestAlgorithms []enumerations.DigestAlgorithm, algorithm enumerations.DigestAlgorithm) bool {
	for _, a := range digestAlgorithms {
		if a == algorithm {
			return true
		}
	}
	return false
}

// isPublicKeySizeSupported is the port of the private isPublicKeySizeSupported(Integer, CertificateToken).
func (v *RevocationDataVerifier) isPublicKeySizeSupported(minKeySize int, certificateToken *model.CertificateToken) bool {
	publicKeySize := spi.DSSPKUtilsPublicKeySize(certificateToken.PublicKey())
	if publicKeySize <= 0 {
		return false
	}
	if publicKeySize < minKeySize {
		return false
	}
	return true
}

// IsRevocationDataSkipNow checks and returns whether the revocation check shall be skipped for
// the given certificate at the current time. Port of the isRevocationDataSkip(CertificateToken)
// overload.
func (v *RevocationDataVerifier) IsRevocationDataSkipNow(certificateToken *model.CertificateToken) bool {
	return v.IsRevocationDataSkip(certificateToken, time.Now())
}

// IsRevocationDataSkip checks and returns whether the revocation check shall be skipped for the
// given certificate at controlTime. Port of the isRevocationDataSkip(CertificateToken, Date)
// overload; exported under this exact two-argument name because
// signature_validation_context.go's "RevocationDataVerifier (VAL-B)" header comment already
// depends on it: "IsRevocationDataSkip(cert *model.CertificateToken, controlTime time.Time) bool".
func (v *RevocationDataVerifier) IsRevocationDataSkip(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	if v.isTrustedAtTime(certificateToken, controlTime) {
		return true
	}
	if certificateToken.IsSelfSigned() {
		return true
	}
	if utils.IsCollectionEmpty(v.revocationSkipCertificateExtensions) {
		return false
	}
	certificateExtensions, err := spi.CertificateExtensionsUtilsCertificateExtensions(certificateToken)
	if err != nil {
		// Java's getCertificateExtensions call is unchecked; the Go port makes ASN.1 parsing
		// failures explicit, but this method has no error return to propagate them through, so
		// - like AbstractTimestampSource's `must` helper elsewhere in this phase - the error is
		// re-raised as the equivalent unchecked DSSException.
		panic(model.NewDSSErrorWithCause(err))
	}
	allCertificateExtensions := certificateExtensions.AllCertificateExtensions()
	if utils.IsCollectionNotEmpty(allCertificateExtensions) {
		oids := make([]string, 0, len(allCertificateExtensions))
		for _, ext := range allCertificateExtensions {
			oids = append(oids, ext.OID())
		}
		if utils.ContainsAny(oids, v.revocationSkipCertificateExtensions) {
			return true
		}
	}
	if utils.IsCollectionEmpty(v.revocationSkipCertificatePolicies) {
		return false
	}
	certificatePolicies := certificateExtensions.CertificatePolicies()
	if certificatePolicies != nil && utils.IsCollectionNotEmpty(certificatePolicies.PolicyList()) {
		oids := make([]string, 0, len(certificatePolicies.PolicyList()))
		for _, policy := range certificatePolicies.PolicyList() {
			oids = append(oids, policy.Oid())
		}
		if utils.ContainsAny(oids, v.revocationSkipCertificatePolicies) {
			return true
		}
	}
	return false
}

// isTrustedAtTime verifies whether the certificateToken is trusted at controlTime. Port of
// isTrustedAtTime(CertificateToken, Date).
func (v *RevocationDataVerifier) isTrustedAtTime(certificateToken *model.CertificateToken, controlTime time.Time) bool {
	currentTrustAnchorVerifier := v.TrustAnchorVerifier()
	if currentTrustAnchorVerifier == nil {
		return false
	}
	return currentTrustAnchorVerifier.IsTrustedAtTime(certificateToken, controlTime, enumerations.Context_REVOCATION)
}

// IsRevocationDataFresh verifies if the revocationToken considered within context is fresh
// enough relative to the given validationTime. Port of
// isRevocationDataFresh(RevocationToken, Date, Context).
func (v *RevocationDataVerifier) IsRevocationDataFresh(revocationToken AnyRevocationToken, validationTime time.Time, context enumerations.Context) bool {
	maximumRevocationFreshness := v.maximumRevocationFreshness(context)
	if maximumRevocationFreshness == nil {
		return v.isRevocationThisUpdateAfterValidationTimeNullConstraint(revocationToken, validationTime)
	}
	return isRevocationThisUpdateAfterValidationTime(revocationToken, validationTime, *maximumRevocationFreshness)
}

// isRevocationThisUpdateAfterValidationTime verifies whether the revocation's thisUpdate time
// is after validationTime minus the acceptable maximumRevocationFreshness. Port of
// isRevocationThisUpdateAfterValidationTime(RevocationToken, Date, long).
func isRevocationThisUpdateAfterValidationTime(revocationToken AnyRevocationToken, validationTime time.Time, maximumRevocationFreshness int64) bool {
	limit := validationTime.Add(-time.Duration(maximumRevocationFreshness) * time.Millisecond)
	thisUpdate := revocationToken.ThisUpdate()
	return !thisUpdate.IsZero() && thisUpdate.After(limit)
}

// isRevocationThisUpdateAfterValidationTimeNullConstraint verifies whether the revocation's
// thisUpdate time is after validationTime minus the difference between nextUpdate and
// thisUpdate field values. Port of
// isRevocationThisUpdateAfterValidationTimeNullConstraint(RevocationToken, Date).
func (v *RevocationDataVerifier) isRevocationThisUpdateAfterValidationTimeNullConstraint(revocationToken AnyRevocationToken, validationTime time.Time) bool {
	if !v.checkRevocationFreshnessNextUpdate {
		// no check to be performed
		return true
	}
	nextUpdate := revocationToken.NextUpdate()
	if nextUpdate.IsZero() {
		return false
	}
	limit := getDifference(nextUpdate, revocationToken.ThisUpdate())
	return isRevocationThisUpdateAfterValidationTime(revocationToken, validationTime, limit)
}

// getDifference is the port of the private getDifference(Date, Date), returning the difference
// in milliseconds.
func getDifference(nextUpdate, thisUpdate time.Time) int64 {
	var nextUpdateTime, thisUpdateTime int64
	if !nextUpdate.IsZero() {
		nextUpdateTime = nextUpdate.UnixMilli()
	}
	if !thisUpdate.IsZero() {
		thisUpdateTime = thisUpdate.UnixMilli()
	}
	return nextUpdateTime - thisUpdateTime
}

// maximumRevocationFreshness is the port of the private getMaximumRevocationFreshness(Context);
// Java's UnsupportedOperationException for an unhandled context becomes a panic.
func (v *RevocationDataVerifier) maximumRevocationFreshness(context enumerations.Context) *int64 {
	switch context {
	case enumerations.Context_SIGNATURE, enumerations.Context_COUNTER_SIGNATURE,
		enumerations.Context_KEY_BINDING_SIGNATURE, enumerations.Context_CERTIFICATE:
		return v.signatureMaximumRevocationFreshness
	case enumerations.Context_TIMESTAMP, enumerations.Context_EVIDENCE_RECORD:
		return v.timestampMaximumRevocationFreshness
	case enumerations.Context_REVOCATION:
		return v.revocationMaximumRevocationFreshness
	default:
		panic("The provided validation context '" + string(context) + "' is not supported!")
	}
}

// CheckCertificateNotRevoked verifies whether a certificate was not revoked at controlTime.
// Port of checkCertificateNotRevoked(RevocationToken, Date).
func (v *RevocationDataVerifier) CheckCertificateNotRevoked(revocationToken AnyRevocationToken, controlTime time.Time) bool {
	status := revocationToken.Status()
	return status.IsKnown() && (!status.IsRevoked() || controlTime.Before(revocationToken.RevocationDate()))
}

// IsAfterThisUpdateAndBeforeNextUpdate verifies whether the given date is within revocation
// data's thisUpdate and nextUpdate times. Port of
// isAfterThisUpdateAndBeforeNextUpdate(RevocationToken, Date).
func (v *RevocationDataVerifier) IsAfterThisUpdateAndBeforeNextUpdate(revocationToken AnyRevocationToken, date time.Time) bool {
	thisUpdate := revocationToken.ThisUpdate()
	nextUpdate := revocationToken.NextUpdate()
	return !thisUpdate.IsZero() && !date.Before(thisUpdate) && (nextUpdate.IsZero() || !date.After(nextUpdate))
}

// IsCertificateChainValid verifies whether the certificate chain is valid at controlTime. Port
// of isCertificateChainValid(List, Date, Context).
func (v *RevocationDataVerifier) IsCertificateChainValid(certificateTokenChain []*model.CertificateToken, controlTime time.Time, context enumerations.Context) bool {
	if v.isAcceptCertificatesWithoutRevocation(context) {
		return true
	}
	for _, certificateToken := range certificateTokenChain {
		if certificateToken.IsSelfSigned() || v.isTrustedAtTime(certificateToken, controlTime) {
			break
		}
		if !certificateToken.IsValid() {
			return false
		}
		if !v.isCertificateValid(certificateToken, certificateTokenChain, controlTime) {
			return false
		}
	}
	return true
}

// isAcceptCertificatesWithoutRevocation is the port of the private isAcceptCertificatesWithoutRevocation(Context).
func (v *RevocationDataVerifier) isAcceptCertificatesWithoutRevocation(context enumerations.Context) bool {
	return (enumerations.Context_TIMESTAMP == context && v.acceptTimestampCertificatesWithoutRevocation) ||
		(enumerations.Context_REVOCATION == context && v.acceptRevocationCertificatesWithoutRevocation)
}

// isCertificateValid verifies if the certificate is valid. Port of
// isCertificateValid(CertificateToken, Collection, Date).
func (v *RevocationDataVerifier) isCertificateValid(certificateToken *model.CertificateToken, certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	if !v.IsRevocationDataSkip(certificateToken, controlTime) {
		if !hasRevocationAccessPoints(certificateToken) {
			return false
		}
		if !v.isCertificateNotRevoked(certificateToken, certificateChain, controlTime) {
			return false
		}
	}
	return true
}

// hasRevocationAccessPoints is the port of the private static hasRevocationAccessPoints(CertificateToken).
func hasRevocationAccessPoints(certificateToken *model.CertificateToken) bool {
	return utils.IsCollectionNotEmpty(spi.CertificateExtensionsUtilsCRLAccessUrls(certificateToken)) ||
		utils.IsCollectionNotEmpty(spi.CertificateExtensionsUtilsOCSPAccessUrls(certificateToken))
}

// isCertificateNotRevoked verifies whether a certificate token is not revoked at controlTime.
// Port of isCertificateNotRevoked(CertificateToken, Collection, Date).
func (v *RevocationDataVerifier) isCertificateNotRevoked(certificateToken *model.CertificateToken, certificateChain []*model.CertificateToken, controlTime time.Time) bool {
	v.populateValidationContext(certificateChain)

	revocationData := v.relatedRevocationTokens(certificateToken)
	for _, revocationToken := range revocationData {
		if !v.isSelfIssuedRevocation(certificateToken, revocationToken) && v.isAcceptableWithIssuer(revocationToken, revocationToken.IssuerCertificateToken(), controlTime) &&
			v.CheckCertificateNotRevoked(revocationToken, controlTime) {
			return true
		}
	}
	return false
}

// isSelfIssuedRevocation verifies whether the verified certificate occurs in the revocation's
// issuer certificate chain. Port of isSelfIssuedRevocation(CertificateToken, RevocationToken).
func (v *RevocationDataVerifier) isSelfIssuedRevocation(certificateToken *model.CertificateToken, revocationData AnyRevocationToken) bool {
	if certificateToken.Equals(revocationData.IssuerCertificateToken()) {
		return true
	}

	// verify certificate occurrence within the chain
	if utils.IsCollectionNotEmpty(revocationData.Certificates()) {
		certificateChain, err := spi.NewCertificateReordererWithSigningCertificate(
			revocationData.IssuerCertificateToken(), revocationData.Certificates()).OrderedCertificates()
		if err == nil {
			for _, c := range certificateChain {
				if c.Equals(certificateToken) {
					return true
				}
			}
		}
		// upstream logs the reorderer failure here, dropped per PORTING.md's slf4j rule; an
		// unresolvable chain is treated the same as "certificate not found in the chain".
	}
	return false
}

// populateValidationContext is the port of the private populateValidationContext(Collection).
func (v *RevocationDataVerifier) populateValidationContext(certificateChain []*model.CertificateToken) {
	if v.validationContext != nil {
		for _, certificateToken := range certificateChain {
			v.validationContext.AddCertificateTokenForVerification(certificateToken)
		}
	}
}

// relatedRevocationTokens is the port of the private getRelatedRevocationTokens(CertificateToken).
func (v *RevocationDataVerifier) relatedRevocationTokens(certificateToken *model.CertificateToken) []AnyRevocationToken {
	if v.validationContext == nil {
		return nil
	}
	revocationData := v.validationContext.GetRevocationData(certificateToken)
	if utils.IsCollectionEmpty(revocationData) {
		return nil
	}
	filtered := make([]AnyRevocationToken, 0, len(revocationData))
	for _, r := range revocationData {
		if utils.AreStringsEqual(certificateToken.DSSIDAsString(), r.RelatedCertificateID()) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
