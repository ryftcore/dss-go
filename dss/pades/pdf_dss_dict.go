// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfDssDict.java,
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/AbstractPdfDssDict.java,
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/SingleDssDict.java,
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/PdfVriDict.java and
// dss-pades/src/main/java/eu/europa/esig/dss/pdf/DSSDictionaryExtractionUtils.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf flattens across several files in this package (see pdf_object.go's
// header). Shapes match those documented by pdf_composite_dss_dict_certificate_source.go /
// pdf_composite_dss_dict_crl_source.go / pdf_composite_dss_dict_ocsp_source.go
// (PdfDssDict.{CRLs,OCSPs,CERTs,VRIs}, PdfVriDict.{Name, TUTime,TSStream}) and
// native_pdf_document_reader.go (SingleDssDictExtract(PdfDict) PdfDssDict).
//
// STRUCTURE DEVIATION: Java expresses SingleDssDict/PdfVriDict as two subclasses of the abstract
// AbstractPdfDssDict, factoring the three token maps and their extraction into the base class.
// Go has no implementation inheritance, so AbstractPdfDssDict becomes dssDictExtraction, a value
// type embedded by both SingleDssDict and PdfVriDict (constructed once, by
// newDssDictExtraction, from the four "which sub-dictionary/array names" abstract-method
// results each subclass supplied in Java - now four plain constructor arguments), as elsewhere
// in this port (e.g. cms_for_pades_builder_helper.go) rather than trying to fake virtual
// dispatch.
//
// DSSDictionaryExtractionUtils.java's four static methods have exactly one call site each
// (AbstractPdfDssDict's constructor, PdfVriDict's constructor) in the whole Java codebase, so
// they are folded into unexported helpers here rather than kept as a separate exported utility
// type; nothing in any landed s5b chunk references DSSDictionaryExtractionUtils by name.
package pades

import (
	"reflect"
	"time"

	"github.com/ryftcore/dss-go/dss/crlparser"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
)

// PdfDssDict represents the DSS dictionary. Port of the PdfDssDict interface.
type PdfDssDict interface {
	// CRLs returns a map of unique identifiers and CRL binaries. Port of getCRLs().
	CRLs() map[PdfObjectKey]*crlparser.CRLBinary

	// OCSPs returns a map of unique identifiers and OCSPResponseBinarys. Port of getOCSPs().
	OCSPs() map[PdfObjectKey]*spi.OCSPResponseBinary

	// CERTs returns a map of unique identifiers and Certificate Tokens. Port of getCERTs().
	CERTs() map[PdfObjectKey]*model.CertificateToken

	// VRIs returns a list of VRI dictionaries. Port of getVRIs().
	VRIs() []*PdfVriDict
}

// dssDictExtraction is the Go stand-in for AbstractPdfDssDict: the three token maps common to
// the DSS dictionary and every VRI dictionary, plus their shared extraction logic.
type dssDictExtraction struct {
	// certMap is the map of certificate objects.
	certMap map[PdfObjectKey]*model.CertificateToken
	// ocspMap is the map of OCSP objects.
	ocspMap map[PdfObjectKey]*spi.OCSPResponseBinary
	// crlMap is the map of CRL objects.
	crlMap map[PdfObjectKey]*crlparser.CRLBinary
}

// newDssDictExtraction extracts the three token maps from dssDictionary. Port of
// AbstractPdfDssDict's constructor, together with DSSDictionaryExtractionUtils's three
// getXFromArray helpers it calls.
func newDssDictExtraction(dssDictionary PdfDict, dictionaryName, certArrayName, crlArrayName, ocspArrayName string) dssDictExtraction {
	return dssDictExtraction{
		certMap: dssDictionaryExtractionUtilsGetCertsFromArray(dssDictionary, dictionaryName, certArrayName),
		ocspMap: dssDictionaryExtractionUtilsGetOCSPsFromArray(dssDictionary, dictionaryName, ocspArrayName),
		crlMap:  dssDictionaryExtractionUtilsGetCRLsFromArray(dssDictionary, dictionaryName, crlArrayName),
	}
}

// CRLs returns the map of CRL objects. Port of AbstractPdfDssDict#getCRLs.
func (d dssDictExtraction) CRLs() map[PdfObjectKey]*crlparser.CRLBinary { return d.crlMap }

// OCSPs returns the map of OCSP objects. Port of AbstractPdfDssDict#getOCSPs.
func (d dssDictExtraction) OCSPs() map[PdfObjectKey]*spi.OCSPResponseBinary { return d.ocspMap }

