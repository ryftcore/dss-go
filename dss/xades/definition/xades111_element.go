// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Element.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES111Element defines elements for a XAdES 1.1.1 schema.
type XAdES111Element string

// XAdES111Element constants, one per XAdES 1.1.1 schema element name.
const (
	XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP        XAdES111Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES111Element_ALL_SIGNED_DATA_OBJECTS           XAdES111Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES111Element_ANY                               XAdES111Element = "ANY"
	XAdES111Element_ARCHIVE_TIMESTAMP                 XAdES111Element = "ARCHIVE_TIMESTAMP"
	XAdES111Element_CERT                              XAdES111Element = "CERT"
	XAdES111Element_CERT_DIGEST                       XAdES111Element = "CERT_DIGEST"
	XAdES111Element_CERT_REFS                         XAdES111Element = "CERT_REFS"
	XAdES111Element_CERTIFICATE_VALUES                XAdES111Element = "CERTIFICATE_VALUES"
	XAdES111Element_CERTIFIED_ROLE                    XAdES111Element = "CERTIFIED_ROLE"
	XAdES111Element_CERTIFIED_ROLES                   XAdES111Element = "CERTIFIED_ROLES"
	XAdES111Element_CITY                              XAdES111Element = "CITY"
	XAdES111Element_CLAIMED_ROLE                      XAdES111Element = "CLAIMED_ROLE"
	XAdES111Element_CLAIMED_ROLES                     XAdES111Element = "CLAIMED_ROLES"
	XAdES111Element_COMMITMENT_TYPE_ID                XAdES111Element = "COMMITMENT_TYPE_ID"
	XAdES111Element_COMMITMENT_TYPE_INDICATION        XAdES111Element = "COMMITMENT_TYPE_INDICATION"
	XAdES111Element_COMMITMENT_TYPE_QUALIFIER         XAdES111Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES111Element_COMMITMENT_TYPE_QUALIFIERS        XAdES111Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES111Element_COMPLETE_CERTIFICATE_REFS         XAdES111Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES111Element_COMPLETE_REVOCATION_REFS          XAdES111Element = "COMPLETE_REVOCATION_REFS"
	XAdES111Element_COUNTER_SIGNATURE                 XAdES111Element = "COUNTER_SIGNATURE"
	XAdES111Element_COUNTRY_NAME                      XAdES111Element = "COUNTRY_NAME"
	XAdES111Element_CRL_IDENTIFIER                    XAdES111Element = "CRL_IDENTIFIER"
	XAdES111Element_CRL_REF                           XAdES111Element = "CRL_REF"
	XAdES111Element_CRL_REFS                          XAdES111Element = "CRL_REFS"
	XAdES111Element_CRL_VALUES                        XAdES111Element = "CRL_VALUES"
	XAdES111Element_DATA_OBJECT_FORMAT                XAdES111Element = "DATA_OBJECT_FORMAT"
	XAdES111Element_DESCRIPTION                       XAdES111Element = "DESCRIPTION"
	XAdES111Element_DIGEST_ALG_AND_VALUE              XAdES111Element = "DIGEST_ALG_AND_VALUE"
	XAdES111Element_DIGEST_METHOD                     XAdES111Element = "DIGEST_METHOD"
	XAdES111Element_DIGEST_VALUE                      XAdES111Element = "DIGEST_VALUE"
	XAdES111Element_DOCUMENTATION_REFERENCE           XAdES111Element = "DOCUMENTATION_REFERENCE"
	XAdES111Element_DOCUMENTATION_REFERENCES          XAdES111Element = "DOCUMENTATION_REFERENCES"
	XAdES111Element_ENCAPSULATED_CRL_VALUE            XAdES111Element = "ENCAPSULATED_CRL_VALUE"
	XAdES111Element_ENCAPSULATED_OCSP_VALUE           XAdES111Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES111Element_ENCAPSULATED_PKI_DATA             XAdES111Element = "ENCAPSULATED_PKI_DATA"
	XAdES111Element_ENCAPSULATED_TIMESTAMP            XAdES111Element = "ENCAPSULATED_TIMESTAMP"
	XAdES111Element_ENCAPSULATED_X509_CERTIFICATE     XAdES111Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES111Element_ENCODING                          XAdES111Element = "ENCODING"
	XAdES111Element_EXPLICIT_TEXT                     XAdES111Element = "EXPLICIT_TEXT"
	XAdES111Element_HASH_DATA_INFO                    XAdES111Element = "HASH_DATA_INFO"
	XAdES111Element_IDENTIFIER                        XAdES111Element = "IDENTIFIER"
	XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP XAdES111Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES111Element_INT                               XAdES111Element = "INT"
	XAdES111Element_ISSUE_TIME                        XAdES111Element = "ISSUE_TIME"
	XAdES111Element_ISSUER                            XAdES111Element = "ISSUER"
	XAdES111Element_ISSUER_SERIAL                     XAdES111Element = "ISSUER_SERIAL"
	XAdES111Element_MIME_TYPE                         XAdES111Element = "MIME_TYPE"
	XAdES111Element_NOTICE_NUMBERS                    XAdES111Element = "NOTICE_NUMBERS"
	XAdES111Element_NOTICE_REF                        XAdES111Element = "NOTICE_REF"
	XAdES111Element_NUMBER                            XAdES111Element = "NUMBER"
	XAdES111Element_OBJECT_IDENTIFIER                 XAdES111Element = "OBJECT_IDENTIFIER"
	XAdES111Element_OBJECT_REFERENCE                  XAdES111Element = "OBJECT_REFERENCE"
	XAdES111Element_OCSP_IDENTIFIER                   XAdES111Element = "OCSP_IDENTIFIER"
	XAdES111Element_OCSP_REF                          XAdES111Element = "OCSP_REF"
	XAdES111Element_OCSP_REFS                         XAdES111Element = "OCSP_REFS"
	XAdES111Element_OCSP_VALUES                       XAdES111Element = "OCSP_VALUES"
	XAdES111Element_ORGANIZATION                      XAdES111Element = "ORGANIZATION"
	XAdES111Element_OTHER_CERTIFICATE                 XAdES111Element = "OTHER_CERTIFICATE"
	XAdES111Element_OTHER_REF                         XAdES111Element = "OTHER_REF"
	XAdES111Element_OTHER_REFS                        XAdES111Element = "OTHER_REFS"
	XAdES111Element_OTHER_VALUE                       XAdES111Element = "OTHER_VALUE"
	XAdES111Element_OTHER_VALUES                      XAdES111Element = "OTHER_VALUES"
	XAdES111Element_POSTAL_CODE                       XAdES111Element = "POSTAL_CODE"
	XAdES111Element_PRODUCED_AT                       XAdES111Element = "PRODUCED_AT"
	XAdES111Element_QUALIFYING_PROPERTIES             XAdES111Element = "QUALIFYING_PROPERTIES"
	XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE   XAdES111Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES111Element_REFS_ONLY_TIMESTAMP               XAdES111Element = "REFS_ONLY_TIMESTAMP"
	XAdES111Element_RESPONDER_ID                      XAdES111Element = "RESPONDER_ID"
	XAdES111Element_REVOCATION_VALUES                 XAdES111Element = "REVOCATION_VALUES"
	XAdES111Element_SIG_AND_REFS_TIMESTAMP            XAdES111Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES111Element_SIG_POLICY_HASH                   XAdES111Element = "SIG_POLICY_HASH"
	XAdES111Element_SIG_POLICY_ID                     XAdES111Element = "SIG_POLICY_ID"
	XAdES111Element_SIG_POLICY_QUALIFIER              XAdES111Element = "SIG_POLICY_QUALIFIER"
	XAdES111Element_SIG_POLICY_QUALIFIERS             XAdES111Element = "SIG_POLICY_QUALIFIERS"
	XAdES111Element_SIGNATURE_POLICY_ID               XAdES111Element = "SIGNATURE_POLICY_ID"
	XAdES111Element_SIGNATURE_POLICY_IDENTIFIER       XAdES111Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES111Element_SIGNATURE_POLICY_IMPLIED          XAdES111Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES111Element_SIGNATURE_PRODUCTION_PLACE        XAdES111Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES111Element_SIGNATURE_TIMESTAMP               XAdES111Element = "SIGNATURE_TIMESTAMP"
	XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES     XAdES111Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES111Element_SIGNED_PROPERTIES                 XAdES111Element = "SIGNED_PROPERTIES"
	XAdES111Element_SIGNED_SIGNATURE_PROPERTIES       XAdES111Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES111Element_SIGNER_ROLE                       XAdES111Element = "SIGNER_ROLE"
	XAdES111Element_SIGNING_CERTIFICATE               XAdES111Element = "SIGNING_CERTIFICATE"
	XAdES111Element_SIGNING_TIME                      XAdES111Element = "SIGNING_TIME"
	XAdES111Element_SP_URI                            XAdES111Element = "SP_URI"
	XAdES111Element_SP_USER_NOTICE                    XAdES111Element = "SP_USER_NOTICE"
	XAdES111Element_STATE_OR_PROVINCE                 XAdES111Element = "STATE_OR_PROVINCE"
	XAdES111Element_TIMESTAMP                         XAdES111Element = "TIMESTAMP"
	XAdES111Element_TRANSFORMS                        XAdES111Element = "TRANSFORMS"
	XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES   XAdES111Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY     XAdES111Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES111Element_UNSIGNED_PROPERTIES               XAdES111Element = "UNSIGNED_PROPERTIES"
	XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES     XAdES111Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES111Element_XML_TIMESTAMP                     XAdES111Element = "XML_TIMESTAMP"
)

