// Ported from
// dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESBaselineRequirementsChecker.java
// (DSS 6.5.RC1).
//
// Performs checks according to EN 319 132-1 v1.1.1 "6.3 Requirements on XAdES signature's
// elements, qualifying properties and services".
//
// slf4j logging is dropped; every LOG.warn/LOG.debug call site is called out in
// the surrounding comment instead.
package xades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/common"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// BaselineRequirementsChecker checks conformance of a XAdES signature to the requested
// baseline format. Port of the class XAdESBaselineRequirementsChecker, extending
// validation.BaselineRequirementsChecker[Signature].
type BaselineRequirementsChecker struct {
	validation.BaselineRequirementsChecker[*Signature]

	// statusMap is the cached reference validation status map.
	statusMap map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus
}

// NewXAdESBaselineRequirementsChecker is the default constructor. Port of the public
// BaselineRequirementsChecker(Signature, CertificateVerifier) constructor.
func NewBaselineRequirementsChecker(signature *Signature, offlineCertificateVerifier validation.CertificateVerifier) *BaselineRequirementsChecker {
	checker := &BaselineRequirementsChecker{
		BaselineRequirementsChecker: validation.NewBaselineRequirementsCheckerBaseWithVerifier[*Signature](signature, offlineCertificateVerifier),
	}
	checker.InitBaselineRequirementsChecker(checker)
	return checker
}

// HasAdESProfile checks if the signature is conformant to the corresponding AdES profile. Port
// of the hasAdESProfile() override.
func (b *BaselineRequirementsChecker) HasAdESProfile() bool {
	return b.HasExtendedBESProfile() || b.HasBaselineBProfile()
}

