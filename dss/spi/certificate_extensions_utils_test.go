package spi

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/utils"
)

// The expectations below are transcribed verbatim from upstream's
// dss-spi/src/test/java/eu/europa/esig/dss/spi/CertificateExtensionUtilsTest.java and
// QcStatementsUtilsTest.java, so that the fidelity of the port stays legible without decoding
// testdata/certificate_extensions/kat.tsv. Each certificate is named by its fixture file; the
// files are stable, and the mechanical known-answer tests cover the whole corpus anyway.
const (
	// certificateExtensionsTestTSPCertificate2014 is src/test/resources/TSP_Certificate_2014.crt,
	// upstream's "certificateWithAIA".
	certificateExtensionsTestTSPCertificate2014 = "cert_32.der"
	// certificateExtensionsTestECEuropaEu is src/test/resources/ec.europa.eu.crt.
	certificateExtensionsTestECEuropaEu = "cert_40.der"
	// certificateExtensionsTestUPCDirectoryName carries a directoryName subject alternative name.
	certificateExtensionsTestUPCDirectoryName = "cert_22.der"
	// certificateExtensionsTestFNMTDirectoryName carries a directoryName with unregistered OIDs.
	certificateExtensionsTestFNMTDirectoryName = "cert_14.der"
	// certificateExtensionsTestCertSignOtherName carries an otherName subject alternative name.
	certificateExtensionsTestCertSignOtherName = "cert_44.der"
	// certificateExtensionsTestExcludedIPAddresses carries excluded iPAddress name constraints.
	certificateExtensionsTestExcludedIPAddresses = "cert_04.der"
	// certificateExtensionsTestPolicyConstraints carries a requireExplicitPolicy constraint.
	certificateExtensionsTestPolicyConstraints = "cert_09.der"
)

func TestCertificateExtensionsUtilsUpstreamSubjectKeyIdentifier(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestTSPCertificate2014)
	subjectKeyIdentifier, err := CertificateExtensionsUtilsSubjectKeyIdentifier(certificateToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subjectKeyIdentifier == nil {
		t.Fatalf("SubjectKeyIdentifier: got nil, want a value")
	}
	if got, want := utils.ToHex(subjectKeyIdentifier.Ski()), "4c4c4cfcacace6bb"; got != want {
		t.Errorf("SKI = %s, want %s", got, want)
	}
}

func TestCertificateExtensionsUtilsUpstreamCertificatePolicies(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestTSPCertificate2014)
	certificatePolicies := CertificateExtensionsUtilsCertificatePolicies(certificateToken)
	if certificatePolicies == nil {
		t.Fatalf("CertificatePolicies: got nil, want a value")
	}
	policyList := certificatePolicies.PolicyList()
	if len(policyList) != 2 {
		t.Fatalf("CertificatePolicies has %d entries, want 2", len(policyList))
	}
	if got, want := policyList[0].Oid(), "1.3.171.1.1.10.8.1"; got != want {
		t.Errorf("policy 0 OID = %q, want %q", got, want)
	}
	if got, want := policyList[0].CpsUrl(), "https://repository.luxtrust.lu"; got != want {
		t.Errorf("policy 0 CPS URL = %q, want %q", got, want)
	}
	if got, want := policyList[1].Oid(), "0.4.0.2042.1.3"; got != want {
		t.Errorf("policy 1 OID = %q, want %q", got, want)
	}
	if got := policyList[1].CpsUrl(); got != "" {
		t.Errorf("policy 1 CPS URL = %q, want the empty string (Java null)", got)
	}
}

func TestCertificateExtensionsUtilsUpstreamAuthorityKeyIdentifier(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestTSPCertificate2014)
	authorityKeyIdentifier, err := CertificateExtensionsUtilsAuthorityKeyIdentifier(certificateToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if authorityKeyIdentifier == nil {
		t.Fatalf("AuthorityKeyIdentifier: got nil, want a value")
	}
	if utils.IsArrayEmpty(authorityKeyIdentifier.KeyIdentifier()) {
		t.Errorf("KeyIdentifier is empty, want a value")
	}
	if utils.IsArrayNotEmpty(authorityKeyIdentifier.AuthorityCertIssuerSerial()) {
		t.Errorf("AuthorityCertIssuerSerial = %v, want nil", authorityKeyIdentifier.AuthorityCertIssuerSerial())
	}
}

