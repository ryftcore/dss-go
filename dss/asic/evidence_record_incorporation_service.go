// Ported from dss-document/src/main/java/eu/europa/esig/dss/evidencerecord/EvidenceRecordIncorporationService.java (DSS 6.5.RC1).
//
// FORWARD DECLARATION (per S7_BRIEF.md's evidence-record note): dss-evidence-record is a
// deferred module not yet ported. AbstractASiCSignatureService implements this interface
// upstream, so a minimal local copy is ported here into package asic - method surface exactly
// as the Java interface declares (Serializable dropped per PORTING.md) - rather than gating the
// whole AbstractASiCSignatureService file on the deferred module. When dss-evidence-record is
// ported in a later phase, this forward declaration should be replaced by an import of the real
// type, and any type assertions against this local interface revisited.
package asic

import "github.com/utain/esig/dss/model"

// EvidenceRecordIncorporationService provides common methods for incorporation of evidence
// records within existing signatures, generic over the ERP implementation of format-related
// parameters for evidence record incorporation.
type EvidenceRecordIncorporationService[ERP model.SerializableEvidenceRecordIncorporationParameters] interface {
	// AddSignatureEvidenceRecord incorporates the Evidence Record as an unsigned property into
	// the signature. Port of addSignatureEvidenceRecord(DSSDocument, DSSDocument, ERP).
	AddSignatureEvidenceRecord(signatureDocument model.DSSDocument, evidenceRecordDocument model.DSSDocument, parameters ERP) model.DSSDocument
}
