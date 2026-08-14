// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DocumentInfo.java (DSS 6.5.RC1).
package job

import (
	"testing"

	"github.com/utain/esig/dss/model"
)

type fakeDocumentInfoIdentifier struct {
	model.IdentifierBase
}

func newFakeDocumentInfoIdentifier(data string) model.Identifier {
	return &fakeDocumentInfoIdentifier{model.NewIdentifierBase("FakeDocumentInfo", "F-", []byte(data))}
}

type fakeDocumentInfo struct {
	id         model.Identifier
	url        string
	parent     *fakeDocumentInfo
	download   DownloadInfoRecord
	parsing    ParsingInfoRecord
	validation ValidationInfoRecord
}

func (f *fakeDocumentInfo) DSSID() model.Identifier                   { return f.id }
func (f *fakeDocumentInfo) DownloadCacheInfo() DownloadInfoRecord     { return f.download }
func (f *fakeDocumentInfo) ParsingCacheInfo() ParsingInfoRecord       { return f.parsing }
func (f *fakeDocumentInfo) ValidationCacheInfo() ValidationInfoRecord { return f.validation }
func (f *fakeDocumentInfo) URL() string                               { return f.url }
func (f *fakeDocumentInfo) Parent() *fakeDocumentInfo                 { return f.parent }
func (f *fakeDocumentInfo) DSSIDAsString() string                     { return f.id.AsXmlID() }

var _ DocumentInfo[*fakeDocumentInfo] = (*fakeDocumentInfo)(nil)

func TestDocumentInfo_RoundTrip(t *testing.T) {
	parent := &fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("parent"), url: "http://parent"}
	child := &fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("child"), url: "http://child", parent: parent}

	var di DocumentInfo[*fakeDocumentInfo] = child
	if di.URL() != "http://child" {
		t.Fatalf("URL() = %q", di.URL())
	}
	if di.Parent() != parent {
		t.Fatalf("Parent() did not round-trip")
	}
	if di.DSSIDAsString() != child.id.AsXmlID() {
		t.Fatalf("DSSIDAsString() mismatch")
	}
}
