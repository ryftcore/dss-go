// Ported from dss-detailed-report-jaxb/src/main/java/eu/europa/esig/dss/detailedreport/DetailedReport.java
// (DSS 6.5.RC1).
//
// DetailedReport wraps the generated jaxb.XmlDetailedReport with the
// block/constraint navigation logic the validation reporting stack drives
// through: locating a token's Basic Building Block, its highest
// completed validation level, its qualification, and the messages attached to
// each. Every public method here is a direct port of the same-named Java
// method; order and null-propagation follow the Java source line for line
// except where Go's nil-pointer-dereference-on-value-receiver rules force an
// explicit nil check Java's plain getter chain did not need (documented at
// each such helper below) - those checks are no-ops on any well-formed
// document, since the fields involved are schema-required in every case they
// guard.
package detailedreport

import (
	"slices"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// DetailedReport represents the detailed report built during the validation
// process. It contains information on each executed constraint. It is
// composed among other of the following building blocks:
//   - Identification of the Signer's Certificate (ISC)
//   - Validation Context Initialization (VCI)
//   - X.509 Certificate Validation (XCV)
//   - Cryptographic Verification (CV)
//   - Signature Acceptance Validation (SAV)
//   - Basic Validation Process
//   - Validation Process for Time-Stamps
//   - Validation Process for AdES-T
//   - Validation of LTV forms
type DetailedReport struct {
	// jaxbDetailedReport is the JAXB Detailed report.
	jaxbDetailedReport *jaxb.XmlDetailedReport

	// messageCollector collects messages of the validation process.
	messageCollector *MessageCollector
}

// NewDetailedReport is the default constructor.
func NewDetailedReport(jaxbDetailedReport *jaxb.XmlDetailedReport) *DetailedReport {
	r := &DetailedReport{jaxbDetailedReport: jaxbDetailedReport}
	// Java creates the collector lazily in getMessageCollector(), an
	// unsynchronised check-then-set. That is harmless on the JVM but is a data
	// race in Go when one report is read by several goroutines, so the
	// collector (which only holds the report pointer) is created up front.
	r.messageCollector = newMessageCollector(r)
	return r
}

// ---------------------------------------------------------------- null helpers
//
// Java's `bbb.getConclusion().getIndication()` never null-checks
// getConclusion(): Conclusion is schema-required, so the chain is safe for
// every valid document and would NPE the same way this port's direct field
// accesses would panic for a malformed one. getSubIndication()/getDateTime()/
// etc. on an *optional* field are safe getters on the Java side purely because
// Java field access never panics; the Go field is a pointer that a
// value-receiver accessor method cannot be called through when nil, so the
// helpers below guard exactly those optional fields, and only those.

func indicationOf(c *jaxb.XmlConclusion) enumerations.Indication {
	if c == nil {
		return ""
	}
	return c.Indication.Indication()
}

func subIndicationOf(c *jaxb.XmlConclusion) enumerations.SubIndication {
	if c == nil || c.SubIndication == nil {
		return ""
	}
	return c.SubIndication.SubIndication()
}

func validationTimeOf(v *jaxb.ValidationTimeValue) enumerations.ValidationTime {
	if v == nil {
		return ""
	}
	return v.ValidationTime()
}

func certificateQualificationOf(v *jaxb.CertificateQualificationValue) enumerations.CertificateQualification {
	if v == nil {
		return ""
	}
	return v.CertificateQualification()
}

// BasicBuildingBlocksIndication returns the result of the Basic Building
// Block for a token (signature, timestamp, revocation).
func (r *DetailedReport) BasicBuildingBlocksIndication(tokenId string) enumerations.Indication {
	bbb := r.BasicBuildingBlockById(tokenId)
	if bbb != nil {
		return indicationOf(bbb.Conclusion)
	}
	return ""
}

// BasicBuildingBlocksSubIndication returns the result of the Basic Building
// Block for a token (signature, timestamp, revocation).
func (r *DetailedReport) BasicBuildingBlocksSubIndication(tokenId string) enumerations.SubIndication {
	bbb := r.BasicBuildingBlockById(tokenId)
	if bbb != nil {
		return subIndicationOf(bbb.Conclusion)
	}
	return ""
}

// BasicBuildingBlocksCertChain returns a list of certificate token ids
// representing the certificate chain of the token in question.
func (r *DetailedReport) BasicBuildingBlocksCertChain(tokenId string) []string {
	certIds := []string{}
	bbb := r.BasicBuildingBlockById(tokenId)
	if bbb != nil && bbb.CertificateChain != nil {
		for _, chainItem := range bbb.CertificateChain.ChainItem {
			certIds = append(certIds, chainItem.Id)
		}
	}
	return certIds
}

// BasicBuildingBlockById returns the full content of the Basic Building Block
// for a token (signature, timestamp, revocation).
func (r *DetailedReport) BasicBuildingBlockById(tokenId string) *jaxb.XmlBasicBuildingBlocks {
	for _, bbb := range r.jaxbDetailedReport.BasicBuildingBlocks {
		if tokenId == bbb.Id {
			return bbb
		}
	}
	return nil
}

// BasicBuildingBlocksNumber returns the number of Basic Building Blocks.
func (r *DetailedReport) BasicBuildingBlocksNumber() int {
	return len(r.jaxbDetailedReport.BasicBuildingBlocks)
}

// BasicBuildingBlocksSignatureId returns the id of the token. The signature is
// identified by its index: 0 for the first one.
func (r *DetailedReport) BasicBuildingBlocksSignatureId(index int) string {
	bbbs := r.jaxbDetailedReport.BasicBuildingBlocks
	if bbbs != nil && len(bbbs) >= index {
		bbb := r.jaxbDetailedReport.BasicBuildingBlocks[index]
		if bbb != nil {
			return bbb.Id
		}
	}
	return ""
}

// SignatureIds returns a list of all signature ids.
func (r *DetailedReport) SignatureIds() []string {
	result := []string{}
	for _, bbb := range r.jaxbDetailedReport.BasicBuildingBlocks {
		t := bbb.Type.Context()
		if enumerations.ContextSignature == t || enumerations.ContextCounterSignature == t || enumerations.ContextKeyBindingSignature == t {
			result = append(result, bbb.Id)
		}
	}
	return result
}

// FirstSignatureId returns the first signature id.
func (r *DetailedReport) FirstSignatureId() string {
	signatureIdList := r.SignatureIds()
	if len(signatureIdList) != 0 {
		return signatureIdList[0]
	}
	return ""
}

// FirstTimestampId returns the first timestamp id.
func (r *DetailedReport) FirstTimestampId() string {
	timestampIdList := r.TimestampIds()
	if len(timestampIdList) != 0 {
		return timestampIdList[0]
	}
	return ""
}

// TimestampIds returns a list of all timestamp ids.
func (r *DetailedReport) TimestampIds() []string {
	result := []string{}
	for _, bbb := range r.jaxbDetailedReport.BasicBuildingBlocks {
		if enumerations.ContextTimestamp == bbb.Type.Context() {
			result = append(result, bbb.Id)
		}
	}
	return result
}

// FirstEvidenceRecordId returns the first evidence record id.
func (r *DetailedReport) FirstEvidenceRecordId() string {
	evidenceRecordIds := r.EvidenceRecordIds()
	if len(evidenceRecordIds) != 0 {
		return evidenceRecordIds[0]
	}
	return ""
}

// EvidenceRecordIds returns a list of all evidence record ids.
//
// NOTE: reproduces a copy-paste artifact in the upstream Java: the final
// loop filters BasicBuildingBlocks by Context.TIMESTAMP, not
// Context.EVIDENCE_RECORD, so it never contributes an id in practice - see
// getTimestampIds for the loop it appears to have been copied from.
func (r *DetailedReport) EvidenceRecordIds() []string {
	result := []string{}
	for _, token := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		switch t := token.(type) {
		case *jaxb.XmlEvidenceRecord:
			result = append(result, derefString(t.Id))
		case *jaxb.XmlSignature:
			for _, er := range t.EvidenceRecord {
				result = append(result, derefString(er.Id))
			}
		}
	}
	for _, bbb := range r.jaxbDetailedReport.BasicBuildingBlocks {
		if enumerations.ContextTimestamp == bbb.Type.Context() {
			result = append(result, bbb.Id)
		}
	}
	return result
}

