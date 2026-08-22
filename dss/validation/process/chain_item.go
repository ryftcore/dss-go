// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/ChainItem.java (DSS 6.5.RC1).
//
// See chain.go for the Result stand-in that replaces Java's
// "T extends XmlConstraintsConclusion" bound, and for the overrides-registration
// pattern this file follows.
//
// Nullability mapping of the overridable getters (Java returns null, Go returns
// the zero value of a string-backed type):
//
//   - getLevel() null           -> Level("")        : check skipped, constraint not defined
//   - getMessageTag() null      -> MessageTag("")   : I18nProvider panics on it, as Java throws
//   - getBlockType() null       -> XmlBlockType("") : the XmlConstraint attribute stays absent
//   - getFailedIndicationForConclusion() / getSuccessIndication() null
//     -> Indication("")  : the generated XmlConclusion Indication member is a
//     plain value, so a null indication and an empty one are indistinguishable
//     in the marshalled output (see notes for SignaturePolicyZeroHashCheck,
//     the only check returning null there, and which is never registered at
//     Level.FAIL, so it never reaches recordConclusion)
//   - getFailedSubIndicationForConclusion() / getSuccessSubIndication() null
//     -> SubIndication("") : the SubIndication member stays absent
//   - buildAdditionalInfo() null -> (*string)(nil)  : the AdditionalInfo element
//     stays absent, which is observable, hence the pointer.
//
// slf4j logging (the "Check skipped", "Unknown level" and undefined-MessageTag
// records) is dropped per PORTING.md.
package process

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
)

// ChainItem is the polymorphic half of the Java abstract class ChainItem<T>: the
// part the chain of responsibility is built and executed through. The state and
// the concrete method bodies live in ChainItemBase.
type ChainItem[T any] interface {
	// SetNextItem allows to build the chain of responsibility, returning the
	// next item. Port of setNextItem(ChainItem).
	SetNextItem(nextItem ChainItem[T]) ChainItem[T]
	// Execute allows to execute the chain of responsibility. Port of execute().
	Execute()
}

// ChainItemOverrides declares the abstract and overridable protected methods of
// ChainItem that the base implementation calls back into. A concrete check
// registers itself through InitChainItem; every method the check does not define
// is supplied by the embedded ChainItemBase through ordinary Go method promotion.
type ChainItemOverrides interface {
	// Process performs the check, returning TRUE if the check succeeds. Port of
	// the abstract process().
	Process() bool
	// Level returns an execution Level of the chain item. Port of getLevel().
	Level() enumerations.Level
	// MessageTag returns an i18n key of a message to get. Port of
	// getMessageTag().
	MessageTag() i18n.MessageTag
	// ErrorMessageTag returns an i18n key of an error message to get. Port of
	// getErrorMessageTag().
	ErrorMessageTag() i18n.MessageTag
	// PreviousErrors returns a list of previous errors occurred in the chain.
	// Port of getPreviousErrors().
	PreviousErrors() []*jaxb.XmlMessage
	// FailedIndicationForConclusion gets an Indication in case of failure. Port
	// of the abstract getFailedIndicationForConclusion().
	FailedIndicationForConclusion() enumerations.Indication
	// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
	// Port of the abstract getFailedSubIndicationForConclusion().
	FailedSubIndicationForConclusion() enumerations.SubIndication
	// BuildErrorMessage builds an error message. Port of buildErrorMessage().
	BuildErrorMessage() *jaxb.XmlMessage
	// BuildConstraintMessage builds a constraint message. Port of
	// buildConstraintMessage().
	BuildConstraintMessage() *jaxb.XmlMessage
	// BlockType returns the validating block type (used for validation result
	// of RAC, RFC, etc.). Port of getBlockType().
	BlockType() jaxb.XmlBlockType
	// BuildAdditionalInfo builds an additional information. Port of
	// buildAdditionalInfo().
	BuildAdditionalInfo() *string
	// AdditionalInfo gets an additional information tag. Port of
	// getAdditionalInfo().
	AdditionalInfo() i18n.MessageTag
	// SuccessIndication gets an Indication if the check succeeds. Port of
	// getSuccessIndication().
	SuccessIndication() enumerations.Indication
	// SuccessSubIndication gets a SubIndication if the check succeeds. Port of
	// getSuccessSubIndication().
	SuccessSubIndication() enumerations.SubIndication
	// ContinueProcessOnFail gets whether the validation process shall be
	// continued on a check failure. Port of continueProcessOnFail().
	ContinueProcessOnFail() bool
	// IsValidConclusion checks if the conclusion has a passed status. Port of
	// isValidConclusion(XmlConclusion).
	IsValidConclusion(conclusion *jaxb.XmlConclusion) bool
}