// HasBaselineBProfile checks if the signature has a corresponding BASELINE-B profile. Port of
// the hasBaselineBProfile() override.
func (b *BaselineRequirementsChecker) HasBaselineBProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// ds:KeyInfo (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, common.XMLDSigPathKeyInfoPath) != 1 {
		// Upstream logs "ds:KeyInfo element shall be present for XAdES-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// ds:SignedInfo/ds:CanonicalizationMethod (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, common.XMLDSigPathSignedInfoCanonicalizationMethod) != 1 {
		// Upstream logs "ds:SignedInfo/ds:CanonicalizationMethod element shall be present for
		// XAdES-BASELINE-B signature (cardinality == 1)!".
		return false
	}
	// ds:SignedInfo/ds:Reference (Cardinality >= 2)
	if xadesBaselineNumberOfOccurrences(signatureElement, common.XMLDSigPathSignedInfoReferencePath) < 2 {
		// Upstream logs "ds:SignedInfo/ds:Reference element shall be present for
		// XAdES-BASELINE-B signature (cardinality >= 2)!".
		return false
	}
	// ds:SignedInfo/ds:Reference/ds:Transforms (Cardinality 0 or 1)
	referenceList, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, common.XMLDSigPathSignedInfoReferencePath)
	for _, reference := range referenceList {
		if xadesBaselineNumberOfOccurrences(reference, common.XMLDSigPathTransformsPath) > 1 {
			// Upstream logs "Only one ds:Reference/ds:Transforms may be present for
			// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
			return false
		}
	}
	// SigningTime (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningTimePath()) != 1 {
		// Upstream logs "SigningTime shall be present for XAdES-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// SigningCertificate/SigningCertificateV2 (Cardinality == 1)
	if !xadesBaselineIsSigningCertificatePresent(signatureElement, xadesPaths) {
		// Upstream logs "SigningCertificate(V2) shall be present for XAdES-BASELINE-B signature
		// (cardinality == 1)!".
		return false
	}
	// DataObjectFormat (Cardinality >= 0)
	referenceStatusMap := b.referenceValidationStatusMap()
	dataObjectFormatList := xadesBaselineDataObjectFormatList(signatureElement, xadesPaths)
	for _, dataObjectFormat := range dataObjectFormatList {
		if !xadesBaselineIsValidXAdESBaselineDataObjectFormat(dataObjectFormat, signature, xadesPaths, referenceStatusMap) {
			return false
		}
	}
	// SignerRole/SignerRoleV2 (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignerRolePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignerRoleV2Path()) > 1 {
		// Upstream logs "Only one SignerRole(V2) may be present for XAdES-BASELINE-B signature
		// (cardinality 0 or 1)!".
		return false
	}
	// CommitmentTypeIndication (Cardinality >= 0)
	// ProductionPlace/SignatureProductionPlaceV2 (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignatureProductionPlacePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignatureProductionPlaceV2Path()) > 1 {
		// Upstream logs "Only one SignatureProductionPlace(V2) may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	// CounterSignature (Cardinality >= 0)
	// AllDataObjectsTimeStamp (Cardinality >= 0)
	// SignaturePolicyIdentifier (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignaturePolicyIdentifierPath()) > 1 {
		// Upstream logs "Only one SignaturePolicyIdentifier may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	// SignaturePolicyStore (Cardinality 0 or 1, conditioned presence requirement (m))
	signaturePolicyStoreAmount := xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignaturePolicyStorePath())
	if signaturePolicyStoreAmount == 1 {
		if !b.IsSignaturePolicyIdentifierHashPresent() {
			// Upstream logs "SignaturePolicyStore shall not be present for XAdES-BASELINE-B
			// signature with not defined SignaturePolicyIdentifier/SigPolicyHash (requirement
			// (m))!".
		}
	} else if signaturePolicyStoreAmount > 1 {
		// Upstream logs "Only one SignaturePolicyIdentifier may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	// ArchiveTimeStamp (defined in namespace whose URI is "http://uri.etsi.org/01903/v1.3.2#")
	// (Cardinality == 0)
	archiveTimestampPath := xadesPaths.ArchiveTimestampPath()
	if archiveTimestampPath != nil {
		archiveTimeStampList, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, archiveTimestampPath)
		for _, archiveTimeStamp := range archiveTimeStampList {
			if definition.XAdESNamespaceXAdES132.Uri() == archiveTimeStamp.Name.Space {
				// Upstream logs "xades132:ArchiveTimeStamp shall not be present for
				// XAdES-BASELINE-B signature (cardinality == 0)!".
				return false
			}
		}
	}
	// Additional requirement (a)
	if !b.ContainsSigningCertificate(signature.CertificateSource().KeyInfoCertificates()) {
		// Upstream logs "Signing certificate shall be present in
		// ds:KeyInfo/ds:X509Data/ds:X509Certificate for XAdES-BASELINE-B signature (requirement
		// (a))!".
		return false
	}
	// Additional requirement (d)
	signedInfo := signature.SignedInfo()
	if signedInfo != nil {
		canonicalizationMethod, _ := xmlutils.XPathUtilsGetValue(signedInfo, common.XMLDSigPathCanonicalizationAlgorithmPath)
		if utils.IsStringNotEmpty(canonicalizationMethod) {
			switch canonicalizationMethod {
			case xmldsig.TransformC14N11WithComments, xmldsig.TransformC14NExclWithComments, xmldsig.TransformC14NWithComments,
				xmldsig.TransformC14N11, xmldsig.TransformC14NExcl, xmldsig.TransformC14N:
				// ok
			default:
				// Upstream logs "ds:SignedInfo/ds:CanonicalizationMethod contains not accepted
				// Algorithm attribute value for XAdES-BASELINE-B signature (requirement (d))!".
				return false
			}
		}
	}
	// Additional requirement (i)
	signingCertificateV2Path := xadesPaths.SigningCertificateV2Path()
	if signingCertificateV2Path != nil {
		signingCertificateV2List, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, signingCertificateV2Path)
		if len(signingCertificateV2List) == 1 {
			signingCertificateV2 := signingCertificateV2List[0]
			certList, _ := xmlutils.XPathUtilsGetNodeList(signingCertificateV2, xadesPaths.CurrentCertChildren())
			for _, cert := range certList {
				if cert.Attr("", definition.XAdES132AttributeURI.AttributeName()) != nil {
					// Upstream logs "SigningCertificateV2/Cert shall not include URI optional
					// attribute for XAdES-BASELINE-B signature (requirement (i))!".
					return false
				}
			}
		}
	}
	// Additional requirement (k)
	for _, referenceStatus := range referenceStatusMap {
		if referenceStatus.dataObjectFormatFound {
			continue
		}
		referenceValidation := referenceStatus.referenceValidation
		// TODO : check whether other reference types should be checked (i.e. KeyInfo, Manifest, etc.)
		uri := referenceValidation.Uri()
		if (uri == "" || xmlutils.DomUtilsStartsFromHash(uri) || xmlutils.DomUtilsIsXPointerQuery(uri)) &&
			(enumerations.DigestMatcherTypeSignedProperties == referenceValidation.Type() ||
				enumerations.DigestMatcherTypeCounterSignature == referenceValidation.Type() ||
				enumerations.DigestMatcherTypeCounterSignedSignatureValue == referenceValidation.Type() ||
				enumerations.DigestMatcherTypeManifest == referenceValidation.Type() ||
				enumerations.DigestMatcherTypeKeyInfo == referenceValidation.Type() ||
				enumerations.DigestMatcherTypeSignatureProperties == referenceValidation.Type()) {
			continue
		}
		// Upstream logs "DataObjectFormat shall be generated for each signed data for
		// XAdES-BASELINE-B signature (requirement (k))!".
		return false
	}
	return true
}

