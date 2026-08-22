// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/evidencerecord/EvidenceRecordValidationProcess.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note.
//
// Package placement deviation: Java's vpfswatsp.evidencerecord
// is a package of its own, distinct from vpfswatsp. Everything else in the
// vpfswatsp tree folds into one Go package, but TimestampsValidationBlock
// extends vpftsp.TimestampsValidationBlock while vpftsp imports vpfswatsp
// (POEExtraction) - an import cycle Go forbids. The five evidence-record classes
// therefore keep Java's own vpfswatsp/evidencerecord package boundary; nothing
// in vpfswatsp, vpftsp or vpftspwatsp imports them (only the validation
// executor's DetailedReportBuilder does), so the edge only ever points upward.
package evidencerecord

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/aov"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/cv"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp"
)

// ValidationProcess performs Evidence Record validation as per
// clause 5.6.3 "Validation Process for Signatures providing Long Term
// Availability and Integrity of Validation Material", step 1) 5.6.3.4
// "Processing".
type ValidationProcess struct {
	*process.ChainBase[*jaxb.XmlValidationProcessEvidenceRecord]

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.Data

	// evidenceRecord is the evidence record being validated.
	evidenceRecord *diagnostic.EvidenceRecordWrapper

	// xmlTimestamps is the collection of timestamps.
	xmlTimestamps []*jaxb.XmlTimestamp

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// policy is the validation policy used to validate evidence records.
	policy policy.ValidationPolicy

	// currentTime is the validation time.
	currentTime time.Time
}