// EAAIds returns a list of all EAA presentation ids.
func (r *DetailedReport) EAAIds() []string {
	result := []string{}
	for _, token := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		if eaa, ok := token.(*jaxb.XmlEAA); ok {
			result = append(result, derefString(eaa.Id))
		}
	}
	return result
}

// FirstEAAId returns the first EAA presentation id.
func (r *DetailedReport) FirstEAAId() string {
	eaaIds := r.EAAIds()
	if len(eaaIds) != 0 {
		return eaaIds[0]
	}
	return ""
}

// RevocationIds returns a list of all revocation data ids.
func (r *DetailedReport) RevocationIds() []string {
	result := []string{}
	for _, bbb := range r.jaxbDetailedReport.BasicBuildingBlocks {
		if enumerations.ContextRevocation == bbb.Type.Context() {
			result = append(result, bbb.Id)
		}
	}
	return result
}

// BestSignatureTime returns the best-signature-time for the signature with id.
func (r *DetailedReport) BestSignatureTime(signatureId string) *time.Time {
	poe := r.BestProofOfExistence(signatureId)
	if poe != nil {
		t := poe.Time.Time()
		return &t
	}
	return nil
}

// BestProofOfExistence gets the best proof-of-existence for the signature with id.
func (r *DetailedReport) BestProofOfExistence(signatureId string) *jaxb.XmlProofOfExistence {
	xmlSignature := r.XmlSignatureById(signatureId)
	if xmlSignature != nil {
		if xmlSignature.ValidationProcessArchivalData != nil {
			return xmlSignature.ValidationProcessArchivalData.ProofOfExistence
		}
		if xmlSignature.ValidationProcessLongTermData != nil {
			return xmlSignature.ValidationProcessLongTermData.ProofOfExistence
		}
		if xmlSignature.ValidationProcessBasicSignature != nil {
			return xmlSignature.ValidationProcessBasicSignature.ProofOfExistence
		}
	}
	return nil
}

// EvidenceRecordLowestPOETime returns the lowest POE of the evidence record
// with the given Id.
func (r *DetailedReport) EvidenceRecordLowestPOETime(evidenceRecordId string) *time.Time {
	xmlEvidenceRecord := r.XmlEvidenceRecordById(evidenceRecordId)
	if xmlEvidenceRecord != nil && xmlEvidenceRecord.ValidationProcessEvidenceRecord != nil {
		poe := xmlEvidenceRecord.ValidationProcessEvidenceRecord.ProofOfExistence
		if poe != nil {
			t := poe.Time.Time()
			return &t
		}
	}
	return nil
}

