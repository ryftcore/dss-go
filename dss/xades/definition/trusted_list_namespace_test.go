// KAT test for trusted_list_namespace.go: TrustedListNamespace_NS URI/prefix transcribed
// verbatim from upstream Java source (dss-xades 6.5.RC1)
// eu.europa.esig.dss.xades.definition.tsl.TrustedListNamespace, per PORTING.md's exhaustive
// table-test rule for registry-like tables.
package definition

import "testing"

func TestTrustedListNamespace_KAT(t *testing.T) {
	if got := TrustedListNamespace_NS.Uri(); got != "http://uri.etsi.org/02231/v2#" {
		t.Errorf("TrustedListNamespace_NS.Uri() = %q, want %q", got, "http://uri.etsi.org/02231/v2#")
	}
	if got := TrustedListNamespace_NS.Prefix(); got != "tl" {
		t.Errorf("TrustedListNamespace_NS.Prefix() = %q, want %q", got, "tl")
	}
}
