// KAT test for trusted_list_namespace.go: TrustedListNamespaceNS URI/prefix transcribed
// verbatim from upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.tsl.TrustedListNamespace, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestTrustedListNamespace_KAT(t *testing.T) {
	if got := TrustedListNamespaceNS.Uri(); got != "http://uri.etsi.org/02231/v2#" {
		t.Errorf("TrustedListNamespaceNS.Uri() = %q, want %q", got, "http://uri.etsi.org/02231/v2#")
	}
	if got := TrustedListNamespaceNS.Prefix(); got != "tl" {
		t.Errorf("TrustedListNamespaceNS.Prefix() = %q, want %q", got, "tl")
	}
}
