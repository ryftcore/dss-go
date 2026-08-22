// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Element.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES111Element defines elements for a XAdES 1.1.1 schema.
type XAdES111Element string

// XAdES111Element constants, one per XAdES 1.1.1 schema element name.
const (
	XAdES111ElementAllDataObjectsTimestamp        XAdES111Element = "ALL_DATA_OBJECTS_TIMESTAMP"
	XAdES111ElementAllSignedDataObjects           XAdES111Element = "ALL_SIGNED_DATA_OBJECTS"
	XAdES111ElementAny                            XAdES111Element = "ANY"
	XAdES111ElementArchiveTimestamp               XAdES111Element = "ARCHIVE_TIMESTAMP"
	XAdES111ElementCert                           XAdES111Element = "CERT"
	XAdES111ElementCertDigest                     XAdES111Element = "CERT_DIGEST"
	XAdES111ElementCertRefs                       XAdES111Element = "CERT_REFS"
	XAdES111ElementCertificateValues              XAdES111Element = "CERTIFICATE_VALUES"
	XAdES111ElementCertifiedRole                  XAdES111Element = "CERTIFIED_ROLE"
	XAdES111ElementCertifiedRoles                 XAdES111Element = "CERTIFIED_ROLES"
	XAdES111ElementCity                           XAdES111Element = "CITY"
	XAdES111ElementClaimedRole                    XAdES111Element = "CLAIMED_ROLE"
	XAdES111ElementClaimedRoles                   XAdES111Element = "CLAIMED_ROLES"
	XAdES111ElementCommitmentTypeID               XAdES111Element = "COMMITMENT_TYPE_ID"
	XAdES111ElementCommitmentTypeIndication       XAdES111Element = "COMMITMENT_TYPE_INDICATION"
	XAdES111ElementCommitmentTypeQualifier        XAdES111Element = "COMMITMENT_TYPE_QUALIFIER"
	XAdES111ElementCommitmentTypeQualifiers       XAdES111Element = "COMMITMENT_TYPE_QUALIFIERS"
	XAdES111ElementCompleteCertificateRefs        XAdES111Element = "COMPLETE_CERTIFICATE_REFS"
	XAdES111ElementCompleteRevocationRefs         XAdES111Element = "COMPLETE_REVOCATION_REFS"
	XAdES111ElementCounterSignature               XAdES111Element = "COUNTER_SIGNATURE"
	XAdES111ElementCountryName                    XAdES111Element = "COUNTRY_NAME"
	XAdES111ElementCRLIdentifier                  XAdES111Element = "CRL_IDENTIFIER"
	XAdES111ElementCRLRef                         XAdES111Element = "CRL_REF"
	XAdES111ElementCRLRefs                        XAdES111Element = "CRL_REFS"
	XAdES111ElementCRLValues                      XAdES111Element = "CRL_VALUES"
	XAdES111ElementDataObjectFormat               XAdES111Element = "DATA_OBJECT_FORMAT"
	XAdES111ElementDescription                    XAdES111Element = "DESCRIPTION"
	XAdES111ElementDigestAlgAndValue              XAdES111Element = "DIGEST_ALG_AND_VALUE"
	XAdES111ElementDigestMethod                   XAdES111Element = "DIGEST_METHOD"
	XAdES111ElementDigestValue                    XAdES111Element = "DIGEST_VALUE"
	XAdES111ElementDocumentationReference         XAdES111Element = "DOCUMENTATION_REFERENCE"
	XAdES111ElementDocumentationReferences        XAdES111Element = "DOCUMENTATION_REFERENCES"
	XAdES111ElementEncapsulatedCRLValue           XAdES111Element = "ENCAPSULATED_CRL_VALUE"
	XAdES111ElementEncapsulatedOCSPValue          XAdES111Element = "ENCAPSULATED_OCSP_VALUE"
	XAdES111ElementEncapsulatedPKIData            XAdES111Element = "ENCAPSULATED_PKI_DATA"
	XAdES111ElementEncapsulatedTimestamp          XAdES111Element = "ENCAPSULATED_TIMESTAMP"
	XAdES111ElementEncapsulatedX509Certificate    XAdES111Element = "ENCAPSULATED_X509_CERTIFICATE"
	XAdES111ElementEncoding                       XAdES111Element = "ENCODING"
	XAdES111ElementExplicitText                   XAdES111Element = "EXPLICIT_TEXT"
	XAdES111ElementHashDataInfo                   XAdES111Element = "HASH_DATA_INFO"
	XAdES111ElementIdentifier                     XAdES111Element = "IDENTIFIER"
	XAdES111ElementIndividualDataObjectsTimestamp XAdES111Element = "INDIVIDUAL_DATA_OBJECTS_TIMESTAMP"
	XAdES111ElementInt                            XAdES111Element = "INT"
	XAdES111ElementIssueTime                      XAdES111Element = "ISSUE_TIME"
	XAdES111ElementIssuer                         XAdES111Element = "ISSUER"
	XAdES111ElementIssuerSerial                   XAdES111Element = "ISSUER_SERIAL"
	XAdES111ElementMIMEType                       XAdES111Element = "MIME_TYPE"
	XAdES111ElementNoticeNumbers                  XAdES111Element = "NOTICE_NUMBERS"
	XAdES111ElementNoticeRef                      XAdES111Element = "NOTICE_REF"
	XAdES111ElementNumber                         XAdES111Element = "NUMBER"
	XAdES111ElementObjectIdentifier               XAdES111Element = "OBJECT_IDENTIFIER"
	XAdES111ElementObjectReference                XAdES111Element = "OBJECT_REFERENCE"
	XAdES111ElementOCSPIdentifier                 XAdES111Element = "OCSP_IDENTIFIER"
	XAdES111ElementOCSPRef                        XAdES111Element = "OCSP_REF"
	XAdES111ElementOCSPRefs                       XAdES111Element = "OCSP_REFS"
	XAdES111ElementOCSPValues                     XAdES111Element = "OCSP_VALUES"
	XAdES111ElementOrganization                   XAdES111Element = "ORGANIZATION"
	XAdES111ElementOtherCertificate               XAdES111Element = "OTHER_CERTIFICATE"
	XAdES111ElementOtherRef                       XAdES111Element = "OTHER_REF"
	XAdES111ElementOtherRefs                      XAdES111Element = "OTHER_REFS"
	XAdES111ElementOtherValue                     XAdES111Element = "OTHER_VALUE"
	XAdES111ElementOtherValues                    XAdES111Element = "OTHER_VALUES"
	XAdES111ElementPostalCode                     XAdES111Element = "POSTAL_CODE"
	XAdES111ElementProducedAt                     XAdES111Element = "PRODUCED_AT"
	XAdES111ElementQualifyingProperties           XAdES111Element = "QUALIFYING_PROPERTIES"
	XAdES111ElementQualifyingPropertiesReference  XAdES111Element = "QUALIFYING_PROPERTIES_REFERENCE"
	XAdES111ElementRefsOnlyTimestamp              XAdES111Element = "REFS_ONLY_TIMESTAMP"
	XAdES111ElementResponderID                    XAdES111Element = "RESPONDER_ID"
	XAdES111ElementRevocationValues               XAdES111Element = "REVOCATION_VALUES"
	XAdES111ElementSigAndRefsTimestamp            XAdES111Element = "SIG_AND_REFS_TIMESTAMP"
	XAdES111ElementSigPolicyHash                  XAdES111Element = "SIG_POLICY_HASH"
	XAdES111ElementSigPolicyID                    XAdES111Element = "SIG_POLICY_ID"
	XAdES111ElementSigPolicyQualifier             XAdES111Element = "SIG_POLICY_QUALIFIER"
	XAdES111ElementSigPolicyQualifiers            XAdES111Element = "SIG_POLICY_QUALIFIERS"
	XAdES111ElementSignaturePolicyID              XAdES111Element = "SIGNATURE_POLICY_ID"
	XAdES111ElementSignaturePolicyIdentifier      XAdES111Element = "SIGNATURE_POLICY_IDENTIFIER"
	XAdES111ElementSignaturePolicyImplied         XAdES111Element = "SIGNATURE_POLICY_IMPLIED"
	XAdES111ElementSignatureProductionPlace       XAdES111Element = "SIGNATURE_PRODUCTION_PLACE"
	XAdES111ElementSignatureTimestamp             XAdES111Element = "SIGNATURE_TIMESTAMP"
	XAdES111ElementSignedDataObjectProperties     XAdES111Element = "SIGNED_DATA_OBJECT_PROPERTIES"
	XAdES111ElementSignedProperties               XAdES111Element = "SIGNED_PROPERTIES"
	XAdES111ElementSignedSignatureProperties      XAdES111Element = "SIGNED_SIGNATURE_PROPERTIES"
	XAdES111ElementSignerRole                     XAdES111Element = "SIGNER_ROLE"
	XAdES111ElementSigningCertificate             XAdES111Element = "SIGNING_CERTIFICATE"
	XAdES111ElementSigningTime                    XAdES111Element = "SIGNING_TIME"
	XAdES111ElementSPURI                          XAdES111Element = "SP_URI"
	XAdES111ElementSPUserNotice                   XAdES111Element = "SP_USER_NOTICE"
	XAdES111ElementStateOrProvince                XAdES111Element = "STATE_OR_PROVINCE"
	XAdES111ElementTimestamp                      XAdES111Element = "TIMESTAMP"
	XAdES111ElementTransforms                     XAdES111Element = "TRANSFORMS"
	XAdES111ElementUnsignedDataObjectProperties   XAdES111Element = "UNSIGNED_DATA_OBJECT_PROPERTIES"
	XAdES111ElementUnsignedDataObjectProperty     XAdES111Element = "UNSIGNED_DATA_OBJECT_PROPERTY"
	XAdES111ElementUnsignedProperties             XAdES111Element = "UNSIGNED_PROPERTIES"
	XAdES111ElementUnsignedSignatureProperties    XAdES111Element = "UNSIGNED_SIGNATURE_PROPERTIES"
	XAdES111ElementXMLTimestamp                   XAdES111Element = "XML_TIMESTAMP"
)