// HasBaselineTProfile checks if the signature has a corresponding BASELINE-T profile. Port of
// the hasBaselineTProfile() override.
func (b *BaselineRequirementsChecker) HasBaselineTProfile() bool {
	if !b.MinimalTRequirement() {
		return false
	}
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// Additional requirement (n)
	signatureTimeStampList, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, xadesPaths.SignatureTimestampPath())
	for _, signatureTimeStamp := range signatureTimeStampList {
		encapsulatedTimestampList, _ := xmlutils.XPathUtilsGetNodeList(signatureTimeStamp, xadesPaths.CurrentEncapsulatedTimestamp())
		if len(encapsulatedTimestampList) != 1 {
			// Upstream logs "SignatureTimeStamp shall contain only one electronic timestamp for
			// XAdES-BASELINE-T signature (requirement (n))!".
			return false
		}
	}
	// Additional requirement (o)
	if !b.SignatureTimestampsCreatedBeforeSignCertExpiration() {
		// Upstream logs "SignatureTimeStamp shall be created before expiration of the
		// signing-certificate for XAdES-BASELINE-T signature (requirement (o))!".
		return false
	}
	return true
}

// HasBaselineLTProfile checks if the signature has a corresponding BASELINE-LT profile. Port of
// the hasBaselineLTProfile() override.
func (b *BaselineRequirementsChecker) HasBaselineLTProfile() bool {
	if !b.MinimalLTRequirement() {
		return false
	}
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// CertificateValues (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CertificateValuesPath()) > 1 {
		// Upstream logs "Only one CertificateValues element may be present for
		// XAdES-BASELINE-LT signature (cardinality 0 or 1)!".
		return false
	}
	// CompleteCertificateRefs/CompleteCertificateRefsV2 (Cardinality == 0)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteCertificateRefsPath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteCertificateRefsV2Path()) > 0 {
		// Upstream logs "CompleteCertificateRefs(V2) shall not be present for
		// XAdES-BASELINE-LT signature (cardinality == 0)!".
		return false
	}
	// AttrAuthoritiesCertValues (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.AttrAuthoritiesCertValuesPath()) > 1 {
		// Upstream logs "Only one AttrAuthoritiesCertValues element may be present for
		// XAdES-BASELINE-LT signature (cardinality 0 or 1)!".
		return false
	}
	// RevocationValues (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.RevocationValuesPath()) > 1 {
		// Upstream logs "Only one RevocationValues element may be present for
		// XAdES-BASELINE-LT signature (cardinality 0 or 1)!".
		return false
	}
	// CompleteRevocationRefs (Cardinality == 0)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteRevocationRefsPath()) > 0 {
		// Upstream logs "CompleteRevocationRefs shall not be present for XAdES-BASELINE-LT
		// signature (cardinality == 0)!".
		return false
	}
	// AttributeRevocationValues (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.AttributeRevocationValuesPath()) > 1 {
		// Upstream logs "AttributeRevocationValues shall not be present for
		// XAdES-BASELINE-LT signature (cardinality 0 or 1)!".
		return false
	}
	// SigAndRefsTimeStamp/SigAndRefsTimeStampV2 (Cardinality == 0)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigAndRefsTimestampPath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigAndRefsTimestampV2Path()) > 0 {
		// Upstream logs "SigAndRefsTimeStamp(V2) shall not be present for XAdES-BASELINE-LT
		// signature (cardinality == 0)!".
		return false
	}
	// RefsOnlyTimeStamp/RefsOnlyTimeStampV2 (Cardinality == 0)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.RefsOnlyTimestampPath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.RefsOnlyTimestampV2Path()) > 0 {
		// Upstream logs "RefsOnlyTimeStampV2 shall not be present for XAdES-BASELINE-LT
		// signature (cardinality == 0)!".
		return false
	}
	return true
}

