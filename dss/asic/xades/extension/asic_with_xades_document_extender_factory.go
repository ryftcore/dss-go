// Ported from dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/extension/ASiCWithXAdESDocumentExtenderFactory.java (DSS 6.5.RC1).
//
// Upstream registers this factory through a META-INF/services entry that
// SignedDocumentExtender#fromDocument's ServiceLoader picks up. Go has no ServiceLoader, so
// dss-document replaced it with an explicit registry; the init() below is this package's
// registration, i.e. the exact counterpart of that services file. Mirrors the
// asic/cades/extension/asic_with_cades_document_extender_factory.go precedent.
package extension

import (
	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/model"
)

// ASiCWithXAdESDocumentExtenderFactory is used to check and load a corresponding
// DocumentExtender implementation for a CAdES signature or signatures augmentation within an
// ASiC container.
type ASiCWithXAdESDocumentExtenderFactory struct{}

// NewASiCWithXAdESDocumentExtenderFactory is the default constructor.
func NewASiCWithXAdESDocumentExtenderFactory() *ASiCWithXAdESDocumentExtenderFactory {
	return &ASiCWithXAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *ASiCWithXAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newASiCWithXAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *ASiCWithXAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewASiCWithXAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewASiCWithXAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*ASiCWithXAdESDocumentExtenderFactory)(nil)