// BasicValidationIndication gets the basic validation indication for a
// signature with id.
func (r *DetailedReport) BasicValidationIndication(signatureId string) enumerations.Indication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessBasicSignature != nil && signature.ValidationProcessBasicSignature.Conclusion != nil {
		return signature.ValidationProcessBasicSignature.Conclusion.Indication.Indication()
	}
	return ""
}

// BasicValidationSubIndication gets the basic validation subIndication for a
// signature with id.
func (r *DetailedReport) BasicValidationSubIndication(signatureId string) enumerations.SubIndication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessBasicSignature != nil && signature.ValidationProcessBasicSignature.Conclusion != nil {
		return subIndicationOf(signature.ValidationProcessBasicSignature.Conclusion)
	}
	return ""
}

// BasicTimestampValidationIndication gets the timestamp basic validation
// indication for a timestamp with id.
func (r *DetailedReport) BasicTimestampValidationIndication(timestampId string) enumerations.Indication {
	timestampValidationById := r.basicTimestampValidationById(timestampId)
	if timestampValidationById != nil && timestampValidationById.Conclusion != nil {
		return timestampValidationById.Conclusion.Indication.Indication()
	}
	return ""
}

// BasicTimestampValidationSubIndication gets the timestamp basic validation
// subIndication for a timestamp with id.
func (r *DetailedReport) BasicTimestampValidationSubIndication(timestampId string) enumerations.SubIndication {
	timestampValidationById := r.basicTimestampValidationById(timestampId)
	if timestampValidationById != nil && timestampValidationById.Conclusion != nil {
		return subIndicationOf(timestampValidationById.Conclusion)
	}
	return ""
}

// ArchiveDataTimestampValidationIndication gets the timestamp validation with
// archive data indication for a timestamp with id.
func (r *DetailedReport) ArchiveDataTimestampValidationIndication(timestampId string) enumerations.Indication {
	timestampValidationById := r.archiveDataTimestampValidationById(timestampId)
	if timestampValidationById != nil && timestampValidationById.Conclusion != nil {
		return timestampValidationById.Conclusion.Indication.Indication()
	}
	return ""
}

// ArchiveDataTimestampValidationSubIndication gets the timestamp validation
// with archive data subIndication for a timestamp with id.
func (r *DetailedReport) ArchiveDataTimestampValidationSubIndication(timestampId string) enumerations.SubIndication {
	timestampValidationById := r.archiveDataTimestampValidationById(timestampId)
	if timestampValidationById != nil && timestampValidationById.Conclusion != nil {
		return subIndicationOf(timestampValidationById.Conclusion)
	}
	return ""
}

// EvidenceRecordValidationIndication gets the evidence record validation
// indication for an evidence record with id.
func (r *DetailedReport) EvidenceRecordValidationIndication(evidenceRecordId string) enumerations.Indication {
	evidenceRecordValidationById := r.evidenceRecordValidationById(evidenceRecordId)
	if evidenceRecordValidationById != nil && evidenceRecordValidationById.Conclusion != nil {
		return evidenceRecordValidationById.Conclusion.Indication.Indication()
	}
	return ""
}

// EvidenceRecordValidationSubIndication gets the evidence record validation
// subIndication for an evidence record with id.
func (r *DetailedReport) EvidenceRecordValidationSubIndication(evidenceRecordId string) enumerations.SubIndication {
	evidenceRecordValidationById := r.evidenceRecordValidationById(evidenceRecordId)
	if evidenceRecordValidationById != nil && evidenceRecordValidationById.Conclusion != nil {
		return subIndicationOf(evidenceRecordValidationById.Conclusion)
	}
	return ""
}

func (r *DetailedReport) evidenceRecordValidationById(evidenceRecordId string) *jaxb.XmlValidationProcessEvidenceRecord {
	evidenceRecord := r.XmlEvidenceRecordById(evidenceRecordId)
	if evidenceRecord != nil {
		return evidenceRecord.ValidationProcessEvidenceRecord
	}
	return nil
}

// XmlEvidenceRecordById returns an XmlEvidenceRecord by the given id. Nil if
// the evidence record is not found.
func (r *DetailedReport) XmlEvidenceRecordById(evidenceRecordId string) *jaxb.XmlEvidenceRecord {
	for _, er := range r.IndependentEvidenceRecords() {
		if derefString(er.Id) == evidenceRecordId {
			return er
		}
	}
	for _, sig := range r.Signatures() {
		for _, er := range sig.EvidenceRecord {
			if derefString(er.Id) == evidenceRecordId {
				return er
			}
		}
	}
	for _, tst := range r.IndependentTimestamps() {
		for _, er := range tst.EvidenceRecord {
			if derefString(er.Id) == evidenceRecordId {
				return er
			}
		}
	}
	return nil
}

// EAAValidationIndication gets the EAA presentation validation indication for
// an EAA presentation with id.
func (r *DetailedReport) EAAValidationIndication(eaaId string) enumerations.Indication {
	eaaValidationById := r.eaaValidationById(eaaId)
	if eaaValidationById != nil && eaaValidationById.Conclusion != nil {
		return eaaValidationById.Conclusion.Indication.Indication()
	}
	return ""
}

