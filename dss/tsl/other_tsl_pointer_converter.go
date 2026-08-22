// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/converter/OtherTSLPointerConverter.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/model"
	tslmodel "github.com/ryftcore/dss-go/dss/model/tsl"
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// Private static element-name constants.
const (
	otherTSLPointerConverterSchemeTerritory          = "SchemeTerritory"
	otherTSLPointerConverterTSLType                  = "TSLType"
	otherTSLPointerConverterMimeType                 = "MimeType"
	otherTSLPointerConverterSchemeOperatorName       = "SchemeOperatorName"
	otherTSLPointerConverterSchemeTypeCommunityRules = "SchemeTypeCommunityRules"
	otherTSLPointerConverterMRA                      = "MutualRecognitionAgreementInformation"
)

// OtherTSLPointerConverter converts an OtherTSLPointerType to an OtherTSLPointer.
type OtherTSLPointerConverter struct {
	// mraSupport defines whether MRA shall be extracted.
	mraSupport bool
}

// NewOtherTSLPointerConverterDefault instantiates an empty object. Port of
// OtherTSLPointerConverter().
//
// Named ...Default (rather than the bare NewOtherTSLPointerConverter PORTING.md's overload
// convention would otherwise give this zero-arg constructor) because tsl/lotl_parsing_task.go
// already calls NewOtherTSLPointerConverter(bool) for the single-arg constructor below - this
// file conforms to that call site instead of the reverse, since the zero-arg constructor has no
// callers anywhere in this module.
func NewOtherTSLPointerConverterDefault() *OtherTSLPointerConverter {
	return &OtherTSLPointerConverter{}
}

// NewOtherTSLPointerConverter is the constructor with a parameter to define the MRA support.
// Port of OtherTSLPointerConverter(boolean). See NewOtherTSLPointerConverterDefault's comment for
// why this (not that) is the bare name.
func NewOtherTSLPointerConverter(mraSupport bool) *OtherTSLPointerConverter {
	return &OtherTSLPointerConverter{mraSupport: mraSupport}
}

// Apply ports apply(OtherTSLPointerType).
func (c *OtherTSLPointerConverter) Apply(original *jaxb.OtherTSLPointerType) *tslmodel.OtherTSLPointer {
	return tslmodel.NewOtherTSLPointerBuilder().
		SetSdiCertificates(c.getCertificates(original.ServiceDigitalIdentities)).
		SetTslLocation(original.TSLLocation).
		SetSchemeTerritory(c.getSchemeTerritory(original.AdditionalInformation)).
		SetTslType(c.getTSLType(original.AdditionalInformation)).
		SetMimeType(c.getMimeType(original.AdditionalInformation)).
		SetSchemeOperatorNames(c.getSchemeOperatorNames(original.AdditionalInformation)).
		SetSchemeTypeCommunityRules(c.getSchemeTypeCommunityRules(original.AdditionalInformation)).
		SetMra(c.getMRA(original.AdditionalInformation)).
		Build()
}

func (c *OtherTSLPointerConverter) getCertificates(serviceDigitalIdentities *jaxb.ServiceDigitalIdentityListType) []*model.CertificateToken {
	var certificates []*model.CertificateToken
	if serviceDigitalIdentities != nil && utils.IsCollectionNotEmpty(serviceDigitalIdentities.ServiceDigitalIdentity) {
		converter := NewDigitalIdentityListTypeConverter()
		for _, digitalIdentityList := range serviceDigitalIdentities.ServiceDigitalIdentity {
			certificates = append(certificates, converter.Apply(digitalIdentityList)...)
		}
	}
	return certificates
}

func (c *OtherTSLPointerConverter) getSchemeTerritory(additionalInformation *jaxb.AdditionalInformationType) string {
	v, _ := otherTSLPointerConverterOtherInformationValue[*string](additionalInformation, otherTSLPointerConverterSchemeTerritory)
	if v == nil {
		return ""
	}
	return *v
}

func (c *OtherTSLPointerConverter) getTSLType(additionalInformation *jaxb.AdditionalInformationType) string {
	v, _ := otherTSLPointerConverterOtherInformationValue[*string](additionalInformation, otherTSLPointerConverterTSLType)
	if v == nil {
		return ""
	}
	return *v
}

func (c *OtherTSLPointerConverter) getMimeType(additionalInformation *jaxb.AdditionalInformationType) string {
	v, _ := otherTSLPointerConverterOtherInformationValue[*string](additionalInformation, otherTSLPointerConverterMimeType)
	if v == nil {
		return ""
	}
	return *v
}

func (c *OtherTSLPointerConverter) getSchemeOperatorNames(additionalInformation *jaxb.AdditionalInformationType) map[string][]string {
	schemeOperatorNames, ok := otherTSLPointerConverterOtherInformationValue[*jaxb.InternationalNamesType](additionalInformation, otherTSLPointerConverterSchemeOperatorName)
	if ok && schemeOperatorNames != nil {
		return NewInternationalNamesTypeConverter().Apply(schemeOperatorNames)
	}
	return nil
}

func (c *OtherTSLPointerConverter) getSchemeTypeCommunityRules(additionalInformation *jaxb.AdditionalInformationType) map[string][]string {
	schemeTypeCommunityRules, ok := otherTSLPointerConverterOtherInformationValue[*jaxb.NonEmptyMultiLangURIListType](additionalInformation, otherTSLPointerConverterSchemeTypeCommunityRules)
	if ok && schemeTypeCommunityRules != nil {
		return NewNonEmptyMultiLangURIListTypeConverter().Apply(schemeTypeCommunityRules)
	}
	return nil
}

func (c *OtherTSLPointerConverter) getMRA(additionalInformation *jaxb.AdditionalInformationType) *tslmodel.MRA {
	if !c.mraSupport {
		return nil
	}
	jaxbMRA, ok := otherTSLPointerConverterOtherInformationValue[*jaxb.MutualRecognitionAgreementInformationType](additionalInformation, otherTSLPointerConverterMRA)
	if ok && jaxbMRA != nil {
		return NewMRAConverter().Apply(jaxbMRA)
	}
	return nil
}

// otherTSLPointerConverterOtherInformationValue ports the private
// <T extends Serializable> T getOtherInformationValue(AdditionalInformationType, Class<T>,
// String), matching on the recognized wildcard element's local name and Go type - see
// abstract_other_tsl_pointer_predicate.go's header for the shape of AnyContent/AnyItem this
// walks.
func otherTSLPointerConverterOtherInformationValue[T any](additionalInformation *jaxb.AdditionalInformationType, elementName string) (T, bool) {
	var zero T
	if additionalInformation == nil {
		return zero, false
	}
	for _, item := range additionalInformation.Items {
		otherInformation := item.OtherInformation
		if otherInformation == nil {
			continue
		}
		for _, content := range otherInformation.Items {
			if content.Elem == nil || content.ElemName.Local != elementName {
				continue
			}
			if v, ok := content.Elem.(T); ok {
				return v, true
			}
		}
	}
	return zero, false
}
