// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/DiagnosticData.java (DSS 6.5.RC1).
//
// Java's getAllSignatures()/getAllCounterSignatures()/getAllKeyBindingSignatures()/
// getAllRevocationData()/getAllEAA()/getAllEAARevocationTokens() return java.util.Set<T>.
// Set<T> normally becomes map[T]struct{}, but every element here is already a
// unique wrapper instance drawn from this type's own cached slices (getSignatures(),
// getAllRevocationData() iterates wrapped.getUsedRevocations() once), so a map would only add
// nondeterministic iteration order without changing membership semantics. This port therefore
// returns a deterministic, insertion-ordered []*T instead; documented here as a deliberate
// deviation from the literal Set<T> mapping rule.
package diagnostic

import (
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// DiagnosticData represents all static data extracted by the process analysing the signature.
// They are independent from the validation policy to be applied.
type DiagnosticData struct {
	// wrapped is the wrapped XmlDiagnosticData jaxb object.
	wrapped *jaxb.XmlDiagnosticData

	// foundSignatures caches the list of found signatures.
	foundSignatures []*SignatureWrapper
	// usedCertificates caches the list of used certificates.
	usedCertificates []*CertificateWrapper
	// usedTimestamps caches the list of found timestamps.
	usedTimestamps []*TimestampWrapper
	// foundEvidenceRecords caches the list of found evidence records.
	foundEvidenceRecords []*EvidenceRecordWrapper
	// foundEAAs caches the list of found EAA presentations.
	foundEAAs []*EAAWrapper
}

// NewDiagnosticData is the default constructor.
func NewDiagnosticData(wrapped *jaxb.XmlDiagnosticData) *DiagnosticData {
	return &DiagnosticData{wrapped: wrapped}
}

// DocumentName returns a name of the validating document. Port of getDocumentName().
func (d *DiagnosticData) DocumentName() string {
	if d.wrapped.DocumentName != nil {
		return *d.wrapped.DocumentName
	}
	return ""
}

// SignatureIdList returns the list of the signature id. Port of getSignatureIdList().
func (d *DiagnosticData) SignatureIdList() []string {
	var signatureIds []string
	signatures := d.wrapped.Signatures.All()
	if signatures != nil {
		for _, xmlSignature := range signatures {
			if xmlSignature.Id != nil {
				signatureIds = append(signatureIds, string(*xmlSignature.Id))
			}
		}
	}
	return signatureIds
}

// FirstSignatureId returns the first signature id. Port of getFirstSignatureId().
func (d *DiagnosticData) FirstSignatureId() string {
	return d.firstSignatureNullSafe().Id()
}

// FirstSignatureDate returns the first signature time. Port of getFirstSignatureDate().
func (d *DiagnosticData) FirstSignatureDate() *time.Time {
	return d.firstSignatureNullSafe().ClaimedSigningTime()
}

// SignatureDate returns the claimed signing time. Port of getSignatureDate(String).
func (d *DiagnosticData) SignatureDate(signatureId string) *time.Time {
	return d.signatureByIdNullSafe(signatureId).ClaimedSigningTime()
}

// FirstSignatureFormat returns the signature format for the first signature. Port of
// getFirstSignatureFormat().
func (d *DiagnosticData) FirstSignatureFormat() enumerations.SignatureLevel {
	return d.firstSignatureNullSafe().SignatureFormat()
}

// SignatureFormat returns the signature format for the given signature. Port of
// getSignatureFormat(String).
func (d *DiagnosticData) SignatureFormat(signatureId string) enumerations.SignatureLevel {
	return d.signatureByIdNullSafe(signatureId).SignatureFormat()
}

// SignedAssertionsInFirstSignature returns the signed assertions for the first signature.
// Port of getSignedAssertionsInFirstSignature().
func (d *DiagnosticData) SignedAssertionsInFirstSignature() []*jaxb.XmlSignerRole {
	return d.firstSignatureNullSafe().SignedAssertions()
}

// SignedAssertions returns the signed assertions for the given signature. Port of
// getSignedAssertions(String).
func (d *DiagnosticData) SignedAssertions(signatureId string) []*jaxb.XmlSignerRole {
	return d.signatureByIdNullSafe(signatureId).SignedAssertions()
}

// FirstSignatureDigestAlgorithm returns the DigestAlgorithm of the first signature. Port of
// getFirstSignatureDigestAlgorithm().
func (d *DiagnosticData) FirstSignatureDigestAlgorithm() enumerations.DigestAlgorithm {
	return d.firstSignatureNullSafe().DigestAlgorithm()
}

// SignatureDigestAlgorithm returns the DigestAlgorithm for the given signature. Port of
// getSignatureDigestAlgorithm(String).
func (d *DiagnosticData) SignatureDigestAlgorithm(signatureId string) enumerations.DigestAlgorithm {
	return d.signatureByIdNullSafe(signatureId).DigestAlgorithm()
}

// FirstSignatureEncryptionAlgorithm returns the EncryptionAlgorithm of the first signature.
// Port of getFirstSignatureEncryptionAlgorithm().
func (d *DiagnosticData) FirstSignatureEncryptionAlgorithm() enumerations.EncryptionAlgorithm {
	return d.firstSignatureNullSafe().EncryptionAlgorithm()
}

// SignatureEncryptionAlgorithm returns the EncryptionAlgorithm for the given signature. Port
// of getSignatureEncryptionAlgorithm(String).
func (d *DiagnosticData) SignatureEncryptionAlgorithm(signatureId string) enumerations.EncryptionAlgorithm {
	return d.signatureByIdNullSafe(signatureId).EncryptionAlgorithm()
}

// SigningCertificateId returns signing certificate dss id for the given signature. Port of
// getSigningCertificateId(String).
func (d *DiagnosticData) SigningCertificateId(signatureId string) string {
	signature := d.signatureByIdNullSafe(signatureId)
	if signature.SigningCertificate() != nil {
		return signature.SigningCertificate().Id()
	}
	return ""
}

// IsSigningCertificateIdentified indicates if the digest value and the issuer and serial
// match for the signing certificate. Port of isSigningCertificateIdentified(String).
func (d *DiagnosticData) IsSigningCertificateIdentified(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsSigningCertificateIdentified()
}

// SignatureCertificateChain returns the list of certificates in the chain of the main
// signature. Port of getSignatureCertificateChain(String).
func (d *DiagnosticData) SignatureCertificateChain(signatureId string) []*CertificateWrapper {
	return d.signatureByIdNullSafe(signatureId).CertificateChain()
}

// SignatureCertificateChainIds returns the list of certificate identifiers in the chain of the
// main signature. Port of getSignatureCertificateChainIds(String).
func (d *DiagnosticData) SignatureCertificateChainIds(signatureId string) []string {
	signature := d.signatureByIdNullSafe(signatureId)
	var result []string
	for _, certWrapper := range signature.CertificateChain() {
		result = append(result, certWrapper.Id())
	}
	return result
}

// FirstPolicyId returns the identifier of the policy of the first signature. Port of
// getFirstPolicyId().
func (d *DiagnosticData) FirstPolicyId() string {
	return d.firstSignatureNullSafe().PolicyId()
}

// PolicyId returns the identifier of the policy. Port of getPolicyId(String).
func (d *DiagnosticData) PolicyId(signatureId string) string {
	return d.signatureByIdNullSafe(signatureId).PolicyId()
}

// PolicyDescription returns the description of the policy. Port of
// getPolicyDescription(String).
func (d *DiagnosticData) PolicyDescription(signatureId string) string {
	return d.signatureByIdNullSafe(signatureId).PolicyDescription()
}

// PolicyDocumentationReferences returns the documentation references of the policy. Port of
// getPolicyDocumentationReferences(String).
func (d *DiagnosticData) PolicyDocumentationReferences(signatureId string) []string {
	return d.signatureByIdNullSafe(signatureId).PolicyDocumentationReferences()
}

// TimestampIdList returns the list of identifier of all timestamps found during the
// validation. Port of getTimestampIdList().
func (d *DiagnosticData) TimestampIdList() []string {
	var timestampIdList []string
	for _, timestampWrapper := range d.TimestampList() {
		timestampIdList = append(timestampIdList, timestampWrapper.Id())
	}
	return timestampIdList
}

// TimestampIdListForSignature returns the list of identifier of the timestamps related to the
// given signature. Port of getTimestampIdList(String).
func (d *DiagnosticData) TimestampIdListForSignature(signatureId string) []string {
	return d.signatureByIdNullSafe(signatureId).TimestampIdsList()
}

// TimestampListForSignature returns the list of timestamps wrappers which cover the given
// signature. Port of getTimestampList(String).
func (d *DiagnosticData) TimestampListForSignature(signatureId string) []*TimestampWrapper {
	return d.signatureByIdNullSafe(signatureId).TimestampList()
}

// IsBLevelTechnicallyValid indicates if the -B level is technically valid. It means that the
// signature value is valid. Port of isBLevelTechnicallyValid(String).
func (d *DiagnosticData) IsBLevelTechnicallyValid(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsBLevelTechnicallyValid()
}

// IsThereTLevel indicates if there is a signature timestamp. Port of isThereTLevel(String).
func (d *DiagnosticData) IsThereTLevel(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsThereTLevel()
}

// IsTLevelTechnicallyValid indicates if the -T level is technically valid. It means that the
// signature and the digest are valid. Port of isTLevelTechnicallyValid(String).
func (d *DiagnosticData) IsTLevelTechnicallyValid(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsTLevelTechnicallyValid()
}

// IsThereXLevel indicates if there is an -X1 or -X2 timestamp. Port of isThereXLevel(String).
func (d *DiagnosticData) IsThereXLevel(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsThereXLevel()
}

// IsXLevelTechnicallyValid indicates if the -X level is technically valid. It means that the
// signature and the digest are valid. Port of isXLevelTechnicallyValid(String).
func (d *DiagnosticData) IsXLevelTechnicallyValid(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsXLevelTechnicallyValid()
}

// IsThereALevel indicates if there is an archive timestamp. Port of isThereALevel(String).
func (d *DiagnosticData) IsThereALevel(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsThereALevel()
}

// IsALevelTechnicallyValid indicates if the -A (-LTA) level is technically valid. It means
// that the signature of the archive timestamps are valid and their imprint is valid too. Port
// of isALevelTechnicallyValid(String).
func (d *DiagnosticData) IsALevelTechnicallyValid(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsALevelTechnicallyValid()
}

// IsThereERSLevel indicates if there is an embedded evidence record. Port of
// isThereERSLevel(String).
func (d *DiagnosticData) IsThereERSLevel(signatureId string) bool {
	return d.signatureByIdNullSafe(signatureId).IsThereERSLevel()
}

// SignerDocuments returns a list of all Signer's documents used to create a signature. NOTE:
// returns a first level documents only (e.g. a signed Manifest for XAdES, when applicable).
// Port of getSignerDocuments(String).
func (d *DiagnosticData) SignerDocuments(signatureId string) []*SignerDataWrapper {
	var result []*SignerDataWrapper
	signatureWrapper := d.signatureByIdNullSafe(signatureId)
	signatureScopes := signatureWrapper.SignatureScopes()
	for _, xmlSignatureScope := range signatureScopes {
		signerData := xmlSignatureScope.SignerData
		// return first level data only
		if signerData.Parent == nil {
			result = append(result, NewSignerDataWrapper(signerData))
		}
	}
	return result
}

// TimestampSigningCertificateId returns the identifier of the timestamp signing certificate.
// Port of getTimestampSigningCertificateId(String).
func (d *DiagnosticData) TimestampSigningCertificateId(timestampId string) string {
	return d.timestampByIdNullSafe(timestampId).SigningCertificate().Id()
}

// TimestampType returns the timestamp type of the given timestamp. Port of
// getTimestampType(String).
func (d *DiagnosticData) TimestampType(timestampId string) enumerations.TimestampType {
	return d.timestampByIdNullSafe(timestampId).Type()
}

// TimestampsByType returns a list of TimestampWrapper for the given TimestampType. Port of
// getTimestampsByType(TimestampType).
func (d *DiagnosticData) TimestampsByType(timestampType enumerations.TimestampType) []*TimestampWrapper {
	var result []*TimestampWrapper
	for _, timestampWrapper := range d.TimestampList() {
		if timestampType != "" && timestampType == timestampWrapper.Type() {
			result = append(result, timestampWrapper)
		}
	}
	return result
}

// IsValidCertificate indicates if the certificate signature is valid and the revocation status
// is valid. Port of isValidCertificate(String).
func (d *DiagnosticData) IsValidCertificate(dssCertificateId string) bool {
	certificate := d.UsedCertificateByIdNullSafe(dssCertificateId)

	signatureValid := certificate.IsSignatureValid()
	latestRevocationData := d.LatestRevocationDataForCertificate(certificate)
	revocationValid := latestRevocationData != nil && latestRevocationData.Status().IsGood()
	trusted := certificate.IsTrusted()
	return signatureValid && (trusted || revocationValid)
}

// CertificateDN returns the subject distinguished name for the given dss certificate
// identifier. Port of getCertificateDN(String).
func (d *DiagnosticData) CertificateDN(dssCertificateId string) string {
	return d.UsedCertificateByIdNullSafe(dssCertificateId).CertificateDN()
}

// CertificateIssuerDN returns the issuer distinguished name for the given dss certificate
// identifier. Port of getCertificateIssuerDN(String).
func (d *DiagnosticData) CertificateIssuerDN(dssCertificateId string) string {
	return d.UsedCertificateByIdNullSafe(dssCertificateId).CertificateIssuerDN()
}

// CertificateSerialNumber returns the serial number of the given dss certificate identifier.
// Port of getCertificateSerialNumber(String).
func (d *DiagnosticData) CertificateSerialNumber(dssCertificateId string) string {
	return d.UsedCertificateByIdNullSafe(dssCertificateId).SerialNumber()
}

// CertificateRevocationSource returns the revocation source for the given certificate. Port
// of getCertificateRevocationSource(String).
func (d *DiagnosticData) CertificateRevocationSource(dssCertificateId string) enumerations.RevocationType {
	certificate := d.UsedCertificateByIdNullSafe(dssCertificateId)
	if certificate.IsRevocationDataAvailable() {
		return d.LatestRevocationDataForCertificate(certificate).RevocationType()
	}
	return ""
}

// CertificateRevocationStatus returns the revocation status for the given certificate. Port
// of getCertificateRevocationStatus(String).
func (d *DiagnosticData) CertificateRevocationStatus(dssCertificateId string) enumerations.CertificateStatus {
	certificate := d.UsedCertificateByIdNullSafe(dssCertificateId)
	if certificate.IsRevocationDataAvailable() {
		return d.LatestRevocationDataForCertificate(certificate).Status()
	}
	return enumerations.CertificateStatusUnknown
}

// CertificateRevocationReason returns the revocation reason for the given certificate. Port
// of getCertificateRevocationReason(String).
func (d *DiagnosticData) CertificateRevocationReason(dssCertificateId string) enumerations.RevocationReason {
	certificate := d.UsedCertificateByIdNullSafe(dssCertificateId)
	if certificate.IsRevocationDataAvailable() {
		return d.LatestRevocationDataForCertificate(certificate).Reason()
	}
	return ""
}

// ErrorMessage retrieves the error message for the given signature id. Port of
// getErrorMessage(String).
func (d *DiagnosticData) ErrorMessage(signatureId string) string {
	return d.signatureByIdNullSafe(signatureId).ErrorMessage()
}

func (d *DiagnosticData) firstSignatureNullSafe() *SignatureWrapper {
	signatures := d.Signatures()
	if len(signatures) != 0 {
		return signatures[0]
	}
	return NewSignatureWrapper(&jaxb.XmlSignature{}) // TODO improve ?
}

// SignatureById returns a signature wrapper for the given signature id. Port of
// getSignatureById(String).
func (d *DiagnosticData) SignatureById(id string) *SignatureWrapper {
	for _, xmlSignature := range d.Signatures() {
		if id == xmlSignature.Id() {
			return xmlSignature
		}
	}
	return nil
}

func (d *DiagnosticData) signatureByIdNullSafe(id string) *SignatureWrapper {
	for _, xmlSignature := range d.Signatures() {
		if id == xmlSignature.Id() {
			return xmlSignature
		}
	}
	return NewSignatureWrapper(&jaxb.XmlSignature{}) // TODO improve ?
}

func (d *DiagnosticData) timestampByIdNullSafe(id string) *TimestampWrapper {
	timestamp := d.TimestampById(id)
	if timestamp != nil {
		return timestamp
	}
	return NewTimestampWrapper(&jaxb.XmlTimestamp{})
}

// TimestampById returns the TimestampWrapper corresponding to the given id. Port of
// getTimestampById(String).
func (d *DiagnosticData) TimestampById(id string) *TimestampWrapper {
	for _, timestampWrapper := range d.TimestampList() {
		if id == timestampWrapper.Id() {
			return timestampWrapper
		}
	}
	return nil
}

// UsedCertificateByIdNullSafe returns a certificate wrapper for the given certificate id.
// Port of getUsedCertificateByIdNullSafe(String).
func (d *DiagnosticData) UsedCertificateByIdNullSafe(id string) *CertificateWrapper {
	cert := d.UsedCertificateById(id)
	if cert != nil {
		return cert
	}
	return NewCertificateWrapper(&jaxb.XmlCertificate{}) // TODO improve ?
}

// UsedCertificateById returns a certificate wrapper for the given certificate id. Port of
// getUsedCertificateById(String).
func (d *DiagnosticData) UsedCertificateById(id string) *CertificateWrapper {
	for _, certificate := range d.UsedCertificates() {
		if id == certificate.Id() {
			return certificate
		}
	}
	return nil
}

// OrphanCertificateById returns an orphan certificate wrapper for the given certificate id.
// Port of getOrphanCertificateById(String).
func (d *DiagnosticData) OrphanCertificateById(id string) *OrphanCertificateTokenWrapper {
	for _, certificate := range d.AllOrphanCertificateObjects() {
		if id == certificate.Id() {
			return certificate
		}
	}
	return nil
}

// CertificatesFromSource returns a list of certificates by their origin source. Port of
// getCertificatesFromSource(CertificateSourceType).
func (d *DiagnosticData) CertificatesFromSource(certificateSourceType enumerations.CertificateSourceType) []*CertificateWrapper {
	var certificates []*CertificateWrapper
	for _, certificate := range d.UsedCertificates() {
		for _, source := range certificate.Sources() {
			if source == certificateSourceType {
				certificates = append(certificates, certificate)
				break
			}
		}
	}
	return certificates
}

// AllOrphanCertificateObjects returns a list of all found OrphanCertificateWrapper values.
// Port of getAllOrphanCertificateObjects().
func (d *DiagnosticData) AllOrphanCertificateObjects() []*OrphanCertificateTokenWrapper {
	var orphanCertificateValues []*OrphanCertificateTokenWrapper
	if d.wrapped.OrphanTokens != nil {
		for _, orphanToken := range d.wrapped.OrphanTokens.OrphanCertificate {
			orphanCertificate := NewOrphanCertificateTokenWrapper(orphanToken)
			if orphanToken.EncapsulationType != nil && *orphanToken.EncapsulationType == jaxb.XmlEncapsulationTypeBinaries &&
				!containsOrphanCertificate(orphanCertificateValues, orphanCertificate) {
				orphanCertificateValues = append(orphanCertificateValues, orphanCertificate)
			}
		}
	}
	return orphanCertificateValues
}

// AllOrphanCertificateReferences returns a list of all found orphan certificate references.
// Port of getAllOrphanCertificateReferences().
func (d *DiagnosticData) AllOrphanCertificateReferences() []*OrphanCertificateTokenWrapper {
	var orphanCertificateRefs []*OrphanCertificateTokenWrapper
	if d.wrapped.OrphanTokens != nil {
		for _, orphanToken := range d.wrapped.OrphanTokens.OrphanCertificate {
			orphanCertificate := NewOrphanCertificateTokenWrapper(orphanToken)
			if orphanToken.EncapsulationType != nil && *orphanToken.EncapsulationType == jaxb.XmlEncapsulationTypeReference &&
				!containsOrphanCertificate(orphanCertificateRefs, orphanCertificate) {
				orphanCertificateRefs = append(orphanCertificateRefs, orphanCertificate)
			}
		}
	}
	return orphanCertificateRefs
}

func containsOrphanCertificate(values []*OrphanCertificateTokenWrapper, candidate *OrphanCertificateTokenWrapper) bool {
	for _, v := range values {
		if v.Id() == candidate.Id() {
			return true
		}
	}
	return false
}

// AllOrphanRevocationObjects returns a list of all found OrphanRevocationWrapper values. Port
// of getAllOrphanRevocationObjects().
func (d *DiagnosticData) AllOrphanRevocationObjects() []*OrphanRevocationTokenWrapper {
	var orphanRevocationValues []*OrphanRevocationTokenWrapper
	if d.wrapped.OrphanTokens != nil {
		for _, orphanToken := range d.wrapped.OrphanTokens.OrphanRevocation {
			orphanRevocation := NewOrphanRevocationTokenWrapper(orphanToken)
			if orphanToken.EncapsulationType != nil && *orphanToken.EncapsulationType == jaxb.XmlEncapsulationTypeBinaries &&
				!containsOrphanRevocation(orphanRevocationValues, orphanRevocation) {
				orphanRevocationValues = append(orphanRevocationValues, orphanRevocation)
			}
		}
	}
	return orphanRevocationValues
}

// AllOrphanRevocationReferences returns a list of all found orphan revocation references. Port
// of getAllOrphanRevocationReferences().
func (d *DiagnosticData) AllOrphanRevocationReferences() []*OrphanRevocationTokenWrapper {
	var orphanRevocationRefs []*OrphanRevocationTokenWrapper
	if d.wrapped.OrphanTokens != nil {
		for _, orphanToken := range d.wrapped.OrphanTokens.OrphanRevocation {
			orphanRevocation := NewOrphanRevocationTokenWrapper(orphanToken)
			if orphanToken.EncapsulationType != nil && *orphanToken.EncapsulationType == jaxb.XmlEncapsulationTypeReference &&
				!containsOrphanRevocation(orphanRevocationRefs, orphanRevocation) {
				orphanRevocationRefs = append(orphanRevocationRefs, orphanRevocation)
			}
		}
	}
	return orphanRevocationRefs
}

func containsOrphanRevocation(values []*OrphanRevocationTokenWrapper, candidate *OrphanRevocationTokenWrapper) bool {
	for _, v := range values {
		if v.Id() == candidate.Id() {
			return true
		}
	}
	return false
}

// CrossCertificates returns a list of cross-certificates. Port of
// getCrossCertificates(CertificateWrapper).
func (d *DiagnosticData) CrossCertificates(certificate *CertificateWrapper) []*CertificateWrapper {
	var crossCertificates []*CertificateWrapper
	for _, candidate := range d.EquivalentCertificates(certificate) {
		if certificate.CertificateDN() != candidate.CertificateDN() || certificate.CertificateIssuerDN() != candidate.CertificateIssuerDN() {
			crossCertificates = append(crossCertificates, candidate)
		}
	}
	return crossCertificates
}

// OrphanCrossCertificates returns a list of orphan cross-certificates. Port of
// getOrphanCrossCertificates(CertificateWrapper).
func (d *DiagnosticData) OrphanCrossCertificates(certificate *CertificateWrapper) []*OrphanCertificateTokenWrapper {
	var crossCertificates []*OrphanCertificateTokenWrapper
	for _, candidate := range d.OrphanEquivalentCertificates(certificate) {
		if certificate.CertificateDN() != candidate.CertificateDN() || certificate.CertificateIssuerDN() != candidate.CertificateIssuerDN() {
			crossCertificates = append(crossCertificates, candidate)
		}
	}
	return crossCertificates
}

// EquivalentCertificates returns a list of equivalent certificates (certificates with the same
// public key). Port of getEquivalentCertificates(CertificateWrapper).
func (d *DiagnosticData) EquivalentCertificates(certificate *CertificateWrapper) []*CertificateWrapper {
	var equivalentCertificates []*CertificateWrapper
	for _, candidate := range d.UsedCertificates() {
		if !certificate.Equals(candidate) && certificate.EntityKey() == candidate.EntityKey() {
			equivalentCertificates = append(equivalentCertificates, candidate)
		}
	}
	return equivalentCertificates
}

// OrphanEquivalentCertificates returns a list of orphan equivalent certificates (certificates
// with the same public key). Port of getOrphanEquivalentCertificates(CertificateWrapper).
func (d *DiagnosticData) OrphanEquivalentCertificates(certificate *CertificateWrapper) []*OrphanCertificateTokenWrapper {
	var equivalentCertificates []*OrphanCertificateTokenWrapper
	for _, candidate := range d.AllOrphanCertificateObjects() {
		if certificate.Id() != candidate.Id() && certificate.EntityKey() == candidate.EntityKey() {
			equivalentCertificates = append(equivalentCertificates, candidate)
		}
	}
	return equivalentCertificates
}

// Signatures retrieves a list of signature wrappers. Port of getSignatures().
func (d *DiagnosticData) Signatures() []*SignatureWrapper {
	if d.foundSignatures == nil {
		xmlSignatures := d.wrapped.Signatures.All()
		for _, xmlSignature := range xmlSignatures {
			d.foundSignatures = append(d.foundSignatures, NewSignatureWrapper(xmlSignature))
		}
	}
	return d.foundSignatures
}

// TimestampList retrieves a list of timestamp wrappers. Port of getTimestampList().
func (d *DiagnosticData) TimestampList() []*TimestampWrapper {
	if d.usedTimestamps == nil {
		xmlTimestamps := d.wrapped.UsedTimestamps.All()
		for _, xmlTimestamp := range xmlTimestamps {
			d.usedTimestamps = append(d.usedTimestamps, NewTimestampWrapper(xmlTimestamp))
		}
	}
	return d.usedTimestamps
}

// NonEvidenceRecordTimestamps returns a list of time-stamp tokens which are not evidence
// record time-stamps. Port of getNonEvidenceRecordTimestamps().
func (d *DiagnosticData) NonEvidenceRecordTimestamps() []*TimestampWrapper {
	var result []*TimestampWrapper
	for _, timestampWrapper := range d.TimestampList() {
		if !timestampWrapper.Type().IsEvidenceRecordTimestamp() {
			result = append(result, timestampWrapper)
		}
	}
	return result
}

// EvidenceRecords retrieves a list of evidence record wrappers. Port of getEvidenceRecords().
func (d *DiagnosticData) EvidenceRecords() []*EvidenceRecordWrapper {
	if d.foundEvidenceRecords == nil {
		xmlEvidenceRecords := d.wrapped.EvidenceRecords.All()
		for _, xmlEvidenceRecord := range xmlEvidenceRecords {
			d.foundEvidenceRecords = append(d.foundEvidenceRecords, NewEvidenceRecordWrapper(xmlEvidenceRecord))
		}
	}
	return d.foundEvidenceRecords
}

// EvidenceRecordById returns the EvidenceRecordWrapper corresponding to the given id. Port of
// getEvidenceRecordById(String).
func (d *DiagnosticData) EvidenceRecordById(id string) *EvidenceRecordWrapper {
	for _, evidenceRecord := range d.EvidenceRecords() {
		if id == evidenceRecord.Id() {
			return evidenceRecord
		}
	}
	return nil
}

// EAAs retrieves a list of EAA wrappers. Port of getEAAs().
func (d *DiagnosticData) EAAs() []*EAAWrapper {
	if d.foundEAAs == nil {
		xmlEAAs := d.wrapped.EAAs.All()
		for _, xmlEAA := range xmlEAAs {
			d.foundEAAs = append(d.foundEAAs, NewEAAWrapper(xmlEAA))
		}
	}
	return d.foundEAAs
}

// EAAById returns the EAAWrapper corresponding to the given id. Port of getEAAById(String).
func (d *DiagnosticData) EAAById(id string) *EAAWrapper {
	for _, eaa := range d.EAAs() {
		if id == eaa.Id() {
			return eaa
		}
	}
	return nil
}

// FirstEAAId returns the first EAA id. Port of getFirstEAAId().
func (d *DiagnosticData) FirstEAAId() string {
	return d.firstEAANullSafe().Id()
}

func (d *DiagnosticData) firstEAANullSafe() *EAAWrapper {
	eaas := d.EAAs()
	if len(eaas) != 0 {
		return eaas[0]
	}
	return NewEAAWrapper(&jaxb.XmlEAA{})
}

// UsedCertificates retrieves a list of certificate wrappers. Port of getUsedCertificates().
func (d *DiagnosticData) UsedCertificates() []*CertificateWrapper {
	if d.usedCertificates == nil {
		xmlCertificates := d.wrapped.UsedCertificates.All()
		for _, certificate := range xmlCertificates {
			d.usedCertificates = append(d.usedCertificates, NewCertificateWrapper(certificate))
		}
	}
	return d.usedCertificates
}

// AllSignatures returns signatures (not countersignatures). Port of getAllSignatures().
func (d *DiagnosticData) AllSignatures() []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, signatureWrapper := range d.Signatures() {
		if !signatureWrapper.IsCounterSignature() && !signatureWrapper.IsKeyBindingSignature() {
			signatures = append(signatures, signatureWrapper)
		}
	}
	return signatures
}

