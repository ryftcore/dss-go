// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/evidencerecord/XAdESEvidenceRecordDigestBuilder.java (DSS 6.5.RC1).
//
// Santuario replacement: every org.apache.xml.security call site (Reference,
// Manifest/ManifestValidator resource-resolver wiring) maps to internal/xmldsig per its doc.go
// table; the resolver-per-digest-algorithm wiring below mirrors manifest_validator.go's
// manifestValidatorInitDetachedContent (same package, not reused directly per PORTING.md's
// no-shared-cross-file-helpers rule for unexported helpers).
//
// FORWARD DEPENDENCY (root package eu.europa.esig.dss.xades, DSSXMLUtils-prefixed
// forward-dependency convention already established by xades_signature.go /
// manifest_validator.go; this file's one addition to that assumed set):
//
//	func DSSXMLUtilsGetDigestOnCanonicalizedBytes(bytes []byte, digestAlgorithm enumerations.DigestAlgorithm, canonicalizationMethod string) (model.DSSMessageDigest, error)
//	  // DSSXMLUtils.getDigestOnCanonicalizedBytes(byte[], DigestAlgorithm, String)
package xades

import (
	"fmt"
	"sort"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/spi/exception"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESEvidenceRecordDigestBuilder computes message-imprint of an XML signature to be protected
// by an evidence-record.
type XAdESEvidenceRecordDigestBuilder struct {
	validation.AbstractSignatureEvidenceRecordDigestBuilder

	// detachedContent is the list of detached documents covered by the signature.
	detachedContent []model.DSSDocument

	// signatureId is the DSS identifier or a signature's identifier of the signature to be
	// covered by an evidence-record.
	signatureId string
}

// NewXAdESEvidenceRecordDigestBuilder instantiates XAdESEvidenceRecordDigestBuilder with a
// SHA-256 digest algorithm. Port of XAdESEvidenceRecordDigestBuilder(DSSDocument).
func NewXAdESEvidenceRecordDigestBuilder(signatureDocument model.DSSDocument) *XAdESEvidenceRecordDigestBuilder {
	return &XAdESEvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilder(signatureDocument),
	}
}

// NewXAdESEvidenceRecordDigestBuilderWithAlgorithm instantiates XAdESEvidenceRecordDigestBuilder
// with a custom digest algorithm. Port of XAdESEvidenceRecordDigestBuilder(DSSDocument,
// DigestAlgorithm).
func NewXAdESEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) *XAdESEvidenceRecordDigestBuilder {
	return &XAdESEvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilderWithAlgorithm(signatureDocument, digestAlgorithm),
	}
}

// newXAdESEvidenceRecordDigestBuilderFromSignature instantiates XAdESEvidenceRecordDigestBuilder
// from a signature for the given evidenceRecordAttribute. Used on validation of an existing
// evidence record. Port of the protected XAdESEvidenceRecordDigestBuilder(AdvancedSignature,
// SignatureAttribute, DigestAlgorithm) constructor.
func newXAdESEvidenceRecordDigestBuilderFromSignature(signature validation.AdvancedSignature,
	evidenceRecordAttribute validation.SignatureAttribute, digestAlgorithm enumerations.DigestAlgorithm) *XAdESEvidenceRecordDigestBuilder {
	return &XAdESEvidenceRecordDigestBuilder{
		AbstractSignatureEvidenceRecordDigestBuilder: *validation.NewAbstractSignatureEvidenceRecordDigestBuilderFromSignature(
			signature, evidenceRecordAttribute, digestAlgorithm),
	}
}

// SetDetachedContent sets a list of detached documents covered by the signature. Port of
// #setDetachedContent, chainable.
func (b *XAdESEvidenceRecordDigestBuilder) SetDetachedContent(detachedContent []model.DSSDocument) *XAdESEvidenceRecordDigestBuilder {
	b.detachedContent = detachedContent
	return b
}

// SetSignatureId sets identifier of the signature to be covered by an evidence-record. Accepts
// a DSS identifier, or an internal signature element's identifier. Note: required for documents
// containing multiple signatures. Port of #setSignatureId, chainable.
func (b *XAdESEvidenceRecordDigestBuilder) SetSignatureId(signatureId string) *XAdESEvidenceRecordDigestBuilder {
	b.signatureId = signatureId
	return b
}

