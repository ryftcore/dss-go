// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/X509CertificateValidation.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.validation.process.bbb.xcv and its checks, crs, rac,
// rac.checks, rfc, rfc.checks, sub and sub.checks subpackages flatten into this
// single Go package xcv (collision-checked), so every check is referenced
// unqualified. See chain.go in validation/process for the Result stand-in that
// replaces Java's "T extends XmlConstraintsConclusion" bound and for the
// overrides-registration pattern the chains follow.
//
// Date mapping. Java's validation times are java.util.Date; the ported
// building blocks take a plain time.Time wherever the value cannot be null
// (currentTime), and a *time.Time wherever Java may hand in null (usageTime,
// which BasicBuildingBlocks fills from CertificateWrapper#getNotBefore(),
// TimestampWrapper#getProductionTime() or RevocationWrapper#getProductionDate()).
// The one Java value that is nominally nullable but cannot be null for a
// schema-valid dump is CertificateWrapper#getNotBefore() as read for lastDate
// below: DiagnosticData.xsd declares Certificate/NotBefore a required element,
// and Java would NPE on a null one further down in
// ValidationProcessUtils#isTrustAnchor. The Go port passes the zero time there.
//
// slf4j logging is dropped per PORTING.md.
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// X509CertificateValidation is 5.2.6 X.509 certificate validation.
//
// This building block validates the signing certificate at current time.
type X509CertificateValidation struct {
	*process.ChainBase[*jaxb.XmlXCV]

	// currentCertificate is the certificate to be validated.
	currentCertificate *diagnostic.CertificateWrapper

	// currentTime is the validation time.
	currentTime time.Time

	// usageTime is the certificate approval status time; nil is Java's null.
	usageTime *time.Time

	// context is the validation context.
	context enumerations.Context

	// aov is the result of cryptographic algorithms validation.
	aov *jaxb.XmlAOV

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy
}

