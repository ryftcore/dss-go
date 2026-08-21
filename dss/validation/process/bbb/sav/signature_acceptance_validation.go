// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/sav/SignatureAcceptanceValidation.java (DSS 6.5.RC1).
//
// UNPORTED DEPENDENCY (flagged per porter brief - do not invent): the
// contentTimestampMessageImprint() method wires
// eu.europa.esig.dss.validation.process.vpfltvd.checks.TimestampMessageImprintWithIdCheck,
// which lives in the Java package eu.europa.esig.dss.validation.process.vpfltvd,
// NOT part of the phase 8c package layout (bbb/{isc,vci,cv,fc,sav} only) and not
// present anywhere in the repository at port time - it in turn extends
// eu.europa.esig.dss.validation.process.vpftspwatsp.checks.TimestampMessageImprintCheck,
// another unported package, so this is a two-level forward dependency. This
// file assumes vpfltvd will land in a sibling package
// "github.com/utain/esig/dss/validation/process/vpfltvd" with a generic
// constructor
// NewTimestampMessageImprintWithIdCheck[T any](i18nProvider *i18n.I18nProvider,
// result *process.Result[T], timestamp *diagnostic.TimestampWrapper,
// constraint policy.LevelRule) *TimestampMessageImprintWithIdCheck[T]
// returning a process.ChainItem[T], mirroring every other check constructor in
// this port. See porter notes: this file does not build until that package
// exists.
package sav

import (
	"time"

	jaxb "github.com/utain/esig/dss/detailedreport/jaxb"
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/i18n"
	"github.com/utain/esig/dss/model/policy"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process"
	"github.com/utain/esig/dss/validation/process/vpfltvd"
)

// SignatureAcceptanceValidation is 5.2.8 Signature acceptance validation (SAV).
// This building block covers any additional verification to be performed on
// the signature itself or on the attributes of the signature ETSI EN 319 132-1.
type SignatureAcceptanceValidation struct {
	*AbstractAcceptanceValidation[*diagnostic.SignatureWrapper]

	// diagnosticData is the Diagnostic Data.
	diagnosticData *diagnostic.DiagnosticData

	// bbbs is a map of BasicBuildingBlocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks
}

// NewSignatureAcceptanceValidation is the default constructor. Port of
// SignatureAcceptanceValidation(I18nProvider, DiagnosticData, Date, SignatureWrapper, Context, Map, XmlAOV, ValidationPolicy).
func NewSignatureAcceptanceValidation(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	currentTime time.Time, signature *diagnostic.SignatureWrapper, context enumerations.Context,
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks, aovResult *jaxb.XmlAOV,
	validationPolicy policy.ValidationPolicy) *SignatureAcceptanceValidation {
	c := &SignatureAcceptanceValidation{
		AbstractAcceptanceValidation: NewAbstractAcceptanceValidation(i18nProvider, signature, currentTime,
			context, aovResult, validationPolicy),
		diagnosticData: diagnosticData,
		bbbs:           bbbs,
	}
	c.InitChainBase(c)
	return c
}

// Title returns the title of the building block. Port of getTitle().
func (c *SignatureAcceptanceValidation) Title() i18n.MessageTag {
	return i18n.MessageTag_SIGNATURE_ACCEPTANCE_VALIDATION
}

