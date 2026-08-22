// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESASiCContentBuilder.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
)

// ASiCWithXAdESASiCContentBuilder builds an ASiCContent for an ASiC with XAdES container.
//
// Java's `extends AbstractASiCContentBuilder` becomes embedding plus the
// InitAbstractASiCContentBuilder(self) registration: Go has no method overriding across
// embedding, so the base dispatches getContainerExtractor through
// asic.AbstractASiCContentBuilderOverrides (S7_BRIEF.md's virtual-dispatch warning).
type ASiCWithXAdESASiCContentBuilder struct {
	*asic.AbstractASiCContentBuilder
}

var _ asic.AbstractASiCContentBuilderOverrides = (*ASiCWithXAdESASiCContentBuilder)(nil)

// NewASiCWithXAdESASiCContentBuilder is the default constructor. Ports the empty
// ASiCWithXAdESASiCContentBuilder().
func NewASiCWithXAdESASiCContentBuilder() *ASiCWithXAdESASiCContentBuilder {
	builder := &ASiCWithXAdESASiCContentBuilder{
		AbstractASiCContentBuilder: asic.NewAbstractASiCContentBuilderBase(),
	}
	builder.InitAbstractASiCContentBuilder(builder)
	return builder
}

// GetContainerExtractor ports the @Override protected getContainerExtractor(DSSDocument).
//
// Cross-chunk assumption (XADVAL): NewASiCWithXAdESContainerExtractor(model.DSSDocument)
// mirrors `new ASiCWithXAdESContainerExtractor(archiveDocument)`; ASiCWithXAdESContainerExtractor
// is ported alongside ASiCWithXAdESSignatureParameters in the sibling XADVAL manifest
// (S7_BRIEF.md), following the same shape as ASiCWithCAdESContainerExtractor in the CADVAL
// chunk - a type embedding *asic.DefaultASiCContainerExtractor.
func (b *ASiCWithXAdESASiCContentBuilder) GetContainerExtractor(archiveDocument model.DSSDocument) asic.ASiCContainerExtractor {
	return NewASiCWithXAdESContainerExtractor(archiveDocument)
}
