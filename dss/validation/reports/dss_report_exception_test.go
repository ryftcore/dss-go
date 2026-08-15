// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/reports/DSSReportException.java
// (DSS 6.5.RC1) — behavior test vectors for the four constructors.

package reports

import (
	"errors"
	"testing"
)

func TestDSSReportException(t *testing.T) {
	if got := NewDSSReportException().Error(); got != "DSSReportException" {
		t.Errorf("empty constructor: got %q", got)
	}
	if got := NewDSSReportExceptionMessage("boom").Error(); got != "boom" {
		t.Errorf("message constructor: got %q", got)
	}
	cause := errors.New("root cause")
	if got := NewDSSReportExceptionCause(cause).Error(); got != "root cause" {
		t.Errorf("cause constructor: got %q", got)
	}
	e := NewDSSReportExceptionMessageCause("boom", cause)
	if got := e.Error(); got != "boom: root cause" {
		t.Errorf("message+cause constructor: got %q", got)
	}
	if !errors.Is(e, cause) {
		t.Errorf("errors.Is should unwrap to cause")
	}
}
