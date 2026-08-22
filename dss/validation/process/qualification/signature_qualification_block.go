// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/signature/SignatureQualificationBlock.java (DSS 6.5.RC1).
//
// HASH-ORDER (closed in phase 8f): as in CertificateQualificationBlock, Java iterates
// HashSet<String> listOfTrustedListUrls/trustedListUrls directly, feeding
// report Constraint order. This port sorts the equivalent Go
// orderedURLSet, whose iterate() reproduces java.util.HashSet's own
// necessarily Java-bucket-identical) result. See the porter brief's hard rule.
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

// SignatureQualificationBlock performs the qualification verification for a
// signature.
type SignatureQualificationBlock struct {
	*process.ChainBase[*jaxb.XmlValidationSignatureQualification]

	// Etsi319102Conclusion is the conclusion of signature validation as in
	// EN 319 102-1.
	Etsi319102Conclusion *jaxb.XmlConclusion

	// BestSignatureTime is the best-signature-time.
	BestSignatureTime time.Time

	// SigningCertificate is the signing certificate.
	SigningCertificate *diagnostic.CertificateWrapper

	// TlAnalysis is the analyses of all available LOTL/TLs.
	TlAnalysis []*jaxb.XmlTLAnalysis

	// relatedTLAnalyses is the list of related LOTL/TL analyses.
	relatedTLAnalyses []*jaxb.XmlTLAnalysis

	// qualificationAtIssuanceTime is the determined signing certificate
	// qualification at its issuance time.
	qualificationAtIssuanceTime enumerations.CertificateQualification

	// qualificationAtSigningTime is the determined signing certificate
	// qualification at best-signature-time.
	qualificationAtSigningTime enumerations.CertificateQualification

	overrides SignatureQualificationBlockOverrides
}

// SignatureQualificationBlockOverrides captures the members Java's
// SignatureQualificationBlock treats virtually:
// getCertQualificationAtIssuanceTimeBlock() and
// getCertQualificationAtSigningTimeBlock(), self-called from initChain().
// TLSBindingSignatureQualificationBlock overrides them and registers itself
// via InitSignatureQualificationBlock.
type SignatureQualificationBlockOverrides interface {
	// CertQualificationAtIssuanceTimeBlock gets a certificate qualification
	// determination process for validation at the certificate issuance
	// time. Port of the overridable protected
	// getCertQualificationAtIssuanceTimeBlock(List).
	CertQualificationAtIssuanceTimeBlock(acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock
	// CertQualificationAtSigningTimeBlock gets a certificate qualification
	// determination process for validation at the certificate signing
	// time. Port of the overridable protected
	// getCertQualificationAtSigningTimeBlock(List, Date).
	CertQualificationAtSigningTimeBlock(acceptableServices []*diagnostic.TrustServiceWrapper, signingTime time.Time) *CertQualificationAtTimeBlock
}

// NewSignatureQualificationBlock is the default constructor. Port of
// SignatureQualificationBlock(I18nProvider, XmlConstraintsConclusionWithProofOfExistence, CertificateWrapper, List).
func NewSignatureQualificationBlock(i18nProvider *i18n.I18nProvider,
	etsi319102validation *jaxb.XmlConstraintsConclusionWithProofOfExistence,
	signingCertificate *diagnostic.CertificateWrapper, tlAnalysis []*jaxb.XmlTLAnalysis) *SignatureQualificationBlock {
	xmlResult := &jaxb.XmlValidationSignatureQualification{}
	c := &SignatureQualificationBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		Etsi319102Conclusion: etsi319102validation.Conclusion,
		SigningCertificate:   signingCertificate,
		TlAnalysis:           tlAnalysis,
	}
	if etsi319102validation.ProofOfExistence != nil {
		c.BestSignatureTime = time.Time(etsi319102validation.ProofOfExistence.Time)
	}

	c.InitSignatureQualificationBlock(c)
	return c
}

// InitSignatureQualificationBlock registers the concrete chain with its
// base so that the base can dispatch to the overridden methods. It must be
// called exactly once, by the concrete chain's constructor, before Execute.
func (c *SignatureQualificationBlock) InitSignatureQualificationBlock(overrides SignatureQualificationBlockOverrides) {
	c.overrides = overrides
	c.InitChainBase(c)
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *SignatureQualificationBlock) Title() i18n.MessageTag {
	return i18n.MessageTag_SIG_QUALIFICATION
}

