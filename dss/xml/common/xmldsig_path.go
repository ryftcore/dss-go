// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigPath.java (DSS 6.5.RC1).
package common

// The "Type" attribute values a ds:Reference element may carry.
const (
	// XMLDSigPath_OBJECT_TYPE is the "Type" attribute value for a ds:Reference element
	// referring to a ds:Object element.
	XMLDSigPath_OBJECT_TYPE = "http://www.w3.org/2000/09/xmldsig#Object"

	// XMLDSigPath_MANIFEST_TYPE is the "Type" attribute value for a ds:Reference element
	// referring to a signed manifest.
	XMLDSigPath_MANIFEST_TYPE = "http://www.w3.org/2000/09/xmldsig#Manifest"

	// XMLDSigPath_COUNTER_SIGNATURE_TYPE is the "Type" attribute value for a ds:Reference
	// element pointing at a ds:SignatureValue of a counter-signed signature.
	XMLDSigPath_COUNTER_SIGNATURE_TYPE = "http://uri.etsi.org/01903#CountersignedSignature"
)

// XPath constants for "http://www.w3.org/2000/09/xmldsig#".
//
// KAT-verified against a Java oracle run over dss-xml-common 6.5.RC1 + dss-alert 6.5.RC1
// (java.lang.reflect over XMLDSigPath.class.getDeclaredFields(), each field's
// getQueryString() printed) - see xmldsig_path_test.go, which asserts every value below
// against the same strings the oracle printed. The exact string produced is reproduced in
// the comment above each field.
var (
	// SIGNATURE_PATH = "./ds:Signature"
	XMLDSigPath_SIGNATURE_PATH = FromCurrentPosition(XMLDSigElement_SIGNATURE)

	// ALL_SIGNATURES_PATH = "//ds:Signature"
	XMLDSigPath_ALL_SIGNATURES_PATH = All(XMLDSigElement_SIGNATURE)

	// ----------------------- From ds:Signature

	// OBJECT_PATH = "./ds:Object"
	XMLDSigPath_OBJECT_PATH = FromCurrentPosition(XMLDSigElement_OBJECT)

	// MANIFEST_PATH = "./ds:Object/ds:Manifest"
	XMLDSigPath_MANIFEST_PATH = FromCurrentPosition(XMLDSigElement_OBJECT, XMLDSigElement_MANIFEST)

	// SIGNED_INFO_PATH = "./ds:SignedInfo"
	XMLDSigPath_SIGNED_INFO_PATH = FromCurrentPosition(XMLDSigElement_SIGNED_INFO)

	// SIGNED_INFO_CANONICALIZATION_METHOD = "./ds:SignedInfo/ds:CanonicalizationMethod"
	XMLDSigPath_SIGNED_INFO_CANONICALIZATION_METHOD = FromCurrentPosition(XMLDSigElement_SIGNED_INFO, XMLDSigElement_CANONICALIZATION_METHOD)

	// SIGNED_INFO_REFERENCE_PATH = "./ds:SignedInfo/ds:Reference"
	XMLDSigPath_SIGNED_INFO_REFERENCE_PATH = FromCurrentPosition(XMLDSigElement_SIGNED_INFO, XMLDSigElement_REFERENCE)

	// SIGNATURE_METHOD_PATH = "./ds:SignedInfo/ds:SignatureMethod"
	XMLDSigPath_SIGNATURE_METHOD_PATH = FromCurrentPosition(XMLDSigElement_SIGNED_INFO, XMLDSigElement_SIGNATURE_METHOD)

	// REFERENCE_PATH = "./ds:Reference"
	XMLDSigPath_REFERENCE_PATH = FromCurrentPosition(XMLDSigElement_REFERENCE)

	// SIGNATURE_VALUE_PATH = "./ds:SignatureValue"
	XMLDSigPath_SIGNATURE_VALUE_PATH = FromCurrentPosition(XMLDSigElement_SIGNATURE_VALUE)

	// SIGNATURE_VALUE_ID_PATH = "./ds:SignatureValue/@Id"
	XMLDSigPath_SIGNATURE_VALUE_ID_PATH = FromCurrentPositionAttribute(XMLDSigElement_SIGNATURE_VALUE, XMLDSigAttribute_ID)

	// ALL_SIGNATURE_VALUES_PATH = "//ds:SignatureValue"
	XMLDSigPath_ALL_SIGNATURE_VALUES_PATH = All(XMLDSigElement_SIGNATURE_VALUE)

	// KEY_INFO_PATH = "./ds:KeyInfo"
	XMLDSigPath_KEY_INFO_PATH = FromCurrentPosition(XMLDSigElement_KEY_INFO)

	// KEY_INFO_X509_DATA = "./ds:KeyInfo/ds:X509Data"
	XMLDSigPath_KEY_INFO_X509_DATA = FromCurrentPosition(XMLDSigElement_KEY_INFO, XMLDSigElement_X509_DATA)

	// KEY_INFO_X509_CERTIFICATE_PATH = "./ds:KeyInfo/ds:X509Data/ds:X509Certificate"
	XMLDSigPath_KEY_INFO_X509_CERTIFICATE_PATH = FromCurrentPosition(XMLDSigElement_KEY_INFO, XMLDSigElement_X509_DATA, XMLDSigElement_X509_CERTIFICATE)

	// SIGNATURE_PROPERTIES_PATH = "./ds:Object/ds:SignatureProperties"
	XMLDSigPath_SIGNATURE_PROPERTIES_PATH = FromCurrentPosition(XMLDSigElement_OBJECT, XMLDSigElement_SIGNATURE_PROPERTIES)

	// SIGNATURE_PROPERTY_PATH = "./ds:Object/ds:SignatureProperties/ds:SignatureProperty"
	XMLDSigPath_SIGNATURE_PROPERTY_PATH = FromCurrentPosition(XMLDSigElement_OBJECT, XMLDSigElement_SIGNATURE_PROPERTIES, XMLDSigElement_SIGNATURE_PROPERTY)

	// ----------------------- For digest

	// DIGEST_METHOD_ALGORITHM_PATH = "./ds:DigestMethod/@Algorithm"
	XMLDSigPath_DIGEST_METHOD_ALGORITHM_PATH = FromCurrentPositionAttribute(XMLDSigElement_DIGEST_METHOD, XMLDSigAttribute_ALGORITHM)

	// DIGEST_VALUE_PATH = "./ds:DigestValue"
	XMLDSigPath_DIGEST_VALUE_PATH = FromCurrentPosition(XMLDSigElement_DIGEST_VALUE)

	// ------------------------- Canonicalization

	// CANONICALIZATION_ALGORITHM_PATH = "./ds:CanonicalizationMethod/@Algorithm"
	XMLDSigPath_CANONICALIZATION_ALGORITHM_PATH = FromCurrentPositionAttribute(XMLDSigElement_CANONICALIZATION_METHOD, XMLDSigAttribute_ALGORITHM)

	// ------------------------- Transforms

	// TRANSFORM_PATH = "./ds:Transform"
	XMLDSigPath_TRANSFORM_PATH = FromCurrentPosition(XMLDSigElement_TRANSFORM)

	// TRANSFORMS_PATH = "./ds:Transforms"
	XMLDSigPath_TRANSFORMS_PATH = FromCurrentPosition(XMLDSigElement_TRANSFORMS)

	// TRANSFORMS_TRANSFORM_PATH = "./ds:Transforms/ds:Transform"
	XMLDSigPath_TRANSFORMS_TRANSFORM_PATH = FromCurrentPosition(XMLDSigElement_TRANSFORMS, XMLDSigElement_TRANSFORM)
)
