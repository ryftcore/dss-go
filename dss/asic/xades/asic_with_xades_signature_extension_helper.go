// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/signature/ASiCWithXAdESSignatureExtensionHelper.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	dssxades "github.com/ryftcore/dss-go/dss/xades"
)

// ASiCWithXAdESSignatureExtensionHelper contains useful methods for ASiC with XAdES counter
// signature creation.
//
// Java's `extends ASiCSignatureExtensionHelper` becomes embedding plus the
// InitASiCSignatureExtensionHelper(self) registration; the base's two constructors are split
// into InitFromDocument/InitFromContent because the document one calls back into
// getASiCContainerExtractor.
type ASiCWithXAdESSignatureExtensionHelper struct {
	*asic.SignatureExtensionHelper
}

var _ asic.SignatureExtensionHelperOverrides = (*ASiCWithXAdESSignatureExtensionHelper)(nil)

// NewASiCWithXAdESSignatureExtensionHelper is the default constructor. Ports the protected
// ASiCWithXAdESSignatureExtensionHelper(DSSDocument).
func NewASiCWithXAdESSignatureExtensionHelper(asicContainer model.DSSDocument) *ASiCWithXAdESSignatureExtensionHelper {
	helper := &ASiCWithXAdESSignatureExtensionHelper{
		SignatureExtensionHelper: asic.NewSignatureExtensionHelperBase(),
	}
	helper.InitASiCSignatureExtensionHelper(helper)
	helper.InitFromDocument(asicContainer)
	return helper
}

// NewASiCWithXAdESSignatureExtensionHelperFromContent creates a helper from an Content.
// Ports the protected ASiCWithXAdESSignatureExtensionHelper(ASiCContent).
func NewASiCWithXAdESSignatureExtensionHelperFromContent(asicContent *asic.Content) *ASiCWithXAdESSignatureExtensionHelper {
	helper := &ASiCWithXAdESSignatureExtensionHelper{
		SignatureExtensionHelper: asic.NewSignatureExtensionHelperBase(),
	}
	helper.InitASiCSignatureExtensionHelper(helper)
	helper.InitFromContent(asicContent)
	return helper
}

// GetASiCContainerExtractor ports the @Override protected
// getASiCContainerExtractor(DSSDocument).
func (h *ASiCWithXAdESSignatureExtensionHelper) GetASiCContainerExtractor(asicContainer model.DSSDocument) asic.ContainerExtractor {
	return NewASiCWithXAdESContainerExtractor(asicContainer)
}

// GetDocumentAnalyzer ports the @Override protected getDocumentAnalyzer(DSSDocument).
//
// Java's `new XMLDocumentAnalyzer(signatureDocument)` throws on an unparsable XML; the Go
// constructor returns the error instead, which is re-raised as a panic here so the control flow
// stays the one the caller (ASiCSignatureExtensionHelper#containsSignatureWithId, which catches
// and returns false) expects.
func (h *ASiCWithXAdESSignatureExtensionHelper) GetDocumentAnalyzer(signatureDocument model.DSSDocument) analyzer.DocumentAnalyzer {
	documentAnalyzer, err := dssxades.NewXMLDocumentAnalyzer(signatureDocument)
	if err != nil {
		panic(err)
	}
	return documentAnalyzer
}

// GetDetachedDocuments ports the @Override getDetachedDocuments(String).
//
// NOTE: returns all found documents (any document can be signed).
func (h *ASiCWithXAdESSignatureExtensionHelper) GetDetachedDocuments(signatureFilename string) []model.DSSDocument {
	return h.GetAsicContent().AllDocuments()
}

// GetManifestFile ports the base ASiCSignatureExtensionHelper#getManifestFile(String) default
// this type does not override (Java: "NOTE: used for ASiC with CAdES only").
func (h *ASiCWithXAdESSignatureExtensionHelper) GetManifestFile(signatureFilename string) *model.ManifestFile {
	return nil
}