// xades111elementTagNames maps each constant to its wire tag name (getTagName()).
var xades111elementTagNames = map[XAdES111Element]string{
	XAdES111ElementAllDataObjectsTimestamp:        "AllDataObjectsTimeStamp",
	XAdES111ElementAllSignedDataObjects:           "AllSignedDataObjects",
	XAdES111ElementAny:                            "Any",
	XAdES111ElementArchiveTimestamp:               "ArchiveTimeStamp",
	XAdES111ElementCert:                           "Cert",
	XAdES111ElementCertDigest:                     "CertDigest",
	XAdES111ElementCertRefs:                       "CertRefs",
	XAdES111ElementCertificateValues:              "CertificateValues",
	XAdES111ElementCertifiedRole:                  "CertifiedRole",
	XAdES111ElementCertifiedRoles:                 "CertifiedRoles",
	XAdES111ElementCity:                           "City",
	XAdES111ElementClaimedRole:                    "ClaimedRole",
	XAdES111ElementClaimedRoles:                   "ClaimedRoles",
	XAdES111ElementCommitmentTypeID:               "CommitmentTypeId",
	XAdES111ElementCommitmentTypeIndication:       "CommitmentTypeIndication",
	XAdES111ElementCommitmentTypeQualifier:        "CommitmentTypeQualifier",
	XAdES111ElementCommitmentTypeQualifiers:       "CommitmentTypeQualifiers",
	XAdES111ElementCompleteCertificateRefs:        "CompleteCertificateRefs",
	XAdES111ElementCompleteRevocationRefs:         "CompleteRevocationRefs",
	XAdES111ElementCounterSignature:               "CounterSignature",
	XAdES111ElementCountryName:                    "CountryName",
	XAdES111ElementCRLIdentifier:                  "CRLIdentifier",
	XAdES111ElementCRLRef:                         "CRLRef",
	XAdES111ElementCRLRefs:                        "CRLRefs",
	XAdES111ElementCRLValues:                      "CRLValues",
	XAdES111ElementDataObjectFormat:               "DataObjectFormat",
	XAdES111ElementDescription:                    "Description",
	XAdES111ElementDigestAlgAndValue:              "DigestAlgAndValue",
	XAdES111ElementDigestMethod:                   "DigestMethod",
	XAdES111ElementDigestValue:                    "DigestValue",
	XAdES111ElementDocumentationReference:         "DocumentationReference",
	XAdES111ElementDocumentationReferences:        "DocumentationReferences",
	XAdES111ElementEncapsulatedCRLValue:           "EncapsulatedCRLValue",
	XAdES111ElementEncapsulatedOCSPValue:          "EncapsulatedOCSPValue",
	XAdES111ElementEncapsulatedPKIData:            "EncapsulatedPKIData",
	XAdES111ElementEncapsulatedTimestamp:          "EncapsulatedTimeStamp",
	XAdES111ElementEncapsulatedX509Certificate:    "EncapsulatedX509Certificate",
	XAdES111ElementEncoding:                       "Encoding",
	XAdES111ElementExplicitText:                   "ExplicitText",
	XAdES111ElementHashDataInfo:                   "HashDataInfo",
	XAdES111ElementIdentifier:                     "Identifier",
	XAdES111ElementIndividualDataObjectsTimestamp: "IndividualDataObjectsTimeStamp",
	XAdES111ElementInt:                            "int",
	XAdES111ElementIssueTime:                      "IssueTime",
	XAdES111ElementIssuer:                         "Issuer",
	XAdES111ElementIssuerSerial:                   "IssuerSerial",
	XAdES111ElementMIMEType:                       "MimeType",
	XAdES111ElementNoticeNumbers:                  "NoticeNumbers",
	XAdES111ElementNoticeRef:                      "NoticeRef",
	XAdES111ElementNumber:                         "Number",
	XAdES111ElementObjectIdentifier:               "ObjectIdentifier",
	XAdES111ElementObjectReference:                "ObjectReference",
	XAdES111ElementOCSPIdentifier:                 "OCSPIdentifier",
	XAdES111ElementOCSPRef:                        "OCSPRef",
	XAdES111ElementOCSPRefs:                       "OCSPRefs",
	XAdES111ElementOCSPValues:                     "OCSPValues",
	XAdES111ElementOrganization:                   "Organization",
	XAdES111ElementOtherCertificate:               "OtherCertificate",
	XAdES111ElementOtherRef:                       "OtherRef",
	XAdES111ElementOtherRefs:                      "OtherRefs",
	XAdES111ElementOtherValue:                     "OtherValue",
	XAdES111ElementOtherValues:                    "OtherValues",
	XAdES111ElementPostalCode:                     "PostalCode",
	XAdES111ElementProducedAt:                     "ProducedAt",
	XAdES111ElementQualifyingProperties:           "QualifyingProperties",
	XAdES111ElementQualifyingPropertiesReference:  "QualifyingPropertiesReference",
	XAdES111ElementRefsOnlyTimestamp:              "RefsOnlyTimeStamp",
	XAdES111ElementResponderID:                    "ResponderID",
	XAdES111ElementRevocationValues:               "RevocationValues",
	XAdES111ElementSigAndRefsTimestamp:            "SigAndRefsTimeStamp",
	XAdES111ElementSigPolicyHash:                  "SigPolicyHash",
	XAdES111ElementSigPolicyID:                    "SigPolicyId",
	XAdES111ElementSigPolicyQualifier:             "SigPolicyQualifier",
	XAdES111ElementSigPolicyQualifiers:            "SigPolicyQualifiers",
	XAdES111ElementSignaturePolicyID:              "SignaturePolicyId",
	XAdES111ElementSignaturePolicyIdentifier:      "SignaturePolicyIdentifier",
	XAdES111ElementSignaturePolicyImplied:         "SignaturePolicyImplied",
	XAdES111ElementSignatureProductionPlace:       "SignatureProductionPlace",
	XAdES111ElementSignatureTimestamp:             "SignatureTimeStamp",
	XAdES111ElementSignedDataObjectProperties:     "SignedDataObjectProperties",
	XAdES111ElementSignedProperties:               "SignedProperties",
	XAdES111ElementSignedSignatureProperties:      "SignedSignatureProperties",
	XAdES111ElementSignerRole:                     "SignerRole",
	XAdES111ElementSigningCertificate:             "SigningCertificate",
	XAdES111ElementSigningTime:                    "SigningTime",
	XAdES111ElementSPURI:                          "SPURI",
	XAdES111ElementSPUserNotice:                   "SPUserNotice",
	XAdES111ElementStateOrProvince:                "StateOrProvince",
	XAdES111ElementTimestamp:                      "TimeStamp",
	XAdES111ElementTransforms:                     "Transforms",
	XAdES111ElementUnsignedDataObjectProperties:   "UnsignedDataObjectProperties",
	XAdES111ElementUnsignedDataObjectProperty:     "UnsignedDataObjectProperty",
	XAdES111ElementUnsignedProperties:             "UnsignedProperties",
	XAdES111ElementUnsignedSignatureProperties:    "UnsignedSignatureProperties",
	XAdES111ElementXMLTimestamp:                   "XMLTimeStamp",
}

