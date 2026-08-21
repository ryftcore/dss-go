// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/executor/AbstractDetailedReportBuilder.java
// (DSS 6.5.RC1).
//
// Java's callers hand process() a LinkedHashMap, and DetailedReportBuilder
// later marshals bbbs.values() straight into the report - so the map's
// insertion order is byte-visible. A Go map has no order, so the builder
// records the first-insertion order of the ids in BBBOrder and the report
// is filled from that (see BasicBuildingBlocksInOrder).

package executor

import (
	"sort"
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
	// process first inserted them, standing in for the insertion order of
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
// blocks map in the order process inserted them, the Go stand-in for Java's
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

// IdentifiedToken is the part of diagnostic.TokenProxy javaHashSetOrder needs:
// Java's AbstractTokenProxy derives both equals() and hashCode() from getId()
// alone.
type IdentifiedToken interface {
	// Id returns the unique identifier of the object. Port of getId().
	Id() string
}

// javaHashSetOrder reorders tokens into the iteration order
// java.util.HashSet yields for the same elements.
//
// Several of the collections the report builders feed to process are
// Set<...TokenProxy> upstream (DiagnosticData.getAllRevocationData(),
// getAllSignatures(), getAllCounterSignatures(), getAllKeyBindingSignatures(),
// getAllEAA(), getAllEAARevocationTokens()), while the ported diagnostic
// wrappers return slices in document order. That difference is BYTE-VISIBLE:
// process fills a LinkedHashMap whose values() the detailed report marshals in
// insertion order, so the <BasicBuildingBlocks> elements come out in the
// HashSet's order upstream.
//
// The order is fully determined, not JVM-dependent: AbstractTokenProxy
// overrides hashCode() as `31 + getId().hashCode()`, so a HashSet of these
// wrappers iterates its table in a pure function of the token ids -
// HashMap.hash(key) = h ^ (h >>> 16), bucket index (n-1) & hash, buckets walked
// in ascending index and, inside a bucket, in insertion order (the split a
// resize performs preserves relative order, and treeification needs 8 entries
// in one bucket, which these corpora never reach). n is the table size the map
// has grown to for the element count, with HashMap's default capacity 16 and
// load factor 0.75.
//
// This is the one place in the port that reproduces a Java hash order rather
// than substituting insertion order (the convention documented in
// utils/ordered_map.go): here the order reaches the marshalled report.
func javaHashSetOrder[T IdentifiedToken](tokens []T) []T {
	if len(tokens) < 2 {
		return tokens
	}
	// HashMap's table size: capacity 16, doubled while size exceeds
	// 0.75 * capacity.
	n := 16
	for len(tokens) > n*3/4 {
		n *= 2
	}
	indexed := make([]struct {
		token  T
		bucket int
		order  int
	}, len(tokens))
	for i, token := range tokens {
		h := int32(31) + javaStringHashCode(token.Id())
		spread := h ^ int32(uint32(h)>>16)
		indexed[i].token = token
		indexed[i].bucket = int(uint32(spread) & uint32(n-1))
		indexed[i].order = i
	}
	sort.SliceStable(indexed, func(i, j int) bool { return indexed[i].bucket < indexed[j].bucket })
	result := make([]T, len(tokens))
	for i := range indexed {
		result[i] = indexed[i].token
	}
	return result
}

// JavaHashSetStringOrder reorders strings into the iteration order
// java.util.HashSet<String> yields for the same elements - the string-keyed
// sibling of javaHashSetOrder, whose doc comment explains the table walk. It
// is needed where a report builder marshals the members of a HashSet<String>
// in iteration order: SimpleReportBuilder/SimpleReportForCertificateBuilder's
// getUniqueServiceNames(), whose result becomes the <trustServiceName>
// sequence of a <trustAnchor>.
func JavaHashSetStringOrder(values []string) []string {
	if len(values) < 2 {
		return values
	}
	n := 16
	for len(values) > n*3/4 {
		n *= 2
	}
	indexed := make([]struct {
		value  string
		bucket int
	}, len(values))
	for i, value := range values {
		h := javaStringHashCode(value)
		spread := h ^ int32(uint32(h)>>16)
		indexed[i].value = value
		indexed[i].bucket = int(uint32(spread) & uint32(n-1))
	}
	sort.SliceStable(indexed, func(i, j int) bool { return indexed[i].bucket < indexed[j].bucket })
	result := make([]string, len(values))
	for i := range indexed {
		result[i] = indexed[i].value
	}
	return result
}

// javaStringHashCode reproduces java.lang.String#hashCode: the ids this port
// hashes are ASCII, so each byte is one UTF-16 code unit.
func javaStringHashCode(s string) int32 {
	var h int32
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}

// process performs the tokens validation. Port of the protected
// process(Collection<? extends AbstractTokenProxy>, Context,
// Map<String, XmlBasicBuildingBlocks>).
func (b *AbstractDetailedReportBuilder) process[T diagnostic.TokenProxy](tokensToProcess []T,
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