func TestCertificateExtensionsUtilsUpstreamAccessLocations(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestECEuropaEu)
	aia := CertificateExtensionsUtilsAuthorityInformationAccess(certificateToken)
	if aia == nil {
		t.Fatalf("AuthorityInformationAccess: got nil, want a value")
	}
	if got, want := aia.Ocsp(), []string{"http://ocsp.luxtrust.lu"}; !certificateExtensionsTestEqual(got, want) {
		t.Errorf("Ocsp() = %v, want %v", got, want)
	}
	if got := CertificateExtensionsUtilsOCSPAccessUrls(certificateToken); !certificateExtensionsTestEqual(got, aia.Ocsp()) {
		t.Errorf("OCSPAccessUrls() = %v, want %v", got, aia.Ocsp())
	}
	if got, want := aia.CaIssuers(), []string{"http://ca.luxtrust.lu/LTQCA.crt"}; !certificateExtensionsTestEqual(got, want) {
		t.Errorf("CaIssuers() = %v, want %v", got, want)
	}
	if got := CertificateExtensionsUtilsCAIssuersAccessUrls(certificateToken); !certificateExtensionsTestEqual(got, aia.CaIssuers()) {
		t.Errorf("CAIssuersAccessUrls() = %v, want %v", got, aia.CaIssuers())
	}

	crlDistributionPoints := CertificateExtensionsUtilsCRLDistributionPoints(certificateToken)
	if crlDistributionPoints == nil {
		t.Fatalf("CRLDistributionPoints: got nil, want a value")
	}
	if got, want := crlDistributionPoints.CrlUrls(), []string{"http://crl.luxtrust.lu/LTQCA.crl"}; !certificateExtensionsTestEqual(got, want) {
		t.Errorf("CrlUrls() = %v, want %v", got, want)
	}
	if got := CertificateExtensionsUtilsCRLAccessUrls(certificateToken); !certificateExtensionsTestEqual(got, crlDistributionPoints.CrlUrls()) {
		t.Errorf("CRLAccessUrls() = %v, want %v", got, crlDistributionPoints.CrlUrls())
	}

	// The extension is present but holds no name at all.
	subjectAlternativeNames := CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken)
	if subjectAlternativeNames == nil {
		t.Fatalf("SubjectAlternativeNames: got nil, want a value")
	}
	if got := len(subjectAlternativeNames.GeneralNames()); got != 0 {
		t.Errorf("GeneralNames() has %d entries, want 0", got)
	}
}

func TestCertificateExtensionsUtilsUpstreamSubjectAlternativeNames(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		file     string
		expected [][2]string
	}{
		{
			name: "directoryName",
			file: certificateExtensionsTestUPCDirectoryName,
			expected: [][2]string{
				{"DIRECTORY_NAME", "CN=4.2.1.10_04_EE2,O=UPC,L=Barcelona,C=ES"},
			},
		},
		{
			// The unregistered attribute types have no RFC 2253 keyword, so the JDK emits their
			// hex-encoded DER; the values survive the RFC 4519 round trip as UTF8Strings.
			name: "directoryNameWithUnregisteredOIDs",
			file: certificateExtensionsTestFNMTDirectoryName,
			expected: [][2]string{
				{"RFC822_NAME", "miguelangel.nafria@seap.minhap.es"},
				{"DIRECTORY_NAME", "1.3.6.1.4.1.5734.1.1=#0c0c4d494755454c20414e47454c," +
					"1.3.6.1.4.1.5734.1.2=#0c064e4146524941," +
					"1.3.6.1.4.1.5734.1.3=#0c094c4153204845524153," +
					"1.3.6.1.4.1.5734.1.4=#0c09373238373432363858"},
			},
		},
		{
			// otherName keeps the complete DER of the [0] tagged GeneralName.
			name: "otherName",
			file: certificateExtensionsTestCertSignOtherName,
			expected: [][2]string{
				{"OTHER_NAME", "#a022060a2b060104018237140203a0140c12636f6e74616374406164722e676f762e726f"},
				{"RFC822_NAME", "contact@adr.gov.ro"},
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			certificateToken := certificateExtensionsKATToken(t, testCase.file)
			subjectAlternativeNames := CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken)
			if subjectAlternativeNames == nil {
				t.Fatalf("SubjectAlternativeNames: got nil, want a value")
			}
			generalNames := subjectAlternativeNames.GeneralNames()
			if len(generalNames) != len(testCase.expected) {
				t.Fatalf("GeneralNames() has %d entries, want %d", len(generalNames), len(testCase.expected))
			}
			for i, expected := range testCase.expected {
				if got := string(generalNames[i].GeneralNameType()); got != expected[0] {
					t.Errorf("general name %d type = %q, want %q", i, got, expected[0])
				}
				if got := generalNames[i].Value(); got != expected[1] {
					t.Errorf("general name %d value = %q, want %q", i, got, expected[1])
				}
			}
		})
	}
}