// AllCounterSignatures returns counter-signatures (not signatures). Port of
// getAllCounterSignatures().
func (d *DiagnosticData) AllCounterSignatures() []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, signatureWrapper := range d.Signatures() {
		if signatureWrapper.IsCounterSignature() {
			signatures = append(signatures, signatureWrapper)
		}
	}
	return signatures
}

// AllCounterSignaturesForMasterSignature returns a set of SignatureWrapper for a given
// masterSignatureWrapper. Port of
// getAllCounterSignaturesForMasterSignature(SignatureWrapper).
func (d *DiagnosticData) AllCounterSignaturesForMasterSignature(masterSignatureWrapper *SignatureWrapper) []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, signatureWrapper := range d.Signatures() {
		if signatureWrapper.IsCounterSignature() && signatureWrapper.Parent().Equals(masterSignatureWrapper) {
			signatures = append(signatures, signatureWrapper)
		}
	}
	return signatures
}

// AllKeyBindingSignatures returns key binding signatures (not EAA signatures). Port of
// getAllKeyBindingSignatures().
func (d *DiagnosticData) AllKeyBindingSignatures() []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, signatureWrapper := range d.Signatures() {
		if signatureWrapper.IsKeyBindingSignature() {
			signatures = append(signatures, signatureWrapper)
		}
	}
	return signatures
}

