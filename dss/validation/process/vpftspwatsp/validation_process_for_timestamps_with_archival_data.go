// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpftspwatsp/ValidationProcessForTimestampsWithArchivalData.java (DSS 6.5.RC1).
//
// CROSS-CHUNK NOTE (phase 8e): Java's BasicTimestampValidationCheck lives in
// vpftsp.checks; the vpftsp chunk relocated it into package vpfbs (see
// vpfbs/basic_timestamp_validation_check.go) to break the vpfbs<->vpftsp import
// cycle, so it is imported from there.
//
// PACKAGE-BOUNDARY NOTE: TimestampMessageImprintCheck, the fourth class of
// Java's vpftspwatsp.checks, stayed in its own Go package
// vpftspwatsp/checks - see that file's header for the cycle that forces it.
package vpftspwatsp

import (
	"time"

	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb/aov"
	"github.com/utain/esig/dss/validation/process/vpfbs"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
)

// ValidationProcessForTimestampsWithArchivalData validates a timestamp with a
// provided archival data (POE).
type ValidationProcessForTimestampsWithArchivalData struct {
	*process.ChainBase[*jaxb.XmlValidationProcessArchivalDataTimestamp]

	// vpftsp is the timestamp validation with long-term data result.
	vpftsp *jaxb.XmlValidationProcessBasicTimestamp

	// timestamp is the timestamp.
	timestamp *diagnostic.TimestampWrapper

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// evidenceRecordValidations is the map of processed evidence records.
	evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord

	// currentTime is the current time of validation.
	currentTime time.Time

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// poe is the POE container.
	poe *vpfswatsp.POEExtraction
}