// ContainsLTLevelCertificates implements validation.BaselineRequirementsCheckerOverrides. Port
// of the protected containsLTLevelCertificates() override.
func (b *BaselineRequirementsChecker) ContainsLTLevelCertificates() bool {
	return b.containsCertificateValues() || b.containsTstOrAnyValDataCertificates()
}

// containsCertificateValues ports the private containsCertificateValues().
func (b *BaselineRequirementsChecker) containsCertificateValues() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()
	return xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CertificateValuesPath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.AttrAuthoritiesCertValuesPath()) != 0
}

// containsTstOrAnyValDataCertificates ports the private containsTstOrAnyValDataCertificates().
func (b *BaselineRequirementsChecker) containsTstOrAnyValDataCertificates() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()
	return xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.EncapsulatedTimeStampValidationDataCertValuesPath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.EncapsulatedAnyValidationDataCertValuesPath()) != 0
}

// HasBaselineLTAProfile checks if the signature has a corresponding BASELINE-LTA profile. Port
// of the hasBaselineLTAProfile() override.
func (b *BaselineRequirementsChecker) HasBaselineLTAProfile() bool {
	return b.MinimalLTARequirement()
}

// HasExtendedBESProfile checks if the signature has a corresponding *AdES-BES profile. Port of
// the hasExtendedBESProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedBESProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// SigningTime (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningTimePath()) > 1 {
		// Upstream logs "Only one SigningTime may be present for XAdES-BES signature
		// (cardinality 0 or 1)!".
		return false
	}
	// SigningCertificate/SigningCertificateV2 (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningCertificatePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningCertificateV2Path()) > 1 {
		// Upstream logs "Only one SigningCertificate(V2) may be present for XAdES-BES signature
		// (cardinality 0 or 1)!".
		return false
	}
	// CommitmentTypeIndication (Cardinality >= 0)
	// DataObjectFormat (Cardinality >= 0)
	referenceStatusMap := b.referenceValidationStatusMap()
	dataObjectFormatList := xadesBaselineDataObjectFormatList(signatureElement, xadesPaths)
	for _, dataObjectFormat := range dataObjectFormatList {
		if !xadesBaselineIsValidXAdESDataObjectFormat(dataObjectFormat, signature, xadesPaths, referenceStatusMap) {
			return false
		}
	}
	// ProductionPlace/SignatureProductionPlaceV2 (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignatureProductionPlacePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignatureProductionPlaceV2Path()) > 1 {
		// Upstream logs "Only one SignatureProductionPlace(V2) may be present for XAdES-BES
		// signature (cardinality 0 or 1)!".
		return false
	}
	// SignerRole/SignerRoleV2 (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignerRolePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignerRoleV2Path()) > 1 {
		// Upstream logs "Only one SignerRole(V2) may be present for XAdES-BES signature
		// (cardinality 0 or 1)!".
		return false
	}
	// CounterSignature (Cardinality >= 0)
	// AllDataObjectsTimeStamp (Cardinality >= 0)
	// IndividualDataObjectsTimeStamp (Cardinality >= 0)
	// Additional requirement (a)
	if !xadesBaselineIsSigningCertificatePresent(signatureElement, xadesPaths) && !b.isSigningCertificateSignedInKeyInfo() {
		// Upstream logs "SigningCertificate(V2) shall be present for XAdES-BES signature or be
		// present in ds:KeyInfo and signed by the signature (requirement (a))!".
		return false
	}

	return true
}

