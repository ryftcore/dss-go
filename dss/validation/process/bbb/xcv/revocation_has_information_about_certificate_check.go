// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rac/checks/RevocationHasInformationAboutCertificateCheck.java (DSS 6.5.RC1).
//
// The slf4j LOG.info records of the three "is not before revocation's thisUpdate"
// branches are dropped per PORTING.md; the branches themselves keep their (empty)
// bodies below so the control flow stays statement-for-statement.
//
// Java declares getNotAfterAfterCertificateNotAfterMessage() protected while the
// five sibling message builders are private. No class in the port extends this
// check, so all six are unexported methods here; a future subclass would export
// that one and route the buildAdditionalInfo self-call through an overrides
// interface, the way Chain and ChainItem do.
package xcv

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// RevocationHasInformationAboutCertificateCheck checks whether the concerned
// certificate has existed at the time of revocation data generation.
type RevocationHasInformationAboutCertificateCheck struct {
	*process.ChainItemBase[*jaxb.XmlRAC]

	// certificate is the certificate in question.
	certificate *diagnostic.CertificateWrapper

	// revocationData is the revocation data to check.
	revocationData *diagnostic.RevocationWrapper

	// notAfterRevoc defines the date after which the revocation issuer ensures
	// the revocation is contained for the certificate. It is computed lazily,
	// exactly as the Java field is: a computation that answers null is repeated
	// on the next call.
	notAfterRevoc *time.Time
}

// NewRevocationHasInformationAboutCertificateCheck is the default constructor.
// Port of RevocationHasInformationAboutCertificateCheck(I18nProvider, XmlRAC,
// CertificateWrapper, RevocationWrapper, LevelRule).
func NewRevocationHasInformationAboutCertificateCheck(i18nProvider *i18n.I18nProvider,
	result *process.Result[*jaxb.XmlRAC], certificate *diagnostic.CertificateWrapper,
	revocationData *diagnostic.RevocationWrapper,
	constraint policy.LevelRule) *RevocationHasInformationAboutCertificateCheck {
	c := &RevocationHasInformationAboutCertificateCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		certificate:    certificate,
		revocationData: revocationData,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *RevocationHasInformationAboutCertificateCheck) Process() bool {
	return c.checkCertHashMatches() || c.checkIssuerHasInformationForExpiredCertificate()
}

// checkIssuerHasInformationForExpiredCertificate ports the private
// checkIssuerHasInformationForExpiredCertificate().
func (c *RevocationHasInformationAboutCertificateCheck) checkIssuerHasInformationForExpiredCertificate() bool {
	certNotAfter := c.certificate.NotAfter()
	revocationIssuerKnowsCertStatusSince := c.getNotAfterRevoc()
	return certNotAfter != nil && revocationIssuerKnowsCertStatusSince != nil &&
		!certNotAfter.Before(*revocationIssuerKnowsCertStatusSince)
}

// checkCertHashMatches ports the private checkCertHashMatches().
func (c *RevocationHasInformationAboutCertificateCheck) checkCertHashMatches() bool {
	/*
	 * certHash extension can be present in an OCSP Response. If present, a digest match indicates the OCSP
	 * responder knows the certificate as we have it, and so also its revocation state
	 */
	return c.revocationData.IsCertHashExtensionPresent() && c.revocationData.IsCertHashExtensionMatch()
}

