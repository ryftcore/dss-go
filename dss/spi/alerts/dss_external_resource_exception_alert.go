// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/alerts/DSSExternalResourceExceptionAlert.java (DSS 6.5.RC1).
package alerts

import (
	"github.com/ryftcore/dss-go/dss/alert"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// dssExternalResourceExceptionAlertHandler reports a DSSExternalResourceException built from
// the Status' error string. Port of the lambda passed to AbstractStatusAlert's constructor.
type dssExternalResourceExceptionAlertHandler struct{}

// Process implements alert.Handler[alert.Status].
func (dssExternalResourceExceptionAlertHandler) Process(object alert.Status) error {
	return exception.NewDSSExternalResourceException(object.ErrorString())
}

// DSSExternalResourceExceptionAlert reports a
// eu.europa.esig.dss.spi.exception.DSSExternalResourceException when the corresponding check
// fails. Ports DSSExternalResourceExceptionAlert.
type DSSExternalResourceExceptionAlert struct {
	*alert.AbstractStatusAlert
}

// NewDSSExternalResourceExceptionAlert creates a DSSExternalResourceExceptionAlert. Port of the
// default constructor.
func NewDSSExternalResourceExceptionAlert() *DSSExternalResourceExceptionAlert {
	return &DSSExternalResourceExceptionAlert{
		AbstractStatusAlert: alert.NewAbstractStatusAlert(dssExternalResourceExceptionAlertHandler{}),
	}
}