// ChainItemBase is an item of the Chain class. That follows the design pattern
// "chain of responsibility". Depending on the Level in LevelRule the Chain will
// continue/stop the current treatment. The ChainItem is a validation constraint
// which allows to collect information, warnings, errors,...
type ChainItemBase[T any] struct {
	// I18nProvider is the internationalization provider.
	I18nProvider *i18n.I18nProvider

	// constraint is the Level constraint for the current chain item.
	constraint policy.LevelRule

	// nextItem is the next item to be executed if the following chainItem
	// succeeds.
	nextItem ChainItem[T]

	// result is the conclusion result.
	result *Result[T]

	// bbbId is the executed BasicBuildingBlock id; nil is Java's null.
	bbbId *string

	// overrides points back at the concrete check; see InitChainItem.
	overrides ChainItemOverrides
}

// NewChainItemBase is the common constructor. Port of
// ChainItem(I18nProvider, T, LevelRule).
//
// Panics with the Java message when i18nProvider is nil (Objects.requireNonNull).
func NewChainItemBase[T any](i18nProvider *i18n.I18nProvider, result *Result[T],
	constraint policy.LevelRule) *ChainItemBase[T] {
	return newChainItemBase(i18nProvider, result, constraint, nil)
}

// NewChainItemBaseWithId is the specific constructor for Basic Building Blocks
// validation, taking the XmlBasicBuildingBlocks id. Port of
// ChainItem(I18nProvider, T, LevelRule, String). A Java caller passing a null id
// maps to NewChainItemBase.
func NewChainItemBaseWithId[T any](i18nProvider *i18n.I18nProvider, result *Result[T],
	constraint policy.LevelRule, bbbId string) *ChainItemBase[T] {
	return newChainItemBase(i18nProvider, result, constraint, &bbbId)
}

func newChainItemBase[T any](i18nProvider *i18n.I18nProvider, result *Result[T],
	constraint policy.LevelRule, bbbId *string) *ChainItemBase[T] {
	if i18nProvider == nil {
		panic("i18nProvider must be defined!")
	}
	return &ChainItemBase[T]{
		I18nProvider: i18nProvider,
		result:       result,
		constraint:   constraint,
		bbbId:        bbbId,
	}
}

// InitChainItem registers the concrete check with its base so that the base can
// dispatch to the overridden methods. It must be called exactly once, by the
// concrete check's constructor, before Execute.
func (c *ChainItemBase[T]) InitChainItem(overrides ChainItemOverrides) {
	c.overrides = overrides
}

// chainItemOverrides returns the registered overrides, panicking when the
// concrete check forgot to call InitChainItem.
func (c *ChainItemBase[T]) chainItemOverrides() ChainItemOverrides {
	if c.overrides == nil {
		panic("ChainItem was not initialised: the concrete check must call InitChainItem in its constructor")
	}
	return c.overrides
}

// SetNextItem allows to build the chain of responsibility. It returns the next
// item. Port of setNextItem(ChainItem).
func (c *ChainItemBase[T]) SetNextItem(nextItem ChainItem[T]) ChainItem[T] {
	c.nextItem = nextItem
	return nextItem
}

// Execute allows to execute the chain of responsibility. It will run all the
// chain until the first Level.FAIL and not valid process. Port of execute().
func (c *ChainItemBase[T]) Execute() {
	level := c.chainItemOverrides().Level()
	if level == "" {
		// Check skipped : constraint not defined
		c.callNext()
	} else {
		switch level {
		case enumerations.LevelIgnore:
			c.ignore()
		case enumerations.LevelFail:
			c.fail()
		case enumerations.LevelInform, enumerations.LevelWarn:
			c.informOrWarn(level)
		default:
			// Unknown level
		}
	}
}

// Level returns an execution Level of the chain item. Port of getLevel().
func (c *ChainItemBase[T]) Level() enumerations.Level {
	if c.constraint != nil {
		return c.constraint.Level()
	}
	return ""
}

// MessageTag returns an i18n key of a message to get. Port of getMessageTag(),
// whose default is null.
func (c *ChainItemBase[T]) MessageTag() i18n.MessageTag {
	return ""
}

// ErrorMessageTag returns an i18n key of an error message to get. Port of
// getErrorMessageTag(), whose default is null.
func (c *ChainItemBase[T]) ErrorMessageTag() i18n.MessageTag {
	return ""
}

// PreviousErrors returns a list of previous errors occurred in the chain. Port
// of getPreviousErrors(), whose default is an empty list.
func (c *ChainItemBase[T]) PreviousErrors() []*jaxb.XmlMessage {
	return nil
}

// recordIgnore ports the private recordIgnore().
func (c *ChainItemBase[T]) recordIgnore() {
	c.recordConstraint(jaxb.XmlStatusIgnored)
}

