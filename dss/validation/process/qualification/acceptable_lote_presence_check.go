// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/checks/AcceptableLoTEPresenceCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AcceptableLoTEPresenceCheck checks if an acceptable LoTE has been found.
type AcceptableLoTEPresenceCheck[T any] struct {
	*process.ChainItemBase[T]

	// validLoTEs is the set of acceptable Lists of Trusted Entities, keyed by
	// object identity (Java's Set<XmlTrustSourceList> is a HashSet relying on
	// XmlTrustSourceList's default identity equals()/hashCode()).
	validLoTEs map[*jaxb.XmlTrustSourceList]struct{}
}

// NewAcceptableLoTEPresenceCheck is the default constructor. Port of
// AcceptableLoTEPresenceCheck(I18nProvider, T, Set, LevelRule).
func NewAcceptableLoTEPresenceCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	validLoTEUrls map[*jaxb.XmlTrustSourceList]struct{}, constraint policy.LevelRule) *AcceptableLoTEPresenceCheck[T] {
	c := &AcceptableLoTEPresenceCheck[T]{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		validLoTEs:    validLoTEUrls,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *AcceptableLoTEPresenceCheck[T]) Process() bool {
	return len(c.validLoTEs) > 0
}

// MessageTag returns the check's message tag. Port of getMessageTag().
func (c *AcceptableLoTEPresenceCheck[T]) MessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageValidLoTEPresent
}

// ErrorMessageTag returns the check's error message tag. Port of getErrorMessageTag().
func (c *AcceptableLoTEPresenceCheck[T]) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTagCertUsageValidLoTEPresentANS
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *AcceptableLoTEPresenceCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *AcceptableLoTEPresenceCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
