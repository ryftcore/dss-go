// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/ASiCWithXAdESSignatureParameters.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/asic"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
)

// ASiCWithXAdESSignatureParameters defines SignatureParameters to deal with ASiC with XAdES
// signature creation/extension.
type ASiCWithXAdESSignatureParameters struct {
	dssxades.SignatureParameters

	// asicParams is the object representing the parameters related to ASiC from of the
	// signature.
	asicParams *asic.Parameters
}

// NewASiCWithXAdESSignatureParameters instantiates object with default ASiCParameters. Port of
// the default constructor.
func NewASiCWithXAdESSignatureParameters() *ASiCWithXAdESSignatureParameters {
	return &ASiCWithXAdESSignatureParameters{
		SignatureParameters: *dssxades.NewSignatureParameters(),
		asicParams:          asic.NewParameters(),
	}
}

// ASiC ports the @Override aSiC().
func (p *ASiCWithXAdESSignatureParameters) ASiC() *asic.Parameters {
	return p.asicParams
}

// String ports #toString.
func (p *ASiCWithXAdESSignatureParameters) String() string {
	return fmt.Sprintf("ASiCWithXAdESSignatureParameters [asicParams=%v] %s",
		p.asicParams, p.SignatureParameters.String())
}

// Equals ports #equals.
func (p *ASiCWithXAdESSignatureParameters) Equals(other *ASiCWithXAdESSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.SignatureParameters.Equals(&other.SignatureParameters) {
		return false
	}
	return asicWithXAdESSignatureParametersASiCParamsEquals(p.asicParams, other.asicParams)
}

// asicWithXAdESSignatureParametersASiCParamsEquals ports the Objects.equals(asicParams,
// that.asicParams) comparison, local to this file per PORTING.md (no cross-file shared helpers).
func asicWithXAdESSignatureParametersASiCParamsEquals(a, b *asic.Parameters) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.IsZipComment() == b.IsZipComment() && a.MimeType() == b.MimeType() && a.ContainerType() == b.ContainerType()
}
