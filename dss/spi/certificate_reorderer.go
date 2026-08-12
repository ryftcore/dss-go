// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateReorderer.java (DSS 6.5.RC1).
package spi

import (
	"github.com/utain/esig/dss/model"
)

// CertificateReorderer reorders a certificate collection to the corresponding certificate
// chain.
type CertificateReorderer struct {
	// signingCertificate is the signing certificate (the last certificate in the chain).
	signingCertificate *model.CertificateToken

	// certificateChain is the collection of certificates.
	certificateChain []*model.CertificateToken
}

// NewCertificateReorderer builds a reorderer over a collection of certificates where DSS
// needs to find the signing certificate.
// Port of the CertificateReorderer(Collection<CertificateToken>) constructor.
func NewCertificateReorderer(certificateChain []*model.CertificateToken) *CertificateReorderer {
	return NewCertificateReordererWithSigningCertificate(nil, certificateChain)
}

// NewCertificateReordererWithSigningCertificate builds a reorderer over a potential signing
// certificate and a certificate chain.
// Port of the CertificateReorderer(CertificateToken, Collection<CertificateToken>) constructor.
func NewCertificateReordererWithSigningCertificate(signingCertificate *model.CertificateToken,
	certificateChain []*model.CertificateToken) *CertificateReorderer {
	return &CertificateReorderer{signingCertificate: signingCertificate, certificateChain: certificateChain}
}

// OrderedCertificates orders the certificates (signing certificate, CA1, CA2 and Root).
// Port of getOrderedCertificates(). Java's DSSExceptions ("No signing certificate found",
// "The certificate chain contains only bridge certificates", "Unable to determine a signing
// certificate : No pertinent input parameters") are returned as errors.
func (r *CertificateReorderer) OrderedCertificates() ([]*model.CertificateToken, error) {
	certificates := r.allCertificatesOnce()
	if len(certificates) == 1 {
		return certificates, nil
	}

	r.initIssuerPublicKeys(certificates)

	identifiedSigningCerts, err := r.signingCertificates(certificates)
	if err != nil {
		return nil, err
	}
	selectedSigningCert, err := r.selectSigningCertificateInList(identifiedSigningCerts)
	if err != nil {
		return nil, err
	}

	rebuiltCertificateChain := certificateReordererBuildCertificateChainForCert(certificates, selectedSigningCert)
	// Some certificates may be ignored when len(certificates) > len(rebuiltCertificateChain);
	// upstream logs the before/after lists here, dropped per PORTING.md's slf4j rule.
	return rebuiltCertificateChain, nil
}

// OrderedCertificateChains orders the certificates (signing certificate, CA1, CA2 and Root)
// into one or more ordered chains, keyed by the identified signing certificate.
// Port of getOrderedCertificateChains().
func (r *CertificateReorderer) OrderedCertificateChains() (map[*model.CertificateToken][]*model.CertificateToken, error) {
	result := make(map[*model.CertificateToken][]*model.CertificateToken)

	certificates := r.allCertificatesOnce()
	if len(certificates) == 1 {
		uniqueCert := certificates[0]
		result[uniqueCert] = []*model.CertificateToken{uniqueCert}
		return result, nil
	}

	r.initIssuerPublicKeys(certificates)

	identifiedSigningCerts, err := r.signingCertificates(certificates)
	if err != nil {
		return nil, err
	}
	for _, identifiedSigningCert := range identifiedSigningCerts {
		result[identifiedSigningCert] = certificateReordererBuildCertificateChainForCert(certificates, identifiedSigningCert)
	}
	return result, nil
}

// initIssuerPublicKeys builds the chain cert -> issuer. Port of the private
// initIssuerPublicKeys(List<CertificateToken>).
func (r *CertificateReorderer) initIssuerPublicKeys(certificates []*model.CertificateToken) {
	for _, token := range certificates {
		if certificateReordererIsIssuerNeeded(token) {
			for _, signer := range certificates {
				if token.IsSignedByToken(signer) {
					break
				}
			}
			// Java logs a warning here if the issuer is still not found; dropped per
			// PORTING.md's slf4j rule.
		}
	}
}

