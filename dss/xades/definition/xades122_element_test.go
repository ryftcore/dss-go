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
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES122Element_ALL_SIGNED_DATA_OBJECTS.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES122Element_ANY.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES122Element_ARCHIVE_TIMESTAMP.TagName(), "ArchiveTimeStamp"},
		{"ATTRIBUTE_CERTIFICATE_REFS", XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS.TagName(), "AttributeCertificateRefs"},
		{"ATTRIBUTE_REVOCATION_REFS", XAdES122Element_ATTRIBUTE_REVOCATION_REFS.TagName(), "AttributeRevocationRefs"},
		{"CERT", XAdES122Element_CERT.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES122Element_CERT_DIGEST.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES122Element_CERT_REFS.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES122Element_CERTIFICATE_VALUES.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES122Element_CERTIFIED_ROLE.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES122Element_CERTIFIED_ROLES.TagName(), "CertifiedRoles"},
		{"CITY", XAdES122Element_CITY.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES122Element_CLAIMED_ROLE.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES122Element_CLAIMED_ROLES.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES122Element_COMMITMENT_TYPE_ID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES122Element_COMMITMENT_TYPE_INDICATION.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES122Element_COMMITMENT_TYPE_QUALIFIER.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES122Element_COMMITMENT_TYPE_QUALIFIERS.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES122Element_COMPLETE_CERTIFICATE_REFS.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES122Element_COMPLETE_REVOCATION_REFS.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES122Element_COUNTER_SIGNATURE.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES122Element_COUNTRY_NAME.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES122Element_CRL_IDENTIFIER.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES122Element_CRL_REF.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES122Element_CRL_REFS.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES122Element_CRL_VALUES.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES122Element_DATA_OBJECT_FORMAT.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES122Element_DESCRIPTION.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES122Element_DIGEST_ALG_AND_VALUE.TagName(), "DigestAlgAndValue"},
		{"DOCUMENTATION_REFERENCE", XAdES122Element_DOCUMENTATION_REFERENCE.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES122Element_DOCUMENTATION_REFERENCES.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES122Element_ENCAPSULATED_CRL_VALUE.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES122Element_ENCAPSULATED_OCSP_VALUE.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES122Element_ENCAPSULATED_PKI_DATA.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES122Element_ENCAPSULATED_TIMESTAMP.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES122Element_ENCAPSULATED_X509_CERTIFICATE.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES122Element_ENCODING.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES122Element_EXPLICIT_TEXT.TagName(), "ExplicitText"},
		{"IDENTIFIER", XAdES122Element_IDENTIFIER.TagName(), "Identifier"},
		{"INCLUDE", XAdES122Element_INCLUDE.TagName(), "Include"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES122Element_INT.TagName(), "int"},
		{"ISSUE_TIME", XAdES122Element_ISSUE_TIME.TagName(), "IssueTime"},
		{"ISSUER", XAdES122Element_ISSUER.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES122Element_ISSUER_SERIAL.TagName(), "IssuerSerial"},
		{"MIME_TYPE", XAdES122Element_MIME_TYPE.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES122Element_NOTICE_NUMBERS.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES122Element_NOTICE_REF.TagName(), "NoticeRef"},
		{"NUMBER", XAdES122Element_NUMBER.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES122Element_OBJECT_IDENTIFIER.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES122Element_OBJECT_REFERENCE.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES122Element_OCSP_IDENTIFIER.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES122Element_OCSP_REF.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES122Element_OCSP_REFS.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES122Element_OCSP_VALUES.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES122Element_ORGANIZATION.TagName(), "Organization"},
		{"OTHER_CERTIFICATE", XAdES122Element_OTHER_CERTIFICATE.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES122Element_OTHER_REF.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES122Element_OTHER_REFS.TagName(), "OtherRefs"},
		{"OTHER_VALUE", XAdES122Element_OTHER_VALUE.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES122Element_OTHER_VALUES.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES122Element_POSTAL_CODE.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES122Element_PRODUCED_AT.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES122Element_QUALIFYING_PROPERTIES.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE.TagName(), "QualifyingPropertiesReference"},
		{"REFS_ONLY_TIMESTAMP", XAdES122Element_REFS_ONLY_TIMESTAMP.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES122Element_RESPONDER_ID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES122Element_REVOCATION_VALUES.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES122Element_SIG_AND_REFS_TIMESTAMP.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES122Element_SIG_POLICY_HASH.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES122Element_SIG_POLICY_ID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES122Element_SIG_POLICY_QUALIFIER.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES122Element_SIG_POLICY_QUALIFIERS.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES122Element_SIGNATURE_POLICY_ID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES122Element_SIGNATURE_POLICY_IDENTIFIER.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES122Element_SIGNATURE_POLICY_IMPLIED.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES122Element_SIGNATURE_PRODUCTION_PLACE.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_TIMESTAMP", XAdES122Element_SIGNATURE_TIMESTAMP.TagName(), "SignatureTimeStamp"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES122Element_SIGNED_PROPERTIES.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES122Element_SIGNED_SIGNATURE_PROPERTIES.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES122Element_SIGNER_ROLE.TagName(), "SignerRole"},
		{"SIGNING_CERTIFICATE", XAdES122Element_SIGNING_CERTIFICATE.TagName(), "SigningCertificate"},
		{"SIGNING_TIME", XAdES122Element_SIGNING_TIME.TagName(), "SigningTime"},
		{"SP_URI", XAdES122Element_SP_URI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES122Element_SP_USER_NOTICE.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES122Element_STATE_OR_PROVINCE.TagName(), "StateOrProvince"},
		{"TIMESTAMP", XAdES122Element_TIMESTAMP.TagName(), "TimeStamp"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES122Element_UNSIGNED_PROPERTIES.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES.TagName(), "UnsignedSignatureProperties"},
		{"XML_TIMESTAMP", XAdES122Element_XML_TIMESTAMP.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES122Element_Namespace(t *testing.T) {
	if got := XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP.URI(); got != XAdESNamespace_XADES_122.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespace_XADES_122.Uri())
	}
	if !XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP.IsSameTagName(XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES122Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES122Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES122Element_ElementGetters(t *testing.T) {
	var e XAdES122Element = XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES122Element_ALL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES122Element_ALL_SIGNED_DATA_OBJECTS {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES122Element_ALL_SIGNED_DATA_OBJECTS", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES122Element_ANY {
			t.Errorf("ElementAny() = %v, want XAdES122Element_ANY", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES122Element_ARCHIVE_TIMESTAMP {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES122Element_ARCHIVE_TIMESTAMP", got)
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
		if got := e.ElementAttributeCertificateRefs(); got != XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS {
			t.Errorf("ElementAttributeCertificateRefs() = %v, want XAdES122Element_ATTRIBUTE_CERTIFICATE_REFS", got)
		}
	})
	t.Run("ElementAttributeRevocationRefs", func(t *testing.T) {
		if got := e.ElementAttributeRevocationRefs(); got != XAdES122Element_ATTRIBUTE_REVOCATION_REFS {
			t.Errorf("ElementAttributeRevocationRefs() = %v, want XAdES122Element_ATTRIBUTE_REVOCATION_REFS", got)
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
		if got := e.ElementCert(); got != XAdES122Element_CERT {
			t.Errorf("ElementCert() = %v, want XAdES122Element_CERT", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES122Element_CERT_DIGEST {
			t.Errorf("ElementCertDigest() = %v, want XAdES122Element_CERT_DIGEST", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES122Element_CERT_REFS {
			t.Errorf("ElementCertRefs() = %v, want XAdES122Element_CERT_REFS", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES122Element_CERTIFICATE_VALUES {
			t.Errorf("ElementCertificateValues() = %v, want XAdES122Element_CERTIFICATE_VALUES", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES122Element_CERTIFIED_ROLE {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES122Element_CERTIFIED_ROLE", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES122Element_CERTIFIED_ROLES {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES122Element_CERTIFIED_ROLES", got)
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
		if got := e.ElementCity(); got != XAdES122Element_CITY {
			t.Errorf("ElementCity() = %v, want XAdES122Element_CITY", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES122Element_CLAIMED_ROLE {
			t.Errorf("ElementClaimedRole() = %v, want XAdES122Element_CLAIMED_ROLE", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES122Element_CLAIMED_ROLES {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES122Element_CLAIMED_ROLES", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES122Element_COMMITMENT_TYPE_ID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES122Element_COMMITMENT_TYPE_ID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES122Element_COMMITMENT_TYPE_INDICATION {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES122Element_COMMITMENT_TYPE_INDICATION", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES122Element_COMMITMENT_TYPE_QUALIFIER {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES122Element_COMMITMENT_TYPE_QUALIFIER", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES122Element_COMMITMENT_TYPE_QUALIFIERS {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES122Element_COMMITMENT_TYPE_QUALIFIERS", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES122Element_COMPLETE_CERTIFICATE_REFS {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES122Element_COMPLETE_CERTIFICATE_REFS", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES122Element_COMPLETE_REVOCATION_REFS {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES122Element_COMPLETE_REVOCATION_REFS", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES122Element_COUNTER_SIGNATURE {
			t.Errorf("ElementCounterSignature() = %v, want XAdES122Element_COUNTER_SIGNATURE", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES122Element_COUNTRY_NAME {
			t.Errorf("ElementCountryName() = %v, want XAdES122Element_COUNTRY_NAME", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES122Element_CRL_IDENTIFIER {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES122Element_CRL_IDENTIFIER", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES122Element_CRL_REF {
			t.Errorf("ElementCRLRef() = %v, want XAdES122Element_CRL_REF", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES122Element_CRL_REFS {
			t.Errorf("ElementCRLRefs() = %v, want XAdES122Element_CRL_REFS", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES122Element_CRL_VALUES {
			t.Errorf("ElementCRLValues() = %v, want XAdES122Element_CRL_VALUES", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES122Element_DATA_OBJECT_FORMAT {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES122Element_DATA_OBJECT_FORMAT", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES122Element_DESCRIPTION {
			t.Errorf("ElementDescription() = %v, want XAdES122Element_DESCRIPTION", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES122Element_DIGEST_ALG_AND_VALUE {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES122Element_DIGEST_ALG_AND_VALUE", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES122Element_DOCUMENTATION_REFERENCE {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES122Element_DOCUMENTATION_REFERENCE", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES122Element_DOCUMENTATION_REFERENCES {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES122Element_DOCUMENTATION_REFERENCES", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES122Element_ENCAPSULATED_CRL_VALUE {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES122Element_ENCAPSULATED_CRL_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES122Element_ENCAPSULATED_OCSP_VALUE {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES122Element_ENCAPSULATED_OCSP_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES122Element_ENCAPSULATED_PKI_DATA {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES122Element_ENCAPSULATED_PKI_DATA", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES122Element_ENCAPSULATED_TIMESTAMP {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES122Element_ENCAPSULATED_TIMESTAMP", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES122Element_ENCAPSULATED_X509_CERTIFICATE {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES122Element_ENCAPSULATED_X509_CERTIFICATE", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES122Element_ENCODING {
			t.Errorf("ElementEncoding() = %v, want XAdES122Element_ENCODING", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES122Element_EXPLICIT_TEXT {
			t.Errorf("ElementExplicitText() = %v, want XAdES122Element_EXPLICIT_TEXT", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES122Element_IDENTIFIER {
			t.Errorf("ElementIdentifier() = %v, want XAdES122Element_IDENTIFIER", got)
		}
	})
	t.Run("ElementInclude", func(t *testing.T) {
		if got := e.ElementInclude(); got != XAdES122Element_INCLUDE {
			t.Errorf("ElementInclude() = %v, want XAdES122Element_INCLUDE", got)
		}
	})
	t.Run("ElementIndividualDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES122Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES122Element_INT {
			t.Errorf("Elementint() = %v, want XAdES122Element_INT", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES122Element_ISSUE_TIME {
			t.Errorf("ElementIssueTime() = %v, want XAdES122Element_ISSUE_TIME", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES122Element_ISSUER {
			t.Errorf("ElementIssuer() = %v, want XAdES122Element_ISSUER", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES122Element_ISSUER_SERIAL {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES122Element_ISSUER_SERIAL", got)
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
		if got := e.ElementMimeType(); got != XAdES122Element_MIME_TYPE {
			t.Errorf("ElementMimeType() = %v, want XAdES122Element_MIME_TYPE", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES122Element_NOTICE_NUMBERS {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES122Element_NOTICE_NUMBERS", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES122Element_NOTICE_REF {
			t.Errorf("ElementNoticeRef() = %v, want XAdES122Element_NOTICE_REF", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES122Element_NUMBER {
			t.Errorf("ElementNumber() = %v, want XAdES122Element_NUMBER", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES122Element_OBJECT_IDENTIFIER {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES122Element_OBJECT_IDENTIFIER", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES122Element_OBJECT_REFERENCE {
			t.Errorf("ElementObjectReference() = %v, want XAdES122Element_OBJECT_REFERENCE", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES122Element_OCSP_IDENTIFIER {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES122Element_OCSP_IDENTIFIER", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES122Element_OCSP_REF {
			t.Errorf("ElementOCSPRef() = %v, want XAdES122Element_OCSP_REF", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES122Element_OCSP_REFS {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES122Element_OCSP_REFS", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES122Element_OCSP_VALUES {
			t.Errorf("ElementOCSPValues() = %v, want XAdES122Element_OCSP_VALUES", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES122Element_ORGANIZATION {
			t.Errorf("ElementOrganization() = %v, want XAdES122Element_ORGANIZATION", got)
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
		if got := e.ElementOtherCertificate(); got != XAdES122Element_OTHER_CERTIFICATE {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES122Element_OTHER_CERTIFICATE", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES122Element_OTHER_REF {
			t.Errorf("ElementOtherRef() = %v, want XAdES122Element_OTHER_REF", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES122Element_OTHER_REFS {
			t.Errorf("ElementOtherRefs() = %v, want XAdES122Element_OTHER_REFS", got)
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
		if got := e.ElementOtherValue(); got != XAdES122Element_OTHER_VALUE {
			t.Errorf("ElementOtherValue() = %v, want XAdES122Element_OTHER_VALUE", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES122Element_OTHER_VALUES {
			t.Errorf("ElementOtherValues() = %v, want XAdES122Element_OTHER_VALUES", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES122Element_POSTAL_CODE {
			t.Errorf("ElementPostalCode() = %v, want XAdES122Element_POSTAL_CODE", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES122Element_PRODUCED_AT {
			t.Errorf("ElementProducedAt() = %v, want XAdES122Element_PRODUCED_AT", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES122Element_QUALIFYING_PROPERTIES {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES122Element_QUALIFYING_PROPERTIES", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES122Element_QUALIFYING_PROPERTIES_REFERENCE", got)
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
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES122Element_REFS_ONLY_TIMESTAMP {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES122Element_REFS_ONLY_TIMESTAMP", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES122Element_RESPONDER_ID {
			t.Errorf("ElementResponderID() = %v, want XAdES122Element_RESPONDER_ID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES122Element_REVOCATION_VALUES {
			t.Errorf("ElementRevocationValues() = %v, want XAdES122Element_REVOCATION_VALUES", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES122Element_SIG_AND_REFS_TIMESTAMP {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES122Element_SIG_AND_REFS_TIMESTAMP", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES122Element_SIG_POLICY_HASH {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES122Element_SIG_POLICY_HASH", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES122Element_SIG_POLICY_ID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES122Element_SIG_POLICY_ID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES122Element_SIG_POLICY_QUALIFIER {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES122Element_SIG_POLICY_QUALIFIER", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES122Element_SIG_POLICY_QUALIFIERS {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES122Element_SIG_POLICY_QUALIFIERS", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES122Element_SIGNATURE_POLICY_ID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES122Element_SIGNATURE_POLICY_ID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES122Element_SIGNATURE_POLICY_IDENTIFIER {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES122Element_SIGNATURE_POLICY_IDENTIFIER", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES122Element_SIGNATURE_POLICY_IMPLIED {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES122Element_SIGNATURE_POLICY_IMPLIED", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES122Element_SIGNATURE_PRODUCTION_PLACE {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES122Element_SIGNATURE_PRODUCTION_PLACE", got)
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
		if got := e.ElementSignatureTimeStamp(); got != XAdES122Element_SIGNATURE_TIMESTAMP {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES122Element_SIGNATURE_TIMESTAMP", got)
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
		if got := e.ElementSignedDataObjectProperties(); got != XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES122Element_SIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES122Element_SIGNED_PROPERTIES {
			t.Errorf("ElementSignedProperties() = %v, want XAdES122Element_SIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES122Element_SIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES122Element_SIGNED_SIGNATURE_PROPERTIES", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES122Element_SIGNER_ROLE {
			t.Errorf("ElementSignerRole() = %v, want XAdES122Element_SIGNER_ROLE", got)
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
		if got := e.ElementSigningCertificate(); got != XAdES122Element_SIGNING_CERTIFICATE {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES122Element_SIGNING_CERTIFICATE", got)
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
		if got := e.ElementSigningTime(); got != XAdES122Element_SIGNING_TIME {
			t.Errorf("ElementSigningTime() = %v, want XAdES122Element_SIGNING_TIME", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES122Element_SP_URI {
			t.Errorf("ElementSPURI() = %v, want XAdES122Element_SP_URI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES122Element_SP_USER_NOTICE {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES122Element_SP_USER_NOTICE", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES122Element_STATE_OR_PROVINCE {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES122Element_STATE_OR_PROVINCE", got)
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
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES122Element_UNSIGNED_DATA_OBJECT_PROPERTY", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES122Element_UNSIGNED_PROPERTIES {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES122Element_UNSIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES122Element_UNSIGNED_SIGNATURE_PROPERTIES", got)
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
		if got := e.ElementXMLTimeStamp(); got != XAdES122Element_XML_TIMESTAMP {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES122Element_XML_TIMESTAMP", got)
		}
	})
}
