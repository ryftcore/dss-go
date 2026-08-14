// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/extension/CAdESDocumentExtenderFactory.java (DSS 6.5.RC1).
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

// CAdESDocumentExtenderFactory is used to check and load a corresponding DocumentExtender
// implementation for a CAdES signature or signatures augmentation.
type CAdESDocumentExtenderFactory struct{}

// NewCAdESDocumentExtenderFactory is the default constructor.
func NewCAdESDocumentExtenderFactory() *CAdESDocumentExtenderFactory {
	return &CAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *CAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newCAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *CAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewCAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewCAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*CAdESDocumentExtenderFactory)(nil)
