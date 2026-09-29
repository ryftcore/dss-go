package blocks

import (
	"fmt"
	"testing"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

func certificateWrapper(id string) *diagnostic.CertificateWrapper {
	cs := jaxb.CollapsedString(id)
	return diagnostic.NewCertificateWrapper(&jaxb.XmlCertificate{XmlAbstractTokenAttrs: jaxb.XmlAbstractTokenAttrs{Id: &cs}})
}

func orphanCertificateWrapper(id string) *diagnostic.OrphanCertificateTokenWrapper {
	cs := jaxb.CollapsedString(id)
	return diagnostic.NewOrphanCertificateTokenWrapper(&jaxb.XmlOrphanCertificateToken{XmlAbstractTokenAttrs: jaxb.XmlAbstractTokenAttrs{Id: &cs}})
}

// TestRemoveAllCertificates covers T32-PERF-BIGO-001: removeAllCertificates /
// removeAllOrphanCertificates are List#removeAll ports (retain the elements that are not
// equal - same Id - to any element of the second list, keeping their order), now built on an id
// set instead of nested loops.
func TestRemoveAllCertificates(t *testing.T) {
	certificates := []*diagnostic.CertificateWrapper{
		certificateWrapper("C-1"), certificateWrapper("C-2"), certificateWrapper("C-3"), certificateWrapper("C-4"),
	}
	toRemove := []*diagnostic.CertificateWrapper{certificateWrapper("C-4"), certificateWrapper("C-2"), certificateWrapper("C-9")}

	got := removeAllCertificates(certificates, toRemove)
	if len(got) != 2 || got[0] != certificates[0] || got[1] != certificates[2] {
		t.Fatalf("removeAllCertificates = %v, want [C-1 C-3] in order", certificateIds(got))
	}

	if got := removeAllCertificates(certificates, nil); len(got) != len(certificates) {
		t.Errorf("nothing to remove: got %v, want all four", certificateIds(got))
	}
	if got := removeAllCertificates(nil, toRemove); len(got) != 0 {
		t.Errorf("nothing to remove from: got %v, want empty", certificateIds(got))
	}
	if got := removeAllCertificates(certificates, certificates); len(got) != 0 {
		t.Errorf("removing everything: got %v, want empty", certificateIds(got))
	}
}

func TestRemoveAllOrphanCertificates(t *testing.T) {
	orphans := []*diagnostic.OrphanCertificateTokenWrapper{
		orphanCertificateWrapper("O-1"), orphanCertificateWrapper("O-2"), orphanCertificateWrapper("O-3"),
	}
	toRemove := []*diagnostic.OrphanCertificateTokenWrapper{orphanCertificateWrapper("O-2")}

	got := removeAllOrphanCertificates(orphans, toRemove)
	if len(got) != 2 || got[0] != orphans[0] || got[1] != orphans[2] {
		t.Fatalf("removeAllOrphanCertificates = %d elements, want [O-1 O-3] in order", len(got))
	}
	if got := removeAllOrphanCertificates(orphans, nil); len(got) != len(orphans) {
		t.Errorf("nothing to remove: got %d, want %d", len(got), len(orphans))
	}
	if got := removeAllOrphanCertificates(nil, toRemove); len(got) != 0 {
		t.Errorf("nothing to remove from: got %d, want 0", len(got))
	}
}

// TestRemoveAllCertificatesLargeSets pins the linear behaviour: attacker-sized cross/equivalent
// certificate sets (tens of thousands of certificates sharing a subject and key) must not cost
// a quadratic number of comparisons (the nested loops made this case hundreds of millions of Equals calls).
func TestRemoveAllCertificatesLargeSets(t *testing.T) {
	const n = 40000
	certificates := make([]*diagnostic.CertificateWrapper, n)
	toRemove := make([]*diagnostic.CertificateWrapper, n/2)
	for i := range certificates {
		certificates[i] = certificateWrapper(fmt.Sprintf("C-%d", i))
	}
	for i := range toRemove {
		toRemove[i] = certificateWrapper(fmt.Sprintf("C-%d", 2*i)) // every even one
	}
	got := removeAllCertificates(certificates, toRemove)
	if len(got) != n/2 {
		t.Fatalf("got %d certificates, want %d", len(got), n/2)
	}
}

func certificateIds(wrappers []*diagnostic.CertificateWrapper) []string {
	ids := make([]string, len(wrappers))
	for i, w := range wrappers {
		ids[i] = w.Id()
	}
	return ids
}
