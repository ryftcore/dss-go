// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/LOTLParsingTask.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see this batch's porter notes): OtherTSLPointerConverter,
// PivotSchemeInformationURI and LOTLSigningCertificatesAnnouncementSchemeInformationURI live in
// dss-tsl-validation's "function" package, which the TSLJOB chunk ports into this same Go package
// (tsl). They are referenced here by their Java names, with this codebase's constructor
// (New<Name>) and functional-interface (Test / Apply) spellings.
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// LOTLParsingTask parses a LOTL and returns a LOTLParsingResult.
type LOTLParsingTask struct {
	AbstractParsingTaskBase

	// lotlSource is the LOTLSource to parse.
	lotlSource *LOTLSource
}

var _ AbstractParsingTaskOverrides = (*LOTLParsingTask)(nil)

// NewLOTLParsingTask is the default constructor, taking the LOTL document to parse and its
// LOTLSource. Port of LOTLParsingTask(DSSDocument, LOTLSource).
//
// Panics with the Java messages when either argument is nil (Objects.requireNonNull; the document
// check lives in the superclass constructor).
func NewLOTLParsingTask(document model.DSSDocument, lotlSource *LOTLSource) *LOTLParsingTask {
	task := &LOTLParsingTask{AbstractParsingTaskBase: NewAbstractParsingTaskBase(document)}
	if lotlSource == nil {
		panic("The LOTLSource is null")
	}
	task.lotlSource = lotlSource
	task.InitAbstractParsingTask(task)
	return task
}

// Get parses the LOTL. Port of the get() override (Supplier#get); the exceptions Java's
// getJAXBObject and verifyTLVersionConformity raise become returned errors, per PORTING.md.
//
// NOTE: Java's covariant return type (LOTLParsingResult) is kept - see tl_parsing_task.go.
func (t *LOTLParsingTask) Get() (*LOTLParsingResult, error) {
	result := NewLOTLParsingResult()
	jaxbObject, err := t.JAXBObject()
	if err != nil {
		return nil, err
	}

	t.parseSchemeInformation(result, jaxbObject.SchemeInformation)
	if err := t.VerifyTLVersionConformity(result, result.Version(), t.lotlSource.TLVersions()); err != nil {
		return nil, err
	}

	return result, nil
}

// parseSchemeInformation ports the private parseSchemeInformation(LOTLParsingResult,
// TSLSchemeInformationType).
func (t *LOTLParsingTask) parseSchemeInformation(result *LOTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	t.CommonParseSchemeInformation(&result.AbstractTLParsingResult, schemeInformation)
	// NOTE: upstream calls the two extractors below unconditionally, i.e. even for a null
	// schemeInformation - which commonParseSchemeInformation guards against but they do not, so
	// a Trusted List missing the (schema-required) SchemeInformation element raises a
	// NullPointerException there. The Go port dereferences the same nil pointer and panics for
	// the same input; no guard is added, so the two implementations agree.
	t.extractOtherTSLPointers(result, schemeInformation)
	t.extractSchemeInformationURI(result, schemeInformation)
}

// extractOtherTSLPointers ports the private extractOtherTSLPointers(LOTLParsingResult,
// TSLSchemeInformationType).
func (t *LOTLParsingTask) extractOtherTSLPointers(result *LOTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	otherTSLPointersType := schemeInformation.PointersToOtherTSL
	if otherTSLPointersType != nil && utils.IsCollectionNotEmpty(otherTSLPointersType.OtherTSLPointer) {
		otherTSLPointers := otherTSLPointersType.OtherTSLPointer
		converter := NewOtherTSLPointerConverter(t.lotlSource.IsMraSupport())

		lotlPointers := make([]*tslmodel.OtherTSLPointer, 0, len(otherTSLPointers))
		for _, otherTSLPointer := range otherTSLPointers {
			if t.lotlSource.LotlPredicate().Test(otherTSLPointer) {
				lotlPointers = append(lotlPointers, converter.Apply(otherTSLPointer))
			}
		}
		result.SetLotlPointers(lotlPointers)

		tlPointers := make([]*tslmodel.OtherTSLPointer, 0, len(otherTSLPointers))
		for _, otherTSLPointer := range otherTSLPointers {
			if t.lotlSource.TlPredicate().Test(otherTSLPointer) {
				tlPointers = append(tlPointers, converter.Apply(otherTSLPointer))
			}
		}
		result.SetTlPointers(tlPointers)
	}
}

