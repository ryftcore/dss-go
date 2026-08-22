// DigestCalculatorProvider is the shared, minimal replacement for
// org.bouncycastle.operator.DigestCalculatorProvider that
// dss-cms/src/main/java/eu/europa/esig/dss/cms/operator/{CustomMessageDigestCalculatorProvider,PrecomputedDigestCalculatorProvider}.java
// (DSS 6.5.RC1) both implement; see custom_message_digest_calculator_provider.go and
// precomputed_digest_calculator_provider.go for the ports of those two classes, and doc.go for
// why the interface collapses BC's Provider/Calculator pair into one method.
package cms

import "github.com/ryftcore/dss-go/dss/internal/asn1ber"

// DigestCalculatorProvider hands back the message-digest value to use for a given digest
// AlgorithmIdentifier. Port of the shared contract of
// org.bouncycastle.operator.DigestCalculatorProvider#get(AlgorithmIdentifier).getDigest():
// BouncyCastle's version streams bytes through a DigestCalculator and reads getDigest() back
// afterward, but neither DSS implementation depends on which bytes were streamed, so Go asks
// for the digest directly.
type DigestCalculatorProvider interface {
	// Digest returns the message-digest value for the given digest algorithm identifier.
	Digest(digestAlgorithmIdentifier *asn1ber.AlgorithmIdentifier) ([]byte, error)
}
