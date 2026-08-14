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
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES132Element_ALL_SIGNED_DATA_OBJECTS.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES132Element_ANY.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES132Element_ARCHIVE_TIMESTAMP.TagName(), "ArchiveTimeStamp"},
		{"ATTR_AUTHORITIES_CERT_VALUES", XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES.TagName(), "AttrAuthoritiesCertValues"},
		{"ATTRIBUTE_CERTIFICATE_REFS", XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS.TagName(), "AttributeCertificateRefs"},
		{"ATTRIBUTE_REVOCATION_REFS", XAdES132Element_ATTRIBUTE_REVOCATION_REFS.TagName(), "AttributeRevocationRefs"},
		{"ATTRIBUTE_REVOCATION_VALUES", XAdES132Element_ATTRIBUTE_REVOCATION_VALUES.TagName(), "AttributeRevocationValues"},
		{"BY_KEY", XAdES132Element_BY_KEY.TagName(), "ByKey"},
		{"BY_NAME", XAdES132Element_BY_NAME.TagName(), "ByName"},
		{"CERT", XAdES132Element_CERT.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES132Element_CERT_DIGEST.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES132Element_CERT_REFS.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES132Element_CERTIFICATE_VALUES.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES132Element_CERTIFIED_ROLE.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES132Element_CERTIFIED_ROLES.TagName(), "CertifiedRoles"},
		{"CERTIFIED_ROLES_V2", XAdES132Element_CERTIFIED_ROLES_V2.TagName(), "CertifiedRolesV2"},
		{"CITY", XAdES132Element_CITY.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES132Element_CLAIMED_ROLE.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES132Element_CLAIMED_ROLES.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES132Element_COMMITMENT_TYPE_ID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES132Element_COMMITMENT_TYPE_INDICATION.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES132Element_COMMITMENT_TYPE_QUALIFIER.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES132Element_COMMITMENT_TYPE_QUALIFIERS.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES132Element_COMPLETE_CERTIFICATE_REFS.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES132Element_COMPLETE_REVOCATION_REFS.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES132Element_COUNTER_SIGNATURE.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES132Element_COUNTRY_NAME.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES132Element_CRL_IDENTIFIER.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES132Element_CRL_REF.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES132Element_CRL_REFS.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES132Element_CRL_VALUES.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES132Element_DATA_OBJECT_FORMAT.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES132Element_DESCRIPTION.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES132Element_DIGEST_ALG_AND_VALUE.TagName(), "DigestAlgAndValue"},
		{"DOCUMENTATION_REFERENCE", XAdES132Element_DOCUMENTATION_REFERENCE.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES132Element_DOCUMENTATION_REFERENCES.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES132Element_ENCAPSULATED_CRL_VALUE.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES132Element_ENCAPSULATED_OCSP_VALUE.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES132Element_ENCAPSULATED_PKI_DATA.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES132Element_ENCAPSULATED_TIMESTAMP.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES132Element_ENCAPSULATED_X509_CERTIFICATE.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES132Element_ENCODING.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES132Element_EXPLICIT_TEXT.TagName(), "ExplicitText"},
		{"IDENTIFIER", XAdES132Element_IDENTIFIER.TagName(), "Identifier"},
		{"INCLUDE", XAdES132Element_INCLUDE.TagName(), "Include"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES132Element_INT.TagName(), "int"},
		{"ISSUE_TIME", XAdES132Element_ISSUE_TIME.TagName(), "IssueTime"},
		{"ISSUER", XAdES132Element_ISSUER.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES132Element_ISSUER_SERIAL.TagName(), "IssuerSerial"},
		{"ISSUER_SERIAL_V2", XAdES132Element_ISSUER_SERIAL_V2.TagName(), "IssuerSerialV2"},
		{"MIME_TYPE", XAdES132Element_MIME_TYPE.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES132Element_NOTICE_NUMBERS.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES132Element_NOTICE_REF.TagName(), "NoticeRef"},
		{"NUMBER", XAdES132Element_NUMBER.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES132Element_OBJECT_IDENTIFIER.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES132Element_OBJECT_REFERENCE.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES132Element_OCSP_IDENTIFIER.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES132Element_OCSP_REF.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES132Element_OCSP_REFS.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES132Element_OCSP_VALUES.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES132Element_ORGANIZATION.TagName(), "Organization"},
		{"OTHER_ATTRIBUTE_CERTIFICATE", XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE.TagName(), "OtherAttributeCertificate"},
		{"OTHER_CERTIFICATE", XAdES132Element_OTHER_CERTIFICATE.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES132Element_OTHER_REF.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES132Element_OTHER_REFS.TagName(), "OtherRefs"},
		{"OTHER_TIMESTAMP", XAdES132Element_OTHER_TIMESTAMP.TagName(), "OtherTimeStamp"},
		{"OTHER_VALUE", XAdES132Element_OTHER_VALUE.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES132Element_OTHER_VALUES.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES132Element_POSTAL_CODE.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES132Element_PRODUCED_AT.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES132Element_QUALIFYING_PROPERTIES.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE.TagName(), "QualifyingPropertiesReference"},
		{"REFERENCE_INFO", XAdES132Element_REFERENCE_INFO.TagName(), "ReferenceInfo"},
		{"REFS_ONLY_TIMESTAMP", XAdES132Element_REFS_ONLY_TIMESTAMP.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES132Element_RESPONDER_ID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES132Element_REVOCATION_VALUES.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES132Element_SIG_AND_REFS_TIMESTAMP.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES132Element_SIG_POLICY_HASH.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES132Element_SIG_POLICY_ID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES132Element_SIG_POLICY_QUALIFIER.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES132Element_SIG_POLICY_QUALIFIERS.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES132Element_SIGNATURE_POLICY_ID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES132Element_SIGNATURE_POLICY_IDENTIFIER.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES132Element_SIGNATURE_POLICY_IMPLIED.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES132Element_SIGNATURE_PRODUCTION_PLACE.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_PRODUCTION_PLACE_V2", XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2.TagName(), "SignatureProductionPlaceV2"},
		{"SIGNATURE_TIMESTAMP", XAdES132Element_SIGNATURE_TIMESTAMP.TagName(), "SignatureTimeStamp"},
		{"SIGNED_ASSERTION", XAdES132Element_SIGNED_ASSERTION.TagName(), "SignedAssertion"},
		{"SIGNED_ASSERTIONS", XAdES132Element_SIGNED_ASSERTIONS.TagName(), "SignedAssertions"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES132Element_SIGNED_PROPERTIES.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES132Element_SIGNED_SIGNATURE_PROPERTIES.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES132Element_SIGNER_ROLE.TagName(), "SignerRole"},
		{"SIGNER_ROLE_V2", XAdES132Element_SIGNER_ROLE_V2.TagName(), "SignerRoleV2"},
		{"SIGNING_CERTIFICATE", XAdES132Element_SIGNING_CERTIFICATE.TagName(), "SigningCertificate"},
		{"SIGNING_CERTIFICATE_V2", XAdES132Element_SIGNING_CERTIFICATE_V2.TagName(), "SigningCertificateV2"},
		{"SIGNING_TIME", XAdES132Element_SIGNING_TIME.TagName(), "SigningTime"},
		{"SP_URI", XAdES132Element_SP_URI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES132Element_SP_USER_NOTICE.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES132Element_STATE_OR_PROVINCE.TagName(), "StateOrProvince"},
		{"STREET_ADDRESS", XAdES132Element_STREET_ADDRESS.TagName(), "StreetAddress"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES132Element_UNSIGNED_PROPERTIES.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES.TagName(), "UnsignedSignatureProperties"},
		{"X509_ATTRIBUTE_CERTIFICATE", XAdES132Element_X509_ATTRIBUTE_CERTIFICATE.TagName(), "X509AttributeCertificate"},
		{"XADES_TIMESTAMP", XAdES132Element_XADES_TIMESTAMP.TagName(), "XAdESTimeStamp"},
		{"XML_TIMESTAMP", XAdES132Element_XML_TIMESTAMP.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES132Element_Namespace(t *testing.T) {
	if got := XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP.URI(); got != XAdESNamespace_XADES_132.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespace_XADES_132.Uri())
	}
	if !XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP.IsSameTagName(XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES132Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES132Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES132Element_ElementGetters(t *testing.T) {
	var e XAdES132Element = XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES132Element_ALL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES132Element_ALL_SIGNED_DATA_OBJECTS {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES132Element_ALL_SIGNED_DATA_OBJECTS", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES132Element_ANY {
			t.Errorf("ElementAny() = %v, want XAdES132Element_ANY", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES132Element_ARCHIVE_TIMESTAMP {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES132Element_ARCHIVE_TIMESTAMP", got)
		}
	})
	t.Run("ElementAttrAuthoritiesCertValues", func(t *testing.T) {
		if got := e.ElementAttrAuthoritiesCertValues(); got != XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES {
			t.Errorf("ElementAttrAuthoritiesCertValues() = %v, want XAdES132Element_ATTR_AUTHORITIES_CERT_VALUES", got)
		}
	})
	t.Run("ElementAttributeCertificateRefs", func(t *testing.T) {
		if got := e.ElementAttributeCertificateRefs(); got != XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS {
			t.Errorf("ElementAttributeCertificateRefs() = %v, want XAdES132Element_ATTRIBUTE_CERTIFICATE_REFS", got)
		}
	})
	t.Run("ElementAttributeRevocationRefs", func(t *testing.T) {
		if got := e.ElementAttributeRevocationRefs(); got != XAdES132Element_ATTRIBUTE_REVOCATION_REFS {
			t.Errorf("ElementAttributeRevocationRefs() = %v, want XAdES132Element_ATTRIBUTE_REVOCATION_REFS", got)
		}
	})
	t.Run("ElementAttributeRevocationValues", func(t *testing.T) {
		if got := e.ElementAttributeRevocationValues(); got != XAdES132Element_ATTRIBUTE_REVOCATION_VALUES {
			t.Errorf("ElementAttributeRevocationValues() = %v, want XAdES132Element_ATTRIBUTE_REVOCATION_VALUES", got)
		}
	})
	t.Run("ElementByKey", func(t *testing.T) {
		if got := e.ElementByKey(); got != XAdES132Element_BY_KEY {
			t.Errorf("ElementByKey() = %v, want XAdES132Element_BY_KEY", got)
		}
	})
	t.Run("ElementByName", func(t *testing.T) {
		if got := e.ElementByName(); got != XAdES132Element_BY_NAME {
			t.Errorf("ElementByName() = %v, want XAdES132Element_BY_NAME", got)
		}
	})
	t.Run("ElementCert", func(t *testing.T) {
		if got := e.ElementCert(); got != XAdES132Element_CERT {
			t.Errorf("ElementCert() = %v, want XAdES132Element_CERT", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES132Element_CERT_DIGEST {
			t.Errorf("ElementCertDigest() = %v, want XAdES132Element_CERT_DIGEST", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES132Element_CERT_REFS {
			t.Errorf("ElementCertRefs() = %v, want XAdES132Element_CERT_REFS", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES132Element_CERTIFICATE_VALUES {
			t.Errorf("ElementCertificateValues() = %v, want XAdES132Element_CERTIFICATE_VALUES", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES132Element_CERTIFIED_ROLE {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES132Element_CERTIFIED_ROLE", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES132Element_CERTIFIED_ROLES {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES132Element_CERTIFIED_ROLES", got)
		}
	})
	t.Run("ElementCertifiedRolesV2", func(t *testing.T) {
		if got := e.ElementCertifiedRolesV2(); got != XAdES132Element_CERTIFIED_ROLES_V2 {
			t.Errorf("ElementCertifiedRolesV2() = %v, want XAdES132Element_CERTIFIED_ROLES_V2", got)
		}
	})
	t.Run("ElementCity", func(t *testing.T) {
		if got := e.ElementCity(); got != XAdES132Element_CITY {
			t.Errorf("ElementCity() = %v, want XAdES132Element_CITY", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES132Element_CLAIMED_ROLE {
			t.Errorf("ElementClaimedRole() = %v, want XAdES132Element_CLAIMED_ROLE", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES132Element_CLAIMED_ROLES {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES132Element_CLAIMED_ROLES", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES132Element_COMMITMENT_TYPE_ID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES132Element_COMMITMENT_TYPE_ID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES132Element_COMMITMENT_TYPE_INDICATION {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES132Element_COMMITMENT_TYPE_INDICATION", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES132Element_COMMITMENT_TYPE_QUALIFIER {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES132Element_COMMITMENT_TYPE_QUALIFIER", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES132Element_COMMITMENT_TYPE_QUALIFIERS {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES132Element_COMMITMENT_TYPE_QUALIFIERS", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES132Element_COMPLETE_CERTIFICATE_REFS {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES132Element_COMPLETE_CERTIFICATE_REFS", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES132Element_COMPLETE_REVOCATION_REFS {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES132Element_COMPLETE_REVOCATION_REFS", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES132Element_COUNTER_SIGNATURE {
			t.Errorf("ElementCounterSignature() = %v, want XAdES132Element_COUNTER_SIGNATURE", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES132Element_COUNTRY_NAME {
			t.Errorf("ElementCountryName() = %v, want XAdES132Element_COUNTRY_NAME", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES132Element_CRL_IDENTIFIER {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES132Element_CRL_IDENTIFIER", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES132Element_CRL_REF {
			t.Errorf("ElementCRLRef() = %v, want XAdES132Element_CRL_REF", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES132Element_CRL_REFS {
			t.Errorf("ElementCRLRefs() = %v, want XAdES132Element_CRL_REFS", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES132Element_CRL_VALUES {
			t.Errorf("ElementCRLValues() = %v, want XAdES132Element_CRL_VALUES", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES132Element_DATA_OBJECT_FORMAT {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES132Element_DATA_OBJECT_FORMAT", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES132Element_DESCRIPTION {
			t.Errorf("ElementDescription() = %v, want XAdES132Element_DESCRIPTION", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES132Element_DIGEST_ALG_AND_VALUE {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES132Element_DIGEST_ALG_AND_VALUE", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES132Element_DOCUMENTATION_REFERENCE {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES132Element_DOCUMENTATION_REFERENCE", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES132Element_DOCUMENTATION_REFERENCES {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES132Element_DOCUMENTATION_REFERENCES", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES132Element_ENCAPSULATED_CRL_VALUE {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES132Element_ENCAPSULATED_CRL_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES132Element_ENCAPSULATED_OCSP_VALUE {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES132Element_ENCAPSULATED_OCSP_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES132Element_ENCAPSULATED_PKI_DATA {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES132Element_ENCAPSULATED_PKI_DATA", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES132Element_ENCAPSULATED_TIMESTAMP {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES132Element_ENCAPSULATED_TIMESTAMP", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES132Element_ENCAPSULATED_X509_CERTIFICATE {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES132Element_ENCAPSULATED_X509_CERTIFICATE", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES132Element_ENCODING {
			t.Errorf("ElementEncoding() = %v, want XAdES132Element_ENCODING", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES132Element_EXPLICIT_TEXT {
			t.Errorf("ElementExplicitText() = %v, want XAdES132Element_EXPLICIT_TEXT", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES132Element_IDENTIFIER {
			t.Errorf("ElementIdentifier() = %v, want XAdES132Element_IDENTIFIER", got)
		}
	})
	t.Run("ElementInclude", func(t *testing.T) {
		if got := e.ElementInclude(); got != XAdES132Element_INCLUDE {
			t.Errorf("ElementInclude() = %v, want XAdES132Element_INCLUDE", got)
		}
	})
	t.Run("ElementIndividualDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES132Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES132Element_INT {
			t.Errorf("Elementint() = %v, want XAdES132Element_INT", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES132Element_ISSUE_TIME {
			t.Errorf("ElementIssueTime() = %v, want XAdES132Element_ISSUE_TIME", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES132Element_ISSUER {
			t.Errorf("ElementIssuer() = %v, want XAdES132Element_ISSUER", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES132Element_ISSUER_SERIAL {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES132Element_ISSUER_SERIAL", got)
		}
	})
	t.Run("ElementIssuerSerialV2", func(t *testing.T) {
		if got := e.ElementIssuerSerialV2(); got != XAdES132Element_ISSUER_SERIAL_V2 {
			t.Errorf("ElementIssuerSerialV2() = %v, want XAdES132Element_ISSUER_SERIAL_V2", got)
		}
	})
	t.Run("ElementMimeType", func(t *testing.T) {
		if got := e.ElementMimeType(); got != XAdES132Element_MIME_TYPE {
			t.Errorf("ElementMimeType() = %v, want XAdES132Element_MIME_TYPE", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES132Element_NOTICE_NUMBERS {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES132Element_NOTICE_NUMBERS", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES132Element_NOTICE_REF {
			t.Errorf("ElementNoticeRef() = %v, want XAdES132Element_NOTICE_REF", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES132Element_NUMBER {
			t.Errorf("ElementNumber() = %v, want XAdES132Element_NUMBER", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES132Element_OBJECT_IDENTIFIER {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES132Element_OBJECT_IDENTIFIER", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES132Element_OBJECT_REFERENCE {
			t.Errorf("ElementObjectReference() = %v, want XAdES132Element_OBJECT_REFERENCE", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES132Element_OCSP_IDENTIFIER {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES132Element_OCSP_IDENTIFIER", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES132Element_OCSP_REF {
			t.Errorf("ElementOCSPRef() = %v, want XAdES132Element_OCSP_REF", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES132Element_OCSP_REFS {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES132Element_OCSP_REFS", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES132Element_OCSP_VALUES {
			t.Errorf("ElementOCSPValues() = %v, want XAdES132Element_OCSP_VALUES", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES132Element_ORGANIZATION {
			t.Errorf("ElementOrganization() = %v, want XAdES132Element_ORGANIZATION", got)
		}
	})
	t.Run("ElementOtherAttributeCertificate", func(t *testing.T) {
		if got := e.ElementOtherAttributeCertificate(); got != XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE {
			t.Errorf("ElementOtherAttributeCertificate() = %v, want XAdES132Element_OTHER_ATTRIBUTE_CERTIFICATE", got)
		}
	})
	t.Run("ElementOtherCertificate", func(t *testing.T) {
		if got := e.ElementOtherCertificate(); got != XAdES132Element_OTHER_CERTIFICATE {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES132Element_OTHER_CERTIFICATE", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES132Element_OTHER_REF {
			t.Errorf("ElementOtherRef() = %v, want XAdES132Element_OTHER_REF", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES132Element_OTHER_REFS {
			t.Errorf("ElementOtherRefs() = %v, want XAdES132Element_OTHER_REFS", got)
		}
	})
	t.Run("ElementOtherTimeStamp", func(t *testing.T) {
		if got := e.ElementOtherTimeStamp(); got != XAdES132Element_OTHER_TIMESTAMP {
			t.Errorf("ElementOtherTimeStamp() = %v, want XAdES132Element_OTHER_TIMESTAMP", got)
		}
	})
	t.Run("ElementOtherValue", func(t *testing.T) {
		if got := e.ElementOtherValue(); got != XAdES132Element_OTHER_VALUE {
			t.Errorf("ElementOtherValue() = %v, want XAdES132Element_OTHER_VALUE", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES132Element_OTHER_VALUES {
			t.Errorf("ElementOtherValues() = %v, want XAdES132Element_OTHER_VALUES", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES132Element_POSTAL_CODE {
			t.Errorf("ElementPostalCode() = %v, want XAdES132Element_POSTAL_CODE", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES132Element_PRODUCED_AT {
			t.Errorf("ElementProducedAt() = %v, want XAdES132Element_PRODUCED_AT", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES132Element_QUALIFYING_PROPERTIES {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES132Element_QUALIFYING_PROPERTIES", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES132Element_QUALIFYING_PROPERTIES_REFERENCE", got)
		}
	})
	t.Run("ElementReferenceInfo", func(t *testing.T) {
		if got := e.ElementReferenceInfo(); got != XAdES132Element_REFERENCE_INFO {
			t.Errorf("ElementReferenceInfo() = %v, want XAdES132Element_REFERENCE_INFO", got)
		}
	})
	t.Run("ElementRefsOnlyTimeStamp", func(t *testing.T) {
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES132Element_REFS_ONLY_TIMESTAMP {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES132Element_REFS_ONLY_TIMESTAMP", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES132Element_RESPONDER_ID {
			t.Errorf("ElementResponderID() = %v, want XAdES132Element_RESPONDER_ID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES132Element_REVOCATION_VALUES {
			t.Errorf("ElementRevocationValues() = %v, want XAdES132Element_REVOCATION_VALUES", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES132Element_SIG_AND_REFS_TIMESTAMP {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES132Element_SIG_AND_REFS_TIMESTAMP", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES132Element_SIG_POLICY_HASH {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES132Element_SIG_POLICY_HASH", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES132Element_SIG_POLICY_ID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES132Element_SIG_POLICY_ID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES132Element_SIG_POLICY_QUALIFIER {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES132Element_SIG_POLICY_QUALIFIER", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES132Element_SIG_POLICY_QUALIFIERS {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES132Element_SIG_POLICY_QUALIFIERS", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES132Element_SIGNATURE_POLICY_ID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES132Element_SIGNATURE_POLICY_ID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES132Element_SIGNATURE_POLICY_IDENTIFIER {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES132Element_SIGNATURE_POLICY_IDENTIFIER", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES132Element_SIGNATURE_POLICY_IMPLIED {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES132Element_SIGNATURE_POLICY_IMPLIED", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES132Element_SIGNATURE_PRODUCTION_PLACE {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES132Element_SIGNATURE_PRODUCTION_PLACE", got)
		}
	})
	t.Run("ElementSignatureProductionPlaceV2", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlaceV2(); got != XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2 {
			t.Errorf("ElementSignatureProductionPlaceV2() = %v, want XAdES132Element_SIGNATURE_PRODUCTION_PLACE_V2", got)
		}
	})
	t.Run("ElementSignatureTimeStamp", func(t *testing.T) {
		if got := e.ElementSignatureTimeStamp(); got != XAdES132Element_SIGNATURE_TIMESTAMP {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES132Element_SIGNATURE_TIMESTAMP", got)
		}
	})
	t.Run("ElementSignedAssertion", func(t *testing.T) {
		if got := e.ElementSignedAssertion(); got != XAdES132Element_SIGNED_ASSERTION {
			t.Errorf("ElementSignedAssertion() = %v, want XAdES132Element_SIGNED_ASSERTION", got)
		}
	})
	t.Run("ElementSignedAssertions", func(t *testing.T) {
		if got := e.ElementSignedAssertions(); got != XAdES132Element_SIGNED_ASSERTIONS {
			t.Errorf("ElementSignedAssertions() = %v, want XAdES132Element_SIGNED_ASSERTIONS", got)
		}
	})
	t.Run("ElementSignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementSignedDataObjectProperties(); got != XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES132Element_SIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES132Element_SIGNED_PROPERTIES {
			t.Errorf("ElementSignedProperties() = %v, want XAdES132Element_SIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES132Element_SIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES132Element_SIGNED_SIGNATURE_PROPERTIES", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES132Element_SIGNER_ROLE {
			t.Errorf("ElementSignerRole() = %v, want XAdES132Element_SIGNER_ROLE", got)
		}
	})
	t.Run("ElementSignerRoleV2", func(t *testing.T) {
		if got := e.ElementSignerRoleV2(); got != XAdES132Element_SIGNER_ROLE_V2 {
			t.Errorf("ElementSignerRoleV2() = %v, want XAdES132Element_SIGNER_ROLE_V2", got)
		}
	})
	t.Run("ElementSigningCertificate", func(t *testing.T) {
		if got := e.ElementSigningCertificate(); got != XAdES132Element_SIGNING_CERTIFICATE {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES132Element_SIGNING_CERTIFICATE", got)
		}
	})
	t.Run("ElementSigningCertificateV2", func(t *testing.T) {
		if got := e.ElementSigningCertificateV2(); got != XAdES132Element_SIGNING_CERTIFICATE_V2 {
			t.Errorf("ElementSigningCertificateV2() = %v, want XAdES132Element_SIGNING_CERTIFICATE_V2", got)
		}
	})
	t.Run("ElementSigningTime", func(t *testing.T) {
		if got := e.ElementSigningTime(); got != XAdES132Element_SIGNING_TIME {
			t.Errorf("ElementSigningTime() = %v, want XAdES132Element_SIGNING_TIME", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES132Element_SP_URI {
			t.Errorf("ElementSPURI() = %v, want XAdES132Element_SP_URI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES132Element_SP_USER_NOTICE {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES132Element_SP_USER_NOTICE", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES132Element_STATE_OR_PROVINCE {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES132Element_STATE_OR_PROVINCE", got)
		}
	})
	t.Run("ElementStreetAddress", func(t *testing.T) {
		if got := e.ElementStreetAddress(); got != XAdES132Element_STREET_ADDRESS {
			t.Errorf("ElementStreetAddress() = %v, want XAdES132Element_STREET_ADDRESS", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperties", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES132Element_UNSIGNED_DATA_OBJECT_PROPERTY", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES132Element_UNSIGNED_PROPERTIES {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES132Element_UNSIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES132Element_UNSIGNED_SIGNATURE_PROPERTIES", got)
		}
	})
	t.Run("ElementX509AttributeCertificate", func(t *testing.T) {
		if got := e.ElementX509AttributeCertificate(); got != XAdES132Element_X509_ATTRIBUTE_CERTIFICATE {
			t.Errorf("ElementX509AttributeCertificate() = %v, want XAdES132Element_X509_ATTRIBUTE_CERTIFICATE", got)
		}
	})
	t.Run("ElementXAdESTimeStamp", func(t *testing.T) {
		if got := e.ElementXAdESTimeStamp(); got != XAdES132Element_XADES_TIMESTAMP {
			t.Errorf("ElementXAdESTimeStamp() = %v, want XAdES132Element_XADES_TIMESTAMP", got)
		}
	})
	t.Run("ElementXMLTimeStamp", func(t *testing.T) {
		if got := e.ElementXMLTimeStamp(); got != XAdES132Element_XML_TIMESTAMP {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES132Element_XML_TIMESTAMP", got)
		}
	})
}
