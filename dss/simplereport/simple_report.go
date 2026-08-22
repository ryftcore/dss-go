// Ported from dss-simple-report-jaxb/src/main/java/eu/europa/esig/dss/simplereport/SimpleReport.java
// (DSS 6.5.RC1).

package simplereport

import (
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/simplereport/jaxb"
)

// SimpleReport is a SimpleReport holder to fetch values from a JAXB
// SimpleReport.
type SimpleReport struct {
	wrapped *jaxb.XmlSimpleReport
}

// NewSimpleReport builds a SimpleReport around a JAXB XmlSimpleReport. Port
// of the SimpleReport(XmlSimpleReport) constructor.
func NewSimpleReport(wrapped *jaxb.XmlSimpleReport) *SimpleReport {
	return &SimpleReport{wrapped: wrapped}
}

// GetValidationTime returns the validation time. Port of getValidationTime().
func (r *SimpleReport) GetValidationTime() *time.Time {
	return xsTime(r.wrapped.ValidationTime)
}

// GetIndication returns the indication obtained after the validation of a
// token with the given DSS unique identifier. Port of getIndication(String).
func (r *SimpleReport) GetIndication(tokenID string) enumerations.Indication {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc == nil {
		return ""
	}
	return tc.Indication.Indication()
}

// GetSubIndication returns the sub-indication obtained after the validation
// of the token. Port of getSubIndication(String).
func (r *SimpleReport) GetSubIndication(tokenID string) enumerations.SubIndication {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc == nil || tc.SubIndication == nil {
		return ""
	}
	return tc.SubIndication.SubIndication()
}

// IsValid checks if a signature is valid (TOTAL_PASSED) or timestamp
// validation PASSED. Port of isValid(String).
func (r *SimpleReport) IsValid(tokenID string) bool {
	ind := r.GetIndication(tokenID)
	return ind == enumerations.IndicationTotalPassed || ind == enumerations.IndicationPassed
}

// GetSignatureIdList retrieves the signature ids. Port of
// getSignatureIdList().
func (r *SimpleReport) GetSignatureIdList() []string {
	var ids []string
	for _, token := range r.wrapped.SignatureOrTimestampOrEvidenceRecord {
		if _, ok := token.(*jaxb.XmlSignature); ok {
			ids = append(ids, tokenID(token))
		}
	}
	return ids
}

// GetTimestampIdList retrieves the timestamp ids. Port of
// getTimestampIdList().
func (r *SimpleReport) GetTimestampIdList() []string {
	var ids []string
	for _, token := range r.wrapped.SignatureOrTimestampOrEvidenceRecord {
		if _, ok := token.(*jaxb.XmlTimestamp); ok {
			ids = append(ids, tokenID(token))
		}
	}
	return ids
}

// GetEvidenceRecordIdList retrieves the evidence record ids. Port of
// getEvidenceRecordIdList().
func (r *SimpleReport) GetEvidenceRecordIdList() []string {
	var ids []string
	for _, token := range r.wrapped.SignatureOrTimestampOrEvidenceRecord {
		if _, ok := token.(*jaxb.XmlEvidenceRecord); ok {
			ids = append(ids, tokenID(token))
		}
	}
	return ids
}

// GetEAAIdList retrieves the EAA ids. Port of getEAAIdList().
func (r *SimpleReport) GetEAAIdList() []string {
	var ids []string
	for _, token := range r.wrapped.SignatureOrTimestampOrEvidenceRecord {
		if _, ok := token.(*jaxb.XmlEAA); ok {
			ids = append(ids, tokenID(token))
		}
	}
	return ids
}

