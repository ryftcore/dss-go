// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DocumentListInfo.java (DSS 6.5.RC1).
package job

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
)

type fakeDocumentListInfo struct {
	fakeDocumentInfo
	children []*fakeDocumentInfo
}

func (f *fakeDocumentListInfo) ChildrenInfos() []*fakeDocumentInfo { return f.children }

var _ DocumentListInfo[*fakeDocumentInfo, *fakeDocumentInfo] = (*fakeDocumentListInfo)(nil)

func TestDocumentListInfo_RoundTrip(t *testing.T) {
	child := &fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("child"), url: "http://child"}
	list := &fakeDocumentListInfo{
		fakeDocumentInfo: fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("list"), url: "http://list"},
		children:         []*fakeDocumentInfo{child},
	}

	var dli DocumentListInfo[*fakeDocumentInfo, *fakeDocumentInfo] = list
	children := dli.ChildrenInfos()
	if len(children) != 1 || children[0] != child {
		t.Fatalf("ChildrenInfos() = %v", children)
	}
	var _ model.IdentifierBasedObject = dli
}