// AllRevocationData returns all revocation data. Port of getAllRevocationData().
func (d *DiagnosticData) AllRevocationData() []*RevocationWrapper {
	var revocationData []*RevocationWrapper
	for _, xmlRevocation := range d.wrapped.UsedRevocations.All() {
		revocationData = append(revocationData, NewRevocationWrapper(xmlRevocation))
	}
	return revocationData
}

// LatestRevocationDataForCertificate returns the last actual revocation for the given
// certificate. Port of getLatestRevocationDataForCertificate(CertificateWrapper).
func (d *DiagnosticData) LatestRevocationDataForCertificate(certificate *CertificateWrapper) *CertificateRevocationWrapper {
	var latest *CertificateRevocationWrapper
	certificateRevocationData := certificate.CertificateRevocationData()
	for _, certRevoc := range certificateRevocationData {
		if latest == nil || (latest.ProductionDate() != nil && certRevoc != nil && certRevoc.ProductionDate() != nil &&
			latest.ProductionDate().Before(*certRevoc.ProductionDate())) {
			latest = certRevoc
		}
	}
	return latest
}

// AllEAA returns all electronic attestation of attributes (EAAs). Port of getAllEAA().
func (d *DiagnosticData) AllEAA() []*EAAWrapper {
	var eaas []*EAAWrapper
	for _, xmlEAA := range d.wrapped.EAAs.All() {
		eaas = append(eaas, NewEAAWrapper(xmlEAA))
	}
	return eaas
}

