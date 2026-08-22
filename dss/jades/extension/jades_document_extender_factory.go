// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/extension/JAdESDocumentExtenderFactory.java (DSS 6.5.RC1).
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

// JAdESDocumentExtenderFactory is used to check and load a corresponding DocumentExtender
// implementation for a JAdES signature or signatures augmentation.
type JAdESDocumentExtenderFactory struct{}

// NewJAdESDocumentExtenderFactory is the default constructor.
func NewJAdESDocumentExtenderFactory() *JAdESDocumentExtenderFactory {
	return &JAdESDocumentExtenderFactory{}
}

// IsSupported ports the overridden isSupported(DSSDocument).
func (f *JAdESDocumentExtenderFactory) IsSupported(doc model.DSSDocument) bool {
	return newJAdESDocumentExtender().IsSupported(doc)
}

// Create ports the overridden create(DSSDocument).
func (f *JAdESDocumentExtenderFactory) Create(doc model.DSSDocument) document.SignedDocumentExtender {
	return NewJAdESDocumentExtender(doc)
}

// init registers the factory with dss-document's SignedDocumentExtender registry, standing in
// for the META-INF/services/eu.europa.esig.dss.extension.SignedDocumentExtenderFactory entry.
func init() {
	document.RegisterSignedDocumentExtenderFactory(NewJAdESDocumentExtenderFactory())
}

// compile-time assertion.
var _ document.SignedDocumentExtenderFactory = (*JAdESDocumentExtenderFactory)(nil)
