// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PdfSignatureDictionary.java
// (DSS 6.5.RC1).
//
// PdfSignatureDictionary is a Java interface with exactly one implementation upstream,
// eu.europa.esig.dss.pdf.PdfSigDictWrapper (dss-pdf root package, flattened into this same
// package per the phase 5b layout, but not itself part of this chunk's manifest). Sibling
// chunks already fixed the collapse of interface+impl into a single concrete Go type -
// pdf_document_reader.go's PdfSignatureDictionaryFields keys on `*PdfSignatureDictionary`, and
// native_pdf_document_reader.go/native_pdf_signature_service.go call
// `NewPdfSigDictWrapperFactory(dictionary).Create()` and `.CheckConsistency(...)` against it -
// so this file plays both roles: the getters PdfSignatureDictionary.java declares, and the
// state + checkConsistency/isConsistent bodies PdfSigDictWrapper.java supplies for them.
//
// FORWARD DEPENDENCIES, closed during integration by the eu.europa.esig.dss.pdf package's files
// (see pdf_object.go's header for why that whole Java package landed in no s5b manifest):
// PdfDict (pdf_object.go), SigFieldPermissions (sig_field_permissions.go),
// PAdESConstantsReferenceName / PAdESConstantsDataName (pades_constants.go).
//
// INTEGRATION CORRECTION: this file originally assumed a PdfObjectModificationsFinder with a
// PdfDict-based Find(PdfDict, PdfDict); the actual Java PdfObjectModificationsFinder interface
// method is find(PdfDocumentReader, PdfDocumentReader) (see pdf_object_modifications_finder.go),
// which native_pdf_signature_service.go was already calling under the name Find - so CheckConsistency
// below now calls the concrete *DefaultPdfObjectModificationsFinder's additional, differently-named
// FindInDicts(PdfDict, PdfDict) overload instead (default_pdf_object_modifications_finder.go's
// header explains the Java find(...) overload pair this de-overloads). ObjectModification.ObjectTree()
// returns *PdfObjectTree (Java's own class name - see pdf_object_tree.go), not the placeholder
// "ObjectTree" this header originally speculated.
//
// cms.CMS stands in for org.bouncycastle CMSSignedData wrapped by eu.europa.esig.dss.cms.CMS.
package pades

import (
	"time"

	"github.com/utain/esig/dss/cms"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
)

// PdfSignatureDictionary contains PDF signature dictionary information. It merges the Java
// interface PdfSignatureDictionary with the state and checkConsistency/isConsistent behaviour
// of its sole implementation, PdfSigDictWrapper (see header).
type PdfSignatureDictionary struct {
	// dictionary is the original PDF dictionary.
	dictionary PdfDict

	// signerName is the name of the signer.
	signerName string

	// contactInfo is the contact info of the signer.
	contactInfo string

	// reason is the reason of signing.
	reason string

	// location is the location of signing.
	location string

	// signingDate is the datetime of signing; the zero value stands for Java's null, matching
	// the shape pades_baseline_requirements_checker.go's header already documents for this
	// getter (SigningDate() time.Time, checked with IsZero()).
	signingDate time.Time

	// dictType is the type of the dictionary.
	dictType string

	// filter is the value of the /Filter parameter.
	filter string

	// subFilter is the value of the /SubFilter parameter.
	subFilter string

	// contents is the value of the /Contents signature parameter.
	contents []byte

	// byteRange is the value of the /ByteRange parameter.
	byteRange *ByteRange

	// docMDP is the value of the /DocMDP parameter.
	docMDP enumerations.CertificationPermission

	// fieldMDP is the value of the /FieldMDP parameter.
	fieldMDP *SigFieldPermissions

	// cmsValue is the CMS signature value.
	cmsValue *cms.CMS

	// consistent identifies whether the signature dictionary is consistent between revisions.
	consistent bool
}

// NewPdfSignatureDictionary is the default constructor. Port of the protected
// PdfSigDictWrapper() no-arg constructor.
func NewPdfSignatureDictionary() *PdfSignatureDictionary {
	return &PdfSignatureDictionary{}
}

// SetDictionary sets the signature field dictionary. Port of the protected
// setDictionary(PdfDict).
func (d *PdfSignatureDictionary) SetDictionary(dictionary PdfDict) {
	d.dictionary = dictionary
}

// SignerName gets the name of the signer. Port of getSignerName().
func (d *PdfSignatureDictionary) SignerName() string {
	return d.signerName
}

// SetSignerName sets the name of the signer. Port of the protected setSignerName(String).
func (d *PdfSignatureDictionary) SetSignerName(signerName string) {
	d.signerName = signerName
}

// ContactInfo gets the signer's contact info. Port of getContactInfo().
func (d *PdfSignatureDictionary) ContactInfo() string {
	return d.contactInfo
}

// SetContactInfo sets the contact info. Port of the protected setContactInfo(String).
func (d *PdfSignatureDictionary) SetContactInfo(contactInfo string) {
	d.contactInfo = contactInfo
}

// Reason gets the signing reason. Port of getReason().
func (d *PdfSignatureDictionary) Reason() string {
	return d.reason
}

// SetReason sets the signing reason. Port of the protected setReason(String).
func (d *PdfSignatureDictionary) SetReason(reason string) {
	d.reason = reason
}

// Location gets the signer's location. Port of getLocation().
func (d *PdfSignatureDictionary) Location() string {
	return d.location
}

// SetLocation sets the signer location. Port of the protected setLocation(String).
func (d *PdfSignatureDictionary) SetLocation(location string) {
	d.location = location
}

// SigningDate gets the claimed signing time. Port of getSigningDate(); the zero time.Time
// stands in for Java's null.
func (d *PdfSignatureDictionary) SigningDate() time.Time {
	return d.signingDate
}