// NewEvidenceRecordValidationProcess is the common constructor. Port of
// ValidationProcess(Provider, Data, EvidenceRecordWrapper, Collection, Map, ValidationPolicy, Date).
func NewEvidenceRecordValidationProcess(i18nProvider *i18n.Provider, diagnosticData *diagnostic.Data,
	evidenceRecord *diagnostic.EvidenceRecordWrapper, xmlTimestamps []*jaxb.XmlTimestamp,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, validationPolicy policy.ValidationPolicy,
	currentTime time.Time) *ValidationProcess {
	xmlResult := &jaxb.XmlValidationProcessEvidenceRecord{}
	c := &ValidationProcess{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),

		diagnosticData: diagnosticData,
		evidenceRecord: evidenceRecord,
		xmlTimestamps:  xmlTimestamps,
		bbbs:           bbbs,
		policy:         validationPolicy,
		currentTime:    currentTime,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the process. Port of getTitle().
func (c *ValidationProcess) Title() i18n.MessageTag {
	return i18n.MessageTagVPER
}

// InitChain initializes the chain. Port of initChain(): the IllegalStateException
// becomes a panic, initChain being called from the base Chain, which cannot
// propagate an error.
func (c *ValidationProcess) InitChain() {

	var item process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord]

	/*
	 * 5.6.3.4 Processing
	 *
	 * a) The process shall take the first ER that was not yet processed.
	 *
	 * b) The process shall verify this ER according to IETF RFC 4998 [i.9] or IETF RFC 6283 [i.10] taking into
	 * account the following additional requirements when validating a time-stamp token at the time of the
	 * following Archive Timestamp:
	 */
	digestMatchers := c.evidenceRecord.DigestMatchers()

	if utils.IsCollectionNotEmpty(digestMatchers) {

		for _, digestMatcher := range digestMatchers {
			// Evidence Records optionally allow additional digests to be present within first data group
			if enumerations.DigestMatcherTypeEvidenceRecordOrphanReference != digestMatcherTypeOf(digestMatcher) {

				referenceDataFound := c.referenceDataFound(digestMatcher)
				if item == nil {
					c.FirstItem = referenceDataFound
					item = referenceDataFound
				} else {
					item = item.SetNextItem(referenceDataFound)
				}

				if digestMatcher.DataFound {
					item = item.SetNextItem(c.referenceDataIntact(digestMatcher))
				}

			}
		}

		atLeastOneDataObjectFound := c.atLeastOneDataObjectFound(digestMatchers)
		if item == nil {
			c.FirstItem = atLeastOneDataObjectFound
			item = atLeastOneDataObjectFound
		} else {
			item = item.SetNextItem(atLeastOneDataObjectFound)
		}

		item = item.SetNextItem(c.referenceDataGroup(digestMatchers))

	}

	if item == nil {
		panic("Evidence record shall contain at least one DigestMatcher!")
	}

	// Embedded or externally provided evidence records
	if (enumerations.EvidenceRecordOriginSignature == c.evidenceRecord.Origin() ||
		enumerations.EvidenceRecordOriginExternal == c.evidenceRecord.Origin()) &&
		utils.IsCollectionNotEmpty(c.evidenceRecord.CoveredSignatures()) {
		item = item.SetNextItem(c.signedFilesCoveredCheck())
	}

	// ASiC container evidence record
	if c.diagnosticData.IsContainerInfoPresent() && c.coversSignatureOrTimestampOrEvidenceRecord(c.evidenceRecord) {
		item = item.SetNextItem(c.signedAndTimestampedFilesCoveredCheck())
	}

	// Algorithm Obsolescence Validation to be included
	var xmlAOV *jaxb.XmlAOV

	/*
	 * i) Before validating a time-stamp the process shall extract POEs (as per clause 5.6.2.3) of the
	 * time-stamp within the next Archive timestamp and initialize the set of temporary POEs with the
	 * extracted POEs.
	 */
	lowestPOE := toXmlProofOfExistenceAtTime(c.currentTime)
	tstChainValid := true

	timestampsList := c.evidenceRecord.TimestampList()
	if utils.IsCollectionNotEmpty(timestampsList) {

		for _, timestamp := range timestampsList {
			/*
			 * ii) The time stamp validation of the time-stamp token shall be performed, as per clause 5.4.
			 *
			 * iii) The past signature validation process for the signature of the time-stamp token as per clause 5.6.2.4
			 * shall be used with the following inputs: the time-stamp, the TSA's certificate, the X.509 validation
			 * parameters, the X.509 validation constraints, the cryptographic constraints, certificate validation
			 * data, the indication/sub-indication returned in step ii) and the set of POEs available so far, and the
			 * set of temporary POEs.
			 */
			bbbTsp := c.bbbs[timestamp.Id()]
			timestampValidation := c.timestampValidation(timestamp)
			if bbbTsp != nil && timestampValidation != nil {

				// Basic and Past time-stamp validations are performed inside
				item = item.SetNextItem(c.timestampValidationConclusive(timestamp, timestampValidation))

				/*
				 * ETSI TS 119 102-2 (4.3.12.7 Crypto Information Element):
				 *
				 * This element shall be present when the main status indication is INDETERMINATE and
				 * the subindication is CRYPTO_CONSTRAINTS_FAILURE. In all other cases, this element may be present.
				 */
				timestampConclusion := timestampValidation.Conclusion
				timestampSubIndication := subIndicationOf(timestampConclusion)
				if enumerations.IndicationIndeterminate == timestampConclusion.Indication.Indication() &&
					(enumerations.SubIndicationCryptoConstraintsFailure == timestampSubIndication ||
						enumerations.SubIndicationCryptoConstraintsFailureNoPOE == timestampSubIndication) {
					tstAOV := bbbTsp.AOV
					if xmlAOV == nil || (c.IsValid(&xmlAOV.XmlConstraintsConclusionContent) && tstAOV != nil &&
						!c.IsValid(&tstAOV.XmlConstraintsConclusionContent)) {
						xmlAOV = tstAOV
					}
				}

				if !c.IsValid(&timestampValidation.XmlConstraintsConclusionContent) {
					tstChainValid = false
					break
				}
			}
		}

		/*
		 * c) If step b) found the ER to be valid, the process shall add a POE for every object covered by the ER at
		 * signing time value of the initial archive time-stamp.
		 */
		if tstChainValid {
			// when valid, conclusion is not yet set
			lowestPOE = toXmlProofOfExistenceForTimestamp(timestampsList[0])
		}
	}

	c.Result.Value.ProofOfExistence = lowestPOE

	// Validate cryptographic constraints of DigestMatchers
	if utils.IsCollectionNotEmpty(digestMatchers) {
		lowestPOETime := lowestPOE.Time.Time()
		algorithmObsolescenceValidation := aov.NewEvidenceRecordAlgorithmObsolescenceValidation(
			c.I18nProvider, c.evidenceRecord, lowestPOETime, c.policy)
		erAOV := algorithmObsolescenceValidation.Execute()

		item = item.SetNextItem(c.algorithmsObsolescenceValidation(erAOV, lowestPOETime)) //nolint:staticcheck // mirrors upstream EvidenceRecordValidationProcess#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

		if xmlAOV == nil || (c.IsValid(&xmlAOV.XmlConstraintsConclusionContent) && erAOV != nil &&
			!c.IsValid(&erAOV.XmlConstraintsConclusionContent)) {
			xmlAOV = erAOV
		}
	}

	c.Result.Value.AOV = xmlAOV

}

