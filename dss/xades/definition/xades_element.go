// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/XAdESElement.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdESElement is a XAdES element definition, shared across the XAdES 1.1.1/1.2.2/1.3.2
// schema versions. Each concrete version (XAdES111Element, XAdES122Element, XAdES132Element)
// implements every method below by returning its own same-named constant, or - for methods
// a given schema version does not define - by panicking with the Java
// UnsupportedOperationException message (see each version's file for its list).
type XAdESElement interface {
	common.DSSElement

	// ElementAllDataObjectsTimeStamp gets the "AllDataObjectsTimeStamp" element.
	ElementAllDataObjectsTimeStamp() common.DSSElement

	// ElementAllSignedDataObjects gets the "AllSignedDataObjects" element.
	ElementAllSignedDataObjects() common.DSSElement

	// ElementAny gets the "Any" element.
	ElementAny() common.DSSElement

	// ElementArchiveTimeStamp gets the "ArchiveTimeStamp" element.
	ElementArchiveTimeStamp() common.DSSElement

	// ElementAttrAuthoritiesCertValues gets the "AttrAuthoritiesCertValues" element.
	ElementAttrAuthoritiesCertValues() common.DSSElement

	// ElementAttributeCertificateRefs gets the "AttributeCertificateRefs" element.
	ElementAttributeCertificateRefs() common.DSSElement

	// ElementAttributeRevocationRefs gets the "AttributeRevocationRefs" element.
	ElementAttributeRevocationRefs() common.DSSElement

	// ElementAttributeRevocationValues gets the "AttributeRevocationValues" element.
	ElementAttributeRevocationValues() common.DSSElement

	// ElementByKey gets the "ByKey" element.
	ElementByKey() common.DSSElement

	// ElementByName gets the "ByName" element.
	ElementByName() common.DSSElement

	// ElementCert gets the "Cert" element.
	ElementCert() common.DSSElement

	// ElementCertDigest gets the "CertDigest" element.
	ElementCertDigest() common.DSSElement

	// ElementCertRefs gets the "CertRefs" element.
	ElementCertRefs() common.DSSElement

	// ElementCertificateValues gets the "CertificateValues" element.
	ElementCertificateValues() common.DSSElement

	// ElementCertifiedRole gets the "CertifiedRole" element.
	ElementCertifiedRole() common.DSSElement

	// ElementCertifiedRoles gets the "CertifiedRoles" element.
	ElementCertifiedRoles() common.DSSElement

	// ElementCertifiedRolesV2 gets the "CertifiedRolesV2" element.
	ElementCertifiedRolesV2() common.DSSElement

	// ElementCity gets the "City" element.
	ElementCity() common.DSSElement

	// ElementClaimedRole gets the "ClaimedRole" element.
	ElementClaimedRole() common.DSSElement

	// ElementClaimedRoles gets the "ClaimedRoles" element.
	ElementClaimedRoles() common.DSSElement

	// ElementCommitmentTypeId gets the "CommitmentTypeId" element.
	ElementCommitmentTypeId() common.DSSElement

	// ElementCommitmentTypeIndication gets the "CommitmentTypeIndication" element.
	ElementCommitmentTypeIndication() common.DSSElement

	// ElementCommitmentTypeQualifier gets the "CommitmentTypeQualifier" element.
	ElementCommitmentTypeQualifier() common.DSSElement

	// ElementCommitmentTypeQualifiers gets the "CommitmentTypeQualifies" element.
	ElementCommitmentTypeQualifiers() common.DSSElement

	// ElementCompleteCertificateRefs gets the "CompleteCertificateRefs" element.
	ElementCompleteCertificateRefs() common.DSSElement

	// ElementCompleteRevocationRefs gets the "CompleteRevocationRefs" element.
	ElementCompleteRevocationRefs() common.DSSElement

	// ElementCounterSignature gets the "CounterSignature" element.
	ElementCounterSignature() common.DSSElement

	// ElementCountryName gets the "CountryName" element.
	ElementCountryName() common.DSSElement

	// ElementCRLIdentifier gets the "CRLIdentifier" element.
	ElementCRLIdentifier() common.DSSElement

	// ElementCRLRef gets the "CRLRef" element.
	ElementCRLRef() common.DSSElement

	// ElementCRLRefs gets the "CRLRefs" element.
	ElementCRLRefs() common.DSSElement

	// ElementCRLValues gets the "CRLValues" element.
	ElementCRLValues() common.DSSElement

	// ElementDataObjectFormat gets the "DataObjectFormat" element.
	ElementDataObjectFormat() common.DSSElement

	// ElementDescription gets the "Description" element.
	ElementDescription() common.DSSElement

	// ElementDigestAlgAndValue gets the "DigestAlgAndValue" element.
	ElementDigestAlgAndValue() common.DSSElement

	// ElementDocumentationReference gets the "DocumentationReference" element.
	ElementDocumentationReference() common.DSSElement

	// ElementDocumentationReferences gets the "DocumentationReferences" element.
	ElementDocumentationReferences() common.DSSElement

	// ElementEncapsulatedCRLValue gets the "EncapsulatedCRLValue" element.
	ElementEncapsulatedCRLValue() common.DSSElement

	// ElementEncapsulatedOCSPValue gets the "EncapsulatedOCSPValue" element.
	ElementEncapsulatedOCSPValue() common.DSSElement

	// ElementEncapsulatedPKIData gets the "EncapsulatedPKIData" element.
	ElementEncapsulatedPKIData() common.DSSElement

	// ElementEncapsulatedTimeStamp gets the "EncapsulatedTimeStamp" element.
	ElementEncapsulatedTimeStamp() common.DSSElement

	// ElementEncapsulatedX509Certificate gets the "EncapsulatedX509Certificate" element.
	ElementEncapsulatedX509Certificate() common.DSSElement

	// ElementEncoding gets the "Encoding" element.
	ElementEncoding() common.DSSElement

	// ElementExplicitText gets the "ExplicitText" element.
	ElementExplicitText() common.DSSElement

	// ElementIdentifier gets the "Identifier" element.
	ElementIdentifier() common.DSSElement

	// ElementInclude gets the "Include" element.
	ElementInclude() common.DSSElement

	// ElementIndividualDataObjectsTimeStamp gets the "IndividualDataObjectsTimeStamp" element.
	ElementIndividualDataObjectsTimeStamp() common.DSSElement

	// Elementint gets the "int" element.
	Elementint() common.DSSElement

	// ElementIssueTime gets the "IssueTime" element.
	ElementIssueTime() common.DSSElement

	// ElementIssuer gets the "IssueTime" element.
	ElementIssuer() common.DSSElement

	// ElementIssuerSerial gets the "IssuerSerial" element.
	ElementIssuerSerial() common.DSSElement

	// ElementIssuerSerialV2 gets the "IssuerSerialV2" element.
	ElementIssuerSerialV2() common.DSSElement

	// ElementMimeType gets the "MimeType" element.
	ElementMimeType() common.DSSElement

	// ElementNoticeNumbers gets the "NoticeNumbers" element.
	ElementNoticeNumbers() common.DSSElement

	// ElementNoticeRef gets the "NoticeRef" element.
	ElementNoticeRef() common.DSSElement

	// ElementNumber gets the "Number" element.
	ElementNumber() common.DSSElement

	// ElementObjectIdentifier gets the "ObjectIdentifier" element.
	ElementObjectIdentifier() common.DSSElement

	// ElementObjectReference gets the "ObjectReference" element.
	ElementObjectReference() common.DSSElement

	// ElementOCSPIdentifier gets the "OCSPIdentifier" element.
	ElementOCSPIdentifier() common.DSSElement

	// ElementOCSPRef gets the "OCSPRef" element.
	ElementOCSPRef() common.DSSElement

	// ElementOCSPRefs gets the "OCSPRefs" element.
	ElementOCSPRefs() common.DSSElement

	// ElementOCSPValues gets the "OCSPValues" element.
	ElementOCSPValues() common.DSSElement

	// ElementOrganization gets the "Organization" element.
	ElementOrganization() common.DSSElement

	// ElementOtherAttributeCertificate gets the "OtherAttributeCertificate" element.
	ElementOtherAttributeCertificate() common.DSSElement

	// ElementOtherCertificate gets the "OtherCertificate" element.
	ElementOtherCertificate() common.DSSElement

	// ElementOtherRef gets the "OtherRef" element.
	ElementOtherRef() common.DSSElement

	// ElementOtherRefs gets the "OtherRefs" element.
	ElementOtherRefs() common.DSSElement

	// ElementOtherTimeStamp gets the "OtherTimeStamp" element.
	ElementOtherTimeStamp() common.DSSElement

	// ElementOtherValue gets the "OtherValue" element.
	ElementOtherValue() common.DSSElement

	// ElementOtherValues gets the "OtherValues" element.
	ElementOtherValues() common.DSSElement

	// ElementPostalCode gets the "PostalCode" element.
	ElementPostalCode() common.DSSElement

	// ElementProducedAt gets the "ProducedAt" element.
	ElementProducedAt() common.DSSElement

	// ElementQualifyingProperties gets the "QualifyingProperties" element.
	ElementQualifyingProperties() common.DSSElement

	// ElementQualifyingPropertiesReference gets the "QualifyingPropertiesReference" element.
	ElementQualifyingPropertiesReference() common.DSSElement

	// ElementReferenceInfo gets the "ReferenceInfo" element.
	ElementReferenceInfo() common.DSSElement

	// ElementRefsOnlyTimeStamp gets the "RefsOnlyTimeStamp" element.
	ElementRefsOnlyTimeStamp() common.DSSElement

	// ElementResponderID gets the "ResponderID" element.
	ElementResponderID() common.DSSElement

	// ElementRevocationValues gets the "RevocationValues" element.
	ElementRevocationValues() common.DSSElement

	// ElementSigAndRefsTimeStamp gets the "SigAndRefsTimeStamp" element.
	ElementSigAndRefsTimeStamp() common.DSSElement

	// ElementSigPolicyHash gets the "SigPolicyHash" element.
	ElementSigPolicyHash() common.DSSElement

	// ElementSigPolicyId gets the "SigPolicyId" element.
	ElementSigPolicyId() common.DSSElement

	// ElementSigPolicyQualifier gets the "SigPolicyQualifier" element.
	ElementSigPolicyQualifier() common.DSSElement

	// ElementSigPolicyQualifiers gets the "SigPolicyQualifiers" element.
	ElementSigPolicyQualifiers() common.DSSElement

	// ElementSignaturePolicyId gets the "SignaturePolicyId" element.
	ElementSignaturePolicyId() common.DSSElement

	// ElementSignaturePolicyIdentifier gets the "SignaturePolicyIdentifier" element.
	ElementSignaturePolicyIdentifier() common.DSSElement

	// ElementSignaturePolicyImplied gets the "SignaturePolicyImplied" element.
	ElementSignaturePolicyImplied() common.DSSElement

	// ElementSignatureProductionPlace gets the "SignatureProductionPlace" element.
	ElementSignatureProductionPlace() common.DSSElement

	// ElementSignatureProductionPlaceV2 gets the "SignatureProductionPlaceV2" element.
	ElementSignatureProductionPlaceV2() common.DSSElement

	// ElementSignatureTimeStamp gets the "SignatureTimeStamp" element.
	ElementSignatureTimeStamp() common.DSSElement

	// ElementSignedAssertion gets the "SignedAssertion" element.
	ElementSignedAssertion() common.DSSElement

	// ElementSignedAssertions gets the "SignedAssertions" element.
	ElementSignedAssertions() common.DSSElement

	// ElementSignedDataObjectProperties gets the "SignedDataObjectProperties" element.
	ElementSignedDataObjectProperties() common.DSSElement

	// ElementSignedProperties gets the "SignedProperties" element.
	ElementSignedProperties() common.DSSElement

	// ElementSignedSignatureProperties gets the "SignedSignatureProperties" element.
	ElementSignedSignatureProperties() common.DSSElement

	// ElementSignerRole gets the "SignerRole" element.
	ElementSignerRole() common.DSSElement

	// ElementSignerRoleV2 gets the "SignerRoleV2" element.
	ElementSignerRoleV2() common.DSSElement

	// ElementSigningCertificate gets the "SigningCertificate" element.
	ElementSigningCertificate() common.DSSElement

	// ElementSigningCertificateV2 gets the "SigningCertificateV2" element.
	ElementSigningCertificateV2() common.DSSElement

	// ElementSigningTime gets the "SigningTime" element.
	ElementSigningTime() common.DSSElement

	// ElementSPURI gets the "SPURI" element.
	ElementSPURI() common.DSSElement

	// ElementSPUserNotice gets the "SPUserNotice" element.
	ElementSPUserNotice() common.DSSElement

	// ElementStateOrProvince gets the "StateOrProvince" element.
	ElementStateOrProvince() common.DSSElement

	// ElementStreetAddress gets the "StreetAddress" element.
	ElementStreetAddress() common.DSSElement

	// ElementUnsignedDataObjectProperties gets the "UnsignedDataObjectProperties" element.
	ElementUnsignedDataObjectProperties() common.DSSElement

	// ElementUnsignedDataObjectProperty gets the "UnsignedDataObjectProperty" element.
	ElementUnsignedDataObjectProperty() common.DSSElement

	// ElementUnsignedProperties gets the "UnsignedProperties" element.
	ElementUnsignedProperties() common.DSSElement

	// ElementUnsignedSignatureProperties gets the "UnsignedSignatureProperties" element.
	ElementUnsignedSignatureProperties() common.DSSElement

	// ElementX509AttributeCertificate gets the "X509AttributeCertificate" element.
	ElementX509AttributeCertificate() common.DSSElement

	// ElementXAdESTimeStamp gets the "XAdESTimeStamp" element.
	ElementXAdESTimeStamp() common.DSSElement

	// ElementXMLTimeStamp gets the "XMLTimeStamp" element.
	ElementXMLTimeStamp() common.DSSElement
}
