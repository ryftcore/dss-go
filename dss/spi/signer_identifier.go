// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/SignerIdentifier.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the type
// keeps its Java name unqualified.
package spi

import (
	"bytes"
	"math/big"

	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// SignerIdentifier represents an ASN.1 SignerId DTO.
type SignerIdentifier struct {
	// issuerName is the X500Principal name of the certificate issuer.
	issuerName *model.X500Principal

	// serialNumber is the certificate's serial number.
	serialNumber *big.Int

	// ski is the SHA-1 hash of the certificate's public key (used in OCSP responses).
	ski []byte

	// current tells whether this is the CertificateIdentifier used for a signature/timestamp.
	current bool
}

// NewSignerIdentifier instantiates the object with null values. Port of the default constructor.
func NewSignerIdentifier() *SignerIdentifier {
	return &SignerIdentifier{}
}

// IssuerName returns the name of the certificate issuer. Port of getIssuerName().
func (s *SignerIdentifier) IssuerName() *model.X500Principal {
	return s.issuerName
}

// SetIssuerName sets the name of the certificate's issuer. Port of setIssuerName(X500Principal).
func (s *SignerIdentifier) SetIssuerName(name *model.X500Principal) {
	s.issuerName = name
}

// SerialNumber returns the serial number of the signer certificate. Port of getSerialNumber().
func (s *SignerIdentifier) SerialNumber() *big.Int {
	return s.serialNumber
}

// SetSerialNumber sets the serial number of the signer certificate.
// Port of setSerialNumber(BigInteger).
func (s *SignerIdentifier) SetSerialNumber(serialNumber *big.Int) {
	s.serialNumber = serialNumber
}

// Ski returns the SHA-1 of the certificate's public key. Port of getSki().
func (s *SignerIdentifier) Ski() []byte {
	return s.ski
}

// SetSki sets the SHA-1 of the certificate's public key. Port of setSki(byte[]).
func (s *SignerIdentifier) SetSki(ski []byte) {
	s.ski = ski
}

// IsCurrent reports whether the SignerIdentifier is related to the current signature.
// Port of isCurrent().
func (s *SignerIdentifier) IsCurrent() bool {
	return s.current
}

// SetCurrent sets whether the SignerIdentifier is related to the current signature.
// Port of setCurrent(boolean).
func (s *SignerIdentifier) SetCurrent(current bool) {
	s.current = current
}

// IssuerSerialEncoded returns the DER-encoded IssuerSerial representation of the object, or nil
// when either the issuer name or the serial number is missing. Port of getIssuerSerialEncoded().
//
//	IssuerSerial ::= SEQUENCE { issuer GeneralNames, serial CertificateSerialNumber }
//
// The issuer is wrapped in a single directoryName GeneralName, whose [4] tag is EXPLICIT because
// Name is a CHOICE. The principal's original DER is reused, never re-encoded.
func (s *SignerIdentifier) IssuerSerialEncoded() []byte {
	if s.issuerName == nil || s.serialNumber == nil {
		return nil
	}
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1(cbasn1.SEQUENCE, func(issuerSerial *cryptobyte.Builder) {
		issuerSerial.AddASN1(cbasn1.SEQUENCE, func(generalNames *cryptobyte.Builder) {
			generalNames.AddASN1(cbasn1.Tag(4).ContextSpecific().Constructed(), func(directoryName *cryptobyte.Builder) {
				directoryName.AddBytes(s.issuerName.Encoded())
			})
		})
		issuerSerial.AddASN1BigInt(s.serialNumber)
	})
	encoded, err := builder.Bytes()
	if err != nil {
		return nil
	}
	return encoded
}

// IsRelatedToCertificate reports whether the SignerIdentifier is related to the given
// certificateToken. Port of isRelatedToCertificate(CertificateToken).
//
// The DSSException CertificateExtensionsUtils#getSubjectKeyIdentifier raises on an unreadable
// subjectKeyIdentifier extension becomes a returned error here.
func (s *SignerIdentifier) IsRelatedToCertificate(certificateToken *model.CertificateToken) (bool, error) {
	id := NewSignerIdentifier()
	id.SetIssuerName(certificateToken.IssuerX500Principal())
	id.SetSerialNumber(certificateToken.SerialNumber())
	certSki, err := CertificateExtensionsUtilsSubjectKeyIdentifier(certificateToken)
	if err != nil {
		return false, err
	}
	if certSki != nil {
		id.SetSki(certSki.Ski())
	}
	return s.IsEquivalent(id), nil
}

// IsEquivalent reports whether the given signerIdentifier is equivalent.
// Port of isEquivalent(SignerIdentifier).
func (s *SignerIdentifier) IsEquivalent(signerIdentifier *SignerIdentifier) bool {
	if signerIdentifier == nil {
		return false
	}
	if s.issuerName != nil && s.serialNumber != nil {
		if !DSSASN1UtilsX500PrincipalAreEquals(s.issuerName, signerIdentifier.IssuerName()) {
			return false
		}
		if signerIdentifier.SerialNumber() == nil || s.serialNumber.Cmp(signerIdentifier.SerialNumber()) != 0 {
			return false
		}
		return true
	}
	return bytes.Equal(s.ski, signerIdentifier.Ski())
}

// IsEmpty reports whether the SignerIdentifier carries no value at all. Port of isEmpty().
func (s *SignerIdentifier) IsEmpty() bool {
	return s.issuerName == nil && s.serialNumber == nil && utils.IsArrayEmpty(s.ski)
}

// String returns a string representation of the identifier. Port of toString().
//
// DEVIATION: Java interpolates X500Principal#toString(), which emits the RFC 1779 flavoured
// "DEFAULT" format; model.X500Principal deliberately renders the RFC 2253 form instead (see its
// own file header), so the issuerName part of this string differs from upstream's.
func (s *SignerIdentifier) String() string {
	if s.issuerName != nil || s.serialNumber != nil {
		return "IssuerSerialInfo [issuerName=" + signerIdentifierToString(s.issuerName) +
			", serialNumber=" + signerIdentifierSerialToString(s.serialNumber) + "]"
	}
	return "IssuerSerialInfo [ski=" + utils.ToBase64(s.ski) + "]"
}

// signerIdentifierToString renders a possibly missing principal the way Java's string
// concatenation does.
func signerIdentifierToString(principal *model.X500Principal) string {
	if principal == nil {
		return "null"
	}
	return principal.String()
}

// signerIdentifierSerialToString renders a possibly missing serial number the way Java's string
// concatenation does.
func signerIdentifierSerialToString(serialNumber *big.Int) string {
	if serialNumber == nil {
		return "null"
	}
	return serialNumber.String()
}

// Equals reports whether both identifiers carry the same issuer name, serial number and SKI.
// Port of equals(Object).
//
// NOTE: hashCode() has no Go counterpart; upstream needs it only to key the JDK hash collections,
// which the port replaces with slices keyed on Equals.
func (s *SignerIdentifier) Equals(other *SignerIdentifier) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	if s.issuerName == nil {
		if other.issuerName != nil {
			return false
		}
	} else if !s.issuerName.Equals(other.issuerName) {
		return false
	}
	if s.serialNumber == nil {
		if other.serialNumber != nil {
			return false
		}
	} else if other.serialNumber == nil || s.serialNumber.Cmp(other.serialNumber) != 0 {
		return false
	}
	return bytes.Equal(s.ski, other.ski)
}
