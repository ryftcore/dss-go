// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/sub/SubX509CertificateValidation.java (DSS 6.5.RC1).
//
// Cross-chunk assumption: this file calls NewCertificateRevocationSelector,
// which XCVA (crs.CertificateRevocationSelector, ported into the same shared
// pkg xcv per the phase 8d brief) is expected to provide with a
// LatestAcceptableCertificateRevocation() *diagnostic.CertificateRevocationWrapper
// accessor and an Execute() *jaxb.XmlCRS method (the process.ChainBase
// convention every other Chain subclass in this port follows); and calls
// aov.NewRevocationDataAlgorithmObsolescenceValidation, which the AOV porter
// is expected to add to package aov with an Execute() *jaxb.XmlAOV method,
// mirroring aov.AlgorithmObsolescenceValidation's Java sibling class of the
// same name. See S8D_BRIEF.md.
package xcv

import (
	"time"

	jaxb "github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/aov"
)

// SubX509CertificateValidation is the sub X509 certificate validation.
type SubX509CertificateValidation struct {
	*process.ChainBase[*jaxb.XmlSubXCV]

	// currentCertificate is the certificate to check.
	currentCertificate *diagnostic.CertificateWrapper

	// validationDate is the validation time.
	validationDate time.Time

	// currentTime is the current time when validation is performed.
	currentTime time.Time

	// context is the validation context.
	context enumerations.Context

	// subContext is the validation subContext.
	subContext enumerations.SubContext

	// aov is the result of cryptographic algorithms validation.
	aov *jaxb.XmlAOV

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy
}

// NewSubX509CertificateValidation is the default constructor. Port of
// SubX509CertificateValidation(I18nProvider, CertificateWrapper, Date, Date, Context, SubContext, XmlAOV, ValidationPolicy).
func NewSubX509CertificateValidation(i18nProvider *i18n.I18nProvider, currentCertificate *diagnostic.CertificateWrapper,
	validationDate time.Time, currentTime time.Time, context enumerations.Context, subContext enumerations.SubContext,
	aovResult *jaxb.XmlAOV, validationPolicy policy.ValidationPolicy) *SubX509CertificateValidation {
	xmlSubXCV := &jaxb.XmlSubXCV{}
	xmlSubXCV.Id = currentCertificate.Id()
	trustAnchor := currentCertificate.IsTrusted()
	xmlSubXCV.TrustAnchor = &trustAnchor
	selfSigned := currentCertificate.IsSelfSigned()
	xmlSubXCV.SelfSigned = &selfSigned

	c := &SubX509CertificateValidation{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlSubXCV,
			&xmlSubXCV.XmlConstraintsConclusionContent, &xmlSubXCV.XmlConstraintsConclusionAttrs)),
		currentCertificate: currentCertificate,
		validationDate:     validationDate,
		currentTime:        currentTime,
		context:            context,
		subContext:         subContext,
		aov:                aovResult,
		validationPolicy:   validationPolicy,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *SubX509CertificateValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SUB_XCV
}

