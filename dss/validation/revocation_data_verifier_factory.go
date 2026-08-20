// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/RevocationDataVerifierFactory.java
// (DSS 6.5.RC1).
//
// Deviations:
//
//   - Java gathers the acceptable digest algorithms and the revocation-skip
//     certificate extensions/policies in HashSets, whose iteration order is
//     unspecified; the ported RevocationDataVerifier takes slices, so this port
//     keeps the first-insertion order instead. The values are only ever looked
//     up by membership, so nothing observable depends on the order.
//
//   - Java's `Long` freshness constraints (null = unset) become *int64.

package validation

import (
	"time"

	"github.com/utain/esig/dss/enumerations"
	modelpolicy "github.com/utain/esig/dss/model/policy"
	spivalidation "github.com/utain/esig/dss/spi/validation"
	validationpolicy "github.com/utain/esig/dss/validation/policy"
)

// RevocationDataVerifierFactory loads a RevocationDataVerifier from a provided
// eu.europa.esig.dss.model.policy.ValidationPolicy. Port of
// RevocationDataVerifierFactory.
type RevocationDataVerifierFactory struct {
	// validationPolicy is the validation policy to load the
	// RevocationDataVerifier from. Port of the private final validationPolicy
	// field.
	validationPolicy modelpolicy.ValidationPolicy

	// validationTime is the used validation time. Port of the private
	// validationTime field.
	validationTime time.Time
}

// NewRevocationDataVerifierFactory is the default constructor. Port of
// RevocationDataVerifierFactory(ValidationPolicy).
func NewRevocationDataVerifierFactory(validationPolicy modelpolicy.ValidationPolicy) *RevocationDataVerifierFactory {
	return &RevocationDataVerifierFactory{validationPolicy: validationPolicy}
}

// ValidationTime gets the validation time, instantiating the value to the
// current time when not provided explicitly. Port of the protected
// getValidationTime().
func (f *RevocationDataVerifierFactory) ValidationTime() time.Time {
	if f.validationTime.IsZero() {
		f.validationTime = time.Now()
	}
	return f.validationTime
}

// SetValidationTime sets the used validation time. Port of
// setValidationTime(Date).
func (f *RevocationDataVerifierFactory) SetValidationTime(validationTime time.Time) *RevocationDataVerifierFactory {
	f.validationTime = validationTime
	return f
}

// Create creates the RevocationDataVerifier. Port of create().
func (f *RevocationDataVerifierFactory) Create() *spivalidation.RevocationDataVerifier {
	revocationDataVerifier := spivalidation.NewEmptyRevocationDataVerifier()
	f.instantiateCryptographicSuite(revocationDataVerifier)
	f.instantiateRevocationSkipConstraints(revocationDataVerifier)
	f.instantiateRevocationFreshnessConstraints(revocationDataVerifier)
	f.instantiateAcceptRevocationIssuersWithoutRevocationConstraint(revocationDataVerifier)
	return revocationDataVerifier
}

// instantiateCryptographicSuite is the port of the private
// instantiateCryptographicSuite(RevocationDataVerifier).
func (f *RevocationDataVerifierFactory) instantiateCryptographicSuite(
	revocationDataVerifier *spivalidation.RevocationDataVerifier) {
	var acceptableDigestAlgorithms []enumerations.DigestAlgorithm
	var acceptableSignatureAlgorithms []*modelpolicy.SignatureAlgorithmWithMinKeySize

	cryptographicSuite := f.revocationCryptographicSuite(f.validationPolicy)
	if cryptographicSuite != nil && enumerations.Level_FAIL == cryptographicSuite.Level() {
		currentTime := f.ValidationTime()
		acceptableDigestAlgorithms = validationpolicy.GetReliableDigestAlgorithmsAtTime(cryptographicSuite, currentTime)
		acceptableSignatureAlgorithms = validationpolicy.GetReliableSignatureAlgorithmsWithMinimalKeyLengthAtTime(
			cryptographicSuite, currentTime)
	} else {
		// LOG.info("No enforced cryptographic constraints have been found in the
		// provided validation policy. Accept all cryptographic algorithms.")
		acceptableDigestAlgorithms = append(acceptableDigestAlgorithms, enumerations.DigestAlgorithmValues()...)
		for _, signatureAlgorithm := range enumerations.SignatureAlgorithmValues() {
			acceptableSignatureAlgorithms = append(acceptableSignatureAlgorithms,
				modelpolicy.NewSignatureAlgorithmWithMinKeySize(signatureAlgorithm, 0))
		}
	}
	revocationDataVerifier.SetAcceptableDigestAlgorithms(acceptableDigestAlgorithms)
	revocationDataVerifier.SetAcceptableSignatureAlgorithmKeyLength(
		toSignatureAlgorithmWithKeySizesMap(acceptableSignatureAlgorithms))
}

