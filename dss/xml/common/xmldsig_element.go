// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigElement.java (DSS 6.5.RC1).
package common

// XMLDSigElement is an element defined in the XMLDSig schema. Ports the Java enum per
// PORTING.md's enum convention: a typed string whose value is the Java name(), with the
// wire tag name (which differs from the Go/Java constant name for most entries) held in a
// lookup table.
type XMLDSigElement string

// XMLDSigElement constants, one per XMLDSig schema element name.
const (
	XMLDSigElement_CANONICALIZATION_METHOD XMLDSigElement = "CANONICALIZATION_METHOD"
	XMLDSigElement_DIGEST_METHOD           XMLDSigElement = "DIGEST_METHOD"
	XMLDSigElement_DIGEST_VALUE            XMLDSigElement = "DIGEST_VALUE"
	XMLDSigElement_DSA_KEY_VALUE           XMLDSigElement = "DSA_KEY_VALUE"
	XMLDSigElement_EXPONENT                XMLDSigElement = "EXPONENT"
	XMLDSigElement_G                       XMLDSigElement = "G"
	XMLDSigElement_HMAC_OUTPUT_LENGTH      XMLDSigElement = "HMAC_OUTPUT_LENGTH"
	XMLDSigElement_J                       XMLDSigElement = "J"
	XMLDSigElement_KEY_INFO                XMLDSigElement = "KEY_INFO"
	XMLDSigElement_KEY_NAME                XMLDSigElement = "KEY_NAME"
	XMLDSigElement_KEY_VALUE               XMLDSigElement = "KEY_VALUE"
	XMLDSigElement_MANIFEST                XMLDSigElement = "MANIFEST"
	XMLDSigElement_MGMT_DATA               XMLDSigElement = "MGMT_DATA"
	XMLDSigElement_MODULUS                 XMLDSigElement = "MODULUS"
	XMLDSigElement_OBJECT                  XMLDSigElement = "OBJECT"
	XMLDSigElement_P                       XMLDSigElement = "P"
	XMLDSigElement_PGEN_COUNTER            XMLDSigElement = "PGEN_COUNTER"
	XMLDSigElement_PGP_DATA                XMLDSigElement = "PGP_DATA"
	XMLDSigElement_PGP_KEY_ID              XMLDSigElement = "PGP_KEY_ID"
	XMLDSigElement_PGP_KEY_PACKET          XMLDSigElement = "PGP_KEY_PACKET"
	XMLDSigElement_Q                       XMLDSigElement = "Q"
	XMLDSigElement_REFERENCE               XMLDSigElement = "REFERENCE"
	XMLDSigElement_RETRIEVAL_METHOD        XMLDSigElement = "RETRIEVAL_METHOD"
	XMLDSigElement_RSA_KEY_VALUE           XMLDSigElement = "RSA_KEY_VALUE"
	XMLDSigElement_SEED                    XMLDSigElement = "SEED"
	XMLDSigElement_SIGNATURE               XMLDSigElement = "SIGNATURE"
	XMLDSigElement_SIGNATURE_METHOD        XMLDSigElement = "SIGNATURE_METHOD"
	XMLDSigElement_SIGNATURE_PROPERTIES    XMLDSigElement = "SIGNATURE_PROPERTIES"
	XMLDSigElement_SIGNATURE_PROPERTY      XMLDSigElement = "SIGNATURE_PROPERTY"
	XMLDSigElement_SIGNATURE_VALUE         XMLDSigElement = "SIGNATURE_VALUE"
	XMLDSigElement_SIGNED_INFO             XMLDSigElement = "SIGNED_INFO"
	XMLDSigElement_SPKI_DATA               XMLDSigElement = "SPKI_DATA"
	XMLDSigElement_SPKI_SEXP               XMLDSigElement = "SPKI_SEXP"
	XMLDSigElement_TRANSFORM               XMLDSigElement = "TRANSFORM"
	XMLDSigElement_TRANSFORMS              XMLDSigElement = "TRANSFORMS"
	XMLDSigElement_X509_CERTIFICATE        XMLDSigElement = "X509_CERTIFICATE"
	XMLDSigElement_X509_CRL                XMLDSigElement = "X509_CRL"
	XMLDSigElement_X509_DATA               XMLDSigElement = "X509_DATA"
	XMLDSigElement_X509_ISSUER_NAME        XMLDSigElement = "X509_ISSUER_NAME"
	XMLDSigElement_X509_ISSUER_SERIAL      XMLDSigElement = "X509_ISSUER_SERIAL"
	XMLDSigElement_X509_SERIAL_NUMBER      XMLDSigElement = "X509_SERIAL_NUMBER"
	XMLDSigElement_X509_SKI                XMLDSigElement = "X509_SKI"
	XMLDSigElement_X509_SUBJECT_NAME       XMLDSigElement = "X509_SUBJECT_NAME"
	XMLDSigElement_XPATH                   XMLDSigElement = "XPATH"
	XMLDSigElement_Y                       XMLDSigElement = "Y"
)

