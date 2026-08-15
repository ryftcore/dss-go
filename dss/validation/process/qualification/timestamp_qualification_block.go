// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/timestamp/TimestampQualificationBlock.java (DSS 6.5.RC1).
//
// FLAGGED HASH-ORDER SITE: as in CertificateQualificationBlock and
// SignatureQualificationBlock, Java iterates HashSet<String>
// listOfTrustedListUrls/trustedListUrls directly, feeding report Constraint
// order. This port sorts the equivalent Go map[string]struct{} sets by URL
// via sortedKeys for a deterministic (if not necessarily Java-bucket-
// identical) result. See the porter brief's hard rule.
package qualification

import (
	"github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/vpfswatsp"
)

// TimestampQualificationBlock performs a qualification verification for a
// timestamp.
type TimestampQualificationBlock struct {
	*process.ChainBase[*jaxb.XmlValidationTimestampQualification]

	// timestamp is the timestamp to be validated.
	timestamp *diagnostic.TimestampWrapper

	// tlAnalysis is the list of all TL analyses.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// poe contains list of all POEs.
	poe *vpfswatsp.POEExtraction

	// relatedTLAnalyses is the list of related LOTL/TL analyses.
	relatedTLAnalyses []*jaxb.XmlTLAnalysis

	// tstQualification is the determined timestamp qualification.
	tstQualification enumerations.TimestampQualification
}

// NewTimestampQualificationBlock is the default constructor. Port of
// TimestampQualificationBlock(I18nProvider, TimestampWrapper, List, POEExtraction).
func NewTimestampQualificationBlock(i18nProvider *i18n.I18nProvider, timestamp *diagnostic.TimestampWrapper,
	tlAnalysis []*jaxb.XmlTLAnalysis, poe *vpfswatsp.POEExtraction) *TimestampQualificationBlock {
	xmlResult := &jaxb.XmlValidationTimestampQualification{}
	c := &TimestampQualificationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		timestamp:        timestamp,
		tlAnalysis:       tlAnalysis,
		poe:              poe,
		tstQualification: enumerations.TimestampQualification_NA,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *TimestampQualificationBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_TST_QUALIFICATION
}

// InitChain initializes the chain. Port of initChain().
func (c *TimestampQualificationBlock) InitChain() {
	signingCertificate := c.timestamp.SigningCertificate()

	item := c.isTrustedListReachedForCertificateChain(signingCertificate)
	c.FirstItem = item

	if signingCertificate != nil && signingCertificate.IsTrustedListReached() {

		originalTSPs := signingCertificate.TrustServices()

		listOfTrustedListUrls := map[string]struct{}{}
		for _, t := range originalTSPs {
			if t.ListOfTrustedLists != nil && t.ListOfTrustedLists.Url != nil {
				listOfTrustedListUrls[*t.ListOfTrustedLists.Url] = struct{}{}
			}
		}

		acceptableLOTLUrls := map[string]struct{}{}
		for _, lotlURL := range sortedKeys(listOfTrustedListUrls) {
			lotlAnalysis := c.getTlAnalysis(lotlURL)
			if lotlAnalysis != nil {
				c.relatedTLAnalyses = append(c.relatedTLAnalyses, lotlAnalysis)

				acceptableLOTL := c.isAcceptableLOTL(lotlAnalysis)
				item = item.SetNextItem(acceptableLOTL)
				if acceptableLOTL.Process() {
					acceptableLOTLUrls[lotlURL] = struct{}{}
				}
			}
		}

		// filter TLs with a found valid set of LOTLs (if assigned)
		trustedListUrls := map[string]struct{}{}
		for _, t := range originalTSPs {
			if t.TrustedList == nil || t.TrustedList.Url == nil {
				continue
			}
			if t.ListOfTrustedLists != nil {
				if t.ListOfTrustedLists.Url == nil {
					continue
				}
				if _, ok := acceptableLOTLUrls[*t.ListOfTrustedLists.Url]; !ok {
					continue
				}
			}
			trustedListUrls[*t.TrustedList.Url] = struct{}{}
		}

		acceptableTLUrls := map[string]struct{}{}
		if len(trustedListUrls) > 0 {
			for _, tlURL := range sortedKeys(trustedListUrls) {
				currentTL := c.getTlAnalysis(tlURL)
				if currentTL != nil {
					c.relatedTLAnalyses = append(c.relatedTLAnalyses, currentTL)

					acceptableTL := c.isAcceptableTL(currentTL)
					item = item.SetNextItem(acceptableTL)
					if acceptableTL.Process() {
						acceptableTLUrls[tlURL] = struct{}{}
					}
				}
			}
		}

		item = item.SetNextItem(c.isAcceptableTLPresent(acceptableTLUrls))

		if len(acceptableTLUrls) > 0 {

			filter := TrustServicesFilterFactoryCreateFilterByUrls(acceptableTLUrls)
			acceptableServices := filter.Filter(originalTSPs)

			tstQualificationAtGenerationTimeBlock := NewTimestampQualificationAtTimeBlockAtGenerationTime(
				c.I18nProvider, enumerations.ValidationTime_TIMESTAMP_GENERATION_TIME, c.timestamp, acceptableServices)
			conclusionAtGenerationTime := tstQualificationAtGenerationTimeBlock.Execute()
			c.Result.Value.ValidationTimestampQualificationAtTime = append(
				c.Result.Value.ValidationTimestampQualificationAtTime, conclusionAtGenerationTime)

			timestampPOE := c.poe.GetLowestPOETime(c.timestamp.Id())
			tstQualificationAtPOETimeBlock := NewTimestampQualificationAtTimeBlock(
				c.I18nProvider, enumerations.ValidationTime_TIMESTAMP_POE_TIME, &timestampPOE, c.timestamp, acceptableServices)
			conclusionAtPOETime := tstQualificationAtPOETimeBlock.Execute()
			c.Result.Value.ValidationTimestampQualificationAtTime = append(
				c.Result.Value.ValidationTimestampQualificationAtTime, conclusionAtPOETime)

			c.determineFinalQualification(conclusionAtGenerationTime.TimestampQualification.TimestampQualification(),
				conclusionAtPOETime.TimestampQualification.TimestampQualification())
		}
	}
}

