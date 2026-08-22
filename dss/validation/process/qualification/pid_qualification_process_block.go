// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/pid/PIDQualificationProcessBlock.java (DSS 6.5.RC1).
//
// FLAGGED HASH-ORDER SITE: as in CertificateApprovalStatusBlock, Java builds
// listsOfLists/lotes as HashSet<XmlTrustSourceList> (identity hashCode/
// equals), iterating them directly to append checks to the report - order
// observable in report output. This port reuses this file's
// sortedTrustSourceLists helper (see certificate_approval_status_block.go)
// for a deterministic (URL-sorted) surrogate order. See the porter brief's
// hard rule on hash-order leaks.
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	diagjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"

	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// PIDQualificationProcessBlock performs verification whether the provided
// token is a PID.
type PIDQualificationProcessBlock struct {
	*process.ChainBase[*jaxb.XmlValidationPIDQualificationProcess]

	// eaa is the EAA to be validated.
	eaa *diagnostic.EAAWrapper

	// eaaConclusion is the conclusion of EAA validation.
	eaaConclusion *jaxb.XmlConclusion

	// loteAnalysis is the list of List of Trusted Entities validations.
	loteAnalysis []*jaxb.XmlLoTEAnalysis

	// currentTime is the validation time.
	currentTime time.Time
}

// NewPIDQualificationProcessBlock is the default constructor. Port of
// PIDQualificationProcessBlock(I18nProvider, EAAWrapper, XmlConclusion, List, Date).
func NewPIDQualificationProcessBlock(i18nProvider *i18n.I18nProvider, eaa *diagnostic.EAAWrapper,
	eaaConclusion *jaxb.XmlConclusion, loteAnalysis []*jaxb.XmlLoTEAnalysis, currentTime time.Time) *PIDQualificationProcessBlock {
	xmlResult := &jaxb.XmlValidationPIDQualificationProcess{}
	c := &PIDQualificationProcessBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		eaa:           eaa,
		eaaConclusion: eaaConclusion,
		loteAnalysis:  loteAnalysis,
		currentTime:   currentTime,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *PIDQualificationProcessBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_PID_QUALIFICATION_PROCESS
}

// InitChain initializes the chain. Port of initChain().
func (c *PIDQualificationProcessBlock) InitChain() {

	if utils.IsCollectionEmpty(c.eaa.EAASignatures()) {
		panic("No signatures found within the EAA token!")
	}

	var certificateApprovalStatusAtIssuanceTime enumerations.CertificateApprovalStatus = enumerations.CertificateApprovalStatusEnum_NA
	var certificateApprovalStatusAtValidationTime enumerations.CertificateApprovalStatus = enumerations.CertificateApprovalStatusEnum_NA

	signature := c.eaa.EAASignatures()[0]
	signingCertificate := signature.SigningCertificate()

	item := c.isListOfTrustedEntitiesReachedForCertificateChain(signingCertificate)
	c.FirstItem = item

	pidDocumentTypeAcceptableCheck := c.pidDocumentTypeAcceptable()
	item = item.SetNextItem(pidDocumentTypeAcceptableCheck)

	if signingCertificate != nil {

		acceptableLoTEs := map[*diagjaxb.XmlTrustSourceList]struct{}{}

		if signingCertificate.IsListOfTrustedEntitiesReached() {

			originalTESs := signingCertificate.TrustedEntityServices()

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

			// filter TLs with a found valid set of LoLoTEs (if assigned)
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

			lotesSorted := sortedTrustSourceLists(lotes)
			if len(lotes) > 0 {
				for _, lote := range lotesSorted {
					currentTL := c.getLoTEAnalysis(lote)
					if currentTL != nil {

						acceptableTL := c.isAcceptableLoTE(currentTL)
						item = item.SetNextItem(acceptableTL)

						var loteType string
						if lote.Type != nil {
							loteType = *lote.Type
						}
						loteForPIDProviders := c.loteForPIDProviders(loteType)
						item = item.SetNextItem(loteForPIDProviders)

						if acceptableTL.Process() && loteForPIDProviders.Process() {
							acceptableLoTEs[lote] = struct{}{}
						}
					}
				}
			}

			item = item.SetNextItem(c.isAcceptableLoTEPresent(acceptableLoTEs))

			if len(acceptableLoTEs) > 0 {

				trustedSourceUrls := getTrustedSourceUrls(lotesSorted)
				filter := TrustedEntitiesFilterFactoryCreateFilterByListUrls(trustedSourceUrls)
				relatedServices := filter.Filter(originalTESs)

				filter = TrustedEntitiesFilterFactoryCreateFilterByServiceTypeIdentifierUri(enumerations.LoTEServiceTypeIdentifierEnum_PID_ISSUANCE.URI())
				filteredServices := filter.Filter(relatedServices)

				item = item.SetNextItem(c.stiForPIDIssuance(filteredServices))

				if utils.IsCollectionNotEmpty(filteredServices) {
					// assign not filtered list to ensure certificate qualification processing
					relatedServices = filteredServices
				}

				if utils.IsCollectionNotEmpty(relatedServices) {

					certApprovalStatusAtIssuanceBlock := c.getCertUsageAtIssuanceTimeBlock(signingCertificate, relatedServices)
					certApprovalStatusAtIssuanceResult := certApprovalStatusAtIssuanceBlock.Execute()
					c.Result.Value.ValidationCertificateApprovalStatus = append(c.Result.Value.ValidationCertificateApprovalStatus,
						certApprovalStatusAtIssuanceResult)

					certApprovalStatusAtValidationTimeBlock := c.getCertUsageAtValidationTimeBlock(signingCertificate, relatedServices)
					certApprovalStatusAtValidationTimeResult := certApprovalStatusAtValidationTimeBlock.Execute()
					c.Result.Value.ValidationCertificateApprovalStatus = append(c.Result.Value.ValidationCertificateApprovalStatus,
						certApprovalStatusAtValidationTimeResult)

					if pidDocumentTypeAcceptableCheck.Process() {

						certificateApprovalStatusAtIssuanceTime = c.getCertificateApprovalStatus(certApprovalStatusAtIssuanceResult)
						item = item.SetNextItem(c.pidProviderAtIssuanceTime(certificateApprovalStatusAtIssuanceTime))

						certificateApprovalStatusAtValidationTime = c.getCertificateApprovalStatus(certApprovalStatusAtValidationTimeResult)
						item = item.SetNextItem(c.pidProviderAtValidationTime(certificateApprovalStatusAtValidationTime)) //nolint:staticcheck // mirrors upstream PIDQualificationProcessBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
					}
				}
			}
		}
	}

	c.determineFinalQualification(certificateApprovalStatusAtIssuanceTime, certificateApprovalStatusAtValidationTime)
}

