// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/eaa/checks/ETSI194721ConformanceCheck.java (DSS 6.5.RC1).
package checks

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// ETSI194721ConformanceCheck verifies whether the issuing authority
// identifier is valid as per TS 119 472-1.
type ETSI194721ConformanceCheck struct {
	*process.ChainItemBase[*jaxb.XmlSAV]

	// eaa is the EAA to check.
	eaa *diagnostic.EAAWrapper

	// validationTime is the validation time.
	validationTime time.Time
}

// NewETSI194721ConformanceCheck is the default constructor.
func NewETSI194721ConformanceCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlSAV],
	eaaWrapper *diagnostic.EAAWrapper, validationTime time.Time, constraint policy.LevelRule) *ETSI194721ConformanceCheck {
	c := &ETSI194721ConformanceCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		eaa:            eaaWrapper,
		validationTime: validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *ETSI194721ConformanceCheck) Process() bool {
	return c.checkVCTPresent() &&
		c.checkVCTIntegrityPresent() &&
		c.checkNowAfterNotBefore() &&
		c.checkNowBeforeExpiration() &&
		c.checkNowAfterAdministrativeDateIssuance() &&
		c.checkNowBeforeAdministrativeDateExpiration() &&
		c.checkSDJWTAdministrativeDateConformance() &&
		c.checkMDOCDocumentNumberPresent() &&
		c.checkMDOCIssuingAuthorityPresent() &&
		c.checkSDJWTIssuingAuthorityAndCountryPresent() &&
		c.checkNoStatusIfShortLived() &&
		c.checkStatusIsPresentIfMandatory() &&
		c.checkSDJWTStatusConformance()
}

