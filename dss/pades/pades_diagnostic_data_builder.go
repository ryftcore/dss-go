//go:build phase8

// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESDiagnosticDataBuilder.java
// (DSS 6.5.RC1).
//
// INTEGRATOR NOTE: gated behind the `phase8` build tag so that `go build ./...` / `go vet ./...`
// / `go test ./...` are green for the rest of the module while dss/validation and
// dss/validation/diagnostic do not exist yet, following the cades_diagnostic_data_builder.go
// precedent (see that file's header). Drop the tag once Phase 8 lands those packages.
//
// BLOCKED FORWARD DEPENDENCIES: this class's Java base,
// eu.europa.esig.dss.cades.validation.CAdESDiagnosticDataBuilder, and the JAXB diagnostic-data
// model it builds (eu.europa.esig.dss.diagnostic.jaxb.{XmlSignature,XmlTimestamp,XmlByteRange,
// XmlDocMDP,XmlModification,XmlModificationDetection,XmlObjectModification,
// XmlObjectModifications,XmlOrphanTokens,XmlPDFLockDictionary,XmlPDFRevision,
// XmlPDFSignatureDictionary,XmlPDFSignatureField}) belong to dss-validation (+*-report-jaxb),
// which PORTING_PLAN.md assigns to the not-yet-ported `validation`/`validation/diagnostic`
// packages (Phase 8). The method bodies below are ported 1:1 against the package path and shape
// PORTING_PLAN.md's table implies (github.com/utain/esig/dss/validation/diagnostic), so that
// this file needs no further changes once Phase 8 lands the package - only its imports need to
// resolve, and the assumed Xml* setter/constructor names below (New<Type>() / Set<Field>(value)
// / a bare exported slice field for a repeated element, mirroring
// cades_diagnostic_data_builder.go's NewXmlArchiveTimestampHashIndex()/SetVersion/SetValid/
// Messages precedent) need to match what Phase 8 actually lands.
//
// FORWARD DEPENDENCIES (not in this chunk's manifest, owned by sibling chunks):
//
//   - PdfRevision (eu.europa.esig.dss.pdf.PdfRevision) - an interface with
//     Fields() []*PdfSignatureField, PdfSigDictInfo() *PdfSignatureDictionary,
//     ModificationDetection() PdfModificationDetection.
//
//   - PAdESSignature.PdfRevision() PdfRevision, PAdESSignature.VRICreationTime() *time.Time.
//
//   - PdfTimestampToken.PdfRevision() *PdfDocTimestampRevision (PdfDocTimestampRevision
//     satisfies PdfRevision too); PdfTimestampTokenOf(*validation.TimestampToken)
//     (*PdfTimestampToken, bool) - same lookup assumed by pades_baseline_requirements_checker.go.
//
//   - SigFieldPermissions (eu.europa.esig.dss.pdf.SigFieldPermissions), already used as a
//     forward dependency by pdf_signature_field.go: Action() enumerations.PdfLockAction,
//     Fields() []string, CertificationPermission() enumerations.CertificationPermission.
//
//   - eu.europa.esig.dss.pdf.modifications, flattened into this same package per the phase 5b
//     layout:
//
//     type PdfModificationDetection interface {
//     AreModificationsDetected() bool
//     AnnotationOverlaps() []PdfModification
//     VisualDifferences() []PdfModification
//     PageDifferences() []PdfModification
//     ObjectModifications() *PdfObjectModifications
//     }
//     type PdfModification interface{ Page() int }
//     type PdfObjectModifications struct{ /* ... */ }
//     func (m *PdfObjectModifications) IsEmpty() bool
//     func (m *PdfObjectModifications) SecureChanges() []*ObjectModification
//     func (m *PdfObjectModifications) FormFillInAndSignatureCreationChanges() []*ObjectModification
//     func (m *PdfObjectModifications) AnnotCreationChanges() []*ObjectModification
//     func (m *PdfObjectModifications) UndefinedChanges() []*ObjectModification
//     type ObjectModification struct{ /* ... */ }
//     func (m *ObjectModification) ObjectTree() fmt.Stringer
//     func (m *ObjectModification) ActionType() enumerations.ObjectModificationType // field name TBD by Phase 8's XmlObjectModification shape
//     func (m *ObjectModification) FieldName() string
//     func (m *ObjectModification) Type() enumerations.ObjectModificationKind       // field name TBD
//
//   - PdfSignatureDictionary additionally exposes, beyond the SubFilter/Type/ByteRange/Contents
//     already confirmed by the landed SIGN chunk (see pades_baseline_requirements_checker.go's
//     header): SignerName() string, Filter() string, ContactInfo() string, Location() string,
//     Reason() string, DocMDP() enumerations.CertificationPermission,
//     FieldMDP() *SigFieldPermissions, IsConsistent() bool.
//
// NAMING RISK: dssdiagnostic.SignedDocumentDiagnosticDataBuilder's fields backing
// documentCertificateSource/documentCRLSource/documentOCSPSource/xmlCertsMap/
// xmlOrphanCertificateTokensMap/xmlRevocationsMap/xmlOrphanRevocationTokensMap are Java
// `protected` (direct field access from CAdESDiagnosticDataBuilder in Java, itself in a
// different Java package from SignedDocumentDiagnosticDataBuilder only by module, not by
// visibility). Go visibility has no "protected" tier, so reaching this state from package pades
// (two packages removed from validation/diagnostic) requires Phase 8 to have exported accessors;
// the names below (DocumentCertificateSource(), XmlCertsMapContains(id string) bool, ...) are
// this chunk's best-effort guess, not confirmed against a landed Phase 8.
package pades

