// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/CertificateQualificationBlock.java (DSS 6.5.RC1).
//
// CROSS-CHUNK ASSUMPTION: AcceptableListOfTrustedListsCheck, AcceptableTrustedListCheck,
// AcceptableTrustedListPresenceCheck (Java package
// qualification.signature.checks) and TrustServiceFilter,
// TrustServicesFilterFactoryCreateFilterByUrls (Java package
// qualification.trust.filter) are owned by sibling porters of this shared
// package and were not present on disk while this file was written. Their
// call sites here follow this package's established naming convention
// (generic [T any] chain items constructed via NewXxx[T], "Factory" static
// methods flattened to TrustServicesFilterFactoryCreateXxx package
// functions) but must be reconciled against the sibling porters' actual
// signatures once available.
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateQualificationBlock verifies the final qualification of a
// certificate, processing its validation at issuance and validation time.
type CertificateQualificationBlock struct {
	*process.ChainBase[*jaxb.XmlCertificateQualificationProcess]

	// BuildingBlocksConclusion is the certificate's BasicBuildingBlock's
	// conclusion.
	BuildingBlocksConclusion *jaxb.XmlConclusion

	// ValidationTime is the validation time.
	ValidationTime time.Time

	// SigningCertificate is the certificate to determine qualification for.
	SigningCertificate *diagnostic.CertificateWrapper

	// TlAnalysis is a list of validation results for all Trusted Lists.
	TlAnalysis []*jaxb.XmlTLAnalysis

	overrides CertificateQualificationBlockOverrides
}

// CertificateQualificationBlockOverrides captures the members Java's
// CertificateQualificationBlock treats virtually: getCertQualificationAtIssuanceTimeBlock()
// and getCertQualificationAtValidationTimeBlock(), self-called from
// initChain(). CertificateQualificationForQWACBlock overrides them and
// registers itself via InitCertificateQualificationBlock.
type CertificateQualificationBlockOverrides interface {
	// CertQualificationAtIssuanceTimeBlock gets a certificate qualification
	// determination process for validation at the certificate issuance
	// time. Port of the overridable protected
	// getCertQualificationAtIssuanceTimeBlock(List).
	CertQualificationAtIssuanceTimeBlock(acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock
	// CertQualificationAtValidationTimeBlock gets a certificate
	// qualification determination process for validation at the validation
	// time. Port of the overridable protected
	// getCertQualificationAtValidationTimeBlock(List).
	CertQualificationAtValidationTimeBlock(acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock
}

// NewCertificateQualificationBlock is the default constructor. Port of
// CertificateQualificationBlock(I18nProvider, XmlConclusion, Date, CertificateWrapper, List).
func NewCertificateQualificationBlock(i18nProvider *i18n.I18nProvider, buildingBlocksConclusion *jaxb.XmlConclusion,
	validationTime time.Time, signingCertificate *diagnostic.CertificateWrapper,
	tlAnalysis []*jaxb.XmlTLAnalysis) *CertificateQualificationBlock {
	xmlResult := &jaxb.XmlCertificateQualificationProcess{}
	c := &CertificateQualificationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		BuildingBlocksConclusion: buildingBlocksConclusion,
		ValidationTime:           validationTime,
		SigningCertificate:       signingCertificate,
		TlAnalysis:               tlAnalysis,
	}

	id := signingCertificate.Id()
	c.Result.Value.Id = &id

	c.InitCertificateQualificationBlock(c)
	return c
}

// InitCertificateQualificationBlock registers the concrete chain with its
// base so that the base can dispatch to the overridden methods. It must be
// called exactly once, by the concrete chain's constructor, before Execute.
func (c *CertificateQualificationBlock) InitCertificateQualificationBlock(overrides CertificateQualificationBlockOverrides) {
	c.overrides = overrides
	c.InitChainBase(c)
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *CertificateQualificationBlock) Title() i18n.MessageTag {
	return i18n.MessageTagCertQualification
}

// InitChain initializes the chain. Port of initChain().
//
// HASH-ORDER (closed in phase 8f): Java builds listOfTrustedListUrls/trustedListUrls
// as HashSet<String> and iterates them directly to append AcceptableListOfTrustedListsCheck/
// AcceptableTrustedListCheck constraints to the report, in HashSet bucket
// order. That order feeds the ordering of report Constraint elements, which
// is observable output. This port iterates the equivalent Go sets
// (map[string]struct{}, whose range order is randomized per the language
// spec - not just non-Java, but non-repeatable within this port too, unless
// sorted explicitly) in sorted URL order instead: deterministic across runs,
// but not necessarily identical to the upstream HashSet bucket order. See
// the porter brief's hard rule on hash-order leaks.
func (c *CertificateQualificationBlock) InitChain() {
	// cover incomplete cert chain / expired/ revoked certs
	item := c.isAcceptableBuildingBlockConclusion(c.BuildingBlocksConclusion)
	c.FirstItem = item

	acceptableTLUrls := map[string]struct{}{}
	originalTSPs := c.SigningCertificate.TrustServices()

	if c.SigningCertificate.IsTrustedListReached() {

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
				currentTL := c.getTlAnalysis(tlURL)
				if currentTL != nil {
					acceptableTL := c.isAcceptableTL(currentTL)
					item = item.SetNextItem(acceptableTL)
					if acceptableTL.Process() {
						acceptableTLUrls[tlURL] = struct{}{}
					}
				}
			}
		}
	}

	item = item.SetNextItem(c.isAcceptableTLPresent(acceptableTLUrls)) //nolint:staticcheck // mirrors upstream CertificateQualificationBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

	if len(acceptableTLUrls) > 0 {

		// 1. filter by service for CAQC
		filter := TrustServicesFilterFactoryCreateFilterByUrls(acceptableTLUrls)
		acceptableServices := filter.Filter(originalTSPs)

		certQualAtIssuanceBlock := c.overrides.CertQualificationAtIssuanceTimeBlock(acceptableServices)
		c.Result.Value.ValidationCertificateQualification = append(c.Result.Value.ValidationCertificateQualification,
			certQualAtIssuanceBlock.Execute())

		certQualAtValidationTimeBlock := c.overrides.CertQualificationAtValidationTimeBlock(acceptableServices)
		c.Result.Value.ValidationCertificateQualification = append(c.Result.Value.ValidationCertificateQualification,
			certQualAtValidationTimeBlock.Execute())
	}
}

