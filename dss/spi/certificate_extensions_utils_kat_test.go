package spi

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/internal/corpustest"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/extension"
	"github.com/utain/esig/dss/utils"
)

// The known answers in testdata/certificate_extensions/kat.tsv were produced by a Java 21 program
// running the very algorithms of CertificateExtensionsUtils / QcStatementUtils / SignerIdentifier
// on top of BouncyCastle 1.78.1 and the JDK X500Principal, over the certificates in
// testdata/certificate_extensions (the certificates upstream's CertificateExtensionUtilsTest and
// QcStatementsUtilsTest exercise, plus the other dss-spi test resources). Certificates are loaded
// through the BouncyCastle JCE provider, exactly as DSSUtils does.
//
// Format: "<certificate index>\t<key>\t<value>" where the value is either "-" for a Java null or
// "b" followed by the base64 of its UTF-8 encoding.

// certificateExtensionsKATUnparseable lists the fixtures crypto/x509 rejects but BouncyCastle
// accepts. They are skipped, and the test fails if the set ever changes. The set depends on the
// Go toolchain — newer releases tighten the parser — so it is defined in build-tagged files:
// certificate_extensions_kat_unparseable_legacy_test.go (pre-1.27) and
// certificate_extensions_kat_unparseable_go127_test.go (1.27+).

type certificateExtensionsKAT map[string]*string

func (k certificateExtensionsKAT) value(t *testing.T, key string) *string {
	t.Helper()
	value, present := k[key]
	if !present {
		t.Fatalf("missing known answer for key %q", key)
	}
	return value
}

// str answers the known answer for key, failing when the fixture holds a Java null.
func (k certificateExtensionsKAT) str(t *testing.T, key string) string {
	t.Helper()
	value := k.value(t, key)
	if value == nil {
		t.Fatalf("known answer for key %q is null", key)
	}
	return *value
}

// orEmpty answers the known answer for key, mapping a Java null to the empty string.
func (k certificateExtensionsKAT) orEmpty(t *testing.T, key string) string {
	t.Helper()
	if value := k.value(t, key); value != nil {
		return *value
	}
	return ""
}

func (k certificateExtensionsKAT) isNull(t *testing.T, key string) bool {
	t.Helper()
	return k.value(t, key) == nil
}

func (k certificateExtensionsKAT) boolean(t *testing.T, key string) bool {
	t.Helper()
	return k.str(t, key) == "true"
}

func (k certificateExtensionsKAT) integer(t *testing.T, key string) int {
	t.Helper()
	value, err := strconv.Atoi(k.str(t, key))
	if err != nil {
		t.Fatalf("known answer for key %q is not an integer: %v", key, err)
	}
	return value
}

// list reads a "<prefix>.count" / "<prefix>.<i>" group.
func (k certificateExtensionsKAT) list(t *testing.T, prefix string) []string {
	t.Helper()
	count := k.integer(t, prefix+".count")
	values := make([]string, 0, count)
	for i := 0; i < count; i++ {
		values = append(values, k.str(t, fmt.Sprintf("%s.%d", prefix, i)))
	}
	return values
}

func certificateExtensionsKATLoad(t *testing.T) []certificateExtensionsKAT {
	t.Helper()
	raw, err := os.ReadFile(corpustest.Path(t, filepath.Join("certificate_extensions", "kat.tsv")))
	if err != nil {
		t.Fatalf("unable to read the known answers: %v", err)
	}
	answers := make([]certificateExtensionsKAT, 0)
	for _, line := range strings.Split(strings.TrimRight(string(raw), "\n"), "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			t.Fatalf("malformed known answer line %q", line)
		}
		index, err := strconv.Atoi(fields[0])
		if err != nil {
			t.Fatalf("malformed known answer index in %q", line)
		}
		for len(answers) <= index {
			answers = append(answers, certificateExtensionsKAT{})
		}
		if fields[2] == "-" {
			answers[index][fields[1]] = nil
			continue
		}
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(fields[2], "b"))
		if err != nil {
			t.Fatalf("malformed known answer value in %q: %v", line, err)
		}
		value := string(decoded)
		answers[index][fields[1]] = &value
	}
	return answers
}