// TagName implements common.DSSElement. Ports getTagName().
func (e XAdES111Element) TagName() string {
	return xades111elementTagNames[e]
}

// Namespace implements common.DSSElement. Ports getNamespace().
func (e XAdES111Element) Namespace() *common.DSSNamespace {
	return XAdESNamespaceXAdES111
}

// URI implements common.DSSElement. Ports getURI().
func (e XAdES111Element) URI() string {
	return XAdESNamespaceXAdES111.Uri()
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
		XAdES111ElementAllDataObjectsTimestamp,
		XAdES111ElementAllSignedDataObjects,
		XAdES111ElementAny,
		XAdES111ElementArchiveTimestamp,
		XAdES111ElementCert,
		XAdES111ElementCertDigest,
		XAdES111ElementCertRefs,
		XAdES111ElementCertificateValues,
		XAdES111ElementCertifiedRole,
		XAdES111ElementCertifiedRoles,
		XAdES111ElementCity,
		XAdES111ElementClaimedRole,
		XAdES111ElementClaimedRoles,
		XAdES111ElementCommitmentTypeID,
		XAdES111ElementCommitmentTypeIndication,
		XAdES111ElementCommitmentTypeQualifier,
		XAdES111ElementCommitmentTypeQualifiers,
		XAdES111ElementCompleteCertificateRefs,
		XAdES111ElementCompleteRevocationRefs,
		XAdES111ElementCounterSignature,
		XAdES111ElementCountryName,
		XAdES111ElementCRLIdentifier,
		XAdES111ElementCRLRef,
		XAdES111ElementCRLRefs,
		XAdES111ElementCRLValues,
		XAdES111ElementDataObjectFormat,
		XAdES111ElementDescription,
		XAdES111ElementDigestAlgAndValue,
		XAdES111ElementDigestMethod,
		XAdES111ElementDigestValue,
		XAdES111ElementDocumentationReference,
		XAdES111ElementDocumentationReferences,
		XAdES111ElementEncapsulatedCRLValue,
		XAdES111ElementEncapsulatedOCSPValue,
		XAdES111ElementEncapsulatedPKIData,
		XAdES111ElementEncapsulatedTimestamp,
		XAdES111ElementEncapsulatedX509Certificate,
		XAdES111ElementEncoding,
		XAdES111ElementExplicitText,
		XAdES111ElementHashDataInfo,
		XAdES111ElementIdentifier,
		XAdES111ElementIndividualDataObjectsTimestamp,
		XAdES111ElementInt,
		XAdES111ElementIssueTime,
		XAdES111ElementIssuer,
		XAdES111ElementIssuerSerial,
		XAdES111ElementMIMEType,
		XAdES111ElementNoticeNumbers,
		XAdES111ElementNoticeRef,
		XAdES111ElementNumber,
		XAdES111ElementObjectIdentifier,
		XAdES111ElementObjectReference,
		XAdES111ElementOCSPIdentifier,
		XAdES111ElementOCSPRef,
		XAdES111ElementOCSPRefs,
		XAdES111ElementOCSPValues,
		XAdES111ElementOrganization,
		XAdES111ElementOtherCertificate,
		XAdES111ElementOtherRef,
		XAdES111ElementOtherRefs,
		XAdES111ElementOtherValue,
		XAdES111ElementOtherValues,
		XAdES111ElementPostalCode,
		XAdES111ElementProducedAt,
		XAdES111ElementQualifyingProperties,
		XAdES111ElementQualifyingPropertiesReference,
		XAdES111ElementRefsOnlyTimestamp,
		XAdES111ElementResponderID,
		XAdES111ElementRevocationValues,
		XAdES111ElementSigAndRefsTimestamp,
		XAdES111ElementSigPolicyHash,
		XAdES111ElementSigPolicyID,
		XAdES111ElementSigPolicyQualifier,
		XAdES111ElementSigPolicyQualifiers,
		XAdES111ElementSignaturePolicyID,
		XAdES111ElementSignaturePolicyIdentifier,
		XAdES111ElementSignaturePolicyImplied,
		XAdES111ElementSignatureProductionPlace,
		XAdES111ElementSignatureTimestamp,
		XAdES111ElementSignedDataObjectProperties,
		XAdES111ElementSignedProperties,
		XAdES111ElementSignedSignatureProperties,
		XAdES111ElementSignerRole,
		XAdES111ElementSigningCertificate,
		XAdES111ElementSigningTime,
		XAdES111ElementSPURI,
		XAdES111ElementSPUserNotice,
		XAdES111ElementStateOrProvince,
		XAdES111ElementTimestamp,
		XAdES111ElementTransforms,
		XAdES111ElementUnsignedDataObjectProperties,
		XAdES111ElementUnsignedDataObjectProperty,
		XAdES111ElementUnsignedProperties,
		XAdES111ElementUnsignedSignatureProperties,
		XAdES111ElementXMLTimestamp,
	}
}

// ElementAllDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementAllDataObjectsTimeStamp().
func (e XAdES111Element) ElementAllDataObjectsTimeStamp() common.DSSElement {
	return XAdES111ElementAllDataObjectsTimestamp
}

// ElementAllSignedDataObjects implements definition.XAdESElement. Ports getElementAllSignedDataObjects().
func (e XAdES111Element) ElementAllSignedDataObjects() common.DSSElement {
	return XAdES111ElementAllSignedDataObjects
}

// ElementAny implements definition.XAdESElement. Ports getElementAny().
func (e XAdES111Element) ElementAny() common.DSSElement {
	return XAdES111ElementAny
}

// ElementArchiveTimeStamp implements definition.XAdESElement. Ports getElementArchiveTimeStamp().
func (e XAdES111Element) ElementArchiveTimeStamp() common.DSSElement {
	return XAdES111ElementArchiveTimestamp
}

// ElementAttrAuthoritiesCertValues implements definition.XAdESElement. Ports getElementAttrAuthoritiesCertValues().
func (e XAdES111Element) ElementAttrAuthoritiesCertValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementAttributeCertificateRefs implements definition.XAdESElement. Ports getElementAttributeCertificateRefs().
func (e XAdES111Element) ElementAttributeCertificateRefs() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementAttributeRevocationRefs implements definition.XAdESElement. Ports getElementAttributeRevocationRefs().
func (e XAdES111Element) ElementAttributeRevocationRefs() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementAttributeRevocationValues implements definition.XAdESElement. Ports getElementAttributeRevocationValues().
func (e XAdES111Element) ElementAttributeRevocationValues() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementByKey implements definition.XAdESElement. Ports getElementByKey().
func (e XAdES111Element) ElementByKey() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementByName implements definition.XAdESElement. Ports getElementByName().
func (e XAdES111Element) ElementByName() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementCert implements definition.XAdESElement. Ports getElementCert().
func (e XAdES111Element) ElementCert() common.DSSElement {
	return XAdES111ElementCert
}

