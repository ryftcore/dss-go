// KAT test for xades132_element.go: every XAdES132Element_* tag name below is transcribed
// verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.xades132.XAdES132Element, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestXAdES132Element_TagNameKAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES132ElementAllDataObjectsTimestamp.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES132ElementAllSignedDataObjects.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES132ElementAny.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES132ElementArchiveTimestamp.TagName(), "ArchiveTimeStamp"},
		{"ATTR_AUTHORITIES_CERT_VALUES", XAdES132ElementAttrAuthoritiesCertValues.TagName(), "AttrAuthoritiesCertValues"},
		{"ATTRIBUTE_CERTIFICATE_REFS", XAdES132ElementAttributeCertificateRefs.TagName(), "AttributeCertificateRefs"},
		{"ATTRIBUTE_REVOCATION_REFS", XAdES132ElementAttributeRevocationRefs.TagName(), "AttributeRevocationRefs"},
		{"ATTRIBUTE_REVOCATION_VALUES", XAdES132ElementAttributeRevocationValues.TagName(), "AttributeRevocationValues"},
		{"BY_KEY", XAdES132ElementByKey.TagName(), "ByKey"},
		{"BY_NAME", XAdES132ElementByName.TagName(), "ByName"},
		{"CERT", XAdES132ElementCert.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES132ElementCertDigest.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES132ElementCertRefs.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES132ElementCertificateValues.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES132ElementCertifiedRole.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES132ElementCertifiedRoles.TagName(), "CertifiedRoles"},
		{"CERTIFIED_ROLES_V2", XAdES132ElementCertifiedRolesV2.TagName(), "CertifiedRolesV2"},
		{"CITY", XAdES132ElementCity.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES132ElementClaimedRole.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES132ElementClaimedRoles.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES132ElementCommitmentTypeID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES132ElementCommitmentTypeIndication.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES132ElementCommitmentTypeQualifier.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES132ElementCommitmentTypeQualifiers.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES132ElementCompleteCertificateRefs.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES132ElementCompleteRevocationRefs.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES132ElementCounterSignature.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES132ElementCountryName.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES132ElementCRLIdentifier.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES132ElementCRLRef.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES132ElementCRLRefs.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES132ElementCRLValues.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES132ElementDataObjectFormat.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES132ElementDescription.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES132ElementDigestAlgAndValue.TagName(), "DigestAlgAndValue"},
		{"DOCUMENTATION_REFERENCE", XAdES132ElementDocumentationReference.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES132ElementDocumentationReferences.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES132ElementEncapsulatedCRLValue.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES132ElementEncapsulatedOCSPValue.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES132ElementEncapsulatedPKIData.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES132ElementEncapsulatedTimestamp.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES132ElementEncapsulatedX509Certificate.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES132ElementEncoding.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES132ElementExplicitText.TagName(), "ExplicitText"},
		{"IDENTIFIER", XAdES132ElementIdentifier.TagName(), "Identifier"},
		{"INCLUDE", XAdES132ElementInclude.TagName(), "Include"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES132ElementIndividualDataObjectsTimestamp.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES132ElementInt.TagName(), "int"},
		{"ISSUE_TIME", XAdES132ElementIssueTime.TagName(), "IssueTime"},
		{"ISSUER", XAdES132ElementIssuer.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES132ElementIssuerSerial.TagName(), "IssuerSerial"},
		{"ISSUER_SERIAL_V2", XAdES132ElementIssuerSerialV2.TagName(), "IssuerSerialV2"},
		{"MIME_TYPE", XAdES132ElementMIMEType.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES132ElementNoticeNumbers.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES132ElementNoticeRef.TagName(), "NoticeRef"},
		{"NUMBER", XAdES132ElementNumber.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES132ElementObjectIdentifier.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES132ElementObjectReference.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES132ElementOCSPIdentifier.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES132ElementOCSPRef.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES132ElementOCSPRefs.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES132ElementOCSPValues.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES132ElementOrganization.TagName(), "Organization"},
		{"OTHER_ATTRIBUTE_CERTIFICATE", XAdES132ElementOtherAttributeCertificate.TagName(), "OtherAttributeCertificate"},
		{"OTHER_CERTIFICATE", XAdES132ElementOtherCertificate.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES132ElementOtherRef.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES132ElementOtherRefs.TagName(), "OtherRefs"},
		{"OTHER_TIMESTAMP", XAdES132ElementOtherTimestamp.TagName(), "OtherTimeStamp"},
		{"OTHER_VALUE", XAdES132ElementOtherValue.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES132ElementOtherValues.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES132ElementPostalCode.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES132ElementProducedAt.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES132ElementQualifyingProperties.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES132ElementQualifyingPropertiesReference.TagName(), "QualifyingPropertiesReference"},
		{"REFERENCE_INFO", XAdES132ElementReferenceInfo.TagName(), "ReferenceInfo"},
		{"REFS_ONLY_TIMESTAMP", XAdES132ElementRefsOnlyTimestamp.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES132ElementResponderID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES132ElementRevocationValues.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES132ElementSigAndRefsTimestamp.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES132ElementSigPolicyHash.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES132ElementSigPolicyID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES132ElementSigPolicyQualifier.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES132ElementSigPolicyQualifiers.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES132ElementSignaturePolicyID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES132ElementSignaturePolicyIdentifier.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES132ElementSignaturePolicyImplied.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES132ElementSignatureProductionPlace.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_PRODUCTION_PLACE_V2", XAdES132ElementSignatureProductionPlaceV2.TagName(), "SignatureProductionPlaceV2"},
		{"SIGNATURE_TIMESTAMP", XAdES132ElementSignatureTimestamp.TagName(), "SignatureTimeStamp"},
		{"SIGNED_ASSERTION", XAdES132ElementSignedAssertion.TagName(), "SignedAssertion"},
		{"SIGNED_ASSERTIONS", XAdES132ElementSignedAssertions.TagName(), "SignedAssertions"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES132ElementSignedDataObjectProperties.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES132ElementSignedProperties.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES132ElementSignedSignatureProperties.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES132ElementSignerRole.TagName(), "SignerRole"},
		{"SIGNER_ROLE_V2", XAdES132ElementSignerRoleV2.TagName(), "SignerRoleV2"},
		{"SIGNING_CERTIFICATE", XAdES132ElementSigningCertificate.TagName(), "SigningCertificate"},
		{"SIGNING_CERTIFICATE_V2", XAdES132ElementSigningCertificateV2.TagName(), "SigningCertificateV2"},
		{"SIGNING_TIME", XAdES132ElementSigningTime.TagName(), "SigningTime"},
		{"SP_URI", XAdES132ElementSPURI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES132ElementSPUserNotice.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES132ElementStateOrProvince.TagName(), "StateOrProvince"},
		{"STREET_ADDRESS", XAdES132ElementStreetAddress.TagName(), "StreetAddress"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES132ElementUnsignedDataObjectProperties.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES132ElementUnsignedDataObjectProperty.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES132ElementUnsignedProperties.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES132ElementUnsignedSignatureProperties.TagName(), "UnsignedSignatureProperties"},
		{"X509_ATTRIBUTE_CERTIFICATE", XAdES132ElementX509AttributeCertificate.TagName(), "X509AttributeCertificate"},
		{"XADES_TIMESTAMP", XAdES132ElementXAdESTimestamp.TagName(), "XAdESTimeStamp"},
		{"XML_TIMESTAMP", XAdES132ElementXMLTimestamp.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES132Element_Namespace(t *testing.T) {
	if got := XAdES132ElementAllDataObjectsTimestamp.URI(); got != XAdESNamespaceXAdES132.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespaceXAdES132.Uri())
	}
	if !XAdES132ElementAllDataObjectsTimestamp.IsSameTagName(XAdES132ElementAllDataObjectsTimestamp.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES132Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES132Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES132Element_ElementGetters(t *testing.T) {
	var e XAdES132Element = XAdES132ElementAllDataObjectsTimestamp
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES132ElementAllDataObjectsTimestamp {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES132ElementAllDataObjectsTimestamp", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES132ElementAllSignedDataObjects {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES132ElementAllSignedDataObjects", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES132ElementAny {
			t.Errorf("ElementAny() = %v, want XAdES132ElementAny", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES132ElementArchiveTimestamp {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES132ElementArchiveTimestamp", got)
		}
	})
	t.Run("ElementAttrAuthoritiesCertValues", func(t *testing.T) {
		if got := e.ElementAttrAuthoritiesCertValues(); got != XAdES132ElementAttrAuthoritiesCertValues {
			t.Errorf("ElementAttrAuthoritiesCertValues() = %v, want XAdES132ElementAttrAuthoritiesCertValues", got)
		}
	})
	t.Run("ElementAttributeCertificateRefs", func(t *testing.T) {
		if got := e.ElementAttributeCertificateRefs(); got != XAdES132ElementAttributeCertificateRefs {
			t.Errorf("ElementAttributeCertificateRefs() = %v, want XAdES132ElementAttributeCertificateRefs", got)
		}
	})
	t.Run("ElementAttributeRevocationRefs", func(t *testing.T) {
		if got := e.ElementAttributeRevocationRefs(); got != XAdES132ElementAttributeRevocationRefs {
			t.Errorf("ElementAttributeRevocationRefs() = %v, want XAdES132ElementAttributeRevocationRefs", got)
		}
	})
	t.Run("ElementAttributeRevocationValues", func(t *testing.T) {
		if got := e.ElementAttributeRevocationValues(); got != XAdES132ElementAttributeRevocationValues {
			t.Errorf("ElementAttributeRevocationValues() = %v, want XAdES132ElementAttributeRevocationValues", got)
		}
	})
	t.Run("ElementByKey", func(t *testing.T) {
		if got := e.ElementByKey(); got != XAdES132ElementByKey {
			t.Errorf("ElementByKey() = %v, want XAdES132ElementByKey", got)
		}
	})
	t.Run("ElementByName", func(t *testing.T) {
		if got := e.ElementByName(); got != XAdES132ElementByName {
			t.Errorf("ElementByName() = %v, want XAdES132ElementByName", got)
		}
	})
	t.Run("ElementCert", func(t *testing.T) {
		if got := e.ElementCert(); got != XAdES132ElementCert {
			t.Errorf("ElementCert() = %v, want XAdES132ElementCert", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES132ElementCertDigest {
			t.Errorf("ElementCertDigest() = %v, want XAdES132ElementCertDigest", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES132ElementCertRefs {
			t.Errorf("ElementCertRefs() = %v, want XAdES132ElementCertRefs", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES132ElementCertificateValues {
			t.Errorf("ElementCertificateValues() = %v, want XAdES132ElementCertificateValues", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES132ElementCertifiedRole {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES132ElementCertifiedRole", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES132ElementCertifiedRoles {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES132ElementCertifiedRoles", got)
		}
	})
	t.Run("ElementCertifiedRolesV2", func(t *testing.T) {
		if got := e.ElementCertifiedRolesV2(); got != XAdES132ElementCertifiedRolesV2 {
			t.Errorf("ElementCertifiedRolesV2() = %v, want XAdES132ElementCertifiedRolesV2", got)
		}
	})
	t.Run("ElementCity", func(t *testing.T) {
		if got := e.ElementCity(); got != XAdES132ElementCity {
			t.Errorf("ElementCity() = %v, want XAdES132ElementCity", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES132ElementClaimedRole {
			t.Errorf("ElementClaimedRole() = %v, want XAdES132ElementClaimedRole", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES132ElementClaimedRoles {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES132ElementClaimedRoles", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES132ElementCommitmentTypeID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES132ElementCommitmentTypeID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES132ElementCommitmentTypeIndication {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES132ElementCommitmentTypeIndication", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES132ElementCommitmentTypeQualifier {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES132ElementCommitmentTypeQualifier", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES132ElementCommitmentTypeQualifiers {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES132ElementCommitmentTypeQualifiers", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES132ElementCompleteCertificateRefs {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES132ElementCompleteCertificateRefs", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES132ElementCompleteRevocationRefs {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES132ElementCompleteRevocationRefs", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES132ElementCounterSignature {
			t.Errorf("ElementCounterSignature() = %v, want XAdES132ElementCounterSignature", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES132ElementCountryName {
			t.Errorf("ElementCountryName() = %v, want XAdES132ElementCountryName", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES132ElementCRLIdentifier {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES132ElementCRLIdentifier", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES132ElementCRLRef {
			t.Errorf("ElementCRLRef() = %v, want XAdES132ElementCRLRef", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES132ElementCRLRefs {
			t.Errorf("ElementCRLRefs() = %v, want XAdES132ElementCRLRefs", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES132ElementCRLValues {
			t.Errorf("ElementCRLValues() = %v, want XAdES132ElementCRLValues", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES132ElementDataObjectFormat {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES132ElementDataObjectFormat", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES132ElementDescription {
			t.Errorf("ElementDescription() = %v, want XAdES132ElementDescription", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES132ElementDigestAlgAndValue {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES132ElementDigestAlgAndValue", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES132ElementDocumentationReference {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES132ElementDocumentationReference", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES132ElementDocumentationReferences {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES132ElementDocumentationReferences", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES132ElementEncapsulatedCRLValue {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES132ElementEncapsulatedCRLValue", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES132ElementEncapsulatedOCSPValue {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES132ElementEncapsulatedOCSPValue", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES132ElementEncapsulatedPKIData {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES132ElementEncapsulatedPKIData", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES132ElementEncapsulatedTimestamp {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES132ElementEncapsulatedTimestamp", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES132ElementEncapsulatedX509Certificate {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES132ElementEncapsulatedX509Certificate", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES132ElementEncoding {
			t.Errorf("ElementEncoding() = %v, want XAdES132ElementEncoding", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES132ElementExplicitText {
			t.Errorf("ElementExplicitText() = %v, want XAdES132ElementExplicitText", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES132ElementIdentifier {
			t.Errorf("ElementIdentifier() = %v, want XAdES132ElementIdentifier", got)
		}
	})
	t.Run("ElementInclude", func(t *testing.T) {
		if got := e.ElementInclude(); got != XAdES132ElementInclude {
			t.Errorf("ElementInclude() = %v, want XAdES132ElementInclude", got)
		}
	})
	t.Run("ElementIndividualDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES132ElementIndividualDataObjectsTimestamp {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES132ElementIndividualDataObjectsTimestamp", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES132ElementInt {
			t.Errorf("Elementint() = %v, want XAdES132ElementInt", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES132ElementIssueTime {
			t.Errorf("ElementIssueTime() = %v, want XAdES132ElementIssueTime", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES132ElementIssuer {
			t.Errorf("ElementIssuer() = %v, want XAdES132ElementIssuer", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES132ElementIssuerSerial {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES132ElementIssuerSerial", got)
		}
	})
	t.Run("ElementIssuerSerialV2", func(t *testing.T) {
		if got := e.ElementIssuerSerialV2(); got != XAdES132ElementIssuerSerialV2 {
			t.Errorf("ElementIssuerSerialV2() = %v, want XAdES132ElementIssuerSerialV2", got)
		}
	})
	t.Run("ElementMimeType", func(t *testing.T) {
		if got := e.ElementMimeType(); got != XAdES132ElementMIMEType {
			t.Errorf("ElementMimeType() = %v, want XAdES132ElementMIMEType", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES132ElementNoticeNumbers {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES132ElementNoticeNumbers", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES132ElementNoticeRef {
			t.Errorf("ElementNoticeRef() = %v, want XAdES132ElementNoticeRef", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES132ElementNumber {
			t.Errorf("ElementNumber() = %v, want XAdES132ElementNumber", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES132ElementObjectIdentifier {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES132ElementObjectIdentifier", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES132ElementObjectReference {
			t.Errorf("ElementObjectReference() = %v, want XAdES132ElementObjectReference", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES132ElementOCSPIdentifier {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES132ElementOCSPIdentifier", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES132ElementOCSPRef {
			t.Errorf("ElementOCSPRef() = %v, want XAdES132ElementOCSPRef", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES132ElementOCSPRefs {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES132ElementOCSPRefs", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES132ElementOCSPValues {
			t.Errorf("ElementOCSPValues() = %v, want XAdES132ElementOCSPValues", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES132ElementOrganization {
			t.Errorf("ElementOrganization() = %v, want XAdES132ElementOrganization", got)
		}
	})
	t.Run("ElementOtherAttributeCertificate", func(t *testing.T) {
		if got := e.ElementOtherAttributeCertificate(); got != XAdES132ElementOtherAttributeCertificate {
			t.Errorf("ElementOtherAttributeCertificate() = %v, want XAdES132ElementOtherAttributeCertificate", got)
		}
	})
	t.Run("ElementOtherCertificate", func(t *testing.T) {
		if got := e.ElementOtherCertificate(); got != XAdES132ElementOtherCertificate {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES132ElementOtherCertificate", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES132ElementOtherRef {
			t.Errorf("ElementOtherRef() = %v, want XAdES132ElementOtherRef", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES132ElementOtherRefs {
			t.Errorf("ElementOtherRefs() = %v, want XAdES132ElementOtherRefs", got)
		}
	})
	t.Run("ElementOtherTimeStamp", func(t *testing.T) {
		if got := e.ElementOtherTimeStamp(); got != XAdES132ElementOtherTimestamp {
			t.Errorf("ElementOtherTimeStamp() = %v, want XAdES132ElementOtherTimestamp", got)
		}
	})
	t.Run("ElementOtherValue", func(t *testing.T) {
		if got := e.ElementOtherValue(); got != XAdES132ElementOtherValue {
			t.Errorf("ElementOtherValue() = %v, want XAdES132ElementOtherValue", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES132ElementOtherValues {
			t.Errorf("ElementOtherValues() = %v, want XAdES132ElementOtherValues", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES132ElementPostalCode {
			t.Errorf("ElementPostalCode() = %v, want XAdES132ElementPostalCode", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES132ElementProducedAt {
			t.Errorf("ElementProducedAt() = %v, want XAdES132ElementProducedAt", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES132ElementQualifyingProperties {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES132ElementQualifyingProperties", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES132ElementQualifyingPropertiesReference {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES132ElementQualifyingPropertiesReference", got)
		}
	})
	t.Run("ElementReferenceInfo", func(t *testing.T) {
		if got := e.ElementReferenceInfo(); got != XAdES132ElementReferenceInfo {
			t.Errorf("ElementReferenceInfo() = %v, want XAdES132ElementReferenceInfo", got)
		}
	})
	t.Run("ElementRefsOnlyTimeStamp", func(t *testing.T) {
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES132ElementRefsOnlyTimestamp {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES132ElementRefsOnlyTimestamp", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES132ElementResponderID {
			t.Errorf("ElementResponderID() = %v, want XAdES132ElementResponderID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES132ElementRevocationValues {
			t.Errorf("ElementRevocationValues() = %v, want XAdES132ElementRevocationValues", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES132ElementSigAndRefsTimestamp {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES132ElementSigAndRefsTimestamp", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES132ElementSigPolicyHash {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES132ElementSigPolicyHash", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES132ElementSigPolicyID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES132ElementSigPolicyID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES132ElementSigPolicyQualifier {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES132ElementSigPolicyQualifier", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES132ElementSigPolicyQualifiers {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES132ElementSigPolicyQualifiers", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES132ElementSignaturePolicyID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES132ElementSignaturePolicyID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES132ElementSignaturePolicyIdentifier {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES132ElementSignaturePolicyIdentifier", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES132ElementSignaturePolicyImplied {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES132ElementSignaturePolicyImplied", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES132ElementSignatureProductionPlace {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES132ElementSignatureProductionPlace", got)
		}
	})
	t.Run("ElementSignatureProductionPlaceV2", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlaceV2(); got != XAdES132ElementSignatureProductionPlaceV2 {
			t.Errorf("ElementSignatureProductionPlaceV2() = %v, want XAdES132ElementSignatureProductionPlaceV2", got)
		}
	})
	t.Run("ElementSignatureTimeStamp", func(t *testing.T) {
		if got := e.ElementSignatureTimeStamp(); got != XAdES132ElementSignatureTimestamp {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES132ElementSignatureTimestamp", got)
		}
	})
	t.Run("ElementSignedAssertion", func(t *testing.T) {
		if got := e.ElementSignedAssertion(); got != XAdES132ElementSignedAssertion {
			t.Errorf("ElementSignedAssertion() = %v, want XAdES132ElementSignedAssertion", got)
		}
	})
	t.Run("ElementSignedAssertions", func(t *testing.T) {
		if got := e.ElementSignedAssertions(); got != XAdES132ElementSignedAssertions {
			t.Errorf("ElementSignedAssertions() = %v, want XAdES132ElementSignedAssertions", got)
		}
	})
	t.Run("ElementSignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementSignedDataObjectProperties(); got != XAdES132ElementSignedDataObjectProperties {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES132ElementSignedDataObjectProperties", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES132ElementSignedProperties {
			t.Errorf("ElementSignedProperties() = %v, want XAdES132ElementSignedProperties", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES132ElementSignedSignatureProperties {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES132ElementSignedSignatureProperties", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES132ElementSignerRole {
			t.Errorf("ElementSignerRole() = %v, want XAdES132ElementSignerRole", got)
		}
	})
	t.Run("ElementSignerRoleV2", func(t *testing.T) {
		if got := e.ElementSignerRoleV2(); got != XAdES132ElementSignerRoleV2 {
			t.Errorf("ElementSignerRoleV2() = %v, want XAdES132ElementSignerRoleV2", got)
		}
	})
	t.Run("ElementSigningCertificate", func(t *testing.T) {
		if got := e.ElementSigningCertificate(); got != XAdES132ElementSigningCertificate {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES132ElementSigningCertificate", got)
		}
	})
	t.Run("ElementSigningCertificateV2", func(t *testing.T) {
		if got := e.ElementSigningCertificateV2(); got != XAdES132ElementSigningCertificateV2 {
			t.Errorf("ElementSigningCertificateV2() = %v, want XAdES132ElementSigningCertificateV2", got)
		}
	})
	t.Run("ElementSigningTime", func(t *testing.T) {
		if got := e.ElementSigningTime(); got != XAdES132ElementSigningTime {
			t.Errorf("ElementSigningTime() = %v, want XAdES132ElementSigningTime", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES132ElementSPURI {
			t.Errorf("ElementSPURI() = %v, want XAdES132ElementSPURI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES132ElementSPUserNotice {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES132ElementSPUserNotice", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES132ElementStateOrProvince {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES132ElementStateOrProvince", got)
		}
	})
	t.Run("ElementStreetAddress", func(t *testing.T) {
		if got := e.ElementStreetAddress(); got != XAdES132ElementStreetAddress {
			t.Errorf("ElementStreetAddress() = %v, want XAdES132ElementStreetAddress", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES132ElementUnsignedDataObjectProperties {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES132ElementUnsignedDataObjectProperties", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES132ElementUnsignedDataObjectProperty {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES132ElementUnsignedDataObjectProperty", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES132ElementUnsignedProperties {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES132ElementUnsignedProperties", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES132ElementUnsignedSignatureProperties {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES132ElementUnsignedSignatureProperties", got)
		}
	})
	t.Run("ElementX509AttributeCertificate", func(t *testing.T) {
		if got := e.ElementX509AttributeCertificate(); got != XAdES132ElementX509AttributeCertificate {
			t.Errorf("ElementX509AttributeCertificate() = %v, want XAdES132ElementX509AttributeCertificate", got)
		}
	})
	t.Run("ElementXAdESTimeStamp", func(t *testing.T) {
		if got := e.ElementXAdESTimeStamp(); got != XAdES132ElementXAdESTimestamp {
			t.Errorf("ElementXAdESTimeStamp() = %v, want XAdES132ElementXAdESTimestamp", got)
		}
	})
	t.Run("ElementXMLTimeStamp", func(t *testing.T) {
		if got := e.ElementXMLTimeStamp(); got != XAdES132ElementXMLTimestamp {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES132ElementXMLTimestamp", got)
		}
	})
}
