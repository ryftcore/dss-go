// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESASiCContentBuilder.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithXAdESASiCContentBuilder builds an Content for an ASiC with XAdES container.
//
// Java's `extends AbstractASiCContentBuilder` becomes embedding plus the
// InitAbstractASiCContentBuilder(self) registration: Go has no method overriding across
// embedding, so the base dispatches getContainerExtractor through
// asic.AbstractContentBuilderOverrides.
type ASiCWithXAdESASiCContentBuilder struct {
	*asic.AbstractContentBuilder
}

var _ asic.AbstractContentBuilderOverrides = (*ASiCWithXAdESASiCContentBuilder)(nil)

// NewASiCWithXAdESASiCContentBuilder is the default constructor. Ports the empty
// ASiCWithXAdESASiCContentBuilder().
func NewASiCWithXAdESASiCContentBuilder() *ASiCWithXAdESASiCContentBuilder {
	builder := &ASiCWithXAdESASiCContentBuilder{
		AbstractContentBuilder: asic.NewAbstractASiCContentBuilderBase(),
	}
	builder.InitAbstractASiCContentBuilder(builder)
	return builder
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(DSSDocument).
func (b *ASiCWithXAdESASiCContentBuilder) GetContainerExtractor(archiveDocument model.DSSDocument) asic.ContainerExtractor {
	return NewASiCWithXAdESContainerExtractor(archiveDocument)
}
