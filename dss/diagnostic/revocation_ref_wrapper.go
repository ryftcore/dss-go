// Ported from dss-diagnostic-jaxb/src/main/java/eu/europa/esig/dss/diagnostic/RevocationRefWrapper.java (DSS 6.5.RC1).
package diagnostic

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
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
	values := w.revocationRef.Origin
	if values == nil {
		return nil
	}
	result := make([]enumerations.RevocationRefOrigin, len(values))
	for i, v := range values {
		result[i] = enumerations.RevocationRefOrigin(v)
	}
	return result
}

// Issuer gets the CRL issuer RDN. NOTE: applicable only for CRL references. Port of
// getIssuer().
func (w *RevocationRefWrapper) Issuer() string {
	if w.revocationRef.Issuer != nil {
		return *w.revocationRef.Issuer
	}
	return ""
}

// IssueTime gets the issue time of the CRL. NOTE: applicable only for CRL references. Port
// of getIssueTime().
func (w *RevocationRefWrapper) IssueTime() *time.Time {
	if w.revocationRef.IssueTime == nil {
		return nil
	}
	t := w.revocationRef.IssueTime.Time()
	return &t
}

// CRLNumber gets the number of the CRL. NOTE: applicable only for CRL references. Port of
// getCRLNumber().
func (w *RevocationRefWrapper) CRLNumber() *big.Int {
	return w.revocationRef.CRLNumber.BigInt()
}

// ProductionTime returns revocation ref production time if present. NOTE: applicable only
// for OCSP response references. Port of getProductionTime().
func (w *RevocationRefWrapper) ProductionTime() *time.Time {
	if w.revocationRef.ProducedAt == nil {
		return nil
	}
	t := w.revocationRef.ProducedAt.Time()
	return &t
}

// ResponderIdName returns responder's ID name if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdName().
func (w *RevocationRefWrapper) ResponderIdName() string {
	if w.revocationRef.ResponderId != nil && w.revocationRef.ResponderId.IssuerName != nil {
		return *w.revocationRef.ResponderId.IssuerName
	}
	return ""
}

// ResponderIdKey returns responder's ID key if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdKey().
func (w *RevocationRefWrapper) ResponderIdKey() []byte {
	if w.revocationRef.ResponderId != nil && w.revocationRef.ResponderId.Ski != nil {
		return []byte(*w.revocationRef.ResponderId.Ski)
	}
	return nil
}

// Uri gets the URI reference to the revocation data, when present. Port of getUri().
func (w *RevocationRefWrapper) Uri() string {
	if w.revocationRef.Uri != nil {
		return *w.revocationRef.Uri
	}
	return ""
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
		if w.revocationRef.ResponderId != nil && w.revocationRef.ResponderId.IssuerName != nil {
			responderIdName = *w.revocationRef.ResponderId.IssuerName
		}
		return fmt.Sprintf("RevocationRefWrapper Origins='%v',  ProductionTime='%v', responderIdName='%s'",
			w.Origins(), w.ProductionTime(), responderIdName)
	}
	return "RevocationRefWrapper revocationRef=null"
}
