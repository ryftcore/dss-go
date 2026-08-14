// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades141/XAdES141Element.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES141Element defines the XAdES 1.4.1 elements (the elements XAdES 1.4.1 adds on
// top of XAdES 1.3.2; this type implements only common.DSSElement, not XAdESElement, since
// upstream's XAdES141Element implements plain DSSElement).
type XAdES141Element string

const (
	XAdES141Element_ANY_VALIDATION_DATA           XAdES141Element = "ANY_VALIDATION_DATA"
	XAdES141Element_ARCHIVE_TIMESTAMP             XAdES141Element = "ARCHIVE_TIMESTAMP"
	XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2 XAdES141Element = "ATTRIBUTE_CERTIFICATE_REFS_V2"
	XAdES141Element_CERT_REFS                     XAdES141Element = "CERT_REFS"
	XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2  XAdES141Element = "COMPLETE_CERTIFICATE_REFS_V2"
	XAdES141Element_NEW_SDO_DIGEST_VALUE          XAdES141Element = "NEW_SDO_DIGEST_VALUE"
	XAdES141Element_ORIGINAL_REF_DIGEST           XAdES141Element = "ORIGINAL_REF_DIGEST"
	XAdES141Element_RECOMPUTED_DIGEST_VALUE       XAdES141Element = "RECOMPUTED_DIGEST_VALUE"
	XAdES141Element_REFS_ONLY_TIMESTAMP_V2        XAdES141Element = "REFS_ONLY_TIMESTAMP_V2"
	XAdES141Element_RENEWED_DIGESTS               XAdES141Element = "RENEWED_DIGESTS"
	XAdES141Element_RENEWED_DIGESTS_V2            XAdES141Element = "RENEWED_DIGESTS_V2"
	XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2     XAdES141Element = "SIG_AND_REFS_TIMESTAMP_V2"
	XAdES141Element_SIGNATURE_POLICY_DOCUMENT     XAdES141Element = "SIGNATURE_POLICY_DOCUMENT"
	XAdES141Element_SIGNATURE_POLICY_STORE        XAdES141Element = "SIGNATURE_POLICY_STORE"
	XAdES141Element_SIG_POL_DOC_LOCAL_URI         XAdES141Element = "SIG_POL_DOC_LOCAL_URI"
	XAdES141Element_SP_DOC_SPECIFICATION          XAdES141Element = "SP_DOC_SPECIFICATION"
	XAdES141Element_TIMESTAMP_VALIDATION_DATA     XAdES141Element = "TIMESTAMP_VALIDATION_DATA"
)

// xades141ElementTagNames maps each constant to its wire tag name (getTagName()).
var xades141ElementTagNames = map[XAdES141Element]string{
	XAdES141Element_ANY_VALIDATION_DATA:           "AnyValidationData",
	XAdES141Element_ARCHIVE_TIMESTAMP:             "ArchiveTimeStamp",
	XAdES141Element_ATTRIBUTE_CERTIFICATE_REFS_V2: "AttributeCertificateRefsV2",
	XAdES141Element_CERT_REFS:                     "CertRefs",
	XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2:  "CompleteCertificateRefsV2",
	XAdES141Element_NEW_SDO_DIGEST_VALUE:          "NewSDODigestValue",
	XAdES141Element_ORIGINAL_REF_DIGEST:           "OriginalRefDigest",
	XAdES141Element_RECOMPUTED_DIGEST_VALUE:       "RecomputedDigestValue",
	XAdES141Element_REFS_ONLY_TIMESTAMP_V2:        "RefsOnlyTimeStampV2",
	XAdES141Element_RENEWED_DIGESTS:               "RenewedDigests",
	XAdES141Element_RENEWED_DIGESTS_V2:            "RenewedDigestsV2",
	XAdES141Element_SIG_AND_REFS_TIMESTAMP_V2:     "SigAndRefsTimeStampV2",
	XAdES141Element_SIGNATURE_POLICY_DOCUMENT:     "SignaturePolicyDocument",
	XAdES141Element_SIGNATURE_POLICY_STORE:        "SignaturePolicyStore",
	XAdES141Element_SIG_POL_DOC_LOCAL_URI:         "SigPolDocLocalURI",
	XAdES141Element_SP_DOC_SPECIFICATION:          "SPDocSpecification",
	XAdES141Element_TIMESTAMP_VALIDATION_DATA:     "TimeStampValidationData",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES141Element) TagName() string {
	return xades141ElementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES141Element) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_141
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES141Element) URI() string {
	return XAdESNamespace_XADES_141.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e XAdES141Element) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ common.DSSElement = XAdES141Element("")
