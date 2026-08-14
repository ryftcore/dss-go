// KAT test for xades132_path.go: every XAdES132Path XPathQuery-returning method below is
// checked against the query string independently re-derived from the same element
// chain as the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.xades132.XAdES132Path; nil-valued upstream overrides
// (schema versions that do not define a given path) are asserted nil.
package definition

import "testing"

func TestXAdES132Path_QueryStringKAT(t *testing.T) {
	p := NewXAdES132Path()

	t.Run("QualifyingPropertiesPath", func(t *testing.T) {
		got := p.QualifyingPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties")
		}
	})
	t.Run("SignedPropertiesPath", func(t *testing.T) {
		got := p.SignedPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties")
		}
	})
	t.Run("SignedSignaturePropertiesPath", func(t *testing.T) {
		got := p.SignedSignaturePropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties")
		}
	})
	t.Run("SigningTimePath", func(t *testing.T) {
		got := p.SigningTimePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningTime" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningTime")
		}
	})
	t.Run("SigningCertificatePath", func(t *testing.T) {
		got := p.SigningCertificatePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificate")
		}
	})
	t.Run("SigningCertificateChildren", func(t *testing.T) {
		got := p.SigningCertificateChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificate/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificate/xades132:Cert")
		}
	})
	t.Run("SigningCertificateV2Path", func(t *testing.T) {
		got := p.SigningCertificateV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificateV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificateV2")
		}
	})
	t.Run("SigningCertificateV2Children", func(t *testing.T) {
		got := p.SigningCertificateV2Children()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificateV2/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SigningCertificateV2/xades132:Cert")
		}
	})
	t.Run("SignatureProductionPlacePath", func(t *testing.T) {
		got := p.SignatureProductionPlacePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignatureProductionPlace" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignatureProductionPlace")
		}
	})
	t.Run("SignatureProductionPlaceV2Path", func(t *testing.T) {
		got := p.SignatureProductionPlaceV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignatureProductionPlaceV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignatureProductionPlaceV2")
		}
	})
	t.Run("SignaturePolicyIdentifierPath", func(t *testing.T) {
		got := p.SignaturePolicyIdentifierPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignaturePolicyIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignaturePolicyIdentifier")
		}
	})
	t.Run("SignerRolePath", func(t *testing.T) {
		got := p.SignerRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole")
		}
	})
	t.Run("ClaimedRolePath", func(t *testing.T) {
		got := p.ClaimedRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole/xades132:ClaimedRoles/xades132:ClaimedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole/xades132:ClaimedRoles/xades132:ClaimedRole")
		}
	})
	t.Run("SignedAssertionPath", func(t *testing.T) {
		got := p.SignedAssertionPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:SignedAssertions/xades132:SignedAssertion" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:SignedAssertions/xades132:SignedAssertion")
		}
	})
	t.Run("SignerRoleV2Path", func(t *testing.T) {
		got := p.SignerRoleV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2")
		}
	})
	t.Run("ClaimedRoleV2Path", func(t *testing.T) {
		got := p.ClaimedRoleV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:ClaimedRoles/xades132:ClaimedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:ClaimedRoles/xades132:ClaimedRole")
		}
	})
	t.Run("CertifiedRolePath", func(t *testing.T) {
		got := p.CertifiedRolePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole/xades132:CertifiedRoles/xades132:CertifiedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRole/xades132:CertifiedRoles/xades132:CertifiedRole")
		}
	})
	t.Run("CertifiedRoleV2Path", func(t *testing.T) {
		got := p.CertifiedRoleV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:CertifiedRolesV2/xades132:CertifiedRole" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedSignatureProperties/xades132:SignerRoleV2/xades132:CertifiedRolesV2/xades132:CertifiedRole")
		}
	})
	t.Run("SignedDataObjectPropertiesPath", func(t *testing.T) {
		got := p.SignedDataObjectPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties")
		}
	})
	t.Run("AllDataObjectsTimestampPath", func(t *testing.T) {
		got := p.AllDataObjectsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:AllDataObjectsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:AllDataObjectsTimeStamp")
		}
	})
	t.Run("IndividualDataObjectsTimestampPath", func(t *testing.T) {
		got := p.IndividualDataObjectsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:IndividualDataObjectsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:IndividualDataObjectsTimeStamp")
		}
	})
	t.Run("DataObjectFormat", func(t *testing.T) {
		got := p.DataObjectFormat()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat")
		}
	})
	t.Run("DataObjectFormatMimeType", func(t *testing.T) {
		got := p.DataObjectFormatMimeType()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat/xades132:MimeType" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat/xades132:MimeType")
		}
	})
	t.Run("DataObjectFormatObjectIdentifier", func(t *testing.T) {
		got := p.DataObjectFormatObjectIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat/xades132:ObjectIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:DataObjectFormat/xades132:ObjectIdentifier")
		}
	})
	t.Run("CommitmentTypeIndicationPath", func(t *testing.T) {
		got := p.CommitmentTypeIndicationPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:CommitmentTypeIndication" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:SignedProperties/xades132:SignedDataObjectProperties/xades132:CommitmentTypeIndication")
		}
	})
	t.Run("UnsignedPropertiesPath", func(t *testing.T) {
		got := p.UnsignedPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties")
		}
	})
	t.Run("UnsignedSignaturePropertiesPath", func(t *testing.T) {
		got := p.UnsignedSignaturePropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties")
		}
	})
	t.Run("CounterSignaturePath", func(t *testing.T) {
		got := p.CounterSignaturePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CounterSignature" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CounterSignature")
		}
	})
	t.Run("AttributeRevocationRefsPath", func(t *testing.T) {
		got := p.AttributeRevocationRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeRevocationRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeRevocationRefs")
		}
	})
	t.Run("CompleteRevocationRefsPath", func(t *testing.T) {
		got := p.CompleteRevocationRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteRevocationRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteRevocationRefs")
		}
	})
	t.Run("CompleteCertificateRefsPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteCertificateRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteCertificateRefs")
		}
	})
	t.Run("CompleteCertificateRefsCertPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsCertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteCertificateRefs/xades132:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CompleteCertificateRefs/xades132:CertRefs/xades132:Cert")
		}
	})
	t.Run("CompleteCertificateRefsV2Path", func(t *testing.T) {
		got := p.CompleteCertificateRefsV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2")
		}
	})
	t.Run("CompleteCertificateRefsV2CertPath", func(t *testing.T) {
		got := p.CompleteCertificateRefsV2CertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2/xades141:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:CompleteCertificateRefsV2/xades141:CertRefs/xades132:Cert")
		}
	})
	t.Run("AttributeCertificateRefsPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeCertificateRefs" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeCertificateRefs")
		}
	})
	t.Run("AttributeCertificateRefsCertPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsCertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeCertificateRefs/xades132:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeCertificateRefs/xades132:CertRefs/xades132:Cert")
		}
	})
	t.Run("AttributeCertificateRefsV2Path", func(t *testing.T) {
		got := p.AttributeCertificateRefsV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2")
		}
	})
	t.Run("AttributeCertificateRefsV2CertPath", func(t *testing.T) {
		got := p.AttributeCertificateRefsV2CertPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2/xades141:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AttributeCertificateRefsV2/xades141:CertRefs/xades132:Cert")
		}
	})
	t.Run("CertificateValuesPath", func(t *testing.T) {
		got := p.CertificateValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CertificateValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CertificateValues")
		}
	})
	t.Run("RevocationValuesPath", func(t *testing.T) {
		got := p.RevocationValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:RevocationValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:RevocationValues")
		}
	})
	t.Run("AttributeRevocationValuesPath", func(t *testing.T) {
		got := p.AttributeRevocationValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeRevocationValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttributeRevocationValues")
		}
	})
	t.Run("EncapsulatedCertificateValuesPath", func(t *testing.T) {
		got := p.EncapsulatedCertificateValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CertificateValues/xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:CertificateValues/xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("AttrAuthoritiesCertValuesPath", func(t *testing.T) {
		got := p.AttrAuthoritiesCertValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttrAuthoritiesCertValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttrAuthoritiesCertValues")
		}
	})
	t.Run("EncapsulatedAttrAuthoritiesCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedAttrAuthoritiesCertValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttrAuthoritiesCertValues/xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:AttrAuthoritiesCertValues/xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("EncapsulatedTimeStampValidationDataCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedTimeStampValidationDataCertValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData/xades132:CertificateValues/xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData/xades132:CertificateValues/xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("TimeStampValidationDataRevocationValuesPath", func(t *testing.T) {
		got := p.TimeStampValidationDataRevocationValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData/xades132:RevocationValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData/xades132:RevocationValues")
		}
	})
	t.Run("AnyValidationDataPath", func(t *testing.T) {
		got := p.AnyValidationDataPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData")
		}
	})
	t.Run("EncapsulatedAnyValidationDataCertValuesPath", func(t *testing.T) {
		got := p.EncapsulatedAnyValidationDataCertValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData/xades132:CertificateValues/xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData/xades132:CertificateValues/xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("AnyValidationDataRevocationValuesPath", func(t *testing.T) {
		got := p.AnyValidationDataRevocationValuesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData/xades132:RevocationValues" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:AnyValidationData/xades132:RevocationValues")
		}
	})
	t.Run("SignatureTimestampPath", func(t *testing.T) {
		got := p.SignatureTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:SignatureTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:SignatureTimeStamp")
		}
	})
	t.Run("SigAndRefsTimestampPath", func(t *testing.T) {
		got := p.SigAndRefsTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:SigAndRefsTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:SigAndRefsTimeStamp")
		}
	})
	t.Run("SigAndRefsTimestampV2Path", func(t *testing.T) {
		got := p.SigAndRefsTimestampV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:SigAndRefsTimeStampV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:SigAndRefsTimeStampV2")
		}
	})
	t.Run("RefsOnlyTimestampPath", func(t *testing.T) {
		got := p.RefsOnlyTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:RefsOnlyTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades132:RefsOnlyTimeStamp")
		}
	})
	t.Run("RefsOnlyTimestampV2Path", func(t *testing.T) {
		got := p.RefsOnlyTimestampV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:RefsOnlyTimeStampV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:RefsOnlyTimeStampV2")
		}
	})
	t.Run("ArchiveTimestampPath", func(t *testing.T) {
		got := p.ArchiveTimestampPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:ArchiveTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:ArchiveTimeStamp")
		}
	})
	t.Run("TimestampValidationDataPath", func(t *testing.T) {
		got := p.TimestampValidationDataPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:TimeStampValidationData")
		}
	})
	t.Run("SignaturePolicyStorePath", func(t *testing.T) {
		got := p.SignaturePolicyStorePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:SignaturePolicyStore" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xades141:SignaturePolicyStore")
		}
	})
	t.Run("SealingEvidenceRecordsPath", func(t *testing.T) {
		got := p.SealingEvidenceRecordsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xadesen:SealingEvidenceRecords" {
			t.Errorf("QueryString() = %q, want %q", qs, "./ds:Object/xades132:QualifyingProperties/xades132:UnsignedProperties/xades132:UnsignedSignatureProperties/xadesen:SealingEvidenceRecords")
		}
	})
	t.Run("CurrentCRLValuesChildren", func(t *testing.T) {
		got := p.CurrentCRLValuesChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLValues/xades132:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLValues/xades132:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentCRLRefsChildren", func(t *testing.T) {
		got := p.CurrentCRLRefsChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLRefs/xades132:CRLRef" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLRefs/xades132:CRLRef")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifier", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLIdentifier")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierIssuer", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierIssuer()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLIdentifier/xades132:Issuer" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLIdentifier/xades132:Issuer")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierIssueTime", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierIssueTime()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLIdentifier/xades132:IssueTime" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLIdentifier/xades132:IssueTime")
		}
	})
	t.Run("CurrentCRLRefCRLIdentifierNumber", func(t *testing.T) {
		got := p.CurrentCRLRefCRLIdentifierNumber()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLIdentifier/xades132:Number" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLIdentifier/xades132:Number")
		}
	})
	t.Run("CurrentOCSPValuesChildren", func(t *testing.T) {
		got := p.CurrentOCSPValuesChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPValues/xades132:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPValues/xades132:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentOCSPRefsChildren", func(t *testing.T) {
		got := p.CurrentOCSPRefsChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPRefs/xades132:OCSPRef" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPRefs/xades132:OCSPRef")
		}
	})
	t.Run("CurrentOCSPRefResponderID", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderID()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPIdentifier/xades132:ResponderID" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPIdentifier/xades132:ResponderID")
		}
	})
	t.Run("CurrentOCSPRefResponderIDByName", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderIDByName()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPIdentifier/xades132:ResponderID/xades132:ByName" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPIdentifier/xades132:ResponderID/xades132:ByName")
		}
	})
	t.Run("CurrentOCSPRefResponderIDByKey", func(t *testing.T) {
		got := p.CurrentOCSPRefResponderIDByKey()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPIdentifier/xades132:ResponderID/xades132:ByKey" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPIdentifier/xades132:ResponderID/xades132:ByKey")
		}
	})
	t.Run("CurrentOCSPRefProducedAt", func(t *testing.T) {
		got := p.CurrentOCSPRefProducedAt()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPIdentifier/xades132:ProducedAt" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPIdentifier/xades132:ProducedAt")
		}
	})
	t.Run("CurrentDigestAlgAndValue", func(t *testing.T) {
		got := p.CurrentDigestAlgAndValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:DigestAlgAndValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:DigestAlgAndValue")
		}
	})
	t.Run("CurrentCertRefsCertChildren", func(t *testing.T) {
		got := p.CurrentCertRefsCertChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CertRefs/xades132:Cert")
		}
	})
	t.Run("CurrentCertRefs141CertChildren", func(t *testing.T) {
		got := p.CurrentCertRefs141CertChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:CertRefs/xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:CertRefs/xades132:Cert")
		}
	})
	t.Run("CurrentCertChildren", func(t *testing.T) {
		got := p.CurrentCertChildren()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:Cert" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:Cert")
		}
	})
	t.Run("CurrentCertDigest", func(t *testing.T) {
		got := p.CurrentCertDigest()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CertDigest" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CertDigest")
		}
	})
	t.Run("CurrentEncapsulatedTimestamp", func(t *testing.T) {
		got := p.CurrentEncapsulatedTimestamp()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:EncapsulatedTimeStamp" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:EncapsulatedTimeStamp")
		}
	})
	t.Run("CurrentEncapsulatedCertificate", func(t *testing.T) {
		got := p.CurrentEncapsulatedCertificate()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("CurrentCertificateValuesEncapsulatedCertificate", func(t *testing.T) {
		got := p.CurrentCertificateValuesEncapsulatedCertificate()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CertificateValues/xades132:EncapsulatedX509Certificate" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CertificateValues/xades132:EncapsulatedX509Certificate")
		}
	})
	t.Run("CurrentRevocationValuesEncapsulatedOCSPValue", func(t *testing.T) {
		got := p.CurrentRevocationValuesEncapsulatedOCSPValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:RevocationValues/xades132:OCSPValues/xades132:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:RevocationValues/xades132:OCSPValues/xades132:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentEncapsulatedOCSPValue", func(t *testing.T) {
		got := p.CurrentEncapsulatedOCSPValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:OCSPValues/xades132:EncapsulatedOCSPValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:OCSPValues/xades132:EncapsulatedOCSPValue")
		}
	})
	t.Run("CurrentRevocationValuesEncapsulatedCRLValue", func(t *testing.T) {
		got := p.CurrentRevocationValuesEncapsulatedCRLValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:RevocationValues/xades132:CRLValues/xades132:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:RevocationValues/xades132:CRLValues/xades132:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentEncapsulatedCRLValue", func(t *testing.T) {
		got := p.CurrentEncapsulatedCRLValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CRLValues/xades132:EncapsulatedCRLValue" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CRLValues/xades132:EncapsulatedCRLValue")
		}
	})
	t.Run("CurrentIssuerSerialIssuerNamePath", func(t *testing.T) {
		got := p.CurrentIssuerSerialIssuerNamePath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:IssuerSerial/ds:X509IssuerName" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:IssuerSerial/ds:X509IssuerName")
		}
	})
	t.Run("CurrentIssuerSerialSerialNumberPath", func(t *testing.T) {
		got := p.CurrentIssuerSerialSerialNumberPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:IssuerSerial/ds:X509SerialNumber" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:IssuerSerial/ds:X509SerialNumber")
		}
	})
	t.Run("CurrentIssuerSerialV2Path", func(t *testing.T) {
		got := p.CurrentIssuerSerialV2Path()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:IssuerSerialV2" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:IssuerSerialV2")
		}
	})
	t.Run("CurrentCommitmentIdentifierPath", func(t *testing.T) {
		got := p.CurrentCommitmentIdentifierPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CommitmentTypeId/xades132:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CommitmentTypeId/xades132:Identifier")
		}
	})
	t.Run("CurrentCommitmentDescriptionPath", func(t *testing.T) {
		got := p.CurrentCommitmentDescriptionPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CommitmentTypeId/xades132:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CommitmentTypeId/xades132:Description")
		}
	})
	t.Run("CurrentCommitmentDocumentationReferencesPath", func(t *testing.T) {
		got := p.CurrentCommitmentDocumentationReferencesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:CommitmentTypeId/xades132:DocumentationReferences" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:CommitmentTypeId/xades132:DocumentationReferences")
		}
	})
	t.Run("CurrentDocumentationReference", func(t *testing.T) {
		got := p.CurrentDocumentationReference()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:DocumentationReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:DocumentationReference")
		}
	})
	t.Run("CurrentDescription", func(t *testing.T) {
		got := p.CurrentDescription()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:Description")
		}
	})
	t.Run("CurrentObjectIdentifier", func(t *testing.T) {
		got := p.CurrentObjectIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:ObjectIdentifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:ObjectIdentifier")
		}
	})
	t.Run("CurrentCommitmentObjectReferencesPath", func(t *testing.T) {
		got := p.CurrentCommitmentObjectReferencesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:ObjectReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:ObjectReference")
		}
	})
	t.Run("CurrentCommitmentAllSignedDataObjectsPath", func(t *testing.T) {
		got := p.CurrentCommitmentAllSignedDataObjectsPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:AllSignedDataObjects" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:AllSignedDataObjects")
		}
	})
	t.Run("CurrentMimeType", func(t *testing.T) {
		got := p.CurrentMimeType()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:MimeType" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:MimeType")
		}
	})
	t.Run("CurrentEncoding", func(t *testing.T) {
		got := p.CurrentEncoding()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:Encoding" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:Encoding")
		}
	})
	t.Run("CurrentSignaturePolicyId", func(t *testing.T) {
		got := p.CurrentSignaturePolicyId()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:Identifier")
		}
	})
	t.Run("CurrentSignaturePolicyDigestAlgAndValue", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDigestAlgAndValue()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyHash" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyHash")
		}
	})
	t.Run("CurrentSignaturePolicySPURI", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPURI()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades132:SPURI" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades132:SPURI")
		}
	})
	t.Run("CurrentSignaturePolicySPUserNotice", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPUserNotice()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades132:SPUserNotice" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades132:SPUserNotice")
		}
	})
	t.Run("CurrentSPUserNoticeNoticeRefOrganization", func(t *testing.T) {
		got := p.CurrentSPUserNoticeNoticeRefOrganization()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:NoticeRef/xades132:Organization" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:NoticeRef/xades132:Organization")
		}
	})
	t.Run("CurrentSPUserNoticeNoticeRefNoticeNumbers", func(t *testing.T) {
		got := p.CurrentSPUserNoticeNoticeRefNoticeNumbers()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:NoticeRef/xades132:NoticeNumbers" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:NoticeRef/xades132:NoticeNumbers")
		}
	})
	t.Run("CurrentSPUserNoticeExplicitText", func(t *testing.T) {
		got := p.CurrentSPUserNoticeExplicitText()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:ExplicitText" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:ExplicitText")
		}
	})
	t.Run("CurrentSignaturePolicySPDocSpecification", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPDocSpecification()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades141:SPDocSpecification" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades141:SPDocSpecification")
		}
	})
	t.Run("CurrentSignaturePolicySPDocSpecificationIdentifier", func(t *testing.T) {
		got := p.CurrentSignaturePolicySPDocSpecificationIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades141:SPDocSpecification/xades132:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers/xades132:SigPolicyQualifier/xades141:SPDocSpecification/xades132:Identifier")
		}
	})
	t.Run("CurrentSignaturePolicyDescription", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDescription()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:Description")
		}
	})
	t.Run("CurrentSignaturePolicyDocumentationReferences", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDocumentationReferences()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:DocumentationReferences" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyId/xades132:DocumentationReferences")
		}
	})
	t.Run("CurrentSignaturePolicyImplied", func(t *testing.T) {
		got := p.CurrentSignaturePolicyImplied()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyImplied" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyImplied")
		}
	})
	t.Run("CurrentSignaturePolicyTransforms", func(t *testing.T) {
		got := p.CurrentSignaturePolicyTransforms()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/ds:Transforms" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/ds:Transforms")
		}
	})
	t.Run("CurrentSignaturePolicyQualifiers", func(t *testing.T) {
		got := p.CurrentSignaturePolicyQualifiers()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:SignaturePolicyId/xades132:SigPolicyQualifiers")
		}
	})
	t.Run("CurrentInclude", func(t *testing.T) {
		got := p.CurrentInclude()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:Include" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:Include")
		}
	})
	t.Run("CurrentQualifyingPropertiesPath", func(t *testing.T) {
		got := p.CurrentQualifyingPropertiesPath()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:QualifyingProperties" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:QualifyingProperties")
		}
	})
	t.Run("CurrentSPDocSpecification", func(t *testing.T) {
		got := p.CurrentSPDocSpecification()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SPDocSpecification" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SPDocSpecification")
		}
	})
	t.Run("CurrentIdentifier", func(t *testing.T) {
		got := p.CurrentIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:Identifier")
		}
	})
	t.Run("CurrentSPDocSpecificationIdentifier", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationIdentifier()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SPDocSpecification/xades132:Identifier" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SPDocSpecification/xades132:Identifier")
		}
	})
	t.Run("CurrentSPDocSpecificationDescription", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationDescription()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SPDocSpecification/xades132:Description" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SPDocSpecification/xades132:Description")
		}
	})
	t.Run("CurrentDocumentationReferenceElements", func(t *testing.T) {
		got := p.CurrentDocumentationReferenceElements()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades132:DocumentationReferences/xades132:DocumentationReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades132:DocumentationReferences/xades132:DocumentationReference")
		}
	})
	t.Run("CurrentSPDocSpecificationDocumentationReferenceElements", func(t *testing.T) {
		got := p.CurrentSPDocSpecificationDocumentationReferenceElements()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SPDocSpecification/xades132:DocumentationReferences/xades132:DocumentationReference" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SPDocSpecification/xades132:DocumentationReferences/xades132:DocumentationReference")
		}
	})
	t.Run("CurrentSignaturePolicyDocument", func(t *testing.T) {
		got := p.CurrentSignaturePolicyDocument()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SignaturePolicyDocument" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SignaturePolicyDocument")
		}
	})
	t.Run("CurrentSigPolDocLocalURI", func(t *testing.T) {
		got := p.CurrentSigPolDocLocalURI()
		if got == nil {
			t.Fatalf("expected non-nil XPathQuery")
		}
		if qs := got.QueryString(); qs != "./xades141:SigPolDocLocalURI" {
			t.Errorf("QueryString() = %q, want %q", qs, "./xades141:SigPolDocLocalURI")
		}
	})

	if got := p.Namespace().Uri(); got != XAdESNamespace_XADES_132.Uri() {
		t.Errorf("Namespace().Uri() = %q, want %q", got, XAdESNamespace_XADES_132.Uri())
	}
	if got := p.SignedPropertiesUri(); got != "http://uri.etsi.org/01903#SignedProperties" {
		t.Errorf("SignedPropertiesUri() = %q, want %q", got, "http://uri.etsi.org/01903#SignedProperties")
	}
	if got := p.CounterSignatureUri(); got != "http://uri.etsi.org/01903#CountersignedSignature" {
		t.Errorf("CounterSignatureUri() = %q", got)
	}
}