// ElementCertDigest implements definition.XAdESElement. Ports getElementCertDigest().
func (e XAdES111Element) ElementCertDigest() common.DSSElement {
	return XAdES111ElementCertDigest
}

// ElementCertRefs implements definition.XAdESElement. Ports getElementCertRefs().
func (e XAdES111Element) ElementCertRefs() common.DSSElement {
	return XAdES111ElementCertRefs
}

// ElementCertificateValues implements definition.XAdESElement. Ports getElementCertificateValues().
func (e XAdES111Element) ElementCertificateValues() common.DSSElement {
	return XAdES111ElementCertificateValues
}

// ElementCertifiedRole implements definition.XAdESElement. Ports getElementCertifiedRole().
func (e XAdES111Element) ElementCertifiedRole() common.DSSElement {
	return XAdES111ElementCertifiedRole
}

// ElementCertifiedRoles implements definition.XAdESElement. Ports getElementCertifiedRoles().
func (e XAdES111Element) ElementCertifiedRoles() common.DSSElement {
	return XAdES111ElementCertifiedRoles
}

// ElementCertifiedRolesV2 implements definition.XAdESElement. Ports getElementCertifiedRolesV2().
func (e XAdES111Element) ElementCertifiedRolesV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementCity implements definition.XAdESElement. Ports getElementCity().
func (e XAdES111Element) ElementCity() common.DSSElement {
	return XAdES111ElementCity
}

// ElementClaimedRole implements definition.XAdESElement. Ports getElementClaimedRole().
func (e XAdES111Element) ElementClaimedRole() common.DSSElement {
	return XAdES111ElementClaimedRole
}

// ElementClaimedRoles implements definition.XAdESElement. Ports getElementClaimedRoles().
func (e XAdES111Element) ElementClaimedRoles() common.DSSElement {
	return XAdES111ElementClaimedRoles
}

// ElementCommitmentTypeId implements definition.XAdESElement. Ports getElementCommitmentTypeId().
func (e XAdES111Element) ElementCommitmentTypeId() common.DSSElement {
	return XAdES111ElementCommitmentTypeID
}

// ElementCommitmentTypeIndication implements definition.XAdESElement. Ports getElementCommitmentTypeIndication().
func (e XAdES111Element) ElementCommitmentTypeIndication() common.DSSElement {
	return XAdES111ElementCommitmentTypeIndication
}

// ElementCommitmentTypeQualifier implements definition.XAdESElement. Ports getElementCommitmentTypeQualifier().
func (e XAdES111Element) ElementCommitmentTypeQualifier() common.DSSElement {
	return XAdES111ElementCommitmentTypeQualifier
}

// ElementCommitmentTypeQualifiers implements definition.XAdESElement. Ports getElementCommitmentTypeQualifiers().
func (e XAdES111Element) ElementCommitmentTypeQualifiers() common.DSSElement {
	return XAdES111ElementCommitmentTypeQualifiers
}

// ElementCompleteCertificateRefs implements definition.XAdESElement. Ports getElementCompleteCertificateRefs().
func (e XAdES111Element) ElementCompleteCertificateRefs() common.DSSElement {
	return XAdES111ElementCompleteCertificateRefs
}

// ElementCompleteRevocationRefs implements definition.XAdESElement. Ports getElementCompleteRevocationRefs().
func (e XAdES111Element) ElementCompleteRevocationRefs() common.DSSElement {
	return XAdES111ElementCompleteRevocationRefs
}

// ElementCounterSignature implements definition.XAdESElement. Ports getElementCounterSignature().
func (e XAdES111Element) ElementCounterSignature() common.DSSElement {
	return XAdES111ElementCounterSignature
}

// ElementCountryName implements definition.XAdESElement. Ports getElementCountryName().
func (e XAdES111Element) ElementCountryName() common.DSSElement {
	return XAdES111ElementCountryName
}

// ElementCRLIdentifier implements definition.XAdESElement. Ports getElementCRLIdentifier().
func (e XAdES111Element) ElementCRLIdentifier() common.DSSElement {
	return XAdES111ElementCRLIdentifier
}

// ElementCRLRef implements definition.XAdESElement. Ports getElementCRLRef().
func (e XAdES111Element) ElementCRLRef() common.DSSElement {
	return XAdES111ElementCRLRef
}

// ElementCRLRefs implements definition.XAdESElement. Ports getElementCRLRefs().
func (e XAdES111Element) ElementCRLRefs() common.DSSElement {
	return XAdES111ElementCRLRefs
}