// InitChain initializes the chain. Port of initChain().
func (c *SubX509CertificateValidation) InitChain() {
	var item process.ChainItem[*jaxb.XmlSubXCV]

	if c.currentCertificate.IsTrusted() {

		if c.currentCertificate.TrustStartDate() != nil || c.currentCertificate.TrustSunsetDate() != nil {

			item = c.validationBeforeSunsetDate(c.currentCertificate, c.subContext, c.currentTime)
			c.FirstItem = item

			if !process.IsTrustAnchor(c.currentCertificate, c.currentTime, c.FailLevelRule()) {
				item = item.SetNextItem(c.otherTrustAnchorAvailable(c.currentCertificate, c.subContext))
			}

		}

		if c.isTrustAnchorReached(c.currentCertificate, c.subContext) {
			// Skip for Trusted Certificate
			return
		}

	}

	if item == nil {
		item = c.serialNumber(c.currentCertificate, c.subContext)
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.serialNumber(c.currentCertificate, c.subContext))
	}

	item = item.SetNextItem(c.surname(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.givenName(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.commonName(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.pseudoUsage(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.pseudonym(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.title(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.email(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.country(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.locality(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.state(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.organizationIdentifier(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.organizationUnit(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.organizationName(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.selfSigned(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.notSelfSigned(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePolicyIds(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePolicyQualifiedIds(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePolicySupportedByQSCDIds(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcCompliance(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcEuLimitValueCurrency(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateMinQcEuLimitValue(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcEuRetentionPeriod(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcSSCD(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcEuPDSLocation(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcType(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcCCLegislation(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateIssuedToNaturalPerson(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateIssuedToLegalPerson(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateSemanticsIdentifier(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePS2DQcRolesOfPSP(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePS2DQcCompetentAuthorityName(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificatePS2DQcCompetentAuthorityId(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcQCSDLegislation(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcIdentificationMethod(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcPSBCountryOfLegislation(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcPSBAuthSourceIdentification(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateQcPSBLegislationIdentification(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.certificateSignatureValid(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.ca(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.issuerName(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.maxPathLength(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.keyUsage(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.extendedKeyUsage(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.aiaPresent(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.policyTree(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.nameConstraints(c.currentCertificate, c.subContext))

	/*
	 * RFC 5280, "4.2.1.1. Authority Key Identifier".
	 * There is one exception; where a CA distributes its public key in the form of a "self-signed"
	 * certificate, the authority key identifier MAY be omitted.
	 */
	if !c.currentCertificate.IsSelfSigned() {
		item = item.SetNextItem(c.authorityKeyIdentifierPresent(c.currentCertificate, c.subContext))
	}

	item = item.SetNextItem(c.subjectKeyIdentifierPresent(c.currentCertificate, c.subContext))

	if c.currentCertificate.IsNoRevAvail() {
		item = item.SetNextItem(c.noRevAvail(c.currentCertificate, c.subContext))
	}

	item = item.SetNextItem(c.supportedCriticalCertificateExtensions(c.currentCertificate, c.subContext))

	item = item.SetNextItem(c.forbiddenCertificateExtensions(c.currentCertificate, c.subContext))

	var latestCertificateRevocation *diagnostic.CertificateRevocationWrapper

	revocationDataRequired := c.revocationDataRequired(c.currentCertificate, c.subContext)

	isRevocationDataRequired := revocationDataRequired.Process()
	if isRevocationDataRequired {

		item = item.SetNextItem(c.revocationInfoAccessPresent(c.currentCertificate, c.subContext))

		item = item.SetNextItem(c.revocationDataPresent(c.currentCertificate, c.subContext))

		if utils.IsCollectionNotEmpty(c.currentCertificate.CertificateRevocationData()) {

			certificateRevocationSelector := NewCertificateRevocationSelector(
				c.I18nProvider, c.currentCertificate, c.validationDate, c.validationPolicy)
			xmlCRS := certificateRevocationSelector.Execute()
			c.Result.Value.CRS = xmlCRS

			item = item.SetNextItem(c.checkCertificateRevocationSelectorResult(xmlCRS))

			latestCertificateRevocation = certificateRevocationSelector.LatestAcceptableCertificateRevocation()

			if latestCertificateRevocation != nil && latestCertificateRevocation.IsRevoked() {
				c.attachRevocationInformation(latestCertificateRevocation)
			}

			if c.IsValid(&xmlCRS.XmlConstraintsConclusionContent) {

				item = item.SetNextItem(c.certificateNotRevoked(latestCertificateRevocation, c.subContext, c.validationDate))

				item = item.SetNextItem(c.certificateNotOnHold(latestCertificateRevocation, c.subContext, c.validationDate))

				rfc := NewRevocationFreshnessChecker(c.I18nProvider, &latestCertificateRevocation.RevocationWrapper,
					c.validationDate, c.context, c.subContext, c.validationPolicy)
				rfcResult := rfc.Execute()
				c.Result.Value.RFC = rfcResult

				item = item.SetNextItem(c.checkRevocationFreshnessCheckerResult(rfcResult))

			}

		}

	} else {
		item = item.SetNextItem(revocationDataRequired)
	}

	// NOTE: cryptographic constraint shall be validated against the current time,
	// and not against time returned by the used validation model
	item = item.SetNextItem(c.certificateCryptographic())

	if latestCertificateRevocation != nil {
		item = item.SetNextItem(c.revocationCryptographic(latestCertificateRevocation))
	}

	if enumerations.SubContextSigningCert == c.subContext {

		item = item.SetNextItem(c.certificateValidityRange(c.currentCertificate,
			latestCertificateRevocation, isRevocationDataRequired, c.subContext, c.currentTime))

		if latestCertificateRevocation != nil {
			revocationIssuerCertificate := latestCertificateRevocation.SigningCertificate()
			if revocationIssuerCertificate != nil {
				if c.isTrustAnchor(revocationIssuerCertificate, enumerations.ContextRevocation, enumerations.SubContextSigningCert) {
					item = item.SetNextItem(c.revocationDataIssuerTrusted(revocationIssuerCertificate)) //nolint:staticcheck // mirrors upstream SubX509CertificateValidation#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
				} else {
					item = item.SetNextItem(c.revocationIssuerValidityRange(latestCertificateRevocation, c.subContext, c.currentTime)) //nolint:staticcheck // mirrors upstream SubX509CertificateValidation#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
				}
			}
		}

	}
}

// attachRevocationInformation ports the private
// attachRevocationInformation(CertificateRevocationWrapper).
func (c *SubX509CertificateValidation) attachRevocationInformation(certificateRevocation *diagnostic.CertificateRevocationWrapper) {
	revocationInfo := &jaxb.XmlRevocationInformation{}
	revocationInfo.CertificateId = c.currentCertificate.Id()
	revocationInfo.RevocationId = certificateRevocation.Id()
	if revocationDate := certificateRevocation.RevocationDate(); revocationDate != nil {
		revocationInfo.RevocationDate = jaxb.XSDateTime(*revocationDate)
	}
	if reason := certificateRevocation.Reason(); reason != "" {
		value := jaxb.RevocationReasonValue(reason)
		revocationInfo.Reason = &value
	}
	c.Result.Value.RevocationInfo = revocationInfo
}

func (c *SubX509CertificateValidation) validationBeforeSunsetDate(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	level, err := process.GetConstraintOrMaxLevel(constraint, enumerations.LevelWarn)
	if err != nil {
		panic(err)
	}
	return NewCertificateValidationBeforeSunsetDateCheck(c.I18nProvider, c.Result, certificate, validationTime, level)
}

func (c *SubX509CertificateValidation) otherTrustAnchorAvailable(certificate *diagnostic.CertificateWrapper,
	subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	return NewOtherTrustAnchorExistsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateValidityRange(certificate *diagnostic.CertificateWrapper,
	usedCertificateRevocation *diagnostic.CertificateRevocationWrapper, revocationDataRequired bool,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNotExpiredConstraint(c.context, subContext)
	isRevocationIssuerTrusted := usedCertificateRevocation != nil && usedCertificateRevocation.SigningCertificate() != nil &&
		c.isTrustAnchor(usedCertificateRevocation.SigningCertificate(), enumerations.ContextRevocation, enumerations.SubContextSigningCert)
	revocationIssuerCheckEnforced := c.revocationIssuerCheckEnforced(c.context, subContext)
	return NewCertificateValidityRangeCheck(c.I18nProvider, c.Result, certificate, usedCertificateRevocation,
		revocationDataRequired, isRevocationIssuerTrusted, revocationIssuerCheckEnforced, validationTime, constraint)
}

func (c *SubX509CertificateValidation) revocationIssuerCheckEnforced(context enumerations.Context, subContext enumerations.SubContext) bool {
	constraint := c.validationPolicy.RevocationIssuerNotExpiredConstraint(context, subContext)
	return constraint != nil && enumerations.LevelFail == constraint.Level()
}

func (c *SubX509CertificateValidation) revocationDataIssuerTrusted(revocationIssuer *diagnostic.CertificateWrapper) process.ChainItem[*jaxb.XmlSubXCV] {
	revocationDataSunsetDate := c.validationPolicy.CertificateSunsetDateConstraint(enumerations.ContextRevocation, enumerations.SubContextSigningCert)
	return NewRevocationIssuerTrustedCheck(c.I18nProvider, c.Result, revocationIssuer, c.currentTime, revocationDataSunsetDate, c.WarnLevelRule())
}

func (c *SubX509CertificateValidation) revocationIssuerValidityRange(usedCertificateRevocation *diagnostic.CertificateRevocationWrapper,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.RevocationIssuerNotExpiredConstraint(c.context, subContext)
	return NewRevocationIssuerValidityRangeCheck(c.I18nProvider, c.Result, &usedCertificateRevocation.RevocationWrapper, validationTime, constraint)
}

func (c *SubX509CertificateValidation) ca(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateCAConstraint(c.context, subContext)
	return NewBasicConstraintsCACheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) issuerName(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateIssuerNameConstraint(c.context, subContext)
	return NewCertificateIssuerNameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) maxPathLength(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateMaxPathLengthConstraint(c.context, subContext)
	return NewBasicConstraintsMaxPathLengthCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) keyUsage(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateKeyUsageConstraint(c.context, subContext)
	return NewKeyUsageCheck(c.I18nProvider, c.Result, certificate, c.context, subContext, constraint)
}

func (c *SubX509CertificateValidation) extendedKeyUsage(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateExtendedKeyUsageConstraint(c.context, subContext)
	return NewExtendedKeyUsageCheck(c.I18nProvider, c.Result, certificate, c.context, subContext, constraint)
}

func (c *SubX509CertificateValidation) aiaPresent(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateAuthorityInfoAccessPresentConstraint(c.context, subContext)
	return NewAuthorityInfoAccessPresentCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) policyTree(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePolicyTreeConstraint(c.context, subContext)
	return NewCertificatePolicyTreeCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) nameConstraints(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNameConstraintsConstraint(c.context, subContext)
	return NewCertificateNameConstraintsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) authorityKeyIdentifierPresent(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateAuthorityKeyIdentifierPresentConstraint(c.context, subContext)
	return NewAuthorityKeyIdentifierPresentCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) subjectKeyIdentifierPresent(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSubjectKeyIdentifierPresentConstraint(c.context, subContext)
	return NewSubjectKeyIdentifierPresentCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) noRevAvail(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNoRevAvailConstraint(c.context, subContext)
	return NewNoRevAvailCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) supportedCriticalCertificateExtensions(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSupportedCriticalExtensionsConstraint(c.context, subContext)
	return NewCertificateSupportedCriticalExtensionsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) forbiddenCertificateExtensions(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateForbiddenExtensionsConstraint(c.context, subContext)
	return NewCertificateForbiddenExtensionsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) revocationDataRequired(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) *RevocationDataRequiredCheck[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.RevocationDataSkipConstraint(c.context, subContext)
	sunsetDateConstraint := c.validationPolicy.CertificateSunsetDateConstraint(c.context, subContext)
	return NewRevocationDataRequiredCheck(c.I18nProvider, c.Result, certificate, c.currentTime, sunsetDateConstraint, constraint)
}

func (c *SubX509CertificateValidation) revocationInfoAccessPresent(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateRevocationInfoAccessPresentConstraint(c.context, subContext)
	return NewRevocationInfoAccessPresentCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) revocationDataPresent(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.RevocationDataAvailableConstraint(c.context, subContext)
	return NewRevocationDataAvailableCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) checkCertificateRevocationSelectorResult(crsResult *jaxb.XmlCRS) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.AcceptableRevocationDataFoundConstraint(c.context, c.subContext)
	return NewCertificateRevocationSelectorResultCheck(c.I18nProvider, c.Result, crsResult, constraint)
}

func (c *SubX509CertificateValidation) checkRevocationFreshnessCheckerResult(rfcResult *jaxb.XmlRFC) process.ChainItem[*jaxb.XmlSubXCV] {
	return NewRevocationFreshnessCheckerResultCheck(c.I18nProvider, c.Result, rfcResult, c.FailLevelRule())
}

func (c *SubX509CertificateValidation) surname(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSurnameConstraint(c.context, subContext)
	return NewSurnameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) givenName(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateGivenNameConstraint(c.context, subContext)
	return NewGivenNameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) commonName(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateCommonNameConstraint(c.context, subContext)
	return NewCommonNameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) pseudonym(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePseudonymConstraint(c.context, subContext)
	return NewPseudonymCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) title(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateTitleConstraint(c.context, subContext)
	return NewTitleCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) email(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateEmailConstraint(c.context, subContext)
	return NewEmailCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) country(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateCountryConstraint(c.context, subContext)
	return NewCountryCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) locality(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateLocalityConstraint(c.context, subContext)
	return NewLocalityCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) state(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateStateConstraint(c.context, subContext)
	return NewStateCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) organizationIdentifier(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateOrganizationIdentifierConstraint(c.context, subContext)
	return NewOrganizationIdentifierCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) organizationName(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateOrganizationNameConstraint(c.context, subContext)
	return NewOrganizationNameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) organizationUnit(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateOrganizationUnitConstraint(c.context, subContext)
	return NewOrganizationUnitCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) serialNumber(signingCertificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSerialNumberConstraint(c.context, subContext)
	return NewSerialNumberCheck(c.I18nProvider, c.Result, signingCertificate, constraint)
}

func (c *SubX509CertificateValidation) pseudoUsage(signingCertificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePseudoUsageConstraint(c.context, subContext)
	return NewPseudoUsageCheck(c.I18nProvider, c.Result, signingCertificate, constraint)
}

func (c *SubX509CertificateValidation) certificateSignatureValid(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSignatureConstraint(c.context, subContext)
	return NewCertificateSignatureValidCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateNotRevoked(latestCertificateRevocation *diagnostic.CertificateRevocationWrapper,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNotRevokedConstraint(c.context, subContext)
	return NewCertificateNotRevokedCheck(c.I18nProvider, c.Result, latestCertificateRevocation, validationTime, constraint, subContext)
}

func (c *SubX509CertificateValidation) certificateNotOnHold(latestCertificateRevocation *diagnostic.CertificateRevocationWrapper,
	subContext enumerations.SubContext, validationTime time.Time) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNotOnHoldConstraint(c.context, subContext)
	return NewCertificateNotOnHoldCheck(c.I18nProvider, c.Result, latestCertificateRevocation, validationTime, constraint)
}

func (c *SubX509CertificateValidation) notSelfSigned(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateNotSelfSignedConstraint(c.context, subContext)
	return NewCertificateNotSelfSignedCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) selfSigned(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSelfSignedConstraint(c.context, subContext)
	return NewCertificateSelfSignedCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePolicyIds(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePolicyIdsConstraint(c.context, subContext)
	return NewCertificatePolicyIdsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePolicyQualifiedIds(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePolicyQualificationIdsConstraint(c.context, subContext)
	return NewCertificatePolicyQualifiedIdsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePolicySupportedByQSCDIds(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePolicySupportedByQSCDIdsConstraint(c.context, subContext)
	return NewCertificatePolicySupportedByQSCDIdsCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcCompliance(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQCComplianceConstraint(c.context, subContext)
	return NewCertificateQcComplianceCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateMinQcEuLimitValue(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateMinQcEuLimitValueConstraint(c.context, subContext)
	return NewCertificateMinQcTransactionLimitCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcEuLimitValueCurrency(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcEuLimitValueCurrencyConstraint(c.context, subContext)
	return NewCertificateQcEuLimitValueCurrencyCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcEuRetentionPeriod(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateMinQcEuRetentionPeriodConstraint(c.context, subContext)
	return NewCertificateMinQcEuRetentionPeriodCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcSSCD(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcSSCDConstraint(c.context, subContext)
	return NewCertificateQcSSCDCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcEuPDSLocation(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcEuPDSLocationConstraint(c.context, subContext)
	return NewCertificateQcEuPDSLocationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcType(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcTypeConstraint(c.context, subContext)
	return NewCertificateQcTypeCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcCCLegislation(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcCCLegislationConstraint(c.context, subContext)
	return NewCertificateQcCCLegislationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateIssuedToNaturalPerson(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateIssuedToNaturalPersonConstraint(c.context, subContext)
	return NewCertificateIssuedToNaturalPersonCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateIssuedToLegalPerson(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateIssuedToLegalPersonConstraint(c.context, subContext)
	return NewCertificateIssuedToLegalPersonCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateSemanticsIdentifier(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateSemanticsIdentifierConstraint(c.context, subContext)
	return NewCertificateSemanticsIdentifierCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePS2DQcRolesOfPSP(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePS2DQcTypeRolesOfPSPConstraint(c.context, subContext)
	return NewCertificatePS2DQcRolesOfPSPCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePS2DQcCompetentAuthorityName(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePS2DQcCompetentAuthorityNameConstraint(c.context, subContext)
	return NewCertificatePS2DQcCompetentAuthorityNameCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificatePS2DQcCompetentAuthorityId(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificatePS2DQcCompetentAuthorityIdConstraint(c.context, subContext)
	return NewCertificatePS2DQcCompetentAuthorityIdCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcQCSDLegislation(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcQSCDLegislationConstraint(c.context, subContext)
	return NewCertificateQcQSCDLegislationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcIdentificationMethod(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcIdentificationMethodConstraint(c.context, subContext)
	return NewCertificateQcIdentificationMethodCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcPSBCountryOfLegislation(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcPSBCountryOfLegislationConstraint(c.context, subContext)
	return NewCertificateQcPSBCountryOfLegislationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcPSBAuthSourceIdentification(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcPSBAuthSourceIdentificationConstraint(c.context, subContext)
	return NewCertificateQcPSBAuthSourceIdentificationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateQcPSBLegislationIdentification(certificate *diagnostic.CertificateWrapper, subContext enumerations.SubContext) process.ChainItem[*jaxb.XmlSubXCV] {
	constraint := c.validationPolicy.CertificateQcPSBLegislationIdentificationConstraint(c.context, subContext)
	return NewCertificateQcPSBLegislationIdentificationCheck(c.I18nProvider, c.Result, certificate, constraint)
}

func (c *SubX509CertificateValidation) certificateCryptographic() process.ChainItem[*jaxb.XmlSubXCV] {
	certificatePosition, err := process.GetSubContextPosition(enumerations.ContextCertificate, c.subContext)
	if err != nil {
		panic(err)
	}
	return NewCertificateAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, c.aov, c.currentTime, certificatePosition, c.currentCertificate.Id())
}

func (c *SubX509CertificateValidation) revocationCryptographic(revocationData *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlSubXCV] {
	// NOTE: we need to execute it explicitly in order to avoid a circular reference on revocation data validation
	revocationAOV := aov.NewRevocationDataAlgorithmObsolescenceValidation(c.I18nProvider, &revocationData.RevocationWrapper, c.currentTime, c.validationPolicy)
	position, err := process.GetCryptoPosition(enumerations.ContextRevocation)
	if err != nil {
		panic(err)
	}
	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, c.Result, revocationAOV.Execute(), c.currentTime, position, revocationData.Id())
}

func (c *SubX509CertificateValidation) isTrustAnchorReached(certificateWrapper *diagnostic.CertificateWrapper, subContext enumerations.SubContext) bool {
	return c.isTrustAnchor(certificateWrapper, c.context, subContext) || !certificateWrapper.IsTrustedChain()
}

func (c *SubX509CertificateValidation) isTrustAnchor(certificateWrapper *diagnostic.CertificateWrapper, context enumerations.Context, subContext enumerations.SubContext) bool {
	constraint := c.validationPolicy.CertificateSunsetDateConstraint(context, subContext)
	return process.IsTrustAnchor(certificateWrapper, c.currentTime, constraint)
}

// CollectAdditionalMessages fills additional messages into the conclusion.
// Port of collectAdditionalMessages(XmlConclusion).
func (c *SubX509CertificateValidation) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	c.ChainBase.CollectAdditionalMessages(conclusion)
	xmlCRS := c.Result.Value.CRS
	if xmlCRS != nil && c.IsValid(&xmlCRS.XmlConstraintsConclusionContent) {
		c.CollectAllMessages(conclusion, xmlCRS.Conclusion)
	}
	xmlRFC := c.Result.Value.RFC
	if xmlRFC != nil && c.IsValid(&xmlRFC.XmlConstraintsConclusionContent) {
		c.CollectAllMessages(conclusion, xmlRFC.Conclusion)
	}
	if c.aov != nil && c.IsValid(&c.aov.XmlConstraintsConclusionContent) {
		c.CollectAllMessages(conclusion, c.aov.Conclusion)
	}
}

// CollectMessages collects required messages from the given xmlConstraint to
// the given conclusion. Port of collectMessages(XmlConclusion, XmlConstraint).
func (c *SubX509CertificateValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraint.BlockType == nil || jaxb.XmlBlockTypeAOV != *constraint.BlockType {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}
