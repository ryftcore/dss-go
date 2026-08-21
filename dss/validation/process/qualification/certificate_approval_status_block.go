// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/usage/CertificateApprovalStatusBlock.java (DSS 6.5.RC1).
//
// CROSS-CHUNK ASSUMPTION: TrustedEntityServiceFilter and
// TrustedEntitiesFilterFactoryCreateFilterByListUrls (Java package
// qualification.trust.filter) are owned by a sibling porter of this shared
// package and were not present on disk while this file was written; see
// cert_qualification_at_time_block.go's header for the established
// flattening convention this file's call sites follow.
package qualification

import (
	"sort"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateApprovalStatusBlock verifies the certificate's approval status.
type CertificateApprovalStatusBlock struct {
	*process.ChainBase[*jaxb.XmlCertificateApprovalStatusProcess]

	// BuildingBlocksConclusion is the certificate's BasicBuildingBlock's
	// conclusion.
	BuildingBlocksConclusion *jaxb.XmlConclusion

	// ValidationTime is the validation time.
	ValidationTime time.Time

	// SigningCertificate is the certificate to determine qualification for.
	SigningCertificate *diagnostic.CertificateWrapper

	// LoteAnalysis is a list of validation results for all Trusted Lists.
	LoteAnalysis []*jaxb.XmlLoTEAnalysis
}

// NewCertificateApprovalStatusBlock is the default constructor. Port of
// CertificateApprovalStatusBlock(I18nProvider, XmlConclusion, Date, CertificateWrapper, List).
func NewCertificateApprovalStatusBlock(i18nProvider *i18n.I18nProvider, buildingBlocksConclusion *jaxb.XmlConclusion,
	validationTime time.Time, signingCertificate *diagnostic.CertificateWrapper,
	loteAnalysis []*jaxb.XmlLoTEAnalysis) *CertificateApprovalStatusBlock {
	xmlResult := &jaxb.XmlCertificateApprovalStatusProcess{}
	c := &CertificateApprovalStatusBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		BuildingBlocksConclusion: buildingBlocksConclusion,
		ValidationTime:           validationTime,
		SigningCertificate:       signingCertificate,
		LoteAnalysis:             loteAnalysis,
	}

	id := signingCertificate.Id()
	c.Result.Value.Id = &id

	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *CertificateApprovalStatusBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_CERT_USAGES
}

