// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/ValidationProcessForSignaturesWithArchivalData.java (DSS 6.5.RC1).
//
// See poe.go for the package-flattening note, and
// past_signature_validation_certificate_revocation_selector.go for the assumed
// Go shape of the cross-chunk vpfltvd classes - here TimestampDelayCheck, a
// plain generic ChainItem whose Go constructor is
//
//	vpfltvd.NewTimestampDelayCheck[T](
//	    *i18n.I18nProvider, *process.Result[T], *diagnostic.SignatureWrapper, time.Time, policy.DurationRule)
package vpfswatsp

import (
	"sort"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/aov"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/sav"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfltvd"
)

// ValidationProcessForSignaturesWithArchivalData performs 5.6 Validation
// process for Signatures with Archival Data.
type ValidationProcessForSignaturesWithArchivalData struct {
	*process.ChainBase[*jaxb.XmlValidationProcessArchivalData]

	// validationProcessLongTermData is the signature validation with long-term
	// data result.
	validationProcessLongTermData *jaxb.XmlValidationProcessLongTermData

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// signature is the signature.
	signature *diagnostic.SignatureWrapper

	// bbbs is the map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// xmlTimestamps is the list of timestamps.
	xmlTimestamps []*jaxb.XmlTimestamp

	// xmlEvidenceRecords is the list of evidence records.
	xmlEvidenceRecords []*jaxb.XmlEvidenceRecord

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// currentTime is the validation time.
	currentTime time.Time

	// poe is the POE container.
	poe *POEExtraction

	// context is the current validation context.
	context enumerations.Context

	// psvResult is the cached instance of PastSignatureValidation.
	psvResult *jaxb.XmlPSV
}