// subIndicationOf reads XmlConclusion#getSubIndication(): the generated member
// is a pointer, whose nil is Java's null.
func subIndicationOf(conclusion *jaxb.XmlConclusion) enumerations.SubIndication {
	if conclusion.SubIndication == nil {
		return ""
	}
	return conclusion.SubIndication.SubIndication()
}

// digestMatcherTypeOf reads XmlDigestMatcher#getType(): the generated member is
// a pointer, whose nil is Java's null.
func digestMatcherTypeOf(digestMatcher *diagnosticjaxb.XmlDigestMatcher) enumerations.DigestMatcherType {
	if digestMatcher.Type == nil {
		return ""
	}
	return digestMatcher.Type.DigestMatcherType()
}

// referenceDataFound ports the private referenceDataFound(XmlDigestMatcher).
func (c *ValidationProcess) referenceDataFound(
	digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordDataObjectExistenceConstraint()
	return cv.NewReferenceDataExistenceCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// referenceDataIntact ports the private referenceDataIntact(XmlDigestMatcher).
func (c *ValidationProcess) referenceDataIntact(
	digestMatcher *diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordDataObjectIntactConstraint()
	return cv.NewReferenceDataIntactCheck(c.I18nProvider, c.Result, digestMatcher, constraint)
}

// atLeastOneDataObjectFound ports the private atLeastOneDataObjectFound(List).
func (c *ValidationProcess) atLeastOneDataObjectFound(
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordDataObjectFoundConstraint()
	return cv.NewAtLeastOneReferenceDataObjectFoundCheck(c.I18nProvider, c.Result, digestMatchers, constraint)
}

// referenceDataGroup ports the private referenceDataGroup(List).
func (c *ValidationProcess) referenceDataGroup(
	digestMatchers []*diagnosticjaxb.XmlDigestMatcher) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordDataObjectGroupConstraint()
	return cv.NewReferenceDataGroupCheck(c.I18nProvider, c.Result, digestMatchers, constraint)
}

// signedFilesCoveredCheck ports the private signedFilesCoveredCheck().
func (c *ValidationProcess) signedFilesCoveredCheck() process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordSignedFilesCoveredConstraint()
	return NewEvidenceRecordSignedFilesCoveredCheck(c.I18nProvider, c.Result, c.evidenceRecord, constraint)
}

