// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/ASiCWithCAdESSignatureParameters.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package cades

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCWithCAdESSignatureParameters defines SignatureParameters to deal with ASiC with CAdES
// signature creation/extension.
type ASiCWithCAdESSignatureParameters struct {
	dsscades.SignatureParameters

	// asicParams is the object representing the parameters related to ASiC for the signature.
	asicParams *asic.Parameters
}

var _ ASiCWithCAdESCommonParameters = (*ASiCWithCAdESSignatureParameters)(nil)

// NewASiCWithCAdESSignatureParameters instantiates object with default ASiCParameters. Port of
// the default constructor.
func NewASiCWithCAdESSignatureParameters() *ASiCWithCAdESSignatureParameters {
	return &ASiCWithCAdESSignatureParameters{
		SignatureParameters: *dsscades.NewSignatureParameters(),
		asicParams:          asic.NewParameters(),
	}
}

// ASiC ports the @Override aSiC().
func (p *ASiCWithCAdESSignatureParameters) ASiC() *asic.Parameters {
	return p.asicParams
}

// SetSignatureLevel ports the @Override setSignatureLevel(SignatureLevel).
//
// Panics with the Java message when signatureLevel is nil or not of the CAdES form
// (IllegalArgumentException).
func (p *ASiCWithCAdESSignatureParameters) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) {
	form, err := signatureLevel.SignatureForm()
	if signatureLevel == "" || err != nil || enumerations.SignatureFormCAdES != form {
		panic("Only CAdES form is allowed !")
	}
	p.SignatureParameters.SetSignatureLevel(signatureLevel)
}

// ZipCreationDate ports the @Override getZipCreationDate().
func (p *ASiCWithCAdESSignatureParameters) ZipCreationDate() time.Time {
	if signingDate := p.BLevel().SigningDate(); signingDate != nil {
		return *signingDate
	}
	return time.Time{}
}

// String ports #toString.
func (p *ASiCWithCAdESSignatureParameters) String() string {
	return fmt.Sprintf("ASiCWithCAdESSignatureParameters [asicParams=%v] %s",
		p.asicParams, p.SignatureParameters.String())
}

// Equals ports #equals.
func (p *ASiCWithCAdESSignatureParameters) Equals(other *ASiCWithCAdESSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.SignatureParameters.Equals(&other.SignatureParameters) {
		return false
	}
	return asicWithCAdESSignatureParametersASiCParamsEquals(p.asicParams, other.asicParams)
}

// asicWithCAdESSignatureParametersASiCParamsEquals ports the Objects.equals(asicParams,
// that.asicParams) comparison, local to this file per PORTING.md (no cross-file shared helpers).
func asicWithCAdESSignatureParametersASiCParamsEquals(a, b *asic.Parameters) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.IsZipComment() == b.IsZipComment() && a.MimeType() == b.MimeType() && a.ContainerType() == b.ContainerType()
}
