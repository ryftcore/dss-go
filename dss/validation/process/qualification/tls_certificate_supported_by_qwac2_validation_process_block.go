// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/qwac/sub/TLSCertificateSupportedByQWAC2ValidationProcessBlock.java (DSS 6.5.RC1).
package qualification

import (
	"time"

	"github.com/ryftcore/dss-go/dss/detailedreport/jaxb"
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/i18n"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// TLSCertificateSupportedByQWAC2ValidationProcessBlock performs validation
// process of a TLS certificate on whether it is supported by a 2-QWAC
// certificate (through the TLS Certificate Binding mechanism).
type TLSCertificateSupportedByQWAC2ValidationProcessBlock struct {
	*process.ChainBase[*jaxb.XmlValidationQWACProcess]

	// diagnosticData is the diagnostic data.
	diagnosticData *diagnostic.DiagnosticData

	// tlsCertificate is the certificate to determine qualification for.
	tlsCertificate *diagnostic.CertificateWrapper

	// tlsCertificateBasicValidationConclusion is the certificate's
	// BasicBuildingBlock's conclusion.
	tlsCertificateBasicValidationConclusion *jaxb.XmlConclusion

	// bindingSignatureBasicValidationConclusion is the Basic Validation
	// conclusion of the TLS Certificate Binding signature.
	bindingSignatureBasicValidationConclusion *jaxb.XmlConclusion

	// bindingCertificateProfile is the QWAC profile of the binding
	// certificate.
	bindingCertificateProfile enumerations.QWACProfile

	// websiteUrl is the URL of the website to validate the QWAC certificate
	// against.
	websiteUrl string
}

// NewTLSCertificateSupportedByQWAC2ValidationProcessBlock is the common
// constructor. Port of
// TLSCertificateSupportedByQWAC2ValidationProcessBlock(I18nProvider, DiagnosticData, CertificateWrapper, XmlConclusion, XmlConclusion, QWACProfile, String).
func NewTLSCertificateSupportedByQWAC2ValidationProcessBlock(i18nProvider *i18n.I18nProvider, diagnosticData *diagnostic.DiagnosticData,
	tlsCertificate *diagnostic.CertificateWrapper, tlsCertificateBasicValidationConclusion *jaxb.XmlConclusion,
	bindingSignatureBasicValidationConclusion *jaxb.XmlConclusion, bindingCertificateProfile enumerations.QWACProfile,
	websiteUrl string) *TLSCertificateSupportedByQWAC2ValidationProcessBlock {
	xmlResult := &jaxb.XmlValidationQWACProcess{}
	c := &TLSCertificateSupportedByQWAC2ValidationProcessBlock{
		ChainBase: process.NewChainBase(i18nProvider, process.NewResult(xmlResult,
			&xmlResult.XmlConstraintsConclusionContent, &xmlResult.XmlConstraintsConclusionAttrs)),
		diagnosticData:                            diagnosticData,
		tlsCertificate:                            tlsCertificate,
		tlsCertificateBasicValidationConclusion:   tlsCertificateBasicValidationConclusion,
		bindingSignatureBasicValidationConclusion: bindingSignatureBasicValidationConclusion,
		bindingCertificateProfile:                 bindingCertificateProfile,
		websiteUrl:                                websiteUrl,
	}
	c.Result.Value.Id = tlsCertificate.Id()
	c.InitChainBase(c)
	return c
}

// BuildChainTitle builds the chain title. Port of buildChainTitle().
func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) BuildChainTitle() string {
	message := i18n.MessageTag_QWAC_VALIDATION_PROFILE
	param, err := process.GetQWACValidationMessageTag(c.QWACProfile())
	if err != nil {
		panic(err)
	}
	return c.I18nProvider.GetMessage(message, param)
}

// QWACProfile gets the current QWAC profile. Port of the public
// QWACProfile getQWACProfile().
func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) QWACProfile() enumerations.QWACProfile {
	return enumerations.QWACProfile_TLS_BY_QWAC_2
}