// HasExtendedEPESProfile checks if the signature has a corresponding *AdES-EPES profile. Port
// of the hasExtendedEPESProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedEPESProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// SignaturePolicyIdentifier (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignaturePolicyIdentifierPath()) != 1 {
		// Upstream logs "SignaturePolicyIdentifier shall be present for XAdES-EPES signature
		// (cardinality == 1)!".
		return false
	}
	// SignaturePolicyStore (Cardinality == 0)
	signaturePolicyStoreOccurrences := xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SignaturePolicyStorePath())
	if signaturePolicyStoreOccurrences > 1 {
		// Upstream logs "Only one SignaturePolicyStore may be present for XAdES-EPES signature
		// (cardinality 0 or 1)!".
		return false
	}
	// Additional requirement (c)
	if signaturePolicyStoreOccurrences == 1 && !b.IsSignaturePolicyIdentifierHashPresent() {
		// Upstream logs "SignaturePolicyStore may be present for XAdES-EPES signature only if
		// SignaturePolicyIdentifier is present and it contains SigPolicyHash element
		// (requirement (c))!".
		return false
	}

	return true
}

// HasExtendedTProfile checks if the signature has a corresponding *AdES-T profile. Port of the
// hasExtendedTProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedTProfile() bool {
	if !b.MinimalTRequirement() {
		return false
	}
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// Additional requirement (d)
	signatureTimeStampList, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, xadesPaths.SignatureTimestampPath())
	for _, signatureTimeStamp := range signatureTimeStampList {
		encapsulatedTimestampList, _ := xmlutils.XPathUtilsGetNodeList(signatureTimeStamp, xadesPaths.CurrentEncapsulatedTimestamp())
		if len(encapsulatedTimestampList) == 0 {
			// Upstream logs "SignatureTimeStamp shall contain one or more electronic timestamp
			// for XAdES-T signature (requirement (d))!".
			return false
		}
	}
	// Additional requirement (e)
	if !b.SignatureTimestampsCreatedBeforeSignCertExpiration() {
		// Upstream logs "SignatureTimeStamp shall be created before expiration of the
		// signing-certificate for XAdES-T signature (requirement (e))!".
		return false
	}

	return true
}

// HasExtendedCProfile checks if the signature has a corresponding *AdES-C profile. Port of the
// hasExtendedCProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedCProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// NOTE: at least complete-certificate-references shall be present for all self-signed
	// certificates
	// CompleteCertificateRefs/CompleteCertificateRefsV2 (Cardinality == 1)
	completeCertificateRefsNumberOfOccurrences := xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteCertificateRefsPath()) +
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteCertificateRefsV2Path())
	if completeCertificateRefsNumberOfOccurrences > 1 || completeCertificateRefsNumberOfOccurrences == 0 {
		// Upstream logs "CompleteCertificateRefs(V2) shall be present for XAdES-C signature
		// (cardinality == 1)!".
		return false
	}

	certificateSources := b.CertificateSourcesExceptLastArchiveTimestamp()
	certificateFound := certificateSources.NumberOfCertificates() > 0
	allSelfSigned := certificateFound && certificateSources.IsAllSelfSigned()

	// CompleteRevocationRefs (Cardinality == 1)
	completeRevocationRefsNumberOfOccurrences := xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.CompleteRevocationRefsPath())
	if completeRevocationRefsNumberOfOccurrences > 1 || (!allSelfSigned && completeRevocationRefsNumberOfOccurrences == 0) {
		// Upstream logs "CompleteRevocationRefs shall be present for XAdES-C signature
		// (cardinality == 1)!".
		return false
	}
	return true
}

// HasExtendedXProfile checks if the signature has a corresponding *AdES-X profile. Port of the
// hasExtendedXProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedXProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	refsOnlyTst := xadesBaselineIsElementPresent(signatureElement, xadesPaths.RefsOnlyTimestampPath())
	refsOnlyTstV2 := xadesBaselineIsElementPresent(signatureElement, xadesPaths.RefsOnlyTimestampV2Path())
	sigAndRefsTst := xadesBaselineIsElementPresent(signatureElement, xadesPaths.SigAndRefsTimestampPath())
	sigAndRefsTstV2 := xadesBaselineIsElementPresent(signatureElement, xadesPaths.SigAndRefsTimestampV2Path())
	if !refsOnlyTst && !refsOnlyTstV2 && !sigAndRefsTst && !sigAndRefsTstV2 {
		// Upstream logs "Either RefsOnlyTimestamp(V2) or SigAndRefsTimestamp(V2) shall be
		// present for XAdES-X signature)!".
		return false
	}
	return true
}

