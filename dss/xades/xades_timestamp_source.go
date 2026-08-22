// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/timestamp/XAdESTimestampSource.java (DSS 6.5.RC1).
//
// # MakeTimestampTokens (plural), not MakeTimestampToken (singular), is the override hook
//
// Java's SignatureTimestampSource#makeTimestampTokens(SA, TimestampType, List) (3-arg, `protected`
// - concrete but overridable) is called unqualified, hence virtually dispatched, from the base's
// own makeTimestampTokensFromSignedAttributes()/makeTimestampTokensFromUnsignedAttributes(). The
// base's plain body wraps the abstract, fully-dispatched makeTimestampToken(SA, TimestampType,
// List) (SINGULAR) in a one-element list; XAdESTimestampSource.java overrides the PLURAL method
// instead (an xades132:XAdESTimeStampType element can carry more than one
// xades132:EncapsulatedTimeStamp, so one signature attribute can produce several TimestampTokens
// at once), and its singular makeTimestampToken override unconditionally throws
// UnsupportedOperationException, matching the fact it is never meant to be called for XAdES.
//
// SignatureTimestampSourceOverrides therefore declares MakeTimestampTokens alongside the singular
// MakeTimestampToken, with a default base implementation on SignatureTimestampSource (promoted by
// embedding to every format that does not shadow it - CAdES, JAdES) that wraps the singular result
// in a one-element slice. This file's MakeTimestampTokens below overrides that default.
//
// # GAP flagged for integrator: getTimestampScopes/getSignatureTimestampReferences/
// # getArchiveTimestampReferences are concrete-but-overridable base methods with no override hook
//
// Three more Java methods this file overrides are declared concrete (not abstract) on the base
// SignatureTimestampSource and relied upon there via ordinary virtual dispatch, exactly the same
// pattern TimestampSource already flags for IncorporateArchiveTimestampReferences/
// GetSignatureSignedDataReferences/GetCounterSignatureReferences (see cades_timestamp_source.go's
// file header for the precedent this follows):
//
//   - getTimestampScopes(TimestampToken) - called unqualified from the base's own
//     validateTimestamps() for every content and archive timestamp. The base's private
//     getTimestampScopes always constructs a bare, XAdES-agnostic
//     validationscope.EncapsulatedTimestampScopeFinder; this file's TimestampScopeFinder
//     override (constructing xades.TimestampScopeFinder instead, to correctly filter
//     IndividualDataObjectsTimestamp includes) is not reached from there. See
//     xades_timestamp_scope_finder.go's own file header for the identical note from that file's
//     perspective.
//   - getSignatureTimestampReferences() and getArchiveTimestampReferences(List<TimestampToken>) -
//     both `protected`, unexported in the Go base (spi/validation/timestamp/
//     signature_timestamp_source.go's getSignatureTimestampReferences/getArchiveTimestampReferences),
//     hence uncallable from this different package at all, let alone overridable. This file
//     reimplements both from exported base accessors (SignerDataReferences, ContentTimestamps,
//     CertificateSource/CRLSource/OCSPSource, UnsignedPropertiesReferences, GetKeyInfoReferences,
//     GetAttributeOrder), mirroring cades_timestamp_source.go's identical
//     signatureTimestampReferences/archiveTimestampReferences reimplementations for the same
//     structural reason.
//
// Since IncorporateArchiveTimestampReferences (part of SignatureTimestampSourceOverrides, promoted
// today from the embedded SignatureTimestampSource[*Signature, *Attribute] since this
// file does not shadow it) itself calls the base's OWN unexported getArchiveTimestampReferences
// (not this file's reimplementation), this file additionally shadows
// IncorporateArchiveTimestampReferences purely to route it through archiveTimestampReferences
// below - even though Java's XAdESTimestampSource.java itself declares no such override (Java's
// virtual dispatch reaches the override of getArchiveTimestampReferences automatically from
// inside the base's own incorporateArchiveTimestampReferences; Go needs the explicit shadow for
// the identical reason TimestampSource needed one, just one level removed).
//
// Every non-abstract Is*/Get*/Make* method below that IS part of SignatureTimestampSourceOverrides
// (all the boolean predicates, MakeEvidenceRecords, GetCertificateRefs/GetCRLRefs/GetOCSPRefs,
// GetEncapsulated*Identifiers, GetArchiveTimestampType, GetCounterSignatures) dispatches correctly
// today - only the four named above need the fixes described.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	mscope "github.com/ryftcore/dss-go/dss/model/scope"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/spi/validation/timestamp"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESTimestampSource is the timestamp source for a XAdES signature. Port of the class
// TimestampSource, extending timestamp.SignatureTimestampSource[*Signature, *Attribute].
//
// @SuppressWarnings("serial")/java.io.Serializable is dropped (no Go counterpart).
type TimestampSource struct {
	timestamp.SignatureTimestampSource[*Signature, *Attribute]

	// signatureElement represents the signature element container. "transient" (Java) has no Go
	// counterpart. Kept alongside the embedded base's own private copy, mirroring
	// TimestampSource.signature's identical precedent (cades_timestamp_source.go): Go has no
	// equivalent of protected field access across packages.
	signatureElement *xmldom.Node

	// xadesPaths is the XAdES XPaths to use.
	xadesPaths definition.XAdESPath

	// signature is being validated. See the signatureElement comment above for why this file keeps
	// its own reference in addition to the embedded base's private copy.
	signature *Signature

	// timestampAttributeMap maps time-stamp tokens to their corresponding XAdES attributes.
	timestampAttributeMap map[*validation.TimestampToken]*Attribute
}