// AllEAARevocationTokens returns all EAA revocation tokens. Port of
// getAllEAARevocationTokens().
func (d *DiagnosticData) AllEAARevocationTokens() []*EAARevocationTokenWrapper {
	var eaaStatusTokens []*EAARevocationTokenWrapper
	for _, xmlEAARevocationToken := range d.wrapped.UsedEAARevocationTokens.All() {
		eaaStatusTokens = append(eaaStatusTokens, NewEAARevocationTokenWrapper(xmlEAARevocationToken))
	}
	return eaaStatusTokens
}

// CertificateById returns CertificateWrapper with the given id. Port of
// getCertificateById(String).
func (d *DiagnosticData) CertificateById(id string) *CertificateWrapper {
	for _, certificateWrapper := range d.UsedCertificates() {
		if id == certificateWrapper.Id() {
			return certificateWrapper
		}
	}
	return nil
}

// RevocationById returns RevocationWrapper with the given id. Port of
// getRevocationById(String).
func (d *DiagnosticData) RevocationById(id string) *RevocationWrapper {
	for _, revocationWrapper := range d.AllRevocationData() {
		if id == revocationWrapper.Id() {
			return revocationWrapper
		}
	}
	return nil
}

// OriginalSignerDocuments returns a complete list of original signer documents signed by all
// signatures. Port of getOriginalSignerDocuments().
func (d *DiagnosticData) OriginalSignerDocuments() []*SignerDataWrapper {
	var signerDocuments []*SignerDataWrapper
	for _, signatureWrapper := range d.Signatures() {
		for _, signatureScope := range signatureWrapper.SignatureScopes() {
			signerData := signatureScope.SignerData
			if signerData != nil {
				wrappedSignedData := NewSignerDataWrapper(signerData)
				if !containsSignerData(signerDocuments, wrappedSignedData) {
					signerDocuments = append(signerDocuments, wrappedSignedData)
				}
			}
		}
	}
	return signerDocuments
}