// InitChain initializes the chain. Port of initChain().
//
// HASH-ORDER (closed in phase 8f): Java builds listsOfLists/lotes as
// HashSet<XmlTrustSourceList> (identity hashCode/equals) and listsBYType as
// a HashMap<String, List<XmlTrustSourceList>>, iterating both directly to
// append checks to the report / to decide the order sub-blocks execute in -
// both observable in the report. This port iterates the Go equivalents
// (map[*jaxb.XmlTrustSourceList]struct{}, map[string][]..., whose range
// order is randomized per the language spec) in a deterministic surrogate
// order instead - sorted by each list's URL - which is stable across runs
// but not necessarily identical to the upstream HashSet/HashMap bucket
// order. See the porter brief's hard rule on hash-order leaks.
func (c *CertificateApprovalStatusBlock) InitChain() {
	// cover incomplete cert chain / expired/ revoked certs
	item := c.isAcceptableBuildingBlockConclusion(c.BuildingBlocksConclusion)
	c.FirstItem = item

	acceptableLoTEs := map[*diagjaxb.XmlTrustSourceList]struct{}{}
	originalTESs := c.SigningCertificate.TrustedEntityServices()

	if c.SigningCertificate.IsListOfTrustedEntitiesReached() {

		listsOfLists := map[*diagjaxb.XmlTrustSourceList]struct{}{}
		for _, t := range originalTESs {
			if t.ListOfTrustedSourceList != nil {
				listsOfLists[t.ListOfTrustedSourceList] = struct{}{}
			}
		}

		acceptableLoLoTEs := map[*diagjaxb.XmlTrustSourceList]struct{}{}
		for _, listOfLists := range sortedTrustSourceLists(listsOfLists) {
			loloteAnalysis := c.getLoTEAnalysis(listOfLists)
			if loloteAnalysis != nil {
				acceptableLOTL := c.isAcceptableLoLoTE(loloteAnalysis)
				item = item.SetNextItem(acceptableLOTL)
				if acceptableLOTL.Process() {
					acceptableLoLoTEs[listOfLists] = struct{}{}
				}
			}
		}

		// filter TLs with a found valid set of LOTLs (if assigned)
		lotes := map[*diagjaxb.XmlTrustSourceList]struct{}{}
		for _, t := range originalTESs {
			if t.TrustedSourceList == nil {
				continue
			}
			if t.ListOfTrustedSourceList != nil {
				if _, ok := acceptableLoLoTEs[t.ListOfTrustedSourceList]; !ok {
					continue
				}
			}
			lotes[t.TrustedSourceList] = struct{}{}
		}

		if len(lotes) > 0 {
			for _, lote := range sortedTrustSourceLists(lotes) {
				currentTL := c.getLoTEAnalysis(lote)
				if currentTL != nil {

					acceptableTL := c.isAcceptableLoTE(currentTL)
					item = item.SetNextItem(acceptableTL)

					item = item.SetNextItem(c.loteTypeKnown(lote.Type))

					if acceptableTL.Process() {
						acceptableLoTEs[lote] = struct{}{}
					}
				}
			}
		}
	}

	item = item.SetNextItem(c.isAcceptableLoTEPresent(acceptableLoTEs)) //nolint:staticcheck // mirrors upstream CertificateApprovalStatusBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

	if len(acceptableLoTEs) > 0 {

		listsByType := mapListsByType(acceptableLoTEs)

		for _, listTypeUri := range sortedStringKeysOfLists(listsByType) {
			lotes := listsByType[listTypeUri]

			trustedSourceUrls := getTrustedSourceUrls(lotes)
			filter := TrustedEntitiesFilterFactoryCreateFilterByListUrls(trustedSourceUrls)
			relatedServices := filter.Filter(originalTESs)

			applicableStiUris := getApplicableStiUris(relatedServices)

			for _, stiUri := range applicableStiUris {
				certApprovalStatusAtIssuanceBlock := c.getCertUsageAtIssuanceTimeBlock(listTypeUri, stiUri, relatedServices)
				c.Result.Value.ValidationCertificateApprovalStatus = append(c.Result.Value.ValidationCertificateApprovalStatus,
					certApprovalStatusAtIssuanceBlock.Execute())

				certApprovalStatusAtValidationTimeBlock := c.getCertUsageAtValidationTimeBlock(listTypeUri, stiUri, relatedServices)
				c.Result.Value.ValidationCertificateApprovalStatus = append(c.Result.Value.ValidationCertificateApprovalStatus,
					certApprovalStatusAtValidationTimeBlock.Execute())
			}
		}
	}
}

// getCertUsageAtIssuanceTimeBlock gets a certificate qualification
// determination process for validation at the certificate issuance time.
// Port of getCertUsageAtIssuanceTimeBlock(String, String, List).
func (c *CertificateApprovalStatusBlock) getCertUsageAtIssuanceTimeBlock(listTypeUri, stiUri string,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	return NewCertificateApprovalStatusAtTimeBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTime_CERTIFICATE_ISSUANCE_TIME,
		c.SigningCertificate, listTypeUri, stiUri, acceptableServices)
}

// getCertUsageAtValidationTimeBlock gets a certificate qualification
// determination process for validation at the validation time. Port of
// getCertUsageAtValidationTimeBlock(String, String, List).
func (c *CertificateApprovalStatusBlock) getCertUsageAtValidationTimeBlock(listTypeUri, stiUri string,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	return NewCertificateApprovalStatusAtTimeBlock(c.I18nProvider, enumerations.ValidationTime_VALIDATION_TIME, &c.ValidationTime,
		c.SigningCertificate, listTypeUri, stiUri, acceptableServices)
}

// getLoTEAnalysis ports the private getLoTEAnalysis(XmlTrustSourceList).
func (c *CertificateApprovalStatusBlock) getLoTEAnalysis(listSource *diagjaxb.XmlTrustSourceList) *jaxb.XmlLoTEAnalysis {
	if listSource.Url == nil {
		return nil
	}
	for _, xmlLoTEAnalysis := range c.LoteAnalysis {
		if *listSource.Url == xmlLoTEAnalysis.URL {
			return xmlLoTEAnalysis
		}
	}
	return nil
}