// NewXAdESTimestampSource is the default constructor. Port of the (XAdESSignature) constructor.
func NewTimestampSource(signature *Signature) *TimestampSource {
	source := &TimestampSource{
		SignatureTimestampSource: timestamp.NewSignatureTimestampSourceBase[*Signature, *Attribute](signature),
		signatureElement:         signature.SignatureElement(),
		xadesPaths:               signature.XAdESPaths(),
		signature:                signature,
		timestampAttributeMap:    make(map[*validation.TimestampToken]*Attribute),
	}
	source.InitSignatureTimestampSource(source)
	return source
}

// BuildSignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides.
// Port of buildSignedSignatureProperties().
func (s *TimestampSource) BuildSignedSignatureProperties() validation.SignatureProperties[*Attribute] {
	return SignedDataObjectPropertiesBuild(s.signatureElement, s.xadesPaths)
}

// BuildUnsignedSignatureProperties implements timestamp.SignatureTimestampSourceOverrides.
// Port of buildUnsignedSignatureProperties().
func (s *TimestampSource) BuildUnsignedSignatureProperties() validation.SignatureProperties[*Attribute] {
	return UnsignedSigPropertiesBuild(s.signatureElement, s.xadesPaths)
}

// GetTimestampMessageImprintDigestBuilderForAlgorithm implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(DigestAlgorithm) override.
func (s *TimestampSource) GetTimestampMessageImprintDigestBuilderForAlgorithm(
	digestAlgorithm enumerations.DigestAlgorithm) timestamp.MessageDigestBuilder {
	return NewTimestampMessageDigestBuilder(s.signature, digestAlgorithm)
}

// GetTimestampMessageImprintDigestBuilderForToken implements
// timestamp.SignatureTimestampSourceOverrides. Port of the
// getTimestampMessageImprintDigestBuilder(TimestampToken) override.
func (s *TimestampSource) GetTimestampMessageImprintDigestBuilderForToken(
	timestampToken *validation.TimestampToken) timestamp.MessageDigestBuilder {
	return NewTimestampMessageDigestBuilderForToken(s.signature, timestampToken).
		SetTimestampAttribute(s.timestampAttributeMap[timestampToken])
}

// GetSignatureTimestampMessageDigest returns message-imprint digest for a SignatureTimestamp.
// Port of the public getSignatureTimestampMessageDigest(DigestAlgorithm, String).
func (s *TimestampSource) GetSignatureTimestampMessageDigest(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string) model.DSSMessageDigest {
	builder := NewTimestampMessageDigestBuilder(s.signature, digestAlgorithm).
		SetCanonicalizationAlgorithm(canonicalizationMethod)
	return builder.SignatureTimestampMessageDigest()
}

// GetTimestampX1MessageDigest returns message-imprint digest for a SigAndRefsTimestamp. Port of
// the public getTimestampX1MessageDigest(DigestAlgorithm, String, boolean); en319132 defines if
// the timestamp shall be created accordingly to ETSI EN 319 132-1 (SigAndRefsTimestampV2).
func (s *TimestampSource) GetTimestampX1MessageDigest(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string, en319132 bool) model.DSSMessageDigest {
	builder := NewTimestampMessageDigestBuilder(s.signature, digestAlgorithm).
		SetCanonicalizationAlgorithm(canonicalizationMethod).
		SetEn319132(en319132)
	return builder.TimestampX1MessageDigest()
}

// GetTimestampX2MessageDigest returns message-imprint digest for a RefsOnlyTimestamp. Port of the
// public getTimestampX2MessageDigest(DigestAlgorithm, String, boolean); en319132 defines if the
// timestamp shall be created accordingly to ETSI EN 319 132-1 (RefsOnlyTimestampV2).
func (s *TimestampSource) GetTimestampX2MessageDigest(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string, en319132 bool) model.DSSMessageDigest {
	builder := NewTimestampMessageDigestBuilder(s.signature, digestAlgorithm).
		SetCanonicalizationAlgorithm(canonicalizationMethod).
		SetEn319132(en319132)
	return builder.TimestampX2MessageDigest()
}

// GetArchiveTimestampData returns message-imprint digest for an ArchiveTimeStamp. Port of the
// public getArchiveTimestampData(DigestAlgorithm, String).
func (s *TimestampSource) GetArchiveTimestampData(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string) model.DSSMessageDigest {
	builder := NewTimestampMessageDigestBuilder(s.signature, digestAlgorithm).
		SetCanonicalizationAlgorithm(canonicalizationMethod)
	return builder.ArchiveTimestampMessageDigest()
}

// IsContentTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not applicable for XAdES. Port of isContentTimestamp(XAdESAttribute).
func (s *TimestampSource) IsContentTimestamp(signedAttribute *Attribute) bool {
	return false
}

// IsAllDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAllDataObjectsTimestamp(XAdESAttribute).
func (s *TimestampSource) IsAllDataObjectsTimestamp(signedAttribute *Attribute) bool {
	return definition.XAdES132ElementAllDataObjectsTimestamp.IsSameTagName(signedAttribute.Name())
}

// IsIndividualDataObjectsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isIndividualDataObjectsTimestamp(XAdESAttribute).
func (s *TimestampSource) IsIndividualDataObjectsTimestamp(signedAttribute *Attribute) bool {
	return definition.XAdES132ElementIndividualDataObjectsTimestamp.IsSameTagName(signedAttribute.Name())
}

// IsSignatureTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignatureTimestamp(XAdESAttribute).
func (s *TimestampSource) IsSignatureTimestamp(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementSignatureTimestamp.IsSameTagName(unsignedAttribute.Name())
}

// IsCompleteCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteCertificateRef(XAdESAttribute).
func (s *TimestampSource) IsCompleteCertificateRef(unsignedAttribute *Attribute) bool {
	localName := unsignedAttribute.Name()
	return definition.XAdES132ElementCompleteCertificateRefs.IsSameTagName(localName) ||
		definition.XAdES141ElementCompleteCertificateRefsV2.IsSameTagName(localName)
}

// IsAttributeCertificateRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeCertificateRef(XAdESAttribute).
func (s *TimestampSource) IsAttributeCertificateRef(unsignedAttribute *Attribute) bool {
	localName := unsignedAttribute.Name()
	return definition.XAdES132ElementAttributeCertificateRefs.IsSameTagName(localName) ||
		definition.XAdES141ElementAttributeCertificateRefsV2.IsSameTagName(localName)
}

// IsCompleteRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCompleteRevocationRef(XAdESAttribute).
func (s *TimestampSource) IsCompleteRevocationRef(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementCompleteRevocationRefs.IsSameTagName(unsignedAttribute.Name())
}

// IsAttributeRevocationRef implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeRevocationRef(XAdESAttribute).
func (s *TimestampSource) IsAttributeRevocationRef(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementAttributeRevocationRefs.IsSameTagName(unsignedAttribute.Name())
}

// IsRefsOnlyTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRefsOnlyTimestamp(XAdESAttribute).
func (s *TimestampSource) IsRefsOnlyTimestamp(unsignedAttribute *Attribute) bool {
	localName := unsignedAttribute.Name()
	return definition.XAdES132ElementRefsOnlyTimestamp.IsSameTagName(localName) ||
		definition.XAdES141ElementRefsOnlyTimestampV2.IsSameTagName(localName)
}

// IsSigAndRefsTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSigAndRefsTimestamp(XAdESAttribute).
func (s *TimestampSource) IsSigAndRefsTimestamp(unsignedAttribute *Attribute) bool {
	localName := unsignedAttribute.Name()
	return definition.XAdES132ElementSigAndRefsTimestamp.IsSameTagName(localName) ||
		definition.XAdES141ElementSigAndRefsTimestampV2.IsSameTagName(localName)
}

// IsCertificateValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCertificateValues(XAdESAttribute).
func (s *TimestampSource) IsCertificateValues(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementCertificateValues.IsSameTagName(unsignedAttribute.Name())
}

// IsRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isRevocationValues(XAdESAttribute).
func (s *TimestampSource) IsRevocationValues(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementRevocationValues.IsSameTagName(unsignedAttribute.Name())
}

// IsAttrAuthoritiesCertValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttrAuthoritiesCertValues(XAdESAttribute).
func (s *TimestampSource) IsAttrAuthoritiesCertValues(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementAttrAuthoritiesCertValues.IsSameTagName(unsignedAttribute.Name())
}

// IsAttributeRevocationValues implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAttributeRevocationValues(XAdESAttribute).
func (s *TimestampSource) IsAttributeRevocationValues(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementAttributeRevocationValues.IsSameTagName(unsignedAttribute.Name())
}

// IsArchiveTimestamp implements timestamp.SignatureTimestampSourceOverrides.
// Port of isArchiveTimestamp(XAdESAttribute).
func (s *TimestampSource) IsArchiveTimestamp(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementArchiveTimestamp.IsSameTagName(unsignedAttribute.Name())
}

// IsTimeStampValidationData implements timestamp.SignatureTimestampSourceOverrides.
// Port of isTimeStampValidationData(XAdESAttribute).
func (s *TimestampSource) IsTimeStampValidationData(unsignedAttribute *Attribute) bool {
	return definition.XAdES141ElementTimestampValidationData.IsSameTagName(unsignedAttribute.Name())
}

// IsAnyValidationData implements timestamp.SignatureTimestampSourceOverrides.
// Port of isAnyValidationData(XAdESAttribute).
func (s *TimestampSource) IsAnyValidationData(unsignedAttribute *Attribute) bool {
	return definition.XAdES141ElementAnyValidationData.IsSameTagName(unsignedAttribute.Name())
}

// IsValidationDataReferences implements timestamp.SignatureTimestampSourceOverrides.
// NOTE: not supported. Port of isValidationDataReferences(XAdESAttribute).
func (s *TimestampSource) IsValidationDataReferences(unsignedAttribute *Attribute) bool {
	return false
}

// IsCounterSignature implements timestamp.SignatureTimestampSourceOverrides.
// Port of isCounterSignature(XAdESAttribute).
func (s *TimestampSource) IsCounterSignature(unsignedAttribute *Attribute) bool {
	return definition.XAdES132ElementCounterSignature.IsSameTagName(unsignedAttribute.Name())
}

// IsSignaturePolicyStore implements timestamp.SignatureTimestampSourceOverrides.
// Port of isSignaturePolicyStore(XAdESAttribute).
func (s *TimestampSource) IsSignaturePolicyStore(unsignedAttribute *Attribute) bool {
	return definition.XAdES141ElementSignaturePolicyStore.IsSameTagName(unsignedAttribute.Name())
}