// InitChain initializes the chain. Port of initChain().
func (c *SignatureQualificationBlock) InitChain() {
	item := c.isAdES(c.Etsi319102Conclusion)
	c.FirstItem = item

	item = item.SetNextItem(c.isTrustedListReachedForCertificateChain(c.SigningCertificate))

	acceptableTLUrls := map[string]struct{}{}

	if c.SigningCertificate != nil && c.SigningCertificate.IsTrustedListReached() {

		originalTSPs := c.SigningCertificate.TrustServices()

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

		item = item.SetNextItem(c.isAcceptableTLPresent(acceptableTLUrls))

		if len(acceptableTLUrls) > 0 {

			// 1. filter by service for CAQC
			filter := TrustServicesFilterFactoryCreateFilterByUrls(acceptableTLUrls)
			acceptableServices := filter.Filter(originalTSPs)

			certQualAtIssuanceBlock := c.overrides.CertQualificationAtIssuanceTimeBlock(acceptableServices)
			certQualAtIssuanceResult := certQualAtIssuanceBlock.Execute()
			c.Result.Value.ValidationCertificateQualification = append(c.Result.Value.ValidationCertificateQualification,
				certQualAtIssuanceResult)
			if certQualAtIssuanceResult.CertificateQualification != nil {
				c.qualificationAtIssuanceTime = certQualAtIssuanceResult.CertificateQualification.CertificateQualification()
			}

			certQualAtSigningTimeBlock := c.overrides.CertQualificationAtSigningTimeBlock(acceptableServices, c.BestSignatureTime)
			certQualAtSigningTimeResult := certQualAtSigningTimeBlock.Execute()
			c.Result.Value.ValidationCertificateQualification = append(c.Result.Value.ValidationCertificateQualification,
				certQualAtSigningTimeResult)
			if certQualAtSigningTimeResult.CertificateQualification != nil {
				c.qualificationAtSigningTime = certQualAtSigningTimeResult.CertificateQualification.CertificateQualification()
			}

			// execute signature qualification determination only when at least one valid trust service is found
			if utils.IsCollectionNotEmpty(certQualAtIssuanceBlock.FilteredServices()) ||
				utils.IsCollectionNotEmpty(certQualAtSigningTimeBlock.FilteredServices()) {
				// Article 32 :
				// (a) the certificate that supports the signature was, at the time of signing, a qualified certificate for
				// electronic signature complying with Annex I;
				item = item.SetNextItem(c.qualifiedCertificateAtSigningTime(c.qualificationAtSigningTime))

				// NOTE: Article 40:
				// Articles 32, 33 and 34 shall apply mutatis mutandis to the validation and preservation of
				// qualified electronic seals.
				item = item.SetNextItem(c.certificateTypeAtSigningTime(c.qualificationAtSigningTime))

				// (b) the qualified certificate
				// 1. was issued by a qualified trust service provider
				item = item.SetNextItem(c.qualifiedCertificateAtIssuance(c.qualificationAtIssuanceTime))

				// 2. was valid at the time of signing;
				// covered in isAdES

				// (c) the signature validation data corresponds to the data provided to the relying party;
				// covered in isAdES

				// (d) the unique set of data representing the signatory in the certificate is correctly provided to the
				// relying party;
				// covered in isAdES

				// (e) the use of any pseudonym is clearly indicated to the relying party if a pseudonym was used at the
				// time of signing;
				// covered in isAdES

				// (f) the electronic signature was created by a qualified electronic signature creation device;
				item = item.SetNextItem(c.qscdAtSigningTime(c.qualificationAtSigningTime)) //nolint:staticcheck // mirrors upstream SignatureQualificationBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.

				// (g) the integrity of the signed data has not been compromised;
				// covered in isAdES
			}
		}
	}
}

// CertQualificationAtIssuanceTimeBlock gets a certificate qualification
// determination process for validation at the certificate issuance time.
// Port of getCertQualificationAtIssuanceTimeBlock(List).
func (c *SignatureQualificationBlock) CertQualificationAtIssuanceTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeBlockAtIssuanceTime(c.I18nProvider, enumerations.ValidationTimeCertificateIssuanceTime,
		c.SigningCertificate, acceptableServices)
}

// CertQualificationAtSigningTimeBlock gets a certificate qualification
// determination process for validation at the certificate signing time.
// Port of getCertQualificationAtSigningTimeBlock(List, Date).
func (c *SignatureQualificationBlock) CertQualificationAtSigningTimeBlock(
	acceptableServices []*diagnostic.TrustServiceWrapper, signingTime time.Time) *CertQualificationAtTimeBlock {
	return NewCertQualificationAtTimeBlock(c.I18nProvider, enumerations.ValidationTimeBESTSignatureTime, &signingTime,
		c.SigningCertificate, acceptableServices)
}

