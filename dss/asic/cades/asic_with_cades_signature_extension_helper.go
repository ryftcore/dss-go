// Ported from dss-asic-cades/src/main/java/eu/europa/esig/dss/asic/cades/signature/ASiCWithCAdESSignatureExtensionHelper.java (DSS 6.5.RC1).
package cades

import (
	"github.com/ryftcore/dss-go/dss/asic"
	dsscades "github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

// ASiCWithCAdESSignatureExtensionHelper contains useful methods for ASiC with CAdES counter
// signature creation.
//
// Java's `extends ASiCSignatureExtensionHelper` becomes embedding plus the
// InitASiCSignatureExtensionHelper(self) registration; the base's two constructors are split
// into InitFromDocument/InitFromContent because the document one calls back into
// getASiCContainerExtractor (S7_BRIEF.md's virtual-dispatch warning).
type ASiCWithCAdESSignatureExtensionHelper struct {
	*asic.ASiCSignatureExtensionHelper
}

var _ asic.ASiCSignatureExtensionHelperOverrides = (*ASiCWithCAdESSignatureExtensionHelper)(nil)

// NewASiCWithCAdESSignatureExtensionHelper is the default constructor. Ports the protected
// ASiCWithCAdESSignatureExtensionHelper(DSSDocument).
func NewASiCWithCAdESSignatureExtensionHelper(asicContainer model.DSSDocument) *ASiCWithCAdESSignatureExtensionHelper {
	helper := &ASiCWithCAdESSignatureExtensionHelper{
		ASiCSignatureExtensionHelper: asic.NewASiCSignatureExtensionHelperBase(),
	}
	helper.InitASiCSignatureExtensionHelper(helper)
	helper.InitFromDocument(asicContainer)
	return helper
}

// NewASiCWithCAdESSignatureExtensionHelperFromContent creates a helper from an ASiCContent.
// Ports the protected ASiCWithCAdESSignatureExtensionHelper(ASiCContent).
func NewASiCWithCAdESSignatureExtensionHelperFromContent(asicContent *asic.ASiCContent) *ASiCWithCAdESSignatureExtensionHelper {
	helper := &ASiCWithCAdESSignatureExtensionHelper{
		ASiCSignatureExtensionHelper: asic.NewASiCSignatureExtensionHelperBase(),
	}
	helper.InitASiCSignatureExtensionHelper(helper)
	helper.InitFromContent(asicContent)
	return helper
}

// GetASiCContainerExtractor ports the @Override protected
// getASiCContainerExtractor(DSSDocument).
func (h *ASiCWithCAdESSignatureExtensionHelper) GetASiCContainerExtractor(asicContainer model.DSSDocument) asic.ASiCContainerExtractor {
	return NewASiCWithCAdESContainerExtractor(asicContainer)
}

// GetDocumentAnalyzer ports the @Override protected getDocumentAnalyzer(DSSDocument).
//
// Java's `new CMSDocumentAnalyzer(signatureDocument)` throws on an unparsable CMS; the Go
// constructor returns the error instead, which is re-raised as a panic here so the control flow
// stays the one the caller (ASiCSignatureExtensionHelper#containsSignatureWithId, which catches
// and returns false) expects.
func (h *ASiCWithCAdESSignatureExtensionHelper) GetDocumentAnalyzer(signatureDocument model.DSSDocument) analyzer.DocumentAnalyzer {
	documentAnalyzer, err := dsscades.NewCMSDocumentAnalyzerFromDocument(signatureDocument)
	if err != nil {
		panic(err)
	}
	return documentAnalyzer
}

// GetDetachedDocuments ports the @Override getDetachedDocuments(String).
func (h *ASiCWithCAdESSignatureExtensionHelper) GetDetachedDocuments(signatureFilename string) []model.DSSDocument {
	signedDocument := ASiCWithCAdESUtilsGetSignedDocument(h.GetAsicContent(), signatureFilename)
	if signedDocument != nil {
		return []model.DSSDocument{signedDocument}
	}
	return []model.DSSDocument{}
}

// GetManifestFile ports the @Override getManifestFile(String).
func (h *ASiCWithCAdESSignatureExtensionHelper) GetManifestFile(signatureFilename string) *model.ManifestFile {
	signatureManifest := asic.ASiCManifestParserGetLinkedManifest(
		h.GetAsicContent().AllManifestDocuments(), signatureFilename)
	if signatureManifest != nil {
		return asic.ASiCManifestParserGetManifestFile(signatureManifest)
	}
	return nil
}