// IsEvidenceRecord implements timestamp.SignatureTimestampSourceOverrides.
// Port of isEvidenceRecord(XAdESAttribute).
func (s *TimestampSource) IsEvidenceRecord(unsignedAttribute *Attribute) bool {
	return definition.XAdESEvidencerecordNamespaceElementSealingEvidenceRecords.IsSameTagName(unsignedAttribute.Name())
}

// MakeTimestampTokens is XAdES's own plural timestamp-token factory: an xades132:XAdESTimeStampType
// element (SignatureTimeStamp, AllDataObjectsTimeStamp, ...) can contain more than one
// xades132:EncapsulatedTimeStamp, so one signature attribute can produce several TimestampTokens
// at once. Port of the protected makeTimestampTokens(XAdESAttribute, TimestampType, List) override.
//
// See the file header GAP note: not yet reachable from the base's own internal population flow
// (populateTimestampTokens's private makeTimestampTokens always wraps
// s.overrides.MakeTimestampToken, singular) until SignatureTimestampSourceOverrides gains this
// plural method.
func (s *TimestampSource) MakeTimestampTokens(signatureAttribute *Attribute, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) []*validation.TimestampToken {
	encapsulatedTimestamps := signatureAttribute.NodeList(s.xadesPaths.CurrentEncapsulatedTimestamp())
	if len(encapsulatedTimestamps) == 0 {
		// Upstream logs "The timestamp(s) {} cannot be extracted from the signature!".
		return nil
	}

	/*
	 * 6.3 Requirements on XAdES signature's elements, qualifying properties and services n)
	 * Requirement for SignatureTimeStamp. Each SignatureTimeStamp element shall contain only one
	 * electronic time-stamp.
	 */
	if s.IsSignatureTimestamp(signatureAttribute) && len(encapsulatedTimestamps) > 1 {
		// Upstream logs "Only one EncapsulatedTimeStamp is allowed in '{}' element!".
	}

	var result []*validation.TimestampToken
	for ii, encapsulatedTimeStamp := range encapsulatedTimestamps {
		timestampToken := s.createTimestampToken(encapsulatedTimeStamp, timestampType, references, signatureAttribute, ii)
		if timestampToken != nil {
			timestampToken.SetCanonicalizationMethod(signatureAttribute.TimestampCanonicalizationMethod())
			timestampToken.SetTimestampIncludes(signatureAttribute.TimestampIncludedReferences())
			if enumerations.TimestampTypeIndividualDataObjectsTimestamp == timestampType {
				xadesTSTimestampAddReferences(timestampToken,
					s.getIndividualDataContentTimestampReferences(signatureAttribute.TimestampIncludedReferences()))
			}
			s.timestampAttributeMap[timestampToken] = signatureAttribute
			result = append(result, timestampToken)
		}
	}
	return result
}

// createTimestampToken ports the private createTimestampToken(Element, TimestampType, List,
// Attribute, Integer); the surrounding try/catch (any failure logged and answered with nil)
// is reproduced with a defer/recover.
func (s *TimestampSource) createTimestampToken(encapsulatedTimeStamp *xmldom.Node, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference, signatureAttribute *Attribute, orderWithinAttribute int) (result *validation.TimestampToken) {
	defer func() {
		if recover() != nil {
			// Upstream logs "Unable to build timestamp token from binaries '{}'! Reason : {}".
			result = nil
		}
	}()

	base64EncodedTimestamp := encapsulatedTimeStamp.TextContent()
	binaries := utils.FromBase64(base64EncodedTimestamp)
	orderWithin := orderWithinAttribute
	identifierBuilder := timestamp.NewSignatureTimestampIdentifierBuilder(binaries).
		SetSignature(s.signature).
		SetAttribute(signatureAttribute).
		SetOrderOfAttribute(s.GetAttributeOrder(signatureAttribute)).
		SetOrderWithinAttribute(&orderWithin)
	token, err := validation.NewTimestampTokenWithIdentifierBuilder(binaries, timestampType, references, identifierBuilder)
	if err != nil {
		panic(err)
	}
	return token
}

// MakeEvidenceRecords implements timestamp.SignatureTimestampSourceOverrides.
// Port of makeEvidenceRecords(XAdESAttribute, List).
//
// Unlike Java's `(Element) element.getChildNodes().item(ii)`, which unconditionally casts EVERY
// child node (including whitespace text nodes) and would throw an uncaught ClassCastException on
// a non-Element child, xadesTimestampSourceMustElement below reproduces that same unconditional,
// uncaught-outside-this-function panic rather than silently skipping non-element children.
func (s *TimestampSource) MakeEvidenceRecords(signatureAttribute *Attribute,
	references []*validation.TimestampedReference) []validation.EvidenceRecord {
	element := signatureAttribute.Element()
	if element == nil || element.FirstChild == nil {
		// Upstream logs "The element containing evidence record(s) is empty!".
		return nil
	}

	var result []validation.EvidenceRecord
	ii := 0
	for child := element.FirstChild; child != nil; child = child.NextSibling {
		encapsulatedEvidenceRecord := xadesTimestampSourceMustElement(child)
		if evidenceRecord := s.createEvidenceRecord(signatureAttribute, encapsulatedEvidenceRecord, ii); evidenceRecord != nil {
			result = append(result, evidenceRecord)
		}
		ii++
	}
	return result
}