// NewValidationProcessForSignaturesWithArchivalData is the default constructor.
// Port of ValidationProcessForSignaturesWithArchivalData(I18nProvider, XmlSignature, SignatureWrapper, DiagnosticData, Map, ValidationPolicy, Date, POEExtraction).
func NewValidationProcessForSignaturesWithArchivalData(i18nProvider *i18n.I18nProvider,
	signatureAnalysis *jaxb.XmlSignature, signature *diagnostic.SignatureWrapper,
	diagnosticData *diagnostic.DiagnosticData, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	validationPolicy policy.ValidationPolicy, currentTime time.Time,
	poe *POEExtraction) *ValidationProcessForSignaturesWithArchivalData {
	xmlResult := &jaxb.XmlValidationProcessArchivalData{}
	c := &ValidationProcessForSignaturesWithArchivalData{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		validationProcessLongTermData: signatureAnalysis.ValidationProcessLongTermData,
		xmlTimestamps:                 signatureAnalysis.Timestamp,
		xmlEvidenceRecords:            signatureAnalysis.EvidenceRecord,
		signature:                     signature,
		diagnosticData:                diagnosticData,
		bbbs:                          bbbs,
		policy:                        validationPolicy,
		currentTime:                   currentTime,
		poe:                           poe,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the process. Port of getTitle().
func (c *ValidationProcessForSignaturesWithArchivalData) Title() i18n.MessageTag {
	return i18n.MessageTag_VPFSWATSP
}

// InitChain initializes the chain. Port of initChain().
func (c *ValidationProcessForSignaturesWithArchivalData) InitChain() {

	c.context = enumerations.Context_SIGNATURE
	if c.signature.IsCounterSignature() {
		c.context = enumerations.Context_COUNTER_SIGNATURE
	} else if c.signature.IsKeyBindingSignature() {
		c.context = enumerations.Context_KEY_BINDING_SIGNATURE
	}

	var item process.ChainItem[*jaxb.XmlValidationProcessArchivalData]

	/*
	 * 5.6.3.4 Processing
	 *
	 * 1) If there is one or more Evidence Records (ERs):
	 *
	 * a) The process shall take the first ER that was not yet processed.
	 * b) The process shall verify this ER according to IETF RFC 4998 [i.9] or IETF RFC 6283 [i.10] taking into
	 * account the following additional requirements when validating a time-stamp token at the time of the
	 * following Archive Timestamp:
	 */
	// steps from b) are performed within ValidationProcessEvidenceRecord
	/*
	 * c) If step b) found the ER to be valid, the process shall add a POE for every object covered by the ER at
	 * signing time value of the initial archive time-stamp.
	 * d) If all ERs have been validated, the process shall continue with step 2).
	 * e) The process shall continue with step 1)a).
	 */
	// TODO : verify evidence records indirectly covering the signature ?
	evidenceRecords := c.signature.EvidenceRecords()
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

	/*
	 * 2) The SVA shall add a POE for each object in the signature at the current time to the set of POEs.
	 */

	// POE provided to the validation

	/*
	 * 3) The SVA shall perform the Validation process for Signatures with Time and Signatures with
	 * Long-Term Validation Material as per clause 5.5 with all the inputs, including the processing
	 * of any signed attributes as specified.
	 *
	 * - If the signature does not contain any attributes for long term availability and integrity of
	 *   validation material, the process shall return the indication/sub-indication and information returned
	 *   by the Validation process for Signatures with Time and Signatures with Long-Term Validation Material.
	 *   Additional information should be included indicating that only the signature-with-time-validation
	 *   process has been performed.
	 *
	 * - If the Validation process for Signatures with Time and Signatures with Long-Term Validation Material
	 *   returned PASSED:
	 * -- If there is no validation constraint mandating the validation of the LTV attributes, the long term
	 *    validation process shall return the indication PASSED.
	 * -- Otherwise, the SVA shall go to step 4.
	 *
	 * - If the Validation process for Signatures with Time and Signatures with Long-Term Validation Material
	 *   returned one of the following indications/sub-indications: INDETERMINATE/REVOKED_NO_POE,
	 *   INDETERMINATE/REVOKED_CA_NO_POE, INDETERMINATE/OUT_OF_BOUNDS_NO_POE,
	 *   INDETERMINATE/OUT_OF_BOUNDS_NOT_REVOKED,
	 *   INDETERMINATE/CRYPTO_CONSTRAINTS_FAILURE_NO_POE,
	 *   INDETERMINATE/REVOCATION_OUT_OF_BOUNDS_NO_POE,
	 *   INDETERMINATE/SIG_CONSTRAINTS_FAILURE or INDETERMINATE/TRY_LATER, the long-term
	 *   validation process shall go to the next step.
	 *
	 * - In all other cases, the process shall return the indication/sub-indication and information
	 *   returned by the Validation process for Signatures with Time and Signatures with
	 *   Long-Term Validation Material.
	 */
	if item == nil {
		item = c.longTermValidation()
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.longTermValidation())
	}

	proofOfExistence := c.validationProcessLongTermData.ProofOfExistence
	c.Result.Value.ProofOfExistence = proofOfExistence

	/*
	 * Return long-term validation process indication
	 */
	if !process.IsAllowedValidationWithLongTermData(c.validationProcessLongTermData.Conclusion) {
		return
	}

	// TODO : this is a workaround, as LTV is not able to recover from NO_CERTIFICATE_CHAIN_FOUND_NO_POE
	if enumerations.Indication_INDETERMINATE != c.validationProcessLongTermData.Conclusion.Indication.Indication() ||
		enumerations.SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE !=
			subIndicationOfConclusion(c.validationProcessLongTermData.Conclusion) {

		item = item.SetNextItem(c.longTermAvailabilityAndIntegrityValidationMaterial())

		if !process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.signature) {
			return
		}

	}

	/*
	 * 4) The process shall initialize best-signature-time to the best-signature-time returned in step 3)
	 * and add this time as POE for the signature to the set of POEs.
	 */
	bestSignatureTime := c.toPOE(proofOfExistence)
	c.poe.AddSignaturePOE(c.signature, bestSignatureTime)

	/*
	 * 5) If there is at least one time-stamp attribute:
	 *
	 * a) The long term validation process shall select the newest time-stamp that has not been processed and
	 * perform the time-stamp validation, as per clause 5.4.
	 */
	timestampsList := c.signature.TimestampList()
	if utils.IsCollectionNotEmpty(timestampsList) {
		// Comparator.comparing(TimestampWrapper::getProductionTime).reversed()
		// over a stable List#sort; a missing production time raises a
		// NullPointerException there and dereferences here in its place.
		sort.SliceStable(timestampsList, func(i, j int) bool {
			return timestampsList[i].ProductionTime().After(*timestampsList[j].ProductionTime())
		})

		for _, newestTimestamp := range timestampsList {
			bbbTsp := c.bbbs[newestTimestamp.Id()]
			timestampValidation := c.timestampValidation(newestTimestamp)
			if bbbTsp != nil && timestampValidation != nil {

				// steps b) and c) are part of ValidationProcessForTimestampsWithArchivalData

				item = item.SetNextItem(c.timestampValidationConclusive(newestTimestamp, timestampValidation))

				/*
				 * d) In all other cases:
				 * - If no specific constraints mandating the validity of the attribute are specified in the
				 * validation constraints, the SVA shall ignore the attribute and shall continue with step 5 using
				 * the next timestamp attribute.
				 * - Otherwise, the process shall fail with the returned indication/sub-indication and associated
				 * explanations.
				 */

			}

			/*
			 * e) If all time-stamp attributes have been processed, the SVA shall continue with step 6.
			 * Otherwise, the SVA shall continue with step 5a.
			 */
		}

	}

	/*
	 * Performs optional checks allowing to fail validation process in case of a missing timestamp
	 */
	item = item.SetNextItem(c.tLevelTimeStamp())

	item = item.SetNextItem(c.ltaLevelTimeStamp())

	/*
	 * 6) The SVA shall determine from the set of POEs the earliest time the existence of the signature
	 * can be proven and set best-signature-time to this new determined time.
	 */
	bestSignatureTime = c.poe.GetLowestPOE(c.signature.Id())
	c.Result.Value.ProofOfExistence = toXmlProofOfExistenceForPOE(bestSignatureTime)

	if c.IsValid(&c.validationProcessLongTermData.XmlConstraintsConclusionContent) {
		// skip past signature validation when basic validation succeeded
		return
	}

	/*
	 * 7) Past signature validation: the long term validation process shall perform the past signature validation
	 * process with the following inputs: the signature, the status indication/sub-indication returned in step 2,
	 * the signing certificate, the X.509 validation parameters, certificate validation data, chain constraints,
	 * cryptographic constraints, the set of POEs and best-signature-time. If it returns PASSED the long term
	 * validation process shall go to the next step. Otherwise, the long term validation process shall return
	 * the indication/sub-indication and associated explanations returned from the past signature validation process.
	 */
	sigBBB := c.bbbs[c.signature.Id()]
	psv := NewPastSignatureValidation(c.I18nProvider, c.signature, c.bbbs,
		c.validationProcessLongTermData.Conclusion, c.poe, c.currentTime, c.policy, c.context)
	c.psvResult = psv.Execute()
	sigBBB.PSV = c.psvResult
	c.enrichBBBWithPSVConclusion(sigBBB, c.psvResult)

	item = item.SetNextItem(c.pastSignatureValidation(c.psvResult))

	/*
	 * 8) Handling time-stamp delay: If the signature contains a signature time stamp token and the validation
	 * specify a time stamp delay:
	 * a) If no signing time property/attribute is present, the process shall return the indication INDETERMINATE
	 *    with the sub indication SIG_CONSTRAINTS_FAILURE.
	 * b) If a signing time property/attribute is present, the process shall check that the claimed time in the
	 *    attribute plus the time stamp delay is after the best-signature-time determined in step 6) above. If
	 *    the check is successful, the process shall go to the next step. Otherwise, the process shall return
	 *    the indication INDETERMINATE with the sub indication SIG_CONSTRAINTS_FAILURE.
	 */
	item = item.SetNextItem(c.timestampDelay(bestSignatureTime.Time()))

	/*
	 * 9) The SVA shall perform the Signature Acceptance Validation process as per clause 5.2.8 with the following
	 * inputs:
	 * a) The Signed Data Object(s).
	 * b) The time determined in step 7 as the validation time parameter.
	 * c) The Cryptographic Constraints.
	 *
	 * If the Signature Acceptance Validation process returns PASSED, the SVA shall go to the next step. Otherwise,
	 * the SVA shall return the indication and sub-indication returned by the Signature Acceptance Validation
	 * Process.
	 */
	item = item.SetNextItem(c.signatureIsAcceptable(bestSignatureTime.Time(), c.context)) //nolint:staticcheck // mirrors upstream ValidationProcessForSignaturesWithArchivalData#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

	/*
	 * 10) Data extraction: the SVA shall return the success indication PASSED. In addition, the SVA should return
	 * additional information extracted from the signature and/or used by the intermediate steps. In particular, the
	 * SVA should return the best-signature-time determined in step 6) as well as intermediate results such as the
	 * validation results of any time-stamp token.
	 */

	/*
	 * Otherwise,
	 * the SVA shall return the indication and sub-indication returned by the Signature Acceptance Validation Process
	 */

}