func (c *ETSI194721ConformanceCheck) checkVCTPresent() bool {
	if enumerations.EAATypeSDJWTVC == c.eaa.EAAType() {
		return c.eaa.EAAVerifiableCredentialsTypeUri() != ""
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkVCTIntegrityPresent() bool {
	if enumerations.EAATypeSDJWTVC == c.eaa.EAAType() {
		return c.eaa.EAAVerifiableCredentialsTypeIntegrityBytes() != nil
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkSDJWTIssuingAuthorityAndCountryPresent() bool {
	if enumerations.EAATypeSDJWTVC == c.eaa.EAAType() {
		eaaSignature := c.eaa.EAASignatures()[0]
		signingCertificate := eaaSignature.SigningCertificate()
		relatedCertificates := eaaSignature.FoundCertificates().RelatedCertificates()

		signCertPresent := signingCertificate != nil && utils.IsCollectionNotEmpty(relatedCertificates) && certificateIn(signingCertificate, relatedCertificates)
		if signCertPresent {
			if signingCertificate.IsQcCompliance() {
				return c.eaa.DocumentIssuingAuthority() == "" && c.eaa.DocumentIssuingAuthorityCountry() == ""
			}
		} else if enumerations.EAAQualificationQEAA == c.eaa.CategoryQualification() ||
			enumerations.EAAQualificationPubEAA == c.eaa.CategoryQualification() {
			// NOTE: TS 119 472-1 v1.2.1 expects a QC for a QEAA/PubEAA, but does
			// not define how to proceed for a not QC. Therefore we accept any
			// certificate in such a case.
			return c.eaa.DocumentIssuingAuthority() != "" && c.eaa.DocumentIssuingAuthorityCountry() != ""
		}
	}

	return true
}

// certificateIn ports the Java stream check
// `relatedCertificates.stream().anyMatch(c -> signingCertificate.getId().equals(c.getId()))`.
func certificateIn(signingCertificate *diagnostic.CertificateWrapper, relatedCertificates []*diagnostic.RelatedCertificateWrapper) bool {
	for _, rc := range relatedCertificates {
		if signingCertificate.Id() == rc.Id() {
			return true
		}
	}
	return false
}

func (c *ETSI194721ConformanceCheck) checkMDOCIssuingAuthorityPresent() bool {
	if enumerations.EAATypeISOIECMDoc == c.eaa.EAAType() {
		return c.eaa.DocumentIssuingAuthority() != ""
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkMDOCDocumentNumberPresent() bool {
	if enumerations.EAATypeISOIECMDoc == c.eaa.EAAType() {
		return c.eaa.DocumentNumber() != ""
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkSDJWTAdministrativeDateConformance() bool {
	if enumerations.EAATypeSDJWTVC == c.eaa.EAAType() {
		return (c.eaa.AdministrativeIssuanceDate() == nil) == (c.eaa.AdministrativeExpirationDate() == nil)
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkNowAfterAdministrativeDateIssuance() bool {
	if c.eaa.AdministrativeIssuanceDate() != nil {
		return !c.validationTime.Before(*c.eaa.AdministrativeIssuanceDate())
	}
	// Administrative date is optional, return true if not present.
	return true
}

func (c *ETSI194721ConformanceCheck) checkNowBeforeAdministrativeDateExpiration() bool {
	if c.eaa.AdministrativeExpirationDate() != nil {
		return c.validationTime.Before(*c.eaa.AdministrativeExpirationDate())
	}
	// Administrative date is optional, return true if not present.
	return true
}

func (c *ETSI194721ConformanceCheck) checkNowAfterNotBefore() bool {
	return c.eaa.EAANotBefore() != nil && !c.validationTime.Before(*c.eaa.EAANotBefore())
}

func (c *ETSI194721ConformanceCheck) checkNowBeforeExpiration() bool {
	return c.eaa.EAAExpiration() != nil && c.validationTime.Before(*c.eaa.EAAExpiration())
}

func (c *ETSI194721ConformanceCheck) checkNoStatusIfShortLived() bool {
	if utils.IsTrue(c.eaa.ShortLived()) {
		return c.eaa.EAAPayload().EAAStatus() == nil
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkStatusIsPresentIfMandatory() bool {
	if (enumerations.EAAQualificationQEAA == c.eaa.CategoryQualification() ||
		enumerations.EAAQualificationPubEAA == c.eaa.CategoryQualification()) &&
		!utils.IsTrue(c.eaa.ShortLived()) {
		return c.eaa.EAAPayload().EAAStatus() != nil
	}
	return true
}

func (c *ETSI194721ConformanceCheck) checkSDJWTStatusConformance() bool {
	// TODO: lax processing until TS 119 472-1 review (ported as-is from Java,
	// which keeps this branch entirely commented out).
	return true
}

// BuildAdditionalInfo builds an additional information. Port of the
// overridden buildAdditionalInfo().
func (c *ETSI194721ConformanceCheck) BuildAdditionalInfo() *string {
	errors := make([]string, 0)
	if !c.checkVCTPresent() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagSDJWTEAAVCTPresent,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.EAANotBefore())))
	}
	if !c.checkVCTIntegrityPresent() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagSDJWTEAAVCTIntPresent,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.EAANotBefore())))
	}
	if !c.checkNowAfterNotBefore() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAANowBeforeNBF,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.EAANotBefore())))
	}
	if !c.checkNowBeforeExpiration() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAANowAfterExp,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.EAAExpiration())))
	}
	if !c.checkNowAfterAdministrativeDateIssuance() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAANowBeforeADI,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.AdministrativeIssuanceDate())))
	}
	if !c.checkNowBeforeAdministrativeDateExpiration() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAANowAfterADE,
			process.GetFormattedDate(&c.validationTime), process.GetFormattedDate(c.eaa.AdministrativeExpirationDate())))
	}
	if !c.checkMDOCDocumentNumberPresent() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAAMDocDocumentNumberAbsent))
	}
	if !c.checkMDOCIssuingAuthorityPresent() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAAMDocIssuingAuthority))
	}
	if !c.checkSDJWTIssuingAuthorityAndCountryPresent() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAASDJWTIssuingAuthority))
	}
	if !c.checkSDJWTAdministrativeDateConformance() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAAADSDJWTConformance))
	}
	if !c.checkNoStatusIfShortLived() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAAShortLivedStatusPresent))
	}
	if !c.checkStatusIsPresentIfMandatory() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAAMandatoryStatusAbsent))
	}
	if !c.checkSDJWTStatusConformance() {
		errors = append(errors, c.I18nProvider.GetMessage(i18n.MessageTagEAARevSDJWTConformance))
	}

	if utils.IsCollectionNotEmpty(errors) {
		message := utils.JoinStrings(errors, " - ")
		return &message
	}
	return nil
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *ETSI194721ConformanceCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAETSI194721
}

// ErrorMessageTag returns the check's error message tag. Port of
// getErrorMessageTag().
func (c *ETSI194721ConformanceCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagEAAETSI194721ANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *ETSI194721ConformanceCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *ETSI194721ConformanceCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationEAAConstraintsFailure
}
