// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigElement.java (DSS 6.5.RC1).
package common

// XMLDSigElement is an element defined in the XMLDSig schema. Ports the Java enum per
// PORTING.md's enum convention: a typed string whose value is the Java name(), with the
// wire tag name (which differs from the Go/Java constant name for most entries) held in a
// lookup table.
type XMLDSigElement string

// XMLDSigElement constants, one per XMLDSig schema element name.
const (
	XMLDSigElementCanonicalizationMethod XMLDSigElement = "CANONICALIZATION_METHOD"
	XMLDSigElementDigestMethod           XMLDSigElement = "DIGEST_METHOD"
	XMLDSigElementDigestValue            XMLDSigElement = "DIGEST_VALUE"
	XMLDSigElementDSAKeyValue            XMLDSigElement = "DSA_KEY_VALUE"
	XMLDSigElementExponent               XMLDSigElement = "EXPONENT"
	XMLDSigElementG                      XMLDSigElement = "G"
	XMLDSigElementHMACOutputLength       XMLDSigElement = "HMAC_OUTPUT_LENGTH"
	XMLDSigElementJ                      XMLDSigElement = "J"
	XMLDSigElementKeyInfo                XMLDSigElement = "KEY_INFO"
	XMLDSigElementKeyName                XMLDSigElement = "KEY_NAME"
	XMLDSigElementKeyValue               XMLDSigElement = "KEY_VALUE"
	XMLDSigElementManifest               XMLDSigElement = "MANIFEST"
	XMLDSigElementMgmtData               XMLDSigElement = "MGMT_DATA"
	XMLDSigElementModulus                XMLDSigElement = "MODULUS"
	XMLDSigElementObject                 XMLDSigElement = "OBJECT"
	XMLDSigElementP                      XMLDSigElement = "P"
	XMLDSigElementPgenCounter            XMLDSigElement = "PGEN_COUNTER"
	XMLDSigElementPGPData                XMLDSigElement = "PGP_DATA"
	XMLDSigElementPGPKeyID               XMLDSigElement = "PGP_KEY_ID"
	XMLDSigElementPGPKeyPacket           XMLDSigElement = "PGP_KEY_PACKET"
	XMLDSigElementQ                      XMLDSigElement = "Q"
	XMLDSigElementReference              XMLDSigElement = "REFERENCE"
	XMLDSigElementRetrievalMethod        XMLDSigElement = "RETRIEVAL_METHOD"
	XMLDSigElementRSAKeyValue            XMLDSigElement = "RSA_KEY_VALUE"
	XMLDSigElementSeed                   XMLDSigElement = "SEED"
	XMLDSigElementSignature              XMLDSigElement = "SIGNATURE"
	XMLDSigElementSignatureMethod        XMLDSigElement = "SIGNATURE_METHOD"
	XMLDSigElementSignatureProperties    XMLDSigElement = "SIGNATURE_PROPERTIES"
	XMLDSigElementSignatureProperty      XMLDSigElement = "SIGNATURE_PROPERTY"
	XMLDSigElementSignatureValue         XMLDSigElement = "SIGNATURE_VALUE"
	XMLDSigElementSignedInfo             XMLDSigElement = "SIGNED_INFO"
	XMLDSigElementSPKIData               XMLDSigElement = "SPKI_DATA"
	XMLDSigElementSPKISexp               XMLDSigElement = "SPKI_SEXP"
	XMLDSigElementTransform              XMLDSigElement = "TRANSFORM"
	XMLDSigElementTransforms             XMLDSigElement = "TRANSFORMS"
	XMLDSigElementX509Certificate        XMLDSigElement = "X509_CERTIFICATE"
	XMLDSigElementX509CRL                XMLDSigElement = "X509_CRL"
	XMLDSigElementX509Data               XMLDSigElement = "X509_DATA"
	XMLDSigElementX509IssuerName         XMLDSigElement = "X509_ISSUER_NAME"
	XMLDSigElementX509IssuerSerial       XMLDSigElement = "X509_ISSUER_SERIAL"
	XMLDSigElementX509SerialNumber       XMLDSigElement = "X509_SERIAL_NUMBER"
	XMLDSigElementX509SKI                XMLDSigElement = "X509_SKI"
	XMLDSigElementX509SubjectName        XMLDSigElement = "X509_SUBJECT_NAME"
	XMLDSigElementXPATH                  XMLDSigElement = "XPATH"
	XMLDSigElementY                      XMLDSigElement = "Y"
)