// pastSignatureValidation ports the private pastSignatureValidation(XmlPSV).
func (c *ValidationProcessForSignaturesWithArchivalData) pastSignatureValidation(
	xmlPSV *jaxb.XmlPSV) process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	return NewPastSignatureValidationCheck(c.I18nProvider, c.Result, c.signature, xmlPSV, c.FailLevelRule())
}

// longTermValidation ports the private longTermValidation().
func (c *ValidationProcessForSignaturesWithArchivalData) longTermValidation() process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	return NewLongTermValidationCheck(c.I18nProvider, c.Result,
		&c.validationProcessLongTermData.XmlConstraintsConclusionContent, c.FailLevelRule())
}

// longTermAvailabilityAndIntegrityValidationMaterial ports the private
// longTermAvailabilityAndIntegrityValidationMaterial().
func (c *ValidationProcessForSignaturesWithArchivalData) longTermAvailabilityAndIntegrityValidationMaterial() process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	var constraint policy.LevelRule
	if c.IsValid(&c.validationProcessLongTermData.XmlConstraintsConclusionContent) {
		constraint = c.InfoLevelRule()
	} else {
		constraint = c.FailLevelRule()
	}
	return NewLongTermAvailabilityAndIntegrityValidationMaterialCheck(c.I18nProvider, c.Result, c.signature,
		&c.validationProcessLongTermData.XmlConstraintsConclusionContent, constraint)
}

