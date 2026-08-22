// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades122/XAdES122Element.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES122Element defines elements for a XAdES 1.2.2 schema.
type XAdES122Element string

// XAdES122Element constants, one per XAdES 1.2.2 schema element name.
const (
	XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP        XAdES122Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES122Element_ALL_SIGNED_DATA_OBJECTS           XAdES122Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES122Element_ANY                               XAdES122Element = "ANY"
	XAdES122Element_ARCHIVE_TIMESTAMP                 XAdES122Element = "ARCHIVE_TIMESTAMP"
	XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS        XAdES122Element = "ATTRIBUTE_CERTIFICATE_REFS"
	XAdES122Element_ATTRIBUTE_REVOCATION_REFS         XAdES122Element = "ATTRIBUTE_REVOCATION_REFS"
	XAdES122Element_CERT                              XAdES122Element = "CERT"
	XAdES122Element_CERT_DIGEST                       XAdES122Element = "CERT_DIGEST"
	XAdES122Element_CERT_REFS                         XAdES122Element = "CERT_REFS"
	XAdES122Element_CERTIFICATE_VALUES                XAdES122Element = "CERTIFICATE_VALUES"
	XAdES122Element_CERTIFIED_ROLE                    XAdES122Element = "CERTIFIED_ROLE"
	XAdES122Element_CERTIFIED_ROLES                   XAdES122Element = "CERTIFIED_ROLES"
	XAdES122Element_CITY                              XAdES122Element = "CITY"
	XAdES122Element_CLAIMED_ROLE                      XAdES122Element = "CLAIMED_ROLE"
	XAdES122Element_CLAIMED_ROLES                     XAdES122Element = "CLAIMED_ROLES"
	XAdES122Element_COMMITMENT_TYPE_ID                XAdES122Element = "COMMITMENT_TYPE_ID"
	XAdES122Element_COMMITMENT_TYPE_INDICATION        XAdES122Element = "COMMITMENT_TYPE_INDICATION"
	XAdES122Element_COMMITMENT_TYPE_QUALIFIER         XAdES122Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES122Element_COMMITMENT_TYPE_QUALIFIERS        XAdES122Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES122Element_COMPLETE_CERTIFICATE_REFS         XAdES122Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES122Element_COMPLETE_REVOCATION_REFS          XAdES122Element = "COMPLETE_REVOCATION_REFS"
	XAdES122Element_COUNTER_SIGNATURE                 XAdES122Element = "COUNTER_SIGNATURE"
	XAdES122Element_COUNTRY_NAME                      XAdES122Element = "COUNTRY_NAME"
	XAdES122Element_CRL_IDENTIFIER                    XAdES122Element = "CRL_IDENTIFIER"
	XAdES122Element_CRL_REF                           XAdES122Element = "CRL_REF"
	XAdES122Element_CRL_REFS                          XAdES122Element = "CRL_REFS"
	XAdES122Element_CRL_VALUES                        XAdES122Element = "CRL_VALUES"
	XAdES122Element_DATA_OBJECT_FORMAT                XAdES122Element = "DATA_OBJECT_FORMAT"
	XAdES122Element_DESCRIPTION                       XAdES122Element = "DESCRIPTION"
	XAdES122Element_DIGEST_ALG_AND_VALUE              XAdES122Element = "DIGEST_ALG_AND_VALUE"
	XAdES122Element_DOCUMENTATION_REFERENCE           XAdES122Element = "DOCUMENTATION_REFERENCE"
	XAdES122Element_DOCUMENTATION_REFERENCES          XAdES122Element = "DOCUMENTATION_REFERENCES"
	XAdES122Element_ENCAPSULATED_CRL_VALUE            XAdES122Element = "ENCAPSULATED_CRL_VALUE"
	XAdES122Element_ENCAPSULATED_OCSP_VALUE           XAdES122Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES122Element_ENCAPSULATED_PKI_DATA             XAdES122Element = "ENCAPSULATED_PKI_DATA"
	XAdES122Element_ENCAPSULATED_TIMESTAMP            XAdES122Element = "ENCAPSULATED_TIMESTAMP"
	XAdES122Element_ENCAPSULATED_X509_CERTIFICATE     XAdES122Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES122Element_ENCODING                          XAdES122Element = "ENCODING"
	XAdES122Element_EXPLICIT_TEXT                     XAdES122Element = "EXPLICIT_TEXT"
	XAdES122Element_IDENTIFIER                        XAdES122Element = "IDENTIFIER"
	XAdES122Element_INCLUDE                           XAdES122Element = "INCLUDE"
	XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP XAdES122Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES122Element_INT                               XAdES122Element = "INT"
	XAdES122Element_ISSUE_TIME                        XAdES122Element = "ISSUE_TIME"
	XAdES122Element_ISSUER                            XAdES122Element = "ISSUER"
	XAdES122Element_ISSUER_SERIAL                     XAdES122Element = "ISSUER_SERIAL"
	XAdES122Element_MIME_TYPE                         XAdES122Element = "MIME_TYPE"
	XAdES122Element_NOTICE_NUMBERS                    XAdES122Element = "NOTICE_NUMBERS"
	XAdES122Element_NOTICE_REF                        XAdES122Element = "NOTICE_REF"
	XAdES122Element_NUMBER                            XAdES122Element = "NUMBER"
	XAdES122Element_OBJECT_IDENTIFIER                 XAdES122Element = "OBJECT_IDENTIFIER"
	XAdES122Element_OBJECT_REFERENCE                  XAdES122Element = "OBJECT_REFERENCE"
	XAdES122Element_OCSP_IDENTIFIER                   XAdES122Element = "OCSP_IDENTIFIER"
	XAdES122Element_OCSP_REF                          XAdES122Element = "OCSP_REF"
	XAdES122Element_OCSP_REFS                         XAdES122Element = "OCSP_REFS"
	XAdES122Element_OCSP_VALUES                       XAdES122Element = "OCSP_VALUES"
	XAdES122Element_ORGANIZATION                      XAdES122Element = "ORGANIZATION"
	XAdES122Element_OTHER_CERTIFICATE                 XAdES122Element = "OTHER_CERTIFICATE"
	XAdES122Element_OTHER_REF                         XAdES122Element = "OTHER_REF"
	XAdES122Element_OTHER_REFS                        XAdES122Element = "OTHER_REFS"
	XAdES122Element_OTHER_VALUE                       XAdES122Element = "OTHER_VALUE"
	XAdES122Element_OTHER_VALUES                      XAdES122Element = "OTHER_VALUES"
	XAdES122Element_POSTAL_CODE                       XAdES122Element = "POSTAL_CODE"
	XAdES122Element_PRODUCED_AT                       XAdES122Element = "PRODUCED_AT"
	XAdES122Element_QUALIFYING_PROPERTIES             XAdES122Element = "QUALIFYING_PROPERTIES"
	XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE   XAdES122Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES122Element_REFS_ONLY_TIMESTAMP               XAdES122Element = "REFS_ONLY_TIMESTAMP"
	XAdES122Element_RESPONDER_ID                      XAdES122Element = "RESPONDER_ID"
	XAdES122Element_REVOCATION_VALUES                 XAdES122Element = "REVOCATION_VALUES"
	XAdES122Element_SIG_AND_REFS_TIMESTAMP            XAdES122Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES122Element_SIG_POLICY_HASH                   XAdES122Element = "SIG_POLICY_HASH"
	XAdES122Element_SIG_POLICY_ID                     XAdES122Element = "SIG_POLICY_ID"
	XAdES122Element_SIG_POLICY_QUALIFIER              XAdES122Element = "SIG_POLICY_QUALIFIER"
	XAdES122Element_SIG_POLICY_QUALIFIERS             XAdES122Element = "SIG_POLICY_QUALIFIERS"
	XAdES122Element_SIGNATURE_POLICY_ID               XAdES122Element = "SIGNATURE_POLICY_ID"
	XAdES122Element_SIGNATURE_POLICY_IDENTIFIER       XAdES122Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES122Element_SIGNATURE_POLICY_IMPLIED          XAdES122Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES122Element_SIGNATURE_PRODUCTION_PLACE        XAdES122Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES122Element_SIGNATURE_TIMESTAMP               XAdES122Element = "SIGNATURE_TIMESTAMP"
	XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES     XAdES122Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES122Element_SIGNED_PROPERTIES                 XAdES122Element = "SIGNED_PROPERTIES"
	XAdES122Element_SIGNED_SIGNATURE_PROPERTIES       XAdES122Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES122Element_SIGNER_ROLE                       XAdES122Element = "SIGNER_ROLE"
	XAdES122Element_SIGNING_CERTIFICATE               XAdES122Element = "SIGNING_CERTIFICATE"
	XAdES122Element_SIGNING_TIME                      XAdES122Element = "SIGNING_TIME"
	XAdES122Element_SP_URI                            XAdES122Element = "SP_URI"
	XAdES122Element_SP_USER_NOTICE                    XAdES122Element = "SP_USER_NOTICE"
	XAdES122Element_STATE_OR_PROVINCE                 XAdES122Element = "STATE_OR_PROVINCE"
	XAdES122Element_TIMESTAMP                         XAdES122Element = "TIMESTAMP"
	XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES   XAdES122Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY     XAdES122Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES122Element_UNSIGNED_PROPERTIES               XAdES122Element = "UNSIGNED_PROPERTIES"
	XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES     XAdES122Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES122Element_XML_TIMESTAMP                     XAdES122Element = "XML_TIMESTAMP"
)

