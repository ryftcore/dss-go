// Ported from dss-crl-parser/src/main/java/eu/europa/esig/dss/crl/CRLBinary.java (DSS 6.5.RC1).
package crlparser

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
)

// CRLBinary represents a DER encoded CRL Binary identifier.
type CRLBinary struct {
	model.EncapsulatedRevocationTokenIdentifier[revocation.CRL]
}

// NewCRLBinary builds a CRLBinary over the given DER encoded binaries.
// Port of the CRLBinary(byte[]) constructor.
func NewCRLBinary(derEncoded []byte) *CRLBinary {
	return &CRLBinary{
		EncapsulatedRevocationTokenIdentifier: *model.NewEncapsulatedRevocationTokenIdentifierWithClassName[revocation.CRL](
			"CRLBinary", derEncoded),
	}
}
