// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/vts/ValidationTimeSlidingCertificateRevocationSelector.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and
// past_signature_validation_certificate_revocation_selector.go for the base
// class vpfltvd.LongTermValidationCertificateRevocationSelector this file
// also subclasses.
package vpfswatsp

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfltvd"
)

// ValidationTimeSlidingCertificateRevocationSelector filters revocation data on
// a "Validation Time Sliding" process.
type ValidationTimeSlidingCertificateRevocationSelector struct {
	*vpfltvd.LongTermValidationCertificateRevocationSelector

	// poe is the POE container.
	poe *POEExtraction

	// certificateRevocationData is the list of acceptable certificate
	// revocation data for VTS processing.
	certificateRevocationData []*diagnostic.CertificateRevocationWrapper
}

// NewValidationTimeSlidingCertificateRevocationSelector is the default
// constructor. Port of
// ValidationTimeSlidingCertificateRevocationSelector(Provider, CertificateWrapper, List, Date, Map, String, POEExtraction, ValidationPolicy).
func NewValidationTimeSlidingCertificateRevocationSelector(i18nProvider *i18n.Provider,
	certificate *diagnostic.CertificateWrapper,
	certificateRevocationData []*diagnostic.CertificateRevocationWrapper, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tokenId string, poe *POEExtraction,
	validationPolicy policy.ValidationPolicy) *ValidationTimeSlidingCertificateRevocationSelector {
	c := &ValidationTimeSlidingCertificateRevocationSelector{
		LongTermValidationCertificateRevocationSelector: &vpfltvd.LongTermValidationCertificateRevocationSelector{},
		certificateRevocationData:                       certificateRevocationData,
		poe:                                             poe,
	}
	c.InitLongTermValidationCertificateRevocationSelectorState(i18nProvider, certificate, currentTime,
		nil, bbbs, tokenId, validationPolicy)
	c.InitLongTermValidationCertificateRevocationSelector(c)
	return c
}

// Title returns the title of the building block. Port of the overridden
// getTitle().
func (c *ValidationTimeSlidingCertificateRevocationSelector) Title() i18n.MessageTag {
	return i18n.MessageTagVTSCRS
}

// CertificateRevocationData returns available certificate revocation data to be
// validated. Port of the overridden getCertificateRevocationData().
func (c *ValidationTimeSlidingCertificateRevocationSelector) CertificateRevocationData() []*diagnostic.CertificateRevocationWrapper {
	return c.certificateRevocationData
}

// VerifyRevocationData verifies the given revocation data and returns the
// resulting ChainItem. Port of the overridden
// verifyRevocationData(ChainItem, CertificateRevocationWrapper).
//
// Java's Boolean validity is a nullable box; the Go form is the comma-ok of the
// validity map lookup, so "Boolean.TRUE.equals(validity)" is "present and true".
func (c *ValidationTimeSlidingCertificateRevocationSelector) VerifyRevocationData(
	item process.ChainItem[*jaxb.XmlCRS],
	revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS] {
	item = c.LongTermValidationCertificateRevocationSelector.VerifyRevocationData(item, revocationWrapper)

	validity, ok := c.RevocationDataValidityMap[revocationWrapper.Id()]
	if ok && validity {
		item = item.SetNextItem(c.revocationIssuedBeforeControlTime(
			&revocationWrapper.RevocationWrapper, c.CurrentTime))

		thisUpdate := revocationWrapper.ThisUpdate()
		validity = thisUpdate != nil && thisUpdate.Before(c.CurrentTime)

		if validity {

			item = item.SetNextItem(c.poeExistsAtOrBeforeControlTime(
				c.Certificate, enumerations.TimestampedObjectTypeCertificate, c.CurrentTime))

			item = item.SetNextItem(c.poeExistsAtOrBeforeControlTime(
				revocationWrapper, enumerations.TimestampedObjectTypeRevocation, c.CurrentTime))

			validity = c.poe.IsPOEExists(c.Certificate.Id(), c.CurrentTime) &&
				c.poe.IsPOEExists(revocationWrapper.Id(), c.CurrentTime)

		}

		// update the validity map
		c.RevocationDataValidityMap[revocationWrapper.Id()] = validity
	}

	return item
}