// EAAValidationSubIndication gets the EAA presentation validation
// subIndication for an EAA presentation with id.
func (r *DetailedReport) EAAValidationSubIndication(eaaId string) enumerations.SubIndication {
	eaaValidationById := r.eaaValidationById(eaaId)
	if eaaValidationById != nil && eaaValidationById.Conclusion != nil {
		return subIndicationOf(eaaValidationById.Conclusion)
	}
	return ""
}

func (r *DetailedReport) eaaValidationById(eaaId string) *jaxb.XmlValidationProcessEAA {
	eaa := r.XmlEAAById(eaaId)
	if eaa != nil {
		return eaa.ValidationProcessEAA
	}
	return nil
}

// XmlEAAById returns an XmlEAA by the given id. Nil if the EAA is not found.
func (r *DetailedReport) XmlEAAById(eaaId string) *jaxb.XmlEAA {
	for _, eaa := range r.EAAs() {
		if derefString(eaa.Id) == eaaId {
			return eaa
		}
	}
	return nil
}

// LongTermValidationIndication gets the long-term validation indication for a
// signature with id.
func (r *DetailedReport) LongTermValidationIndication(signatureId string) enumerations.Indication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessLongTermData != nil && signature.ValidationProcessLongTermData.Conclusion != nil {
		return signature.ValidationProcessLongTermData.Conclusion.Indication.Indication()
	}
	return ""
}

// LongTermValidationSubIndication gets the long-term validation subIndication
// for a signature with id.
func (r *DetailedReport) LongTermValidationSubIndication(signatureId string) enumerations.SubIndication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessLongTermData != nil && signature.ValidationProcessLongTermData.Conclusion != nil {
		return subIndicationOf(signature.ValidationProcessLongTermData.Conclusion)
	}
	return ""
}

// ArchiveDataValidationIndication gets the validation with archive data
// indication for a signature with id.
func (r *DetailedReport) ArchiveDataValidationIndication(signatureId string) enumerations.Indication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessArchivalData != nil && signature.ValidationProcessArchivalData.Conclusion != nil {
		return signature.ValidationProcessArchivalData.Conclusion.Indication.Indication()
	}
	return ""
}

// ArchiveDataValidationSubIndication gets the validation with archive data
// subIndication for a signature with id.
func (r *DetailedReport) ArchiveDataValidationSubIndication(signatureId string) enumerations.SubIndication {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationProcessArchivalData != nil && signature.ValidationProcessArchivalData.Conclusion != nil {
		return subIndicationOf(signature.ValidationProcessArchivalData.Conclusion)
	}
	return ""
}

// SignatureQualification gets the qualification for a signature with id.
func (r *DetailedReport) SignatureQualification(signatureId string) enumerations.SignatureQualification {
	signature := r.XmlSignatureById(signatureId)
	if signature != nil && signature.ValidationSignatureQualification != nil {
		return signature.ValidationSignatureQualification.SignatureQualification.SignatureQualification()
	}
	return ""
}

// TimestampQualification gets the final qualification result for a timestamp
// with id.
func (r *DetailedReport) TimestampQualification(timestampId string) enumerations.TimestampQualification {
	timestampQualif := r.xmlTimestampQualificationById(timestampId)
	if timestampQualif != nil {
		return timestampQualif.TimestampQualification.TimestampQualification()
	}
	return ""
}

// TimestampQualificationAtTstGenerationTime gets the qualification for a
// timestamp with the given id at the timestamp generation time.
func (r *DetailedReport) TimestampQualificationAtTstGenerationTime(timestampId string) enumerations.TimestampQualification {
	return r.timestampQualificationAtValidationTime(enumerations.ValidationTimeTimestampGenerationTime, timestampId)
}

// TimestampQualificationAtBestPoeTime gets the qualification for a timestamp
// with the given id at its best available POE time.
func (r *DetailedReport) TimestampQualificationAtBestPoeTime(timestampId string) enumerations.TimestampQualification {
	return r.timestampQualificationAtValidationTime(enumerations.ValidationTimeTimestampPOETime, timestampId)
}

func (r *DetailedReport) timestampQualificationAtValidationTime(validationTime enumerations.ValidationTime, timestampId string) enumerations.TimestampQualification {
	tstQualificationValidation := r.xmlTimestampQualificationById(timestampId)
	if tstQualificationValidation != nil {
		for _, tstQualAtTime := range tstQualificationValidation.ValidationTimestampQualificationAtTime {
			if validationTime == validationTimeOf(tstQualAtTime.ValidationTime) {
				return tstQualAtTime.TimestampQualification.TimestampQualification()
			}
		}
	}
	return ""
}

func (r *DetailedReport) xmlTimestampQualificationById(timestampId string) *jaxb.XmlValidationTimestampQualification {
	timestamp := r.XmlTimestampById(timestampId)
	if timestamp != nil {
		return timestamp.ValidationTimestampQualification
	}
	return nil
}

func (r *DetailedReport) basicTimestampValidationById(timestampId string) *jaxb.XmlValidationProcessBasicTimestamp {
	timestamp := r.XmlTimestampById(timestampId)
	if timestamp != nil {
		return timestamp.ValidationProcessBasicTimestamp
	}
	return nil
}

func (r *DetailedReport) archiveDataTimestampValidationById(timestampId string) *jaxb.XmlValidationProcessArchivalDataTimestamp {
	timestamp := r.XmlTimestampById(timestampId)
	if timestamp != nil {
		return timestamp.ValidationProcessArchivalDataTimestamp
	}
	return nil
}

