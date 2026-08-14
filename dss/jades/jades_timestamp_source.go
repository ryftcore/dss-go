// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/timestamp/JAdESTimestampSource.java (DSS 6.5.RC1).
//
// SCC flattening: Java's eu.europa.esig.dss.jades.validation.timestamp package folds into this
// one Go package per the phase-6 package layout (S6_BRIEF.md).
//
// # DEVIATION: SA generic parameter is *EtsiUComponent, not *JAdESAttribute
//
// Java's JAdESTimestampSource extends SignatureTimestampSource<JAdESSignature, JAdESAttribute>:
// buildUnsignedSignatureProperties() returns the JAdESEtsiUHeader itself, via an unchecked raw-type
// cast `(SignatureProperties) signature.getEtsiUHeader()` - JAdESEtsiUHeader actually implements
// SignatureProperties<EtsiUComponent>, and Java's type erasure lets that pass at runtime, later
// relying on getCounterSignatures' own `unsignedAttribute instanceof EtsiUComponent` guard to
// recover the concrete type.
//
// Go generics are invariant with no type erasure: timestamp.SignatureTimestampSource[AS, SA] can
// only be instantiated with ONE SA for both signed and unsigned properties, and *JAdESEtsiUHeader
// already satisfies validation.SignatureProperties[*EtsiUComponent] exactly (see
// jades_etsi_u_header.go) - there is no unchecked-cast escape hatch to pretend it is
// SignatureProperties[*JAdESAttribute] instead. This file therefore instantiates the generic base
// with SA = *EtsiUComponent (a strict widening: every *EtsiUComponent already satisfies
// validation.SignatureAttribute via its embedded JAdESAttribute), and BuildSignedSignatureProperties
// wraps the *JAdESAttribute-typed JAdESSignedProperties in the small
// jadesSignedPropertiesAsEtsiUComponents adapter below so both sides of the interface line up.
// Every signed-attribute predicate this file implements only ever reads promoted JAdESAttribute
// state (HeaderName/Value), never EtsiUComponent-specific state (Component/IsBase64UrlEncoded), so
// the adapter's synthetic EtsiUComponent wrapping (Component=nil, base64UrlEncoded=false) is never
// observed. GetCounterSignatures below is simpler than Java's for this same reason: unsignedAttribute
// already IS a *EtsiUComponent, with no `instanceof` check needed.
//
// FORWARD DEPENDENCY: *JAdESSignature - see abstract_jws_document_analyzer.go's file header. This
// file additionally needs:
//
//	func (s *JAdESSignature) Jws() *JWS                              // getJws() (already used elsewhere)
//	func (s *JAdESSignature) EtsiUHeader() *JAdESEtsiUHeader         // getEtsiUHeader() (already used elsewhere)
//	func (s *JAdESSignature) MasterCSigComponent() *EtsiUComponent  // (already referenced by dss_json_utils.go's DSSJsonUtilsExtractJAdESCounterSignature via SetMasterCSigComponent)
//
// FORWARD DEPENDENCY: JAdESSignedProperties (Java eu.europa.esig.dss.jades.validation.
// JAdESSignedProperties) is not in this manifest. Assumed shape, matching its Java source (read
// for accuracy, out of this manifest's scope):
//
//	type JAdESSignedProperties struct { ... } // implements validation.SignatureProperties[*JAdESAttribute]
//	func NewJAdESSignedProperties(headers *jose.Headers) *JAdESSignedProperties
//
// FORWARD DEPENDENCY: JAdESCertificateRefExtractionUtils / JAdESRevocationRefExtractionUtils
// (Java eu.europa.esig.dss.jades.validation, static utility classes, flattened per PORTING.md's
// "<JavaClass><MethodName>" convention). Assumed shapes:
//
//	func JAdESCertificateRefExtractionUtilsCreateCertificateRef(certId *jose.Object) *spi.CertificateRef
//	func JAdESRevocationRefExtractionUtilsCreateCRLRef(crlRefMap *jose.Object) *spi.CRLRef
//	func JAdESRevocationRefExtractionUtilsCreateOCSPRef(ocspRefMap *jose.Object) *spi.OCSPRef
package jades

import (
	"github.com/utain/esig/dss/crlparser"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/jose"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/timestamp"
	"github.com/utain/esig/dss/utils"
)

// jadesSignedPropertiesAsEtsiUComponents adapts a validation.SignatureProperties[*JAdESAttribute]
// (JAdESSignedProperties) to validation.SignatureProperties[*EtsiUComponent], boxing each
// *JAdESAttribute the wrapped value returns into a synthetic *EtsiUComponent. See the file
// header's "DEVIATION" note.
type jadesSignedPropertiesAsEtsiUComponents struct {
	signedProperties validation.SignatureProperties[*JAdESAttribute]
}

