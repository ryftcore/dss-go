// Ported from the generated JAXB classes ValidationReportType.java and
// SignatureValidationReportType.java (specs-validation-report, DSS 6.5.RC1),
// plus the modelTypes registry jaxb_content_model.go's self-closing
// heuristic needs (see dss/diagnostic/jaxb/jaxb_model.go, Phase 8a, for the
// same registry's rationale): every complexType of 1910202xmlSchema.xsd,
// listed exactly once.
package jaxb

import (
	"encoding/xml"
	"reflect"
)

// ValidationReportType is the Go form of the generated JAXB class
// ValidationReportType (complexType ValidationReportType, the document
// element).
//
// # Root namespace declarations
//
// The JAXB reference implementation pre-declares every namespace the tree
// uses once, at this document element, with a "nsN" prefix assigned in
// document-traversal order (see jaxb_crossns.go's header for the fuller
// story and why "ns2" is always XMLDSig). encoding/xml has no equivalent:
// it declares a namespace locally, at first use, wherever a field tag
// needs one. So this type captures the source document's root xmlns:*
// declarations verbatim on Unmarshal (extraNamespaces) and replays them,
// in the same order, on Marshal - which is sufficient for byte-parity
// because everything that *uses* one of those prefixes (DigestMethodType,
// DSDigestValue, SignatureValueType, SignatureType's own "ns2:", and the
// literal "ns3:"/"ns4:" tokens preserved inside SignaturePolicyIdentifierType/
// DigitalIdentityType/TSPInformationType's raw-captured content) only
// needs the prefix declared somewhere in scope, not any particular value.
type ValidationReportType struct {
	XMLName                    xml.Name
	SignatureValidationReport  []*SignatureValidationReportType `xml:"SignatureValidationReport"`
	SignatureValidationObjects *ValidationObjectListType        `xml:"SignatureValidationObjects,omitempty"`
	SignatureValidator         *SignatureValidatorType          `xml:"SignatureValidator,omitempty"`
	Signature                  *SignatureType                   `xml:"http://www.w3.org/2000/09/xmldsig# Signature,omitempty"`

	// extraNamespaces holds the root element's captured xmlns:* declarations
	// (prefixed ones only; the default "xmlns" for Namespace itself is
	// re-derived on Marshal), in source document order.
	extraNamespaces []xml.Attr
}

// UnmarshalXML captures the root element's xmlns:* declarations (see the
// type's doc comment) before decoding the rest of the document normally.
func (v *ValidationReportType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, a := range start.Attr {
		if a.Name.Space == "xmlns" {
			v.extraNamespaces = append(v.extraNamespaces, a)
		}
	}
	type alias ValidationReportType
	aux := (*alias)(v)
	return d.DecodeElement(aux, &start)
}