// getNotAfterRevoc ports the private getNotAfterRevoc().
func (c *RevocationHasInformationAboutCertificateCheck) getNotAfterRevoc() *time.Time {
	if c.notAfterRevoc == nil {
		c.notAfterRevoc = c.revocationData.ThisUpdate()

		/*
		 * If a CRL contains the extension expiredCertsOnCRL defined in [i.12], it shall prevail over the TL
		 * extension value but only for that specific CRL.
		 */
		expiredCertsOnCRL := c.revocationData.ExpiredCertsOnCRL()
		if expiredCertsOnCRL != nil {
			if expiredCertsOnCRL.Before(*c.notAfterRevoc) {
				c.notAfterRevoc = expiredCertsOnCRL
			}
			// else: upstream only logs
		}

		/*
		 * If an OCSP response contains the extension ArchiveCutoff defined in section 4.4.4 of
		 * IETF RFC 6960 [i.11], it shall prevail over the TL extension value but only for that specific OCSP
		 * response.
		 */
		archiveCutOff := c.revocationData.ArchiveCutOff()
		if archiveCutOff != nil {
			if archiveCutOff.Before(*c.notAfterRevoc) {
				c.notAfterRevoc = archiveCutOff
			}
			// else: upstream only logs
		}

		/* expiredCertsRevocationInfo Extension from TL */
		if expiredCertsOnCRL == nil && archiveCutOff == nil {
			expiredCertsRevocationInfo := c.getExpiredCertsRevocationInfo(c.revocationData)
			if expiredCertsRevocationInfo != nil {
				if expiredCertsRevocationInfo.Before(*c.notAfterRevoc) {
					c.notAfterRevoc = expiredCertsRevocationInfo
				}
				// else: upstream only logs
			}
		}
	}
	return c.notAfterRevoc
}

// getExpiredCertsRevocationInfo ports the private
// getExpiredCertsRevocationInfo(RevocationWrapper).
func (c *RevocationHasInformationAboutCertificateCheck) getExpiredCertsRevocationInfo(
	revocationData *diagnostic.RevocationWrapper) *time.Time {
	revocCert := revocationData.SigningCertificate()
	if revocCert != nil {
		return revocCert.CertificateTSPServiceExpiredCertsRevocationInfo()
	}
	return nil
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *RevocationHasInformationAboutCertificateCheck) BuildAdditionalInfo() *string {
	var message string
	if !c.Process() {
		message = c.getNotAfterAfterCertificateNotAfterMessage()

	} else if c.checkRevocationThisUpdateIsInCertificateValidityRange() {
		message = c.getRevocationConsistentMessage()

	} else if c.checkCertHashMatches() {
		message = c.getRevocationCertHashOkMessage()

	} else if c.checkExpiredCertsOnCRLPresent() {
		message = c.getRevocationConsistentWithExpiredCertsOnCRLMessage()

	} else if c.checkArchiveCutOffPresent() {
		message = c.getRevocationConsistentWithArchiveCutoffMessage()

	} else if c.checkExpiredCertsRevocationInfoPresent() {
		message = c.getRevocationConsistentWithExpiredCertsRevocationInfoMessage()

	} else {
		message = c.getRevocationInfoMessage()
	}
	return &message
}

// checkRevocationThisUpdateIsInCertificateValidityRange ports the private
// checkRevocationThisUpdateIsInCertificateValidityRange().
func (c *RevocationHasInformationAboutCertificateCheck) checkRevocationThisUpdateIsInCertificateValidityRange() bool {
	thisUpdate := c.revocationData.ThisUpdate()
	return !thisUpdate.Before(*c.certificate.NotBefore()) &&
		!thisUpdate.After(*c.certificate.NotAfter())
}

// checkExpiredCertsOnCRLPresent ports the private checkExpiredCertsOnCRLPresent().
func (c *RevocationHasInformationAboutCertificateCheck) checkExpiredCertsOnCRLPresent() bool {
	return c.revocationData.ExpiredCertsOnCRL() != nil
}

// checkArchiveCutOffPresent ports the private checkArchiveCutOffPresent().
func (c *RevocationHasInformationAboutCertificateCheck) checkArchiveCutOffPresent() bool {
	return c.revocationData.ArchiveCutOff() != nil
}

// checkExpiredCertsRevocationInfoPresent ports the private
// checkExpiredCertsRevocationInfoPresent().
func (c *RevocationHasInformationAboutCertificateCheck) checkExpiredCertsRevocationInfoPresent() bool {
	return c.getExpiredCertsRevocationInfo(c.revocationData) != nil
}

