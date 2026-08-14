// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/timestamp/CAdESTimestampSource.java (DSS 6.5.RC1).
//
// # FORWARD DEPENDENCY: CAdESSignature, CAdESAttribute, CAdESSignedAttributes, CAdESUnsignedAttributes
//
// *CAdESSignature (Java eu.europa.esig.dss.cades.validation.CAdESSignature) is assigned to a
// sibling chunk of this manifest, per PORTING.md item 18 implementing spi/validation.AdvancedSignature
// in full; every other landed file in this package already forward-references it the same way
// (see cades_level_baseline_lta_timestamp_extractor.go). This file additionally calls:
//
//	func (s *CAdESSignature) SignerInformation() *cmscore.SignerInfo // getSignerInformation()
//	func (s *CAdESSignature) CMS() *cms.CMS                          // getCMS()
//	func (s *CAdESSignature) DetachedContents() []model.DSSDocument  // getDetachedContents()
//	func (s *CAdESSignature) CertificateSource() *spi.SignatureCertificateSource
//	func (s *CAdESSignature) CRLSource() spi.OfflineRevocationSource[revocation.CRL]
//	func (s *CAdESSignature) OCSPSource() spi.OfflineRevocationSource[revocation.OCSP]
//	func (s *CAdESSignature) CounterSignatures() []validation.AdvancedSignature
//	func (s *CAdESSignature) ID() string
//
// all already part of validation.AdvancedSignature (frozen) or already assumed by this package's
// landed files.
//
// CAdESAttribute (Java eu.europa.esig.dss.cades.validation.CAdESAttribute) is likewise assigned
// to a sibling chunk (the same one owning CAdESSignature), and satisfies
// validation.SignatureAttribute + comparable, per SignatureTimestampSource's SA constraint. The
// subset of its API this file calls, inferred from every CAdESAttribute call in the ported Java
// source:
//
//	func (a *CAdESAttribute) ASN1Oid() asn1.ObjectIdentifier      // getASN1Oid()
//	func (a *CAdESAttribute) ASN1Object() *asn1ber.Element        // getASN1Object()
//	func (a *CAdESAttribute) AttrValues() []*asn1ber.Element      // getAttrValues() (ASN1Set)
//	func (a *CAdESAttribute) Attribute() *cmscore.Attribute       // getAttribute()
//	func (a *CAdESAttribute) ToTimeStampToken() *cmscore.TimeStampToken // toTimeStampToken()
//	func (a *CAdESAttribute) IsTimeStampToken() bool              // isTimeStampToken()
//	func (a *CAdESAttribute) IsEvidenceRecord() bool              // isEvidenceRecord()
//
// plus Identifier() identifier.SignatureAttributeIdentifier from validation.SignatureAttribute
// itself. ASN1Object()/AttrValues() are assumed to hand back *asn1ber.Element (the same
// representation DSSASN1Utils.getAsn1Encodable(Attribute) already resolves to for a single-valued
// attribute, see spi.DSSASN1UtilsAsn1Encodable), not a re-parsed byte slice: the callers below
// need DER-normalised re-encoding (asn1ber.Element#DEREncoded()), not just raw octets.
//
// CAdESSignedAttributes.build(SignerInformation) / CAdESUnsignedAttributes.build(SignerInformation)
// are Java static factories on two more sibling-chunk classes (eu.europa.esig.dss.cades.validation),
// flattened per PORTING.md's "<JavaClass><MethodName>" convention for static factories:
//
//	func CAdESSignedAttributesBuild(signerInformation *cmscore.SignerInfo) validation.SignatureProperties[*CAdESAttribute]
//	func CAdESUnsignedAttributesBuild(signerInformation *cmscore.SignerInfo) validation.SignatureProperties[*CAdESAttribute]
//
// CAdESEmbeddedEvidenceRecordHelper (Java eu.europa.esig.dss.cades.evidencerecord.
// CAdESEmbeddedEvidenceRecordHelper) is a further sibling-chunk type (a different dss-cades
// subpackage, also flattened into this same Go package per PORTING.md's cades layout), assumed
// to embed validation.AbstractEmbeddedEvidenceRecordHelper (see spi/validation/
// abstract_embedded_evidence_record_helper.go) and therefore promote SetDetachedContents,
// SetOrderOfAttribute and SetOrderWithinAttribute:
//
//	func NewCAdESEmbeddedEvidenceRecordHelper(signature *CAdESSignature, unsignedAttribute *CAdESAttribute) *CAdESEmbeddedEvidenceRecordHelper
//
// # GAP flagged for integrator: SignatureTimestampSource needs four additive accessors
//
// This file needs read access to the merged certificate/CRL/OCSP sources
// SignatureTimestampSource[AS, SA] (spi/validation/timestamp, frozen) builds up during
// populateTimestampTokens() - Java's protected `certificateSource`/`crlSource`/`ocspSource`
// fields, which every subclass (CAdESTimestampSource included) reads directly - and to the
// position-lookup its protected getAttributeOrder(SA) provides. Neither is exported today. The
// four one-method, additive, non-breaking accessors below (mirroring the already-documented
// SetTimestampedReferences GAP in signature_timestamp_source.go) are assumed to exist:
//
//	func (s *SignatureTimestampSource[AS, SA]) CertificateSource() *spi.ListCertificateSource {
//		return s.certificateSource
//	}
//	func (s *SignatureTimestampSource[AS, SA]) CRLSource() *spi.ListRevocationSource[revocation.CRL] {
//		return s.crlSource
//	}
//	func (s *SignatureTimestampSource[AS, SA]) OCSPSource() *spi.ListRevocationSource[revocation.OCSP] {
//		return s.ocspSource
//	}
//	func (s *SignatureTimestampSource[AS, SA]) GetAttributeOrder(signatureAttribute SA) *int {
//		return s.getAttributeOrder(signatureAttribute)
//	}
//
// This file calls them as CAdESTimestampSource.CertificateSource()/CRLSource()/OCSPSource()/
// GetAttributeOrder(), promoted through the embedded SignatureTimestampSource[*CAdESSignature,
// *CAdESAttribute].
//
// # GAP flagged for integrator: three CAdES overrides of "concrete but overridable" base hooks
// # are currently unreachable
//
// Unlike the fully-abstract methods of SignatureTimestampSourceOverrides (all correctly
// dispatched today), CAdESTimestampSource.java additionally @Overrides three methods that the
// abstract Java class declares CONCRETE (with a default body), relying on ordinary virtual
// dispatch for the override to take effect wherever the base class's own methods call them
// unqualified:
//
//   - incorporateArchiveTimestampReferences(TimestampToken, List<TimestampToken>) - the base's
//     own populateTimestampTokens() -> makeTimestampTokensFromUnsignedAttributes() ->
//     incorporateArchiveTimestampReferencesForTokens() calls this UNQUALIFIED for every newly
//     created archive-timestamp token, so in Java this is where EVERY CAdES archive timestamp
//     gets its correct references (ats-hash-index-v3 cross-referencing, or the V2 extra
//     getSignatureSignedDataReferences augmentation - see below).
//   - getSignatureSignedDataReferences() - called unqualified from the base's own
//     getArchiveTimestampReferences(List), processExternalTimestamp(TimestampToken) and
//     processExternalEvidenceRecord(EvidenceRecord).
//   - getCounterSignatureReferences(AdvancedSignature) - called unqualified, once per element,
//     from the base's own getCounterSignaturesReferences(List) (plural).
//
// Go has no virtual dispatch through struct embedding: every one of these three base methods -
// IncorporateArchiveTimestampReferences (exported), the private getArchiveTimestampReferences,
// GetSignatureSignedDataReferences (exported, concrete, currently returns nil) and the private
// getCounterSignatureReferences - is invoked from *within the base type's own methods* as `s.X(...)`,
// which Go binds statically to the base's own implementation regardless of what a struct
// embedding SignatureTimestampSource additionally defines under the same name. Reaching a
// concrete subclass's override from inside the base therefore requires going through the
// `overrides` field, exactly like every already-abstract method does - meaning these three
// methods need to move onto SignatureTimestampSourceOverrides[AS, SA], and their three call
// sites within signature_timestamp_source.go need to change from `s.x(...)` to
// `s.overrides.X(...)`:
//
//	IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken, previousTimestamps []*validation.TimestampToken)
//	GetSignatureSignedDataReferences() []*validation.TimestampedReference
//	GetCounterSignatureReferences(counterSignature validation.AdvancedSignature) []*validation.TimestampedReference
//
// Until this lands, this file's CAdESTimestampSource.IncorporateArchiveTimestampReferences /
// GetSignatureSignedDataReferences / GetCounterSignatureReferences are correct or best-effort
// faithful ports of the Java overrides (each fully self-contained, calling only already-exported
// base accessors, the GAP-flagged four above, and this package's own exported/reusable helpers -
// no other base access needed), reachable if invoked directly through the concrete
// *CAdESTimestampSource type, but NOT yet reachable from the base's own internal population
// flow: archive timestamps created today get the base's plain (non-CAdES-augmented) reference
// set, and counter-signature/external-timestamp/external-evidence-record processing does not
// yet pick up the merged CMS SignedData references this file's GetSignatureSignedDataReferences
// computes. Flagged prominently in the porter report; the integrator arbitrates. Every other
// override in this file (all fully abstract in Java) dispatches correctly today.
package cades

