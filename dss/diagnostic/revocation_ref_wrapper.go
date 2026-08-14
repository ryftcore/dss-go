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

// GetOrigins returns a list of revocation reference origins. Port of getOrigins().
func (w *RevocationRefWrapper) GetOrigins() []enumerations.RevocationRefOrigin {
	return w.revocationRef.Origins
}

// GetIssuer gets the CRL issuer RDN. NOTE: applicable only for CRL references. Port of
// getIssuer().
func (w *RevocationRefWrapper) GetIssuer() string {
	return w.revocationRef.Issuer
}

// GetIssueTime gets the issue time of the CRL. NOTE: applicable only for CRL references. Port
// of getIssueTime().
func (w *RevocationRefWrapper) GetIssueTime() *time.Time {
	return w.revocationRef.IssueTime
}

// GetCRLNumber gets the number of the CRL. NOTE: applicable only for CRL references. Port of
// getCRLNumber().
func (w *RevocationRefWrapper) GetCRLNumber() *big.Int {
	return w.revocationRef.CRLNumber
}

// GetProductionTime returns revocation ref production time if present. NOTE: applicable only
// for OCSP response references. Port of getProductionTime().
func (w *RevocationRefWrapper) GetProductionTime() *time.Time {
	return w.revocationRef.ProducedAt
}

// GetResponderIdName returns responder's ID name if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdName().
func (w *RevocationRefWrapper) GetResponderIdName() string {
	if w.revocationRef.ResponderId != nil {
		return w.revocationRef.ResponderId.IssuerName
	}
	return ""
}

// GetResponderIdKey returns responder's ID key if present. NOTE: applicable only for OCSP
// response references. Port of getResponderIdKey().
func (w *RevocationRefWrapper) GetResponderIdKey() []byte {
	if w.revocationRef.ResponderId != nil {
		return w.revocationRef.ResponderId.Ski
	}
	return nil
}

// GetUri gets the URI reference to the revocation data, when present. Port of getUri().
func (w *RevocationRefWrapper) GetUri() string {
	return w.revocationRef.Uri
}

// GetDigestAlgoAndValue returns digest algo and value. Port of getDigestAlgoAndValue().
func (w *RevocationRefWrapper) GetDigestAlgoAndValue() *jaxb.XmlDigestAlgoAndValue {
	return w.revocationRef.DigestAlgoAndValue
}

// GetRevocationId returns an Id of the related revocation token, when present. Returns Id of
// the reference otherwise. Port of getRevocationId().
func (w *RevocationRefWrapper) GetRevocationId() string {
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
