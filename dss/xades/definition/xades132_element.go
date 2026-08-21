// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades132/XAdES132Element.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES132Element defines elements for a XAdES 1.3.2 schema.
type XAdES132Element string

// XAdES132Element constants, one per XAdES 1.3.2 schema element name.
const (
	XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP        XAdES132Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES132Element_ALL_SIGNED_DATA_OBJECTS           XAdES132Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES132Element_ANY                               XAdES132Element = "ANY"
	XAdES132Element_ARCHIVE_TIMESTAMP                 XAdES132Element = "ARCHIVE_TIMESTAMP"
	XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES      XAdES132Element = "ATTR_AUTHORITIES_CERT_VALUES"
	XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS        XAdES132Element = "ATTRIBUTE_CERTIFICATE_REFS"
	XAdES132Element_ATTRIBUTE_REVOCATION_REFS         XAdES132Element = "ATTRIBUTE_REVOCATION_REFS"
	XAdES132Element_ATTRIBUTE_REVOCATION_VALUES       XAdES132Element = "ATTRIBUTE_REVOCATION_VALUES"
	XAdES132Element_BY_KEY                            XAdES132Element = "BY_KEY"
	XAdES132Element_BY_NAME                           XAdES132Element = "BY_NAME"
	XAdES132Element_CERT                              XAdES132Element = "CERT"
	XAdES132Element_CERT_DIGEST                       XAdES132Element = "CERT_DIGEST"
	XAdES132Element_CERT_REFS                         XAdES132Element = "CERT_REFS"
	XAdES132Element_CERTIFICATE_VALUES                XAdES132Element = "CERTIFICATE_VALUES"
	XAdES132Element_CERTIFIED_ROLE                    XAdES132Element = "CERTIFIED_ROLE"
	XAdES132Element_CERTIFIED_ROLES                   XAdES132Element = "CERTIFIED_ROLES"
	XAdES132Element_CERTIFIED_ROLES_V2                XAdES132Element = "CERTIFIED_ROLES_V2"
	XAdES132Element_CITY                              XAdES132Element = "CITY"
	XAdES132Element_CLAIMED_ROLE                      XAdES132Element = "CLAIMED_ROLE"
	XAdES132Element_CLAIMED_ROLES                     XAdES132Element = "CLAIMED_ROLES"
	XAdES132Element_COMMITMENT_TYPE_ID                XAdES132Element = "COMMITMENT_TYPE_ID"
	XAdES132Element_COMMITMENT_TYPE_INDICATION        XAdES132Element = "COMMITMENT_TYPE_INDICATION"
	XAdES132Element_COMMITMENT_TYPE_QUALIFIER         XAdES132Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES132Element_COMMITMENT_TYPE_QUALIFIERS        XAdES132Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES132Element_COMPLETE_CERTIFICATE_REFS         XAdES132Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES132Element_COMPLETE_REVOCATION_REFS          XAdES132Element = "COMPLETE_REVOCATION_REFS"
	XAdES132Element_COUNTER_SIGNATURE                 XAdES132Element = "COUNTER_SIGNATURE"
	XAdES132Element_COUNTRY_NAME                      XAdES132Element = "COUNTRY_NAME"
	XAdES132Element_CRL_IDENTIFIER                    XAdES132Element = "CRL_IDENTIFIER"
	XAdES132Element_CRL_REF                           XAdES132Element = "CRL_REF"
	XAdES132Element_CRL_REFS                          XAdES132Element = "CRL_REFS"
	XAdES132Element_CRL_VALUES                        XAdES132Element = "CRL_VALUES"
	XAdES132Element_DATA_OBJECT_FORMAT                XAdES132Element = "DATA_OBJECT_FORMAT"
	XAdES132Element_DESCRIPTION                       XAdES132Element = "DESCRIPTION"
	XAdES132Element_DIGEST_ALG_AND_VALUE              XAdES132Element = "DIGEST_ALG_AND_VALUE"
	XAdES132Element_DOCUMENTATION_REFERENCE           XAdES132Element = "DOCUMENTATION_REFERENCE"
	XAdES132Element_DOCUMENTATION_REFERENCES          XAdES132Element = "DOCUMENTATION_REFERENCES"
	XAdES132Element_ENCAPSULATED_CRL_VALUE            XAdES132Element = "ENCAPSULATED_CRL_VALUE"
	XAdES132Element_ENCAPSULATED_OCSP_VALUE           XAdES132Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES132Element_ENCAPSULATED_PKI_DATA             XAdES132Element = "ENCAPSULATED_PKI_DATA"
	XAdES132Element_ENCAPSULATED_TIMESTAMP            XAdES132Element = "ENCAPSULATED_TIMESTAMP"
	XAdES132Element_ENCAPSULATED_X509_CERTIFICATE     XAdES132Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES132Element_ENCODING                          XAdES132Element = "ENCODING"
	XAdES132Element_EXPLICIT_TEXT                     XAdES132Element = "EXPLICIT_TEXT"
	XAdES132Element_IDENTIFIER                        XAdES132Element = "IDENTIFIER"
	XAdES132Element_INCLUDE                           XAdES132Element = "INCLUDE"
	XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP XAdES132Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES132Element_INT                               XAdES132Element = "INT"
	XAdES132Element_ISSUE_TIME                        XAdES132Element = "ISSUE_TIME"
	XAdES132Element_ISSUER                            XAdES132Element = "ISSUER"
	XAdES132Element_ISSUER_SERIAL                     XAdES132Element = "ISSUER_SERIAL"
	XAdES132Element_ISSUER_SERIAL_V2                  XAdES132Element = "ISSUER_SERIAL_V2"
	XAdES132Element_MIME_TYPE                         XAdES132Element = "MIME_TYPE"
	XAdES132Element_NOTICE_NUMBERS                    XAdES132Element = "NOTICE_NUMBERS"
	XAdES132Element_NOTICE_REF                        XAdES132Element = "NOTICE_REF"
	XAdES132Element_NUMBER                            XAdES132Element = "NUMBER"
	XAdES132Element_OBJECT_IDENTIFIER                 XAdES132Element = "OBJECT_IDENTIFIER"
	XAdES132Element_OBJECT_REFERENCE                  XAdES132Element = "OBJECT_REFERENCE"
	XAdES132Element_OCSP_IDENTIFIER                   XAdES132Element = "OCSP_IDENTIFIER"
	XAdES132Element_OCSP_REF                          XAdES132Element = "OCSP_REF"
	XAdES132Element_OCSP_REFS                         XAdES132Element = "OCSP_REFS"
	XAdES132Element_OCSP_VALUES                       XAdES132Element = "OCSP_VALUES"
	XAdES132Element_ORGANIZATION                      XAdES132Element = "ORGANIZATION"
	XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE       XAdES132Element = "OTHER_ATTRIBUTE_CERTIFICATE"
	XAdES132Element_OTHER_CERTIFICATE                 XAdES132Element = "OTHER_CERTIFICATE"
	XAdES132Element_OTHER_REF                         XAdES132Element = "OTHER_REF"
	XAdES132Element_OTHER_REFS                        XAdES132Element = "OTHER_REFS"
	XAdES132Element_OTHER_TIMESTAMP                   XAdES132Element = "OTHER_TIMESTAMP"
	XAdES132Element_OTHER_VALUE                       XAdES132Element = "OTHER_VALUE"
	XAdES132Element_OTHER_VALUES                      XAdES132Element = "OTHER_VALUES"
	XAdES132Element_POSTAL_CODE                       XAdES132Element = "POSTAL_CODE"
	XAdES132Element_PRODUCED_AT                       XAdES132Element = "PRODUCED_AT"
	XAdES132Element_QUALIFYING_PROPERTIES             XAdES132Element = "QUALIFYING_PROPERTIES"
	XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE   XAdES132Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES132Element_REFERENCE_INFO                    XAdES132Element = "REFERENCE_INFO"
	XAdES132Element_REFS_ONLY_TIMESTAMP               XAdES132Element = "REFS_ONLY_TIMESTAMP"
	XAdES132Element_RESPONDER_ID                      XAdES132Element = "RESPONDER_ID"
	XAdES132Element_REVOCATION_VALUES                 XAdES132Element = "REVOCATION_VALUES"
	XAdES132Element_SIG_AND_REFS_TIMESTAMP            XAdES132Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES132Element_SIG_POLICY_HASH                   XAdES132Element = "SIG_POLICY_HASH"
	XAdES132Element_SIG_POLICY_ID                     XAdES132Element = "SIG_POLICY_ID"
	XAdES132Element_SIG_POLICY_QUALIFIER              XAdES132Element = "SIG_POLICY_QUALIFIER"
	XAdES132Element_SIG_POLICY_QUALIFIERS             XAdES132Element = "SIG_POLICY_QUALIFIERS"
	XAdES132Element_SIGNATURE_POLICY_ID               XAdES132Element = "SIGNATURE_POLICY_ID"
	XAdES132Element_SIGNATURE_POLICY_IDENTIFIER       XAdES132Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES132Element_SIGNATURE_POLICY_IMPLIED          XAdES132Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES132Element_SIGNATURE_PRODUCTION_PLACE        XAdES132Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2     XAdES132Element = "SIGNATURE_PRODUCTION_PLACE_V2"
	XAdES132Element_SIGNATURE_TIMESTAMP               XAdES132Element = "SIGNATURE_TIMESTAMP"
	XAdES132Element_SIGNED_ASSERTION                  XAdES132Element = "SIGNED_ASSERTION"
	XAdES132Element_SIGNED_ASSERTIONS                 XAdES132Element = "SIGNED_ASSERTIONS"
	XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES     XAdES132Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES132Element_SIGNED_PROPERTIES                 XAdES132Element = "SIGNED_PROPERTIES"
	XAdES132Element_SIGNED_SIGNATURE_PROPERTIES       XAdES132Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES132Element_SIGNER_ROLE                       XAdES132Element = "SIGNER_ROLE"
	XAdES132Element_SIGNER_ROLE_V2                    XAdES132Element = "SIGNER_ROLE_V2"
	XAdES132Element_SIGNING_CERTIFICATE               XAdES132Element = "SIGNING_CERTIFICATE"
	XAdES132Element_SIGNING_CERTIFICATE_V2            XAdES132Element = "SIGNING_CERTIFICATE_V2"
	XAdES132Element_SIGNING_TIME                      XAdES132Element = "SIGNING_TIME"
	XAdES132Element_SP_URI                            XAdES132Element = "SP_URI"
	XAdES132Element_SP_USER_NOTICE                    XAdES132Element = "SP_USER_NOTICE"
	XAdES132Element_STATE_OR_PROVINCE                 XAdES132Element = "STATE_OR_PROVINCE"
	XAdES132Element_STREET_ADDRESS                    XAdES132Element = "STREET_ADDRESS"
	XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES   XAdES132Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY     XAdES132Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES132Element_UNSIGNED_PROPERTIES               XAdES132Element = "UNSIGNED_PROPERTIES"
	XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES     XAdES132Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES132Element_X509_ATTRIBUTE_CERTIFICATE        XAdES132Element = "X509_ATTRIBUTE_CERTIFICATE"
	XAdES132Element_XADES_TIMESTAMP                   XAdES132Element = "XADES_TIMESTAMP"
	XAdES132Element_XML_TIMESTAMP                     XAdES132Element = "XML_TIMESTAMP"
)

