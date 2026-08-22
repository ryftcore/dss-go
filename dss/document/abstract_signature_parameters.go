// Ported from dss-document/src/main/java/eu/europa/esig/dss/signature/AbstractSignatureParameters.java (DSS 6.5.RC1).
package document

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// AbstractSignatureParameters holds parameters for a Signature creation/extension, generic over
// the TP implementation of certain format timestamp parameters.
type AbstractSignatureParameters[TP model.SerializableTimestampParameters] struct {
	model.AbstractSerializableSignatureParameters[TP]

	// context is the internal signature processing variable.
	context *ProfileParameters

	// detachedContents are the documents to be signed.
	detachedContents []model.DSSDocument

	// signingCertificate contains the signing certificate.
	signingCertificate *model.CertificateToken

	// signedData is an optional parameter that contains the actual canonicalized data that was
	// used when creating the signature value. This allows scenarios where ToBeSigned was
	// externally updated before signature value was created (i.e. signature certificate was
	// appended). If this parameter is specified it will be used in the signed document.
	signedData []byte

	// certificateChain is the List of chain of certificates. It includes the signing
	// certificate.
	certificateChain []*model.CertificateToken

	// contentTimestamps is here because that's a signed attribute. It must be computed before
	// GetDataToSign/SignDocument.
	contentTimestamps []*validation.TimestampToken
}

// NewAbstractSignatureParameters instantiates the object with null values. Port of the protected
// no-arg constructor.
func NewAbstractSignatureParameters[TP model.SerializableTimestampParameters]() AbstractSignatureParameters[TP] {
	return AbstractSignatureParameters[TP]{
		AbstractSerializableSignatureParameters: model.NewAbstractSerializableSignatureParameters[TP](),
		certificateChain:                        []*model.CertificateToken{},
	}
}

// ContentTimestamps returns the list of TimestampToken to be incorporated within the signature
// and representing the content-timestamp. Port of #getContentTimestamps.
func (p *AbstractSignatureParameters[TP]) ContentTimestamps() []*validation.TimestampToken {
	return p.contentTimestamps
}

// SetContentTimestamps sets a list of content timestamps to be included into the signature. Port
// of #setContentTimestamps.
func (p *AbstractSignatureParameters[TP]) SetContentTimestamps(contentTimestamps []*validation.TimestampToken) {
	p.contentTimestamps = contentTimestamps
}

// DetachedContents returns the documents to sign. In the case of the DETACHED signature this is
// the detached document. Port of #getDetachedContents.
func (p *AbstractSignatureParameters[TP]) DetachedContents() []model.DSSDocument {
	if p.context != nil && len(p.context.DetachedContents()) > 0 {
		return p.context.DetachedContents()
	}
	if len(p.detachedContents) > 0 {
		return p.detachedContents
	}
	return []model.DSSDocument{}
}

// SetDetachedContents sets the documents to sign. When signing, this method is internally
// invoked by AbstractSignatureService and the related variable detachedContent is overwritten by
// the service parameter. In the case of the DETACHED signature this is the detached document. In
// the case of ASiC-S this is the document to be signed.
//
// When extending, this method must be invoked to indicate the detachedContent. Port of
// #setDetachedContents.
func (p *AbstractSignatureParameters[TP]) SetDetachedContents(detachedContents []model.DSSDocument) {
	p.detachedContents = detachedContents
}

// SigningCertificate gets the signing certificate. Port of #getSigningCertificate.
func (p *AbstractSignatureParameters[TP]) SigningCertificate() *model.CertificateToken {
	return p.signingCertificate
}

// SetSigningCertificate sets the signing certificate. The encryption algorithm is also set from
// the public key.
//
// NOTE: This method overwrites the encryptionAlgorithm value by extracting the encryption
// algorithm from the provided signing-certificate. In order to enforce a specific encryption
// algorithm (when supported by the key), please call SetEncryptionAlgorithm after this method.
// Port of #setSigningCertificate.
func (p *AbstractSignatureParameters[TP]) SetSigningCertificate(signingCertificate *model.CertificateToken) {
	p.signingCertificate = signingCertificate
	encryptionAlgorithm, err := documentEncryptionAlgorithmForKey(signingCertificate.PublicKey())
	if err != nil {
		panic(err)
	}
	p.SetEncryptionAlgorithm(encryptionAlgorithm)
}

// documentEncryptionAlgorithmForKey ports EncryptionAlgorithm#forKey(Key), the sole call site of
// which lives in this manifest's SetSigningCertificate. forKey delegates to forName(key.
// getAlgorithm()); enumerations.EncryptionAlgorithmForName (already landed in phase 1a) is that
// forName, and model.PublicKey#Algorithm is the Go port of Key#getAlgorithm - so this local
// helper reproduces forKey's one-line body without needing forKey itself ported out-of-manifest.
func documentEncryptionAlgorithmForKey(publicKey *model.PublicKey) (enumerations.EncryptionAlgorithm, error) {
	return enumerations.EncryptionAlgorithmForName(publicKey.Algorithm())
}

// SignedData gets signed data. Port of #getSignedData.
func (p *AbstractSignatureParameters[TP]) SignedData() []byte {
	return p.signedData
}

