// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/timestamp/XAdESTimestampMessageDigestBuilder.java (DSS 6.5.RC1).
//
// # Santuario replacement
//
// org.apache.xml.security.signature.{Reference,XMLSignatureInput} calls map to internal/xmldsig
// per its doc.go table (frozen; see S4D_BRIEF.md's "Santuario-replacement rules"): Reference is
// *xmldsig.Reference (already used throughout this package, e.g. dss_reference.go,
// reference_processor.go), and XMLSignatureInput's getOctetStream() is *xmldsig.Data.Bytes().
// XMLCanonicalizer canonicalization goes through xml/utils.XMLCanonicalizer (frozen, byte-parity
// with the Java Transformer per S4D_BRIEF.md).
//
// # Exception-scope deviation
//
// Java's writeReferenceBytes/writeDigestValueOnCanonicalizedBinaries declare/throw
// XMLSecurityException|IOException, caught narrowly by getAllDataObjectsTimestampMessageDigest/
// getIndividualDataObjectsTimestampMessageDigest's `catch (XMLSecurityException | IOException e)`;
// writeCanonicalizedValue/writeDigestValueOnCanonicalizedNode instead throw an unchecked
// DSSException, caught only by the broader `catch (Exception e)` of
// getSignatureTimestampMessageDigest/getTimestampX1MessageDigest/getTimestampX2MessageDigest/
// getArchiveTimestampMessageDigest. This port keeps that same split (writeReferenceBytes/
// writeDigestValueOnCanonicalizedBinaries return a plain error consumed directly by
// allDataObjectsTimestampMessageDigest/individualDataObjectsTimestampMessageDigest;
// writeCanonicalizedValue/writeDigestValueOnCanonicalizedNode panic, caught by a defer/recover at
// the top of the four broader methods - which also wrap any writeReferenceBytes error they call
// into the same panic, matching getArchiveTimestampMessageDigest's own broader catch of that
// method too) rather than threading Go's single untyped error through every call, since Go has no
// typed catch to mirror precisely otherwise.
//
// # DSSMessageDigest null vs "empty digest" collapse
//
// getTimestampX2MessageDigest/getArchiveTimestampMessageDigest return Java `null` on failure,
// while getSignatureTimestampMessageDigest/getTimestampX1MessageDigest instead return
// DSSMessageDigest.createEmptyDigest(); this port's model.DSSMessageDigest is a value type whose
// zero value IS model.CreateEmptyDSSMessageDigest()'s result (see model/dss_message_digest.go),
// so both collapse to the identical Go value - no special-casing needed.
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/timestamp"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESTimestampMessageDigestBuilder builds a message-imprint for XAdES timestamps. Port of the
// class XAdESTimestampMessageDigestBuilder, implementing timestamp.TimestampMessageDigestBuilder.
type XAdESTimestampMessageDigestBuilder struct {
	// references is the list of XAdES signature references.
	references []*xmldsig.Reference

	// signature is the signature element.
	signature *xmldom.Node

	// xadesPaths is the XAdES XPaths to use.
	xadesPaths definition.XAdESPath

	// digestAlgorithm is the digest algorithm to be used for message-imprint digest computation.
	digestAlgorithm enumerations.DigestAlgorithm

	// timestampToken is the timestamp token to compute message-digest for.
	timestampToken *validation.TimestampToken

	// canonicalizationAlgorithm is the canonicalization algorithm to be used for message-imprint
	// computation.
	canonicalizationAlgorithm string

	// en319132 identifies whether message-imprint shall be built as per EN 319 132-1 standard
	// (against the old version).
	en319132 bool

	// timestampAttribute is a signature attribute corresponding to the time-stamp.
	timestampAttribute *XAdESAttribute
}

// NewXAdESTimestampMessageDigestBuilder is the default constructor to be used for a new
// timestamp creation. This constructor requires certain properties to be provided for
// message-digest computation (see available setters). Port of the (XAdESSignature,
// DigestAlgorithm) constructor.
//
// Panics with the Java message when digestAlgorithm is empty (Objects.requireNonNull).
func NewXAdESTimestampMessageDigestBuilder(signature *XAdESSignature,
	digestAlgorithm enumerations.DigestAlgorithm) *XAdESTimestampMessageDigestBuilder {
	b := newXAdESTimestampMessageDigestBuilderBase(signature)
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	b.digestAlgorithm = digestAlgorithm
	return b
}

