// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/AbstractGetDataToSignHelper.java (DSS 6.5.RC1).
package asic

// AbstractGetDataToSignHelper is an abstract base helping to extract a data to be signed for
// ASiC containers.
type AbstractGetDataToSignHelper struct {
	// asicContent is the content of the ASiC container. Java declares the field protected
	// final; a Go subclass in another package reaches it through AsicContent().
	asicContent *ASiCContent
}

// NewAbstractGetDataToSignHelper is the default constructor. Ports
// AbstractGetDataToSignHelper(ASiCContent).
func NewAbstractGetDataToSignHelper(asicContent *ASiCContent) AbstractGetDataToSignHelper {
	return AbstractGetDataToSignHelper{asicContent: asicContent}
}

// AsicContent returns the content of the ASiC container. Accessor for the protected
// asicContent field.
func (h *AbstractGetDataToSignHelper) AsicContent() *ASiCContent {
	return h.asicContent
}