// CERTs returns the map of certificate objects. Port of AbstractPdfDssDict#getCERTs.
func (d dssDictExtraction) CERTs() map[PdfObjectKey]*model.CertificateToken { return d.certMap }

// dssDictionaryExtractionUtilsGetCertsFromArray extracts the certificate object map.
// Port of DSSDictionaryExtractionUtils#getCertsFromArray. Upstream logs and skips an entry that
// fails to parse; this stays silent.
func dssDictionaryExtractionUtilsGetCertsFromArray(dict PdfDict, dictionaryName, arrayName string) map[PdfObjectKey]*model.CertificateToken {
	certMap := make(map[PdfObjectKey]*model.CertificateToken)
	certsArray := dict.AsArray(arrayName)
	if certsArray == nil {
		return certMap
	}
	for i := 0; i < certsArray.Size(); i++ {
		objectKey := certsArray.ObjectKey(i)
		if objectKey == nil {
			continue
		}
		if _, found := certMap[objectKey]; found {
			continue
		}
		streamBytes, err := certsArray.StreamBytes(i)
		if err != nil {
			continue
		}
		certToken, err := spi.DSSUtilsLoadCertificateFromBinary(streamBytes)
		if err != nil {
			continue
		}
		certMap[objectKey] = certToken
	}
	return certMap
}

// dssDictionaryExtractionUtilsGetCRLsFromArray extracts the CRL object map.
// Port of DSSDictionaryExtractionUtils#getCRLsFromArray.
func dssDictionaryExtractionUtilsGetCRLsFromArray(dict PdfDict, dictionaryName, arrayName string) map[PdfObjectKey]*crlparser.CRLBinary {
	crlMap := make(map[PdfObjectKey]*crlparser.CRLBinary)
	crlArray := dict.AsArray(arrayName)
	if crlArray == nil {
		return crlMap
	}
	for i := 0; i < crlArray.Size(); i++ {
		objectKey := crlArray.ObjectKey(i)
		if objectKey == nil {
			continue
		}
		if _, found := crlMap[objectKey]; found {
			continue
		}
		streamBytes, err := crlArray.StreamBytes(i)
		if err != nil {
			continue
		}
		crlBinary, err := crlparser.CRLUtilsBuildCRLBinary(streamBytes)
		if err != nil {
			continue
		}
		crlMap[objectKey] = crlBinary
	}
	return crlMap
}

// dssDictionaryExtractionUtilsGetOCSPsFromArray extracts the OCSP object map.
// Port of DSSDictionaryExtractionUtils#getOCSPsFromArray.
func dssDictionaryExtractionUtilsGetOCSPsFromArray(dict PdfDict, dictionaryName, arrayName string) map[PdfObjectKey]*spi.OCSPResponseBinary {
	ocspMap := make(map[PdfObjectKey]*spi.OCSPResponseBinary)
	ocspArray := dict.AsArray(arrayName)
	if ocspArray == nil {
		return ocspMap
	}
	for i := 0; i < ocspArray.Size(); i++ {
		objectKey := ocspArray.ObjectKey(i)
		if objectKey == nil {
			continue
		}
		if _, found := ocspMap[objectKey]; found {
			continue
		}
		streamBytes, err := ocspArray.StreamBytes(i)
		if err != nil {
			continue
		}
		ocspResp, err := spi.NewOCSPRespFromBinaries(streamBytes)
		if err != nil {
			continue
		}
		basicOCSPResp, err := ocspResp.ResponseObject()
		if err != nil || basicOCSPResp == nil {
			continue
		}
		ocspResponseBinary, err := spi.OCSPResponseBinaryBuild(basicOCSPResp)
		if err != nil {
			continue
		}
		ocspMap[objectKey] = ocspResponseBinary
	}
	return ocspMap
}

// dssDictionaryExtractionUtilsGetDictionaryCreationTime returns the 'TU' field value when
// present, nil otherwise. Port of DSSDictionaryExtractionUtils#getDictionaryCreationTime.
func dssDictionaryExtractionUtilsGetDictionaryCreationTime(dict PdfDict) *time.Time {
	value := dict.DateValue(PAdESConstantsTuDictionaryNameVri)
	if value.IsZero() {
		return nil
	}
	return &value
}

// dssDictionaryExtractionUtilsGetTimestampBinaries returns the timestamp binaries extracted
// from the 'TS' field, when present. Port of DSSDictionaryExtractionUtils#getTimestampBinaries.
func dssDictionaryExtractionUtilsGetTimestampBinaries(dict PdfDict) []byte {
	tsDict := dict.AsDict(PAdESConstantsTsDictionaryNameVri)
	if tsDict == nil {
		return nil
	}
	streamBytes, err := tsDict.StreamBytes()
	if err != nil {
		// Upstream logs "Unable to extract 'TS' stream : {}".
		return nil
	}
	return streamBytes
}

