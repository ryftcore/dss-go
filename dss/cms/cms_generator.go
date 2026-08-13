// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/CMSGenerator.java (DSS 6.5.RC1).
//
// Java's static loadCMSGenerator() ServiceLoader factory - "choose between 'dss-cms-object' and
// 'dss-cms-stream'" - has no port: this package has exactly one native CMS generator (see
// doc.go), constructed directly with NewAbstractCMSGenerator.
package cms

import (
	"github.com/utain/esig/dss/internal/asn1ber"
	"github.com/utain/esig/dss/internal/cmscore"
	"github.com/utain/esig/dss/model"
)

// CMSGenerator generates a CMS with the given input data. Port of the CMSGenerator interface;
// *AbstractCMSGenerator is its (only) implementation.
type CMSGenerator interface {
	// SetSignerInfoGenerator adds a SignerInfoGenerator containing information about a new
	// signer to be embedded within CMS. Port of #setSignerInfoGenerator.
	SetSignerInfoGenerator(signerInfoGenerator *SignerInfoGenerator)

	// SetCertificates adds certificates to be embedded within SignedData.certificates field,
	// each as its own DER encoding. Port of #setCertificates.
	SetCertificates(certificateStore [][]byte)

	// SetSigners adds existing SignerInformation's. Port of #setSigners.
	SetSigners(signers []*cmscore.SignerInfo)

	// SetAttributeCertificates adds attribute certificates, each as its own DER encoding
	// (a [2] IMPLICIT AttributeCertificateV2). Port of #setAttributeCertificates.
	SetAttributeCertificates(attributeCertificates [][]byte)

	// SetCRLs adds CRLs, each as its own DER encoding. Port of #setCRLs.
	SetCRLs(crls [][]byte)

	// SetOcspBasicStore adds a collection of OCSP basic responses, each as its own DER encoding.
	// Port of #setOcspBasicStore.
	SetOcspBasicStore(ocspBasicStore [][]byte)

	// SetOcspResponsesStore adds a collection of OCSP responses, each as its own DER encoding.
	// Port of #setOcspResponsesStore.
	SetOcspResponsesStore(ocspResponsesStore [][]byte)

	// SetDigestAlgorithmIDs adds a collection of digest algorithm IDs. Port of #setDigestAlgorithmIDs.
	SetDigestAlgorithmIDs(digestAlgorithmIDs []*asn1ber.AlgorithmIdentifier)

	// SetToBeSignedDocument adds a document to be signed. Port of #setToBeSignedDocument.
	SetToBeSignedDocument(document model.DSSDocument)

	// SetEncapsulate sets whether the document shall be encapsulated within CMS. Port of
	// #setEncapsulate.
	SetEncapsulate(encapsulate bool)

	// Generate generates the CMS. Port of #generate.
	Generate() (*CMS, error)
}

var _ CMSGenerator = (*AbstractCMSGenerator)(nil)
