// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/evidencerecord/digest/DataObjectDigestBuilder.java (DSS 6.5.RC1).
package digest

import "github.com/ryftcore/dss-go/dss/model"

// DataObjectDigestBuilder is a common interface for the classes providing a
// functionality build digest for data objects to be protected by an
// evidence record preservation service.
type DataObjectDigestBuilder interface {
	// Build generates a hash value. Returns a Digest containing the hash
	// value of the data object and the used digest algorithm. Ports
	// #build.
	Build() *model.Digest
}
