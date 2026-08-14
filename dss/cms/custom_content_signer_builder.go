// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/operator/CustomContentSignerBuilder.java
// (DSS 6.5.RC1).
package cms

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
)

// CustomContentSignerBuilder is used to create an instance of CustomContentSigner.
// Port of the CustomContentSignerBuilder class.
type CustomContentSignerBuilder struct{}

// NewCustomContentSignerBuilder is the default constructor. Port of the no-arg constructor.
func NewCustomContentSignerBuilder() *CustomContentSignerBuilder {
	return &CustomContentSignerBuilder{}
}

// Build builds a CustomContentSigner for the CMS signature creation. This method creates a
// CustomContentSigner with an absent SignatureValue. Method is normally used for message-digest
// computation. Port of #build(SignatureAlgorithm).
//
// Panics with the Java message when signatureAlgorithm is empty (Objects.requireNonNull).
func (b *CustomContentSignerBuilder) Build(signatureAlgorithm enumerations.SignatureAlgorithm) (*CustomContentSigner, error) {
	if signatureAlgorithm == "" {
		panic("SignatureAlgorithm cannot be null!")
	}
	return NewCustomContentSigner(signatureAlgorithm.JCEID())
}

// BuildWithSignatureValue builds a CustomContentSigner for the CMS signature creation using the
// given SignatureValue. Port of #build(SignatureAlgorithm, SignatureValue).
//
// Panics with the Java message when signatureAlgorithm is empty or signatureValue is nil
// (Objects.requireNonNull); the IllegalArgumentException raised for a SignatureValue computed
// with a different SignatureAlgorithm is a returned error, since the whole creation path
// already has an error channel.
func (b *CustomContentSignerBuilder) BuildWithSignatureValue(signatureAlgorithm enumerations.SignatureAlgorithm,
	signatureValue *model.SignatureValue) (*CustomContentSigner, error) {
	if signatureAlgorithm == "" {
		panic("SignatureAlgorithm cannot be null!")
	}
	if signatureValue == nil {
		panic("signatureValue cannot be null!")
	}
	if signatureAlgorithm != signatureValue.Algorithm() {
		return nil, fmt.Errorf("The defined SignatureAlgorithm '%s' does not match the SignatureAlgorithm '%s' "+
			"used on SignatureValue computation!", signatureAlgorithm, signatureValue.Algorithm())
	}
	return NewCustomContentSignerWithSignature(signatureAlgorithm.JCEID(), signatureValue.Value())
}
