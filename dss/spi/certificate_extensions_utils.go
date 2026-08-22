// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/CertificateExtensionsUtils.java (DSS 6.5.RC1).
//
// Naming: the Java class is a holder of static methods. eu.europa.esig.dss.spi,
// spi.x509 and spi.x509.revocation.* all flatten into the single Go package spi, and the
// method names of the different utility classes collide there (QcStatementUtils exports
// getQcStatements() as well), so every exported function keeps the owning Java class name
// as its prefix, following the "static factories keep their contract" rule of PORTING.md
// (SignatureAlgorithm#forOID -> SignatureAlgorithmForOID).
//
// BouncyCastle replacement: upstream loads every certificate through the BouncyCastle JCE
// provider (DSSUtils -> DSSCertificateTokenSecurityFactory -> DSSSecurityProvider), therefore
// X509Certificate#getSubjectAlternativeNames(), #getExtensionValue(), #getKeyUsage(),
// #getExtendedKeyUsage() and #getBasicConstraints() are BouncyCastle's implementations and NOT
// the JDK's. The code below reproduces BouncyCastle's behaviour; the two differ substantially
// for subjectAlternativeName (directoryName rendering, iPAddress rendering, otherName bytes).
// ASN.1 is read with cryptobyte so the original DER is preserved wherever the model exposes it.
package spi

import (
	"crypto/x509"
	encasn1 "encoding/asn1"
	"fmt"
	"math/big"
	"strings"
	"unicode/utf16"

	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/x509/extension"
	"github.com/ryftcore/dss-go/dss/utils"
)

// certificateExtensionsUtilsOIDIdADCaIssuers is X509ObjectIdentifiers.id_ad_caIssuers.
const certificateExtensionsUtilsOIDIdADCaIssuers = "1.3.6.1.5.5.7.48.2"

// certificateExtensionsUtilsOIDIdADOCSP is X509ObjectIdentifiers.id_ad_ocsp.
const certificateExtensionsUtilsOIDIdADOCSP = "1.3.6.1.5.5.7.48.1"

// certificateExtensionsUtilsOIDIdQtCps is PolicyQualifierId.id_qt_cps.
const certificateExtensionsUtilsOIDIdQtCps = "1.3.6.1.5.5.7.2.1"

// CertificateExtensionsUtilsCertificateExtensions extracts the certificate extensions from the
// given certificateToken. Port of getCertificateExtensions(CertificateToken).
//
// DEVIATION: Java walks getCriticalExtensionOIDs() then getNonCriticalExtensionOIDs(), which are
// BouncyCastle HashSets, so the resulting order of CertificateExtensions#getAllCertificateExtensions()
// and #getOtherExtensions() is the (arbitrary but stable) hash order of the OID strings. The Go
// port walks the certificate's extensions in DER order, critical ones first, which yields the same
// contents in a deterministic order.
func CertificateExtensionsUtilsCertificateExtensions(certificateToken *model.CertificateToken) (*extension.CertificateExtensions, error) {
	certificateExtensions := extension.NewCertificateExtensions()
	for _, critical := range []bool{true, false} {
		if err := certificateExtensionsUtilsSetCertificateExtensions(certificateExtensions, certificateToken, critical); err != nil {
			return nil, err
		}
	}
	return certificateExtensions, nil
}

// certificateExtensionsUtilsSetCertificateExtensions ports the private
// setCertificateExtensions(CertificateExtensions, CertificateToken, Collection<String>) for the
// critical or the non-critical extension OIDs of the token.
func certificateExtensionsUtilsSetCertificateExtensions(certificateExtensions *extension.CertificateExtensions,
	certificateToken *model.CertificateToken, critical bool) error {
	for _, ext := range certificateToken.Certificate().Extensions {
		if ext.Critical != critical {
			continue
		}
		oid := ext.Id.String()
		switch {
		case CertificateExtensionsUtilsIsSubjectAlternativeNames(oid):
			certificateExtensions.SetSubjectAlternativeNames(CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken))
		case CertificateExtensionsUtilsIsAuthorityKeyIdentifier(oid):
			authorityKeyIdentifier, err := CertificateExtensionsUtilsAuthorityKeyIdentifier(certificateToken)
			if err != nil {
				return err
			}
			certificateExtensions.SetAuthorityKeyIdentifier(authorityKeyIdentifier)
		case CertificateExtensionsUtilsIsSubjectKeyIdentifier(oid):
			subjectKeyIdentifier, err := CertificateExtensionsUtilsSubjectKeyIdentifier(certificateToken)
			if err != nil {
				return err
			}
			certificateExtensions.SetSubjectKeyIdentifier(subjectKeyIdentifier)
		case CertificateExtensionsUtilsIsAuthorityInformationAccess(oid):
			certificateExtensions.SetAuthorityInformationAccess(CertificateExtensionsUtilsAuthorityInformationAccess(certificateToken))
		case CertificateExtensionsUtilsIsCRLDistributionPoints(oid):
			certificateExtensions.SetCRLDistributionPoints(CertificateExtensionsUtilsCRLDistributionPoints(certificateToken))
		case CertificateExtensionsUtilsIsBasicConstraints(oid):
			certificateExtensions.SetBasicConstraints(CertificateExtensionsUtilsBasicConstraints(certificateToken))
		case CertificateExtensionsUtilsIsNameConstraints(oid):
			certificateExtensions.SetNameConstraints(CertificateExtensionsUtilsNameConstraints(certificateToken))
		case CertificateExtensionsUtilsIsPolicyConstraints(oid):
			certificateExtensions.SetPolicyConstraints(CertificateExtensionsUtilsPolicyConstraints(certificateToken))
		case CertificateExtensionsUtilsIsInhibitAnyPolicy(oid):
			certificateExtensions.SetInhibitAnyPolicy(CertificateExtensionsUtilsInhibitAnyPolicy(certificateToken))
		case CertificateExtensionsUtilsIsFreshestCRL(oid):
			certificateExtensions.SetFreshestCRL(CertificateExtensionsUtilsFreshestCRL(certificateToken))
		case CertificateExtensionsUtilsIsKeyUsage(oid):
			certificateExtensions.SetKeyUsage(CertificateExtensionsUtilsKeyUsage(certificateToken))
		case CertificateExtensionsUtilsIsExtendedKeyUsage(oid):
			certificateExtensions.SetExtendedKeyUsage(CertificateExtensionsUtilsExtendedKeyUsage(certificateToken))
		case CertificateExtensionsUtilsIsCertificatePolicies(oid):
			certificateExtensions.SetCertificatePolicies(CertificateExtensionsUtilsCertificatePolicies(certificateToken))
		case CertificateExtensionsUtilsIsOcspNoCheck(oid):
			certificateExtensions.SetOcspNoCheck(CertificateExtensionsUtilsOcspNoCheck(certificateToken))
		case CertificateExtensionsUtilsIsValidityAssuredShortTerm(oid):
			certificateExtensions.SetValidityAssuredShortTerm(CertificateExtensionsUtilsValAssuredSTCerts(certificateToken))
		case CertificateExtensionsUtilsIsQcStatements(oid):
			certificateExtensions.SetQcStatements(CertificateExtensionsUtilsQcStatements(certificateToken))
		case CertificateExtensionsUtilsIsNoRevocationAvailable(oid):
			certificateExtensions.SetNoRevAvail(CertificateExtensionsUtilsNoRevAvail(certificateToken))
		default:
			certificateExtensions.AddOtherExtension(certificateExtensionsUtilsOtherCertificateExtension(certificateToken, oid))
		}
	}
	return nil
}

// CertificateExtensionsUtilsIsSubjectAlternativeNames reports whether oid is the subject
// alternative names extension OID. Port of isSubjectAlternativeNames(String).
func CertificateExtensionsUtilsIsSubjectAlternativeNames(oid string) bool {
	return enumerations.CertificateExtensionEnumSubjectAlternativeName.OID() == oid
}

// CertificateExtensionsUtilsIsAuthorityKeyIdentifier reports whether oid is the authority key
// identifier extension OID. Port of isAuthorityKeyIdentifier(String).
func CertificateExtensionsUtilsIsAuthorityKeyIdentifier(oid string) bool {
	return enumerations.CertificateExtensionEnumAuthorityKeyIdentifier.OID() == oid
}

