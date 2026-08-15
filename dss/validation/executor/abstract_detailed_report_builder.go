// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/AbstractDetailedReportBuilder.java
// (DSS 6.5.RC1).
//
// Two shape differences against Java, both forced by the target language:
//
//   - Java's process(Collection<? extends AbstractTokenProxy>, ...) is an
//     inherited protected method taking a wildcard-typed collection. Go methods
//     cannot declare type parameters, so the port is the package-level generic
//     function Process, whose first argument is the builder.
//
//   - Java's callers hand process() a LinkedHashMap, and DetailedReportBuilder
//     later marshals bbbs.values() straight into the report - so the map's
//     insertion order is byte-visible. A Go map has no order, so the builder
//     records the first-insertion order of the ids in BBBOrder and the report
//     is filled from that (see BasicBuildingBlocksInOrder).

package executor

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	diagnosticjaxb "github.com/utain/esig/dss/diagnostic/jaxb"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process/blocks"
	"github.com/utain/esig/dss/validation/process/qualification"
)

// AbstractDetailedReportBuilder is the abstract code for a DetailedReport
// builder. Port of AbstractDetailedReportBuilder; the Java protected final
// fields are exported so the concrete builders embedding it read them the
// way the Java subclasses do.
type AbstractDetailedReportBuilder struct {
	// I18nProvider is the i18n provider. Port of the protected final
	// i18nProvider field.
	I18nProvider *i18n.I18nProvider

	// DiagnosticData is the DiagnosticData to use. Port of the protected
	// final diagnosticData field.
	DiagnosticData *diagnostic.DiagnosticData

	// Policy is the validation policy. Port of the protected final policy
	// field.
	Policy policy.ValidationPolicy

	// CurrentTime is the validation time. Port of the protected final
	// currentTime field.
	CurrentTime time.Time

	// BBBOrder records the ids of the basic building blocks in the order
	// Process first inserted them, standing in for the insertion order of
	// Java's LinkedHashMap. See the file header.
	BBBOrder []string
}

// NewAbstractDetailedReportBuilder is the default constructor. Port of the
// protected AbstractDetailedReportBuilder(I18nProvider, Date,
// ValidationPolicy, DiagnosticData) constructor.
func NewAbstractDetailedReportBuilder(i18nProvider *i18n.I18nProvider, currentTime time.Time,
	validationPolicy policy.ValidationPolicy,
	diagnosticData *diagnostic.DiagnosticData) AbstractDetailedReportBuilder {
	return AbstractDetailedReportBuilder{
		I18nProvider:   i18nProvider,
		CurrentTime:    currentTime,
		Policy:         validationPolicy,
		DiagnosticData: diagnosticData,
	}
}

// Init initializes the XmlDetailedReport by adding the TL analysis. Port
// of the protected init().
func (b *AbstractDetailedReportBuilder) Init() *jaxb.XmlDetailedReport {
	detailedReport := &jaxb.XmlDetailedReport{}
	detailedReport.ValidationTime = jaxb.NewXSDateTime(b.CurrentTime)

	if b.Policy.EIDASConstraintPresent() {
		detailedReport.TLAnalysis = append(detailedReport.TLAnalysis,
			b.ExecuteAllTLAnalysis(b.DiagnosticData, b.Policy, b.CurrentTime)...)
		detailedReport.LoTEAnalysis = append(detailedReport.LoTEAnalysis,
			b.ExecuteAllLoTEAnalysis(b.DiagnosticData, b.Policy, b.CurrentTime)...)
	}

	return detailedReport
}

// ExecuteAllTLAnalysis executes the TL analysis. Port of the protected
// executeAllTLAnalysis(DiagnosticData, ValidationPolicy, Date).
func (b *AbstractDetailedReportBuilder) ExecuteAllTLAnalysis(diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time) []*jaxb.XmlTLAnalysis {
	var result []*jaxb.XmlTLAnalysis
	result = append(result, b.validateTL(validationPolicy, currentTime, diagnosticData.ListOfTrustedLists())...)
	result = append(result, b.validateTL(validationPolicy, currentTime, diagnosticData.TrustedLists())...)
	return result
}

// validateTL is the port of the private validateTL(ValidationPolicy, Date,
// List<XmlTrustedList>).
func (b *AbstractDetailedReportBuilder) validateTL(validationPolicy policy.ValidationPolicy, currentTime time.Time,
	trustedLists []*diagnosticjaxb.XmlTrustedList) []*jaxb.XmlTLAnalysis {
	var result []*jaxb.XmlTLAnalysis
	if len(trustedLists) > 0 {
		for _, xmlTrustedList := range trustedLists {
			tlValidation := qualification.NewTLValidationBlock(b.I18nProvider, xmlTrustedList, currentTime, validationPolicy)
			result = append(result, tlValidation.Execute())
		}
	}
	return result
}

