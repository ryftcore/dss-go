// Ported from dss-validation/.../validation/process/bbb/fc/AbstractSignatureFormatChecking.java (DSS 6.5.RC1).
package fc

import (
	"strings"

	drjaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	policy "github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractSignatureFormatCheckingOverrides captures the two members Java's
// AbstractSignatureFormatChecking declares abstract (filenameAdherenceCheck(),
// manifestFilenameAdherenceCheck()) and self-calls from getASiCContainerValidationChain().
// SignatureFormatChecking and TimestampFormatChecking implement it and register themselves
// via InitAbstractSignatureFormatChecking.
type AbstractSignatureFormatCheckingOverrides interface {
	// FilenameAdherenceCheck creates a filename adherence check for the token type.
	FilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC]
	// ManifestFilenameAdherenceCheck creates a filename adherence check for the manifest
	// file related to the token.
	ManifestFilenameAdherenceCheck() process.ChainItem[*drjaxb.XmlFC]
}

// AbstractSignatureFormatChecking performs validation for signature-based tokens (signatures
// and timestamps). S is an AbstractSignatureWrapper implementation.
type AbstractSignatureFormatChecking[S diagnostic.AbstractSignatureWrapperOverrides] struct {
	AbstractFormatChecking[S]

	overrides AbstractSignatureFormatCheckingOverrides
}

// InitAbstractSignatureFormatChecking wires the shared state; called by the concrete
// constructor before InitChain.
func (c *AbstractSignatureFormatChecking[S]) InitAbstractSignatureFormatChecking(i18nProvider *i18n.I18nProvider,
	diagnosticData *diagnostic.DiagnosticData, token S, context enumerations.Context, pol policy.ValidationPolicy,
	overrides AbstractSignatureFormatCheckingOverrides) {
	c.InitAbstractFormatChecking(i18nProvider, diagnosticData, token, context, pol)
	c.overrides = overrides
}

// GetPDFRevisionValidationChain chains all PDF revision related checks to the given item
// chain, when applicable.
func (c *AbstractSignatureFormatChecking[S]) GetPDFRevisionValidationChain(item process.ChainItem[*drjaxb.XmlFC]) process.ChainItem[*drjaxb.XmlFC] {
	pdfRevision := c.Token.PDFRevision()
	if pdfRevision == nil {
		return item
	}

	if item == nil {
		item = c.byteRangeCheck()
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.byteRangeCheck())
	}

	item = item.SetNextItem(c.byteRangeCollisionCheck())
	item = item.SetNextItem(c.byteRangeAllDocumentCheck())
	item = item.SetNextItem(c.pdfSignatureDictionaryCheck())
	item = item.SetNextItem(c.pdfPageDifferenceCheck())
	item = item.SetNextItem(c.pdfAnnotationOverlapCheck())
	item = item.SetNextItem(c.pdfVisualDifferenceCheck())

	// /DocMDP check
	if pdfRevision.DocMDPPermissions() != "" {
		item = item.SetNextItem(c.docMDPCheck())
	}
	// /FieldMDP
	if pdfRevision.FieldMDP() != nil {
		item = item.SetNextItem(c.fieldMDPCheck())
	}
	// /SigFieldLock
	if pdfRevision.SigFieldLock() != nil {
		item = item.SetNextItem(c.sigFieldLockCheck())
	}

	item = item.SetNextItem(c.formFillChangesCheck())
	item = item.SetNextItem(c.annotationChangesCheck())
	item = item.SetNextItem(c.undefinedChangesCheck())

	return item
}

// GetPdfaValidationChain chains all PDF/A related checks to the given item chain, when applicable.
func (c *AbstractSignatureFormatChecking[S]) GetPdfaValidationChain(item process.ChainItem[*drjaxb.XmlFC]) process.ChainItem[*drjaxb.XmlFC] {
	// executed only when dss-pdfa module has been loaded
	if !c.DiagnosticData.IsPDFAValidationPerformed() {
		return item
	}

	if item == nil {
		item = c.pdfaProfileCheck()
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.pdfaProfileCheck())
	}

	item = item.SetNextItem(c.pdfaCompliantCheck())

	return item
}

