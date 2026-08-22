// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/signature/ASiCSignatureExtensionHelper.java (DSS 6.5.RC1).
package asic

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/exception"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
	"github.com/ryftcore/dss-go/dss/utils"
)

// SignatureExtensionHelperOverrides declares the operations SignatureExtensionHelper
// calls back into virtually. Every concrete helper must call InitASiCSignatureExtensionHelper
// with itself before use - including from within its own constructor, since the base
// constructor logic (extraction from a raw container document)
// itself needs GetASiCContainerExtractor and is split into InitFromDocument/InitFromContent for
// that reason (see those methods' doc comments).
type SignatureExtensionHelperOverrides interface {
	// GetDetachedDocuments returns a list of detached documents for a signature with a given
	// filename. Port of the public abstract getDetachedDocuments(String).
	GetDetachedDocuments(signatureFilename string) []model.DSSDocument

	// GetManifestFile returns a related manifest file for a signature with the given filename.
	// NOTE: used for ASiC with CAdES only. Port of the public getManifestFile(String)
	// (default: not applicable, returns nil - overridable).
	GetManifestFile(signatureFilename string) *model.ManifestFile

	// GetASiCContainerExtractor gets an ASiC container extractor relative to the current
	// implementation. Port of the protected abstract getASiCContainerExtractor(DSSDocument).
	GetASiCContainerExtractor(asicContainer model.DSSDocument) ContainerExtractor

	// GetDocumentAnalyzer gets a Document Analyzer relative to the current implementation.
	// Port of the protected abstract getDocumentAnalyzer(DSSDocument).
	GetDocumentAnalyzer(signatureDocument model.DSSDocument) analyzer.DocumentAnalyzer
}

// SignatureExtensionHelper contains methods for document extraction in order to expand the
// signature with additional elements.
type SignatureExtensionHelper struct {
	// overrides points back at the concrete helper; see InitASiCSignatureExtensionHelper.
	overrides SignatureExtensionHelperOverrides

	// asicContent represents a cached instance of ASiC container extraction result.
	asicContent *Content
}

// NewASiCSignatureExtensionHelperBase builds the empty base state a subclass embeds. Port of
// the shared construction logic; the subclass constructor must follow it with
// InitASiCSignatureExtensionHelper and then InitFromDocument or InitFromContent.
func NewASiCSignatureExtensionHelperBase() *SignatureExtensionHelper {
	return &SignatureExtensionHelper{}
}

// InitASiCSignatureExtensionHelper registers the concrete helper with its base so the base can
// dispatch GetDetachedDocuments/GetManifestFile/GetASiCContainerExtractor/GetDocumentAnalyzer.
// Every concrete helper constructor must call this once, before InitFromDocument (which calls
// back into GetASiCContainerExtractor).
func (h *SignatureExtensionHelper) InitASiCSignatureExtensionHelper(overrides SignatureExtensionHelperOverrides) {
	h.overrides = overrides
}

func (h *SignatureExtensionHelper) requireOverrides() SignatureExtensionHelperOverrides {
	if h.overrides == nil {
		panic("ASiCSignatureExtensionHelper was not initialised: the concrete helper must call InitASiCSignatureExtensionHelper in its constructor")
	}
	return h.overrides
}

// InitFromDocument ports the protected ASiCSignatureExtensionHelper(DSSDocument) constructor,
// split out because it calls back into GetASiCContainerExtractor (an overrides method) - see
// InitASiCSignatureExtensionHelper's doc comment.
//
// Panics with an *exception.IllegalInputException when asicContainer is not an ASiC container.
func (h *SignatureExtensionHelper) InitFromDocument(asicContainer model.DSSDocument) {
	isASiC, err := UtilsIsASiC(asicContainer)
	if err != nil {
		panic(err)
	}
	if !isASiC {
		panic(exception.NewIllegalInputException("The provided file shall be an ASiC container with signatures inside!"))
	}
	extractor := h.requireOverrides().GetASiCContainerExtractor(asicContainer)
	content, err := extractor.Extract()
	if err != nil {
		panic(err)
	}
	h.asicContent = content
}

// InitFromContent ports the protected ASiCSignatureExtensionHelper(ASiCContent) constructor.
func (h *SignatureExtensionHelper) InitFromContent(asicContent *Content) {
	h.asicContent = asicContent
}