// SetParallelEvidenceRecord overrides AbstractSignatureEvidenceRecordDigestBuilder, keeping the
// concrete return type. Port of the covariant-return #setParallelEvidenceRecord(boolean)
// override.
func (b *XAdESEvidenceRecordDigestBuilder) SetParallelEvidenceRecord(parallelEvidenceRecord bool) *XAdESEvidenceRecordDigestBuilder {
	b.AbstractSignatureEvidenceRecordDigestBuilder.SetParallelEvidenceRecord(parallelEvidenceRecord)
	return b
}

// Build generates the hash value for the signature enveloping the evidence-record. Port of
// #build().
func (b *XAdESEvidenceRecordDigestBuilder) Build() (model.Digest, error) {
	xadesSignature, err := b.getXAdESSignature()
	if err != nil {
		return model.Digest{}, err
	}
	messageDigest, err := b.getXmlSignatureMessageImprint(xadesSignature)
	if err != nil {
		return model.Digest{}, err
	}
	return messageDigest.Digest, nil
}

// getXAdESSignature returns a signature to compute evidence record's digest for. Port of the
// protected #getXAdESSignature().
//
// Panics when neither an AdvancedSignature nor a signature document was provided (Java's
// IllegalStateException).
func (b *XAdESEvidenceRecordDigestBuilder) getXAdESSignature() (*XAdESSignature, error) {
	if b.Signature() != nil {
		xadesSignature, ok := b.Signature().(*XAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", b.Signature())
		}
		return xadesSignature, nil
	} else if b.SignatureDocument() != nil {
		documentAnalyzer, err := NewXMLDocumentAnalyzer(b.SignatureDocument())
		if err != nil {
			return nil, err
		}
		documentAnalyzer.SetDetachedContents(b.detachedContent)

		signatures := documentAnalyzer.Signatures()
		var signature validation.AdvancedSignature
		switch {
		case utils.CollectionSize(signatures) == 0:
			return nil, exception.NewIllegalInputException("The provided document does not contain any signature! " +
				"Unable to compute message-imprint for an integrated evidence-record.")

		case utils.IsStringNotEmpty(b.signatureId):
			signature = documentAnalyzer.SignatureByID(b.signatureId)
			if signature == nil {
				return nil, fmt.Errorf("no signature with Id '%s' found in the document!", b.signatureId)
			}

		case utils.CollectionSize(signatures) > 1:
			return nil, exception.NewIllegalInputException("The provided document contains multiple signatures! " +
				"Please use #setSignatureId method in order to provide the identifier.")

		default:
			signature = signatures[0]
		}

		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			return nil, fmt.Errorf("unexpected signature type %T", signature)
		}
		return xadesSignature, nil
	}
	panic("Either DSSDocument containing the signature or AdvancedSignature shall be defined!")
}

