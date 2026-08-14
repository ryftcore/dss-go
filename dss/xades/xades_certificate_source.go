// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESCertificateSource.java
// (DSS 6.5.RC1).
//
// DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey (Java DSSXMLUtils.getKeyInfoSigningCertificatePublicKey(Element))
// is a forward dependency of the root eu.europa.esig.dss.xades package (not in this manifest),
// per the same "DSSXMLUtils"-prefixed forward-dependency convention xades_signature.go's header
// already documents for its own sibling call sites; its assumed shape is
// func DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey(signatureElement *xmldom.Node) *model.PublicKey.
package xades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/spi"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/common"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESCertificateSource provides the mechanism to retrieve certificates contained in a XAdES
// signature. Port of the class XAdESCertificateSource, extending spi.SignatureCertificateSource.
type XAdESCertificateSource struct {
	spi.SignatureCertificateSource

	// signatureElement is the Signature element.
	signatureElement *xmldom.Node

	// xadesPaths contains a list of XAdES path corresponding to the signature.
	xadesPaths definition.XAdESPath
}

// NewXAdESCertificateSource is the port of the constructor XAdESCertificateSource(Element,
// XAdESPath). All certificates are extracted during instantiation.
//
// Panics with the Java messages when signatureElement or xadesPaths is missing
// (Objects.requireNonNull).
func NewXAdESCertificateSource(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *XAdESCertificateSource {
	if signatureElement == nil {
		panic("Element signature must not be null")
	}
	if xadesPaths == nil {
		panic("XAdESPaths must not be null")
	}

	s := &XAdESCertificateSource{
		signatureElement: signatureElement,
		xadesPaths:       xadesPaths,
	}
	s.InitSignatureCertificateSource(s)

	// init
	s.extractCertificates(common.XMLDSigPath_KEY_INFO_X509_CERTIFICATE_PATH, enumerations.CertificateOrigin_KEY_INFO)
	s.extractCertificates(xadesPaths.EncapsulatedCertificateValuesPath(), enumerations.CertificateOrigin_CERTIFICATE_VALUES)
	s.extractCertificates(xadesPaths.EncapsulatedAttrAuthoritiesCertValuesPath(), enumerations.CertificateOrigin_ATTR_AUTHORITIES_CERT_VALUES)
	s.extractCertificates(xadesPaths.EncapsulatedTimeStampValidationDataCertValuesPath(), enumerations.CertificateOrigin_TIMESTAMP_VALIDATION_DATA)
	s.extractCertificates(xadesPaths.EncapsulatedAnyValidationDataCertValuesPath(), enumerations.CertificateOrigin_ANY_VALIDATION_DATA)

	s.extractCertificateRefs(xadesPaths.SigningCertificateChildren(), xadesPaths.SigningCertificateV2Children(),
		enumerations.CertificateRefOrigin_SIGNING_CERTIFICATE)
	s.extractCertificateRefs(xadesPaths.CompleteCertificateRefsCertPath(), xadesPaths.CompleteCertificateRefsV2CertPath(),
		enumerations.CertificateRefOrigin_COMPLETE_CERTIFICATE_REFS)
	s.extractCertificateRefs(xadesPaths.AttributeCertificateRefsCertPath(), xadesPaths.AttributeCertificateRefsV2CertPath(),
		enumerations.CertificateRefOrigin_ATTRIBUTE_CERTIFICATE_REFS)

	// Upstream logs "+XAdESCertificateSource".

	return s
}

// extractCertificates ports the private extractCertificates(XPathQuery, CertificateOrigin).
func (s *XAdESCertificateSource) extractCertificates(xPathQuery common.XPathQuery, origin enumerations.CertificateOrigin) {
	if xPathQuery == nil {
		return
	}
	nodeList, err := xmlutils.XPathUtilsGetNodeList(s.signatureElement, xPathQuery)
	if err != nil {
		return
	}
	for _, certificateElement := range nodeList {
		base64EncodedCertificate := certificateElement.TextContent()
		derEncoded := utils.FromBase64(base64EncodedCertificate)
		cert, err := spi.DSSUtilsLoadCertificateFromBinary(derEncoded)
		if err != nil {
			// Upstream logs "Unable to parse certificate '{}' : {}".
			continue
		}
		s.AddCertificateWithOrigin(cert, origin)
	}
}

// extractCertificateRefs ports the private extractCertificateRefs(XPathQuery, XPathQuery,
// CertificateRefOrigin).
func (s *XAdESCertificateSource) extractCertificateRefs(xpathV1, xpathV2 common.XPathQuery, origin enumerations.CertificateRefOrigin) {
	if xpathV1 != nil {
		certRefNodeList, err := xmlutils.XPathUtilsGetNodeList(s.signatureElement, xpathV1)
		if err == nil {
			s.extractXAdESCertsV1(certRefNodeList, origin)
		}
	}
	if xpathV2 != nil {
		certRefNodeList, err := xmlutils.XPathUtilsGetNodeList(s.signatureElement, xpathV2)
		if err == nil {
			s.extractXAdESCertsV2(certRefNodeList, origin)
		}
	}
}

// extractXAdESCertsV1 ports the private extractXAdESCertsV1(NodeList, CertificateRefOrigin).
func (s *XAdESCertificateSource) extractXAdESCertsV1(certNodeList []*xmldom.Node, origin enumerations.CertificateRefOrigin) {
	for _, certRefElement := range certNodeList {
		certificateRef := XAdESCertificateRefExtractionUtilsCreateCertificateRefFromV1(certRefElement, s.xadesPaths)
		if certificateRef != nil {
			s.AddCertificateRef(certificateRef, origin)
		}
	}
}

// extractXAdESCertsV2 ports the private extractXAdESCertsV2(NodeList, CertificateRefOrigin).
func (s *XAdESCertificateSource) extractXAdESCertsV2(certNodeList []*xmldom.Node, origin enumerations.CertificateRefOrigin) {
	for _, certRefElement := range certNodeList {
		certificateRef := XAdESCertificateRefExtractionUtilsCreateCertificateRefFromV2(certRefElement, s.xadesPaths)
		if certificateRef != nil {
			s.AddCertificateRef(certificateRef, origin)
		}
	}
}

// ExtractCandidatesForSigningCertificate implements spi.SignatureCertificateSourceOverrides.
// Port of the protected extractCandidatesForSigningCertificate(CertificateSource) override.
//
// 5.1.4.1 XAdES processing: candidates for the signing certificate extracted from ds:KeyInfo
// shall be checked against all references present in the ds:SigningCertificate property, if
// present, since one of these references shall be a reference to the signing certificate.
func (s *XAdESCertificateSource) ExtractCandidatesForSigningCertificate(signingCertificateSource spi.CertificateSource) *spi.CandidatesForSigningCertificate {
	candidatesForSigningCertificate := spi.NewCandidatesForSigningCertificate()

	for _, certificateToken := range s.KeyInfoCertificates() {
		candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
	}

	// if KeyInfo does not contain certificates, check other certificates embedded into the
	// signature
	if candidatesForSigningCertificate.IsEmpty() {
		publicKey := DSSXMLUtilsGetKeyInfoSigningCertificatePublicKey(s.signatureElement)
		if publicKey != nil {
			// try to find out the signing certificate token by provided public key
			certsByPublicKey := s.ByPublicKey(publicKey)

			if utils.IsMapNotEmpty(certsByPublicKey) {
				for _, certificateToken := range certsByPublicKey {
					candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
				}
			} else {
				// process public key only if no certificates found
				candidatesForSigningCertificate.Add(spi.NewCertificateValidityFromPublicKey(publicKey))
			}

		} else {
			// Add all found certificates
			for _, certificateToken := range s.Certificates() {
				candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certificateToken))
			}
		}
	}

	if signingCertificateSource != nil {
		s.resolveFromSource(signingCertificateSource, candidatesForSigningCertificate)
	}

	s.checkCandidatesAgainstSigningCertificateRef(candidatesForSigningCertificate)

	return candidatesForSigningCertificate
}