// NewXAdESTimestampMessageDigestBuilderForToken is the constructor to be used for existing
// timestamp message-imprint computation. Port of the (XAdESSignature, TimestampToken)
// constructor.
//
// Panics with the Java message when timestampToken is nil (Objects.requireNonNull).
func NewXAdESTimestampMessageDigestBuilderForToken(signature *XAdESSignature,
	timestampToken *validation.TimestampToken) *XAdESTimestampMessageDigestBuilder {
	b := newXAdESTimestampMessageDigestBuilderBase(signature)
	if timestampToken == nil {
		panic("TimestampToken cannot be null!")
	}
	b.timestampToken = timestampToken
	b.digestAlgorithm = timestampToken.DigestAlgorithm()
	b.canonicalizationAlgorithm = timestampToken.CanonicalizationMethod()
	return b
}

// newXAdESTimestampMessageDigestBuilderBase is the internal constructor to instantiate required
// values from a signature object. Port of the private (XAdESSignature) constructor.
//
// Panics with the Java message when signature is nil (Objects.requireNonNull).
func newXAdESTimestampMessageDigestBuilderBase(signature *XAdESSignature) *XAdESTimestampMessageDigestBuilder {
	if signature == nil {
		panic("Signature cannot be null!")
	}
	return &XAdESTimestampMessageDigestBuilder{
		signature:  signature.SignatureElement(),
		references: signature.References(),
		xadesPaths: signature.XAdESPaths(),
	}
}

// SetCanonicalizationAlgorithm sets the canonicalization algorithm to be used for message-digest
// computation. Port of setCanonicalizationAlgorithm(String).
func (b *XAdESTimestampMessageDigestBuilder) SetCanonicalizationAlgorithm(canonicalizationAlgorithm string) *XAdESTimestampMessageDigestBuilder {
	b.canonicalizationAlgorithm = canonicalizationAlgorithm
	return b
}

// SetEn319132 sets whether the message-digest should be computed for an EN 319 132-1 standard
// timestamp token. Port of setEn319132(boolean).
func (b *XAdESTimestampMessageDigestBuilder) SetEn319132(en319132 bool) *XAdESTimestampMessageDigestBuilder {
	b.en319132 = en319132
	return b
}

// SetTimestampAttribute sets a signature attribute corresponding to the time-stamp token.
// Defines also en319132 based on the provided timestamp attribute. Port of
// setTimestampAttribute(XAdESAttribute).
func (b *XAdESTimestampMessageDigestBuilder) SetTimestampAttribute(timestampAttribute *XAdESAttribute) *XAdESTimestampMessageDigestBuilder {
	b.timestampAttribute = timestampAttribute
	if timestampAttribute != nil {
		b.en319132 = xadesTimestampMessageDigestBuilderIsEn319132TimestampToken(timestampAttribute)
	}
	return b
}

// ContentTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getContentTimestampMessageDigest().
//
// Panics with the Java message when checkSignatureIntegrity (i.e. References()) has not been
// invoked first (IllegalStateException), or with an UnsupportedOperationException-equivalent
// message for a content timestamp type this method does not support.
func (b *XAdESTimestampMessageDigestBuilder) ContentTimestampMessageDigest() model.DSSMessageDigest {
	// all data timestamp is considered by default
	timeStampType := enumerations.TimestampType_ALL_DATA_OBJECTS_TIMESTAMP
	if b.timestampToken != nil {
		timeStampType = b.timestampToken.TimeStampType()
	}
	if len(b.references) == 0 {
		panic("The method 'checkSignatureIntegrity' must be invoked first!")
	}

	switch timeStampType {
	case enumerations.TimestampType_ALL_DATA_OBJECTS_TIMESTAMP:
		return b.allDataObjectsTimestampMessageDigest()
	case enumerations.TimestampType_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP:
		return b.individualDataObjectsTimestampMessageDigest()
	default:
		panic(fmt.Sprintf("The content timestamp of type '%s' is not supported!", timeStampType))
	}
}