// SetSignedData sets signed data: data that was used when creating the signature value. Port of
// #setSignedData.
func (p *AbstractSignatureParameters[TP]) SetSignedData(signedData []byte) {
	p.signedData = signedData
}

// CertificateChain gets the certificate chain. Port of #getCertificateChain.
func (p *AbstractSignatureParameters[TP]) CertificateChain() []*model.CertificateToken {
	return p.certificateChain
}

// SetCertificateChain sets the certificate chain. Port of the List overload of
// #setCertificateChain.
func (p *AbstractSignatureParameters[TP]) SetCertificateChain(certificateChain []*model.CertificateToken) {
	p.certificateChain = certificateChain
}

// SetCertificateChainFromTokens sets the list of certificates which constitute the chain. If the
// certificate is already present in the array then it is ignored. Port of the varargs overload
// of #setCertificateChain.
func (p *AbstractSignatureParameters[TP]) SetCertificateChainFromTokens(certificateChainArray ...*model.CertificateToken) {
	for _, certificate := range certificateChainArray {
		if certificate == nil {
			continue
		}
		found := false
		for _, existing := range p.certificateChain {
			if existing == certificate {
				found = true
				break
			}
		}
		if !found {
			p.certificateChain = append(p.certificateChain, certificate)
		}
	}
}

// GetDeterministicId returns the deterministic identifier used for unique identification of a
// created signature (used in XAdES and PAdES). The identifier is built in a deterministic way to
// ensure the same value on both method calls during the signature creation. Port of
// #getDeterministicId.
func (p *AbstractSignatureParameters[TP]) GetDeterministicId() string {
	deterministicId := p.GetContext().DeterministicId()
	if deterministicId == "" {
		var identifier *model.TokenIdentifier
		if p.signingCertificate != nil {
			identifier = p.signingCertificate.BuildTokenIdentifier()
		}
		var signingDate time.Time
		if sd := p.BLevel().SigningDate(); sd != nil {
			signingDate = *sd
		}
		newDeterministicId, err := spi.DSSUtilsDeterministicID(signingDate, identifier)
		if err != nil {
			panic(err)
		}
		deterministicId = newDeterministicId
		p.GetContext().SetDeterministicId(deterministicId)
	}
	return deterministicId
}

// GetContext gets the signature creation context (internal variable). Port of #getContext.
func (p *AbstractSignatureParameters[TP]) GetContext() *ProfileParameters {
	if p.context == nil {
		p.context = NewProfileParameters()
	}
	return p.context
}

// Reinit re-inits signature parameters to clean temporary settings. Port of #reinit.
func (p *AbstractSignatureParameters[TP]) Reinit() {
	p.context = nil
}

// String ports #toString.
func (p *AbstractSignatureParameters[TP]) String() string {
	return fmt.Sprintf("AbstractSignatureParameters [context=%v, detachedContents=%v, signingCertificate=%v, signedData=%v, certificateChain=%v, contentTimestamps=%v] %s",
		p.context, p.detachedContents, p.signingCertificate, p.signedData, p.certificateChain, p.contentTimestamps,
		p.AbstractSerializableSignatureParameters.String())
}

// Equals ports #equals.
func (p *AbstractSignatureParameters[TP]) Equals(other *AbstractSignatureParameters[TP]) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.AbstractSerializableSignatureParameters.Equals(&other.AbstractSerializableSignatureParameters) {
		return false
	}
	if !abstractSignatureParametersContextEquals(p.context, other.context) {
		return false
	}
	if !abstractSignatureParametersDocumentsEqual(p.detachedContents, other.detachedContents) {
		return false
	}
	if !abstractSignatureParametersCertificateEquals(p.signingCertificate, other.signingCertificate) {
		return false
	}
	if !bytes.Equal(p.signedData, other.signedData) {
		return false
	}
	if !abstractSignatureParametersCertificateChainEquals(p.certificateChain, other.certificateChain) {
		return false
	}
	return abstractSignatureParametersTimestampsEqual(p.contentTimestamps, other.contentTimestamps)
}

// abstractSignatureParametersContextEquals ports Objects.equals(context, ...).
func abstractSignatureParametersContextEquals(a, b *ProfileParameters) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// abstractSignatureParametersDocumentsEqual ports Objects.equals(detachedContents, ...).
func abstractSignatureParametersDocumentsEqual(a, b []model.DSSDocument) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// abstractSignatureParametersCertificateEquals ports Objects.equals(signingCertificate, ...).
func abstractSignatureParametersCertificateEquals(a, b *model.CertificateToken) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// abstractSignatureParametersCertificateChainEquals ports Objects.equals(certificateChain, ...).
func abstractSignatureParametersCertificateChainEquals(a, b []*model.CertificateToken) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !abstractSignatureParametersCertificateEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

// abstractSignatureParametersTimestampsEqual ports Objects.equals(contentTimestamps, ...).
func abstractSignatureParametersTimestampsEqual(a, b []*validation.TimestampToken) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// compile-time interface assertion.
var _ model.SerializableSignatureParameters = (*AbstractSignatureParameters[model.SerializableTimestampParameters])(nil)