// ElementCRLValues implements definition.XAdESElement. Ports getElementCRLValues().
func (e XAdES111Element) ElementCRLValues() common.DSSElement {
	return XAdES111ElementCRLValues
}

// ElementDataObjectFormat implements definition.XAdESElement. Ports getElementDataObjectFormat().
func (e XAdES111Element) ElementDataObjectFormat() common.DSSElement {
	return XAdES111ElementDataObjectFormat
}

// ElementDescription implements definition.XAdESElement. Ports getElementDescription().
func (e XAdES111Element) ElementDescription() common.DSSElement {
	return XAdES111ElementDescription
}

// ElementDigestAlgAndValue implements definition.XAdESElement. Ports getElementDigestAlgAndValue().
func (e XAdES111Element) ElementDigestAlgAndValue() common.DSSElement {
	return XAdES111ElementDigestAlgAndValue
}

// ElementDocumentationReference implements definition.XAdESElement. Ports getElementDocumentationReference().
func (e XAdES111Element) ElementDocumentationReference() common.DSSElement {
	return XAdES111ElementDocumentationReference
}

// ElementDocumentationReferences implements definition.XAdESElement. Ports getElementDocumentationReferences().
func (e XAdES111Element) ElementDocumentationReferences() common.DSSElement {
	return XAdES111ElementDocumentationReferences
}

// ElementEncapsulatedCRLValue implements definition.XAdESElement. Ports getElementEncapsulatedCRLValue().
func (e XAdES111Element) ElementEncapsulatedCRLValue() common.DSSElement {
	return XAdES111ElementEncapsulatedCRLValue
}

// ElementEncapsulatedOCSPValue implements definition.XAdESElement. Ports getElementEncapsulatedOCSPValue().
func (e XAdES111Element) ElementEncapsulatedOCSPValue() common.DSSElement {
	return XAdES111ElementEncapsulatedOCSPValue
}

// ElementEncapsulatedPKIData implements definition.XAdESElement. Ports getElementEncapsulatedPKIData().
func (e XAdES111Element) ElementEncapsulatedPKIData() common.DSSElement {
	return XAdES111ElementEncapsulatedPKIData
}

// ElementEncapsulatedTimeStamp implements definition.XAdESElement. Ports getElementEncapsulatedTimeStamp().
func (e XAdES111Element) ElementEncapsulatedTimeStamp() common.DSSElement {
	return XAdES111ElementEncapsulatedTimestamp
}

// ElementEncapsulatedX509Certificate implements definition.XAdESElement. Ports getElementEncapsulatedX509Certificate().
func (e XAdES111Element) ElementEncapsulatedX509Certificate() common.DSSElement {
	return XAdES111ElementEncapsulatedX509Certificate
}

// ElementEncoding implements definition.XAdESElement. Ports getElementEncoding().
func (e XAdES111Element) ElementEncoding() common.DSSElement {
	return XAdES111ElementEncoding
}

// ElementExplicitText implements definition.XAdESElement. Ports getElementExplicitText().
func (e XAdES111Element) ElementExplicitText() common.DSSElement {
	return XAdES111ElementExplicitText
}

// ElementIdentifier implements definition.XAdESElement. Ports getElementIdentifier().
func (e XAdES111Element) ElementIdentifier() common.DSSElement {
	return XAdES111ElementIdentifier
}

// ElementInclude implements definition.XAdESElement. Ports getElementInclude().
func (e XAdES111Element) ElementInclude() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementIndividualDataObjectsTimeStamp implements definition.XAdESElement. Ports getElementIndividualDataObjectsTimeStamp().
func (e XAdES111Element) ElementIndividualDataObjectsTimeStamp() common.DSSElement {
	return XAdES111ElementIndividualDataObjectsTimestamp
}

// Elementint implements definition.XAdESElement. Ports getElementint().
func (e XAdES111Element) Elementint() common.DSSElement {
	return XAdES111ElementInt
}

// ElementIssueTime implements definition.XAdESElement. Ports getElementIssueTime().
func (e XAdES111Element) ElementIssueTime() common.DSSElement {
	return XAdES111ElementIssueTime
}

// ElementIssuer implements definition.XAdESElement. Ports getElementIssuer().
func (e XAdES111Element) ElementIssuer() common.DSSElement {
	return XAdES111ElementIssuer
}

// ElementIssuerSerial implements definition.XAdESElement. Ports getElementIssuerSerial().
func (e XAdES111Element) ElementIssuerSerial() common.DSSElement {
	return XAdES111ElementIssuerSerial
}

// ElementIssuerSerialV2 implements definition.XAdESElement. Ports getElementIssuerSerialV2().
func (e XAdES111Element) ElementIssuerSerialV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementMimeType implements definition.XAdESElement. Ports getElementMimeType().
func (e XAdES111Element) ElementMimeType() common.DSSElement {
	return XAdES111ElementMIMEType
}

// ElementNoticeNumbers implements definition.XAdESElement. Ports getElementNoticeNumbers().
func (e XAdES111Element) ElementNoticeNumbers() common.DSSElement {
	return XAdES111ElementNoticeNumbers
}

// ElementNoticeRef implements definition.XAdESElement. Ports getElementNoticeRef().
func (e XAdES111Element) ElementNoticeRef() common.DSSElement {
	return XAdES111ElementNoticeRef
}

// ElementNumber implements definition.XAdESElement. Ports getElementNumber().
func (e XAdES111Element) ElementNumber() common.DSSElement {
	return XAdES111ElementNumber
}

// ElementObjectIdentifier implements definition.XAdESElement. Ports getElementObjectIdentifier().
func (e XAdES111Element) ElementObjectIdentifier() common.DSSElement {
	return XAdES111ElementObjectIdentifier
}