// GetFirstSignatureId returns the first signature id. Port of
// getFirstSignatureId().
func (r *SimpleReport) GetFirstSignatureId() string {
	ids := r.GetSignatureIdList()
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

// GetFirstTimestampId returns the first timestamp id. Port of
// getFirstTimestampId().
func (r *SimpleReport) GetFirstTimestampId() string {
	ids := r.GetTimestampIdList()
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

// GetFirstEvidenceRecordId returns the first evidence record id. Port of
// getFirstEvidenceRecordId().
func (r *SimpleReport) GetFirstEvidenceRecordId() string {
	ids := r.GetEvidenceRecordIdList()
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

// GetFirstEAAId returns the first EAA id. Port of getFirstEAAId().
func (r *SimpleReport) GetFirstEAAId() string {
	ids := r.GetEAAIdList()
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

// GetDocumentFilename returns a file name for the validated document. Port
// of getDocumentFilename().
func (r *SimpleReport) GetDocumentFilename() string {
	if r.wrapped.DocumentName != nil {
		return *r.wrapped.DocumentName
	}
	return ""
}

// GetTokenFilename returns a file name for a given tokenId. Port of
// getTokenFilename(String).
func (r *SimpleReport) GetTokenFilename(tokenID string) string {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.Filename != nil {
		return *tc.Filename
	}
	return ""
}

// GetCertificateChain returns a certificate chain for a given tokenId. Port
// of getCertificateChain(String).
func (r *SimpleReport) GetCertificateChain(tokenID string) *jaxb.XmlCertificateChain {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil {
		return tc.CertificateChain
	}
	return nil
}

// GetAdESValidationErrors retrieves the ETSI EN 319 102-1 AdES validation
// errors for a given token by id. Port of getAdESValidationErrors(String).
func (r *SimpleReport) GetAdESValidationErrors(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.AdESValidationDetails != nil {
		return convertMessages(tc.AdESValidationDetails.Error)
	}
	return nil
}

// GetAdESValidationWarnings retrieves the ETSI EN 319 102-1 AdES validation
// warnings for a given token by id. Port of getAdESValidationWarnings(String).
func (r *SimpleReport) GetAdESValidationWarnings(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.AdESValidationDetails != nil {
		return convertMessages(tc.AdESValidationDetails.Warning)
	}
	return nil
}

// GetAdESValidationInfo retrieves the ETSI EN 319 102-1 AdES validation
// information for a given token by id. Port of getAdESValidationInfo(String).
func (r *SimpleReport) GetAdESValidationInfo(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.AdESValidationDetails != nil {
		return convertMessages(tc.AdESValidationDetails.Info)
	}
	return nil
}

// GetQualificationErrors retrieves the qualification process's errors for a
// given token by id. Port of getQualificationErrors(String).
func (r *SimpleReport) GetQualificationErrors(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.QualificationDetails != nil {
		return convertMessages(tc.QualificationDetails.Error)
	}
	return nil
}

// GetQualificationWarnings retrieves the qualification process's warnings
// for a given token by id. Port of getQualificationWarnings(String).
func (r *SimpleReport) GetQualificationWarnings(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.QualificationDetails != nil {
		return convertMessages(tc.QualificationDetails.Warning)
	}
	return nil
}

// GetQualificationInfo retrieves the qualification process's information
// for a given token by id. Port of getQualificationInfo(String).
func (r *SimpleReport) GetQualificationInfo(tokenID string) []Message {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil && tc.QualificationDetails != nil {
		return convertMessages(tc.QualificationDetails.Info)
	}
	return nil
}

// GetSignatureQualification returns the signature type: QES, AdES, AdESqc,
// NA. Port of getSignatureQualification(String).
func (r *SimpleReport) GetSignatureQualification(signatureID string) enumerations.SignatureQualification {
	qualif := enumerations.SignatureQualificationNA
	sig := r.signatureByID(signatureID)
	if sig != nil && sig.SignatureLevel != nil {
		qualif = sig.SignatureLevel.Value.SignatureQualification()
	}
	return qualif
}

// GetSignatureFormat returns the signature format (XAdES_BASELINE_B...).
// Port of getSignatureFormat(String).
func (r *SimpleReport) GetSignatureFormat(signatureID string) enumerations.SignatureLevel {
	sig := r.signatureByID(signatureID)
	if sig != nil {
		return sig.SignatureFormat.SignatureLevel()
	}
	return ""
}

// GetBestSignatureTime returns the best-signature-time. Port of
// getBestSignatureTime(String).
func (r *SimpleReport) GetBestSignatureTime(signatureID string) *time.Time {
	sig := r.signatureByID(signatureID)
	if sig != nil {
		return xsTime(sig.BestSignatureTime)
	}
	return nil
}

// GetSigningTime returns the signature time. Port of getSigningTime(String).
func (r *SimpleReport) GetSigningTime(signatureID string) *time.Time {
	sig := r.signatureByID(signatureID)
	if sig != nil {
		return xsTime(sig.SigningTime)
	}
	return nil
}

// GetExtensionPeriodMin returns the minimal useful date for the token
// extension, when the token validation is TOTAL_PASSED or PASSED. Port of
// getExtensionPeriodMin(String).
func (r *SimpleReport) GetExtensionPeriodMin(tokenID string) *time.Time {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil {
		return xsTime(tc.ExtensionPeriodMin)
	}
	return nil
}

// GetExtensionPeriodMax returns the maximum useful date for the token
// extension, when the token validation is TOTAL_PASSED or PASSED. Port of
// getExtensionPeriodMax(String).
func (r *SimpleReport) GetExtensionPeriodMax(tokenID string) *time.Time {
	tc := tokenContent(r.tokenByID(tokenID))
	if tc != nil {
		return xsTime(tc.ExtensionPeriodMax)
	}
	return nil
}

// GetSignedBy returns the signature's signer name. Port of
// getSignedBy(String).
func (r *SimpleReport) GetSignedBy(signatureID string) string {
	sig := r.signatureByID(signatureID)
	if sig != nil && sig.SignedBy != nil {
		return *sig.SignedBy
	}
	return ""
}

// GetSignaturesCount returns the number of signatures. Port of
// getSignaturesCount().
func (r *SimpleReport) GetSignaturesCount() int {
	return r.wrapped.SignaturesCount
}

// GetValidSignaturesCount returns the number of valid signatures. Port of
// getValidSignaturesCount().
func (r *SimpleReport) GetValidSignaturesCount() int {
	return r.wrapped.ValidSignaturesCount
}

// GetProductionTime returns the timestamp production time. Port of
// getProductionTime(String).
func (r *SimpleReport) GetProductionTime(timestampID string) *time.Time {
	ts := r.timestampByID(timestampID)
	if ts != nil {
		return xsTime(ts.ProductionTime)
	}
	return nil
}

// GetProducedBy returns the timestamp's producer name. Port of
// getProducedBy(String).
func (r *SimpleReport) GetProducedBy(timestampID string) string {
	ts := r.timestampByID(timestampID)
	if ts != nil && ts.ProducedBy != nil {
		return *ts.ProducedBy
	}
	return ""
}

// GetTimestampQualification returns the timestamp's qualification. Port of
// getTimestampQualification(String).
func (r *SimpleReport) GetTimestampQualification(timestampID string) enumerations.TimestampQualification {
	ts := r.timestampByID(timestampID)
	if ts != nil && ts.TimestampLevel != nil {
		return ts.TimestampLevel.Value.TimestampQualification()
	}
	return ""
}

// GetEAAQualification returns the first determined EAA's qualification.
// This method could be used for a simple EAAQualification result
// extraction, suitable for the most use cases. Should you need a more
// comprehensive validation output, use GetEAAQualifications. Port of
// getEAAQualification(String).
func (r *SimpleReport) GetEAAQualification(eaaPresentationID string) enumerations.EAAQualification {
	eaa := r.GetEAAById(eaaPresentationID)
	if eaa != nil && len(eaa.EAALevel) > 0 {
		return eaa.EAALevel[0].Value.EAAQualification()
	}
	return ""
}

// GetEAAQualifications returns a list of determined EAA's qualifications.
// This list should be used if a comprehensive result of EAA validation is
// required, as potentially a token may be qualified with different outputs
// during the validation process, even though it should not happen in
// production environments. Port of getEAAQualifications(String).
func (r *SimpleReport) GetEAAQualifications(eaaPresentationID string) []enumerations.EAAQualification {
	eaa := r.GetEAAById(eaaPresentationID)
	if eaa == nil || len(eaa.EAALevel) == 0 {
		return nil
	}
	out := make([]enumerations.EAAQualification, 0, len(eaa.EAALevel))
	for _, level := range eaa.EAALevel {
		out = append(out, level.Value.EAAQualification())
	}
	return out
}

// tokenByID returns a wrapper for the given token id. Port of the
// private getTokenById(String).
func (r *SimpleReport) tokenByID(tokenID string) jaxb.XmlTokenItem {
	return embeddedTokenByID(r.wrapped.SignatureOrTimestampOrEvidenceRecord, tokenID)
}

// signatureByID returns a wrapper for the given signature. Port of the
// private getSignatureById(String).
func (r *SimpleReport) signatureByID(signatureID string) *jaxb.XmlSignature {
	if sig, ok := r.tokenByID(signatureID).(*jaxb.XmlSignature); ok {
		return sig
	}
	return nil
}

// timestampByID returns a wrapper for the given timestamp. Port of the
// private getTimestampById(String).
func (r *SimpleReport) timestampByID(timestampID string) *jaxb.XmlTimestamp {
	if ts, ok := r.tokenByID(timestampID).(*jaxb.XmlTimestamp); ok {
		return ts
	}
	return nil
}

// GetEvidenceRecordById returns a wrapper for the given evidence record.
// Port of getEvidenceRecordById(String).
func (r *SimpleReport) GetEvidenceRecordById(evidenceRecordID string) *jaxb.XmlEvidenceRecord {
	if er, ok := r.tokenByID(evidenceRecordID).(*jaxb.XmlEvidenceRecord); ok {
		return er
	}
	return nil
}

// GetEAAById returns a wrapper for the given EAA. Port of getEAAById(String).
func (r *SimpleReport) GetEAAById(eaaID string) *jaxb.XmlEAA {
	if eaa, ok := r.tokenByID(eaaID).(*jaxb.XmlEAA); ok {
		return eaa
	}
	return nil
}

// GetSignatureTimestamps returns a list of timestamps for a signature with
// the given id. Port of getSignatureTimestamps(String).
func (r *SimpleReport) GetSignatureTimestamps(signatureID string) []*jaxb.XmlTimestamp {
	sig := r.signatureByID(signatureID)
	if sig != nil && sig.Timestamps != nil {
		return sig.Timestamps.Timestamp
	}
	return nil
}

// GetSignatureEvidenceRecords returns a list of evidence records for a
// signature with the given id. Port of getSignatureEvidenceRecords(String).
func (r *SimpleReport) GetSignatureEvidenceRecords(signatureID string) []*jaxb.XmlEvidenceRecord {
	sig := r.signatureByID(signatureID)
	if sig != nil && sig.EvidenceRecords != nil {
		return sig.EvidenceRecords.EvidenceRecord
	}
	return nil
}

// GetTimestampEvidenceRecords returns a list of evidence records for a
// time-stamp with the given id. Port of getTimestampEvidenceRecords(String).
func (r *SimpleReport) GetTimestampEvidenceRecords(timestampID string) []*jaxb.XmlEvidenceRecord {
	ts := r.timestampByID(timestampID)
	if ts != nil && ts.EvidenceRecords != nil {
		return ts.EvidenceRecords.EvidenceRecord
	}
	return nil
}

// GetEvidenceRecordTimestamps returns a list of timestamps for an evidence
// record with the given id. Port of getEvidenceRecordTimestamps(String).
func (r *SimpleReport) GetEvidenceRecordTimestamps(evidenceRecordID string) []*jaxb.XmlTimestamp {
	er := r.GetEvidenceRecordById(evidenceRecordID)
	if er != nil && er.Timestamps != nil {
		return er.Timestamps.Timestamp
	}
	return nil
}

// GetEAASignatures returns a list of signatures used to create the EAA with
// the given Id. NOTE: This method does not return the key binding
// signature. To extract the latest, use GetEAAKeyBindingSignature. Port of
// getEAASignatures(String).
func (r *SimpleReport) GetEAASignatures(eaaID string) []*jaxb.XmlSignature {
	eaa := r.GetEAAById(eaaID)
	if eaa != nil && eaa.EAASignature != nil {
		return eaa.EAASignature
	}
	return nil
}

// GetEAAKeyBindingSignature returns a key binding signature for the EAA,
// when present. Port of getEAAKeyBindingSignature(String).
func (r *SimpleReport) GetEAAKeyBindingSignature(eaaPresentationID string) *jaxb.XmlSignature {
	eaa := r.GetEAAById(eaaPresentationID)
	if eaa != nil {
		return eaa.KeyBindingSignature
	}
	return nil
}

// GetEvidenceRecordPOE returns the lowest POE of an evidence record. Port of
// getEvidenceRecordPOE(String).
func (r *SimpleReport) GetEvidenceRecordPOE(evidenceRecordID string) *time.Time {
	er := r.GetEvidenceRecordById(evidenceRecordID)
	if er != nil {
		return xsTime(er.POETime)
	}
	return nil
}

// GetSignatureScopes returns a list of XmlSignatureScopes for the token
// (signature, timestamp or evidence record) with a given Id. Port of
// getSignatureScopes(String); panics for a token class the port does not
// recognise, mirroring the UnsupportedOperationException Java throws.
func (r *SimpleReport) GetSignatureScopes(tokenID string) []*jaxb.XmlSignatureScope {
	switch t := r.tokenByID(tokenID).(type) {
	case *jaxb.XmlSignature:
		return t.SignatureScope
	case *jaxb.XmlTimestamp:
		return t.TimestampScope
	case *jaxb.XmlEvidenceRecord:
		return t.EvidenceRecordScope
	case nil:
		return nil
	default:
		panic("simplereport: signature scope extraction is not supported for an object of this class")
	}
}

// GetContainerType returns a container type, when applicable (i.e. ASiC
// validation). Port of getContainerType().
func (r *SimpleReport) GetContainerType() enumerations.ASiCContainerType {
	if r.wrapped.ContainerType != nil {
		return r.wrapped.ContainerType.ASiCContainerType()
	}
	return ""
}

// GetPDFAProfile returns a PDF/A Profile name. Port of getPDFAProfile().
func (r *SimpleReport) GetPDFAProfile() string {
	if r.wrapped.PDFAInfo != nil && r.wrapped.PDFAInfo.PDFAProfile != nil {
		return *r.wrapped.PDFAInfo.PDFAProfile
	}
	return ""
}

// IsPDFACompliant returns whether the PDF document is compliant to the
// determined PDF/A specification. Returns false for all non-PDF documents.
// Port of isPDFACompliant(); the RI unboxes PDFAInfo.valid without a null
// check, an upstream NPE risk this port preserves via a nil-pointer
// dereference for the same input shape (a PDFAInfo element present without
// its optional valid attribute) rather than papering over it - see
// jaxb.XmlPDFAInfo's Valid field comment.
func (r *SimpleReport) IsPDFACompliant() bool {
	if r.wrapped.PDFAInfo != nil {
		return *r.wrapped.PDFAInfo.Valid
	}
	return false
}

// GetJaxbModel returns the jaxb model of the simple report. Port of
// getJaxbModel().
func (r *SimpleReport) GetJaxbModel() *jaxb.XmlSimpleReport {
	return r.wrapped
}

// ---------------------------------------------------------- token traversal

// tokenID returns the Id of the given choice-group alternative, or "" for
// nil / an unrecognised type.
func tokenID(item jaxb.XmlTokenItem) string {
	switch t := item.(type) {
	case *jaxb.XmlSignature:
		return t.Id
	case *jaxb.XmlTimestamp:
		return t.Id
	case *jaxb.XmlEvidenceRecord:
		return t.Id
	case *jaxb.XmlEAA:
		return t.Id
	default:
		return ""
	}
}

// tokenContent returns the shared Token content of the given choice-group
// alternative, or nil for nil / an unrecognised type.
func tokenContent(item jaxb.XmlTokenItem) *jaxb.XmlTokenContent {
	switch t := item.(type) {
	case *jaxb.XmlSignature:
		return &t.XmlTokenContent
	case *jaxb.XmlTimestamp:
		return &t.XmlTokenContent
	case *jaxb.XmlEvidenceRecord:
		return &t.XmlTokenContent
	case *jaxb.XmlEAA:
		return &t.XmlTokenContent
	default:
		return nil
	}
}

// embeddedTokenByID searches tokens and their embedded
// timestamps/evidence-records/signatures for a token with the given id.
// Port of the private getEmbeddedTokenById(List, String).
func embeddedTokenByID(tokens []jaxb.XmlTokenItem, tokenIDWanted string) jaxb.XmlTokenItem {
	for _, token := range tokens {
		if tokenID(token) == tokenIDWanted {
			return token
		}
		switch t := token.(type) {
		case *jaxb.XmlSignature:
			if found := signatureTimestampByID(t, tokenIDWanted); found != nil {
				return found
			}
			if found := signatureEvidenceRecordByID(t, tokenIDWanted); found != nil {
				return found
			}
		case *jaxb.XmlTimestamp:
			if found := timestampEvidenceRecordByID(t, tokenIDWanted); found != nil {
				return found
			}
		case *jaxb.XmlEvidenceRecord:
			if found := evidenceRecordTimestampByID(t, tokenIDWanted); found != nil {
				return found
			}
		case *jaxb.XmlEAA:
			if found := eAASignatureByID(t, tokenIDWanted); found != nil {
				return found
			}
		}
	}
	return nil
}

func signatureTimestampByID(sig *jaxb.XmlSignature, tokenIDWanted string) jaxb.XmlTokenItem {
	if sig.Timestamps == nil {
		return nil
	}
	return embeddedTokenByID(timestampsToItems(sig.Timestamps.Timestamp), tokenIDWanted)
}

func signatureEvidenceRecordByID(sig *jaxb.XmlSignature, tokenIDWanted string) jaxb.XmlTokenItem {
	if sig.EvidenceRecords == nil {
		return nil
	}
	return embeddedTokenByID(evidenceRecordsToItems(sig.EvidenceRecords.EvidenceRecord), tokenIDWanted)
}

func timestampEvidenceRecordByID(ts *jaxb.XmlTimestamp, tokenIDWanted string) jaxb.XmlTokenItem {
	if ts.EvidenceRecords == nil {
		return nil
	}
	return embeddedTokenByID(evidenceRecordsToItems(ts.EvidenceRecords.EvidenceRecord), tokenIDWanted)
}

func evidenceRecordTimestampByID(er *jaxb.XmlEvidenceRecord, tokenIDWanted string) jaxb.XmlTokenItem {
	if er.Timestamps == nil {
		return nil
	}
	return embeddedTokenByID(timestampsToItems(er.Timestamps.Timestamp), tokenIDWanted)
}

func eAASignatureByID(eaa *jaxb.XmlEAA, tokenIDWanted string) jaxb.XmlTokenItem {
	if len(eaa.EAASignature) > 0 {
		if found := embeddedTokenByID(signaturesToItems(eaa.EAASignature), tokenIDWanted); found != nil {
			return found
		}
	}
	if eaa.KeyBindingSignature != nil {
		if found := embeddedTokenByID([]jaxb.XmlTokenItem{eaa.KeyBindingSignature}, tokenIDWanted); found != nil {
			return found
		}
	}
	return nil
}

func timestampsToItems(list []*jaxb.XmlTimestamp) []jaxb.XmlTokenItem {
	if len(list) == 0 {
		return nil
	}
	out := make([]jaxb.XmlTokenItem, len(list))
	for i, t := range list {
		out[i] = t
	}
	return out
}

func evidenceRecordsToItems(list []*jaxb.XmlEvidenceRecord) []jaxb.XmlTokenItem {
	if len(list) == 0 {
		return nil
	}
	out := make([]jaxb.XmlTokenItem, len(list))
	for i, t := range list {
		out[i] = t
	}
	return out
}

func signaturesToItems(list []*jaxb.XmlSignature) []jaxb.XmlTokenItem {
	if len(list) == 0 {
		return nil
	}
	out := make([]jaxb.XmlTokenItem, len(list))
	for i, t := range list {
		out[i] = t
	}
	return out
}

// ------------------------------------------------------------------ helpers

// xsTime converts an optional jaxb.XSDateTime field to *time.Time,
// tolerating nil.
func xsTime(d *jaxb.XSDateTime) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time()
	return &t
}

// convertMessages converts a list of XmlMessage into a list of Message.
// Port of the private convert(Collection<XmlMessage>).
func convertMessages(msgs []*jaxb.XmlMessage) []Message {
	out := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		if m == nil {
			continue
		}
		key := ""
		if m.Key != nil {
			key = *m.Key
		}
		out = append(out, Message{Key: key, Value: m.Value})
	}
	return out
}
