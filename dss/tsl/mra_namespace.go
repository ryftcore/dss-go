// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/definition/mra/MRANamespace.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/xml/common"

// MRANamespaceNS is the Trusted List MRA XSD namespace. Port of MRANamespace.NS; the Java
// class is a static-only utils holder with a private constructor, which Go renders as this
// single package-level variable (same shape as xades/definition's TrustedListNamespaceNS).
var MRANamespaceNS = common.NewDSSNamespace("http://ec.europa.eu/tools/lotl/mra/schema/v2#", "mra")