// GetAsicContent returns the ASiCContent. Ports getAsicContent().
func (h *SignatureExtensionHelper) GetAsicContent() *Content {
	return h.asicContent
}

// ExtractSignatureDocument returns a file containing a signature with the given id. Ports
// extractSignatureDocument(String).
//
// Panics with an *exception.IllegalInputException / Java's IllegalArgumentException message on
// invalid input, matching Java's thrown exceptions.
func (h *SignatureExtensionHelper) ExtractSignatureDocument(signatureId string) model.DSSDocument {
	signatureDocuments := h.asicContent.SignatureDocuments()
	if utils.IsCollectionEmpty(signatureDocuments) {
		panic(exception.NewIllegalInputException("No signatures found to be extended!"))
	}

	// single signature processing
	if len(signatureDocuments) == 1 {
		signatureDocument := signatureDocuments[0]
		if h.containsSignatureWithId(signatureDocument, signatureId, true) {
			h.CheckSignatureExtensionPossible(signatureDocument)
			return signatureDocument
		}
	}

	// multiple signatures container processing
	if signatureId == "" {
		panic("More than one signature found in a document! " +
			"Please provide a signatureId within the parameters.")
	}
	for _, signatureDocument := range signatureDocuments {
		if h.containsSignatureWithId(signatureDocument, signatureId, false) {
			h.CheckSignatureExtensionPossible(signatureDocument)
			return signatureDocument
		}
	}
	panic(fmt.Sprintf("A signature with id '%s' has not been found!", signatureId))
}

// GetSignatureDocuments gets a list of signature documents. This method performs a validation
// of the signature file, whether its extension is possible. Panics with an
// *exception.IllegalInputException in case of invalid extension configuration. Ports
// getSignatureDocuments().
func (h *SignatureExtensionHelper) GetSignatureDocuments() []model.DSSDocument {
	signatureDocuments := h.asicContent.SignatureDocuments()
	if utils.IsCollectionEmpty(signatureDocuments) {
		panic(exception.NewIllegalInputException("No supported signature documents found! Unable to extend the container."))
	}
	for _, signatureDocument := range signatureDocuments {
		h.CheckSignatureExtensionPossible(signatureDocument)
	}
	return signatureDocuments
}

// containsSignatureWithId ports the private containsSignatureWithId(DSSDocument, String,
// boolean). Logging (LOG.warn) is dropped per PORTING.md; the caught exception is swallowed
// exactly as Java swallows it (returns false), matching the original control flow.
func (h *SignatureExtensionHelper) containsSignatureWithId(signatureDocument model.DSSDocument, signatureId string, acceptSingleSignature bool) (result bool) {
	defer func() {
		if recover() != nil {
			result = false
		}
	}()

	overrides := h.requireOverrides()
	validator := overrides.GetDocumentAnalyzer(signatureDocument)
	validator.SetDetachedContents(overrides.GetDetachedDocuments(signatureDocument.Name()))
	validator.SetManifestFile(overrides.GetManifestFile(signatureDocument.Name()))

	signatures := validator.Signatures()
	if acceptSingleSignature && len(signatures) == 1 {
		return true
	}
	for _, signature := range signatures {
		if h.containsSignatureWithIdRecursive(signature, signatureId) {
			return true
		}
	}
	return false
}

// containsSignatureWithIdRecursive ports the private containsSignatureWithId(AdvancedSignature,
// String).
func (h *SignatureExtensionHelper) containsSignatureWithIdRecursive(signature validation.AdvancedSignature, signatureId string) bool {
	if signatureId == signature.ID() || signatureId == signature.DAIdentifier() {
		return true
	}
	for _, counterSignature := range signature.CounterSignatures() {
		if h.containsSignatureWithIdRecursive(counterSignature, signatureId) {
			return true
		}
	}
	return false
}

// CheckSignatureExtensionPossible verifies if a signatureDocument can be modified. Panics with
// an *exception.IllegalInputException when an extension is not possible. Ports
// checkSignatureExtensionPossible(DSSDocument).
func (h *SignatureExtensionHelper) CheckSignatureExtensionPossible(signatureDocument model.DSSDocument) {
	if UtilsIsCoveredByManifest(h.GetAsicContent().AllManifestDocuments(), signatureDocument.Name()) {
		panic(exception.NewIllegalInputException(fmt.Sprintf("The modification of the signature is not possible! "+
			"Reason : a signature with a filename '%s' is covered by another manifest.", signatureDocument.Name())))
	}
}
