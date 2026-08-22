// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/bbb/BasicBuildingBlocks.java (DSS 6.5.RC1).
//
// STRUCTURAL FIX A (phase 8d AOV porter, per S8D_BRIEF.md): this dispatcher
// instantiates every one of the 5.2.x building blocks - isc, vci, cv (siblings
// under bbb, ported in phase 8c), fc, sav (also bbb siblings), and aov, xcv
// (ported in phase 8d). Every one of fc/sav/isc/vci/cv/aov/xcv itself imports
// package bbb for the shared AbstractValueCheckItem/AbstractMultiValuesCheckItem/
// AbstractCertificateCheckItem base types (bbb/abstract_*_check_item.go), so a
// dispatcher importing all of them cannot itself live in package bbb - Go
// (unlike Java, which resolves per-class) rejects the resulting import cycle:
// bbb -> {fc,sav,...} -> bbb. This file therefore moves out of
// dss/validation/process/bbb into this new, sibling package
// dss/validation/process/blocks, breaking the cycle. It was previously gated
// behind a //go:build phase8d tag for the same reason (before the aov/xcv
// packages existed to import); now that they do, this file is un-tagged and
// builds for real. See basic_building_blocks_eaa.go / _noeaa.go for the
// companion STRUCTURAL FIX B (tag composition for the EAA-dependent calls
// still gated behind fc/sav's own "eaa" build tag).
package blocks

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/model/policy"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/aov"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/cv"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/fc"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/isc"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/sav"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/vci"
	"github.com/ryftcore/dss-go/dss/validation/process/bbb/xcv"
)

// aovBlock is the Go form of Java's AlgorithmObsolescenceValidation<?> local.
type aovBlock interface {
	Execute() *jaxb.XmlAOV
}

// savBlock is the Go form of Java's AbstractAcceptanceValidation<?> local.
type savBlock interface {
	Execute() *jaxb.XmlSAV
}

// BasicBuildingBlocks is 5.2 Basic building blocks.
type BasicBuildingBlocks struct {
	// i18nProvider is the internationalization provider.
	i18nProvider *i18n.I18nProvider

	// diagnosticData is the Diagnostic Data.
	diagnosticData *diagnostic.DiagnosticData

	// token is the validating token.
	token diagnostic.TokenProxy

	// policy is the validation policy.
	policy policy.ValidationPolicy

	// bbbs is the map of basic building blocks.
	bbbs map[string]*jaxb.XmlBasicBuildingBlocks

	// currentTime is the validation time.
	currentTime time.Time

	// context is the validation context.
	context enumerations.Context
}

// NewBasicBuildingBlocks is the default constructor. Port of
// BasicBuildingBlocks(I18nProvider, DiagnosticData, TokenProxy, Date, Map, ValidationPolicy, Context).
func NewBasicBuildingBlocks(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	token diagnostic.TokenProxy, currentTime time.Time, bbbs map[string]*jaxb.XmlBasicBuildingBlocks,
	validationPolicy policy.ValidationPolicy, context enumerations.Context) *BasicBuildingBlocks {
	return &BasicBuildingBlocks{
		i18nProvider:   i18nProvider,
		diagnosticData: diagnosticData,
		token:          token,
		currentTime:    currentTime,
		bbbs:           bbbs,
		policy:         validationPolicy,
		context:        context,
	}
}

