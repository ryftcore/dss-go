// Ported from DetailedReport.xsd (DSS 6.5.RC1): the registry of the generated
// model types. Every complexType of the schema - i.e. every generated JAXB
// class of eu.europa.esig.dss.detailedreport.jaxb - is listed here exactly
// once (the two simpleType enums, XmlStatus and XmlBlockType, are Go string
// types rather than structs and so have no place in a struct registry).
//
// The registry is what lets the marshaller derive, by reflection over the very
// structs that define the port, which elements carry character data and which
// are complex (see jaxb_content_model.go), and what lets the schema sweep
// check that no element or attribute of the XSD is missing from the model.

package jaxb

import "reflect"

// modelTypes is the set of generated model types, one per schema complexType.
var modelTypes = []reflect.Type{
	reflect.TypeOf(XmlAOV{}),
	reflect.TypeOf(XmlBasicBuildingBlocks{}),
	reflect.TypeOf(XmlCC{}),
	reflect.TypeOf(XmlCRS{}),
	reflect.TypeOf(XmlCV{}),
	reflect.TypeOf(XmlCertificate{}),
	reflect.TypeOf(XmlCertificateApprovalStatus{}),
	reflect.TypeOf(XmlCertificateApprovalStatusProcess{}),
	reflect.TypeOf(XmlCertificateChain{}),
	reflect.TypeOf(XmlCertificateChainCryptographicValidation{}),
	reflect.TypeOf(XmlCertificateQualificationProcess{}),
	reflect.TypeOf(XmlChainItem{}),
	reflect.TypeOf(XmlConclusion{}),
	reflect.TypeOf(XmlConstraint{}),
	reflect.TypeOf(XmlConstraintsConclusion{}),
	reflect.TypeOf(XmlConstraintsConclusionWithControlTime{}),
	reflect.TypeOf(XmlConstraintsConclusionWithProofOfExistence{}),
	reflect.TypeOf(XmlCryptographicAlgorithm{}),
	reflect.TypeOf(XmlCryptographicValidation{}),
	reflect.TypeOf(XmlDetailedReport{}),
	reflect.TypeOf(XmlEAA{}),
	reflect.TypeOf(XmlEvidenceRecord{}),
	reflect.TypeOf(XmlFC{}),
	reflect.TypeOf(XmlISC{}),
	reflect.TypeOf(XmlLoTEAnalysis{}),
	reflect.TypeOf(XmlMessage{}),
	reflect.TypeOf(XmlPCV{}),
	reflect.TypeOf(XmlPSV{}),
	reflect.TypeOf(XmlProofOfExistence{}),
	reflect.TypeOf(XmlQWACProcess{}),
	reflect.TypeOf(XmlRAC{}),
	reflect.TypeOf(XmlRFC{}),
	reflect.TypeOf(XmlRevocationBasicValidation{}),
	reflect.TypeOf(XmlRevocationInformation{}),
	reflect.TypeOf(XmlSAV{}),
	reflect.TypeOf(XmlSemantic{}),
	reflect.TypeOf(XmlSignature{}),
	reflect.TypeOf(XmlSubXCV{}),
	reflect.TypeOf(XmlTLAnalysis{}),
	reflect.TypeOf(XmlTimestamp{}),
	reflect.TypeOf(XmlVCI{}),
	reflect.TypeOf(XmlVTS{}),
	reflect.TypeOf(XmlValidationCertificateApprovalStatus{}),
	reflect.TypeOf(XmlValidationCertificateQualification{}),
	reflect.TypeOf(XmlValidationEAAQualification{}),
	reflect.TypeOf(XmlValidationEAAQualificationProcess{}),
	reflect.TypeOf(XmlValidationPIDQualificationProcess{}),
	reflect.TypeOf(XmlValidationProcessArchivalData{}),
	reflect.TypeOf(XmlValidationProcessArchivalDataTimestamp{}),
	reflect.TypeOf(XmlValidationProcessBasicSignature{}),
	reflect.TypeOf(XmlValidationProcessBasicTimestamp{}),
	reflect.TypeOf(XmlValidationProcessEAA{}),
	reflect.TypeOf(XmlValidationProcessEvidenceRecord{}),
	reflect.TypeOf(XmlValidationProcessLongTermData{}),
	reflect.TypeOf(XmlValidationQWACProcess{}),
	reflect.TypeOf(XmlValidationSignatureQualification{}),
	reflect.TypeOf(XmlValidationTimestampQualification{}),
	reflect.TypeOf(XmlValidationTimestampQualificationAtTime{}),
	reflect.TypeOf(XmlXCV{}),
}
