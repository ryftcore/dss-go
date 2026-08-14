// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/TimestampWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// TimestampWrapper provides a user-friendly interface for dealing with the JAXB jaxb.XmlTimestamp
// object.
type TimestampWrapper struct {
	AbstractSignatureWrapperBase

	// timestamp is the wrapped XmlTimestamp.
	timestamp *jaxb.XmlTimestamp
}

// NewTimestampWrapper is the default constructor. Port of TimestampWrapper(XmlTimestamp);
// panics per Objects.requireNonNull(timestamp, "XmlTimestamp cannot be null!").
func NewTimestampWrapper(timestamp *jaxb.XmlTimestamp) *TimestampWrapper {
	if timestamp == nil {
		panic("XmlTimestamp cannot be null!")
	}
	w := &TimestampWrapper{timestamp: timestamp}
	w.InitSignatureWrapper(w)
	return w
}

// Id is the AbstractTokenProxy override. Port of getId().
func (w *TimestampWrapper) Id() string {
	if w.timestamp.Id != nil {
		return string(*w.timestamp.Id)
	}
	return ""
}

// IsTimestampDuplicated checks if the time-stamp's Id is duplicated within the validating
// document. Port of isTimestampDuplicated().
func (w *TimestampWrapper) IsTimestampDuplicated() bool {
	return w.timestamp.Duplicated != nil && *w.timestamp.Duplicated
}

// CurrentBasicSignature is the AbstractTokenProxy override. Port of getCurrentBasicSignature().
func (w *TimestampWrapper) CurrentBasicSignature() *jaxb.XmlBasicSignature {
	return w.timestamp.BasicSignature
}

// CurrentCertificateChain is the AbstractTokenProxy override. Port of
// getCurrentCertificateChain().
func (w *TimestampWrapper) CurrentCertificateChain() []*jaxb.XmlChainItem {
	return w.timestamp.CertificateChain.All()
}

// CurrentSigningCertificate is the AbstractTokenProxy override. Port of
// getCurrentSigningCertificate().
func (w *TimestampWrapper) CurrentSigningCertificate() *jaxb.XmlSigningCertificate {
	return w.timestamp.SigningCertificate
}

// FoundCertificates is the AbstractTokenProxy override. Port of foundCertificates().
func (w *TimestampWrapper) FoundCertificates() *FoundCertificatesProxy {
	return NewFoundCertificatesProxy(w.timestamp.FoundCertificates)
}

// FoundRevocations is the AbstractTokenProxy override. Port of foundRevocations().
func (w *TimestampWrapper) FoundRevocations() *FoundRevocationsProxy {
	return NewFoundRevocationsProxy(w.timestamp.FoundRevocations)
}

// EvidenceRecords returns a list of evidence records covering the time-stamp file (applicable
// for detached time-stamps only). Port of getEvidenceRecords().
func (w *TimestampWrapper) EvidenceRecords() []*EvidenceRecordWrapper {
	var result []*EvidenceRecordWrapper
	for _, xmlEvidenceRecord := range w.timestamp.FoundEvidenceRecords.All() {
		result = append(result, NewEvidenceRecordWrapper(xmlEvidenceRecord.EvidenceRecord))
	}
	return result
}

// EvidenceRecordIdsList returns a list of associated evidence record identifiers. Port of
// getEvidenceRecordIdsList().
func (w *TimestampWrapper) EvidenceRecordIdsList() []string {
	var result []string
	for _, evidenceRecordWrapper := range w.EvidenceRecords() {
		result = append(result, evidenceRecordWrapper.Id())
	}
	return result
}

// EvidenceRecordTimestampIds returns identifiers of all covering evidence record time-stamps.
// Port of getEvidenceRecordTimestampIds().
func (w *TimestampWrapper) EvidenceRecordTimestampIds() []string {
	var result []string
	for _, evidenceRecordWrapper := range w.EvidenceRecords() {
		result = append(result, evidenceRecordWrapper.TimestampIdsList()...)
	}
	return result
}

// Type returns the type of the timestamp. Port of getType().
func (w *TimestampWrapper) Type() enumerations.TimestampType {
	if w.timestamp.Type != nil {
		return enumerations.TimestampType(*w.timestamp.Type)
	}
	return ""
}

// ArchiveTimestampType returns archive timestamp type, if applicable. NOTE: returns "" for non
// archive timestamps. Port of getArchiveTimestampType().
func (w *TimestampWrapper) ArchiveTimestampType() enumerations.ArchiveTimestampType {
	if w.timestamp.ArchiveTimestampType != nil {
		return enumerations.ArchiveTimestampType(*w.timestamp.ArchiveTimestampType)
	}
	return ""
}