// CertificateExtensionsUtilsIsSubjectKeyIdentifier reports whether oid is the subject key
// identifier extension OID. Port of isSubjectKeyIdentifier(String).
func CertificateExtensionsUtilsIsSubjectKeyIdentifier(oid string) bool {
	return enumerations.CertificateExtensionEnumSubjectKeyIdentifier.OID() == oid
}

// CertificateExtensionsUtilsIsAuthorityInformationAccess reports whether oid is the authority
// information access extension OID. Port of isAuthorityInformationAccess(String).
func CertificateExtensionsUtilsIsAuthorityInformationAccess(oid string) bool {
	return enumerations.CertificateExtensionEnumAuthorityInformationAccess.OID() == oid
}

// CertificateExtensionsUtilsIsCRLDistributionPoints reports whether oid is the CRL distribution
// points extension OID. Port of isCRLDistributionPoints(String).
func CertificateExtensionsUtilsIsCRLDistributionPoints(oid string) bool {
	return enumerations.CertificateExtensionEnumCRLDistributionPoints.OID() == oid
}

// CertificateExtensionsUtilsIsBasicConstraints reports whether oid is the basic constraints
// extension OID. Port of isBasicConstraints(String).
func CertificateExtensionsUtilsIsBasicConstraints(oid string) bool {
	return enumerations.CertificateExtensionEnumBasicConstraints.OID() == oid
}

// CertificateExtensionsUtilsIsNameConstraints reports whether oid is the name constraints
// extension OID. Port of isNameConstraints(String).
func CertificateExtensionsUtilsIsNameConstraints(oid string) bool {
	return enumerations.CertificateExtensionEnumNameConstraints.OID() == oid
}

// CertificateExtensionsUtilsIsPolicyConstraints reports whether oid is the policy constraints
// extension OID. Port of isPolicyConstraints(String).
func CertificateExtensionsUtilsIsPolicyConstraints(oid string) bool {
	return enumerations.CertificateExtensionEnumPolicyConstraints.OID() == oid
}

// CertificateExtensionsUtilsIsKeyUsage reports whether oid is the key usage extension OID.
// Port of isKeyUsage(String).
func CertificateExtensionsUtilsIsKeyUsage(oid string) bool {
	return enumerations.CertificateExtensionEnumKeyUsage.OID() == oid
}

// CertificateExtensionsUtilsIsExtendedKeyUsage reports whether oid is the extended key usage
// extension OID. Port of isExtendedKeyUsage(String).
func CertificateExtensionsUtilsIsExtendedKeyUsage(oid string) bool {
	return enumerations.CertificateExtensionEnumExtendedKeyUsage.OID() == oid
}

// CertificateExtensionsUtilsIsInhibitAnyPolicy reports whether oid is the inhibit anyPolicy
// extension OID. Port of isInhibitAnyPolicy(String).
func CertificateExtensionsUtilsIsInhibitAnyPolicy(oid string) bool {
	return enumerations.CertificateExtensionEnumInhibitAnyPolicy.OID() == oid
}

// CertificateExtensionsUtilsIsFreshestCRL reports whether oid is the Freshest CRL (a.k.a. Delta
// CRL) extension OID. Port of isFreshestCRL(String).
func CertificateExtensionsUtilsIsFreshestCRL(oid string) bool {
	return enumerations.CertificateExtensionEnumFreshestCRL.OID() == oid
}

// CertificateExtensionsUtilsIsCertificatePolicies reports whether oid is the certificate policies
// extension OID. Port of isCertificatePolicies(String).
func CertificateExtensionsUtilsIsCertificatePolicies(oid string) bool {
	return enumerations.CertificateExtensionEnumCertificatePolicies.OID() == oid
}

// CertificateExtensionsUtilsIsOcspNoCheck reports whether oid is the ocsp-nocheck extension OID.
// Port of isOcspNoCheck(String).
func CertificateExtensionsUtilsIsOcspNoCheck(oid string) bool {
	return enumerations.CertificateExtensionEnumOCSPNoCheck.OID() == oid
}

// CertificateExtensionsUtilsIsValidityAssuredShortTerm reports whether oid is the
// ext-etsi-valassured-ST-certs extension OID. Port of isValidityAssuredShortTerm(String).
func CertificateExtensionsUtilsIsValidityAssuredShortTerm(oid string) bool {
	return enumerations.CertificateExtensionEnumValidityAssuredShortTerm.OID() == oid
}

// CertificateExtensionsUtilsIsQcStatements reports whether oid is the qc-statements extension OID.
// Port of isQcStatements(String).
func CertificateExtensionsUtilsIsQcStatements(oid string) bool {
	return enumerations.CertificateExtensionEnumQCStatements.OID() == oid
}

// CertificateExtensionsUtilsIsNoRevocationAvailable reports whether oid is the noRevAvail
// extension OID. Port of isNoRevocationAvailable(String).
func CertificateExtensionsUtilsIsNoRevocationAvailable(oid string) bool {
	return enumerations.CertificateExtensionEnumNoRevocationAvailable.OID() == oid
}

// CertificateExtensionsUtilsSubjectAlternativeNames returns the subject alternative names, when
// present, or nil when the extension cannot be read. Port of getSubjectAlternativeNames(CertificateToken).
func CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken *model.CertificateToken) *extension.SubjectAlternativeNames {
	certificate := certificateToken.Certificate()
	subjectAlternateNames := extension.NewSubjectAlternativeNames()
	subjectAlternateNames.SetOctets(certificateExtensionsUtilsExtensionValue(certificate, subjectAlternateNames.OID()))

	alternativeNames, ok := certificateExtensionsUtilsAlternativeNames(certificate, subjectAlternateNames.OID())
	if !ok {
		// BouncyCastle raises a CertificateParsingException, which upstream logs and swallows.
		return nil
	}
	result := make([]*extension.GeneralName, 0, len(alternativeNames))
	for _, alternativeName := range alternativeNames {
		generalName := certificateExtensionsUtilsGeneralNameOf(alternativeName)
		if generalName != nil {
			result = append(result, generalName)
		}
	}
	subjectAlternateNames.SetGeneralNames(result)
	subjectAlternateNames.CheckCritical(certificateToken)
	return subjectAlternateNames
}

// certificateExtensionsUtilsAltName is one entry of the Collection<List<?>> returned by
// BouncyCastle's X509Certificate#getSubjectAlternativeNames(): a GeneralName tag number plus
// either a String or a byte array.
type certificateExtensionsUtilsAltName struct {
	// tagNo is the GeneralName tag number.
	tagNo int
	// value holds the String entry BouncyCastle produces for the string-valued general names.
	value string
	// binaries holds the byte[] entry it produces for otherName, x400Address and ediPartyName.
	binaries []byte
	// directoryName holds the retyped Name of a directoryName entry; see
	// certificateExtensionsUtilsRFC4519Name for why the intermediate string is not materialised.
	directoryName []byte
}

// certificateExtensionsUtilsAlternativeNames ports BouncyCastle's
// X509CertificateImpl#getAlternativeNames(Certificate, String). It reports ok=false where
// BouncyCastle raises a CertificateParsingException, and an empty result where it returns null.
func certificateExtensionsUtilsAlternativeNames(certificate *x509.Certificate, oid string) ([]certificateExtensionsUtilsAltName, bool) {
	extensionContent := certificateExtensionsUtilsExtensionContent(certificate, oid)
	if extensionContent == nil {
		return nil, true
	}
	elements, parsed := certificateExtensionsUtilsSequenceElements(extensionContent)
	if !parsed {
		return nil, false
	}
	result := make([]certificateExtensionsUtilsAltName, 0, len(elements))
	for _, element := range elements {
		generalName, ok := certificateExtensionsUtilsParseTagged(element)
		if !ok {
			return nil, false
		}
		switch generalName.tagNo {
		case 0, 3, 5: // otherName, x400Address, ediPartyName: the encoded GeneralName
			result = append(result, certificateExtensionsUtilsAltName{tagNo: generalName.tagNo, binaries: generalName.full})
		case 4: // directoryName: X500Name.getInstance(RFC4519Style.INSTANCE, name).toString()
			retyped, converted := certificateExtensionsUtilsRFC4519Name(generalName.content)
			if !converted {
				return nil, false
			}
			result = append(result, certificateExtensionsUtilsAltName{tagNo: 4, directoryName: retyped})
		case 1, 2, 6: // rfc822Name, dNSName, uniformResourceIdentifier
			if generalName.constructed {
				return nil, false
			}
			result = append(result, certificateExtensionsUtilsAltName{tagNo: generalName.tagNo, value: string(generalName.content)})
		case 8: // registeredID
			objectIdentifier, ok := certificateExtensionsUtilsObjectIdentifier(generalName.content)
			if !ok {
				return nil, false
			}
			result = append(result, certificateExtensionsUtilsAltName{tagNo: 8, value: objectIdentifier})
		case 7: // iPAddress; an address of an unsupported length is skipped by BouncyCastle
			address, ok := certificateExtensionsUtilsHostAddress(generalName.content)
			if !ok {
				continue
			}
			result = append(result, certificateExtensionsUtilsAltName{tagNo: 7, value: address})
		default:
			return nil, false // BouncyCastle: IOException("Bad tag number")
		}
	}
	return result, true
}