// Execute executes the 5.2 token validation process. Port of execute().
func (b *BasicBuildingBlocks) Execute() *jaxb.XmlBasicBuildingBlocks {
	result := &jaxb.XmlBasicBuildingBlocks{}
	result.Id = b.token.Id()
	result.Type = jaxb.ContextValue(b.context)
	result.Conclusion = &jaxb.XmlConclusion{}

	/*
	 * 5.2.2 Format Checking
	 */
	xmlFC := b.executeFormatChecking()
	if xmlFC != nil {
		result.FC = xmlFC
		b.updateFinalConclusion(result, xmlFC.Conclusion)
	}

	/*
	 * 5.2.3 Identification of the signing certificate
	 */
	xmlISC := b.executeIdentificationOfTheSigningCertificate()
	if xmlISC != nil {
		result.ISC = xmlISC
		result.CertificateChain = xmlISC.CertificateChain
		b.updateFinalConclusion(result, xmlISC.Conclusion)
	}

	/*
	 * 5.2.4 Validation context initialization (only for signature)
	 */
	xmlVCI := b.executeValidationContextInitialization()
	if xmlVCI != nil {
		result.VCI = xmlVCI
		b.updateFinalConclusion(result, xmlVCI.Conclusion)
	}

	/*
	 * Algorithm Obsolescence Validation
	 * NOTE: out of standard, but executed to receive a harmonized input for 5.2.6 and 5.2.8 building blocks
	 */
	xmlAOV := b.executeAlgorithmObsolescenceValidation()
	if xmlAOV != nil {
		result.AOV = xmlAOV
	}

	/*
	 * 5.2.6 X.509 certificate validation
	 */
	xmlXCV := b.executeX509CertificateValidation(xmlAOV)
	if xmlXCV != nil {
		result.XCV = xmlXCV
		b.addAdditionalInfo(xmlXCV)
		b.updateFinalConclusion(result, xmlXCV.Conclusion)
	}

	/*
	 * 5.2.7 Cryptographic verification
	 */
	xmlCV := b.executeCryptographicVerification()
	if xmlCV != nil {
		result.CV = xmlCV
		b.updateFinalConclusion(result, xmlCV.Conclusion)
	}

	/*
	 * 5.2.8 Signature acceptance validation (SAV)
	 */
	xmlSAV := b.executeSignatureAcceptanceValidation(xmlAOV)
	if xmlSAV != nil {
		result.SAV = xmlSAV
		b.updateFinalConclusion(result, xmlSAV.Conclusion)
	}

	if result.Conclusion.Indication.Indication() == "" {
		result.Conclusion.Indication = jaxb.IndicationValue(enumerations.Indication_PASSED)
	}

	return result
}

// updateFinalConclusion ports the private
// updateFinalConclusion(XmlBasicBuildingBlocks, XmlConstraintsConclusion): the
// Java parameter is the generated base class, of which only the conclusion is
// read, so the conclusion is passed directly here.
func (b *BasicBuildingBlocks) updateFinalConclusion(result *jaxb.XmlBasicBuildingBlocks,
	currentConclusion *jaxb.XmlConclusion) {
	finalConclusion := result.Conclusion

	if enumerations.Indication_PASSED != currentConclusion.Indication.Indication() {
		finalConclusion.Indication = currentConclusion.Indication
		finalConclusion.SubIndication = currentConclusion.SubIndication
		finalConclusion.Errors = append(finalConclusion.Errors, currentConclusion.Errors...)
	}
	finalConclusion.Warnings = append(finalConclusion.Warnings, currentConclusion.Warnings...)
	finalConclusion.Infos = append(finalConclusion.Infos, currentConclusion.Infos...)
}

// executeFormatChecking ports the private executeFormatChecking().
func (b *BasicBuildingBlocks) executeFormatChecking() *jaxb.XmlFC {
	if enumerations.Context_SIGNATURE == b.context || enumerations.Context_COUNTER_SIGNATURE == b.context ||
		enumerations.Context_KEY_BINDING_SIGNATURE == b.context {
		block := fc.NewSignatureFormatChecking(b.i18nProvider, b.diagnosticData,
			b.token.(*diagnostic.SignatureWrapper), b.context, b.policy)
		return block.Execute()
	} else if enumerations.Context_TIMESTAMP == b.context {
		block := fc.NewTimestampFormatChecking(b.i18nProvider, b.diagnosticData,
			b.token.(*diagnostic.TimestampWrapper), b.context, b.policy)
		xmlFC := block.Execute()
		if utils.IsCollectionNotEmpty(xmlFC.Constraint) {
			return xmlFC
		}
	} else if enumerations.Context_EAA == b.context || enumerations.Context_EAA_REVOCATION == b.context {
		return b.executeEAAFormatChecking()
	}
	return nil
}