// xades122elementTagNames maps each constant to its wire tag name (getTagName()).
var xades122elementTagNames = map[XAdES122Element]string{
	XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP:        "AllDataObjectsTimeStamp",
	XAdES122Element_ALL_SIGNED_DATA_OBJECTS:           "AllSignedDataObjects",
	XAdES122Element_ANY:                               "Any",
	XAdES122Element_ARCHIVE_TIMESTAMP:                 "ArchiveTimeStamp",
	XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS:        "AttributeCertificateRefs",
	XAdES122Element_ATTRIBUTE_REVOCATION_REFS:         "AttributeRevocationRefs",
	XAdES122Element_CERT:                              "Cert",
	XAdES122Element_CERT_DIGEST:                       "CertDigest",
	XAdES122Element_CERT_REFS:                         "CertRefs",
	XAdES122Element_CERTIFICATE_VALUES:                "CertificateValues",
	XAdES122Element_CERTIFIED_ROLE:                    "CertifiedRole",
	XAdES122Element_CERTIFIED_ROLES:                   "CertifiedRoles",
	XAdES122Element_CITY:                              "City",
	XAdES122Element_CLAIMED_ROLE:                      "ClaimedRole",
	XAdES122Element_CLAIMED_ROLES:                     "ClaimedRoles",
	XAdES122Element_COMMITMENT_TYPE_ID:                "CommitmentTypeId",
	XAdES122Element_COMMITMENT_TYPE_INDICATION:        "CommitmentTypeIndication",
	XAdES122Element_COMMITMENT_TYPE_QUALIFIER:         "CommitmentTypeQualifier",
	XAdES122Element_COMMITMENT_TYPE_QUALIFIERS:        "CommitmentTypeQualifiers",
	XAdES122Element_COMPLETE_CERTIFICATE_REFS:         "CompleteCertificateRefs",
	XAdES122Element_COMPLETE_REVOCATION_REFS:          "CompleteRevocationRefs",
	XAdES122Element_COUNTER_SIGNATURE:                 "CounterSignature",
	XAdES122Element_COUNTRY_NAME:                      "CountryName",
	XAdES122Element_CRL_IDENTIFIER:                    "CRLIdentifier",
	XAdES122Element_CRL_REF:                           "CRLRef",
	XAdES122Element_CRL_REFS:                          "CRLRefs",
	XAdES122Element_CRL_VALUES:                        "CRLValues",
	XAdES122Element_DATA_OBJECT_FORMAT:                "DataObjectFormat",
	XAdES122Element_DESCRIPTION:                       "Description",
	XAdES122Element_DIGEST_ALG_AND_VALUE:              "DigestAlgAndValue",
	XAdES122Element_DOCUMENTATION_REFERENCE:           "DocumentationReference",
	XAdES122Element_DOCUMENTATION_REFERENCES:          "DocumentationReferences",
	XAdES122Element_ENCAPSULATED_CRL_VALUE:            "EncapsulatedCRLValue",
	XAdES122Element_ENCAPSULATED_OCSP_VALUE:           "EncapsulatedOCSPValue",
	XAdES122Element_ENCAPSULATED_PKI_DATA:             "EncapsulatedPKIData",
	XAdES122Element_ENCAPSULATED_TIMESTAMP:            "EncapsulatedTimeStamp",
	XAdES122Element_ENCAPSULATED_X509_CERTIFICATE:     "EncapsulatedX509Certificate",
	XAdES122Element_ENCODING:                          "Encoding",
	XAdES122Element_EXPLICIT_TEXT:                     "ExplicitText",
	XAdES122Element_IDENTIFIER:                        "Identifier",
	XAdES122Element_INCLUDE:                           "Include",
	XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP: "IndividualDataObjectsTimeStamp",
	XAdES122Element_INT:                               "int",
	XAdES122Element_ISSUE_TIME:                        "IssueTime",
	XAdES122Element_ISSUER:                            "Issuer",
	XAdES122Element_ISSUER_SERIAL:                     "IssuerSerial",
	XAdES122Element_MIME_TYPE:                         "MimeType",
	XAdES122Element_NOTICE_NUMBERS:                    "NoticeNumbers",
	XAdES122Element_NOTICE_REF:                        "NoticeRef",
	XAdES122Element_NUMBER:                            "Number",
	XAdES122Element_OBJECT_IDENTIFIER:                 "ObjectIdentifier",
	XAdES122Element_OBJECT_REFERENCE:                  "ObjectReference",
	XAdES122Element_OCSP_IDENTIFIER:                   "OCSPIdentifier",
	XAdES122Element_OCSP_REF:                          "OCSPRef",
	XAdES122Element_OCSP_REFS:                         "OCSPRefs",
	XAdES122Element_OCSP_VALUES:                       "OCSPValues",
	XAdES122Element_ORGANIZATION:                      "Organization",
	XAdES122Element_OTHER_CERTIFICATE:                 "OtherCertificate",
	XAdES122Element_OTHER_REF:                         "OtherRef",
	XAdES122Element_OTHER_REFS:                        "OtherRefs",
	XAdES122Element_OTHER_VALUE:                       "OtherValue",
	XAdES122Element_OTHER_VALUES:                      "OtherValues",
	XAdES122Element_POSTAL_CODE:                       "PostalCode",
	XAdES122Element_PRODUCED_AT:                       "ProducedAt",
	XAdES122Element_QUALIFYING_PROPERTIES:             "QualifyingProperties",
	XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE:   "QualifyingPropertiesReference",
	XAdES122Element_REFS_ONLY_TIMESTAMP:               "RefsOnlyTimeStamp",
	XAdES122Element_RESPONDER_ID:                      "ResponderID",
	XAdES122Element_REVOCATION_VALUES:                 "RevocationValues",
	XAdES122Element_SIG_AND_REFS_TIMESTAMP:            "SigAndRefsTimeStamp",
	XAdES122Element_SIG_POLICY_HASH:                   "SigPolicyHash",
	XAdES122Element_SIG_POLICY_ID:                     "SigPolicyId",
	XAdES122Element_SIG_POLICY_QUALIFIER:              "SigPolicyQualifier",
	XAdES122Element_SIG_POLICY_QUALIFIERS:             "SigPolicyQualifiers",
	XAdES122Element_SIGNATURE_POLICY_ID:               "SignaturePolicyId",
	XAdES122Element_SIGNATURE_POLICY_IDENTIFIER:       "SignaturePolicyIdentifier",
	XAdES122Element_SIGNATURE_POLICY_IMPLIED:          "SignaturePolicyImplied",
	XAdES122Element_SIGNATURE_PRODUCTION_PLACE:        "SignatureProductionPlace",
	XAdES122Element_SIGNATURE_TIMESTAMP:               "SignatureTimeStamp",
	XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES:     "SignedDataObjectProperties",
	XAdES122Element_SIGNED_PROPERTIES:                 "SignedProperties",
	XAdES122Element_SIGNED_SIGNATURE_PROPERTIES:       "SignedSignatureProperties",
	XAdES122Element_SIGNER_ROLE:                       "SignerRole",
	XAdES122Element_SIGNING_CERTIFICATE:               "SigningCertificate",
	XAdES122Element_SIGNING_TIME:                      "SigningTime",
	XAdES122Element_SP_URI:                            "SPURI",
	XAdES122Element_SP_USER_NOTICE:                    "SPUserNotice",
	XAdES122Element_STATE_OR_PROVINCE:                 "StateOrProvince",
	XAdES122Element_TIMESTAMP:                         "TimeStamp",
	XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES:   "UnsignedDataObjectProperties",
	XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY:     "UnsignedDataObjectProperty",
	XAdES122Element_UNSIGNED_PROPERTIES:               "UnsignedProperties",
	XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES:     "UnsignedSignatureProperties",
	XAdES122Element_XML_TIMESTAMP:                     "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES122Element) TagName() string {
	return xades122elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES122Element) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_122
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES122Element) URI() string {
	return XAdESNamespace_XADES_122.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e XAdES122Element) IsSameTagName(value string) bool {
	return e.TagName() == value
}