// EAAQualifications gets the final qualification result for an EAA
// presentation with id.
func (r *DetailedReport) EAAQualifications(eaaId string) []enumerations.EAAQualification {
	eaaQualification := r.xmlEAAQualificationById(eaaId)
	if eaaQualification != nil {
		out := make([]enumerations.EAAQualification, 0, len(eaaQualification.EAAQualification))
		for _, q := range eaaQualification.EAAQualification {
			out = append(out, q.EAAQualification())
		}
		return out
	}
	return nil
}

func (r *DetailedReport) xmlEAAQualificationById(eaaId string) *jaxb.XmlValidationEAAQualification {
	eaa := r.XmlEAAById(eaaId)
	if eaa != nil {
		return eaa.ValidationEAAQualification
	}
	return nil
}

// XmlTimestampById returns an XmlTimestamp by the given id. Nil if the
// timestamp is not found.
func (r *DetailedReport) XmlTimestampById(timestampId string) *jaxb.XmlTimestamp {
	for _, xmlTimestamp := range r.IndependentTimestamps() {
		if derefString(xmlTimestamp.Id) == timestampId {
			return xmlTimestamp
		}
		for _, er := range xmlTimestamp.EvidenceRecord {
			for _, ert := range er.Timestamp {
				if derefString(ert.Id) == timestampId {
					return ert
				}
			}
		}
	}

	for _, xmlSignature := range r.Signatures() {
		for _, xmlTimestamp := range xmlSignature.Timestamp {
			if derefString(xmlTimestamp.Id) == timestampId {
				return xmlTimestamp
			}
		}
		for _, er := range xmlSignature.EvidenceRecord {
			for _, xmlTimestamp := range er.Timestamp {
				if derefString(xmlTimestamp.Id) == timestampId {
					return xmlTimestamp
				}
			}
		}
	}

	for _, er := range r.IndependentEvidenceRecords() {
		for _, xmlTimestamp := range er.Timestamp {
			if derefString(xmlTimestamp.Id) == timestampId {
				return xmlTimestamp
			}
		}
	}
	return nil
}

// XmlSignatureById returns an XmlSignature by the given id. Nil if the
// signature is not found.
func (r *DetailedReport) XmlSignatureById(signatureId string) *jaxb.XmlSignature {
	for _, xmlSignature := range r.Signatures() {
		if signatureId == derefString(xmlSignature.Id) {
			return xmlSignature
		}
	}
	return nil
}

// XmlCertificateById returns an XmlCertificate by id if it exists, nil
// otherwise. NOTE: should be used only for certificate validation process.
func (r *DetailedReport) XmlCertificateById(certificateId string) *jaxb.XmlCertificate {
	for _, xmlCertificate := range r.Certificates() {
		if certificateId == derefString(xmlCertificate.Id) {
			return xmlCertificate
		}
	}
	return nil
}

// Signatures returns a list of all signatures.
func (r *DetailedReport) Signatures() []*jaxb.XmlSignature {
	result := []*jaxb.XmlSignature{}
	for _, element := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		switch t := element.(type) {
		case *jaxb.XmlSignature:
			result = append(result, t)
		case *jaxb.XmlEAA:
			result = append(result, t.Signature...)
			if t.KeyBindingSignature != nil {
				result = append(result, t.KeyBindingSignature)
			}
		}
	}
	return result
}

// IndependentTimestamps returns a list of all independent (detached)
// timestamps.
func (r *DetailedReport) IndependentTimestamps() []*jaxb.XmlTimestamp {
	result := []*jaxb.XmlTimestamp{}
	for _, element := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		if t, ok := element.(*jaxb.XmlTimestamp); ok {
			result = append(result, t)
		}
	}
	return result
}

// IndependentEvidenceRecords returns a list of all independent (detached)
// evidence records.
func (r *DetailedReport) IndependentEvidenceRecords() []*jaxb.XmlEvidenceRecord {
	result := []*jaxb.XmlEvidenceRecord{}
	for _, element := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		if t, ok := element.(*jaxb.XmlEvidenceRecord); ok {
			result = append(result, t)
		}
	}
	return result
}

// EAAs returns a list of all EAA presentations.
func (r *DetailedReport) EAAs() []*jaxb.XmlEAA {
	result := []*jaxb.XmlEAA{}
	for _, element := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		if t, ok := element.(*jaxb.XmlEAA); ok {
			result = append(result, t)
		}
	}
	return result
}

// Certificates returns a list of processed XmlCertificates. NOTE: the method
// returns a non-empty list only for a certificate validation process.
func (r *DetailedReport) Certificates() []*jaxb.XmlCertificate {
	result := []*jaxb.XmlCertificate{}
	for _, element := range r.jaxbDetailedReport.SignatureOrTimestampOrEvidenceRecord {
		if t, ok := element.(*jaxb.XmlCertificate); ok {
			result = append(result, t)
		}
	}
	return result
}

// TLAnalysisById returns a complete block of a TL validation.
func (r *DetailedReport) TLAnalysisById(tlId string) *jaxb.XmlTLAnalysis {
	for _, xmlTLAnalysis := range r.jaxbDetailedReport.TLAnalysis {
		if tlId == xmlTLAnalysis.Id {
			return xmlTLAnalysis
		}
	}
	return nil
}

// JAXBModel returns the JAXB Detailed Report.
func (r *DetailedReport) JAXBModel() *jaxb.XmlDetailedReport {
	return r.jaxbDetailedReport
}

