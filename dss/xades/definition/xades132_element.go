// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades132/XAdES132Element.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES132Element defines elements for a XAdES 1.3.2 schema.
type XAdES132Element string

// XAdES132Element constants, one per XAdES 1.3.2 schema element name.
const (
	XAdES132ElementAllDataObjectsTimestamp        XAdES132Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES132ElementAllSignedDataObjects           XAdES132Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES132ElementAny                            XAdES132Element = "ANY"
	XAdES132ElementArchiveTimestamp               XAdES132Element = "ARCHIVE_TIMESTAMP"
	XAdES132ElementAttrAuthoritiesCertValues      XAdES132Element = "ATTR_AUTHORITIES_CERT_VALUES"
	XAdES132ElementAttributeCertificateRefs       XAdES132Element = "ATTRIBUTE_CERTIFICATE_REFS"
	XAdES132ElementAttributeRevocationRefs        XAdES132Element = "ATTRIBUTE_REVOCATION_REFS"
	XAdES132ElementAttributeRevocationValues      XAdES132Element = "ATTRIBUTE_REVOCATION_VALUES"
	XAdES132ElementByKey                          XAdES132Element = "BY_KEY"
	XAdES132ElementByName                         XAdES132Element = "BY_NAME"
	XAdES132ElementCert                           XAdES132Element = "CERT"
	XAdES132ElementCertDigest                     XAdES132Element = "CERT_DIGEST"
	XAdES132ElementCertRefs                       XAdES132Element = "CERT_REFS"
	XAdES132ElementCertificateValues              XAdES132Element = "CERTIFICATE_VALUES"
	XAdES132ElementCertifiedRole                  XAdES132Element = "CERTIFIED_ROLE"
	XAdES132ElementCertifiedRoles                 XAdES132Element = "CERTIFIED_ROLES"
	XAdES132ElementCertifiedRolesV2               XAdES132Element = "CERTIFIED_ROLES_V2"
	XAdES132ElementCity                           XAdES132Element = "CITY"
	XAdES132ElementClaimedRole                    XAdES132Element = "CLAIMED_ROLE"
	XAdES132ElementClaimedRoles                   XAdES132Element = "CLAIMED_ROLES"
	XAdES132ElementCommitmentTypeID               XAdES132Element = "COMMITMENT_TYPE_ID"
	XAdES132ElementCommitmentTypeIndication       XAdES132Element = "COMMITMENT_TYPE_INDICATION"
	XAdES132ElementCommitmentTypeQualifier        XAdES132Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES132ElementCommitmentTypeQualifiers       XAdES132Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES132ElementCompleteCertificateRefs        XAdES132Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES132ElementCompleteRevocationRefs         XAdES132Element = "COMPLETE_REVOCATION_REFS"
	XAdES132ElementCounterSignature               XAdES132Element = "COUNTER_SIGNATURE"
	XAdES132ElementCountryName                    XAdES132Element = "COUNTRY_NAME"
	XAdES132ElementCRLIdentifier                  XAdES132Element = "CRL_IDENTIFIER"
	XAdES132ElementCRLRef                         XAdES132Element = "CRL_REF"
	XAdES132ElementCRLRefs                        XAdES132Element = "CRL_REFS"
	XAdES132ElementCRLValues                      XAdES132Element = "CRL_VALUES"
	XAdES132ElementDataObjectFormat               XAdES132Element = "DATA_OBJECT_FORMAT"
	XAdES132ElementDescription                    XAdES132Element = "DESCRIPTION"
	XAdES132ElementDigestAlgAndValue              XAdES132Element = "DIGEST_ALG_AND_VALUE"
	XAdES132ElementDocumentationReference         XAdES132Element = "DOCUMENTATION_REFERENCE"
	XAdES132ElementDocumentationReferences        XAdES132Element = "DOCUMENTATION_REFERENCES"
	XAdES132ElementEncapsulatedCRLValue           XAdES132Element = "ENCAPSULATED_CRL_VALUE"
	XAdES132ElementEncapsulatedOCSPValue          XAdES132Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES132ElementEncapsulatedPKIData            XAdES132Element = "ENCAPSULATED_PKI_DATA"
	XAdES132ElementEncapsulatedTimestamp          XAdES132Element = "ENCAPSULATED_TIMESTAMP"
	XAdES132ElementEncapsulatedX509Certificate    XAdES132Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES132ElementEncoding                       XAdES132Element = "ENCODING"
	XAdES132ElementExplicitText                   XAdES132Element = "EXPLICIT_TEXT"
	XAdES132ElementIdentifier                     XAdES132Element = "IDENTIFIER"
	XAdES132ElementInclude                        XAdES132Element = "INCLUDE"
	XAdES132ElementIndividualDataObjectsTimestamp XAdES132Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES132ElementInt                            XAdES132Element = "INT"
	XAdES132ElementIssueTime                      XAdES132Element = "ISSUE_TIME"
	XAdES132ElementIssuer                         XAdES132Element = "ISSUER"
	XAdES132ElementIssuerSerial                   XAdES132Element = "ISSUER_SERIAL"
	XAdES132ElementIssuerSerialV2                 XAdES132Element = "ISSUER_SERIAL_V2"
	XAdES132ElementMIMEType                       XAdES132Element = "MIME_TYPE"
	XAdES132ElementNoticeNumbers                  XAdES132Element = "NOTICE_NUMBERS"
	XAdES132ElementNoticeRef                      XAdES132Element = "NOTICE_REF"
	XAdES132ElementNumber                         XAdES132Element = "NUMBER"
	XAdES132ElementObjectIdentifier               XAdES132Element = "OBJECT_IDENTIFIER"
	XAdES132ElementObjectReference                XAdES132Element = "OBJECT_REFERENCE"
	XAdES132ElementOCSPIdentifier                 XAdES132Element = "OCSP_IDENTIFIER"
	XAdES132ElementOCSPRef                        XAdES132Element = "OCSP_REF"
	XAdES132ElementOCSPRefs                       XAdES132Element = "OCSP_REFS"
	XAdES132ElementOCSPValues                     XAdES132Element = "OCSP_VALUES"
	XAdES132ElementOrganization                   XAdES132Element = "ORGANIZATION"
	XAdES132ElementOtherAttributeCertificate      XAdES132Element = "OTHER_ATTRIBUTE_CERTIFICATE"
	XAdES132ElementOtherCertificate               XAdES132Element = "OTHER_CERTIFICATE"
	XAdES132ElementOtherRef                       XAdES132Element = "OTHER_REF"
	XAdES132ElementOtherRefs                      XAdES132Element = "OTHER_REFS"
	XAdES132ElementOtherTimestamp                 XAdES132Element = "OTHER_TIMESTAMP"
	XAdES132ElementOtherValue                     XAdES132Element = "OTHER_VALUE"
	XAdES132ElementOtherValues                    XAdES132Element = "OTHER_VALUES"
	XAdES132ElementPostalCode                     XAdES132Element = "POSTAL_CODE"
	XAdES132ElementProducedAt                     XAdES132Element = "PRODUCED_AT"
	XAdES132ElementQualifyingProperties           XAdES132Element = "QUALIFYING_PROPERTIES"
	XAdES132ElementQualifyingPropertiesReference  XAdES132Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES132ElementReferenceInfo                  XAdES132Element = "REFERENCE_INFO"
	XAdES132ElementRefsOnlyTimestamp              XAdES132Element = "REFS_ONLY_TIMESTAMP"
	XAdES132ElementResponderID                    XAdES132Element = "RESPONDER_ID"
	XAdES132ElementRevocationValues               XAdES132Element = "REVOCATION_VALUES"
	XAdES132ElementSigAndRefsTimestamp            XAdES132Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES132ElementSigPolicyHash                  XAdES132Element = "SIG_POLICY_HASH"
	XAdES132ElementSigPolicyID                    XAdES132Element = "SIG_POLICY_ID"
	XAdES132ElementSigPolicyQualifier             XAdES132Element = "SIG_POLICY_QUALIFIER"
	XAdES132ElementSigPolicyQualifiers            XAdES132Element = "SIG_POLICY_QUALIFIERS"
	XAdES132ElementSignaturePolicyID              XAdES132Element = "SIGNATURE_POLICY_ID"
	XAdES132ElementSignaturePolicyIdentifier      XAdES132Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES132ElementSignaturePolicyImplied         XAdES132Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES132ElementSignatureProductionPlace       XAdES132Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES132ElementSignatureProductionPlaceV2     XAdES132Element = "SIGNATURE_PRODUCTION_PLACE_V2"
	XAdES132ElementSignatureTimestamp             XAdES132Element = "SIGNATURE_TIMESTAMP"
	XAdES132ElementSignedAssertion                XAdES132Element = "SIGNED_ASSERTION"
	XAdES132ElementSignedAssertions               XAdES132Element = "SIGNED_ASSERTIONS"
	XAdES132ElementSignedDataObjectProperties     XAdES132Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES132ElementSignedProperties               XAdES132Element = "SIGNED_PROPERTIES"
	XAdES132ElementSignedSignatureProperties      XAdES132Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES132ElementSignerRole                     XAdES132Element = "SIGNER_ROLE"
	XAdES132ElementSignerRoleV2                   XAdES132Element = "SIGNER_ROLE_V2"
	XAdES132ElementSigningCertificate             XAdES132Element = "SIGNING_CERTIFICATE"
	XAdES132ElementSigningCertificateV2           XAdES132Element = "SIGNING_CERTIFICATE_V2"
	XAdES132ElementSigningTime                    XAdES132Element = "SIGNING_TIME"
	XAdES132ElementSPURI                          XAdES132Element = "SP_URI"
	XAdES132ElementSPUserNotice                   XAdES132Element = "SP_USER_NOTICE"
	XAdES132ElementStateOrProvince                XAdES132Element = "STATE_OR_PROVINCE"
	XAdES132ElementStreetAddress                  XAdES132Element = "STREET_ADDRESS"
	XAdES132ElementUnsignedDataObjectProperties   XAdES132Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES132ElementUnsignedDataObjectProperty     XAdES132Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES132ElementUnsignedProperties             XAdES132Element = "UNSIGNED_PROPERTIES"
	XAdES132ElementUnsignedSignatureProperties    XAdES132Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES132ElementX509AttributeCertificate       XAdES132Element = "X509_ATTRIBUTE_CERTIFICATE"
	XAdES132ElementXAdESTimestamp                 XAdES132Element = "XADES_TIMESTAMP"
	XAdES132ElementXMLTimestamp                   XAdES132Element = "XML_TIMESTAMP"
)