// executeIdentificationOfTheSigningCertificate ports the private
// executeIdentificationOfTheSigningCertificate().
func (b *BasicBuildingBlocks) executeIdentificationOfTheSigningCertificate() *jaxb.XmlISC {
	if enumerations.Context_CERTIFICATE != b.context && enumerations.Context_EAA != b.context {
		block := isc.NewIdentificationOfTheSigningCertificate(b.i18nProvider, b.token, b.context, b.policy)
		return block.Execute()
	} else {
		return nil
	}
}

// executeValidationContextInitialization ports the private
// executeValidationContextInitialization().
func (b *BasicBuildingBlocks) executeValidationContextInitialization() *jaxb.XmlVCI {
	if enumerations.Context_SIGNATURE == b.context || enumerations.Context_COUNTER_SIGNATURE == b.context ||
		enumerations.Context_KEY_BINDING_SIGNATURE == b.context {
		block := vci.NewValidationContextInitialization(b.i18nProvider,
			b.token.(*diagnostic.SignatureWrapper), b.context, b.policy)
		return block.Execute()
	}
	return nil
}

// executeAlgorithmObsolescenceValidation ports the private
// executeAlgorithmObsolescenceValidation().
func (b *BasicBuildingBlocks) executeAlgorithmObsolescenceValidation() *jaxb.XmlAOV {
	var block aovBlock
	if enumerations.Context_SIGNATURE == b.context || enumerations.Context_COUNTER_SIGNATURE == b.context ||
		enumerations.Context_KEY_BINDING_SIGNATURE == b.context {
		block = aov.NewSignatureAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.SignatureWrapper), b.context, b.currentTime, b.policy)
	} else if enumerations.Context_TIMESTAMP == b.context {
		block = aov.NewTimestampAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.TimestampWrapper), b.currentTime, b.policy)
	} else if enumerations.Context_REVOCATION == b.context {
		block = aov.NewRevocationDataAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.RevocationWrapper), b.currentTime, b.policy)
	} else if enumerations.Context_CERTIFICATE == b.context {
		block = aov.NewCertificateAndChainAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.CertificateWrapper), b.context, b.currentTime, b.policy)
	} else if enumerations.Context_EAA == b.context {
		block = aov.NewEAAAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.EAAWrapper), b.currentTime, b.policy)
	} else if enumerations.Context_EAA_REVOCATION == b.context {
		block = aov.NewEAARevocationAlgorithmObsolescenceValidation(
			b.i18nProvider, b.token.(*diagnostic.EAARevocationTokenWrapper), b.currentTime, b.policy)
	}
	if block != nil {
		return block.Execute()
	}
	return nil
}

// executeX509CertificateValidation ports the private
// executeX509CertificateValidation(XmlAOV).
func (b *BasicBuildingBlocks) executeX509CertificateValidation(xmlAOV *jaxb.XmlAOV) *jaxb.XmlXCV {
	x509CertificateValidation := b.getX509CertificateValidation(xmlAOV)
	if x509CertificateValidation != nil {
		return x509CertificateValidation.Execute()
	}
	return nil
}

