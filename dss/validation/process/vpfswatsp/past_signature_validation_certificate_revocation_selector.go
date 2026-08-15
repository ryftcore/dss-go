// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/checks/psv/PastSignatureValidationCertificateRevocationSelector.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// CROSS-CHUNK DEPENDENCY (phase 8e): the Java base class
// eu.europa.esig.dss.validation.process.vpfltvd.LongTermValidationCertificateRevocationSelector
// belongs to the vpfltvd chunk, ported in parallel. This file is written
// against the Go shape that chunk's other classes dictate, which is the one
// bbb/xcv's CertificateRevocationSelector already established for a
// Chain-derived, extension-designed class:
//
//	type LongTermValidationCertificateRevocationSelectorOverrides interface {
//	    xcv.CertificateRevocationSelectorOverrides
//	    RevocationBBBConclusion(*diagnostic.CertificateRevocationWrapper) *jaxb.XmlConclusion
//	}
//	type LongTermValidationCertificateRevocationSelector struct {
//	    *xcv.CertificateRevocationSelector
//	    Bbbs    map[string]*jaxb.XmlBasicBuildingBlocks // Java: protected bbbs
//	    TokenId string                                  // Java: protected tokenId
//	}
//	func (s *LongTermValidationCertificateRevocationSelector) InitLongTermValidationCertificateRevocationSelectorState(
//	    i18nProvider *i18n.I18nProvider, certificate *diagnostic.CertificateWrapper, currentTime time.Time,
//	    diagnosticData *diagnostic.DiagnosticData, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
//	    tokenId string, validationPolicy policy.ValidationPolicy)
//	func (s *LongTermValidationCertificateRevocationSelector) InitLongTermValidationCertificateRevocationSelector(
//	    overrides LongTermValidationCertificateRevocationSelectorOverrides)
//	func (s *LongTermValidationCertificateRevocationSelector) VerifyRevocationData(
//	    item process.ChainItem[*jaxb.XmlCRS], revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS]
//
// The protected Java constructor that passes a null DiagnosticData is the state
// initializer called with a nil one below.
package vpfswatsp

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb/xcv"
	"github.com/utain/esig/dss/validation/process/vpfltvd"
)

// PastSignatureValidationCertificateRevocationSelector filters revocation data
// on a "Past Signature Validation" process.
type PastSignatureValidationCertificateRevocationSelector struct {
	*vpfltvd.LongTermValidationCertificateRevocationSelector

	// poe is the POE container.
	poe *POEExtraction

	// acceptableCertificateRevocations is a list of acceptable revocation data
	// for the given certificate.
	acceptableCertificateRevocations []*diagnostic.CertificateRevocationWrapper
}

// NewPastSignatureValidationCertificateRevocationSelector is the default
// constructor. Port of
// PastSignatureValidationCertificateRevocationSelector(I18nProvider, CertificateWrapper, Date, Map, String, POEExtraction, ValidationPolicy).
func NewPastSignatureValidationCertificateRevocationSelector(i18nProvider *i18n.I18nProvider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, tokenId string, poe *POEExtraction,
	validationPolicy policy.ValidationPolicy) *PastSignatureValidationCertificateRevocationSelector {
	c := &PastSignatureValidationCertificateRevocationSelector{
		LongTermValidationCertificateRevocationSelector: &vpfltvd.LongTermValidationCertificateRevocationSelector{},
		poe:                              poe,
		acceptableCertificateRevocations: make([]*diagnostic.CertificateRevocationWrapper, 0),
	}
	c.InitLongTermValidationCertificateRevocationSelectorState(i18nProvider, certificate, currentTime,
		nil, bbbs, tokenId, validationPolicy)
	c.InitLongTermValidationCertificateRevocationSelector(c)
	return c
}

// Title returns the title of the building block. Port of the overridden
// getTitle().
func (c *PastSignatureValidationCertificateRevocationSelector) Title() i18n.MessageTag {
	return i18n.MessageTag_PSV_CRS
}