import (
	"bytes"
	"encoding/asn1"

	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/analyzer"
	"github.com/utain/esig/dss/spi/validation/timestamp"
	"github.com/utain/esig/dss/utils"
)

// CAdESTimestampSource is the timestamp source for a CAdES signature.
type CAdESTimestampSource struct {
	timestamp.SignatureTimestampSource[*CAdESSignature, *CAdESAttribute]

	// signature is being validated. Kept alongside the embedded base's own private copy: Java's
	// `signature` field, declared on the abstract SignatureTimestampSource class, is directly
	// visible to this subclass; Go has no equivalent of protected field access across packages,
	// so this file keeps its own reference, exactly as CadesLevelBaselineLTATimestampExtractor
	// does in this same package for the identical reason.
	signature *CAdESSignature
}

// NewCAdESTimestampSource is the default constructor. Port of the (CAdESSignature) constructor.
func NewCAdESTimestampSource(signature *CAdESSignature) *CAdESTimestampSource {
	source := &CAdESTimestampSource{
		SignatureTimestampSource: timestamp.NewSignatureTimestampSourceBase[*CAdESSignature, *CAdESAttribute](signature),
		signature:                signature,
	}
	source.InitSignatureTimestampSource(source)
	return source
}

// GetTimestampMessageImprintDigestBuilderForAlgorithm implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(DigestAlgorithm) override.
func (c *CAdESTimestampSource) GetTimestampMessageImprintDigestBuilderForAlgorithm(
	digestAlgorithm enumerations.DigestAlgorithm) timestamp.TimestampMessageDigestBuilder {
	return NewCAdESTimestampMessageDigestBuilder(c.signature, digestAlgorithm)
}

// GetTimestampMessageImprintDigestBuilderForToken implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(TimestampToken) override.
func (c *CAdESTimestampSource) GetTimestampMessageImprintDigestBuilderForToken(
	timestampToken *validation.TimestampToken) timestamp.TimestampMessageDigestBuilder {
	return NewCAdESTimestampMessageDigestBuilderForToken(c.signature, c.CertificateSource(), timestampToken)
}

// BuildSignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides.
// Port of buildSignedSignatureProperties().
func (c *CAdESTimestampSource) BuildSignedSignatureProperties() validation.SignatureProperties[*CAdESAttribute] {
	return CAdESSignedAttributesBuild(c.signature.SignerInformation())
}

// BuildUnsignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides.
// Port of buildUnsignedSignatureProperties().
func (c *CAdESTimestampSource) BuildUnsignedSignatureProperties() validation.SignatureProperties[*CAdESAttribute] {
	return CAdESUnsignedAttributesBuild(c.signature.SignerInformation())
}

// IsContentTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isContentTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsContentTimestamp(signedAttribute *CAdESAttribute) bool {
	return OID_id_aa_ets_contentTimestamp.Equal(signedAttribute.ASN1Oid())
}

// IsAllDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isAllDataObjectsTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsAllDataObjectsTimestamp(signedAttribute *CAdESAttribute) bool {
	return false
}

// IsIndividualDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isIndividualDataObjectsTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsIndividualDataObjectsTimestamp(signedAttribute *CAdESAttribute) bool {
	return false
}

// IsSignatureTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignatureTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsSignatureTimestamp(unsignedAttribute *CAdESAttribute) bool {
	return OID_id_aa_signatureTimeStampToken.Equal(unsignedAttribute.ASN1Oid())
}

// IsCompleteCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteCertificateRef(CAdESAttribute).
func (c *CAdESTimestampSource) IsCompleteCertificateRef(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_certificateRefs.Equal(unsignedAttribute.ASN1Oid())
}

// IsAttributeCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeCertificateRef(CAdESAttribute).
func (c *CAdESTimestampSource) IsAttributeCertificateRef(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_attributeCertificateRefsOid.Equal(unsignedAttribute.ASN1Oid())
}

// IsCompleteRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteRevocationRef(CAdESAttribute).
func (c *CAdESTimestampSource) IsCompleteRevocationRef(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_revocationRefs.Equal(unsignedAttribute.ASN1Oid())
}

// IsAttributeRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeRevocationRef(CAdESAttribute).
func (c *CAdESTimestampSource) IsAttributeRevocationRef(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_attributeRevocationRefsOid.Equal(unsignedAttribute.ASN1Oid())
}

// IsRefsOnlyTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRefsOnlyTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsRefsOnlyTimestamp(unsignedAttribute *CAdESAttribute) bool {
	return OID_id_aa_ets_certCRLTimestamp.Equal(unsignedAttribute.ASN1Oid())
}

// IsSigAndRefsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSigAndRefsTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsSigAndRefsTimestamp(unsignedAttribute *CAdESAttribute) bool {
	return OID_id_aa_ets_escTimeStamp.Equal(unsignedAttribute.ASN1Oid())
}

// IsCertificateValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCertificateValues(CAdESAttribute).
func (c *CAdESTimestampSource) IsCertificateValues(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_certValues.Equal(unsignedAttribute.ASN1Oid())
}

// IsRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRevocationValues(CAdESAttribute).
func (c *CAdESTimestampSource) IsRevocationValues(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_revocationValues.Equal(unsignedAttribute.ASN1Oid())
}

// IsAttrAuthoritiesCertValues implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isAttrAuthoritiesCertValues(CAdESAttribute).
func (c *CAdESTimestampSource) IsAttrAuthoritiesCertValues(unsignedAttribute *CAdESAttribute) bool {
	return false
}

// IsAttributeRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isAttributeRevocationValues(CAdESAttribute).
func (c *CAdESTimestampSource) IsAttributeRevocationValues(unsignedAttribute *CAdESAttribute) bool {
	return false
}

// IsArchiveTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isArchiveTimestamp(CAdESAttribute).
func (c *CAdESTimestampSource) IsArchiveTimestamp(unsignedAttribute *CAdESAttribute) bool {
	return c.isArchiveTimestampV2Attribute(unsignedAttribute) || c.isArchiveTimestampV3Attribute(unsignedAttribute)
}

// isArchiveTimestampV2Attribute ports the private isArchiveTimestampV2(CAdESAttribute); named
// distinctly from isArchiveTimestampV2Token (the private (TimestampToken) overload) since Go has
// no overloading.
func (c *CAdESTimestampSource) isArchiveTimestampV2Attribute(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_archiveTimestampV2.Equal(unsignedAttribute.ASN1Oid())
}

// isArchiveTimestampV3Attribute ports the private isArchiveTimestampV3(CAdESAttribute).
func (c *CAdESTimestampSource) isArchiveTimestampV3Attribute(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_archiveTimestampV3.Equal(unsignedAttribute.ASN1Oid())
}

// IsTimeStampValidationData implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isTimeStampValidationData(CAdESAttribute).
func (c *CAdESTimestampSource) IsTimeStampValidationData(unsignedAttribute *CAdESAttribute) bool {
	return false
}

// IsAnyValidationData implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isAnyValidationData(CAdESAttribute).
func (c *CAdESTimestampSource) IsAnyValidationData(unsignedAttribute *CAdESAttribute) bool {
	return false
}

// IsValidationDataReferences implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for CAdES. Port of isValidationDataReferences(CAdESAttribute).
func (c *CAdESTimestampSource) IsValidationDataReferences(unsignedAttribute *CAdESAttribute) bool {
	return false
}

// IsCounterSignature implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCounterSignature(CAdESAttribute).
func (c *CAdESTimestampSource) IsCounterSignature(unsignedAttribute *CAdESAttribute) bool {
	return OID_counterSignature.Equal(unsignedAttribute.ASN1Oid())
}

// OID_counterSignature is org.bouncycastle.asn1.cms.CMSAttributes#counterSignature
// (1.2.840.113549.1.9.6), a PKCS#9 attribute type with no eu.europa.esig.dss.spi.OID
// counterpart; named per PORTING.md's "OID_<javaFieldName>" convention, matching how
// cades_utils.go already names other BouncyCastle-only OIDs this package needs.
var OID_counterSignature = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 6}

// IsSignaturePolicyStore implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignaturePolicyStore(CAdESAttribute).
func (c *CAdESTimestampSource) IsSignaturePolicyStore(unsignedAttribute *CAdESAttribute) bool {
	return spi.OID_id_aa_ets_sigPolicyStore.Equal(unsignedAttribute.ASN1Oid())
}

// IsEvidenceRecord implements timestamp.SignatureTimestampSourceOverrides.
// Port of isEvidenceRecord(CAdESAttribute).
func (c *CAdESTimestampSource) IsEvidenceRecord(unsignedAttribute *CAdESAttribute) bool {
	return unsignedAttribute.IsEvidenceRecord()
}

// MakeTimestampToken implements timestamp.SignatureTimestampSourceOverrides.
// Port of makeTimestampToken(CAdESAttribute, TimestampType, List).
//
// Panics on a construction failure (model.NewDSSErrorWithCause), matching Java's unchecked
// propagation of `new TimestampToken(...)`'s failure through a method with no throws clause; see
// the SignatureTimestampSourceOverrides doc comment on this convention.
func (c *CAdESTimestampSource) MakeTimestampToken(signatureAttribute *CAdESAttribute, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) *validation.TimestampToken {
	timeStampToken := signatureAttribute.ToTimeStampToken()
	if timeStampToken == nil {
		return nil
	}
	identifierBuilder := timestamp.NewSignatureTimestampIdentifierBuilderFromToken(timeStampToken).
		SetSignature(c.signature).
		SetAttribute(signatureAttribute).
		SetOrderOfAttribute(c.GetAttributeOrder(signatureAttribute))
	token, err := validation.NewTimestampTokenFromTimeStampTokenWithIdentifierBuilder(
		timeStampToken, timestampType, references, identifierBuilder)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return token
}

// MakeEvidenceRecords implements timestamp.SignatureTimestampSourceOverrides.
// Port of makeEvidenceRecords(CAdESAttribute, List).
func (c *CAdESTimestampSource) MakeEvidenceRecords(unsignedAttribute *CAdESAttribute,
	references []*validation.TimestampedReference) []validation.EvidenceRecord {
	var result []validation.EvidenceRecord
	attrValues := unsignedAttribute.AttrValues()
	for i, erASN1Encodable := range attrValues {
		if evidenceRecord := c.createEvidenceRecord(unsignedAttribute, erASN1Encodable, i); evidenceRecord != nil {
			result = append(result, evidenceRecord)
		}
	}
	return result
}

// createEvidenceRecord ports the private createEvidenceRecord(CAdESAttribute, ASN1Encodable, int),
// whose surrounding try/catch (any failure logged and answered with null) is reproduced by
// returning nil on the first failing step, per PORTING.md.
func (c *CAdESTimestampSource) createEvidenceRecord(unsignedAttribute *CAdESAttribute,
	erASN1Encodable *asn1ber.Element, orderWithinAttribute int) validation.EvidenceRecord {
	erEncoded := erASN1Encodable.DEREncoded()
	erDocument := model.NewInMemoryDocument(erEncoded)
	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(erDocument)
	if err != nil {
		// Upstream logs "Unable to build an embedded evidence record. Reason : {}".
		return nil
	}
	evidenceRecordAnalyzer.SetEvidenceRecordOrigin(enumerations.EvidenceRecordOrigin_SIGNATURE)

	incorporationType := CAdESUtilsEvidenceRecordIncorporationType(unsignedAttribute.ASN1Oid())
	evidenceRecordAnalyzer.SetEvidenceRecordIncorporationType(incorporationType)

	embeddedEvidenceRecordHelper := NewCAdESEmbeddedEvidenceRecordHelper(c.signature, unsignedAttribute)
	if enumerations.EvidenceRecordIncorporationType_EXTERNAL_EVIDENCE_RECORD == incorporationType {
		if utils.CollectionSize(c.signature.DetachedContents()) == 1 {
			embeddedEvidenceRecordHelper.SetDetachedContents(c.signature.DetachedContents())
		}
		// else: Upstream logs "Detached document has not been provided to the validation of an
		// external-evidence-record!".
	}
	embeddedEvidenceRecordHelper.SetOrderOfAttribute(c.GetAttributeOrder(unsignedAttribute))
	orderWithin := orderWithinAttribute
	embeddedEvidenceRecordHelper.SetOrderWithinAttribute(&orderWithin)

	evidenceRecordAnalyzer.SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper)
	return evidenceRecordAnalyzer.EvidenceRecord()
}