// SingleDssDict is a representation of a DSS (Document Security Store) Dictionary embedded in a
// PDF file. The dictionary is unique in a PDF file and can contain: VRI dictionary, certificates
// (Certs), OCSP responses (OCSPs) and CRLs. This dictionary is filled in the PAdES-BASELINE-LT
// extension process. Port of the SingleDssDict class.
type SingleDssDict struct {
	dssDictExtraction

	// vris represents a list of VRI dictionaries incorporated into the DSS dictionary.
	vris []*PdfVriDict
}

// SingleDssDictExtract extracts a DSS dictionary from a PdfDict, nil when absent.
// Port of the static #extract(PdfDict).
func SingleDssDictExtract(documentDict PdfDict) PdfDssDict {
	if documentDict == nil {
		// Upstream logs "No DSS dictionary found".
		return nil
	}
	dssDict := documentDict.AsDict(PAdESConstantsDssDictionaryName)
	if dssDict == nil {
		// Upstream logs "No DSS dictionary found".
		return nil
	}
	return newSingleDssDict(dssDict)
}

// newSingleDssDict builds a SingleDssDict from the /DSS dictionary itself. Port of the
// protected constructor SingleDssDict(PdfDict).
func newSingleDssDict(dssDictionary PdfDict) *SingleDssDict {
	return &SingleDssDict{
		dssDictExtraction: newDssDictExtraction(dssDictionary, PAdESConstantsDssDictionaryName,
			PAdESConstantsCertArrayNameDss, PAdESConstantsCrlArrayNameDss, PAdESConstantsOcspArrayNameDss),
		vris: singleDssDictExtractVRIs(dssDictionary),
	}
}

// singleDssDictExtractVRIs extracts the VRI dictionaries embedded in the DSS dictionary.
// Port of SingleDssDict#extractVRIs. Upstream logs and swallows any exception raised while
// walking the /VRI dictionary; this stays silent.
func singleDssDictExtractVRIs(dssDictionary PdfDict) []*PdfVriDict {
	vriDict := dssDictionary.AsDict(PAdESConstantsVriDictionaryName)
	if vriDict == nil {
		return nil
	}
	names := vriDict.List()
	var result []*PdfVriDict
	for _, name := range names {
		if !spi.DSSUtilsIsSHA1Digest(name) {
			continue
		}
		result = append(result, NewPdfVriDict(name, vriDict.AsDict(name)))
	}
	return result
}

// VRIs returns the list of VRI dictionaries. Port of SingleDssDict#getVRIs.
func (d *SingleDssDict) VRIs() []*PdfVriDict { return d.vris }

// PdfVriDict represents a VRI dictionary. Port of the PdfVriDict class.
type PdfVriDict struct {
	dssDictExtraction

	// name is the VRI dictionary key (SHA-1 value of a signature).
	name string

	// tuTime represents a 'TU' time value, nil when absent.
	tuTime *time.Time

	// tsStream represents a 'TS' timestamp binary value, nil when absent.
	tsStream []byte
}

// NewPdfVriDict builds a VRI dictionary from its key and its content dictionary.
// Port of the PdfVriDict(String, PdfDict) constructor.
func NewPdfVriDict(name string, vriDict PdfDict) *PdfVriDict {
	return &PdfVriDict{
		dssDictExtraction: newDssDictExtraction(vriDict, PAdESConstantsVriDictionaryName,
			PAdESConstantsCertArrayNameVri, PAdESConstantsCrlArrayNameVri, PAdESConstantsOcspArrayNameVri),
		name:     name,
		tuTime:   dssDictionaryExtractionUtilsGetDictionaryCreationTime(vriDict),
		tsStream: dssDictionaryExtractionUtilsGetTimestampBinaries(vriDict),
	}
}

// Name returns the key of the VRI dictionary. Port of #getName.
func (v *PdfVriDict) Name() string { return v.name }

// VRIs is not applicable for a VRI dictionary; it returns nil, the empty-list counterpart of
// Java's Collections.emptyList(). Port of PdfVriDict#getVRIs.
func (v *PdfVriDict) VRIs() []*PdfVriDict { return nil }

// TUTime returns the 'TU' time, nil when absent. Port of #getTUTime.
func (v *PdfVriDict) TUTime() *time.Time { return v.tuTime }