// ElementObjectReference implements definition.XAdESElement. Ports getElementObjectReference().
func (e XAdES111Element) ElementObjectReference() common.DSSElement {
	return XAdES111ElementObjectReference
}

// ElementOCSPIdentifier implements definition.XAdESElement. Ports getElementOCSPIdentifier().
func (e XAdES111Element) ElementOCSPIdentifier() common.DSSElement {
	return XAdES111ElementOCSPIdentifier
}

// ElementOCSPRef implements definition.XAdESElement. Ports getElementOCSPRef().
func (e XAdES111Element) ElementOCSPRef() common.DSSElement {
	return XAdES111ElementOCSPRef
}

// ElementOCSPRefs implements definition.XAdESElement. Ports getElementOCSPRefs().
func (e XAdES111Element) ElementOCSPRefs() common.DSSElement {
	return XAdES111ElementOCSPRefs
}

// ElementOCSPValues implements definition.XAdESElement. Ports getElementOCSPValues().
func (e XAdES111Element) ElementOCSPValues() common.DSSElement {
	return XAdES111ElementOCSPValues
}

// ElementOrganization implements definition.XAdESElement. Ports getElementOrganization().
func (e XAdES111Element) ElementOrganization() common.DSSElement {
	return XAdES111ElementOrganization
}

// ElementOtherAttributeCertificate implements definition.XAdESElement. Ports getElementOtherAttributeCertificate().
func (e XAdES111Element) ElementOtherAttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementOtherCertificate implements definition.XAdESElement. Ports getElementOtherCertificate().
func (e XAdES111Element) ElementOtherCertificate() common.DSSElement {
	return XAdES111ElementOtherCertificate
}

// ElementOtherRef implements definition.XAdESElement. Ports getElementOtherRef().
func (e XAdES111Element) ElementOtherRef() common.DSSElement {
	return XAdES111ElementOtherRef
}

// ElementOtherRefs implements definition.XAdESElement. Ports getElementOtherRefs().
func (e XAdES111Element) ElementOtherRefs() common.DSSElement {
	return XAdES111ElementOtherRefs
}

// ElementOtherTimeStamp implements definition.XAdESElement. Ports getElementOtherTimeStamp().
func (e XAdES111Element) ElementOtherTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementOtherValue implements definition.XAdESElement. Ports getElementOtherValue().
func (e XAdES111Element) ElementOtherValue() common.DSSElement {
	return XAdES111ElementOtherValue
}

// ElementOtherValues implements definition.XAdESElement. Ports getElementOtherValues().
func (e XAdES111Element) ElementOtherValues() common.DSSElement {
	return XAdES111ElementOtherValues
}

// ElementPostalCode implements definition.XAdESElement. Ports getElementPostalCode().
func (e XAdES111Element) ElementPostalCode() common.DSSElement {
	return XAdES111ElementPostalCode
}

// ElementProducedAt implements definition.XAdESElement. Ports getElementProducedAt().
func (e XAdES111Element) ElementProducedAt() common.DSSElement {
	return XAdES111ElementProducedAt
}

// ElementQualifyingProperties implements definition.XAdESElement. Ports getElementQualifyingProperties().
func (e XAdES111Element) ElementQualifyingProperties() common.DSSElement {
	return XAdES111ElementQualifyingProperties
}

// ElementQualifyingPropertiesReference implements definition.XAdESElement. Ports getElementQualifyingPropertiesReference().
func (e XAdES111Element) ElementQualifyingPropertiesReference() common.DSSElement {
	return XAdES111ElementQualifyingPropertiesReference
}

// ElementReferenceInfo implements definition.XAdESElement. Ports getElementReferenceInfo().
func (e XAdES111Element) ElementReferenceInfo() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementRefsOnlyTimeStamp implements definition.XAdESElement. Ports getElementRefsOnlyTimeStamp().
func (e XAdES111Element) ElementRefsOnlyTimeStamp() common.DSSElement {
	return XAdES111ElementRefsOnlyTimestamp
}

// ElementResponderID implements definition.XAdESElement. Ports getElementResponderID().
func (e XAdES111Element) ElementResponderID() common.DSSElement {
	return XAdES111ElementResponderID
}

// ElementRevocationValues implements definition.XAdESElement. Ports getElementRevocationValues().
func (e XAdES111Element) ElementRevocationValues() common.DSSElement {
	return XAdES111ElementRevocationValues
}

// ElementSigAndRefsTimeStamp implements definition.XAdESElement. Ports getElementSigAndRefsTimeStamp().
func (e XAdES111Element) ElementSigAndRefsTimeStamp() common.DSSElement {
	return XAdES111ElementSigAndRefsTimestamp
}

// ElementSigPolicyHash implements definition.XAdESElement. Ports getElementSigPolicyHash().
func (e XAdES111Element) ElementSigPolicyHash() common.DSSElement {
	return XAdES111ElementSigPolicyHash
}

// ElementSigPolicyId implements definition.XAdESElement. Ports getElementSigPolicyId().
func (e XAdES111Element) ElementSigPolicyId() common.DSSElement {
	return XAdES111ElementSigPolicyID
}

// ElementSigPolicyQualifier implements definition.XAdESElement. Ports getElementSigPolicyQualifier().
func (e XAdES111Element) ElementSigPolicyQualifier() common.DSSElement {
	return XAdES111ElementSigPolicyQualifier
}

// ElementSigPolicyQualifiers implements definition.XAdESElement. Ports getElementSigPolicyQualifiers().
func (e XAdES111Element) ElementSigPolicyQualifiers() common.DSSElement {
	return XAdES111ElementSigPolicyQualifiers
}

// ElementSignaturePolicyId implements definition.XAdESElement. Ports getElementSignaturePolicyId().
func (e XAdES111Element) ElementSignaturePolicyId() common.DSSElement {
	return XAdES111ElementSignaturePolicyID
}