// xades111elementTagNames maps each constant to its wire tag name (getTagName()).
var xades111elementTagNames = map[XAdES111Element]string{
	XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP:        "AllDataObjectsTimeStamp",
	XAdES111Element_ALL_SIGNED_DATA_OBJECTS:           "AllSignedDataObjects",
	XAdES111Element_ANY:                               "Any",
	XAdES111Element_ARCHIVE_TIMESTAMP:                 "ArchiveTimeStamp",
	XAdES111Element_CERT:                              "Cert",
	XAdES111Element_CERT_DIGEST:                       "CertDigest",
	XAdES111Element_CERT_REFS:                         "CertRefs",
	XAdES111Element_CERTIFICATE_VALUES:                "CertificateValues",
	XAdES111Element_CERTIFIED_ROLE:                    "CertifiedRole",
	XAdES111Element_CERTIFIED_ROLES:                   "CertifiedRoles",
	XAdES111Element_CITY:                              "City",
	XAdES111Element_CLAIMED_ROLE:                      "ClaimedRole",
	XAdES111Element_CLAIMED_ROLES:                     "ClaimedRoles",
	XAdES111Element_COMMITMENT_TYPE_ID:                "CommitmentTypeId",
	XAdES111Element_COMMITMENT_TYPE_INDICATION:        "CommitmentTypeIndication",
	XAdES111Element_COMMITMENT_TYPE_QUALIFIER:         "CommitmentTypeQualifier",
	XAdES111Element_COMMITMENT_TYPE_QUALIFIERS:        "CommitmentTypeQualifiers",
	XAdES111Element_COMPLETE_CERTIFICATE_REFS:         "CompleteCertificateRefs",
	XAdES111Element_COMPLETE_REVOCATION_REFS:          "CompleteRevocationRefs",
	XAdES111Element_COUNTER_SIGNATURE:                 "CounterSignature",
	XAdES111Element_COUNTRY_NAME:                      "CountryName",
	XAdES111Element_CRL_IDENTIFIER:                    "CRLIdentifier",
	XAdES111Element_CRL_REF:                           "CRLRef",
	XAdES111Element_CRL_REFS:                          "CRLRefs",
	XAdES111Element_CRL_VALUES:                        "CRLValues",
	XAdES111Element_DATA_OBJECT_FORMAT:                "DataObjectFormat",
	XAdES111Element_DESCRIPTION:                       "Description",
	XAdES111Element_DIGEST_ALG_AND_VALUE:              "DigestAlgAndValue",
	XAdES111Element_DIGEST_METHOD:                     "DigestMethod",
	XAdES111Element_DIGEST_VALUE:                      "DigestValue",
	XAdES111Element_DOCUMENTATION_REFERENCE:           "DocumentationReference",
	XAdES111Element_DOCUMENTATION_REFERENCES:          "DocumentationReferences",
	XAdES111Element_ENCAPSULATED_CRL_VALUE:            "EncapsulatedCRLValue",
	XAdES111Element_ENCAPSULATED_OCSP_VALUE:           "EncapsulatedOCSPValue",
	XAdES111Element_ENCAPSULATED_PKI_DATA:             "EncapsulatedPKIData",
	XAdES111Element_ENCAPSULATED_TIMESTAMP:            "EncapsulatedTimeStamp",
	XAdES111Element_ENCAPSULATED_X509_CERTIFICATE:     "EncapsulatedX509Certificate",
	XAdES111Element_ENCODING:                          "Encoding",
	XAdES111Element_EXPLICIT_TEXT:                     "ExplicitText",
	XAdES111Element_HASH_DATA_INFO:                    "HashDataInfo",
	XAdES111Element_IDENTIFIER:                        "Identifier",
	XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP: "IndividualDataObjectsTimeStamp",
	XAdES111Element_INT:                               "int",
	XAdES111Element_ISSUE_TIME:                        "IssueTime",
	XAdES111Element_ISSUER:                            "Issuer",
	XAdES111Element_ISSUER_SERIAL:                     "IssuerSerial",
	XAdES111Element_MIME_TYPE:                         "MimeType",
	XAdES111Element_NOTICE_NUMBERS:                    "NoticeNumbers",
	XAdES111Element_NOTICE_REF:                        "NoticeRef",
	XAdES111Element_NUMBER:                            "Number",
	XAdES111Element_OBJECT_IDENTIFIER:                 "ObjectIdentifier",
	XAdES111Element_OBJECT_REFERENCE:                  "ObjectReference",
	XAdES111Element_OCSP_IDENTIFIER:                   "OCSPIdentifier",
	XAdES111Element_OCSP_REF:                          "OCSPRef",
	XAdES111Element_OCSP_REFS:                         "OCSPRefs",
	XAdES111Element_OCSP_VALUES:                       "OCSPValues",
	XAdES111Element_ORGANIZATION:                      "Organization",
	XAdES111Element_OTHER_CERTIFICATE:                 "OtherCertificate",
	XAdES111Element_OTHER_REF:                         "OtherRef",
	XAdES111Element_OTHER_REFS:                        "OtherRefs",
	XAdES111Element_OTHER_VALUE:                       "OtherValue",
	XAdES111Element_OTHER_VALUES:                      "OtherValues",
	XAdES111Element_POSTAL_CODE:                       "PostalCode",
	XAdES111Element_PRODUCED_AT:                       "ProducedAt",
	XAdES111Element_QUALIFYING_PROPERTIES:             "QualifyingProperties",
	XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE:   "QualifyingPropertiesReference",
	XAdES111Element_REFS_ONLY_TIMESTAMP:               "RefsOnlyTimeStamp",
	XAdES111Element_RESPONDER_ID:                      "ResponderID",
	XAdES111Element_REVOCATION_VALUES:                 "RevocationValues",
	XAdES111Element_SIG_AND_REFS_TIMESTAMP:            "SigAndRefsTimeStamp",
	XAdES111Element_SIG_POLICY_HASH:                   "SigPolicyHash",
	XAdES111Element_SIG_POLICY_ID:                     "SigPolicyId",
	XAdES111Element_SIG_POLICY_QUALIFIER:              "SigPolicyQualifier",
	XAdES111Element_SIG_POLICY_QUALIFIERS:             "SigPolicyQualifiers",
	XAdES111Element_SIGNATURE_POLICY_ID:               "SignaturePolicyId",
	XAdES111Element_SIGNATURE_POLICY_IDENTIFIER:       "SignaturePolicyIdentifier",
	XAdES111Element_SIGNATURE_POLICY_IMPLIED:          "SignaturePolicyImplied",
	XAdES111Element_SIGNATURE_PRODUCTION_PLACE:        "SignatureProductionPlace",
	XAdES111Element_SIGNATURE_TIMESTAMP:               "SignatureTimeStamp",
	XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES:     "SignedDataObjectProperties",
	XAdES111Element_SIGNED_PROPERTIES:                 "SignedProperties",
	XAdES111Element_SIGNED_SIGNATURE_PROPERTIES:       "SignedSignatureProperties",
	XAdES111Element_SIGNER_ROLE:                       "SignerRole",
	XAdES111Element_SIGNING_CERTIFICATE:               "SigningCertificate",
	XAdES111Element_SIGNING_TIME:                      "SigningTime",
	XAdES111Element_SP_URI:                            "SPURI",
	XAdES111Element_SP_USER_NOTICE:                    "SPUserNotice",
	XAdES111Element_STATE_OR_PROVINCE:                 "StateOrProvince",
	XAdES111Element_TIMESTAMP:                         "TimeStamp",
	XAdES111Element_TRANSFORMS:                        "Transforms",
	XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES:   "UnsignedDataObjectProperties",
	XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY:     "UnsignedDataObjectProperty",
	XAdES111Element_UNSIGNED_PROPERTIES:               "UnsignedProperties",
	XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES:     "UnsignedSignatureProperties",
	XAdES111Element_XML_TIMESTAMP:                     "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES111Element) TagName() string {
	return xades111elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES111Element) Namespace() *common.DSSNamespace {
	return XAdESNamespace_XADES_111
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES111Element) URI() string {
	return XAdESNamespace_XADES_111.Uri()
}