func containsSignerData(values []*SignerDataWrapper, candidate *SignerDataWrapper) bool {
	for _, v := range values {
		if v.Equals(candidate) {
			return true
		}
	}
	return false
}

// AllSignerDocuments returns a list of all covered documents, including the ones covering by
// timestamp(s), when applicable. Port of getAllSignerDocuments().
func (d *DiagnosticData) AllSignerDocuments() []*SignerDataWrapper {
	var signerDocuments []*SignerDataWrapper
	for _, signerData := range d.wrapped.OriginalDocuments.All() {
		signerDocuments = append(signerDocuments, NewSignerDataWrapper(signerData))
	}
	return signerDocuments
}

// JaxbModel returns the jaxb model of the diagnostic data. Port of getJaxbModel().
func (d *DiagnosticData) JaxbModel() *jaxb.XmlDiagnosticData {
	return d.wrapped
}

// IsContainerInfoPresent checks if the document is a container (ASiC). Port of
// isContainerInfoPresent().
func (d *DiagnosticData) IsContainerInfoPresent() bool {
	return d.wrapped.ContainerInfo != nil
}

// ContainerType returns the container type. Port of getContainerType().
func (d *DiagnosticData) ContainerType() enumerations.ASiCContainerType {
	containerInfo := d.wrapped.ContainerInfo
	if containerInfo != nil && containerInfo.ContainerType != nil {
		return enumerations.ASiCContainerType(*containerInfo.ContainerType)
	}
	return ""
}

