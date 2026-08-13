// CMS, the entry point of this package: a ContentInfo carrying a SignedData. It replaces
// org.bouncycastle.cms.CMSSignedData and answers what eu.europa.esig.dss.cms.CMS asks for, so
// that the port of the dss-cms module is a thin wrapper over it.
package cmscore

import (
	"encoding/asn1"
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// CMS is a parsed or built CMS document: an id-signedData ContentInfo and the SignedData it
// carries.
type CMS struct {
	contentInfo *ContentInfo
	signedData  *SignedData
}

// ParseCMS decodes a complete CMS document. BER is accepted, and the document's original bytes
// are preserved throughout - see the package documentation. Trailing bytes after the document
// are rejected; use ParseCMSTolerateTrailingBytes for a whole-file parse that must not be
// (see that function's doc comment for why the two differ).
func ParseCMS(input []byte) (*CMS, error) {
	contentInfo, err := ParseContentInfo(input)
	if err != nil {
		return nil, err
	}
	return CMSFromContentInfo(contentInfo)
}

// ParseCMSTolerateTrailingBytes decodes a complete CMS document the same way ParseCMS does,
// except it does not require the input to be fully consumed - see
// ParseContentInfoTolerateTrailingBytes's doc comment. cms.CMSUtilsParseToCMSBinaries, the
// entry point for parsing a whole document (as opposed to a CMS structure embedded in a fixed
// field, e.g. an OCTET STRING's content), uses this rather than ParseCMS.
func ParseCMSTolerateTrailingBytes(input []byte) (*CMS, error) {
	contentInfo, err := ParseContentInfoTolerateTrailingBytes(input)
	if err != nil {
		return nil, err
	}
	return CMSFromContentInfo(contentInfo)
}

// CMSFromContentInfo decodes the SignedData of an already parsed ContentInfo.
func CMSFromContentInfo(contentInfo *ContentInfo) (*CMS, error) {
	if !contentInfo.ContentType.Equal(OIDSignedData) {
		return nil, fmt.Errorf("cmscore: the ContentInfo carries %s, not id-signedData", contentInfo.ContentType)
	}
	if contentInfo.Content() == nil {
		return nil, errNoContent
	}
	signedData, err := SignedDataFromElement(contentInfo.Content())
	if err != nil {
		return nil, err
	}
	return &CMS{contentInfo: contentInfo, signedData: signedData}, nil
}

// NewCMS wraps a built SignedData, so that a signer can hand out a complete document.
func NewCMS(signedData *SignedData) *CMS {
	return &CMS{contentInfo: NewContentInfo(OIDSignedData, signedData.DER()), signedData: signedData}
}

// ContentInfo returns the outermost structure.
func (c *CMS) ContentInfo() *ContentInfo { return c.contentInfo }

// SignedData returns the carried SignedData.
func (c *CMS) SignedData() *SignedData { return c.signedData }

// Version returns SignedData.version. Port of CMS#getVersion.
func (c *CMS) Version() int { return c.signedData.Version }

// DigestAlgorithmIDs returns the members of SignedData.digestAlgorithms. Port of
// CMS#getDigestAlgorithmIDs, whose Java signature returns a Set; the received order is kept
// here, the caller deduplicating if it cares.
func (c *CMS) DigestAlgorithmIDs() []*asn1ber.AlgorithmIdentifier {
	return c.signedData.DigestAlgorithms
}

// IsDetachedSignature reports whether encapContentInfo.eContent is absent. Port of
// CMS#isDetachedSignature.
func (c *CMS) IsDetachedSignature() bool { return c.signedData.IsDetached() }

// SignedContentType returns encapContentInfo.eContentType. Port of CMS#getSignedContentType.
func (c *CMS) SignedContentType() asn1.ObjectIdentifier {
	return c.signedData.EncapContentInfo.EContentType
}

// SignedContent returns the signed content, nil for a detached signature. Port of
// CMS#getSignedContent, minus the DSSDocument wrapping.
func (c *CMS) SignedContent() []byte { return c.signedData.EncapContentInfo.Content() }

// Certificates returns the encoding of every plain X.509 certificate of SignedData.certificates.
// Port of CMS#getCertificates.
func (c *CMS) Certificates() [][]byte { return c.signedData.CertificateDERs() }

// AttributeCertificates returns the encoding of every version 2 attribute certificate of
// SignedData.certificates. Port of CMS#getAttributeCertificates.
func (c *CMS) AttributeCertificates() [][]byte {
	return c.signedData.Certificates.AttributeCertificatesV2()
}

// CRLs returns the encoding of every CertificateList of SignedData.crls, OCSP responses
// excluded. Port of CMS#getCRLs.
func (c *CMS) CRLs() [][]byte { return c.signedData.CRLs.CRLs() }

// OCSPResponses returns the OCSPResponse encodings carried in SignedData.crls under
// id-ri-ocsp-response. Port of CMS#getOcspResponseStore.
func (c *CMS) OCSPResponses() [][]byte { return c.signedData.CRLs.OCSPResponses() }

// OCSPBasicResponses returns the BasicOCSPResponse encodings carried in SignedData.crls under
// id-pkix-ocsp-basic. Port of CMS#getOcspBasicStore.
func (c *CMS) OCSPBasicResponses() [][]byte { return c.signedData.CRLs.OCSPBasicResponses() }

// SignerInfos returns the members of SignedData.signerInfos. Port of CMS#getSignerInfos.
func (c *CMS) SignerInfos() []*SignerInfo { return c.signedData.SignerInfos }

// Encoded returns the document's original encoding, BER included. Port of CMS#getEncoded.
func (c *CMS) Encoded() []byte { return c.contentInfo.Encoded() }

// DEREncoded returns the DER encoding of the document, normalising a BER input. Port of
// CMS#getDEREncoded.
func (c *CMS) DEREncoded() []byte {
	if element := c.contentInfo.Element(); element != nil {
		return element.DEREncoded()
	}
	return c.contentInfo.DER()
}

// IsDefiniteLength reports whether the ContentInfo used the definite-length form, which is
// what DSS turns into the choice between the DL and the BER encoding.
func (c *CMS) IsDefiniteLength() bool { return c.contentInfo.IsDefiniteLength() }
