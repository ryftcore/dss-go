// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/ValidationContext.java (DSS 6.5.RC1).
//
// This interface is implemented by SignatureValidationContext; every method below is named to
// match SignatureValidationContext's actual exported method set exactly (confirmed by reading
// signature_validation_context.go directly, not only its header-comment) so that
// *SignatureValidationContext continues to satisfy this interface, which is required for
// signature_validation_context.go:264's `c.revocationDataVerifier.SetValidationContext(c)` and
// baseline_requirements_checker.go's `b.validationContext = NewSignatureValidationContext()`
// (assigned to a field/return typed Context) to compile.
//
// Per the Java 6.2+ behavior change documented in the javadoc of every checkXxx() method below
// ("returning only the boolean validation result, without alerts handling"), these methods
// return bool - not TokenStatus/SignatureStatus - matching SignatureValidationContext's
// CheckXxx() bool methods exactly.
//
// getValidationData has two Java overloads (AdvancedSignature, TimestampToken); Go has no
// overloading, so SignatureValidationContext exposes GetValidationData (signature overload) and
// GetValidationDataForTimestamp (timestamp overload) - both are included here under those exact
// names.
//
// getRevocationData(CertificateToken) is exposed as GetRevocationData (get-prefix kept, unlike
// most of this port's dropped getters, for consistency with SignatureValidationContext's other
// exported getters: GetCurrentTime, GetProcessedSignatures, GetAllCertificateSources,
// GetDocumentCertificateSource, GetValidationData, ...).
package validation

import (
	"time"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi"
)

