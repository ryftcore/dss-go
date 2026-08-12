// Ported from dss-enumerations/.../loader/LoTELoader.java (DSS 6.5.RC1).
package enumerations

import "testing"

type fakeLoTELoader struct {
	name               string
	listType           ListType
	sti                LoTEServiceTypeIdentifier
	status             LoTEServiceStatus
	certApprovalStatus CertificateApprovalStatus
}

func (f *fakeLoTELoader) ListTypeFromURI(uri string) ListType { return f.listType }
func (f *fakeLoTELoader) ServiceTypeIdentifierFromURI(uri string) LoTEServiceTypeIdentifier {
	return f.sti
}
func (f *fakeLoTELoader) ServiceStatusFromURI(uri string) LoTEServiceStatus { return f.status }
func (f *fakeLoTELoader) CertificateApprovalStatusFromLabel(label string) CertificateApprovalStatus {
	return f.certApprovalStatus
}
func (f *fakeLoTELoader) CertificateApprovalStatusFromDefinition(listType ListType, sti LoTEServiceTypeIdentifier, status LoTEServiceStatus) CertificateApprovalStatus {
	return f.certApprovalStatus
}

func withLoTELoaders(t *testing.T, loaders ...LoTELoader) {
	t.Helper()
	saved := loTELoaderRegistry
	loTELoaderRegistry = nil
	for _, l := range loaders {
		RegisterLoTELoader(l)
	}
	t.Cleanup(func() { loTELoaderRegistry = saved })
}

func TestRegisterLoTELoader_PreservesOrder(t *testing.T) {
	a := &fakeLoTELoader{name: "a"}
	b := &fakeLoTELoader{name: "b"}
	withLoTELoaders(t, a, b)

	loaders := loTELoaders()
	if len(loaders) != 2 || loaders[0].(*fakeLoTELoader).name != "a" || loaders[1].(*fakeLoTELoader).name != "b" {
		t.Errorf("loTELoaders() = %v, want [a, b] in registration order", loaders)
	}
}

func TestLoTELoaders_FirstMatchWins(t *testing.T) {
	miss := &fakeLoTELoader{name: "miss"}
	hitListType := &loteEmptyListType{uri: "urn:hit"}
	hit := &fakeLoTELoader{name: "hit", listType: hitListType}
	withLoTELoaders(t, miss, hit)

	got := ListTypeFromURI("urn:example")
	if got != ListType(hitListType) {
		t.Errorf("ListTypeFromURI = %v, want %v (first non-nil loader result)", got, hitListType)
	}
}

// TestDefaultLoTELoaders_MatchUpstreamServiceRegistration pins the default registry to
// what upstream declares in
// META-INF/services/eu.europa.esig.dss.enumerations.loader.LoTELoader: LoTEEnumLoader
// first, LoTEEmptyLoader second. The order is load-bearing — LoTEEmptyLoader answers every
// lookup with a non-nil echo value, so registering it first would shadow LoTEEnumLoader and
// stop any URI from resolving to its enum constant.
func TestDefaultLoTELoaders_MatchUpstreamServiceRegistration(t *testing.T) {
	loaders := loTELoaders()
	if len(loaders) != 2 {
		t.Fatalf("default loTELoaders() has %d entries, want 2", len(loaders))
	}
	if _, ok := loaders[0].(*LoTEEnumLoader); !ok {
		t.Errorf("default loader[0] = %T, want *LoTEEnumLoader", loaders[0])
	}
	if _, ok := loaders[1].(*LoTEEmptyLoader); !ok {
		t.Errorf("default loader[1] = %T, want *LoTEEmptyLoader", loaders[1])
	}
}

// TestDefaultLoTELoaders_Behavior checks the lookups a caller sees with the default
// registry, matching a JDK run of the upstream enumerations with their services file on the
// classpath: a known URI resolves to its enum constant, and an unknown one falls through to
// LoTEEmptyLoader, which echoes the URI back with an empty label rather than returning nil.
func TestDefaultLoTELoaders_Behavior(t *testing.T) {
	lt := ListTypeFromURI("http://uri.etsi.org/19602/LoTEType/EUPIDProvidersList")
	if got, ok := lt.(LoTETypeEnum); !ok || got != LoTETypeEnum_EUPIDProvidersList {
		t.Errorf("ListTypeFromURI(known) = %#v, want LoTETypeEnum_EUPIDProvidersList", lt)
	}
	unknown := ListTypeFromURI("bogus")
	if unknown == nil {
		t.Fatal("ListTypeFromURI(bogus) = nil, want a LoTEEmptyLoader echo value")
	}
	if unknown.URI() != "bogus" || unknown.Label() != "" {
		t.Errorf("ListTypeFromURI(bogus) = {label=%q,uri=%q}, want {label=\"\",uri=\"bogus\"}", unknown.Label(), unknown.URI())
	}

	sti := LoTEServiceTypeIdentifierFromURI("http://uri.etsi.org/19602/SvcType/Register")
	if got, ok := sti.(LoTEServiceTypeIdentifierEnum); !ok || got != LoTEServiceTypeIdentifierEnum_REGISTER {
		t.Errorf("LoTEServiceTypeIdentifierFromURI(known) = %#v, want REGISTER", sti)
	}
	status := LoTEServiceStatusFromURI("http://uri.etsi.org/19602/PubEAAProvidersList/SvcStatus/notified")
	if got, ok := status.(LoTEServiceStatusEnum); !ok || got != LoTEServiceStatusEnum_PUB_EAA_PROVIDER_NOTIFIED {
		t.Errorf("LoTEServiceStatusFromURI(known) = %#v, want PUB_EAA_PROVIDER_NOTIFIED", status)
	}

	// LoTEEnumLoader returns CERT_FOR_UNKNOWN rather than nil, so LoTEEmptyLoader (whose
	// fromDefinition is unsupported and returns nil) is never consulted here.
	cas := CertificateApprovalStatusFromDefinition(nil, nil, nil)
	if got, ok := cas.(CertificateApprovalStatusEnum); !ok || got != CertificateApprovalStatusEnum_CERT_FOR_UNKNOWN {
		t.Errorf("CertificateApprovalStatusFromDefinition(nil,nil,nil) = %#v, want CERT_FOR_UNKNOWN", cas)
	}
}