// getXmlSignatureMessageImprint generates message-imprint for the given XAdESSignature. Port of
// the protected #getXmlSignatureMessageImprint(XAdESSignature).
func (b *XAdESEvidenceRecordDigestBuilder) getXmlSignatureMessageImprint(signature *XAdESSignature) (model.DSSMessageDigest, error) {
	/*
	 * The initial time-stamp token encapsulated within the first ArchiveTimeStamp of any of
	 * the evidence-records enclosed within the xadesen:SealingEvidenceRecords unsigned
	 * qualifying property, shall incorporate a HashTree, whose first child shall contain the
	 * digest value of the group of data objects listed below, concatenated in the order
	 * specified in IETF RFC 6283 if the xadesen:SealingEvidenceRecords unsigned qualifying
	 * property contains XMLERS evidence-records, or in IETF RFC 4998 if the
	 * xadesen:SealingEvidenceRecords unsigned qualifying property contains ERS
	 * evidence-records:
	 */
	var digestObjectsGroup [][]byte

	// TODO : not clear which ds:CanonicalizationMethod to use. Using
	// ds:SignedInfo/ds:CanonicalizationMethod for now
	canonicalizationAlgorithm, err := b.getCanonicalizationAlgorithm(signature)
	if err != nil {
		return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
	}

	/*
	 * 1) The data objects resulting of processing each ds:Reference element within
	 * ds:SignedInfo as specified below:
	 * - Process the ds:Reference element according to the reference processing model of
	 *   XMLDSIG, clause 4.4.3.2.
	 * - If the result is a XML node set, canonicalize using the canonicalization algorithm
	 *   present in ds:CanonicalizationMethod element.
	 */
	for _, reference := range signature.References() {
		digestValue, err := b.getReferenceBytesDigestValue(reference, canonicalizationAlgorithm)
		if err != nil {
			return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
		}
		digestObjectsGroup = append(digestObjectsGroup, digestValue)
	}

	/*
	 * 2) The data objects resulting of taking the XMLDSIG elements listed below, and
	 * canonicalizing each one using the canonicalization algorithm present in
	 * ds:CanonicalizationMethod element:
	 * - The ds:SignedInfo element.
	 * - The ds:SignatureValue element.
	 * - The ds:KeyInfo element, if present.
	 */
	digestValue, err := b.getDigestValueOnCanonicalizedNodeByQuery(signature, common.XMLDSigPath_SIGNED_INFO_PATH, canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
	}
	digestObjectsGroup = append(digestObjectsGroup, digestValue)

	digestValue, err = b.getDigestValueOnCanonicalizedNodeByQuery(signature, common.XMLDSigPath_SIGNATURE_VALUE_PATH, canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
	}
	digestObjectsGroup = append(digestObjectsGroup, digestValue)

	digestValue, err = b.getDigestValueOnCanonicalizedNodeByQuery(signature, common.XMLDSigPath_KEY_INFO_PATH, canonicalizationAlgorithm)
	if err != nil {
		return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
	}
	digestObjectsGroup = append(digestObjectsGroup, digestValue)

	// Steps 3) and 4) are done together (signature is expected to be prepared)
	/*
	 * 3) The data objects resulting of taking all the unsigned qualifying properties
	 * incorporated into the XAdES signature except the xadesen:SealingEvidenceRecords element
	 * under construction, and canonicalizing each one as specified in clause 4.5 of ETSI EN
	 * 319 132-1.
	 */
	/*
	 * 4) As many xadesv141:TimeStampValidationData qualifying properties will be added as
	 * required for incorporating the validation data, not already present in the XAdES
	 * signature, that are required for validating all the time-stamp tokens incorporated
	 * (within signed or unsigned qualifying properties) into the XAdES qualifying properties
	 * different than xadesen:SealingEvidenceRecords. Each xadesv141:TimeStampValidationData
	 * qualifying property shall be generated following the specifications of ETSI EN 319
	 * 132-1. For every xadesv141:TimeStampValidationData qualifying property incorporated, the
	 * corresponding data object resulting of canonicalizing this qualifying property as
	 * specified in clause 4.5 of ETSI EN 319 132-1 will be generated and added to the group of
	 * data objects to be time stamped.
	 */
	unsignedProperties, err := b.getUnsignedSignaturePropertiesList(signature)
	if err != nil {
		return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
	}
	for _, xadesAttribute := range unsignedProperties {
		digestValue, err = b.getDigestValueOnCanonicalizedNode(xadesAttribute.Element(), canonicalizationAlgorithm)
		if err != nil {
			return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
		}
		digestObjectsGroup = append(digestObjectsGroup, digestValue)
	}

	/*
	 * 5) All the ds:Object elements except the one containing QualifyingProperties element, as
	 * specified in step 5) of clause 5.5.2.2 of ETSI EN 319 132-1.
	 */
	for _, object := range signature.Objects() {
		if !b.containsQualifyingProperties(object, signature.XAdESPaths()) {
			digestValue, err = b.getDigestValueOnCanonicalizedNode(object, canonicalizationAlgorithm)
			if err != nil {
				return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
			}
			digestObjectsGroup = append(digestObjectsGroup, digestValue)
		}
	}

	/*
	 * 6) The objects derived from the presence of signed ds:Manifest elements. These objects
	 * shall be generated as it is specified below:
	 */
	for _, reference := range signature.References() {
		if reference.TypeIsReferenceToManifest() {
			manifestDataObjectDigests, err := b.getManifestDataObjectDigests(signature, reference, canonicalizationAlgorithm)
			if err != nil {
				return model.DSSMessageDigest{}, fmt.Errorf("unable to compute message-imprint for an evidence-record. Reason : %s", err.Error())
			}
			digestObjectsGroup = append(digestObjectsGroup, manifestDataObjectDigests...)
		}
	}

	// compute final digest
	return b.computeDigestValueGroupHash(digestObjectsGroup)
}