// getX509CertificateValidation ports the private getX509CertificateValidation(XmlAOV).
func (b *BasicBuildingBlocks) getX509CertificateValidation(xmlAOV *jaxb.XmlAOV) *xcv.X509CertificateValidation {
	if enumerations.Context_CERTIFICATE == b.context {
		certificate := b.token.(*diagnostic.CertificateWrapper)
		return xcv.NewX509CertificateValidation(b.i18nProvider, certificate, b.currentTime,
			certificate.NotBefore(), b.context, xmlAOV, b.policy)

	} else {
		signingCertificate := b.token.SigningCertificate()
		if signingCertificate != nil {
			if enumerations.Context_SIGNATURE == b.context || enumerations.Context_COUNTER_SIGNATURE == b.context ||
				enumerations.Context_KEY_BINDING_SIGNATURE == b.context {
				return xcv.NewX509CertificateValidation(b.i18nProvider, signingCertificate, b.currentTime,
					signingCertificate.NotBefore(), b.context, xmlAOV, b.policy)

			} else if enumerations.Context_TIMESTAMP == b.context {
				return xcv.NewX509CertificateValidation(b.i18nProvider, signingCertificate, b.currentTime,
					b.token.(*diagnostic.TimestampWrapper).ProductionTime(), b.context, xmlAOV, b.policy)

			} else if enumerations.Context_REVOCATION == b.context {
				return xcv.NewX509CertificateValidation(b.i18nProvider, signingCertificate, b.currentTime,
					b.token.(*diagnostic.RevocationWrapper).ProductionDate(), b.context, xmlAOV, b.policy)

			} else if enumerations.Context_EAA_REVOCATION == b.context {
				return xcv.NewX509CertificateValidation(b.i18nProvider, signingCertificate, b.currentTime,
					b.token.(*diagnostic.EAARevocationTokenWrapper).IssuedAt(), b.context, xmlAOV, b.policy)
			}
		}
	}
	return nil
}

// addAdditionalInfo ports the private addAdditionalInfo(XmlXCV).
func (b *BasicBuildingBlocks) addAdditionalInfo(xmlXCV *jaxb.XmlXCV) {
	for _, subXCV := range xmlXCV.SubXCV {
		cert := b.diagnosticData.UsedCertificateById(subXCV.Id)
		crossCertificates := b.diagnosticData.CrossCertificates(cert)
		if utils.IsCollectionNotEmpty(crossCertificates) {
			addCrossCertificates(subXCV, getCertificateWrapperIds(crossCertificates))
		}
		orphanCrossCertificates := b.diagnosticData.OrphanCrossCertificates(cert)
		if utils.IsCollectionNotEmpty(orphanCrossCertificates) {
			addCrossCertificates(subXCV, getOrphanCertificateWrapperIds(orphanCrossCertificates))
		}

		equivalentCertificates := removeAllCertificates(b.diagnosticData.EquivalentCertificates(cert), crossCertificates)
		if utils.IsCollectionNotEmpty(equivalentCertificates) {
			addEquivalentCertificates(subXCV, getCertificateWrapperIds(equivalentCertificates))
		}
		orphanEquivalentCertificates := removeAllOrphanCertificates(
			b.diagnosticData.OrphanEquivalentCertificates(cert), orphanCrossCertificates)
		if utils.IsCollectionNotEmpty(orphanEquivalentCertificates) {
			addEquivalentCertificates(subXCV, getOrphanCertificateWrapperIds(orphanEquivalentCertificates))
		}
	}
}

// addCrossCertificates ports subXCV.getCrossCertificates().addAll(ids); the
// generated member is a StringList wrapper, created on first use.
func addCrossCertificates(subXCV *jaxb.XmlSubXCV, ids []string) {
	if subXCV.CrossCertificate == nil {
		subXCV.CrossCertificate = &jaxb.StringList{}
	}
	*subXCV.CrossCertificate = append(*subXCV.CrossCertificate, ids...)
}

// addEquivalentCertificates ports subXCV.getEquivalentCertificates().addAll(ids).
func addEquivalentCertificates(subXCV *jaxb.XmlSubXCV, ids []string) {
	if subXCV.EquivalentCertificate == nil {
		subXCV.EquivalentCertificate = &jaxb.StringList{}
	}
	*subXCV.EquivalentCertificate = append(*subXCV.EquivalentCertificate, ids...)
}

