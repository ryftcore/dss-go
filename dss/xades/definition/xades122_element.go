// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades122/XAdES122Element.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES122Element defines elements for a XAdES 1.2.2 schema.
type XAdES122Element string

// XAdES122Element constants, one per XAdES 1.2.2 schema element name.
const (
	XAdES122ElementAllDataObjectsTimestamp        XAdES122Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES122ElementAllSignedDataObjects           XAdES122Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES122ElementAny                            XAdES122Element = "ANY"
	XAdES122ElementArchiveTimestamp               XAdES122Element = "ARCHIVE_TIMESTAMP"
	XAdES122ElementAttributeCertificateRefs       XAdES122Element = "ATTRIBUTE_CERTIFICATE_REFS"
	XAdES122ElementAttributeRevocationRefs        XAdES122Element = "ATTRIBUTE_REVOCATION_REFS"
	XAdES122ElementCert                           XAdES122Element = "CERT"
	XAdES122ElementCertDigest                     XAdES122Element = "CERT_DIGEST"
	XAdES122ElementCertRefs                       XAdES122Element = "CERT_REFS"
	XAdES122ElementCertificateValues              XAdES122Element = "CERTIFICATE_VALUES"
	XAdES122ElementCertifiedRole                  XAdES122Element = "CERTIFIED_ROLE"
	XAdES122ElementCertifiedRoles                 XAdES122Element = "CERTIFIED_ROLES"
	XAdES122ElementCity                           XAdES122Element = "CITY"
	XAdES122ElementClaimedRole                    XAdES122Element = "CLAIMED_ROLE"
	XAdES122ElementClaimedRoles                   XAdES122Element = "CLAIMED_ROLES"
	XAdES122ElementCommitmentTypeID               XAdES122Element = "COMMITMENT_TYPE_ID"
	XAdES122ElementCommitmentTypeIndication       XAdES122Element = "COMMITMENT_TYPE_INDICATION"
	XAdES122ElementCommitmentTypeQualifier        XAdES122Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES122ElementCommitmentTypeQualifiers       XAdES122Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES122ElementCompleteCertificateRefs        XAdES122Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES122ElementCompleteRevocationRefs         XAdES122Element = "COMPLETE_REVOCATION_REFS"
	XAdES122ElementCounterSignature               XAdES122Element = "COUNTER_SIGNATURE"
	XAdES122ElementCountryName                    XAdES122Element = "COUNTRY_NAME"
	XAdES122ElementCRLIdentifier                  XAdES122Element = "CRL_IDENTIFIER"
	XAdES122ElementCRLRef                         XAdES122Element = "CRL_REF"
	XAdES122ElementCRLRefs                        XAdES122Element = "CRL_REFS"
	XAdES122ElementCRLValues                      XAdES122Element = "CRL_VALUES"
	XAdES122ElementDataObjectFormat               XAdES122Element = "DATA_OBJECT_FORMAT"
	XAdES122ElementDescription                    XAdES122Element = "DESCRIPTION"
	XAdES122ElementDigestAlgAndValue              XAdES122Element = "DIGEST_ALG_AND_VALUE"
	XAdES122ElementDocumentationReference         XAdES122Element = "DOCUMENTATION_REFERENCE"
	XAdES122ElementDocumentationReferences        XAdES122Element = "DOCUMENTATION_REFERENCES"
	XAdES122ElementEncapsulatedCRLValue           XAdES122Element = "ENCAPSULATED_CRL_VALUE"
	XAdES122ElementEncapsulatedOCSPValue          XAdES122Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES122ElementEncapsulatedPKIData            XAdES122Element = "ENCAPSULATED_PKI_DATA"
	XAdES122ElementEncapsulatedTimestamp          XAdES122Element = "ENCAPSULATED_TIMESTAMP"
	XAdES122ElementEncapsulatedX509Certificate    XAdES122Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES122ElementEncoding                       XAdES122Element = "ENCODING"
	XAdES122ElementExplicitText                   XAdES122Element = "EXPLICIT_TEXT"
	XAdES122ElementIdentifier                     XAdES122Element = "IDENTIFIER"
	XAdES122ElementInclude                        XAdES122Element = "INCLUDE"
	XAdES122ElementIndividualDataObjectsTimestamp XAdES122Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES122ElementInt                            XAdES122Element = "INT"
	XAdES122ElementIssueTime                      XAdES122Element = "ISSUE_TIME"
	XAdES122ElementIssuer                         XAdES122Element = "ISSUER"
	XAdES122ElementIssuerSerial                   XAdES122Element = "ISSUER_SERIAL"
	XAdES122ElementMIMEType                       XAdES122Element = "MIME_TYPE"
	XAdES122ElementNoticeNumbers                  XAdES122Element = "NOTICE_NUMBERS"
	XAdES122ElementNoticeRef                      XAdES122Element = "NOTICE_REF"
	XAdES122ElementNumber                         XAdES122Element = "NUMBER"
	XAdES122ElementObjectIdentifier               XAdES122Element = "OBJECT_IDENTIFIER"
	XAdES122ElementObjectReference                XAdES122Element = "OBJECT_REFERENCE"
	XAdES122ElementOCSPIdentifier                 XAdES122Element = "OCSP_IDENTIFIER"
	XAdES122ElementOCSPRef                        XAdES122Element = "OCSP_REF"
	XAdES122ElementOCSPRefs                       XAdES122Element = "OCSP_REFS"
	XAdES122ElementOCSPValues                     XAdES122Element = "OCSP_VALUES"
	XAdES122ElementOrganization                   XAdES122Element = "ORGANIZATION"
	XAdES122ElementOtherCertificate               XAdES122Element = "OTHER_CERTIFICATE"
	XAdES122ElementOtherRef                       XAdES122Element = "OTHER_REF"
	XAdES122ElementOtherRefs                      XAdES122Element = "OTHER_REFS"
	XAdES122ElementOtherValue                     XAdES122Element = "OTHER_VALUE"
	XAdES122ElementOtherValues                    XAdES122Element = "OTHER_VALUES"
	XAdES122ElementPostalCode                     XAdES122Element = "POSTAL_CODE"
	XAdES122ElementProducedAt                     XAdES122Element = "PRODUCED_AT"
	XAdES122ElementQualifyingProperties           XAdES122Element = "QUALIFYING_PROPERTIES"
	XAdES122ElementQualifyingPropertiesReference  XAdES122Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES122ElementRefsOnlyTimestamp              XAdES122Element = "REFS_ONLY_TIMESTAMP"
	XAdES122ElementResponderID                    XAdES122Element = "RESPONDER_ID"
	XAdES122ElementRevocationValues               XAdES122Element = "REVOCATION_VALUES"
	XAdES122ElementSigAndRefsTimestamp            XAdES122Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES122ElementSigPolicyHash                  XAdES122Element = "SIG_POLICY_HASH"
	XAdES122ElementSigPolicyID                    XAdES122Element = "SIG_POLICY_ID"
	XAdES122ElementSigPolicyQualifier             XAdES122Element = "SIG_POLICY_QUALIFIER"
	XAdES122ElementSigPolicyQualifiers            XAdES122Element = "SIG_POLICY_QUALIFIERS"
	XAdES122ElementSignaturePolicyID              XAdES122Element = "SIGNATURE_POLICY_ID"
	XAdES122ElementSignaturePolicyIdentifier      XAdES122Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES122ElementSignaturePolicyImplied         XAdES122Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES122ElementSignatureProductionPlace       XAdES122Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES122ElementSignatureTimestamp             XAdES122Element = "SIGNATURE_TIMESTAMP"
	XAdES122ElementSignedDataObjectProperties     XAdES122Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES122ElementSignedProperties               XAdES122Element = "SIGNED_PROPERTIES"
	XAdES122ElementSignedSignatureProperties      XAdES122Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES122ElementSignerRole                     XAdES122Element = "SIGNER_ROLE"
	XAdES122ElementSigningCertificate             XAdES122Element = "SIGNING_CERTIFICATE"
	XAdES122ElementSigningTime                    XAdES122Element = "SIGNING_TIME"
	XAdES122ElementSPURI                          XAdES122Element = "SP_URI"
	XAdES122ElementSPUserNotice                   XAdES122Element = "SP_USER_NOTICE"
	XAdES122ElementStateOrProvince                XAdES122Element = "STATE_OR_PROVINCE"
	XAdES122ElementTimestamp                      XAdES122Element = "TIMESTAMP"
	XAdES122ElementUnsignedDataObjectProperties   XAdES122Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES122ElementUnsignedDataObjectProperty     XAdES122Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES122ElementUnsignedProperties             XAdES122Element = "UNSIGNED_PROPERTIES"
	XAdES122ElementUnsignedSignatureProperties    XAdES122Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES122ElementXMLTimestamp                   XAdES122Element = "XML_TIMESTAMP"
)

