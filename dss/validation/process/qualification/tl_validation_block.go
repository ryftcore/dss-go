// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/TLValidationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	dssjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLValidationBlock is used to perform validation of a Trusted List.
type TLValidationBlock struct {
	*process.ChainBase[*jaxb.XmlTLAnalysis]

	// currentTL is the Trusted List to be validated.
	currentTL *dssjaxb.XmlTrustedList

	// currentTime is the validation time.
	currentTime time.Time

	// policy is the signature validation policy.
	policy policy.ValidationPolicy
}

// NewTLValidationBlock is the default constructor. Port of
// TLValidationBlock(I18nProvider, XmlTrustedList, Date, ValidationPolicy).
func NewTLValidationBlock(i18nProvider *i18n.I18nProvider, currentTL *dssjaxb.XmlTrustedList, currentTime time.Time,
	validationPolicy policy.ValidationPolicy) *TLValidationBlock {
	xmlTLAnalysis := &jaxb.XmlTLAnalysis{}
	b := &TLValidationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlTLAnalysis,
			&xmlTLAnalysis.XmlConstraintsConclusionContent, &xmlTLAnalysis.XmlConstraintsConclusionAttrs)),
		currentTL:   currentTL,
		currentTime: currentTime,
		policy:      validationPolicy,
	}

	b.Result.Value.CountryCode = currentTL.CountryCode
	if currentTL.Url != nil {
		b.Result.Value.URL = *currentTL.Url
	}
	if currentTL.Id != nil {
		b.Result.Value.Id = string(*currentTL.Id)
	}

	b.InitChainBase(b)
	return b
}

// BuildChainTitle builds the chain title. Port of buildChainTitle(): Java's
// java.text.MessageFormat renders a null getCountryCode() as the literal "null", so a nil
// CountryCode is rendered as the string "null" here rather than being dropped, matching
// the BasicValidationProcessCheck.BuildAdditionalInfo precedent.
func (b *TLValidationBlock) BuildChainTitle() string {
	countryCode := "null"
	if b.currentTL.CountryCode != nil {
		countryCode = *b.currentTL.CountryCode
	}
	tag := i18n.MessageTagTL
	if utils.IsTrue(b.currentTL.LOTL) {
		tag = i18n.MessageTagLOTL
	}
	if title := process.BuildStringMessage(b.I18nProvider, tag, countryCode); title != nil {
		return *title
	}
	return ""
}

// Title returns the title of the building block. Port of getTitle(), whose default is
// null - the empty MessageTag here (Chain never sets one).
func (b *TLValidationBlock) Title() i18n.MessageTag {
	return ""
}

// AddAdditionalInfo adds additional info to the chain. Port of addAdditionalInfo(), which
// is empty by default.
func (b *TLValidationBlock) AddAdditionalInfo() {
	// default is empty
}

// InitChain initializes the chain. Port of initChain().
func (b *TLValidationBlock) InitChain() {

	item := b.tlFreshness()
	b.FirstItem = item

	if !b.isLastTL() {
		item = item.SetNextItem(b.tlNotExpired())
	}

	item = item.SetNextItem(b.tlVersion())

	item = item.SetNextItem(b.tlStructure())

	item = item.SetNextItem(b.tlWellSigned())

	if b.currentTL.Mra != nil && *b.currentTL.Mra {
		item.SetNextItem(b.tlMRAEnacted())
	}
}

// isLastTL ports the private isLastTL().
func (b *TLValidationBlock) isLastTL() bool {
	return b.currentTL.NextUpdate == nil
}

// tlFreshness ports the private tlFreshness().
func (b *TLValidationBlock) tlFreshness() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.TLFreshnessConstraint()
	return NewTLFreshnessCheck(b.I18nProvider, b.Result, &b.currentTL.XmlTrustSourceListContent, b.currentTime, constraint)
}

// tlNotExpired ports the private tlNotExpired().
func (b *TLValidationBlock) tlNotExpired() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.TLNotExpiredConstraint()
	return NewTLNotExpiredCheck(b.I18nProvider, b.Result, &b.currentTL.XmlTrustSourceListContent, b.currentTime, constraint)
}

// tlVersion ports the private tlVersion().
func (b *TLValidationBlock) tlVersion() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.TLVersionConstraint()
	return NewTLVersionCheck(b.I18nProvider, b.Result, &b.currentTL.XmlTrustSourceListContent, b.currentTime, constraint)
}

// tlStructure ports the private tlStructure().
func (b *TLValidationBlock) tlStructure() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.TLStructureConstraint()
	return NewTLStructureCheck(b.I18nProvider, b.Result, &b.currentTL.XmlTrustSourceListContent, constraint)
}

// tlWellSigned ports the private tlWellSigned().
func (b *TLValidationBlock) tlWellSigned() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.TLWellSignedConstraint()
	return NewTLWellSignedCheck(b.I18nProvider, b.Result, &b.currentTL.XmlTrustSourceListContent, constraint)
}

// tlMRAEnacted ports the private tlMRAEnacted().
func (b *TLValidationBlock) tlMRAEnacted() process.ChainItem[*jaxb.XmlTLAnalysis] {
	return NewTLMRACheck(b.I18nProvider, b.Result, b.currentTL, b.InfoLevelRule())
}