// getCanonicalizationAlgorithm returns the corresponding ds:CanonicalizationMethod used within
// the signature. Port of the protected #getCanonicalizationAlgorithm(XAdESSignature).
func (b *XAdESEvidenceRecordDigestBuilder) getCanonicalizationAlgorithm(signature *XAdESSignature) (string, error) {
	signedInfo := signature.SignedInfo()
	if signedInfo == nil {
		return "", fmt.Errorf("ds:SignedInfo element shall be defined within a signature")
	}
	canonicalizationMethod, err := xmlutils.XPathUtilsGetValue(signedInfo, common.XMLDSigPath_CANONICALIZATION_ALGORITHM_PATH)
	if err != nil {
		return "", err
	}
	if utils.IsStringEmpty(canonicalizationMethod) {
		// Upstream warns "No canonicalization method found within ds:SignedInfo element. Re-use
		// the default canonicalization algorithm
		// 'http://www.w3.org/TR/2001/REC-xml-c14n-20010315'".
		canonicalizationMethod = xmlutils.XMLCanonicalizerDefaultXMLDSigC14NMethod
	}
	return canonicalizationMethod, nil
}

// getReferenceBytesDigestValue ports the private
// getReferenceBytesDigestValue(Reference, String).
func (b *XAdESEvidenceRecordDigestBuilder) getReferenceBytesDigestValue(reference *xmldsig.Reference, canonicalizationAlgorithm string) ([]byte, error) {
	/*
	 * 1) process the retrieved ds:Reference element according to the reference-processing
	 * model of XMLDSIG clause 4.4.3.2;
	 */
	/*
	 * 2) If the result is a XML node set, canonicalize using the canonicalization algorithm
	 * present in ds:CanonicalizationMethod element.
	 */
	if b.isResultXmlNodeSet(reference) {
		referencedBytes, err := reference.ReferencedBytes()
		if err != nil {
			return nil, fmt.Errorf("an error occurred on ds:Reference processing. In case of detached signature, "+
				"please use #setDetachedContent method to provide original documents. More information : %s", err.Error())
		}
		if xmlutils.DomUtilsIsDOMBytes(referencedBytes) {
			messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedBytes(referencedBytes, b.DigestAlgorithm(), canonicalizationAlgorithm)
			if err != nil {
				return nil, err
			}
			return messageDigest.Value(), nil
		}
		return spi.DSSUtilsDigest(b.DigestAlgorithm(), referencedBytes)
	}

	input, err := reference.ContentsAfterTransformation()
	if err != nil {
		return nil, fmt.Errorf("an error occurred on ds:Reference processing. In case of detached signature, "+
			"please use #setDetachedContent method to provide original documents. More information : %s", err.Error())
	}
	octets, err := input.Bytes()
	if err != nil {
		return nil, err
	}
	return b.getDigestValueOnBytes(octets)
}

// getDigestValueOnBytes ports the private getDigestValueOnInputStream(InputStream).
func (b *XAdESEvidenceRecordDigestBuilder) getDigestValueOnBytes(data []byte) ([]byte, error) {
	messageDigestCalculator, err := spi.NewDSSMessageDigestCalculator(b.DigestAlgorithm())
	if err != nil {
		return nil, err
	}
	messageDigestCalculator.Update(data)
	return messageDigestCalculator.MessageDigest(b.DigestAlgorithm()).Value(), nil
}

// getDigestValueOnCanonicalizedNodeByQuery ports the private
// getDigestValueOnCanonicalizedNode(XAdESSignature, XPathQuery, String).
func (b *XAdESEvidenceRecordDigestBuilder) getDigestValueOnCanonicalizedNodeByQuery(signature *XAdESSignature,
	xPathQuery common.XPathQuery, canonicalizationAlgorithm string) ([]byte, error) {
	element, err := xmlutils.XPathUtilsGetElement(signature.SignatureElement(), xPathQuery)
	if err != nil {
		return nil, err
	}
	return b.getDigestValueOnCanonicalizedNode(element, canonicalizationAlgorithm)
}

// getDigestValueOnCanonicalizedNode ports the private getDigestValueOnCanonicalizedNode(Node,
// String).
func (b *XAdESEvidenceRecordDigestBuilder) getDigestValueOnCanonicalizedNode(node *xmldom.Node, canonicalizationAlgorithm string) ([]byte, error) {
	messageDigest, err := DSSXMLUtilsGetDigestOnCanonicalizedNode(node, b.DigestAlgorithm(), canonicalizationAlgorithm)
	if err != nil {
		return nil, err
	}
	return messageDigest.Value(), nil
}