// TSStream returns the 'TS' stream value, nil when absent. Port of #getTSStream.
func (v *PdfVriDict) TSStream() []byte { return v.tsStream }

// PdfDssDictEquals reports whether a and b represent the same DSS dictionary content.
// Port of AbstractPdfDssDict#equals / SingleDssDict#equals / PdfVriDict#equals, collapsed into
// one function since Go interface values compared with == can panic when the dynamic type is
// uncomparable (a PdfDssDict backed by maps and slices is): unlike Java's instance methods,
// which dispatch on the receiver, this free function does the getClass()-equality check itself
// via reflect.TypeOf before comparing exactly the fields the three Java equals() methods do.
//
// It must compare EXACTLY those fields and no others. Java's AbstractPdfDssDict#equals compares
// only certMap/crlMap/ocspMap; SingleDssDict#equals adds vris; PdfVriDict#equals adds name -
// and NEITHER compares a VRI dictionary's 'TU' time or 'TS' stream, nor the underlying PdfDict
// the maps were extracted from. A blanket reflect.DeepEqual over the whole struct (which this
// used to be) is therefore strictly stricter than Java, and that difference is not cosmetic:
// AbstractPDFSignatureService#getRevisions decides whether to emit a PdfDocDssRevision by
// asking whether the previous revision's DSS dictionary equals this one, and the FIRST such
// revision is what flips containsDSSRevisions() - which in turn decides whether an EARLIER
// signature in the same document is handed the document's /DSS dictionary at all. Comparing one
// field too many there makes two identical /DSS dictionaries look different, invents a
// PdfDocDssRevision, and leaks the document's validation data onto a signature that must not
// see it (upstream's own Signature-P-SK-6.pdf and pades-lt-extended-abde.pdf: the second
// signature came out PAdES-BASELINE-LT where upstream says PAdES-BASELINE-T).
//
// The three token maps are compared by their entries' own Java equality - CertificateToken and
// EncapsulatedTokenIdentifier (CRLBinary/OCSPResponseBinary) both define equals() as identity
// of their DSS Id, the digest of the encoded token - rather than by deep structural equality of
// the parsed objects, which would additionally (and wrongly) compare lazily-populated caches.
func PdfDssDictEquals(a, b PdfDssDict) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	// AbstractPdfDssDict#equals: certMap, crlMap and ocspMap.
	if !pdfDssDictTokenMapsEqual(a, b) {
		return false
	}
	switch left := a.(type) {
	case *SingleDssDict:
		// SingleDssDict#equals: super.equals(obj) && vris.equals(other.vris) - a java.util.List
		// equality, i.e. same size, same order, element-wise equals (PdfVriDict#equals below).
		right := b.(*SingleDssDict)
		if len(left.vris) != len(right.vris) {
			return false
		}
		for i := range left.vris {
			if !PdfDssDictEquals(left.vris[i], right.vris[i]) {
				return false
			}
		}
		return true
	case *PdfVriDict:
		// PdfVriDict#equals: super.equals(obj) && name.equals(other.name). Deliberately NOT
		// tuTime/tsStream - see this function's doc comment.
		return left.name == b.(*PdfVriDict).name
	default:
		return true
	}
}

// pdfDssDictTokenMapsEqual ports AbstractPdfDssDict#equals's certMap/crlMap/ocspMap comparison:
// three java.util.Map equalities, i.e. same key set and, per key, values equal by their own
// equals() - the DSS Id for all three value types (see PdfDssDictEquals's doc comment).
func pdfDssDictTokenMapsEqual(a, b PdfDssDict) bool {
	aCerts, bCerts := a.CERTs(), b.CERTs()
	if len(aCerts) != len(bCerts) {
		return false
	}
	for key, left := range aCerts {
		right, found := bCerts[key]
		if !found || left.DSSIDAsString() != right.DSSIDAsString() {
			return false
		}
	}
	aCRLs, bCRLs := a.CRLs(), b.CRLs()
	if len(aCRLs) != len(bCRLs) {
		return false
	}
	for key, left := range aCRLs {
		right, found := bCRLs[key]
		if !found || !left.Equals(&right.MultipleDigestIdentifier) {
			return false
		}
	}
	aOCSPs, bOCSPs := a.OCSPs(), b.OCSPs()
	if len(aOCSPs) != len(bOCSPs) {
		return false
	}
	for key, left := range aOCSPs {
		right, found := bOCSPs[key]
		if !found || !left.Equals(&right.MultipleDigestIdentifier) {
			return false
		}
	}
	return true
}
