// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/LoTEValidationBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	dssjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
)

// LoTEValidationBlock performs validation of a LoTE.
type LoTEValidationBlock struct {
	*process.ChainBase[*jaxb.XmlTLAnalysis]

	// currentList is the Trusted List to be validated.
	currentList *dssjaxb.XmlTrustSourceList

	// currentTime is the validation time.
	currentTime time.Time

	// policy is the signature validation policy.
	policy policy.ValidationPolicy
}

// NewLoTEValidationBlock is the default constructor. Port of
// LoTEValidationBlock(I18nProvider, XmlTrustSourceList, Date, ValidationPolicy).
func NewLoTEValidationBlock(i18nProvider *i18n.I18nProvider, currentList *dssjaxb.XmlTrustSourceList, currentTime time.Time,
	validationPolicy policy.ValidationPolicy) *LoTEValidationBlock {
	xmlTLAnalysis := &jaxb.XmlTLAnalysis{}
	b := &LoTEValidationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlTLAnalysis,
			&xmlTLAnalysis.XmlConstraintsConclusionContent, &xmlTLAnalysis.XmlConstraintsConclusionAttrs)),
		currentList: currentList,
		currentTime: currentTime,
		policy:      validationPolicy,
	}

	b.Result.Value.CountryCode = currentList.CountryCode
	if currentList.Url != nil {
		b.Result.Value.URL = *currentList.Url
	}
	if currentList.Id != nil {
		b.Result.Value.Id = string(*currentList.Id)
	}

	b.InitChainBase(b)
	return b
}

// BuildChainTitle builds the chain title. Port of buildChainTitle(): Java's
// java.text.MessageFormat renders a null getCountryCode() as the literal "null", so a nil
// CountryCode is rendered as the string "null" here rather than being dropped, matching
// the BasicValidationProcessCheck.BuildAdditionalInfo precedent.
func (b *LoTEValidationBlock) BuildChainTitle() string {
	countryCode := "null"
	if b.currentList.CountryCode != nil {
		countryCode = *b.currentList.CountryCode
	}
	if title := process.BuildStringMessage(b.I18nProvider, i18n.MessageTag_LOTE, countryCode); title != nil {
		return *title
	}
	return ""
}

// Title returns the title of the building block. Port of getTitle(), whose default is
// null - the empty MessageTag here (Chain never sets one).
func (b *LoTEValidationBlock) Title() i18n.MessageTag {
	return ""
}

// AddAdditionalInfo adds additional info to the chain. Port of addAdditionalInfo(), which
// is empty by default.
func (b *LoTEValidationBlock) AddAdditionalInfo() {
	// default is empty
}

// InitChain initializes the chain. Port of initChain().
func (b *LoTEValidationBlock) InitChain() {

	item := b.loteFreshness()
	b.FirstItem = item

	item = item.SetNextItem(b.loteNotExpired())

	item = item.SetNextItem(b.loteVersion())

	item = item.SetNextItem(b.loteStructure())

	item.SetNextItem(b.loteWellSigned())
}

// loteFreshness ports the private loteFreshness().
func (b *LoTEValidationBlock) loteFreshness() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.LoTEFreshnessConstraint()
	return NewTLFreshnessCheck(b.I18nProvider, b.Result, &b.currentList.XmlTrustSourceListContent, b.currentTime, constraint)
}

// loteNotExpired ports the private loteNotExpired().
func (b *LoTEValidationBlock) loteNotExpired() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.LoTENotExpiredConstraint()
	return NewTLNotExpiredCheck(b.I18nProvider, b.Result, &b.currentList.XmlTrustSourceListContent, b.currentTime, constraint)
}

// loteVersion ports the private loteVersion().
func (b *LoTEValidationBlock) loteVersion() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.LoTEVersionConstraint()
	return NewTLVersionCheck(b.I18nProvider, b.Result, &b.currentList.XmlTrustSourceListContent, b.currentTime, constraint)
}

// loteStructure ports the private loteStructure().
func (b *LoTEValidationBlock) loteStructure() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.LoTEStructureConstraint()
	return NewTLStructureCheck(b.I18nProvider, b.Result, &b.currentList.XmlTrustSourceListContent, constraint)
}

// loteWellSigned ports the private loteWellSigned().
func (b *LoTEValidationBlock) loteWellSigned() process.ChainItem[*jaxb.XmlTLAnalysis] {
	constraint := b.policy.LoTEWellSignedConstraint()
	return NewTLWellSignedCheck(b.I18nProvider, b.Result, &b.currentList.XmlTrustSourceListContent, constraint)
}