// revocationCryptographicSuite is the port of the private
// getRevocationCryptographicSuite(ValidationPolicy).
func (f *RevocationDataVerifierFactory) revocationCryptographicSuite(
	validationPolicy modelpolicy.ValidationPolicy) modelpolicy.CryptographicSuite {
	return validationPolicy.SignatureCryptographicConstraint(enumerations.Context_REVOCATION)
}

// toSignatureAlgorithmWithKeySizesMap is the port of the private
// toSignatureAlgorithmWithKeySizesMap(Collection<SignatureAlgorithmWithMinKeySize>).
func toSignatureAlgorithmWithKeySizesMap(
	signatureAlgorithms []*modelpolicy.SignatureAlgorithmWithMinKeySize) map[enumerations.SignatureAlgorithm]int {
	signatureAlgorithmsMap := make(map[enumerations.SignatureAlgorithm]int)
	for _, encryptionAlgorithmWithMinKeySize := range signatureAlgorithms {
		signatureAlgorithm := encryptionAlgorithmWithMinKeySize.SignatureAlgorithm()
		minKeySize, present := signatureAlgorithmsMap[signatureAlgorithm]
		keySize := encryptionAlgorithmWithMinKeySize.MinKeySize()
		if !present || minKeySize > keySize {
			signatureAlgorithmsMap[signatureAlgorithm] = keySize
		}
	}
	return signatureAlgorithmsMap
}

// instantiateRevocationSkipConstraints is the port of the private
// instantiateRevocationSkipConstraints(RevocationDataVerifier).
func (f *RevocationDataVerifierFactory) instantiateRevocationSkipConstraints(
	revocationDataVerifier *spivalidation.RevocationDataVerifier) {
	certificateExtensions := newOrderedStringSet()
	certificatePolicies := newOrderedStringSet()

	f.populateRevocationSkipFromBasicSignatureConstraints(certificateExtensions, certificatePolicies, enumerations.Context_SIGNATURE)
	f.populateRevocationSkipFromBasicSignatureConstraints(certificateExtensions, certificatePolicies, enumerations.Context_COUNTER_SIGNATURE)
	f.populateRevocationSkipFromBasicSignatureConstraints(certificateExtensions, certificatePolicies, enumerations.Context_KEY_BINDING_SIGNATURE)
	f.populateRevocationSkipFromBasicSignatureConstraints(certificateExtensions, certificatePolicies, enumerations.Context_REVOCATION)
	f.populateRevocationSkipFromBasicSignatureConstraints(certificateExtensions, certificatePolicies, enumerations.Context_TIMESTAMP)

	revocationDataVerifier.SetRevocationSkipCertificateExtensions(certificateExtensions.values)
	revocationDataVerifier.SetRevocationSkipCertificatePolicies(certificatePolicies.values)
}

// orderedStringSet stands in for Java's HashSet<String>: it keeps membership
// semantics while giving the ported (slice-taking) RevocationDataVerifier a
// deterministic order. See the file header.
type orderedStringSet struct {
	seen   map[string]struct{}
	values []string
}

func newOrderedStringSet() *orderedStringSet {
	return &orderedStringSet{seen: make(map[string]struct{})}
}

func (s *orderedStringSet) addAll(values []string) {
	for _, value := range values {
		if _, exists := s.seen[value]; exists {
			continue
		}
		s.seen[value] = struct{}{}
		s.values = append(s.values, value)
	}
}

