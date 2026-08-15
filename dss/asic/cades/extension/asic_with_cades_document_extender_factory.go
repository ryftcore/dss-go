// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/extension/ASiCWithCAdESDocumentExtenderFactory.java (DSS 6.5.RC1).
//
// Upstream registers this factory through a META-INF/services entry that
// SignedDocumentExtender#fromDocument's ServiceLoader picks up. Go has no ServiceLoader, so
// dss-document replaced it with an explicit registry; the init() below is this package's
// registration, i.e. the exact counterpart of that services file.
package extension

import (
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/model"
)

// ASiCWithCAdESDocumentExtenderFactory is used to check and load a corresponding DocumentExtender
// implementation for a CAdES signature or signatures augmentation within an ASiC container.
type ASiCWithCAdESDocumentExtenderFactory struct{}

// NewASiCWithCAdESDocumentExtenderFactory is the default constructor.
func NewASiCWithCAdESDocumentExtenderFactory() *ASiCWithCAdESDocumentExtenderFactory {
	return &ASiCWithCAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *ASiCWithCAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newASiCWithCAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *ASiCWithCAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewASiCWithCAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewASiCWithCAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*ASiCWithCAdESDocumentExtenderFactory)(nil)