// getTlAnalysis ports the private getTlAnalysis(String).
func (c *TimestampQualificationBlock) getTlAnalysis(url string) *jaxb.XmlTLAnalysis {
	for _, xmlTLAnalysis := range c.tlAnalysis {
		if url == xmlTLAnalysis.URL {
			return xmlTLAnalysis
		}
	}
	return nil
}

// determineFinalQualification ports the private
// determineFinalQualification(TimestampQualification, TimestampQualification).
func (c *TimestampQualificationBlock) determineFinalQualification(qualAtGenerationTime, qualAtPOETime enumerations.TimestampQualification) {
	if enumerations.TimestampQualification_QTSA == qualAtGenerationTime && enumerations.TimestampQualification_QTSA == qualAtPOETime {
		c.tstQualification = enumerations.TimestampQualification_QTSA
	} else {
		c.tstQualification = enumerations.TimestampQualification_TSA
	}
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *TimestampQualificationBlock) AddAdditionalInfo() {
	c.setIndication()
	c.setTimestampQualification()
}

// setIndication ports the private setIndication().
func (c *TimestampQualificationBlock) setIndication() {
	conclusion := c.Result.Conclusion()
	if conclusion == nil {
		return
	}
	indication := jaxb.IndicationValue(enumerations.Indication_PASSED)
	if len(conclusion.Errors) > 0 {
		indication = jaxb.IndicationValue(enumerations.Indication_FAILED)
	} else if len(conclusion.Warnings) > 0 {
		indication = jaxb.IndicationValue(enumerations.Indication_INDETERMINATE)
	}
	conclusion.Indication = indication
}

// setTimestampQualification ports the private setTimestampQualification().
func (c *TimestampQualificationBlock) setTimestampQualification() {
	c.Result.Value.TimestampQualification = jaxb.TimestampQualificationValue(c.tstQualification)
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden protected void
// collectAdditionalMessages(XmlConclusion).
func (c *TimestampQualificationBlock) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	for _, tstQualAtTime := range c.Result.Value.ValidationTimestampQualificationAtTime {
		c.CollectAllMessages(conclusion, tstQualAtTime.Conclusion)
	}
	for _, relatedTLAnalysis := range c.relatedTLAnalyses {
		c.CollectAllMessages(conclusion, relatedTLAnalysis.Conclusion)
	}
}

func (c *TimestampQualificationBlock) isTrustedListReachedForCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlValidationTimestampQualification] {
	return NewTrustedListReachedForCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, c.FailLevelRule())
}

func (c *TimestampQualificationBlock) isAcceptableLOTL(xmlLOTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableListOfTrustedListsCheck[*jaxb.XmlValidationTimestampQualification] {
	return NewAcceptableListOfTrustedListsCheck(c.I18nProvider, c.Result, xmlLOTLAnalysis, c.WarnLevelRule())
}

func (c *TimestampQualificationBlock) isAcceptableTL(xmlTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableTrustedListCheck[*jaxb.XmlValidationTimestampQualification] {
	return NewAcceptableTrustedListCheck(c.I18nProvider, c.Result, xmlTLAnalysis, c.WarnLevelRule())
}

func (c *TimestampQualificationBlock) isAcceptableTLPresent(acceptableUrls map[string]struct{}) process.ChainItem[*jaxb.XmlValidationTimestampQualification] {
	return NewAcceptableTrustedListPresenceCheck(c.I18nProvider, c.Result, acceptableUrls, c.FailLevelRule())
}
