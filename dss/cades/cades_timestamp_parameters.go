// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESTimestampParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// CAdESTimestampParameters defines TimestampParameters to deal with CAdES timestamp creation.
type CAdESTimestampParameters struct {
	model.TimestampParameters
}

var _ model.SerializableTimestampParameters = (*CAdESTimestampParameters)(nil)

// NewCAdESTimestampParameters instantiates the object with the default digest algorithm. Port of
// the empty constructor.
func NewCAdESTimestampParameters() *CAdESTimestampParameters {
	return &CAdESTimestampParameters{TimestampParameters: model.NewTimestampParameters()}
}

// NewCAdESTimestampParametersWithDigestAlgorithm instantiates the object with the given digest
// algorithm to use for timestamping data. Port of CAdESTimestampParameters(DigestAlgorithm).
func NewCAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *CAdESTimestampParameters {
	return &CAdESTimestampParameters{TimestampParameters: model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm)}
}

// String ports #toString.
func (p *CAdESTimestampParameters) String() string {
	return fmt.Sprintf("CAdESTimestampParameters %s", p.TimestampParameters.String())
}