// IsExist implements validation.SignatureProperties.
func (p *jadesSignedPropertiesAsEtsiUComponents) IsExist() bool {
	return p.signedProperties.IsExist()
}

// Attributes implements validation.SignatureProperties, boxing each *JAdESAttribute.
func (p *jadesSignedPropertiesAsEtsiUComponents) Attributes() []*EtsiUComponent {
	attributes := p.signedProperties.Attributes()
	result := make([]*EtsiUComponent, len(attributes))
	for i, a := range attributes {
		result[i] = &EtsiUComponent{JAdESAttribute: *a}
	}
	return result
}

// JAdESTimestampSource extracts timestamps from a JAdES signature. Port of the class
// JAdESTimestampSource, extending timestamp.SignatureTimestampSource[*JAdESSignature,
// *EtsiUComponent] (see the file header's DEVIATION note on the SA parameter).
//
// @SuppressWarnings("serial")/java.io.Serializable is dropped (no Go counterpart).
type JAdESTimestampSource struct {
	timestamp.SignatureTimestampSource[*JAdESSignature, *EtsiUComponent]

	// signature is being validated. See cades_timestamp_source.go's identical precedent (the
	// embedded base keeps its own private copy; this file needs its own reference too, since Go
	// has no protected field access across packages).
	signature *JAdESSignature

	// timestampAttributeMap maps time-stamp tokens to corresponding JAdES attributes.
	timestampAttributeMap map[*validation.TimestampToken]*EtsiUComponent
}

// NewJAdESTimestampSource is the default constructor. Port of the (JAdESSignature) constructor.
func NewJAdESTimestampSource(signature *JAdESSignature) *JAdESTimestampSource {
	source := &JAdESTimestampSource{
		SignatureTimestampSource: timestamp.NewSignatureTimestampSourceBase[*JAdESSignature, *EtsiUComponent](signature),
		signature:                signature,
		timestampAttributeMap:    make(map[*validation.TimestampToken]*EtsiUComponent),
	}
	source.InitSignatureTimestampSource(source)
	return source
}

// BuildSignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides. Port of
// buildSignedSignatureProperties(). See the file header's DEVIATION note.
func (s *JAdESTimestampSource) BuildSignedSignatureProperties() validation.SignatureProperties[*EtsiUComponent] {
	return &jadesSignedPropertiesAsEtsiUComponents{signedProperties: NewJAdESSignedProperties(s.signature.Jws().Headers())}
}

// BuildUnsignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides. Port
// of buildUnsignedSignatureProperties(). See the file header's DEVIATION note: no cast is needed,
// unlike Java's unchecked raw-type cast, since *JAdESEtsiUHeader already implements
// validation.SignatureProperties[*EtsiUComponent] exactly.
func (s *JAdESTimestampSource) BuildUnsignedSignatureProperties() validation.SignatureProperties[*EtsiUComponent] {
	return s.signature.EtsiUHeader()
}

// IsContentTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isContentTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsContentTimestamp(signedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesAdoTst == signedAttribute.HeaderName()
}

// IsAllDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported. Port of isAllDataObjectsTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsAllDataObjectsTimestamp(signedAttribute *EtsiUComponent) bool {
	return false
}

// IsIndividualDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported. Port of isIndividualDataObjectsTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsIndividualDataObjectsTimestamp(signedAttribute *EtsiUComponent) bool {
	return false
}

// IsSignatureTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignatureTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsSignatureTimestamp(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesSigTst == unsignedAttribute.HeaderName()
}

// IsCompleteCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteCertificateRef(JAdESAttribute).
func (s *JAdESTimestampSource) IsCompleteCertificateRef(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesXRefs == unsignedAttribute.HeaderName()
}

// IsAttributeCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeCertificateRef(JAdESAttribute).
func (s *JAdESTimestampSource) IsAttributeCertificateRef(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesAxRefs == unsignedAttribute.HeaderName()
}

// IsCompleteRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteRevocationRef(JAdESAttribute).
func (s *JAdESTimestampSource) IsCompleteRevocationRef(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesRRefs == unsignedAttribute.HeaderName()
}

// IsAttributeRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeRevocationRef(JAdESAttribute).
func (s *JAdESTimestampSource) IsAttributeRevocationRef(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesArRefs == unsignedAttribute.HeaderName()
}

// IsRefsOnlyTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRefsOnlyTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsRefsOnlyTimestamp(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesRfsTst == unsignedAttribute.HeaderName()
}

// IsSigAndRefsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSigAndRefsTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsSigAndRefsTimestamp(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesSigRTst == unsignedAttribute.HeaderName()
}

// IsCertificateValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCertificateValues(JAdESAttribute).
func (s *JAdESTimestampSource) IsCertificateValues(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesXVals == unsignedAttribute.HeaderName()
}

// IsRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRevocationValues(JAdESAttribute).
func (s *JAdESTimestampSource) IsRevocationValues(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesRVals == unsignedAttribute.HeaderName()
}

// IsArchiveTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isArchiveTimestamp(JAdESAttribute).
func (s *JAdESTimestampSource) IsArchiveTimestamp(unsignedAttribute *EtsiUComponent) bool {
	return jadesTimestampSourceIsArchiveTimestamp(unsignedAttribute.HeaderName())
}

// jadesTimestampSourceIsArchiveTimestamp ports the private isArchiveTimestamp(String).
func jadesTimestampSourceIsArchiveTimestamp(headerName string) bool {
	return JAdESHeaderParameterNamesArcTst == headerName
}

// IsTimeStampValidationData implements timestamp.SignatureTimestampSourceOverrides.
// Port of isTimeStampValidationData(JAdESAttribute).
func (s *JAdESTimestampSource) IsTimeStampValidationData(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesTstVd == unsignedAttribute.HeaderName()
}

// IsAnyValidationData implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAnyValidationData(JAdESAttribute).
func (s *JAdESTimestampSource) IsAnyValidationData(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesAnyValData == unsignedAttribute.HeaderName()
}

// IsValidationDataReferences implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported. Port of isValidationDataReferences(JAdESAttribute).
func (s *JAdESTimestampSource) IsValidationDataReferences(unsignedAttribute *EtsiUComponent) bool {
	return false
}

// IsCounterSignature implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCounterSignature(JAdESAttribute).
func (s *JAdESTimestampSource) IsCounterSignature(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesCSig == unsignedAttribute.HeaderName()
}

// IsSignaturePolicyStore implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignaturePolicyStore(JAdESAttribute).
func (s *JAdESTimestampSource) IsSignaturePolicyStore(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesSigPst == unsignedAttribute.HeaderName()
}

// IsAttrAuthoritiesCertValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttrAuthoritiesCertValues(JAdESAttribute).
func (s *JAdESTimestampSource) IsAttrAuthoritiesCertValues(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesAxVals == unsignedAttribute.HeaderName()
}

// IsAttributeRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeRevocationValues(JAdESAttribute).
func (s *JAdESTimestampSource) IsAttributeRevocationValues(unsignedAttribute *EtsiUComponent) bool {
	return JAdESHeaderParameterNamesArVals == unsignedAttribute.HeaderName()
}

// IsEvidenceRecord implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported. Port of isEvidenceRecord(JAdESAttribute).
func (s *JAdESTimestampSource) IsEvidenceRecord(unsignedAttribute *EtsiUComponent) bool {
	return false
}

// # GAP flagged for integrator: getSignatureTimestampReferences is a concrete-but-overridable
// # base method with no override hook
//
// Java's JAdESTimestampSource overrides the base's protected getSignatureTimestampReferences(),
// adding getKeyInfoReferences() to the inherited result:
//
//	protected List<TimestampedReference> getSignatureTimestampReferences() {
//	    List<TimestampedReference> timestampedReferences = super.getSignatureTimestampReferences();
//	    addReferences(timestampedReferences, getKeyInfoReferences());
//	    return timestampedReferences;
//	}
//
// getSignatureTimestampReferences() is `protected`, unexported in the Go base
// (spi/validation/timestamp/signature_timestamp_source.go's own getSignatureTimestampReferences),
// hence uncallable and unshadowable from this different package - exactly the same structural gap
// cades_timestamp_source.go and xades_timestamp_source.go already flag and work around for their
// own formats (see cades_timestamp_source.go's file header for the precedent this follows).
// signatureTimestampReferences below reimplements it from exported base accessors, then
// IncorporateArchiveTimestampReferences is shadowed purely to route the base's own
// incorporateArchiveTimestampReferences flow through it (archiveTimestampReferences), matching
// the actually-observable behaviour Java's virtual dispatch produces even though
// JAdESTimestampSource.java itself declares no such IncorporateArchiveTimestampReferences
// override.
func (s *JAdESTimestampSource) IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken,
	previousTimestamps []*validation.TimestampToken) {
	jadesTSTimestampAddReferences(timestampToken, s.archiveTimestampReferences(previousTimestamps))
}