// mapListsByType ports the private mapListsByType(Collection). See
// InitChain's header for the flagged HashSet/HashMap iteration-order
// deviation: both the key order (consumed via sortedStringKeysOfLists) and
// the per-key slice order here are sorted by URL for determinism.
func mapListsByType(trustSourceLists map[*diagjaxb.XmlTrustSourceList]struct{}) map[string][]*diagjaxb.XmlTrustSourceList {
	result := map[string][]*diagjaxb.XmlTrustSourceList{}
	for _, trustSourceList := range sortedTrustSourceLists(trustSourceLists) {
		var listType string
		if trustSourceList.Type != nil {
			listType = *trustSourceList.Type
		}
		result[listType] = append(result[listType], trustSourceList)
	}
	return result
}

// sortedTrustSourceLists returns a *jaxb.XmlTrustSourceList set's members in
// a deterministic order (sorted by URL, empty/nil URL first), standing in
// for Java's HashSet<XmlTrustSourceList> (identity hashCode/equals)
// iteration order; see InitChain's header for the hash-order caveat.
func sortedTrustSourceLists(m map[*diagjaxb.XmlTrustSourceList]struct{}) []*diagjaxb.XmlTrustSourceList {
	lists := make([]*diagjaxb.XmlTrustSourceList, 0, len(m))
	for l := range m {
		lists = append(lists, l)
	}
	sort.Slice(lists, func(i, j int) bool {
		var ui, uj string
		if lists[i].Url != nil {
			ui = *lists[i].Url
		}
		if lists[j].Url != nil {
			uj = *lists[j].Url
		}
		return ui < uj
	})
	return lists
}

// sortedStringKeysOfLists returns a map's string keys in sorted order,
// standing in for Java's HashMap<String, List<...>> iteration order; see
// InitChain's header for the hash-order caveat.
func sortedStringKeysOfLists(m map[string][]*diagjaxb.XmlTrustSourceList) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// getApplicableStiUris ports the private getApplicableStiUris(List).
func getApplicableStiUris(services []*diagnostic.TrustedEntityServiceWrapper) []string {
	var result []string
	for _, s := range services {
		result = append(result, s.Type)
	}
	return result
}

// getTrustedSourceUrls ports the private getTrustedSourceUrls(Collection).
func getTrustedSourceUrls(lotes []*diagjaxb.XmlTrustSourceList) []string {
	var result []string
	for _, lote := range lotes {
		if lote.Url != nil {
			result = append(result, *lote.Url)
		}
	}
	return result
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *CertificateApprovalStatusBlock) AddAdditionalInfo() {
	c.setIndication()
}

// setIndication ports the private setIndication().
func (c *CertificateApprovalStatusBlock) setIndication() {
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

func (c *CertificateApprovalStatusBlock) isAcceptableLoLoTE(xmlLoLoTEAnalysis *jaxb.XmlLoTEAnalysis) *AcceptableLoLoTECheck[*jaxb.XmlCertificateApprovalStatusProcess] {
	return NewAcceptableLoLoTECheck(c.I18nProvider, c.Result, xmlLoLoTEAnalysis, c.WarnLevelRule())
}

func (c *CertificateApprovalStatusBlock) isAcceptableLoTE(xmlLoTEAnalysis *jaxb.XmlLoTEAnalysis) *AcceptableLoTECheck[*jaxb.XmlCertificateApprovalStatusProcess] {
	return NewAcceptableLoTECheck(c.I18nProvider, c.Result, xmlLoTEAnalysis, c.WarnLevelRule())
}

func (c *CertificateApprovalStatusBlock) loteTypeKnown(loteType *string) process.ChainItem[*jaxb.XmlCertificateApprovalStatusProcess] {
	var t string
	if loteType != nil {
		t = *loteType
	}
	return NewListTypeKnownCheck(c.I18nProvider, c.Result, t, c.WarnLevelRule())
}

func (c *CertificateApprovalStatusBlock) isAcceptableLoTEPresent(acceptableLoTEs map[*diagjaxb.XmlTrustSourceList]struct{}) process.ChainItem[*jaxb.XmlCertificateApprovalStatusProcess] {
	return NewAcceptableLoTEPresenceCheck(c.I18nProvider, c.Result, acceptableLoTEs, c.FailLevelRule())
}

func (c *CertificateApprovalStatusBlock) isAcceptableBuildingBlockConclusion(buildingBlocksConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlCertificateApprovalStatusProcess] {
	return NewAcceptableBuildingBlockConclusionCheck(c.I18nProvider, c.Result, buildingBlocksConclusion, c.WarnLevelRule())
}