// XAdES122ElementFromTagName returns the XAdES122Element constant with the given tag name, or nil
// if none matches. Ports the static factory fromTagName(String).
func XAdES122ElementFromTagName(tagName string) common.DSSElement {
	for _, e := range XAdES122ElementValues() {
		if e.TagName() == tagName {
			return e
		}
	}
	return nil
}

// XAdES122ElementValues returns every XAdES122Element constant, in declaration order. Ports values().
func XAdES122ElementValues() []XAdES122Element {
	return []XAdES122Element{
		XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP,
		XAdES122Element_ALL_SIGNED_DATA_OBJECTS,
		XAdES122Element_ANY,
		XAdES122Element_ARCHIVE_TIMESTAMP,
		XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS,
		XAdES122Element_ATTRIBUTE_REVOCATION_REFS,
		XAdES122Element_CERT,
		XAdES122Element_CERT_DIGEST,
		XAdES122Element_CERT_REFS,
		XAdES122Element_CERTIFICATE_VALUES,
		XAdES122Element_CERTIFIED_ROLE,
		XAdES122Element_CERTIFIED_ROLES,
		XAdES122Element_CITY,
		XAdES122Element_CLAIMED_ROLE,
		XAdES122Element_CLAIMED_ROLES,
		XAdES122Element_COMMITMENT_TYPE_ID,
		XAdES122Element_COMMITMENT_TYPE_INDICATION,
		XAdES122Element_COMMITMENT_TYPE_QUALIFIER,
		XAdES122Element_COMMITMENT_TYPE_QUALIFIERS,
		XAdES122Element_COMPLETE_CERTIFICATE_REFS,
		XAdES122Element_COMPLETE_REVOCATION_REFS,
		XAdES122Element_COUNTER_SIGNATURE,
		XAdES122Element_COUNTRY_NAME,
		XAdES122Element_CRL_IDENTIFIER,
		XAdES122Element_CRL_REF,
		XAdES122Element_CRL_REFS,
		XAdES122Element_CRL_VALUES,
		XAdES122Element_DATA_OBJECT_FORMAT,
		XAdES122Element_DESCRIPTION,
		XAdES122Element_DIGEST_ALG_AND_VALUE,
		XAdES122Element_DOCUMENTATION_REFERENCE,
		XAdES122Element_DOCUMENTATION_REFERENCES,
		XAdES122Element_ENCAPSULATED_CRL_VALUE,
		XAdES122Element_ENCAPSULATED_OCSP_VALUE,
		XAdES122Element_ENCAPSULATED_PKI_DATA,
		XAdES122Element_ENCAPSULATED_TIMESTAMP,
		XAdES122Element_ENCAPSULATED_X509_CERTIFICATE,
		XAdES122Element_ENCODING,
		XAdES122Element_EXPLICIT_TEXT,
		XAdES122Element_IDENTIFIER,
		XAdES122Element_INCLUDE,
		XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP,
		XAdES122Element_INT,
		XAdES122Element_ISSUE_TIME,
		XAdES122Element_ISSUER,
		XAdES122Element_ISSUER_SERIAL,
		XAdES122Element_MIME_TYPE,
		XAdES122Element_NOTICE_NUMBERS,
		XAdES122Element_NOTICE_REF,
		XAdES122Element_NUMBER,
		XAdES122Element_OBJECT_IDENTIFIER,
		XAdES122Element_OBJECT_REFERENCE,
		XAdES122Element_OCSP_IDENTIFIER,
		XAdES122Element_OCSP_REF,
		XAdES122Element_OCSP_REFS,
		XAdES122Element_OCSP_VALUES,
		XAdES122Element_ORGANIZATION,
		XAdES122Element_OTHER_CERTIFICATE,
		XAdES122Element_OTHER_REF,
		XAdES122Element_OTHER_REFS,
		XAdES122Element_OTHER_VALUE,
		XAdES122Element_OTHER_VALUES,
		XAdES122Element_POSTAL_CODE,
		XAdES122Element_PRODUCED_AT,
		XAdES122Element_QUALIFYING_PROPERTIES,
		XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE,
		XAdES122Element_REFS_ONLY_TIMESTAMP,
		XAdES122Element_RESPONDER_ID,
		XAdES122Element_REVOCATION_VALUES,
		XAdES122Element_SIG_AND_REFS_TIMESTAMP,
		XAdES122Element_SIG_POLICY_HASH,
		XAdES122Element_SIG_POLICY_ID,
		XAdES122Element_SIG_POLICY_QUALIFIER,
		XAdES122Element_SIG_POLICY_QUALIFIERS,
		XAdES122Element_SIGNATURE_POLICY_ID,
		XAdES122Element_SIGNATURE_POLICY_IDENTIFIER,
		XAdES122Element_SIGNATURE_POLICY_IMPLIED,
		XAdES122Element_SIGNATURE_PRODUCTION_PLACE,
		XAdES122Element_SIGNATURE_TIMESTAMP,
		XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES,
		XAdES122Element_SIGNED_PROPERTIES,
		XAdES122Element_SIGNED_SIGNATURE_PROPERTIES,
		XAdES122Element_SIGNER_ROLE,
		XAdES122Element_SIGNING_CERTIFICATE,
		XAdES122Element_SIGNING_TIME,
		XAdES122Element_SP_URI,
		XAdES122Element_SP_USER_NOTICE,
		XAdES122Element_STATE_OR_PROVINCE,
		XAdES122Element_TIMESTAMP,
		XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES,
		XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY,
		XAdES122Element_UNSIGNED_PROPERTIES,
		XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES,
		XAdES122Element_XML_TIMESTAMP,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES122Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES122Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES122Element_ALL_SIGNED_DATA_OBJECTS
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES122Element) ElementAny() common.DSSElement {
	return XAdES122Element_ANY
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES122Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES122Element_ARCHIVE_TIMESTAMP
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES122Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES122Element) ElementAttributeCertificateRefs() common.DSSElement {
	return XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES122Element) ElementAttributeRevocationRefs() common.DSSElement {
	return XAdES122Element_ATTRIBUTE_REVOCATION_REFS
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES122Element) ElementAttributeRevocationValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES122Element) ElementByKey() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES122Element) ElementByName() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES122Element) ElementCert() common.DSSElement {
	return XAdES122Element_CERT
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES122Element) ElementCertDigest() common.DSSElement {
	return XAdES122Element_CERT_DIGEST
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES122Element) ElementCertRefs() common.DSSElement {
	return XAdES122Element_CERT_REFS
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES122Element) ElementCertificateValues() common.DSSElement {
	return XAdES122Element_CERTIFICATE_VALUES
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES122Element) ElementCertifiedRole() common.DSSElement {
	return XAdES122Element_CERTIFIED_ROLE
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES122Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES122Element_CERTIFIED_ROLES
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES122Element) ElementCertifiedRolesV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES122Element) ElementCity() common.DSSElement {
	return XAdES122Element_CITY
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES122Element) ElementClaimedRole() common.DSSElement {
	return XAdES122Element_CLAIMED_ROLE
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES122Element) ElementClaimedRoles() common.DSSElement {
	return XAdES122Element_CLAIMED_ROLES
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES122Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES122Element_COMMITMENT_TYPE_ID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES122Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES122Element_COMMITMENT_TYPE_INDICATION
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES122Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES122Element_COMMITMENT_TYPE_QUALIFIER
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES122Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES122Element_COMMITMENT_TYPE_QUALIFIERS
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES122Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES122Element_COMPLETE_CERTIFICATE_REFS
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES122Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES122Element_COMPLETE_REVOCATION_REFS
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES122Element) ElementCounterSignature() common.DSSElement {
	return XAdES122Element_COUNTER_SIGNATURE
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES122Element) ElementCountryName() common.DSSElement {
	return XAdES122Element_COUNTRY_NAME
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES122Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES122Element_CRL_IDENTIFIER
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES122Element) ElementCRLRef() common.DSSElement {
	return XAdES122Element_CRL_REF
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES122Element) ElementCRLRefs() common.DSSElement {
	return XAdES122Element_CRL_REFS
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES122Element) ElementCRLValues() common.DSSElement {
	return XAdES122Element_CRL_VALUES
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES122Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES122Element_DATA_OBJECT_FORMAT
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES122Element) ElementDescription() common.DSSElement {
	return XAdES122Element_DESCRIPTION
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES122Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES122Element_DIGEST_ALG_AND_VALUE
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES122Element) ElementDocumentationReference() common.DSSElement {
	return XAdES122Element_DOCUMENTATION_REFERENCE
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES122Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES122Element_DOCUMENTATION_REFERENCES
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES122Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES122Element_ENCAPSULATED_CRL_VALUE
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES122Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES122Element_ENCAPSULATED_OCSP_VALUE
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES122Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES122Element_ENCAPSULATED_PKI_DATA
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES122Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES122Element_ENCAPSULATED_TIMESTAMP
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES122Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES122Element_ENCAPSULATED_X509_CERTIFICATE
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES122Element) ElementEncoding() common.DSSElement {
	return XAdES122Element_ENCODING
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES122Element) ElementExplicitText() common.DSSElement {
	return XAdES122Element_EXPLICIT_TEXT
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES122Element) ElementIdentifier() common.DSSElement {
	return XAdES122Element_IDENTIFIER
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES122Element) ElementInclude() common.DSSElement {
	return XAdES122Element_INCLUDE
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES122Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES122Element) Elementint() common.DSSElement {
	return XAdES122Element_INT
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES122Element) ElementIssueTime() common.DSSElement {
	return XAdES122Element_ISSUE_TIME
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES122Element) ElementIssuer() common.DSSElement {
	return XAdES122Element_ISSUER
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES122Element) ElementIssuerSerial() common.DSSElement {
	return XAdES122Element_ISSUER_SERIAL
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES122Element) ElementIssuerSerialV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES122Element) ElementMimeType() common.DSSElement {
	return XAdES122Element_MIME_TYPE
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES122Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES122Element_NOTICE_NUMBERS
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES122Element) ElementNoticeRef() common.DSSElement {
	return XAdES122Element_NOTICE_REF
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES122Element) ElementNumber() common.DSSElement {
	return XAdES122Element_NUMBER
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES122Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES122Element_OBJECT_IDENTIFIER
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES122Element) ElementObjectReference() common.DSSElement {
	return XAdES122Element_OBJECT_REFERENCE
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES122Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES122Element_OCSP_IDENTIFIER
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES122Element) ElementOCSPRef() common.DSSElement {
	return XAdES122Element_OCSP_REF
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES122Element) ElementOCSPRefs() common.DSSElement {
	return XAdES122Element_OCSP_REFS
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES122Element) ElementOCSPValues() common.DSSElement {
	return XAdES122Element_OCSP_VALUES
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES122Element) ElementOrganization() common.DSSElement {
	return XAdES122Element_ORGANIZATION
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES122Element) ElementOtherAttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES122Element) ElementOtherCertificate() common.DSSElement {
	return XAdES122Element_OTHER_CERTIFICATE
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES122Element) ElementOtherRef() common.DSSElement {
	return XAdES122Element_OTHER_REF
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES122Element) ElementOtherRefs() common.DSSElement {
	return XAdES122Element_OTHER_REFS
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES122Element) ElementOtherTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES122Element) ElementOtherValue() common.DSSElement {
	return XAdES122Element_OTHER_VALUE
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES122Element) ElementOtherValues() common.DSSElement {
	return XAdES122Element_OTHER_VALUES
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES122Element) ElementPostalCode() common.DSSElement {
	return XAdES122Element_POSTAL_CODE
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES122Element) ElementProducedAt() common.DSSElement {
	return XAdES122Element_PRODUCED_AT
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES122Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES122Element_QUALIFYING_PROPERTIES
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES122Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES122Element) ElementReferenceInfo() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES122Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES122Element_REFS_ONLY_TIMESTAMP
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES122Element) ElementResponderID() common.DSSElement {
	return XAdES122Element_RESPONDER_ID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES122Element) ElementRevocationValues() common.DSSElement {
	return XAdES122Element_REVOCATION_VALUES
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES122Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES122Element_SIG_AND_REFS_TIMESTAMP
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES122Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES122Element_SIG_POLICY_HASH
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES122Element) ElementSigPolicyId() common.DSSElement {
	return XAdES122Element_SIG_POLICY_ID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES122Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES122Element_SIG_POLICY_QUALIFIER
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES122Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES122Element_SIG_POLICY_QUALIFIERS
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES122Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES122Element_SIGNATURE_POLICY_ID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES122Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES122Element_SIGNATURE_POLICY_IDENTIFIER
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES122Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES122Element_SIGNATURE_POLICY_IMPLIED
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES122Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES122Element_SIGNATURE_PRODUCTION_PLACE
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES122Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES122Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES122Element_SIGNATURE_TIMESTAMP
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES122Element) ElementSignedAssertion() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES122Element) ElementSignedAssertions() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES122Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES122Element) ElementSignedProperties() common.DSSElement {
	return XAdES122Element_SIGNED_PROPERTIES
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES122Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES122Element_SIGNED_SIGNATURE_PROPERTIES
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES122Element) ElementSignerRole() common.DSSElement {
	return XAdES122Element_SIGNER_ROLE
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES122Element) ElementSignerRoleV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES122Element) ElementSigningCertificate() common.DSSElement {
	return XAdES122Element_SIGNING_CERTIFICATE
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES122Element) ElementSigningCertificateV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES122Element) ElementSigningTime() common.DSSElement {
	return XAdES122Element_SIGNING_TIME
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES122Element) ElementSPURI() common.DSSElement {
	return XAdES122Element_SP_URI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES122Element) ElementSPUserNotice() common.DSSElement {
	return XAdES122Element_SP_USER_NOTICE
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES122Element) ElementStateOrProvince() common.DSSElement {
	return XAdES122Element_STATE_OR_PROVINCE
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES122Element) ElementStreetAddress() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES122Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES122Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES122Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES122Element_UNSIGNED_PROPERTIES
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES122Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES122Element) ElementX509AttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES122Element) ElementXAdESTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_122.Uri())
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES122Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES122Element_XML_TIMESTAMP
}

var _ XAdESElement = XAdES122Element("")
