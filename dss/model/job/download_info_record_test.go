// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/DownloadInfoRecord.java (DSS 6.5.RC1).
package job

import (
	"testing"
	"time"

	"github.com/ryftcore/dss-go/dss/model"
)

type fakeDownloadInfoRecord struct {
	fakeInfoRecord
	document                model.DSSDocument
	lastDownloadAttemptTime time.Time
}

func (f *fakeDownloadInfoRecord) Document() model.DSSDocument { return f.document }
func (f *fakeDownloadInfoRecord) LastDownloadAttemptTime() time.Time {
	return f.lastDownloadAttemptTime
}

var _ DownloadInfoRecord = (*fakeDownloadInfoRecord)(nil)

func TestDownloadInfoRecord_RoundTrip(t *testing.T) {
	now := time.Now()
	doc := model.NewInMemoryDocument([]byte("hello"))
	rec := &fakeDownloadInfoRecord{document: doc, lastDownloadAttemptTime: now}

	var dir DownloadInfoRecord = rec
	if dir.Document() != doc {
		t.Fatalf("Document() did not round-trip")
	}
	if !dir.LastDownloadAttemptTime().Equal(now) {
		t.Fatalf("LastDownloadAttemptTime() did not round-trip")
	}
}