// removeAllCertificates ports List#removeAll(Collection) over certificate
// wrappers, whose equals() compares the wrapper type and the token Id.
func removeAllCertificates(certificates []*diagnostic.CertificateWrapper,
	toRemove []*diagnostic.CertificateWrapper) []*diagnostic.CertificateWrapper {
	var result []*diagnostic.CertificateWrapper
	for _, certificate := range certificates {
		found := false
		for _, removed := range toRemove {
			if certificate.Equals(removed) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, certificate)
		}
	}
	return result
}

// removeAllOrphanCertificates ports List#removeAll(Collection) over orphan
// certificate token wrappers.
func removeAllOrphanCertificates(certificates []*diagnostic.OrphanCertificateTokenWrapper,
	toRemove []*diagnostic.OrphanCertificateTokenWrapper) []*diagnostic.OrphanCertificateTokenWrapper {
	var result []*diagnostic.OrphanCertificateTokenWrapper
	for _, certificate := range certificates {
		found := false
		for _, removed := range toRemove {
			if certificate.Id() == removed.Id() {
				found = true
				break
			}
		}
		if !found {
			result = append(result, certificate)
		}
	}
	return result
}

// getCertificateWrapperIds returns a list of token ids. Port of the private
// static getCertificateWrapperIds(Collection).
func getCertificateWrapperIds(tokens []*diagnostic.CertificateWrapper) []string {
	ids := make([]string, 0, len(tokens))
	for _, token := range tokens {
		ids = append(ids, token.Id())
	}
	return ids
}

// getOrphanCertificateWrapperIds returns a list of orphan token ids. Port of the
// private static getOrphanCertificateWrapperIds(Collection).
func getOrphanCertificateWrapperIds(tokens []*diagnostic.OrphanCertificateTokenWrapper) []string {
	ids := make([]string, 0, len(tokens))
	for _, token := range tokens {
		ids = append(ids, token.Id())
	}
	return ids
}

// executeCryptographicVerification ports the private executeCryptographicVerification().
func (b *BasicBuildingBlocks) executeCryptographicVerification() *jaxb.XmlCV {
	if enumerations.Context_CERTIFICATE != b.context {
		block := cv.NewCryptographicVerification(b.i18nProvider, b.diagnosticData, b.token, b.context, b.policy)
		return block.Execute()
	} else {
		return nil
	}
}

// executeSignatureAcceptanceValidation ports the private
// executeSignatureAcceptanceValidation(XmlAOV).
func (b *BasicBuildingBlocks) executeSignatureAcceptanceValidation(xmlAOV *jaxb.XmlAOV) *jaxb.XmlSAV {
	var block savBlock
	if enumerations.Context_SIGNATURE == b.context || enumerations.Context_COUNTER_SIGNATURE == b.context ||
		enumerations.Context_KEY_BINDING_SIGNATURE == b.context {
		block = sav.NewSignatureAcceptanceValidation(b.i18nProvider, b.diagnosticData, b.currentTime,
			b.token.(*diagnostic.SignatureWrapper), b.context, b.bbbs, xmlAOV, b.policy)
	} else if enumerations.Context_TIMESTAMP == b.context {
		block = sav.NewTimestampAcceptanceValidation(b.i18nProvider, b.currentTime,
			b.token.(*diagnostic.TimestampWrapper), xmlAOV, b.policy)
	} else if enumerations.Context_REVOCATION == b.context {
		block = sav.NewRevocationAcceptanceValidation(b.i18nProvider, b.currentTime,
			b.token.(*diagnostic.RevocationWrapper), xmlAOV, b.policy)
	} else if enumerations.Context_EAA == b.context || enumerations.Context_EAA_REVOCATION == b.context {
		block = b.eaaAcceptanceValidationBlock(xmlAOV)
	}
	if block != nil {
		return block.Execute()
	}
	return nil
}