// getUnsignedSignaturePropertiesList ports the private
// getUnsignedSignaturePropertiesList(XAdESSignature).
func (b *XAdESEvidenceRecordDigestBuilder) getUnsignedSignaturePropertiesList(signature *XAdESSignature) ([]*XAdESAttribute, error) {
	// NOTE : only direct incorporation is supported
	unsignedSigProperties := XAdESUnsignedSigPropertiesBuild(signature.SignatureElement(), signature.XAdESPaths())
	if !unsignedSigProperties.IsExist() {
		// Upstream logs "No xades:UnsignedSignatureProperties is present to compute the
		// message-imprint for an evidence-record".
		return nil, nil
	}
	if b.EvidenceRecordAttribute() != nil {
		return b.getPrecedingAttributes(unsignedSigProperties, b.EvidenceRecordAttribute()), nil
	}

	if b.IsParallelEvidenceRecord() {
		erAttribute := XAdESSignatureUtilsGetLastSealingEvidenceRecordAttribute(unsignedSigProperties)
		if erAttribute != nil {
			return b.getPrecedingAttributes(unsignedSigProperties, erAttribute), nil
		}
	}
	return unsignedSigProperties.Attributes(), nil
}

// getPrecedingAttributes ports the private getPrecedingAttributes(XAdESUnsignedSigProperties,
// SignatureAttribute).
func (b *XAdESEvidenceRecordDigestBuilder) getPrecedingAttributes(unsignedSigProperties *XAdESUnsignedSigProperties,
	attribute validation.SignatureAttribute) []*XAdESAttribute {
	target, _ := attribute.(*XAdESAttribute)
	var attributes []*XAdESAttribute
	for _, currentAttribute := range unsignedSigProperties.Attributes() {
		if target != nil && target.Equals(currentAttribute) {
			break
		}
		attributes = append(attributes, currentAttribute)
	}
	return attributes
}

// containsQualifyingProperties ports the private containsQualifyingProperties(Node, XAdESPath).
func (b *XAdESEvidenceRecordDigestBuilder) containsQualifyingProperties(node *xmldom.Node, xadesPath definition.XAdESPath) bool {
	qualifyingProperties, err := xmlutils.XPathUtilsGetNode(node, xadesPath.CurrentQualifyingPropertiesPath())
	if err != nil {
		return false
	}
	return qualifyingProperties != nil
}

// getManifestDataObjectDigests ports the private getManifestDataObjectDigests(XAdESSignature,
// Reference, String).
func (b *XAdESEvidenceRecordDigestBuilder) getManifestDataObjectDigests(signature *XAdESSignature, referenceToManifest *xmldsig.Reference,
	canonicalizationAlgorithm string) ([][]byte, error) {
	var digestObjectsGroup [][]byte
	if err := b.getManifestDataObjectDigestsRecursively(signature, referenceToManifest, canonicalizationAlgorithm, &digestObjectsGroup); err != nil {
		return nil, err
	}
	return digestObjectsGroup, nil
}

// getManifestDataObjectDigestsRecursively ports the private
// getManifestDataObjectDigestsRecursively(XAdESSignature, Reference, String, List<byte[]>).
func (b *XAdESEvidenceRecordDigestBuilder) getManifestDataObjectDigestsRecursively(signature *XAdESSignature, referenceToManifest *xmldsig.Reference,
	canonicalizationAlgorithm string, digestObjectsGroup *[][]byte) error {
	manifestReferences, err := b.getManifestReferences(signature, referenceToManifest)
	if err != nil {
		return err
	}
	for _, manifestReference := range manifestReferences {
		/*
		 * a) For each ds:Reference child element of each signed ds:Manifest element retrieve
		 *    the data object referenced by its URI attribute.
		 */
		if !b.isResultXmlNodeSet(manifestReference) || !manifestReference.TypeIsReferenceToManifest() {
			/*
			 * b) If the retrieved data object is not a XML node set, or it is a XML node set
			 *    different than a ds:Manifest element, process it as specified by the reference
			 *    processing model of XMLDSIG, clause 4.4.3.2. The resulting data object shall be
			 *    added to the group of data objects to be digested.
			 */
			bytes, err := b.getReferenceBytesDigestValue(manifestReference, canonicalizationAlgorithm)
			if err != nil {
				return err
			}
			*digestObjectsGroup = append(*digestObjectsGroup, bytes)

		} else {
			/*
			 * c) If the retrieved data object is a ds:Manifest element, apply the steps 6) a) to
			 *    6) c) recursively for generating the objects to be added to the group of data
			 *    objects to be digested.
			 */
			if err := b.getManifestDataObjectDigestsRecursively(signature, manifestReference, canonicalizationAlgorithm, digestObjectsGroup); err != nil {
				return err
			}
		}
	}
	return nil
}

