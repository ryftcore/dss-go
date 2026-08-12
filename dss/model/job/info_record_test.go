// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/InfoRecord.java (DSS 6.5.RC1).
package job

import (
	"testing"
	"time"
)

type fakeInfoRecord struct {
	refreshNeeded            bool
	desynchronized           bool
	synchronized             bool
	errorFlag                bool
	toBeDeleted              bool
	statusName               string
	lastStateTransitionTime  time.Time
	lastSuccessSyncTime      time.Time
	exceptionMessage         string
	exceptionStackTrace      string
	exceptionFirstOccurrence time.Time
	exceptionLastOccurrence  time.Time
	resultExist              bool
}

func (f *fakeInfoRecord) IsRefreshNeeded() bool                     { return f.refreshNeeded }
func (f *fakeInfoRecord) IsDesynchronized() bool                    { return f.desynchronized }
func (f *fakeInfoRecord) IsSynchronized() bool                      { return f.synchronized }
func (f *fakeInfoRecord) IsError() bool                             { return f.errorFlag }
func (f *fakeInfoRecord) IsToBeDeleted() bool                       { return f.toBeDeleted }
func (f *fakeInfoRecord) StatusName() string                        { return f.statusName }
func (f *fakeInfoRecord) LastStateTransitionTime() time.Time        { return f.lastStateTransitionTime }
func (f *fakeInfoRecord) LastSuccessSynchronizationTime() time.Time { return f.lastSuccessSyncTime }
func (f *fakeInfoRecord) ExceptionMessage() string                  { return f.exceptionMessage }
func (f *fakeInfoRecord) ExceptionStackTrace() string               { return f.exceptionStackTrace }
func (f *fakeInfoRecord) ExceptionFirstOccurrenceTime() time.Time   { return f.exceptionFirstOccurrence }
func (f *fakeInfoRecord) ExceptionLastOccurrenceTime() time.Time    { return f.exceptionLastOccurrence }
func (f *fakeInfoRecord) IsResultExist() bool                       { return f.resultExist }

var _ InfoRecord = (*fakeInfoRecord)(nil)

func TestInfoRecord_RoundTrip(t *testing.T) {
	now := time.Now()
	rec := &fakeInfoRecord{
		refreshNeeded:            true,
		desynchronized:           true,
		synchronized:             false,
		errorFlag:                true,
		toBeDeleted:              true,
		statusName:               "ERROR",
		lastStateTransitionTime:  now,
		lastSuccessSyncTime:      now.Add(-time.Hour),
		exceptionMessage:         "boom",
		exceptionStackTrace:      "trace",
		exceptionFirstOccurrence: now.Add(-2 * time.Hour),
		exceptionLastOccurrence:  now,
		resultExist:              true,
	}

	var ir InfoRecord = rec
	if !ir.IsRefreshNeeded() || !ir.IsDesynchronized() || ir.IsSynchronized() {
		t.Fatalf("boolean getters did not round-trip")
	}
	if ir.StatusName() != "ERROR" {
		t.Fatalf("StatusName() = %q", ir.StatusName())
	}
	if !ir.LastStateTransitionTime().Equal(now) {
		t.Fatalf("LastStateTransitionTime() mismatch")
	}
	if ir.ExceptionMessage() != "boom" || ir.ExceptionStackTrace() != "trace" {
		t.Fatalf("exception fields did not round-trip")
	}
	if !ir.IsResultExist() {
		t.Fatalf("IsResultExist() = false, want true")
	}
}