// timestampValidationConclusive ports the private
// timestampValidationConclusive(TimestampWrapper, XmlValidationProcessArchivalDataTimestamp).
func (c *ValidationProcessForSignaturesWithArchivalData) timestampValidationConclusive(
	timestampWrapper *diagnostic.TimestampWrapper,
	timestampValidationResult *jaxb.XmlValidationProcessArchivalDataTimestamp) process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	return NewTimestampValidationCheck(c.I18nProvider, c.Result, timestampWrapper, timestampValidationResult,
		c.timestampValidationConstraintLevel())
}

// evidenceRecordValidationConclusive ports the private
// evidenceRecordValidationConclusive(EvidenceRecordWrapper, XmlValidationProcessEvidenceRecord).
func (c *ValidationProcessForSignaturesWithArchivalData) evidenceRecordValidationConclusive(
	evidenceRecordWrapper *diagnostic.EvidenceRecordWrapper,
	erValidationResult *jaxb.XmlValidationProcessEvidenceRecord) process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	return NewEvidenceRecordValidationCheck(c.I18nProvider, c.Result, evidenceRecordWrapper, erValidationResult,
		c.evidenceRecordValidationConstraintLevel())
}

// timestampValidationConstraintLevel ports the private
// getTimestampValidationConstraintLevel().
func (c *ValidationProcessForSignaturesWithArchivalData) timestampValidationConstraintLevel() policy.LevelRule {
	constraint := c.policy.TimestampValidConstraint()
	if constraint == nil {
		constraint = c.WarnLevelRule()
	}
	return constraint
}

