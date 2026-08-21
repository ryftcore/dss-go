// Ported from dss-document/src/main/java/eu/europa/esig/dss/evidencerecord/EvidenceRecordIncorporationService.java (DSS 6.5.RC1).
package document

import "github.com/ryftcore/dss-go/dss/model"

// EvidenceRecordIncorporationService provides common methods for incorporation of evidence
// records within existing signatures, generic over the ERP implementation of format related
// parameters for evidence record incorporation.
type EvidenceRecordIncorporationService[ERP model.SerializableEvidenceRecordIncorporationParameters] interface {
	// AddSignatureEvidenceRecord incorporates the Evidence Record as an unsigned property into
	// the signature. Port of #addSignatureEvidenceRecord.
	AddSignatureEvidenceRecord(signatureDocument, evidenceRecordDocument model.DSSDocument, parameters ERP) model.DSSDocument
}