// InitChain initializes the chain. Port of initChain().
func (c *SignatureAcceptanceValidation) InitChain() {

	signatureForm, err := c.token.SignatureFormat().SignatureForm()
	if err != nil {
		panic(err)
	}

	item := c.structuralValidation()
	c.FirstItem = item

	if c.token.SigningCertificate() != nil {

		item = item.SetNextItem(c.signingCertificateAttributePresent())

		if c.token.IsSigningCertificateReferencePresent() {
			/*
			 * 5.2.8.4.2.1 Processing signing certificate reference constraint
			 *
			 * If the Signing Certificate Identifier attribute contains references to
			 * other certificates in the path, the building block shall check each of
			 * the certificates in the certification path against these references.
			 *
			 * When this property contains one or more references to certificates other than
			 * those present in the certification path, the building block shall return
			 * the indication INDETERMINATE with the sub-indication SIG_CONSTRAINTS_FAILURE.
			 */
			item = item.SetNextItem(c.unicitySigningCertificateAttribute())

			item = item.SetNextItem(c.signingCertificateReferencesValidity())

			/*
			 * When one or more certificates in the certification path are not referenced
			 * by this property, and the signature policy mandates references to all
			 * the certificates in the certification path to be present, the building block shall
			 * return the indication INDETERMINATE with the sub-indication SIG_CONSTRAINTS_FAILURE.
			 */
			item = item.SetNextItem(c.allCertificatesInPathReferenced())
		}

		// verification for JAdES / CB-AdES
		if enumerations.SignatureForm_JAdES == signatureForm || enumerations.SignatureForm_CBAdES == signatureForm {

			item = item.SetNextItem(c.keyIdentifierPresent())

			if c.token.KeyIdentifierReference() != nil {
				item = item.SetNextItem(c.keyIdentifierMatch())
			}

			item = item.SetNextItem(c.x509UrlPresent())

			if utils.IsCollectionNotEmpty(c.token.X509UrlReferences()) {
				item = item.SetNextItem(c.x509UrlMatch())
			}

		}

	}

	// signing-time
	item = item.SetNextItem(c.signingTime())

	if c.token.ClaimedSigningTime() != nil && c.token.SigningCertificate() != nil {
		item = item.SetNextItem(c.signingTimeInCertificateValidityRange())
	}

	if enumerations.SignatureForm_JAdES == signatureForm {
		item = item.SetNextItem(c.signatureType())
	}

	// content-type
	item = item.SetNextItem(c.contentType())

	// content-hints
	item = item.SetNextItem(c.contentHints())

	// message-digest for CAdES/PAdES and SignedProperties for XAdES are present
	if enumerations.SignatureForm_JAdES != signatureForm && enumerations.SignatureForm_CBAdES != signatureForm {
		item = item.SetNextItem(c.messageDigestOrSignedProperties())
	}

	// TODO content-reference

	// content-identifier
	item = item.SetNextItem(c.contentIdentifier())

	// commitment-type-indication
	item = item.SetNextItem(c.commitmentTypeIndications())

	// signer-location
	item = item.SetNextItem(c.signerLocation())

	// claimed-roles
	item = item.SetNextItem(c.claimedRoles())

	// certified-roles
	item = item.SetNextItem(c.certifiedRoles())

	// TODO signer-attributes

	// content-timestamp
	item = item.SetNextItem(c.contentTimeStamp())

	// content-timestamp
	for _, contentTimestamp := range c.token.ContentTimestamps() {

		contentTimestampBBB := c.bbbs[contentTimestamp.Id()]
		if contentTimestampBBB != nil {
			// NOTE: if TIMESTAMP validation level has been reached
			item = item.SetNextItem(c.contentTimestampBasicValidation(contentTimestamp, contentTimestampBBB.Conclusion))
		}

		item = item.SetNextItem(c.contentTimestampMessageImprint(contentTimestamp))

	}

	// counter-signature
	item = item.SetNextItem(c.counterSignature())

	// signature-time-stamp
	item = item.SetNextItem(c.signatureTimeStamp())

	// validation-data-time-stamp
	item = item.SetNextItem(c.validationDataTimeStamp())

	// validation-data-refs-only-time-stamp
	item = item.SetNextItem(c.validationDataRefsOnlyTimeStamp())

	// archive-time-stamp
	item = item.SetNextItem(c.archiveTimeStamp())

	// document-time-stamp (PAdES only)
	if enumerations.SignatureForm_PAdES == signatureForm {
		item = item.SetNextItem(c.documentTimeStamp())
	}

	// cryptographic check
	item = c.cryptographic(item) //nolint:staticcheck // mirrors upstream SignatureAcceptanceValidation#initChain: Java closes the chain with the same dead store `item = cryptographic(item);` - the helper links and returns the new tail, which nothing reads.
}

