// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/evidencerecord/digest/DataObjectDigestBuilderFactory.java (DSS 6.5.RC1).
package digest

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// DataObjectDigestBuilderFactory creates an instance of
// DataObjectDigestBuilder.
type DataObjectDigestBuilderFactory interface {
	// Create creates an instance of DataObjectDigestBuilder to build hash
	// for the document, according to the given implementation, using a
	// default digest algorithm. Ports #create(DSSDocument).
	Create(document model.DSSDocument) DataObjectDigestBuilder

	// CreateWithAlgorithm creates an instance of DataObjectDigestBuilder
	// to build hash for the document, according to the given
	// implementation, using a provided digestAlgorithm. Ports
	// #create(DSSDocument, DigestAlgorithm).
	CreateWithAlgorithm(document model.DSSDocument, digestAlgorithm enumerations.DigestAlgorithm) DataObjectDigestBuilder
}