// ExecuteAllLoTEAnalysis executes the LoTE analysis. Port of the protected
// executeAllLoTEAnalysis(DiagnosticData, ValidationPolicy, Date).
func (b *AbstractDetailedReportBuilder) ExecuteAllLoTEAnalysis(diagnosticData *diagnostic.DiagnosticData,
	validationPolicy policy.ValidationPolicy, currentTime time.Time) []*jaxb.XmlLoTEAnalysisEntry {
	var result []*jaxb.XmlLoTEAnalysisEntry
	result = append(result, b.validateLoTE(validationPolicy, currentTime, diagnosticData.ListsOfListsOfTrustedEntities())...)
	result = append(result, b.validateLoTE(validationPolicy, currentTime, diagnosticData.ListsOfTrustedEntities())...)
	return result
}

// validateLoTE is the port of the private
// validateLoTE(ValidationPolicy, Date, List<XmlListOfTrustedEntities>).
func (b *AbstractDetailedReportBuilder) validateLoTE(validationPolicy policy.ValidationPolicy, currentTime time.Time,
	trustedLists []*diagnosticjaxb.XmlListOfTrustedEntities) []*jaxb.XmlLoTEAnalysisEntry {
	var result []*jaxb.XmlLoTEAnalysisEntry
	if len(trustedLists) > 0 {
		for _, xmlTrustedList := range trustedLists {
			// Java passes the XmlListOfTrustedEntities straight through: it
			// extends XmlTrustSourceList. The generated Go model has no
			// inheritance, so the supertype view is rebuilt from the two base
			// structs both types embed. LoTEValidationBlock only reads it.
			trustSourceList := &diagnosticjaxb.XmlTrustSourceList{
				XmlTrustSourceListContent: xmlTrustedList.XmlTrustSourceListContent,
				XmlTrustSourceListAttrs:   xmlTrustedList.XmlTrustSourceListAttrs,
			}
			loteValidation := qualification.NewLoTEValidationBlock(b.I18nProvider, trustSourceList, currentTime, validationPolicy)
			// Java's LoTEValidationBlock is a Chain<XmlTLAnalysis>, so the
			// object added to the List<XmlLoTEAnalysis> field carries the
			// runtime type XmlTLAnalysis and JAXB emits xsi:type="TLAnalysis"
			// (see testdata/oracle/dr-eaa-pid.xml). AsTLAnalysis records that.
			tlAnalysis := loteValidation.Execute()
			result = append(result, &jaxb.XmlLoTEAnalysisEntry{
				XmlLoTEAnalysis: &jaxb.XmlLoTEAnalysis{
					XmlConstraintsConclusionContent: tlAnalysis.XmlConstraintsConclusionContent,
					XmlLoTEAnalysisAttrs:            tlAnalysis.XmlLoTEAnalysisAttrs,
				},
				AsTLAnalysis: true,
			})
		}
	}
	return result
}

// BasicBuildingBlocksInOrder returns the values of the basic building
// blocks map in the order Process inserted them, the Go stand-in for Java's
// LinkedHashMap.values(). See the file header.
func (b *AbstractDetailedReportBuilder) BasicBuildingBlocksInOrder(
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks) []*jaxb.XmlBasicBuildingBlocks {
	result := make([]*jaxb.XmlBasicBuildingBlocks, 0, len(b.BBBOrder))
	for _, id := range b.BBBOrder {
		if value, ok := bbbs[id]; ok {
			result = append(result, value)
		}
	}
	return result
}

// Process performs the tokens validation. Port of the protected
// process(Collection<? extends AbstractTokenProxy>, Context,
// Map<String, XmlBasicBuildingBlocks>); see the file header for why it is a
// function rather than a method.
func Process[T diagnostic.TokenProxy](b *AbstractDetailedReportBuilder, tokensToProcess []T,
	context enumerations.Context, bbbs map[string]*jaxb.XmlBasicBuildingBlocks) {
	for _, token := range tokensToProcess {
		bbb := blocks.NewBasicBuildingBlocks(
			b.I18nProvider, b.DiagnosticData, token, b.CurrentTime, bbbs, b.Policy, context)
		result := bbb.Execute()
		if _, exists := bbbs[token.Id()]; !exists {
			b.BBBOrder = append(b.BBBOrder, token.Id())
		}
		bbbs[token.Id()] = result
	}
}