// certificateExtensionsUtilsGeneralNameOf ports the private getGeneralName(List<?>).
func certificateExtensionsUtilsGeneralNameOf(altName certificateExtensionsUtilsAltName) *extension.GeneralName {
	generalName := extension.NewGeneralName()
	generalNameType := enumerations.GeneralNameTypeFromIndex(altName.tagNo)
	generalName.SetGeneralNameType(generalNameType)
	if altName.binaries != nil {
		generalName.SetValue(certificateExtensionsUtilsToHexEncoded(altName.binaries))
		return generalName
	}
	value := altName.value
	if generalNameType == enumerations.GeneralNameTypeDirectoryName {
		value = certificateExtensionsUtilsToRFC2253RDN(altName.directoryName)
	}
	generalName.SetValue(value)
	return generalName
}

// CertificateExtensionsUtilsAuthorityInformationAccess returns the authority information access,
// when present. Port of getAuthorityInformationAccess(CertificateToken).
func CertificateExtensionsUtilsAuthorityInformationAccess(certificateToken *model.CertificateToken) *extension.AuthorityInformationAccess {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumAuthorityInformationAccess.OID()
	authInfoAccessExtensionValue := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(authInfoAccessExtensionValue) == 0 {
		return nil
	}

	sequence, ok := certificateExtensionsUtilsSequenceFromDEROctetString(authInfoAccessExtensionValue)
	if !ok || len(sequence) == 0 {
		return nil
	}

	caIssuers := make([]string, 0)
	ocsp := make([]string, 0)
	for _, accessDescription := range sequence {
		// AccessDescription ::= SEQUENCE { accessMethod OBJECT IDENTIFIER, accessLocation GeneralName }
		fields, parsed := certificateExtensionsUtilsSequenceElements(accessDescription)
		if !parsed || len(fields) != 2 {
			return nil
		}
		accessMethod, ok := certificateExtensionsUtilsObjectIdentifierElement(fields[0])
		if !ok {
			return nil
		}
		generalName, ok := certificateExtensionsUtilsParseTagged(fields[1])
		if !ok {
			return nil
		}
		location, isURI := certificateExtensionsUtilsParseGn(generalName)
		if !isURI {
			continue
		}
		switch accessMethod {
		case certificateExtensionsUtilsOIDIdADCaIssuers:
			caIssuers = append(caIssuers, location)
		case certificateExtensionsUtilsOIDIdADOCSP:
			ocsp = append(ocsp, location)
		}
	}

	authorityInformationAccess := extension.NewAuthorityInformationAccess()
	authorityInformationAccess.SetOctets(authInfoAccessExtensionValue)
	authorityInformationAccess.SetCaIssuers(caIssuers)
	authorityInformationAccess.SetOcsp(ocsp)
	authorityInformationAccess.CheckCritical(certificateToken)
	return authorityInformationAccess
}

// CertificateExtensionsUtilsCAIssuersAccessUrls returns the CA issuers URIs extracted from the
// authorityInfoAccess.caIssuers field, or an empty slice. Port of getCAIssuersAccessUrls(CertificateToken).
func CertificateExtensionsUtilsCAIssuersAccessUrls(certificate *model.CertificateToken) []string {
	aia := CertificateExtensionsUtilsAuthorityInformationAccess(certificate)
	if aia != nil {
		return aia.CaIssuers()
	}
	return []string{}
}

// CertificateExtensionsUtilsOCSPAccessUrls returns the OCSP URIs extracted from the
// authorityInfoAccess.ocsp field, or an empty slice. Port of getOCSPAccessUrls(CertificateToken).
func CertificateExtensionsUtilsOCSPAccessUrls(certificate *model.CertificateToken) []string {
	aia := CertificateExtensionsUtilsAuthorityInformationAccess(certificate)
	if aia != nil {
		return aia.Ocsp()
	}
	return []string{}
}

// CertificateExtensionsUtilsAuthorityKeyIdentifier returns the authority key identifier, when
// present. Port of getAuthorityKeyIdentifier(CertificateToken).
//
// The Java method wraps the IOException of JcaX509ExtensionUtils#parseExtensionValue in a
// DSSException; here that becomes a returned model.DSSError. A structurally invalid
// AuthorityKeyIdentifier, which Java would surface as an uncaught IllegalArgumentException, is
// reported through the same error.
func CertificateExtensionsUtilsAuthorityKeyIdentifier(certificateToken *model.CertificateToken) (*extension.AuthorityKeyIdentifier, error) {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumAuthorityKeyIdentifier.OID()
	extensionValue := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(extensionValue) == 0 {
		return nil, nil
	}

	content := certificateExtensionsUtilsExtensionContent(certificate, oid)
	elements, parsed := certificateExtensionsUtilsSequenceElements(content)
	if !parsed {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to retrieve authority key identifier of a certificate. "+
			"Reason : %s", "unable to parse the AuthorityKeyIdentifier sequence"))
	}

	// AuthorityKeyIdentifier ::= SEQUENCE { [0] keyIdentifier OCTET STRING OPTIONAL,
	//     [1] authorityCertIssuer GeneralNames OPTIONAL, [2] authorityCertSerialNumber INTEGER OPTIONAL }
	var keyIdentifier []byte
	var authorityCertIssuer []byte
	var authorityCertSerialNumber []byte
	for _, element := range elements {
		tagged, ok := certificateExtensionsUtilsParseTagged(element)
		if !ok {
			return nil, model.NewDSSError(fmt.Sprintf("Unable to retrieve authority key identifier of a certificate. "+
				"Reason : %s", "unexpected AuthorityKeyIdentifier element"))
		}
		switch tagged.tagNo {
		case 0:
			keyIdentifier = tagged.content
		case 1:
			authorityCertIssuer = tagged.content
		case 2:
			authorityCertSerialNumber = tagged.content
		default:
			return nil, model.NewDSSError(fmt.Sprintf("Unable to retrieve authority key identifier of a certificate. "+
				"Reason : %s", "illegal tag"))
		}
	}

	authorityKeyIdentifier := extension.NewAuthorityKeyIdentifier()
	authorityKeyIdentifier.SetOctets(extensionValue)
	authorityKeyIdentifier.SetKeyIdentifier(keyIdentifier)
	if authorityCertIssuer != nil && authorityCertSerialNumber != nil {
		// IssuerSerial ::= SEQUENCE { issuer GeneralNames, serial CertificateSerialNumber }
		builder := cryptobyte.NewBuilder(nil)
		builder.AddASN1(cbasn1.SEQUENCE, func(issuerSerial *cryptobyte.Builder) {
			issuerSerial.AddASN1(cbasn1.SEQUENCE, func(generalNames *cryptobyte.Builder) {
				generalNames.AddBytes(authorityCertIssuer)
			})
			issuerSerial.AddASN1(cbasn1.INTEGER, func(serialNumber *cryptobyte.Builder) {
				serialNumber.AddBytes(authorityCertSerialNumber)
			})
		})
		encoded, err := builder.Bytes()
		if err != nil {
			return nil, model.NewDSSErrorMessageCause("Unable to retrieve authority key identifier of a certificate.", err)
		}
		authorityKeyIdentifier.SetAuthorityCertIssuerSerial(encoded)
	}
	authorityKeyIdentifier.CheckCritical(certificateToken)
	return authorityKeyIdentifier, nil
}

