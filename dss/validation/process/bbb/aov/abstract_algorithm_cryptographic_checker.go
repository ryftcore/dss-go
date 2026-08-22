// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/AbstractAlgorithmCryptographicChecker.java (DSS 6.5.RC1).
//
// Package placement deviation: Java's eu.europa.esig.dss.validation.process.bbb.aov.cc
// package (the CryptographicChecker family) is flattened into this pkg aov.
package aov

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// algorithmUnidentified is the name string for a unidentified (unsupported)
// algorithm. Port of the protected static final String ALGORITHM_UNIDENTIFIED.
const algorithmUnidentified = "UNIDENTIFIED"

// algorithmUnidentifiedURN is the urn for a not identified (unsupported)
// algorithm. Port of the protected static final String ALGORITHM_UNIDENTIFIED_URN.
const algorithmUnidentifiedURN = "urn:etsi:019102:algorithm:unidentified"

// AbstractAlgorithmCryptographicCheckerOverrides captures the members Java's
// AbstractAlgorithmCryptographicChecker treats virtually: the abstract
// getAlgorithm() and getNotAfter(), self-called from addAdditionalInfo().
// SignatureAlgorithmCryptographicChecker and DigestAlgorithmCryptographicChecker
// implement it and register themselves via InitAbstractAlgorithmCryptographicChecker.
type AbstractAlgorithmCryptographicCheckerOverrides interface {
	// Algorithm builds and returns the validated algorithm. Port of the
	// abstract protected XmlCryptographicAlgorithm getAlgorithm().
	Algorithm() *jaxb.XmlCryptographicAlgorithm
	// NotAfter returns time after which the used cryptographic algorithm(s)
	// is no longer considered secure. Port of the abstract protected Date
	// getNotAfter().
	NotAfter() *time.Time
}

// AbstractAlgorithmCryptographicChecker is the Go form of the abstract Java
// class AbstractAlgorithmCryptographicChecker.
type AbstractAlgorithmCryptographicChecker struct {
	*process.ChainBase[*jaxb.XmlCC]

	// validationDate is the validation time.
	validationDate time.Time

	// cryptographicSuite is the cryptographic constraint.
	cryptographicSuite policy.CryptographicSuite

	// position is the validation constraint position.
	position i18n.MessageTag

	// cryptographicAlgorithm is the verified cryptographic algorithm, cached
	// on first computation by the concrete Algorithm() implementations.
	cryptographicAlgorithm *jaxb.XmlCryptographicAlgorithm

	overrides AbstractAlgorithmCryptographicCheckerOverrides
}

// InitAbstractAlgorithmCryptographicChecker wires the shared state; called by
// the concrete constructor before InitChainBase. Port of the common
// constructor AbstractAlgorithmCryptographicChecker(I18nProvider, Date,
// MessageTag, CryptographicSuite).
func (c *AbstractAlgorithmCryptographicChecker) InitAbstractAlgorithmCryptographicChecker(i18nProvider *i18n.I18nProvider,
	validationDate time.Time, position i18n.MessageTag, cryptographicSuite policy.CryptographicSuite,
	overrides AbstractAlgorithmCryptographicCheckerOverrides) {
	xmlCC := &jaxb.XmlCC{}
	result := process.NewResult(xmlCC, &xmlCC.XmlConstraintsConclusionContent, &xmlCC.XmlConstraintsConclusionAttrs)
	c.ChainBase = process.NewChainBase(i18nProvider, result)
	c.validationDate = validationDate
	c.position = position
	c.cryptographicSuite = cryptographicSuite
	c.overrides = overrides
}

// Title returns the title of the chain (i.e. the BasicBuildingBlock title).
// Port of the overridden protected MessageTag getTitle().
func (c *AbstractAlgorithmCryptographicChecker) Title() i18n.MessageTag {
	return i18n.MessageTagCC
}

// AddAdditionalInfo adds additional info to the chain. Port of the overridden
// protected void addAdditionalInfo(). Java's super.addAdditionalInfo() call is
// the Chain default, which is empty, and is therefore not ported.
func (c *AbstractAlgorithmCryptographicChecker) AddAdditionalInfo() {
	c.Result.Value.CryptographicValidation = c.getCryptographicValidation()
}

// getCryptographicValidation builds a XmlCryptographicValidation information.
// Port of the protected XmlCryptographicValidation getCryptographicValidation().
func (c *AbstractAlgorithmCryptographicChecker) getCryptographicValidation() *jaxb.XmlCryptographicValidation {
	xmlCryptographicValidation := &jaxb.XmlCryptographicValidation{}
	xmlCryptographicValidation.Algorithm = c.overrides.Algorithm()
	if notAfter := c.overrides.NotAfter(); notAfter != nil {
		xmlCryptographicValidation.NotAfter = jaxb.NewXSDateTime(*notAfter)
	}
	xmlCryptographicValidation.Conclusion = c.Result.Conclusion()
	return xmlCryptographicValidation
}
