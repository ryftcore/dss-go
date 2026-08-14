// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/EvidenceRecordWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"fmt"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// EvidenceRecordWrapper provides a user-friendly interface for dealing with the JAXB
// jaxb.XmlEvidenceRecord object. Unlike most wrappers in this package, EvidenceRecordWrapper
// does not extend AbstractTokenProxy in Java: it is a standalone wrapper.
type EvidenceRecordWrapper struct {
	// evidenceRecord is the wrapped XmlEvidenceRecord.
	evidenceRecord *jaxb.XmlEvidenceRecord
}

// NewEvidenceRecordWrapper is the default constructor. Port of
// EvidenceRecordWrapper(XmlEvidenceRecord); panics per
// Objects.requireNonNull(evidenceRecord, "XmlEvidenceRecord cannot be null!").
func NewEvidenceRecordWrapper(evidenceRecord *jaxb.XmlEvidenceRecord) *EvidenceRecordWrapper {
	if evidenceRecord == nil {
		panic("XmlEvidenceRecord cannot be null!")
	}
	return &EvidenceRecordWrapper{evidenceRecord: evidenceRecord}
}

// Id gets unique identifier. Port of getId().
func (w *EvidenceRecordWrapper) Id() string {
	return w.evidenceRecord.Id
}

// IsEvidenceRecordDuplicated checks if the evidence record's Id is duplicated within the
// validating document. Port of isEvidenceRecordDuplicated().
func (w *EvidenceRecordWrapper) IsEvidenceRecordDuplicated() bool {
	return w.evidenceRecord.Duplicated != nil && *w.evidenceRecord.Duplicated
}

// Filename returns name of the evidence record's document, when applicable. Port of
// getFilename().
func (w *EvidenceRecordWrapper) Filename() string {
	return w.evidenceRecord.DocumentName
}

// DigestMatchers gets a list of digest matchers representing the associated archival data
// objects validation status. Port of getDigestMatchers().
func (w *EvidenceRecordWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return w.evidenceRecord.DigestMatchers
}

// FirstTimestamp returns initial time-stamp of the evidence record. Port of
// getFirstTimestamp().
func (w *EvidenceRecordWrapper) FirstTimestamp() *TimestampWrapper {
	timestampList := w.TimestampList()
	if len(timestampList) != 0 {
		return timestampList[0]
	}
	return nil
}

// TimestampList gets a list of time-stamp tokens associated with the evidence record. Port
// of getTimestampList().
func (w *EvidenceRecordWrapper) TimestampList() []*TimestampWrapper {
	var tsps []*TimestampWrapper
	for _, xmlFoundTimestamp := range w.evidenceRecord.EvidenceRecordTimestamps {
		tsps = append(tsps, NewTimestampWrapper(xmlFoundTimestamp.Timestamp))
	}
	return tsps
}

// TimestampIdsList returns a list of time-stamp identifiers associated with the Evidence
// Record. Port of getTimestampIdsList().
func (w *EvidenceRecordWrapper) TimestampIdsList() []string {
	var result []string
	for _, tsp := range w.TimestampList() {
		result = append(result, tsp.Id())
	}
	return result
}

// FoundCertificates returns a collection of certificate tokens embedded within Evidence
// Record. Port of foundCertificates().
func (w *EvidenceRecordWrapper) FoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(w.evidenceRecord.FoundCertificates)
}

// FoundRevocations returns a collection of revocation tokens embedded within Evidence Record.
// Port of foundRevocations().
func (w *EvidenceRecordWrapper) FoundRevocations() *FoundRevocationsProxy {
	return NewFoundRevocationsProxy(w.evidenceRecord.FoundRevocations)
}

// EvidenceRecordType gets the evidence record format type. Port of getEvidenceRecordType().
func (w *EvidenceRecordWrapper) EvidenceRecordType() enumerations.EvidenceRecordTypeEnum {
	return w.evidenceRecord.Type
}

// Origin gets the origin of the evidence record. Port of getOrigin().
func (w *EvidenceRecordWrapper) Origin() enumerations.EvidenceRecordOrigin {
	return w.evidenceRecord.Origin
}

