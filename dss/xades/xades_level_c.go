// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/XAdESLevelC.java (DSS 6.5.RC1).
//
// Java extends XAdESLevelBaselineT and overrides extendSignatures(List); the Go port embeds the
// -T level, and "super.extendSignatures(signatures)" is the explicit
// c.XAdESLevelBaselineT.ExtendSignatures call. XAdESLevelX embeds this type in turn.
//
// BouncyCastle's BasicOCSPResp/RespID are replaced by their spi counterparts (PORTING.md);
// XMLGregorianCalendar.toXMLFormat() is folded into DomUtilsCreateXMLGregorianCalendar, which
// already returns the lexical form every upstream call site immediately asks for.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	xmlutils "github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESLevelC contains the XAdES-C profile aspects.
type XAdESLevelC struct {
	XAdESLevelBaselineT
}

// NewXAdESLevelC is the default constructor for XAdESLevelC.
// Port of XAdESLevelC(CertificateVerifier).
func NewXAdESLevelC(certificateVerifier validation.CertificateVerifier) *XAdESLevelC {
	extension := &XAdESLevelC{}
	extension.InitXAdESLevelC(extension, certificateVerifier)
	return extension
}

// InitXAdESLevelC registers the concrete extension level with this base and forwards to the -T
// level. Port of the super(certificateVerifier) call of XAdESLevelC(CertificateVerifier).
func (c *XAdESLevelC) InitXAdESLevelC(self XAdESSignatureExtensionOverrides,
	certificateVerifier validation.CertificateVerifier) {
	c.InitXAdESLevelBaselineT(self, certificateVerifier)
}

// ExtendSignatures builds up taking the XAdES-T signature and incorporating additional data
// required for validation: the sequence of references to the full set of CA certificates that
// have been used to validate the electronic signature up to (but not including) the signer's
// certificate, and a full set of references to the revocation data that have been used in the
// validation of the signer and CA certificates. Adds <CompleteCertificateRefs> and
// <CompleteRevocationRefs> segments into the <UnsignedSignatureProperties> element.
//
// There SHALL be at most one occurrence of the CompleteRevocationRefs and
// CompleteCertificateRefs properties in the signature; old references must be removed.
// Port of the overridden protected #extendSignatures(List).
func (c *XAdESLevelC) ExtendSignatures(signatures []validation.AdvancedSignature) error {
	if err := c.XAdESLevelBaselineT.ExtendSignatures(signatures); err != nil {
		return err
	}

	signaturesToExtend := c.extendToCLevelSignatures(signatures)
	if utils.IsCollectionEmpty(signaturesToExtend) {
		return nil
	}

	// Reset sources
	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := c.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}

		// Data sources can already be loaded in memory (force reload)
		c.XadesSignature.ResetCertificateSource()
		c.XadesSignature.ResetRevocationSources()
		c.XadesSignature.ResetTimestampSource()
	}

	signatureRequirementsChecker := c.SignatureRequirementsChecker()
	if enumerations.SignatureLevel_XAdES_C == c.Params.SignatureLevel() {
		signatureRequirementsChecker.AssertExtendToCLevelPossible(signaturesToExtend)
	}
	signatureRequirementsChecker.AssertSignaturesValid(signaturesToExtend)
	signatureRequirementsChecker.AssertCertificateChainValidForCLevel(signaturesToExtend)

	// Perform signature validation
	validationDataContainer, err := c.DocumentAnalyzer.GetValidationData(signaturesToExtend)
	if err != nil {
		return err
	}

	// Append ValidationData
	for _, signature := range signaturesToExtend {
		xadesSignature, ok := signature.(*XAdESSignature)
		if !ok {
			// Java's (XAdESSignature) cast; a non-XAdES signature would raise a ClassCastException.
			return fmt.Errorf("unexpected signature type %T", signature)
		}
		if _, err := c.InitializeSignatureBuilder(xadesSignature); err != nil {
			return err
		}
		if signatureRequirementsChecker.HasXLevelOrHigher(signature) {
			// Unable to extend due to higher levels covering the current C-level
			continue
		}

		indent, err := c.removeOldCertificateRefs()
		if err != nil {
			return err
		}
		if err := c.removeOldRevocationRefs(); err != nil {
			return err
		}

		validationDataForInclusion := c.validationDataForCLevelInclusion(validationDataContainer, signature)

		levelTUnsignedProperties := c.UnsignedSignaturePropertiesDom.Clone(true)

		// XAdES-C: complete certificate references
		// <xades:CompleteCertificateRefs>
		// ...<xades:CertRefs>
		// ......<xades:Cert>
		// .........<xades:CertDigest>
		if err := c.incorporateCertificateRefs(c.UnsignedSignaturePropertiesDom,
			validationDataForInclusion.CertificateTokens(), indent); err != nil {
			return err
		}

		// XAdES-C: complete revocation references
		// <xades:CompleteRevocationRefs>
		if utils.IsCollectionNotEmpty(validationDataForInclusion.CrlTokens()) ||
			utils.IsCollectionNotEmpty(validationDataForInclusion.OcspTokens()) {
			currentElements, err := c.CurrentXAdESElements()
			if err != nil {
				return err
			}
			completeRevocationRefsDom := xmlutils.DomUtilsAddElement(c.DocumentDom,
				c.UnsignedSignaturePropertiesDom, c.XadesNamespace(),
				currentElements.ElementCompleteRevocationRefs())
			if err := c.incorporateCRLRefs(completeRevocationRefsDom,
				validationDataForInclusion.CrlTokens()); err != nil {
				return err
			}
			if err := c.incorporateOCSPRefs(completeRevocationRefsDom,
				validationDataForInclusion.OcspTokens()); err != nil {
				return err
			}
		}

		indented, err := c.IndentIfPrettyPrint(c.UnsignedSignaturePropertiesDom, levelTUnsignedProperties)
		if err != nil {
			return err
		}
		c.UnsignedSignaturePropertiesDom = indented
	}
	return nil
}