// CertificateExtensionsUtilsSubjectKeyIdentifier returns the subject key identifier, when present.
// Port of getSubjectKeyIdentifier(CertificateToken).
//
// As for the authority key identifier, the Java DSSException becomes a returned model.DSSError.
func CertificateExtensionsUtilsSubjectKeyIdentifier(certificateToken *model.CertificateToken) (*extension.SubjectKeyIdentifier, error) {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumSubjectKeyIdentifier.OID()
	extensionValue := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(extensionValue) == 0 {
		return nil, nil
	}

	// SubjectKeyIdentifier ::= KeyIdentifier ::= OCTET STRING
	content := certificateExtensionsUtilsExtensionContent(certificate, oid)
	input := cryptobyte.String(content)
	var ski cryptobyte.String
	if !input.ReadASN1(&ski, cbasn1.OCTET_STRING) {
		return nil, model.NewDSSError(fmt.Sprintf("Unable to retrieve subject key identifier of a certificate. "+
			"Reason : %s", "the extension value is not an OCTET STRING"))
	}

	subjectKeyIdentifier := extension.NewSubjectKeyIdentifier()
	subjectKeyIdentifier.SetOctets(extensionValue)
	subjectKeyIdentifier.SetSki([]byte(ski))
	subjectKeyIdentifier.CheckCritical(certificateToken)
	return subjectKeyIdentifier, nil
}

// CertificateExtensionsUtilsCRLDistributionPoints returns the CRL distribution points, when
// present. Port of getCRLDistributionPoints(CertificateToken).
func CertificateExtensionsUtilsCRLDistributionPoints(certificateToken *model.CertificateToken) *extension.CRLDistributionPoints {
	oid := enumerations.CertificateExtensionEnumCRLDistributionPoints.OID()
	crlDistributionPointsBytes := certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid)
	if crlDistributionPointsBytes == nil {
		return nil
	}
	crlDistributionPoints := extension.NewCRLDistributionPoints()
	crlDistributionPoints.SetOctets(crlDistributionPointsBytes)
	crlDistributionPoints.SetCrlUrls(certificateExtensionsUtilsCRLDistributionPointsUrls(crlDistributionPointsBytes))
	crlDistributionPoints.CheckCritical(certificateToken)
	return crlDistributionPoints
}

// certificateExtensionsUtilsCRLDistributionPointsUrls ports the private
// getCRLDistributionPointsUrls(byte[]). Java discards the partially collected URLs when the
// parsing raises, so a failure yields an empty slice.
func certificateExtensionsUtilsCRLDistributionPointsUrls(crlDistributionPointsBytes []byte) []string {
	distributionPoints, ok := certificateExtensionsUtilsSequenceFromDEROctetString(crlDistributionPointsBytes)
	if !ok {
		return []string{}
	}
	urls := make([]string, 0)
	for _, distributionPoint := range distributionPoints {
		// DistributionPoint ::= SEQUENCE { [0] EXPLICIT distributionPoint DistributionPointName OPTIONAL, ... }
		fields, parsed := certificateExtensionsUtilsSequenceElements(distributionPoint)
		if !parsed {
			return []string{}
		}
		var distributionPointName *certificateExtensionsUtilsTagged
		for _, field := range fields {
			tagged, ok := certificateExtensionsUtilsParseTagged(field)
			if !ok {
				return []string{}
			}
			if tagged.tagNo == 0 {
				name, ok := certificateExtensionsUtilsParseTagged(tagged.content)
				if !ok {
					return []string{}
				}
				distributionPointName = &name
			}
		}
		if distributionPointName == nil {
			// Java dereferences a null DistributionPointName and aborts the whole extraction.
			return []string{}
		}
		if distributionPointName.tagNo != 0 { // DistributionPointName.FULL_NAME
			continue
		}
		names := cryptobyte.String(distributionPointName.content)
		for !names.Empty() {
			var element cryptobyte.String
			var tag cbasn1.Tag
			if !names.ReadAnyASN1Element(&element, &tag) {
				return []string{}
			}
			generalName, ok := certificateExtensionsUtilsParseTagged(element)
			if !ok {
				return []string{}
			}
			if location, isURI := certificateExtensionsUtilsParseGn(generalName); isURI {
				urls = append(urls, location)
			}
		}
	}
	return urls
}

// CertificateExtensionsUtilsCRLAccessUrls returns the CRL distribution URIs extracted from the
// cRLDistributionPoints field, or an empty slice. Port of getCRLAccessUrls(CertificateToken).
func CertificateExtensionsUtilsCRLAccessUrls(certificate *model.CertificateToken) []string {
	crlDistributionPoints := CertificateExtensionsUtilsCRLDistributionPoints(certificate)
	if crlDistributionPoints != nil {
		return crlDistributionPoints.CrlUrls()
	}
	return []string{}
}

// CertificateExtensionsUtilsBasicConstraints returns the basic constraints extension. Like
// upstream it always returns an object, even when the certificate carries no such extension.
// Port of getBasicConstraints(CertificateToken).
func CertificateExtensionsUtilsBasicConstraints(certificateToken *model.CertificateToken) *extension.BasicConstraints {
	basicConstraints := extension.NewBasicConstraints()
	basicConstraints.SetOctets(certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), basicConstraints.OID()))

	// CertificateToken#getPathLenConstraint() is the port of X509Certificate#getBasicConstraints().
	value := certificateToken.PathLenConstraint()
	basicConstraints.SetCa(value != -1)
	basicConstraints.SetPathLenConstraint(value)
	basicConstraints.CheckCritical(certificateToken)
	return basicConstraints
}

// CertificateExtensionsUtilsNameConstraints returns the name constraints extension, when present.
// Port of getNameConstraints(CertificateToken).
func CertificateExtensionsUtilsNameConstraints(certificateToken *model.CertificateToken) *extension.NameConstraints {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumNameConstraints.OID()
	nameConstraintsBinaries := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(nameConstraintsBinaries) == 0 {
		return nil
	}

	// NameConstraints ::= SEQUENCE { [0] permittedSubtrees GeneralSubtrees OPTIONAL,
	//     [1] excludedSubtrees GeneralSubtrees OPTIONAL }
	elements, ok := certificateExtensionsUtilsSequenceFromDEROctetString(nameConstraintsBinaries)
	if !ok {
		return nil
	}
	permittedSubtrees := make([]*extension.GeneralSubtree, 0)
	excludedSubtrees := make([]*extension.GeneralSubtree, 0)
	for _, element := range elements {
		tagged, parsed := certificateExtensionsUtilsParseTagged(element)
		if !parsed {
			return nil
		}
		subtrees, converted := certificateExtensionsUtilsGeneralSubtrees(tagged.content)
		if !converted {
			return nil
		}
		switch tagged.tagNo {
		case 0:
			permittedSubtrees = subtrees
		case 1:
			excludedSubtrees = subtrees
		}
	}

	nameConstraints := extension.NewNameConstraints()
	nameConstraints.SetOctets(nameConstraintsBinaries)
	nameConstraints.SetPermittedSubtrees(permittedSubtrees)
	nameConstraints.SetExcludedSubtrees(excludedSubtrees)
	nameConstraints.CheckCritical(certificateToken)
	return nameConstraints
}