// GetASiCContainerValidationChain chains all ASiC container related checks to the given item
// chain, when applicable.
func (c *AbstractSignatureFormatChecking[S]) GetASiCContainerValidationChain(item process.ChainItem[*drjaxb.XmlFC]) process.ChainItem[*drjaxb.XmlFC] {
	if !c.DiagnosticData.IsContainerInfoPresent() {
		return item
	}

	if item == nil {
		item = c.containerTypeCheck()
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.containerTypeCheck())
	}

	item = item.SetNextItem(c.zipCommentPresentCheck())

	if strings.TrimSpace(c.DiagnosticData.ZipComment()) != "" {
		item = item.SetNextItem(c.acceptableZipCommentCheck())
	}

	item = item.SetNextItem(c.mimetypeFilePresentCheck())

	if c.DiagnosticData.IsMimetypeFilePresent() {
		item = item.SetNextItem(c.mimetypeFileContentCheck())
	}

	item = item.SetNextItem(c.manifestFilePresentCheck())
	item = item.SetNextItem(c.signedFilesPresentCheck())
	item = item.SetNextItem(c.overrides.FilenameAdherenceCheck())

	if c.manifestExistsForToken() {
		item = item.SetNextItem(c.overrides.ManifestFilenameAdherenceCheck())
	}

	return item
}

func (c *AbstractSignatureFormatChecking[S]) byteRangeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.ByteRangeConstraint(c.Context)
	return NewByteRangeCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) byteRangeCollisionCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.ByteRangeCollisionConstraint(c.Context)
	return NewByteRangeCollisionCheck(c.I18nProvider, c.Result, c.Token, c.DiagnosticData, constraint)
}

func (c *AbstractSignatureFormatChecking[S]) byteRangeAllDocumentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.ByteRangeAllDocumentConstraint(c.Context)
	return NewByteRangeAllDocumentCheck(c.I18nProvider, c.Result, c.DiagnosticData, constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfSignatureDictionaryCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.PdfSignatureDictionaryConstraint(c.Context)
	return NewPdfSignatureDictionaryCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfPageDifferenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.PdfPageDifferenceConstraint(c.Context)
	return NewPdfPageDifferenceCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfAnnotationOverlapCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.PdfAnnotationOverlapConstraint(c.Context)
	return NewPdfAnnotationOverlapCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfVisualDifferenceCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.PdfVisualDifferenceConstraint(c.Context)
	return NewPdfVisualDifferenceCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) docMDPCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.DocMDPConstraint(c.Context)
	return NewDocMDPCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) fieldMDPCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FieldMDPConstraint(c.Context)
	return NewFieldMDPCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) sigFieldLockCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.SigFieldLockConstraint(c.Context)
	return NewSigFieldLockCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) formFillChangesCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.FormFillChangesConstraint(c.Context)
	return NewFormFillChangesCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) annotationChangesCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AnnotationChangesConstraint(c.Context)
	return NewAnnotationChangesCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) undefinedChangesCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.UndefinedChangesConstraint(c.Context)
	return NewUndefinedChangesCheck(c.I18nProvider, c.Result, c.Token.PDFRevision(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfaProfileCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AcceptablePDFAProfilesConstraint()
	return NewPDFAProfileCheck(c.I18nProvider, c.Result, c.DiagnosticData.PDFAProfileId(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) pdfaCompliantCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.PDFACompliantConstraint()
	return NewPDFAComplianceCheck(c.I18nProvider, c.Result, c.DiagnosticData.IsPDFACompliant(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) containerTypeCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AcceptedContainerTypesConstraint()
	return NewContainerTypeCheck(c.I18nProvider, c.Result, c.DiagnosticData.ContainerType(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) zipCommentPresentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.ZipCommentPresentConstraint()
	return NewZipCommentPresentCheck(c.I18nProvider, c.Result, c.DiagnosticData.ZipComment(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) acceptableZipCommentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AcceptedZipCommentsConstraint()
	return NewAcceptableZipCommentCheck(c.I18nProvider, c.Result, c.DiagnosticData.ZipComment(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) mimetypeFilePresentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.MimeTypeFilePresentConstraint()
	return NewMimeTypeFilePresentCheck(c.I18nProvider, c.Result, c.DiagnosticData.IsMimetypeFilePresent(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) mimetypeFileContentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.AcceptedMimeTypeContentsConstraint()
	return NewAcceptableMimetypeFileContentCheck(c.I18nProvider, c.Result, c.DiagnosticData.MimetypeFileContent(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) manifestFilePresentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.ManifestFilePresentConstraint()
	return NewManifestFilePresentCheck(c.I18nProvider, c.Result, c.DiagnosticData.ContainerInfo(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) signedFilesPresentCheck() process.ChainItem[*drjaxb.XmlFC] {
	constraint := c.Policy.SignedFilesPresentConstraint()
	return NewSignedFilesPresentCheck(c.I18nProvider, c.Result, c.DiagnosticData.ContainerInfo(), constraint)
}

func (c *AbstractSignatureFormatChecking[S]) manifestExistsForToken() bool {
	return c.DiagnosticData.ManifestFileForFilename(c.Token.Filename()) != nil
}
