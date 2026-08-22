// Ported from dss-document/src/main/java/eu/europa/esig/dss/evidencerecord/EvidenceRecordIncorporationService.java (DSS 6.5.RC1).
//
// dss-evidence-record is a deferred module not yet ported. AbstractASiCSignatureService
// implements this interface upstream, so a minimal local copy is ported here into package asic -
// method surface exactly as the Java interface declares (Serializable dropped per PORTING.md) -
// rather than gating the whole AbstractASiCSignatureService file on the deferred module. Replace
// this local copy with an import of the real type once dss/evidencerecord is ported, and revisit
// any type assertions against this local interface.
package asic

import "github.com/ryftcore/dss-go/dss/model"

// EvidenceRecordIncorporationService provides common methods for incorporation of evidence
// records within existing signatures, generic over the ERP implementation of format-related
// parameters for evidence record incorporation.
type EvidenceRecordIncorporationService[ERP model.SerializableEvidenceRecordIncorporationParameters] interface {
	// AddSignatureEvidenceRecord incorporates the Evidence Record as an unsigned property into
	// the signature. Port of addSignatureEvidenceRecord(DSSDocument, DSSDocument, ERP).
	AddSignatureEvidenceRecord(signatureDocument model.DSSDocument, evidenceRecordDocument model.DSSDocument, parameters ERP) model.DSSDocument
}