// certificateExtensionsUtilsGeneralSubtrees ports the private
// getGeneralSubtrees(org.bouncycastle.asn1.x509.GeneralSubtree[]) together with BouncyCastle's
// GeneralSubtree decoding. General names of an unsupported type are skipped.
func certificateExtensionsUtilsGeneralSubtrees(subtreesContent []byte) ([]*extension.GeneralSubtree, bool) {
	result := make([]*extension.GeneralSubtree, 0)
	input := cryptobyte.String(subtreesContent)
	for !input.Empty() {
		var subtree cryptobyte.String
		if !input.ReadASN1(&subtree, cbasn1.SEQUENCE) {
			return nil, false
		}
		// GeneralSubtree ::= SEQUENCE { base GeneralName,
		//     [0] minimum BaseDistance DEFAULT 0, [1] maximum BaseDistance OPTIONAL }
		var base cryptobyte.String
		var baseTag cbasn1.Tag
		if !subtree.ReadAnyASN1Element(&base, &baseTag) {
			return nil, false
		}
		generalName, ok := certificateExtensionsUtilsParseTagged(base)
		if !ok || generalName.tagNo > 8 {
			// BouncyCastle's GeneralName raises IllegalArgumentException("unknown tag") here, which
			// upstream logs before answering a null NameConstraints.
			return nil, false
		}
		minimum := big.NewInt(0)
		var maximum *big.Int
		count := 0
		for !subtree.Empty() {
			var element cryptobyte.String
			var tag cbasn1.Tag
			if !subtree.ReadAnyASN1Element(&element, &tag) {
				return nil, false
			}
			tagged, parsed := certificateExtensionsUtilsParseTagged(element)
			if !parsed {
				return nil, false
			}
			count++
			switch {
			case tagged.tagNo == 0 && count == 1:
				minimum = certificateExtensionsUtilsInteger(tagged.content)
			case tagged.tagNo == 1 && count <= 2:
				maximum = certificateExtensionsUtilsInteger(tagged.content)
			default:
				return nil, false // BouncyCastle: IllegalArgumentException("Bad tag number")
			}
		}

		generalNameType := enumerations.GeneralNameTypeFromIndex(generalName.tagNo)
		if generalNameType == "" {
			// Unreachable: BouncyCastle already rejected every tag GeneralNameType does not cover.
			// Ported all the same, as upstream logs and skips the subtree here.
			continue
		}
		generalSubtree := extension.NewGeneralSubtree()
		generalSubtree.SetGeneralNameType(generalNameType)
		generalSubtree.SetMinimum(minimum)
		generalSubtree.SetMaximum(maximum)
		generalSubtree.SetValue(certificateExtensionsUtilsStringValue(generalNameType, generalName))
		result = append(result, generalSubtree)
	}
	return result, true
}

// certificateExtensionsUtilsStringValue ports the private
// getStringValue(GeneralNameType, ASN1Encodable). The ASN1Encodable is BouncyCastle's decoding of
// the GeneralName body: a SEQUENCE for otherName/x400Address/ediPartyName, an IA5String for
// rfc822Name/dNSName/uniformResourceIdentifier, an X500Name for directoryName, an OCTET STRING for
// iPAddress and an OBJECT IDENTIFIER for registeredID.
func certificateExtensionsUtilsStringValue(generalNameType enumerations.GeneralNameType, generalName certificateExtensionsUtilsTagged) string {
	switch generalNameType {
	case enumerations.GeneralNameTypeOtherName,
		enumerations.GeneralNameTypeEDIPartyName,
		enumerations.GeneralNameTypeX400Address:
		return certificateExtensionsUtilsToHexEncoded(certificateExtensionsUtilsRetag(generalName.content, cbasn1.SEQUENCE))

	case enumerations.GeneralNameTypeRFC822Name,
		enumerations.GeneralNameTypeDNSName,
		enumerations.GeneralNameTypeUniformResourceIdentifier:
		if !generalName.constructed {
			return string(generalName.content)
		}
		return certificateExtensionsUtilsToHexEncoded(generalName.full)

	case enumerations.GeneralNameTypeDirectoryName:
		principal, err := model.NewX500Principal(generalName.content)
		if err != nil {
			return certificateExtensionsUtilsToHexEncoded(generalName.content)
		}
		return model.NewX500PrincipalHelper(principal).RFC2253()

	case enumerations.GeneralNameTypeIPAddress:
		return certificateExtensionsUtilsToHexEncoded(generalName.content)

	case enumerations.GeneralNameTypeRegisteredID:
		objectIdentifier, ok := certificateExtensionsUtilsObjectIdentifier(generalName.content)
		if !ok {
			return certificateExtensionsUtilsToHexEncoded(generalName.content)
		}
		return objectIdentifier

	default:
		return certificateExtensionsUtilsToHexEncoded(generalName.full)
	}
}

// certificateExtensionsUtilsToHexEncoded ports the private toHexEncoded(byte[]).
func certificateExtensionsUtilsToHexEncoded(binaries []byte) string {
	return "#" + utils.ToHex(binaries)
}

// CertificateExtensionsUtilsPolicyConstraints returns the policy constraints extension, when
// present. Port of getPolicyConstraints(CertificateToken).
func CertificateExtensionsUtilsPolicyConstraints(certificateToken *model.CertificateToken) *extension.PolicyConstraints {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumPolicyConstraints.OID()
	policyConstraintsBinaries := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(policyConstraintsBinaries) == 0 {
		return nil
	}

	// PolicyConstraints ::= SEQUENCE { [0] requireExplicitPolicy SkipCerts OPTIONAL,
	//     [1] inhibitPolicyMapping SkipCerts OPTIONAL }
	elements, ok := certificateExtensionsUtilsSequenceFromDEROctetString(policyConstraintsBinaries)
	if !ok {
		return nil
	}
	var requireExplicitPolicy, inhibitPolicyMapping *big.Int
	for _, element := range elements {
		tagged, parsed := certificateExtensionsUtilsParseTagged(element)
		if !parsed {
			return nil
		}
		switch tagged.tagNo {
		case 0:
			requireExplicitPolicy = certificateExtensionsUtilsInteger(tagged.content)
		case 1:
			inhibitPolicyMapping = certificateExtensionsUtilsInteger(tagged.content)
		default:
			return nil // BouncyCastle: IllegalArgumentException("Unknown tag encountered.")
		}
	}

	policyConstraints := extension.NewPolicyConstraints()
	policyConstraints.SetOctets(policyConstraintsBinaries)
	if inhibitPolicyMapping != nil {
		policyConstraints.SetInhibitPolicyMapping(int(inhibitPolicyMapping.Int64()))
	}
	if requireExplicitPolicy != nil {
		policyConstraints.SetRequireExplicitPolicy(int(requireExplicitPolicy.Int64()))
	}
	policyConstraints.CheckCritical(certificateToken)
	return policyConstraints
}

// CertificateExtensionsUtilsInhibitAnyPolicy returns the inhibit anyPolicy extension, when
// present. Port of getInhibitAnyPolicy(CertificateToken).
func CertificateExtensionsUtilsInhibitAnyPolicy(certificateToken *model.CertificateToken) *extension.InhibitAnyPolicy {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumInhibitAnyPolicy.OID()
	inhibitAnyPolicyBinaries := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(inhibitAnyPolicyBinaries) == 0 {
		return nil
	}

	content := certificateExtensionsUtilsExtensionContent(certificate, oid)
	input := cryptobyte.String(content)
	var value cryptobyte.String
	if !input.ReadASN1(&value, cbasn1.INTEGER) {
		return nil
	}

	inhibitAnyPolicy := extension.NewInhibitAnyPolicy()
	inhibitAnyPolicy.SetOctets(inhibitAnyPolicyBinaries)
	inhibitAnyPolicy.SetValue(int(certificateExtensionsUtilsInteger(value).Int64()))
	inhibitAnyPolicy.CheckCritical(certificateToken)
	return inhibitAnyPolicy
}

// CertificateExtensionsUtilsFreshestCRL returns the Freshest CRL, when present.
// Port of getFreshestCRL(CertificateToken).
func CertificateExtensionsUtilsFreshestCRL(certificateToken *model.CertificateToken) *extension.FreshestCRL {
	oid := enumerations.CertificateExtensionEnumFreshestCRL.OID()
	freshestCrlBytes := certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid)
	if freshestCrlBytes == nil {
		return nil
	}
	freshestCRL := extension.NewFreshestCRL()
	freshestCRL.SetOctets(freshestCrlBytes)
	freshestCRL.SetCrlUrls(certificateExtensionsUtilsCRLDistributionPointsUrls(freshestCrlBytes))
	freshestCRL.CheckCritical(certificateToken)
	return freshestCRL
}

// CertificateExtensionsUtilsKeyUsage returns the key usage, when present.
// Port of getKeyUsage(CertificateToken).
func CertificateExtensionsUtilsKeyUsage(certificateToken *model.CertificateToken) *extension.KeyUsage {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumKeyUsage.OID()
	// X509Certificate#getKeyUsage() answers null only when the extension is absent.
	if certificateExtensionsUtilsExtensionContent(certificate, oid) == nil {
		return nil
	}
	keyUsage := extension.NewKeyUsage()
	keyUsage.SetOctets(certificateExtensionsUtilsExtensionValue(certificate, oid))
	keyUsage.SetKeyUsageBits(certificateToken.KeyUsageBits())
	keyUsage.CheckCritical(certificateToken)
	return keyUsage
}