func (c *SignatureAcceptanceValidation) structuralValidation() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.StructuralValidationConstraint(c.context)
	return NewStructuralValidationCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) keyIdentifierPresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.KeyIdentifierPresent(c.context)
	return NewKeyIdentifierPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) keyIdentifierMatch() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.KeyIdentifierMatch(c.context)
	return NewKeyIdentifierMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) x509UrlPresent() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.X509UrlPresent(c.context)
	return NewX509UrlPresentCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) x509UrlMatch() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.X509UrlMatch(c.context)
	return NewX509UrlMatchCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) signingTime() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SigningTimeConstraint(c.context)
	return NewSigningTimeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) signingTimeInCertificateValidityRange() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SigningTimeInCertRangeConstraint(c.context)
	return NewSigningTimeInCertificateValidityRangeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) signatureType() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SignatureTypeConstraint(c.context)
	return NewSignatureTypeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) contentType() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ContentTypeConstraint(c.context)
	return NewContentTypeCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) contentHints() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ContentHintsConstraint(c.context)
	return NewContentHintsCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) contentIdentifier() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ContentIdentifierConstraint(c.context)
	return NewContentIdentifierCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) messageDigestOrSignedProperties() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.MessageDigestOrSignedPropertiesConstraint(c.context)
	return NewMessageDigestOrSignedPropertiesCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) commitmentTypeIndications() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.CommitmentTypeIndicationConstraint(c.context)
	return NewCommitmentTypeIndicationsCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) signerLocation() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SignerLocationConstraint(c.context)
	return NewSignerLocationCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) contentTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ContentTimeStampConstraint(c.context)
	return NewContentTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) contentTimestampBasicValidation(timestamp *diagnostic.TimestampWrapper,
	xmlConclusion *jaxb.XmlConclusion) process.ChainItem[*jaxb.XmlSAV] {
	return NewContentTimestampBasicValidationCheck(c.I18nProvider, c.Result, timestamp, xmlConclusion,
		c.getTimestampBasicValidationConstraintLevel())
}

func (c *SignatureAcceptanceValidation) contentTimestampMessageImprint(contentTimestamp *diagnostic.TimestampWrapper) process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ContentTimeStampMessageImprintConstraint(c.context)
	return vpfltvd.NewTimestampMessageImprintWithIdCheck(c.I18nProvider, c.Result, contentTimestamp, constraint)
}

func (c *SignatureAcceptanceValidation) claimedRoles() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ClaimedRoleConstraint(c.context)
	return NewClaimedRolesCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) certifiedRoles() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.CertifiedRolesConstraint(c.context)
	return NewCertifiedRolesCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) counterSignature() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.CounterSignatureConstraint(c.context)
	return NewCounterSignatureCheck(c.I18nProvider, c.Result, c.diagnosticData, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) signatureTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.SignatureTimeStampConstraint(c.context)
	return NewSignatureTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) validationDataTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ValidationDataTimeStampConstraint(c.context)
	return NewValidationDataTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) validationDataRefsOnlyTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ValidationDataRefsOnlyTimeStampConstraint(c.context)
	return NewValidationDataRefsOnlyTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) archiveTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.ArchiveTimeStampConstraint(c.context)
	return NewArchiveTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

func (c *SignatureAcceptanceValidation) documentTimeStamp() process.ChainItem[*jaxb.XmlSAV] {
	constraint := c.validationPolicy.DocumentTimeStampConstraint(c.context)
	return NewDocumentTimeStampCheck(c.I18nProvider, c.Result, c.token, constraint)
}

// getTimestampBasicValidationConstraintLevel ports the private
// getTimestampBasicValidationConstraintLevel().
func (c *SignatureAcceptanceValidation) getTimestampBasicValidationConstraintLevel() policy.LevelRule {
	constraint := c.validationPolicy.TimestampValidConstraint()
	// continue if LTA is present
	if constraint == nil || process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.token) {
		constraint = c.WarnLevelRule()
	}
	return constraint
}

// CollectMessages collects required messages from the given constraint to the
// given conclusion. Port of collectMessages(XmlConclusion, XmlConstraint).
func (c *SignatureAcceptanceValidation) CollectMessages(conclusion *jaxb.XmlConclusion, constraint *jaxb.XmlConstraint) {
	if constraintBlockType(constraint) == jaxb.XmlBlockType_TST_BBB &&
		(c.validationPolicy.TimestampValidConstraint() == nil ||
			process.IsLongTermAvailabilityAndIntegrityMaterialPresent(c.token)) {
		// skip validation messages for content TSTs
	} else {
		c.ChainBase.CollectMessages(conclusion, constraint)
	}
}

// constraintBlockType reads XmlConstraint#getBlockType(): the generated member
// is a pointer, whose nil is Java's null.
func constraintBlockType(constraint *jaxb.XmlConstraint) jaxb.XmlBlockType {
	if constraint.BlockType == nil {
		return ""
	}
	return *constraint.BlockType
}