// resolveFromSource ports the private resolveFromSource(CertificateSource,
// CandidatesForSigningCertificate).
func (s *XAdESCertificateSource) resolveFromSource(certificateSource spi.CertificateSource, candidatesForSigningCertificate *spi.CandidatesForSigningCertificate) {
	signingCertificateRefs := s.SigningCertificateRefs()
	if utils.IsCollectionNotEmpty(signingCertificateRefs) {
		for _, certificateRef := range signingCertificateRefs {
			s.resolveForReference(certificateRef, certificateSource, candidatesForSigningCertificate)
		}
	} else {
		certificates := certificateSource.Certificates()
		// Upstream logs "No signing certificate reference found. Resolve all {} certificates
		// from the provided certificate source as signing candidates.".
		for _, certCandidate := range certificates {
			candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certCandidate))
		}
	}
}

// resolveForReference ports the private resolveForReference(CertificateRef, CertificateSource,
// CandidatesForSigningCertificate).
func (s *XAdESCertificateSource) resolveForReference(certificateRef *spi.CertificateRef, certificateSource spi.CertificateSource, candidatesForSigningCertificate *spi.CandidatesForSigningCertificate) {
	signerIdentifier := certificateRef.CertificateIdentifier()
	if signerIdentifier != nil {
		certificatesByIdentifier := certificateSource.BySignerIdentifier(signerIdentifier)
		if utils.IsMapNotEmpty(certificatesByIdentifier) {
			// Upstream logs "Resolved certificate by certificate identifier".
			for _, certCandidate := range certificatesByIdentifier {
				candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certCandidate))
			}
			return
		}
	}

	certDigest := certificateRef.CertDigest()
	if !certDigest.IsEmpty() {
		certificatesByDigest := certificateSource.ByCertificateDigest(certDigest)
		if utils.IsMapNotEmpty(certificatesByDigest) {
			// Upstream logs "Resolved certificate by digest".
			for _, certCandidate := range certificatesByDigest {
				candidatesForSigningCertificate.Add(spi.NewCertificateValidity(certCandidate))
			}
		}
	}
}

