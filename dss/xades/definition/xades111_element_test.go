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
		{"ALL_DATA_OBJECTS_TIMESTAMP", XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName(), "AllDataObjectsTimeStamp"},
		{"ALL_SIGNED_DATA_OBJECTS", XAdES111Element_ALL_SIGNED_DATA_OBJECTS.TagName(), "AllSignedDataObjects"},
		{"ANY", XAdES111Element_ANY.TagName(), "Any"},
		{"ARCHIVE_TIMESTAMP", XAdES111Element_ARCHIVE_TIMESTAMP.TagName(), "ArchiveTimeStamp"},
		{"CERT", XAdES111Element_CERT.TagName(), "Cert"},
		{"CERT_DIGEST", XAdES111Element_CERT_DIGEST.TagName(), "CertDigest"},
		{"CERT_REFS", XAdES111Element_CERT_REFS.TagName(), "CertRefs"},
		{"CERTIFICATE_VALUES", XAdES111Element_CERTIFICATE_VALUES.TagName(), "CertificateValues"},
		{"CERTIFIED_ROLE", XAdES111Element_CERTIFIED_ROLE.TagName(), "CertifiedRole"},
		{"CERTIFIED_ROLES", XAdES111Element_CERTIFIED_ROLES.TagName(), "CertifiedRoles"},
		{"CITY", XAdES111Element_CITY.TagName(), "City"},
		{"CLAIMED_ROLE", XAdES111Element_CLAIMED_ROLE.TagName(), "ClaimedRole"},
		{"CLAIMED_ROLES", XAdES111Element_CLAIMED_ROLES.TagName(), "ClaimedRoles"},
		{"COMMITMENT_TYPE_ID", XAdES111Element_COMMITMENT_TYPE_ID.TagName(), "CommitmentTypeId"},
		{"COMMITMENT_TYPE_INDICATION", XAdES111Element_COMMITMENT_TYPE_INDICATION.TagName(), "CommitmentTypeIndication"},
		{"COMMITMENT_TYPE_QUALIFIER", XAdES111Element_COMMITMENT_TYPE_QUALIFIER.TagName(), "CommitmentTypeQualifier"},
		{"COMMITMENT_TYPE_QUALIFIERS", XAdES111Element_COMMITMENT_TYPE_QUALIFIERS.TagName(), "CommitmentTypeQualifiers"},
		{"COMPLETE_CERTIFICATE_REFS", XAdES111Element_COMPLETE_CERTIFICATE_REFS.TagName(), "CompleteCertificateRefs"},
		{"COMPLETE_REVOCATION_REFS", XAdES111Element_COMPLETE_REVOCATION_REFS.TagName(), "CompleteRevocationRefs"},
		{"COUNTER_SIGNATURE", XAdES111Element_COUNTER_SIGNATURE.TagName(), "CounterSignature"},
		{"COUNTRY_NAME", XAdES111Element_COUNTRY_NAME.TagName(), "CountryName"},
		{"CRL_IDENTIFIER", XAdES111Element_CRL_IDENTIFIER.TagName(), "CRLIdentifier"},
		{"CRL_REF", XAdES111Element_CRL_REF.TagName(), "CRLRef"},
		{"CRL_REFS", XAdES111Element_CRL_REFS.TagName(), "CRLRefs"},
		{"CRL_VALUES", XAdES111Element_CRL_VALUES.TagName(), "CRLValues"},
		{"DATA_OBJECT_FORMAT", XAdES111Element_DATA_OBJECT_FORMAT.TagName(), "DataObjectFormat"},
		{"DESCRIPTION", XAdES111Element_DESCRIPTION.TagName(), "Description"},
		{"DIGEST_ALG_AND_VALUE", XAdES111Element_DIGEST_ALG_AND_VALUE.TagName(), "DigestAlgAndValue"},
		{"DIGEST_METHOD", XAdES111Element_DIGEST_METHOD.TagName(), "DigestMethod"},
		{"DIGEST_VALUE", XAdES111Element_DIGEST_VALUE.TagName(), "DigestValue"},
		{"DOCUMENTATION_REFERENCE", XAdES111Element_DOCUMENTATION_REFERENCE.TagName(), "DocumentationReference"},
		{"DOCUMENTATION_REFERENCES", XAdES111Element_DOCUMENTATION_REFERENCES.TagName(), "DocumentationReferences"},
		{"ENCAPSULATED_CRL_VALUE", XAdES111Element_ENCAPSULATED_CRL_VALUE.TagName(), "EncapsulatedCRLValue"},
		{"ENCAPSULATED_OCSP_VALUE", XAdES111Element_ENCAPSULATED_OCSP_VALUE.TagName(), "EncapsulatedOCSPValue"},
		{"ENCAPSULATED_PKI_DATA", XAdES111Element_ENCAPSULATED_PKI_DATA.TagName(), "EncapsulatedPKIData"},
		{"ENCAPSULATED_TIMESTAMP", XAdES111Element_ENCAPSULATED_TIMESTAMP.TagName(), "EncapsulatedTimeStamp"},
		{"ENCAPSULATED_X509_CERTIFICATE", XAdES111Element_ENCAPSULATED_X509_CERTIFICATE.TagName(), "EncapsulatedX509Certificate"},
		{"ENCODING", XAdES111Element_ENCODING.TagName(), "Encoding"},
		{"EXPLICIT_TEXT", XAdES111Element_EXPLICIT_TEXT.TagName(), "ExplicitText"},
		{"HASH_DATA_INFO", XAdES111Element_HASH_DATA_INFO.TagName(), "HashDataInfo"},
		{"IDENTIFIER", XAdES111Element_IDENTIFIER.TagName(), "Identifier"},
		{"INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP.TagName(), "IndividualDataObjectsTimeStamp"},
		{"INT", XAdES111Element_INT.TagName(), "int"},
		{"ISSUE_TIME", XAdES111Element_ISSUE_TIME.TagName(), "IssueTime"},
		{"ISSUER", XAdES111Element_ISSUER.TagName(), "Issuer"},
		{"ISSUER_SERIAL", XAdES111Element_ISSUER_SERIAL.TagName(), "IssuerSerial"},
		{"MIME_TYPE", XAdES111Element_MIME_TYPE.TagName(), "MimeType"},
		{"NOTICE_NUMBERS", XAdES111Element_NOTICE_NUMBERS.TagName(), "NoticeNumbers"},
		{"NOTICE_REF", XAdES111Element_NOTICE_REF.TagName(), "NoticeRef"},
		{"NUMBER", XAdES111Element_NUMBER.TagName(), "Number"},
		{"OBJECT_IDENTIFIER", XAdES111Element_OBJECT_IDENTIFIER.TagName(), "ObjectIdentifier"},
		{"OBJECT_REFERENCE", XAdES111Element_OBJECT_REFERENCE.TagName(), "ObjectReference"},
		{"OCSP_IDENTIFIER", XAdES111Element_OCSP_IDENTIFIER.TagName(), "OCSPIdentifier"},
		{"OCSP_REF", XAdES111Element_OCSP_REF.TagName(), "OCSPRef"},
		{"OCSP_REFS", XAdES111Element_OCSP_REFS.TagName(), "OCSPRefs"},
		{"OCSP_VALUES", XAdES111Element_OCSP_VALUES.TagName(), "OCSPValues"},
		{"ORGANIZATION", XAdES111Element_ORGANIZATION.TagName(), "Organization"},
		{"OTHER_CERTIFICATE", XAdES111Element_OTHER_CERTIFICATE.TagName(), "OtherCertificate"},
		{"OTHER_REF", XAdES111Element_OTHER_REF.TagName(), "OtherRef"},
		{"OTHER_REFS", XAdES111Element_OTHER_REFS.TagName(), "OtherRefs"},
		{"OTHER_VALUE", XAdES111Element_OTHER_VALUE.TagName(), "OtherValue"},
		{"OTHER_VALUES", XAdES111Element_OTHER_VALUES.TagName(), "OtherValues"},
		{"POSTAL_CODE", XAdES111Element_POSTAL_CODE.TagName(), "PostalCode"},
		{"PRODUCED_AT", XAdES111Element_PRODUCED_AT.TagName(), "ProducedAt"},
		{"QUALIFYING_PROPERTIES", XAdES111Element_QUALIFYING_PROPERTIES.TagName(), "QualifyingProperties"},
		{"QUALIFYING_PROPERTIES_REFERENCE", XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE.TagName(), "QualifyingPropertiesReference"},
		{"REFS_ONLY_TIMESTAMP", XAdES111Element_REFS_ONLY_TIMESTAMP.TagName(), "RefsOnlyTimeStamp"},
		{"RESPONDER_ID", XAdES111Element_RESPONDER_ID.TagName(), "ResponderID"},
		{"REVOCATION_VALUES", XAdES111Element_REVOCATION_VALUES.TagName(), "RevocationValues"},
		{"SIG_AND_REFS_TIMESTAMP", XAdES111Element_SIG_AND_REFS_TIMESTAMP.TagName(), "SigAndRefsTimeStamp"},
		{"SIG_POLICY_HASH", XAdES111Element_SIG_POLICY_HASH.TagName(), "SigPolicyHash"},
		{"SIG_POLICY_ID", XAdES111Element_SIG_POLICY_ID.TagName(), "SigPolicyId"},
		{"SIG_POLICY_QUALIFIER", XAdES111Element_SIG_POLICY_QUALIFIER.TagName(), "SigPolicyQualifier"},
		{"SIG_POLICY_QUALIFIERS", XAdES111Element_SIG_POLICY_QUALIFIERS.TagName(), "SigPolicyQualifiers"},
		{"SIGNATURE_POLICY_ID", XAdES111Element_SIGNATURE_POLICY_ID.TagName(), "SignaturePolicyId"},
		{"SIGNATURE_POLICY_IDENTIFIER", XAdES111Element_SIGNATURE_POLICY_IDENTIFIER.TagName(), "SignaturePolicyIdentifier"},
		{"SIGNATURE_POLICY_IMPLIED", XAdES111Element_SIGNATURE_POLICY_IMPLIED.TagName(), "SignaturePolicyImplied"},
		{"SIGNATURE_PRODUCTION_PLACE", XAdES111Element_SIGNATURE_PRODUCTION_PLACE.TagName(), "SignatureProductionPlace"},
		{"SIGNATURE_TIMESTAMP", XAdES111Element_SIGNATURE_TIMESTAMP.TagName(), "SignatureTimeStamp"},
		{"SIGNED_DATA_OBJECT_PROPERTIES", XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES.TagName(), "SignedDataObjectProperties"},
		{"SIGNED_PROPERTIES", XAdES111Element_SIGNED_PROPERTIES.TagName(), "SignedProperties"},
		{"SIGNED_SIGNATURE_PROPERTIES", XAdES111Element_SIGNED_SIGNATURE_PROPERTIES.TagName(), "SignedSignatureProperties"},
		{"SIGNER_ROLE", XAdES111Element_SIGNER_ROLE.TagName(), "SignerRole"},
		{"SIGNING_CERTIFICATE", XAdES111Element_SIGNING_CERTIFICATE.TagName(), "SigningCertificate"},
		{"SIGNING_TIME", XAdES111Element_SIGNING_TIME.TagName(), "SigningTime"},
		{"SP_URI", XAdES111Element_SP_URI.TagName(), "SPURI"},
		{"SP_USER_NOTICE", XAdES111Element_SP_USER_NOTICE.TagName(), "SPUserNotice"},
		{"STATE_OR_PROVINCE", XAdES111Element_STATE_OR_PROVINCE.TagName(), "StateOrProvince"},
		{"TIMESTAMP", XAdES111Element_TIMESTAMP.TagName(), "TimeStamp"},
		{"TRANSFORMS", XAdES111Element_TRANSFORMS.TagName(), "Transforms"},
		{"UNSIGNED_DATA_OBJECT_PROPERTIES", XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES.TagName(), "UnsignedDataObjectProperties"},
		{"UNSIGNED_DATA_OBJECT_PROPERTY", XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY.TagName(), "UnsignedDataObjectProperty"},
		{"UNSIGNED_PROPERTIES", XAdES111Element_UNSIGNED_PROPERTIES.TagName(), "UnsignedProperties"},
		{"UNSIGNED_SIGNATURE_PROPERTIES", XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES.TagName(), "UnsignedSignatureProperties"},
		{"XML_TIMESTAMP", XAdES111Element_XML_TIMESTAMP.TagName(), "XMLTimeStamp"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestXAdES111Element_Namespace(t *testing.T) {
	if got := XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP.URI(); got != XAdESNamespace_XADES_111.Uri() {
		t.Errorf("URI() = %q, want %q", got, XAdESNamespace_XADES_111.Uri())
	}
	if !XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP.IsSameTagName(XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP.TagName()) {
		t.Errorf("IsSameTagName should match its own tag name")
	}
}

// TestXAdES111Element_ElementGetters exercises every XAdESElement getter implemented on
// XAdES111Element, checking it returns the same-named local constant (or panics with the
// Java UnsupportedOperationException message, for versions that do not define the element).
func TestXAdES111Element_ElementGetters(t *testing.T) {
	var e XAdES111Element = XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP
	t.Run("ElementAllDataObjectsTimeStamp", func(t *testing.T) {
		if got := e.ElementAllDataObjectsTimeStamp(); got != XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementAllDataObjectsTimeStamp() = %v, want XAdES111Element_ALL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("ElementAllSignedDataObjects", func(t *testing.T) {
		if got := e.ElementAllSignedDataObjects(); got != XAdES111Element_ALL_SIGNED_DATA_OBJECTS {
			t.Errorf("ElementAllSignedDataObjects() = %v, want XAdES111Element_ALL_SIGNED_DATA_OBJECTS", got)
		}
	})
	t.Run("ElementAny", func(t *testing.T) {
		if got := e.ElementAny(); got != XAdES111Element_ANY {
			t.Errorf("ElementAny() = %v, want XAdES111Element_ANY", got)
		}
	})
	t.Run("ElementArchiveTimeStamp", func(t *testing.T) {
		if got := e.ElementArchiveTimeStamp(); got != XAdES111Element_ARCHIVE_TIMESTAMP {
			t.Errorf("ElementArchiveTimeStamp() = %v, want XAdES111Element_ARCHIVE_TIMESTAMP", got)
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
		if got := e.ElementCert(); got != XAdES111Element_CERT {
			t.Errorf("ElementCert() = %v, want XAdES111Element_CERT", got)
		}
	})
	t.Run("ElementCertDigest", func(t *testing.T) {
		if got := e.ElementCertDigest(); got != XAdES111Element_CERT_DIGEST {
			t.Errorf("ElementCertDigest() = %v, want XAdES111Element_CERT_DIGEST", got)
		}
	})
	t.Run("ElementCertRefs", func(t *testing.T) {
		if got := e.ElementCertRefs(); got != XAdES111Element_CERT_REFS {
			t.Errorf("ElementCertRefs() = %v, want XAdES111Element_CERT_REFS", got)
		}
	})
	t.Run("ElementCertificateValues", func(t *testing.T) {
		if got := e.ElementCertificateValues(); got != XAdES111Element_CERTIFICATE_VALUES {
			t.Errorf("ElementCertificateValues() = %v, want XAdES111Element_CERTIFICATE_VALUES", got)
		}
	})
	t.Run("ElementCertifiedRole", func(t *testing.T) {
		if got := e.ElementCertifiedRole(); got != XAdES111Element_CERTIFIED_ROLE {
			t.Errorf("ElementCertifiedRole() = %v, want XAdES111Element_CERTIFIED_ROLE", got)
		}
	})
	t.Run("ElementCertifiedRoles", func(t *testing.T) {
		if got := e.ElementCertifiedRoles(); got != XAdES111Element_CERTIFIED_ROLES {
			t.Errorf("ElementCertifiedRoles() = %v, want XAdES111Element_CERTIFIED_ROLES", got)
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
		if got := e.ElementCity(); got != XAdES111Element_CITY {
			t.Errorf("ElementCity() = %v, want XAdES111Element_CITY", got)
		}
	})
	t.Run("ElementClaimedRole", func(t *testing.T) {
		if got := e.ElementClaimedRole(); got != XAdES111Element_CLAIMED_ROLE {
			t.Errorf("ElementClaimedRole() = %v, want XAdES111Element_CLAIMED_ROLE", got)
		}
	})
	t.Run("ElementClaimedRoles", func(t *testing.T) {
		if got := e.ElementClaimedRoles(); got != XAdES111Element_CLAIMED_ROLES {
			t.Errorf("ElementClaimedRoles() = %v, want XAdES111Element_CLAIMED_ROLES", got)
		}
	})
	t.Run("ElementCommitmentTypeId", func(t *testing.T) {
		if got := e.ElementCommitmentTypeId(); got != XAdES111Element_COMMITMENT_TYPE_ID {
			t.Errorf("ElementCommitmentTypeId() = %v, want XAdES111Element_COMMITMENT_TYPE_ID", got)
		}
	})
	t.Run("ElementCommitmentTypeIndication", func(t *testing.T) {
		if got := e.ElementCommitmentTypeIndication(); got != XAdES111Element_COMMITMENT_TYPE_INDICATION {
			t.Errorf("ElementCommitmentTypeIndication() = %v, want XAdES111Element_COMMITMENT_TYPE_INDICATION", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifier", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifier(); got != XAdES111Element_COMMITMENT_TYPE_QUALIFIER {
			t.Errorf("ElementCommitmentTypeQualifier() = %v, want XAdES111Element_COMMITMENT_TYPE_QUALIFIER", got)
		}
	})
	t.Run("ElementCommitmentTypeQualifiers", func(t *testing.T) {
		if got := e.ElementCommitmentTypeQualifiers(); got != XAdES111Element_COMMITMENT_TYPE_QUALIFIERS {
			t.Errorf("ElementCommitmentTypeQualifiers() = %v, want XAdES111Element_COMMITMENT_TYPE_QUALIFIERS", got)
		}
	})
	t.Run("ElementCompleteCertificateRefs", func(t *testing.T) {
		if got := e.ElementCompleteCertificateRefs(); got != XAdES111Element_COMPLETE_CERTIFICATE_REFS {
			t.Errorf("ElementCompleteCertificateRefs() = %v, want XAdES111Element_COMPLETE_CERTIFICATE_REFS", got)
		}
	})
	t.Run("ElementCompleteRevocationRefs", func(t *testing.T) {
		if got := e.ElementCompleteRevocationRefs(); got != XAdES111Element_COMPLETE_REVOCATION_REFS {
			t.Errorf("ElementCompleteRevocationRefs() = %v, want XAdES111Element_COMPLETE_REVOCATION_REFS", got)
		}
	})
	t.Run("ElementCounterSignature", func(t *testing.T) {
		if got := e.ElementCounterSignature(); got != XAdES111Element_COUNTER_SIGNATURE {
			t.Errorf("ElementCounterSignature() = %v, want XAdES111Element_COUNTER_SIGNATURE", got)
		}
	})
	t.Run("ElementCountryName", func(t *testing.T) {
		if got := e.ElementCountryName(); got != XAdES111Element_COUNTRY_NAME {
			t.Errorf("ElementCountryName() = %v, want XAdES111Element_COUNTRY_NAME", got)
		}
	})
	t.Run("ElementCRLIdentifier", func(t *testing.T) {
		if got := e.ElementCRLIdentifier(); got != XAdES111Element_CRL_IDENTIFIER {
			t.Errorf("ElementCRLIdentifier() = %v, want XAdES111Element_CRL_IDENTIFIER", got)
		}
	})
	t.Run("ElementCRLRef", func(t *testing.T) {
		if got := e.ElementCRLRef(); got != XAdES111Element_CRL_REF {
			t.Errorf("ElementCRLRef() = %v, want XAdES111Element_CRL_REF", got)
		}
	})
	t.Run("ElementCRLRefs", func(t *testing.T) {
		if got := e.ElementCRLRefs(); got != XAdES111Element_CRL_REFS {
			t.Errorf("ElementCRLRefs() = %v, want XAdES111Element_CRL_REFS", got)
		}
	})
	t.Run("ElementCRLValues", func(t *testing.T) {
		if got := e.ElementCRLValues(); got != XAdES111Element_CRL_VALUES {
			t.Errorf("ElementCRLValues() = %v, want XAdES111Element_CRL_VALUES", got)
		}
	})
	t.Run("ElementDataObjectFormat", func(t *testing.T) {
		if got := e.ElementDataObjectFormat(); got != XAdES111Element_DATA_OBJECT_FORMAT {
			t.Errorf("ElementDataObjectFormat() = %v, want XAdES111Element_DATA_OBJECT_FORMAT", got)
		}
	})
	t.Run("ElementDescription", func(t *testing.T) {
		if got := e.ElementDescription(); got != XAdES111Element_DESCRIPTION {
			t.Errorf("ElementDescription() = %v, want XAdES111Element_DESCRIPTION", got)
		}
	})
	t.Run("ElementDigestAlgAndValue", func(t *testing.T) {
		if got := e.ElementDigestAlgAndValue(); got != XAdES111Element_DIGEST_ALG_AND_VALUE {
			t.Errorf("ElementDigestAlgAndValue() = %v, want XAdES111Element_DIGEST_ALG_AND_VALUE", got)
		}
	})
	t.Run("ElementDocumentationReference", func(t *testing.T) {
		if got := e.ElementDocumentationReference(); got != XAdES111Element_DOCUMENTATION_REFERENCE {
			t.Errorf("ElementDocumentationReference() = %v, want XAdES111Element_DOCUMENTATION_REFERENCE", got)
		}
	})
	t.Run("ElementDocumentationReferences", func(t *testing.T) {
		if got := e.ElementDocumentationReferences(); got != XAdES111Element_DOCUMENTATION_REFERENCES {
			t.Errorf("ElementDocumentationReferences() = %v, want XAdES111Element_DOCUMENTATION_REFERENCES", got)
		}
	})
	t.Run("ElementEncapsulatedCRLValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedCRLValue(); got != XAdES111Element_ENCAPSULATED_CRL_VALUE {
			t.Errorf("ElementEncapsulatedCRLValue() = %v, want XAdES111Element_ENCAPSULATED_CRL_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedOCSPValue", func(t *testing.T) {
		if got := e.ElementEncapsulatedOCSPValue(); got != XAdES111Element_ENCAPSULATED_OCSP_VALUE {
			t.Errorf("ElementEncapsulatedOCSPValue() = %v, want XAdES111Element_ENCAPSULATED_OCSP_VALUE", got)
		}
	})
	t.Run("ElementEncapsulatedPKIData", func(t *testing.T) {
		if got := e.ElementEncapsulatedPKIData(); got != XAdES111Element_ENCAPSULATED_PKI_DATA {
			t.Errorf("ElementEncapsulatedPKIData() = %v, want XAdES111Element_ENCAPSULATED_PKI_DATA", got)
		}
	})
	t.Run("ElementEncapsulatedTimeStamp", func(t *testing.T) {
		if got := e.ElementEncapsulatedTimeStamp(); got != XAdES111Element_ENCAPSULATED_TIMESTAMP {
			t.Errorf("ElementEncapsulatedTimeStamp() = %v, want XAdES111Element_ENCAPSULATED_TIMESTAMP", got)
		}
	})
	t.Run("ElementEncapsulatedX509Certificate", func(t *testing.T) {
		if got := e.ElementEncapsulatedX509Certificate(); got != XAdES111Element_ENCAPSULATED_X509_CERTIFICATE {
			t.Errorf("ElementEncapsulatedX509Certificate() = %v, want XAdES111Element_ENCAPSULATED_X509_CERTIFICATE", got)
		}
	})
	t.Run("ElementEncoding", func(t *testing.T) {
		if got := e.ElementEncoding(); got != XAdES111Element_ENCODING {
			t.Errorf("ElementEncoding() = %v, want XAdES111Element_ENCODING", got)
		}
	})
	t.Run("ElementExplicitText", func(t *testing.T) {
		if got := e.ElementExplicitText(); got != XAdES111Element_EXPLICIT_TEXT {
			t.Errorf("ElementExplicitText() = %v, want XAdES111Element_EXPLICIT_TEXT", got)
		}
	})
	t.Run("ElementIdentifier", func(t *testing.T) {
		if got := e.ElementIdentifier(); got != XAdES111Element_IDENTIFIER {
			t.Errorf("ElementIdentifier() = %v, want XAdES111Element_IDENTIFIER", got)
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
		if got := e.ElementIndividualDataObjectsTimeStamp(); got != XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP {
			t.Errorf("ElementIndividualDataObjectsTimeStamp() = %v, want XAdES111Element_INDIVIDUAL_DATA_OBJECTS_TIMESTAMP", got)
		}
	})
	t.Run("Elementint", func(t *testing.T) {
		if got := e.Elementint(); got != XAdES111Element_INT {
			t.Errorf("Elementint() = %v, want XAdES111Element_INT", got)
		}
	})
	t.Run("ElementIssueTime", func(t *testing.T) {
		if got := e.ElementIssueTime(); got != XAdES111Element_ISSUE_TIME {
			t.Errorf("ElementIssueTime() = %v, want XAdES111Element_ISSUE_TIME", got)
		}
	})
	t.Run("ElementIssuer", func(t *testing.T) {
		if got := e.ElementIssuer(); got != XAdES111Element_ISSUER {
			t.Errorf("ElementIssuer() = %v, want XAdES111Element_ISSUER", got)
		}
	})
	t.Run("ElementIssuerSerial", func(t *testing.T) {
		if got := e.ElementIssuerSerial(); got != XAdES111Element_ISSUER_SERIAL {
			t.Errorf("ElementIssuerSerial() = %v, want XAdES111Element_ISSUER_SERIAL", got)
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
		if got := e.ElementMimeType(); got != XAdES111Element_MIME_TYPE {
			t.Errorf("ElementMimeType() = %v, want XAdES111Element_MIME_TYPE", got)
		}
	})
	t.Run("ElementNoticeNumbers", func(t *testing.T) {
		if got := e.ElementNoticeNumbers(); got != XAdES111Element_NOTICE_NUMBERS {
			t.Errorf("ElementNoticeNumbers() = %v, want XAdES111Element_NOTICE_NUMBERS", got)
		}
	})
	t.Run("ElementNoticeRef", func(t *testing.T) {
		if got := e.ElementNoticeRef(); got != XAdES111Element_NOTICE_REF {
			t.Errorf("ElementNoticeRef() = %v, want XAdES111Element_NOTICE_REF", got)
		}
	})
	t.Run("ElementNumber", func(t *testing.T) {
		if got := e.ElementNumber(); got != XAdES111Element_NUMBER {
			t.Errorf("ElementNumber() = %v, want XAdES111Element_NUMBER", got)
		}
	})
	t.Run("ElementObjectIdentifier", func(t *testing.T) {
		if got := e.ElementObjectIdentifier(); got != XAdES111Element_OBJECT_IDENTIFIER {
			t.Errorf("ElementObjectIdentifier() = %v, want XAdES111Element_OBJECT_IDENTIFIER", got)
		}
	})
	t.Run("ElementObjectReference", func(t *testing.T) {
		if got := e.ElementObjectReference(); got != XAdES111Element_OBJECT_REFERENCE {
			t.Errorf("ElementObjectReference() = %v, want XAdES111Element_OBJECT_REFERENCE", got)
		}
	})
	t.Run("ElementOCSPIdentifier", func(t *testing.T) {
		if got := e.ElementOCSPIdentifier(); got != XAdES111Element_OCSP_IDENTIFIER {
			t.Errorf("ElementOCSPIdentifier() = %v, want XAdES111Element_OCSP_IDENTIFIER", got)
		}
	})
	t.Run("ElementOCSPRef", func(t *testing.T) {
		if got := e.ElementOCSPRef(); got != XAdES111Element_OCSP_REF {
			t.Errorf("ElementOCSPRef() = %v, want XAdES111Element_OCSP_REF", got)
		}
	})
	t.Run("ElementOCSPRefs", func(t *testing.T) {
		if got := e.ElementOCSPRefs(); got != XAdES111Element_OCSP_REFS {
			t.Errorf("ElementOCSPRefs() = %v, want XAdES111Element_OCSP_REFS", got)
		}
	})
	t.Run("ElementOCSPValues", func(t *testing.T) {
		if got := e.ElementOCSPValues(); got != XAdES111Element_OCSP_VALUES {
			t.Errorf("ElementOCSPValues() = %v, want XAdES111Element_OCSP_VALUES", got)
		}
	})
	t.Run("ElementOrganization", func(t *testing.T) {
		if got := e.ElementOrganization(); got != XAdES111Element_ORGANIZATION {
			t.Errorf("ElementOrganization() = %v, want XAdES111Element_ORGANIZATION", got)
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
		if got := e.ElementOtherCertificate(); got != XAdES111Element_OTHER_CERTIFICATE {
			t.Errorf("ElementOtherCertificate() = %v, want XAdES111Element_OTHER_CERTIFICATE", got)
		}
	})
	t.Run("ElementOtherRef", func(t *testing.T) {
		if got := e.ElementOtherRef(); got != XAdES111Element_OTHER_REF {
			t.Errorf("ElementOtherRef() = %v, want XAdES111Element_OTHER_REF", got)
		}
	})
	t.Run("ElementOtherRefs", func(t *testing.T) {
		if got := e.ElementOtherRefs(); got != XAdES111Element_OTHER_REFS {
			t.Errorf("ElementOtherRefs() = %v, want XAdES111Element_OTHER_REFS", got)
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
		if got := e.ElementOtherValue(); got != XAdES111Element_OTHER_VALUE {
			t.Errorf("ElementOtherValue() = %v, want XAdES111Element_OTHER_VALUE", got)
		}
	})
	t.Run("ElementOtherValues", func(t *testing.T) {
		if got := e.ElementOtherValues(); got != XAdES111Element_OTHER_VALUES {
			t.Errorf("ElementOtherValues() = %v, want XAdES111Element_OTHER_VALUES", got)
		}
	})
	t.Run("ElementPostalCode", func(t *testing.T) {
		if got := e.ElementPostalCode(); got != XAdES111Element_POSTAL_CODE {
			t.Errorf("ElementPostalCode() = %v, want XAdES111Element_POSTAL_CODE", got)
		}
	})
	t.Run("ElementProducedAt", func(t *testing.T) {
		if got := e.ElementProducedAt(); got != XAdES111Element_PRODUCED_AT {
			t.Errorf("ElementProducedAt() = %v, want XAdES111Element_PRODUCED_AT", got)
		}
	})
	t.Run("ElementQualifyingProperties", func(t *testing.T) {
		if got := e.ElementQualifyingProperties(); got != XAdES111Element_QUALIFYING_PROPERTIES {
			t.Errorf("ElementQualifyingProperties() = %v, want XAdES111Element_QUALIFYING_PROPERTIES", got)
		}
	})
	t.Run("ElementQualifyingPropertiesReference", func(t *testing.T) {
		if got := e.ElementQualifyingPropertiesReference(); got != XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE {
			t.Errorf("ElementQualifyingPropertiesReference() = %v, want XAdES111Element_QUALIFYING_PROPERTIES_REFERENCE", got)
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
		if got := e.ElementRefsOnlyTimeStamp(); got != XAdES111Element_REFS_ONLY_TIMESTAMP {
			t.Errorf("ElementRefsOnlyTimeStamp() = %v, want XAdES111Element_REFS_ONLY_TIMESTAMP", got)
		}
	})
	t.Run("ElementResponderID", func(t *testing.T) {
		if got := e.ElementResponderID(); got != XAdES111Element_RESPONDER_ID {
			t.Errorf("ElementResponderID() = %v, want XAdES111Element_RESPONDER_ID", got)
		}
	})
	t.Run("ElementRevocationValues", func(t *testing.T) {
		if got := e.ElementRevocationValues(); got != XAdES111Element_REVOCATION_VALUES {
			t.Errorf("ElementRevocationValues() = %v, want XAdES111Element_REVOCATION_VALUES", got)
		}
	})
	t.Run("ElementSigAndRefsTimeStamp", func(t *testing.T) {
		if got := e.ElementSigAndRefsTimeStamp(); got != XAdES111Element_SIG_AND_REFS_TIMESTAMP {
			t.Errorf("ElementSigAndRefsTimeStamp() = %v, want XAdES111Element_SIG_AND_REFS_TIMESTAMP", got)
		}
	})
	t.Run("ElementSigPolicyHash", func(t *testing.T) {
		if got := e.ElementSigPolicyHash(); got != XAdES111Element_SIG_POLICY_HASH {
			t.Errorf("ElementSigPolicyHash() = %v, want XAdES111Element_SIG_POLICY_HASH", got)
		}
	})
	t.Run("ElementSigPolicyId", func(t *testing.T) {
		if got := e.ElementSigPolicyId(); got != XAdES111Element_SIG_POLICY_ID {
			t.Errorf("ElementSigPolicyId() = %v, want XAdES111Element_SIG_POLICY_ID", got)
		}
	})
	t.Run("ElementSigPolicyQualifier", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifier(); got != XAdES111Element_SIG_POLICY_QUALIFIER {
			t.Errorf("ElementSigPolicyQualifier() = %v, want XAdES111Element_SIG_POLICY_QUALIFIER", got)
		}
	})
	t.Run("ElementSigPolicyQualifiers", func(t *testing.T) {
		if got := e.ElementSigPolicyQualifiers(); got != XAdES111Element_SIG_POLICY_QUALIFIERS {
			t.Errorf("ElementSigPolicyQualifiers() = %v, want XAdES111Element_SIG_POLICY_QUALIFIERS", got)
		}
	})
	t.Run("ElementSignaturePolicyId", func(t *testing.T) {
		if got := e.ElementSignaturePolicyId(); got != XAdES111Element_SIGNATURE_POLICY_ID {
			t.Errorf("ElementSignaturePolicyId() = %v, want XAdES111Element_SIGNATURE_POLICY_ID", got)
		}
	})
	t.Run("ElementSignaturePolicyIdentifier", func(t *testing.T) {
		if got := e.ElementSignaturePolicyIdentifier(); got != XAdES111Element_SIGNATURE_POLICY_IDENTIFIER {
			t.Errorf("ElementSignaturePolicyIdentifier() = %v, want XAdES111Element_SIGNATURE_POLICY_IDENTIFIER", got)
		}
	})
	t.Run("ElementSignaturePolicyImplied", func(t *testing.T) {
		if got := e.ElementSignaturePolicyImplied(); got != XAdES111Element_SIGNATURE_POLICY_IMPLIED {
			t.Errorf("ElementSignaturePolicyImplied() = %v, want XAdES111Element_SIGNATURE_POLICY_IMPLIED", got)
		}
	})
	t.Run("ElementSignatureProductionPlace", func(t *testing.T) {
		if got := e.ElementSignatureProductionPlace(); got != XAdES111Element_SIGNATURE_PRODUCTION_PLACE {
			t.Errorf("ElementSignatureProductionPlace() = %v, want XAdES111Element_SIGNATURE_PRODUCTION_PLACE", got)
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
		if got := e.ElementSignatureTimeStamp(); got != XAdES111Element_SIGNATURE_TIMESTAMP {
			t.Errorf("ElementSignatureTimeStamp() = %v, want XAdES111Element_SIGNATURE_TIMESTAMP", got)
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
		if got := e.ElementSignedDataObjectProperties(); got != XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementSignedDataObjectProperties() = %v, want XAdES111Element_SIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedProperties", func(t *testing.T) {
		if got := e.ElementSignedProperties(); got != XAdES111Element_SIGNED_PROPERTIES {
			t.Errorf("ElementSignedProperties() = %v, want XAdES111Element_SIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementSignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementSignedSignatureProperties(); got != XAdES111Element_SIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementSignedSignatureProperties() = %v, want XAdES111Element_SIGNED_SIGNATURE_PROPERTIES", got)
		}
	})
	t.Run("ElementSignerRole", func(t *testing.T) {
		if got := e.ElementSignerRole(); got != XAdES111Element_SIGNER_ROLE {
			t.Errorf("ElementSignerRole() = %v, want XAdES111Element_SIGNER_ROLE", got)
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
		if got := e.ElementSigningCertificate(); got != XAdES111Element_SIGNING_CERTIFICATE {
			t.Errorf("ElementSigningCertificate() = %v, want XAdES111Element_SIGNING_CERTIFICATE", got)
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
		if got := e.ElementSigningTime(); got != XAdES111Element_SIGNING_TIME {
			t.Errorf("ElementSigningTime() = %v, want XAdES111Element_SIGNING_TIME", got)
		}
	})
	t.Run("ElementSPURI", func(t *testing.T) {
		if got := e.ElementSPURI(); got != XAdES111Element_SP_URI {
			t.Errorf("ElementSPURI() = %v, want XAdES111Element_SP_URI", got)
		}
	})
	t.Run("ElementSPUserNotice", func(t *testing.T) {
		if got := e.ElementSPUserNotice(); got != XAdES111Element_SP_USER_NOTICE {
			t.Errorf("ElementSPUserNotice() = %v, want XAdES111Element_SP_USER_NOTICE", got)
		}
	})
	t.Run("ElementStateOrProvince", func(t *testing.T) {
		if got := e.ElementStateOrProvince(); got != XAdES111Element_STATE_OR_PROVINCE {
			t.Errorf("ElementStateOrProvince() = %v, want XAdES111Element_STATE_OR_PROVINCE", got)
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
		if got := e.ElementUnsignedDataObjectProperties(); got != XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES {
			t.Errorf("ElementUnsignedDataObjectProperties() = %v, want XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedDataObjectProperty", func(t *testing.T) {
		if got := e.ElementUnsignedDataObjectProperty(); got != XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY {
			t.Errorf("ElementUnsignedDataObjectProperty() = %v, want XAdES111Element_UNSIGNED_DATA_OBJECT_PROPERTY", got)
		}
	})
	t.Run("ElementUnsignedProperties", func(t *testing.T) {
		if got := e.ElementUnsignedProperties(); got != XAdES111Element_UNSIGNED_PROPERTIES {
			t.Errorf("ElementUnsignedProperties() = %v, want XAdES111Element_UNSIGNED_PROPERTIES", got)
		}
	})
	t.Run("ElementUnsignedSignatureProperties", func(t *testing.T) {
		if got := e.ElementUnsignedSignatureProperties(); got != XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES {
			t.Errorf("ElementUnsignedSignatureProperties() = %v, want XAdES111Element_UNSIGNED_SIGNATURE_PROPERTIES", got)
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
		if got := e.ElementXMLTimeStamp(); got != XAdES111Element_XML_TIMESTAMP {
			t.Errorf("ElementXMLTimeStamp() = %v, want XAdES111Element_XML_TIMESTAMP", got)
		}
	})
}