// IncorporateArchiveTimestampReferences is CAdES's override of the base's concrete
// incorporateArchiveTimestampReferences(TimestampToken, List<TimestampToken>) hook.
// Port of that override.
//
// See the file header GAP note: not yet reachable from the base's own internal population flow
// (incorporateArchiveTimestampReferencesForTokens calls the base's own implementation, not this
// one) until SignatureTimestampSourceOverrides gains this method.
func (c *CAdESTimestampSource) IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken,
	previousTimestamps []*validation.TimestampToken) {
	if c.isArchiveTimestampV2Token(timestampToken) {
		// For an ATSTv2 all the incorporated unsigned properties are covered.
		//
		// This reproduces `super.incorporateArchiveTimestampReferences(...)` as Java's virtual
		// dispatch actually executes it: the base method's own body calls
		// getArchiveTimestampReferences(previousTimestamps) unqualified, which - being a virtual
		// call in Java - resolves to CAdES's own override (this file's archiveTimestampReferences,
		// not the base's plain default), even though the call originates from "super". Calling
		// archiveTimestampReferences directly here reproduces that actual behaviour.
		cadesTSTimestampAddReferences(timestampToken, c.archiveTimestampReferences(previousTimestamps))
	}
	// else archive-timestamp-v3.
	var timestampedReferences []*validation.TimestampedReference
	cadesTSAddReferences(&timestampedReferences, c.signatureTimestampReferences())

	unsignedAttributes := timestampToken.UnsignedAttributes()
	atsHashIndexVersionIdentifier := CAdESUtilsAtsHashIndexVersionIdentifier(unsignedAttributes)
	atsHashIndex := CAdESUtilsAtsHashIndexByVersion(unsignedAttributes, atsHashIndexVersionIdentifier)
	if atsHashIndex != nil {
		digestAlgorithm := c.hashIndexDigestAlgorithm(atsHashIndex)

		certsHashIndex := CAdESUtilsCertificatesHashIndex(atsHashIndex)
		crlHashIndex := CAdESUtilsCRLHashIndex(atsHashIndex)
		cadesTSAddReferences(&timestampedReferences, c.signedDataCertificateReferences(certsHashIndex, digestAlgorithm))
		cadesTSAddReferences(&timestampedReferences, c.signedDataRevocationReferences(crlHashIndex, digestAlgorithm))

		unsignedAttrsHashIndex := CAdESUtilsUnsignedAttributesHashIndex(atsHashIndex)
		cadesTSAddReferences(&timestampedReferences, c.unsignedAttributesReferences(
			atsHashIndexVersionIdentifier, unsignedAttrsHashIndex, digestAlgorithm, previousTimestamps))
	}
	cadesTSTimestampAddReferences(timestampToken, timestampedReferences)
}

// isArchiveTimestampV2Token ports the private isArchiveTimestampV2(TimestampToken); named
// distinctly from isArchiveTimestampV2Attribute (the private (CAdESAttribute) overload) since Go
// has no overloading.
func (c *CAdESTimestampSource) isArchiveTimestampV2Token(timestampToken *validation.TimestampToken) bool {
	return timestampToken.ArchiveTimestampType() == enumerations.ArchiveTimestampType_CAdES_V2 ||
		timestampToken.ArchiveTimestampType() == enumerations.ArchiveTimestampType_CAdES
}

// archiveTimestampReferences ports getArchiveTimestampReferences(List<TimestampToken>): the
// base's own three-reference build (see spi/validation/timestamp/signature_timestamp_source.go's
// getArchiveTimestampReferences), reimplemented here from exported accessors since this file
// cannot reach the base's private version, followed by CAdES's own extra addition ("executed for
// ArchiveTimestampV2 only" upstream). Port of the override.
func (c *CAdESTimestampSource) archiveTimestampReferences(previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	var timestampedReferences []*validation.TimestampedReference
	cadesTSAddReferences(&timestampedReferences, c.signatureTimestampReferences())
	cadesTSAddReferences(&timestampedReferences, c.encapsulatedReferencesFromTimestamps(previousTimestamps))
	cadesTSAddReferences(&timestampedReferences, c.UnsignedPropertiesReferences())
	// executed for ArchiveTimestampV2 only.
	cadesTSAddReferences(&timestampedReferences, c.GetSignatureSignedDataReferences())
	return timestampedReferences
}

// signatureTimestampReferences reimplements the base's private getSignatureTimestampReferences(),
// which this file cannot call directly, from exported base accessors.
func (c *CAdESTimestampSource) signatureTimestampReferences() []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	cadesTSAddReferences(&references, c.encapsulatedReferencesFromTimestamps(c.ContentTimestamps()))
	cadesTSAddReferences(&references, c.SignerDataReferences())
	cadesTSAddReference(&references, validation.NewTimestampedReference(c.signature.ID(), enumerations.TimestampedObjectType_SIGNATURE))
	cadesTSAddReferences(&references, c.signingCertificateTimestampReferences())
	return references
}

// signingCertificateTimestampReferences reimplements the base's private
// getSigningCertificateTimestampReferences().
func (c *CAdESTimestampSource) signingCertificateTimestampReferences() []*validation.TimestampedReference {
	signatureCertificateSource := c.signature.CertificateSource()
	return timestamp.CreateReferencesForCertificateRefs(
		signatureCertificateSource.SigningCertificateRefs(), signatureCertificateSource, c.CertificateSource())
}

// encapsulatedReferencesFromTimestamps reimplements the base's private
// getEncapsulatedReferencesFromTimestamps(List), from the exported timestamp.ReferencesFromTimestamp
// and this file's GAP-flagged merged-source accessors.
func (c *CAdESTimestampSource) encapsulatedReferencesFromTimestamps(timestampTokens []*validation.TimestampToken) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	for _, timestampToken := range timestampTokens {
		refs := cadesTSMust(timestamp.ReferencesFromTimestamp(timestampToken, c.CertificateSource(), c.CRLSource(), c.OCSPSource()))
		cadesTSAddReferences(&references, refs)
	}
	return references
}

