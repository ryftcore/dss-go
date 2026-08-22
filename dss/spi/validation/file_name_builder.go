// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/FileNameBuilder.java (DSS 6.5.RC1).
package validation

import (
	"fmt"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

const (
	// fileNameBuilderContainerPrefix represents a container prefix string.
	fileNameBuilderContainerPrefix = "container"
	// fileNameBuilderEAAPrefix represents an EAA prefix string.
	fileNameBuilderEAAPrefix = "eaa"
	// fileNameBuilderDocumentPrefix represents a document prefix string.
	fileNameBuilderDocumentPrefix = "document"
	// fileNameBuilderSignedSuffix represents a signed document suffix string.
	fileNameBuilderSignedSuffix = "-signed"
	// fileNameBuilderCounterSignedSuffix represents a counter-signed document suffix string.
	fileNameBuilderCounterSignedSuffix = "-counter-signed"
	// fileNameBuilderTimestampedSuffix represents a timestamped document suffix string.
	fileNameBuilderTimestampedSuffix = "-timestamped"
	// fileNameBuilderExtendedSuffix represents an extended document suffix string.
	fileNameBuilderExtendedSuffix = "-extended"
	// fileNameBuilderSignaturePolicyStoreSuffix represents a document with added
	// signature-policy-store suffix string.
	fileNameBuilderSignaturePolicyStoreSuffix = "-sig-policy-store"
	// fileNameBuilderEvidenceRecordSuffix represents a document with added evidence-record
	// suffix string.
	fileNameBuilderEvidenceRecordSuffix = "-preserved"
	// fileNameBuilderEAAPresentationSuffix represents an issued EAA presentation document.
	fileNameBuilderEAAPresentationSuffix = "-presentation"
	// fileNameBuilderP7MExtension is the filename extension for an enveloping CMS signature.
	fileNameBuilderP7MExtension = "p7m"
	// fileNameBuilderP7SExtension is the filename extension for a detached CMS signature.
	fileNameBuilderP7SExtension = "p7s"
)

// FileNameBuilder is used to create a meaningful name for a document depending on its original
// name and the signing operation.
type FileNameBuilder struct {
	// originalFilename is the original document filename.
	originalFilename string
	// signingOperation is the performed signing-operation.
	signingOperation enumerations.SigningOperation
	// signatureLevel is the final signature level.
	signatureLevel enumerations.SignatureLevel
	// signaturePackaging is the signature packaging.
	signaturePackaging enumerations.SignaturePackaging
	// mimeType is the target document MimeType (used for extension definition).
	mimeType enumerations.MimeType
}

// NewFileNameBuilder instantiates the builder. Port of the default constructor.
func NewFileNameBuilder() *FileNameBuilder {
	return &FileNameBuilder{}
}

// SetOriginalFilename sets the original filename of the document. Port of setOriginalFilename(String).
func (b *FileNameBuilder) SetOriginalFilename(originalFilename string) *FileNameBuilder {
	b.originalFilename = originalFilename
	return b
}

// SetSigningOperation sets the performed signing operation type. Port of setSigningOperation(SigningOperation).
func (b *FileNameBuilder) SetSigningOperation(signingOperation enumerations.SigningOperation) *FileNameBuilder {
	b.signingOperation = signingOperation
	return b
}

// SetSignatureLevel sets the final signature level. Port of setSignatureLevel(SignatureLevel).
func (b *FileNameBuilder) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) *FileNameBuilder {
	b.signatureLevel = signatureLevel
	return b
}

// SetSignaturePackaging sets the signature packaging. Port of setSignaturePackaging(SignaturePackaging).
func (b *FileNameBuilder) SetSignaturePackaging(signaturePackaging enumerations.SignaturePackaging) *FileNameBuilder {
	b.signaturePackaging = signaturePackaging
	return b
}

// SetMimeType sets the document mimetype. Port of setMimeType(MimeType).
func (b *FileNameBuilder) SetMimeType(mimeType enumerations.MimeType) *FileNameBuilder {
	b.mimeType = mimeType
	return b
}