// extendToCLevelSignatures ports the private getExtendToCLevelSignatures.
func (c *XAdESLevelC) extendToCLevelSignatures(
	signatures []validation.AdvancedSignature) []validation.AdvancedSignature {
	signaturesToExtend := make([]validation.AdvancedSignature, 0)
	for _, signature := range signatures {
		if c.cLevelExtensionRequired(signature) {
			signaturesToExtend = append(signaturesToExtend, signature)
		}
	}
	return signaturesToExtend
}

// cLevelExtensionRequired ports the private cLevelExtensionRequired.
func (c *XAdESLevelC) cLevelExtensionRequired(signature validation.AdvancedSignature) bool {
	return enumerations.SignatureLevel_XAdES_C == c.Params.SignatureLevel() ||
		enumerations.SignatureLevel_XAdES_XL == c.Params.SignatureLevel() ||
		!signature.HasXProfile()
}

// removeOldCertificateRefs ports the private removeOldCertificateRefs, including its exact
// two-step removal: the V1 element is removed first and the V2 element is removed whenever it
// exists or the V1 removal yielded no indent.
func (c *XAdESLevelC) removeOldCertificateRefs() (string, error) {
	text := ""
	certRefs, err := xmlutils.XPathUtilsGetElement(c.XadesSignature.SignatureElement(),
		c.XadesPath.CompleteCertificateRefsPath())
	if err != nil {
		return "", err
	}
	certRefsV2, err := xmlutils.XPathUtilsGetElement(c.XadesSignature.SignatureElement(),
		c.XadesPath.CompleteCertificateRefsV2Path())
	if err != nil {
		return "", err
	}
	if certRefs != nil || certRefsV2 != nil {
		text = c.RemoveNode(certRefs)
		if text == "" || certRefsV2 != nil {
			text = c.RemoveNode(certRefsV2)
		}
		/* Because the element was removed, the certificate source needs to be reset */
		c.XadesSignature.ResetCertificateSource()
	}
	return text, nil
}

// removeOldRevocationRefs ports the private removeOldRevocationRefs.
func (c *XAdESLevelC) removeOldRevocationRefs() error {
	toRemove, err := xmlutils.XPathUtilsGetElement(c.XadesSignature.SignatureElement(),
		c.XadesPath.CompleteRevocationRefsPath())
	if err != nil {
		return err
	}
	if toRemove != nil {
		c.RemoveNode(toRemove)
		/* Because the element was removed, the revocation sources need to be reset */
		c.XadesSignature.ResetRevocationSources()
	}
	return nil
}

// incorporateCertificateRefs ports the private incorporateCertificateRefs. Upstream carries a
// "TODO : review indent usage" here: the indent argument is accepted but never applied.
func (c *XAdESLevelC) incorporateCertificateRefs(parentDom *xmldom.Node,
	certificatesToBeAdded []*model.CertificateToken, indent string) error {
	_ = indent
	if utils.IsCollectionNotEmpty(certificatesToBeAdded) {
		completeCertificateRefsDom, err := c.createCompleteCertificateRefsDom(parentDom)
		if err != nil {
			return err
		}
		certRefsDom, err := c.createCertRefsDom(completeCertificateRefsDom)
		if err != nil {
			return err
		}

		tokenReferencesDigestAlgorithm := c.Params.TokenReferencesDigestAlgorithm()
		for _, certificateToken := range certificatesToBeAdded {
			if _, err := c.IncorporateCert(certRefsDom, certificateToken,
				tokenReferencesDigestAlgorithm); err != nil {
				return err
			}
		}
	}
	return nil
}