// xadesTimestampSourceMustElement ports Java's unconditional `(Element) node` cast: panics when
// node is not an Element, matching the unchecked ClassCastException upstream would raise.
func xadesTimestampSourceMustElement(node *xmldom.Node) *xmldom.Node {
	if node.Kind != xmldom.Element {
		panic(fmt.Sprintf("ClassCastException: node of kind %v cannot be cast to Element", node.Kind))
	}
	return node
}

// createEvidenceRecord ports the private createEvidenceRecord(XAdESAttribute, Element, int),
// whose surrounding try/catch (any failure logged and answered with nil) is reproduced by a
// defer/recover.
func (s *TimestampSource) createEvidenceRecord(signatureAttribute *Attribute,
	encapsulatedEvidenceRecord *xmldom.Node, orderWithinAttribute int) (result validation.EvidenceRecord) {
	defer func() {
		if recover() != nil {
			// Upstream logs "Unable to build an embedded evidence record. Reason : {}".
			result = nil
		}
	}()

	erDocument := xadesTimestampSourceEvidenceRecordDocument(encapsulatedEvidenceRecord)
	evidenceRecordAnalyzer, err := analyzer.EvidenceRecordAnalyzerFromDocument(erDocument)
	if err != nil {
		return nil
	}
	evidenceRecordAnalyzer.SetEvidenceRecordOrigin(enumerations.EvidenceRecordOriginSignature)

	embeddedEvidenceRecordHelper := NewEmbeddedEvidenceRecordHelper(s.signature, signatureAttribute)
	embeddedEvidenceRecordHelper.SetDetachedContents(s.signature.DetachedContents())
	embeddedEvidenceRecordHelper.SetOrderOfAttribute(s.GetAttributeOrder(signatureAttribute))
	orderWithin := orderWithinAttribute
	embeddedEvidenceRecordHelper.SetOrderWithinAttribute(&orderWithin)
	evidenceRecordAnalyzer.SetEmbeddedEvidenceRecordHelper(embeddedEvidenceRecordHelper)

	return evidenceRecordAnalyzer.EvidenceRecord()
}

// xadesTimestampSourceEvidenceRecordDocument ports the private
// getEvidenceRecordDocument(Element). Panics with the Java UnsupportedOperationException message
// for an unsupported format, caught by createEvidenceRecord's defer/recover.
func xadesTimestampSourceEvidenceRecordDocument(encapsulatedEvidenceRecord *xmldom.Node) model.DSSDocument {
	localName := encapsulatedEvidenceRecord.Name.Local
	switch {
	case definition.XAdESEvidencerecordNamespaceElementEvidenceRecord.IsSameTagName(localName):
		return xmlutils.NewDOMDocument(encapsulatedEvidenceRecord)

	case definition.XAdESEvidencerecordNamespaceElementASN1EvidenceRecord.IsSameTagName(localName):
		base64EncodedEvidenceRecord := encapsulatedEvidenceRecord.TextContent()
		if utils.IsBase64Encoded(base64EncodedEvidenceRecord) {
			return model.NewInMemoryDocument(utils.FromBase64(base64EncodedEvidenceRecord))
		}
	}
	panic(fmt.Sprintf("The provided format of an evidence record within "+
		"the '%s' element is not supported.", localName))
}

// GetTimestampScopes generates timestamp token scopes using this format's TimestampScopeFinder
// (filtering by IndividualDataObjectsTimestamp includes when present), rather than the base's
// plain EncapsulatedTimestampScopeFinder. Port of the protected getTimestampScopes(TimestampToken)
// override.
//
// See the file header GAP note: not yet reachable from the base's own internal validateTimestamps()
// calls until SignatureTimestampSourceOverrides gains this method.
func (s *TimestampSource) GetTimestampScopes(timestampToken *validation.TimestampToken) []mscope.SignatureScope {
	timestampScopeFinder := NewTimestampScopeFinder()
	timestampScopeFinder.SetSignature(s.signature)
	return timestampScopeFinder.FindTimestampScope(timestampToken)
}

// getIndividualDataContentTimestampReferences ports the private
// getIndividualDataContentTimestampReferences(List<TimestampInclude>).
func (s *TimestampSource) getIndividualDataContentTimestampReferences(
	timestampIncludes []*validation.TimestampInclude) []*validation.TimestampedReference {
	var result []mscope.SignatureScope
	signatureScopes := s.signature.SignatureScopes()
	if utils.IsCollectionNotEmpty(signatureScopes) {
		for _, xadesReferenceValidation := range s.signature.XAdESReferenceValidations() {
			if xadesTimestampSourceIsContentTimestampedReference(xadesReferenceValidation, timestampIncludes) {
				for _, signatureScope := range signatureScopes {
					if utils.EndsWithIgnoreCase(xadesReferenceValidation.Uri(), signatureScope.DocumentName()) {
						result = append(result, signatureScope)
					}
				}
			}
		}
	}
	return timestamp.SignerDataTimestampedReferences(result)
}

// xadesTimestampSourceIsContentTimestampedReference ports the private
// isContentTimestampedReference(ReferenceValidation, List<TimestampInclude>); named
// distinctly from xades_timestamp_scope_finder.go's identical-bodied private method of the same
// Java name, since Java declares its own private copy in each class and this port mirrors that
// duplication rather than sharing a helper.
func xadesTimestampSourceIsContentTimestampedReference(xadesReferenceValidation *ReferenceValidation,
	includes []*validation.TimestampInclude) bool {
	if xadesReferenceValidation.Id() != "" {
		for _, timestampInclude := range includes {
			if xadesReferenceValidation.Id() == timestampInclude.URI() {
				return true
			}
		}
	}
	return false
}