// signedDataCertificateReferences ports the private
// getSignedDataCertificateReferences(ASN1Sequence, DigestAlgorithm).
//
// DEVIATION: drops the `signatureCertificateSource instanceof CMSCertificateSource` guard.
// spi.SignatureCertificateSource#SignedDataCertificates() (the Go port of
// SignatureCertificateSource#getSignedDataCertificates()) is already declared, and answers
// correctly, on the base type itself in both Java and Go - the instanceof guard is defensive and
// redundant even upstream; see cades_timestamp_message_digest_builder.go's file header for the
// identical judgment call.
func (c *CAdESTimestampSource) signedDataCertificateReferences(certsHashIndex []byte,
	digestAlgorithm enumerations.DigestAlgorithm) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	certsHashList := cadesTSMust(spi.DSSASN1UtilsDEROctetStrings(certsHashIndex))
	for _, certificate := range c.signature.CertificateSource().SignedDataCertificates() {
		digest := cadesTSMust(certificate.Digest(digestAlgorithm))
		if cadesTSIsDigestValuePresent(digest, certsHashList) {
			cadesTSAddReference(&references, timestamp.CreateReferenceForCertificate(certificate))
		}
		// Upstream logs whether the certificate is present in the timestamp.
	}
	return references
}

// signedDataRevocationReferences ports the private
// getSignedDataRevocationReferences(ASN1Sequence, DigestAlgorithm).
func (c *CAdESTimestampSource) signedDataRevocationReferences(crlsHashIndex []byte,
	digestAlgorithm enumerations.DigestAlgorithm) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	crlsHashList := cadesTSMust(spi.DSSASN1UtilsDEROctetStrings(crlsHashIndex))
	cadesTSAddReferences(&references, timestamp.CreateReferencesForCRLBinaries(c.signedDataCRLBinaries(crlsHashList, digestAlgorithm)))
	ocspRefs := cadesTSMust(timestamp.CreateReferencesForOCSPBinaries(
		c.signedDataOCSPResponseBinaries(crlsHashList, digestAlgorithm), c.CertificateSource()))
	cadesTSAddReferences(&references, ocspRefs)
	return references
}

// signedDataCRLBinaries ports the private getSignedDataCRLBinaries(List, DigestAlgorithm),
// reusing cadesLTACMSSignedDataRevocationBinaries (defined in this package by
// cades_level_baseline_lta_timestamp_extractor.go) as the Go equivalent of the
// `instanceof CMSCRLSource` guard - see that function's own doc comment.
func (c *CAdESTimestampSource) signedDataCRLBinaries(crlsHashList [][]byte,
	digestAlgorithm enumerations.DigestAlgorithm) []*crlparser.CRLBinary {
	var crlBinaries []*crlparser.CRLBinary
	for _, identifier := range cadesLTACMSSignedDataRevocationBinaries[revocation.CRL](c.signature.CRLSource()) {
		crlBinary, ok := identifier.(*crlparser.CRLBinary)
		if !ok {
			continue
		}
		digest, err := crlBinary.DigestValue(digestAlgorithm)
		if err != nil {
			continue
		}
		if cadesTSIsDigestValuePresent(digest, crlsHashList) {
			crlBinaries = append(crlBinaries, crlBinary)
		}
	}
	return crlBinaries
}

// signedDataOCSPResponseBinaries ports the private
// getSignedDataOCSPResponseBinaries(List, DigestAlgorithm).
func (c *CAdESTimestampSource) signedDataOCSPResponseBinaries(crlsHashList [][]byte,
	digestAlgorithm enumerations.DigestAlgorithm) []*spi.OCSPResponseBinary {
	var ocspBinaries []*spi.OCSPResponseBinary
	for _, identifier := range cadesLTACMSSignedDataRevocationBinaries[revocation.OCSP](c.signature.OCSPSource()) {
		binary, ok := identifier.(*spi.OCSPResponseBinary)
		if !ok {
			continue
		}
		if c.isOCSPResponsePresent(binary, crlsHashList, digestAlgorithm) {
			ocspBinaries = append(ocspBinaries, binary)
		}
	}
	return ocspBinaries
}

// isOCSPResponsePresent ports the private isOCSPResponsePresent(OCSPResponseBinary, List,
// DigestAlgorithm): "This method tries the best option to verify the digest of an OCSP response,
// if it fails, tries the alternative encoding".
func (c *CAdESTimestampSource) isOCSPResponsePresent(binary *spi.OCSPResponseBinary, crlsHashList [][]byte,
	digestAlgorithm enumerations.DigestAlgorithm) bool {
	objectIdentifier := binary.ASN1ObjectIdentifier()
	if cmscore.OIDPKIXOCSPBasic.Equal(objectIdentifier) {
		return cadesTSIsOCSPDigestValueMatch(binary.BasicOCSPRespContent(), objectIdentifier, crlsHashList, digestAlgorithm)
	}
	// OCSPObjectIdentifiers.id_ri_ocsp_response case.
	if objectIdentifier == nil {
		objectIdentifier = cadesLTAPKIXOCSPResponseOID
	}
	return cadesTSIsOCSPDigestValueMatch(binary.Binaries(), objectIdentifier, crlsHashList, digestAlgorithm) ||
		cadesTSIsOCSPDigestValueMatch(binary.BasicOCSPRespContent(), objectIdentifier, crlsHashList, digestAlgorithm)
}

// cadesTSIsOCSPDigestValueMatch ports the private isOCSPDigestValueMatch(byte[],
// ASN1ObjectIdentifier, List, DigestAlgorithm).
func cadesTSIsOCSPDigestValueMatch(binaries []byte, objectIdentifier asn1.ObjectIdentifier, crlsHashList [][]byte,
	digestAlgorithm enumerations.DigestAlgorithm) bool {
	encoded := cadesTSMust(CAdESUtilsSignedDataEncodedOCSPResponse(binaries, objectIdentifier))
	digest := cadesTSMust(spi.DSSUtilsDigest(digestAlgorithm, encoded))
	return cadesTSIsDigestValuePresent(digest, crlsHashList)
}

// unsignedAttributesReferences ports the private getUnsignedAttributesReferences
// (ASN1ObjectIdentifier, ASN1Sequence, DigestAlgorithm, List).
//
// Rebuilds the unsigned-signature-properties directly via CAdESUnsignedAttributesBuild rather
// than reaching the base's cached, identity-comparable copy (unreachable from this package - see
// the file header GAP note on getAttributeOrder for the one place that distinction *does*
// matter): this method only ever compares attribute *content* (a digest of its octets) against
// timestampUnsignedAttributesHashesList, never object identity, so a freshly parsed copy is
// behaviourally identical.
func (c *CAdESTimestampSource) unsignedAttributesReferences(atsHashIndexVersionIdentifier asn1.ObjectIdentifier,
	unsignedAttrsHashIndex []byte, digestAlgorithm enumerations.DigestAlgorithm,
	previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	timestampUnsignedAttributesHashesList := cadesTSMust(spi.DSSASN1UtilsDEROctetStrings(unsignedAttrsHashIndex))

	unsignedSignatureProperties := CAdESUnsignedAttributesBuild(c.signature.SignerInformation())
	for _, unsignedAttribute := range unsignedSignatureProperties.Attributes() {
		octets := cadesTSMust(CAdESUtilsOctetStringForAtsHashIndex(unsignedAttribute.Attribute(), atsHashIndexVersionIdentifier))
		for _, bytesToDigest := range octets {
			digest := cadesTSMust(spi.DSSUtilsDigest(digestAlgorithm, bytesToDigest))
			if cadesTSIsDigestValuePresent(digest, timestampUnsignedAttributesHashesList) {
				cadesTSAddReferences(&references, c.referencesFromUnsignedProperty(unsignedAttribute, previousTimestamps))
			}
		}
	}
	return references
}

