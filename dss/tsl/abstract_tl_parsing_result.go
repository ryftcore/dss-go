// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/AbstractTLParsingResult.java (DSS 6.5.RC1).
//
// CROSS-CHUNK DEPENDENCY (see this batch's porter notes): this file embeds
// eu.europa.esig.dss.validation.job.parsing.AbstractParsingResult, ported by the dss-validation-job
// chunk into Go package dss/validation/job, which had not landed when this file was written. The
// embedded name, its constructor and its two members below follow PORTING.md's "exported Go
// identifiers keep the Java name" rule literally (AbstractParsingResult /
// NewAbstractParsingResult / StructureValidationMessages / SetStructureValidationMessages); if
// that chunk chose different spellings, the fix is a mechanical rename here.
package tsl

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/validation/job"
)

// AbstractTLParsingResult is the abstract base for a Trusted Lists parsing analysis result.
// Concrete results (TLParsingResult, LOTLParsingResult) embed it.
type AbstractTLParsingResult struct {
	*job.AbstractParsingResult

	// tslType is the LOTL/TL TSLType.
	tslType enumerations.TSLType

	// sequenceNumber is the LOTL/TL sequence number. Nil stands in for Java's null Integer.
	sequenceNumber *int

	// version is the LOTL/TL version. Nil stands in for Java's null Integer.
	version *int

	// territory is the country (territory).
	territory string

	// issueDate is the issuance date of the LOTL/TL. The zero time.Time stands in for Java's
	// null Date.
	issueDate time.Time

	// nextUpdateDate is the next update date. The zero time.Time stands in for Java's null
	// Date.
	nextUpdateDate time.Time

	// distributionPoints holds the distribution point urls.
	distributionPoints []string
}

// NewAbstractTLParsingResult instantiates the base state of a Trusted List parsing result with
// the Java default (null) values. Port of the protected default constructor.
func NewAbstractTLParsingResult() AbstractTLParsingResult {
	return AbstractTLParsingResult{AbstractParsingResult: job.NewAbstractParsingResult()}
}

// TSLType gets the TSLType. Port of getTSLType().
func (r *AbstractTLParsingResult) TSLType() enumerations.TSLType {
	return r.tslType
}

// SetTSLType sets the TSLType. Port of setTSLType(TSLType).
func (r *AbstractTLParsingResult) SetTSLType(tslType enumerations.TSLType) {
	r.tslType = tslType
}

// SequenceNumber gets the sequence number. Port of getSequenceNumber().
func (r *AbstractTLParsingResult) SequenceNumber() *int {
	return r.sequenceNumber
}

// SetSequenceNumber sets the sequence number. Port of setSequenceNumber(Integer).
func (r *AbstractTLParsingResult) SetSequenceNumber(sequenceNumber *int) {
	r.sequenceNumber = sequenceNumber
}

// Version gets the version. Port of getVersion().
func (r *AbstractTLParsingResult) Version() *int {
	return r.version
}

// SetVersion sets the version. Port of setVersion(Integer).
func (r *AbstractTLParsingResult) SetVersion(version *int) {
	r.version = version
}

// Territory gets the territory (country). Port of getTerritory().
func (r *AbstractTLParsingResult) Territory() string {
	return r.territory
}

// SetTerritory sets the territory. Port of setTerritory(String).
func (r *AbstractTLParsingResult) SetTerritory(territory string) {
	r.territory = territory
}

// IssueDate gets the issuing date. Port of getIssueDate().
func (r *AbstractTLParsingResult) IssueDate() time.Time {
	return r.issueDate
}

// SetIssueDate sets the issue date. Port of setIssueDate(Date).
func (r *AbstractTLParsingResult) SetIssueDate(issueDate time.Time) {
	r.issueDate = issueDate
}

// NextUpdateDate gets the next update date. Port of getNextUpdateDate().
func (r *AbstractTLParsingResult) NextUpdateDate() time.Time {
	return r.nextUpdateDate
}

// SetNextUpdateDate sets the next update date. Port of setNextUpdateDate(Date).
func (r *AbstractTLParsingResult) SetNextUpdateDate(nextUpdateDate time.Time) {
	r.nextUpdateDate = nextUpdateDate
}

// DistributionPoints gets the distribution points. Port of getDistributionPoints().
func (r *AbstractTLParsingResult) DistributionPoints() []string {
	return r.distributionPoints
}

// SetDistributionPoints sets the list of distribution point urls. Port of
// setDistributionPoints(List).
func (r *AbstractTLParsingResult) SetDistributionPoints(distributionPoints []string) {
	r.distributionPoints = distributionPoints
}
