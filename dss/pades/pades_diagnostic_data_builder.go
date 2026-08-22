// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESDiagnosticDataBuilder.java
// (DSS 6.5.RC1).
//
// Virtual dispatch: buildDetachedXmlTimestamp() and buildXmlOrphanTokens() are both called back
// into from the base SignedDocumentDiagnosticDataBuilder/DataBuilder rather than
// invoked directly on a concretely-typed receiver, so both are exported here
// (BuildDetachedXmlTimestamp, BuildXmlOrphanTokens) to satisfy
// dssdiagnostic.SignedDocumentDiagnosticDataBuilderOverrides, which includes a
// BuildXmlOrphanTokens hook (validation/reports/diagnostic/signed_document_diagnostic_data_builder.go)
// so this override can reach it; see that file's doc comment.
package pades

import (
	"time"

	"github.com/ryftcore/dss-go/dss/cades"
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model/x509/revocation"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	dssdiagnostic "github.com/ryftcore/dss-go/dss/validation/reports/diagnostic"
)

// PAdESDiagnosticDataBuilder is the DiagnosticDataBuilder for a PDF signature. Port of the class
// DiagnosticDataBuilder, extending cades.DiagnosticDataBuilder.
type DiagnosticDataBuilder struct {
	cades.DiagnosticDataBuilder
}

// NewPAdESDiagnosticDataBuilder is the port of the default constructor.
func NewDiagnosticDataBuilder() *DiagnosticDataBuilder {
	b := &DiagnosticDataBuilder{
		DiagnosticDataBuilder: cades.DiagnosticDataBuilder{
			SignedDocumentDiagnosticDataBuilder: *dssdiagnostic.NewSignedDocumentDiagnosticDataBuilder(),
		},
	}
	b.InitSignedDocumentDiagnosticDataBuilder(b)
	return b
}

// BuildDetachedXmlSignature builds the XmlSignature, adding PDF-specific PDFRevision and
// VRIDictionaryCreationTime data. Port of the buildDetachedXmlSignature(AdvancedSignature)
// override.
func (b *DiagnosticDataBuilder) BuildDetachedXmlSignature(signature validation.AdvancedSignature) *jaxb.XmlSignature {
	xmlSignature := b.DiagnosticDataBuilder.BuildDetachedXmlSignature(signature)
	padesSignature := signature.(*Signature)
	xmlSignature.PDFRevision = b.xmlPDFRevision(padesSignature.PdfRevision())
	xmlSignature.VRIDictionaryCreationTime = xsDateTimeOrNil(padesSignature.VRICreationTime())
	return xmlSignature
}

// BuildDetachedXmlTimestamp builds the XmlTimestamp, adding the PDFRevision for
// DOCUMENT_TIMESTAMPs. Port of the protected buildDetachedXmlTimestamp(TimestampToken) override.
func (b *DiagnosticDataBuilder) BuildDetachedXmlTimestamp(timestampToken *validation.TimestampToken) *jaxb.XmlTimestamp {
	xmlTimestamp := b.DiagnosticDataBuilder.BuildDetachedXmlTimestamp(timestampToken)
	if pdfTimestampToken, ok := PdfTimestampTokenOf(timestampToken); ok {
		xmlTimestamp.PDFRevision = b.xmlPDFRevision(pdfTimestampToken.PdfRevision())
	}
	return xmlTimestamp
}

// xsDateTimeOrNil wraps a *time.Time for an XSDateTime-typed field, nil-vs-empty preserved:
// Java's setVRIDictionaryCreationTime(Date) accepts the null PAdESSignature.getVRICreationTime()
// may return.
func xsDateTimeOrNil(t *time.Time) *jaxb.XSDateTime {
	if t == nil {
		return nil
	}
	return jaxb.NewXSDateTime(*t)
}

// xmlPDFRevision ports the private getXmlPDFRevision(PdfRevision).
func (b *DiagnosticDataBuilder) xmlPDFRevision(pdfRevision PdfRevision) *jaxb.XmlPDFRevision {
	if pdfRevision == nil {
		return nil
	}
	xmlPDFRevision := &jaxb.XmlPDFRevision{}
	fields := pdfRevision.Fields()
	if utils.IsCollectionNotEmpty(fields) {
		for _, field := range fields {
			xmlPDFRevision.SignatureField = append(xmlPDFRevision.SignatureField, b.xmlPDFSignatureField(field))
		}
	}
	xmlPDFRevision.PDFSignatureDictionary = b.xmlPDFSignatureDictionary(pdfRevision.PdfSigDictInfo())
	xmlPDFRevision.ModificationDetection = b.xmlModificationDetection(pdfRevision.ModificationDetection())
	return xmlPDFRevision
}

