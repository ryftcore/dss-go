// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/extension/XAdESDocumentExtenderFactory.java (DSS 6.5.RC1).
//
// Upstream registers this factory through a META-INF/services entry that
// SignedDocumentExtender#fromDocument's ServiceLoader picks up. Go has no ServiceLoader, so
// dss-document replaced it with an explicit registry; the init() below is this package's
// registration, i.e. the exact counterpart of that services file.
package extension

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/model"
)

// XAdESDocumentExtenderFactory is used to check and load a corresponding DocumentExtender
// implementation for a XAdES signature or signatures augmentation.
type XAdESDocumentExtenderFactory struct{}

// NewXAdESDocumentExtenderFactory is the default constructor.
func NewXAdESDocumentExtenderFactory() *XAdESDocumentExtenderFactory {
	return &XAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *XAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newXAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *XAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewXAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewXAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*XAdESDocumentExtenderFactory)(nil)
