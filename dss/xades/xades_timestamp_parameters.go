// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/XAdESTimestampParameters.java (DSS 6.5.RC1).
package xades

import (
	"fmt"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
	xmlutils "github.com/utain/esig/dss/xml/utils"
)

// XAdESTimestampParameters holds the parameters for a XAdES timestamp creation.
type XAdESTimestampParameters struct {
	model.TimestampParameters

	// canonicalizationMethod is the canonicalization method to use for the message-imprint.
	canonicalizationMethod string
}

var _ model.SerializableTimestampParameters = (*XAdESTimestampParameters)(nil)

// NewXAdESTimestampParameters instantiates the object with null values, including the
// canonicalizationMethod field initializer. Port of the empty constructor.
func NewXAdESTimestampParameters() *XAdESTimestampParameters {
	return &XAdESTimestampParameters{
		TimestampParameters:    model.NewTimestampParameters(),
		canonicalizationMethod: xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
	}
}

// NewXAdESTimestampParametersWithDigestAlgorithm is a constructor with a digest algorithm to
// use for message-imprint digest calculation, including the canonicalizationMethod field
// initializer. Port of XAdESTimestampParameters(DigestAlgorithm).
func NewXAdESTimestampParametersWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *XAdESTimestampParameters {
	return &XAdESTimestampParameters{
		TimestampParameters:    model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm),
		canonicalizationMethod: xmlutils.XMLCanonicalizerDefaultDSSC14NMethod,
	}
}

// NewXAdESTimestampParametersWithCanonicalization is the default constructor with a digest
// algorithm and canonicalization method to use for the message-imprint. Port of
// XAdESTimestampParameters(DigestAlgorithm, String).
func NewXAdESTimestampParametersWithCanonicalization(digestAlgorithm enumerations.DigestAlgorithm,
	canonicalizationMethod string) *XAdESTimestampParameters {
	return &XAdESTimestampParameters{
		TimestampParameters:    model.NewTimestampParametersWithDigestAlgorithm(digestAlgorithm),
		canonicalizationMethod: canonicalizationMethod,
	}
}

// CanonicalizationMethod gets the canonicalization method. Ports getCanonicalizationMethod().
func (p *XAdESTimestampParameters) CanonicalizationMethod() string {
	return p.canonicalizationMethod
}

// SetCanonicalizationMethod sets the canonicalization method. Panics with the Java message when
// canonicalizationMethod is empty (IllegalArgumentException upstream). Ports
// setCanonicalizationMethod(String).
func (p *XAdESTimestampParameters) SetCanonicalizationMethod(canonicalizationMethod string) {
	if utils.IsStringEmpty(canonicalizationMethod) {
		panic("Canonicalization cannot be empty! See EN 319 132-1: 4.5 Managing canonicalization of XML nodesets.")
	}
	p.canonicalizationMethod = canonicalizationMethod
}

// String ports toString().
func (p *XAdESTimestampParameters) String() string {
	return fmt.Sprintf("XAdESTimestampParameters [canonicalizationMethod='%s'] %s",
		p.canonicalizationMethod, p.TimestampParameters.String())
}

// Equals ports equals(Object).
func (p *XAdESTimestampParameters) Equals(other *XAdESTimestampParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.TimestampParameters.Equals(&other.TimestampParameters) {
		return false
	}
	return p.canonicalizationMethod == other.canonicalizationMethod
}
