// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/eaa/EAAQualificationProcessBlock.java (DSS 6.5.RC1).
//
// Java iterates HashSet<String> listOfTrustedListUrls/trustedListUrls
// directly, feeding report Constraint order, as in
// CertificateQualificationBlock/SignatureQualificationBlock/TimestampQualificationBlock.
// This port reproduces that HashSet iteration order via orderedURLSet (see
// certificate_qualification_block.go), whose iterate() calls
// utils.JavaHashMapStringKeyOrder.
package qualification

import (
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// EAAQualificationProcessBlock performs validation of an EAA according to
// the EAA types defined in Regulation EU 2024/1183 (eIDAS 2.0).
type EAAQualificationProcessBlock struct {
	*process.ChainBase[*jaxb.XmlValidationEAAQualificationProcess]

	// eaa is the EAA to be validated.
	eaa *diagnostic.EAAWrapper

	// eaaConclusion is the conclusion of EAA validation.
	eaaConclusion *jaxb.XmlConclusion

	// signatureMap is a map of signature validation processes.
	signatureMap map[string]*jaxb.XmlSignature

	// tlAnalysis is the list of all TL analyses.
	tlAnalysis []*jaxb.XmlTLAnalysis

	// currentTime is the validation time.
	currentTime time.Time
}

// NewEAAQualificationProcessBlock is the default constructor. Port of
// EAAQualificationProcessBlock(Provider, EAAWrapper, XmlConclusion, Map, List, Date).
func NewEAAQualificationProcessBlock(i18nProvider *i18n.Provider, eaa *diagnostic.EAAWrapper,
	eaaConclusion *jaxb.XmlConclusion, signatureMap map[string]*jaxb.XmlSignature, tlAnalysis []*jaxb.XmlTLAnalysis,
	currentTime time.Time) *EAAQualificationProcessBlock {
	xmlResult := &jaxb.XmlValidationEAAQualificationProcess{}
	c := &EAAQualificationProcessBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		eaa:           eaa,
		eaaConclusion: eaaConclusion,
		signatureMap:  signatureMap,
		tlAnalysis:    tlAnalysis,
		currentTime:   currentTime,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *EAAQualificationProcessBlock) Title() i18n.MessageTag {
	return i18n.MessageTagEAAQualificationProcess
}

// InitChain initializes the chain. Port of initChain().
func (c *EAAQualificationProcessBlock) InitChain() {

	if utils.IsCollectionEmpty(c.eaa.EAASignatures()) {
		panic("No signatures found within the EAA token!")
	}

	signature := c.eaa.EAASignatures()[0]
	signingCertificate := signature.SigningCertificate()

	item := c.isTrustedListReachedForCertificateChain(signingCertificate)
	c.FirstItem = item

	item = item.SetNextItem(c.categoryPresent())

	claimedQualification := c.getClaimedQualification()

	if enumerations.EAAQualificationQEAA == claimedQualification {
		item = item.SetNextItem(c.categoryForQEAA())
	} else if enumerations.EAAQualificationPubEAA == claimedQualification {
		item = item.SetNextItem(c.categoryForPubEAA())
	}

	signatureQualification := enumerations.SignatureQualificationNA

	if signingCertificate != nil {

		if signingCertificate.IsTrustedListReached() {

			acceptableTLUrls := map[string]struct{}{}

			originalTSPs := signingCertificate.TrustServices()

			listOfTrustedListUrls := newOrderedURLSet()
			for _, t := range originalTSPs {
				if t.ListOfTrustedLists != nil && t.ListOfTrustedLists.Url != nil {
					listOfTrustedListUrls.add(*t.ListOfTrustedLists.Url)
				}
			}

			acceptableLOTLUrls := map[string]struct{}{}
			for _, lotlURL := range listOfTrustedListUrls.iterate() {
				lotlAnalysis := c.getTLAnalysis(lotlURL)
				if lotlAnalysis != nil {
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

			if trustedListUrls.len() > 0 {
				for _, tlURL := range trustedListUrls.iterate() {
					currentTL := c.getTLAnalysis(tlURL)
					if currentTL != nil {
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
				filteredServices := filter.Filter(originalTSPs)

				// Execute only for Trusted Lists with defined MRA
				if c.isMRAEnactedForTrustedList(filteredServices) {
					filter = TrustServicesFilterFactoryCreateMRAEnactedFilter()
					filteredServices = filter.Filter(filteredServices)

					filter = TrustServicesFilterFactoryCreateFilterByMRAEquivalenceStartingDate(&c.currentTime)
					filteredServices = filter.Filter(filteredServices)

					item = c.hasMraEnactedTrustService(filteredServices)
					c.FirstItem = item
				}

				// TODO : add filter for EAA and Pub-EAA ?

				if enumerations.EAAQualificationQEAA == claimedQualification {

					// 1. filter by service for EAA/Q
					filter = TrustServicesFilterFactoryCreateFilterByQEAA()
					filteredServices = filter.Filter(filteredServices)

					item = item.SetNextItem(c.hasQEAA(filteredServices))
				}

				// 2. filter by granted
				filter = TrustServicesFilterFactoryCreateFilterByGranted()
				filteredServices = filter.Filter(filteredServices)

				item = item.SetNextItem(c.hasGrantedStatus(filteredServices))

				// 3. filter by date (validation time)
				filter = TrustServicesFilterFactoryCreateFilterByDate(&c.currentTime)
				filteredServices = filter.Filter(filteredServices)

				item = item.SetNextItem(c.hasGrantedStatusAtValidationTime(filteredServices))

				if utils.IsCollectionEmpty(filteredServices) {
					claimedQualification = c.toNotQualifiedEAA(claimedQualification)
				}
			}
		}

		xmlSignature, ok := c.signatureMap[signature.Id()]
		if !ok || xmlSignature == nil {
			panic(fmt.Sprintf("Signature validation is not found for Id '%s'", signature.Id()))
		}
		validationSignatureQualification := xmlSignature.ValidationSignatureQualification
		if validationSignatureQualification == nil {
			panic(fmt.Sprintf("Signature qualification validation is not found for Id '%s'", signature.Id()))
		}

		signatureQualification = validationSignatureQualification.SignatureQualification.SignatureQualification()

		if enumerations.EAAQualificationQEAA == claimedQualification || enumerations.EAAQualificationPubEAA == claimedQualification {
			item = item.SetNextItem(c.isSignatureQualificationStatusAcceptable(signature, signatureQualification))
		}

		if enumerations.EAAQualificationPubEAA == claimedQualification {
			psbEaa := c.psbEaa(signingCertificate)
			item = item.SetNextItem(psbEaa) //nolint:staticcheck // mirrors upstream EAAQualificationProcessBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

			if !psbEaa.Process() {
				claimedQualification = c.toNotQualifiedEAA(claimedQualification)
			}
		}
	}

	c.determineFinalQualification(claimedQualification, signatureQualification)
}

func (c *EAAQualificationProcessBlock) isTrustedListReachedForCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewTrustedListReachedForCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) categoryPresent() process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewEAACategoryForEAAPresenceCheck(c.I18nProvider, c.Result, c.eaa, c.InfoLevelRule())
}

func (c *EAAQualificationProcessBlock) categoryForQEAA() process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewEAACategoryForQEAACheck(c.I18nProvider, c.Result, c.eaa, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) categoryForPubEAA() process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewEAACategoryForPubEAACheck(c.I18nProvider, c.Result, c.eaa, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) isSignatureQualificationStatusAcceptable(signature *diagnostic.SignatureWrapper,
	signatureQualification enumerations.SignatureQualification) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewEAAQualifiedSignatureOrSealCheck(c.I18nProvider, c.Result, signature, signatureQualification, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) isAcceptableLOTL(xmlLOTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableListOfTrustedListsCheck[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewAcceptableListOfTrustedListsCheck(c.I18nProvider, c.Result, xmlLOTLAnalysis, c.WarnLevelRule())
}

func (c *EAAQualificationProcessBlock) isAcceptableTL(xmlTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableTrustedListCheck[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewAcceptableTrustedListCheck(c.I18nProvider, c.Result, xmlTLAnalysis, c.WarnLevelRule())
}

func (c *EAAQualificationProcessBlock) isAcceptableTLPresent(acceptableUrls map[string]struct{}) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewAcceptableTrustedListPresenceCheck(c.I18nProvider, c.Result, acceptableUrls, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) hasMraEnactedTrustService(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewRelatedToMraEnactedTrustServiceCheck(c.I18nProvider, c.Result, services, c.WarnLevelRule())
}

func (c *EAAQualificationProcessBlock) hasQEAA(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewQEAACheck(c.I18nProvider, c.Result, services, c.WarnLevelRule())
}

func (c *EAAQualificationProcessBlock) hasGrantedStatus(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewGrantedStatusCheck(c.I18nProvider, c.Result, services, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) hasGrantedStatusAtValidationTime(services []*diagnostic.TrustServiceWrapper) process.ChainItem[*jaxb.XmlValidationEAAQualificationProcess] {
	return NewGrantedStatusAtTimeCheck(c.I18nProvider, c.Result, services, enumerations.ValidationTimeValidationTime, c.FailLevelRule())
}

func (c *EAAQualificationProcessBlock) psbEaa(certificateWrapper *diagnostic.CertificateWrapper) *EAAIssuerQcPSBPresentCheck {
	return NewEAAIssuerQcPSBPresentCheck(c.I18nProvider, c.Result, certificateWrapper, c.WarnLevelRule())
}

// isMRAEnactedForTrustedList ports the private isMRAEnactedForTrustedList(List).
func (c *EAAQualificationProcessBlock) isMRAEnactedForTrustedList(trustServices []*diagnostic.TrustServiceWrapper) bool {
	for _, trustService := range trustServices {
		if trustService.TrustedList != nil && utils.IsTrue(trustService.TrustedList.Mra) {
			return true
		}
	}
	return false
}

// getTLAnalysis ports the private getTLAnalysis(String).
func (c *EAAQualificationProcessBlock) getTLAnalysis(url string) *jaxb.XmlTLAnalysis {
	for _, xmlTLAnalysis := range c.tlAnalysis {
		if url == xmlTLAnalysis.URL {
			return xmlTLAnalysis
		}
	}
	return nil
}

// getClaimedQualification ports the private getClaimedQualification().
func (c *EAAQualificationProcessBlock) getClaimedQualification() enumerations.EAAQualification {
	eaaCategory := c.eaa.EAACategory()
	if enumerations.EAACategoryEUQEAA.URN() == eaaCategory {
		return enumerations.EAAQualificationQEAA
	} else if enumerations.EAACategoryEUPubEAA.URN() == eaaCategory {
		return enumerations.EAAQualificationPubEAA
	} else if eaaCategory == "" {
		// EAA-5.2.2.1-01: SD-JWT VC EAAs issued by EAAs issuers registered in the European Union,
		// which are neither SD-JWT VC QEAAs nor SD-JWT VC PuB-EAAs, shall not include the category claim.
		return enumerations.EAAQualificationEAA
	}
	return enumerations.EAAQualificationUnknown
}

// toNotQualifiedEAA ports the private toNotQualifiedEAA(EAAQualification).
func (c *EAAQualificationProcessBlock) toNotQualifiedEAA(qualification enumerations.EAAQualification) enumerations.EAAQualification {
	if enumerations.EAAQualificationQEAA == qualification || enumerations.EAAQualificationPubEAA == qualification {
		return enumerations.EAAQualificationEAA
	}
	return qualification
}

// determineFinalQualification ports the private
// determineFinalQualification(EAAQualification, SignatureQualification).
func (c *EAAQualificationProcessBlock) determineFinalQualification(claimedQualification enumerations.EAAQualification,
	signatureQualification enumerations.SignatureQualification) {
	finalQualification := EAAQualificationMatrixGetEAAQualification(c.eaaConclusion.Indication.Indication(), claimedQualification, signatureQualification)
	c.Result.Value.EAAQualification = jaxb.EAAQualificationValue(finalQualification)
}
