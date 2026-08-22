// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/XAdESNamespace.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdESNamespace defines the list of used XAdES namespaces. Ports the Java
// namespace-only holder class (private constructor, no other members) as package-level vars.
var (
	// XAdESNamespaceXMLDSIGFilter2 is The XMLDSIG Filter 2.0 namespace.
	XAdESNamespaceXMLDSIGFilter2 = common.NewDSSNamespace("http://www.w3.org/2002/06/xmldsig-filter2", "dsig-filter2")

	// XAdESNamespaceXAdES111 is XAdES 1.1.1.
	XAdESNamespaceXAdES111 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.1.1#", "xades111")

	// XAdESNamespaceXAdES122 is XAdES 1.2.2.
	XAdESNamespaceXAdES122 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.2.2#", "xades122")

	// XAdESNamespaceXAdES132 is XAdES 1.3.2.
	XAdESNamespaceXAdES132 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.3.2#", "xades132")

	// XAdESNamespaceXAdES141 is XAdES 1.4.1.
	XAdESNamespaceXAdES141 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.4.1#", "xades141")

	// XAdESNamespaceXAdESEvidencerecordNamespace is XAdES EN 1.1.1 (Evidence record).
	XAdESNamespaceXAdESEvidencerecordNamespace = common.NewDSSNamespace("http://uri.etsi.org/19132/v1.1.1#", "xadesen")
)