// xmldsigElementTagNames maps each constant to its wire tag name (getTagName()).
var xmldsigElementTagNames = map[XMLDSigElement]string{
	XMLDSigElement_CANONICALIZATION_METHOD: "CanonicalizationMethod",
	XMLDSigElement_DIGEST_METHOD:           "DigestMethod",
	XMLDSigElement_DIGEST_VALUE:            "DigestValue",
	XMLDSigElement_DSA_KEY_VALUE:           "DSAKeyValue",
	XMLDSigElement_EXPONENT:                "Exponent",
	XMLDSigElement_G:                       "G",
	XMLDSigElement_HMAC_OUTPUT_LENGTH:      "HMACOutputLength",
	XMLDSigElement_J:                       "J",
	XMLDSigElement_KEY_INFO:                "KeyInfo",
	XMLDSigElement_KEY_NAME:                "KeyName",
	XMLDSigElement_KEY_VALUE:               "KeyValue",
	XMLDSigElement_MANIFEST:                "Manifest",
	XMLDSigElement_MGMT_DATA:               "MgmtData",
	XMLDSigElement_MODULUS:                 "Modulus",
	XMLDSigElement_OBJECT:                  "Object",
	XMLDSigElement_P:                       "P",
	XMLDSigElement_PGEN_COUNTER:            "PgenCounter",
	XMLDSigElement_PGP_DATA:                "PGPData",
	XMLDSigElement_PGP_KEY_ID:              "PGPKeyID",
	XMLDSigElement_PGP_KEY_PACKET:          "PGPKeyPacket",
	XMLDSigElement_Q:                       "Q",
	XMLDSigElement_REFERENCE:               "Reference",
	XMLDSigElement_RETRIEVAL_METHOD:        "RetrievalMethod",
	XMLDSigElement_RSA_KEY_VALUE:           "RSAKeyValue",
	XMLDSigElement_SEED:                    "Seed",
	XMLDSigElement_SIGNATURE:               "Signature",
	XMLDSigElement_SIGNATURE_METHOD:        "SignatureMethod",
	XMLDSigElement_SIGNATURE_PROPERTIES:    "SignatureProperties",
	XMLDSigElement_SIGNATURE_PROPERTY:      "SignatureProperty",
	XMLDSigElement_SIGNATURE_VALUE:         "SignatureValue",
	XMLDSigElement_SIGNED_INFO:             "SignedInfo",
	XMLDSigElement_SPKI_DATA:               "SPKIData",
	XMLDSigElement_SPKI_SEXP:               "SPKISexp",
	XMLDSigElement_TRANSFORM:               "Transform",
	XMLDSigElement_TRANSFORMS:              "Transforms",
	XMLDSigElement_X509_CERTIFICATE:        "X509Certificate",
	XMLDSigElement_X509_CRL:                "X509CRL",
	XMLDSigElement_X509_DATA:               "X509Data",
	XMLDSigElement_X509_ISSUER_NAME:        "X509IssuerName",
	XMLDSigElement_X509_ISSUER_SERIAL:      "X509IssuerSerial",
	XMLDSigElement_X509_SERIAL_NUMBER:      "X509SerialNumber",
	XMLDSigElement_X509_SKI:                "X509SKI",
	XMLDSigElement_X509_SUBJECT_NAME:       "X509SubjectName",
	XMLDSigElement_XPATH:                   "XPath",
	XMLDSigElement_Y:                       "Y",
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