// populateRevocationSkipFromBasicSignatureConstraints is the port of the
// private populateRevocationSkipFromBasicSignatureConstraints(Set, Set, Context).
func (f *RevocationDataVerifierFactory) populateRevocationSkipFromBasicSignatureConstraints(
	certificateExtensions, certificatePolicies *orderedStringSet, context enumerations.Context) {
	f.populateRevocationSkipFromCertificateConstraints(certificateExtensions, certificatePolicies,
		context, enumerations.SubContext_SIGNING_CERT)
	f.populateRevocationSkipFromCertificateConstraints(certificateExtensions, certificatePolicies,
		context, enumerations.SubContext_CA_CERTIFICATE)
}

// populateRevocationSkipFromCertificateConstraints is the port of the private
// populateRevocationSkipFromCertificateConstraints(Set, Set, Context, SubContext).
func (f *RevocationDataVerifierFactory) populateRevocationSkipFromCertificateConstraints(
	certificateExtensions, certificatePolicies *orderedStringSet,
	context enumerations.Context, subContext enumerations.SubContext) {
	revocationDataSkipConstraint := f.validationPolicy.RevocationDataSkipConstraint(context, subContext)
	if revocationDataSkipConstraint == nil {
		return
	}

	certificateExtensionsConstraint := revocationDataSkipConstraint.CertificateExtensions()
	if certificateExtensionsConstraint != nil {
		certificateExtensions.addAll(certificateExtensionsConstraint.Values())
	}

	certificatePoliciesConstraint := revocationDataSkipConstraint.CertificatePolicies()
	if certificatePoliciesConstraint != nil {
		certificatePolicies.addAll(certificatePoliciesConstraint.Values())
	}
}

// instantiateRevocationFreshnessConstraints is the port of the private
// instantiateRevocationFreshnessConstraints(RevocationDataVerifier).
func (f *RevocationDataVerifierFactory) instantiateRevocationFreshnessConstraints(
	revocationDataVerifier *spivalidation.RevocationDataVerifier) {
	revocationDataVerifier.SetSignatureMaximumRevocationFreshness(f.signatureRevocationFreshnessConstraint())
	revocationDataVerifier.SetTimestampMaximumRevocationFreshness(f.revocationFreshnessConstraint(enumerations.Context_TIMESTAMP))
	revocationDataVerifier.SetRevocationMaximumRevocationFreshness(f.revocationFreshnessConstraint(enumerations.Context_REVOCATION))

	revocationFreshnessNextUpdateConstraint := f.revocationFreshnessNextUpdateConstraintPresent(enumerations.Context_SIGNATURE)
	revocationFreshnessNextUpdateConstraint = revocationFreshnessNextUpdateConstraint || f.revocationFreshnessNextUpdateConstraintPresent(enumerations.Context_COUNTER_SIGNATURE)
	revocationFreshnessNextUpdateConstraint = revocationFreshnessNextUpdateConstraint || f.revocationFreshnessNextUpdateConstraintPresent(enumerations.Context_KEY_BINDING_SIGNATURE)
	revocationFreshnessNextUpdateConstraint = revocationFreshnessNextUpdateConstraint || f.revocationFreshnessNextUpdateConstraintPresent(enumerations.Context_TIMESTAMP)
	revocationFreshnessNextUpdateConstraint = revocationFreshnessNextUpdateConstraint || f.revocationFreshnessNextUpdateConstraintPresent(enumerations.Context_REVOCATION)

	revocationDataVerifier.SetCheckRevocationFreshnessNextUpdate(revocationFreshnessNextUpdateConstraint)
}

// signatureRevocationFreshnessConstraint is the port of the private
// getSignatureRevocationFreshnessConstraint().
func (f *RevocationDataVerifierFactory) signatureRevocationFreshnessConstraint() *int64 {
	maximumRevocationFreshness := f.revocationFreshnessConstraint(enumerations.Context_SIGNATURE)

	counterSignatureRevocationFreshnessConstraint := f.revocationFreshnessConstraint(enumerations.Context_COUNTER_SIGNATURE)
	if maximumRevocationFreshness == nil || (counterSignatureRevocationFreshnessConstraint != nil &&
		*counterSignatureRevocationFreshnessConstraint < *maximumRevocationFreshness) {
		maximumRevocationFreshness = counterSignatureRevocationFreshnessConstraint
	}
	keyBindingSignatureRevocationFreshnessConstraint := f.revocationFreshnessConstraint(enumerations.Context_KEY_BINDING_SIGNATURE)
	if maximumRevocationFreshness == nil || (keyBindingSignatureRevocationFreshnessConstraint != nil &&
		*keyBindingSignatureRevocationFreshnessConstraint < *maximumRevocationFreshness) {
		maximumRevocationFreshness = keyBindingSignatureRevocationFreshnessConstraint
	}

	return maximumRevocationFreshness
}