// NewX509CertificateValidation is the default constructor with usage time. Port
// of X509CertificateValidation(I18nProvider, CertificateWrapper, Date, Date,
// Context, XmlAOV, ValidationPolicy).
func NewX509CertificateValidation(i18nProvider *i18n.I18nProvider,
	currentCertificate *diagnostic.CertificateWrapper, currentTime time.Time, usageTime *time.Time,
	context enumerations.Context, aov *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *X509CertificateValidation {
	xmlXCV := &jaxb.XmlXCV{}
	c := &X509CertificateValidation{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlXCV,
			&xmlXCV.XmlConstraintsConclusionContent, &xmlXCV.XmlConstraintsConclusionAttrs)),
		currentCertificate: currentCertificate,
		currentTime:        currentTime,
		usageTime:          usageTime,
		context:            context,
		aov:                aov,
		validationPolicy:   validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *X509CertificateValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_X509_CERTIFICATE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *X509CertificateValidation) InitChain() {

	/*
	 * 1) If the signing certificate represents a trust anchor, then:
	 *
	 * a) If, in the X.509 Validation Constraints, a sunset date is associated to that trust anchor,
	 *    the building block shall check whether validation is before the sunset date. If validation time is
	 *    at or after the sunset date, the building block shall set the current status to
	 *    INDETERMINATE/NO_CERTIFICATE_CHAIN_FOUND_NO_POE and shall go to step 2).
	 * b) Else, the building block may, based on signature policy or local configuration, return with
	 *    the indication PASSED. Otherwise, the building block shall go to the next step.
	 *
	 * 2) The building block shall build a new prospective certificate chain that has not yet been evaluated.
	 * If the "Other Certificates" parameter is present, only certificates contained in that set of certificates
	 * may be used to build the chain. The chain shall satisfy the conditions of a prospective certificate chain:
	 *
	 * a) If no new chain can be built, the building block shall return the current status, the last chain built
	 *    and any additional information saved in step 4-a) or, if no chain has been built, the indication
	 *    INDETERMINATE with the sub-indication NO_CERTIFICATE_CHAIN_FOUND.
	 * b) Otherwise, the building block shall add this chain to the set of prospected chains and shall go to step 3).
	 *
	 * 3) If, in the X.509 Validation Constraints, a sunset date is associated to the trust anchor from which
	 * the current chain has been built, the building block shall check whether validation is before
	 * the sunset date. If validation time is at or after the sunset date, the building block shall set
	 * the current status to INDETERMINATE/NO_CERTIFICATE_CHAIN_FOUND_NO_POE and shall go to step 2).
	 */
	certificateChain := c.currentCertificate.CertificateChain()

	item := c.prospectiveCertificateChain(c.currentCertificate)
	c.FirstItem = item

	subContext := enumerations.SubContext_SIGNING_CERT
	// Java holds an Iterator over the certificate chain, or null when the chain
	// is empty; the Go port walks the same slice by index, so
	// "certChainIt != null && certChainIt.hasNext()" is the bounds test below -
	// an empty chain never enters the branch either way.
	certChainIt := 0

	trustAnchorCandidate := c.currentCertificate
	var trustAnchor *diagnostic.CertificateWrapper

	for {
		if trustAnchorCandidate.IsTrusted() && trustAnchorCandidate.TrustStartDate() != nil ||
			trustAnchorCandidate.TrustSunsetDate() != nil {

			item = item.SetNextItem(c.validationBeforeSunsetDate(trustAnchorCandidate, subContext, c.currentTime))

		}

		if c.isTrustAnchorReached(trustAnchorCandidate, subContext) {

			item = item.SetNextItem(c.prospectiveCertificateChainValidAtValidationTime(trustAnchorCandidate, subContext, c.currentTime))

			trustAnchor = trustAnchorCandidate
			break
		}

		if certChainIt < len(certificateChain) {
			trustAnchorCandidate = certificateChain[certChainIt]
			certChainIt++
		} else {
			trustAnchorCandidate = nil
		}
		subContext = enumerations.SubContext_CA_CERTIFICATE

		if trustAnchorCandidate == nil {
			break
		}
	}

	/*
	 * 4) The building block shall perform validation of the prospective certificate chain with the following inputs:
	 * the prospective chain built in the previous step, the trust anchor used in the previous step, the X.509 parameters
	 * provided in the inputs and the validation time. The validation shall be following the PKIX Certification Path
	 * Validation of IETF RFC 5280 [1], clause 6.1 with the exception of the validity model and the verification of
	 * whether the validation time is during the validity period of the signing certificate.
	 */
	if c.currentCertificate.IsTrusted() || c.currentCertificate.IsTrustedChain() || !c.prospectiveCertificateChainCheckEnforced() {

		item = item.SetNextItem(c.trustServiceWithExpectedTypeIdentifier(c.currentCertificate))

		item = item.SetNextItem(c.trustServiceWithExpectedStatus(c.currentCertificate))

		certificateValidation := NewSubX509CertificateValidation(c.I18nProvider,
			c.currentCertificate, c.currentTime, c.currentTime, c.context, enumerations.SubContext_SIGNING_CERT, c.aov, c.validationPolicy)
		subXCV := certificateValidation.Execute()
		c.Result.Value.SubXCV = append(c.Result.Value.SubXCV, subXCV)

		if c.currentCertificate.IsTrusted() {
			item = item.SetNextItem(c.checkTrustAnchorSubXCVResult(subXCV))
		} else {
			item = item.SetNextItem(c.checkSubXCVResult(subXCV))
		}

		if trustAnchor != nil && trustAnchor == c.currentCertificate {
			return
		}

		model := c.validationPolicy.ValidationModel()

		// Check CA_CERTIFICATEs
		var lastDate time.Time
		if enumerations.ValidationModel_SHELL == model {
			lastDate = c.currentTime
		} else if notBefore := c.currentCertificate.NotBefore(); notBefore != nil {
			lastDate = *notBefore
		}
		if utils.IsCollectionNotEmpty(certificateChain) {
			for _, certificate := range certificateChain {
				certificateValidation = NewSubX509CertificateValidation(c.I18nProvider,
					certificate, lastDate, c.currentTime, c.context, enumerations.SubContext_CA_CERTIFICATE, c.aov, c.validationPolicy)
				subXCV = certificateValidation.Execute()
				c.Result.Value.SubXCV = append(c.Result.Value.SubXCV, subXCV)

				if certificate.IsTrusted() {
					item = item.SetNextItem(c.checkTrustAnchorSubXCVResult(subXCV))
				} else {
					item = item.SetNextItem(c.checkSubXCVResult(subXCV))
				}

				if enumerations.ValidationModel_CHAIN == model {
					if notBefore := certificate.NotBefore(); notBefore != nil {
						lastDate = *notBefore
					} else {
						lastDate = time.Time{}
					}
				}
				// keep same time for SHELL and HYBRID

				if trustAnchor != nil && trustAnchor == certificate {
					return
				}
			}
		}

	}
}