// xmlPDFSignatureField ports the private getXmlPDFSignatureField(PdfSignatureField).
func (b *DiagnosticDataBuilder) xmlPDFSignatureField(pdfSignatureField *PdfSignatureField) *jaxb.XmlPDFSignatureField {
	name := pdfSignatureField.FieldName()
	return &jaxb.XmlPDFSignatureField{
		Name:         &name,
		SigFieldLock: b.xmlPDFLockDictionary(pdfSignatureField.LockDictionary()),
	}
}

// xmlPDFLockDictionary ports the private getXmlPDFLockDictionary(SigFieldPermissions).
func (b *DiagnosticDataBuilder) xmlPDFLockDictionary(lockDictionary *SigFieldPermissions) *jaxb.XmlPDFLockDictionary {
	if lockDictionary == nil {
		return nil
	}
	xmlPDFLockDictionary := &jaxb.XmlPDFLockDictionary{}
	if action := lockDictionary.Action(); action != "" {
		v := jaxb.PdfLockActionValue(action)
		xmlPDFLockDictionary.Action = &v
	}
	if utils.IsCollectionNotEmpty(lockDictionary.Fields()) {
		xmlPDFLockDictionary.Field = append(xmlPDFLockDictionary.Field, lockDictionary.Fields()...)
	}
	if permission := lockDictionary.CertificationPermission(); permission != "" {
		v := jaxb.CertificationPermissionValue(permission)
		xmlPDFLockDictionary.Permissions = &v
	}
	return xmlPDFLockDictionary
}

// xmlPDFSignatureDictionary ports the private getXmlPDFSignatureDictionary(PdfSignatureDictionary).
func (b *DiagnosticDataBuilder) xmlPDFSignatureDictionary(pdfSigDict *PdfSignatureDictionary) *jaxb.XmlPDFSignatureDictionary {
	if pdfSigDict == nil {
		return nil
	}
	return &jaxb.XmlPDFSignatureDictionary{
		SignerName:         padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.SignerName()),
		Type:               padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Type()),
		Filter:             padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Filter()),
		SubFilter:          padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.SubFilter()),
		ContactInfo:        padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.ContactInfo()),
		Location:           padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Location()),
		Reason:             padesDiagnosticDataBuilderEmptyToNil(pdfSigDict.Reason()),
		SignatureByteRange: b.xmlByteRange(pdfSigDict.ByteRange()),
		DocMDP:             b.xmlDocMDP(pdfSigDict.DocMDP()),
		FieldMDP:           b.xmlPDFLockDictionary(pdfSigDict.FieldMDP()),
		Consistent:         pdfSigDict.IsConsistent(),
	}
}

// padesDiagnosticDataBuilderEmptyToNil ports the private emptyToNull(String) inherited from the
// diagnostic-data builder base: an empty Go string already stands for Java's null throughout
// this port, so a *string is only introduced here where the Xml field needs to distinguish
// "absent" (nil) from "present but empty" (a pointer to "").
func padesDiagnosticDataBuilderEmptyToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// xmlByteRange ports the private getXmlByteRange(ByteRange).
func (b *DiagnosticDataBuilder) xmlByteRange(byteRange *ByteRange) *jaxb.XmlByteRange {
	return &jaxb.XmlByteRange{
		Value: byteRange.ToBigIntegerList(),
		Valid: byteRange.IsValid(),
	}
}

// xmlDocMDP ports the private getXmlDocMDP(CertificationPermission).
func (b *DiagnosticDataBuilder) xmlDocMDP(certificationPermission enumerations.CertificationPermission) *jaxb.XmlDocMDP {
	if certificationPermission == "" {
		return nil
	}
	v := jaxb.CertificationPermissionValue(certificationPermission)
	return &jaxb.XmlDocMDP{Permissions: &v}
}

