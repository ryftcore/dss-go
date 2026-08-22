// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/xcv/crs/CertificateRevocationSelector.java (DSS 6.5.RC1).
//
// See x509_certificate_validation.go for the package-flattening note.
//
// The class is designed for extension - vpfltvd.LongTermValidationCertificateRevocationSelector
// overrides verifyRevocationData, getRevocationAcceptanceValidationResult,
// acceptableRevocationDataAvailable and collectMessages - so every self-call on
// an overridable protected method is routed through
// CertificateRevocationSelectorOverrides, the way Chain routes its own through
// ChainOverrides. A subclass registers itself once, with
// InitCertificateRevocationSelector, which also performs the Chain registration.
package xcv

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// CertificateRevocationSelectorOverrides declares the overridable protected
// methods of CertificateRevocationSelector that the base implementation calls
// back into, on top of the Chain ones. Every method a subclass does not define
// is supplied by the embedded CertificateRevocationSelector through ordinary Go
// method promotion.
type CertificateRevocationSelectorOverrides interface {
	process.ChainOverrides

	// CertificateRevocationData returns available certificate revocation data
	// to be validated. Port of getCertificateRevocationData().
	CertificateRevocationData() []*diagnostic.CertificateRevocationWrapper
	// VerifyRevocationData verifies the given revocation data and returns the
	// resulting ChainItem. Port of
	// verifyRevocationData(ChainItem, CertificateRevocationWrapper).
	VerifyRevocationData(item process.ChainItem[*jaxb.XmlCRS],
		revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS]
	// RevocationAcceptanceValidationResult returns a RevocationAcceptanceValidation
	// result for the given revocation token. Port of
	// getRevocationAcceptanceValidationResult(CertificateRevocationWrapper).
	RevocationAcceptanceValidationResult(revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlRAC
	// RevocationAcceptable checks whether the revocation data has passed the
	// Revocation Acceptance Validation. Port of revocationAcceptable(XmlRAC).
	RevocationAcceptable(racResult *jaxb.XmlRAC) process.ChainItem[*jaxb.XmlCRS]
	// AcceptableRevocationDataAvailable checks whether the acceptable revocation
	// data is available. Port of acceptableRevocationDataAvailable().
	AcceptableRevocationDataAvailable() process.ChainItem[*jaxb.XmlCRS]
}

// CertificateRevocationSelector validates revocation data for a given
// certificate and returns the latest valid entry.
type CertificateRevocationSelector struct {
	*process.ChainBase[*jaxb.XmlCRS]

	// Certificate is the certificate to get a latest valid revocation data for.
	// Exported because Java declares the field protected.
	Certificate *diagnostic.CertificateWrapper

	// CurrentTime is the validation time. Exported because Java declares the
	// field protected.
	CurrentTime time.Time

	// ValidationPolicy is the validation policy. Exported because Java declares
	// the field protected.
	ValidationPolicy policy.ValidationPolicy

	// RevocationDataValidityMap contains validation results of the revocation
	// data processing. Exported because Java declares the field protected.
	//
	// Java keys the map with the RevocationWrapper itself, whose equals/hashCode
	// (AbstractTokenProxy, refined by RevocationWrapper) compare nothing but the
	// token id; the Go map is therefore keyed by that id directly. Java's
	// nullable Boolean value is the comma-ok of a lookup.
	RevocationDataValidityMap map[string]bool

	// validatedTokens is the internal set of processed tokens (avoids infinite
	// loop). Java's Set<String> is shared by reference between the selector and
	// the acceptance checkers it spawns; the Go map is the same shared object.
	validatedTokens map[string]struct{}

	// latestCertificateRevocation is the latest acceptable certificate
	// revocation, to be returned after the selector execution.
	latestCertificateRevocation *diagnostic.CertificateRevocationWrapper

	// overrides points back at the concrete selector; see
	// InitCertificateRevocationSelector.
	overrides CertificateRevocationSelectorOverrides
}

// NewCertificateRevocationSelector is the default constructor. Port of
// CertificateRevocationSelector(I18nProvider, CertificateWrapper, Date, ValidationPolicy),
// which delegates to the constructor below with a fresh empty set.
func NewCertificateRevocationSelector(i18nProvider *i18n.I18nProvider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	validationPolicy policy.ValidationPolicy) *CertificateRevocationSelector {
	return NewCertificateRevocationSelectorWithValidatedTokens(i18nProvider, certificate, currentTime,
		validationPolicy, make(map[string]struct{}))
}

// NewCertificateRevocationSelectorWithValidatedTokens is the constructor with a
// set of validated tokens. Port of
// CertificateRevocationSelector(I18nProvider, CertificateWrapper, Date, ValidationPolicy, Set).
func NewCertificateRevocationSelectorWithValidatedTokens(i18nProvider *i18n.I18nProvider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, validatedTokens map[string]struct{}) *CertificateRevocationSelector {
	c := &CertificateRevocationSelector{}
	c.InitCertificateRevocationSelectorState(i18nProvider, certificate, currentTime, validationPolicy, validatedTokens)
	c.InitCertificateRevocationSelector(c)
	return c
}

// InitCertificateRevocationSelectorState wires the shared state, the way the
// Java constructor body does. A subclass calls it before
// InitCertificateRevocationSelector, in place of the Java super(...) call.
func (c *CertificateRevocationSelector) InitCertificateRevocationSelectorState(i18nProvider *i18n.I18nProvider,
	certificate *diagnostic.CertificateWrapper, currentTime time.Time,
	validationPolicy policy.ValidationPolicy, validatedTokens map[string]struct{}) {
	xmlCRS := &jaxb.XmlCRS{}
	c.ChainBase = process.NewChainBase(i18nProvider, process.NewResult(xmlCRS,
		&xmlCRS.XmlConstraintsConclusionContent, &xmlCRS.XmlConstraintsConclusionAttrs))
	c.Certificate = certificate
	c.CurrentTime = currentTime
	c.ValidationPolicy = validationPolicy
	c.RevocationDataValidityMap = make(map[string]bool)
	c.validatedTokens = validatedTokens

	id := certificate.Id()
	xmlCRS.Id = &id
	c.validatedTokens[certificate.Id()] = struct{}{}
}

// InitCertificateRevocationSelector registers the concrete selector with its
// base so that the base can dispatch to the overridden methods, Chain's
// included. It must be called exactly once, by the concrete selector's
// constructor, before Execute.
func (c *CertificateRevocationSelector) InitCertificateRevocationSelector(overrides CertificateRevocationSelectorOverrides) {
	c.overrides = overrides
	c.InitChainBase(overrides)
}

// crsOverrides returns the registered overrides, panicking when the concrete
// selector forgot to call InitCertificateRevocationSelector.
func (c *CertificateRevocationSelector) crsOverrides() CertificateRevocationSelectorOverrides {
	if c.overrides == nil {
		panic("CertificateRevocationSelector was not initialised: the concrete selector must call InitCertificateRevocationSelector in its constructor")
	}
	return c.overrides
}

// Title returns the title of the building block. Port of getTitle().
func (c *CertificateRevocationSelector) Title() i18n.MessageTag {
	return i18n.MessageTag_CRS
}

// InitChain initializes the chain. Port of initChain().
func (c *CertificateRevocationSelector) InitChain() {
	overrides := c.crsOverrides()

	var item process.ChainItem[*jaxb.XmlCRS]

	for _, revocationWrapper := range overrides.CertificateRevocationData() {

		item = overrides.VerifyRevocationData(item, revocationWrapper)

		if valid, ok := c.RevocationDataValidityMap[revocationWrapper.Id()]; ok && valid &&
			(c.latestCertificateRevocation == nil || (revocationWrapper.ProductionDate() != nil &&
				c.latestCertificateRevocation.ProductionDate().Before(*revocationWrapper.ProductionDate()))) {
			c.latestCertificateRevocation = revocationWrapper
		}
	}

	if c.latestCertificateRevocation != nil {
		id := c.latestCertificateRevocation.Id()
		c.Result.Value.LatestAcceptableRevocationId = &id
	}

	if item == nil {
		item = overrides.AcceptableRevocationDataAvailable()
		c.FirstItem = item
	} else {
		item = item.SetNextItem(overrides.AcceptableRevocationDataAvailable()) //nolint:staticcheck // mirrors upstream CertificateRevocationSelector#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
}

// CertificateRevocationData returns available certificate revocation data to be
// validated. Port of getCertificateRevocationData().
func (c *CertificateRevocationSelector) CertificateRevocationData() []*diagnostic.CertificateRevocationWrapper {
	return c.Certificate.CertificateRevocationData()
}

// VerifyRevocationData verifies the given revocation data and returns the
// resulting ChainItem. item is the last initialized chain item to be processed
// in prior to the revocation validation. Port of
// verifyRevocationData(ChainItem, CertificateRevocationWrapper).
func (c *CertificateRevocationSelector) VerifyRevocationData(item process.ChainItem[*jaxb.XmlCRS],
	revocationWrapper *diagnostic.CertificateRevocationWrapper) process.ChainItem[*jaxb.XmlCRS] {
	overrides := c.crsOverrides()

	racResult := overrides.RevocationAcceptanceValidationResult(revocationWrapper)

	if racResult != nil {
		if item == nil {
			item = overrides.RevocationAcceptable(racResult)
			c.FirstItem = item
		} else {
			item = item.SetNextItem(overrides.RevocationAcceptable(racResult))
		}
	}

	var racContent *jaxb.XmlConstraintsConclusionContent
	if racResult != nil {
		racContent = &racResult.XmlConstraintsConclusionContent
	}
	c.RevocationDataValidityMap[revocationWrapper.Id()] = c.IsValid(racContent)

	return item
}

// RevocationAcceptanceValidationResult returns a RevocationAcceptanceValidation
// result for the given revocation token. Port of
// getRevocationAcceptanceValidationResult(CertificateRevocationWrapper).
func (c *CertificateRevocationSelector) RevocationAcceptanceValidationResult(
	revocationWrapper *diagnostic.CertificateRevocationWrapper) *jaxb.XmlRAC {
	rac := NewRevocationAcceptanceChecker(
		c.I18nProvider, c.Certificate, revocationWrapper, c.CurrentTime, c.ValidationPolicy, c.validatedTokens)
	racResult := rac.Execute()

	c.Result.Value.RAC = append(c.Result.Value.RAC, racResult)

	return racResult
}

// LatestAcceptableCertificateRevocation returns the latest acceptable
// certificate revocation data. Port of getLatestAcceptableCertificateRevocation().
//
// NOTE: method Execute() shall be called before.
func (c *CertificateRevocationSelector) LatestAcceptableCertificateRevocation() *diagnostic.CertificateRevocationWrapper {
	return c.latestCertificateRevocation
}

// RevocationAcceptable checks whether the revocation data has passed the
// Revocation Acceptance Validation. Port of revocationAcceptable(XmlRAC).
func (c *CertificateRevocationSelector) RevocationAcceptable(racResult *jaxb.XmlRAC) process.ChainItem[*jaxb.XmlCRS] {
	return NewRevocationAcceptanceCheckerResultCheck(c.I18nProvider, c.Result, racResult, c.WarnLevelRule())
}

// AcceptableRevocationDataAvailable checks whether the acceptable revocation
// data is available. Port of acceptableRevocationDataAvailable().
func (c *CertificateRevocationSelector) AcceptableRevocationDataAvailable() process.ChainItem[*jaxb.XmlCRS] {
	// Java hands the CertificateRevocationWrapper straight to a RevocationWrapper
	// parameter; the Go wrapper embeds its base by value, so the base is addressed
	// out of it, and a null latest revocation stays a nil *RevocationWrapper.
	var acceptableRevocationData *diagnostic.RevocationWrapper
	if c.latestCertificateRevocation != nil {
		acceptableRevocationData = &c.latestCertificateRevocation.RevocationWrapper
	}
	return NewAcceptableRevocationDataAvailableCheck(c.I18nProvider, c.Result, acceptableRevocationData, c.FailLevelRule())
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of the overridden
// collectMessages(XmlConclusion, XmlConstraint).
func (c *CertificateRevocationSelector) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	// collect all messages from not RAC checks, collect from RAC only when all of them failed
	// XmlBlockType.RAC.equals(getBlockType()): the generated Go BlockType member
	// is a *XmlBlockType, whose nil is Java's null.
	if !(constraint.BlockType != nil && jaxb.XmlBlockTypeRAC == *constraint.BlockType) ||
		!c.IsValid(&c.Result.Value.XmlConstraintsConclusionContent) {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// CollectAdditionalMessages fills additional messages into the conclusion. Port
// of the overridden collectAdditionalMessages(XmlConclusion).
func (c *CertificateRevocationSelector) CollectAdditionalMessages(conclusion *jaxb.XmlConclusion) {
	if !c.IsValid(&c.Result.Value.XmlConstraintsConclusionContent) {
		for _, rac := range c.Result.Value.RAC {
			c.CollectAllMessages(conclusion, rac.Conclusion)
		}
	} else {
		// collect additional messages for the valid RAC(s)
		for _, rac := range c.Result.Value.RAC {
			if c.IsValid(&rac.XmlConstraintsConclusionContent) {
				c.CollectAllMessages(conclusion, rac.Conclusion)
			}
		}
	}
}
