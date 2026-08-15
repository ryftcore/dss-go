// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/ASiCWithXAdESSignatureParameters.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/asic"
	dssxades "github.com/utain/esig/dss/xades"
)

// ASiCWithXAdESSignatureParameters defines SignatureParameters to deal with ASiC with XAdES
// signature creation/extension.
type ASiCWithXAdESSignatureParameters struct {
	dssxades.XAdESSignatureParameters

	// asicParams is the object representing the parameters related to ASiC from of the
	// signature.
	asicParams *asic.ASiCParameters
}

// NewASiCWithXAdESSignatureParameters instantiates object with default ASiCParameters. Port of
// the default constructor.
func NewASiCWithXAdESSignatureParameters() *ASiCWithXAdESSignatureParameters {
	return &ASiCWithXAdESSignatureParameters{
		XAdESSignatureParameters: *dssxades.NewXAdESSignatureParameters(),
		asicParams:               asic.NewASiCParameters(),
	}
}

// ASiC ports the @Override aSiC().
func (p *ASiCWithXAdESSignatureParameters) ASiC() *asic.ASiCParameters {
	return p.asicParams
}

// String ports #toString.
func (p *ASiCWithXAdESSignatureParameters) String() string {
	return fmt.Sprintf("ASiCWithXAdESSignatureParameters [asicParams=%v] %s",
		p.asicParams, p.XAdESSignatureParameters.String())
}

// Equals ports #equals.
func (p *ASiCWithXAdESSignatureParameters) Equals(other *ASiCWithXAdESSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.XAdESSignatureParameters.Equals(&other.XAdESSignatureParameters) {
		return false
	}
	return asicWithXAdESSignatureParametersASiCParamsEquals(p.asicParams, other.asicParams)
}

// asicWithXAdESSignatureParametersASiCParamsEquals ports the Objects.equals(asicParams,
// that.asicParams) comparison, local to this file per PORTING.md (no cross-file shared helpers).
func asicWithXAdESSignatureParametersASiCParamsEquals(a, b *asic.ASiCParameters) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.IsZipComment() == b.IsZipComment() && a.MimeType() == b.MimeType() && a.ContainerType() == b.ContainerType()
}