// VerifyRevocationData verifies the given revocation data and returns the
// resulting ChainItem. Port of the overridden
// verifyRevocationData(ChainItem, CertificateRevocationWrapper).
//
// Java's Boolean validity is a nullable box; the Go form is the comma-ok of the
// validity map lookup, so "Boolean.TRUE.equals(validity)" is "present and true".
func (c *PastSignatureValidationCertificateRevocationSelector) VerifyRevocationData(
	item process.ChainItem[*jaxb.XmlCRS],
	revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS] {
	item = c.LongTermValidationCertificateRevocationSelector.VerifyRevocationData(item, revocationWrapper)

	validity, ok := c.RevocationDataValidityMap[revocationWrapper.Id()]
	if ok && validity {
		revocationIssuer := revocationWrapper.SigningCertificate()

		if revocationIssuer != nil {

			if c.isRevocationIssuerTrusted(revocationIssuer) {

				item = item.SetNextItem(c.revocationDataIssuerTrusted(revocationIssuer))

				c.acceptableCertificateRevocations = append(c.acceptableCertificateRevocations, revocationWrapper)
				c.addAcceptableRevocationId(revocationWrapper.Id())

				c.RevocationDataValidityMap[revocationWrapper.Id()] = true

			} else {

				item = item.SetNextItem(c.poeForRevocationDataIssuerExists(revocationIssuer))

				validity = c.poe.IsPOEExistInRange(revocationIssuer.Id(),
					revocationIssuer.NotBefore(), revocationIssuer.NotAfter())

				if validity {
					c.acceptableCertificateRevocations = append(c.acceptableCertificateRevocations, revocationWrapper)
					c.addAcceptableRevocationId(revocationWrapper.Id())
				}

				// update the validity map
				c.RevocationDataValidityMap[revocationWrapper.Id()] = validity

			}

		}
	}

	return item
}

// addAcceptableRevocationId ports XmlCRS#getAcceptableRevocationId().add(String):
// the generated member is a *StringList, which JAXB's getter would have created
// on first access.
func (c *PastSignatureValidationCertificateRevocationSelector) addAcceptableRevocationId(revocationId string) {
	if c.Result.Value.AcceptableRevocationId == nil {
		c.Result.Value.AcceptableRevocationId = &jaxb.StringList{}
	}
	*c.Result.Value.AcceptableRevocationId = append(*c.Result.Value.AcceptableRevocationId, revocationId)
}

// RevocationBBBConclusion returns a conclusion of the revocation basic building
// block execution process. Port of the overridden
// getRevocationBBBConclusion(CertificateRevocationWrapper).
func (c *PastSignatureValidationCertificateRevocationSelector) RevocationBBBConclusion(
	revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlConclusion {
	revocationBBB, ok := c.Bbbs[revocationWrapper.Id()]
	if ok && revocationBBB != nil {
		return revocationBBB.Conclusion
	}
	return nil
}

// revocationDataIssuerTrusted ports the private
// revocationDataIssuerTrusted(CertificateWrapper).
func (c *PastSignatureValidationCertificateRevocationSelector) revocationDataIssuerTrusted(
	revocationIssuer *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlCRS] {
	sunsetDateConstraint := c.ValidationPolicy.CertificateSunsetDateConstraint(
		enumerations.Context_REVOCATION, enumerations.SubContext_SIGNING_CERT)
	return xcv.NewRevocationIssuerTrustedCheck(c.I18nProvider, c.Result, revocationIssuer, c.CurrentTime,
		sunsetDateConstraint, c.WarnLevelRule())
}

// poeForRevocationDataIssuerExists ports the private
// poeForRevocationDataIssuerExists(CertificateWrapper).
func (c *PastSignatureValidationCertificateRevocationSelector) poeForRevocationDataIssuerExists(
	revocationIssuer *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlCRS] {
	return NewPOEExistsWithinCertificateValidityRangeCheck(c.I18nProvider, c.Result, revocationIssuer, c.poe,
		c.WarnLevelRule())
}

// AcceptableRevocationDataAvailable checks whether the acceptable revocation
// data is available. Port of the overridden acceptableRevocationDataAvailable().
func (c *PastSignatureValidationCertificateRevocationSelector) AcceptableRevocationDataAvailable() process.ChainItem[*jaxb.XmlCRS] {
	return NewPastValidationAcceptableRevocationDataAvailable(c.I18nProvider, c.Result,
		c.acceptableCertificateRevocations, c.Result.Value.RAC, c.FailLevelRule())
}

// isRevocationIssuerTrusted ports the private
// isRevocationIssuerTrusted(CertificateWrapper).
func (c *PastSignatureValidationCertificateRevocationSelector) isRevocationIssuerTrusted(
	certificateWrapper *diagnostic.CertificateWrapper) bool {
	constraint := c.ValidationPolicy.CertificateSunsetDateConstraint(
		enumerations.Context_REVOCATION, enumerations.SubContext_SIGNING_CERT)
	return process.IsTrustAnchor(certificateWrapper, c.CurrentTime, constraint)
}

// AcceptableCertificateRevocations returns a list of acceptable certificate
// revocation data in the past validation process. Port of
// getAcceptableCertificateRevocations().
func (c *PastSignatureValidationCertificateRevocationSelector) AcceptableCertificateRevocations() []*diagnostic.CertificateRevocationWrapper {
	return c.acceptableCertificateRevocations
}