// archiveTimestampReferences reimplements the base's private getArchiveTimestampReferences(List),
// which this file cannot call directly, from exported base accessors. Port of the base's own
// getArchiveTimestampReferences(List<TimestampToken>) body (JAdES adds nothing beyond what
// signatureTimestampReferences below already folds in).
func (s *JAdESTimestampSource) archiveTimestampReferences(previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	var timestampedReferences []*validation.TimestampedReference
	jadesTSAddReferences(&timestampedReferences, s.signatureTimestampReferences())
	jadesTSAddReferences(&timestampedReferences, s.encapsulatedReferencesFromTimestamps(previousTimestamps))
	jadesTSAddReferences(&timestampedReferences, s.UnsignedPropertiesReferences())
	jadesTSAddReferences(&timestampedReferences, s.GetKeyInfoReferences())
	return timestampedReferences
}

// signatureTimestampReferences reimplements the base's private getSignatureTimestampReferences(),
// then applies JAdES's own getKeyInfoReferences() addition. Port of the protected
// getSignatureTimestampReferences() override.
func (s *JAdESTimestampSource) signatureTimestampReferences() []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	jadesTSAddReferences(&references, s.encapsulatedReferencesFromTimestamps(s.ContentTimestamps()))
	jadesTSAddReferences(&references, s.SignerDataReferences())
	jadesTSAddReference(&references, validation.NewTimestampedReference(s.signature.ID(), enumerations.TimestampedObjectType_SIGNATURE))
	jadesTSAddReferences(&references, s.signingCertificateTimestampReferences())
	jadesTSAddReferences(&references, s.GetKeyInfoReferences())
	return references
}

// signingCertificateTimestampReferences reimplements the base's private
// getSigningCertificateTimestampReferences().
func (s *JAdESTimestampSource) signingCertificateTimestampReferences() []*validation.TimestampedReference {
	signatureCertificateSource := s.signature.CertificateSource()
	return timestamp.CreateReferencesForCertificateRefs(
		signatureCertificateSource.SigningCertificateRefs(), signatureCertificateSource, s.CertificateSource())
}

// encapsulatedReferencesFromTimestamps reimplements the base's private
// getEncapsulatedReferencesFromTimestamps(List), from the exported timestamp.ReferencesFromTimestamp
// and this file's merged-source accessors.
func (s *JAdESTimestampSource) encapsulatedReferencesFromTimestamps(timestampTokens []*validation.TimestampToken) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	for _, timestampToken := range timestampTokens {
		refs, err := timestamp.ReferencesFromTimestamp(timestampToken, s.CertificateSource(), s.CRLSource(), s.OCSPSource())
		if err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
		jadesTSAddReferences(&references, refs)
	}
	return references
}

// -----------------------------------------------------------------------------
// Small reference-list helpers, mirroring the frozen spi/validation/timestamp package's own
// unexported addReference/addReferences/timestampAddReferences (abstract_timestamp_source.go),
// which this file cannot call directly (different package). Prefixed distinctly (jadesTS) to
// avoid colliding with any sibling chunk writing to this same package concurrently, mirroring
// cades_timestamp_source.go's identical cadesTS-prefixed set.
// -----------------------------------------------------------------------------

// jadesTSAddReference adds referenceToAdd to *referenceList without duplicates.
func jadesTSAddReference(referenceList *[]*validation.TimestampedReference, referenceToAdd *validation.TimestampedReference) {
	jadesTSAddReferences(referenceList, []*validation.TimestampedReference{referenceToAdd})
}

// jadesTSAddReferences adds referencesToAdd to *referenceList without duplicates (by
// TimestampedReference.Equals).
func jadesTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
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

// jadesTSTimestampAddReferences enriches timestampToken's TimestampedReferences with
// referencesToAdd, without duplicates, via TimestampToken.SetTimestampedReferences.
func jadesTSTimestampAddReferences(timestampToken *validation.TimestampToken, referencesToAdd []*validation.TimestampedReference) {
	merged := append([]*validation.TimestampedReference(nil), timestampToken.TimestampedReferences()...)
	jadesTSAddReferences(&merged, referencesToAdd)
	timestampToken.SetTimestampedReferences(merged)
}

// GetCertificateRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCertificateRefs(JAdESAttribute).
func (s *JAdESTimestampSource) GetCertificateRefs(unsignedAttribute *EtsiUComponent) []*spi.CertificateRef {
	var result []*spi.CertificateRef
	certRefs := DSSJsonUtilsToListValue(unsignedAttribute.Value())
	if utils.IsCollectionNotEmpty(certRefs) {
		for _, item := range certRefs {
			certId := DSSJsonUtilsToMap(item, JAdESHeaderParameterNamesCertId)
			if certId != nil && certId.Size() != 0 {
				if certificateRef := JAdESCertificateRefExtractionUtilsCreateCertificateRef(certId); certificateRef != nil {
					result = append(result, certificateRef)
				}
			}
		}
	}
	return result
}

// GetCRLRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCRLRefs(JAdESAttribute).
func (s *JAdESTimestampSource) GetCRLRefs(unsignedAttribute *EtsiUComponent) []*spi.CRLRef {
	var result []*spi.CRLRef
	refsValueMap := DSSJsonUtilsToMapValue(unsignedAttribute.Value())
	if refsValueMap != nil && refsValueMap.Size() != 0 {
		crlRefs := DSSJsonUtilsGetAsList(refsValueMap, JAdESHeaderParameterNamesCrlRefs)
		for _, item := range crlRefs {
			crlRefMap := DSSJsonUtilsToMapValue(item)
			if crlRefMap != nil && crlRefMap.Size() != 0 {
				if crlRef := JAdESRevocationRefExtractionUtilsCreateCRLRef(crlRefMap); crlRef != nil {
					result = append(result, crlRef)
				}
			}
		}
	}
	return result
}

// GetOCSPRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getOCSPRefs(JAdESAttribute).
func (s *JAdESTimestampSource) GetOCSPRefs(unsignedAttribute *EtsiUComponent) []*spi.OCSPRef {
	var result []*spi.OCSPRef
	refsValueMap := DSSJsonUtilsToMapValue(unsignedAttribute.Value())
	if refsValueMap != nil && refsValueMap.Size() != 0 {
		ocsp := DSSJsonUtilsGetAsList(refsValueMap, JAdESHeaderParameterNamesOcspRefs)
		for _, item := range ocsp {
			ocspRefMap := DSSJsonUtilsToMapValue(item)
			if ocspRefMap != nil && ocspRefMap.Size() != 0 {
				if ocspRef := JAdESRevocationRefExtractionUtilsCreateOCSPRef(ocspRefMap); ocspRef != nil {
					result = append(result, ocspRef)
				}
			}
		}
	}
	return result
}

// GetEncapsulatedCertificateIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCertificateIdentifiers(JAdESAttribute).
func (s *JAdESTimestampSource) GetEncapsulatedCertificateIdentifiers(unsignedAttribute *EtsiUComponent) []model.Identifier {
	var xVals []any
	switch {
	case s.IsTimeStampValidationData(unsignedAttribute):
		tstVd := DSSJsonUtilsToMap(unsignedAttribute.Value(), JAdESHeaderParameterNamesTstVd)
		if tstVd != nil && tstVd.Size() != 0 {
			xVals = DSSJsonUtilsGetAsList(tstVd, JAdESHeaderParameterNamesXVals)
		}
	case s.IsAnyValidationData(unsignedAttribute):
		anyVD := DSSJsonUtilsToMap(unsignedAttribute.Value(), JAdESHeaderParameterNamesAnyValData)
		if anyVD != nil && anyVD.Size() != 0 {
			xVals = DSSJsonUtilsGetAsList(anyVD, JAdESHeaderParameterNamesXVals)
		}
	default:
		xVals = DSSJsonUtilsToList(unsignedAttribute.Value(), JAdESHeaderParameterNamesXVals)
	}

	if utils.IsCollectionEmpty(xVals) {
		return nil
	}
	var certificateIdentifiers []model.Identifier
	for _, encapsulatedCert := range xVals {
		if certificateToken := jadesTimestampSourceToCertificateToken(encapsulatedCert); certificateToken != nil {
			certificateIdentifiers = append(certificateIdentifiers, certificateToken.DSSID())
		}
	}
	return certificateIdentifiers
}

// jadesTimestampSourceToCertificateToken ports the private toCertificateToken(Object),
// swallowing any parse failure and returning nil, matching the Java catch-all.
func jadesTimestampSourceToCertificateToken(encapsulatedCert any) (result *model.CertificateToken) {
	defer func() {
		if recover() != nil {
			// Upstream logs "An error occurred during parsing a certificate. Reason : {}".
			result = nil
		}
	}()

	m := DSSJsonUtilsToMapValue(encapsulatedCert)
	if m == nil || m.Size() == 0 {
		return nil
	}
	x509Cert := DSSJsonUtilsGetAsMap(m, JAdESHeaderParameterNamesX509Cert)
	otherCert := DSSJsonUtilsGetAsMap(m, JAdESHeaderParameterNamesOtherCert)
	if x509Cert != nil && x509Cert.Size() != 0 {
		base64Cert := DSSJsonUtilsGetAsString(x509Cert, JAdESHeaderParameterNamesVal)
		if utils.IsStringNotBlank(base64Cert) {
			binaries := utils.FromBase64(base64Cert)
			token, err := spi.DSSUtilsLoadCertificateFromBinary(binaries)
			if err != nil {
				return nil
			}
			return token
		}
	} else if otherCert != nil && otherCert.Size() != 0 {
		// Upstream logs "The header '{}' is not supported! The entry is skipped.".
	}
	return nil
}

// GetEncapsulatedCRLIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCRLIdentifiers(JAdESAttribute).
func (s *JAdESTimestampSource) GetEncapsulatedCRLIdentifiers(unsignedAttribute *EtsiUComponent) []*crlparser.CRLBinary {
	rVals := jadesTimestampSourceRVals(s, unsignedAttribute)
	if rVals == nil {
		return nil
	}
	var crlIdentifiers []*crlparser.CRLBinary
	crlVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesCrlVals)
	for _, item := range crlVals {
		if crlBinary := jadesTimestampSourceToCRLBinary(item); crlBinary != nil {
			crlIdentifiers = append(crlIdentifiers, crlBinary)
		}
	}
	return crlIdentifiers
}

// jadesTimestampSourceRVals ports the `rVals` local-variable computation shared by
// getEncapsulatedCRLIdentifiers/getEncapsulatedOCSPIdentifiers.
func jadesTimestampSourceRVals(s *JAdESTimestampSource, unsignedAttribute *EtsiUComponent) *jose.Object {
	switch {
	case s.IsTimeStampValidationData(unsignedAttribute):
		tstVd := DSSJsonUtilsToMap(unsignedAttribute.Value(), JAdESHeaderParameterNamesTstVd)
		if tstVd != nil && tstVd.Size() != 0 {
			return DSSJsonUtilsGetAsMap(tstVd, JAdESHeaderParameterNamesRVals)
		}
		return nil
	case s.IsAnyValidationData(unsignedAttribute):
		anyVD := DSSJsonUtilsToMap(unsignedAttribute.Value(), JAdESHeaderParameterNamesAnyValData)
		if anyVD != nil && anyVD.Size() != 0 {
			return DSSJsonUtilsGetAsMap(anyVD, JAdESHeaderParameterNamesRVals)
		}
		return nil
	default:
		return DSSJsonUtilsToMap(unsignedAttribute.Value(), JAdESHeaderParameterNamesRVals)
	}
}

// jadesTimestampSourceToCRLBinary ports the private toCRLBinary(Object), swallowing any parse
// failure and returning nil, matching the Java catch-all.
func jadesTimestampSourceToCRLBinary(crlVal any) *crlparser.CRLBinary {
	encapsulatedCrl := DSSJsonUtilsToMapValue(crlVal)
	if encapsulatedCrl == nil || encapsulatedCrl.Size() == 0 {
		return nil
	}
	base64Crl := DSSJsonUtilsGetAsString(encapsulatedCrl, JAdESHeaderParameterNamesVal)
	if !utils.IsStringNotBlank(base64Crl) {
		return nil
	}
	binaries := utils.FromBase64(base64Crl)
	crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(binaries)
	if err != nil {
		// Upstream logs "An error occurred during parsing a CRL. Reason : {}".
		return nil
	}
	return crlBinary
}

// GetEncapsulatedOCSPIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedOCSPIdentifiers(JAdESAttribute).
func (s *JAdESTimestampSource) GetEncapsulatedOCSPIdentifiers(unsignedAttribute *EtsiUComponent) []*spi.OCSPResponseBinary {
	rVals := jadesTimestampSourceRVals(s, unsignedAttribute)
	if rVals == nil {
		return nil
	}
	var ocspIdentifiers []*spi.OCSPResponseBinary
	ocspVals := DSSJsonUtilsGetAsList(rVals, JAdESHeaderParameterNamesOcspVals)
	for _, item := range ocspVals {
		if ocspResponseBinary := jadesTimestampSourceToOCSPResponseBinary(item); ocspResponseBinary != nil {
			ocspIdentifiers = append(ocspIdentifiers, ocspResponseBinary)
		}
	}
	return ocspIdentifiers
}

// jadesTimestampSourceToOCSPResponseBinary ports the private toOCSPResponseBinary(Object),
// swallowing any parse failure and returning nil, matching the Java catch-all (which logs "An
// error occurred during parsing a CRL. Reason : {}" - reproduced verbatim, including the
// upstream copy/paste "CRL" wording, since this is only a dropped log message).
func jadesTimestampSourceToOCSPResponseBinary(ocspVal any) *spi.OCSPResponseBinary {
	encapsulatedOcsp := DSSJsonUtilsToMapValue(ocspVal)
	if encapsulatedOcsp == nil || encapsulatedOcsp.Size() == 0 {
		return nil
	}
	base64Ocsp := DSSJsonUtilsGetAsString(encapsulatedOcsp, JAdESHeaderParameterNamesVal)
	if !utils.IsStringNotBlank(base64Ocsp) {
		return nil
	}
	binaries := utils.FromBase64(base64Ocsp)
	basicOCSPResp, err := spi.DSSRevocationUtilsLoadOCSPFromBinaries(binaries)
	if err != nil {
		return nil
	}
	binary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
	if err != nil {
		return nil
	}
	return binary
}