// allDataObjectsTimestampMessageDigest returns the computed message-imprint digest for
// xades132:AllDataObjectsTimestamp token. Port of the protected
// getAllDataObjectsTimestampMessageDigest().
func (b *XAdESTimestampMessageDigestBuilder) allDataObjectsTimestampMessageDigest() model.DSSMessageDigest {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err == nil {
		for _, reference := range b.references {
			if !DSSXMLUtilsIsSignedProperties(reference, b.xadesPaths) {
				if werr := b.writeReferenceBytes(digestCalculator, reference, b.canonicalizationAlgorithm); werr != nil {
					err = werr
					break
				}
			}
		}
	}
	if err != nil {
		// Upstream logs "Unable to extract AllDataObjectsTimestampData. Reason : {}".
		return model.CreateEmptyDSSMessageDigest()
	}
	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// individualDataObjectsTimestampMessageDigest returns the computed message-imprint digest for
// xades132:IndividualDataObjectsTimestamp token. Port of the protected
// getIndividualDataObjectsTimestampMessageDigest().
//
// Panics with the Java message when the Included referencedData attribute is either not present
// or set to false (IllegalArgumentException).
func (b *XAdESTimestampMessageDigestBuilder) individualDataObjectsTimestampMessageDigest() model.DSSMessageDigest {
	if !xadesTimestampMessageDigestBuilderCheckTimestampTokenIncludes(b.timestampToken) {
		panic("The Included referencedData attribute is either not present or set to false!")
	}

	includes := b.timestampToken.TimestampIncludes()
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.digestAlgorithm)
	if err == nil {
		for _, include := range includes {
			reference := xadesTimestampMessageDigestBuilderCorrespondingReference(include, b.references)
			if reference != nil {
				if werr := b.writeReferenceBytes(digestCalculator, reference, b.canonicalizationAlgorithm); werr != nil {
					err = werr
					break
				}
			}
			// else Upstream logs "No ds:Reference found corresponding to an
			// IndividualDataObjectsTimestamp include with URI '{}'!".
		}
	}
	if err != nil {
		// Upstream logs "Unable to extract IndividualDataObjectsTimestampData. Reason : {}".
		return model.CreateEmptyDSSMessageDigest()
	}
	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// writeReferenceBytes processes the retrieved ds:Reference element according to the
// reference-processing model of XMLDSIG [1] clause 4.4.3.2, and if the result is an XML node
// set, canonicalizes it as specified in clause 4.5. Port of the private writeReferenceBytes
// (DSSMessageDigestCalculator, Reference, String).
func (b *XAdESTimestampMessageDigestBuilder) writeReferenceBytes(digestCalculator *spi.DSSMessageDigestCalculator,
	reference *xmldsig.Reference, canonicalizationMethod string) error {
	if ReferenceOutputType_NODE_SET == DSSXMLUtilsGetReferenceOutputType(reference) {
		referencedBytes, err := reference.ReferencedBytes()
		if err != nil {
			return err
		}
		if xmlutils.DomUtilsIsDOMBytes(referencedBytes) {
			return xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedBinaries(digestCalculator, referencedBytes, canonicalizationMethod)
		}
		digestCalculator.Update(referencedBytes)
		return nil
	}

	input, err := reference.ContentsAfterTransformation()
	if err != nil {
		return err
	}
	octetStream, err := input.Bytes()
	if err != nil {
		return err
	}
	digestCalculator.Update(octetStream)
	return nil
}

// xadesTimestampMessageDigestBuilderCheckTimestampTokenIncludes ensures that all Include elements
// referring to the Reference elements have a referencedData attribute set to "true". Port of the
// private checkTimestampTokenIncludes(TimestampToken).
func xadesTimestampMessageDigestBuilderCheckTimestampTokenIncludes(timestampToken *validation.TimestampToken) bool {
	timestampIncludes := timestampToken.TimestampIncludes()
	for _, timestampInclude := range timestampIncludes {
		if !timestampInclude.IsReferencedData() {
			return false
		}
	}
	return true
}

// xadesTimestampMessageDigestBuilderCorrespondingReference ports the private
// getCorrespondingReference(TimestampInclude, List<Reference>).
func xadesTimestampMessageDigestBuilderCorrespondingReference(timestampInclude *validation.TimestampInclude,
	references []*xmldsig.Reference) *xmldsig.Reference {
	uri := timestampInclude.URI()
	for _, reference := range references {
		if uri == reference.ID() {
			return reference
		}
	}
	return nil
}

// SignatureTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getSignatureTimestampMessageDigest().
func (b *XAdESTimestampMessageDigestBuilder) SignatureTimestampMessageDigest() (result model.DSSMessageDigest) {
	defer func() {
		if recover() != nil {
			// Upstream logs MESSAGE_IMPRINT_ERROR/MESSAGE_IMPRINT_ERROR_WITH_ID.
			result = model.CreateEmptyDSSMessageDigest()
		}
	}()

	digestCalculator := xadesTimestampMessageDigestBuilderMustNewCalculator(b.digestAlgorithm)
	b.writeCanonicalizedValue(digestCalculator, common.XMLDSigPath_SIGNATURE_VALUE_PATH, b.canonicalizationAlgorithm)
	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// TimestampX1MessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getTimestampX1MessageDigest().
//
// A.1.5.1 The SigAndRefsTimeStampV2 qualifying property (A.1.5.1.2 Not distributed case)
//
// The input to the electronic time-stamp's message imprint computation input shall be the result
// of taking in order each of the XAdES components listed below, canonicalizing each one as
// specified in clause 4.5, and concatenating the resulting octet streams.
func (b *XAdESTimestampMessageDigestBuilder) TimestampX1MessageDigest() (result model.DSSMessageDigest) {
	defer func() {
		if recover() != nil {
			// Upstream logs MESSAGE_IMPRINT_ERROR/MESSAGE_IMPRINT_ERROR_WITH_ID.
			result = model.CreateEmptyDSSMessageDigest()
		}
	}()

	digestCalculator := xadesTimestampMessageDigestBuilderMustNewCalculator(b.digestAlgorithm)

	// 1) The ds:SignatureValue element.
	b.writeCanonicalizedValue(digestCalculator, common.XMLDSigPath_SIGNATURE_VALUE_PATH, b.canonicalizationAlgorithm)

	// 2) Those among the following unsigned qualifying properties that appear before
	// SigAndRefsTimeStampV2, in their order of appearance within the UnsignedSignatureProperties
	// element. Canonicalization copy is used in order to allow XL/A levels creation.
	unsignedProperties := b.unsignedSignaturePropertiesCanonicalizationCopy()
	if unsignedProperties == nil {
		panic("UnsignedSignatureProperties are not initialized!")
	}

	xadesUnsignedSigProperties := NewXAdESUnsignedSigProperties(unsignedProperties, b.xadesPaths)
	for _, xadesAttribute := range xadesUnsignedSigProperties.Attributes() {
		if b.timestampAttribute != nil && b.timestampAttribute.Equals(xadesAttribute) {
			break
		}

		if b.en319132 {
			/*
			 * - The SignatureTimeStamp qualifying properties.
			 * - The CompleteCertificateRefsV2 qualifying property.
			 * - The CompleteRevocationRefs qualifying property.
			 * - The AttributeCertificateRefsV2 qualifying property if it is present. And
			 * - The AttributeRevocationRefs qualifying property if it is present.
			 */
			if xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(xadesAttribute,
				definition.XAdES132Element_SIGNATURE_TIMESTAMP, definition.XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2,
				definition.XAdES132Element_COMPLETE_REVOCATION_REFS, definition.XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2,
				definition.XAdES132Element_ATTRIBUTE_REVOCATION_REFS) {
				b.writeCanonicalizedAttribute(digestCalculator, xadesAttribute, b.canonicalizationAlgorithm)
			}

		} else {
			/*
			 * TS 101 903 v1.4.2 : 7.5.1 The SigAndRefsTimeStamp element (7.5.1.1 Not distributed
			 * case)
			 *
			 * 2) Those among the following unsigned properties that appear before
			 * SigAndRefsTimeStamp, in their order of appearance within the
			 * UnsignedSignatureProperties element:
			 * - SignatureTimeStamp elements.
			 * - The CompleteCertificateRefs element.
			 * - The CompleteRevocationRefs element.
			 * - The AttributeCertificateRefs element if this property is present.
			 * - The AttributeRevocationRefs element if this property is present.
			 */
			if xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(xadesAttribute,
				definition.XAdES132Element_SIGNATURE_TIMESTAMP, definition.XAdES132Element_COMPLETE_CERTIFICATE_REFS,
				definition.XAdES132Element_COMPLETE_REVOCATION_REFS, definition.XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS,
				definition.XAdES132Element_ATTRIBUTE_REVOCATION_REFS) {
				b.writeCanonicalizedAttribute(digestCalculator, xadesAttribute, b.canonicalizationAlgorithm)
			}
		}
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// TimestampX2MessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getTimestampX2MessageDigest().
//
// A.1.5.2 The RefsOnlyTimeStampV2 qualifying property (A.1.5.2.2 Not distributed case)
//
// The electronic time-stamp's message imprint computation input shall be the result of taking
// those of the qualifying unsigned properties listed below that appear before the
// RefsOnlyTimeStampV2 in their order of appearance within the UnsignedSignatureProperties
// element, canonicalizing each one as specified in clause 4.5, and concatenating the resulting
// octet streams.
func (b *XAdESTimestampMessageDigestBuilder) TimestampX2MessageDigest() (result model.DSSMessageDigest) {
	defer func() {
		if recover() != nil {
			// Upstream logs MESSAGE_IMPRINT_ERROR/MESSAGE_IMPRINT_ERROR_WITH_ID.
			result = model.CreateEmptyDSSMessageDigest()
		}
	}()

	digestCalculator := xadesTimestampMessageDigestBuilderMustNewCalculator(b.digestAlgorithm)

	// Canonicalization copy is used in order to allow XL/A level creation.
	unsignedProperties := b.unsignedSignaturePropertiesCanonicalizationCopy()
	if unsignedProperties == nil {
		panic("UnsignedSignatureProperties are not initialized!")
	}

	xadesUnsignedSigProperties := NewXAdESUnsignedSigProperties(unsignedProperties, b.xadesPaths)
	for _, xadesAttribute := range xadesUnsignedSigProperties.Attributes() {
		if b.timestampAttribute != nil && b.timestampAttribute.Equals(xadesAttribute) {
			break
		}

		// Use RefsOnlyTimeStampV2 on signature creation/extension.
		if b.en319132 {
			/*
			 * - The CompleteCertificateRefsV2 qualifying property.
			 * - The CompleteRevocationRefs qualifying property.
			 * - The AttributeCertificateRefsV2 qualifying property if it is present. And
			 * - The AttributeRevocationRefs qualifying property if it is present.
			 */
			if xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(xadesAttribute,
				definition.XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2, definition.XAdES132Element_COMPLETE_REVOCATION_REFS,
				definition.XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2, definition.XAdES132Element_ATTRIBUTE_REVOCATION_REFS) {
				b.writeCanonicalizedAttribute(digestCalculator, xadesAttribute, b.canonicalizationAlgorithm)
			}

		} else {
			/*
			 * TS 101 903 v1.4.2 : 7.5.1 The SigAndRefsTimeStamp element (7.5.1.1 Not distributed
			 * case)
			 *
			 * - The CompleteCertificateRefs element.
			 * - The CompleteRevocationRefs element.
			 * - The AttributeCertificateRefs element if this property is present.
			 * - The AttributeRevocationRefs element if this property is present.
			 */
			if xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(xadesAttribute,
				definition.XAdES132Element_COMPLETE_CERTIFICATE_REFS, definition.XAdES132Element_COMPLETE_REVOCATION_REFS,
				definition.XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS, definition.XAdES132Element_ATTRIBUTE_REVOCATION_REFS) {
				b.writeCanonicalizedAttribute(digestCalculator, xadesAttribute, b.canonicalizationAlgorithm)
			}
		}
	}

	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// ArchiveTimestampMessageDigest implements timestamp.TimestampMessageDigestBuilder.
// Port of getArchiveTimestampMessageDigest().
//
// 8.2.1 Not distributed case
//
// When xadesv141:ArchiveTimeStamp and all the unsigned properties covered by its time-stamp
// certificateToken have the same parent, this property uses the Implicit mechanism for all the
// time-stamped data objects. The input to the computation of the digest value MUST be built as
// follows:
//
//  1. Initialize the final octet stream as an empty octet stream.
func (b *XAdESTimestampMessageDigestBuilder) ArchiveTimestampMessageDigest() (result model.DSSMessageDigest) {
	defer func() {
		if recover() != nil {
			// Upstream logs MESSAGE_IMPRINT_ERROR/MESSAGE_IMPRINT_ERROR_WITH_ID.
			result = model.CreateEmptyDSSMessageDigest()
		}
	}()

	digestCalculator := xadesTimestampMessageDigestBuilderMustNewCalculator(b.digestAlgorithm)

	/*
	 * 2) Take all the ds:Reference elements in their order of appearance within ds:SignedInfo
	 * referencing whatever the signer wants to sign including the SignedProperties element.
	 * Process each one as indicated below:
	 * - Process the retrieved ds:Reference element according to the reference processing model
	 *   of XMLDSIG.
	 * - If the result is a XML node set, canonicalize it. If ds:Canonicalization is present, the
	 *   algorithm indicated by this element is used. If not, the standard canonicalization method
	 *   specified by XMLDSIG is used.
	 * - Concatenate the resulting octets to the final octet stream.
	 *
	 * The references are already calculated (see CheckSignatureIntegrity()).
	 */
	referenceURIs := make(map[string]struct{}, len(b.references))
	for _, reference := range b.references {
		referenceURIs[xmlutils.DomUtilsGetId(reference.URI())] = struct{}{}
		if err := b.writeReferenceBytes(digestCalculator, reference, b.canonicalizationAlgorithm); err != nil {
			panic(model.NewDSSErrorWithCause(err))
		}
	}

	/*
	 * 3) Take the following XMLDSIG elements in the order they are listed below, canonicalize
	 * each one and concatenate each resulting octet stream to the final octet stream:
	 * - The ds:SignedInfo element.
	 * - The ds:SignatureValue element.
	 * - The ds:KeyInfo element, if present.
	 */
	b.writeCanonicalizedValue(digestCalculator, common.XMLDSigPath_SIGNED_INFO_PATH, b.canonicalizationAlgorithm)
	b.writeCanonicalizedValue(digestCalculator, common.XMLDSigPath_SIGNATURE_VALUE_PATH, b.canonicalizationAlgorithm)
	b.writeCanonicalizedValue(digestCalculator, common.XMLDSigPath_KEY_INFO_PATH, b.canonicalizationAlgorithm)

	/*
	 * 4) Take the unsigned signature properties that appear before the current
	 * xadesv141:ArchiveTimeStamp in the order they appear within the
	 * xades:UnsignedSignatureProperties, canonicalize each one and concatenate each resulting
	 * octet stream to the final octet stream.
	 */
	b.writeTimestampedUnsignedProperties(digestCalculator, b.timestampToken, b.canonicalizationAlgorithm)

	/*
	 * 5) Take all the ds:Object elements except the one containing xades:QualifyingProperties
	 * element. Canonicalize each one and concatenate each resulting octet stream to the final
	 * octet stream. If ds:Canonicalization is present, the algorithm indicated by this element is
	 * used. If not, the standard canonicalization method specified by XMLDSIG is used.
	 */
	objects := b.objects()
	b.writeObjectBytes(digestCalculator, objects, referenceURIs, b.canonicalizationAlgorithm)

	return digestCalculator.MessageDigest(b.digestAlgorithm)
}

// writeCanonicalizedValue ports the private writeCanonicalizedValue(DSSMessageDigestCalculator,
// XPathQuery, String) overload. Panics (via
// xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode) on a canonicalization
// failure, caught by the defer/recover of the four broader callers above.
func (b *XAdESTimestampMessageDigestBuilder) writeCanonicalizedValue(digestCalculator *spi.DSSMessageDigestCalculator,
	xPathQuery common.XPathQuery, canonicalizationMethod string) {
	element, err := xmlutils.XPathUtilsGetElement(b.signature, xPathQuery)
	if err != nil || element == nil {
		return
	}
	xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode(digestCalculator, element, canonicalizationMethod)
}

// writeCanonicalizedAttribute ports the private writeCanonicalizedValue(DSSMessageDigestCalculator,
// XAdESAttribute, String) overload; named distinctly from writeCanonicalizedValue since Go has no
// overloading.
func (b *XAdESTimestampMessageDigestBuilder) writeCanonicalizedAttribute(digestCalculator *spi.DSSMessageDigestCalculator,
	attribute *XAdESAttribute, canonicalizationMethod string) {
	xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode(digestCalculator, attribute.Element(), canonicalizationMethod)
}

// xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedBinaries ports the private
// writeDigestValueOnCanonicalizedBinaries(DSSMessageDigestCalculator, byte[], String), returning
// an error (rather than panicking) to match the narrower XMLSecurityException|IOException catch
// scope of its two callers - see the file header.
func xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedBinaries(digestCalculator *spi.DSSMessageDigestCalculator,
	binaries []byte, canonicalizationMethod string) error {
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationMethod)
	if err != nil {
		return model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", err)
	}
	os := digestCalculator.Writer()
	if cErr := canonicalizer.CanonicalizeBytesTo(binaries, os); cErr != nil {
		_ = os.Close()
		return model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", cErr)
	}
	if cErr := os.Close(); cErr != nil {
		return model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", cErr)
	}
	return nil
}

// xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode ports the private
// writeDigestValueOnCanonicalizedNode(DSSMessageDigestCalculator, Node, String). Panics with a
// model.DSSError on a canonicalization failure, matching Java's `throw new DSSException(...)`.
func xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode(digestCalculator *spi.DSSMessageDigestCalculator,
	node *xmldom.Node, canonicalizationMethod string) {
	canonicalizer, err := xmlutils.XMLCanonicalizerCreateInstanceWithMethod(canonicalizationMethod)
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", err))
	}
	os := digestCalculator.Writer()
	if cErr := canonicalizer.CanonicalizeNodeTo(node, os); cErr != nil {
		_ = os.Close()
		panic(model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", cErr))
	}
	if cErr := os.Close(); cErr != nil {
		panic(model.NewDSSErrorMessageCause("Cannot build an AllDataObjectsTimestamp : An error occurred on canonicalization", cErr))
	}
}

// unsignedSignaturePropertiesDom ports the private getUnsignedSignaturePropertiesDom().
func (b *XAdESTimestampMessageDigestBuilder) unsignedSignaturePropertiesDom() *xmldom.Node {
	element, err := xmlutils.XPathUtilsGetElement(b.signature, b.xadesPaths.UnsignedSignaturePropertiesPath())
	if err != nil {
		return nil
	}
	return element
}

// unsignedSignaturePropertiesCanonicalizationCopy ports the private
// getUnsignedSignaturePropertiesCanonicalizationCopy(). Panics on failure, folding into the same
// broader-catch defer/recover as writeCanonicalizedValue.
func (b *XAdESTimestampMessageDigestBuilder) unsignedSignaturePropertiesCanonicalizationCopy() *xmldom.Node {
	signatureID := DSSXMLUtilsGetIDIdentifier(b.signature)
	node, err := DSSXMLUtilsEnsureNamespacesDefinedWithQuery(b.signature.OwnerDocument(), signatureID, b.xadesPaths.UnsignedSignaturePropertiesPath())
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return node
}

// writeTimestampedUnsignedProperties ports the private writeTimestampedUnsignedProperties
// (DSSMessageDigestCalculator, TimestampToken, String).
//
// In the SD-DSS implementation, when validating the signature, the framework does not add
// missing data (CertificateValues/RevocationValues/AttrAuthoritiesCertValues/
// AttributeRevocationValues MUST-be-added rules from the spec); to do so the signature must be
// extended - matching upstream's commented-out branches, reproduced here only as comments.
func (b *XAdESTimestampMessageDigestBuilder) writeTimestampedUnsignedProperties(digestCalculator *spi.DSSMessageDigestCalculator,
	timestampToken *validation.TimestampToken, canonicalizationMethod string) {
	xadesUnsignedSigProperties := b.xadesUnsignedSignatureProperties(timestampToken)
	for _, xadesAttribute := range xadesUnsignedSigProperties.Attributes() {
		if b.timestampAttribute != nil && b.timestampAttribute.Equals(xadesAttribute) {
			break
		}
		b.writeCanonicalizedAttribute(digestCalculator, xadesAttribute, canonicalizationMethod)
	}
}

// xadesUnsignedSignatureProperties ports the private getXAdESUnsignedSignatureProperties
// (TimestampToken). Panics with the Java message when UnsignedSignatureProperties are not
// initialized (IllegalStateException).
func (b *XAdESTimestampMessageDigestBuilder) xadesUnsignedSignatureProperties(timestampToken *validation.TimestampToken) *XAdESUnsignedSigProperties {
	var unsignedProperties *xmldom.Node
	if timestampToken == nil {
		// timestamp creation
		unsignedProperties = b.unsignedSignaturePropertiesCanonicalizationCopy()
	} else {
		unsignedProperties = b.unsignedSignaturePropertiesDom()
	}
	if unsignedProperties == nil {
		panic("UnsignedSignatureProperties are not initialized!")
	}
	return NewXAdESUnsignedSigProperties(unsignedProperties, b.xadesPaths)
}

// xadesTimestampMessageDigestBuilderIsEn319132TimestampToken ports the private
// isEn319132TimestampToken(XAdESAttribute).
func xadesTimestampMessageDigestBuilderIsEn319132TimestampToken(timestampAttribute *XAdESAttribute) bool {
	return xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(timestampAttribute,
		definition.XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP, definition.XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP,
		definition.XAdES132Element_SIGNATURE_TIMESTAMP, definition.XAdES141Element_REFS_ONLY_TIMESTAMP_V2,
		definition.XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2, definition.XAdES141Element_ARCHIVE_TIMESTAMP)
}

// xadesTimestampMessageDigestBuilderCheckAttributeNameMatches ports the private
// checkAttributeNameMatches(XAdESAttribute, DSSElement...).
func xadesTimestampMessageDigestBuilderCheckAttributeNameMatches(attribute *XAdESAttribute, elements ...common.DSSElement) bool {
	if attribute == nil {
		return false
	}
	for _, element := range elements {
		if element.TagName() == attribute.Name() {
			return true
		}
	}
	return false
}

// objects returns the list of ds:Object elements for the current signature element. Port of the
// private getObjects().
func (b *XAdESTimestampMessageDigestBuilder) objects() []*xmldom.Node {
	nodeList, err := xmlutils.XPathUtilsGetNodeList(b.signature, common.XMLDSigPath_OBJECT_PATH)
	if err != nil {
		return nil
	}
	return nodeList
}

// writeObjectBytes ports the private writeObjectBytes(DSSMessageDigestCalculator, NodeList,
// Set<String>, String).
func (b *XAdESTimestampMessageDigestBuilder) writeObjectBytes(digestCalculator *spi.DSSMessageDigestCalculator,
	objects []*xmldom.Node, referenceURIs map[string]struct{}, canonicalizationMethod string) {
	xades141 := b.timestampToken == nil || enumerations.ArchiveTimestampType_XAdES != b.timestampToken.ArchiveTimestampType()
	for _, node := range objects {
		qualifyingProperties, err := xmlutils.XPathUtilsGetElement(node, b.xadesPaths.CurrentQualifyingPropertiesPath())
		if err == nil && qualifyingProperties != nil {
			continue
		}
		if !xades141 {
			/*
			 * !!! ETSI TS 101 903 V1.3.2 (2006-03)
			 * 5) Take any ds:Object element in the signature that is not referenced by any
			 * ds:Reference within ds:SignedInfo, except that one containing the
			 * QualifyingProperties element. Canonicalize each one and concatenate each resulting
			 * octet stream to the final octet stream. If ds:Canonicalization is present, the
			 * algorithm indicated by this element is used. If not, the standard canonicalization
			 * method specified by XMLDSIG is used.
			 */
			id := DSSXMLUtilsGetIDIdentifier(node)
			if id != "" {
				if _, contains := referenceURIs[id]; contains {
					continue
				}
			}
		}
		xadesTimestampMessageDigestBuilderWriteDigestValueOnCanonicalizedNode(digestCalculator, node, canonicalizationMethod)
	}
}

// xadesTimestampMessageDigestBuilderMustNewCalculator constructs a DSSMessageDigestCalculator for
// digestAlgorithm, panicking with a model.DSSError on failure (an unsupported digest algorithm),
// mirroring the unchecked propagation Java's `new DSSMessageDigestCalculator(digestAlgorithm)`
// exhibits within the broader `catch (Exception e)` blocks that construct it.
func xadesTimestampMessageDigestBuilderMustNewCalculator(digestAlgorithm enumerations.DigestAlgorithm) *spi.DSSMessageDigestCalculator {
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(digestAlgorithm)
	if err != nil {
		panic(model.NewDSSErrorWithCause(err))
	}
	return digestCalculator
}

// compile-time assertion: *XAdESTimestampMessageDigestBuilder satisfies
// timestamp.TimestampMessageDigestBuilder, matching Java's "implements TimestampMessageDigestBuilder".
var _ timestamp.TimestampMessageDigestBuilder = (*XAdESTimestampMessageDigestBuilder)(nil)
