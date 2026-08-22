// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/ReferenceIdProvider.java (DSS 6.5.RC1).
package xades

import (
	"strconv"

	"github.com/ryftcore/dss-go/dss/utils"
)

// ReferenceIdProvider is used to generate a deterministic reference identifier.
type ReferenceIdProvider struct {
	// signatureParameters are the signature parameters used to create the signature.
	signatureParameters *SignatureParameters

	// referenceIdPrefix is the id-prefix for the ds:Reference element. Java's field
	// initializer is "r" (its javadoc says "r-", which the code contradicts).
	referenceIdPrefix string

	// index is the internal reference id counter.
	index int
}

// NewReferenceIdProvider ports the default constructor, including the referenceIdPrefix field
// initializer.
func NewReferenceIdProvider() *ReferenceIdProvider {
	return &ReferenceIdProvider{referenceIdPrefix: "r"}
}

// SetSignatureParameters sets the signature parameters used to build a deterministic
// identifier. Ports setSignatureParameters(XAdESSignatureParameters).
func (p *ReferenceIdProvider) SetSignatureParameters(signatureParameters *SignatureParameters) {
	p.signatureParameters = signatureParameters
}

// SetReferenceIdPrefix sets the reference id prefix to be used on reference creation. Ports
// setReferenceIdPrefix(String), whose IllegalArgumentException becomes a panic with the same
// message: the argument is a caller-side constant at every call site (ManifestBuilder), so this
// is an unchecked programming error rather than an input-validation failure, and the setter
// keeps its bare return.
func (p *ReferenceIdProvider) SetReferenceIdPrefix(referenceIdPrefix string) {
	if utils.IsStringBlank(referenceIdPrefix) {
		panic("The reference id prefix cannot be blank!")
	}
	p.referenceIdPrefix = referenceIdPrefix
}

// ReferenceId returns the following signature reference identifier. Ports getReferenceId().
func (p *ReferenceIdProvider) ReferenceId() string {
	p.increaseIndex()

	referenceId := p.referenceIdPrefix
	referenceId += "-"
	if p.signatureParameters != nil {
		referenceId += p.signatureParameters.GetDeterministicId()
		referenceId += "-"
	}
	referenceId += strconv.Itoa(p.index)
	return referenceId
}

// increaseIndex ports the private increaseIndex().
func (p *ReferenceIdProvider) increaseIndex() {
	p.index++
}