// revocationIssuedBeforeControlTime ports the private
// revocationIssuedBeforeControlTime(RevocationWrapper, Date).
func (c *ValidationTimeSlidingCertificateRevocationSelector) revocationIssuedBeforeControlTime(
	revocation *diagnostic.RevocationWrapper, controlTime time.Time) process.ChainItem[*jaxb.XmlCRS] {
	return NewRevocationIssuedBeforeControlTimeCheck(c.I18nProvider, c.Result, revocation, controlTime,
		c.WarnLevelRule())
}

// poeExistsAtOrBeforeControlTime ports the private
// poeExistsAtOrBeforeControlTime(TokenProxy, TimestampedObjectType, Date).
func (c *ValidationTimeSlidingCertificateRevocationSelector) poeExistsAtOrBeforeControlTime(
	token diagnostic.TokenProxy, objectType enumerations.TimestampedObjectType,
	controlTime time.Time) process.ChainItem[*jaxb.XmlCRS] {
	return NewPOEExistsAtOrBeforeControlTimeCheck(c.I18nProvider, c.Result, token, objectType, controlTime, c.poe,
		c.WarnLevelRule())
}

// AcceptableRevocationDataAvailable checks whether the acceptable revocation
// data is available. Port of the overridden acceptableRevocationDataAvailable():
//
// If at least one revocation status information is selected, the building block
// shall go to the next step. If there is no such information, the building
// block shall return the indication INDETERMINATE with the sub-indication
// NO_POE.
func (c *ValidationTimeSlidingCertificateRevocationSelector) AcceptableRevocationDataAvailable() process.ChainItem[*jaxb.XmlCRS] {
	// Java hands the CertificateRevocationWrapper straight to a RevocationWrapper
	// parameter; the Go wrapper embeds its base by value, so the base is addressed
	// out of it, and a null latest revocation stays a nil *RevocationWrapper.
	var acceptableRevocationData *diagnostic.RevocationWrapper
	if latest := c.LatestAcceptableCertificateRevocation(); latest != nil {
		acceptableRevocationData = &latest.RevocationWrapper
	}
	return newVTSAcceptableRevocationDataAvailableCheck(c.I18nProvider, c.Result, acceptableRevocationData,
		c.FailLevelRule())
}

// vtsAcceptableRevocationDataAvailableCheck is the Go form of the anonymous
// AcceptableRevocationDataAvailableCheck<XmlCRS> subclass declared inside
// acceptableRevocationDataAvailable().
type vtsAcceptableRevocationDataAvailableCheck struct {
	*xcv.AcceptableRevocationDataAvailableCheck[*jaxb.XmlCRS]
}

// newVTSAcceptableRevocationDataAvailableCheck instantiates the anonymous
// subclass.
func newVTSAcceptableRevocationDataAvailableCheck(i18nProvider *i18n.Provider,
	result *process.Result[*jaxb.XmlCRS], acceptableRevocationData *diagnostic.RevocationWrapper,
	constraint policy.LevelRule) *vtsAcceptableRevocationDataAvailableCheck {
	c := &vtsAcceptableRevocationDataAvailableCheck{
		AcceptableRevocationDataAvailableCheck: xcv.NewAcceptableRevocationDataAvailableCheck(
			i18nProvider, result, acceptableRevocationData, constraint),
	}
	// Re-register with the outer type so the two overrides below dispatch.
	c.InitChainItem(c)
	return c
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the anonymous subclass's getFailedIndicationForConclusion().
func (c *vtsAcceptableRevocationDataAvailableCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the anonymous subclass's getFailedSubIndicationForConclusion().
func (c *vtsAcceptableRevocationDataAvailableCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationNoPOE
}