// MakeTimestampToken implements timestamp.SignatureTimestampSourceOverrides.
// Port of makeTimestampToken(XAdESAttribute, TimestampType, List): unconditionally unsupported,
// since an xades132:XAdESTimeStampType element can contain more than one timestamp - see
// MakeTimestampTokens (plural) above, which XAdES uses instead.
func (s *TimestampSource) MakeTimestampToken(signatureAttribute *Attribute, timestampType enumerations.TimestampType,
	references []*validation.TimestampedReference) *validation.TimestampToken {
	panic("XAdESTimeStampType element can contain more than one timestamp")
}

// IncorporateArchiveTimestampReferences shadows the promoted embedded-base method so that this
// port's required substitute for Java's virtual dispatch reaches this file's own
// archiveTimestampReferences (which folds in the extra ds:KeyInfo references XAdES's
// getArchiveTimestampReferences override adds) - see the file header GAP note.
func (s *TimestampSource) IncorporateArchiveTimestampReferences(timestampToken *validation.TimestampToken,
	previousTimestamps []*validation.TimestampToken) {
	xadesTSTimestampAddReferences(timestampToken, s.archiveTimestampReferences(previousTimestamps))
}

// archiveTimestampReferences reimplements the base's private getArchiveTimestampReferences(List),
// which this file cannot call directly (unexported, different package), from exported base
// accessors, then applies XAdES's own extra ds:KeyInfo addition. Port of the protected
// getArchiveTimestampReferences(List<TimestampToken>) override, folding in - via
// signatureTimestampReferences below - the identical addition Java's own base method picks up
// through virtual dispatch of getSignatureTimestampReferences().
func (s *TimestampSource) archiveTimestampReferences(previousTimestamps []*validation.TimestampToken) []*validation.TimestampedReference {
	var timestampedReferences []*validation.TimestampedReference
	xadesTSAddReferences(&timestampedReferences, s.signatureTimestampReferences())
	xadesTSAddReferences(&timestampedReferences, s.encapsulatedReferencesFromTimestamps(previousTimestamps))
	xadesTSAddReferences(&timestampedReferences, s.UnsignedPropertiesReferences())
	xadesTSAddReferences(&timestampedReferences, s.GetKeyInfoReferences())
	return timestampedReferences
}

// signatureTimestampReferences reimplements the base's private getSignatureTimestampReferences(),
// which this file cannot call directly, from exported base accessors, then applies XAdES's own
// extra ds:KeyInfo addition. Port of the protected getSignatureTimestampReferences() override.
func (s *TimestampSource) signatureTimestampReferences() []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	xadesTSAddReferences(&references, s.encapsulatedReferencesFromTimestamps(s.ContentTimestamps()))
	xadesTSAddReferences(&references, s.SignerDataReferences())
	xadesTSAddReference(&references, validation.NewTimestampedReference(s.signature.ID(), enumerations.TimestampedObjectTypeSignature))
	xadesTSAddReferences(&references, s.signingCertificateTimestampReferences())
	if SignatureUtilsIsKeyInfoCovered(s.signature) {
		xadesTSAddReferences(&references, s.GetKeyInfoReferences())
	}
	return references
}

// signingCertificateTimestampReferences reimplements the base's private
// getSigningCertificateTimestampReferences().
func (s *TimestampSource) signingCertificateTimestampReferences() []*validation.TimestampedReference {
	signatureCertificateSource := s.signature.CertificateSource()
	return timestamp.CreateReferencesForCertificateRefs(
		signatureCertificateSource.SigningCertificateRefs(), signatureCertificateSource, s.CertificateSource())
}

// encapsulatedReferencesFromTimestamps reimplements the base's private
// getEncapsulatedReferencesFromTimestamps(List), from the exported timestamp.ReferencesFromTimestamp
// and this file's merged-source accessors.
func (s *TimestampSource) encapsulatedReferencesFromTimestamps(timestampTokens []*validation.TimestampToken) []*validation.TimestampedReference {
	var references []*validation.TimestampedReference
	for _, timestampToken := range timestampTokens {
		refs := xadesTSMust(timestamp.ReferencesFromTimestamp(timestampToken, s.CertificateSource(), s.CRLSource(), s.OCSPSource()))
		xadesTSAddReferences(&references, refs)
	}
	return references
}

// GetCertificateRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCertificateRefs(XAdESAttribute).
func (s *TimestampSource) GetCertificateRefs(unsignedAttribute *Attribute) []*spi.CertificateRef {
	var certRefs []*spi.CertificateRef
	certificateRefV1 := xadesTimestampSourceIsCertificateRefV1(unsignedAttribute)

	var certRefsNodeList []*xmldom.Node
	if certificateRefV1 {
		if currentCertRefsCertChildrenPath := s.xadesPaths.CurrentCertRefsCertChildren(); currentCertRefsCertChildrenPath != nil {
			certRefsNodeList = unsignedAttribute.NodeList(currentCertRefsCertChildrenPath)
		}
	} else {
		if currentCertRefs141CertChildrenPath := s.xadesPaths.CurrentCertRefs141CertChildren(); currentCertRefs141CertChildrenPath != nil {
			certRefsNodeList = unsignedAttribute.NodeList(currentCertRefs141CertChildrenPath)
		}
	}

	for _, certRefElement := range certRefsNodeList {
		var certificateRef *spi.CertificateRef
		if certificateRefV1 {
			certificateRef = CertificateRefExtractionUtilsCreateCertificateRefFromV1(certRefElement, s.xadesPaths)
		} else {
			certificateRef = CertificateRefExtractionUtilsCreateCertificateRefFromV2(certRefElement, s.xadesPaths)
		}
		if certificateRef != nil {
			certRefs = append(certRefs, certificateRef)
		}
	}
	return certRefs
}