// xades122elementTagNames maps each constant to its wire tag name (getTagName()).
var xades122elementTagNames = map[XAdES122Element]string{
	XAdES122ElementAllDataObjectsTimestamp:        "AllDataObjectsTimeStamp",
	XAdES122ElementAllSignedDataObjects:           "AllSignedDataObjects",
	XAdES122ElementAny:                            "Any",
	XAdES122ElementArchiveTimestamp:               "ArchiveTimeStamp",
	XAdES122ElementAttributeCertificateRefs:       "AttributeCertificateRefs",
	XAdES122ElementAttributeRevocationRefs:        "AttributeRevocationRefs",
	XAdES122ElementCert:                           "Cert",
	XAdES122ElementCertDigest:                     "CertDigest",
	XAdES122ElementCertRefs:                       "CertRefs",
	XAdES122ElementCertificateValues:              "CertificateValues",
	XAdES122ElementCertifiedRole:                  "CertifiedRole",
	XAdES122ElementCertifiedRoles:                 "CertifiedRoles",
	XAdES122ElementCity:                           "City",
	XAdES122ElementClaimedRole:                    "ClaimedRole",
	XAdES122ElementClaimedRoles:                   "ClaimedRoles",
	XAdES122ElementCommitmentTypeID:               "CommitmentTypeId",
	XAdES122ElementCommitmentTypeIndication:       "CommitmentTypeIndication",
	XAdES122ElementCommitmentTypeQualifier:        "CommitmentTypeQualifier",
	XAdES122ElementCommitmentTypeQualifiers:       "CommitmentTypeQualifiers",
	XAdES122ElementCompleteCertificateRefs:        "CompleteCertificateRefs",
	XAdES122ElementCompleteRevocationRefs:         "CompleteRevocationRefs",
	XAdES122ElementCounterSignature:               "CounterSignature",
	XAdES122ElementCountryName:                    "CountryName",
	XAdES122ElementCRLIdentifier:                  "CRLIdentifier",
	XAdES122ElementCRLRef:                         "CRLRef",
	XAdES122ElementCRLRefs:                        "CRLRefs",
	XAdES122ElementCRLValues:                      "CRLValues",
	XAdES122ElementDataObjectFormat:               "DataObjectFormat",
	XAdES122ElementDescription:                    "Description",
	XAdES122ElementDigestAlgAndValue:              "DigestAlgAndValue",
	XAdES122ElementDocumentationReference:         "DocumentationReference",
	XAdES122ElementDocumentationReferences:        "DocumentationReferences",
	XAdES122ElementEncapsulatedCRLValue:           "EncapsulatedCRLValue",
	XAdES122ElementEncapsulatedOCSPValue:          "EncapsulatedOCSPValue",
	XAdES122ElementEncapsulatedPKIData:            "EncapsulatedPKIData",
	XAdES122ElementEncapsulatedTimestamp:          "EncapsulatedTimeStamp",
	XAdES122ElementEncapsulatedX509Certificate:    "EncapsulatedX509Certificate",
	XAdES122ElementEncoding:                       "Encoding",
	XAdES122ElementExplicitText:                   "ExplicitText",
	XAdES122ElementIdentifier:                     "Identifier",
	XAdES122ElementInclude:                        "Include",
	XAdES122ElementIndividualDataObjectsTimestamp: "IndividualDataObjectsTimeStamp",
	XAdES122ElementInt:                            "int",
	XAdES122ElementIssueTime:                      "IssueTime",
	XAdES122ElementIssuer:                         "Issuer",
	XAdES122ElementIssuerSerial:                   "IssuerSerial",
	XAdES122ElementMIMEType:                       "MimeType",
	XAdES122ElementNoticeNumbers:                  "NoticeNumbers",
	XAdES122ElementNoticeRef:                      "NoticeRef",
	XAdES122ElementNumber:                         "Number",
	XAdES122ElementObjectIdentifier:               "ObjectIdentifier",
	XAdES122ElementObjectReference:                "ObjectReference",
	XAdES122ElementOCSPIdentifier:                 "OCSPIdentifier",
	XAdES122ElementOCSPRef:                        "OCSPRef",
	XAdES122ElementOCSPRefs:                       "OCSPRefs",
	XAdES122ElementOCSPValues:                     "OCSPValues",
	XAdES122ElementOrganization:                   "Organization",
	XAdES122ElementOtherCertificate:               "OtherCertificate",
	XAdES122ElementOtherRef:                       "OtherRef",
	XAdES122ElementOtherRefs:                      "OtherRefs",
	XAdES122ElementOtherValue:                     "OtherValue",
	XAdES122ElementOtherValues:                    "OtherValues",
	XAdES122ElementPostalCode:                     "PostalCode",
	XAdES122ElementProducedAt:                     "ProducedAt",
	XAdES122ElementQualifyingProperties:           "QualifyingProperties",
	XAdES122ElementQualifyingPropertiesReference:  "QualifyingPropertiesReference",
	XAdES122ElementRefsOnlyTimestamp:              "RefsOnlyTimeStamp",
	XAdES122ElementResponderID:                    "ResponderID",
	XAdES122ElementRevocationValues:               "RevocationValues",
	XAdES122ElementSigAndRefsTimestamp:            "SigAndRefsTimeStamp",
	XAdES122ElementSigPolicyHash:                  "SigPolicyHash",
	XAdES122ElementSigPolicyID:                    "SigPolicyId",
	XAdES122ElementSigPolicyQualifier:             "SigPolicyQualifier",
	XAdES122ElementSigPolicyQualifiers:            "SigPolicyQualifiers",
	XAdES122ElementSignaturePolicyID:              "SignaturePolicyId",
	XAdES122ElementSignaturePolicyIdentifier:      "SignaturePolicyIdentifier",
	XAdES122ElementSignaturePolicyImplied:         "SignaturePolicyImplied",
	XAdES122ElementSignatureProductionPlace:       "SignatureProductionPlace",
	XAdES122ElementSignatureTimestamp:             "SignatureTimeStamp",
	XAdES122ElementSignedDataObjectProperties:     "SignedDataObjectProperties",
	XAdES122ElementSignedProperties:               "SignedProperties",
	XAdES122ElementSignedSignatureProperties:      "SignedSignatureProperties",
	XAdES122ElementSignerRole:                     "SignerRole",
	XAdES122ElementSigningCertificate:             "SigningCertificate",
	XAdES122ElementSigningTime:                    "SigningTime",
	XAdES122ElementSPURI:                          "SPURI",
	XAdES122ElementSPUserNotice:                   "SPUserNotice",
	XAdES122ElementStateOrProvince:                "StateOrProvince",
	XAdES122ElementTimestamp:                      "TimeStamp",
	XAdES122ElementUnsignedDataObjectProperties:   "UnsignedDataObjectProperties",
	XAdES122ElementUnsignedDataObjectProperty:     "UnsignedDataObjectProperty",
	XAdES122ElementUnsignedProperties:             "UnsignedProperties",
	XAdES122ElementUnsignedSignatureProperties:    "UnsignedSignatureProperties",
	XAdES122ElementXMLTimestamp:                   "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES122Element) TagName() string {
	return xades122elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES122Element) Namespace() *common.DSSNamespace {
	return XAdESNamespaceXAdES122
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES122Element) URI() string {
	return XAdESNamespaceXAdES122.Uri()
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
		XAdES122ElementAllDataObjectsTimestamp,
		XAdES122ElementAllSignedDataObjects,
		XAdES122ElementAny,
		XAdES122ElementArchiveTimestamp,
		XAdES122ElementAttributeCertificateRefs,
		XAdES122ElementAttributeRevocationRefs,
		XAdES122ElementCert,
		XAdES122ElementCertDigest,
		XAdES122ElementCertRefs,
		XAdES122ElementCertificateValues,
		XAdES122ElementCertifiedRole,
		XAdES122ElementCertifiedRoles,
		XAdES122ElementCity,
		XAdES122ElementClaimedRole,
		XAdES122ElementClaimedRoles,
		XAdES122ElementCommitmentTypeID,
		XAdES122ElementCommitmentTypeIndication,
		XAdES122ElementCommitmentTypeQualifier,
		XAdES122ElementCommitmentTypeQualifiers,
		XAdES122ElementCompleteCertificateRefs,
		XAdES122ElementCompleteRevocationRefs,
		XAdES122ElementCounterSignature,
		XAdES122ElementCountryName,
		XAdES122ElementCRLIdentifier,
		XAdES122ElementCRLRef,
		XAdES122ElementCRLRefs,
		XAdES122ElementCRLValues,
		XAdES122ElementDataObjectFormat,
		XAdES122ElementDescription,
		XAdES122ElementDigestAlgAndValue,
		XAdES122ElementDocumentationReference,
		XAdES122ElementDocumentationReferences,
		XAdES122ElementEncapsulatedCRLValue,
		XAdES122ElementEncapsulatedOCSPValue,
		XAdES122ElementEncapsulatedPKIData,
		XAdES122ElementEncapsulatedTimestamp,
		XAdES122ElementEncapsulatedX509Certificate,
		XAdES122ElementEncoding,
		XAdES122ElementExplicitText,
		XAdES122ElementIdentifier,
		XAdES122ElementInclude,
		XAdES122ElementIndividualDataObjectsTimestamp,
		XAdES122ElementInt,
		XAdES122ElementIssueTime,
		XAdES122ElementIssuer,
		XAdES122ElementIssuerSerial,
		XAdES122ElementMIMEType,
		XAdES122ElementNoticeNumbers,
		XAdES122ElementNoticeRef,
		XAdES122ElementNumber,
		XAdES122ElementObjectIdentifier,
		XAdES122ElementObjectReference,
		XAdES122ElementOCSPIdentifier,
		XAdES122ElementOCSPRef,
		XAdES122ElementOCSPRefs,
		XAdES122ElementOCSPValues,
		XAdES122ElementOrganization,
		XAdES122ElementOtherCertificate,
		XAdES122ElementOtherRef,
		XAdES122ElementOtherRefs,
		XAdES122ElementOtherValue,
		XAdES122ElementOtherValues,
		XAdES122ElementPostalCode,
		XAdES122ElementProducedAt,
		XAdES122ElementQualifyingProperties,
		XAdES122ElementQualifyingPropertiesReference,
		XAdES122ElementRefsOnlyTimestamp,
		XAdES122ElementResponderID,
		XAdES122ElementRevocationValues,
		XAdES122ElementSigAndRefsTimestamp,
		XAdES122ElementSigPolicyHash,
		XAdES122ElementSigPolicyID,
		XAdES122ElementSigPolicyQualifier,
		XAdES122ElementSigPolicyQualifiers,
		XAdES122ElementSignaturePolicyID,
		XAdES122ElementSignaturePolicyIdentifier,
		XAdES122ElementSignaturePolicyImplied,
		XAdES122ElementSignatureProductionPlace,
		XAdES122ElementSignatureTimestamp,
		XAdES122ElementSignedDataObjectProperties,
		XAdES122ElementSignedProperties,
		XAdES122ElementSignedSignatureProperties,
		XAdES122ElementSignerRole,
		XAdES122ElementSigningCertificate,
		XAdES122ElementSigningTime,
		XAdES122ElementSPURI,
		XAdES122ElementSPUserNotice,
		XAdES122ElementStateOrProvince,
		XAdES122ElementTimestamp,
		XAdES122ElementUnsignedDataObjectProperties,
		XAdES122ElementUnsignedDataObjectProperty,
		XAdES122ElementUnsignedProperties,
		XAdES122ElementUnsignedSignatureProperties,
		XAdES122ElementXMLTimestamp,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES122Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES122ElementAllDataObjectsTimestamp
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES122Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES122ElementAllSignedDataObjects
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES122Element) ElementAny() common.DSSElement {
	return XAdES122ElementAny
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES122Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES122ElementArchiveTimestamp
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES122Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES122Element) ElementAttributeCertificateRefs() common.DSSElement {
	return XAdES122ElementAttributeCertificateRefs
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES122Element) ElementAttributeRevocationRefs() common.DSSElement {
	return XAdES122ElementAttributeRevocationRefs
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES122Element) ElementAttributeRevocationValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES122Element) ElementByKey() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES122Element) ElementByName() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES122Element) ElementCert() common.DSSElement {
	return XAdES122ElementCert
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES122Element) ElementCertDigest() common.DSSElement {
	return XAdES122ElementCertDigest
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES122Element) ElementCertRefs() common.DSSElement {
	return XAdES122ElementCertRefs
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES122Element) ElementCertificateValues() common.DSSElement {
	return XAdES122ElementCertificateValues
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES122Element) ElementCertifiedRole() common.DSSElement {
	return XAdES122ElementCertifiedRole
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES122Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES122ElementCertifiedRoles
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES122Element) ElementCertifiedRolesV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES122Element) ElementCity() common.DSSElement {
	return XAdES122ElementCity
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES122Element) ElementClaimedRole() common.DSSElement {
	return XAdES122ElementClaimedRole
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES122Element) ElementClaimedRoles() common.DSSElement {
	return XAdES122ElementClaimedRoles
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES122Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES122ElementCommitmentTypeID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES122Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES122ElementCommitmentTypeIndication
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES122Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES122ElementCommitmentTypeQualifier
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES122Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES122ElementCommitmentTypeQualifiers
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES122Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES122ElementCompleteCertificateRefs
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES122Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES122ElementCompleteRevocationRefs
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES122Element) ElementCounterSignature() common.DSSElement {
	return XAdES122ElementCounterSignature
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES122Element) ElementCountryName() common.DSSElement {
	return XAdES122ElementCountryName
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES122Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES122ElementCRLIdentifier
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES122Element) ElementCRLRef() common.DSSElement {
	return XAdES122ElementCRLRef
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES122Element) ElementCRLRefs() common.DSSElement {
	return XAdES122ElementCRLRefs
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES122Element) ElementCRLValues() common.DSSElement {
	return XAdES122ElementCRLValues
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES122Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES122ElementDataObjectFormat
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES122Element) ElementDescription() common.DSSElement {
	return XAdES122ElementDescription
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES122Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES122ElementDigestAlgAndValue
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES122Element) ElementDocumentationReference() common.DSSElement {
	return XAdES122ElementDocumentationReference
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES122Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES122ElementDocumentationReferences
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES122Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES122ElementEncapsulatedCRLValue
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES122Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES122ElementEncapsulatedOCSPValue
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES122Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES122ElementEncapsulatedPKIData
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES122Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES122ElementEncapsulatedTimestamp
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES122Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES122ElementEncapsulatedX509Certificate
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES122Element) ElementEncoding() common.DSSElement {
	return XAdES122ElementEncoding
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES122Element) ElementExplicitText() common.DSSElement {
	return XAdES122ElementExplicitText
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES122Element) ElementIdentifier() common.DSSElement {
	return XAdES122ElementIdentifier
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES122Element) ElementInclude() common.DSSElement {
	return XAdES122ElementInclude
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES122Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES122ElementIndividualDataObjectsTimestamp
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES122Element) Elementint() common.DSSElement {
	return XAdES122ElementInt
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES122Element) ElementIssueTime() common.DSSElement {
	return XAdES122ElementIssueTime
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES122Element) ElementIssuer() common.DSSElement {
	return XAdES122ElementIssuer
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES122Element) ElementIssuerSerial() common.DSSElement {
	return XAdES122ElementIssuerSerial
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES122Element) ElementIssuerSerialV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES122Element) ElementMimeType() common.DSSElement {
	return XAdES122ElementMIMEType
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES122Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES122ElementNoticeNumbers
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES122Element) ElementNoticeRef() common.DSSElement {
	return XAdES122ElementNoticeRef
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES122Element) ElementNumber() common.DSSElement {
	return XAdES122ElementNumber
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES122Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES122ElementObjectIdentifier
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES122Element) ElementObjectReference() common.DSSElement {
	return XAdES122ElementObjectReference
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES122Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES122ElementOCSPIdentifier
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES122Element) ElementOCSPRef() common.DSSElement {
	return XAdES122ElementOCSPRef
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES122Element) ElementOCSPRefs() common.DSSElement {
	return XAdES122ElementOCSPRefs
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES122Element) ElementOCSPValues() common.DSSElement {
	return XAdES122ElementOCSPValues
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES122Element) ElementOrganization() common.DSSElement {
	return XAdES122ElementOrganization
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES122Element) ElementOtherAttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES122Element) ElementOtherCertificate() common.DSSElement {
	return XAdES122ElementOtherCertificate
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES122Element) ElementOtherRef() common.DSSElement {
	return XAdES122ElementOtherRef
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES122Element) ElementOtherRefs() common.DSSElement {
	return XAdES122ElementOtherRefs
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES122Element) ElementOtherTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES122Element) ElementOtherValue() common.DSSElement {
	return XAdES122ElementOtherValue
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES122Element) ElementOtherValues() common.DSSElement {
	return XAdES122ElementOtherValues
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES122Element) ElementPostalCode() common.DSSElement {
	return XAdES122ElementPostalCode
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES122Element) ElementProducedAt() common.DSSElement {
	return XAdES122ElementProducedAt
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES122Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES122ElementQualifyingProperties
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES122Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES122ElementQualifyingPropertiesReference
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES122Element) ElementReferenceInfo() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES122Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES122ElementRefsOnlyTimestamp
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES122Element) ElementResponderID() common.DSSElement {
	return XAdES122ElementResponderID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES122Element) ElementRevocationValues() common.DSSElement {
	return XAdES122ElementRevocationValues
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES122Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES122ElementSigAndRefsTimestamp
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES122Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES122ElementSigPolicyHash
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES122Element) ElementSigPolicyId() common.DSSElement {
	return XAdES122ElementSigPolicyID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES122Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES122ElementSigPolicyQualifier
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES122Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES122ElementSigPolicyQualifiers
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES122Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES122ElementSignaturePolicyID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES122Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES122ElementSignaturePolicyIdentifier
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES122Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES122ElementSignaturePolicyImplied
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES122Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES122ElementSignatureProductionPlace
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES122Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES122Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES122ElementSignatureTimestamp
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES122Element) ElementSignedAssertion() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES122Element) ElementSignedAssertions() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES122Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES122ElementSignedDataObjectProperties
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES122Element) ElementSignedProperties() common.DSSElement {
	return XAdES122ElementSignedProperties
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES122Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES122ElementSignedSignatureProperties
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES122Element) ElementSignerRole() common.DSSElement {
	return XAdES122ElementSignerRole
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES122Element) ElementSignerRoleV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES122Element) ElementSigningCertificate() common.DSSElement {
	return XAdES122ElementSigningCertificate
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES122Element) ElementSigningCertificateV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES122Element) ElementSigningTime() common.DSSElement {
	return XAdES122ElementSigningTime
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES122Element) ElementSPURI() common.DSSElement {
	return XAdES122ElementSPURI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES122Element) ElementSPUserNotice() common.DSSElement {
	return XAdES122ElementSPUserNotice
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES122Element) ElementStateOrProvince() common.DSSElement {
	return XAdES122ElementStateOrProvince
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES122Element) ElementStreetAddress() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES122Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES122ElementUnsignedDataObjectProperties
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES122Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES122ElementUnsignedDataObjectProperty
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES122Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES122ElementUnsignedProperties
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES122Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES122ElementUnsignedSignatureProperties
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES122Element) ElementX509AttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES122Element) ElementXAdESTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES122.Uri())
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES122Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES122ElementXMLTimestamp
}

var _ XAdESElement = XAdES122Element("")