// evidenceRecordValidationConstraintLevel ports the private
// getEvidenceRecordValidationConstraintLevel().
func (c *ValidationProcessForSignaturesWithArchivalData) evidenceRecordValidationConstraintLevel() policy.LevelRule {
	constraint := c.policy.EvidenceRecordValidConstraint()
	if constraint == nil {
		constraint = c.WarnLevelRule()
	}
	return constraint
}

// tLevelTimeStamp ports the private tLevelTimeStamp().
func (c *ValidationProcessForSignaturesWithArchivalData) tLevelTimeStamp() process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	constraint := c.policy.TLevelTimeStampConstraint(c.context)
	return sav.NewTLevelTimeStampCheck(c.I18nProvider, c.Result, c.signature, c.bbbs, c.xmlTimestamps, constraint)
}

// ltaLevelTimeStamp ports the private ltaLevelTimeStamp().
func (c *ValidationProcessForSignaturesWithArchivalData) ltaLevelTimeStamp() process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	constraint := c.policy.LTALevelTimeStampConstraint(c.context)
	return sav.NewLTALevelTimeStampCheck(c.I18nProvider, c.Result, c.signature, c.bbbs, c.xmlTimestamps, constraint)
}

// timestampDelay ports the private timestampDelay(Date).
func (c *ValidationProcessForSignaturesWithArchivalData) timestampDelay(
	bestSignatureTime time.Time) process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	// The vpfltvd check keeps Java's nullable Date; the best-signature-time
	// reached here is always the non-null time of the lowest POE.
	return vpfltvd.NewTimestampDelayCheck(c.I18nProvider, c.Result, c.signature, &bestSignatureTime,
		c.policy.TimestampDelayConstraint())
}

// signatureIsAcceptable ports the private signatureIsAcceptable(Date, Context):
// the unsupported-Context error becomes a panic, the caller being initChain,
// which cannot propagate one.
//
// NOTE: we execute AlgorithmObsolescenceValidation check, as the only time
// sensitive part of SAV.
func (c *ValidationProcessForSignaturesWithArchivalData) signatureIsAcceptable(bestSignatureTime time.Time,
	context enumerations.Context) process.ChainItem[*jaxb.XmlValidationProcessArchivalData] {
	algorithmObsolescenceValidation := aov.NewSignatureAlgorithmObsolescenceValidation(
		c.I18nProvider, c.signature, context, bestSignatureTime, c.policy)
	aovResult := algorithmObsolescenceValidation.Execute()

	position, err := process.GetCryptoPosition(context)
	if err != nil {
		panic(err)
	}

	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, aovResult, bestSignatureTime,
		position, c.signature.Id())
}

// enrichBBBWithPSVConclusion ports the private
// enrichBBBWithPSVConclusion(XmlBasicBuildingBlocks, XmlPSV).
func (c *ValidationProcessForSignaturesWithArchivalData) enrichBBBWithPSVConclusion(
	bbb *jaxb.XmlBasicBuildingBlocks, psv *jaxb.XmlPSV) {
	bbbConclusion := bbb.Conclusion
	psvConclusion := psv.Conclusion
	bbbConclusion.Indication = psvConclusion.Indication
	bbbConclusion.SubIndication = psvConclusion.SubIndication
	bbbConclusion.Errors = append(bbbConclusion.Errors, psvConclusion.Errors...)
	bbbConclusion.Warnings = append(bbbConclusion.Warnings, psvConclusion.Warnings...)
	bbbConclusion.Infos = append(bbbConclusion.Infos, psvConclusion.Infos...)
}

// toXmlProofOfExistenceForPOE ports the private toXmlProofOfExistence(POE).
func toXmlProofOfExistenceForPOE(poe POE) *jaxb.XmlProofOfExistence {
	xmlPoe := &jaxb.XmlProofOfExistence{}
	xmlPoe.Time = jaxb.XSDateTime(poe.Time())
	xmlPoe.TimestampId = poe.POEProviderId()
	return xmlPoe
}