// xadesTimestampSourceIsCertificateRefV1 ports the private isCertificateRefV1(XAdESAttribute).
func xadesTimestampSourceIsCertificateRefV1(unsignedAttribute *Attribute) bool {
	localName := unsignedAttribute.Name()
	return definition.XAdES132ElementAttributeCertificateRefs.IsSameTagName(localName) ||
		definition.XAdES132ElementCompleteCertificateRefs.IsSameTagName(localName)
}

// GetCRLRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCRLRefs(XAdESAttribute).
func (s *TimestampSource) GetCRLRefs(unsignedAttribute *Attribute) []*spi.CRLRef {
	var crlRefs []*spi.CRLRef
	nodeList := unsignedAttribute.NodeList(s.xadesPaths.CurrentCRLRefsChildren())
	for _, element := range nodeList {
		crlRef := RevocationRefExtractionUtilsCreateCRLRef(s.xadesPaths, element)
		if crlRef != nil {
			crlRefs = append(crlRefs, crlRef)
		}
	}
	return crlRefs
}

// GetOCSPRefs implements timestamp.SignatureTimestampSourceOverrides.
// Port of getOCSPRefs(XAdESAttribute).
func (s *TimestampSource) GetOCSPRefs(unsignedAttribute *Attribute) []*spi.OCSPRef {
	var ocspRefs []*spi.OCSPRef
	nodeList := unsignedAttribute.NodeList(s.xadesPaths.CurrentOCSPRefsChildren())
	for _, element := range nodeList {
		ocspRef := RevocationRefExtractionUtilsCreateOCSPRef(s.xadesPaths, element)
		if ocspRef != nil {
			ocspRefs = append(ocspRefs, ocspRef)
		}
	}
	return ocspRefs
}

// GetEncapsulatedCertificateIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCertificateIdentifiers(XAdESAttribute).
func (s *TimestampSource) GetEncapsulatedCertificateIdentifiers(unsignedAttribute *Attribute) []model.Identifier {
	var certificateIdentifiers []model.Identifier
	xPathQuery := s.xadesPaths.CurrentEncapsulatedCertificate()
	if s.IsTimeStampValidationData(unsignedAttribute) || s.IsAnyValidationData(unsignedAttribute) {
		xPathQuery = s.xadesPaths.CurrentCertificateValuesEncapsulatedCertificate()
	}
	encapsulatedNodes := unsignedAttribute.NodeList(xPathQuery)
	for _, element := range encapsulatedNodes {
		func() {
			defer func() { recover() }()
			binaries := xadesTimestampSourceEncapsulatedTokenBinaries(element)
			certificateToken, err := spi.DSSUtilsLoadCertificateFromBinary(binaries)
			if err != nil {
				// Upstream logs "Unable to parse an encapsulated certificate : {}".
				return
			}
			certificateIdentifiers = append(certificateIdentifiers, certificateToken.DSSID())
		}()
	}
	return certificateIdentifiers
}

// GetEncapsulatedCRLIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedCRLIdentifiers(XAdESAttribute).
func (s *TimestampSource) GetEncapsulatedCRLIdentifiers(unsignedAttribute *Attribute) []*crlparser.CRLBinary {
	var crlIdentifiers []*crlparser.CRLBinary
	xPathQuery := s.xadesPaths.CurrentEncapsulatedCRLValue()
	if s.IsTimeStampValidationData(unsignedAttribute) || s.IsAnyValidationData(unsignedAttribute) {
		xPathQuery = s.xadesPaths.CurrentRevocationValuesEncapsulatedCRLValue()
	}
	encapsulatedNodes := unsignedAttribute.NodeList(xPathQuery)
	for _, element := range encapsulatedNodes {
		func() {
			defer func() { recover() }()
			binaries := xadesTimestampSourceEncapsulatedTokenBinaries(element)
			crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(binaries)
			if err != nil {
				// Upstream logs "Unable to parse CRL binaries : {}".
				return
			}
			crlIdentifiers = append(crlIdentifiers, crlBinary)
		}()
	}
	return crlIdentifiers
}

// GetEncapsulatedOCSPIdentifiers implements timestamp.SignatureTimestampSourceOverrides.
// Port of getEncapsulatedOCSPIdentifiers(XAdESAttribute).
func (s *TimestampSource) GetEncapsulatedOCSPIdentifiers(unsignedAttribute *Attribute) []*spi.OCSPResponseBinary {
	var ocspIdentifiers []*spi.OCSPResponseBinary
	xPathQuery := s.xadesPaths.CurrentEncapsulatedOCSPValue()
	if s.IsTimeStampValidationData(unsignedAttribute) || s.IsAnyValidationData(unsignedAttribute) {
		xPathQuery = s.xadesPaths.CurrentRevocationValuesEncapsulatedOCSPValue()
	}
	encapsulatedNodes := unsignedAttribute.NodeList(xPathQuery)
	for _, element := range encapsulatedNodes {
		func() {
			defer func() { recover() }()
			binaries := xadesTimestampSourceEncapsulatedTokenBinaries(element)
			basicOCSPResp, err := spi.DSSRevocationUtilsLoadOCSPFromBinaries(binaries)
			if err != nil {
				// Upstream logs "Unable to parse OCSP response binaries : {}".
				return
			}
			binary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
			if err != nil {
				return
			}
			ocspIdentifiers = append(ocspIdentifiers, binary)
		}()
	}
	return ocspIdentifiers
}

