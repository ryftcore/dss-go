// Cross-validation harness: determinism guard for the trusted-source assembly.
//
// TLValidationJob's synchronization step accumulates certificates into maps that Java keys by
// CertificateToken (whose equals/hashCode is the certificate's DSS id) and drains through
// entrySet(). A java.util.HashMap's entry order is arbitrary but, for a given key set, identical
// on every run, so upstream's resulting TrustedListsCertificateSource content comes out in a
// stable order; Go randomises map iteration per range, so the naive translation - a
// map[*model.CertificateToken]... ranged directly - makes getCertificates(), and everything
// downstream that consumes it up to diagnostic data, a different permutation on every run.
//
// This test pins that down where it is observable: it builds the certificate source from the same
// real Slovak trusted list contract item (D) uses, several times in one process, and requires the
// full certificate list AND each certificate's TrustProperties to come out identically every
// time. It is deliberately NOT a comparison against a Java dump (Java's own arbitrary order is a
// function of the certificate digest's byte-array hashCode and is not reproducible here) - what
// must hold is that the Go side is stable, which contract items (C) and (D) cannot see because
// both sort their dumps before comparing.
package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ryftcore/dss-go/dss/spi"
	spitsl "github.com/ryftcore/dss-go/dss/spi/tsl"
	"github.com/ryftcore/dss-go/dss/tsl"
)

// tsdBuildCertificateSource runs contract item (D)'s offline TLValidationJob once and returns the
// resulting trusted certificate source.
func tsdBuildCertificateSource(t *testing.T) *spitsl.TrustedListsCertificateSource {
	t.Helper()

	issuerSource := spi.NewCommonCertificateSource()
	issuerSource.AddCertificate(cqLoadCertificate(t, cqTLIssuer))

	tlSource := tsl.NewTLSource()
	tlSource.SetUrl("sk-tl.xml")
	tlSource.SetTLVersions([]int{5, 6})
	tlSource.SetCertificateSource(&issuerSource)

	tlValidationJob := tsl.NewTLValidationJob()
	tlValidationJob.SetTrustedListSources(tlSource)
	tlValidationJob.SetOfflineDataLoader(newCqFileLoader(t))
	certificateSource := spitsl.NewTrustedListsCertificateSource()
	tlValidationJob.SetTrustedListCertificateSource(certificateSource)

	if err := tlValidationJob.OfflineRefresh(); err != nil {
		t.Fatalf("OfflineRefresh: %v", err)
	}
	return certificateSource
}

// tsdFingerprint renders the whole source - certificate order included - as one comparable
// string.
func tsdFingerprint(source *spitsl.TrustedListsCertificateSource) string {
	var builder strings.Builder
	for _, certificate := range source.Certificates() {
		sum := sha256.Sum256(certificate.Encoded())
		builder.WriteString(hex.EncodeToString(sum[:]))
		for _, trustProperties := range source.TrustServices(certificate) {
			builder.WriteByte('|')
			if tlInfo := trustProperties.TLInfo(); tlInfo != nil {
				builder.WriteString(tlInfo.Url())
			}
			builder.WriteByte('#')
			if trustService := trustProperties.TrustService(); trustService != nil {
				for entry := range trustService.Iterator() {
					builder.WriteString(entry.Type())
					builder.WriteByte(',')
					builder.WriteString(entry.Status())
					builder.WriteByte(';')
				}
			}
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

// TestTrustedSourceAssemblyIsDeterministic requires repeated identical runs to produce a
// byte-identical trusted certificate source.
func TestTrustedSourceAssemblyIsDeterministic(t *testing.T) {
	first := tsdFingerprint(tsdBuildCertificateSource(t))
	if strings.Count(first, "\n") < 200 {
		t.Fatalf("fingerprint covers only %d certificates; the fixture is not exercising the "+
			"synchronizer any more", strings.Count(first, "\n"))
	}
	for run := 1; run < 8; run++ {
		got := tsdFingerprint(tsdBuildCertificateSource(t))
		if got == first {
			continue
		}
		// Report whether the difference is a mere permutation, which is what a map-iteration
		// order leak looks like, as opposed to genuinely different content.
		counts := make(map[string]int)
		for _, line := range strings.Split(first, "\n") {
			counts[line]++
		}
		for _, line := range strings.Split(got, "\n") {
			counts[line]--
		}
		permutation := true
		for _, count := range counts {
			if count != 0 {
				permutation = false
				break
			}
		}
		t.Fatalf("run %d differs from run 0 (same content, different order: %v) - a map "+
			"iteration order is leaking into the trusted certificate source", run, permutation)
	}
}