// NewValidationProcessForTimestampsWithArchivalData is the default constructor.
// Port of ValidationProcessForTimestampsWithArchivalData(I18nProvider, TimestampWrapper, XmlValidationProcessBasicTimestamp, Map, Map, Date, ValidationPolicy, POEExtraction).
func NewValidationProcessForTimestampsWithArchivalData(i18nProvider *i18n.I18nProvider,
	timestamp *diagnostic.TimestampWrapper, vpftspResult *jaxb.XmlValidationProcessBasicTimestamp,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, evidenceRecordValidations map[string]*jaxb.XmlEvidenceRecord,
	currentTime time.Time, validationPolicy policy.ValidationPolicy,
	poe *vpfswatsp.POEExtraction) *ValidationProcessForTimestampsWithArchivalData {
	xmlResult := &jaxb.XmlValidationProcessArchivalDataTimestamp{}
	c := &ValidationProcessForTimestampsWithArchivalData{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		vpftsp:                    vpftspResult,
		timestamp:                 timestamp,
		bbbs:                      bbbs,
		evidenceRecordValidations: evidenceRecordValidations,
		currentTime:               currentTime,
		policy:                    validationPolicy,
		poe:                       poe,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the process. Port of getTitle().
func (c *ValidationProcessForTimestampsWithArchivalData) Title() i18n.MessageTag {
	return i18n.MessageTag_VPFTSPWATSP
}

// InitChain initializes the chain. Port of initChain().
func (c *ValidationProcessForTimestampsWithArchivalData) InitChain() {
	lowestPOE := c.lowestPOE()
	lowestPOETime := lowestPOE.Time.Time()
	c.Result.Value.ProofOfExistence = lowestPOE

	basicTimestampConclusion := c.vpftsp.Conclusion

	var item process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp]

	/* Step 0. Execute detached evidence records processing, when applicable */
	evidenceRecords := c.timestamp.EvidenceRecords()
	if utils.IsCollectionNotEmpty(evidenceRecords) {
		for _, evidenceRecord := range evidenceRecords {
			evidenceRecordValidation := c.evidenceRecordValidation(evidenceRecord)
			if evidenceRecordValidation != nil {

				evidenceRecordValidationConclusive := c.evidenceRecordValidationConclusive(evidenceRecord,
					evidenceRecordValidation)

				if item == nil {
					item = evidenceRecordValidationConclusive
					c.FirstItem = item
				} else {
					item = item.SetNextItem(evidenceRecordValidationConclusive)
				}

				if c.IsValid(&evidenceRecordValidation.XmlConstraintsConclusionContent) {
					c.poe.ExtractEvidenceRecordPOE(evidenceRecord)
				}

			}
		}
	}

	timestampBasicSignatureValidationAcceptable := c.timestampBasicSignatureValidationAcceptable(c.vpftsp)
	if item == nil {
		item = timestampBasicSignatureValidationAcceptable
		c.FirstItem = item
	} else {
		item = item.SetNextItem(timestampBasicSignatureValidationAcceptable)
	}

	if process.IsAllowedBasicTimestampValidation(basicTimestampConclusion) {

		item = item.SetNextItem(c.timestampBasicSignatureValidationConclusive(c.timestamp, c.vpftsp))

		/*
		 * b) If PASSED is returned and a POE exists for the time-stamp for a time when the cryptographic hash
		 * function used in the time-stamp (messageImprint.hashAlgorithm) has been considered reliable, the SVA
		 * shall perform the POE extraction process (clause 5.6.2.3) with the signature, the time-stamp and the
		 * cryptographic constraints as inputs. The SVA shall add the returned POEs to the set of POEs.
		 */
		if c.IsValid(&c.vpftsp.XmlConstraintsConclusionContent) {

			aovResult := c.basicTimestampAlgorithmObsolescenceValidation(c.timestamp)

			digestMatchersValidation := aovResult.DigestMatchersValidation
			item = item.SetNextItem(c.timestampDigestAlgorithm(c.timestamp, digestMatchersValidation, lowestPOETime)) //nolint:staticcheck // mirrors upstream ValidationProcessForTimestampsWithArchivalData#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

			// NOTE: POE is extracted outside the class

		} else if process.IsAllowedBasicTimestampValidation(basicTimestampConclusion) {
			/*
			 * c) If the output of the validation is INDETERMINATE/REVOKED_NO_POE,
			 * INDETERMINATE/REVOKED_CA_NO_POE, INDETERMINATE/OUT_OF_BOUNDS_NO_POE,
			 * INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED,
			 * INDETERMINATE/CRYPTO_CONSTRAINTS_FAILURE_NO_POE or
			 * INDETERMINATE/REVOCATION_OUT_OF_BOUNDS_NO_POE,, the SVA shall perform past
			 * signature validation process (as per clause 5.6.2.4) with the following inputs: the time-stamp, the
			 * indication/sub-indication returned by the time-stamp validation process in step 5)a), the TSA's certificate,
			 * the X.509 validation parameters, X.509 validation constraints, cryptographic constraints, certificate
			 * validation data and the set of POEs. Then:
			 *
			 * i)    If it returns PASSED the SVA shall determine from the set of POEs the earliest time the existence
			 *       of the time-stamp can be proven.
			 *
			 * ii)   The SVA shall perform the Signature Acceptance Validation process as per clause 5.2.8 with the
			 *       following inputs:
			 *       - The Signed Data Object(s).
			 *       - The time determined in step i) above as the validation time parameter.
			 *       - The Cryptographic Constraints.
			 *       If the Signature Acceptance Validation process returns PASSED, the SVA shall go to the next step.
			 *       Otherwise, the SVA shall go to step d).
			 *
			 * iii) If a POE exists for the time-stamp for a time when the cryptographic hash function used in the
			 *      time-stamp has been considered reliable, the SVA shall perform the POE extraction process
			 *      (clause 5.6.2.3) and shall add the returned POEs to the set of POEs, and shall continue with
			 *      step 5)a) using the next time-stamp attribute.
			 */

			psv := vpfswatsp.NewPastSignatureValidation(c.I18nProvider, c.timestamp, c.bbbs,
				basicTimestampConclusion, c.poe, c.currentTime, c.policy, enumerations.Context_TIMESTAMP)
			psvResult := psv.Execute()

			tstBBB := c.bbbs[c.timestamp.Id()]
			c.enrichBBBWithPSVConclusion(tstBBB, psvResult)

			item = item.SetNextItem(c.pastTimestampValidation(c.timestamp, psvResult))

			/*
			 * If it returns PASSED and the cryptographic hash function used in the time-stamp is considered
			 * reliable at the generation time of the time-stamp, the long term validation process shall
			 * perform the POE extraction process and shall add the returned POEs to the set of POEs
			 * continue with step 5a using the next timestamp attribute.
			 */
			if c.IsValid(&psvResult.XmlConstraintsConclusionContent) {

				// NOTE: we execute AlgorithmObsolescenceValidation check, as the only time sensitive part of SAV
				algorithmObsolescenceValidation := aov.NewTimestampAlgorithmObsolescenceValidation(
					c.I18nProvider, c.timestamp, lowestPOETime, c.policy)
				aovResult := algorithmObsolescenceValidation.Execute()

				item = item.SetNextItem(c.timestampIsAcceptable(aovResult, lowestPOETime))

				digestMatchersValidation := aovResult.DigestMatchersValidation
				item = item.SetNextItem(c.timestampDigestAlgorithm(c.timestamp, digestMatchersValidation, //nolint:staticcheck // mirrors upstream ValidationProcessForTimestampsWithArchivalData#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
					lowestPOETime))

				// NOTE: POE is extracted outside the class

			}
		}
	}

}

// lowestPOE ports the private getLowestPOE(): Java dereferences the lowest POE
// unguarded, so a time-stamp with no POE at all raises a NullPointerException;
// the nil POE interface panics here in its place.
func (c *ValidationProcessForTimestampsWithArchivalData) lowestPOE() *jaxb.XmlProofOfExistence {
	lowestPOE := c.poe.GetLowestPOE(c.timestamp.Id())
	xmlProofOfExistence := &jaxb.XmlProofOfExistence{}
	xmlProofOfExistence.Time = jaxb.XSDateTime(lowestPOE.Time())
	return xmlProofOfExistence
}

// timestampBasicSignatureValidationAcceptable ports the private
// timestampBasicSignatureValidationAcceptable(XmlValidationProcessBasicTimestamp).
func (c *ValidationProcessForTimestampsWithArchivalData) timestampBasicSignatureValidationAcceptable(
	timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	return NewAcceptableBasicTimestampValidationCheck(c.I18nProvider, c.Result,
		&timestampValidationResult.XmlConstraintsConclusionContent, c.FailLevelRule())
}

// enrichBBBWithPSVConclusion ports the private
// enrichBBBWithPSVConclusion(XmlBasicBuildingBlocks, XmlPSV).
func (c *ValidationProcessForTimestampsWithArchivalData) enrichBBBWithPSVConclusion(
	bbb *jaxb.XmlBasicBuildingBlocks, psv *jaxb.XmlPSV) {
	bbb.PSV = psv

	bbbConclusion := bbb.Conclusion
	psvConclusion := psv.Conclusion
	bbbConclusion.Indication = psvConclusion.Indication
	bbbConclusion.SubIndication = psvConclusion.SubIndication
	bbbConclusion.Errors = append(bbbConclusion.Errors, psvConclusion.Errors...)
	bbbConclusion.Warnings = append(bbbConclusion.Warnings, psvConclusion.Warnings...)
	bbbConclusion.Infos = append(bbbConclusion.Infos, psvConclusion.Infos...)
}

// evidenceRecordValidation ports the private
// getEvidenceRecordValidation(EvidenceRecordWrapper): Java dereferences the map
// lookup unguarded, so an unknown evidence record raises a
// NullPointerException; the nil pointer panics here in its place.
func (c *ValidationProcessForTimestampsWithArchivalData) evidenceRecordValidation(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.XmlValidationProcessEvidenceRecord {
	xmlEvidenceRecord := c.evidenceRecordValidations[evidenceRecord.Id()]
	return xmlEvidenceRecord.ValidationProcessEvidenceRecord
}

// evidenceRecordValidationConclusive ports the private
// evidenceRecordValidationConclusive(EvidenceRecordWrapper, XmlValidationProcessEvidenceRecord).
func (c *ValidationProcessForTimestampsWithArchivalData) evidenceRecordValidationConclusive(
	evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper,
	erValidationResult *jaxb.XmlValidationProcessEvidenceRecord) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	return vpfswatsp.NewEvidenceRecordValidationCheck(c.I18nProvider, c.Result, evidenceRecordWrapper,
		erValidationResult, c.evidenceRecordValidationConstraintLevel())
}

// timestampBasicSignatureValidationConclusive ports the private
// timestampBasicSignatureValidationConclusive(TimestampWrapper, XmlValidationProcessBasicTimestamp).
func (c *ValidationProcessForTimestampsWithArchivalData) timestampBasicSignatureValidationConclusive(
	timestampWrapper *diagnostic.TimestampWrapper,
	timestampValidationResult *jaxb.XmlValidationProcessBasicTimestamp) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	return vpfbs.NewBasicTimestampValidationCheck(c.I18nProvider, c.Result, timestampWrapper,
		timestampValidationResult, c.WarnLevelRule())
}

