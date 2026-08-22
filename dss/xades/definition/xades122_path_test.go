// KAT test for xades122_path.go: every XAdES122Path XPathQuery-returning method below is
// checked against the query string independently re-derived from the same element
// chain as the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.xades122.XAdES122Path; nil-valued upstream overrides
// (schema versions that do not define a given path) are asserted nil.
package definition

import "testing"

func TestXAdES122Path_QueryStringKAT(t *testing.T) {
	p := NewXAdES122Path()

	t.Run("QualifyingPropertiesPath", func(t *testing.T) {
		got := p.QualifyingPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties")
		}
	})
	t.Run("SignedPropertiesPath", func(t *testing.T) {
		got := p.SignedPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties")
		}
	})
	t.Run("SignedSignaturePropertiesPath", func(t *testing.T) {
		got := p.SignedSignaturePropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties")
		}
	})
	t.Run("SigningTimePath", func(t *testing.T) {
		got := p.SigningTimePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningTime" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningTime")
		}
	})
	t.Run("SigningCertificatePath", func(t *testing.T) {
		got := p.SigningCertificatePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningCertificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningCertificate")
		}
	})
	t.Run("SigningCertificateChildren", func(t *testing.T) {
		got := p.SigningCertificateChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningCertificate/xades122:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SigningCertificate/xades122:Cert")
		}
	})
	t.Run("SigningCertificateV2Path", func(t *testing.T) {
		got := p.SigningCertificateV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SigningCertificateV2Children", func(t *testing.T) {
		got := p.SigningCertificateV2Children()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignatureProductionPlacePath", func(t *testing.T) {
		got := p.SignatureProductionPlacePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignatureProductionPlace" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignatureProductionPlace")
		}
	})
	t.Run("SignatureProductionPlaceV2Path", func(t *testing.T) {
		got := p.SignatureProductionPlaceV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignaturePolicyIdentifierPath", func(t *testing.T) {
		got := p.SignaturePolicyIdentifierPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignaturePolicyIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignaturePolicyIdentifier")
		}
	})
	t.Run("SignerRolePath", func(t *testing.T) {
		got := p.SignerRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole")
		}
	})
	t.Run("ClaimedRolePath", func(t *testing.T) {
		got := p.ClaimedRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole/xades122:ClaimedRoles/xades122:ClaimedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole/xades122:ClaimedRoles/xades122:ClaimedRole")
		}
	})
	t.Run("SignedAssertionPath", func(t *testing.T) {
		got := p.SignedAssertionPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignerRoleV2Path", func(t *testing.T) {
		got := p.SignerRoleV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("ClaimedRoleV2Path", func(t *testing.T) {
		got := p.ClaimedRoleV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CertifiedRolePath", func(t *testing.T) {
		got := p.CertifiedRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole/xades122:CertifiedRoles/xades122:CertifiedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedSignatureProperties/xades122:SignerRole/xades122:CertifiedRoles/xades122:CertifiedRole")
		}
	})
	t.Run("CertifiedRoleV2Path", func(t *testing.T) {
		got := p.CertifiedRoleV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignedDataObjectPropertiesPath", func(t *testing.T) {
		got := p.SignedDataObjectPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties")
		}
	})
	t.Run("AllDataObjectsTimestampPath", func(t *testing.T) {
		got := p.AllDataObjectsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:AllDataObjectsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:AllDataObjectsTimeStamp")
		}
	})
	t.Run("IndividualDataObjectsTimestampPath", func(t *testing.T) {
		got := p.IndividualDataObjectsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:IndividualDataObjectsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:IndividualDataObjectsTimeStamp")
		}
	})
	t.Run("DataObjectFormat", func(t *testing.T) {
		got := p.DataObjectFormat()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat")
		}
	})
	t.Run("DataObjectFormatMimeType", func(t *testing.T) {
		got := p.DataObjectFormatMimeType()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat/xades122:MimeType" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat/xades122:MimeType")
		}
	})
	t.Run("DataObjectFormatObjectIdentifier", func(t *testing.T) {
		got := p.DataObjectFormatObjectIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat/xades122:ObjectIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:DataObjectFormat/xades122:ObjectIdentifier")
		}
	})
	t.Run("CommitmentTypeIndicationPath", func(t *testing.T) {
		got := p.CommitmentTypeIndicationPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:CommitmentTypeIndication" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:SignedProperties/xades122:SignedDataObjectProperties/xades122:CommitmentTypeIndication")
		}
	})
	t.Run("UnsignedPropertiesPath", func(t *testing.T) {
		got := p.UnsignedPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties")
		}
	})
	t.Run("UnsignedSignaturePropertiesPath", func(t *testing.T) {
		got := p.UnsignedSignaturePropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties")
		}
	})
	t.Run("CounterSignaturePath", func(t *testing.T) {
		got := p.CounterSignaturePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CounterSignature" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CounterSignature")
		}
	})
	t.Run("AttributeRevocationRefsPath", func(t *testing.T) {
		got := p.AttributeRevocationRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeRevocationRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeRevocationRefs")
		}
	})
	t.Run("CompleteRevocationRefsPath", func(t *testing.T) {
		got := p.CompleteRevocationRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteRevocationRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteRevocationRefs")
		}
	})
	t.Run("CompleteCertificateRefsPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteCertificateRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteCertificateRefs")
		}
	})
	t.Run("CompleteCertificateRefsCertPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsCertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteCertificateRefs/xades122:CertRefs/xades122:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CompleteCertificateRefs/xades122:CertRefs/xades122:Cert")
		}
	})
	t.Run("CompleteCertificateRefsV2Path", func(t *testing.T) {
		got := p.CompleteCertificateRefsV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CompleteCertificateRefsV2CertPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsV2CertPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("AttributeCertificateRefsPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeCertificateRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeCertificateRefs")
		}
	})
	t.Run("AttributeCertificateRefsCertPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsCertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeCertificateRefs/xades122:CertRefs/xades122:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:AttributeCertificateRefs/xades122:CertRefs/xades122:Cert")
		}
	})
	t.Run("AttributeCertificateRefsV2Path", func(t *testing.T) {
		got := p.AttributeCertificateRefsV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("AttributeCertificateRefsV2CertPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsV2CertPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CertificateValuesPath", func(t *testing.T) {
		got := p.CertificateValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CertificateValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CertificateValues")
		}
	})
	t.Run("RevocationValuesPath", func(t *testing.T) {
		got := p.RevocationValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:RevocationValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:RevocationValues")
		}
	})
	t.Run("AttributeRevocationValuesPath", func(t *testing.T) {
		got := p.AttributeRevocationValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("EncapsulatedCertificateValuesPath", func(t *testing.T) {
		got := p.EncapsulatedCertificateValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CertificateValues/xades122:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:CertificateValues/xades122:EncapsulatedX509Certificate")
		}
	})
	t.Run("AttrAuthoritiesCertValuesPath", func(t *testing.T) {
		got := p.AttrAuthoritiesCertValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("EncapsulatedAttrAuthoritiesCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedAttrAuthoritiesCertValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("EncapsulatedTimeStampValidationDataCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedTimeStampValidationDataCertValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("TimeStampValidationDataRevocationValuesPath", func(t *testing.T) {
		got := p.TimeStampValidationDataRevocationValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("AnyValidationDataPath", func(t *testing.T) {
		got := p.AnyValidationDataPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("EncapsulatedAnyValidationDataCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedAnyValidationDataCertValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("AnyValidationDataRevocationValuesPath", func(t *testing.T) {
		got := p.AnyValidationDataRevocationValuesPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignatureTimestampPath", func(t *testing.T) {
		got := p.SignatureTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:SignatureTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:SignatureTimeStamp")
		}
	})
	t.Run("SigAndRefsTimestampPath", func(t *testing.T) {
		got := p.SigAndRefsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:SigAndRefsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades122:QualifyingProperties/xades122:UnsignedProperties/xades122:UnsignedSignatureProperties/xades122:SigAndRefsTimeStamp")
		}
	})
	t.Run("SigAndRefsTimestampV2Path", func(t *testing.T) {
		got := p.SigAndRefsTimestampV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("RefsOnlyTimestampPath", func(t *testing.T) {
		got := p.RefsOnlyTimestampPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("RefsOnlyTimestampV2Path", func(t *testing.T) {
		got := p.RefsOnlyTimestampV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("ArchiveTimestampPath", func(t *testing.T) {
		got := p.ArchiveTimestampPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("TimestampValidationDataPath", func(t *testing.T) {
		got := p.TimestampValidationDataPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SignaturePolicyStorePath", func(t *testing.T) {
		got := p.SignaturePolicyStorePath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("SealingEvidenceRecordsPath", func(t *testing.T) {
		got := p.SealingEvidenceRecordsPath()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentCRLValuesChildren", func(t *testing.T) {
		got := p.CurrentCRLValuesChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLValues/xades122:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLValues/xades122:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentCRLRefsChildren", func(t *testing.T) {
		got := p.CurrentCRLRefsChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLRefs/xades122:CRLRef" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLRefs/xades122:CRLRef")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifier", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLIdentifier")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierIssuer", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierIssuer()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLIdentifier/xades122:Issuer" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLIdentifier/xades122:Issuer")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierIssueTime", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierIssueTime()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLIdentifier/xades122:IssueTime" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLIdentifier/xades122:IssueTime")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierNumber", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierNumber()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLIdentifier/xades122:Number" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLIdentifier/xades122:Number")
		}
	})
	t.Run("CurrentOCSPValuesChildren", func(t *testing.T) {
		got := p.CurrentOCSPValuesChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:OCSPValues/xades122:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:OCSPValues/xades122:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentOCSPRefsChildren", func(t *testing.T) {
		got := p.CurrentOCSPRefsChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:OCSPRefs/xades122:OCSPRef" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:OCSPRefs/xades122:OCSPRef")
		}
	})
	t.Run("CurrentOCSPRefResponderID", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderID()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:OCSPIdentifier/xades122:ResponderID" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:OCSPIdentifier/xades122:ResponderID")
		}
	})
	t.Run("CurrentOCSPRefResponderIDByName", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderIDByName()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentOCSPRefResponderIDByKey", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderIDByKey()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentOCSPRefProducedAt", func(t *testing.T) {
		got := p.CurrentOCSPRefProducedAt()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:OCSPIdentifier/xades122:ProducedAt" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:OCSPIdentifier/xades122:ProducedAt")
		}
	})
	t.Run("CurrentDigestAlgAndValue", func(t *testing.T) {
		got := p.CurrentDigestAlgAndValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:DigestAlgAndValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:DigestAlgAndValue")
		}
	})
	t.Run("CurrentCertRefsCertChildren", func(t *testing.T) {
		got := p.CurrentCertRefsCertChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CertRefs/xades122:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CertRefs/xades122:Cert")
		}
	})
	t.Run("CurrentCertRefs141CertChildren", func(t *testing.T) {
		got := p.CurrentCertRefs141CertChildren()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentCertChildren", func(t *testing.T) {
		got := p.CurrentCertChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:Cert")
		}
	})
	t.Run("CurrentCertDigest", func(t *testing.T) {
		got := p.CurrentCertDigest()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CertDigest" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CertDigest")
		}
	})
	t.Run("CurrentEncapsulatedTimestamp", func(t *testing.T) {
		got := p.CurrentEncapsulatedTimestamp()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:EncapsulatedTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:EncapsulatedTimeStamp")
		}
	})
	t.Run("CurrentEncapsulatedCertificate", func(t *testing.T) {
		got := p.CurrentEncapsulatedCertificate()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:EncapsulatedX509Certificate")
		}
	})
	t.Run("CurrentCertificateValuesEncapsulatedCertificate", func(t *testing.T) {
		got := p.CurrentCertificateValuesEncapsulatedCertificate()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CertificateValues/xades122:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CertificateValues/xades122:EncapsulatedX509Certificate")
		}
	})
	t.Run("CurrentRevocationValuesEncapsulatedOCSPValue", func(t *testing.T) {
		got := p.CurrentRevocationValuesEncapsulatedOCSPValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:RevocationValues/xades122:OCSPValues/xades122:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:RevocationValues/xades122:OCSPValues/xades122:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentEncapsulatedOCSPValue", func(t *testing.T) {
		got := p.CurrentEncapsulatedOCSPValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:OCSPValues/xades122:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:OCSPValues/xades122:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentRevocationValuesEncapsulatedCRLValue", func(t *testing.T) {
		got := p.CurrentRevocationValuesEncapsulatedCRLValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:RevocationValues/xades122:CRLValues/xades122:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:RevocationValues/xades122:CRLValues/xades122:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentEncapsulatedCRLValue", func(t *testing.T) {
		got := p.CurrentEncapsulatedCRLValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CRLValues/xades122:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CRLValues/xades122:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentIssuerSerialIssuerNamePath", func(t *testing.T) {
		got := p.CurrentIssuerSerialIssuerNamePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:IssuerSerial/ds:X509IssuerName" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:IssuerSerial/ds:X509IssuerName")
		}
	})
	t.Run("CurrentIssuerSerialSerialNumberPath", func(t *testing.T) {
		got := p.CurrentIssuerSerialSerialNumberPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:IssuerSerial/ds:X509SerialNumber" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:IssuerSerial/ds:X509SerialNumber")
		}
	})
	t.Run("CurrentIssuerSerialV2Path", func(t *testing.T) {
		got := p.CurrentIssuerSerialV2Path()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentCommitmentIdentifierPath", func(t *testing.T) {
		got := p.CurrentCommitmentIdentifierPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CommitmentTypeId/xades122:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CommitmentTypeId/xades122:Identifier")
		}
	})
	t.Run("CurrentCommitmentDescriptionPath", func(t *testing.T) {
		got := p.CurrentCommitmentDescriptionPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CommitmentTypeId/xades122:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CommitmentTypeId/xades122:Description")
		}
	})
	t.Run("CurrentCommitmentDocumentationReferencesPath", func(t *testing.T) {
		got := p.CurrentCommitmentDocumentationReferencesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:CommitmentTypeId/xades122:DocumentationReferences" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:CommitmentTypeId/xades122:DocumentationReferences")
		}
	})
	t.Run("CurrentDocumentationReference", func(t *testing.T) {
		got := p.CurrentDocumentationReference()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:DocumentationReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:DocumentationReference")
		}
	})
	t.Run("CurrentDescription", func(t *testing.T) {
		got := p.CurrentDescription()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:Description")
		}
	})
	t.Run("CurrentObjectIdentifier", func(t *testing.T) {
		got := p.CurrentObjectIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:ObjectIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:ObjectIdentifier")
		}
	})
	t.Run("CurrentCommitmentObjectReferencesPath", func(t *testing.T) {
		got := p.CurrentCommitmentObjectReferencesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:ObjectReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:ObjectReference")
		}
	})
	t.Run("CurrentCommitmentAllSignedDataObjectsPath", func(t *testing.T) {
		got := p.CurrentCommitmentAllSignedDataObjectsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:AllSignedDataObjects" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:AllSignedDataObjects")
		}
	})
	t.Run("CurrentMimeType", func(t *testing.T) {
		got := p.CurrentMimeType()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:MimeType" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:MimeType")
		}
	})
	t.Run("CurrentEncoding", func(t *testing.T) {
		got := p.CurrentEncoding()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:Encoding" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:Encoding")
		}
	})
	t.Run("CurrentSignaturePolicyId", func(t *testing.T) {
		got := p.CurrentSignaturePolicyId()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:Identifier")
		}
	})
	t.Run("CurrentSignaturePolicyDigestAlgAndValue", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDigestAlgAndValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyHash" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyHash")
		}
	})
	t.Run("CurrentSignaturePolicySPURI", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPURI()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers/xades122:SigPolicyQualifier/xades122:SPURI" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers/xades122:SigPolicyQualifier/xades122:SPURI")
		}
	})
	t.Run("CurrentSignaturePolicySPUserNotice", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPUserNotice()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers/xades122:SigPolicyQualifier/xades122:SPUserNotice" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers/xades122:SigPolicyQualifier/xades122:SPUserNotice")
		}
	})
	t.Run("CurrentSPUserNoticeNoticeRefOrganization", func(t *testing.T) {
		got := p.CurrentSPUserNoticeNoticeRefOrganization()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:NoticeRef/xades122:Organization" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:NoticeRef/xades122:Organization")
		}
	})
	t.Run("CurrentSPUserNoticeNoticeRefNoticeNumbers", func(t *testing.T) {
		got := p.CurrentSPUserNoticeNoticeRefNoticeNumbers()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:NoticeRef/xades122:NoticeNumbers" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:NoticeRef/xades122:NoticeNumbers")
		}
	})
	t.Run("CurrentSPUserNoticeExplicitText", func(t *testing.T) {
		got := p.CurrentSPUserNoticeExplicitText()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:ExplicitText" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:ExplicitText")
		}
	})
	t.Run("CurrentSignaturePolicySPDocSpecification", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPDocSpecification()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSignaturePolicySPDocSpecificationIdentifier", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPDocSpecificationIdentifier()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSignaturePolicyDescription", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDescription()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:Description")
		}
	})
	t.Run("CurrentSignaturePolicyDocumentationReferences", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDocumentationReferences()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:DocumentationReferences" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyId/xades122:DocumentationReferences")
		}
	})
	t.Run("CurrentSignaturePolicyImplied", func(t *testing.T) {
		got := p.CurrentSignaturePolicyImplied()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyImplied" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyImplied")
		}
	})
	t.Run("CurrentSignaturePolicyTransforms", func(t *testing.T) {
		got := p.CurrentSignaturePolicyTransforms()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/ds:Transforms" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/ds:Transforms")
		}
	})
	t.Run("CurrentSignaturePolicyQualifiers", func(t *testing.T) {
		got := p.CurrentSignaturePolicyQualifiers()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:SignaturePolicyId/xades122:SigPolicyQualifiers")
		}
	})
	t.Run("CurrentInclude", func(t *testing.T) {
		got := p.CurrentInclude()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:Include" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:Include")
		}
	})
	t.Run("CurrentQualifyingPropertiesPath", func(t *testing.T) {
		got := p.CurrentQualifyingPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:QualifyingProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:QualifyingProperties")
		}
	})
	t.Run("CurrentSPDocSpecification", func(t *testing.T) {
		got := p.CurrentSPDocSpecification()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentIdentifier", func(t *testing.T) {
		got := p.CurrentIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades122:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades122:Identifier")
		}
	})
	t.Run("CurrentSPDocSpecificationIdentifier", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationIdentifier()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSPDocSpecificationDescription", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationDescription()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentDocumentationReferenceElements", func(t *testing.T) {
		got := p.CurrentDocumentationReferenceElements()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSPDocSpecificationDocumentationReferenceElements", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationDocumentationReferenceElements()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSignaturePolicyDocument", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDocument()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
	t.Run("CurrentSigPolDocLocalURI", func(t *testing.T) {
		got := p.CurrentSigPolDocLocalURI()
		if got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})

	if got := p.Namespace().Uri(); got != XAdESNamespaceXAdES122.Uri() {
		t.Errorf("Namespace().Uri() = %q, want %q", got, XAdESNamespaceXAdES122.Uri())
	}
	if got := p.SignedPropertiesUri(); got != "http://uri.etsi.org/01903/v1.2.2#SignedProperties" {
		t.Errorf("SignedPropertiesUri() = %q, want %q", got, "http://uri.etsi.org/01903/v1.2.2#SignedProperties")
	}
	if got := p.CounterSignatureUri(); got != "http://uri.etsi.org/01903#CountersignedSignature" {
		t.Errorf("CounterSignatureUri() = %q", got)
	}
}
