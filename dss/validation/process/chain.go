// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/Chain.java (DSS 6.5.RC1).
//
// Java's Chain<T extends XmlConstraintsConclusion> is an abstract class whose
// type parameter is bounded by the generated JAXB base class
// XmlConstraintsConclusion. The generated Go model has no class hierarchy: every
// schema type that extends ConstraintsConclusion embeds
// jaxb.XmlConstraintsConclusionContent (Constraint/Conclusion) and
// jaxb.XmlConstraintsConclusionAttrs (Title) instead of inheriting them, and Go
// generics cannot reach an embedded field through a type parameter. The Result
// type below is that missing bound: it pairs the concrete result object (Java's
// "T result") with pointers to the two embedded base structs, so Chain and
// ChainItem can read and write the inherited Constraint/Conclusion/Title members
// exactly where Java does. The same stand-in is used by
// detailedreport.DetailedReport#HighestConclusion, which types its return to
// *jaxb.XmlConstraintsConclusionContent for the same reason.
//
// The abstract/overridable protected methods are declared in ChainOverrides; a
// concrete chain registers itself with InitChain so that the base's self-calls
// dispatch to the subclass, the way model.TokenBase dispatches through
// model.TokenOverrides.
package process

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
)

// Result binds a generated detailed-report result object (Java's
// "T extends XmlConstraintsConclusion") to the two JAXB base structs it embeds.
// It is the Go stand-in for the XmlConstraintsConclusion supertype: the members
// Chain and ChainItem inherit from it are reached through the pointers captured
// here.
//
// A concrete chain builds it once, next to the result object itself:
//
//	xmlISC := &jaxb.XmlISC{}
//	result := process.NewResult(xmlISC,
//	    &xmlISC.XmlConstraintsConclusionContent, &xmlISC.XmlConstraintsConclusionAttrs)
type Result[T any] struct {
	// Value is the concrete result object, Java's "protected final T result".
	Value T

	// content is the embedded XmlConstraintsConclusionContent of Value,
	// carrying the inherited Constraint and Conclusion members.
	content *jaxb.XmlConstraintsConclusionContent

	// attrs is the embedded XmlConstraintsConclusionAttrs of Value, carrying
	// the inherited Title member.
	attrs *jaxb.XmlConstraintsConclusionAttrs
}

// NewResult binds a result object to its embedded JAXB base structs.
func NewResult[T any](value T, content *jaxb.XmlConstraintsConclusionContent,
	attrs *jaxb.XmlConstraintsConclusionAttrs) *Result[T] {
	return &Result[T]{Value: value, content: content, attrs: attrs}
}

// Conclusion returns the result's conclusion. Port of getConclusion().
func (r *Result[T]) Conclusion() *jaxb.XmlConclusion {
	return r.content.Conclusion
}

// SetConclusion sets the result's conclusion. Port of setConclusion(XmlConclusion).
func (r *Result[T]) SetConclusion(conclusion *jaxb.XmlConclusion) {
	r.content.Conclusion = conclusion
}

// Constraint returns the result's constraint list. Port of getConstraint().
func (r *Result[T]) Constraint() []*jaxb.XmlConstraint {
	return r.content.Constraint
}

// AddConstraint appends to the result's constraint list. Port of
// getConstraint().add(XmlConstraint).
func (r *Result[T]) AddConstraint(constraint *jaxb.XmlConstraint) {
	r.content.Constraint = append(r.content.Constraint, constraint)
}

// SetTitle sets the result's title. Port of setTitle(String). Java's null title
// maps to the empty string: the generated Title member is a plain (non-pointer)
// Go string, so both render identically.
func (r *Result[T]) SetTitle(title string) {
	r.attrs.Title = title
}

// Chain is the polymorphic half of the Java abstract class Chain<T>: the part
// callers hold. The state and the concrete method bodies live in ChainBase.
type Chain[T any] interface {
	// Execute allows initialization and execution of the complete chain until
	// the first failure. Port of execute().
	Execute() T
}