// createCompleteCertificateRefsDom ports the private createCompleteCertificateRefsDom.
func (c *XAdESLevelC) createCompleteCertificateRefsDom(parentDom *xmldom.Node) (*xmldom.Node, error) {
	if c.Params.IsEn319132() {
		return xmlutils.DomUtilsAddElement(c.DocumentDom, parentDom, c.Xades141Namespace(),
			definition.XAdES141Element_COMPLETE_CERTIFICATE_REFS_V2), nil
	}
	currentElements, err := c.CurrentXAdESElements()
	if err != nil {
		return nil, err
	}
	return xmlutils.DomUtilsAddElement(c.DocumentDom, parentDom, c.XadesNamespace(),
		currentElements.ElementCompleteCertificateRefs()), nil
}

// createCertRefsDom ports the private createCertRefsDom.
func (c *XAdESLevelC) createCertRefsDom(parentDom *xmldom.Node) (*xmldom.Node, error) {
	if c.Params.IsEn319132() {
		return xmlutils.DomUtilsAddElement(c.DocumentDom, parentDom, c.Xades141Namespace(),
			definition.XAdES141Element_CERT_REFS), nil
	}
	currentElements, err := c.CurrentXAdESElements()
	if err != nil {
		return nil, err
	}
	return xmlutils.DomUtilsAddElement(c.DocumentDom, parentDom, c.XadesNamespace(),
		currentElements.ElementCertRefs()), nil
}

// validationDataForCLevelInclusion ports the private getValidationDataForCLevelInclusion.
func (c *XAdESLevelC) validationDataForCLevelInclusion(
	validationDataContainer *validation.ValidationDataContainer,
	signature validation.AdvancedSignature) *validation.ValidationData {
	validationData := validationDataContainer.AllValidationDataForSignature(signature)
	validationData.ExcludeCertificateTokens(c.certificateTokensForExclusion())
	return validationData
}

// certificateTokensForExclusion ports the private getCertificateTokensForExclusion.
func (c *XAdESLevelC) certificateTokensForExclusion() []*model.CertificateToken {
	/*
	 * A.1.1 The CompleteCertificateRefsV2 qualifying property
	 *
	 * The CompleteCertificateRefsV2 qualifying property:
	 * ...
	 * 2) Shall not contain the reference to the signing certificate.
	 * ...
	 */
	signingCertificateToken := c.XadesSignature.SigningCertificateToken()
	if signingCertificateToken != nil {
		return []*model.CertificateToken{signingCertificateToken}
	}
	return []*model.CertificateToken{}
}

// incorporateCRLRefs incorporates CRL references like
//
//	<xades:CRLRefs>
//	    <xades:CRLRef>
//	        <xades:DigestAlgAndValue>
//	            <ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>
//	            <ds:DigestValue>G+z+DaZ6X44wEOueVYvZGmTh4dBkjjctKxcJYEV4HmU=</ds:DigestValue>
//	        </xades:DigestAlgAndValue>
//	        <xades:CRLIdentifier URI="LevelACAOK.crl">
//	            <xades:Issuer>CN=LevelACAOK,OU=Plugtests_STF-428_2011-2012,O=ETSI,C=FR</xades:Issuer>
//	            <xades:IssueTime>2012-03-13T13:58:28.000-03:00</xades:IssueTime>
//	        <xades:Number>4415260066222</xades:Number>
//
// Port of the private incorporateCRLRefs.
func (c *XAdESLevelC) incorporateCRLRefs(completeRevocationRefsDom *xmldom.Node,
	crlTokens []*spi.CRLToken) error {
	if len(crlTokens) == 0 {
		return nil
	}
	currentElements, err := c.CurrentXAdESElements()
	if err != nil {
		return err
	}

	crlRefsDom := xmlutils.DomUtilsAddElement(c.DocumentDom, completeRevocationRefsDom,
		c.XadesNamespace(), currentElements.ElementCRLRefs())

	for _, crlToken := range crlTokens {

		crlRefDom := xmlutils.DomUtilsAddElement(c.DocumentDom, crlRefsDom,
			c.XadesNamespace(), currentElements.ElementCRLRef())

		digestAlgorithm := c.Params.TokenReferencesDigestAlgorithm()
		digestAlgAndValueDom := xmlutils.DomUtilsAddElement(c.DocumentDom, crlRefDom,
			c.XadesNamespace(), currentElements.ElementDigestAlgAndValue())
		c.IncorporateDigestMethod(digestAlgAndValueDom, digestAlgorithm)
		if err := c.IncorporateDigestValueOfToken(digestAlgAndValueDom, digestAlgorithm, crlToken); err != nil {
			return err
		}

		crlIdentifierDom := xmlutils.DomUtilsAddElement(c.DocumentDom, crlRefDom,
			c.XadesNamespace(), currentElements.ElementCRLIdentifier())
		// crlIdentifierDom.setAttribute("URI",".crl");
		issuerX500PrincipalName := crlToken.IssuerX500Principal().RFC2253Name()
		xmlutils.DomUtilsAddTextElement(c.DocumentDom, crlIdentifierDom, c.XadesNamespace(),
			currentElements.ElementIssuer(), issuerX500PrincipalName)

		thisUpdate := crlToken.ThisUpdate()
		thisUpdateAsXmlFormat := xmlutils.DomUtilsCreateXMLGregorianCalendar(thisUpdate)
		xmlutils.DomUtilsAddTextElement(c.DocumentDom, crlIdentifierDom, c.XadesNamespace(),
			currentElements.ElementIssueTime(), thisUpdateAsXmlFormat)

		crlNumber := crlToken.CRLNumber()
		if crlNumber != nil {
			xmlutils.DomUtilsAddTextElement(c.DocumentDom, crlIdentifierDom, c.XadesNamespace(),
				currentElements.ElementNumber(), crlNumber.String())
		}

	}
	return nil
}