func certificateExtensionsKATToken(t *testing.T, name string) *model.CertificateToken {
	t.Helper()
	der, err := os.ReadFile(filepath.Join("testdata", "certificate_extensions", name))
	if err != nil {
		t.Fatalf("unable to read %s: %v", name, err)
	}
	certificate, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("unable to parse %s: %v", name, err)
	}
	certificateToken, err := model.NewCertificateToken(certificate)
	if err != nil {
		t.Fatalf("unable to build a CertificateToken for %s: %v", name, err)
	}
	return certificateToken
}

// TestCertificateExtensionsKATCorpus guards the fixture corpus itself: every certificate must be
// readable by crypto/x509 except the two documented ones.
func TestCertificateExtensionsKATCorpus(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "certificate_extensions", "*.der"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no certificate fixtures found: %v", err)
	}
	unparseable := map[string]bool{}
	for _, file := range files {
		der, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("unable to read %s: %v", file, err)
		}
		if _, err := x509.ParseCertificate(der); err != nil {
			unparseable[filepath.Base(file)] = true
		}
	}
	if len(unparseable) != len(certificateExtensionsKATUnparseable) {
		t.Fatalf("unexpected set of certificates crypto/x509 rejects: %v", unparseable)
	}
	for name := range certificateExtensionsKATUnparseable {
		if !unparseable[name] {
			t.Errorf("%s is expected to be rejected by crypto/x509 but was parsed", name)
		}
	}
}

func TestCertificateExtensionsUtilsKnownAnswers(t *testing.T) {
	for index, answers := range certificateExtensionsKATLoad(t) {
		name := answers.str(t, "file")
		if certificateExtensionsKATUnparseable[name] {
			continue
		}
		t.Run(fmt.Sprintf("%02d_%s", index, name), func(t *testing.T) {
			certificateToken := certificateExtensionsKATToken(t, name)
			certificateExtensionsKATAssertSubjectAlternativeNames(t, certificateToken, answers)
			certificateExtensionsKATAssertAuthorityInformationAccess(t, certificateToken, answers)
			certificateExtensionsKATAssertKeyIdentifiers(t, certificateToken, answers)
			certificateExtensionsKATAssertCRLDistributionPoints(t, certificateToken, answers)
			certificateExtensionsKATAssertBasicConstraints(t, certificateToken, answers)
			certificateExtensionsKATAssertNameConstraints(t, certificateToken, answers)
			certificateExtensionsKATAssertPolicies(t, certificateToken, answers)
			certificateExtensionsKATAssertKeyUsages(t, certificateToken, answers)
			certificateExtensionsKATAssertNullIdentifiedFlags(t, certificateToken, answers)
			certificateExtensionsKATAssertCertificateExtensions(t, certificateToken, answers)
		})
	}
}

func certificateExtensionsKATAssertSubjectAlternativeNames(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	subjectAlternativeNames := CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken)
	if answers.isNull(t, "san") {
		if subjectAlternativeNames != nil {
			t.Fatalf("SubjectAlternativeNames: got a value, want nil")
		}
		return
	}
	if subjectAlternativeNames == nil {
		t.Fatalf("SubjectAlternativeNames: got nil, want a value")
	}
	certificateExtensionsKATAssertOptionalHex(t, "san.octets", subjectAlternativeNames.Octets(), answers)
	if got, want := subjectAlternativeNames.IsCritical(), answers.boolean(t, "san.critical"); got != want {
		t.Errorf("SubjectAlternativeNames.IsCritical() = %v, want %v", got, want)
	}
	generalNames := subjectAlternativeNames.GeneralNames()
	if got, want := len(generalNames), answers.integer(t, "san.count"); got != want {
		t.Fatalf("SubjectAlternativeNames.GeneralNames() has %d entries, want %d", got, want)
	}
	for i, generalName := range generalNames {
		if got, want := string(generalName.GeneralNameType()), answers.str(t, fmt.Sprintf("san.%d.type", i)); got != want {
			t.Errorf("subject alternative name %d type = %q, want %q", i, got, want)
		}
		if got, want := generalName.Value(), answers.str(t, fmt.Sprintf("san.%d.value", i)); got != want {
			t.Errorf("subject alternative name %d value = %q, want %q", i, got, want)
		}
	}
}