// getTlAnalysis ports the private getTlAnalysis(String).
func (c *SignatureQualificationBlock) getTlAnalysis(url string) *jaxb.XmlTLAnalysis {
	for _, xmlTLAnalysis := range c.TlAnalysis {
		if url == xmlTLAnalysis.URL {
			return xmlTLAnalysis
		}
	}
	return nil
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo().
func (c *SignatureQualificationBlock) AddAdditionalInfo() {
	c.setIndication()
	c.determineFinalQualification()
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of the overridden protected void
// collectAdditionalMessages(XmlConclusion).
func (c *SignatureQualificationBlock) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	for _, certQualAtTime := range c.Result.Value.ValidationCertificateQualification {
		c.CollectAllMessages(conclusion, certQualAtTime.Conclusion)
	}
	for _, relatedTLAnalysis := range c.relatedTLAnalyses {
		c.CollectAllMessages(conclusion, relatedTLAnalysis.Conclusion)
	}
}

// determineFinalQualification ports the private determineFinalQualification().
func (c *SignatureQualificationBlock) determineFinalQualification() {
	sigQualif := enumerations.SignatureQualificationNA

	if c.Etsi319102Conclusion != nil && c.qualificationAtIssuanceTime != "" && c.qualificationAtSigningTime != "" {
		certificateQualificationCalculator := NewFinalCertificateQualificationCalculator(
			c.qualificationAtIssuanceTime, c.qualificationAtSigningTime)
		finalCertQualification := certificateQualificationCalculator.FinalQualification()
		sigQualif = SigQualificationMatrixGetSignatureQualification(c.Etsi319102Conclusion.Indication.Indication(), finalCertQualification)
	}

	c.Result.Value.SignatureQualification = jaxb.SignatureQualificationValue(sigQualif)
}

// setIndication ports the private setIndication().
func (c *SignatureQualificationBlock) setIndication() {
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

func (c *SignatureQualificationBlock) isTrustedListReachedForCertificateChain(
	signingCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewTrustedListReachedForCertificateChainCheck(c.I18nProvider, c.Result, signingCertificate, c.FailLevelRule())
}

func (c *SignatureQualificationBlock) isAcceptableLOTL(xmlLOTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableListOfTrustedListsCheck[*jaxb.XmlValidationSignatureQualification] {
	return NewAcceptableListOfTrustedListsCheck(c.I18nProvider, c.Result, xmlLOTLAnalysis, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) isAcceptableTL(xmlTLAnalysis *jaxb.XmlTLAnalysis) *AcceptableTrustedListCheck[*jaxb.XmlValidationSignatureQualification] {
	return NewAcceptableTrustedListCheck(c.I18nProvider, c.Result, xmlTLAnalysis, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) isAcceptableTLPresent(acceptableUrls map[string]struct{}) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewAcceptableTrustedListPresenceCheck(c.I18nProvider, c.Result, acceptableUrls, c.FailLevelRule())
}

func (c *SignatureQualificationBlock) isAdES(etsi319102Conclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewAdESAcceptableCheck(c.I18nProvider, c.Result, etsi319102Conclusion, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) qualifiedCertificateAtSigningTime(qualificationAtSigningTime enumerations.CertificateQualification) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewQualifiedCertificateAtSigningTimeCheck(c.I18nProvider, c.Result, qualificationAtSigningTime, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) certificateTypeAtSigningTime(qualificationAtSigningTime enumerations.CertificateQualification) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewCertificateTypeAtSigningTimeCheck(c.I18nProvider, c.Result, qualificationAtSigningTime, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) qualifiedCertificateAtIssuance(qualificationAtIssuance enumerations.CertificateQualification) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewQualifiedCertificateAtCertificateIssuanceCheck(c.I18nProvider, c.Result, qualificationAtIssuance, c.WarnLevelRule())
}

func (c *SignatureQualificationBlock) qscdAtSigningTime(qualificationAtSigningTime enumerations.CertificateQualification) process.ChainItem[*jaxb.XmlValidationSignatureQualification] {
	return NewQSCDCertificateAtSigningTimeCheck(c.I18nProvider, c.Result, qualificationAtSigningTime, c.WarnLevelRule())
}