// xadesTimestampSourceEncapsulatedTokenBinaries returns the encapsulated byte array from the
// given encapsulatedElement. Port of the private getEncapsulatedTokenBinaries(Element). Panics
// with the Java DSSException message when the element does not contain an encapsulated base64
// text-node value.
func xadesTimestampSourceEncapsulatedTokenBinaries(encapsulatedElement *xmldom.Node) []byte {
	if encapsulatedElement.FirstChild != nil {
		firstChild := encapsulatedElement.FirstChild
		if firstChild.Kind == xmldom.Text {
			base64String := firstChild.TextContent()
			if utils.IsBase64Encoded(base64String) {
				return utils.FromBase64(base64String)
			}
		}
	}
	panic(model.NewDSSError(fmt.Sprintf("Cannot create the token reference. "+
		"The element with local name [%s] must contain an encapsulated base64 token value! "+
		"The found value is not a text node!", encapsulatedElement.Name.Local)))
}

// GetArchiveTimestampType implements timestamp.SignatureTimestampSourceOverrides.
// Port of getArchiveTimestampType(XAdESAttribute).
func (s *TimestampSource) GetArchiveTimestampType(unsignedAttribute *Attribute) enumerations.ArchiveTimestampType {
	if definition.XAdESNamespaceXAdES141.IsSameUri(unsignedAttribute.Namespace()) {
		return enumerations.ArchiveTimestampTypeXAdES141
	}
	return enumerations.ArchiveTimestampTypeXAdES
}

// GetCounterSignatures implements timestamp.SignatureTimestampSourceOverrides.
// Port of getCounterSignatures(XAdESAttribute).
func (s *TimestampSource) GetCounterSignatures(unsignedAttribute *Attribute) []validation.AdvancedSignature {
	counterSignatureNode, err := xmlutils.XPathUtilsGetNode(unsignedAttribute.Element(), common.XMLDSigPathSignaturePath)
	if err == nil && counterSignatureNode != nil {
		counterSignatures := s.signature.CounterSignatures()
		for _, counterSignature := range counterSignatures {
			// Java casts unconditionally: `(XAdESSignature) counterSignature`. Every concrete
			// AdvancedSignature an XAdES source's CounterSignatures() hands out is a
			// *Signature, so this assertion always succeeds in practice; a genuine mismatch
			// panics, exactly as Java's ClassCastException would.
			xadesCounterSignature := counterSignature.(*Signature)
			if counterSignatureNode == xadesCounterSignature.SignatureElement() {
				// NOTE: only one counter signature is allowed within the CounterSignature
				// qualifying property.
				return []validation.AdvancedSignature{counterSignature}
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Small reference-list helpers, mirroring the frozen spi/validation/timestamp package's own
// unexported addReference/addReferences/timestampAddReferences (abstract_timestamp_source.go),
// which this file cannot call directly (different package). Prefixed distinctly (xadesTS) to
// avoid colliding with any sibling chunk writing to this same package concurrently, mirroring
// cades_timestamp_source.go's identical cadesTS-prefixed set.
// -----------------------------------------------------------------------------

// xadesTSAddReference adds referenceToAdd to *referenceList without duplicates.
// Port of the AbstractTimestampSource addReference(List, TimestampedReference) helper.
func xadesTSAddReference(referenceList *[]*validation.TimestampedReference, referenceToAdd *validation.TimestampedReference) {
	xadesTSAddReferences(referenceList, []*validation.TimestampedReference{referenceToAdd})
}

// xadesTSAddReferences adds referencesToAdd to *referenceList without duplicates (by
// TimestampedReference.Equals). Port of the AbstractTimestampSource addReferences(List, List)
// helper.
func xadesTSAddReferences(referenceList *[]*validation.TimestampedReference, referencesToAdd []*validation.TimestampedReference) {
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

// xadesTSTimestampAddReferences enriches timestampToken's TimestampedReferences with
// referencesToAdd, without duplicates, via TimestampToken.SetTimestampedReferences - the Go
// counterpart of Java's `addReferences(timestampToken.getTimestampedReferences(), ...)` List
// reference-semantics mutation. Mirrors the frozen timestamp package's own unexported
// timestampAddReferences.
func xadesTSTimestampAddReferences(timestampToken *validation.TimestampToken, referencesToAdd []*validation.TimestampedReference) {
	merged := append([]*validation.TimestampedReference(nil), timestampToken.TimestampedReferences()...)
	xadesTSAddReferences(&merged, referencesToAdd)
	timestampToken.SetTimestampedReferences(merged)
}

// xadesTSMust panics with a model.DSSError wrapping err when err is non-nil, otherwise returning
// value unchanged. Mirrors spi/validation/timestamp/abstract_timestamp_source.go's own unexported
// `must` helper (this file cannot call it directly, different package).
func xadesTSMust[T any](value T, err error) T {
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return value
}

// compile-time assertion: *TimestampSource implements
// timestamp.SignatureTimestampSourceOverrides[*XAdESSignature, *XAdESAttribute], matching Java's
// "extends SignatureTimestampSource<Signature, Attribute>".
var _ timestamp.SignatureTimestampSourceOverrides[*Signature, *Attribute] = (*TimestampSource)(nil)