// xades132elementTagNames maps each constant to its wire tag name (getTagName()).
var xades132elementTagNames = map[XAdES132Element]string{
	XAdES132ElementAllDataObjectsTimestamp:        "AllDataObjectsTimeStamp",
	XAdES132ElementAllSignedDataObjects:           "AllSignedDataObjects",
	XAdES132ElementAny:                            "Any",
	XAdES132ElementArchiveTimestamp:               "ArchiveTimeStamp",
	XAdES132ElementAttrAuthoritiesCertValues:      "AttrAuthoritiesCertValues",
	XAdES132ElementAttributeCertificateRefs:       "AttributeCertificateRefs",
	XAdES132ElementAttributeRevocationRefs:        "AttributeRevocationRefs",
	XAdES132ElementAttributeRevocationValues:      "AttributeRevocationValues",
	XAdES132ElementByKey:                          "ByKey",
	XAdES132ElementByName:                         "ByName",
	XAdES132ElementCert:                           "Cert",
	XAdES132ElementCertDigest:                     "CertDigest",
	XAdES132ElementCertRefs:                       "CertRefs",
	XAdES132ElementCertificateValues:              "CertificateValues",
	XAdES132ElementCertifiedRole:                  "CertifiedRole",
	XAdES132ElementCertifiedRoles:                 "CertifiedRoles",
	XAdES132ElementCertifiedRolesV2:               "CertifiedRolesV2",
	XAdES132ElementCity:                           "City",
	XAdES132ElementClaimedRole:                    "ClaimedRole",
	XAdES132ElementClaimedRoles:                   "ClaimedRoles",
	XAdES132ElementCommitmentTypeID:               "CommitmentTypeId",
	XAdES132ElementCommitmentTypeIndication:       "CommitmentTypeIndication",
	XAdES132ElementCommitmentTypeQualifier:        "CommitmentTypeQualifier",
	XAdES132ElementCommitmentTypeQualifiers:       "CommitmentTypeQualifiers",
	XAdES132ElementCompleteCertificateRefs:        "CompleteCertificateRefs",
	XAdES132ElementCompleteRevocationRefs:         "CompleteRevocationRefs",
	XAdES132ElementCounterSignature:               "CounterSignature",
	XAdES132ElementCountryName:                    "CountryName",
	XAdES132ElementCRLIdentifier:                  "CRLIdentifier",
	XAdES132ElementCRLRef:                         "CRLRef",
	XAdES132ElementCRLRefs:                        "CRLRefs",
	XAdES132ElementCRLValues:                      "CRLValues",
	XAdES132ElementDataObjectFormat:               "DataObjectFormat",
	XAdES132ElementDescription:                    "Description",
	XAdES132ElementDigestAlgAndValue:              "DigestAlgAndValue",
	XAdES132ElementDocumentationReference:         "DocumentationReference",
	XAdES132ElementDocumentationReferences:        "DocumentationReferences",
	XAdES132ElementEncapsulatedCRLValue:           "EncapsulatedCRLValue",
	XAdES132ElementEncapsulatedOCSPValue:          "EncapsulatedOCSPValue",
	XAdES132ElementEncapsulatedPKIData:            "EncapsulatedPKIData",
	XAdES132ElementEncapsulatedTimestamp:          "EncapsulatedTimeStamp",
	XAdES132ElementEncapsulatedX509Certificate:    "EncapsulatedX509Certificate",
	XAdES132ElementEncoding:                       "Encoding",
	XAdES132ElementExplicitText:                   "ExplicitText",
	XAdES132ElementIdentifier:                     "Identifier",
	XAdES132ElementInclude:                        "Include",
	XAdES132ElementIndividualDataObjectsTimestamp: "IndividualDataObjectsTimeStamp",
	XAdES132ElementInt:                            "int",
	XAdES132ElementIssueTime:                      "IssueTime",
	XAdES132ElementIssuer:                         "Issuer",
	XAdES132ElementIssuerSerial:                   "IssuerSerial",
	XAdES132ElementIssuerSerialV2:                 "IssuerSerialV2",
	XAdES132ElementMIMEType:                       "MimeType",
	XAdES132ElementNoticeNumbers:                  "NoticeNumbers",
	XAdES132ElementNoticeRef:                      "NoticeRef",
	XAdES132ElementNumber:                         "Number",
	XAdES132ElementObjectIdentifier:               "ObjectIdentifier",
	XAdES132ElementObjectReference:                "ObjectReference",
	XAdES132ElementOCSPIdentifier:                 "OCSPIdentifier",
	XAdES132ElementOCSPRef:                        "OCSPRef",
	XAdES132ElementOCSPRefs:                       "OCSPRefs",
	XAdES132ElementOCSPValues:                     "OCSPValues",
	XAdES132ElementOrganization:                   "Organization",
	XAdES132ElementOtherAttributeCertificate:      "OtherAttributeCertificate",
	XAdES132ElementOtherCertificate:               "OtherCertificate",
	XAdES132ElementOtherRef:                       "OtherRef",
	XAdES132ElementOtherRefs:                      "OtherRefs",
	XAdES132ElementOtherTimestamp:                 "OtherTimeStamp",
	XAdES132ElementOtherValue:                     "OtherValue",
	XAdES132ElementOtherValues:                    "OtherValues",
	XAdES132ElementPostalCode:                     "PostalCode",
	XAdES132ElementProducedAt:                     "ProducedAt",
	XAdES132ElementQualifyingProperties:           "QualifyingProperties",
	XAdES132ElementQualifyingPropertiesReference:  "QualifyingPropertiesReference",
	XAdES132ElementReferenceInfo:                  "ReferenceInfo",
	XAdES132ElementRefsOnlyTimestamp:              "RefsOnlyTimeStamp",
	XAdES132ElementResponderID:                    "ResponderID",
	XAdES132ElementRevocationValues:               "RevocationValues",
	XAdES132ElementSigAndRefsTimestamp:            "SigAndRefsTimeStamp",
	XAdES132ElementSigPolicyHash:                  "SigPolicyHash",
	XAdES132ElementSigPolicyID:                    "SigPolicyId",
	XAdES132ElementSigPolicyQualifier:             "SigPolicyQualifier",
	XAdES132ElementSigPolicyQualifiers:            "SigPolicyQualifiers",
	XAdES132ElementSignaturePolicyID:              "SignaturePolicyId",
	XAdES132ElementSignaturePolicyIdentifier:      "SignaturePolicyIdentifier",
	XAdES132ElementSignaturePolicyImplied:         "SignaturePolicyImplied",
	XAdES132ElementSignatureProductionPlace:       "SignatureProductionPlace",
	XAdES132ElementSignatureProductionPlaceV2:     "SignatureProductionPlaceV2",
	XAdES132ElementSignatureTimestamp:             "SignatureTimeStamp",
	XAdES132ElementSignedAssertion:                "SignedAssertion",
	XAdES132ElementSignedAssertions:               "SignedAssertions",
	XAdES132ElementSignedDataObjectProperties:     "SignedDataObjectProperties",
	XAdES132ElementSignedProperties:               "SignedProperties",
	XAdES132ElementSignedSignatureProperties:      "SignedSignatureProperties",
	XAdES132ElementSignerRole:                     "SignerRole",
	XAdES132ElementSignerRoleV2:                   "SignerRoleV2",
	XAdES132ElementSigningCertificate:             "SigningCertificate",
	XAdES132ElementSigningCertificateV2:           "SigningCertificateV2",
	XAdES132ElementSigningTime:                    "SigningTime",
	XAdES132ElementSPURI:                          "SPURI",
	XAdES132ElementSPUserNotice:                   "SPUserNotice",
	XAdES132ElementStateOrProvince:                "StateOrProvince",
	XAdES132ElementStreetAddress:                  "StreetAddress",
	XAdES132ElementUnsignedDataObjectProperties:   "UnsignedDataObjectProperties",
	XAdES132ElementUnsignedDataObjectProperty:     "UnsignedDataObjectProperty",
	XAdES132ElementUnsignedProperties:             "UnsignedProperties",
	XAdES132ElementUnsignedSignatureProperties:    "UnsignedSignatureProperties",
	XAdES132ElementX509AttributeCertificate:       "X509AttributeCertificate",
	XAdES132ElementXAdESTimestamp:                 "XAdESTimeStamp",
	XAdES132ElementXMLTimestamp:                   "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES132Element) TagName() string {
	return xades132elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES132Element) Namespace() *common.DSSNamespace {
	return XAdESNamespaceXAdES132
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES132Element) URI() string {
	return XAdESNamespaceXAdES132.Uri()
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
		XAdES132ElementAllDataObjectsTimestamp,
		XAdES132ElementAllSignedDataObjects,
		XAdES132ElementAny,
		XAdES132ElementArchiveTimestamp,
		XAdES132ElementAttrAuthoritiesCertValues,
		XAdES132ElementAttributeCertificateRefs,
		XAdES132ElementAttributeRevocationRefs,
		XAdES132ElementAttributeRevocationValues,
		XAdES132ElementByKey,
		XAdES132ElementByName,
		XAdES132ElementCert,
		XAdES132ElementCertDigest,
		XAdES132ElementCertRefs,
		XAdES132ElementCertificateValues,
		XAdES132ElementCertifiedRole,
		XAdES132ElementCertifiedRoles,
		XAdES132ElementCertifiedRolesV2,
		XAdES132ElementCity,
		XAdES132ElementClaimedRole,
		XAdES132ElementClaimedRoles,
		XAdES132ElementCommitmentTypeID,
		XAdES132ElementCommitmentTypeIndication,
		XAdES132ElementCommitmentTypeQualifier,
		XAdES132ElementCommitmentTypeQualifiers,
		XAdES132ElementCompleteCertificateRefs,
		XAdES132ElementCompleteRevocationRefs,
		XAdES132ElementCounterSignature,
		XAdES132ElementCountryName,
		XAdES132ElementCRLIdentifier,
		XAdES132ElementCRLRef,
		XAdES132ElementCRLRefs,
		XAdES132ElementCRLValues,
		XAdES132ElementDataObjectFormat,
		XAdES132ElementDescription,
		XAdES132ElementDigestAlgAndValue,
		XAdES132ElementDocumentationReference,
		XAdES132ElementDocumentationReferences,
		XAdES132ElementEncapsulatedCRLValue,
		XAdES132ElementEncapsulatedOCSPValue,
		XAdES132ElementEncapsulatedPKIData,
		XAdES132ElementEncapsulatedTimestamp,
		XAdES132ElementEncapsulatedX509Certificate,
		XAdES132ElementEncoding,
		XAdES132ElementExplicitText,
		XAdES132ElementIdentifier,
		XAdES132ElementInclude,
		XAdES132ElementIndividualDataObjectsTimestamp,
		XAdES132ElementInt,
		XAdES132ElementIssueTime,
		XAdES132ElementIssuer,
		XAdES132ElementIssuerSerial,
		XAdES132ElementIssuerSerialV2,
		XAdES132ElementMIMEType,
		XAdES132ElementNoticeNumbers,
		XAdES132ElementNoticeRef,
		XAdES132ElementNumber,
		XAdES132ElementObjectIdentifier,
		XAdES132ElementObjectReference,
		XAdES132ElementOCSPIdentifier,
		XAdES132ElementOCSPRef,
		XAdES132ElementOCSPRefs,
		XAdES132ElementOCSPValues,
		XAdES132ElementOrganization,
		XAdES132ElementOtherAttributeCertificate,
		XAdES132ElementOtherCertificate,
		XAdES132ElementOtherRef,
		XAdES132ElementOtherRefs,
		XAdES132ElementOtherTimestamp,
		XAdES132ElementOtherValue,
		XAdES132ElementOtherValues,
		XAdES132ElementPostalCode,
		XAdES132ElementProducedAt,
		XAdES132ElementQualifyingProperties,
		XAdES132ElementQualifyingPropertiesReference,
		XAdES132ElementReferenceInfo,
		XAdES132ElementRefsOnlyTimestamp,
		XAdES132ElementResponderID,
		XAdES132ElementRevocationValues,
		XAdES132ElementSigAndRefsTimestamp,
		XAdES132ElementSigPolicyHash,
		XAdES132ElementSigPolicyID,
		XAdES132ElementSigPolicyQualifier,
		XAdES132ElementSigPolicyQualifiers,
		XAdES132ElementSignaturePolicyID,
		XAdES132ElementSignaturePolicyIdentifier,
		XAdES132ElementSignaturePolicyImplied,
		XAdES132ElementSignatureProductionPlace,
		XAdES132ElementSignatureProductionPlaceV2,
		XAdES132ElementSignatureTimestamp,
		XAdES132ElementSignedAssertion,
		XAdES132ElementSignedAssertions,
		XAdES132ElementSignedDataObjectProperties,
		XAdES132ElementSignedProperties,
		XAdES132ElementSignedSignatureProperties,
		XAdES132ElementSignerRole,
		XAdES132ElementSignerRoleV2,
		XAdES132ElementSigningCertificate,
		XAdES132ElementSigningCertificateV2,
		XAdES132ElementSigningTime,
		XAdES132ElementSPURI,
		XAdES132ElementSPUserNotice,
		XAdES132ElementStateOrProvince,
		XAdES132ElementStreetAddress,
		XAdES132ElementUnsignedDataObjectProperties,
		XAdES132ElementUnsignedDataObjectProperty,
		XAdES132ElementUnsignedProperties,
		XAdES132ElementUnsignedSignatureProperties,
		XAdES132ElementX509AttributeCertificate,
		XAdES132ElementXAdESTimestamp,
		XAdES132ElementXMLTimestamp,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES132Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES132ElementAllDataObjectsTimestamp
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES132Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES132ElementAllSignedDataObjects
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES132Element) ElementAny() common.DSSElement {
	return XAdES132ElementAny
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES132Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES132ElementArchiveTimestamp
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES132Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	return XAdES132ElementAttrAuthoritiesCertValues
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES132Element) ElementAttributeCertificateRefs() common.DSSElement {
	return XAdES132ElementAttributeCertificateRefs
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES132Element) ElementAttributeRevocationRefs() common.DSSElement {
	return XAdES132ElementAttributeRevocationRefs
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES132Element) ElementAttributeRevocationValues() common.DSSElement {
	return XAdES132ElementAttributeRevocationValues
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES132Element) ElementByKey() common.DSSElement {
	return XAdES132ElementByKey
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES132Element) ElementByName() common.DSSElement {
	return XAdES132ElementByName
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES132Element) ElementCert() common.DSSElement {
	return XAdES132ElementCert
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES132Element) ElementCertDigest() common.DSSElement {
	return XAdES132ElementCertDigest
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES132Element) ElementCertRefs() common.DSSElement {
	return XAdES132ElementCertRefs
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES132Element) ElementCertificateValues() common.DSSElement {
	return XAdES132ElementCertificateValues
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES132Element) ElementCertifiedRole() common.DSSElement {
	return XAdES132ElementCertifiedRole
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES132Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES132ElementCertifiedRoles
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES132Element) ElementCertifiedRolesV2() common.DSSElement {
	return XAdES132ElementCertifiedRolesV2
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES132Element) ElementCity() common.DSSElement {
	return XAdES132ElementCity
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES132Element) ElementClaimedRole() common.DSSElement {
	return XAdES132ElementClaimedRole
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES132Element) ElementClaimedRoles() common.DSSElement {
	return XAdES132ElementClaimedRoles
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES132Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES132ElementCommitmentTypeID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES132Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES132ElementCommitmentTypeIndication
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES132Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES132ElementCommitmentTypeQualifier
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES132Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES132ElementCommitmentTypeQualifiers
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES132Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES132ElementCompleteCertificateRefs
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES132Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES132ElementCompleteRevocationRefs
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES132Element) ElementCounterSignature() common.DSSElement {
	return XAdES132ElementCounterSignature
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES132Element) ElementCountryName() common.DSSElement {
	return XAdES132ElementCountryName
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES132Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES132ElementCRLIdentifier
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES132Element) ElementCRLRef() common.DSSElement {
	return XAdES132ElementCRLRef
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES132Element) ElementCRLRefs() common.DSSElement {
	return XAdES132ElementCRLRefs
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES132Element) ElementCRLValues() common.DSSElement {
	return XAdES132ElementCRLValues
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES132Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES132ElementDataObjectFormat
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES132Element) ElementDescription() common.DSSElement {
	return XAdES132ElementDescription
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES132Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES132ElementDigestAlgAndValue
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES132Element) ElementDocumentationReference() common.DSSElement {
	return XAdES132ElementDocumentationReference
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES132Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES132ElementDocumentationReferences
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES132Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES132ElementEncapsulatedCRLValue
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES132Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES132ElementEncapsulatedOCSPValue
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES132Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES132ElementEncapsulatedPKIData
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES132Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES132ElementEncapsulatedTimestamp
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES132Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES132ElementEncapsulatedX509Certificate
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES132Element) ElementEncoding() common.DSSElement {
	return XAdES132ElementEncoding
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES132Element) ElementExplicitText() common.DSSElement {
	return XAdES132ElementExplicitText
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES132Element) ElementIdentifier() common.DSSElement {
	return XAdES132ElementIdentifier
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES132Element) ElementInclude() common.DSSElement {
	return XAdES132ElementInclude
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES132Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES132ElementIndividualDataObjectsTimestamp
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES132Element) Elementint() common.DSSElement {
	return XAdES132ElementInt
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES132Element) ElementIssueTime() common.DSSElement {
	return XAdES132ElementIssueTime
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES132Element) ElementIssuer() common.DSSElement {
	return XAdES132ElementIssuer
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES132Element) ElementIssuerSerial() common.DSSElement {
	return XAdES132ElementIssuerSerial
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES132Element) ElementIssuerSerialV2() common.DSSElement {
	return XAdES132ElementIssuerSerialV2
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES132Element) ElementMimeType() common.DSSElement {
	return XAdES132ElementMIMEType
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES132Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES132ElementNoticeNumbers
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES132Element) ElementNoticeRef() common.DSSElement {
	return XAdES132ElementNoticeRef
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES132Element) ElementNumber() common.DSSElement {
	return XAdES132ElementNumber
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES132Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES132ElementObjectIdentifier
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES132Element) ElementObjectReference() common.DSSElement {
	return XAdES132ElementObjectReference
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES132Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES132ElementOCSPIdentifier
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES132Element) ElementOCSPRef() common.DSSElement {
	return XAdES132ElementOCSPRef
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES132Element) ElementOCSPRefs() common.DSSElement {
	return XAdES132ElementOCSPRefs
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES132Element) ElementOCSPValues() common.DSSElement {
	return XAdES132ElementOCSPValues
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES132Element) ElementOrganization() common.DSSElement {
	return XAdES132ElementOrganization
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES132Element) ElementOtherAttributeCertificate() common.DSSElement {
	return XAdES132ElementOtherAttributeCertificate
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES132Element) ElementOtherCertificate() common.DSSElement {
	return XAdES132ElementOtherCertificate
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES132Element) ElementOtherRef() common.DSSElement {
	return XAdES132ElementOtherRef
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES132Element) ElementOtherRefs() common.DSSElement {
	return XAdES132ElementOtherRefs
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES132Element) ElementOtherTimeStamp() common.DSSElement {
	return XAdES132ElementOtherTimestamp
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES132Element) ElementOtherValue() common.DSSElement {
	return XAdES132ElementOtherValue
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES132Element) ElementOtherValues() common.DSSElement {
	return XAdES132ElementOtherValues
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES132Element) ElementPostalCode() common.DSSElement {
	return XAdES132ElementPostalCode
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES132Element) ElementProducedAt() common.DSSElement {
	return XAdES132ElementProducedAt
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES132Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES132ElementQualifyingProperties
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES132Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES132ElementQualifyingPropertiesReference
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES132Element) ElementReferenceInfo() common.DSSElement {
	return XAdES132ElementReferenceInfo
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES132Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES132ElementRefsOnlyTimestamp
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES132Element) ElementResponderID() common.DSSElement {
	return XAdES132ElementResponderID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES132Element) ElementRevocationValues() common.DSSElement {
	return XAdES132ElementRevocationValues
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES132Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES132ElementSigAndRefsTimestamp
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES132Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES132ElementSigPolicyHash
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES132Element) ElementSigPolicyId() common.DSSElement {
	return XAdES132ElementSigPolicyID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES132Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES132ElementSigPolicyQualifier
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES132Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES132ElementSigPolicyQualifiers
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES132Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES132ElementSignaturePolicyID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES132Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES132ElementSignaturePolicyIdentifier
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES132Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES132ElementSignaturePolicyImplied
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES132Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES132ElementSignatureProductionPlace
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES132Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	return XAdES132ElementSignatureProductionPlaceV2
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES132Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES132ElementSignatureTimestamp
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES132Element) ElementSignedAssertion() common.DSSElement {
	return XAdES132ElementSignedAssertion
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES132Element) ElementSignedAssertions() common.DSSElement {
	return XAdES132ElementSignedAssertions
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES132Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES132ElementSignedDataObjectProperties
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES132Element) ElementSignedProperties() common.DSSElement {
	return XAdES132ElementSignedProperties
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES132Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES132ElementSignedSignatureProperties
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES132Element) ElementSignerRole() common.DSSElement {
	return XAdES132ElementSignerRole
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES132Element) ElementSignerRoleV2() common.DSSElement {
	return XAdES132ElementSignerRoleV2
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES132Element) ElementSigningCertificate() common.DSSElement {
	return XAdES132ElementSigningCertificate
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES132Element) ElementSigningCertificateV2() common.DSSElement {
	return XAdES132ElementSigningCertificateV2
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES132Element) ElementSigningTime() common.DSSElement {
	return XAdES132ElementSigningTime
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES132Element) ElementSPURI() common.DSSElement {
	return XAdES132ElementSPURI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES132Element) ElementSPUserNotice() common.DSSElement {
	return XAdES132ElementSPUserNotice
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES132Element) ElementStateOrProvince() common.DSSElement {
	return XAdES132ElementStateOrProvince
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES132Element) ElementStreetAddress() common.DSSElement {
	return XAdES132ElementStreetAddress
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES132Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES132ElementUnsignedDataObjectProperties
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES132Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES132ElementUnsignedDataObjectProperty
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES132Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES132ElementUnsignedProperties
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES132Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES132ElementUnsignedSignatureProperties
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES132Element) ElementX509AttributeCertificate() common.DSSElement {
	return XAdES132ElementX509AttributeCertificate
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES132Element) ElementXAdESTimeStamp() common.DSSElement {
	return XAdES132ElementXAdESTimestamp
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES132Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES132ElementXMLTimestamp
}

var _ XAdESElement = XAdES132Element("")