// recordValid ports the private recordValid().
func (c *ChainItemBase[T]) recordValid() {
	c.recordConstraint(jaxb.XmlStatusOK)
}

// recordInvalid ports the private recordInvalid().
func (c *ChainItemBase[T]) recordInvalid() {
	c.recordConstraint(jaxb.XmlStatusNotOK)
}

// recordCustomSuccessConclusion ports the private recordCustomSuccessConclusion().
func (c *ChainItemBase[T]) recordCustomSuccessConclusion() {
	overrides := c.chainItemOverrides()
	conclusion := &jaxb.XmlConclusion{}
	conclusion.Indication = jaxb.IndicationValue(overrides.SuccessIndication())
	setSubIndication(conclusion, overrides.SuccessSubIndication())
	c.result.SetConclusion(conclusion)
}

// recordConclusion ports the private recordConclusion().
func (c *ChainItemBase[T]) recordConclusion() {
	overrides := c.chainItemOverrides()

	var conclusion *jaxb.XmlConclusion
	if c.result.Conclusion() != nil {
		// For uninterrupted chains
		conclusion = c.result.Conclusion()
	} else {
		conclusion = &jaxb.XmlConclusion{}
		conclusion.Indication = jaxb.IndicationValue(overrides.FailedIndicationForConclusion())
		setSubIndication(conclusion, overrides.FailedSubIndicationForConclusion())
	}

	previousErrors := overrides.PreviousErrors()
	if utils.IsCollectionNotEmpty(previousErrors) {
		conclusion.Errors = append(conclusion.Errors, previousErrors...)
	} else {
		conclusion.Errors = append(conclusion.Errors, overrides.BuildErrorMessage())
	}

	c.result.SetConclusion(conclusion)
}

// setSubIndication ports XmlConclusion#setSubIndication: Java's null leaves the
// member unset, which the generated pointer member expresses as nil.
func setSubIndication(conclusion *jaxb.XmlConclusion, subIndication enumerations.SubIndication) {
	if subIndication == "" {
		conclusion.SubIndication = nil
		return
	}
	value := jaxb.SubIndicationValue(subIndication)
	conclusion.SubIndication = &value
}

// recordInfosOrWarns ports the private recordInfosOrWarns(Level).
func (c *ChainItemBase[T]) recordInfosOrWarns(level enumerations.Level) {
	if enumerations.LevelInform == level {
		c.recordConstraint(jaxb.XmlStatusInformation)
	} else if enumerations.LevelWarn == level {
		c.recordConstraint(jaxb.XmlStatusWarning)
	}
}

// recordConstraint ports the private recordConstraint(XmlStatus).
func (c *ChainItemBase[T]) recordConstraint(status jaxb.XmlStatus) {
	overrides := c.chainItemOverrides()

	xmlConstraint := &jaxb.XmlConstraint{}
	xmlConstraint.Name = overrides.BuildConstraintMessage()
	xmlConstraint.Status = status
	if c.bbbId != nil {
		bbbId := *c.bbbId
		xmlConstraint.Id = &bbbId
	}
	if blockType := overrides.BlockType(); blockType != "" {
		xmlConstraint.BlockType = &blockType
	}

	if jaxb.XmlStatusNotOK == status {
		xmlConstraint.Error = overrides.BuildErrorMessage()
	} else if jaxb.XmlStatusWarning == status {
		xmlConstraint.Warning = overrides.BuildErrorMessage()
	} else if jaxb.XmlStatusInformation == status {
		xmlConstraint.Info = overrides.BuildErrorMessage()
	}

	if jaxb.XmlStatusIgnored != status {
		xmlConstraint.AdditionalInfo = overrides.BuildAdditionalInfo()
	}
	c.addConstraint(xmlConstraint)
}

// BuildErrorMessage builds an error message. Port of buildErrorMessage().
func (c *ChainItemBase[T]) BuildErrorMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.chainItemOverrides().ErrorMessageTag())
}

// BuildConstraintMessage builds a constraint message. Port of
// buildConstraintMessage().
func (c *ChainItemBase[T]) BuildConstraintMessage() *jaxb.XmlMessage {
	return c.BuildXmlMessage(c.chainItemOverrides().MessageTag())
}

// BlockType returns the validating block type (used for validation result of
// RAC, RFC, etc.). Port of getBlockType(), whose default is null.
func (c *ChainItemBase[T]) BlockType() jaxb.XmlBlockType {
	// by default
	return ""
}

// BuildAdditionalInfo builds an additional information. Port of
// buildAdditionalInfo().
func (c *ChainItemBase[T]) BuildAdditionalInfo() *string {
	return BuildStringMessage(c.I18nProvider, c.chainItemOverrides().AdditionalInfo())
}