// xades132elementTagNames maps each constant to its wire tag name (getTagName()).
var xades132elementTagNames = map[XAdES132Element]string{
	XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP:        "AllDataObjectsTimeStamp",
	XAdES132Element_ALL_SIGNED_DATA_OBJECTS:           "AllSignedDataObjects",
	XAdES132Element_ANY:                               "Any",
	XAdES132Element_ARCHIVE_TIMESTAMP:                 "ArchiveTimeStamp",
	XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES:      "AttrAuthoritiesCertValues",
	XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS:        "AttributeCertificateRefs",
	XAdES132Element_ATTRIBUTE_REVOCATION_REFS:         "AttributeRevocationRefs",
	XAdES132Element_ATTRIBUTE_REVOCATION_VALUES:       "AttributeRevocationValues",
	XAdES132Element_BY_KEY:                            "ByKey",
	XAdES132Element_BY_NAME:                           "ByName",
	XAdES132Element_CERT:                              "Cert",
	XAdES132Element_CERT_DIGEST:                       "CertDigest",
	XAdES132Element_CERT_REFS:                         "CertRefs",
	XAdES132Element_CERTIFICATE_VALUES:                "CertificateValues",
	XAdES132Element_CERTIFIED_ROLE:                    "CertifiedRole",
	XAdES132Element_CERTIFIED_ROLES:                   "CertifiedRoles",
	XAdES132Element_CERTIFIED_ROLES_V2:                "CertifiedRolesV2",
	XAdES132Element_CITY:                              "City",
	XAdES132Element_CLAIMED_ROLE:                      "ClaimedRole",
	XAdES132Element_CLAIMED_ROLES:                     "ClaimedRoles",
	XAdES132Element_COMMITMENT_TYPE_ID:                "CommitmentTypeId",
	XAdES132Element_COMMITMENT_TYPE_INDICATION:        "CommitmentTypeIndication",
	XAdES132Element_COMMITMENT_TYPE_QUALIFIER:         "CommitmentTypeQualifier",
	XAdES132Element_COMMITMENT_TYPE_QUALIFIERS:        "CommitmentTypeQualifiers",
	XAdES132Element_COMPLETE_CERTIFICATE_REFS:         "CompleteCertificateRefs",
	XAdES132Element_COMPLETE_REVOCATION_REFS:          "CompleteRevocationRefs",
	XAdES132Element_COUNTER_SIGNATURE:                 "CounterSignature",
	XAdES132Element_COUNTRY_NAME:                      "CountryName",
	XAdES132Element_CRL_IDENTIFIER:                    "CRLIdentifier",
	XAdES132Element_CRL_REF:                           "CRLRef",
	XAdES132Element_CRL_REFS:                          "CRLRefs",
	XAdES132Element_CRL_VALUES:                        "CRLValues",
	XAdES132Element_DATA_OBJECT_FORMAT:                "DataObjectFormat",
	XAdES132Element_DESCRIPTION:                       "Description",
	XAdES132Element_DIGEST_ALG_AND_VALUE:              "DigestAlgAndValue",
	XAdES132Element_DOCUMENTATION_REFERENCE:           "DocumentationReference",
	XAdES132Element_DOCUMENTATION_REFERENCES:          "DocumentationReferences",
	XAdES132Element_ENCAPSULATED_CRL_VALUE:            "EncapsulatedCRLValue",
	XAdES132Element_ENCAPSULATED_OCSP_VALUE:           "EncapsulatedOCSPValue",
	XAdES132Element_ENCAPSULATED_PKI_DATA:             "EncapsulatedPKIData",
	XAdES132Element_ENCAPSULATED_TIMESTAMP:            "EncapsulatedTimeStamp",
	XAdES132Element_ENCAPSULATED_X509_CERTIFICATE:     "EncapsulatedX509Certificate",
	XAdES132Element_ENCODING:                          "Encoding",
	XAdES132Element_EXPLICIT_TEXT:                     "ExplicitText",
	XAdES132Element_IDENTIFIER:                        "Identifier",
	XAdES132Element_INCLUDE:                           "Include",
	XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP: "IndividualDataObjectsTimeStamp",
	XAdES132Element_INT:                               "int",
	XAdES132Element_ISSUE_TIME:                        "IssueTime",
	XAdES132Element_ISSUER:                            "Issuer",
	XAdES132Element_ISSUER_SERIAL:                     "IssuerSerial",
	XAdES132Element_ISSUER_SERIAL_V2:                  "IssuerSerialV2",
	XAdES132Element_MIME_TYPE:                         "MimeType",
	XAdES132Element_NOTICE_NUMBERS:                    "NoticeNumbers",
	XAdES132Element_NOTICE_REF:                        "NoticeRef",
	XAdES132Element_NUMBER:                            "Number",
	XAdES132Element_OBJECT_IDENTIFIER:                 "ObjectIdentifier",
	XAdES132Element_OBJECT_REFERENCE:                  "ObjectReference",
	XAdES132Element_OCSP_IDENTIFIER:                   "OCSPIdentifier",
	XAdES132Element_OCSP_REF:                          "OCSPRef",
	XAdES132Element_OCSP_REFS:                         "OCSPRefs",
	XAdES132Element_OCSP_VALUES:                       "OCSPValues",
	XAdES132Element_ORGANIZATION:                      "Organization",
	XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE:       "OtherAttributeCertificate",
	XAdES132Element_OTHER_CERTIFICATE:                 "OtherCertificate",
	XAdES132Element_OTHER_REF:                         "OtherRef",
	XAdES132Element_OTHER_REFS:                        "OtherRefs",
	XAdES132Element_OTHER_TIMESTAMP:                   "OtherTimeStamp",
	XAdES132Element_OTHER_VALUE:                       "OtherValue",
	XAdES132Element_OTHER_VALUES:                      "OtherValues",
	XAdES132Element_POSTAL_CODE:                       "PostalCode",
	XAdES132Element_PRODUCED_AT:                       "ProducedAt",
	XAdES132Element_QUALIFYING_PROPERTIES:             "QualifyingProperties",
	XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE:   "QualifyingPropertiesReference",
	XAdES132Element_REFERENCE_INFO:                    "ReferenceInfo",
	XAdES132Element_REFS_ONLY_TIMESTAMP:               "RefsOnlyTimeStamp",
	XAdES132Element_RESPONDER_ID:                      "ResponderID",
	XAdES132Element_REVOCATION_VALUES:                 "RevocationValues",
	XAdES132Element_SIG_AND_REFS_TIMESTAMP:            "SigAndRefsTimeStamp",
	XAdES132Element_SIG_POLICY_HASH:                   "SigPolicyHash",
	XAdES132Element_SIG_POLICY_ID:                     "SigPolicyId",
	XAdES132Element_SIG_POLICY_QUALIFIER:              "SigPolicyQualifier",
	XAdES132Element_SIG_POLICY_QUALIFIERS:             "SigPolicyQualifiers",
	XAdES132Element_SIGNATURE_POLICY_ID:               "SignaturePolicyId",
	XAdES132Element_SIGNATURE_POLICY_IDENTIFIER:       "SignaturePolicyIdentifier",
	XAdES132Element_SIGNATURE_POLICY_IMPLIED:          "SignaturePolicyImplied",
	XAdES132Element_SIGNATURE_PRODUCTION_PLACE:        "SignatureProductionPlace",
	XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2:     "SignatureProductionPlaceV2",
	XAdES132Element_SIGNATURE_TIMESTAMP:               "SignatureTimeStamp",
	XAdES132Element_SIGNED_ASSERTION:                  "SignedAssertion",
	XAdES132Element_SIGNED_ASSERTIONS:                 "SignedAssertions",
	XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES:     "SignedDataObjectProperties",
	XAdES132Element_SIGNED_PROPERTIES:                 "SignedProperties",
	XAdES132Element_SIGNED_SIGNATURE_PROPERTIES:       "SignedSignatureProperties",
	XAdES132Element_SIGNER_ROLE:                       "SignerRole",
	XAdES132Element_SIGNER_ROLE_V2:                    "SignerRoleV2",
	XAdES132Element_SIGNING_CERTIFICATE:               "SigningCertificate",
	XAdES132Element_SIGNING_CERTIFICATE_V2:            "SigningCertificateV2",
	XAdES132Element_SIGNING_TIME:                      "SigningTime",
	XAdES132Element_SP_URI:                            "SPURI",
	XAdES132Element_SP_USER_NOTICE:                    "SPUserNotice",
	XAdES132Element_STATE_OR_PROVINCE:                 "StateOrProvince",
	XAdES132Element_STREET_ADDRESS:                    "StreetAddress",
	XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES:   "UnsignedDataObjectProperties",
	XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY:     "UnsignedDataObjectProperty",
	XAdES132Element_UNSIGNED_PROPERTIES:               "UnsignedProperties",
	XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES:     "UnsignedSignatureProperties",
	XAdES132Element_X509_ATTRIBUTE_CERTIFICATE:        "X509AttributeCertificate",
	XAdES132Element_XADES_TIMESTAMP:                   "XAdESTimeStamp",
	XAdES132Element_XML_TIMESTAMP:                     "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES132Element) TagName() string {
	return xades132elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES132Element) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_132
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES132Element) URI() string {
	return XAdESNamespace_XADES_132.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e XAdES132Element) IsSameTagName(value string) bool {
	return e.TagName() == value
}

