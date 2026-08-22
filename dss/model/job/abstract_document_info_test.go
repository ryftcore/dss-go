// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/AbstractDocumentInfo.java (DSS 6.5.RC1).
package job

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

type fakeAbstractDocInfoIdentifier struct {
	model.IdentifierBase
}

func newFakeAbstractDocInfoIdentifier(data string) model.Identifier {
	return &fakeAbstractDocInfoIdentifier{model.NewIdentifierBase("FakeDocumentInfoIdentifier", "F-", []byte(data))}
}

type fakeAbstractDocInfo struct {
	AbstractDocumentInfoBase[*fakeAbstractDocInfo]
	buildCount int
}

func newFakeAbstractDocInfo(url string, parent *fakeAbstractDocInfo) *fakeAbstractDocInfo {
	d := &fakeAbstractDocInfo{
		AbstractDocumentInfoBase: NewAbstractDocumentInfoBaseWithParent[*fakeAbstractDocInfo](nil, nil, nil, url, parent),
	}
	d.InitAbstractDocumentInfo(d)
	return d
}

func (d *fakeAbstractDocInfo) BuildIdentifier() model.Identifier {
	d.buildCount++
	return newFakeAbstractDocInfoIdentifier(d.URL())
}

func TestAbstractDocumentInfoBase_DSSIDCachesIdentifier(t *testing.T) {
	d := newFakeAbstractDocInfo("http://example.org/tl.xml", nil)

	first := d.DSSID()
	second := d.DSSID()
	if first != second {
		t.Fatalf("DSSID() should return the cached identifier on subsequent calls")
	}
	if d.buildCount != 1 {
		t.Fatalf("BuildIdentifier should be invoked exactly once, got %d", d.buildCount)
	}
	if got, want := d.DSSIDAsString(), first.AsXmlID(); got != want {
		t.Fatalf("DSSIDAsString() = %q, want %q", got, want)
	}
	if got, want := d.URL(), "http://example.org/tl.xml"; got != want {
		t.Fatalf("URL() = %q, want %q", got, want)
	}
}

func TestAbstractDocumentInfoBase_Parent(t *testing.T) {
	parent := newFakeAbstractDocInfo("http://parent.example.org", nil)
	child := newFakeAbstractDocInfo("http://child.example.org", parent)

	if child.Parent() != parent {
		t.Fatalf("Parent() did not round-trip the constructor argument")
	}
	if got := newFakeAbstractDocInfo("http://no-parent.example.org", nil).Parent(); got != nil {
		t.Fatalf("Parent() = %v, want nil for the single-argument constructor", got)
	}
}

func TestAbstractDocumentInfoBase_PanicsOnEmptyURL(t *testing.T) {
	defer func() {
		if r := recover(); r != "URL String shall be provided!" {
			t.Fatalf("recover() = %v, want the Java requireNonNull message", r)
		}
	}()
	NewAbstractDocumentInfoBase[*fakeAbstractDocInfo](nil, nil, nil, "")
}
