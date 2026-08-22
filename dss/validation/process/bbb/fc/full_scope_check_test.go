package fc

import (
	"testing"

	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

func fullScope(name string, kind enumerations.SignatureScopeType) *diagjaxb.XmlSignatureScope {
	k := diagjaxb.SignatureScopeTypeValue(kind)
	return &diagjaxb.XmlSignatureScope{Scope: &k, Name: &name}
}

func TestFullScopeCheck_Process(t *testing.T) {
	tests := []struct {
		name   string
		scopes []*diagjaxb.XmlSignatureScope
		want   bool
	}{
		{"happy: no scopes", nil, true},
		{"happy: all FULL", []*diagjaxb.XmlSignatureScope{
			fullScope("doc1.xml", enumerations.SignatureScopeTypeFull),
			fullScope("doc2.xml", enumerations.SignatureScopeTypeFull),
		}, true},
		{"failure: one PARTIAL", []*diagjaxb.XmlSignatureScope{
			fullScope("doc1.xml", enumerations.SignatureScopeTypeFull),
			fullScope("doc2.xml", enumerations.SignatureScopeTypePartial),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFullScopeCheck(testI18nProvider(t), newTestFCResult(), tt.scopes,
				process.GetLevelRule(enumerations.LevelFail))
			if got := c.Process(); got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
