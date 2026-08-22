// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigPath.java (DSS 6.5.RC1).
package common

// The "Type" attribute values a ds:Reference element may carry.
const (
	// XMLDSigPathObjectType is the "Type" attribute value for a ds:Reference element
	// referring to a ds:Object element.
	XMLDSigPathObjectType = "http://www.w3.org/2000/09/xmldsig#Object"

	// XMLDSigPathManifestType is the "Type" attribute value for a ds:Reference element
	// referring to a signed manifest.
	XMLDSigPathManifestType = "http://www.w3.org/2000/09/xmldsig#Manifest"

	// XMLDSigPathCounterSignatureType is the "Type" attribute value for a ds:Reference
	// element pointing at a ds:SignatureValue of a counter-signed signature.
	XMLDSigPathCounterSignatureType = "http://uri.etsi.org/01903#CountersignedSignature"
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
	XMLDSigPathSignaturePath = FromCurrentPosition(XMLDSigElementSignature)

	// ALL_SIGNATURES_PATH = "//ds:Signature"
	XMLDSigPathAllSignaturesPath = All(XMLDSigElementSignature)

	// ----------------------- From ds:Signature

	// OBJECT_PATH = "./ds:Object"
	XMLDSigPathObjectPath = FromCurrentPosition(XMLDSigElementObject)

	// MANIFEST_PATH = "./ds:Object/ds:Manifest"
	XMLDSigPathManifestPath = FromCurrentPosition(XMLDSigElementObject, XMLDSigElementManifest)

	// SIGNED_INFO_PATH = "./ds:SignedInfo"
	XMLDSigPathSignedInfoPath = FromCurrentPosition(XMLDSigElementSignedInfo)

	// SIGNED_INFO_CANONICALIZATION_METHOD = "./ds:SignedInfo/ds:CanonicalizationMethod"
	XMLDSigPathSignedInfoCanonicalizationMethod = FromCurrentPosition(XMLDSigElementSignedInfo, XMLDSigElementCanonicalizationMethod)

	// SIGNED_INFO_REFERENCE_PATH = "./ds:SignedInfo/ds:Reference"
	XMLDSigPathSignedInfoReferencePath = FromCurrentPosition(XMLDSigElementSignedInfo, XMLDSigElementReference)

	// SIGNATURE_METHOD_PATH = "./ds:SignedInfo/ds:SignatureMethod"
	XMLDSigPathSignatureMethodPath = FromCurrentPosition(XMLDSigElementSignedInfo, XMLDSigElementSignatureMethod)

	// REFERENCE_PATH = "./ds:Reference"
	XMLDSigPathReferencePath = FromCurrentPosition(XMLDSigElementReference)

	// SIGNATURE_VALUE_PATH = "./ds:SignatureValue"
	XMLDSigPathSignatureValuePath = FromCurrentPosition(XMLDSigElementSignatureValue)

	// SIGNATURE_VALUE_ID_PATH = "./ds:SignatureValue/@Id"
	XMLDSigPathSignatureValueIDPath = FromCurrentPositionAttribute(XMLDSigElementSignatureValue, XMLDSigAttributeID)

	// ALL_SIGNATURE_VALUES_PATH = "//ds:SignatureValue"
	XMLDSigPathAllSignatureValuesPath = All(XMLDSigElementSignatureValue)

	// KEY_INFO_PATH = "./ds:KeyInfo"
	XMLDSigPathKeyInfoPath = FromCurrentPosition(XMLDSigElementKeyInfo)

	// KEY_INFO_X509_DATA = "./ds:KeyInfo/ds:X509Data"
	XMLDSigPathKeyInfoX509Data = FromCurrentPosition(XMLDSigElementKeyInfo, XMLDSigElementX509Data)

	// KEY_INFO_X509_CERTIFICATE_PATH = "./ds:KeyInfo/ds:X509Data/ds:X509Certificate"
	XMLDSigPathKeyInfoX509CertificatePath = FromCurrentPosition(XMLDSigElementKeyInfo, XMLDSigElementX509Data, XMLDSigElementX509Certificate)

	// SIGNATURE_PROPERTIES_PATH = "./ds:Object/ds:SignatureProperties"
	XMLDSigPathSignaturePropertiesPath = FromCurrentPosition(XMLDSigElementObject, XMLDSigElementSignatureProperties)

	// SIGNATURE_PROPERTY_PATH = "./ds:Object/ds:SignatureProperties/ds:SignatureProperty"
	XMLDSigPathSignaturePropertyPath = FromCurrentPosition(XMLDSigElementObject, XMLDSigElementSignatureProperties, XMLDSigElementSignatureProperty)

	// ----------------------- For digest

	// DIGEST_METHOD_ALGORITHM_PATH = "./ds:DigestMethod/@Algorithm"
	XMLDSigPathDigestMethodAlgorithmPath = FromCurrentPositionAttribute(XMLDSigElementDigestMethod, XMLDSigAttributeAlgorithm)

	// DIGEST_VALUE_PATH = "./ds:DigestValue"
	XMLDSigPathDigestValuePath = FromCurrentPosition(XMLDSigElementDigestValue)

	// ------------------------- Canonicalization

	// CANONICALIZATION_ALGORITHM_PATH = "./ds:CanonicalizationMethod/@Algorithm"
	XMLDSigPathCanonicalizationAlgorithmPath = FromCurrentPositionAttribute(XMLDSigElementCanonicalizationMethod, XMLDSigAttributeAlgorithm)

	// ------------------------- Transforms

	// TRANSFORM_PATH = "./ds:Transform"
	XMLDSigPathTransformPath = FromCurrentPosition(XMLDSigElementTransform)

	// TRANSFORMS_PATH = "./ds:Transforms"
	XMLDSigPathTransformsPath = FromCurrentPosition(XMLDSigElementTransforms)

	// TRANSFORMS_TRANSFORM_PATH = "./ds:Transforms/ds:Transform"
	XMLDSigPathTransformsTransformPath = FromCurrentPosition(XMLDSigElementTransforms, XMLDSigElementTransform)
)