// toPOE ports the private toPOE(XmlProofOfExistence): the generated TimestampId
// member is a *string, whose nil is Java's null.
func (c *ValidationProcessForSignaturesWithArchivalData) toPOE(xmlProofOfExistence *jaxb.XmlProofOfExistence) POE {
	timestampId := xmlProofOfExistence.TimestampId
	if timestampId != nil {
		for _, timestamp := range c.diagnosticData.TimestampList() {
			if *timestampId == timestamp.Id() {
				return NewTimestampPOE(timestamp)
			}
		}
		// Should not happen, as current revision of the standard does not handle ERs within LTV process.
		for _, evidenceRecordWrapper := range c.diagnosticData.EvidenceRecords() {
			if *timestampId == evidenceRecordWrapper.Id() {
				return NewEvidenceRecordPOE(evidenceRecordWrapper)
			}
		}
	}
	return NewPOE(xmlProofOfExistence.Time.Time())
}

// timestampValidation ports the private getTimestampValidation(TimestampWrapper):
// the generated XmlTimestamp Id member is a *string, whose nil is Java's null,
// and Utils.areStringsEqual(null, x) is false unless both are null.
func (c *ValidationProcessForSignaturesWithArchivalData) timestampValidation(
	newestTimestamp *diagnostic.TimestampWrapper) *jaxb.XmlValidationProcessArchivalDataTimestamp {
	for _, xmlTimestamp := range c.xmlTimestamps {
		if xmlTimestamp.Id != nil && utils.AreStringsEqual(*xmlTimestamp.Id, newestTimestamp.Id()) {
			return xmlTimestamp.ValidationProcessArchivalDataTimestamp
		}
	}
	return nil
}

// evidenceRecordValidation ports the private
// getEvidenceRecordValidation(EvidenceRecordWrapper).
func (c *ValidationProcessForSignaturesWithArchivalData) evidenceRecordValidation(
	evidenceRecord *diagnostic.EvidenceRecordWrapper) *jaxb.XmlValidationProcessEvidenceRecord {
	for _, xmlEvidenceRecord := range c.xmlEvidenceRecords {
		if xmlEvidenceRecord.Id != nil && utils.AreStringsEqual(*xmlEvidenceRecord.Id, evidenceRecord.Id()) {
			return xmlEvidenceRecord.ValidationProcessEvidenceRecord
		}
	}
	return nil
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint): the generated BlockType member
// is a *XmlBlockType, whose nil is Java's null.
func (c *ValidationProcessForSignaturesWithArchivalData) CollectMessages(conclusion *jaxb.XmlConclusion,
	constraint *jaxb.XmlConstraint) {
	blockType := jaxb.XmlBlockType("")
	if constraint.BlockType != nil {
		blockType = *constraint.BlockType
	}
	if jaxb.XmlBlockType_TST == blockType && c.policy.TimestampValidConstraint() == nil {
		// skip propagating of validation messages for TSTs in default processing
	} else if jaxb.XmlBlockType_LTA == blockType {
		// skip LTA data missing message
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion. Port
// of the overridden collectAdditionalMessages(XmlConclusion).
func (c *ValidationProcessForSignaturesWithArchivalData) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	if !process.IsAllowedValidationWithLongTermData(c.validationProcessLongTermData.Conclusion) ||
		!process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.signature) {
		conclusion.Warnings = append(conclusion.Warnings, c.validationProcessLongTermData.Conclusion.Warnings...)
		conclusion.Infos = append(conclusion.Infos, c.validationProcessLongTermData.Conclusion.Infos...)
	}
	if c.psvResult != nil {
		if !c.IsValid(&c.psvResult.XmlConstraintsConclusionContent) &&
			subIndicationOfConclusion(c.psvResult.Conclusion) !=
				subIndicationOfConclusion(c.validationProcessLongTermData.Conclusion) {
			conclusion.Errors = append(conclusion.Errors, c.psvResult.Conclusion.Errors...)
		}
		conclusion.Warnings = append(conclusion.Warnings, c.psvResult.Conclusion.Warnings...)
		conclusion.Infos = append(conclusion.Infos, c.psvResult.Conclusion.Infos...)
	}
}