// GetTimestampMessageImprintDigestBuilderForToken implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(TimestampToken) override.
func (s *JAdESTimestampSource) GetTimestampMessageImprintDigestBuilderForToken(
	timestampToken *validation.TimestampToken) timestamp.TimestampMessageDigestBuilder {
	var timestampAttribute *JAdESAttribute
	if etsiUComponent := s.timestampAttributeMap[timestampToken]; etsiUComponent != nil {
		timestampAttribute = &etsiUComponent.JAdESAttribute
	}
	return NewJAdESTimestampMessageDigestBuilderForToken(s.signature, timestampToken).
		SetTimestampAttribute(timestampAttribute)
}

// GetTimestampMessageImprintDigestBuilderForAlgorithm implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(DigestAlgorithm) override.
func (s *JAdESTimestampSource) GetTimestampMessageImprintDigestBuilderForAlgorithm(
	digestAlgorithm enumerations.DigestAlgorithm) timestamp.TimestampMessageDigestBuilder {
	return NewJAdESTimestampMessageDigestBuilder(s.signature, digestAlgorithm)
}

// GetCounterSignatures implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCounterSignatures(JAdESAttribute). See the file header's DEVIATION note: no
// `instanceof EtsiUComponent` check is needed, unsignedAttribute already is one.
func (s *JAdESTimestampSource) GetCounterSignatures(unsignedAttribute *EtsiUComponent) []validation.AdvancedSignature {
	counterSignatures := s.signature.CounterSignatures()
	for _, counterSignature := range counterSignatures {
		// Java casts unconditionally: `(JAdESSignature) counterSignature`. Every concrete
		// AdvancedSignature a JAdES source's CounterSignatures() hands out is a *JAdESSignature.
		jadesCounterSignature := counterSignature.(*JAdESSignature)
		if unsignedAttribute == jadesCounterSignature.MasterCSigComponent() {
			// NOTE: only one counter signature is allowed within the CounterSignature unprotected
			// header.
			return []validation.AdvancedSignature{counterSignature}
		}
	}
	return nil
}

// GetSignatureTimestampData returns the message-imprint digest for a SignatureTimestamp
// (BASE64URL(JWS Signature Value)). Port of the public getSignatureTimestampData(DigestAlgorithm).
func (s *JAdESTimestampSource) GetSignatureTimestampData(digestAlgorithm enumerations.DigestAlgorithm) model.DSSMessageDigest {
	builder := NewJAdESTimestampMessageDigestBuilder(s.signature, digestAlgorithm)
	return builder.SignatureTimestampMessageDigest()
}

// GetArchiveTimestampData returns the message-imprint digest for an ArchiveTimestamp. Port of the
// public getArchiveTimestampData(DigestAlgorithm, String).
func (s *JAdESTimestampSource) GetArchiveTimestampData(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string) model.DSSMessageDigest {
	builder := NewJAdESTimestampMessageDigestBuilder(s.signature, digestAlgorithm).
		SetCanonicalizationAlgorithm(canonicalizationMethod)
	return builder.ArchiveTimestampMessageDigest()
}

// MakeTimestampToken implements timestamp.SignatureTimestampSourceOverrides. Port of
// makeTimestampToken(JAdESAttribute, TimestampType, List): unconditionally unsupported, since an
// 'etsiU' timestamp-container attribute can contain more than one timestamp token - see
// MakeTimestampTokens (plural) below, which JAdES uses instead.
func (s *JAdESTimestampSource) MakeTimestampToken(signatureAttribute *EtsiUComponent, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) *validation.TimestampToken {
	panic("Attribute can contain more than one timestamp")
}

// MakeTimestampTokens is JAdES's own plural timestamp-token factory. Port of the protected
// makeTimestampTokens(JAdESAttribute, TimestampType, List) override.
func (s *JAdESTimestampSource) MakeTimestampTokens(signatureAttribute *EtsiUComponent, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) []*validation.TimestampToken {
	if enumerations.TimestampType_ARCHIVE_TIMESTAMP == timestampType {
		return s.extractArchiveTimestampTokens(signatureAttribute, references)
	}
	tstContainer := DSSJsonUtilsToMap(signatureAttribute.Value(), JAdESHeaderParameterNamesTstContainer)
	return s.extractTimestampTokens(signatureAttribute, tstContainer, timestampType, references)
}