// AdditionalInfo gets an additional information. Port of getAdditionalInfo(),
// whose default is null.
func (c *ChainItemBase[T]) AdditionalInfo() i18n.MessageTag {
	return ""
}

// addConstraint ports the private addConstraint(XmlConstraint).
func (c *ChainItemBase[T]) addConstraint(constraint *jaxb.XmlConstraint) {
	c.result.AddConstraint(constraint)
}

// BuildXmlMessage builds the XmlMessage. Port of
// buildXmlMessage(MessageTag, Object...).
//
// Java guards the formatted message against null and logs an slf4j error
// instead of filling the message in; I18nProvider#getMessage never returns null
// (an undefined key resolves to the tag id itself), so the guard is unreachable
// and only its taken branch is ported.
func (c *ChainItemBase[T]) BuildXmlMessage(messageTag i18n.MessageTag, args ...interface{}) *jaxb.XmlMessage {
	xmlMessage := &jaxb.XmlMessage{}
	message := c.I18nProvider.GetMessage(messageTag, args...)
	key := messageTag.Id()
	xmlMessage.Key = &key
	xmlMessage.Value = message
	return xmlMessage
}

// fail ports the private fail(): this method skips next elements.
func (c *ChainItemBase[T]) fail() {
	overrides := c.chainItemOverrides()
	valid := overrides.Process()
	if valid {
		c.recordValid()
		if !c.isCustomSuccessConclusion() {
			c.callNext()
		} else {
			c.recordCustomSuccessConclusion()
		}
	} else {
		c.recordInvalid()
		c.recordConclusion()
		if overrides.ContinueProcessOnFail() {
			c.callNext()
		}
	}
}

// isCustomSuccessConclusion ports the private isCustomSuccessConclusion().
func (c *ChainItemBase[T]) isCustomSuccessConclusion() bool {
	return c.chainItemOverrides().SuccessIndication() != ""
}

// SuccessIndication gets an Indication if the check succeeds. Port of
// getSuccessIndication(), whose default is null.
func (c *ChainItemBase[T]) SuccessIndication() enumerations.Indication {
	return ""
}

// SuccessSubIndication gets a SubIndication if the check succeeds. Port of
// getSuccessSubIndication(), whose default is null.
func (c *ChainItemBase[T]) SuccessSubIndication() enumerations.SubIndication {
	return ""
}

// ContinueProcessOnFail gets whether the validation process shall be continued
// on a check failure. Default : FALSE (break the validation process in case of a
// check failure). Port of continueProcessOnFail().
func (c *ChainItemBase[T]) ContinueProcessOnFail() bool {
	return false
}

// informOrWarn ports the private informOrWarn(Level).
func (c *ChainItemBase[T]) informOrWarn(level enumerations.Level) {
	valid := c.chainItemOverrides().Process()
	if valid {
		c.recordValid()
	} else {
		c.recordInfosOrWarns(level)
	}
	c.callNext()
}

// ignore ports the private ignore().
func (c *ChainItemBase[T]) ignore() {
	c.recordIgnore()
	c.callNext()
}

// callNext ports the private callNext().
func (c *ChainItemBase[T]) callNext() {
	if c.nextItem != nil {
		c.nextItem.Execute()
	}
}

// IsValid checks if the conclusion is valid. Port of
// isValid(XmlConstraintsConclusion): the Java parameter is the generated base
// class, whose Go stand-in is the embedded content struct every extension type
// carries.
func (c *ChainItemBase[T]) IsValid(constraintConclusion *jaxb.XmlConstraintsConclusionContent) bool {
	return constraintConclusion != nil && c.chainItemOverrides().IsValidConclusion(constraintConclusion.Conclusion)
}

// IsValidConclusion checks if the conclusion has a PASSED indication. Port of
// isValidConclusion(XmlConclusion).
func (c *ChainItemBase[T]) IsValidConclusion(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationPassed == conclusion.Indication.Indication() ||
		enumerations.IndicationTotalPassed == conclusion.Indication.Indication())
}

// IsInvalidConclusion checks if the conclusion has a FAILED indication. Port of
// isInvalidConclusion(XmlConclusion).
func (c *ChainItemBase[T]) IsInvalidConclusion(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && (enumerations.IndicationFailed == conclusion.Indication.Indication() ||
		enumerations.IndicationTotalFailed == conclusion.Indication.Indication())
}

// IsIndeterminateConclusion checks if the conclusion has an INDETERMINATE
// indication. Port of isIndeterminateConclusion(XmlConclusion).
func (c *ChainItemBase[T]) IsIndeterminateConclusion(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && enumerations.IndicationIndeterminate == conclusion.Indication.Indication()
}