// incorporateOCSPRefs adds OCSP references like
//
//	<xades:CRLRefs/>
//	<xades:OCSPRefs>
//	    <xades:OCSPRef>
//	        <xades:OCSPIdentifier>
//	            <xades:ResponderID>
//	                <xades:ByName>C=AA,O=DSS,CN=OCSP A</xades:ByName>
//	            </xades:ResponderID>
//	            <xades:ProducedAt>2013-11-25T12:33:34.000+01:00</xades:ProducedAt>
//	        </xades:OCSPIdentifier>
//	        <xades:DigestAlgAndValue>
//	            <ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
//	            <ds:DigestValue>O1uHdchN+zFzbGrBg2FP3/idD0k=</ds:DigestValue>
//	            ...
//
// Port of the private incorporateOCSPRefs.
func (c *XAdESLevelC) incorporateOCSPRefs(completeRevocationRefsDom *xmldom.Node,
	ocspTokens []*spi.OCSPToken) error {
	if len(ocspTokens) == 0 {
		return nil
	}
	currentElements, err := c.CurrentXAdESElements()
	if err != nil {
		return err
	}

	ocspRefsDom := xmlutils.DomUtilsAddElement(c.DocumentDom, completeRevocationRefsDom,
		c.XadesNamespace(), currentElements.ElementOCSPRefs())

	for _, ocspToken := range ocspTokens {

		basicOcspResp := ocspToken.BasicOCSPResp()
		if basicOcspResp != nil {

			ocspRefDom := xmlutils.DomUtilsAddElement(c.DocumentDom, ocspRefsDom,
				c.XadesNamespace(), currentElements.ElementOCSPRef())

			ocspIdentifierDom := xmlutils.DomUtilsAddElement(c.DocumentDom, ocspRefDom,
				c.XadesNamespace(), currentElements.ElementOCSPIdentifier())
			responderIDDom := xmlutils.DomUtilsAddElement(c.DocumentDom, ocspIdentifierDom,
				c.XadesNamespace(), currentElements.ElementResponderID())

			respID := basicOcspResp.ResponderID()
			responderId, err := spi.DSSRevocationUtilsDSSResponderIDFromRespID(respID)
			if err != nil {
				return err
			}

			if responderId.X500Principal() != nil {
				xmlutils.DomUtilsAddTextElement(c.DocumentDom, responderIDDom, c.XadesNamespace(),
					currentElements.ElementByName(), responderId.X500Principal().String())
			} else {
				base64EncodedKeyHashOctetStringBytes := utils.ToBase64(responderId.Ski())
				xmlutils.DomUtilsAddTextElement(c.DocumentDom, responderIDDom, c.XadesNamespace(),
					currentElements.ElementByKey(), base64EncodedKeyHashOctetStringBytes)
			}

			producedAt := basicOcspResp.ProducedAt()
			producedAtXmlEncoded := xmlutils.DomUtilsCreateXMLGregorianCalendar(producedAt)
			xmlutils.DomUtilsAddTextElement(c.DocumentDom, ocspIdentifierDom, c.XadesNamespace(),
				currentElements.ElementProducedAt(), producedAtXmlEncoded)

			digestAlgorithm := c.Params.TokenReferencesDigestAlgorithm()
			digestAlgAndValueDom := xmlutils.DomUtilsAddElement(c.DocumentDom, ocspRefDom,
				c.XadesNamespace(), currentElements.ElementDigestAlgAndValue())
			c.IncorporateDigestMethod(digestAlgAndValueDom, digestAlgorithm)
			if err := c.IncorporateDigestValueOfToken(digestAlgAndValueDom, digestAlgorithm,
				ocspToken); err != nil {
				return err
			}
		}

	}
	return nil
}
