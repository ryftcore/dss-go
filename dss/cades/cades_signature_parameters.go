// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/CAdESSignatureParameters.java (DSS 6.5.RC1).
package cades

import (
	"fmt"

	"github.com/utain/esig/dss/document"
)

// CAdESSignatureParameters defines SignatureParameters to deal with CAdES signature
// creation/extension.
type CAdESSignatureParameters struct {
	document.AbstractSignatureParameters[*CAdESTimestampParameters]

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
func NewCAdESSignatureParameters() *CAdESSignatureParameters {
	return &CAdESSignatureParameters{
		AbstractSignatureParameters: document.NewAbstractSignatureParameters[*CAdESTimestampParameters](),
		en319122:                    true,
		parallelSignature:           true,
	}
}

// IsEn319122 returns if the signature shall be created according to ETSI EN 319 122, otherwise as
// an old format. Port of #isEn319122.
func (p *CAdESSignatureParameters) IsEn319122() bool {
	return p.en319122
}

// SetEn319122 sets if the signature shall be created according to ETSI EN 319 122, otherwise will
// be created with respect to the old standard.
//
// Default: true. Port of #setEn319122.
func (p *CAdESSignatureParameters) SetEn319122(en319122 bool) {
	p.en319122 = en319122
}

// ContentHintsType gets content hints type. Port of #getContentHintsType.
func (p *CAdESSignatureParameters) ContentHintsType() string {
	return p.contentHintsType
}

// SetContentHintsType sets content hints type. Port of #setContentHintsType.
func (p *CAdESSignatureParameters) SetContentHintsType(contentHintsType string) {
	p.contentHintsType = contentHintsType
}

// ContentHintsDescription gets content hints description. Port of #getContentHintsDescription.
func (p *CAdESSignatureParameters) ContentHintsDescription() string {
	return p.contentHintsDescription
}

// SetContentHintsDescription sets content hints description. Port of
// #setContentHintsDescription.
func (p *CAdESSignatureParameters) SetContentHintsDescription(contentHintsDescription string) {
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
func (p *CAdESSignatureParameters) ContentIdentifierSuffix() string {
	return p.contentIdentifierSuffix
}

// SetContentIdentifierSuffix sets content identifier suffix.
//
// NOTE: THIS VALUE WILL BE SET AUTOMATICALLY IF LEFT BLANK
//
// Port of #setContentIdentifierSuffix.
func (p *CAdESSignatureParameters) SetContentIdentifierSuffix(contentIdentifierSuffix string) {
	p.contentIdentifierSuffix = contentIdentifierSuffix
}

// ContentIdentifierPrefix returns the content identifier prefix. See ContentIdentifierSuffix's
// doc. Port of #getContentIdentifierPrefix.
func (p *CAdESSignatureParameters) ContentIdentifierPrefix() string {
	return p.contentIdentifierPrefix
}

// SetContentIdentifierPrefix sets content identifier prefix. Port of #setContentIdentifierPrefix.
func (p *CAdESSignatureParameters) SetContentIdentifierPrefix(contentIdentifierPrefix string) {
	p.contentIdentifierPrefix = contentIdentifierPrefix
}

// IsParallelSignature returns whether a parallel signature should be created when an original
// document is represented by a CMSSignedData (i.e. another CMS signature document). Port of
// #isParallelSignature.
func (p *CAdESSignatureParameters) IsParallelSignature() bool {
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
func (p *CAdESSignatureParameters) SetParallelSignature(parallelSignature bool) {
	p.parallelSignature = parallelSignature
}

// GetContentTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getContentTimestampParameters.
func (p *CAdESSignatureParameters) GetContentTimestampParameters() *CAdESTimestampParameters {
	if p.ContentTimestampParameters == nil {
		p.ContentTimestampParameters = NewCAdESTimestampParameters()
	}
	return p.ContentTimestampParameters
}

// GetSignatureTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getSignatureTimestampParameters.
func (p *CAdESSignatureParameters) GetSignatureTimestampParameters() *CAdESTimestampParameters {
	if p.SignatureTimestampParameters == nil {
		p.SignatureTimestampParameters = NewCAdESTimestampParameters()
	}
	return p.SignatureTimestampParameters
}

// GetArchiveTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating CAdESTimestampParameters. Port of #getArchiveTimestampParameters.
func (p *CAdESSignatureParameters) GetArchiveTimestampParameters() *CAdESTimestampParameters {
	if p.ArchiveTimestampParameters == nil {
		p.ArchiveTimestampParameters = NewCAdESTimestampParameters()
	}
	return p.ArchiveTimestampParameters
}

// String ports #toString.
func (p *CAdESSignatureParameters) String() string {
	return fmt.Sprintf("CAdESSignatureParameters [en319122=%v, contentHintsType='%v', contentHintsDescription='%v', contentIdentifierPrefix='%v', contentIdentifierSuffix='%v', parallelSignature=%v] %s",
		p.en319122, p.contentHintsType, p.contentHintsDescription, p.contentIdentifierPrefix, p.contentIdentifierSuffix, p.parallelSignature,
		p.AbstractSignatureParameters.String())
}

// Equals ports #equals.
func (p *CAdESSignatureParameters) Equals(other *CAdESSignatureParameters) bool {
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