// extractTimestampTokens ports the private extractTimestampTokens(JAdESAttribute, Map,
// TimestampType, List).
func (s *JAdESTimestampSource) extractTimestampTokens(signatureAttribute *EtsiUComponent, tstContainer *jose.Object,
	timestampType enumerations.TimestampType, references []*validation.TimestampedReference) []*validation.TimestampToken {
	var result []*validation.TimestampToken
	if tstContainer != nil && tstContainer.Size() != 0 {
		tstTokens := DSSJsonUtilsGetAsList(tstContainer, JAdESHeaderParameterNamesTstTokens)
		if utils.IsCollectionNotEmpty(tstTokens) {
			for i, tstToken := range tstTokens {
				timestampToken := s.toTimestampToken(tstToken, signatureAttribute, i, timestampType, references)
				if timestampToken != nil {
					s.timestampAttributeMap[timestampToken] = signatureAttribute
					result = append(result, timestampToken)
				}
			}
		} else {
			// Upstream logs "'{}' element is not found! Returns an empty array if timestamps.".
		}
	}
	return result
}

// toTimestampToken ports the private toTimestampToken(Object, JAdESAttribute, Integer,
// TimestampType, List).
func (s *JAdESTimestampSource) toTimestampToken(tstToken any, signatureAttribute *EtsiUComponent, orderWithinAttribute int,
	timestampType enumerations.TimestampType, references []*validation.TimestampedReference) *validation.TimestampToken {
	tstTokenMap := DSSJsonUtilsToMapValue(tstToken)
	if tstTokenMap == nil || tstTokenMap.Size() == 0 {
		return nil
	}
	encoding := DSSJsonUtilsGetAsString(tstTokenMap, JAdESHeaderParameterNamesEncoding)
	if encoding != "" && encoding != enumerations.PKIEncoding_DER.URI() {
		// Upstream logs "Unsupported encoding {}".
		return nil
	}
	tstBase64 := DSSJsonUtilsGetAsString(tstTokenMap, JAdESHeaderParameterNamesVal)
	if tstBase64 == "" {
		return nil
	}

	binaries := utils.FromBase64(tstBase64)
	orderWithin := orderWithinAttribute
	identifierBuilder := timestamp.NewSignatureTimestampIdentifierBuilder(binaries).
		SetSignature(s.signature).
		SetAttribute(signatureAttribute).
		SetOrderOfAttribute(s.GetAttributeOrder(signatureAttribute)).
		SetOrderWithinAttribute(&orderWithin)
	token, err := validation.NewTimestampTokenWithIdentifierBuilder(binaries, timestampType, references, identifierBuilder)
	if err != nil {
		// Upstream logs "Unable to create timestamp from base64-encoded string '{}'. Reason :
		// {}".
		return nil
	}
	return token
}

// extractArchiveTimestampTokens ports the private extractArchiveTimestampTokens(JAdESAttribute,
// List).
func (s *JAdESTimestampSource) extractArchiveTimestampTokens(signatureAttribute *EtsiUComponent,
	references []*validation.TimestampedReference) []*validation.TimestampToken {
	arcTst := DSSJsonUtilsToMap(signatureAttribute.Value(), JAdESHeaderParameterNamesArcTst)
	return s.extractTimestampTokens(signatureAttribute, arcTst, enumerations.TimestampType_ARCHIVE_TIMESTAMP, references)
}

// GetArchiveTimestampType implements timestamp.SignatureTimestampSourceOverrides.
// Port of getArchiveTimestampType(JAdESAttribute).
func (s *JAdESTimestampSource) GetArchiveTimestampType(unsignedAttribute *EtsiUComponent) enumerations.ArchiveTimestampType {
	return enumerations.ArchiveTimestampType_JAdES
}

// MakeEvidenceRecords implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: embedded evidence records are not supported within JAdES format. Port of
// makeEvidenceRecords(JAdESAttribute, List).
func (s *JAdESTimestampSource) MakeEvidenceRecords(signatureAttribute *EtsiUComponent,
	references []*validation.TimestampedReference) []validation.EvidenceRecord {
	// Upstream logs "Embedded evidence records are not supported within JAdES format! The
	// unsigned attribute is skipped." when signatureAttribute != nil.
	return nil
}

// compile-time assertion: *JAdESTimestampSource implements
// timestamp.SignatureTimestampSourceOverrides[*JAdESSignature, *EtsiUComponent] (see the file
// header's DEVIATION note on why the SA parameter differs from Java's literal type argument).
// Will not compile until JAdESSignature/JAdESSignedProperties/JAdESCertificateRefExtractionUtils/
// JAdESRevocationRefExtractionUtils land (forward dependencies of sibling chunks); see
// PORTING.md's "chunks may NOT build mid-port" rule.
var _ timestamp.SignatureTimestampSourceOverrides[*JAdESSignature, *EtsiUComponent] = (*JAdESTimestampSource)(nil)
