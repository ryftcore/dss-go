// Ported from the generated JAXB classes:
//   - ConstraintStatusType.java
//   - ValidationStatusType.java
//   - ValidationConstraintsEvaluationReportType.java
//   - SignatureValidationPolicyType.java
//   - IndividualValidationConstraintReportType.java
//
// (specs-validation-report, DSS 6.5.RC1).
package jaxb

// ConstraintStatusType is the Go form of the generated JAXB class
// ConstraintStatusType (complexType ConstraintStatusType).
type ConstraintStatusType struct {
	Status       ConstraintStatus `xml:"Status"`
	OverriddenBy *string          `xml:"OverriddenBy,omitempty"`
}

// ValidationStatusType is the Go form of the generated JAXB class
// ValidationStatusType (complexType ValidationStatusType).
type ValidationStatusType struct {
	MainIndication                 URIIndication               `xml:"MainIndication"`
	SubIndication                  []URISubIndication          `xml:"SubIndication,omitempty"`
	AssociatedValidationReportData []*ValidationReportDataType `xml:"AssociatedValidationReportData,omitempty"`
}

// ValidationConstraintsEvaluationReportType is the Go form of the generated
// JAXB class ValidationConstraintsEvaluationReportType (complexType
// ValidationConstraintsEvaluationReportType).
type ValidationConstraintsEvaluationReportType struct {
	SignatureValidationPolicy *SignatureValidationPolicyType              `xml:"SignatureValidationPolicy,omitempty"`
	ValidationConstraint      []*IndividualValidationConstraintReportType `xml:"ValidationConstraint,omitempty"`
}

// SignatureValidationPolicyType is the Go form of the generated JAXB class
// SignatureValidationPolicyType (complexType SignatureValidationPolicyType).
type SignatureValidationPolicyType struct {
	SignaturePolicyIdentifier SignaturePolicyIdentifierType `xml:"SignaturePolicyIdentifier"`
	PolicyName                *string                       `xml:"PolicyName,omitempty"`
	FormalPolicyURI           *string                       `xml:"FormalPolicyURI,omitempty"`
	ReadablePolicyURI         *string                       `xml:"ReadablePolicyURI,omitempty"`
	FormalPolicyObject        *VOReferenceType              `xml:"FormalPolicyObject,omitempty"`
}

// IndividualValidationConstraintReportType is the Go form of the generated
// JAXB class IndividualValidationConstraintReportType (complexType
// IndividualValidationConstraintReportType).
type IndividualValidationConstraintReportType struct {
	ValidationConstraintIdentifier string                `xml:"ValidationConstraintIdentifier"`
	ValidationConstraintParameter  []*TypedDataType      `xml:"ValidationConstraintParameter,omitempty"`
	ConstraintStatus               ConstraintStatusType  `xml:"ConstraintStatus"`
	ValidationStatus               *ValidationStatusType `xml:"ValidationStatus,omitempty"`
	Indications                    *RawContent           `xml:"Indications,omitempty"`
}
