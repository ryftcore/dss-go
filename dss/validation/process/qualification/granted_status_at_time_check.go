// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/timestamp/checks/GrantedStatusAtTimeCheck.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// GrantedStatusAtTimeCheck verifies whether the certificate has related
// TrustServices which have been 'granted' at the given validation time.
type GrantedStatusAtTimeCheck[T any] struct {
	*process.ChainItemBase[T]

	// trustServicesAtTime is the list of granted TrustServices at
	// timestamp's production time.
	trustServicesAtTime []*diagnostic.TrustServiceWrapper

	// validationTime is the validation time of the TSP.
	validationTime enumerations.ValidationTime
}

// NewGrantedStatusAtTimeCheck is the default constructor. Port of
// GrantedStatusAtTimeCheck(I18nProvider, T, List, ValidationTime, LevelRule).
func NewGrantedStatusAtTimeCheck[T any](i18nProvider *i18n.I18nProvider, result *process.Result[T],
	trustServicesAtTime []*diagnostic.TrustServiceWrapper, validationTime enumerations.ValidationTime,
	constraint policy.LevelRule) *GrantedStatusAtTimeCheck[T] {
	c := &GrantedStatusAtTimeCheck[T]{
		ChainItemBase:       process.NewChainItemBase(i18nProvider, result, constraint),
		trustServicesAtTime: trustServicesAtTime,
		validationTime:      validationTime,
	}
	c.InitChainItem(c)
	return c
}

// Process performs the check. Port of process().
func (c *GrantedStatusAtTimeCheck[T]) Process() bool {
	return utils.IsCollectionNotEmpty(c.trustServicesAtTime)
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *GrantedStatusAtTimeCheck[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTagQualHasGrantedAt, tag)
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *GrantedStatusAtTimeCheck[T]) BuildErrorMessage() *jaxb.XmlMessage {
	tag, err := process.GetValidationTimeMessageTag(c.validationTime)
	if err != nil {
		panic(err)
	}
	return c.BuildXmlMessage(i18n.MessageTagQualHasGrantedAtANS, tag)
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// getFailedIndicationForConclusion().
func (c *GrantedStatusAtTimeCheck[T]) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.IndicationFailed
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of getFailedSubIndicationForConclusion(), whose default is null.
func (c *GrantedStatusAtTimeCheck[T]) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return ""
}
