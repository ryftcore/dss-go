// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMS.java (DSS 6.5.RC1).
//
// Java's CMS is an interface implemented once per ICMSUtils backend (CMSSignedDataObject in
// dss-cms-object, a streaming counterpart in dss-cms-stream). This port has one native
// implementation (see doc.go), so CMS is the concrete struct directly, wrapping the
// *cmscore.CMS internal/cmscore already builds and parses byte-exactly. Every getter below
// keeps its Java name (get- dropped) and return shape translated per doc.go's BC-replacement
// table: a BouncyCastle Store<X509CertificateHolder/X509CRLHolder/...> becomes [][]byte, the
// DER encoding of each member, and SignerInformationStore becomes []*cmscore.SignerInfo.
package cms

import (
	"encoding/asn1"

	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
)

// CMS represents a content of a CMS Signed Data object. Port of the CMS interface.
type CMS struct {
	core *cmscore.CMS
}

// newCMS wraps an already parsed or built *cmscore.CMS.
func newCMS(core *cmscore.CMS) *CMS {
	if core == nil {
		return nil
	}
	return &CMS{core: core}
}

// Core returns the internal/cmscore representation this CMS wraps. It is exported so that
// cades and other same-phase callers needing cmscore-level access (e.g. to feed
// spi/validation.NewTimestampTokenFromCMS) do not have to re-parse the document; it has no
// Java counterpart.
func (c *CMS) Core() *cmscore.CMS { return c.core }

// Version returns value of SignedData.version field. Port of #getVersion.
func (c *CMS) Version() int { return c.core.Version() }

// DigestAlgorithmIDs returns a set of algorithm identifiers (OIDs) incorporated within
// SignedData.digestAlgorithms field of CMS. Port of #getDigestAlgorithmIDs; Java returns a
// Set, the received order is kept here since every caller either de-duplicates itself
// (CMSUtilsPopulateDigestAlgorithmSet) or does not care about duplicates.
func (c *CMS) DigestAlgorithmIDs() []*asn1ber.AlgorithmIdentifier { return c.core.DigestAlgorithmIDs() }

// IsDetachedSignature returns whether the signature is detached (i.e.
// SignedData.encapContentInfo.eContent is null). Port of #isDetachedSignature.
func (c *CMS) IsDetachedSignature() bool { return c.core.IsDetachedSignature() }

// SignedContentType gets signed content type, present within the
// SignedData.encapContentInfo.eContentType field. Port of #getSignedContentType.
func (c *CMS) SignedContentType() asn1.ObjectIdentifier { return c.core.SignedContentType() }

// SignedContent gets the signed content incorporated within the
// SignedData.encapContentInfo.eContent field, nil for a detached signature. Port of
// #getSignedContent.
func (c *CMS) SignedContent() model.DSSDocument {
	content := c.core.SignedContent()
	if content == nil {
		return nil
	}
	return model.NewInMemoryDocument(content)
}

// Certificates gets the certificates store, representing the value of SignedData.certificates
// field: the DER encoding of every plain X.509 certificate it holds. Port of #getCertificates.
func (c *CMS) Certificates() [][]byte { return c.core.Certificates() }

// AttributeCertificates gets attribute certificates incorporated within CMS: the DER encoding
// of every version 2 attribute certificate of SignedData.certificates. Port of
// #getAttributeCertificates.
func (c *CMS) AttributeCertificates() [][]byte { return c.core.AttributeCertificates() }

// CRLs gets the CRLs store (OCSP excluded), representing the value of SignedData.crls field.
// Port of #getCRLs.
func (c *CMS) CRLs() [][]byte { return c.core.CRLs() }

// OcspResponseStore gets the OCSP Responses Store, incorporated within the SignedData.crls
// field under id-ri-ocsp-response. Port of #getOcspResponseStore.
func (c *CMS) OcspResponseStore() [][]byte { return c.core.OCSPResponses() }

// OcspBasicStore gets the OCSP Basic Store, incorporated within the SignedData.crls field
// under id-pkix-ocsp-basic. Port of #getOcspBasicStore.
func (c *CMS) OcspBasicStore() [][]byte { return c.core.OCSPBasicResponses() }

// SignerInfos gets the signers of the signature, incorporated within the
// SignedData.signerInfos field. Port of #getSignerInfos.
func (c *CMS) SignerInfos() []*cmscore.SignerInfo { return c.core.SignerInfos() }

// DEREncoded gets DER-encoded content of the CMS SignedData.
// NOTE: This method returns the encoded value using an in-memory byte array (see doc.go: not
// applicable for large CMS processing, and this port has no other kind). Port of
// #getDEREncoded.
func (c *CMS) DEREncoded() []byte { return c.core.DEREncoded() }

// Encoded gets encoded content of the CMS SignedData, keeping the original encoding. Port of
// #getEncoded.
func (c *CMS) Encoded() []byte { return c.core.Encoded() }