// revocationFreshnessConstraint is the port of the private
// getRevocationFreshnessConstraint(Context).
func (f *RevocationDataVerifierFactory) revocationFreshnessConstraint(context enumerations.Context) *int64 {
	maximumRevocationFreshness := f.revocationFreshnessConstraintValue(context, enumerations.SubContext_SIGNING_CERT)
	caCertRevocationFreshness := f.revocationFreshnessConstraintValue(context, enumerations.SubContext_CA_CERTIFICATE)
	if maximumRevocationFreshness == nil || (caCertRevocationFreshness != nil &&
		*caCertRevocationFreshness < *maximumRevocationFreshness) {
		maximumRevocationFreshness = caCertRevocationFreshness
	}
	return maximumRevocationFreshness
}

// revocationFreshnessConstraintValue is the port of the private
// getRevocationFreshnessConstraintValue(Context, SubContext).
func (f *RevocationDataVerifierFactory) revocationFreshnessConstraintValue(context enumerations.Context,
	subContext enumerations.SubContext) *int64 {
	revocationFreshness := f.validationPolicy.RevocationFreshnessConstraint(context, subContext)
	if revocationFreshness != nil {
		duration := revocationFreshness.Duration()
		return &duration
	}
	return nil
}

// revocationFreshnessNextUpdateConstraintPresent is the port of the private
// getRevocationFreshnessNextUpdateConstraintPresent(Context).
func (f *RevocationDataVerifierFactory) revocationFreshnessNextUpdateConstraintPresent(context enumerations.Context) bool {
	revocationFreshnessNextUpdateConstraint := f.validationPolicy.RevocationFreshnessNextUpdateConstraint(
		context, enumerations.SubContext_SIGNING_CERT)
	if revocationFreshnessNextUpdateConstraint != nil {
		return true
	}
	revocationFreshnessNextUpdateConstraint = f.validationPolicy.RevocationFreshnessNextUpdateConstraint(
		context, enumerations.SubContext_CA_CERTIFICATE)
	return revocationFreshnessNextUpdateConstraint != nil
}

// instantiateAcceptRevocationIssuersWithoutRevocationConstraint is the port of
// the private instantiateAcceptRevocationIssuersWithoutRevocationConstraint(RevocationDataVerifier).
func (f *RevocationDataVerifierFactory) instantiateAcceptRevocationIssuersWithoutRevocationConstraint(
	revocationDataVerifier *spivalidation.RevocationDataVerifier) {
	revocationDataAvailableConstraint := f.revocationDataAvailablePresent(enumerations.Context_REVOCATION)
	revocationDataVerifier.SetAcceptRevocationCertificatesWithoutRevocation(!revocationDataAvailableConstraint)

	revocationDataAvailableConstraint = f.revocationDataAvailablePresent(enumerations.Context_TIMESTAMP)
	revocationDataVerifier.SetAcceptTimestampCertificatesWithoutRevocation(!revocationDataAvailableConstraint)
}

// revocationDataAvailablePresent is the port of the private
// getRevocationDataAvailablePresent(Context).
func (f *RevocationDataVerifierFactory) revocationDataAvailablePresent(context enumerations.Context) bool {
	signingCertificateRule := f.validationPolicy.RevocationDataAvailableConstraint(context, enumerations.SubContext_SIGNING_CERT)
	if signingCertificateRule != nil && enumerations.Level_FAIL == signingCertificateRule.Level() {
		return true
	}
	caCertificateRule := f.validationPolicy.RevocationDataAvailableConstraint(context, enumerations.SubContext_CA_CERTIFICATE)
	return caCertificateRule != nil && enumerations.Level_FAIL == caCertificateRule.Level()
}