func (c *PIDQualificationProcessBlock) isListOfTrustedEntitiesReachedForCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewListOfTrustedEntitiesReachedForCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, c.FailLevelRule())
}

func (c *PIDQualificationProcessBlock) isAcceptableLoLoTE(xmlLoLoTEAnalysis *jaxb.XmlLoTEAnalysis) *AcceptableLoLoTECheck[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewAcceptableLoLoTECheck(c.I18nProvider, c.Result, xmlLoLoTEAnalysis, c.WarnLevelRule())
}

func (c *PIDQualificationProcessBlock) isAcceptableLoTE(xmlLoTEAnalysis *jaxb.XmlLoTEAnalysis) *AcceptableLoTECheck[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewAcceptableLoTECheck(c.I18nProvider, c.Result, xmlLoTEAnalysis, c.WarnLevelRule())
}

func (c *PIDQualificationProcessBlock) isAcceptableLoTEPresent(acceptableLoTEs map[*diagjaxb.XmlTrustSourceList]struct{}) process.ChainItem[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewAcceptableLoTEPresenceCheck(c.I18nProvider, c.Result, acceptableLoTEs, c.FailLevelRule())
}

func (c *PIDQualificationProcessBlock) loteForPIDProviders(loteType string) *PIDProviderListCheck {
	return NewPIDProviderListCheck(c.I18nProvider, c.Result, loteType, c.WarnLevelRule())
}

func (c *PIDQualificationProcessBlock) stiForPIDIssuance(relatedServices []*diagnostic.TrustedEntityServiceWrapper) process.ChainItem[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewPIDIssuanceTrustedEntityServicesCheck(c.I18nProvider, c.Result, relatedServices, c.FailLevelRule())
}

func (c *PIDQualificationProcessBlock) pidDocumentTypeAcceptable() *PIDDocumentTypeAcceptableCheck {
	return NewPIDDocumentTypeAcceptableCheck(c.I18nProvider, c.Result, c.eaa, c.FailLevelRule())
}

func (c *PIDQualificationProcessBlock) pidProviderAtIssuanceTime(certificateApprovalStatus enumerations.CertificateApprovalStatus) process.ChainItem[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewPIDProviderCertificateAtIssuanceTimeCheck(c.I18nProvider, c.Result, certificateApprovalStatus, c.FailLevelRule())
}

func (c *PIDQualificationProcessBlock) pidProviderAtValidationTime(certificateApprovalStatus enumerations.CertificateApprovalStatus) process.ChainItem[*jaxb.XmlValidationPIDQualificationProcess] {
	return NewPIDProviderCertificateAtValidationTimeCheck(c.I18nProvider, c.Result, certificateApprovalStatus, c.FailLevelRule())
}