import (
	"math/big"

	"github.com/utain/esig/dss/cades"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
	dssdiagnostic "github.com/utain/esig/dss/validation/diagnostic"
)

// PAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a PDF signature.
type PAdESDiagnosticDataBuilder struct {
	cades.CAdESDiagnosticDataBuilder
}

// NewPAdESDiagnosticDataBuilder is the port of the default constructor.
func NewPAdESDiagnosticDataBuilder() *PAdESDiagnosticDataBuilder {
	return &PAdESDiagnosticDataBuilder{}
}

// BuildDetachedXmlSignature builds the XmlSignature, adding PDF-specific PDFRevision and
// VRIDictionaryCreationTime data. Port of the buildDetachedXmlSignature(AdvancedSignature)
// override.
//
// Shadows the embedded base's method of the same name; see the file header's forward-dependency
// note and PORTING.md's "Virtual dispatch" precedent (analyzer/default_document_analyzer.go) on
// why the override has to be reproduced this way rather than relying on embedding alone.
func (b *PAdESDiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *dssdiagnostic.XmlSignature {
	xmlSignature := b.CAdESDiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	padesSignature := signature.(*PAdESSignature)
	xmlSignature.SetPDFRevision(b.xmlPDFRevision(padesSignature.PdfRevision()))
	xmlSignature.SetVRIDictionaryCreationTime(padesSignature.VRICreationTime())
	return xmlSignature
}

// buildDetachedXmlTimestamp builds the XmlTimestamp, adding the PDFRevision for
// DOCUMENT_TIMESTAMPs. Port of the protected buildDetachedXmlTimestamp(TimestampToken) override.
func (b *PAdESDiagnosticDataBuilder) buildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *dssdiagnostic.XmlTimestamp {
	xmlTimestamp := b.CAdESDiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
	if pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken); ok {
		xmlTimestamp.SetPDFRevision(b.xmlPDFRevision(pdfTimestampToken.PdfRevision()))
	}
	return xmlTimestamp
}

// xmlPDFRevision ports the private getXmlPDFRevision(PdfRevision).
func (b *PAdESDiagnosticDataBuilder) xmlPDFRevision(pdfRevision PdfRevision) *dssdiagnostic.XmlPDFRevision {
	if pdfRevision != nil {
		xmlPDFRevision := dssdiagnostic.NewXmlPDFRevision()
		fields := pdfRevision.Fields()
		if utils.IsCollectionNotEmpty(fields) {
			for _, field := range fields {
				xmlPDFRevision.Fields = append(xmlPDFRevision.Fields, b.xmlPDFSignatureField(field))
			}
		}
		xmlPDFRevision.SetPDFSignatureDictionary(b.xmlPDFSignatureDictionary(pdfRevision.PdfSigDictInfo()))
		xmlPDFRevision.SetModificationDetection(b.xmlModificationDetection(pdfRevision.ModificationDetection()))
		return xmlPDFRevision
	}
	return nil
}

// xmlPDFSignatureField ports the private getXmlPDFSignatureField(PdfSignatureField).
func (b *PAdESDiagnosticDataBuilder) xmlPDFSignatureField(pdfSignatureField *PdfSignatureField) *dssdiagnostic.XmlPDFSignatureField {
	xmlPdfSignatureField := dssdiagnostic.NewXmlPDFSignatureField()
	xmlPdfSignatureField.SetName(pdfSignatureField.FieldName())
	xmlPdfSignatureField.SetSigFieldLock(b.xmlPDFLockDictionary(pdfSignatureField.LockDictionary()))
	return xmlPdfSignatureField
}