// ChainOverrides declares the abstract and overridable protected methods of
// Chain that the base implementation calls back into. A concrete chain
// registers itself through InitChain; every method the subclass does not define
// is supplied by the embedded ChainBase through ordinary Go method promotion.
type ChainOverrides interface {
	// BuildChainTitle builds the chain title. Port of buildChainTitle().
	BuildChainTitle() string
	// Title returns the title of a Chain (i.e. BasicBuildingBlock title).
	// Port of getTitle().
	Title() i18n.MessageTag
	// AddAdditionalInfo adds additional info to the chain. Port of
	// addAdditionalInfo().
	AddAdditionalInfo()
	// InitChain initializes the chain. Port of the abstract initChain().
	InitChain()
	// IsValidConclusion checks if the conclusion is valid. Port of
	// isValidConclusion(XmlConclusion).
	IsValidConclusion(conclusion *jaxb.XmlConclusion) bool
	// CollectMessages collects required messages from the given constraint to
	// the given conclusion. Port of
	// collectMessages(XmlConclusion, XmlConstraint).
	CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint)
	// CollectAdditionalMessages fills additional messages into the conclusion.
	// Port of collectAdditionalMessages(XmlConclusion).
	CollectAdditionalMessages(conclusion *jaxb.XmlConclusion)
}

// ChainBase is part of the design pattern "Chain of responsibility".
//
// All sub-classes need to implement the method InitChain() which will define
// the ChainItem (constraints) to execute.
//
// The chain is built as follows with the method ChainItem.SetNextItem.
type ChainBase[T any] struct {
	// Result is the result object: Java's "protected final T result", bound to
	// the JAXB base structs it embeds.
	Result *Result[T]

	// I18nProvider is the internationalization provider.
	I18nProvider *i18n.I18nProvider

	// FirstItem is the first item to execute the chain.
	FirstItem ChainItem[T]

	// overrides points back at the concrete chain; see InitChain.
	overrides ChainOverrides
}

// NewChainBase is the common constructor. Port of
// Chain(I18nProvider, T newInstance): the Java constructor takes the new
// instance of the result object, this one takes it bound to its JAXB base
// structs (see Result).
func NewChainBase[T any](i18nProvider *i18n.I18nProvider, newInstance *Result[T]) *ChainBase[T] {
	return &ChainBase[T]{
		I18nProvider: i18nProvider,
		Result:       newInstance,
	}
}

// InitChainBase registers the concrete chain with its base so that the base can
// dispatch to the overridden methods. It must be called exactly once, by the
// concrete chain's constructor, before Execute.
//
// The registration method is named InitChainBase, not InitChain after the
// Init<TypeName> convention the other override-registration pairs in this port
// follow, because Chain already has an abstract initChain() a concrete chain
// must define - a same-named method on the concrete type would shadow the
// promoted registrar.
func (c *ChainBase[T]) InitChainBase(overrides ChainOverrides) {
	c.overrides = overrides
}

// chainOverrides returns the registered overrides, panicking when the concrete
// chain forgot to call InitChain.
func (c *ChainBase[T]) chainOverrides() ChainOverrides {
	if c.overrides == nil {
		panic("Chain was not initialised: the concrete chain must call InitChain in its constructor")
	}
	return c.overrides
}

// Execute allows initialization and execution of the complete chain until the
// first failure. It returns the complete result with constraints and final
// conclusion for the chain. Port of execute().
func (c *ChainBase[T]) Execute() T {
	overrides := c.chainOverrides()
	overrides.InitChain()

	if c.FirstItem != nil {
		c.FirstItem.Execute()
	}

	c.Result.SetTitle(overrides.BuildChainTitle())

	if c.Result.Conclusion() == nil {
		conclusion := &jaxb.XmlConclusion{}
		conclusion.Indication = jaxb.IndicationValue(enumerations.IndicationPassed)
		c.Result.SetConclusion(conclusion)
	}

	c.collectMessages()
	overrides.AddAdditionalInfo()

	return c.Result.Value
}

// BuildChainTitle builds the chain title. Port of buildChainTitle(): Java's null
// (returned when no title MessageTag is defined) maps to the empty string, which
// the generated non-pointer Title member renders identically.
func (c *ChainBase[T]) BuildChainTitle() string {
	if title := BuildStringMessage(c.I18nProvider, c.chainOverrides().Title()); title != nil {
		return *title
	}
	return ""
}