// IsCertificateValidation returns whether the certificate validation has been
// performed (therefore the certificate corresponding data can be retrieved).
func (r *DetailedReport) IsCertificateValidation() bool {
	certificates := r.Certificates()
	return len(certificates) != 0
}

// CertificateQualificationAtIssuance gets the qualification for certificate
// with id at its issuance time.
func (r *DetailedReport) CertificateQualificationAtIssuance(certificateId string) enumerations.CertificateQualification {
	return r.certificateQualificationAtTime(enumerations.ValidationTimeCertificateIssuanceTime, certificateId)
}

// CertificateQualificationAtValidation gets the qualification for certificate
// with id at the validation time.
func (r *DetailedReport) CertificateQualificationAtValidation(certificateId string) enumerations.CertificateQualification {
	return r.certificateQualificationAtTime(enumerations.ValidationTimeValidationTime, certificateId)
}

func (r *DetailedReport) certificateQualificationAtTime(validationTime enumerations.ValidationTime, certificateId string) enumerations.CertificateQualification {
	if certificateId == "" {
		return enumerations.CertificateQualificationNA
	}

	certificate := r.XmlCertificateById(certificateId)
	if certificate != nil {
		certificateQualificationProcess := certificate.CertificateQualificationProcess
		if certificateQualificationProcess != nil {
			for _, vcq := range certificateQualificationProcess.ValidationCertificateQualification {
				if validationTime == validationTimeOf(vcq.ValidationTime) {
					return certificateQualificationOf(vcq.CertificateQualification)
				}
			}
		}
	} else {
		signatures := r.Signatures()
		for _, xmlSignature := range signatures {
			signatureQualification := xmlSignature.ValidationSignatureQualification
			if signatureQualification != nil && signatureQualification.ValidationCertificateQualification != nil {
				for _, cq := range signatureQualification.ValidationCertificateQualification {
					if certificateId == cq.Id && validationTime == validationTimeOf(cq.ValidationTime) {
						return certificateQualificationOf(cq.CertificateQualification)
					}
				}
			}
		}
	}

	return enumerations.CertificateQualificationNA
}

// CertificateQWACProfile gets the QWAC Profile of the given certificate, if
// the validation has been performed.
// NOTE: applicable only on QWAC validation.
func (r *DetailedReport) CertificateQWACProfile(certificateId string) enumerations.QWACProfile {
	if certificateId == "" {
		return ""
	}

	var qwacProcess *jaxb.XmlQWACProcess

	xmlCertificate := r.XmlCertificateById(certificateId)
	if xmlCertificate != nil {
		qwacProcess = xmlCertificate.QWACProcess
	} else {
		signatures := r.Signatures()
		for _, xmlSignature := range signatures {
			signatureQualification := xmlSignature.ValidationSignatureQualification
			if signatureQualification != nil {
				if signatureQualification.QWACProcess != nil && certificateId == signatureQualification.QWACProcess.Id {
					qwacProcess = signatureQualification.QWACProcess
				}
			}
		}
	}

	if qwacProcess != nil && qwacProcess.QWACType != nil {
		return qwacProcess.QWACType.QWACProfile()
	}
	return ""
}

// CertificateApprovalStatussAtIssuanceTime gets certificate approval statuses
// obtained on TS 119 602 List(s) of Trusted Entities processing for the
// certificate with the given identifier at the certificate issuance time.
func (r *DetailedReport) CertificateApprovalStatussAtIssuanceTime(certificateId string) []enumerations.CertificateApprovalStatus {
	return r.certificateApprovalStatussAtTime(certificateId, enumerations.ValidationTimeCertificateIssuanceTime)
}

// CertificateApprovalStatussAtValidationTime gets certificate approval
// statuses obtained on TS 119 602 List(s) of Trusted Entities processing for
// the certificate with the given identifier at the certificate validation
// time.
func (r *DetailedReport) CertificateApprovalStatussAtValidationTime(certificateId string) []enumerations.CertificateApprovalStatus {
	return r.certificateApprovalStatussAtTime(certificateId, enumerations.ValidationTimeValidationTime)
}

func (r *DetailedReport) certificateApprovalStatussAtTime(certificateId string, validationTime enumerations.ValidationTime) []enumerations.CertificateApprovalStatus {
	if certificateId == "" {
		return []enumerations.CertificateApprovalStatus{}
	}

	result := []enumerations.CertificateApprovalStatus{}

	certificate := r.XmlCertificateById(certificateId)
	if certificate != nil {
		certificateApprovalStatusProcess := certificate.CertificateApprovalStatusProcess
		if certificateApprovalStatusProcess != nil {
			for _, vcas := range certificateApprovalStatusProcess.ValidationCertificateApprovalStatus {
				if validationTime == validationTimeOf(vcas.ValidationTime) {
					result = append(result, buildFromXmlCertificateApprovalStatus(vcas.CertificateApprovalStatus))
				}
			}
		}
	}

	return result
}

func buildFromXmlCertificateApprovalStatus(xmlCertificateApprovalStatus *jaxb.XmlCertificateApprovalStatus) enumerations.CertificateApprovalStatus {
	if xmlCertificateApprovalStatus == nil {
		return nil
	}
	result := enumerations.CertificateApprovalStatusFromDefinition(xmlCertificateApprovalStatus.ListType,
		xmlCertificateApprovalStatus.ServiceTypeIdentifier, xmlCertificateApprovalStatus.ServiceStatus)
	if result != nil && result.Label() != "" &&
		result != enumerations.CertificateApprovalStatus(enumerations.CertificateApprovalStatusEnumCertForUnknown) {
		return result
	}
	return enumerations.NewCertificateApprovalStatus(enumerations.CertificateApprovalStatusEnumCertForUnknown.Label(),
		xmlCertificateApprovalStatus.ListType, xmlCertificateApprovalStatus.ServiceTypeIdentifier, xmlCertificateApprovalStatus.ServiceStatus)
}

