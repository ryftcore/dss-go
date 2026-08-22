// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/timestamp/TimestampQualificationBlock.java (DSS 6.5.RC1).
//
// Java iterates HashSet<String> listOfTrustedListUrls/trustedListUrls
// directly, feeding report Constraint order, as in
// CertificateQualificationBlock and SignatureQualificationBlock. This port
// reproduces that HashSet iteration order via orderedURLSet (see
// certificate_qualification_block.go), whose iterate() calls
// utils.JavaHashMapStringKeyOrder.
package qualification

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/vpfswatsp"
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
		tstQualification: enumerations.TimestampQualificationNA,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *TimestampQualificationBlock) Title() i18n.MessageTag {
	return i18n.MessageTagTSTQualification
}

// InitChain initializes the chain. Port of initChain().
func (c *TimestampQualificationBlock) InitChain() {
	signingCertificate := c.timestamp.SigningCertificate()

	item := c.isTrustedListReachedForCertificateChain(signingCertificate)
	c.FirstItem = item

	if signingCertificate != nil && signingCertificate.IsTrustedListReached() {

		originalTSPs := signingCertificate.TrustServices()

		listOfTrustedListUrls := newOrderedURLSet()
		for _, t := range originalTSPs {
			if t.ListOfTrustedLists != nil && t.ListOfTrustedLists.Url != nil {
				listOfTrustedListUrls.add(*t.ListOfTrustedLists.Url)
			}
		}

		acceptableLOTLUrls := map[string]struct{}{}
		for _, lotlURL := range listOfTrustedListUrls.iterate() {
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
		trustedListUrls := newOrderedURLSet()
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
			trustedListUrls.add(*t.TrustedList.Url)
		}

		acceptableTLUrls := map[string]struct{}{}
		if trustedListUrls.len() > 0 {
			for _, tlURL := range trustedListUrls.iterate() {
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

		item = item.SetNextItem(c.isAcceptableTLPresent(acceptableTLUrls)) //nolint:staticcheck // mirrors upstream TimestampQualificationBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

		if len(acceptableTLUrls) > 0 {

			filter := TrustServicesFilterFactoryCreateFilterByUrls(acceptableTLUrls)
			acceptableServices := filter.Filter(originalTSPs)

			tstQualificationAtGenerationTimeBlock := NewTimestampQualificationAtTimeBlockAtGenerationTime(
				c.I18nProvider, enumerations.ValidationTimeTimestampGenerationTime, c.timestamp, acceptableServices)
			conclusionAtGenerationTime := tstQualificationAtGenerationTimeBlock.Execute()
			c.Result.Value.ValidationTimestampQualificationAtTime = append(
				c.Result.Value.ValidationTimestampQualificationAtTime, conclusionAtGenerationTime)

			timestampPOE := c.poe.GetLowestPOETime(c.timestamp.Id())
			tstQualificationAtPOETimeBlock := NewTimestampQualificationAtTimeBlock(
				c.I18nProvider, enumerations.ValidationTimeTimestampPOETime, &timestampPOE, c.timestamp, acceptableServices)
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
	if enumerations.TimestampQualificationQTSA == qualAtGenerationTime && enumerations.TimestampQualificationQTSA == qualAtPOETime {
		c.tstQualification = enumerations.TimestampQualificationQTSA
	} else {
		c.tstQualification = enumerations.TimestampQualificationTSA
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
	indication := jaxb.IndicationValue(enumerations.IndicationPassed)
	if len(conclusion.Errors) > 0 {
		indication = jaxb.IndicationValue(enumerations.IndicationFailed)
	} else if len(conclusion.Warnings) > 0 {
		indication = jaxb.IndicationValue(enumerations.IndicationIndeterminate)
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
