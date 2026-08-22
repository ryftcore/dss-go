// KAT test for xadesen_evidencerecord_namespace_element.go: every
// XAdESEvidencerecordNamespaceElement_* tag name below is transcribed verbatim from the
// upstream Java source (dss-xades 6.5.RC1).
package definition

import "testing"

func TestXAdESEvidencerecordNamespaceElement_KAT(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ASN1_EVIDENCE_RECORD", XAdESEvidencerecordNamespaceElementASN1EvidenceRecord.TagName(), "ASN1EvidenceRecord"},
		{"EVIDENCE_RECORD", XAdESEvidencerecordNamespaceElementEvidenceRecord.TagName(), "EvidenceRecord"},
		{"SEALING_EVIDENCE_RECORDS", XAdESEvidencerecordNamespaceElementSealingEvidenceRecords.TagName(), "SealingEvidenceRecords"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s.TagName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}