// basicTimestampAlgorithmObsolescenceValidation ports the private
// getBasicTimestampAlgorithmObsolescenceValidation(TimestampWrapper).
func (c *ValidationProcessForTimestampsWithArchivalData) basicTimestampAlgorithmObsolescenceValidation(
	newestTimestamp *diagnostic.TimestampWrapper) *jaxb.XmlAOV {
	return c.bbbs[newestTimestamp.Id()].AOV
}

// pastTimestampValidation ports the private
// pastTimestampValidation(TimestampWrapper, XmlPSV).
func (c *ValidationProcessForTimestampsWithArchivalData) pastTimestampValidation(
	timestamp *diagnostic.TimestampWrapper,
	xmlPSV *jaxb.XmlPSV) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	return NewPastTimestampValidationCheck(c.I18nProvider, c.Result, timestamp, xmlPSV, c.FailLevelRule())
}

// timestampDigestAlgorithm ports the private
// timestampDigestAlgorithm(TimestampWrapper, XmlCryptographicValidation, Date).
func (c *ValidationProcessForTimestampsWithArchivalData) timestampDigestAlgorithm(
	timestampWrapper *diagnostic.TimestampWrapper, cvResult *jaxb.XmlCryptographicValidation,
	poeTime time.Time) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	// TODO : modify check to validate TST's POE
	return NewTimestampDigestAlgorithmValidationCheck(c.I18nProvider, c.Result, timestampWrapper, cvResult, poeTime,
		c.WarnLevelRule())
}