// MarshalXML writes the document element with Namespace's default xmlns
// declaration followed by the captured extraNamespaces declarations, then
// the document's content in schema order.
func (v *ValidationReportType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "ValidationReport"}
	start.Attr = make([]xml.Attr, 0, 1+len(v.extraNamespaces))
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: Namespace})
	extraNamespaces := v.extraNamespaces
	if len(extraNamespaces) == 0 {
		// A report BUILT from scratch (dss/validation/executor's
		// ETSIValidationReportBuilder, the Go counterpart of Java's) has no
		// captured declarations. The JAXB RI still pre-declares every namespace
		// of the bound packages at the document element, in a fixed order, for
		// exactly such a marshal - so reproduce that set here. Phase 8f
		// addition, flagged for the integrator: it only applies when nothing was
		// captured, so round-tripping an unmarshalled report is unaffected.
		extraNamespaces = jaxbRootNamespaces
	}
	for _, a := range extraNamespaces {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:" + a.Name.Local}, Value: a.Value})
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for _, sr := range v.SignatureValidationReport {
		if err := e.EncodeElement(sr, xml.StartElement{Name: xml.Name{Local: "SignatureValidationReport"}}); err != nil {
			return err
		}
	}
	if v.SignatureValidationObjects != nil {
		if err := e.EncodeElement(v.SignatureValidationObjects, xml.StartElement{Name: xml.Name{Local: "SignatureValidationObjects"}}); err != nil {
			return err
		}
	}
	if v.SignatureValidator != nil {
		if err := e.EncodeElement(v.SignatureValidator, xml.StartElement{Name: xml.Name{Local: "SignatureValidator"}}); err != nil {
			return err
		}
	}
	if v.Signature != nil {
		// SignatureType.MarshalXML overrides the element name to "ns2:Signature"
		// regardless of what is passed here.
		if err := e.EncodeElement(v.Signature, xml.StartElement{Name: xml.Name{Local: "Signature"}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// jaxbRootNamespaces is the prefix set the JAXB RI declares on the document
// element of every marshalled validation report, in the order it writes them
// (verified against the 144-file upstream oracle corpus: every
// <ValidationReport> start tag carries exactly these, ns2 then ns4 then ns3).
var jaxbRootNamespaces = []xml.Attr{
	{Name: xml.Name{Local: "ns2"}, Value: "http://www.w3.org/2000/09/xmldsig#"},
	{Name: xml.Name{Local: "ns4"}, Value: "http://uri.etsi.org/02231/v2#"},
	{Name: xml.Name{Local: "ns3"}, Value: "http://uri.etsi.org/01903/v1.3.2#"},
}

// SignatureValidationReportType is the Go form of the generated JAXB class
// SignatureValidationReportType (complexType SignatureValidationReportType).
type SignatureValidationReportType struct {
	SignatureIdentifier                   *SignatureIdentifierType                   `xml:"SignatureIdentifier,omitempty"`
	ValidationConstraintsEvaluationReport *ValidationConstraintsEvaluationReportType `xml:"ValidationConstraintsEvaluationReport,omitempty"`
	ValidationTimeInfo                    *ValidationTimeInfoType                    `xml:"ValidationTimeInfo,omitempty"`
	SignersDocument                       *SignersDocumentType                       `xml:"SignersDocument,omitempty"`
	SignatureAttributes                   *SignatureAttributesType                   `xml:"SignatureAttributes,omitempty"`
	SignerInformation                     *SignerInformationType                     `xml:"SignerInformation,omitempty"`
	SignatureQuality                      *SignatureQualityType                      `xml:"SignatureQuality,omitempty"`
	SignatureValidationProcess            *SignatureValidationProcessType            `xml:"SignatureValidationProcess,omitempty"`
	SignatureValidationStatus             ValidationStatusType                       `xml:"SignatureValidationStatus"`
	OtherInformation                      *AnyType                                   `xml:"OtherInformation,omitempty"`
}

// modelTypes is the set of generated model types, one per schema
// complexType, that carry element content via ordinary struct-tag
// reflection. The types with hand-written MarshalXML/UnmarshalXML
// (SignatureAttributesType, ValidationObjectRepresentationType,
// SignersDocumentType, SARevIDListType) are still listed: their own fields
// (or, for the choice types, the fixed Items slice) do not drive the
// self-closing heuristic, but omitting them would leave the types they
// nest (e.g. SACertIDListType via SigningCertificate) unreachable by the
// registry walk if nothing else referenced them - every one here is also
// reachable some other way, so this is a completeness safety net, not load
// bearing today.
var modelTypes = []reflect.Type{
	reflect.TypeOf(AnyType{}),
	reflect.TypeOf(RawContent{}),
	reflect.TypeOf(VOReferenceType{}),
	reflect.TypeOf(NsPrefixMappingType{}),
	reflect.TypeOf(TypedDataType{}),
	reflect.TypeOf(AdditionalValidationReportDataType{}),

	reflect.TypeOf(DigestMethodType{}),
	reflect.TypeOf(DigestAlgAndValueType{}),
	reflect.TypeOf(SignatureValueType{}),
	reflect.TypeOf(SignatureType{}),
	reflect.TypeOf(SignaturePolicyIdentifierType{}),
	reflect.TypeOf(DigitalIdentityType{}),
	reflect.TypeOf(TSPInformationType{}),

	reflect.TypeOf(ConstraintStatusType{}),
	reflect.TypeOf(ValidationStatusType{}),
	reflect.TypeOf(ValidationConstraintsEvaluationReportType{}),
	reflect.TypeOf(SignatureValidationPolicyType{}),
	reflect.TypeOf(IndividualValidationConstraintReportType{}),

	reflect.TypeOf(SignatureIdentifierType{}),
	reflect.TypeOf(SignatureValidatorType{}),
	reflect.TypeOf(SignatureReferenceType{}),
	reflect.TypeOf(XAdESSignaturePtrType{}),

	reflect.TypeOf(SignerInformationType{}),
	reflect.TypeOf(SignatureQualityType{}),
	reflect.TypeOf(SignatureValidationProcessType{}),
	reflect.TypeOf(ValidationTimeInfoType{}),
	reflect.TypeOf(POEType{}),
	reflect.TypeOf(POEProvisioningType{}),
	reflect.TypeOf(CertificateChainType{}),
	reflect.TypeOf(RevocationStatusInformationType{}),
	reflect.TypeOf(CryptoInformationType{}),
	reflect.TypeOf(ValidationReportDataType{}),

	reflect.TypeOf(AttributeBaseType{}),
	reflect.TypeOf(SASigningTimeType{}),
	reflect.TypeOf(SACertIDType{}),
	reflect.TypeOf(SACertIDListType{}),
	reflect.TypeOf(SADataObjectFormatType{}),
	reflect.TypeOf(SACommitmentTypeIndicationType{}),
	reflect.TypeOf(SASigPolicyIdentifierType{}),
	reflect.TypeOf(SASignatureProductionPlaceType{}),
	reflect.TypeOf(SAOneSignerRoleType{}),
	reflect.TypeOf(SASignerRoleType{}),
	reflect.TypeOf(SACounterSignatureType{}),
	reflect.TypeOf(SACRLIDType{}),
	reflect.TypeOf(SAOCSPIDType{}),
	reflect.TypeOf(SARevIDListType{}),
	reflect.TypeOf(SADSSType{}),
	reflect.TypeOf(SATimestampType{}),
	reflect.TypeOf(SAVRIType{}),
	reflect.TypeOf(SAReasonType{}),
	reflect.TypeOf(SANameType{}),
	reflect.TypeOf(SAContactInfoType{}),
	reflect.TypeOf(SASubFilterType{}),
	reflect.TypeOf(SAFilterType{}),
	reflect.TypeOf(SAMessageDigestType{}),
	reflect.TypeOf(SignatureAttributesType{}),

	reflect.TypeOf(ValidationObjectListType{}),
	reflect.TypeOf(ValidationObjectType{}),
	reflect.TypeOf(ValidationObjectRepresentationType{}),
	reflect.TypeOf(SignersDocumentType{}),

	reflect.TypeOf(ValidationReportType{}),
	reflect.TypeOf(SignatureValidationReportType{}),
}