// signedAndTimestampedFilesCoveredCheck ports the private
// signedAndTimestampedFilesCoveredCheck().
func (c *ValidationProcess) signedAndTimestampedFilesCoveredCheck() process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	constraint := c.policy.EvidenceRecordContainerSignedAndTimestampedFilesCoveredConstraint()
	return NewEvidenceRecordSignedAndTimestampedFilesCoveredCheck(c.I18nProvider, c.Result, c.diagnosticData,
		c.evidenceRecord, constraint)
}

// timestampValidationConclusive ports the private
// timestampValidationConclusive(TimestampWrapper, XmlValidationProcessArchivalDataTimestamp).
func (c *ValidationProcess) timestampValidationConclusive(timestampWrapper *diagnostic.TimestampWrapper,
	timestampValidationResult *jaxb.XmlValidationProcessArchivalDataTimestamp) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	return vpfswatsp.NewTimestampValidationCheck(c.I18nProvider, c.Result, timestampWrapper, timestampValidationResult,
		c.FailLevelRule())
}

// algorithmsObsolescenceValidation ports the private
// algorithmsObsolescenceValidation(XmlAOV, Date): the unsupported-Context error
// becomes a panic, the caller being initChain, which cannot propagate one.
func (c *ValidationProcess) algorithmsObsolescenceValidation(aovResult *jaxb.XmlAOV,
	lowestPOETime time.Time) process.ChainItem[*jaxb.XmlValidationProcessEvidenceRecord] {
	position, err := process.GetCryptoPosition(enumerations.ContextEvidenceRecord)
	if err != nil {
		panic(err)
	}
	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, lowestPOETime, position,
		c.evidenceRecord.Id())
}

// timestampValidation ports the private getTimestampValidation(TimestampWrapper):
// the generated XmlTimestamp Id member is a *string, whose nil is Java's null,
// and Utils.areStringsEqual(null, x) is false unless both are null.
func (c *ValidationProcess) timestampValidation(
	newestTimestamp *diagnostic.TimestampWrapper) *jaxb.XmlValidationProcessArchivalDataTimestamp {
	for _, xmlTimestamp := range c.xmlTimestamps {
		if xmlTimestamp.Id != nil && utils.AreStringsEqual(*xmlTimestamp.Id, newestTimestamp.Id()) {
			return xmlTimestamp.ValidationProcessArchivalDataTimestamp
		}
	}
	return nil
}

// toXmlProofOfExistenceAtTime ports the private toXmlProofOfExistence(Date).
func toXmlProofOfExistenceAtTime(date time.Time) *jaxb.XmlProofOfExistence {
	xmlPoe := &jaxb.XmlProofOfExistence{}
	xmlPoe.Time = jaxb.XSDateTime(date)
	return xmlPoe
}

// toXmlProofOfExistenceForTimestamp ports the private
// toXmlProofOfExistence(TimestampWrapper). Java dereferences getProductionTime()
// into a Date member that accepts null; the generated Go member is a plain
// XSDateTime, so a missing production time maps to the zero time.
func toXmlProofOfExistenceForTimestamp(timestampWrapper *diagnostic.TimestampWrapper) *jaxb.XmlProofOfExistence {
	var productionTime time.Time
	if t := timestampWrapper.ProductionTime(); t != nil {
		productionTime = *t
	}
	xmlPoe := toXmlProofOfExistenceAtTime(productionTime)
	id := timestampWrapper.Id()
	xmlPoe.TimestampId = &id
	return xmlPoe
}

// coversSignatureOrTimestampOrEvidenceRecord ports the private
// coversSignatureOrTimestampOrEvidenceRecord(EvidenceRecordWrapper).
func (c *ValidationProcess) coversSignatureOrTimestampOrEvidenceRecord(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) bool {
	return utils.IsCollectionNotEmpty(evidenceRecord.CoveredSignatures()) ||
		utils.IsCollectionNotEmpty(evidenceRecord.CoveredTimestamps()) ||
		utils.IsCollectionNotEmpty(evidenceRecord.CoveredEvidenceRecords())
}
