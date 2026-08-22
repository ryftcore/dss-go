// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TSPSource.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// TSPSource is the abstraction of a Time Stamping authority which delivers RFC 3161 Time Stamp
// Responses containing tokens, from Time Stamp Requests.
//
// java.io.Serializable is dropped (no Go counterpart).
type TSPSource interface {
	// TimeStampResponse gets a TimeStampResponse relevant to the provided digest.
	// digestAlgorithm is the used digest algorithm and digest the computed digest to be
	// timestamped; the returned TimestampBinary holds the binaries of a signed timestamp token.
	// Port of getTimeStampResponse(DigestAlgorithm, byte[]), whose DSSException becomes the
	// returned error.
	TimeStampResponse(digestAlgorithm enumerations.DigestAlgorithm, digest []byte) (*model.TimestampBinary, error)
}