// getNotAfterAfterCertificateNotAfterMessage returns the additional information
// message in case if computed time 'notAfter' is after the certificate's
// notAfter. Port of getNotAfterAfterCertificateNotAfterMessage(), which reads
// the field and not the lazy getter.
func (c *RevocationHasInformationAboutCertificateCheck) getNotAfterAfterCertificateNotAfterMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_NOT_AFTER_AFTER,
		c.formattedDate(c.notAfterRevoc),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// getRevocationConsistentMessage returns the additional information message when
// the revocation is consistent. Port of getRevocationConsistentMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationConsistentMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CONSISTENT,
		c.formattedDate(c.revocationData.ThisUpdate()),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// getRevocationCertHashOkMessage returns the additional information message when
// certHash matches. Port of getRevocationCertHashOkMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationCertHashOkMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CERT_HASH_OK)
}

// getRevocationConsistentWithExpiredCertsOnCRLMessage returns the additional
// information message when the revocation is consistent with expiredCertsOnCRL.
// Port of getRevocationConsistentWithExpiredCertsOnCRLMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationConsistentWithExpiredCertsOnCRLMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CONSISTENT_CRL,
		c.formattedDate(c.revocationData.ThisUpdate()),
		c.formattedDate(c.revocationData.ExpiredCertsOnCRL()),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// getRevocationConsistentWithArchiveCutoffMessage returns the additional
// information message when the revocation is consistent with archiveCutoff.
// Port of getRevocationConsistentWithArchiveCutoffMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationConsistentWithArchiveCutoffMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CONSISTENT_OCSP,
		c.formattedDate(c.revocationData.ThisUpdate()),
		c.formattedDate(c.revocationData.ArchiveCutOff()),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// getRevocationConsistentWithExpiredCertsRevocationInfoMessage returns the
// additional information message when the revocation is consistent with the
// trusted list's expiredCertsRevocationInfo. Port of
// getRevocationConsistentWithExpiredCertsRevocationInfoMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationConsistentWithExpiredCertsRevocationInfoMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_CONSISTENT_TL,
		c.formattedDate(c.revocationData.ThisUpdate()),
		c.formattedDate(c.getExpiredCertsRevocationInfo(c.revocationData)),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// getRevocationInfoMessage returns the additional information message for
// revocation data in case of other events. Port of getRevocationInfoMessage().
func (c *RevocationHasInformationAboutCertificateCheck) getRevocationInfoMessage() string {
	return c.I18nProvider.GetMessage(i18n.MessageTag_REVOCATION_INFO,
		c.formattedDate(c.revocationData.ThisUpdate()),
		c.formattedDate(c.certificate.NotBefore()),
		c.formattedDate(c.certificate.NotAfter()))
}

// formattedDate renders a date as an I18nProvider argument the way Java does.
//
// ValidationProcessUtils#getFormattedDate answers null for a null Date, and
// java.text.MessageFormat renders that null as the four characters "null"; the
// Go port of that helper answers the empty string instead (a deliberate choice
// recorded in its header, so that the result stays usable as an argument), which
// loses those four characters from the message. This check reaches the case -
// notAfterRevoc is null whenever the revocation carries no thisUpdate, which is
// exactly what ThisUpdatePresenceCheck exists to catch - so it restores the Java
// rendering here.
func (c *RevocationHasInformationAboutCertificateCheck) formattedDate(date *time.Time) string {
	if date == nil {
		return "null"
	}
	return process.GetFormattedDate(date)
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *RevocationHasInformationAboutCertificateCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *RevocationHasInformationAboutCertificateCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_REVOC_HAS_CERT_INFO_ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *RevocationHasInformationAboutCertificateCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of getFailedSubIndicationForConclusion().
func (c *RevocationHasInformationAboutCertificateCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CERTIFICATE_CHAIN_GENERAL_FAILURE
}
