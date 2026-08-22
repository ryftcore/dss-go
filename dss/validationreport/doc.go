// Package validationreport is the root of the ETSI TS 119 102-2 Validation
// Report port: the facade, utils, enums and URI-based parser of
// specs-validation-report (DSS 6.5.RC1), flattened into one package. The
// generated JAXB model lives in the jaxb subpackage.
//
// # Enum and parser re-exports
//
// ObjectType, ConstraintStatus, TypeOfProof, SignatureValidationProcessID
// and the UriBasedEnumParser functions are defined in package jaxb, not
// here, and re-exported from this package as type aliases and thin
// forwarding functions (see object_type.go, constraint_status.go,
// type_of_proof.go, signature_validation_process_id.go and
// uri_based_enum_parser.go). This split is forced by Go's import-cycle
// rule: ValidationReportFacade (this package) must import jaxb for
// ValidationReportType, and the generated model's ConstraintStatusType/
// POEType/SignatureValidationProcessType/ValidationObjectType fields are
// typed with these very enums, so jaxb must be able to use them too - Java's
// classpath has no such constraint, but a Go package cannot both import a
// package and be imported by it. Placing the concrete definitions in jaxb
// and aliasing them here keeps every symbol in dss/validationreport working
// and identical, following the same alias-on-extraction idiom PORTING.md
// documents for internal/ packages, applied in the opposite direction. See
// jaxb/jaxb_enums.go's header for the full rationale.
package validationreport
