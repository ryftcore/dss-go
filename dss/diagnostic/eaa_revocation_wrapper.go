// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/EAARevocationWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// EAARevocationWrapper contains information about the validity of an EAA.
type EAARevocationWrapper struct {
	EAARevocationTokenWrapper

	// xmlEAARevocationStatus is the wrapped XmlEAARevocationStatus.
	xmlEAARevocationStatus *jaxb.XmlEAARevocationStatus
}

// NewEAARevocationWrapper is the default constructor.
func NewEAARevocationWrapper(xmlEAARevocationStatus *jaxb.XmlEAARevocationStatus) *EAARevocationWrapper {
	w := &EAARevocationWrapper{
		EAARevocationTokenWrapper: *NewEAARevocationTokenWrapper(xmlEAARevocationStatus.EAARevocationToken),
		xmlEAARevocationStatus:    xmlEAARevocationStatus,
	}
	w.InitTokenProxy(w)
	return w
}

// Status returns the status of the concerned EAA. Port of getStatus().
func (w *EAARevocationWrapper) Status() enumerations.EAAStatus {
	if w.xmlEAARevocationStatus.Status != nil {
		return enumerations.EAAStatus(*w.xmlEAARevocationStatus.Status)
	}
	return ""
}
