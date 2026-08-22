// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/ASiCWithCAdESCommonParameters.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped silently (no Go counterpart).
package cades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// ASiCWithCAdESCommonParameters defines common parameters for an ASiC with CAdES container for
// signature/timestamp creation.
type ASiCWithCAdESCommonParameters interface {
	// ASiC returns ASiC container parameters. Port of aSiC().
	ASiC() *asic.ASiCParameters

	// DigestAlgorithm returns a DigestAlgorithm to be used to hash a data to be timestamped.
	// Port of getDigestAlgorithm().
	DigestAlgorithm() enumerations.DigestAlgorithm

	// ZipCreationDate returns a signing date. Port of getZipCreationDate().
	ZipCreationDate() time.Time
}