func certificateExtensionsKATAssertAuthorityInformationAccess(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	aia := CertificateExtensionsUtilsAuthorityInformationAccess(certificateToken)
	caIssuersUrls := CertificateExtensionsUtilsCAIssuersAccessUrls(certificateToken)
	ocspUrls := CertificateExtensionsUtilsOCSPAccessUrls(certificateToken)
	certificateExtensionsKATAssertStrings(t, "caIssuersUrls", caIssuersUrls, answers)
	certificateExtensionsKATAssertStrings(t, "ocspUrls", ocspUrls, answers)
	if answers.isNull(t, "aia") {
		if aia != nil {
			t.Fatalf("AuthorityInformationAccess: got a value, want nil")
		}
		return
	}
	if aia == nil {
		t.Fatalf("AuthorityInformationAccess: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "aia.octets", aia.Octets(), answers)
	if got, want := aia.IsCritical(), answers.boolean(t, "aia.critical"); got != want {
		t.Errorf("AuthorityInformationAccess.IsCritical() = %v, want %v", got, want)
	}
	certificateExtensionsKATAssertStrings(t, "aia.caIssuers", aia.CaIssuers(), answers)
	certificateExtensionsKATAssertStrings(t, "aia.ocsp", aia.Ocsp(), answers)
}

func certificateExtensionsKATAssertKeyIdentifiers(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	aki, err := CertificateExtensionsUtilsAuthorityKeyIdentifier(certificateToken)
	if err != nil {
		t.Fatalf("AuthorityKeyIdentifier: unexpected error %v", err)
	}
	if answers.isNull(t, "aki") {
		if aki != nil {
			t.Errorf("AuthorityKeyIdentifier: got a value, want nil")
		}
	} else {
		if aki == nil {
			t.Fatalf("AuthorityKeyIdentifier: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "aki.octets", aki.Octets(), answers)
		if got, want := aki.IsCritical(), answers.boolean(t, "aki.critical"); got != want {
			t.Errorf("AuthorityKeyIdentifier.IsCritical() = %v, want %v", got, want)
		}
		certificateExtensionsKATAssertOptionalHex(t, "aki.keyIdentifier", aki.KeyIdentifier(), answers)
		certificateExtensionsKATAssertOptionalHex(t, "aki.issuerSerial", aki.AuthorityCertIssuerSerial(), answers)
	}

	ski, err := CertificateExtensionsUtilsSubjectKeyIdentifier(certificateToken)
	if err != nil {
		t.Fatalf("SubjectKeyIdentifier: unexpected error %v", err)
	}
	if answers.isNull(t, "ski") {
		if ski != nil {
			t.Errorf("SubjectKeyIdentifier: got a value, want nil")
		}
		return
	}
	if ski == nil {
		t.Fatalf("SubjectKeyIdentifier: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "ski.octets", ski.Octets(), answers)
	if got, want := ski.IsCritical(), answers.boolean(t, "ski.critical"); got != want {
		t.Errorf("SubjectKeyIdentifier.IsCritical() = %v, want %v", got, want)
	}
	certificateExtensionsKATAssertHex(t, "ski.ski", ski.Ski(), answers)
}

func certificateExtensionsKATAssertCRLDistributionPoints(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	crlDistributionPoints := CertificateExtensionsUtilsCRLDistributionPoints(certificateToken)
	certificateExtensionsKATAssertStrings(t, "crlUrls", CertificateExtensionsUtilsCRLAccessUrls(certificateToken), answers)
	if answers.isNull(t, "crldp") {
		if crlDistributionPoints != nil {
			t.Errorf("CRLDistributionPoints: got a value, want nil")
		}
	} else {
		if crlDistributionPoints == nil {
			t.Fatalf("CRLDistributionPoints: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "crldp.octets", crlDistributionPoints.Octets(), answers)
		if got, want := crlDistributionPoints.IsCritical(), answers.boolean(t, "crldp.critical"); got != want {
			t.Errorf("CRLDistributionPoints.IsCritical() = %v, want %v", got, want)
		}
		certificateExtensionsKATAssertStrings(t, "crldp.urls", crlDistributionPoints.CrlUrls(), answers)
	}

	freshestCRL := CertificateExtensionsUtilsFreshestCRL(certificateToken)
	if answers.isNull(t, "freshest") {
		if freshestCRL != nil {
			t.Errorf("FreshestCRL: got a value, want nil")
		}
		return
	}
	if freshestCRL == nil {
		t.Fatalf("FreshestCRL: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "freshest.octets", freshestCRL.Octets(), answers)
	if got, want := freshestCRL.IsCritical(), answers.boolean(t, "freshest.critical"); got != want {
		t.Errorf("FreshestCRL.IsCritical() = %v, want %v", got, want)
	}
	certificateExtensionsKATAssertStrings(t, "freshest.urls", freshestCRL.CrlUrls(), answers)
}

func certificateExtensionsKATAssertBasicConstraints(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	basicConstraints := CertificateExtensionsUtilsBasicConstraints(certificateToken)
	if basicConstraints == nil {
		t.Fatalf("BasicConstraints: got nil, want a value")
	}
	certificateExtensionsKATAssertOptionalHex(t, "bc.octets", basicConstraints.Octets(), answers)
	if got, want := basicConstraints.IsCritical(), answers.boolean(t, "bc.critical"); got != want {
		t.Errorf("BasicConstraints.IsCritical() = %v, want %v", got, want)
	}
	value := answers.integer(t, "bc.value")
	if got, want := basicConstraints.PathLenConstraint(), value; got != want {
		t.Errorf("BasicConstraints.PathLenConstraint() = %d, want %d", got, want)
	}
	if got, want := basicConstraints.IsCa(), value != -1; got != want {
		t.Errorf("BasicConstraints.IsCa() = %v, want %v", got, want)
	}
}

func certificateExtensionsKATAssertNameConstraints(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	nameConstraints := CertificateExtensionsUtilsNameConstraints(certificateToken)
	if answers.isNull(t, "nc") {
		if nameConstraints != nil {
			t.Fatalf("NameConstraints: got a value, want nil")
		}
		return
	}
	if nameConstraints == nil {
		t.Fatalf("NameConstraints: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "nc.octets", nameConstraints.Octets(), answers)
	if got, want := nameConstraints.IsCritical(), answers.boolean(t, "nc.critical"); got != want {
		t.Errorf("NameConstraints.IsCritical() = %v, want %v", got, want)
	}

	for _, subtrees := range []struct {
		prefix string
		values []*extensionGeneralSubtree
	}{
		{"nc.permitted", certificateExtensionsKATSubtrees(nameConstraints.PermittedSubtrees())},
		{"nc.excluded", certificateExtensionsKATSubtrees(nameConstraints.ExcludedSubtrees())},
	} {
		if got, want := len(subtrees.values), answers.integer(t, subtrees.prefix+".count"); got != want {
			t.Fatalf("%s has %d entries, want %d", subtrees.prefix, got, want)
		}
		for i, subtree := range subtrees.values {
			if got, want := subtree.generalNameType, answers.str(t, fmt.Sprintf("%s.%d.type", subtrees.prefix, i)); got != want {
				t.Errorf("%s.%d type = %q, want %q", subtrees.prefix, i, got, want)
			}
			if got, want := subtree.minimum, answers.orEmpty(t, fmt.Sprintf("%s.%d.min", subtrees.prefix, i)); got != want {
				t.Errorf("%s.%d minimum = %q, want %q", subtrees.prefix, i, got, want)
			}
			if got, want := subtree.maximum, answers.orEmpty(t, fmt.Sprintf("%s.%d.max", subtrees.prefix, i)); got != want {
				t.Errorf("%s.%d maximum = %q, want %q", subtrees.prefix, i, got, want)
			}
			if got, want := subtree.value, answers.str(t, fmt.Sprintf("%s.%d.value", subtrees.prefix, i)); got != want {
				t.Errorf("%s.%d value = %q, want %q", subtrees.prefix, i, got, want)
			}
		}
	}
}

// extensionGeneralSubtree flattens a GeneralSubtree for comparison with the fixture, where a null
// BigInteger is rendered as the empty string.
type extensionGeneralSubtree struct {
	generalNameType string
	minimum         string
	maximum         string
	value           string
}

func certificateExtensionsKATSubtrees(subtrees []*extension.GeneralSubtree) []*extensionGeneralSubtree {
	flattened := make([]*extensionGeneralSubtree, 0, len(subtrees))
	for _, subtree := range subtrees {
		flat := &extensionGeneralSubtree{
			generalNameType: string(subtree.GeneralNameType()),
			value:           subtree.Value(),
		}
		if minimum := subtree.Minimum(); minimum != nil {
			flat.minimum = minimum.String()
		}
		if maximum := subtree.Maximum(); maximum != nil {
			flat.maximum = maximum.String()
		}
		flattened = append(flattened, flat)
	}
	return flattened
}

func certificateExtensionsKATAssertPolicies(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	policyConstraints := CertificateExtensionsUtilsPolicyConstraints(certificateToken)
	if answers.isNull(t, "pc") {
		if policyConstraints != nil {
			t.Errorf("PolicyConstraints: got a value, want nil")
		}
	} else {
		if policyConstraints == nil {
			t.Fatalf("PolicyConstraints: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "pc.octets", policyConstraints.Octets(), answers)
		if got, want := policyConstraints.IsCritical(), answers.boolean(t, "pc.critical"); got != want {
			t.Errorf("PolicyConstraints.IsCritical() = %v, want %v", got, want)
		}
		if got, want := policyConstraints.RequireExplicitPolicy(), answers.integer(t, "pc.requireExplicitPolicy"); got != want {
			t.Errorf("PolicyConstraints.RequireExplicitPolicy() = %d, want %d", got, want)
		}
		if got, want := policyConstraints.InhibitPolicyMapping(), answers.integer(t, "pc.inhibitPolicyMapping"); got != want {
			t.Errorf("PolicyConstraints.InhibitPolicyMapping() = %d, want %d", got, want)
		}
	}

	inhibitAnyPolicy := CertificateExtensionsUtilsInhibitAnyPolicy(certificateToken)
	if answers.isNull(t, "iap") {
		if inhibitAnyPolicy != nil {
			t.Errorf("InhibitAnyPolicy: got a value, want nil")
		}
	} else {
		if inhibitAnyPolicy == nil {
			t.Fatalf("InhibitAnyPolicy: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "iap.octets", inhibitAnyPolicy.Octets(), answers)
		if got, want := inhibitAnyPolicy.IsCritical(), answers.boolean(t, "iap.critical"); got != want {
			t.Errorf("InhibitAnyPolicy.IsCritical() = %v, want %v", got, want)
		}
		if got, want := inhibitAnyPolicy.Value(), answers.integer(t, "iap.value"); got != want {
			t.Errorf("InhibitAnyPolicy.Value() = %d, want %d", got, want)
		}
	}

	certificatePolicies := CertificateExtensionsUtilsCertificatePolicies(certificateToken)
	if answers.isNull(t, "cp") {
		if certificatePolicies != nil {
			t.Errorf("CertificatePolicies: got a value, want nil")
		}
		return
	}
	if certificatePolicies == nil {
		t.Fatalf("CertificatePolicies: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "cp.octets", certificatePolicies.Octets(), answers)
	if got, want := certificatePolicies.IsCritical(), answers.boolean(t, "cp.critical"); got != want {
		t.Errorf("CertificatePolicies.IsCritical() = %v, want %v", got, want)
	}
	policyList := certificatePolicies.PolicyList()
	if got, want := len(policyList), answers.integer(t, "cp.count"); got != want {
		t.Fatalf("CertificatePolicies.PolicyList() has %d entries, want %d", got, want)
	}
	for i, policy := range policyList {
		if got, want := policy.Oid(), answers.str(t, fmt.Sprintf("cp.%d.oid", i)); got != want {
			t.Errorf("certificate policy %d OID = %q, want %q", i, got, want)
		}
		if got, want := policy.CpsUrl(), answers.orEmpty(t, fmt.Sprintf("cp.%d.cpsUrl", i)); got != want {
			t.Errorf("certificate policy %d CPS URL = %q, want %q", i, got, want)
		}
	}
}

func certificateExtensionsKATAssertKeyUsages(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	keyUsage := CertificateExtensionsUtilsKeyUsage(certificateToken)
	if answers.isNull(t, "ku") {
		if keyUsage != nil {
			t.Errorf("KeyUsage: got a value, want nil")
		}
	} else {
		if keyUsage == nil {
			t.Fatalf("KeyUsage: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "ku.octets", keyUsage.Octets(), answers)
		if got, want := keyUsage.IsCritical(), answers.boolean(t, "ku.critical"); got != want {
			t.Errorf("KeyUsage.IsCritical() = %v, want %v", got, want)
		}
		bits := make([]string, 0)
		for _, bit := range keyUsage.KeyUsageBits() {
			bits = append(bits, string(bit))
		}
		certificateExtensionsKATAssertStrings(t, "ku.bits", bits, answers)
	}

	extendedKeyUsage := CertificateExtensionsUtilsExtendedKeyUsage(certificateToken)
	if answers.isNull(t, "eku") {
		if extendedKeyUsage != nil {
			t.Errorf("ExtendedKeyUsage: got a value, want nil")
		}
		return
	}
	if extendedKeyUsage == nil {
		t.Fatalf("ExtendedKeyUsage: got nil, want a value")
	}
	certificateExtensionsKATAssertOptionalHex(t, "eku.octets", extendedKeyUsage.Octets(), answers)
	if got, want := extendedKeyUsage.IsCritical(), answers.boolean(t, "eku.critical"); got != want {
		t.Errorf("ExtendedKeyUsage.IsCritical() = %v, want %v", got, want)
	}
	if answers.isNull(t, "eku.oids.count") {
		if extendedKeyUsage.Oids() != nil {
			t.Errorf("ExtendedKeyUsage.Oids() = %v, want nil", extendedKeyUsage.Oids())
		}
		return
	}
	certificateExtensionsKATAssertStrings(t, "eku.oids", extendedKeyUsage.Oids(), answers)
}

func certificateExtensionsKATAssertNullIdentifiedFlags(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	ocspNoCheck := CertificateExtensionsUtilsOcspNoCheck(certificateToken)
	if answers.isNull(t, "onc") {
		if ocspNoCheck != nil {
			t.Errorf("OcspNoCheck: got a value, want nil")
		}
		if CertificateExtensionsUtilsHasOcspNoCheckExtension(certificateToken) {
			t.Errorf("HasOcspNoCheckExtension() = true, want false")
		}
	} else {
		if ocspNoCheck == nil {
			t.Fatalf("OcspNoCheck: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "onc.octets", ocspNoCheck.Octets(), answers)
		if got, want := ocspNoCheck.IsCritical(), answers.boolean(t, "onc.critical"); got != want {
			t.Errorf("OcspNoCheck.IsCritical() = %v, want %v", got, want)
		}
		want := answers.boolean(t, "onc.value")
		if got := ocspNoCheck.IsOcspNoCheck(); got != want {
			t.Errorf("OcspNoCheck.IsOcspNoCheck() = %v, want %v", got, want)
		}
		if got := CertificateExtensionsUtilsHasOcspNoCheckExtension(certificateToken); got != want {
			t.Errorf("HasOcspNoCheckExtension() = %v, want %v", got, want)
		}
	}

	validityAssuredShortTerm := CertificateExtensionsUtilsValAssuredSTCerts(certificateToken)
	if answers.isNull(t, "vast") {
		if validityAssuredShortTerm != nil {
			t.Errorf("ValAssuredSTCerts: got a value, want nil")
		}
		if CertificateExtensionsUtilsHasValAssuredShortTermCertsExtension(certificateToken) {
			t.Errorf("HasValAssuredShortTermCertsExtension() = true, want false")
		}
	} else {
		if validityAssuredShortTerm == nil {
			t.Fatalf("ValAssuredSTCerts: got nil, want a value")
		}
		certificateExtensionsKATAssertHex(t, "vast.octets", validityAssuredShortTerm.Octets(), answers)
		if got, want := validityAssuredShortTerm.IsCritical(), answers.boolean(t, "vast.critical"); got != want {
			t.Errorf("ValAssuredSTCerts.IsCritical() = %v, want %v", got, want)
		}
		want := answers.boolean(t, "vast.value")
		if got := validityAssuredShortTerm.IsValAssuredSTCerts(); got != want {
			t.Errorf("ValAssuredSTCerts.IsValAssuredSTCerts() = %v, want %v", got, want)
		}
		if got := CertificateExtensionsUtilsHasValAssuredShortTermCertsExtension(certificateToken); got != want {
			t.Errorf("HasValAssuredShortTermCertsExtension() = %v, want %v", got, want)
		}
	}

	noRevAvail := CertificateExtensionsUtilsNoRevAvail(certificateToken)
	if answers.isNull(t, "nra") {
		if noRevAvail != nil {
			t.Errorf("NoRevAvail: got a value, want nil")
		}
		return
	}
	if noRevAvail == nil {
		t.Fatalf("NoRevAvail: got nil, want a value")
	}
	certificateExtensionsKATAssertHex(t, "nra.octets", noRevAvail.Octets(), answers)
	if got, want := noRevAvail.IsCritical(), answers.boolean(t, "nra.critical"); got != want {
		t.Errorf("NoRevAvail.IsCritical() = %v, want %v", got, want)
	}
	if got, want := noRevAvail.IsNoRevAvail(), answers.boolean(t, "nra.value"); got != want {
		t.Errorf("NoRevAvail.IsNoRevAvail() = %v, want %v", got, want)
	}
}

// certificateExtensionsKATAssertCertificateExtensions checks getCertificateExtensions() against
// the OIDs Java collected. Their order is BouncyCastle's HashSet order upstream and DER order here,
// so only the contents are compared (see the deviation noted on the ported function). An extension
// whose getter answers null is recorded on neither the typed field nor the aggregate list, so the
// expectation is reconstructed from the per-extension known answers.
func certificateExtensionsKATAssertCertificateExtensions(t *testing.T, certificateToken *model.CertificateToken, answers certificateExtensionsKAT) {
	t.Helper()
	certificateExtensions, err := CertificateExtensionsUtilsCertificateExtensions(certificateToken)
	if err != nil {
		t.Fatalf("CertificateExtensions: unexpected error %v", err)
	}
	got := make([]string, 0)
	for _, certificateExtension := range certificateExtensions.AllCertificateExtensions() {
		got = append(got, certificateExtension.OID())
	}

	// The key of the known answer telling whether the getter of a supported extension answered null.
	// basicConstraints has no entry: upstream always builds it.
	supported := map[string]string{
		"2.5.29.17":            "san",
		"2.5.29.35":            "aki",
		"2.5.29.14":            "ski",
		"1.3.6.1.5.5.7.1.1":    "aia",
		"2.5.29.31":            "crldp",
		"2.5.29.19":            "",
		"2.5.29.30":            "nc",
		"2.5.29.36":            "pc",
		"2.5.29.54":            "iap",
		"2.5.29.46":            "freshest",
		"2.5.29.15":            "ku",
		"2.5.29.37":            "eku",
		"2.5.29.32":            "cp",
		"1.3.6.1.5.5.7.48.1.5": "onc",
		"0.4.0.194121.2.1":     "vast",
		"1.3.6.1.5.5.7.1.3":    "qc",
		"2.5.29.56":            "nra",
	}
	want := make([]string, 0)
	for _, oid := range answers.list(t, "allOids") {
		key, isSupported := supported[oid]
		if !isSupported || key == "" || !answers.isNull(t, key) {
			want = append(want, oid)
		}
	}
	sort.Strings(got)
	sort.Strings(want)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("CertificateExtensions OIDs = %v, want %v", got, want)
	}

	// The aggregate must expose the very same objects the individual getters build.
	if subjectAlternativeNames := CertificateExtensionsUtilsSubjectAlternativeNames(certificateToken); subjectAlternativeNames != nil &&
		certificateExtensionsKATHasOID(want, enumerations.CertificateExtensionEnum_SUBJECT_ALTERNATIVE_NAME.OID()) {
		aggregated := certificateExtensions.SubjectAlternativeNames()
		if aggregated == nil || len(aggregated.GeneralNames()) != len(subjectAlternativeNames.GeneralNames()) {
			t.Errorf("CertificateExtensions.SubjectAlternativeNames() disagrees with the direct getter")
		}
	}
	if certificateExtensionsKATHasOID(want, enumerations.CertificateExtensionEnum_QC_STATEMENTS.OID()) {
		if certificateExtensions.QcStatements() == nil && CertificateExtensionsUtilsQcStatements(certificateToken) != nil {
			t.Errorf("CertificateExtensions.QcStatements() disagrees with the direct getter")
		}
	}
}

func certificateExtensionsKATHasOID(oids []string, oid string) bool {
	for _, value := range oids {
		if value == oid {
			return true
		}
	}
	return false
}

func certificateExtensionsKATAssertHex(t *testing.T, key string, got []byte, answers certificateExtensionsKAT) {
	t.Helper()
	if want := answers.str(t, key); utils.ToHex(got) != want {
		t.Errorf("%s = %s, want %s", key, utils.ToHex(got), want)
	}
}

func certificateExtensionsKATAssertOptionalHex(t *testing.T, key string, got []byte, answers certificateExtensionsKAT) {
	t.Helper()
	if answers.isNull(t, key) {
		if got != nil {
			t.Errorf("%s = %s, want nil", key, utils.ToHex(got))
		}
		return
	}
	certificateExtensionsKATAssertHex(t, key, got, answers)
}

func certificateExtensionsKATAssertStrings(t *testing.T, prefix string, got []string, answers certificateExtensionsKAT) {
	t.Helper()
	want := answers.list(t, prefix)
	if len(got) != len(want) {
		t.Fatalf("%s = %v, want %v", prefix, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d] = %q, want %q", prefix, i, got[i], want[i])
		}
	}
}
