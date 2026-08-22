// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/mra/MRAUtils.java (DSS 6.5.RC1).

package trustedlist

import _ "embed"

// MRASchemaLocation is MRAUtils.MRA_SCHEMA_LOCATION. See
// trusted_list_utils.go's header: no runtime XSD validation is performed
// and no JAXBContext is built from it.
const MRASchemaLocation = "/xsd/mra/mra_schema_v2_19612v020401.xsd"

// MRASchema is mra_schema_v2_19612v020401.xsd, copied verbatim from
// specs-trusted-list/src/main/resources/xsd/mra.
//
//go:embed xsd/mra_schema_v2_19612v020401.xsd
var MRASchema []byte

// MRABaseSchema is mra_schema_v2.xsd, the base MRA schema
// mra_schema_v2_19612v020401.xsd imports (not referenced by MRAUtils
// itself, which only names the combined one above, but copied alongside it
// since XSDs ts_119612 + mra are copied+embedded as a set).
//
//go:embed xsd/mra_schema_v2.xsd
var MRABaseSchema []byte
