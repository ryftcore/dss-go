// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/RevocationRefWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"fmt"
	"math/big"
	"time"

	"github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
)

// RevocationRefWrapper represents a revocation data reference wrapper.
type RevocationRefWrapper struct {
	// revocationRef is the wrapped XmlRevocationRef.
	revocationRef *jaxb.XmlRevocationRef

	// revocationId is the Id of the related revocation token.
	revocationId string
}

// NewRevocationRefWrapper is the default constructor.
func NewRevocationRefWrapper(revocationRef *jaxb.XmlRevocationRef, revocationId string) *RevocationRefWrapper {
	return &RevocationRefWrapper{revocationRef: revocationRef, revocationId: revocationId}
}

// Origins returns a list of revocation reference origins. Port of getOrigins().
func (w *RevocationRefWrapper) Origins() []enumerations.RevocationRefOrigin {
	return w.revocationRef.Origins
}

// Issuer gets the CRL issuer RDN. NOTE: applicable only for CRL references. Port of
// getIssuer().
func (w *RevocationRefWrapper) Issuer() string {
	return w.revocationRef.Issuer
}

// IssueTime gets the issue time of the CRL. NOTE: applicable only for CRL references. Port
// of getIssueTime().
func (w *RevocationRefWrapper) IssueTime() *time.Time {
	return w.revocationRef.IssueTime
}

// CRLNumber gets the number of the CRL. NOTE: applicable only for CRL references. Port of
// getCRLNumber().
func (w *RevocationRefWrapper) CRLNumber() *big.Int {
	return w.revocationRef.CRLNumber
}

// ProductionTime returns revocation ref production time if present. NOTE: applicable only
// for OCSP response references. Port of getProductionTime().
func (w *RevocationRefWrapper) ProductionTime() *time.Time {
	return w.revocationRef.ProducedAt
}

// ResponderIdName returns responder's ID name if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdName().
func (w *RevocationRefWrapper) ResponderIdName() string {
	if w.revocationRef.ResponderId != nil {
		return w.revocationRef.ResponderId.IssuerName
	}
	return ""
}

// ResponderIdKey returns responder's ID key if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdKey().
func (w *RevocationRefWrapper) ResponderIdKey() []byte {
	if w.revocationRef.ResponderId != nil {
		return w.revocationRef.ResponderId.Ski
	}
	return nil
}

// Uri gets the URI reference to the revocation data, when present. Port of getUri().
func (w *RevocationRefWrapper) Uri() string {
	return w.revocationRef.Uri
}

// DigestAlgoAndValue returns digest algo and value. Port of getDigestAlgoAndValue().
func (w *RevocationRefWrapper) DigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.revocationRef.DigestAlgoAndValue
}

// RevocationId returns an Id of the related revocation token, when present. Returns Id of
// the reference otherwise. Port of getRevocationId().
func (w *RevocationRefWrapper) RevocationId() string {
	return w.revocationId
}

// String returns a string representation of the wrapper. Port of toString().
func (w *RevocationRefWrapper) String() string {
	if w.revocationRef != nil {
		responderIdName := ""
		if w.revocationRef.ResponderId != nil {
			responderIdName = w.revocationRef.ResponderId.IssuerName
		}
		return fmt.Sprintf("RevocationRefWrapper Origins='%v',  ProductionTime='%v', responderIdName='%s'",
			w.revocationRef.Origins, w.revocationRef.ProducedAt, responderIdName)
	}
	return "RevocationRefWrapper revocationRef=<nil>"
}