// xmlPDFLockDictionary ports the private getXmlPDFLockDictionary(SigFieldPermissions).
func (b *PAdESDiagnosticDataBuilder) xmlPDFLockDictionary(lockDictionary *SigFieldPermissions) *dssdiagnostic.XmlPDFLockDictionary {
	if lockDictionary != nil {
		xmlPDFLockDictionary := dssdiagnostic.NewXmlPDFLockDictionary()
		xmlPDFLockDictionary.SetAction(lockDictionary.Action())
		if utils.IsCollectionNotEmpty(lockDictionary.Fields()) {
			xmlPDFLockDictionary.Fields = append(xmlPDFLockDictionary.Fields, lockDictionary.Fields()...)
		}
		if lockDictionary.CertificationPermission() != "" {
			xmlPDFLockDictionary.SetPermissions(lockDictionary.CertificationPermission())
		}
		return xmlPDFLockDictionary
	}
	return nil
}

// xmlPDFSignatureDictionary ports the private getXmlPDFSignatureDictionary(PdfSignatureDictionary).
func (b *PAdESDiagnosticDataBuilder) xmlPDFSignatureDictionary(pdfSigDict *PdfSignatureDictionary) *dssdiagnostic.XmlPDFSignatureDictionary {
	if pdfSigDict != nil {
		pdfSignatureDictionary := dssdiagnostic.NewXmlPDFSignatureDictionary()
		pdfSignatureDictionary.SetSignerName(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.SignerName()))
		pdfSignatureDictionary.SetType(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Type()))
		pdfSignatureDictionary.SetFilter(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Filter()))
		pdfSignatureDictionary.SetSubFilter(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.SubFilter()))
		pdfSignatureDictionary.SetContactInfo(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.ContactInfo()))
		pdfSignatureDictionary.SetLocation(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Location()))
		pdfSignatureDictionary.SetReason(padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Reason()))
		pdfSignatureDictionary.SetSignatureByteRange(b.xmlByteRange(pdfSigDict.ByteRange()))
		pdfSignatureDictionary.SetDocMDP(b.xmlDocMDP(pdfSigDict.DocMDP()))
		pdfSignatureDictionary.SetFieldMDP(b.xmlPDFLockDictionary(pdfSigDict.FieldMDP()))
		pdfSignatureDictionary.SetConsistent(pdfSigDict.IsConsistent())
		return pdfSignatureDictionary
	}
	return nil
}

// padesDiagnosticDataBuilderEmptyToNil ports the private emptyToNull(String) inherited from the
// diagnostic-data builder base: an empty Go string already stands for Java's null throughout
// this port, so a *string is only introduced here where the Xml setter needs to distinguish
// "absent" (nil) from "present but empty" (a pointer to "").
func padesDiagnosticDataBuilderEmptyToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// xmlByteRange ports the private getXmlByteRange(ByteRange).
func (b *PAdESDiagnosticDataBuilder) xmlByteRange(byteRange *ByteRange) *dssdiagnostic.XmlByteRange {
	xmlByteRange := dssdiagnostic.NewXmlByteRange()
	xmlByteRange.Value = append(xmlByteRange.Value, byteRange.ToBigIntegerList()...)
	xmlByteRange.SetValid(byteRange.IsValid())
	return xmlByteRange
}

// xmlDocMDP ports the private getXmlDocMDP(CertificationPermission).
func (b *PAdESDiagnosticDataBuilder) xmlDocMDP(certificationPermission enumerations.CertificationPermission) *dssdiagnostic.XmlDocMDP {
	if certificationPermission != "" {
		xmlDocMDP := dssdiagnostic.NewXmlDocMDP()
		xmlDocMDP.SetPermissions(certificationPermission)
		return xmlDocMDP
	}
	return nil
}

// xmlModificationDetection ports the private getXmlModificationDetection(PdfModificationDetection).
func (b *PAdESDiagnosticDataBuilder) xmlModificationDetection(modificationDetection PdfModificationDetection) *dssdiagnostic.XmlModificationDetection {
	if modificationDetection != nil && modificationDetection.AreModificationsDetected() {
		xmlModificationDetection := dssdiagnostic.NewXmlModificationDetection()

		annotationOverlaps := modificationDetection.AnnotationOverlaps()
		if utils.IsCollectionNotEmpty(annotationOverlaps) {
			xmlModificationDetection.AnnotationOverlap = append(xmlModificationDetection.AnnotationOverlap,
				b.xmlModifications(annotationOverlaps)...)
		}

		visualDifferences := modificationDetection.VisualDifferences()
		if utils.IsCollectionNotEmpty(visualDifferences) {
			xmlModificationDetection.VisualDifference = append(xmlModificationDetection.VisualDifference,
				b.xmlModifications(visualDifferences)...)
		}

		pageDifferences := modificationDetection.PageDifferences()
		if utils.IsCollectionNotEmpty(pageDifferences) {
			xmlModificationDetection.PageDifference = append(xmlModificationDetection.PageDifference,
				b.xmlModifications(pageDifferences)...)
		}

		objectModifications := modificationDetection.ObjectModifications()
		if !objectModifications.IsEmpty() {
			xmlModificationDetection.SetObjectModifications(b.xmlObjectModifications(objectModifications))
		}

		return xmlModificationDetection
	}
	return nil
}