// HasExtendedXLProfile checks if the signature has a corresponding *AdES-XL profile. Port of the
// hasExtendedXLProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedXLProfile() bool {
	return b.MinimalLTRequirement()
}

// HasExtendedAProfile checks if the signature has a corresponding *AdES-A profile. Port of the
// hasExtendedAProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedAProfile() bool {
	return b.MinimalLTARequirement()
}

// HasExtendedERSProfile checks if the signature has a corresponding *AdES-E-ERS profile. Port of
// the hasExtendedERSProfile() override.
func (b *BaselineRequirementsChecker) HasExtendedERSProfile() bool {
	signature := b.Signature()
	signatureElement := signature.SignatureElement()
	xadesPaths := signature.XAdESPaths()

	// SigningTime (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningTimePath()) != 1 {
		// Upstream logs "SigningTime shall be present for XAdES-E-ERS signature
		// (cardinality == 1)!".
		return false
	}
	// TODO : verify validation data presence for all except last timestamp (minimalLTRequirement) ?
	// other requirements are skipped, as same as for XAdES-LT/XAdES-XL
	// xadesen:SealingEvidenceRecords (Cardinality >= 1)
	if xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SealingEvidenceRecordsPath()) == 0 {
		// Upstream logs "xadesen:SealingEvidenceRecords shall be present for XAdES-E-ERS
		// signature (cardinality >= 1)!".
		return false
	}
	return true
}

// xadesBaselineIsSigningCertificatePresent ports the private isSigningCertificatePresent(Element, XAdESPath).
func xadesBaselineIsSigningCertificatePresent(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) bool {
	return xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningCertificatePath())+
		xadesBaselineNumberOfOccurrences(signatureElement, xadesPaths.SigningCertificateV2Path()) == 1
}

// xadesBaselineDataObjectFormatList ports the private getDataObjectFormatList(Element, XAdESPath).
func xadesBaselineDataObjectFormatList(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) []*xmldom.Node {
	nodeList, _ := xmlutils.XPathUtilsGetNodeList(signatureElement, xadesPaths.DataObjectFormat())
	return nodeList
}

// xadesBaselineIsValidXAdESDataObjectFormat ports the private isValidXAdESDataObjectFormat(Element,
// Signature, XAdESPath, Map).
func xadesBaselineIsValidXAdESDataObjectFormat(dataObjectFormat *xmldom.Node, signature *Signature, xadesPaths definition.XAdESPath,
	referenceStatusMap map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus) bool {
	// 5.2.4 The DataObjectFormat qualifying property
	signatureElement := signature.SignatureElement()
	// This qualifying property shall contain at least one of the following elements:
	// Description, ObjectIdentifier and MimeType.
	if xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentDescription()) != 1 &&
		xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentObjectIdentifier()) != 1 &&
		xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentMimeType()) != 1 {
		// Upstream logs "At least one of the following elements: Description, ObjectIdentifier
		// and MimeType shall be present within DataObjectFormat element for XAdES signature!".
		return false
	}
	/*
	 * The ObjectReference attribute shall reference the ds:Reference child of the
	 * ds:SignedInfo or a signed ds:Manifest element referencing the signed data object
	 * qualified by this qualifying property.
	 */
	objectReferenceAttr := dataObjectFormat.Attr("", definition.XAdES132AttributeObjectReference.AttributeName())
	objectReference := ""
	if objectReferenceAttr != nil {
		objectReference = objectReferenceAttr.Value
	}
	if utils.IsStringEmpty(objectReference) {
		// Upstream logs "DataObjectFormat's ObjectReference attribute shall be present for
		// XAdES signature (cardinality == 1)!".
		return false
	}
	if !xmlutils.DomUtilsIsElementReference(objectReference) {
		// Upstream logs "DataObjectFormat's ObjectReference attribute shall refer to an element
		// within the document!".
		return false
	}
	id := xmlutils.DomUtilsGetId(objectReference)
	matchingReference := xadesBaselineMatchingReference(id, referenceStatusMap)
	if matchingReference == nil {
		// Upstream logs "DataObjectFormat's ObjectReference attribute shall refer to a signed
		// data object within the document!".
		return false
	}
	/*
	 * If the DataObjectFormat qualifying property references a ds:Reference that in turn
	 * references a ds:Object within the XAdES signature, and if this ds:Object element has the
	 * MimeType or (and) the Encoding attribute(s), then DataObjectFormat's children MimeType
	 * and Encoding shall have exactly the same values, if they are present.
	 */
	return xadesBaselineIsDataObjectFormatValuesCompliant(dataObjectFormat, matchingReference, signatureElement, xadesPaths)
}