// CertQualificationAtIssuanceTimeBlock gets a certificate qualification
// determination process for validation at the certificate issuance time.
// Port of getCertQualificationAtIssuanceTimeBlock(List).
func (c *CertificateQualificationBlock) CertQualificationAtIssuanceTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTimeCertificateIssuanceTime,
		c.SigningCertificate, acceptableServices)
}

// CertQualificationAtValidationTimeBlock gets a certificate qualification
// determination process for validation at the validation time. Port of
// getCertQualificationAtValidationTimeBlock(List).
func (c *CertificateQualificationBlock) CertQualificationAtValidationTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeBlock(c.I18nProvider, enumerations.ValidationTimeValidationTime, &c.ValidationTime,
		c.SigningCertificate, acceptableServices)
}

// getTlAnalysis ports the private getTlAnalysis(String).
func (c *CertificateQualificationBlock) getTlAnalysis(url string) *jaxb.XmlTLAnalysis {
	for _, xmlTLAnalysis := range c.TlAnalysis {
		if url == xmlTLAnalysis.URL {
			return xmlTLAnalysis
		}
	}
	return nil
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *CertificateQualificationBlock) AddAdditionalInfo() {
	c.setIndication()
}

// setIndication ports the private setIndication().
func (c *CertificateQualificationBlock) setIndication() {
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

func (c *CertificateQualificationBlock) isAcceptableLOTL(xmlLOTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableListOfTrustedListsCheck[*jaxb.XmlCertificateQualificationProcess] {
	return NewAcceptableListOfTrustedListsCheck(c.I18nProvider, c.Result, xmlLOTLAnalysis, c.WarnLevelRule())
}

func (c *CertificateQualificationBlock) isAcceptableTL(xmlTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableTrustedListCheck[*jaxb.XmlCertificateQualificationProcess] {
	return NewAcceptableTrustedListCheck(c.I18nProvider, c.Result, xmlTLAnalysis, c.WarnLevelRule())
}

func (c *CertificateQualificationBlock) isAcceptableTLPresent(acceptableUrls map[string]struct{}) process.ChainItem[*jaxb.XmlCertificateQualificationProcess] {
	return NewAcceptableTrustedListPresenceCheck(c.I18nProvider, c.Result, acceptableUrls, c.FailLevelRule())
}

func (c *CertificateQualificationBlock) isAcceptableBuildingBlockConclusion(buildingBlocksConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlCertificateQualificationProcess] {
	return NewAcceptableBuildingBlockConclusionCheck(c.I18nProvider, c.Result, buildingBlocksConclusion, c.WarnLevelRule())
}

// orderedURLSet collects trusted-list / list-of-trusted-lists URLs the way
// Java's Collectors.toSet() does - into a java.util.HashSet<String> - and
// records the insertion order so iterate() can reproduce that HashSet's real
// iteration order. Upstream appends one report <Constraint> per element while
// iterating these sets, so the order is OBSERVABLE, byte-compared output; a
// sorted-by-URL substitute (this port's earlier, explicitly flagged reading)
// is deterministic but not upstream's. Found by the phase-8f full-corpus
// report byte-parity run on DSS-2049/dss2049-doubleTL.xml.
type orderedURLSet struct {
	seen  map[string]struct{}
	order []string
}

func newOrderedURLSet() *orderedURLSet {
	return &orderedURLSet{seen: map[string]struct{}{}}
}

func (s *orderedURLSet) add(url string) {
	if _, exists := s.seen[url]; exists {
		return
	}
	s.seen[url] = struct{}{}
	s.order = append(s.order, url)
}

func (s *orderedURLSet) len() int { return len(s.order) }

// iterate answers the URLs in java.util.HashSet iteration order.
func (s *orderedURLSet) iterate() []string {
	return utils.JavaHashMapStringKeyOrder(s.order)
}
