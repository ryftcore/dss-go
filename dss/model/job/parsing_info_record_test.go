// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/job/ParsingInfoRecord.java (DSS 6.5.RC1).
package job

import "testing"

type fakeParsingInfoRecord struct {
	fakeInfoRecord
	structureValidationMessages []string
}

func (f *fakeParsingInfoRecord) StructureValidationMessages() []string {
	return f.structureValidationMessages
}

var _ ParsingInfoRecord = (*fakeParsingInfoRecord)(nil)

func TestParsingInfoRecord_RoundTrip(t *testing.T) {
	rec := &fakeParsingInfoRecord{structureValidationMessages: []string{"unexpected element"}}

	var pir ParsingInfoRecord = rec
	msgs := pir.StructureValidationMessages()
	if len(msgs) != 1 || msgs[0] != "unexpected element" {
		t.Fatalf("StructureValidationMessages() = %v", msgs)
	}
}