// xmlModifications ports the private getXmlModifications(List<PdfModification>).
func (b *PAdESDiagnosticDataBuilder) xmlModifications(modifications []PdfModification) []*dssdiagnostic.XmlModification {
	xmlModifications := make([]*dssdiagnostic.XmlModification, 0)
	if utils.IsCollectionNotEmpty(modifications) {
		for _, pdfModification := range modifications {
			xmlModifications = append(xmlModifications, b.xmlModification(pdfModification))
		}
	}
	return xmlModifications
}

// xmlModification ports the private getXmlModification(PdfModification).
func (b *PAdESDiagnosticDataBuilder) xmlModification(pdfModification PdfModification) *dssdiagnostic.XmlModification {
	xmlModification := dssdiagnostic.NewXmlModification()
	xmlModification.SetPage(big.NewInt(int64(pdfModification.Page())))
	return xmlModification
}

// xmlObjectModifications ports the private getXmlObjectModifications(PdfObjectModifications).
func (b *PAdESDiagnosticDataBuilder) xmlObjectModifications(objectModifications *PdfObjectModifications) *dssdiagnostic.XmlObjectModifications {
	xmlObjectModifications := dssdiagnostic.NewXmlObjectModifications()
	for _, modification := range objectModifications.SecureChanges() {
		xmlObjectModifications.ExtensionChanges = append(xmlObjectModifications.ExtensionChanges, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.FormFillInAndSignatureCreationChanges() {
		xmlObjectModifications.SignatureOrFormFill = append(xmlObjectModifications.SignatureOrFormFill, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.AnnotCreationChanges() {
		xmlObjectModifications.AnnotationChanges = append(xmlObjectModifications.AnnotationChanges, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.UndefinedChanges() {
		xmlObjectModifications.Undefined = append(xmlObjectModifications.Undefined, b.xmlObjectModification(modification))
	}
	return xmlObjectModifications
}

// xmlObjectModification ports the private getXmlObjectModification(ObjectModification).
func (b *PAdESDiagnosticDataBuilder) xmlObjectModification(objectModification *ObjectModification) *dssdiagnostic.XmlObjectModification {
	xmlObjectModification := dssdiagnostic.NewXmlObjectModification()
	xmlObjectModification.SetValue(objectModification.ObjectTree().String())
	xmlObjectModification.SetAction(objectModification.ActionType())
	xmlObjectModification.SetFieldName(objectModification.FieldName())
	xmlObjectModification.SetType(objectModification.Type())
	return xmlObjectModification
}

// buildXmlOrphanTokens ports the protected buildXmlOrphanTokens() override.
func (b *PAdESDiagnosticDataBuilder) buildXmlOrphanTokens() *dssdiagnostic.XmlOrphanTokens {
	b.buildOrphanTokensFromDocumentSources() // necessary to collect all data from DSS PDF revisions
	return b.CAdESDiagnosticDataBuilder.BuildXmlOrphanTokens()
}

// buildOrphanTokensFromDocumentSources ports the private buildOrphanTokensFromDocumentSources().
func (b *PAdESDiagnosticDataBuilder) buildOrphanTokensFromDocumentSources() {
	for _, certificateToken := range b.DocumentCertificateSource().Certificates() {
		id := certificateToken.DSSIDAsString()
		if !b.XmlCertsMapContains(id) && !b.XmlOrphanCertificateTokensMapContains(id) {
			b.BuildXmlOrphanCertificateToken(certificateToken)
		}
	}
	for _, revocationIdentifier := range b.DocumentCRLSource().AllRevocationBinaries() {
		id := revocationIdentifier.AsXmlID()
		if !b.XmlRevocationsMapContains(id) && !b.XmlOrphanRevocationTokensMapContains(id) {
			b.CreateOrphanTokenFromRevocationIdentifier(revocationIdentifier)
		}
	}
	for _, revocationIdentifier := range b.DocumentOCSPSource().AllRevocationBinaries() {
		id := revocationIdentifier.AsXmlID()
		if !b.XmlRevocationsMapContains(id) && !b.XmlOrphanRevocationTokensMapContains(id) {
			b.CreateOrphanTokenFromRevocationIdentifier(revocationIdentifier)
		}
	}
}