// ElementSignaturePolicyIdentifier implements definition.XAdESElement. Ports getElementSignaturePolicyIdentifier().
func (e XAdES111Element) ElementSignaturePolicyIdentifier() common.DSSElement {
	return XAdES111ElementSignaturePolicyIdentifier
}

// ElementSignaturePolicyImplied implements definition.XAdESElement. Ports getElementSignaturePolicyImplied().
func (e XAdES111Element) ElementSignaturePolicyImplied() common.DSSElement {
	return XAdES111ElementSignaturePolicyImplied
}

// ElementSignatureProductionPlace implements definition.XAdESElement. Ports getElementSignatureProductionPlace().
func (e XAdES111Element) ElementSignatureProductionPlace() common.DSSElement {
	return XAdES111ElementSignatureProductionPlace
}

// ElementSignatureProductionPlaceV2 implements definition.XAdESElement. Ports getElementSignatureProductionPlaceV2().
func (e XAdES111Element) ElementSignatureProductionPlaceV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementSignatureTimeStamp implements definition.XAdESElement. Ports getElementSignatureTimeStamp().
func (e XAdES111Element) ElementSignatureTimeStamp() common.DSSElement {
	return XAdES111ElementSignatureTimestamp
}

// ElementSignedAssertion implements definition.XAdESElement. Ports getElementSignedAssertion().
func (e XAdES111Element) ElementSignedAssertion() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementSignedAssertions implements definition.XAdESElement. Ports getElementSignedAssertions().
func (e XAdES111Element) ElementSignedAssertions() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementSignedDataObjectProperties implements definition.XAdESElement. Ports getElementSignedDataObjectProperties().
func (e XAdES111Element) ElementSignedDataObjectProperties() common.DSSElement {
	return XAdES111ElementSignedDataObjectProperties
}

// ElementSignedProperties implements definition.XAdESElement. Ports getElementSignedProperties().
func (e XAdES111Element) ElementSignedProperties() common.DSSElement {
	return XAdES111ElementSignedProperties
}

// ElementSignedSignatureProperties implements definition.XAdESElement. Ports getElementSignedSignatureProperties().
func (e XAdES111Element) ElementSignedSignatureProperties() common.DSSElement {
	return XAdES111ElementSignedSignatureProperties
}

// ElementSignerRole implements definition.XAdESElement. Ports getElementSignerRole().
func (e XAdES111Element) ElementSignerRole() common.DSSElement {
	return XAdES111ElementSignerRole
}

// ElementSignerRoleV2 implements definition.XAdESElement. Ports getElementSignerRoleV2().
func (e XAdES111Element) ElementSignerRoleV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementSigningCertificate implements definition.XAdESElement. Ports getElementSigningCertificate().
func (e XAdES111Element) ElementSigningCertificate() common.DSSElement {
	return XAdES111ElementSigningCertificate
}

// ElementSigningCertificateV2 implements definition.XAdESElement. Ports getElementSigningCertificateV2().
func (e XAdES111Element) ElementSigningCertificateV2() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementSigningTime implements definition.XAdESElement. Ports getElementSigningTime().
func (e XAdES111Element) ElementSigningTime() common.DSSElement {
	return XAdES111ElementSigningTime
}

// ElementSPURI implements definition.XAdESElement. Ports getElementSPURI().
func (e XAdES111Element) ElementSPURI() common.DSSElement {
	return XAdES111ElementSPURI
}

// ElementSPUserNotice implements definition.XAdESElement. Ports getElementSPUserNotice().
func (e XAdES111Element) ElementSPUserNotice() common.DSSElement {
	return XAdES111ElementSPUserNotice
}

// ElementStateOrProvince implements definition.XAdESElement. Ports getElementStateOrProvince().
func (e XAdES111Element) ElementStateOrProvince() common.DSSElement {
	return XAdES111ElementStateOrProvince
}

// ElementStreetAddress implements definition.XAdESElement. Ports getElementStreetAddress().
func (e XAdES111Element) ElementStreetAddress() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementUnsignedDataObjectProperties implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperties().
func (e XAdES111Element) ElementUnsignedDataObjectProperties() common.DSSElement {
	return XAdES111ElementUnsignedDataObjectProperties
}

// ElementUnsignedDataObjectProperty implements definition.XAdESElement. Ports getElementUnsignedDataObjectProperty().
func (e XAdES111Element) ElementUnsignedDataObjectProperty() common.DSSElement {
	return XAdES111ElementUnsignedDataObjectProperty
}

// ElementUnsignedProperties implements definition.XAdESElement. Ports getElementUnsignedProperties().
func (e XAdES111Element) ElementUnsignedProperties() common.DSSElement {
	return XAdES111ElementUnsignedProperties
}

// ElementUnsignedSignatureProperties implements definition.XAdESElement. Ports getElementUnsignedSignatureProperties().
func (e XAdES111Element) ElementUnsignedSignatureProperties() common.DSSElement {
	return XAdES111ElementUnsignedSignatureProperties
}

// ElementX509AttributeCertificate implements definition.XAdESElement. Ports getElementX509AttributeCertificate().
func (e XAdES111Element) ElementX509AttributeCertificate() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementXAdESTimeStamp implements definition.XAdESElement. Ports getElementXAdESTimeStamp().
func (e XAdES111Element) ElementXAdESTimeStamp() common.DSSElement {
	panic("Element not supported by " + XAdESNamespaceXAdES111.Uri())
}

// ElementXMLTimeStamp implements definition.XAdESElement. Ports getElementXMLTimeStamp().
func (e XAdES111Element) ElementXMLTimeStamp() common.DSSElement {
	return XAdES111ElementXMLTimestamp
}

var _ XAdESElement = XAdES111Element("")