// referencesFromUnsignedProperty ports the private getReferencesFromUnsignedProperty
// (CAdESAttribute, List).
func (c *CAdESTimestampSource) referencesFromUnsignedProperty(unsignedAttribute *CAdESAttribute,
	previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	switch {
	case unsignedAttribute.IsTimeStampToken():
		references := c.referencesFromMatchingTimestamp(unsignedAttribute, previousTimestamps)
		if len(references) == 0 {
			// Upstream logs "The timestamp order is broken! Unable to find a covered timestamp.".
		}
		return references

	case c.IsCompleteCertificateRef(unsignedAttribute) || c.IsAttributeCertificateRef(unsignedAttribute):
		return c.timestampedCertificateRefs(unsignedAttribute)

	case c.IsCompleteRevocationRef(unsignedAttribute) || c.IsAttributeRevocationRef(unsignedAttribute):
		return c.timestampedRevocationRefs(unsignedAttribute)

	case c.IsCertificateValues(unsignedAttribute):
		return c.timestampedCertificateValues(unsignedAttribute)

	case c.IsRevocationValues(unsignedAttribute):
		return c.timestampedRevocationValues(unsignedAttribute)

	case c.IsCounterSignature(unsignedAttribute):
		counterSignatures := c.GetCounterSignatures(unsignedAttribute)
		return c.counterSignaturesReferences(counterSignatures)

	default:
		// Upstream logs "Unable to find an unsigned attribute with the digest from
		// ats-hash-index-v3".
		return nil
	}
}

// referencesFromMatchingTimestamp ports the private getReferencesFromMatchingTimestamp
// (CAdESAttribute, List).
func (c *CAdESTimestampSource) referencesFromMatchingTimestamp(unsignedAttribute *CAdESAttribute,
	previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	asn1Object := unsignedAttribute.ASN1Object()
	derEncoded := asn1Object.DEREncoded()
	for _, timestampToken := range previousTimestamps {
		if bytes.Equal(derEncoded, timestampToken.Encoded()) {
			return cadesTSMust(timestamp.ReferencesFromTimestamp(timestampToken, c.CertificateSource(), c.CRLSource(), c.OCSPSource()))
		}
	}
	return nil
}

// timestampedCertificateRefs reimplements the base's private getTimestampedCertificateRefs(SA),
// which this file cannot call directly, from the exported timestamp.CreateReferencesForCertificateRefs.
func (c *CAdESTimestampSource) timestampedCertificateRefs(unsignedAttribute *CAdESAttribute) []*validation.TimestampedReference {
	return timestamp.CreateReferencesForCertificateRefs(c.GetCertificateRefs(unsignedAttribute), c.signature.CertificateSource(), c.CertificateSource())
}

// timestampedRevocationRefs reimplements the base's private getTimestampedRevocationRefs(SA).
func (c *CAdESTimestampSource) timestampedRevocationRefs(unsignedAttribute *CAdESAttribute) []*validation.TimestampedReference {
	var timestampedReferences []*validation.TimestampedReference
	cadesTSAddReferences(&timestampedReferences,
		timestamp.CreateReferencesForCRLRefs(c.GetCRLRefs(unsignedAttribute), c.signature.CRLSource(), c.CRLSource()))
	ocspRefs := cadesTSMust(timestamp.CreateReferencesForOCSPRefs(
		c.GetOCSPRefs(unsignedAttribute), c.signature.OCSPSource(), c.CertificateSource(), c.OCSPSource()))
	cadesTSAddReferences(&timestampedReferences, ocspRefs)
	return timestampedReferences
}

// timestampedCertificateValues reimplements the base's private
// getTimestampedCertificateValues(SA).
func (c *CAdESTimestampSource) timestampedCertificateValues(unsignedAttribute *CAdESAttribute) []*validation.TimestampedReference {
	return timestamp.CreateReferencesForIdentifiers(
		c.GetEncapsulatedCertificateIdentifiers(unsignedAttribute), enumerations.TimestampedObjectType_CERTIFICATE)
}

// timestampedRevocationValues reimplements the base's private
// getTimestampedRevocationValues(SA).
func (c *CAdESTimestampSource) timestampedRevocationValues(unsignedAttribute *CAdESAttribute) []*validation.TimestampedReference {
	var timestampedReferences []*validation.TimestampedReference
	cadesTSAddReferences(&timestampedReferences, timestamp.CreateReferencesForCRLBinaries(c.GetEncapsulatedCRLIdentifiers(unsignedAttribute)))
	ocspRefs := cadesTSMust(timestamp.CreateReferencesForOCSPBinaries(c.GetEncapsulatedOCSPIdentifiers(unsignedAttribute), c.CertificateSource()))
	cadesTSAddReferences(&timestampedReferences, ocspRefs)
	return timestampedReferences
}

// counterSignaturesReferences reimplements the base's private
// getCounterSignaturesReferences(List) (plural), calling this file's own GetCounterSignatureReferences
// (singular) directly rather than through the base (see the file header GAP note): this is, if
// anything, more correct than what the base's own dispatch-broken plural helper would currently
// produce for a CAdES signature.
func (c *CAdESTimestampSource) counterSignaturesReferences(counterSignatures []validation.AdvancedSignature) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	for _, counterSignature := range counterSignatures {
		references = append(references, c.GetCounterSignatureReferences(counterSignature)...)
	}
	return references
}

// GetSignatureSignedDataReferences is CAdES's override of the base's concrete
// getSignatureSignedDataReferences() hook. Port of that override.
//
// DEVIATION: drops the `signatureCertificateSource instanceof CMSCertificateSource` guard, per
// the identical judgment call in signedDataCertificateReferences above. The CRL/OCSP branches
// keep the Go equivalent of their instanceof guards via cadesLTACMSSignedDataRevocationBinaries.
//
// See the file header GAP note: not yet reachable from the base's own internal callers
// (getArchiveTimestampReferences, processExternalTimestamp, processExternalEvidenceRecord) until
// SignatureTimestampSourceOverrides gains this method.
func (c *CAdESTimestampSource) GetSignatureSignedDataReferences() []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	cadesTSAddReferences(&references, timestamp.CreateReferencesForCertificates(c.signature.CertificateSource().SignedDataCertificates()))
	cadesTSAddReferences(&references, timestamp.CreateReferencesForCRLBinaries(
		cadesLTACMSSignedDataRevocationBinaries[revocation.CRL](c.signature.CRLSource())))
	ocspRefs := cadesTSMust(timestamp.CreateReferencesForOCSPBinaries(
		cadesLTACMSSignedDataRevocationBinaries[revocation.OCSP](c.signature.OCSPSource()), c.CertificateSource()))
	cadesTSAddReferences(&references, ocspRefs)
	return references
}