// IsEmbedded gets whether the evidence record has been embedded into a signature (supported
// for XAdES and CAdES). Port of isEmbedded().
func (w *EvidenceRecordWrapper) IsEmbedded() bool {
	return w.evidenceRecord.Embedded != nil && *w.evidenceRecord.Embedded
}

// Parent returns a master-signature in case of a counter-signature. Port of getParent().
func (w *EvidenceRecordWrapper) Parent() *SignatureWrapper {
	parent := w.evidenceRecord.Parent
	if parent != nil {
		return NewSignatureWrapper(parent)
	}
	return nil
}

// IncorporationType gets the incorporation of the evidence record within an embedding
// signature. NOTE: applicable only for attached evidence records in CAdES. Port of
// getIncorporationType().
func (w *EvidenceRecordWrapper) IncorporationType() enumerations.EvidenceRecordIncorporationType {
	return w.evidenceRecord.IncorporationType
}

// IsStructuralValidationValid gets if a structural validation of the evidence record is valid.
// Port of isStructuralValidationValid().
func (w *EvidenceRecordWrapper) IsStructuralValidationValid() bool {
	return w.evidenceRecord.StructuralValidation != nil && w.evidenceRecord.StructuralValidation.Valid
}

// StructuralValidationMessages returns structural validation error messages, when
// applicable. Port of getStructuralValidationMessages().
func (w *EvidenceRecordWrapper) StructuralValidationMessages() []string {
	structuralValidation := w.evidenceRecord.StructuralValidation
	if structuralValidation != nil {
		return structuralValidation.Messages
	}
	return nil
}

// CoveredObjects returns a list of objects covered by the evidence record. Port of
// getCoveredObjects().
func (w *EvidenceRecordWrapper) CoveredObjects() []*jaxb.XmlTimestampedObject {
	return w.evidenceRecord.TimestampedObjects
}

// CoveredSignatures returns a list of SignatureWrapper covered by the current evidence
// record. Port of getCoveredSignatures().
func (w *EvidenceRecordWrapper) CoveredSignatures() []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_SIGNATURE) {
		xmlSignature, ok := token.(*jaxb.XmlSignature)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_SIGNATURE))
		}
		signatures = append(signatures, NewSignatureWrapper(xmlSignature))
	}
	return signatures
}

// CoveredCertificates returns a list of certificates covered by the current evidence
// record. Port of getCoveredCertificates().
func (w *EvidenceRecordWrapper) CoveredCertificates() []*CertificateWrapper {
	var certificates []*CertificateWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_CERTIFICATE) {
		xmlCertificate, ok := token.(*jaxb.XmlCertificate)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_CERTIFICATE))
		}
		certificates = append(certificates, NewCertificateWrapper(xmlCertificate))
	}
	return certificates
}

// CoveredRevocations returns a list of revocation data covered by the current evidence
// record. Port of getCoveredRevocations().
func (w *EvidenceRecordWrapper) CoveredRevocations() []*RevocationWrapper {
	var revocations []*RevocationWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_REVOCATION) {
		xmlRevocation, ok := token.(*jaxb.XmlRevocation)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_REVOCATION))
		}
		revocations = append(revocations, NewRevocationWrapper(xmlRevocation))
	}
	return revocations
}

// CoveredTimestamps returns a list of timestamps covered by the current evidence record.
// Port of getCoveredTimestamps().
func (w *EvidenceRecordWrapper) CoveredTimestamps() []*TimestampWrapper {
	var timestamps []*TimestampWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_TIMESTAMP) {
		xmlTimestamp, ok := token.(*jaxb.XmlTimestamp)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_TIMESTAMP))
		}
		timestamps = append(timestamps, NewTimestampWrapper(xmlTimestamp))
	}
	return timestamps
}

// CoveredEvidenceRecords returns a list of evidence records covered by the current evidence
// record. Port of getCoveredEvidenceRecords().
func (w *EvidenceRecordWrapper) CoveredEvidenceRecords() []*EvidenceRecordWrapper {
	var evidenceRecords []*EvidenceRecordWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_EVIDENCE_RECORD) {
		xmlEvidenceRecord, ok := token.(*jaxb.XmlEvidenceRecord)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_EVIDENCE_RECORD))
		}
		evidenceRecords = append(evidenceRecords, NewEvidenceRecordWrapper(xmlEvidenceRecord))
	}
	return evidenceRecords
}