// certificateReordererBuildCertificateChainForCert ports the private
// buildCertificateChainForCert(List<CertificateToken>, CertificateToken).
func certificateReordererBuildCertificateChainForCert(certificates []*model.CertificateToken,
	certToAdd *model.CertificateToken) []*model.CertificateToken {
	var result []*model.CertificateToken
	for certToAdd != nil && !certificateReordererContains(result, certToAdd) {
		result = append(result, certToAdd)
		certToAdd = certificateReordererByPubKey(certificates, certToAdd.PublicKeyOfTheSigner())
	}
	return result
}

// certificateReordererByPubKey ports the private getCertificateByPubKey(List<CertificateToken>, PublicKey).
func certificateReordererByPubKey(certificates []*model.CertificateToken, publicKeyOfTheSigner *model.PublicKey) *model.CertificateToken {
	if publicKeyOfTheSigner == nil {
		return nil
	}
	for _, certificateToken := range certificates {
		if certificateToken.PublicKey().Equals(publicKeyOfTheSigner) {
			return certificateToken
		}
	}
	return nil
}

// certificateReordererIsIssuerNeeded ports the private isIssuerNeeded(CertificateToken).
func certificateReordererIsIssuerNeeded(token *model.CertificateToken) bool {
	return !token.IsSelfSigned() && token.PublicKeyOfTheSigner() == nil
}

// selectSigningCertificateInList ports the private selectSigningCertificateInList(List<CertificateToken>).
func (r *CertificateReorderer) selectSigningCertificateInList(identifiedSigningCerts []*model.CertificateToken) (*model.CertificateToken, error) {
	if len(identifiedSigningCerts) == 1 {
		return identifiedSigningCerts[0], nil
	}
	// More than one chain detected (upstream logs a warning; dropped per PORTING.md).
	if r.signingCertificate != nil && certificateReordererContains(identifiedSigningCerts, r.signingCertificate) {
		return r.signingCertificate, nil
	}
	return nil, model.NewDSSError("Unable to determine a signing certificate : No pertinent input parameters")
}

// allCertificatesOnce ports the private getAllCertificatesOnce(), avoiding duplicate entries.
func (r *CertificateReorderer) allCertificatesOnce() []*model.CertificateToken {
	var result []*model.CertificateToken
	if r.signingCertificate != nil {
		result = append(result, r.signingCertificate)
	}
	for _, certificateToken := range r.certificateChain {
		if certificateToken != nil && !certificateReordererContains(result, certificateToken) {
			result = append(result, certificateToken)
		}
	}
	return result
}

// signingCertificates identifies the signing certificates (the certificates which did not
// sign any other certificate). Port of the private getSigningCertificates(List<CertificateToken>).
func (r *CertificateReorderer) signingCertificates(certificates []*model.CertificateToken) ([]*model.CertificateToken, error) {
	if len(certificates) == 0 {
		return nil, model.NewDSSError("No signing certificate found")
	}

	var potentialSigners []*model.CertificateToken
	for _, signer := range certificates {
		caCert := false
		for _, token := range certificates {
			if signer.PublicKey().Equals(token.PublicKeyOfTheSigner()) {
				caCert = true
				break
			}
		}
		if !caCert {
			potentialSigners = append(potentialSigners, signer)
		}
	}

	if len(potentialSigners) == 0 {
		return nil, model.NewDSSError("The certificate chain contains only bridge certificates")
	}
	return potentialSigners, nil
}

// certificateReordererContains reports whether the list contains the given certificate token.
func certificateReordererContains(list []*model.CertificateToken, token *model.CertificateToken) bool {
	for _, item := range list {
		if item == token {
			return true
		}
	}
	return false
}