// InitChain initializes the chain. Port of initChain().
//
// 6.2.2 Usage of 2-QWACs with TLS Certificate Binding (i.e. "Approach #2")
//
// When using 2-QWACs with secure TLS connections to websites, web browsers
// shall:
//
//  1. Establish a secure TLS connection with the site using the web
//     browsers' procedures and configuration, and evaluate the presented
//     TLS Certificate with the security requirements of the web browser
//     vendor and their policies for web security, domain authentication and
//     the encryption of web traffic as outlined in Recital 65 of the
//     Regulation (EU) 2024/1183 [i.3].
//     - If this step fails, the procedure finishes negatively.
func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) InitChain() {

	item := c.isAcceptableBuildingBlockConclusion()
	c.FirstItem = item

	item = item.SetNextItem(c.qwacDomainName())

	// 2) Examine the HTTP headers included in any main frame navigation
	// response from the server (relating to navigation by the web browser
	// to the address as displayed in the address bar) for a HTTP 'Link'
	// response header (as defined in IETF RFC 8288 [6]) with a rel value of
	// tls-certificate-binding.
	// - If this step is absent, the procedure finishes negatively.
	item = item.SetNextItem(c.tlsCertificateBindingUrlPresent())

	// 3) Fetch the resource located at this link and evaluate it for
	// conformance with the profile laid out in Annex B.
	// - If this step fails or the resource is non-conformant, the procedure
	// finishes negatively.
	item = item.SetNextItem(c.tlsCertificateBindingSignatureFound())

	if c.diagnosticData.TLSCertificateBindingSignature() != nil {

		item = item.SetNextItem(c.tlsCertificateBindingSignatureFormat())

		item = item.SetNextItem(c.tlsCertificateBindingSignatureSerializationType())

		// TODO : no verification of present headers -> ETSI TS 119 411-5
		// v1.2.1 has a sigD/crit dictionaries conflict

		// NOTE: header requirements are to be indirectly checked through the
		// validation policy on signature validation

		item = item.SetNextItem(c.tlsCertificateBindingSignatureExpProtectedHeaderPresent())

		item = item.SetNextItem(c.tlsCertificateBindingSignatureExpiryDate())

		// 4) Examine the QWAC presented in the binding with the validation
		// criteria laid out in clause 6.1.2 of the present document.
		// - If this step fails or the certificate is not considered a
		// '2-QWAC' under clause 6.1.2 of the present document, the
		// procedure finishes negatively.
		item = item.SetNextItem(c.tlsCertificateBindingIsQWAC())

		// 5) Validate the JAdES signature on the TLS Certificate binding
		// according to ETSI EN 319 102-1 [2].
		// - If this step fails or the TLS Certificate binding is not
		// considered valid, the procedure finishes negatively
		item = item.SetNextItem(c.tlsCertificateBindingSignatureValid())

		// 6) Validate that the TLS Certificate used to establish this
		// connection in Step 1 appears in the list contained in the
		// validated binding.
		// - If this step fails or the list does not contain the
		// certificate, the procedure finishes negatively.
		item = item.SetNextItem(c.tlsCertificateBindingCertificateAppearInSignature()) //nolint:staticcheck // mirrors upstream TLSCertificateSupportedByQWAC2ValidationProcessBlock#initChain: Java's trailing `item = item.setNextItem(...)` is the same dead store - setNextItem links the item and returns it, and nothing reads the tail afterwards.
	}
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) isAcceptableBuildingBlockConclusion() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewAcceptableBuildingBlockConclusionCheck(c.I18nProvider, c.Result, c.tlsCertificateBasicValidationConclusion, c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) qwacDomainName() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewQWACDomainNameCheck(c.I18nProvider, c.Result, c.tlsCertificate, c.websiteUrl, c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingUrlPresent() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingUrlPresentCheck(c.I18nProvider, c.Result, c.diagnosticData.TLSCertificateBindingUrl(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureFound() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingSignatureFoundCheck(c.I18nProvider, c.Result, c.diagnosticData.TLSCertificateBindingSignature(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureFormat() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingSignatureFormatCheck(c.I18nProvider, c.Result, c.diagnosticData.TLSCertificateBindingSignature(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureSerializationType() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingSignatureSerializationTypeCheck(c.I18nProvider, c.Result, c.diagnosticData.TLSCertificateBindingSignature(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureExpProtectedHeaderPresent() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingSignatureExpProtectedHeaderPresentCheck(c.I18nProvider, c.Result, c.diagnosticData.TLSCertificateBindingSignature(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureExpiryDate() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	var currentTime time.Time
	if vd := c.diagnosticData.ValidationDate(); vd != nil {
		currentTime = *vd
	}
	return NewTLSCertificateBindingSignatureExpiryDateCheck(c.I18nProvider, c.Result, currentTime,
		c.diagnosticData.TLSCertificateBindingSignature(), c.diagnosticData.UsedCertificates(), c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingIsQWAC() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingQWACCheck(c.I18nProvider, c.Result, c.bindingCertificateProfile, c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingSignatureValid() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingSignatureValidationResultCheck(c.I18nProvider, c.Result, c.bindingSignatureBasicValidationConclusion, c.FailLevelRule())
}

func (c *TLSCertificateSupportedByQWAC2ValidationProcessBlock) tlsCertificateBindingCertificateAppearInSignature() process.ChainItem[*jaxb.XmlValidationQWACProcess] {
	return NewTLSCertificateBindingPresentInSignatureCheck(c.I18nProvider, c.Result, c.tlsCertificate,
		c.diagnosticData.TLSCertificateBindingSignature(), c.FailLevelRule())
}