// ZipComment returns the zip comment (if the document is a container). Port of
// getZipComment().
func (d *DiagnosticData) ZipComment() string {
	containerInfo := d.wrapped.ContainerInfo
	if containerInfo != nil && containerInfo.ZipComment != nil {
		return *containerInfo.ZipComment
	}
	return ""
}

// IsMimetypeFilePresent checks if the container has a mimetype file. Port of
// isMimetypeFilePresent().
func (d *DiagnosticData) IsMimetypeFilePresent() bool {
	containerInfo := d.wrapped.ContainerInfo
	if containerInfo != nil {
		return containerInfo.MimeTypeFilePresent != nil && *containerInfo.MimeTypeFilePresent
	}
	return false
}

// MimetypeFileContent returns the content of the mimetype file (if container). Port of
// getMimetypeFileContent().
func (d *DiagnosticData) MimetypeFileContent() string {
	containerInfo := d.wrapped.ContainerInfo
	if containerInfo != nil && containerInfo.MimeTypeContent != nil {
		return *containerInfo.MimeTypeContent
	}
	return ""
}

// ContainerInfo returns information about ASiC container (when applicable). Port of
// getContainerInfo().
func (d *DiagnosticData) ContainerInfo() *jaxb.XmlContainerInfo {
	return d.wrapped.ContainerInfo
}

