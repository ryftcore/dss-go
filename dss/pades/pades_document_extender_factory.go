// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/extension/PAdESDocumentExtenderFactory.java
// (DSS 6.5.RC1).
//
// Upstream registers this factory through a META-INF/services entry that
// SignedDocumentExtender#fromDocument's ServiceLoader picks up. Go has no ServiceLoader, so
// dss-document replaced it with an explicit registry; the init() below is this package's
// registration, i.e. the exact counterpart of that services file, per the
// cades/extension/cades_document_extender_factory.go precedent.
package pades

import (
	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/model"
)

// PAdESDocumentExtenderFactory is used to check and load a corresponding DocumentExtender
// implementation for a PAdES signature or signatures augmentation.
type PAdESDocumentExtenderFactory struct{}

// NewPAdESDocumentExtenderFactory is the default constructor.
func NewPAdESDocumentExtenderFactory() *PAdESDocumentExtenderFactory {
	return &PAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *PAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newPAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *PAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewPAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewPAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*PAdESDocumentExtenderFactory)(nil)