// IsSameTagName implements common.DSSElement. Ports isSameTagName(String).
func (e XAdES111Element) IsSameTagName(value string) bool {
	return e.TagName() == value
}

// XAdES111ElementFromTagName returns the XAdES111Element constant with the given tag name, or nil
// if none matches. Ports the static factory fromTagName(String).
func XAdES111ElementFromTagName(tagName string) common.DSSElement {
	for _, e := range XAdES111ElementValues() {
		if e.TagName() == tagName {
			return e
		}
	}
	return nil
}

// XAdES111ElementValues returns every XAdES111Element constant, in declaration order. Ports values().
func XAdES111ElementValues() []XAdES111Element {
	return []XAdES111Element{
		XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP,
		XAdES111Element_ALL_SIGNED_DATA_OBJECTS,
		XAdES111Element_ANY,
		XAdES111Element_ARCHIVE_TIMESTAMP,
		XAdES111Element_CERT,
		XAdES111Element_CERT_DIGEST,
		XAdES111Element_CERT_REFS,
		XAdES111Element_CERTIFICATE_VALUES,
		XAdES111Element_CERTIFIED_ROLE,
		XAdES111Element_CERTIFIED_ROLES,
		XAdES111Element_CITY,
		XAdES111Element_CLAIMED_ROLE,
		XAdES111Element_CLAIMED_ROLES,
		XAdES111Element_COMMITMENT_TYPE_ID,
		XAdES111Element_COMMITMENT_TYPE_INDICATION,
		XAdES111Element_COMMITMENT_TYPE_QUALIFIER,
		XAdES111Element_COMMITMENT_TYPE_QUALIFIERS,
		XAdES111Element_COMPLETE_CERTIFICATE_REFS,
		XAdES111Element_COMPLETE_REVOCATION_REFS,
		XAdES111Element_COUNTER_SIGNATURE,
		XAdES111Element_COUNTRY_NAME,
		XAdES111Element_CRL_IDENTIFIER,
		XAdES111Element_CRL_REF,
		XAdES111Element_CRL_REFS,
		XAdES111Element_CRL_VALUES,
		XAdES111Element_DATA_OBJECT_FORMAT,
		XAdES111Element_DESCRIPTION,
		XAdES111Element_DIGEST_ALG_AND_VALUE,
		XAdES111Element_DIGEST_METHOD,
		XAdES111Element_DIGEST_VALUE,
		XAdES111Element_DOCUMENTATION_REFERENCE,
		XAdES111Element_DOCUMENTATION_REFERENCES,
		XAdES111Element_ENCAPSULATED_CRL_VALUE,
		XAdES111Element_ENCAPSULATED_OCSP_VALUE,
		XAdES111Element_ENCAPSULATED_PKI_DATA,
		XAdES111Element_ENCAPSULATED_TIMESTAMP,
		XAdES111Element_ENCAPSULATED_X509_CERTIFICATE,
		XAdES111Element_ENCODING,
		XAdES111Element_EXPLICIT_TEXT,
		XAdES111Element_HASH_DATA_INFO,
		XAdES111Element_IDENTIFIER,
		XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP,
		XAdES111Element_INT,
		XAdES111Element_ISSUE_TIME,
		XAdES111Element_ISSUER,
		XAdES111Element_ISSUER_SERIAL,
		XAdES111Element_MIME_TYPE,
		XAdES111Element_NOTICE_NUMBERS,
		XAdES111Element_NOTICE_REF,
		XAdES111Element_NUMBER,
		XAdES111Element_OBJECT_IDENTIFIER,
		XAdES111Element_OBJECT_REFERENCE,
		XAdES111Element_OCSP_IDENTIFIER,
		XAdES111Element_OCSP_REF,
		XAdES111Element_OCSP_REFS,
		XAdES111Element_OCSP_VALUES,
		XAdES111Element_ORGANIZATION,
		XAdES111Element_OTHER_CERTIFICATE,
		XAdES111Element_OTHER_REF,
		XAdES111Element_OTHER_REFS,
		XAdES111Element_OTHER_VALUE,
		XAdES111Element_OTHER_VALUES,
		XAdES111Element_POSTAL_CODE,
		XAdES111Element_PRODUCED_AT,
		XAdES111Element_QUALIFYING_PROPERTIES,
		XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE,
		XAdES111Element_REFS_ONLY_TIMESTAMP,
		XAdES111Element_RESPONDER_ID,
		XAdES111Element_REVOCATION_VALUES,
		XAdES111Element_SIG_AND_REFS_TIMESTAMP,
		XAdES111Element_SIG_POLICY_HASH,
		XAdES111Element_SIG_POLICY_ID,
		XAdES111Element_SIG_POLICY_QUALIFIER,
		XAdES111Element_SIG_POLICY_QUALIFIERS,
		XAdES111Element_SIGNATURE_POLICY_ID,
		XAdES111Element_SIGNATURE_POLICY_IDENTIFIER,
		XAdES111Element_SIGNATURE_POLICY_IMPLIED,
		XAdES111Element_SIGNATURE_PRODUCTION_PLACE,
		XAdES111Element_SIGNATURE_TIMESTAMP,
		XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES,
		XAdES111Element_SIGNED_PROPERTIES,
		XAdES111Element_SIGNED_SIGNATURE_PROPERTIES,
		XAdES111Element_SIGNER_ROLE,
		XAdES111Element_SIGNING_CERTIFICATE,
		XAdES111Element_SIGNING_TIME,
		XAdES111Element_SP_URI,
		XAdES111Element_SP_USER_NOTICE,
		XAdES111Element_STATE_OR_PROVINCE,
		XAdES111Element_TIMESTAMP,
		XAdES111Element_TRANSFORMS,
		XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES,
		XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY,
		XAdES111Element_UNSIGNED_PROPERTIES,
		XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES,
		XAdES111Element_XML_TIMESTAMP,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES111Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES111Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES111Element_ALL_SIGNED_DATA_OBJECTS
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES111Element) ElementAny() common.DSSElement {
	return XAdES111Element_ANY
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES111Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES111Element_ARCHIVE_TIMESTAMP
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES111Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES111Element) ElementAttributeCertificateRefs() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES111Element) ElementAttributeRevocationRefs() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES111Element) ElementAttributeRevocationValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES111Element) ElementByKey() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES111Element) ElementByName() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES111Element) ElementCert() common.DSSElement {
	return XAdES111Element_CERT
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES111Element) ElementCertDigest() common.DSSElement {
	return XAdES111Element_CERT_DIGEST
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES111Element) ElementCertRefs() common.DSSElement {
	return XAdES111Element_CERT_REFS
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES111Element) ElementCertificateValues() common.DSSElement {
	return XAdES111Element_CERTIFICATE_VALUES
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES111Element) ElementCertifiedRole() common.DSSElement {
	return XAdES111Element_CERTIFIED_ROLE
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES111Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES111Element_CERTIFIED_ROLES
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES111Element) ElementCertifiedRolesV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES111Element) ElementCity() common.DSSElement {
	return XAdES111Element_CITY
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES111Element) ElementClaimedRole() common.DSSElement {
	return XAdES111Element_CLAIMED_ROLE
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES111Element) ElementClaimedRoles() common.DSSElement {
	return XAdES111Element_CLAIMED_ROLES
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES111Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES111Element_COMMITMENT_TYPE_ID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES111Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES111Element_COMMITMENT_TYPE_INDICATION
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES111Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES111Element_COMMITMENT_TYPE_QUALIFIER
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES111Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES111Element_COMMITMENT_TYPE_QUALIFIERS
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES111Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES111Element_COMPLETE_CERTIFICATE_REFS
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES111Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES111Element_COMPLETE_REVOCATION_REFS
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES111Element) ElementCounterSignature() common.DSSElement {
	return XAdES111Element_COUNTER_SIGNATURE
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES111Element) ElementCountryName() common.DSSElement {
	return XAdES111Element_COUNTRY_NAME
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES111Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES111Element_CRL_IDENTIFIER
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES111Element) ElementCRLRef() common.DSSElement {
	return XAdES111Element_CRL_REF
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES111Element) ElementCRLRefs() common.DSSElement {
	return XAdES111Element_CRL_REFS
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES111Element) ElementCRLValues() common.DSSElement {
	return XAdES111Element_CRL_VALUES
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES111Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES111Element_DATA_OBJECT_FORMAT
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES111Element) ElementDescription() common.DSSElement {
	return XAdES111Element_DESCRIPTION
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES111Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES111Element_DIGEST_ALG_AND_VALUE
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES111Element) ElementDocumentationReference() common.DSSElement {
	return XAdES111Element_DOCUMENTATION_REFERENCE
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES111Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES111Element_DOCUMENTATION_REFERENCES
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES111Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES111Element_ENCAPSULATED_CRL_VALUE
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES111Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES111Element_ENCAPSULATED_OCSP_VALUE
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES111Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES111Element_ENCAPSULATED_PKI_DATA
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES111Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES111Element_ENCAPSULATED_TIMESTAMP
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES111Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES111Element_ENCAPSULATED_X509_CERTIFICATE
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES111Element) ElementEncoding() common.DSSElement {
	return XAdES111Element_ENCODING
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES111Element) ElementExplicitText() common.DSSElement {
	return XAdES111Element_EXPLICIT_TEXT
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES111Element) ElementIdentifier() common.DSSElement {
	return XAdES111Element_IDENTIFIER
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES111Element) ElementInclude() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES111Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES111Element) Elementint() common.DSSElement {
	return XAdES111Element_INT
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES111Element) ElementIssueTime() common.DSSElement {
	return XAdES111Element_ISSUE_TIME
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES111Element) ElementIssuer() common.DSSElement {
	return XAdES111Element_ISSUER
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES111Element) ElementIssuerSerial() common.DSSElement {
	return XAdES111Element_ISSUER_SERIAL
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES111Element) ElementIssuerSerialV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES111Element) ElementMimeType() common.DSSElement {
	return XAdES111Element_MIME_TYPE
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES111Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES111Element_NOTICE_NUMBERS
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES111Element) ElementNoticeRef() common.DSSElement {
	return XAdES111Element_NOTICE_REF
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES111Element) ElementNumber() common.DSSElement {
	return XAdES111Element_NUMBER
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES111Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES111Element_OBJECT_IDENTIFIER
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES111Element) ElementObjectReference() common.DSSElement {
	return XAdES111Element_OBJECT_REFERENCE
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES111Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES111Element_OCSP_IDENTIFIER
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES111Element) ElementOCSPRef() common.DSSElement {
	return XAdES111Element_OCSP_REF
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES111Element) ElementOCSPRefs() common.DSSElement {
	return XAdES111Element_OCSP_REFS
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES111Element) ElementOCSPValues() common.DSSElement {
	return XAdES111Element_OCSP_VALUES
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES111Element) ElementOrganization() common.DSSElement {
	return XAdES111Element_ORGANIZATION
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES111Element) ElementOtherAttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES111Element) ElementOtherCertificate() common.DSSElement {
	return XAdES111Element_OTHER_CERTIFICATE
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES111Element) ElementOtherRef() common.DSSElement {
	return XAdES111Element_OTHER_REF
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES111Element) ElementOtherRefs() common.DSSElement {
	return XAdES111Element_OTHER_REFS
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES111Element) ElementOtherTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES111Element) ElementOtherValue() common.DSSElement {
	return XAdES111Element_OTHER_VALUE
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES111Element) ElementOtherValues() common.DSSElement {
	return XAdES111Element_OTHER_VALUES
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES111Element) ElementPostalCode() common.DSSElement {
	return XAdES111Element_POSTAL_CODE
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES111Element) ElementProducedAt() common.DSSElement {
	return XAdES111Element_PRODUCED_AT
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES111Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES111Element_QUALIFYING_PROPERTIES
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES111Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES111Element) ElementReferenceInfo() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES111Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES111Element_REFS_ONLY_TIMESTAMP
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES111Element) ElementResponderID() common.DSSElement {
	return XAdES111Element_RESPONDER_ID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES111Element) ElementRevocationValues() common.DSSElement {
	return XAdES111Element_REVOCATION_VALUES
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES111Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES111Element_SIG_AND_REFS_TIMESTAMP
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES111Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES111Element_SIG_POLICY_HASH
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES111Element) ElementSigPolicyId() common.DSSElement {
	return XAdES111Element_SIG_POLICY_ID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES111Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES111Element_SIG_POLICY_QUALIFIER
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES111Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES111Element_SIG_POLICY_QUALIFIERS
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES111Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES111Element_SIGNATURE_POLICY_ID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES111Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES111Element_SIGNATURE_POLICY_IDENTIFIER
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES111Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES111Element_SIGNATURE_POLICY_IMPLIED
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES111Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES111Element_SIGNATURE_PRODUCTION_PLACE
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES111Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES111Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES111Element_SIGNATURE_TIMESTAMP
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES111Element) ElementSignedAssertion() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES111Element) ElementSignedAssertions() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES111Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES111Element) ElementSignedProperties() common.DSSElement {
	return XAdES111Element_SIGNED_PROPERTIES
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES111Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES111Element_SIGNED_SIGNATURE_PROPERTIES
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES111Element) ElementSignerRole() common.DSSElement {
	return XAdES111Element_SIGNER_ROLE
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES111Element) ElementSignerRoleV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES111Element) ElementSigningCertificate() common.DSSElement {
	return XAdES111Element_SIGNING_CERTIFICATE
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES111Element) ElementSigningCertificateV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES111Element) ElementSigningTime() common.DSSElement {
	return XAdES111Element_SIGNING_TIME
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES111Element) ElementSPURI() common.DSSElement {
	return XAdES111Element_SP_URI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES111Element) ElementSPUserNotice() common.DSSElement {
	return XAdES111Element_SP_USER_NOTICE
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES111Element) ElementStateOrProvince() common.DSSElement {
	return XAdES111Element_STATE_OR_PROVINCE
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES111Element) ElementStreetAddress() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES111Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES111Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES111Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES111Element_UNSIGNED_PROPERTIES
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES111Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES111Element) ElementX509AttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES111Element) ElementXAdESTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespace_XADES_111.Uri())
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES111Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES111Element_XML_TIMESTAMP
}

var _ XAdESElement = XAdES111Element("")
