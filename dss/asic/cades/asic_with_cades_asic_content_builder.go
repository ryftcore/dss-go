// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESASiCContentBuilder.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithCAdESASiCContentBuilder builds an Content for an ASiC with CAdES container.
//
// Java's `extends AbstractASiCContentBuilder` becomes embedding plus the
// InitAbstractASiCContentBuilder(self) registration: Go has no method overriding across
// embedding, so the base dispatches getContainerExtractor through
// asic.AbstractContentBuilderOverrides.
type ASiCWithCAdESASiCContentBuilder struct {
	*asic.AbstractContentBuilder
}

var _ asic.AbstractContentBuilderOverrides = (*ASiCWithCAdESASiCContentBuilder)(nil)

// NewASiCWithCAdESASiCContentBuilder is the default constructor. Ports the empty
// ASiCWithCAdESASiCContentBuilder().
func NewASiCWithCAdESASiCContentBuilder() *ASiCWithCAdESASiCContentBuilder {
	builder := &ASiCWithCAdESASiCContentBuilder{
		AbstractContentBuilder: asic.NewAbstractASiCContentBuilderBase(),
	}
	builder.InitAbstractASiCContentBuilder(builder)
	return builder
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(DSSDocument).
func (b *ASiCWithCAdESASiCContentBuilder) GetContainerExtractor(archiveDocument model.DSSDocument) asic.ContainerExtractor {
	return NewASiCWithCAdESContainerExtractor(archiveDocument)
}