// EvidenceRecordTimestampType returns an evidence record archive timestamp type, if
// applicable. NOTE: returns "" for non evidence record archive timestamps. Port of
// getEvidenceRecordTimestampType().
func (w *TimestampWrapper) EvidenceRecordTimestampType() enumerations.EvidenceRecordTimestampType {
	if w.timestamp.EvidenceRecordTimestampType != nil {
		return enumerations.EvidenceRecordTimestampType(*w.timestamp.EvidenceRecordTimestampType)
	}
	return ""
}

// AtsHashIndexVersion gets the version of the ats-hash-index attribute, when present. NOTE:
// applicable only for CAdES archive-time-stamp-v3. Port of getAtsHashIndexVersion().
func (w *TimestampWrapper) AtsHashIndexVersion() enumerations.ArchiveTimestampHashIndexVersion {
	if w.timestamp.ArchiveTimestampHashIndex != nil && w.timestamp.ArchiveTimestampHashIndex.Version != nil {
		return enumerations.ArchiveTimestampHashIndexVersion(*w.timestamp.ArchiveTimestampHashIndex.Version)
	}
	return ""
}

// IsAtsHashIndexValid returns whether the ats-hash-index(-v3) attribute is valid, when present
// (all hashes match). NOTE: applicable only for CAdES archive-time-stamp-v3. Port of
// isAtsHashIndexValid().
func (w *TimestampWrapper) IsAtsHashIndexValid() bool {
	return w.timestamp.ArchiveTimestampHashIndex != nil && w.timestamp.ArchiveTimestampHashIndex.Valid
}

// AtsHashIndexValidationMessages returns a list of error messages occurred in the result of
// ats-hash-index(-v3) attribute validation, if any. NOTE: applicable only for CAdES
// archive-time-stamp-v3. Port of getAtsHashIndexValidationMessages().
func (w *TimestampWrapper) AtsHashIndexValidationMessages() []string {
	if w.timestamp.ArchiveTimestampHashIndex != nil {
		return w.timestamp.ArchiveTimestampHashIndex.Message
	}
	return nil
}

// ProductionTime returns the indicated production time of the timestamp. Port of
// getProductionTime().
func (w *TimestampWrapper) ProductionTime() *time.Time {
	if w.timestamp.ProductionTime != nil {
		t := w.timestamp.ProductionTime.Time()
		return &t
	}
	return nil
}

// MessageImprint returns message-imprint jaxb.XmlDigestMatcher. Port of getMessageImprint().
func (w *TimestampWrapper) MessageImprint() *jaxb.XmlDigestMatcher {
	for _, digestMatcher := range w.DigestMatchers() {
		if digestMatcher.Type != nil && enumerations.DigestMatcherType_MESSAGE_IMPRINT == enumerations.DigestMatcherType(*digestMatcher.Type) {
			return digestMatcher
		}
	}
	return nil
}

// IsMessageImprintDataFound indicates if the message-imprint is found (all the required data
// for message-imprint computation is present). Port of isMessageImprintDataFound().
func (w *TimestampWrapper) IsMessageImprintDataFound() bool {
	if messageImprint := w.MessageImprint(); messageImprint != nil {
		return messageImprint.DataFound
	}
	return false
}

// IsMessageImprintDataIntact indicates if the message-imprint is intact (matches the computed
// message-imprint). Port of isMessageImprintDataIntact().
func (w *TimestampWrapper) IsMessageImprintDataIntact() bool {
	if messageImprint := w.MessageImprint(); messageImprint != nil {
		return messageImprint.DataIntact
	}
	return false
}

// Filename is the AbstractSignatureWrapper override. Port of getFilename().
func (w *TimestampWrapper) Filename() string {
	if w.timestamp.TimestampFilename != nil {
		return *w.timestamp.TimestampFilename
	}
	return ""
}

// DigestMatchers is the AbstractTokenProxy override. Port of getDigestMatchers().
func (w *TimestampWrapper) DigestMatchers() []*jaxb.XmlDigestMatcher {
	return w.timestamp.DigestMatcher
}

// TimestampedObjects returns a complete list of all jaxb.XmlTimestampedObject covered by the
// timestamp. Port of getTimestampedObjects().
func (w *TimestampWrapper) TimestampedObjects() []*jaxb.XmlTimestampedObject {
	return w.timestamp.TimestampedObjects.All()
}