// CertificateXCVConclusion gets the XCV building block conclusion for a
// certificate with id.
func (r *DetailedReport) CertificateXCVConclusion(certificateId string) *jaxb.XmlConclusion {
	certificates := r.Certificates()
	if len(certificates) == 0 {
		panic("Only supported in report for certificate")
	}

	// process cert chain for signatures
	signatureIds := r.SignatureIds()
	basicBuildingBlocks := r.jaxbDetailedReport.BasicBuildingBlocks
	for _, xmlBasicBuildingBlocks := range basicBuildingBlocks {
		if !slices.Contains(signatureIds, xmlBasicBuildingBlocks.Id) {
			continue // skip for signature
		}

		xcv := xmlBasicBuildingBlocks.XCV
		if xcv != nil {
			for _, xmlSubXCV := range xcv.SubXCV {
				if certificateId == xmlSubXCV.Id {
					return xmlSubXCV.Conclusion
				}
			}
		}
	}

	// process other certificates (certificate validation only)
	for _, xmlBasicBuildingBlocks := range basicBuildingBlocks {
		if slices.Contains(signatureIds, xmlBasicBuildingBlocks.Id) {
			continue // skip for signature
		}

		xcv := xmlBasicBuildingBlocks.XCV
		if xcv != nil {
			for _, xmlSubXCV := range xcv.SubXCV {
				if certificateId == xmlSubXCV.Id {
					return xmlSubXCV.Conclusion
				}
			}
		}
	}
	return nil
}

// FinalConclusion returns the final validation conclusion for a token with a
// given Id.
func (r *DetailedReport) FinalConclusion(tokenId string) *jaxb.XmlConclusion {
	if signatureById := r.XmlSignatureById(tokenId); signatureById != nil {
		return signatureById.Conclusion
	}
	if timestampById := r.XmlTimestampById(tokenId); timestampById != nil {
		return timestampById.Conclusion
	}
	if evidenceRecordById := r.XmlEvidenceRecordById(tokenId); evidenceRecordById != nil {
		return evidenceRecordById.Conclusion
	}
	if eaaById := r.XmlEAAById(tokenId); eaaById != nil {
		return eaaById.Conclusion
	}
	if bbb := r.BasicBuildingBlockById(tokenId); bbb != nil {
		return bbb.Conclusion
	}
	return nil
}

// FinalIndication gets the validation indication to a token with id
// corresponding to the highest validation level.
func (r *DetailedReport) FinalIndication(tokenId string) enumerations.Indication {
	return indicationOf(r.FinalConclusion(tokenId))
}

// FinalSubIndication gets the validation subIndication to a token with id
// corresponding to the highest validation level.
func (r *DetailedReport) FinalSubIndication(tokenId string) enumerations.SubIndication {
	return subIndicationOf(r.FinalConclusion(tokenId))
}

// HighestConclusion gets the validation conclusion to a signature with id
// corresponding to the highest validation level. Java returns a reference
// typed to the XmlConstraintsConclusion superclass shared by
// ValidationProcessArchivalData/LongTermData/BasicSignature; Go has no
// upcast between the three distinct generated struct types (see
// jaxb_process.go's Content/Attrs embedding), but every caller only ever
// reads .Conclusion off the result, so returning the shared
// XmlConstraintsConclusionContent they all embed serves the same purpose.
//
// Like Java, it answers nil for a signature that carries none of the three
// blocks (Java's last branch returns getValidationProcessBasicSignature(),
// which is null then) instead of dereferencing the absent basic block; the
// message collector tolerates that nil the way Java's getMessages does.
func (r *DetailedReport) HighestConclusion(signatureId string) *jaxb.XmlConstraintsConclusionContent {
	xmlSignature := r.XmlSignatureById(signatureId)
	if xmlSignature.ValidationProcessArchivalData != nil {
		return &xmlSignature.ValidationProcessArchivalData.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent
	} else if xmlSignature.ValidationProcessLongTermData != nil {
		return &xmlSignature.ValidationProcessLongTermData.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent
	} else if xmlSignature.ValidationProcessBasicSignature != nil {
		return &xmlSignature.ValidationProcessBasicSignature.XmlConstraintsConclusionWithProofOfExistenceContent.XmlConstraintsConclusionContent
	}
	return nil
}

// SigningCertificate gets the signing certificate validation block for the
// given BasicBuildingBlock.
func (r *DetailedReport) SigningCertificate(bbbId string) *jaxb.XmlSubXCV {
	basicBuildingBlocks := r.BasicBuildingBlockById(bbbId)
	if basicBuildingBlocks != nil {
		xcv := basicBuildingBlocks.XCV
		if xcv != nil {
			subXCVs := xcv.SubXCV
			if len(subXCVs) != 0 {
				return subXCVs[0]
			}
		}
	}
	return nil
}

// MessageCollector gets the used MessageCollector. NewDetailedReport creates it
// eagerly, so concurrent readers never race on a lazy initialisation; the
// fallback only serves a zero-value DetailedReport.
func (r *DetailedReport) MessageCollector() *MessageCollector {
	if r.messageCollector == nil {
		r.messageCollector = newMessageCollector(r)
	}
	return r.messageCollector
}