// CoveredSignedData returns a list of Signed data covered by the current evidence record.
// Port of getCoveredSignedData().
func (w *EvidenceRecordWrapper) CoveredSignedData() []*SignerDataWrapper {
	var signerData []*SignerDataWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_SIGNED_DATA) {
		xmlSignerData, ok := token.(*jaxb.XmlSignerData)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_SIGNED_DATA))
		}
		signerData = append(signerData, NewSignerDataWrapper(xmlSignerData))
	}
	return signerData
}

// AllCoveredOrphanTokens returns a list of all OrphanTokens covered by the evidence record.
// Port of getAllCoveredOrphanTokens().
func (w *EvidenceRecordWrapper) AllCoveredOrphanTokens() []OrphanTokenWrapperOverrides {
	var timestampedObjectIds []OrphanTokenWrapperOverrides
	for _, c := range w.CoveredOrphanCertificates() {
		timestampedObjectIds = append(timestampedObjectIds, c)
	}
	for _, r := range w.CoveredOrphanRevocations() {
		timestampedObjectIds = append(timestampedObjectIds, r)
	}
	return timestampedObjectIds
}

// CoveredOrphanCertificates returns a list of OrphanCertificateTokens covered by the
// evidence record. Port of getCoveredOrphanCertificates().
func (w *EvidenceRecordWrapper) CoveredOrphanCertificates() []*OrphanCertificateTokenWrapper {
	var orphanCertificates []*OrphanCertificateTokenWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_ORPHAN_CERTIFICATE) {
		xmlOrphanCertificateToken, ok := token.(*jaxb.XmlOrphanCertificateToken)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_ORPHAN_CERTIFICATE))
		}
		orphanCertificates = append(orphanCertificates, NewOrphanCertificateTokenWrapper(xmlOrphanCertificateToken))
	}
	return orphanCertificates
}

// CoveredOrphanRevocations returns a list of OrphanRevocationTokens covered by the evidence
// record. Port of getCoveredOrphanRevocations().
func (w *EvidenceRecordWrapper) CoveredOrphanRevocations() []*OrphanRevocationTokenWrapper {
	var orphanRevocations []*OrphanRevocationTokenWrapper
	for _, token := range w.getCoveredObjectsByCategory(enumerations.TimestampedObjectType_ORPHAN_REVOCATION) {
		xmlOrphanRevocationToken, ok := token.(*jaxb.XmlOrphanRevocationToken)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_ORPHAN_REVOCATION))
		}
		orphanRevocations = append(orphanRevocations, NewOrphanRevocationTokenWrapper(xmlOrphanRevocationToken))
	}
	return orphanRevocations
}

func (w *EvidenceRecordWrapper) getCoveredObjectsByCategory(category enumerations.TimestampedObjectType) []jaxb.XmlAbstractToken {
	var coveredObjectIds []jaxb.XmlAbstractToken
	for _, coveredObject := range w.CoveredObjects() {
		if category == coveredObject.Category {
			coveredObjectIds = append(coveredObjectIds, coveredObject.Token)
		}
	}
	return coveredObjectIds
}

// EvidenceRecordScopes returns Evidence record's Signature Scopes. Port of
// getEvidenceRecordScopes().
func (w *EvidenceRecordWrapper) EvidenceRecordScopes() []*jaxb.XmlSignatureScope {
	return w.evidenceRecord.EvidenceRecordScopes
}

// Binaries returns binaries of the evidence record. Port of getBinaries().
func (w *EvidenceRecordWrapper) Binaries() []byte {
	return w.evidenceRecord.Base64Encoded
}

// DigestAlgoAndValue returns digest algorithm and value of the timestamp token binaries,
// when defined. Port of getDigestAlgoAndValue().
func (w *EvidenceRecordWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.evidenceRecord.DigestAlgoAndValue
}
