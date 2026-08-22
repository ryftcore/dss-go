// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/ASiCWithCAdESTimestampParameters.java (DSS 6.5.RC1).
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

// ASiCWithCAdESTimestampParameters defines TimestampParameters to deal with ASiC with CAdES
// timestamp creation.
type ASiCWithCAdESTimestampParameters struct {
	dsscades.TimestampParameters

	// zipCreationDate is used to set the DateTime for created ZIP entries.
	zipCreationDate time.Time

	// asicParams is the object representing the parameters related to ASiC for the timestamp.
	asicParams *asic.Parameters
}

var _ ASiCWithCAdESCommonParameters = (*ASiCWithCAdESTimestampParameters)(nil)

// ASiC ports the @Override aSiC().
func (p *ASiCWithCAdESTimestampParameters) ASiC() *asic.Parameters {
	return p.asicParams
}

// NewASiCWithCAdESTimestampParameters is the empty constructor. Port of the empty constructor.
func NewASiCWithCAdESTimestampParameters() *ASiCWithCAdESTimestampParameters {
	return &ASiCWithCAdESTimestampParameters{
		TimestampParameters: *dsscades.NewCAdESTimestampParameters(),
		zipCreationDate:     time.Now(),
		asicParams:          asic.NewASiCParameters(),
	}
}

// NewASiCWithCAdESTimestampParametersWithDigestAlgorithm is the constructor defining a
// DigestAlgorithm. Port of ASiCWithCAdESTimestampParameters(DigestAlgorithm).
func NewASiCWithCAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *ASiCWithCAdESTimestampParameters {
	return &ASiCWithCAdESTimestampParameters{
		TimestampParameters: *dsscades.NewCAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm),
		zipCreationDate:     time.Now(),
		asicParams:          asic.NewASiCParameters(),
	}
}

// NewASiCWithCAdESTimestampParametersWithDigestAlgorithmAndASiCParams is the constructor defining
// a DigestAlgorithm and ASiCParameters. Port of ASiCWithCAdESTimestampParameters(DigestAlgorithm,
// Parameters).
func NewASiCWithCAdESTimestampParametersWithDigestAlgorithmAndASiCParams(digestAlgorithm enumerations.DigestAlgorithm, asicParams *asic.Parameters) *ASiCWithCAdESTimestampParameters {
	return &ASiCWithCAdESTimestampParameters{
		TimestampParameters: *dsscades.NewCAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm),
		zipCreationDate:     time.Now(),
		asicParams:          asicParams,
	}
}

// ZipCreationDate ports the @Override getZipCreationDate().
func (p *ASiCWithCAdESTimestampParameters) ZipCreationDate() time.Time {
	return p.zipCreationDate
}

// SetZipCreationDate sets ZIP creation date, used to define a creation date for ZIP container
// entries. Port of setZipCreationDate(Date).
func (p *ASiCWithCAdESTimestampParameters) SetZipCreationDate(zipCreationDate time.Time) {
	p.zipCreationDate = zipCreationDate
}

// String ports #toString.
func (p *ASiCWithCAdESTimestampParameters) String() string {
	return fmt.Sprintf("ASiCWithCAdESTimestampParameters [zipCreationDate=%v, asicParams=%v] %s",
		p.zipCreationDate, p.asicParams, p.TimestampParameters.String())
}

// Equals ports #equals.
func (p *ASiCWithCAdESTimestampParameters) Equals(other *ASiCWithCAdESTimestampParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	return p.zipCreationDate.Equal(other.zipCreationDate) &&
		asicWithCAdESSignatureParametersASiCParamsEquals(p.asicParams, other.asicParams)
}
