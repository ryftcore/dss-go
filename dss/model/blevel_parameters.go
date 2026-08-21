// Ported from dss-model/.../BLevelParameters.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// BLevelParameters is used to define common b-level parameters for a
// signature creation.
//
// Policy and SignerLocation are outside this manifest; assumed to already
// exist in this package.
type BLevelParameters struct {
	// trustAnchorBPPolicy indicates if the Baseline profile's trust anchor
	// policy shall be followed: ETSI TS 103 171 V2.1.1 (2012-03) 6.2.1
	// Placement of the signing certificate. When -B level is constructed
	// the trust anchor is not included, when -LT level is constructed the
	// trust anchor is included. NOTE: when trust anchor baseline profile
	// policy is defined only the certificates previous to the trust anchor
	// are included when -B level is constructed.
	trustAnchorBPPolicy bool

	// signingDate is the claimed signing time.
	signingDate *time.Time

	// claimedSignerRoles are the claimed signer roles.
	claimedSignerRoles []string

	// signedAssertions are the signed assertions.
	signedAssertions []string

	// signaturePolicy is the Signature Policy Identifier.
	signaturePolicy *Policy

	// commitmentTypeIndications is the list of commitment type
	// indications.
	commitmentTypeIndications []enumerations.CommitmentType

	// signerLocation is the SignerLocation.
	signerLocation *SignerLocation
}

// NewBLevelParameters instantiates the object with default values. Ports
// the default constructor.
func NewBLevelParameters() *BLevelParameters {
	now := time.Now()
	return &BLevelParameters{
		trustAnchorBPPolicy: true,
		signingDate:         &now,
	}
}

// IsTrustAnchorBPPolicy checks if the trust anchor policy shall be used
// when creating -B and -LT levels.
func (b *BLevelParameters) IsTrustAnchorBPPolicy() bool { return b.trustAnchorBPPolicy }

// SetTrustAnchorBPPolicy allows setting the trust anchor policy to use
// when creating -B and -LT levels.
//
// NOTE: when trust anchor baseline profile policy is defined only the
// certificates previous to the trust anchor are included when building
// -B level.
func (b *BLevelParameters) SetTrustAnchorBPPolicy(trustAnchorBPPolicy bool) {
	b.trustAnchorBPPolicy = trustAnchorBPPolicy
}

// SignaturePolicy gets the signature policy to use during the signature
// creation process.
func (b *BLevelParameters) SignaturePolicy() *Policy { return b.signaturePolicy }

// SetSignaturePolicy indicates the signature policy to use.
func (b *BLevelParameters) SetSignaturePolicy(signaturePolicy *Policy) {
	b.signaturePolicy = signaturePolicy
}

// SigningDate gets the signing date.
func (b *BLevelParameters) SigningDate() *time.Time { return b.signingDate }

// SetSigningDate sets the signing date. Panics if signingDate is nil (Java
// Objects.requireNonNull("SigningDate cannot be null!")).
func (b *BLevelParameters) SetSigningDate(signingDate *time.Time) {
	if signingDate == nil {
		panic("SigningDate cannot be null!")
	}
	b.signingDate = signingDate
}

// ClaimedSignerRoles gets the list of claimed roles.
func (b *BLevelParameters) ClaimedSignerRoles() []string { return b.claimedSignerRoles }

// SetClaimedSignerRoles sets a list of claimed signer roles.
func (b *BLevelParameters) SetClaimedSignerRoles(claimedSignerRoles []string) {
	b.claimedSignerRoles = claimedSignerRoles
}

// SignedAssertions gets the signed assertions.
func (b *BLevelParameters) SignedAssertions() []string { return b.signedAssertions }

// SetSignedAssertions sets signed assertions.
func (b *BLevelParameters) SetSignedAssertions(signedAssertions []string) {
	b.signedAssertions = signedAssertions
}

// CommitmentTypeIndications gets the commitment type indications.
func (b *BLevelParameters) CommitmentTypeIndications() []enumerations.CommitmentType {
	return b.commitmentTypeIndications
}

// SetCommitmentTypeIndications sets the commitment type indications
// (predefined values are available as enumerations.CommitmentType
// implementations).
func (b *BLevelParameters) SetCommitmentTypeIndications(commitmentTypeIndications []enumerations.CommitmentType) {
	b.commitmentTypeIndications = commitmentTypeIndications
}

// SignerLocation gets the signer location.
func (b *BLevelParameters) SignerLocation() *SignerLocation { return b.signerLocation }

// SetSignerLocation sets the signer location.
func (b *BLevelParameters) SetSignerLocation(signerLocation *SignerLocation) {
	b.signerLocation = signerLocation
}

// Equals ports BLevelParameters#equals.
func (b *BLevelParameters) Equals(other *BLevelParameters) bool {
	if b == other {
		return true
	}
	if other == nil {
		return false
	}
	if !reflect.DeepEqual(b.claimedSignerRoles, other.claimedSignerRoles) {
		return false
	}
	if !reflect.DeepEqual(b.signedAssertions, other.signedAssertions) {
		return false
	}
	if !reflect.DeepEqual(b.commitmentTypeIndications, other.commitmentTypeIndications) {
		return false
	}
	if !bLevelParametersPolicyEqual(b.signaturePolicy, other.signaturePolicy) {
		return false
	}
	if !bLevelParametersLocationEqual(b.signerLocation, other.signerLocation) {
		return false
	}
	if !bLevelParametersDateEqual(b.signingDate, other.signingDate) {
		return false
	}
	return b.trustAnchorBPPolicy == other.trustAnchorBPPolicy
}

// bLevelParametersPolicyEqual compares two possibly-nil *Policy values.
func bLevelParametersPolicyEqual(a, b *Policy) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// bLevelParametersLocationEqual compares two possibly-nil *SignerLocation
// values.
func bLevelParametersLocationEqual(a, b *SignerLocation) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// bLevelParametersDateEqual compares two possibly-nil *time.Time values.
func bLevelParametersDateEqual(a, b *time.Time) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equal(*b)
}

// String ports BLevelParameters#toString.
func (b *BLevelParameters) String() string {
	return fmt.Sprintf("BLevelParameters [trustAnchorBPPolicy=%v, signingDate=%v, claimedSignerRoles=%v, signedAssertions=%v, signaturePolicy=%v, commitmentTypeIndication=%v, signerLocation=%v]",
		b.trustAnchorBPPolicy, b.signingDate, b.claimedSignerRoles, b.signedAssertions, b.signaturePolicy, b.commitmentTypeIndications, b.signerLocation)
}