// Build generates and returns a final name for the document to create. Port of build().
//
// Java's build() is declared without a checked exception (DSSException is a RuntimeException);
// its two DSSException throw sites (an unsupported SigningOperation, an unsupported
// SignatureForm) become the returned error, per PORTING.md's "data-dependent throw -> (T, error)"
// convention.
func (b *FileNameBuilder) Build() (string, error) {
	var finalName strings.Builder

	var originalName string
	originalExtension := utils.EmptyString
	if b.isContainerMimeType(b.mimeType) {
		originalName = fileNameBuilderContainerPrefix
	} else {
		originalName = b.originalFilename
	}

	if utils.IsStringNotEmpty(originalName) {
		originalExtension = utils.GetFileNameExtension(originalName)
		if utils.IsStringNotEmpty(originalExtension) {
			// remove extension
			originalName = originalName[:len(originalName)-len(originalExtension)-1]
		}
		originalName = spi.DSSUtilsReplaceAllNonAlphanumericCharacters(originalName, "-")

		finalName.WriteString(originalName)

	} else if b.isEAA() {
		finalName.WriteString(fileNameBuilderEAAPrefix)
	} else {
		finalName.WriteString(fileNameBuilderDocumentPrefix)
	}

	if b.signingOperation != "" {
		switch b.signingOperation {
		case enumerations.SigningOperationSign:
			finalName.WriteString(fileNameBuilderSignedSuffix)
		case enumerations.SigningOperationCounterSign:
			finalName.WriteString(fileNameBuilderCounterSignedSuffix)
		case enumerations.SigningOperationTimestamp:
			finalName.WriteString(fileNameBuilderTimestampedSuffix)
		case enumerations.SigningOperationExtend:
			finalName.WriteString(fileNameBuilderExtendedSuffix)
		case enumerations.SigningOperationAddSigPolicyStore:
			finalName.WriteString(fileNameBuilderSignaturePolicyStoreSuffix)
		case enumerations.SigningOperationAddEvidenceRecord:
			finalName.WriteString(fileNameBuilderEvidenceRecordSuffix)
		case enumerations.SigningOperationEAAPresentation:
			finalName.WriteString(fileNameBuilderEAAPresentationSuffix)
		default:
			return "", model.NewDSSError(fmt.Sprintf("The following operation '%s' is not supported!", b.signingOperation))
		}
	}

	if b.signatureLevel != "" {
		finalName.WriteByte('-')
		finalName.WriteString(utils.LowerCase(strings.ReplaceAll(string(b.signatureLevel), "_", "-")))
	}

	extension, err := b.fileExtensionString(b.signatureLevel, b.signaturePackaging, b.mimeType)
	if err != nil {
		return "", err
	}
	if !utils.IsStringNotBlank(extension) {
		extension = originalExtension
	}
	if utils.IsStringNotBlank(extension) {
		finalName.WriteByte('.')
		finalName.WriteString(extension)
	}

	return finalName.String(), nil
}

// isContainerMimeType is the private isContainerMimeType(MimeType).
func (b *FileNameBuilder) isContainerMimeType(mimeType enumerations.MimeType) bool {
	return mimeType == enumerations.MimeTypeEnumASiCS || mimeType == enumerations.MimeTypeEnumASiCE
}

// isEAA is the private isEAA().
func (b *FileNameBuilder) isEAA() bool {
	return enumerations.SigningOperationEAAPresentation == b.signingOperation
}

// fileExtensionString is the private getFileExtensionString(SignatureLevel, SignaturePackaging,
// MimeType).
func (b *FileNameBuilder) fileExtensionString(level enumerations.SignatureLevel, packaging enumerations.SignaturePackaging, mimeType enumerations.MimeType) (string, error) {
	if mimeType != nil {
		return mimeType.Extension(), nil

	} else if level != "" {
		signatureForm, err := level.SignatureForm()
		if err != nil {
			return "", err
		}
		switch signatureForm {
		case enumerations.SignatureFormXAdES:
			return enumerations.MimeTypeEnumXML.Extension(), nil
		case enumerations.SignatureFormCAdES:
			if packaging != "" {
				if enumerations.SignaturePackagingDetached == packaging {
					return fileNameBuilderP7SExtension, nil
				}
				return fileNameBuilderP7MExtension, nil
			}
			// return empty (break)
		case enumerations.SignatureFormPAdES:
			return enumerations.MimeTypeEnumPDF.Extension(), nil
		case enumerations.SignatureFormJAdES:
			return enumerations.MimeTypeEnumJSON.Extension(), nil
		case enumerations.SignatureFormCBAdES:
			return enumerations.MimeTypeEnumCose.Extension(), nil
		default:
			return "", model.NewDSSError(fmt.Sprintf(
				"Unable to generate a full document name! The SignatureForm %s is not supported.", signatureForm))
		}
	}
	return utils.EmptyString, nil
}