// TimestampedSignatures returns a list of SignatureWrapper covered be the current timestamp.
// Port of getTimestampedSignatures().
func (w *TimestampWrapper) TimestampedSignatures() []*SignatureWrapper {
	var signatures []*SignatureWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_SIGNATURE) {
		xmlSignature, ok := token.(*jaxb.XmlSignature)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_SIGNATURE))
		}
		signatures = append(signatures, NewSignatureWrapper(xmlSignature))
	}
	return signatures
}

// TimestampedCertificates returns a list of certificates covered be the current timestamp.
// Port of getTimestampedCertificates().
func (w *TimestampWrapper) TimestampedCertificates() []*CertificateWrapper {
	var certificates []*CertificateWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_CERTIFICATE) {
		xmlCertificate, ok := token.(*jaxb.XmlCertificate)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_CERTIFICATE))
		}
		certificates = append(certificates, NewCertificateWrapper(xmlCertificate))
	}
	return certificates
}

// TimestampedRevocations returns a list of revocation data covered be the current timestamp.
// Port of getTimestampedRevocations().
func (w *TimestampWrapper) TimestampedRevocations() []*RevocationWrapper {
	var revocations []*RevocationWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_REVOCATION) {
		xmlRevocation, ok := token.(*jaxb.XmlRevocation)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_REVOCATION))
		}
		revocations = append(revocations, NewRevocationWrapper(xmlRevocation))
	}
	return revocations
}

// TimestampedTimestamps returns a list of timestamps covered be the current timestamp. Port of
// getTimestampedTimestamps().
func (w *TimestampWrapper) TimestampedTimestamps() []*TimestampWrapper {
	var timestamps []*TimestampWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_TIMESTAMP) {
		xmlTimestamp, ok := token.(*jaxb.XmlTimestamp)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_TIMESTAMP))
		}
		timestamps = append(timestamps, NewTimestampWrapper(xmlTimestamp))
	}
	return timestamps
}

// TimestampedEvidenceRecords returns a list of evidence records covered be the current
// timestamp. Port of getTimestampedEvidenceRecords().
func (w *TimestampWrapper) TimestampedEvidenceRecords() []*EvidenceRecordWrapper {
	var evidenceRecords []*EvidenceRecordWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_EVIDENCE_RECORD) {
		xmlEvidenceRecord, ok := token.(*jaxb.XmlEvidenceRecord)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_EVIDENCE_RECORD))
		}
		evidenceRecords = append(evidenceRecords, NewEvidenceRecordWrapper(xmlEvidenceRecord))
	}
	return evidenceRecords
}

// TimestampedSignedData returns a list of Signed data covered be the current timestamp. Port
// of getTimestampedSignedData().
func (w *TimestampWrapper) TimestampedSignedData() []*SignerDataWrapper {
	var signerData []*SignerDataWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_SIGNED_DATA) {
		xmlSignerData, ok := token.(*jaxb.XmlSignerData)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_SIGNED_DATA))
		}
		signerData = append(signerData, NewSignerDataWrapper(xmlSignerData))
	}
	return signerData
}

// IsSigningCertificateIdentified indicates if the signing certificate reference is present
// within the timestamp token and matches the actual signing certificate. Port of
// isSigningCertificateIdentified().
func (w *TimestampWrapper) IsSigningCertificateIdentified() bool {
	signingCertificate := w.SigningCertificate()
	signingCertificateReference := w.SigningCertificateReference()
	if signingCertificate != nil && signingCertificateReference != nil {
		return signingCertificateReference.IsDigestValueMatch() &&
			(!signingCertificateReference.IsIssuerSerialPresent() || signingCertificateReference.IsIssuerSerialMatch())
	}
	return false
}

// timestampedObjectsByCategory is the private helper backing every TimestampedXxx() accessor.
// Port of the private getTimestampedObjectsByCategory(TimestampedObjectType).
func (w *TimestampWrapper) timestampedObjectsByCategory(category enumerations.TimestampedObjectType) []jaxb.XmlToken {
	var timestampedObjectIds []jaxb.XmlToken
	for _, timestampedObject := range w.TimestampedObjects() {
		if timestampedObject.Category != nil && category == enumerations.TimestampedObjectType(*timestampedObject.Category) &&
			timestampedObject.Token != nil {
			timestampedObjectIds = append(timestampedObjectIds, timestampedObject.Token.Token)
		}
	}
	return timestampedObjectIds
}

// AllTimestampedOrphanTokens returns a list of all OrphanTokens. Port of
// getAllTimestampedOrphanTokens().
func (w *TimestampWrapper) AllTimestampedOrphanTokens() []OrphanTokenWrapperOverrides {
	var timestampedObjectIds []OrphanTokenWrapperOverrides
	for _, c := range w.TimestampedOrphanCertificates() {
		timestampedObjectIds = append(timestampedObjectIds, c)
	}
	for _, r := range w.TimestampedOrphanRevocations() {
		timestampedObjectIds = append(timestampedObjectIds, r)
	}
	return timestampedObjectIds
}