// AdESValidationErrors returns a list of ETSI EN 319 102-1 AdES validation
// error messages for a token with the given id.
func (r *DetailedReport) AdESValidationErrors(tokenId string) []Message {
	return r.MessageCollector().AdESValidationErrors(tokenId)
}

// AdESValidationWarnings returns a list of ETSI EN 319 102-1 AdES validation
// warning messages for a token with the given id.
func (r *DetailedReport) AdESValidationWarnings(tokenId string) []Message {
	return r.MessageCollector().AdESValidationWarnings(tokenId)
}

// AdESValidationInfos returns a list of ETSI EN 319 102-1 AdES validation info
// messages for a token with the given id.
func (r *DetailedReport) AdESValidationInfos(tokenId string) []Message {
	return r.MessageCollector().AdESValidationInfos(tokenId)
}

// QualificationErrors returns a list of qualification validation errors for a
// token with the given id.
func (r *DetailedReport) QualificationErrors(tokenId string) []Message {
	return r.MessageCollector().QualificationErrors(tokenId)
}

// QualificationWarnings returns a list of qualification validation warnings
// for a token with the given id.
func (r *DetailedReport) QualificationWarnings(tokenId string) []Message {
	return r.MessageCollector().QualificationWarnings(tokenId)
}

// QualificationInfos returns a list of qualification validation infos for a
// token with the given id.
func (r *DetailedReport) QualificationInfos(tokenId string) []Message {
	return r.MessageCollector().QualificationInfos(tokenId)
}

// CertificateQualificationErrorsAtIssuanceTime returns a list of
// qualification validation errors for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationErrorsAtIssuanceTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationErrorsAtIssuanceTime(certificateId)
}

// CertificateQualificationWarningsAtIssuanceTime returns a list of
// qualification validation warnings for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationWarningsAtIssuanceTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationWarningsAtIssuanceTime(certificateId)
}

// CertificateQualificationInfosAtIssuanceTime returns a list of qualification
// validation information messages for a certificate with the given id at
// certificate issuance time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationInfosAtIssuanceTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationInfosAtIssuanceTime(certificateId)
}

// CertificateQualificationErrorsAtValidationTime returns a list of
// qualification validation errors for a certificate with the given id at
// validation time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationErrorsAtValidationTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationErrorsAtValidationTime(certificateId)
}

// CertificateQualificationWarningsAtValidationTime returns a list of
// qualification validation warnings for a certificate with the given id at
// validation time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationWarningsAtValidationTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationWarningsAtValidationTime(certificateId)
}

// CertificateQualificationInfosAtValidationTime returns a list of
// qualification validation information messages for a certificate with the
// given id at validation time.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateQualificationInfosAtValidationTime(certificateId string) []Message {
	return r.MessageCollector().CertificateQualificationInfosAtValidationTime(certificateId)
}

// QWACValidationErrors returns a list of QWAC validation errors for a
// certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (r *DetailedReport) QWACValidationErrors(certificateId string) []Message {
	return r.MessageCollector().QWACValidationErrors(certificateId)
}

// QWACValidationWarnings returns a list of QWAC validation warnings for a
// certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (r *DetailedReport) QWACValidationWarnings(certificateId string) []Message {
	return r.MessageCollector().QWACValidationWarnings(certificateId)
}

// QWACValidationInfos returns a list of QWAC validation information messages
// for a certificate with the given id at certificate issuance time.
// NOTE: applicable only on QWAC validation.
func (r *DetailedReport) QWACValidationInfos(certificateId string) []Message {
	return r.MessageCollector().QWACValidationInfos(certificateId)
}

// CertificateApprovalStatusErrorsAtIssuanceTime returns a list of
// qualification validation errors for a certificate with the given id at
// certificate issuance time for the given certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusErrorsAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusErrorsAtIssuanceTime(certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusWarningsAtIssuanceTime returns a list of
// qualification validation warnings for a certificate with the given id at
// certificate issuance time for the given certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusWarningsAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusWarningsAtIssuanceTime(certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusInfosAtIssuanceTime returns a list of
// qualification validation information messages for a certificate with the
// given id at certificate issuance time for the given
// certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusInfosAtIssuanceTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusInfosAtIssuanceTime(certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusErrorsAtValidationTime returns a list of
// qualification validation errors for a certificate with the given id at
// validation time for the given certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusErrorsAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusErrorsAtValidationTime(certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusWarningsAtValidationTime returns a list of
// qualification validation warnings for a certificate with the given id at
// validation time for the given certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusWarningsAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusWarningsAtValidationTime(certificateId, certificateApprovalStatus)
}

// CertificateApprovalStatusInfosAtValidationTime returns a list of
// qualification validation information messages for a certificate with the
// given id at validation time for the given certificateApprovalStatus.
// NOTE: applicable only on certificate validation.
func (r *DetailedReport) CertificateApprovalStatusInfosAtValidationTime(certificateId string, certificateApprovalStatus enumerations.CertificateApprovalStatus) []Message {
	return r.MessageCollector().CertificateApprovalStatusInfosAtValidationTime(certificateId, certificateApprovalStatus)
}

// -------------------------------------------------------------- small helpers

// derefString returns *s, or "" for a nil s, matching Java's ability to
// compare/hold a null String where this port holds an optional *string.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
