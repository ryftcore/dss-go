// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/client/http/DSSFileLoader.java (DSS 6.5.RC1).
package http

import "github.com/utain/esig/dss/model"

// DSSFileLoader loads a model.DSSDocument instead of raw binaries.
type DSSFileLoader interface {
	// GetDocument returns the DSSDocument from the provided url. Returns a
	// *model.DSSError in case of DataLoader error, mirroring the Java
	// method's DSSException.
	GetDocument(url string) (model.DSSDocument, error)
}
