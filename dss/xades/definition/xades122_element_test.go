// KAT test for xades122_element.go: every XAdES122Element_* tag name below is transcribed
// verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.xades122.XAdES122Element, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestXAdES122Element_TagNameKAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES122ElementAllDataObjectsTimestamp.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES122ElementAllSignedDataObjects.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES122ElementAny.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES122ElementArchiveTimestamp.TagName(), "ArchiveTimeStamp"},
		{"ATTRIBUTE_CERTIFICATE_REFS", XAdES122ElementAttributeCertificateRefs.TagName(), "AttributeCertificateRefs"},
		{"ATTRIBUTE_REVOCATION_REFS", XAdES122ElementAttributeRevocationRefs.TagName(), "AttributeRevocationRefs"},
		{"CERT", XAdES122ElementCert.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES122ElementCertDigest.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES122ElementCertRefs.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES122ElementCertificateValues.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES122ElementCertifiedRole.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES122ElementCertifiedRoles.TagName(), "CertifiedRoles"},
		{"CITY", XAdES122ElementCity.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES122ElementClaimedRole.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES122ElementClaimedRoles.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES122ElementCommitmentTypeID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES122ElementCommitmentTypeIndication.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES122ElementCommitmentTypeQualifier.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES122ElementCommitmentTypeQualifiers.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES122ElementCompleteCertificateRefs.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES122ElementCompleteRevocationRefs.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES122ElementCounterSignature.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES122ElementCountryName.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES122ElementCRLIdentifier.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES122ElementCRLRef.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES122ElementCRLRefs.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES122ElementCRLValues.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES122ElementDataObjectFormat.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES122ElementDescription.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES122ElementDigestAlgAndValue.TagName(), "DigestAlgAndValue"},
		{"DOCUMENTATION_REFERENCE", XAdES122ElementDocumentationReference.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES122ElementDocumentationReferences.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES122ElementEncapsulatedCRLValue.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES122ElementEncapsulatedOCSPValue.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES122ElementEncapsulatedPKIData.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES122ElementEncapsulatedTimestamp.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES122ElementEncapsulatedX509Certificate.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES122ElementEncoding.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES122ElementExplicitText.TagName(), "ExplicitText"},
		{"IDENTIFIER", XAdES122ElementIdentifier.TagName(), "Identifier"},
		{"INCLUDE", XAdES122ElementInclude.TagName(), "Include"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES122ElementIndividualDataObjectsTimestamp.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES122ElementInt.TagName(), "int"},
		{"ISSUE_TIME", XAdES122ElementIssueTime.TagName(), "IssueTime"},
		{"ISSUER", XAdES122ElementIssuer.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES122ElementIssuerSerial.TagName(), "IssuerSerial"},
		{"MIME_TYPE", XAdES122ElementMIMEType.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES122ElementNoticeNumbers.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES122ElementNoticeRef.TagName(), "NoticeRef"},
		{"NUMBER", XAdES122ElementNumber.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES122ElementObjectIdentifier.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES122ElementObjectReference.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES122ElementOCSPIdentifier.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES122ElementOCSPRef.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES122ElementOCSPRefs.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES122ElementOCSPValues.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES122ElementOrganization.TagName(), "Organization"},
		{"OTHER_CERTIFICATE", XAdES122ElementOtherCertificate.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES122ElementOtherRef.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES122ElementOtherRefs.TagName(), "OtherRefs"},
		{"OTHER_VALUE", XAdES122ElementOtherValue.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES122ElementOtherValues.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES122ElementPostalCode.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES122ElementProducedAt.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES122ElementQualifyingProperties.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES122ElementQualifyingPropertiesReference.TagName(), "QualifyingPropertiesReference"},
		{"REFS_ONLY_TIMESTAMP", XAdES122ElementRefsOnlyTimestamp.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES122ElementResponderID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES122ElementRevocationValues.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES122ElementSigAndRefsTimestamp.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES122ElementSigPolicyHash.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES122ElementSigPolicyID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES122ElementSigPolicyQualifier.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES122ElementSigPolicyQualifiers.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES122ElementSignaturePolicyID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES122ElementSignaturePolicyIdentifier.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES122ElementSignaturePolicyImplied.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES122ElementSignatureProductionPlace.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_TIMESTAMP", XAdES122ElementSignatureTimestamp.TagName(), "SignatureTimeStamp"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES122ElementSignedDataObjectProperties.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES122ElementSignedProperties.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES122ElementSignedSignatureProperties.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES122ElementSignerRole.TagName(), "SignerRole"},
		{"SIGNING_CERTIFICATE", XAdES122ElementSigningCertificate.TagName(), "SigningCertificate"},
		{"SIGNING_TIME", XAdES122ElementSigningTime.TagName(), "SigningTime"},
		{"SP_URI", XAdES122ElementSPURI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES122ElementSPUserNotice.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES122ElementStateOrProvince.TagName(), "StateOrProvince"},
		{"TIMESTAMP", XAdES122ElementTimestamp.TagName(), "TimeStamp"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES122ElementUnsignedDataObjectProperties.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES122ElementUnsignedDataObjectProperty.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES122ElementUnsignedProperties.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES122ElementUnsignedSignatureProperties.TagName(), "UnsignedSignatureProperties"},
		{"XML_TIMESTAMP", XAdES122ElementXMLTimestamp.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES122Element_Namespace(t *testing.T) {
	if got := XAdES122ElementAllDataObjectsTimestamp.URI(); got != XAdESNamespaceXAdES122.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespaceXAdES122.Uri())
	}
	if !XAdES122ElementAllDataObjectsTimestamp.IsSameTagName(XAdES122ElementAllDataObjectsTimestamp.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES122Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES122Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES122Element_ElementGetters(t *testing.T) {
	var e XAdES122Element = XAdES122ElementAllDataObjectsTimestamp
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES122ElementAllDataObjectsTimestamp {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES122ElementAllDataObjectsTimestamp", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES122ElementAllSignedDataObjects {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES122ElementAllSignedDataObjects", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES122ElementAny {
			t.Errorf("ElementAny() = %v, want XAdES122ElementAny", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES122ElementArchiveTimestamp {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES122ElementArchiveTimestamp", got)
		}
	})
	t.Run("ElementAttrAuthoritiesCertValues", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttrAuthoritiesCertValues() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementAttrAuthoritiesCertValues()
	})
	t.Run("ElementAttributeCertificateRefs", func(t *testing.T) {
		if got := e.ElementAttributeCertificateRefs(); got != XAdES122ElementAttributeCertificateRefs {
			t.Errorf("ElementAttributeCertificateRefs() = %v, want XAdES122ElementAttributeCertificateRefs", got)
		}
	})
	t.Run("ElementAttributeRevocationRefs", func(t *testing.T) {
		if got := e.ElementAttributeRevocationRefs(); got != XAdES122ElementAttributeRevocationRefs {
			t.Errorf("ElementAttributeRevocationRefs() = %v, want XAdES122ElementAttributeRevocationRefs", got)
		}
	})
	t.Run("ElementAttributeRevocationValues", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttributeRevocationValues() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementAttributeRevocationValues()
	})
	t.Run("ElementByKey", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementByKey() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementByKey()
	})
	t.Run("ElementByName", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementByName() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementByName()
	})
	t.Run("ElementCert", func(t *testing.T) {
		if got := e.ElementCert(); got != XAdES122ElementCert {
			t.Errorf("ElementCert() = %v, want XAdES122ElementCert", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES122ElementCertDigest {
			t.Errorf("ElementCertDigest() = %v, want XAdES122ElementCertDigest", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES122ElementCertRefs {
			t.Errorf("ElementCertRefs() = %v, want XAdES122ElementCertRefs", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES122ElementCertificateValues {
			t.Errorf("ElementCertificateValues() = %v, want XAdES122ElementCertificateValues", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES122ElementCertifiedRole {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES122ElementCertifiedRole", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES122ElementCertifiedRoles {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES122ElementCertifiedRoles", got)
		}
	})
	t.Run("ElementCertifiedRolesV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementCertifiedRolesV2() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementCertifiedRolesV2()
	})
	t.Run("ElementCity", func(t *testing.T) {
		if got := e.ElementCity(); got != XAdES122ElementCity {
			t.Errorf("ElementCity() = %v, want XAdES122ElementCity", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES122ElementClaimedRole {
			t.Errorf("ElementClaimedRole() = %v, want XAdES122ElementClaimedRole", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES122ElementClaimedRoles {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES122ElementClaimedRoles", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES122ElementCommitmentTypeID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES122ElementCommitmentTypeID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES122ElementCommitmentTypeIndication {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES122ElementCommitmentTypeIndication", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES122ElementCommitmentTypeQualifier {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES122ElementCommitmentTypeQualifier", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES122ElementCommitmentTypeQualifiers {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES122ElementCommitmentTypeQualifiers", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES122ElementCompleteCertificateRefs {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES122ElementCompleteCertificateRefs", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES122ElementCompleteRevocationRefs {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES122ElementCompleteRevocationRefs", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES122ElementCounterSignature {
			t.Errorf("ElementCounterSignature() = %v, want XAdES122ElementCounterSignature", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES122ElementCountryName {
			t.Errorf("ElementCountryName() = %v, want XAdES122ElementCountryName", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES122ElementCRLIdentifier {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES122ElementCRLIdentifier", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES122ElementCRLRef {
			t.Errorf("ElementCRLRef() = %v, want XAdES122ElementCRLRef", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES122ElementCRLRefs {
			t.Errorf("ElementCRLRefs() = %v, want XAdES122ElementCRLRefs", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES122ElementCRLValues {
			t.Errorf("ElementCRLValues() = %v, want XAdES122ElementCRLValues", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES122ElementDataObjectFormat {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES122ElementDataObjectFormat", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES122ElementDescription {
			t.Errorf("ElementDescription() = %v, want XAdES122ElementDescription", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES122ElementDigestAlgAndValue {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES122ElementDigestAlgAndValue", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES122ElementDocumentationReference {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES122ElementDocumentationReference", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES122ElementDocumentationReferences {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES122ElementDocumentationReferences", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES122ElementEncapsulatedCRLValue {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES122ElementEncapsulatedCRLValue", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES122ElementEncapsulatedOCSPValue {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES122ElementEncapsulatedOCSPValue", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES122ElementEncapsulatedPKIData {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES122ElementEncapsulatedPKIData", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES122ElementEncapsulatedTimestamp {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES122ElementEncapsulatedTimestamp", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES122ElementEncapsulatedX509Certificate {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES122ElementEncapsulatedX509Certificate", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES122ElementEncoding {
			t.Errorf("ElementEncoding() = %v, want XAdES122ElementEncoding", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES122ElementExplicitText {
			t.Errorf("ElementExplicitText() = %v, want XAdES122ElementExplicitText", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES122ElementIdentifier {
			t.Errorf("ElementIdentifier() = %v, want XAdES122ElementIdentifier", got)
		}
	})
	t.Run("ElementInclude", func(t *testing.T) {
		if got := e.ElementInclude(); got != XAdES122ElementInclude {
			t.Errorf("ElementInclude() = %v, want XAdES122ElementInclude", got)
		}
	})
	t.Run("ElementIndividualDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES122ElementIndividualDataObjectsTimestamp {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES122ElementIndividualDataObjectsTimestamp", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES122ElementInt {
			t.Errorf("Elementint() = %v, want XAdES122ElementInt", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES122ElementIssueTime {
			t.Errorf("ElementIssueTime() = %v, want XAdES122ElementIssueTime", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES122ElementIssuer {
			t.Errorf("ElementIssuer() = %v, want XAdES122ElementIssuer", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES122ElementIssuerSerial {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES122ElementIssuerSerial", got)
		}
	})
	t.Run("ElementIssuerSerialV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementIssuerSerialV2() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementIssuerSerialV2()
	})
	t.Run("ElementMimeType", func(t *testing.T) {
		if got := e.ElementMimeType(); got != XAdES122ElementMIMEType {
			t.Errorf("ElementMimeType() = %v, want XAdES122ElementMIMEType", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES122ElementNoticeNumbers {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES122ElementNoticeNumbers", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES122ElementNoticeRef {
			t.Errorf("ElementNoticeRef() = %v, want XAdES122ElementNoticeRef", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES122ElementNumber {
			t.Errorf("ElementNumber() = %v, want XAdES122ElementNumber", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES122ElementObjectIdentifier {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES122ElementObjectIdentifier", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES122ElementObjectReference {
			t.Errorf("ElementObjectReference() = %v, want XAdES122ElementObjectReference", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES122ElementOCSPIdentifier {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES122ElementOCSPIdentifier", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES122ElementOCSPRef {
			t.Errorf("ElementOCSPRef() = %v, want XAdES122ElementOCSPRef", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES122ElementOCSPRefs {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES122ElementOCSPRefs", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES122ElementOCSPValues {
			t.Errorf("ElementOCSPValues() = %v, want XAdES122ElementOCSPValues", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES122ElementOrganization {
			t.Errorf("ElementOrganization() = %v, want XAdES122ElementOrganization", got)
		}
	})
	t.Run("ElementOtherAttributeCertificate", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementOtherAttributeCertificate() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementOtherAttributeCertificate()
	})
	t.Run("ElementOtherCertificate", func(t *testing.T) {
		if got := e.ElementOtherCertificate(); got != XAdES122ElementOtherCertificate {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES122ElementOtherCertificate", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES122ElementOtherRef {
			t.Errorf("ElementOtherRef() = %v, want XAdES122ElementOtherRef", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES122ElementOtherRefs {
			t.Errorf("ElementOtherRefs() = %v, want XAdES122ElementOtherRefs", got)
		}
	})
	t.Run("ElementOtherTimeStamp", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementOtherTimeStamp() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementOtherTimeStamp()
	})
	t.Run("ElementOtherValue", func(t *testing.T) {
		if got := e.ElementOtherValue(); got != XAdES122ElementOtherValue {
			t.Errorf("ElementOtherValue() = %v, want XAdES122ElementOtherValue", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES122ElementOtherValues {
			t.Errorf("ElementOtherValues() = %v, want XAdES122ElementOtherValues", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES122ElementPostalCode {
			t.Errorf("ElementPostalCode() = %v, want XAdES122ElementPostalCode", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES122ElementProducedAt {
			t.Errorf("ElementProducedAt() = %v, want XAdES122ElementProducedAt", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES122ElementQualifyingProperties {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES122ElementQualifyingProperties", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES122ElementQualifyingPropertiesReference {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES122ElementQualifyingPropertiesReference", got)
		}
	})
	t.Run("ElementReferenceInfo", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementReferenceInfo() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementReferenceInfo()
	})
	t.Run("ElementRefsOnlyTimeStamp", func(t *testing.T) {
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES122ElementRefsOnlyTimestamp {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES122ElementRefsOnlyTimestamp", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES122ElementResponderID {
			t.Errorf("ElementResponderID() = %v, want XAdES122ElementResponderID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES122ElementRevocationValues {
			t.Errorf("ElementRevocationValues() = %v, want XAdES122ElementRevocationValues", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES122ElementSigAndRefsTimestamp {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES122ElementSigAndRefsTimestamp", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES122ElementSigPolicyHash {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES122ElementSigPolicyHash", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES122ElementSigPolicyID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES122ElementSigPolicyID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES122ElementSigPolicyQualifier {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES122ElementSigPolicyQualifier", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES122ElementSigPolicyQualifiers {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES122ElementSigPolicyQualifiers", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES122ElementSignaturePolicyID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES122ElementSignaturePolicyID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES122ElementSignaturePolicyIdentifier {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES122ElementSignaturePolicyIdentifier", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES122ElementSignaturePolicyImplied {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES122ElementSignaturePolicyImplied", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES122ElementSignatureProductionPlace {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES122ElementSignatureProductionPlace", got)
		}
	})
	t.Run("ElementSignatureProductionPlaceV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignatureProductionPlaceV2() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementSignatureProductionPlaceV2()
	})
	t.Run("ElementSignatureTimeStamp", func(t *testing.T) {
		if got := e.ElementSignatureTimeStamp(); got != XAdES122ElementSignatureTimestamp {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES122ElementSignatureTimestamp", got)
		}
	})
	t.Run("ElementSignedAssertion", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignedAssertion() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementSignedAssertion()
	})
	t.Run("ElementSignedAssertions", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignedAssertions() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementSignedAssertions()
	})
	t.Run("ElementSignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementSignedDataObjectProperties(); got != XAdES122ElementSignedDataObjectProperties {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES122ElementSignedDataObjectProperties", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES122ElementSignedProperties {
			t.Errorf("ElementSignedProperties() = %v, want XAdES122ElementSignedProperties", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES122ElementSignedSignatureProperties {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES122ElementSignedSignatureProperties", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES122ElementSignerRole {
			t.Errorf("ElementSignerRole() = %v, want XAdES122ElementSignerRole", got)
		}
	})
	t.Run("ElementSignerRoleV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignerRoleV2() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementSignerRoleV2()
	})
	t.Run("ElementSigningCertificate", func(t *testing.T) {
		if got := e.ElementSigningCertificate(); got != XAdES122ElementSigningCertificate {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES122ElementSigningCertificate", got)
		}
	})
	t.Run("ElementSigningCertificateV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSigningCertificateV2() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementSigningCertificateV2()
	})
	t.Run("ElementSigningTime", func(t *testing.T) {
		if got := e.ElementSigningTime(); got != XAdES122ElementSigningTime {
			t.Errorf("ElementSigningTime() = %v, want XAdES122ElementSigningTime", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES122ElementSPURI {
			t.Errorf("ElementSPURI() = %v, want XAdES122ElementSPURI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES122ElementSPUserNotice {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES122ElementSPUserNotice", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES122ElementStateOrProvince {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES122ElementStateOrProvince", got)
		}
	})
	t.Run("ElementStreetAddress", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementStreetAddress() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementStreetAddress()
	})
	t.Run("ElementUnsignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES122ElementUnsignedDataObjectProperties {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES122ElementUnsignedDataObjectProperties", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES122ElementUnsignedDataObjectProperty {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES122ElementUnsignedDataObjectProperty", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES122ElementUnsignedProperties {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES122ElementUnsignedProperties", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES122ElementUnsignedSignatureProperties {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES122ElementUnsignedSignatureProperties", got)
		}
	})
	t.Run("ElementX509AttributeCertificate", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementX509AttributeCertificate() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementX509AttributeCertificate()
	})
	t.Run("ElementXAdESTimeStamp", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementXAdESTimeStamp() should panic for XAdES122Element")
			}
		}()
		_ = e.ElementXAdESTimeStamp()
	})
	t.Run("ElementXMLTimeStamp", func(t *testing.T) {
		if got := e.ElementXMLTimeStamp(); got != XAdES122ElementXMLTimestamp {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES122ElementXMLTimestamp", got)
		}
	})
}
