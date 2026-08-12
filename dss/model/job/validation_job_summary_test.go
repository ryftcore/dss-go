// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ValidationJobSummary.java (DSS 6.5.RC1).
package job

import "testing"

type fakeValidationJobSummary struct {
	documentListInfos  []*fakeDocumentListInfo
	otherDocumentInfos []*fakeDocumentInfo
}

func (f *fakeValidationJobSummary) DocumentListInfos() []*fakeDocumentListInfo {
	return f.documentListInfos
}
func (f *fakeValidationJobSummary) OtherDocumentInfos() []*fakeDocumentInfo {
	return f.otherDocumentInfos
}

var _ ValidationJobSummary[*fakeDocumentInfo, *fakeDocumentListInfo] = (*fakeValidationJobSummary)(nil)

func TestValidationJobSummary_RoundTrip(t *testing.T) {
	list := &fakeDocumentListInfo{fakeDocumentInfo: fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("list"), url: "http://list"}}
	other := &fakeDocumentInfo{id: newFakeDocumentInfoIdentifier("other"), url: "http://other"}
	summary := &fakeValidationJobSummary{
		documentListInfos:  []*fakeDocumentListInfo{list},
		otherDocumentInfos: []*fakeDocumentInfo{other},
	}

	var vjs ValidationJobSummary[*fakeDocumentInfo, *fakeDocumentListInfo] = summary
	if got := vjs.DocumentListInfos(); len(got) != 1 || got[0] != list {
		t.Fatalf("DocumentListInfos() = %v", got)
	}
	if got := vjs.OtherDocumentInfos(); len(got) != 1 || got[0] != other {
		t.Fatalf("OtherDocumentInfos() = %v", got)
	}
}
