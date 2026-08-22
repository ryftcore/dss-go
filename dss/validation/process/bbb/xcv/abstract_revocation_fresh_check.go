// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/rfc/checks/AbstractRevocationFreshCheck.java (DSS 6.5.RC1).
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractRevocationFreshCheck is the abstract revocation check class. A
// concrete check embeds it instead of process.ChainItemBase and registers
// itself with InitChainItem the same way.
type AbstractRevocationFreshCheck struct {
	*process.ChainItemBase[*jaxb.XmlRFC]

	// RevocationData is the revocation data to check.
	RevocationData *diagnostic.RevocationWrapper

	// validationDate is the validation time.
	validationDate time.Time
}

// NewAbstractRevocationFreshCheck is the default constructor. Port of
// AbstractRevocationFreshCheck(Provider, XmlRFC, RevocationWrapper, Date, LevelRule).
func NewAbstractRevocationFreshCheck(i18nProvider *i18n.Provider, result *process.Result[*jaxb.XmlRFC],
	revocationData *diagnostic.RevocationWrapper, validationDate time.Time,
	constraint policy.LevelRule) *AbstractRevocationFreshCheck {
	return &AbstractRevocationFreshCheck{
		ChainItemBase:  process.NewChainItemBase(i18nProvider, result, constraint),
		RevocationData: revocationData,
		validationDate: validationDate,
	}
}

// IsThisUpdateTimeAfterValidationTime returns if the revocation production
// data is after validation time with the allowed freshness. Port of
// isThisUpdateTimeAfterValidationTime(), taking getMaxFreshness() as a
// parameter since Go has no protected-abstract-method dispatch back into the
// base without the overrides-registration machinery ChainItem itself uses.
func (c *AbstractRevocationFreshCheck) IsThisUpdateTimeAfterValidationTime(maxFreshness int64) bool {
	limit := c.validationDate.Add(-time.Duration(maxFreshness) * time.Millisecond)

	if c.RevocationData == nil {
		return false
	}
	thisUpdate := c.RevocationData.ThisUpdate()
	return thisUpdate != nil && thisUpdate.After(limit)
}

// BuildAdditionalInfo builds an additional information. Port of buildAdditionalInfo().
func (c *AbstractRevocationFreshCheck) BuildAdditionalInfo() *string {
	thisUpdateString := "not defined"
	nextUpdateString := "not defined"
	if c.RevocationData != nil {
		if c.RevocationData.ThisUpdate() != nil {
			thisUpdateString = process.GetFormattedDate(c.RevocationData.ThisUpdate())
		}
		if c.RevocationData.NextUpdate() != nil {
			nextUpdateString = process.GetFormattedDate(c.RevocationData.NextUpdate())
		}
	}
	validationDate := c.validationDate
	message := c.I18nProvider.GetMessage(i18n.MessageTagRevocationCheck,
		process.GetFormattedDate(&validationDate), thisUpdateString, nextUpdateString)
	return &message
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AbstractRevocationFreshCheck) MessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBRFCIRIF
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AbstractRevocationFreshCheck) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagBBBRFCIRIFANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AbstractRevocationFreshCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationIndeterminate
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion().
func (c *AbstractRevocationFreshCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndicationTryLater
}