// extractSchemeInformationURI ports the private extractSchemeInformationURI(LOTLParsingResult,
// TSLSchemeInformationType).
func (t *LOTLParsingTask) extractSchemeInformationURI(result *LOTLParsingResult,
	schemeInformation *jaxb.TSLSchemeInformationType) {
	schemeInformationURI := schemeInformation.SchemeInformationURI
	if schemeInformationURI != nil {
		t.extractSigningCertificatesAnnouncementURL(result, schemeInformationURI)
		t.extractPivotURLs(result, schemeInformationURI)
	}
}

// extractSigningCertificatesAnnouncementURL ports the private
// extractSigningCertificatesAnnouncementURL(LOTLParsingResult, NonEmptyMultiLangURIListType).
func (t *LOTLParsingTask) extractSigningCertificatesAnnouncementURL(result *LOTLParsingResult,
	schemeInformationURI *jaxb.NonEmptyMultiLangURIListType) {
	signingCertificatesAnnouncementPredicate := t.lotlSource.SigningCertificatesAnnouncementPredicate()
	if signingCertificatesAnnouncementPredicate != nil {
		var uris []string
		for _, uri := range schemeInformationURI.URI {
			if signingCertificatesAnnouncementPredicate.Test(uri) {
				uris = append(uris, uri.Value)
			}
		}
		if utils.IsCollectionNotEmpty(uris) {
			// Upstream compares uris.get(0) against the predicate's own URI only to log
			// "LOTLSigningCertificatesAnnouncement URI change detected. New URI : {}"; with
			// the logging dropped the comparison has no effect, so only the assignment
			// survives.
			result.SetSigningCertificateAnnouncementURL(uris[0])
		}
	}
}

// extractPivotURLs ports the private extractPivotURLs(LOTLParsingResult,
// NonEmptyMultiLangURIListType), including its break as soon as the OJ URL is reached.
func (t *LOTLParsingTask) extractPivotURLs(result *LOTLParsingResult,
	schemeInformationURI *jaxb.NonEmptyMultiLangURIListType) {
	if t.lotlSource.IsPivotSupport() {
		signCertAnnouncementPredicate := t.lotlSource.SigningCertificatesAnnouncementPredicate()
		signCertAnnouncementURL := ""
		hasSignCertAnnouncementURL := false
		if signCertAnnouncementPredicate != nil {
			signCertAnnouncementURL = signCertAnnouncementPredicate.Uri()
			hasSignCertAnnouncementURL = true
		}

		filteredPivots := []string{}
		pivotSchemeInformationURI := NewPivotSchemeInformationURI()
		for _, nonEmptyMultiLangURIType := range schemeInformationURI.URI {
			if pivotSchemeInformationURI.Test(nonEmptyMultiLangURIType) {
				filteredPivots = append(filteredPivots, nonEmptyMultiLangURIType.Value)
			}
			// check if OJ URL is reached
			if hasSignCertAnnouncementURL && signCertAnnouncementURL == nonEmptyMultiLangURIType.Value {
				break
			}
		}
		result.SetPivotURLs(filteredPivots)
	}
}

// CreateTrustedListFacade loads the MRA facade when the LOTL source supports MRA, and the plain
// Trusted List facade otherwise. Port of the protected createTrustedListFacade() override.
func (t *LOTLParsingTask) CreateTrustedListFacade() trustedListFacade {
	if t.lotlSource.IsMraSupport() {
		return trustedlist.NewMRAFacade()
	}
	return t.AbstractParsingTaskBase.CreateTrustedListFacade()
}
