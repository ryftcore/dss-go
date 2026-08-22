// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/CAdESSignatureParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
)

// SignatureParameters defines SignatureParameters to deal with CAdES signature
// creation/extension.
type SignatureParameters struct {
	document.AbstractSignatureParameters[*TimestampParameters]

	// en319122 defines if the signature shall be created according to ETSI EN 319 122.
	en319122 bool

	// contentHintsType is the Content Hints type.
	contentHintsType string

	// contentHintsDescription is the Content Hints description.
	contentHintsDescription string

	// contentIdentifierPrefix is the content identifier prefix.
	contentIdentifierPrefix string

	// contentIdentifierSuffix is the content identifier suffix.
	contentIdentifierSuffix string

	// parallelSignature defines if a parallel signature should be created when a CMS signature is
	// provided as a signed document.
	parallelSignature bool
}

// NewCAdESSignatureParameters instantiates object with null values. Port of the default
// constructor.
func NewSignatureParameters() *SignatureParameters {
	return &SignatureParameters{
		AbstractSignatureParameters: document.NewAbstractSignatureParameters[*TimestampParameters](),
		en319122:                    true,
		parallelSignature:           true,
	}
}

// IsEn319122 returns if the signature shall be created according to ETSI EN 319 122, otherwise as
// an old format. Port of #isEn319122.
func (p *SignatureParameters) IsEn319122() bool {
	return p.en319122
}

// SetEn319122 sets if the signature shall be created according to ETSI EN 319 122, otherwise will
// be created with respect to the old standard.
//
// Default: true. Port of #setEn319122.
func (p *SignatureParameters) SetEn319122(en319122 bool) {
	p.en319122 = en319122
}

// ContentHintsType gets content hints type. Port of #getContentHintsType.
func (p *SignatureParameters) ContentHintsType() string {
	return p.contentHintsType
}

// SetContentHintsType sets content hints type. Port of #setContentHintsType.
func (p *SignatureParameters) SetContentHintsType(contentHintsType string) {
	p.contentHintsType = contentHintsType
}

// ContentHintsDescription gets content hints description. Port of #getContentHintsDescription.
func (p *SignatureParameters) ContentHintsDescription() string {
	return p.contentHintsDescription
}

// SetContentHintsDescription sets content hints description. Port of
// #setContentHintsDescription.
func (p *SignatureParameters) SetContentHintsDescription(contentHintsDescription string) {
	p.contentHintsDescription = contentHintsDescription
}

// ContentIdentifierSuffix returns the content identifier suffix.
//
// ETSI TS 101 733 V2.2.1 (2013-04)
//
// 5.10.2 content-identifier Attribute
// The content-identifier attribute provides an identifier for the signed content, for use when a
// reference may be later required to that content; for example, in the content-reference
// attribute in other signed data sent later. The content-identifier shall be a signed attribute.
// content-identifier attribute type values for the ES have an ASN.1 type ContentIdentifier, as
// defined in ESS (RFC 2634).
//
// The minimal content-identifier attribute should contain a concatenation of user-specific
// identification information (such as a user name or public keying material identification
// information), a GeneralizedTime string, and a random number.
//
// Port of #getContentIdentifierSuffix.
func (p *SignatureParameters) ContentIdentifierSuffix() string {
	return p.contentIdentifierSuffix
}

// SetContentIdentifierSuffix sets content identifier suffix.
//
// NOTE: THIS VALUE WILL BE SET AUTOMATICALLY IF LEFT BLANK
//
// Port of #setContentIdentifierSuffix.
func (p *SignatureParameters) SetContentIdentifierSuffix(contentIdentifierSuffix string) {
	p.contentIdentifierSuffix = contentIdentifierSuffix
}

// ContentIdentifierPrefix returns the content identifier prefix. See ContentIdentifierSuffix's
// doc. Port of #getContentIdentifierPrefix.
func (p *SignatureParameters) ContentIdentifierPrefix() string {
	return p.contentIdentifierPrefix
}

// SetContentIdentifierPrefix sets content identifier prefix. Port of #setContentIdentifierPrefix.
func (p *SignatureParameters) SetContentIdentifierPrefix(contentIdentifierPrefix string) {
	p.contentIdentifierPrefix = contentIdentifierPrefix
}

// IsParallelSignature returns whether a parallel signature should be created when an original
// document is represented by a CMSSignedData (i.e. another CMS signature document). Port of
// #isParallelSignature.
func (p *SignatureParameters) IsParallelSignature() bool {
	return p.parallelSignature
}

// SetParallelSignature sets whether a parallel signature should be created in case the provided
// documentToSign is represented by a CMSSignedData (i.e. another CMS signature document). When
// enabled, the created produced document will be build based on the original CMSSignedData
// content, and the new signature will be included inside the CMSSignedData as a parallel
// signature, next to existing signatures. When disabled, the new signature will cover the
// binaries of the original CMSSignedData, as in case of a not-signed original document. This
// allows creation of nested CMS signatures. When a non CMSSignedData document is provided, the
// parameter does not impact the processing.
//
// Default: TRUE (creates a parallel signature, when applicable). Port of
// #setParallelSignature.
func (p *SignatureParameters) SetParallelSignature(parallelSignature bool) {
	p.parallelSignature = parallelSignature
}

// GetContentTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getContentTimestampParameters.
func (p *SignatureParameters) GetContentTimestampParameters() *TimestampParameters {
	if p.ContentTimestampParameters == nil {
		p.ContentTimestampParameters = NewTimestampParameters()
	}
	return p.ContentTimestampParameters
}

// GetSignatureTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getSignatureTimestampParameters.
func (p *SignatureParameters) GetSignatureTimestampParameters() *TimestampParameters {
	if p.SignatureTimestampParameters == nil {
		p.SignatureTimestampParameters = NewTimestampParameters()
	}
	return p.SignatureTimestampParameters
}

// GetArchiveTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getArchiveTimestampParameters.
func (p *SignatureParameters) GetArchiveTimestampParameters() *TimestampParameters {
	if p.ArchiveTimestampParameters == nil {
		p.ArchiveTimestampParameters = NewTimestampParameters()
	}
	return p.ArchiveTimestampParameters
}

// String ports #toString.
func (p *SignatureParameters) String() string {
	return fmt.Sprintf("CAdESSignatureParameters [en319122=%v, contentHintsType='%v', contentHintsDescription='%v', contentIdentifierPrefix='%v', contentIdentifierSuffix='%v', parallelSignature=%v] %s",
		p.en319122, p.contentHintsType, p.contentHintsDescription, p.contentIdentifierPrefix, p.contentIdentifierSuffix, p.parallelSignature,
		p.AbstractSignatureParameters.String())
}

// Equals ports #equals.
func (p *SignatureParameters) Equals(other *SignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.AbstractSignatureParameters.Equals(&other.AbstractSignatureParameters) {
		return false
	}
	return p.en319122 == other.en319122 &&
		p.parallelSignature == other.parallelSignature &&
		p.contentHintsType == other.contentHintsType &&
		p.contentHintsDescription == other.contentHintsDescription &&
		p.contentIdentifierPrefix == other.contentIdentifierPrefix &&
		p.contentIdentifierSuffix == other.contentIdentifierSuffix
}
