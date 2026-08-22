// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CandidatesForSigningCertificate.java (DSS 6.5.RC1).
package spi

import "github.com/ryftcore/dss-go/dss/model"

// CandidatesForSigningCertificate holds the list of the candidates for the signing
// certificate of the main signature.
type CandidatesForSigningCertificate struct {
	// theCertificateValidity is the reference to the signing certificate with its validity.
	// Set after the signature verification.
	theCertificateValidity *CertificateValidity

	// certificateValidityList contains the candidates for the signing certificate.
	certificateValidityList []*CertificateValidity
}

// NewCandidatesForSigningCertificate instantiates the object with null or empty values.
// Port of the default constructor.
func NewCandidatesForSigningCertificate() *CandidatesForSigningCertificate {
	return &CandidatesForSigningCertificate{}
}

// CertificateValidityList gets the list of candidates for the signing certificate.
// Port of getCertificateValidityList().
func (c *CandidatesForSigningCertificate) CertificateValidityList() []*CertificateValidity {
	return c.certificateValidityList
}

// IsEmpty tests if any candidate is known. Port of isEmpty().
func (c *CandidatesForSigningCertificate) IsEmpty() bool {
	return len(c.certificateValidityList) == 0
}

// Add allows to add a candidate for the signing certificate.
// Port of add(CertificateValidity).
func (c *CandidatesForSigningCertificate) Add(certificateValidity *CertificateValidity) {
	c.certificateValidityList = append(c.certificateValidityList, certificateValidity)
}

// SetTheCertificateValidity allows to set the CertificateValidity object after the
// verification of its signature. theCertificateValidity must be in the list of the
// candidates. Port of setTheCertificateValidity(CertificateValidity).
//
// Panics with the Java message when theCertificateValidity is missing (Objects.requireNonNull)
// and returns an error when it is not part of the candidates.
func (c *CandidatesForSigningCertificate) SetTheCertificateValidity(theCertificateValidity *CertificateValidity) error {
	if theCertificateValidity == nil {
		panic("The CertificateValidity cannot be null")
	}
	if !candidatesForSigningCertificateContains(c.certificateValidityList, theCertificateValidity) {
		return model.NewDSSError("theSigningCertificateValidity must be the part of the candidates!")
	}
	c.theCertificateValidity = theCertificateValidity
	return nil
}

// TheCertificateValidity returns the signing certificate validity, or nil if such a
// certificate was not identified. TheCertificateValidity must be set before.
// Port of getTheCertificateValidity().
func (c *CandidatesForSigningCertificate) TheCertificateValidity() *CertificateValidity {
	return c.theCertificateValidity
}

// TheBestCandidate returns the best candidate for the signing certificate. The only way to
// be sure that it is the right one is to validate the signature.
// Port of getTheBestCandidate().
func (c *CandidatesForSigningCertificate) TheBestCandidate() *CertificateValidity {
	var firstCandidate *CertificateValidity
	for _, certificateValidity := range c.certificateValidityList {
		if firstCandidate == nil {
			firstCandidate = certificateValidity
		}
		if certificateValidity.IsValid() {
			return certificateValidity
		}
	}

	if signerIdMatchCandidate := c.bySignerIdMatch(); signerIdMatchCandidate != nil {
		return signerIdMatchCandidate
	}

	if c.theCertificateValidity != nil {
		return c.theCertificateValidity
	}
	return firstCandidate
}

// bySignerIdMatch returns the signing certificate which was identified with the CMS SID.
// Port of the private getBySignerIdMatch().
func (c *CandidatesForSigningCertificate) bySignerIdMatch() *CertificateValidity {
	for _, certificateValidity := range c.certificateValidityList {
		if certificateValidity.IsSignerIdMatch() {
			return certificateValidity
		}
	}
	return nil
}

// candidatesForSigningCertificateContains reports whether the list contains the given
// candidate, standing in for Java's List#contains (identity-based equality; CertificateValidity
// has no Equals override in the Java source either).
func candidatesForSigningCertificateContains(list []*CertificateValidity, candidate *CertificateValidity) bool {
	for _, item := range list {
		if item == candidate {
			return true
		}
	}
	return false
}