func TestCertificateExtensionsUtilsUpstreamNameConstraints(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestExcludedIPAddresses)
	nameConstraints := CertificateExtensionsUtilsNameConstraints(certificateToken)
	if nameConstraints == nil {
		t.Fatalf("NameConstraints: got nil, want a value")
	}
	if len(nameConstraints.PermittedSubtrees()) != 0 {
		t.Errorf("PermittedSubtrees() has %d entries, want 0", len(nameConstraints.PermittedSubtrees()))
	}
	excludedSubtrees := nameConstraints.ExcludedSubtrees()
	if len(excludedSubtrees) != 2 {
		t.Fatalf("ExcludedSubtrees() has %d entries, want 2", len(excludedSubtrees))
	}
	// An iPAddress name constraint carries the address AND its mask, so it is rendered as the raw
	// octets - unlike a subject alternative name, which BouncyCastle renders as a host address.
	for i, want := range []string{
		"#0000000000000000",
		"#0000000000000000000000000000000000000000000000000000000000000000",
	} {
		if got := excludedSubtrees[i].GeneralNameType(); got != enumerations.GeneralNameTypeIPAddress {
			t.Errorf("excluded subtree %d type = %q, want IP_ADDRESS", i, got)
		}
		if got := excludedSubtrees[i].Value(); got != want {
			t.Errorf("excluded subtree %d value = %q, want %q", i, got, want)
		}
		// GeneralSubtree#getMinimum() defaults to zero, #getMaximum() stays null.
		if minimum := excludedSubtrees[i].Minimum(); minimum == nil || minimum.Sign() != 0 {
			t.Errorf("excluded subtree %d minimum = %v, want 0", i, minimum)
		}
		if maximum := excludedSubtrees[i].Maximum(); maximum != nil {
			t.Errorf("excluded subtree %d maximum = %v, want nil", i, maximum)
		}
	}
}