// prospectiveCertificateChain ports the private
// prospectiveCertificateChain(CertificateWrapper).
func (c *X509CertificateValidation) prospectiveCertificateChain(
	currentCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlXCV] {
	constraint := c.validationPolicy.ProspectiveCertificateChainConstraint(c.context)
	return NewProspectiveCertificateChainCheck(c.I18nProvider, c.Result, currentCertificate, c.context, constraint)
}

// validationBeforeSunsetDate ports the private
// validationBeforeSunsetDate(CertificateWrapper, SubContext, Date).
//
// Java's getConstraintOrMaxLevel throws UnsupportedOperationException for a
// Level it does not know; the Go port returns that as an error, which this
// unreachable-by-construction call site raises as a panic, the way
// sav.AbstractAcceptanceValidation#cryptographic does.
func (c *X509CertificateValidation) validationBeforeSunsetDate(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlXCV] {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	levelRule, err := process.GetConstraintOrMaxLevel(constraint, enumerations.Level_WARN)
	if err != nil {
		panic(err)
	}
	return NewCertificateValidationBeforeSunsetDateWithIdCheck(c.I18nProvider, c.Result, certificate, validationTime,
		levelRule)
}

// prospectiveCertificateChainValidAtValidationTime ports the private
// prospectiveCertificateChainValidAtValidationTime(CertificateWrapper, SubContext, Date).
func (c *X509CertificateValidation) prospectiveCertificateChainValidAtValidationTime(
	certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext,
	validationTime time.Time) process.ChainItem[*jaxb.XmlXCV] {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	return NewProspectiveCertificateChainAtValidationTimeCheck(c.I18nProvider, c.Result, certificate, validationTime, constraint)
}

// trustServiceWithExpectedTypeIdentifier ports the private
// trustServiceWithExpectedTypeIdentifier(CertificateWrapper).
func (c *X509CertificateValidation) trustServiceWithExpectedTypeIdentifier(
	currentCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlXCV] {
	constraint := c.validationPolicy.TrustServiceTypeIdentifierConstraint(c.context)
	return NewTrustServiceTypeIdentifierCheck(c.I18nProvider, c.Result, currentCertificate, c.usageTime, c.context, constraint)
}

// trustServiceWithExpectedStatus ports the private
// trustServiceWithExpectedStatus(CertificateWrapper).
func (c *X509CertificateValidation) trustServiceWithExpectedStatus(
	currentCertificate *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlXCV] {
	constraint := c.validationPolicy.TrustServiceStatusConstraint(c.context)
	return NewTrustServiceStatusCheck(c.I18nProvider, c.Result, currentCertificate, c.usageTime, c.context, constraint)
}

// checkSubXCVResult ports the private checkSubXCVResult(XmlSubXCV).
func (c *X509CertificateValidation) checkSubXCVResult(subXCVResult *jaxb.XmlSubXCV) process.ChainItem[*jaxb.XmlXCV] {
	return NewCheckSubXCVResult(c.I18nProvider, c.Result, subXCVResult, c.FailLevelRule())
}

// checkTrustAnchorSubXCVResult ports the private
// checkTrustAnchorSubXCVResult(XmlSubXCV), whose body is an anonymous subclass
// of CheckSubXCVResult; see trustAnchorCheckSubXCVResult below.
func (c *X509CertificateValidation) checkTrustAnchorSubXCVResult(subXCVResult *jaxb.XmlSubXCV) process.ChainItem[*jaxb.XmlXCV] {
	return newTrustAnchorCheckSubXCVResult(c.I18nProvider, c.Result, subXCVResult, c.FailLevelRule())
}

// trustAnchorCheckSubXCVResult is the Go form of the anonymous CheckSubXCVResult
// subclass returned by checkTrustAnchorSubXCVResult(XmlSubXCV): the same check
// reporting a trust-anchor specific message and conclusion.
type trustAnchorCheckSubXCVResult struct {
	*CheckSubXCVResult
}

// newTrustAnchorCheckSubXCVResult builds the anonymous subclass and re-registers
// the overrides with the outer type, so that the base's self-calls reach the
// three methods overridden here rather than those of CheckSubXCVResult.
func newTrustAnchorCheckSubXCVResult(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlXCV],
	subXCVResult *jaxb.XmlSubXCV, constraint policy.LevelRule) *trustAnchorCheckSubXCVResult {
	c := &trustAnchorCheckSubXCVResult{
		CheckSubXCVResult: NewCheckSubXCVResult(i18nProvider, result, subXCVResult, constraint),
	}
	c.InitChainItem(c)
	return c
}