// TimestampedOrphanCertificates returns a list of OrphanCertificateTokens. Port of
// getTimestampedOrphanCertificates().
func (w *TimestampWrapper) TimestampedOrphanCertificates() []*OrphanCertificateTokenWrapper {
	var orphanCertificates []*OrphanCertificateTokenWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_ORPHAN_CERTIFICATE) {
		xmlOrphanCertificateToken, ok := token.(*jaxb.XmlOrphanCertificateToken)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_ORPHAN_CERTIFICATE))
		}
		orphanCertificates = append(orphanCertificates, NewOrphanCertificateTokenWrapper(xmlOrphanCertificateToken))
	}
	return orphanCertificates
}

// TimestampedOrphanRevocations returns a list of OrphanRevocationTokens. Port of
// getTimestampedOrphanRevocations().
func (w *TimestampWrapper) TimestampedOrphanRevocations() []*OrphanRevocationTokenWrapper {
	var orphanRevocations []*OrphanRevocationTokenWrapper
	for _, token := range w.timestampedObjectsByCategory(enumerations.TimestampedObjectType_ORPHAN_REVOCATION) {
		xmlOrphanRevocationToken, ok := token.(*jaxb.XmlOrphanRevocationToken)
		if !ok {
			panic(fmt.Sprintf("Unexpected token of type [%T] found. Expected : %s", token, enumerations.TimestampedObjectType_ORPHAN_REVOCATION))
		}
		orphanRevocations = append(orphanRevocations, NewOrphanRevocationTokenWrapper(xmlOrphanRevocationToken))
	}
	return orphanRevocations
}

// Binaries is the AbstractTokenProxy override. Port of getBinaries().
func (w *TimestampWrapper) Binaries() []byte {
	if w.timestamp.Base64Encoded != nil {
		return []byte(*w.timestamp.Base64Encoded)
	}
	return nil
}

// DigestAlgoAndValue returns digest algorithm and value of the timestamp token binaries, when
// defined. Port of getDigestAlgoAndValue().
func (w *TimestampWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.timestamp.DigestAlgoAndValue
}

/* -------- PAdES RFC3161 Specific parameters --------- */

// PDFRevision returns a PAdES-specific PDF Revision info. NOTE: applicable only for PDF
// Document Timestamp. Port of getPDFRevision().
func (w *TimestampWrapper) PDFRevision() *PDFRevisionWrapper {
	if w.timestamp.PDFRevision != nil {
		return NewPDFRevisionWrapper(w.timestamp.PDFRevision)
	}
	return nil
}

// SignatureInformationStore returns a list if Signer Infos (Signer Information Store) from
// CAdES CMS Signed Data. Port of getSignatureInformationStore().
func (w *TimestampWrapper) SignatureInformationStore() []*jaxb.XmlSignerInfo {
	return w.timestamp.SignerInformationStore.All()
}

// IsTSAGeneralNamePresent checks if the tsa field of TSTInfo is present. Port of
// isTSAGeneralNamePresent().
func (w *TimestampWrapper) IsTSAGeneralNamePresent() bool {
	return w.timestamp.TSAGeneralName != nil
}

// TSAGeneralNameValue gets TSA General Name value. Port of getTSAGeneralNameValue().
func (w *TimestampWrapper) TSAGeneralNameValue() string {
	if w.IsTSAGeneralNamePresent() {
		return w.timestamp.TSAGeneralName.Value
	}
	return ""
}

// IsTSAGeneralNameMatch checks if the content of TSTInfo.tsa field matches the timestamp's
// issuer distinguishing name, without taking order into account. Port of
// isTSAGeneralNameMatch().
func (w *TimestampWrapper) IsTSAGeneralNameMatch() bool {
	return w.IsTSAGeneralNamePresent() && w.timestamp.TSAGeneralName.ContentMatch
}

// IsTSAGeneralNameOrderMatch checks if the content and the order of TSTInfo.tsa field match
// the timestamp's issuer distinguishing name. Port of isTSAGeneralNameOrderMatch().
func (w *TimestampWrapper) IsTSAGeneralNameOrderMatch() bool {
	return w.IsTSAGeneralNamePresent() && w.timestamp.TSAGeneralName.OrderMatch
}

// TimestampScopes returns Timestamp's Signature Scopes. Port of getTimestampScopes().
func (w *TimestampWrapper) TimestampScopes() []*jaxb.XmlSignatureScope {
	return w.timestamp.TimestampScopes.All()
}