func TestCertificateExtensionsUtilsUpstreamPolicyConstraints(t *testing.T) {
	certificateToken := certificateExtensionsKATToken(t, certificateExtensionsTestPolicyConstraints)
	policyConstraints := CertificateExtensionsUtilsPolicyConstraints(certificateToken)
	if policyConstraints == nil {
		t.Fatalf("PolicyConstraints: got nil, want a value")
	}
	if got := policyConstraints.RequireExplicitPolicy(); got != 0 {
		t.Errorf("RequireExplicitPolicy() = %d, want 0", got)
	}
	if got := policyConstraints.InhibitPolicyMapping(); got != -1 {
		t.Errorf("InhibitPolicyMapping() = %d, want -1", got)
	}

	certificateExtensions, err := CertificateExtensionsUtilsCertificateExtensions(certificateToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fromExtensions := certificateExtensions.PolicyConstraints()
	if fromExtensions == nil {
		t.Fatalf("CertificateExtensions.PolicyConstraints(): got nil, want a value")
	}
	if fromExtensions.RequireExplicitPolicy() != 0 || fromExtensions.InhibitPolicyMapping() != -1 {
		t.Errorf("CertificateExtensions.PolicyConstraints() = (%d, %d), want (0, -1)",
			fromExtensions.RequireExplicitPolicy(), fromExtensions.InhibitPolicyMapping())
	}
}

// TestCertificateExtensionsUtilsPredicates covers every OID predicate against the enumeration the
// extension model is built on.
func TestCertificateExtensionsUtilsPredicates(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		predicate func(string) bool
		value     enumerations.CertificateExtensionEnum
	}{
		{"SubjectAlternativeNames", CertificateExtensionsUtilsIsSubjectAlternativeNames, enumerations.CertificateExtensionEnumSubjectAlternativeName},
		{"AuthorityKeyIdentifier", CertificateExtensionsUtilsIsAuthorityKeyIdentifier, enumerations.CertificateExtensionEnumAuthorityKeyIdentifier},
		{"SubjectKeyIdentifier", CertificateExtensionsUtilsIsSubjectKeyIdentifier, enumerations.CertificateExtensionEnumSubjectKeyIdentifier},
		{"AuthorityInformationAccess", CertificateExtensionsUtilsIsAuthorityInformationAccess, enumerations.CertificateExtensionEnumAuthorityInformationAccess},
		{"CRLDistributionPoints", CertificateExtensionsUtilsIsCRLDistributionPoints, enumerations.CertificateExtensionEnumCRLDistributionPoints},
		{"BasicConstraints", CertificateExtensionsUtilsIsBasicConstraints, enumerations.CertificateExtensionEnumBasicConstraints},
		{"NameConstraints", CertificateExtensionsUtilsIsNameConstraints, enumerations.CertificateExtensionEnumNameConstraints},
		{"PolicyConstraints", CertificateExtensionsUtilsIsPolicyConstraints, enumerations.CertificateExtensionEnumPolicyConstraints},
		{"KeyUsage", CertificateExtensionsUtilsIsKeyUsage, enumerations.CertificateExtensionEnumKeyUsage},
		{"ExtendedKeyUsage", CertificateExtensionsUtilsIsExtendedKeyUsage, enumerations.CertificateExtensionEnumExtendedKeyUsage},
		{"InhibitAnyPolicy", CertificateExtensionsUtilsIsInhibitAnyPolicy, enumerations.CertificateExtensionEnumInhibitAnyPolicy},
		{"FreshestCRL", CertificateExtensionsUtilsIsFreshestCRL, enumerations.CertificateExtensionEnumFreshestCRL},
		{"CertificatePolicies", CertificateExtensionsUtilsIsCertificatePolicies, enumerations.CertificateExtensionEnumCertificatePolicies},
		{"OcspNoCheck", CertificateExtensionsUtilsIsOcspNoCheck, enumerations.CertificateExtensionEnumOCSPNoCheck},
		{"ValidityAssuredShortTerm", CertificateExtensionsUtilsIsValidityAssuredShortTerm, enumerations.CertificateExtensionEnumValidityAssuredShortTerm},
		{"QcStatements", CertificateExtensionsUtilsIsQcStatements, enumerations.CertificateExtensionEnumQCStatements},
		{"NoRevocationAvailable", CertificateExtensionsUtilsIsNoRevocationAvailable, enumerations.CertificateExtensionEnumNoRevocationAvailable},
	} {
		if !testCase.predicate(testCase.value.OID()) {
			t.Errorf("Is%s(%q) = false, want true", testCase.name, testCase.value.OID())
		}
		if testCase.predicate("1.2.3.4.5") {
			t.Errorf("Is%s(1.2.3.4.5) = true, want false", testCase.name)
		}
	}

	// The OIDs upstream reads from BouncyCastle and from OID.java rather than from the enumeration.
	for _, testCase := range []struct {
		name  string
		value enumerations.CertificateExtensionEnum
		oid   string
	}{
		{"ocsp-nocheck", enumerations.CertificateExtensionEnumOCSPNoCheck, "1.3.6.1.5.5.7.48.1.5"},
		{"ext-etsi-valassured-ST-certs", enumerations.CertificateExtensionEnumValidityAssuredShortTerm, "0.4.0.194121.2.1"},
		{"qcStatements", enumerations.CertificateExtensionEnumQCStatements, "1.3.6.1.5.5.7.1.3"},
		{"noRevAvail", enumerations.CertificateExtensionEnumNoRevocationAvailable, "2.5.29.56"},
	} {
		if got := testCase.value.OID(); got != testCase.oid {
			t.Errorf("%s OID = %q, want %q", testCase.name, got, testCase.oid)
		}
	}
}

// TestCertificateExtensionsUtilsBasicConstraintsAlwaysBuilt records that, like upstream,
// getBasicConstraints answers an object even for a certificate without the extension.
func TestCertificateExtensionsUtilsBasicConstraintsAlwaysBuilt(t *testing.T) {
	var withoutBasicConstraints *model.CertificateToken
	for _, answers := range certificateExtensionsKATLoad(t) {
		name := answers.str(t, "file")
		if certificateExtensionsKATUnparseable[name] {
			continue
		}
		if answers.isNull(t, "bc.octets") {
			withoutBasicConstraints = certificateExtensionsKATToken(t, name)
			break
		}
	}
	if withoutBasicConstraints == nil {
		t.Skip("the fixture corpus holds no certificate without a basicConstraints extension")
	}
	basicConstraints := CertificateExtensionsUtilsBasicConstraints(withoutBasicConstraints)
	if basicConstraints == nil {
		t.Fatalf("BasicConstraints: got nil, want a value")
	}
	if basicConstraints.Octets() != nil {
		t.Errorf("Octets() = %v, want nil", basicConstraints.Octets())
	}
	if basicConstraints.IsCa() {
		t.Errorf("IsCa() = true, want false")
	}
	if got := basicConstraints.PathLenConstraint(); got != -1 {
		t.Errorf("PathLenConstraint() = %d, want -1", got)
	}
}

func certificateExtensionsTestEqual(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