// ErrorMessageTag returns an i18n key of an error message to get. Port of the
// overridden getErrorMessageTag().
func (c *trustAnchorCheckSubXCVResult) ErrorMessageTag() i18n.MessageTag {
	return i18n.MessageTag_BBB_XCV_SUB_ANS_2
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port of
// the overridden getFailedIndicationForConclusion().
func (c *trustAnchorCheckSubXCVResult) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure. Port
// of the overridden getFailedSubIndicationForConclusion().
func (c *trustAnchorCheckSubXCVResult) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_NO_CERTIFICATE_CHAIN_FOUND_NO_POE
}

// prospectiveCertificateChainCheckEnforced ports the private
// prospectiveCertificateChainCheckEnforced().
func (c *X509CertificateValidation) prospectiveCertificateChainCheckEnforced() bool {
	constraint := c.validationPolicy.ProspectiveCertificateChainConstraint(c.context)
	return constraint != nil && enumerations.Level_FAIL == constraint.Level()
}

// isTrustAnchor ports the private
// isTrustAnchor(CertificateWrapper, Context, SubContext).
func (c *X509CertificateValidation) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper,
	context enumerations.Context, subContext enumerations.SubContext) bool {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(context, subContext)
	return process.IsTrustAnchor(certificateWrapper, c.currentTime, constraint)
}

// isTrustAnchorReached ports the private
// isTrustAnchorReached(CertificateWrapper, SubContext).
func (c *X509CertificateValidation) isTrustAnchorReached(certificateWrapper *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) bool {
	return c.isTrustAnchor(certificateWrapper, c.context, subContext) ||
		(certificateWrapper.IsTrusted() && !certificateWrapper.IsTrustedChain()) // second part is to filter only prospective certificate chains
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint).
func (c *X509CertificateValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	// collect all messages, except prospective certificate chain expiration warning (only final message should be returned)
	// XmlBlockType.SUB_XCV_TA.equals(getBlockType()): the generated Go BlockType
	// member is a *XmlBlockType, whose nil is Java's null.
	if !(constraint.BlockType != nil && jaxb.XmlBlockType_SUB_XCV_TA == *constraint.BlockType) {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion. Port
// of the overridden collectAdditionalMessages(XmlConclusion).
func (c *X509CertificateValidation) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	for _, subXCV := range c.Result.Value.SubXCV {
		c.CollectAllMessages(conclusion, subXCV.Conclusion)
		for _, constraint := range subXCV.Constraint {
			if constraint.BlockType != nil && jaxb.XmlBlockType_SUB_XCV_TA == *constraint.BlockType {
				if constraint.Error != nil {
					conclusion.Errors = c.removeMessage(conclusion.Errors, constraint.Error.Key)
				}
				if constraint.Warning != nil {
					conclusion.Warnings = c.removeMessage(conclusion.Warnings, constraint.Warning.Key)
				}
				if constraint.Info != nil {
					conclusion.Infos = c.removeMessage(conclusion.Infos, constraint.Info.Key)
				}
			}
		}
	}
}

// removeMessage ports the private removeMessage(List<XmlMessage>, String). Java
// mutates the list in place through removeIf; the Go model's message lists are
// plain slices, so the retained members are collected into a fresh slice the
// caller assigns back - the two are indistinguishable to every reader, and a
// fresh slice cannot disturb a list that happens to share the backing array.
//
// The Java message key is a String that may be null; the generated Go Key
// member is a *string, and messageKey.equals(other) is false for a null other,
// true only when both texts match. messageKey itself is never null here:
// getError()/getWarning()/getInfo() always carry a key.
func (c *X509CertificateValidation) removeMessage(messages []*jaxb.XmlMessage, messageKey *string) []*jaxb.XmlMessage {
	if utils.IsCollectionEmpty(messages) {
		return messages
	}
	kept := make([]*jaxb.XmlMessage, 0, len(messages))
	for _, xmlMessage := range messages {
		if messageKey != nil && xmlMessage.Key != nil && *messageKey == *xmlMessage.Key {
			continue
		}
		kept = append(kept, xmlMessage)
	}
	return kept
}