// Context allows the implementation of the validators for: certificates, timestamps
// and revocation data.
type Context interface {
	// Initialize initializes the Context by retrieving the relevant data from
	// certificateVerifier. Port of initialize(CertificateVerifier).
	Initialize(certificateVerifier CertificateVerifier)

	// GetCurrentTime gets the current validation time. Port of getCurrentTime().
	GetCurrentTime() time.Time

	// AddSignatureForVerification adds a new signature to collect the information to verify.
	// Port of addSignatureForVerification(AdvancedSignature).
	AddSignatureForVerification(signature AdvancedSignature)

	// AddRevocationTokenForVerification adds a new revocation token to the list of tokens to
	// verify. If the revocation token has already been added then it is ignored. Port of
	// addRevocationTokenForVerification(RevocationToken).
	AddRevocationTokenForVerification(revocationToken AnyRevocationToken)

	// AddCertificateTokenForVerification adds a new certificate token to the list of tokens to
	// verify. If the certificate token has already been added then it is ignored. Port of
	// addCertificateTokenForVerification(CertificateToken).
	AddCertificateTokenForVerification(certificateToken *model.CertificateToken)

	// AddTimestampTokenForVerification adds a new timestamp token to the list of tokens to
	// verify. If the timestamp token has already been added then it is ignored. Port of
	// addTimestampTokenForVerification(TimestampToken).
	AddTimestampTokenForVerification(timestampToken *TimestampToken)

	// AddEvidenceRecordForVerification adds Evidence Record's content to proceed with
	// validation. Port of addEvidenceRecordForVerification(EvidenceRecord).
	AddEvidenceRecordForVerification(evidenceRecord EvidenceRecord)

	// AddDocumentCertificateSource adds an extracted certificate source to the used list of
	// sources. Port of addDocumentCertificateSource(CertificateSource).
	AddDocumentCertificateSource(certificateSource spi.CertificateSource)

	// AddDocumentCertificateSourceFromList adds a list certificate source to the used list of
	// sources. Port of the addDocumentCertificateSource(ListCertificateSource) overload.
	AddDocumentCertificateSourceFromList(listCertificateSource *spi.ListCertificateSource)

	// AddDocumentCRLSource adds an extracted CRL source to the used list of sources. Port of
	// addDocumentCRLSource(OfflineRevocationSource).
	AddDocumentCRLSource(crlSource spi.OfflineRevocationSource[revocation.CRL])

	// AddDocumentCRLSourceFromList adds a list CRL source to the used list of sources. Port of
	// the addDocumentCRLSource(ListRevocationSource) overload.
	AddDocumentCRLSourceFromList(crlSource *spi.ListRevocationSource[revocation.CRL])

	// AddDocumentOCSPSource adds an extracted OCSP source to the used list of sources. Port of
	// addDocumentOCSPSource(OfflineRevocationSource).
	AddDocumentOCSPSource(ocspSource spi.OfflineRevocationSource[revocation.OCSP])

	// AddDocumentOCSPSourceFromList adds a list OCSP source to the used list of sources. Port
	// of the addDocumentOCSPSource(ListRevocationSource) overload.
	AddDocumentOCSPSourceFromList(ocspSource *spi.ListRevocationSource[revocation.OCSP])

	// Validate carries out the validation process in recursive manner for not yet checked
	// tokens. Port of validate().
	Validate()

	// CheckAllRequiredRevocationDataPresent returns whether all processed certificates have a
	// revocation data. Port of checkAllRequiredRevocationDataPresent().
	CheckAllRequiredRevocationDataPresent() bool

	// CheckAllPOECoveredByRevocationData returns whether all POE (timestamp tokens) are covered
	// by a revocation data. Port of checkAllPOECoveredByRevocationData().
	CheckAllPOECoveredByRevocationData() bool

	// CheckAllTimestampsValid returns whether all processed timestamps are valid and intact.
	// Port of checkAllTimestampsValid().
	CheckAllTimestampsValid() bool

	// CheckCertificateNotRevoked returns whether the certificate is not revoked. Port of
	// checkCertificateNotRevoked(CertificateToken).
	CheckCertificateNotRevoked(certificateToken *model.CertificateToken) bool

	// CheckAllSignatureCertificatesNotRevoked returns whether none of the signature's
	// certificate chain certificates are not revoked, validating recursively. Port of
	// checkAllSignatureCertificatesNotRevoked().
	CheckAllSignatureCertificatesNotRevoked() bool

	// CheckAllSignatureCertificateHaveFreshRevocationData returns whether for all signature's
	// certificate chain certificates there is a fresh revocation data, after the earliest
	// available timestamp token production time. Port of
	// checkAllSignatureCertificateHaveFreshRevocationData().
	CheckAllSignatureCertificateHaveFreshRevocationData() bool

	// CheckAllSignaturesNotExpired returns whether all signatures added to the Context
	// are not yet expired. Port of checkAllSignaturesNotExpired().
	CheckAllSignaturesNotExpired() bool

	// CheckCertificateNotExpired returns whether the certificate token is not yet expired. Port
	// of checkCertificateNotExpired(CertificateToken).
	CheckCertificateNotExpired(certificateToken *model.CertificateToken) bool

	// CheckAllSignaturesAreYetValid returns whether all signatures added to the
	// Context have been produced with yet valid certificates at the time of signing.
	// Port of checkAllSignaturesAreYetValid().
	CheckAllSignaturesAreYetValid() bool

	// CheckCertificateIsYetValid returns whether the certificate token is yet valid. Port of
	// checkCertificateIsYetValid(CertificateToken).
	CheckCertificateIsYetValid(certificateToken *model.CertificateToken) bool

	// GetProcessedSignatures returns signatures added to the validation context. Port of
	// getProcessedSignatures().
	GetProcessedSignatures() []AdvancedSignature

	// GetProcessedCertificates returns a read only list of all certificates used in the process
	// of the validation of all signatures from the given document. This list includes the
	// certificate to check, certification chain certificates, OCSP response certificate...
	// Port of getProcessedCertificates().
	GetProcessedCertificates() []*model.CertificateToken

	// GetProcessedRevocations returns a read only list of all revocations used in the process of
	// the validation of all signatures from the given document. Port of
	// getProcessedRevocations().
	GetProcessedRevocations() []AnyRevocationToken

	// GetProcessedTimestamps returns a read only list of all timestamps processed during the
	// validation of all signatures from the given document. Port of getProcessedTimestamps().
	GetProcessedTimestamps() []*TimestampToken

	// GetProcessedEvidenceRecords returns evidence records added to the validation context.
	// Port of getProcessedEvidenceRecords().
	GetProcessedEvidenceRecords() []EvidenceRecord

	// GetAllCertificateSources returns a list of all CertificateSources used during the
	// validation process. It is represented by sources extracted from the provided document
	// (e.g. signatures, timestamps) as well as the sources obtained during the validation
	// process (e.g. AIA, OCSP). Port of getAllCertificateSources().
	GetAllCertificateSources() *spi.ListCertificateSource

	// GetDocumentCertificateSource returns a list of all CertificateSources extracted from a
	// validating document (signature(s), timestamp(s)). Port of
	// getDocumentCertificateSource().
	GetDocumentCertificateSource() *spi.ListCertificateSource

	// GetDocumentCRLSource returns a list of all CRL OfflineRevocationSources extracted from a
	// validating document. Port of getDocumentCRLSource().
	GetDocumentCRLSource() *spi.ListRevocationSource[revocation.CRL]

	// GetDocumentOCSPSource returns a list of all OCSP OfflineRevocationSources extracted from a
	// validating document. Port of getDocumentOCSPSource().
	GetDocumentOCSPSource() *spi.ListRevocationSource[revocation.OCSP]

	// GetValidationData returns a validation data for the given signature's certificate chain.
	// Port of getValidationData(AdvancedSignature).
	GetValidationData(signature AdvancedSignature) *Data

	// GetValidationDataForTimestamp returns a validation data for the given timestampToken's
	// certificate chain. Port of the getValidationData(TimestampToken) overload.
	GetValidationDataForTimestamp(timestampToken *TimestampToken) *Data

	// GetRevocationData returns revocation data for the given certificateToken, whether
	// extracted from a signature file or obtained online. Port of
	// getRevocationData(CertificateToken).
	GetRevocationData(certificateToken *model.CertificateToken) []AnyRevocationToken
}

// compile-time assertion: *SignatureValidationContext satisfies Context.
var _ Context = (*SignatureValidationContext)(nil)