// xmlModificationDetection ports the private getXmlModificationDetection(PdfModificationDetection).
func (b *DiagnosticDataBuilder) xmlModificationDetection(modificationDetection *PdfModificationDetection) *jaxb.XmlModificationDetection {
	if modificationDetection == nil || !modificationDetection.AreModificationsDetected() {
		return nil
	}
	xmlModificationDetection := &jaxb.XmlModificationDetection{}

	if annotationOverlaps := modificationDetection.AnnotationOverlaps(); utils.IsCollectionNotEmpty(annotationOverlaps) {
		xmlModificationDetection.AnnotationOverlap = b.xmlModifications(annotationOverlaps)
	}
	if visualDifferences := modificationDetection.VisualDifferences(); utils.IsCollectionNotEmpty(visualDifferences) {
		xmlModificationDetection.VisualDifference = b.xmlModifications(visualDifferences)
	}
	if pageDifferences := modificationDetection.PageDifferences(); utils.IsCollectionNotEmpty(pageDifferences) {
		xmlModificationDetection.PageDifference = b.xmlModifications(pageDifferences)
	}

	objectModifications := modificationDetection.ObjectModifications()
	if !objectModifications.IsEmpty() {
		xmlModificationDetection.ObjectModifications = b.xmlObjectModifications(objectModifications)
	}

	return xmlModificationDetection
}

// xmlModifications ports the private getXmlModifications(List<PdfModification>).
func (b *DiagnosticDataBuilder) xmlModifications(modifications []PdfModification) []*jaxb.XmlModification {
	xmlModifications := make([]*jaxb.XmlModification, 0, len(modifications))
	for _, pdfModification := range modifications {
		xmlModifications = append(xmlModifications, b.xmlModification(pdfModification))
	}
	return xmlModifications
}

// xmlModification ports the private getXmlModification(PdfModification).
func (b *DiagnosticDataBuilder) xmlModification(pdfModification PdfModification) *jaxb.XmlModification {
	return &jaxb.XmlModification{Page: jaxb.NewBigIntegerFromInt64(int64(pdfModification.Page()))}
}

// xmlObjectModifications ports the private getXmlObjectModifications(PdfObjectModifications).
func (b *DiagnosticDataBuilder) xmlObjectModifications(objectModifications PdfObjectModifications) *jaxb.XmlObjectModifications {
	xmlObjectModifications := &jaxb.XmlObjectModifications{}
	for _, modification := range objectModifications.SecureChanges() {
		xmlObjectModifications.ExtensionChange = append(xmlObjectModifications.ExtensionChange, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.FormFillInAndSignatureCreationChanges() {
		xmlObjectModifications.SignatureOrFormFill = append(xmlObjectModifications.SignatureOrFormFill, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.AnnotCreationChanges() {
		xmlObjectModifications.AnnotationChange = append(xmlObjectModifications.AnnotationChange, b.xmlObjectModification(modification))
	}
	for _, modification := range objectModifications.UndefinedChanges() {
		xmlObjectModifications.Undefined = append(xmlObjectModifications.Undefined, b.xmlObjectModification(modification))
	}
	return xmlObjectModifications
}

// xmlObjectModification ports the private getXmlObjectModification(ObjectModification).
func (b *DiagnosticDataBuilder) xmlObjectModification(objectModification ObjectModification) *jaxb.XmlObjectModification {
	action := jaxb.PdfObjectModificationTypeValue(objectModification.ActionType())
	return &jaxb.XmlObjectModification{
		Value:     objectModification.ObjectTree().String(),
		Action:    &action,
		FieldName: padesDiagnosticDataBuilderEmptyToNil(objectModification.FieldName()),
		Type:      padesDiagnosticDataBuilderEmptyToNil(objectModification.Type()),
	}
}

// BuildXmlOrphanTokens ports the protected @Override buildXmlOrphanTokens().
func (b *DiagnosticDataBuilder) BuildXmlOrphanTokens() *jaxb.XmlOrphanTokens {
	b.buildOrphanTokensFromDocumentSources() // necessary to collect all data from DSS PDF revisions
	return b.DataBuilder.BuildXmlOrphanTokens()
}

// buildOrphanTokensFromDocumentSources ports the private buildOrphanTokensFromDocumentSources().
func (b *DiagnosticDataBuilder) buildOrphanTokensFromDocumentSources() {
	for _, certificateToken := range b.GetDocumentCertificateSource().Certificates() {
		id := certificateToken.DSSIDAsString()
		if !b.IsKnownCertificate(id) {
			b.BuildXmlOrphanCertificateToken(certificateToken)
		}
	}
	for _, revocationIdentifier := range b.GetDocumentCRLSource().AllRevocationBinaries() {
		id := revocationIdentifier.AsXmlID()
		if !b.IsKnownRevocation(id) {
			b.CreateOrphanTokenFromRevocationIdentifier[revocation.CRL](revocationIdentifier)
		}
	}
	for _, revocationIdentifier := range b.GetDocumentOCSPSource().AllRevocationBinaries() {
		id := revocationIdentifier.AsXmlID()
		if !b.IsKnownRevocation(id) {
			b.CreateOrphanTokenFromRevocationIdentifier[revocation.OCSP](revocationIdentifier)
		}
	}
}
