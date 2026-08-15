// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/AbstractAcceptanceValidation.java (DSS 6.5.RC1).
//
// Integration note (phase 8c pass): cryptographic()/algorithmObsolescenceValidationCheck()
// need eu.europa.esig.dss.validation.process.bbb.aov.checks.AlgorithmObsolescenceValidationCheck,
// which the porter flagged as unported (bbb/aov is not part of the phase 8c
// package layout). That one check class - the sole consumer of an
// already-built XmlAOV result, with no dependency on the rest of the aov
// tree - has since been ported minimally into
// github.com/utain/esig/dss/validation/process/bbb/aov (see that package's
// header for the exact scope), so this file now builds as originally
// written.
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/bbb/aov"
)

// AbstractAcceptanceValidation is 5.2.8 Signature acceptance validation (SAV).
// This building block covers any additional verification to be performed on
// the signature itself or on the attributes of the signature ETSI EN 319 132-1.
// TK is the validation token wrapper type.
type AbstractAcceptanceValidation[TK diagnostic.TokenProxy] struct {
	*process.ChainBase[*jaxb.XmlSAV]

	// token is the token to be validated.
	token TK

	// currentTime is the validation time.
	currentTime time.Time

	// context is the validation context.
	context enumerations.Context

	// aov is the result of the Algorithm Obsolescence Validation block.
	aov *jaxb.XmlAOV

	// validationPolicy is the validation policy.
	validationPolicy policy.ValidationPolicy
}

// NewAbstractAcceptanceValidation is the default constructor. Port of
// AbstractAcceptanceValidation(I18nProvider, T, Date, Context, XmlAOV, ValidationPolicy).
func NewAbstractAcceptanceValidation[TK diagnostic.TokenProxy](i18nProvider *i18n.I18nProvider, token TK,
	currentTime time.Time, context enumerations.Context, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *AbstractAcceptanceValidation[TK] {
	xmlSAV := &jaxb.XmlSAV{}
	return &AbstractAcceptanceValidation[TK]{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlSAV,
			&xmlSAV.XmlConstraintsConclusionContent, &xmlSAV.XmlConstraintsConclusionAttrs)),
		token:            token,
		currentTime:      currentTime,
		context:          context,
		aov:              aovResult,
		validationPolicy: validationPolicy,
	}
}

// signingCertificateAttributePresent checks whether a signing-certificate
// signed attribute is present. Port of signingCertificateAttributePresent().
func (c *AbstractAcceptanceValidation[TK]) signingCertificateAttributePresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SigningCertificateAttributePresentConstraint(c.context)
	return NewSigningCertificateAttributePresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// unicitySigningCertificateAttribute checks if only one signing-certificate
// signed attribute is present. Port of unicitySigningCertificateAttribute().
func (c *AbstractAcceptanceValidation[TK]) unicitySigningCertificateAttribute() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.UnicitySigningCertificateAttributeConstraint(c.context)
	return NewUnicitySigningCertificateAttributeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// signingCertificateReferencesValidity checks whether a signing-certificate
// signed attribute is valid to the determined signing certificate. Port of
// signingCertificateReferencesValidity().
func (c *AbstractAcceptanceValidation[TK]) signingCertificateReferencesValidity() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SigningCertificateRefersCertificateChainConstraint(c.context)
	return NewSigningCertificateReferencesValidityCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// allCertificatesInPathReferenced checks if all certificates in a signing
// certificate chain are references within signing-certificate signed
// attribute. Port of allCertificatesInPathReferenced().
func (c *AbstractAcceptanceValidation[TK]) allCertificatesInPathReferenced() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ReferencesToAllCertificateChainPresentConstraint(c.context)
	return NewAllCertificatesInPathReferencedCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// cryptographic verifies cryptographic validity of signature references and
// signing-certificate signed attribute. item is the last initialized chain
// item to be processed. Port of cryptographic(ChainItem).
func (c *AbstractAcceptanceValidation[TK]) cryptographic(item process.ChainItem[*jaxb.XmlSAV]) process.ChainItem[*jaxb.XmlSAV] {
	// The basic signature constraints validation
	position, err := process.GetCryptoPosition(c.context)
	if err != nil {
		panic(err)
	}

	if item == nil {
		item = c.algorithmObsolescenceValidationCheck(c.Result, c.aov, position)
		c.FirstItem = item
	} else {
		item = item.SetNextItem(c.algorithmObsolescenceValidationCheck(c.Result, c.aov, position))
	}

	return item
}

// algorithmObsolescenceValidationCheck verifies the result of the Algorithm
// Obsolescence Validation building block. Port of
// algorithmObsolescenceValidationCheck(XmlSAV, XmlAOV, MessageTag).
func (c *AbstractAcceptanceValidation[TK]) algorithmObsolescenceValidationCheck(result *process.Result[*jaxb.XmlSAV],
	aovResult *jaxb.XmlAOV, position i18n.MessageTag) process.ChainItem[*jaxb.XmlSAV] {
	return aov.NewAlgorithmObsolescenceValidationCheck(c.I18nProvider, result, aovResult, c.currentTime, position, c.token.Id())
}
