// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/signature/CAdESTimestampParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TimestampParameters defines TimestampParameters to deal with CAdES timestamp creation.
type TimestampParameters struct {
	model.TimestampParameters
}

var _ model.SerializableTimestampParameters = (*TimestampParameters)(nil)

// NewCAdESTimestampParameters instantiates the object with the default digest algorithm. Port of
// the empty constructor.
func NewTimestampParameters() *TimestampParameters {
	return &TimestampParameters{TimestampParameters: model.NewTimestampParameters()}
}

// NewTimestampParametersWithDigestAlgorithm instantiates the object with the given digest
// algorithm to use for timestamping data. Port of CAdESTimestampParameters(DigestAlgorithm).
func NewTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *TimestampParameters {
	return &TimestampParameters{TimestampParameters: model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm)}
}

// String ports #toString.
func (p *TimestampParameters) String() string {
	return fmt.Sprintf("CAdESTimestampParameters %s", p.TimestampParameters.String())
}