// SetSigningDate sets the date of signing. Port of the protected setSigningDate(Date).
func (d *PdfSignatureDictionary) SetSigningDate(signingDate time.Time) {
	d.signingDate = signingDate
}

// Type gets the type of the dictionary. Port of getType().
func (d *PdfSignatureDictionary) Type() string {
	return d.dictType
}

// SetType sets the type of the dictionary. Port of the protected setType(String).
func (d *PdfSignatureDictionary) SetType(dictType string) {
	d.dictType = dictType
}

// Filter gets the Filter value. Port of getFilter().
func (d *PdfSignatureDictionary) Filter() string {
	return d.filter
}

// SetFilter sets the /Filter value. Port of the protected setFilter(String).
func (d *PdfSignatureDictionary) SetFilter(filter string) {
	d.filter = filter
}

// SubFilter gets the SubFilter value. Port of getSubFilter().
func (d *PdfSignatureDictionary) SubFilter() string {
	return d.subFilter
}

// SetSubFilter sets the /SubFilter value. Port of the protected setSubFilter(String).
func (d *PdfSignatureDictionary) SetSubFilter(subFilter string) {
	d.subFilter = subFilter
}

// Contents gets /Contents binaries (CMSSignedData). Port of getContents(): byte[].
func (d *PdfSignatureDictionary) Contents() []byte {
	return d.contents
}

// SetContents sets the /Contents signature value. Port of the protected setContents(byte[]).
func (d *PdfSignatureDictionary) SetContents(contents []byte) {
	d.contents = contents
}

// ByteRange gets the signed/timestamped ByteRange. Port of getByteRange().
func (d *PdfSignatureDictionary) ByteRange() *ByteRange {
	return d.byteRange
}

// SetByteRange sets the /ByteRange value. Port of the protected setByteRange(ByteRange).
func (d *PdfSignatureDictionary) SetByteRange(byteRange *ByteRange) {
	d.byteRange = byteRange
}

// DocMDP returns a /DocMDP dictionary, when present. Port of getDocMDP().
func (d *PdfSignatureDictionary) DocMDP() enumerations.CertificationPermission {
	return d.docMDP
}

// SetDocMDP sets the /DocMPD dictionary value. Port of the protected
// setDocMDP(CertificationPermission).
func (d *PdfSignatureDictionary) SetDocMDP(docMDP enumerations.CertificationPermission) {
	d.docMDP = docMDP
}

// FieldMDP returns a /FieldMDP dictionary, when present. Port of getFieldMDP().
func (d *PdfSignatureDictionary) FieldMDP() *SigFieldPermissions {
	return d.fieldMDP
}

// SetFieldMDP sets the /FieldMDP dictionary value. Port of the protected
// setFieldMDP(SigFieldPermissions).
func (d *PdfSignatureDictionary) SetFieldMDP(fieldMDP *SigFieldPermissions) {
	d.fieldMDP = fieldMDP
}

// CMS gets the CMS from /Contents. Port of getCMS().
func (d *PdfSignatureDictionary) CMS() *cms.CMS {
	return d.cmsValue
}

// SetCMS sets the CMS value read from /Contents. Port of the protected setCMS(CMS).
func (d *PdfSignatureDictionary) SetCMS(cmsValue *cms.CMS) {
	d.cmsValue = cmsValue
}

// CheckConsistency verifies the equality of the current PdfSignatureDictionary with the
// provided signatureDictionary. NOTE: this method is similar to an equals(PdfSignatureDictionary)
// method, but also modifies the state of the object accessible from IsConsistent(). If no
// signature dictionary was found in the signed revision, nil may be provided.
// Port of checkConsistency(PdfSignatureDictionary).
func (d *PdfSignatureDictionary) CheckConsistency(signatureDictionary *PdfSignatureDictionary) bool {
	if signatureDictionary == nil {
		// Upstream logs "PdfSignatureDictionary from signed revision is null!".
		d.consistent = false

	} else {
		modificationsFinder := NewDefaultPdfObjectModificationsFinder()
		pdfObjectModifications := modificationsFinder.FindInDicts(signatureDictionary.dictionary, d.dictionary)
		undefinedChanges := pdfObjectModifications.UndefinedChanges()
		undefinedChanges = removeReferenceData(undefinedChanges)
		d.consistent = utils.IsCollectionEmpty(undefinedChanges)
		// Upstream logs "The signature dictionary from final PDF revision is not equal to the
		// signed revision version!" and, at debug level, the undefined modifications' object
		// trees, when consistent is false.
	}

	return d.consistent
}

// removeReferenceData ports the private removeReferenceData(List<ObjectModification>): the
// /Reference /Data dictionary contains references to PDF objects covered by the signature. The
// changes inside do not impact signature validity directly.
func removeReferenceData(modifications []ObjectModification) []ObjectModification {
	if utils.IsCollectionEmpty(modifications) {
		return modifications
	}
	filtered := modifications[:0]
	for _, objectModification := range modifications {
		keyChain := objectModification.ObjectTree().KeyChain()
		if containsString(keyChain, PAdESConstantsReferenceName) && containsString(keyChain, PAdESConstantsDataName) {
			continue
		}
		filtered = append(filtered, objectModification)
	}
	return filtered
}

// containsString reports whether needle is present in haystack, standing in for
// List#contains(Object) over the object tree's key chain.
func containsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// IsConsistent checks if the signature dictionary is consistent.
// NOTE: CheckConsistency(*PdfSignatureDictionary) shall be executed before!
// Port of isConsistent().
func (d *PdfSignatureDictionary) IsConsistent() bool {
	return d.consistent
}
