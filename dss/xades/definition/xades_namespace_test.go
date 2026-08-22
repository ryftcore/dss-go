// KAT test for xades_namespace.go: every XAdESNamespace_* URI/prefix pair below is
// transcribed verbatim from the upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.XAdESNamespace, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestXAdESNamespace_KAT(t *testing.T) {
	cases := []struct {
		name   string
		uri    string
		prefix string
	}{
		{"XMLDSIG_FILTER2", "http://www.w3.org/2002/06/xmldsig-filter2", "dsig-filter2"},
		{"XADES_111", "http://uri.etsi.org/01903/v1.1.1#", "xades111"},
		{"XADES_122", "http://uri.etsi.org/01903/v1.2.2#", "xades122"},
		{"XADES_132", "http://uri.etsi.org/01903/v1.3.2#", "xades132"},
		{"XADES_141", "http://uri.etsi.org/01903/v1.4.1#", "xades141"},
		{"XADES_EVIDENCERECORD_NAMESPACE", "http://uri.etsi.org/19132/v1.1.1#", "xadesen"},
	}
	for _, c := range cases {
		if c.uri == "" || c.prefix == "" {
			t.Fatalf("%s: empty uri/prefix", c.name)
		}
	}
	// Direct field checks (also exercises Uri()/Prefix()).
	if got := XAdESNamespaceXMLDSIGFilter2.Uri(); got != "http://www.w3.org/2002/06/xmldsig-filter2" {
		t.Errorf("XMLDSIG_FILTER2.Uri() = %q, want %q", got, "http://www.w3.org/2002/06/xmldsig-filter2")
	}
	if got := XAdESNamespaceXMLDSIGFilter2.Prefix(); got != "dsig-filter2" {
		t.Errorf("XMLDSIG_FILTER2.Prefix() = %q, want %q", got, "dsig-filter2")
	}
	if got := XAdESNamespaceXAdES111.Uri(); got != "http://uri.etsi.org/01903/v1.1.1#" {
		t.Errorf("XADES_111.Uri() = %q, want %q", got, "http://uri.etsi.org/01903/v1.1.1#")
	}
	if got := XAdESNamespaceXAdES111.Prefix(); got != "xades111" {
		t.Errorf("XADES_111.Prefix() = %q, want %q", got, "xades111")
	}
	if got := XAdESNamespaceXAdES122.Uri(); got != "http://uri.etsi.org/01903/v1.2.2#" {
		t.Errorf("XADES_122.Uri() = %q, want %q", got, "http://uri.etsi.org/01903/v1.2.2#")
	}
	if got := XAdESNamespaceXAdES122.Prefix(); got != "xades122" {
		t.Errorf("XADES_122.Prefix() = %q, want %q", got, "xades122")
	}
	if got := XAdESNamespaceXAdES132.Uri(); got != "http://uri.etsi.org/01903/v1.3.2#" {
		t.Errorf("XADES_132.Uri() = %q, want %q", got, "http://uri.etsi.org/01903/v1.3.2#")
	}
	if got := XAdESNamespaceXAdES132.Prefix(); got != "xades132" {
		t.Errorf("XADES_132.Prefix() = %q, want %q", got, "xades132")
	}
	if got := XAdESNamespaceXAdES141.Uri(); got != "http://uri.etsi.org/01903/v1.4.1#" {
		t.Errorf("XADES_141.Uri() = %q, want %q", got, "http://uri.etsi.org/01903/v1.4.1#")
	}
	if got := XAdESNamespaceXAdES141.Prefix(); got != "xades141" {
		t.Errorf("XADES_141.Prefix() = %q, want %q", got, "xades141")
	}
	if got := XAdESNamespaceXAdESEvidencerecordNamespace.Uri(); got != "http://uri.etsi.org/19132/v1.1.1#" {
		t.Errorf("XADES_EVIDENCERECORD_NAMESPACE.Uri() = %q, want %q", got, "http://uri.etsi.org/19132/v1.1.1#")
	}
	if got := XAdESNamespaceXAdESEvidencerecordNamespace.Prefix(); got != "xadesen" {
		t.Errorf("XADES_EVIDENCERECORD_NAMESPACE.Prefix() = %q, want %q", got, "xadesen")
	}
}