// CertificateExtensionsUtilsExtendedKeyUsage returns the extended key usage, when present.
// Port of getExtendedKeyUsage(CertificateToken).
func CertificateExtensionsUtilsExtendedKeyUsage(certificateToken *model.CertificateToken) *extension.ExtendedKeyUsages {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumExtendedKeyUsage.OID()
	extendedKeyUsage := extension.NewExtendedKeyUsages()
	extendedKeyUsage.SetOctets(certificateExtensionsUtilsExtensionValue(certificate, oid))

	// ExtKeyUsageSyntax ::= SEQUENCE SIZE (1..MAX) OF KeyPurposeId. crypto/x509 splits the purpose
	// identifiers over ExtKeyUsage and UnknownExtKeyUsage and loses their order, so the extension is
	// decoded here instead.
	if content := certificateExtensionsUtilsExtensionContent(certificate, oid); content != nil {
		elements, parsed := certificateExtensionsUtilsSequenceElements(content)
		if !parsed {
			return nil // CertificateParsingException
		}
		oids := make([]string, 0, len(elements))
		for _, element := range elements {
			keyPurposeID, ok := certificateExtensionsUtilsObjectIdentifierElement(element)
			if !ok {
				return nil
			}
			oids = append(oids, keyPurposeID)
		}
		extendedKeyUsage.SetOids(oids)
	}
	extendedKeyUsage.CheckCritical(certificateToken)
	return extendedKeyUsage
}

// CertificateExtensionsUtilsCertificatePolicies returns the certificate policies, when present.
// Port of getCertificatePolicies(CertificateToken).
func CertificateExtensionsUtilsCertificatePolicies(certificateToken *model.CertificateToken) *extension.CertificatePolicies {
	certificate := certificateToken.Certificate()
	oid := enumerations.CertificateExtensionEnumCertificatePolicies.OID()
	certificatePoliciesBinaries := certificateExtensionsUtilsExtensionValue(certificate, oid)
	if len(certificatePoliciesBinaries) == 0 {
		return nil
	}

	elements, ok := certificateExtensionsUtilsSequenceFromDEROctetString(certificatePoliciesBinaries)
	if !ok {
		return nil
	}
	policiesList := make([]*extension.CertificatePolicy, 0, len(elements))
	for _, element := range elements {
		certificatePolicy, converted := certificateExtensionsUtilsCertificatePolicy(element)
		if !converted {
			return nil
		}
		policiesList = append(policiesList, certificatePolicy)
	}

	certificatePolicies := extension.NewCertificatePolicies()
	certificatePolicies.SetOctets(certificatePoliciesBinaries)
	certificatePolicies.SetPolicyList(policiesList)
	certificatePolicies.CheckCritical(certificateToken)
	return certificatePolicies
}

// certificateExtensionsUtilsCertificatePolicy ports the private getCertificatePolicy(ASN1Encodable)
// together with BouncyCastle's PolicyInformation/PolicyQualifierInfo decoding.
func certificateExtensionsUtilsCertificatePolicy(policyObject []byte) (*extension.CertificatePolicy, bool) {
	// PolicyInformation ::= SEQUENCE { policyIdentifier CertPolicyId,
	//     policyQualifiers SEQUENCE SIZE (1..MAX) OF PolicyQualifierInfo OPTIONAL }
	fields, parsed := certificateExtensionsUtilsSequenceElements(policyObject)
	if !parsed || len(fields) < 1 || len(fields) > 2 {
		return nil, false
	}
	policyIdentifier, ok := certificateExtensionsUtilsObjectIdentifierElement(fields[0])
	if !ok {
		return nil, false
	}

	certificatePolicy := extension.NewCertificatePolicy()
	certificatePolicy.SetOid(policyIdentifier)
	if len(fields) == 2 {
		policyQualifiers, parsed := certificateExtensionsUtilsSequenceElements(fields[1])
		if !parsed {
			return nil, false
		}
		for _, policyQualifier := range policyQualifiers {
			// PolicyQualifierInfo ::= SEQUENCE { policyQualifierId PolicyQualifierId, qualifier ANY }
			qualifierFields, parsed := certificateExtensionsUtilsSequenceElements(policyQualifier)
			if !parsed || len(qualifierFields) != 2 {
				return nil, false
			}
			policyQualifierID, ok := certificateExtensionsUtilsObjectIdentifierElement(qualifierFields[0])
			if !ok {
				return nil, false
			}
			if policyQualifierID == certificateExtensionsUtilsOIDIdQtCps {
				certificatePolicy.SetCpsUrl(DSSASN1UtilsString(qualifierFields[1]))
			}
		}
	}
	return certificatePolicy, true
}

// CertificateExtensionsUtilsOcspNoCheck returns the ocsp-nocheck extension value, when present.
// Port of getOcspNoCheck(CertificateToken).
func CertificateExtensionsUtilsOcspNoCheck(certificateToken *model.CertificateToken) *extension.OCSPNoCheck {
	oid := enumerations.CertificateExtensionEnumOCSPNoCheck.OID()
	extensionValue := certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid)
	if extensionValue == nil {
		return nil
	}
	ocspNoCheck := extension.NewOCSPNoCheck()
	ocspNoCheck.SetOctets(extensionValue)
	ocspNoCheck.SetOcspNoCheck(certificateExtensionsUtilsIsNullIdentifiedValuePresent(extensionValue))
	ocspNoCheck.CheckCritical(certificateToken)
	return ocspNoCheck
}

// CertificateExtensionsUtilsHasOcspNoCheckExtension reports whether the certificate carries the
// ocsp-nocheck extension, indicating whether revocation data should be checked for an OCSP signing
// certificate (RFC 6960). Port of hasOcspNoCheckExtension(CertificateToken).
func CertificateExtensionsUtilsHasOcspNoCheckExtension(certificateToken *model.CertificateToken) bool {
	ocspNoCheck := CertificateExtensionsUtilsOcspNoCheck(certificateToken)
	return ocspNoCheck != nil && ocspNoCheck.IsOcspNoCheck()
}

// CertificateExtensionsUtilsValAssuredSTCerts returns the ext-etsi-valassured-ST-certs extension
// value, when present. Port of getValAssuredSTCerts(CertificateToken).
func CertificateExtensionsUtilsValAssuredSTCerts(certificateToken *model.CertificateToken) *extension.ValidityAssuredShortTerm {
	oid := enumerations.CertificateExtensionEnumValidityAssuredShortTerm.OID()
	extensionValue := certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid)
	if extensionValue == nil {
		return nil
	}
	validityAssuredShortTerm := extension.NewValidityAssuredShortTerm()
	validityAssuredShortTerm.SetOctets(extensionValue)
	validityAssuredShortTerm.SetValAssuredSTCerts(certificateExtensionsUtilsIsNullIdentifiedValuePresent(extensionValue))
	validityAssuredShortTerm.CheckCritical(certificateToken)
	return validityAssuredShortTerm
}

// certificateExtensionsUtilsIsNullIdentifiedValuePresent ports the private
// isNullIdentifiedValuePresent(byte[]): the extension value is a DER OCTET STRING wrapping a DER NULL.
func certificateExtensionsUtilsIsNullIdentifiedValuePresent(extensionValue []byte) bool {
	input := cryptobyte.String(extensionValue)
	var content cryptobyte.String
	if !input.ReadASN1(&content, cbasn1.OCTET_STRING) {
		return false
	}
	var null cryptobyte.String
	if !content.ReadASN1(&null, cbasn1.NULL) {
		return false
	}
	return len(null) == 0
}

// CertificateExtensionsUtilsNoRevAvail returns the noRevAvail extension value, when present.
// Port of getNoRevAvail(CertificateToken).
func CertificateExtensionsUtilsNoRevAvail(certificateToken *model.CertificateToken) *extension.NoRevAvail {
	oid := enumerations.CertificateExtensionEnumNoRevocationAvailable.OID()
	extensionValue := certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid)
	if extensionValue == nil {
		return nil
	}
	noRevAvail := extension.NewNoRevAvail()
	noRevAvail.SetOctets(extensionValue)
	noRevAvail.SetNoRevAvail(certificateExtensionsUtilsIsNullIdentifiedValuePresent(extensionValue))
	noRevAvail.CheckCritical(certificateToken)
	return noRevAvail
}