// xadesBaselineMatchingReference ports the private getMatchingReference(String, Map).
func xadesBaselineMatchingReference(id string, referenceStatusMap map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus) *model.ReferenceValidation {
	referenceValidationStatus, ok := referenceStatusMap[id]
	if ok {
		referenceValidationStatus.dataObjectFormatFound = true
		return referenceValidationStatus.referenceValidation
	}
	return nil
}

// xadesBaselineIsDataObjectFormatValuesCompliant ports the private
// isDataObjectFormatValuesCompliant(Element, ReferenceValidation, Element, XAdESPath).
func xadesBaselineIsDataObjectFormatValuesCompliant(dataObjectFormat *xmldom.Node, reference *model.ReferenceValidation, signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) bool {
	if xmlutils.DomUtilsIsElementReference(reference.Uri()) {
		dataObjectFormatMimeType, _ := xmlutils.XPathUtilsGetElement(dataObjectFormat, xadesPaths.CurrentMimeType())
		if dataObjectFormatMimeType != nil {
			object := DSSXMLUtilsGetObjectById(signatureElement, reference.Uri())
			if object != nil {
				objectMimeType := object.AttrValue("", common.XMLDSigAttributeMIMEType.AttributeName())
				if utils.IsStringNotEmpty(objectMimeType) && objectMimeType != dataObjectFormatMimeType.TextContent() {
					// Upstream logs "DataObjectFormat's MimeType attribute shall have the same
					// value as the corresponding signed ds:Object element, when present!".
					return false
				}
			}
		}
		dataObjectFormatEncoding, _ := xmlutils.XPathUtilsGetElement(dataObjectFormat, xadesPaths.CurrentEncoding())
		if dataObjectFormatEncoding != nil {
			object := DSSXMLUtilsGetObjectById(signatureElement, reference.Uri())
			if object != nil {
				objectEncoding := object.AttrValue("", common.XMLDSigAttributeEncoding.AttributeName())
				if utils.IsStringNotEmpty(objectEncoding) && objectEncoding != dataObjectFormatEncoding.TextContent() {
					// Upstream logs "DataObjectFormat's Encoding attribute shall have the same
					// value as the corresponding signed ds:Object element, when present!".
					return false
				}
			}
		}
	}
	return true
}

// xadesBaselineIsValidXAdESBaselineDataObjectFormat ports the private
// isValidXAdESBaselineDataObjectFormat(Element, Signature, XAdESPath, Map).
func xadesBaselineIsValidXAdESBaselineDataObjectFormat(dataObjectFormat *xmldom.Node, signature *Signature, xadesPaths definition.XAdESPath,
	referenceStatusMap map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus) bool {
	if !xadesBaselineIsValidXAdESDataObjectFormat(dataObjectFormat, signature, xadesPaths, referenceStatusMap) {
		return false
	}
	// DataObjectFormat/Description (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentDescription()) > 1 {
		// Upstream logs "Only one DataObjectFormat/Description may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	// DataObjectFormat/ObjectIdentifier (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentObjectIdentifier()) > 1 {
		// Upstream logs "Only one DataObjectFormat/ObjectIdentifier may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	// DataObjectFormat/MimeType (Cardinality == 1)
	if xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentMimeType()) != 1 {
		// Upstream logs "DataObjectFormat/MimeType shall be present for XAdES-BASELINE-B
		// signature (cardinality == 1)!".
		return false
	}
	// DataObjectFormat/Encoding (Cardinality 0 or 1)
	if xadesBaselineNumberOfOccurrences(dataObjectFormat, xadesPaths.CurrentEncoding()) > 1 {
		// Upstream logs "Only one DataObjectFormat/Encoding may be present for
		// XAdES-BASELINE-B signature (cardinality 0 or 1)!".
		return false
	}
	return true
}