// xmldsigElementTagNames maps each constant to its wire tag name (getTagName()).
var xmldsigElementTagNames = map[XMLDSigElement]string{
	XMLDSigElementCanonicalizationMethod: "CanonicalizationMethod",
	XMLDSigElementDigestMethod:           "DigestMethod",
	XMLDSigElementDigestValue:            "DigestValue",
	XMLDSigElementDSAKeyValue:            "DSAKeyValue",
	XMLDSigElementExponent:               "Exponent",
	XMLDSigElementG:                      "G",
	XMLDSigElementHMACOutputLength:       "HMACOutputLength",
	XMLDSigElementJ:                      "J",
	XMLDSigElementKeyInfo:                "KeyInfo",
	XMLDSigElementKeyName:                "KeyName",
	XMLDSigElementKeyValue:               "KeyValue",
	XMLDSigElementManifest:               "Manifest",
	XMLDSigElementMgmtData:               "MgmtData",
	XMLDSigElementModulus:                "Modulus",
	XMLDSigElementObject:                 "Object",
	XMLDSigElementP:                      "P",
	XMLDSigElementPgenCounter:            "PgenCounter",
	XMLDSigElementPGPData:                "PGPData",
	XMLDSigElementPGPKeyID:               "PGPKeyID",
	XMLDSigElementPGPKeyPacket:           "PGPKeyPacket",
	XMLDSigElementQ:                      "Q",
	XMLDSigElementReference:              "Reference",
	XMLDSigElementRetrievalMethod:        "RetrievalMethod",
	XMLDSigElementRSAKeyValue:            "RSAKeyValue",
	XMLDSigElementSeed:                   "Seed",
	XMLDSigElementSignature:              "Signature",
	XMLDSigElementSignatureMethod:        "SignatureMethod",
	XMLDSigElementSignatureProperties:    "SignatureProperties",
	XMLDSigElementSignatureProperty:      "SignatureProperty",
	XMLDSigElementSignatureValue:         "SignatureValue",
	XMLDSigElementSignedInfo:             "SignedInfo",
	XMLDSigElementSPKIData:               "SPKIData",
	XMLDSigElementSPKISexp:               "SPKISexp",
	XMLDSigElementTransform:              "Transform",
	XMLDSigElementTransforms:             "Transforms",
	XMLDSigElementX509Certificate:        "X509Certificate",
	XMLDSigElementX509CRL:                "X509CRL",
	XMLDSigElementX509Data:               "X509Data",
	XMLDSigElementX509IssuerName:         "X509IssuerName",
	XMLDSigElementX509IssuerSerial:       "X509IssuerSerial",
	XMLDSigElementX509SerialNumber:       "X509SerialNumber",
	XMLDSigElementX509SKI:                "X509SKI",
	XMLDSigElementX509SubjectName:        "X509SubjectName",
	XMLDSigElementXPATH:                  "XPath",
	XMLDSigElementY:                      "Y",
}

// TagName implements DSSElement. Ports getTagName().
func (e XMLDSigElement) TagName() string {
	return xmldsigElementTagNames[e]
}

// Namespace implements DSSElement. Every XMLDSigElement carries XMLDSigNS. Ports
// getNamespace().
func (e XMLDSigElement) Namespace() *DSSNamespace {
	return XMLDSigNS
}

// URI implements DSSElement. Ports getURI().
func (e XMLDSigElement) URI() string {
	return XMLDSigNS.Uri()
}

// IsSameTagName implements DSSElement. Ports isSameTagName(String).
func (e XMLDSigElement) IsSameTagName(value string) bool {
	return e.TagName() == value
}

var _ DSSElement = XMLDSigElement("")