// Title returns the title of a Chain (i.e. BasicBuildingBlock title). Port of
// getTitle(), whose default is null - the empty MessageTag here.
func (c *ChainBase[T]) Title() i18n.MessageTag {
	return ""
}

// AddAdditionalInfo adds additional info to the chain. Port of
// addAdditionalInfo(), which is empty by default.
func (c *ChainBase[T]) AddAdditionalInfo() {
	// default is empty
}

// IsValid checks if the given constraintConclusion has a successful validation
// result. Port of isValid(XmlConstraintsConclusion): the Java parameter is the
// generated base class, whose Go stand-in is the embedded content struct every
// extension type carries.
func (c *ChainBase[T]) IsValid(constraintConclusion *jaxb.XmlConstraintsConclusionContent) bool {
	return constraintConclusion != nil && c.chainOverrides().IsValidConclusion(constraintConclusion.Conclusion)
}

// IsValidConclusion checks if the conclusion is valid, i.e. has a PASSED
// Indication. Port of isValidConclusion(XmlConclusion).
func (c *ChainBase[T]) IsValidConclusion(conclusion *jaxb.XmlConclusion) bool {
	return conclusion != nil && enumerations.IndicationPassed == conclusion.Indication.Indication()
}

// FailLevelRule returns the FAIL level constraint. Port of getFailLevelRule().
func (c *ChainBase[T]) FailLevelRule() policy.LevelRule {
	return GetLevelRule(enumerations.LevelFail)
}

// WarnLevelRule returns the WARN level constraint. Port of getWarnLevelRule().
func (c *ChainBase[T]) WarnLevelRule() policy.LevelRule {
	return GetLevelRule(enumerations.LevelWarn)
}

// InfoLevelRule returns the INFORM level constraint. Port of getInfoLevelRule().
func (c *ChainBase[T]) InfoLevelRule() policy.LevelRule {
	return GetLevelRule(enumerations.LevelInform)
}

// collectMessages collects all required messages. Port of the private
// collectMessages().
func (c *ChainBase[T]) collectMessages() {
	overrides := c.chainOverrides()
	conclusion := c.Result.Conclusion()

	constraints := c.Result.Constraint()
	for _, constraint := range constraints {
		overrides.CollectMessages(conclusion, constraint)
	}
	overrides.CollectAdditionalMessages(conclusion)
}

// CollectMessages collects required messages from the given xmlConstraint to
// the given conclusion. Port of collectMessages(XmlConclusion, XmlConstraint).
//
// NOTE: by default the only one error is already collected in the chain (no more
// possible), therefore no need to collect it again.
func (c *ChainBase[T]) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	warning := constraint.Warning
	if warning != nil {
		conclusion.Warnings = append(conclusion.Warnings, warning)
	}
	info := constraint.Info
	if info != nil {
		conclusion.Infos = append(conclusion.Infos, info)
	}
}

// CollectAllMessages fills all messages from conclusionToFillFrom into
// conclusionToFill. Port of collectAllMessages(XmlConclusion, XmlConclusion).
func (c *ChainBase[T]) CollectAllMessages(conclusionToFill *jaxb.XmlConclusion, conclusionToFillFrom *jaxb.XmlConclusion) {
	errors := conclusionToFillFrom.Errors
	if errors != nil {
		conclusionToFill.Errors = append(conclusionToFill.Errors, errors...)
	}
	warnings := conclusionToFillFrom.Warnings
	if warnings != nil {
		conclusionToFill.Warnings = append(conclusionToFill.Warnings, warnings...)
	}
	infos := conclusionToFillFrom.Infos
	if infos != nil {
		conclusionToFill.Infos = append(conclusionToFill.Infos, infos...)
	}
}

// CollectAdditionalMessages allows to fill up additional messages into the
// conclusion. Port of collectAdditionalMessages(XmlConclusion), empty by default.
func (c *ChainBase[T]) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	// empty by default
}