// ManifestFiles gets a list of all manifest files extracted from the ASiC container. Port of
// getManifestFiles().
func (d *DiagnosticData) ManifestFiles() []*jaxb.XmlManifestFile {
	if d.wrapped.ContainerInfo != nil {
		return d.wrapped.ContainerInfo.ManifestFiles.All()
	}
	return nil
}

// ManifestFileForFilename gets an XmlManifestFile for the given filename document. Port of
// getManifestFileForFilename(String).
func (d *DiagnosticData) ManifestFileForFilename(filename string) *jaxb.XmlManifestFile {
	if filename != "" {
		for _, manifestFile := range d.ManifestFiles() {
			if manifestFile.SignatureFilename != nil && filename == *manifestFile.SignatureFilename {
				return manifestFile
			}
		}
	}
	return nil
}

// ContainerContentFilenames gets a list of all original signed document filenames. Port of
// getContainerContentFilenames().
func (d *DiagnosticData) ContainerContentFilenames() []string {
	if d.wrapped.ContainerInfo != nil {
		return d.wrapped.ContainerInfo.ContentFiles.All()
	}
	return nil
}

// IsPDFAValidationPerformed returns whether a document has been validated against PDF/A
// compliance. Port of isPDFAValidationPerformed().
func (d *DiagnosticData) IsPDFAValidationPerformed() bool {
	return d.wrapped.PDFAInfo != nil
}

// PDFAProfileId returns evaluated PDF/A profile Id. Port of getPDFAProfileId().
func (d *DiagnosticData) PDFAProfileId() string {
	if d.wrapped.PDFAInfo != nil && d.wrapped.PDFAInfo.ProfileId != nil {
		return *d.wrapped.PDFAInfo.ProfileId
	}
	return ""
}

// IsPDFACompliant returns whether the document is a PDF/A compliant (PDF/A validation shall be
// performed!). Port of isPDFACompliant().
func (d *DiagnosticData) IsPDFACompliant() bool {
	if d.wrapped.PDFAInfo != nil {
		return d.wrapped.PDFAInfo.Compliant
	}
	return false
}

