// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/aov/cc/checks/AbstractCryptographicCheck.java (DSS 6.5.RC1).
//
// Package placement deviation: Java's eu.europa.esig.dss.validation.process.bbb.aov.cc.checks
// package is flattened into this pkg aov, per the phase 8d porter brief ("cc +
// both checks subpackages flattened").
package aov

import (
	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// AbstractCryptographicCheck is the Go form of the abstract Java class
// AbstractCryptographicCheck. A concrete check embeds it instead of
// process.ChainItemBase and registers itself with InitChainItem, the same
// pattern as bbb.AbstractMultiValuesCheckItem.
type AbstractCryptographicCheck struct {
	*process.ChainItemBase[*jaxb.XmlCC]

	// position is the validating constraint position.
	position i18n.MessageTag
}

// NewAbstractCryptographicCheck is the default constructor. Port of
// AbstractCryptographicCheck(I18nProvider, XmlCC, MessageTag, LevelRule).
func NewAbstractCryptographicCheck(i18nProvider *i18n.I18nProvider, result *process.Result[*jaxb.XmlCC],
	position i18n.MessageTag, constraint policy.LevelRule) *AbstractCryptographicCheck {
	return &AbstractCryptographicCheck{
		ChainItemBase: process.NewChainItemBase(i18nProvider, result, constraint),
		position:      position,
	}
}

// FailedIndicationForConclusion gets an Indication in case of failure. Port
// of the overridden getFailedIndicationForConclusion().
func (c *AbstractCryptographicCheck) FailedIndicationForConclusion() enumerations.Indication {
	return enumerations.Indication_INDETERMINATE
}

// FailedSubIndicationForConclusion gets a SubIndication in case of failure.
// Port of the overridden getFailedSubIndicationForConclusion().
func (c *AbstractCryptographicCheck) FailedSubIndicationForConclusion() enumerations.SubIndication {
	return enumerations.SubIndication_CRYPTO_CONSTRAINTS_FAILURE_NO_POE
}

// Name returns name for a DigestAlgorithm safely. Port of the overloaded
// protected String getName(DigestAlgorithm).
func (c *AbstractCryptographicCheck) Name(digestAlgorithm enumerations.DigestAlgorithm) string {
	if digestAlgorithm != "" {
		return digestAlgorithm.Name()
	}
	return "?"
}

// SignatureAlgorithmName returns name for a SignatureAlgorithm safely. Port
// of the overloaded protected String getName(SignatureAlgorithm), renamed
// since Go has no overloading.
func (c *AbstractCryptographicCheck) SignatureAlgorithmName(signatureAlgorithm enumerations.SignatureAlgorithm) string {
	if signatureAlgorithm != "" {
		return signatureAlgorithm.Name()
	}
	return "?"
}