// hashIndexDigestAlgorithm ports the private getHashIndexDigestAlgorithm(ASN1Sequence), which is
// identical to cadesLTAHashIndexDigestAlgorithm (defined by
// cades_level_baseline_lta_timestamp_extractor.go in this package, given the algorithm
// identifier rather than the raw ats-hash-index sequence) applied to the same
// AlgorithmIdentifierFromATSHashIndex extraction; reused directly (both Java classes carry their
// own private copy of the identical logic).
func (c *CAdESTimestampSource) hashIndexDigestAlgorithm(atsHashIndex []byte) enumerations.DigestAlgorithm {
	return cadesTSMust(cadesLTAHashIndexDigestAlgorithm(spi.DSSASN1UtilsAlgorithmIdentifierFromATSHashIndex(atsHashIndex)))
}

// cadesTSIsDigestValuePresent ports the private isDigestValuePresent(byte[], List<DEROctetString>).
func cadesTSIsDigestValuePresent(digestValue []byte, hashList [][]byte) bool {
	for _, hash := range hashList {
		if bytes.Equal(hash, digestValue) {
			return true
		}
	}
	return false
}

// GetCertificateRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCertificateRefs(CAdESAttribute).
func (c *CAdESTimestampSource) GetCertificateRefs(unsignedAttribute *CAdESAttribute) []*spi.CertificateRef {
	var certRefs []*spi.CertificateRef
	seq := unsignedAttribute.ASN1Object()
	if seq == nil {
		return certRefs
	}
	for _, child := range seq.Children() {
		otherCertID, err := spi.ParseOtherCertID(child.Encoded())
		if err != nil {
			// Upstream logs "Unable to parse encapsulated OtherCertID : {}".
			continue
		}
		certRef, err := spi.DSSASN1UtilsCertificateRef(otherCertID)
		if err != nil {
			continue
		}
		certRefs = append(certRefs, certRef)
	}
	return certRefs
}

// GetCRLRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCRLRefs(CAdESAttribute). Unlike GetCertificateRefs, Java declares no per-iteration
// try/catch here: a parse failure panics, matching unchecked propagation.
func (c *CAdESTimestampSource) GetCRLRefs(unsignedAttribute *CAdESAttribute) []*spi.CRLRef {
	var refs []*spi.CRLRef
	seq := unsignedAttribute.ASN1Object()
	if seq == nil {
		return refs
	}
	for _, child := range seq.Children() {
		otherRefId := cadesTSMust(spi.ParseCrlOcspRef(child.Encoded()))
		if otherRefId.Crlids != nil {
			for _, id := range cadesTSMust(otherRefId.Crlids.Crls()) {
				refs = append(refs, cadesTSMust(spi.NewCRLRefFromCrlValidatedID(id)))
			}
		}
	}
	return refs
}

// GetOCSPRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getOCSPRefs(CAdESAttribute); see GetCRLRefs on the absence of a per-iteration try/catch.
func (c *CAdESTimestampSource) GetOCSPRefs(unsignedAttribute *CAdESAttribute) []*spi.OCSPRef {
	var refs []*spi.OCSPRef
	seq := unsignedAttribute.ASN1Object()
	if seq == nil {
		return refs
	}
	for _, child := range seq.Children() {
		otherCertId := cadesTSMust(spi.ParseCrlOcspRef(child.Encoded()))
		if otherCertId.Ocspids != nil {
			for _, ocspResponsesID := range cadesTSMust(otherCertId.Ocspids.OcspResponses()) {
				refs = append(refs, cadesTSMust(spi.NewOCSPRefFromOcspResponsesID(ocspResponsesID)))
			}
		}
	}
	return refs
}

// GetEncapsulatedCertificateIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCertificateIdentifiers(CAdESAttribute).
func (c *CAdESTimestampSource) GetEncapsulatedCertificateIdentifiers(unsignedAttribute *CAdESAttribute) []model.Identifier {
	var certificateIdentifiers []model.Identifier
	seq := unsignedAttribute.ASN1Object()
	if seq == nil {
		return certificateIdentifiers
	}
	for _, child := range seq.Children() {
		certificateToken, err := spi.DSSUtilsLoadCertificateFromBinary(child.Encoded())
		if err != nil {
			// Upstream logs "Unable to parse an encapsulated certificate : {}".
			continue
		}
		certificateIdentifiers = append(certificateIdentifiers, certificateToken.DSSID())
	}
	return certificateIdentifiers
}

// GetEncapsulatedCRLIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCRLIdentifiers(CAdESAttribute).
func (c *CAdESTimestampSource) GetEncapsulatedCRLIdentifiers(unsignedAttribute *CAdESAttribute) []*crlparser.CRLBinary {
	asn1Object := unsignedAttribute.ASN1Object()
	if asn1Object == nil {
		return nil
	}
	revocationValues := spi.DSSASN1UtilsRevocationValues(asn1Object.Encoded())
	if revocationValues == nil {
		return nil
	}
	return c.buildCRLIdentifiers(revocationValues.CrlVals)
}

// buildCRLIdentifiers ports the protected buildCRLIdentifiers(CertificateList...).
func (c *CAdESTimestampSource) buildCRLIdentifiers(crlVals [][]byte) []*crlparser.CRLBinary {
	var crlBinaryIdentifiers []*crlparser.CRLBinary
	for _, crlVal := range crlVals {
		crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(crlVal)
		if err != nil {
			// Upstream logs "Unable to parse CRL binaries : {}".
			continue
		}
		crlBinaryIdentifiers = append(crlBinaryIdentifiers, crlBinary)
	}
	return crlBinaryIdentifiers
}

// GetEncapsulatedOCSPIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedOCSPIdentifiers(CAdESAttribute).
func (c *CAdESTimestampSource) GetEncapsulatedOCSPIdentifiers(unsignedAttribute *CAdESAttribute) []*spi.OCSPResponseBinary {
	asn1Object := unsignedAttribute.ASN1Object()
	if asn1Object == nil {
		return nil
	}
	revocationValues := spi.DSSASN1UtilsRevocationValues(asn1Object.Encoded())
	if revocationValues == nil {
		return nil
	}
	return c.buildOCSPIdentifiers(revocationValues.OcspVals)
}

// buildOCSPIdentifiers ports the protected buildOCSPIdentifiers(BasicOCSPResponse...).
func (c *CAdESTimestampSource) buildOCSPIdentifiers(ocspVals [][]byte) []*spi.OCSPResponseBinary {
	var ocspIdentifiers []*spi.OCSPResponseBinary
	for _, ocspVal := range ocspVals {
		basicOCSPResp := spi.DSSRevocationUtilsBasicOcspResp(ocspVal)
		if basicOCSPResp == nil {
			continue
		}
		binary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
		if err != nil {
			// Upstream logs "Unable to parse OCSP response binaries : {}".
			continue
		}
		ocspIdentifiers = append(ocspIdentifiers, binary)
	}
	return ocspIdentifiers
}