// timestampIsAcceptable ports the private timestampIsAcceptable(XmlAOV, Date):
// the unsupported-Context error becomes a panic, the caller being initChain,
// which cannot propagate one.
func (c *ValidationProcessForTimestampsWithArchivalData) timestampIsAcceptable(aovResult *jaxb.XmlAOV,
	lowestPOE time.Time) process.ChainItem[*jaxb.XmlValidationProcessArchivalDataTimestamp] {
	position, err := process.GetCryptoPosition(enumerations.Context_TIMESTAMP)
	if err != nil {
		panic(err)
	}
	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, lowestPOE, position,
		c.timestamp.Id())
}

// evidenceRecordValidationConstraintLevel ports the private
// getEvidenceRecordValidationConstraintLevel().
func (c *ValidationProcessForTimestampsWithArchivalData) evidenceRecordValidationConstraintLevel() policy.LevelRule {
	constraint := c.policy.EvidenceRecordValidConstraint()
	if constraint == nil {
		constraint = c.WarnLevelRule()
	}
	return constraint
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint): the generated BlockType member
// is a *XmlBlockType, whose nil is Java's null.
func (c *ValidationProcessForTimestampsWithArchivalData) CollectMessages(conclusion *jaxb.XmlConclusion,
	constraint *jaxb.XmlConstraint) {
	blockType := jaxb.XmlBlockType("")
	if constraint.BlockType != nil {
		blockType = *constraint.BlockType
	}
	if (jaxb.XmlBlockType_TST_BBB == blockType || jaxb.XmlBlockType_TST_PSV == blockType) &&
		c.policy.TimestampValidConstraint() == nil {
		// skip propagating of validation messages for TSTs in default processing
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion. Port
// of the overridden collectAdditionalMessages(XmlConclusion).
func (c *ValidationProcessForTimestampsWithArchivalData) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	if !process.IsAllowedBasicTimestampValidation(c.vpftsp.Conclusion) {
		conclusion.Warnings = append(conclusion.Warnings, c.vpftsp.Conclusion.Warnings...)
		conclusion.Infos = append(conclusion.Infos, c.vpftsp.Conclusion.Infos...)
	}
}