// CertificateExtensionsUtilsHasValAssuredShortTermCertsExtension reports whether the certificate
// carries the ext-etsi-valassured-ST-certs extension, i.e. whether its validity is assured because
// it is a "short-term certificate". Port of hasValAssuredShortTermCertsExtension(CertificateToken).
func CertificateExtensionsUtilsHasValAssuredShortTermCertsExtension(certificateToken *model.CertificateToken) bool {
	valAssuredSTCerts := CertificateExtensionsUtilsValAssuredSTCerts(certificateToken)
	return valAssuredSTCerts != nil && valAssuredSTCerts.IsValAssuredSTCerts()
}

// CertificateExtensionsUtilsQcStatements returns the qc-statements extension value, when present.
// Port of getQcStatements(CertificateToken).
func CertificateExtensionsUtilsQcStatements(certificateToken *model.CertificateToken) *extension.QcStatements {
	qcStatements := QcStatementUtilsQcStatements(certificateToken)
	if qcStatements != nil {
		qcStatements.CheckCritical(certificateToken)
	}
	return qcStatements
}

// certificateExtensionsUtilsOtherCertificateExtension ports the private
// getOtherCertificateExtension(CertificateToken, String) for an unsupported extension OID.
func certificateExtensionsUtilsOtherCertificateExtension(certificateToken *model.CertificateToken, oid string) *extension.CertificateExtension {
	var certificateExtension extension.CertificateExtension
	if value := enumerations.CertificateExtensionEnumForOID(oid); value != "" {
		certificateExtension = extension.NewCertificateExtensionFromEnum(value)
	} else {
		certificateExtension = extension.NewCertificateExtension(oid)
	}
	certificateExtension.SetOctets(certificateExtensionsUtilsExtensionValue(certificateToken.Certificate(), oid))
	certificateExtension.CheckCritical(certificateToken)
	return &certificateExtension
}

// ---------------------------------------------------------------------------------------------
// ASN.1 helpers replacing BouncyCastle. They are private to this file per the porting conventions.
// ---------------------------------------------------------------------------------------------

// certificateExtensionsUtilsExtensionValue reproduces X509Certificate#getExtensionValue(String):
// the DER-encoded OCTET STRING wrapping the extension value, or nil when the extension is absent.
func certificateExtensionsUtilsExtensionValue(certificate *x509.Certificate, oid string) []byte {
	content := certificateExtensionsUtilsExtensionContent(certificate, oid)
	if content == nil {
		return nil
	}
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1OctetString(content)
	encoded, err := builder.Bytes()
	if err != nil {
		return nil
	}
	return encoded
}

// certificateExtensionsUtilsExtensionContent reproduces BouncyCastle's
// X509CertificateImpl#getExtensionBytes(Certificate, String): the extension value octets, unwrapped.
func certificateExtensionsUtilsExtensionContent(certificate *x509.Certificate, oid string) []byte {
	for _, ext := range certificate.Extensions {
		if ext.Id.String() == oid {
			return ext.Value
		}
	}
	return nil
}

// certificateExtensionsUtilsSequenceFromDEROctetString ports
// DSSASN1Utils#getAsn1SequenceFromDerOctetString(byte[]), returning the elements of the SEQUENCE
// encapsulated in the DER OCTET STRING.
func certificateExtensionsUtilsSequenceFromDEROctetString(binaries []byte) ([][]byte, bool) {
	input := cryptobyte.String(binaries)
	var content cryptobyte.String
	if !input.ReadASN1(&content, cbasn1.OCTET_STRING) {
		return nil, false
	}
	return certificateExtensionsUtilsSequenceElements(content)
}

// certificateExtensionsUtilsSequenceElements returns the complete DER encoding of every element of
// the SEQUENCE starting at der. Trailing bytes are ignored, as ASN1InputStream#readObject() does.
func certificateExtensionsUtilsSequenceElements(der []byte) ([][]byte, bool) {
	input := cryptobyte.String(der)
	var sequence cryptobyte.String
	if !input.ReadASN1(&sequence, cbasn1.SEQUENCE) {
		return nil, false
	}
	elements := make([][]byte, 0)
	for !sequence.Empty() {
		var element cryptobyte.String
		var tag cbasn1.Tag
		if !sequence.ReadAnyASN1Element(&element, &tag) {
			return nil, false
		}
		elements = append(elements, element)
	}
	return elements, true
}

// certificateExtensionsUtilsTagged is a context-tagged ASN.1 element, i.e. a GeneralName or any
// other CHOICE/OPTIONAL member of the extensions parsed here.
type certificateExtensionsUtilsTagged struct {
	// tagNo is the context tag number.
	tagNo int
	// full is the complete DER encoding, tag and length included.
	full []byte
	// content holds the content octets.
	content []byte
	// constructed reports whether the element is constructed.
	constructed bool
}

// certificateExtensionsUtilsParseTagged decodes a context-tagged element.
func certificateExtensionsUtilsParseTagged(der []byte) (certificateExtensionsUtilsTagged, bool) {
	input := cryptobyte.String(der)
	var element cryptobyte.String
	var tag cbasn1.Tag
	if !input.ReadAnyASN1Element(&element, &tag) {
		return certificateExtensionsUtilsTagged{}, false
	}
	// 0x80 is cryptobyte/asn1's (unexported) classContextSpecific, 0x20 its classConstructed.
	if uint8(tag)&0xc0 != 0x80 {
		return certificateExtensionsUtilsTagged{}, false
	}
	body := cryptobyte.String(element)
	var content cryptobyte.String
	var contentTag cbasn1.Tag
	if !body.ReadAnyASN1(&content, &contentTag) {
		return certificateExtensionsUtilsTagged{}, false
	}
	return certificateExtensionsUtilsTagged{
		tagNo:       int(uint8(tag) & 0x1f),
		full:        element,
		content:     content,
		constructed: uint8(tag)&0x20 != 0,
	}, true
}

// certificateExtensionsUtilsParseGn ports the private parseGn(GeneralName): only a
// uniformResourceIdentifier yields a location.
func certificateExtensionsUtilsParseGn(generalName certificateExtensionsUtilsTagged) (string, bool) {
	if generalName.tagNo == 6 && !generalName.constructed {
		return string(generalName.content), true
	}
	return "", false
}

// certificateExtensionsUtilsRetag re-encodes content under the given universal tag, standing in for
// BouncyCastle's implicit re-tagging of a context-tagged element (ASN1Sequence.getInstance(o, false)).
func certificateExtensionsUtilsRetag(content []byte, tag cbasn1.Tag) []byte {
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1(tag, func(child *cryptobyte.Builder) {
		child.AddBytes(content)
	})
	encoded, err := builder.Bytes()
	if err != nil {
		return content
	}
	return encoded
}

// certificateExtensionsUtilsInteger decodes the two's-complement content octets of an INTEGER.
func certificateExtensionsUtilsInteger(content []byte) *big.Int {
	value := new(big.Int).SetBytes(content)
	if len(content) > 0 && content[0]&0x80 != 0 {
		value.Sub(value, new(big.Int).Lsh(big.NewInt(1), uint(len(content))*8))
	}
	return value
}

// certificateExtensionsUtilsObjectIdentifier decodes the content octets of an OBJECT IDENTIFIER.
func certificateExtensionsUtilsObjectIdentifier(content []byte) (string, bool) {
	return certificateExtensionsUtilsObjectIdentifierElement(certificateExtensionsUtilsRetag(content, cbasn1.OBJECT_IDENTIFIER))
}

// certificateExtensionsUtilsObjectIdentifierElement decodes a complete OBJECT IDENTIFIER element.
func certificateExtensionsUtilsObjectIdentifierElement(der []byte) (string, bool) {
	input := cryptobyte.String(der)
	var objectIdentifier encasn1.ObjectIdentifier
	if !input.ReadASN1ObjectIdentifier(&objectIdentifier) {
		return "", false
	}
	return objectIdentifier.String(), true
}

