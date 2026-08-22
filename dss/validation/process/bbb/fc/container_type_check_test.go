package fc

import (
	"testing"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// testMultiValuesRule is a minimal policy.MultiValuesRule for tests.
type testMultiValuesRule struct {
	level  enumerations.Level
	values []string
}

func (r testMultiValuesRule) Level() enumerations.Level { return r.level }
func (r testMultiValuesRule) Values() []string          { return r.values }

func newTestFCResult() *process.Result[*drjaxb.XmlFC] {
	xmlFC := &drjaxb.XmlFC{}
	return process.NewResult(xmlFC, &xmlFC.XmlConstraintsConclusionContent, &xmlFC.XmlConstraintsConclusionAttrs)
}

func testI18nProvider(t *testing.T) *i18n.Provider {
	t.Helper()
	return i18n.NewProvider()
}

func TestContainerTypeCheck_Process(t *testing.T) {
	tests := []struct {
		name          string
		containerType enumerations.ASiCContainerType
		accepted      []string
		want          bool
	}{
		// The policy ids carry ASiCContainerType#toString(), which replaces '_' with '-'
		// ("ASiC-E"), not the enum name - see ContainerTypeCheck.Process.
		{"happy: accepted container type", enumerations.ASiCContainerTypeASiCE, []string{"ASiC-E", "ASiC-S"}, true},
		{"failure: not in accepted list", enumerations.ASiCContainerTypeASiCS, []string{"ASiC-E"}, false},
		{"failure: enum-name spelling is not accepted", enumerations.ASiCContainerTypeASiCE, []string{"ASiC_E"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContainerTypeCheck(testI18nProvider(t), newTestFCResult(), tt.containerType,
				testMultiValuesRule{level: enumerations.LevelFail, values: tt.accepted})
			if got := c.Process(); got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