// XAdES132ElementFromTagName returns the XAdES132Element constant with the given tag name, or nil
// if none matches. Ports the static factory fromTagName(String).
func XAdES132ElementFromTagName(tagName string) common.DSSElement {
	for _, e := range XAdES132ElementValues() {
		if e.TagName() == tagName {
			return e
		}
	}
	return nil
}

// XAdES132ElementValues returns every XAdES132Element constant, in declaration order. Ports values().
func XAdES132ElementValues() []XAdES132Element {
	return []XAdES132Element{
		XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP,
		XAdES132Element_ALL_SIGNED_DATA_OBJECTS,
		XAdES132Element_ANY,
		XAdES132Element_ARCHIVE_TIMESTAMP,
		XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES,
		XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS,
		XAdES132Element_ATTRIBUTE_REVOCATION_REFS,
		XAdES132Element_ATTRIBUTE_REVOCATION_VALUES,
		XAdES132Element_BY_KEY,
		XAdES132Element_BY_NAME,
		XAdES132Element_CERT,
		XAdES132Element_CERT_DIGEST,
		XAdES132Element_CERT_REFS,
		XAdES132Element_CERTIFICATE_VALUES,
		XAdES132Element_CERTIFIED_ROLE,
		XAdES132Element_CERTIFIED_ROLES,
		XAdES132Element_CERTIFIED_ROLES_V2,
		XAdES132Element_CITY,
		XAdES132Element_CLAIMED_ROLE,
		XAdES132Element_CLAIMED_ROLES,
		XAdES132Element_COMMITMENT_TYPE_ID,
		XAdES132Element_COMMITMENT_TYPE_INDICATION,
		XAdES132Element_COMMITMENT_TYPE_QUALIFIER,
		XAdES132Element_COMMITMENT_TYPE_QUALIFIERS,
		XAdES132Element_COMPLETE_CERTIFICATE_REFS,
		XAdES132Element_COMPLETE_REVOCATION_REFS,
		XAdES132Element_COUNTER_SIGNATURE,
		XAdES132Element_COUNTRY_NAME,
		XAdES132Element_CRL_IDENTIFIER,
		XAdES132Element_CRL_REF,
		XAdES132Element_CRL_REFS,
		XAdES132Element_CRL_VALUES,
		XAdES132Element_DATA_OBJECT_FORMAT,
		XAdES132Element_DESCRIPTION,
		XAdES132Element_DIGEST_ALG_AND_VALUE,
		XAdES132Element_DOCUMENTATION_REFERENCE,
		XAdES132Element_DOCUMENTATION_REFERENCES,
		XAdES132Element_ENCAPSULATED_CRL_VALUE,
		XAdES132Element_ENCAPSULATED_OCSP_VALUE,
		XAdES132Element_ENCAPSULATED_PKI_DATA,
		XAdES132Element_ENCAPSULATED_TIMESTAMP,
		XAdES132Element_ENCAPSULATED_X509_CERTIFICATE,
		XAdES132Element_ENCODING,
		XAdES132Element_EXPLICIT_TEXT,
		XAdES132Element_IDENTIFIER,
		XAdES132Element_INCLUDE,
		XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP,
		XAdES132Element_INT,
		XAdES132Element_ISSUE_TIME,
		XAdES132Element_ISSUER,
		XAdES132Element_ISSUER_SERIAL,
		XAdES132Element_ISSUER_SERIAL_V2,
		XAdES132Element_MIME_TYPE,
		XAdES132Element_NOTICE_NUMBERS,
		XAdES132Element_NOTICE_REF,
		XAdES132Element_NUMBER,
		XAdES132Element_OBJECT_IDENTIFIER,
		XAdES132Element_OBJECT_REFERENCE,
		XAdES132Element_OCSP_IDENTIFIER,
		XAdES132Element_OCSP_REF,
		XAdES132Element_OCSP_REFS,
		XAdES132Element_OCSP_VALUES,
		XAdES132Element_ORGANIZATION,
		XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE,
		XAdES132Element_OTHER_CERTIFICATE,
		XAdES132Element_OTHER_REF,
		XAdES132Element_OTHER_REFS,
		XAdES132Element_OTHER_TIMESTAMP,
		XAdES132Element_OTHER_VALUE,
		XAdES132Element_OTHER_VALUES,
		XAdES132Element_POSTAL_CODE,
		XAdES132Element_PRODUCED_AT,
		XAdES132Element_QUALIFYING_PROPERTIES,
		XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE,
		XAdES132Element_REFERENCE_INFO,
		XAdES132Element_REFS_ONLY_TIMESTAMP,
		XAdES132Element_RESPONDER_ID,
		XAdES132Element_REVOCATION_VALUES,
		XAdES132Element_SIG_AND_REFS_TIMESTAMP,
		XAdES132Element_SIG_POLICY_HASH,
		XAdES132Element_SIG_POLICY_ID,
		XAdES132Element_SIG_POLICY_QUALIFIER,
		XAdES132Element_SIG_POLICY_QUALIFIERS,
		XAdES132Element_SIGNATURE_POLICY_ID,
		XAdES132Element_SIGNATURE_POLICY_IDENTIFIER,
		XAdES132Element_SIGNATURE_POLICY_IMPLIED,
		XAdES132Element_SIGNATURE_PRODUCTION_PLACE,
		XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2,
		XAdES132Element_SIGNATURE_TIMESTAMP,
		XAdES132Element_SIGNED_ASSERTION,
		XAdES132Element_SIGNED_ASSERTIONS,
		XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES,
		XAdES132Element_SIGNED_PROPERTIES,
		XAdES132Element_SIGNED_SIGNATURE_PROPERTIES,
		XAdES132Element_SIGNER_ROLE,
		XAdES132Element_SIGNER_ROLE_V2,
		XAdES132Element_SIGNING_CERTIFICATE,
		XAdES132Element_SIGNING_CERTIFICATE_V2,
		XAdES132Element_SIGNING_TIME,
		XAdES132Element_SP_URI,
		XAdES132Element_SP_USER_NOTICE,
		XAdES132Element_STATE_OR_PROVINCE,
		XAdES132Element_STREET_ADDRESS,
		XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES,
		XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY,
		XAdES132Element_UNSIGNED_PROPERTIES,
		XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES,
		XAdES132Element_X509_ATTRIBUTE_CERTIFICATE,
		XAdES132Element_XADES_TIMESTAMP,
		XAdES132Element_XML_TIMESTAMP,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES132Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES132Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES132Element_ALL_SIGNED_DATA_OBJECTS
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES132Element) ElementAny() common.DSSElement {
	return XAdES132Element_ANY
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES132Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES132Element_ARCHIVE_TIMESTAMP
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES132Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	return XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES132Element) ElementAttributeCertificateRefs() common.DSSElement {
	return XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES132Element) ElementAttributeRevocationRefs() common.DSSElement {
	return XAdES132Element_ATTRIBUTE_REVOCATION_REFS
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES132Element) ElementAttributeRevocationValues() common.DSSElement {
	return XAdES132Element_ATTRIBUTE_REVOCATION_VALUES
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES132Element) ElementByKey() common.DSSElement {
	return XAdES132Element_BY_KEY
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES132Element) ElementByName() common.DSSElement {
	return XAdES132Element_BY_NAME
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES132Element) ElementCert() common.DSSElement {
	return XAdES132Element_CERT
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES132Element) ElementCertDigest() common.DSSElement {
	return XAdES132Element_CERT_DIGEST
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES132Element) ElementCertRefs() common.DSSElement {
	return XAdES132Element_CERT_REFS
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES132Element) ElementCertificateValues() common.DSSElement {
	return XAdES132Element_CERTIFICATE_VALUES
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES132Element) ElementCertifiedRole() common.DSSElement {
	return XAdES132Element_CERTIFIED_ROLE
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES132Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES132Element_CERTIFIED_ROLES
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES132Element) ElementCertifiedRolesV2() common.DSSElement {
	return XAdES132Element_CERTIFIED_ROLES_V2
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES132Element) ElementCity() common.DSSElement {
	return XAdES132Element_CITY
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES132Element) ElementClaimedRole() common.DSSElement {
	return XAdES132Element_CLAIMED_ROLE
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES132Element) ElementClaimedRoles() common.DSSElement {
	return XAdES132Element_CLAIMED_ROLES
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES132Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES132Element_COMMITMENT_TYPE_ID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES132Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES132Element_COMMITMENT_TYPE_INDICATION
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES132Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES132Element_COMMITMENT_TYPE_QUALIFIER
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES132Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES132Element_COMMITMENT_TYPE_QUALIFIERS
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES132Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES132Element_COMPLETE_CERTIFICATE_REFS
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES132Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES132Element_COMPLETE_REVOCATION_REFS
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES132Element) ElementCounterSignature() common.DSSElement {
	return XAdES132Element_COUNTER_SIGNATURE
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES132Element) ElementCountryName() common.DSSElement {
	return XAdES132Element_COUNTRY_NAME
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES132Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES132Element_CRL_IDENTIFIER
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES132Element) ElementCRLRef() common.DSSElement {
	return XAdES132Element_CRL_REF
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES132Element) ElementCRLRefs() common.DSSElement {
	return XAdES132Element_CRL_REFS
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES132Element) ElementCRLValues() common.DSSElement {
	return XAdES132Element_CRL_VALUES
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES132Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES132Element_DATA_OBJECT_FORMAT
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES132Element) ElementDescription() common.DSSElement {
	return XAdES132Element_DESCRIPTION
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES132Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES132Element_DIGEST_ALG_AND_VALUE
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES132Element) ElementDocumentationReference() common.DSSElement {
	return XAdES132Element_DOCUMENTATION_REFERENCE
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES132Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES132Element_DOCUMENTATION_REFERENCES
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES132Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES132Element_ENCAPSULATED_CRL_VALUE
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES132Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES132Element_ENCAPSULATED_OCSP_VALUE
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES132Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES132Element_ENCAPSULATED_PKI_DATA
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES132Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES132Element_ENCAPSULATED_TIMESTAMP
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES132Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES132Element_ENCAPSULATED_X509_CERTIFICATE
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES132Element) ElementEncoding() common.DSSElement {
	return XAdES132Element_ENCODING
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES132Element) ElementExplicitText() common.DSSElement {
	return XAdES132Element_EXPLICIT_TEXT
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES132Element) ElementIdentifier() common.DSSElement {
	return XAdES132Element_IDENTIFIER
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES132Element) ElementInclude() common.DSSElement {
	return XAdES132Element_INCLUDE
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES132Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES132Element) Elementint() common.DSSElement {
	return XAdES132Element_INT
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES132Element) ElementIssueTime() common.DSSElement {
	return XAdES132Element_ISSUE_TIME
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES132Element) ElementIssuer() common.DSSElement {
	return XAdES132Element_ISSUER
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES132Element) ElementIssuerSerial() common.DSSElement {
	return XAdES132Element_ISSUER_SERIAL
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES132Element) ElementIssuerSerialV2() common.DSSElement {
	return XAdES132Element_ISSUER_SERIAL_V2
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES132Element) ElementMimeType() common.DSSElement {
	return XAdES132Element_MIME_TYPE
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES132Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES132Element_NOTICE_NUMBERS
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES132Element) ElementNoticeRef() common.DSSElement {
	return XAdES132Element_NOTICE_REF
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES132Element) ElementNumber() common.DSSElement {
	return XAdES132Element_NUMBER
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES132Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES132Element_OBJECT_IDENTIFIER
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES132Element) ElementObjectReference() common.DSSElement {
	return XAdES132Element_OBJECT_REFERENCE
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES132Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES132Element_OCSP_IDENTIFIER
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES132Element) ElementOCSPRef() common.DSSElement {
	return XAdES132Element_OCSP_REF
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES132Element) ElementOCSPRefs() common.DSSElement {
	return XAdES132Element_OCSP_REFS
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES132Element) ElementOCSPValues() common.DSSElement {
	return XAdES132Element_OCSP_VALUES
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES132Element) ElementOrganization() common.DSSElement {
	return XAdES132Element_ORGANIZATION
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES132Element) ElementOtherAttributeCertificate() common.DSSElement {
	return XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES132Element) ElementOtherCertificate() common.DSSElement {
	return XAdES132Element_OTHER_CERTIFICATE
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES132Element) ElementOtherRef() common.DSSElement {
	return XAdES132Element_OTHER_REF
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES132Element) ElementOtherRefs() common.DSSElement {
	return XAdES132Element_OTHER_REFS
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES132Element) ElementOtherTimeStamp() common.DSSElement {
	return XAdES132Element_OTHER_TIMESTAMP
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES132Element) ElementOtherValue() common.DSSElement {
	return XAdES132Element_OTHER_VALUE
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES132Element) ElementOtherValues() common.DSSElement {
	return XAdES132Element_OTHER_VALUES
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES132Element) ElementPostalCode() common.DSSElement {
	return XAdES132Element_POSTAL_CODE
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES132Element) ElementProducedAt() common.DSSElement {
	return XAdES132Element_PRODUCED_AT
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES132Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES132Element_QUALIFYING_PROPERTIES
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES132Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES132Element) ElementReferenceInfo() common.DSSElement {
	return XAdES132Element_REFERENCE_INFO
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES132Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES132Element_REFS_ONLY_TIMESTAMP
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES132Element) ElementResponderID() common.DSSElement {
	return XAdES132Element_RESPONDER_ID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES132Element) ElementRevocationValues() common.DSSElement {
	return XAdES132Element_REVOCATION_VALUES
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES132Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES132Element_SIG_AND_REFS_TIMESTAMP
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES132Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES132Element_SIG_POLICY_HASH
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES132Element) ElementSigPolicyId() common.DSSElement {
	return XAdES132Element_SIG_POLICY_ID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES132Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES132Element_SIG_POLICY_QUALIFIER
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES132Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES132Element_SIG_POLICY_QUALIFIERS
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES132Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES132Element_SIGNATURE_POLICY_ID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES132Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES132Element_SIGNATURE_POLICY_IDENTIFIER
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES132Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES132Element_SIGNATURE_POLICY_IMPLIED
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES132Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES132Element_SIGNATURE_PRODUCTION_PLACE
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES132Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	return XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES132Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES132Element_SIGNATURE_TIMESTAMP
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES132Element) ElementSignedAssertion() common.DSSElement {
	return XAdES132Element_SIGNED_ASSERTION
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES132Element) ElementSignedAssertions() common.DSSElement {
	return XAdES132Element_SIGNED_ASSERTIONS
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES132Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES132Element) ElementSignedProperties() common.DSSElement {
	return XAdES132Element_SIGNED_PROPERTIES
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES132Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES132Element_SIGNED_SIGNATURE_PROPERTIES
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES132Element) ElementSignerRole() common.DSSElement {
	return XAdES132Element_SIGNER_ROLE
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES132Element) ElementSignerRoleV2() common.DSSElement {
	return XAdES132Element_SIGNER_ROLE_V2
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES132Element) ElementSigningCertificate() common.DSSElement {
	return XAdES132Element_SIGNING_CERTIFICATE
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES132Element) ElementSigningCertificateV2() common.DSSElement {
	return XAdES132Element_SIGNING_CERTIFICATE_V2
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES132Element) ElementSigningTime() common.DSSElement {
	return XAdES132Element_SIGNING_TIME
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES132Element) ElementSPURI() common.DSSElement {
	return XAdES132Element_SP_URI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES132Element) ElementSPUserNotice() common.DSSElement {
	return XAdES132Element_SP_USER_NOTICE
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES132Element) ElementStateOrProvince() common.DSSElement {
	return XAdES132Element_STATE_OR_PROVINCE
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES132Element) ElementStreetAddress() common.DSSElement {
	return XAdES132Element_STREET_ADDRESS
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES132Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES132Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES132Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES132Element_UNSIGNED_PROPERTIES
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES132Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES132Element) ElementX509AttributeCertificate() common.DSSElement {
	return XAdES132Element_X509_ATTRIBUTE_CERTIFICATE
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES132Element) ElementXAdESTimeStamp() common.DSSElement {
	return XAdES132Element_XADES_TIMESTAMP
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES132Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES132Element_XML_TIMESTAMP
}

var _ XAdESElement = XAdES132Element("")