// PDFAValidationErrors returns a collection of PDF/A validation errors occurred during the
// validation. Port of getPDFAValidationErrors().
func (d *DiagnosticData) PDFAValidationErrors() []string {
	if d.wrapped.PDFAInfo != nil {
		return d.wrapped.PDFAInfo.ValidationMessages.All()
	}
	return nil
}

// WebsiteUrl gets the remote website URL used to establish a TLS/SSL secure connection. NOTE:
// this method is used on QWAC validation. Port of getWebsiteUrl().
func (d *DiagnosticData) WebsiteUrl() string {
	if d.wrapped.ConnectionInfo != nil && d.wrapped.ConnectionInfo.Url != nil {
		return *d.wrapped.ConnectionInfo.Url
	}
	return ""
}

// TLSCertificateBindingUrl gets the TLS Certificate Binding URL, when present (i.e. URL under
// the 'Link' response header). NOTE: this method is used on QWAC validation. Port of
// getTLSCertificateBindingUrl().
func (d *DiagnosticData) TLSCertificateBindingUrl() string {
	if d.wrapped.ConnectionInfo != nil && d.wrapped.ConnectionInfo.TLSCertificateBindingUrl != nil {
		return *d.wrapped.ConnectionInfo.TLSCertificateBindingUrl
	}
	return ""
}

// TLSCertificate gets a TLS certificate used to establish a secure connection during the
// TLS/SSL handshake. NOTE: this method is used on QWAC validation. Port of
// getTLSCertificate().
func (d *DiagnosticData) TLSCertificate() *CertificateWrapper {
	if d.wrapped.ConnectionInfo != nil && d.wrapped.ConnectionInfo.TLSCertificate != nil &&
		d.wrapped.ConnectionInfo.TLSCertificate.Certificate != nil {
		return NewCertificateWrapper(d.wrapped.ConnectionInfo.TLSCertificate.Certificate)
	}
	return nil
}

// TLSCertificateBindingSignature gets the TLS Certificate Binding signature, when present
// (i.e. accessed from the URL under the 'Link' response header). NOTE: this method is used on
// QWAC validation. Port of getTLSCertificateBindingSignature().
func (d *DiagnosticData) TLSCertificateBindingSignature() *SignatureWrapper {
	if d.wrapped.ConnectionInfo != nil && d.wrapped.ConnectionInfo.TLSCertificateBindingSignature != nil &&
		d.wrapped.ConnectionInfo.TLSCertificateBindingSignature.Signature != nil {
		return NewSignatureWrapper(d.wrapped.ConnectionInfo.TLSCertificateBindingSignature.Signature)
	}
	return nil
}

// EAAPresentationInfo returns information about EAA Presentation document. Port of
// getEAAPresentationInfo().
func (d *DiagnosticData) EAAPresentationInfo() *jaxb.XmlEAAPresentationInfo {
	return d.wrapped.EAAPresentationInfo
}

// EAAPresentationType gets type of the EAA Presentation document. Port of
// getEAAPresentationType().
func (d *DiagnosticData) EAAPresentationType() enumerations.EAAPresentationType {
	eaaPresentationInfo := d.EAAPresentationInfo()
	if eaaPresentationInfo != nil && eaaPresentationInfo.EAAPresentationType != nil {
		return enumerations.EAAPresentationType(*eaaPresentationInfo.EAAPresentationType)
	}
	return ""
}

// TrustedLists returns the JAXB model of the used trusted lists. Port of getTrustedLists().
func (d *DiagnosticData) TrustedLists() []*jaxb.XmlTrustedList {
	var result []*jaxb.XmlTrustedList
	for _, xmlTrustedList := range d.wrapped.TrustedLists.All() {
		if xmlTrustedList.LOTL == nil || !*xmlTrustedList.LOTL {
			result = append(result, xmlTrustedList)
		}
	}
	return result
}

// ListOfTrustedLists returns the JAXB model of the LOTL. Port of getListOfTrustedLists().
func (d *DiagnosticData) ListOfTrustedLists() []*jaxb.XmlTrustedList {
	var result []*jaxb.XmlTrustedList
	for _, xmlTrustedList := range d.wrapped.TrustedLists.All() {
		if xmlTrustedList.LOTL != nil && *xmlTrustedList.LOTL {
			result = append(result, xmlTrustedList)
		}
	}
	return result
}

// ListsOfTrustedEntities returns the JAXB model of the used lists of trusted entities. Port of
// getListsOfTrustedEntities().
func (d *DiagnosticData) ListsOfTrustedEntities() []*jaxb.XmlListOfTrustedEntities {
	var result []*jaxb.XmlListOfTrustedEntities
	for _, lote := range d.wrapped.ListsOfTrustedEntities.All() {
		if lote.LoLoTE == nil || !*lote.LoLoTE {
			result = append(result, lote)
		}
	}
	return result
}

// ListsOfListsOfTrustedEntities returns the JAXB model of the used lists of lists of trusted
// entities. Port of getListsOfListsOfTrustedEntities().
func (d *DiagnosticData) ListsOfListsOfTrustedEntities() []*jaxb.XmlListOfTrustedEntities {
	var result []*jaxb.XmlListOfTrustedEntities
	for _, lote := range d.wrapped.ListsOfTrustedEntities.All() {
		if lote.LoLoTE != nil && *lote.LoLoTE {
			result = append(result, lote)
		}
	}
	return result
}

// ValidationDate returns the validation time. Port of getValidationDate().
func (d *DiagnosticData) ValidationDate() *time.Time {
	if d.wrapped.ValidationDate == nil {
		return nil
	}
	t := d.wrapped.ValidationDate.Time()
	return &t
}