// getManifestReferences ports the private getManifestReferences(XAdESSignature, Reference).
func (b *XAdESEvidenceRecordDigestBuilder) getManifestReferences(signature *XAdESSignature, referenceToManifest *xmldsig.Reference) ([]*xmldsig.Reference, error) {
	uri := referenceToManifest.URI()
	manifestElement := DSSXMLUtilsGetManifestById(signature.SignatureElement(), uri)
	manifest, err := xmldsig.NewManifest(manifestElement, nil)
	if err != nil {
		return nil, err
	}
	xadesEvidenceRecordDigestBuilderInitManifestDetachedContent(manifest, b.detachedContent)
	return manifest.References()
}

// xadesEvidenceRecordDigestBuilderInitManifestDetachedContent registers one
// DetachedSignatureResolver per distinct digest algorithm found among the manifest's
// ds:Reference/ds:DigestMethod elements, mirroring manifest_validator.go's
// manifestValidatorInitDetachedContent (this file's own copy, per PORTING.md's
// no-shared-cross-file-helpers rule for unexported helpers).
func xadesEvidenceRecordDigestBuilderInitManifestDetachedContent(manifest *xmldsig.Manifest, detachedContents []model.DSSDocument) {
	if utils.IsCollectionEmpty(detachedContents) {
		return
	}
	for _, digestAlgorithm := range DSSXMLUtilsGetReferenceDigestAlgos(manifest.Element()) {
		manifest.AddResourceResolver(&xmldsig.DetachedSignatureResolver{
			Documents:       detachedContents,
			DigestAlgorithm: digestAlgorithm,
		})
	}
}

// isResultXmlNodeSet ports the private isResultXmlNodeSet(Reference).
func (b *XAdESEvidenceRecordDigestBuilder) isResultXmlNodeSet(reference *xmldsig.Reference) bool {
	return ReferenceOutputType_NODE_SET == DSSXMLUtilsGetReferenceOutputType(reference)
}

// computeDigestValueGroupHash ports the private computeDigestValueGroupHash(List<byte[]>).
func (b *XAdESEvidenceRecordDigestBuilder) computeDigestValueGroupHash(digestValueGroup [][]byte) (model.DSSMessageDigest, error) {
	/*
	 * The algorithm by which a root hash value is generated from the <HashTree> element is as
	 * follows: the content of each <DigestValue> element within the first <Sequence> element is
	 * base64 ([RFC4648], using the base64 alphabet not the base64url alphabet) decoded to
	 * obtain a binary value (representing the hash value). All collected hash values from the
	 * sequence are ordered in binary ascending order, concatenated and a new hash value is
	 * generated from that string. With one exception to this rule: when the first <Sequence>
	 * element has only one <DigestValue> element, then its binary value is added to the next
	 * list obtained from the next <Sequence> element.
	 */
	// 1. Group together items
	// NOTE: byte array already contains digest only

	// 2a. Exception
	if utils.CollectionSize(digestValueGroup) == 1 {
		return model.NewDSSMessageDigestWithValue(b.DigestAlgorithm(), digestValueGroup[0]), nil
	}
	// 2b. Binary ascending sort
	sort.Slice(digestValueGroup, func(i, j int) bool {
		return validation.ByteArrayComparatorCompare(digestValueGroup[i], digestValueGroup[j]) < 0
	})
	// 3. Concatenate
	digestCalculator, err := spi.NewDSSMessageDigestCalculator(b.DigestAlgorithm())
	if err != nil {
		return model.DSSMessageDigest{}, err
	}
	for _, hashValue := range digestValueGroup {
		digestCalculator.Update(hashValue)
	}
	// 4. Calculate hash value
	return digestCalculator.MessageDigest(b.DigestAlgorithm()), nil
}

// compile-time interface assertion.
var _ validation.SignatureEvidenceRecordDigestBuilder = (*XAdESEvidenceRecordDigestBuilder)(nil)