// isSigningCertificateSignedInKeyInfo ports the private isSigningCertificateSignedInKeyInfo().
func (b *BaselineRequirementsChecker) isSigningCertificateSignedInKeyInfo() bool {
	signature := b.Signature()
	signingCertificate := signature.SigningCertificateToken()
	if signingCertificate != nil && SignatureUtilsIsKeyInfoCovered(signature) {
		keyInfoCertificates := signature.CertificateSource().KeyInfoCertificates()
		for _, keyInfoCertificate := range keyInfoCertificates {
			if signingCertificate.Equals(keyInfoCertificate) {
				return true
			}
		}
	}
	return false
}

// xadesBaselineNumberOfOccurrences ports the private getNumberOfOccurrences(Element, XPathQuery).
func xadesBaselineNumberOfOccurrences(element *xmldom.Node, xPath common.XPathQuery) int {
	if element == nil || xPath == nil {
		return 0
	}
	amount, err := xmlutils.XPathUtilsGetNodesAmount(element, xPath)
	if err != nil {
		return 0
	}
	return amount
}

// xadesBaselineIsElementPresent ports the private isElementPresent(Element, XPathQuery).
func xadesBaselineIsElementPresent(element *xmldom.Node, xPathString common.XPathQuery) bool {
	return xadesBaselineNumberOfOccurrences(element, xPathString) > 0
}

// referenceValidationStatusMap ports the private getReferenceValidationStatusMap().
func (b *BaselineRequirementsChecker) referenceValidationStatusMap() map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus {
	if b.statusMap == nil {
		b.statusMap = make(map[string]*xadesBaselineRequirementsCheckerReferenceValidationStatus)
		for _, referenceValidation := range b.Signature().ReferenceValidations() {
			b.statusMap[referenceValidation.Id()] = newXadesBaselineRequirementsCheckerReferenceValidationStatus(referenceValidation)
			if utils.IsCollectionNotEmpty(referenceValidation.DependentValidations()) {
				for _, dependentReference := range referenceValidation.DependentValidations() {
					b.statusMap[dependentReference.Id()] = newXadesBaselineRequirementsCheckerReferenceValidationStatusManifestEntry(dependentReference)
				}
			}
		}
	}
	return b.statusMap
}

// xadesBaselineRequirementsCheckerReferenceValidationStatus is a helper to evaluate
// ReferenceValidation status efficiently. Port of the private static final class
// ReferenceValidationStatus.
type xadesBaselineRequirementsCheckerReferenceValidationStatus struct {
	// referenceValidation contains the validation result of the signature ds:Reference element.
	referenceValidation *model.ReferenceValidation

	// dataObjectFormatFound reports whether the associated DataObjectFormat element has been
	// found.
	dataObjectFormatFound bool
}

// newXadesBaselineRequirementsCheckerReferenceValidationStatus is the port of the private
// ReferenceValidationStatus(ReferenceValidation) constructor.
func newXadesBaselineRequirementsCheckerReferenceValidationStatus(referenceValidation *model.ReferenceValidation) *xadesBaselineRequirementsCheckerReferenceValidationStatus {
	return &xadesBaselineRequirementsCheckerReferenceValidationStatus{referenceValidation: referenceValidation}
}

// newXadesBaselineRequirementsCheckerReferenceValidationStatusManifestEntry is the port of the
// private ReferenceValidationStatus(ReferenceValidation, boolean) constructor used for a
// manifest entry (dependent reference validation).
func newXadesBaselineRequirementsCheckerReferenceValidationStatusManifestEntry(referenceValidation *model.ReferenceValidation) *xadesBaselineRequirementsCheckerReferenceValidationStatus {
	return &xadesBaselineRequirementsCheckerReferenceValidationStatus{referenceValidation: referenceValidation, dataObjectFormatFound: true}
}

// compile-time assertion: an BaselineRequirementsChecker satisfies its own overrides
// contract and the non-generic BaselineRequirementsCheckerContract DefaultAdvancedSignature
// stores it as.
var (
	_ validation.BaselineRequirementsCheckerOverrides = (*BaselineRequirementsChecker)(nil)
	_ validation.BaselineRequirementsCheckerContract  = (*BaselineRequirementsChecker)(nil)
)
