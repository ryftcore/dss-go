package fc

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

func TestZipCommentPresentCheck_Process(t *testing.T) {
	tests := []struct {
		name       string
		zipComment string
		want       bool
	}{
		{"happy: comment present", "some comment", true},
		{"failure: comment blank", "   ", false},
		{"failure: comment empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewZipCommentPresentCheck(testI18nProvider(t), newTestFCResult(), tt.zipComment,
				process.GetLevelRule(enumerations.Level_FAIL))
			if got := c.Process(); got != tt.want {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