// getLoTEAnalysis ports the private getLoTEAnalysis(XmlTrustSourceList).
func (c *PIDQualificationProcessBlock) getLoTEAnalysis(listSource *diagjaxb.XmlTrustSourceList) *jaxb.XmlLoTEAnalysis {
	if listSource.Url == nil {
		return nil
	}
	for _, xmlTLAnalysis := range c.loteAnalysis {
		if *listSource.Url == xmlTLAnalysis.URL {
			return xmlTLAnalysis
		}
	}
	return nil
}

// getCertUsageAtIssuanceTimeBlock gets a certificate qualification
// determination process for validation at the certificate issuance time.
// Port of the overridable protected getCertUsageAtIssuanceTimeBlock(CertificateWrapper, List).
func (c *PIDQualificationProcessBlock) getCertUsageAtIssuanceTimeBlock(certificate *diagnostic.CertificateWrapper,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	return NewCertificateApprovalStatusAtTimeBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTime_CERTIFICATE_ISSUANCE_TIME,
		certificate, enumerations.LoTETypeEnum_EUPIDProvidersList.URI(), enumerations.LoTEServiceTypeIdentifierEnum_PID_ISSUANCE.URI(), acceptableServices)
}

// getCertUsageAtValidationTimeBlock gets a certificate qualification
// determination process for validation at the validation time. Port of the
// overridable protected getCertUsageAtValidationTimeBlock(CertificateWrapper, List).
func (c *PIDQualificationProcessBlock) getCertUsageAtValidationTimeBlock(certificate *diagnostic.CertificateWrapper,
	acceptableServices []*diagnostic.TrustedEntityServiceWrapper) *CertificateApprovalStatusAtTimeBlock {
	return NewCertificateApprovalStatusAtTimeBlock(c.I18nProvider, enumerations.ValidationTime_VALIDATION_TIME, &c.currentTime,
		certificate, enumerations.LoTETypeEnum_EUPIDProvidersList.URI(), enumerations.LoTEServiceTypeIdentifierEnum_PID_ISSUANCE.URI(), acceptableServices)
}

// getCertificateApprovalStatus ports the private
// getCertificateApprovalStatus(XmlValidationCertificateApprovalStatus).
func (c *PIDQualificationProcessBlock) getCertificateApprovalStatus(
	xmlValidationCertificateApprovalStatus *jaxb.XmlValidationCertificateApprovalStatus) enumerations.CertificateApprovalStatus {
	xmlCertificateApprovalStatus := xmlValidationCertificateApprovalStatus.CertificateApprovalStatus
	return enumerations.CertificateApprovalStatusFromDefinition(xmlCertificateApprovalStatus.ListType,
		xmlCertificateApprovalStatus.ServiceTypeIdentifier, xmlCertificateApprovalStatus.ServiceStatus)
}

// determineFinalQualification ports the private
// determineFinalQualification(CertificateApprovalStatus, CertificateApprovalStatus).
func (c *PIDQualificationProcessBlock) determineFinalQualification(
	certificateApprovalStatusAtIssuanceTime, certificateApprovalStatusAtValidationTime enumerations.CertificateApprovalStatus) {
	certificateApprovalStatus := c.determinedFinalCertificateApprovalStatus(certificateApprovalStatusAtIssuanceTime, certificateApprovalStatusAtValidationTime)
	finalQualification := EAAQualificationMatrixGetPIDQualification(c.eaaConclusion.Indication.Indication(), certificateApprovalStatus)
	c.Result.Value.EAAQualification = jaxb.EAAQualificationValue(finalQualification)
}

// determinedFinalCertificateApprovalStatus ports the private
// determinedFinalCertificateApprovalStatus(CertificateApprovalStatus, CertificateApprovalStatus).
func (c *PIDQualificationProcessBlock) determinedFinalCertificateApprovalStatus(
	certificateApprovalStatusAtIssuanceTime, certificateApprovalStatusAtValidationTime enumerations.CertificateApprovalStatus) enumerations.CertificateApprovalStatus {
	if certificateApprovalStatusAtIssuanceTime == certificateApprovalStatusAtValidationTime {
		return certificateApprovalStatusAtIssuanceTime
	}
	return nil
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden protected void
// collectAdditionalMessages(XmlConclusion).
func (c *PIDQualificationProcessBlock) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	for _, certificateApprovalStatus := range c.Result.Value.ValidationCertificateApprovalStatus {
		c.CollectAllMessages(conclusion, certificateApprovalStatus.Conclusion)
	}
}