// GetArchiveTimestampType implements timestamp.SignatureTimestampSourceOverrides.
// Port of getArchiveTimestampType(CAdESAttribute).
func (c *CAdESTimestampSource) GetArchiveTimestampType(unsignedAttribute *CAdESAttribute) enumerations.ArchiveTimestampType {
	switch {
	case spi.OID_id_aa_ets_archiveTimestampV2.Equal(unsignedAttribute.ASN1Oid()):
		return enumerations.ArchiveTimestampType_CAdES_V2
	case spi.OID_id_aa_ets_archiveTimestampV3.Equal(unsignedAttribute.ASN1Oid()):
		return enumerations.ArchiveTimestampType_CAdES_V3
	}
	return enumerations.ArchiveTimestampType_CAdES
}

// GetCounterSignatures implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCounterSignatures(CAdESAttribute).
//
// Unable to build a SignerInformation with BC (protected constructor): extract all found
// CounterSignatures and compare with a found SignerInfo(s).
func (c *CAdESTimestampSource) GetCounterSignatures(unsignedAttribute *CAdESAttribute) []validation.AdvancedSignature {
	var cadesResult []validation.AdvancedSignature
	allCounterSignatures := c.signature.CounterSignatures()

	for _, element := range unsignedAttribute.AttrValues() {
		si := cadesTSMust(cmscore.SignerInfoFromElement(element))
		encodedSI := si.DER()

		for _, counterSignature := range allCounterSignatures {
			// Java casts unconditionally: `(CAdESSignature) counterSignature`. Every concrete
			// AdvancedSignature a CAdES source's CounterSignatures() hands out is a
			// *CAdESSignature, so this assertion always succeeds in practice; a genuine mismatch
			// panics, exactly as Java's ClassCastException would.
			cadesCounterSignature := counterSignature.(*CAdESSignature)
			encodedSignerInfo := cadesCounterSignature.SignerInformation().DER()
			if bytes.Equal(encodedSI, encodedSignerInfo) {
				cadesResult = append(cadesResult, counterSignature)
			}
		}
	}
	return cadesResult
}

// GetCounterSignatureReferences is CAdES's override of the base's concrete
// getCounterSignatureReferences(AdvancedSignature) hook.
//
// The reason to override: CAdES counter signature does not have a private SignedData
// certificates/revocations.
//
// See the file header GAP note: not yet reachable from the base's own internal
// getCounterSignaturesReferences (plural) until SignatureTimestampSourceOverrides gains this
// method; this file's own counterSignaturesReferences above calls it directly instead.
func (c *CAdESTimestampSource) GetCounterSignatureReferences(counterSignature validation.AdvancedSignature) []*validation.TimestampedReference {
	var counterSigReferences []*validation.TimestampedReference

	counterSigReferences = append(counterSigReferences,
		validation.NewTimestampedReference(counterSignature.ID(), enumerations.TimestampedObjectType_SIGNATURE))

	signingCertificateRefs := counterSignature.CertificateSource().SigningCertificateRefs()
	cadesTSAddReferences(&counterSigReferences, timestamp.CreateReferencesForCertificateRefs(
		signingCertificateRefs, counterSignature.CertificateSource(), c.CertificateSource()))

	counterSignatureTimestampSource := counterSignature.TimestampSource()
	cadesTSAddReferences(&counterSigReferences, counterSignatureTimestampSource.SignerDataReferences())
	cadesTSAddReferences(&counterSigReferences, counterSignatureTimestampSource.UnsignedPropertiesReferences())
	cadesTSAddReferences(&counterSigReferences, c.encapsulatedReferencesFromTimestamps(counterSignatureTimestampSource.AllTimestamps()))

	return counterSigReferences
}

// -----------------------------------------------------------------------------
// Small reference-list helpers, mirroring the frozen spi/validation/timestamp package's own
// unexported addReference/addReferences/timestampAddReferences (abstract_timestamp_source.go),
// which this file cannot call directly (different package). Prefixed distinctly (cadesTS) to
// avoid colliding with any sibling chunk writing to this same package concurrently.
// -----------------------------------------------------------------------------

// cadesTSAddReference adds referenceToAdd to *referenceList without duplicates.
// Port of the AbstractTimestampSource addReference(List, TimestampedReference) helper.
func cadesTSAddReference(referenceList *[]*validation.TimestampedReference, referenceToAdd *validation.TimestampedReference) {
	cadesTSAddReferences(referenceList, []*validation.TimestampedReference{referenceToAdd})
}

// cadesTSAddReferences adds referencesToAdd to *referenceList without duplicates (by
// TimestampedReference.Equals). Port of the AbstractTimestampSource addReferences(List, List)
// helper.
func cadesTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
	for _, candidate := range referencesToAdd {
		found := false
		for _, existing := range *referenceList {
			if existing.Equals(candidate) {
				found = true
				break
			}
		}
		if !found {
			*referenceList = append(*referenceList, candidate)
		}
	}
}

// cadesTSTimestampAddReferences enriches timestampToken's TimestampedReferences with
// referencesToAdd, without duplicates, via TimestampToken.SetTimestampedReferences - the Go
// counterpart of Java's `addReferences(timestampToken.getTimestampedReferences(), ...)` List
// reference-semantics mutation. Mirrors the frozen timestamp package's own unexported
// timestampAddReferences.
func cadesTSTimestampAddReferences(timestampToken *validation.TimestampToken, referencesToAdd []*validation.TimestampedReference) {
	merged := append([]*validation.TimestampedReference(nil), timestampToken.TimestampedReferences()...)
	cadesTSAddReferences(&merged, referencesToAdd)
	timestampToken.SetTimestampedReferences(merged)
}

// cadesTSMust panics with a model.DSSError wrapping err when err is non-nil, otherwise returning
// value unchanged. Mirrors spi/validation/timestamp/abstract_timestamp_source.go's own unexported
// `must` helper (this file cannot call it directly, different package): the boundary this file
// uses wherever a method must satisfy a frozen, error-free Java signature (no throws clause) but
// internally calls a helper this port made fallible.
func cadesTSMust[T any](value T, err error) T {
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return value
}

// compile-time assertion: *CAdESTimestampSource implements
// timestamp.SignatureTimestampSourceOverrides[*CAdESSignature, *CAdESAttribute], matching Java's
// "extends SignatureTimestampSource<CAdESSignature, CAdESAttribute>". Will not compile until
// *CAdESSignature/*CAdESAttribute land (forward dependency of a sibling chunk); see PORTING.md's
// "cades chunks may NOT build mid-port" note.
var _ timestamp.SignatureTimestampSourceOverrides[*CAdESSignature, *CAdESAttribute] = (*CAdESTimestampSource)(nil)