// certificateExtensionsUtilsHostAddress ports java.net.InetAddress#getByAddress(byte[]) followed by
// #getHostAddress(), as BouncyCastle renders an iPAddress general name. An address of any other
// length raises an UnknownHostException there, which is reported as ok=false.
func certificateExtensionsUtilsHostAddress(address []byte) (string, bool) {
	switch len(address) {
	case 4:
		return fmt.Sprintf("%d.%d.%d.%d", address[0], address[1], address[2], address[3]), true
	case 16:
		// getByAddress() answers an Inet4Address for an IPv4-mapped address.
		if mapped := certificateExtensionsUtilsIPv4Mapped(address); mapped != nil {
			return fmt.Sprintf("%d.%d.%d.%d", mapped[0], mapped[1], mapped[2], mapped[3]), true
		}
		// Inet6Address#numericToTextFormat: eight lowercase hexadecimal groups, never compressed.
		groups := make([]string, 0, 8)
		for i := 0; i < 8; i++ {
			groups = append(groups, fmt.Sprintf("%x", int(address[i*2])<<8|int(address[i*2+1])))
		}
		return strings.Join(groups, ":"), true
	default:
		return "", false
	}
}

// certificateExtensionsUtilsIPv4Mapped ports sun.net.util.IPAddressUtil#convertFromIPv4MappedAddress.
func certificateExtensionsUtilsIPv4Mapped(address []byte) []byte {
	for i := 0; i < 10; i++ {
		if address[i] != 0 {
			return nil
		}
	}
	if address[10] != 0xff || address[11] != 0xff {
		return nil
	}
	return address[12:16]
}

// certificateExtensionsUtilsRFC4519Name ports BouncyCastle's
// X500Name.getInstance(RFC4519Style.INSTANCE, name).toString(), whose result upstream immediately
// feeds back to RFC4519Style#fromString in toRFC2253RDN. Rather than materialising the intermediate
// string, the round trip is applied directly to the DER: RFC4519Style re-encodes every attribute
// value it rendered as a string (IA5String for domainComponent, PrintableString for country,
// serialNumber, dnQualifier and telephoneNumber, UTF8String otherwise), and returns the original
// encoding verbatim for every value it rendered hex-encoded.
//
// The returned value is the retyped Name in DER; certificateExtensionsUtilsToRFC2253RDN turns it
// into the RFC 2253 string.
func certificateExtensionsUtilsRFC4519Name(nameDER []byte) ([]byte, bool) {
	input := cryptobyte.String(nameDER)
	var sequence cryptobyte.String
	if !input.ReadASN1(&sequence, cbasn1.SEQUENCE) {
		return nil, false
	}
	relativeDistinguishedNames := make([][]byte, 0)
	for !sequence.Empty() {
		var set cryptobyte.String
		if !sequence.ReadASN1(&set, cbasn1.SET) {
			return nil, false
		}
		attributes := make([][]byte, 0)
		for !set.Empty() {
			var attribute cryptobyte.String
			if !set.ReadASN1(&attribute, cbasn1.SEQUENCE) {
				return nil, false
			}
			var attributeType encasn1.ObjectIdentifier
			if !attribute.ReadASN1ObjectIdentifier(&attributeType) {
				return nil, false
			}
			var value cryptobyte.String
			var valueTag cbasn1.Tag
			if !attribute.ReadAnyASN1Element(&value, &valueTag) {
				return nil, false
			}
			encoded, ok := certificateExtensionsUtilsEncodeAttributeValue(attributeType.String(), value)
			if !ok {
				return nil, false
			}
			builder := cryptobyte.NewBuilder(nil)
			builder.AddASN1(cbasn1.SEQUENCE, func(child *cryptobyte.Builder) {
				child.AddASN1ObjectIdentifier(attributeType)
				child.AddBytes(encoded)
			})
			attributeDER, err := builder.Bytes()
			if err != nil {
				return nil, false
			}
			attributes = append(attributes, attributeDER)
		}
		builder := cryptobyte.NewBuilder(nil)
		builder.AddASN1(cbasn1.SET, func(child *cryptobyte.Builder) {
			for _, attribute := range attributes {
				child.AddBytes(attribute)
			}
		})
		setDER, err := builder.Bytes()
		if err != nil {
			return nil, false
		}
		relativeDistinguishedNames = append(relativeDistinguishedNames, setDER)
	}
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1(cbasn1.SEQUENCE, func(child *cryptobyte.Builder) {
		for _, relativeDistinguishedName := range relativeDistinguishedNames {
			child.AddBytes(relativeDistinguishedName)
		}
	})
	retyped, err := builder.Bytes()
	if err != nil {
		return nil, false
	}
	return retyped, true
}

// certificateExtensionsUtilsEncodeAttributeValue ports RFC4519Style#stringToValue applied to the
// string RFC4519Style#toString produced for the same value.
func certificateExtensionsUtilsEncodeAttributeValue(attributeType string, value []byte) ([]byte, bool) {
	text, isString := certificateExtensionsUtilsASN1String(value)
	if !isString {
		return value, true
	}
	tag := cbasn1.UTF8String
	switch attributeType {
	case "0.9.2342.19200300.100.1.25": // domainComponent
		tag = cbasn1.IA5String
	case "2.5.4.6", "2.5.4.5", "2.5.4.46", "2.5.4.20": // c, serialNumber, dnQualifier, telephoneNumber
		tag = cbasn1.PrintableString
	}
	var content []byte
	if tag == cbasn1.UTF8String {
		content = []byte(text)
	} else {
		// org.bouncycastle.util.Strings#toByteArray truncates every UTF-16 code unit to a byte.
		units := utf16.Encode([]rune(text))
		content = make([]byte, len(units))
		for i, unit := range units {
			content[i] = byte(unit)
		}
	}
	builder := cryptobyte.NewBuilder(nil)
	builder.AddASN1(tag, func(child *cryptobyte.Builder) {
		child.AddBytes(content)
	})
	encoded, err := builder.Bytes()
	if err != nil {
		return nil, false
	}
	return encoded, true
}

// certificateExtensionsUtilsToRFC2253RDN ports the private toRFC2253RDN(String): the RFC 4519
// rendering is parsed back into an X500Principal, whose RFC 2253 name is returned. The Java method
// answers its input when the conversion fails; here the input is the retyped DER produced by
// certificateExtensionsUtilsRFC4519Name, so a failure can only come from the principal parsing and
// yields the hex-encoded name.
func certificateExtensionsUtilsToRFC2253RDN(retypedNameDER []byte) string {
	principal, err := model.NewX500Principal(retypedNameDER)
	if err != nil {
		return certificateExtensionsUtilsToHexEncoded(retypedNameDER)
	}
	return model.NewX500PrincipalHelper(principal).RFC2253()
}

// certificateExtensionsUtilsASN1String decodes a complete DER element whose class implements
// BouncyCastle's ASN1String, reporting isString=false otherwise. DERUniversalString implements
// ASN1String but is explicitly excluded by IETFUtils#valueToString, and neither ASN1UTCTime,
// ASN1GeneralizedTime nor ASN1ObjectDescriptor implement it.
func certificateExtensionsUtilsASN1String(der []byte) (string, bool) {
	input := cryptobyte.String(der)
	var content cryptobyte.String
	var tag cbasn1.Tag
	if !input.ReadAnyASN1(&content, &tag) {
		return "", false
	}
	switch tag {
	case cbasn1.UTF8String: // org.bouncycastle.util.Strings#fromUTF8ByteArray
		return string(content), true
	case cbasn1.Tag(30): // BMPString, decoded as UTF-16BE
		if len(content)%2 != 0 {
			return "", false
		}
		units := make([]uint16, len(content)/2)
		for i := range units {
			units[i] = uint16(content[i*2])<<8 | uint16(content[i*2+1])
		}
		return string(utf16.Decode(units)), true
	case cbasn1.PrintableString, cbasn1.IA5String, cbasn1.T61String,
		cbasn1.Tag(18), cbasn1.Tag(21), cbasn1.Tag(25), cbasn1.Tag(26), cbasn1.Tag(27):
		// NumericString, VideotexString, GraphicString, VisibleString and GeneralString join
		// PrintableString, IA5String and T61String in org.bouncycastle.util.Strings#fromByteArray.
		runes := make([]rune, len(content))
		for i, b := range content {
			runes[i] = rune(b)
		}
		return string(runes), true
	default:
		return "", false
	}
}