// checkCandidatesAgainstSigningCertificateRef ports the private
// checkCandidatesAgainstSigningCertificateRef(CandidatesForSigningCertificate): checks the
// protection of the certificates included within the signature (XAdES: KeyInfo) against the
// substitution attack.
func (s *XAdESCertificateSource) checkCandidatesAgainstSigningCertificateRef(candidates *spi.CandidatesForSigningCertificate) {
	potentialSigningCertificates := s.SigningCertificateRefs()
	if utils.IsCollectionNotEmpty(potentialSigningCertificates) {
		// first reference shall be a reference to a signing certificate
		signingCert := potentialSigningCertificates[0]

		var bestCertificateValidity *spi.CertificateValidity
		// check all certificates against the signingCert ref and find the best one
		certificateValidityList := candidates.CertificateValidityList()
		for _, certificateValidity := range certificateValidityList {
			if s.isValid(certificateValidity, signingCert) {
				bestCertificateValidity = certificateValidity
			}
		}
		if bestCertificateValidity != nil {
			// SetTheCertificateValidity's error indicates the candidate is not in the list, which
			// cannot happen here since bestCertificateValidity was read from candidates itself.
			_ = candidates.SetTheCertificateValidity(bestCertificateValidity)
		}
	}
}

// isValid ports the private isValid(CertificateValidity, CertificateRef).
func (s *XAdESCertificateSource) isValid(certificateValidity *spi.CertificateValidity, signingCert *spi.CertificateRef) bool {
	certificateValidity.SetDigestPresent(!signingCert.CertDigest().IsEmpty())
	certificateValidity.SetIssuerSerialPresent(signingCert.CertificateIdentifier() != nil)

	certificateToken := certificateValidity.CertificateToken()
	if certificateToken != nil {
		certificateMatcher := s.CertificateMatcher()
		certificateValidity.SetDigestEqual(certificateMatcher.MatchByDigest(certificateToken, signingCert))
		certificateValidity.SetSerialNumberEqual(certificateMatcher.MatchBySerialNumber(certificateToken, signingCert))
		certificateValidity.SetDistinguishedNameEqual(certificateMatcher.MatchByIssuerName(certificateToken, signingCert))
	}
	return certificateValidity.IsValid()
}

// compile-time assertion: an XAdESCertificateSource satisfies its own overrides contract.
var _ spi.SignatureCertificateSourceOverrides = (*XAdESCertificateSource)(nil)
