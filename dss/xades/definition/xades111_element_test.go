// KAT test for xades111_element.go: every XAdES111Element_* tag name below is transcribed
// verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.xades111.XAdES111Element, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestXAdES111Element_TagNameKAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES111ElementAllDataObjectsTimestamp.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES111ElementAllSignedDataObjects.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES111ElementAny.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES111ElementArchiveTimestamp.TagName(), "ArchiveTimeStamp"},
		{"CERT", XAdES111ElementCert.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES111ElementCertDigest.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES111ElementCertRefs.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES111ElementCertificateValues.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES111ElementCertifiedRole.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES111ElementCertifiedRoles.TagName(), "CertifiedRoles"},
		{"CITY", XAdES111ElementCity.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES111ElementClaimedRole.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES111ElementClaimedRoles.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES111ElementCommitmentTypeID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES111ElementCommitmentTypeIndication.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES111ElementCommitmentTypeQualifier.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES111ElementCommitmentTypeQualifiers.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES111ElementCompleteCertificateRefs.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES111ElementCompleteRevocationRefs.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES111ElementCounterSignature.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES111ElementCountryName.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES111ElementCRLIdentifier.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES111ElementCRLRef.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES111ElementCRLRefs.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES111ElementCRLValues.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES111ElementDataObjectFormat.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES111ElementDescription.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES111ElementDigestAlgAndValue.TagName(), "DigestAlgAndValue"},
		{"DIGEST_METHOD", XAdES111ElementDigestMethod.TagName(), "DigestMethod"},
		{"DIGEST_VALUE", XAdES111ElementDigestValue.TagName(), "DigestValue"},
		{"DOCUMENTATION_REFERENCE", XAdES111ElementDocumentationReference.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES111ElementDocumentationReferences.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES111ElementEncapsulatedCRLValue.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES111ElementEncapsulatedOCSPValue.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES111ElementEncapsulatedPKIData.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES111ElementEncapsulatedTimestamp.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES111ElementEncapsulatedX509Certificate.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES111ElementEncoding.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES111ElementExplicitText.TagName(), "ExplicitText"},
		{"HASH_DATA_INFO", XAdES111ElementHashDataInfo.TagName(), "HashDataInfo"},
		{"IDENTIFIER", XAdES111ElementIdentifier.TagName(), "Identifier"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES111ElementIndividualDataObjectsTimestamp.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES111ElementInt.TagName(), "int"},
		{"ISSUE_TIME", XAdES111ElementIssueTime.TagName(), "IssueTime"},
		{"ISSUER", XAdES111ElementIssuer.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES111ElementIssuerSerial.TagName(), "IssuerSerial"},
		{"MIME_TYPE", XAdES111ElementMIMEType.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES111ElementNoticeNumbers.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES111ElementNoticeRef.TagName(), "NoticeRef"},
		{"NUMBER", XAdES111ElementNumber.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES111ElementObjectIdentifier.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES111ElementObjectReference.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES111ElementOCSPIdentifier.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES111ElementOCSPRef.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES111ElementOCSPRefs.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES111ElementOCSPValues.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES111ElementOrganization.TagName(), "Organization"},
		{"OTHER_CERTIFICATE", XAdES111ElementOtherCertificate.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES111ElementOtherRef.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES111ElementOtherRefs.TagName(), "OtherRefs"},
		{"OTHER_VALUE", XAdES111ElementOtherValue.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES111ElementOtherValues.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES111ElementPostalCode.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES111ElementProducedAt.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES111ElementQualifyingProperties.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES111ElementQualifyingPropertiesReference.TagName(), "QualifyingPropertiesReference"},
		{"REFS_ONLY_TIMESTAMP", XAdES111ElementRefsOnlyTimestamp.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES111ElementResponderID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES111ElementRevocationValues.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES111ElementSigAndRefsTimestamp.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES111ElementSigPolicyHash.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES111ElementSigPolicyID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES111ElementSigPolicyQualifier.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES111ElementSigPolicyQualifiers.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES111ElementSignaturePolicyID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES111ElementSignaturePolicyIdentifier.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES111ElementSignaturePolicyImplied.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES111ElementSignatureProductionPlace.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_TIMESTAMP", XAdES111ElementSignatureTimestamp.TagName(), "SignatureTimeStamp"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES111ElementSignedDataObjectProperties.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES111ElementSignedProperties.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES111ElementSignedSignatureProperties.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES111ElementSignerRole.TagName(), "SignerRole"},
		{"SIGNING_CERTIFICATE", XAdES111ElementSigningCertificate.TagName(), "SigningCertificate"},
		{"SIGNING_TIME", XAdES111ElementSigningTime.TagName(), "SigningTime"},
		{"SP_URI", XAdES111ElementSPURI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES111ElementSPUserNotice.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES111ElementStateOrProvince.TagName(), "StateOrProvince"},
		{"TIMESTAMP", XAdES111ElementTimestamp.TagName(), "TimeStamp"},
		{"TRANSFORMS", XAdES111ElementTransforms.TagName(), "Transforms"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES111ElementUnsignedDataObjectProperties.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES111ElementUnsignedDataObjectProperty.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES111ElementUnsignedProperties.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES111ElementUnsignedSignatureProperties.TagName(), "UnsignedSignatureProperties"},
		{"XML_TIMESTAMP", XAdES111ElementXMLTimestamp.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES111Element_Namespace(t *testing.T) {
	if got := XAdES111ElementAllDataObjectsTimestamp.URI(); got != XAdESNamespaceXAdES111.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespaceXAdES111.Uri())
	}
	if !XAdES111ElementAllDataObjectsTimestamp.IsSameTagName(XAdES111ElementAllDataObjectsTimestamp.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES111Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES111Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES111Element_ElementGetters(t *testing.T) {
	var e XAdES111Element = XAdES111ElementAllDataObjectsTimestamp
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES111ElementAllDataObjectsTimestamp {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES111ElementAllDataObjectsTimestamp", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES111ElementAllSignedDataObjects {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES111ElementAllSignedDataObjects", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES111ElementAny {
			t.Errorf("ElementAny() = %v, want XAdES111ElementAny", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES111ElementArchiveTimestamp {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES111ElementArchiveTimestamp", got)
		}
	})
	t.Run("ElementAttrAuthoritiesCertValues", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttrAuthoritiesCertValues() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementAttrAuthoritiesCertValues()
	})
	t.Run("ElementAttributeCertificateRefs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttributeCertificateRefs() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementAttributeCertificateRefs()
	})
	t.Run("ElementAttributeRevocationRefs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttributeRevocationRefs() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementAttributeRevocationRefs()
	})
	t.Run("ElementAttributeRevocationValues", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementAttributeRevocationValues() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementAttributeRevocationValues()
	})
	t.Run("ElementByKey", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementByKey() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementByKey()
	})
	t.Run("ElementByName", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementByName() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementByName()
	})
	t.Run("ElementCert", func(t *testing.T) {
		if got := e.ElementCert(); got != XAdES111ElementCert {
			t.Errorf("ElementCert() = %v, want XAdES111ElementCert", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES111ElementCertDigest {
			t.Errorf("ElementCertDigest() = %v, want XAdES111ElementCertDigest", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES111ElementCertRefs {
			t.Errorf("ElementCertRefs() = %v, want XAdES111ElementCertRefs", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES111ElementCertificateValues {
			t.Errorf("ElementCertificateValues() = %v, want XAdES111ElementCertificateValues", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES111ElementCertifiedRole {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES111ElementCertifiedRole", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES111ElementCertifiedRoles {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES111ElementCertifiedRoles", got)
		}
	})
	t.Run("ElementCertifiedRolesV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementCertifiedRolesV2() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementCertifiedRolesV2()
	})
	t.Run("ElementCity", func(t *testing.T) {
		if got := e.ElementCity(); got != XAdES111ElementCity {
			t.Errorf("ElementCity() = %v, want XAdES111ElementCity", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES111ElementClaimedRole {
			t.Errorf("ElementClaimedRole() = %v, want XAdES111ElementClaimedRole", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES111ElementClaimedRoles {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES111ElementClaimedRoles", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES111ElementCommitmentTypeID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES111ElementCommitmentTypeID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES111ElementCommitmentTypeIndication {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES111ElementCommitmentTypeIndication", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES111ElementCommitmentTypeQualifier {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES111ElementCommitmentTypeQualifier", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES111ElementCommitmentTypeQualifiers {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES111ElementCommitmentTypeQualifiers", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES111ElementCompleteCertificateRefs {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES111ElementCompleteCertificateRefs", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES111ElementCompleteRevocationRefs {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES111ElementCompleteRevocationRefs", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES111ElementCounterSignature {
			t.Errorf("ElementCounterSignature() = %v, want XAdES111ElementCounterSignature", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES111ElementCountryName {
			t.Errorf("ElementCountryName() = %v, want XAdES111ElementCountryName", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES111ElementCRLIdentifier {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES111ElementCRLIdentifier", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES111ElementCRLRef {
			t.Errorf("ElementCRLRef() = %v, want XAdES111ElementCRLRef", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES111ElementCRLRefs {
			t.Errorf("ElementCRLRefs() = %v, want XAdES111ElementCRLRefs", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES111ElementCRLValues {
			t.Errorf("ElementCRLValues() = %v, want XAdES111ElementCRLValues", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES111ElementDataObjectFormat {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES111ElementDataObjectFormat", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES111ElementDescription {
			t.Errorf("ElementDescription() = %v, want XAdES111ElementDescription", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES111ElementDigestAlgAndValue {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES111ElementDigestAlgAndValue", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES111ElementDocumentationReference {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES111ElementDocumentationReference", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES111ElementDocumentationReferences {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES111ElementDocumentationReferences", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES111ElementEncapsulatedCRLValue {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES111ElementEncapsulatedCRLValue", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES111ElementEncapsulatedOCSPValue {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES111ElementEncapsulatedOCSPValue", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES111ElementEncapsulatedPKIData {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES111ElementEncapsulatedPKIData", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES111ElementEncapsulatedTimestamp {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES111ElementEncapsulatedTimestamp", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES111ElementEncapsulatedX509Certificate {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES111ElementEncapsulatedX509Certificate", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES111ElementEncoding {
			t.Errorf("ElementEncoding() = %v, want XAdES111ElementEncoding", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES111ElementExplicitText {
			t.Errorf("ElementExplicitText() = %v, want XAdES111ElementExplicitText", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES111ElementIdentifier {
			t.Errorf("ElementIdentifier() = %v, want XAdES111ElementIdentifier", got)
		}
	})
	t.Run("ElementInclude", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementInclude() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementInclude()
	})
	t.Run("ElementIndividualDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES111ElementIndividualDataObjectsTimestamp {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES111ElementIndividualDataObjectsTimestamp", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES111ElementInt {
			t.Errorf("Elementint() = %v, want XAdES111ElementInt", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES111ElementIssueTime {
			t.Errorf("ElementIssueTime() = %v, want XAdES111ElementIssueTime", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES111ElementIssuer {
			t.Errorf("ElementIssuer() = %v, want XAdES111ElementIssuer", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES111ElementIssuerSerial {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES111ElementIssuerSerial", got)
		}
	})
	t.Run("ElementIssuerSerialV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementIssuerSerialV2() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementIssuerSerialV2()
	})
	t.Run("ElementMimeType", func(t *testing.T) {
		if got := e.ElementMimeType(); got != XAdES111ElementMIMEType {
			t.Errorf("ElementMimeType() = %v, want XAdES111ElementMIMEType", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES111ElementNoticeNumbers {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES111ElementNoticeNumbers", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES111ElementNoticeRef {
			t.Errorf("ElementNoticeRef() = %v, want XAdES111ElementNoticeRef", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES111ElementNumber {
			t.Errorf("ElementNumber() = %v, want XAdES111ElementNumber", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES111ElementObjectIdentifier {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES111ElementObjectIdentifier", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES111ElementObjectReference {
			t.Errorf("ElementObjectReference() = %v, want XAdES111ElementObjectReference", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES111ElementOCSPIdentifier {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES111ElementOCSPIdentifier", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES111ElementOCSPRef {
			t.Errorf("ElementOCSPRef() = %v, want XAdES111ElementOCSPRef", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES111ElementOCSPRefs {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES111ElementOCSPRefs", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES111ElementOCSPValues {
			t.Errorf("ElementOCSPValues() = %v, want XAdES111ElementOCSPValues", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES111ElementOrganization {
			t.Errorf("ElementOrganization() = %v, want XAdES111ElementOrganization", got)
		}
	})
	t.Run("ElementOtherAttributeCertificate", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementOtherAttributeCertificate() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementOtherAttributeCertificate()
	})
	t.Run("ElementOtherCertificate", func(t *testing.T) {
		if got := e.ElementOtherCertificate(); got != XAdES111ElementOtherCertificate {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES111ElementOtherCertificate", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES111ElementOtherRef {
			t.Errorf("ElementOtherRef() = %v, want XAdES111ElementOtherRef", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES111ElementOtherRefs {
			t.Errorf("ElementOtherRefs() = %v, want XAdES111ElementOtherRefs", got)
		}
	})
	t.Run("ElementOtherTimeStamp", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementOtherTimeStamp() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementOtherTimeStamp()
	})
	t.Run("ElementOtherValue", func(t *testing.T) {
		if got := e.ElementOtherValue(); got != XAdES111ElementOtherValue {
			t.Errorf("ElementOtherValue() = %v, want XAdES111ElementOtherValue", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES111ElementOtherValues {
			t.Errorf("ElementOtherValues() = %v, want XAdES111ElementOtherValues", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES111ElementPostalCode {
			t.Errorf("ElementPostalCode() = %v, want XAdES111ElementPostalCode", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES111ElementProducedAt {
			t.Errorf("ElementProducedAt() = %v, want XAdES111ElementProducedAt", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES111ElementQualifyingProperties {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES111ElementQualifyingProperties", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES111ElementQualifyingPropertiesReference {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES111ElementQualifyingPropertiesReference", got)
		}
	})
	t.Run("ElementReferenceInfo", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementReferenceInfo() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementReferenceInfo()
	})
	t.Run("ElementRefsOnlyTimeStamp", func(t *testing.T) {
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES111ElementRefsOnlyTimestamp {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES111ElementRefsOnlyTimestamp", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES111ElementResponderID {
			t.Errorf("ElementResponderID() = %v, want XAdES111ElementResponderID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES111ElementRevocationValues {
			t.Errorf("ElementRevocationValues() = %v, want XAdES111ElementRevocationValues", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES111ElementSigAndRefsTimestamp {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES111ElementSigAndRefsTimestamp", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES111ElementSigPolicyHash {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES111ElementSigPolicyHash", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES111ElementSigPolicyID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES111ElementSigPolicyID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES111ElementSigPolicyQualifier {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES111ElementSigPolicyQualifier", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES111ElementSigPolicyQualifiers {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES111ElementSigPolicyQualifiers", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES111ElementSignaturePolicyID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES111ElementSignaturePolicyID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES111ElementSignaturePolicyIdentifier {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES111ElementSignaturePolicyIdentifier", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES111ElementSignaturePolicyImplied {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES111ElementSignaturePolicyImplied", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES111ElementSignatureProductionPlace {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES111ElementSignatureProductionPlace", got)
		}
	})
	t.Run("ElementSignatureProductionPlaceV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignatureProductionPlaceV2() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementSignatureProductionPlaceV2()
	})
	t.Run("ElementSignatureTimeStamp", func(t *testing.T) {
		if got := e.ElementSignatureTimeStamp(); got != XAdES111ElementSignatureTimestamp {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES111ElementSignatureTimestamp", got)
		}
	})
	t.Run("ElementSignedAssertion", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignedAssertion() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementSignedAssertion()
	})
	t.Run("ElementSignedAssertions", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignedAssertions() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementSignedAssertions()
	})
	t.Run("ElementSignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementSignedDataObjectProperties(); got != XAdES111ElementSignedDataObjectProperties {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES111ElementSignedDataObjectProperties", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES111ElementSignedProperties {
			t.Errorf("ElementSignedProperties() = %v, want XAdES111ElementSignedProperties", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES111ElementSignedSignatureProperties {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES111ElementSignedSignatureProperties", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES111ElementSignerRole {
			t.Errorf("ElementSignerRole() = %v, want XAdES111ElementSignerRole", got)
		}
	})
	t.Run("ElementSignerRoleV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSignerRoleV2() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementSignerRoleV2()
	})
	t.Run("ElementSigningCertificate", func(t *testing.T) {
		if got := e.ElementSigningCertificate(); got != XAdES111ElementSigningCertificate {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES111ElementSigningCertificate", got)
		}
	})
	t.Run("ElementSigningCertificateV2", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementSigningCertificateV2() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementSigningCertificateV2()
	})
	t.Run("ElementSigningTime", func(t *testing.T) {
		if got := e.ElementSigningTime(); got != XAdES111ElementSigningTime {
			t.Errorf("ElementSigningTime() = %v, want XAdES111ElementSigningTime", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES111ElementSPURI {
			t.Errorf("ElementSPURI() = %v, want XAdES111ElementSPURI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES111ElementSPUserNotice {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES111ElementSPUserNotice", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES111ElementStateOrProvince {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES111ElementStateOrProvince", got)
		}
	})
	t.Run("ElementStreetAddress", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementStreetAddress() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementStreetAddress()
	})
	t.Run("ElementUnsignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES111ElementUnsignedDataObjectProperties {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES111ElementUnsignedDataObjectProperties", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES111ElementUnsignedDataObjectProperty {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES111ElementUnsignedDataObjectProperty", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES111ElementUnsignedProperties {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES111ElementUnsignedProperties", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES111ElementUnsignedSignatureProperties {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES111ElementUnsignedSignatureProperties", got)
		}
	})
	t.Run("ElementX509AttributeCertificate", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementX509AttributeCertificate() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementX509AttributeCertificate()
	})
	t.Run("ElementXAdESTimeStamp", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("ElementXAdESTimeStamp() should panic for XAdES111Element")
			}
		}()
		_ = e.ElementXAdESTimeStamp()
	})
	t.Run("ElementXMLTimeStamp", func(t *testing.T) {
		if got := e.ElementXMLTimeStamp(); got != XAdES111ElementXMLTimestamp {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES111ElementXMLTimestamp", got)
		}
	})
}
