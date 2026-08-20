// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/TrustedListUtils.java (DSS 6.5.RC1).

package trustedlist

import _ "embed"

// Schema location constants, matching TrustedListUtils'
// TRUSTED_LIST_SCHEMA_LOCATION/TRUSTED_LIST_SIE_SCHEMA_LOCATION/
// TRUSTED_LIST_ADDITIONALTYPES_SCHEMA_LOCATION. Java resolves these against
// the classpath to build a javax.xml.validation.Schema for
// AbstractJaxbFacade to validate against; this port does not perform
// runtime XSD validation (nothing in TLMODEL's scope consumes it - see
// dss/trustedlist/jaxb's doc.go), so the constants and the embedded bytes
// below exist for provenance/documentation and for any future validator,
// not because anything in this package reads them today.
const (
	TrustedListSchemaLocation                = "/xsd/ts_119612v020401_xsd.xsd"
	TrustedListSIESchemaLocation             = "/xsd/ts_119612v020401_sie_xsd.xsd"
	TrustedListAdditionalTypesSchemaLocation = "/xsd/ts_119612v020401_additionaltypes_xsd.xsd"
)

// TrustedListSchema is ts_119612v020401_xsd.xsd, copied verbatim from
// specs-trusted-list/src/main/resources/xsd.
//
//go:embed xsd/ts_119612v020401_xsd.xsd
var TrustedListSchema []byte

// TrustedListSIESchema is ts_119612v020401_sie_xsd.xsd (the ecc/"SvcInfoExt"
// extension schema), copied verbatim.
//
//go:embed xsd/ts_119612v020401_sie_xsd.xsd
var TrustedListSIESchema []byte

// TrustedListAdditionalTypesSchema is
// ts_119612v020401_additionaltypes_xsd.xsd (the tslx extension schema),
// copied verbatim.
//
//go:embed xsd/ts_119612v020401_additionaltypes_xsd.xsd
var TrustedListAdditionalTypesSchema []byte
