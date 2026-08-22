// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/XAdESNamespace.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdESNamespace defines the list of used XAdES namespaces. Ports the Java
// namespace-only holder class (private constructor, no other members) as package-level vars.
var (
	// XAdESNamespace_XMLDSIG_FILTER2 is The XMLDSIG Filter 2.0 namespace.
	XAdESNamespace_XMLDSIG_FILTER2 = common.NewDSSNamespace("http://www.w3.org/2002/06/xmldsig-filter2", "dsig-filter2")

	// XAdESNamespace_XADES_111 is XAdES 1.1.1.
	XAdESNamespace_XADES_111 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.1.1#", "xades111")

	// XAdESNamespace_XADES_122 is XAdES 1.2.2.
	XAdESNamespace_XADES_122 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.2.2#", "xades122")

	// XAdESNamespace_XADES_132 is XAdES 1.3.2.
	XAdESNamespace_XADES_132 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.3.2#", "xades132")

	// XAdESNamespace_XADES_141 is XAdES 1.4.1.
	XAdESNamespace_XADES_141 = common.NewDSSNamespace("http://uri.etsi.org/01903/v1.4.1#", "xades141")

	// XAdESNamespace_XADES_EVIDENCERECORD_NAMESPACE is XAdES EN 1.1.1 (Evidence record).
	XAdESNamespace_XADES_EVIDENCERECORD_NAMESPACE = common.NewDSSNamespace("http://uri.etsi.org/19132/v1.1.1#", "xadesen")
)
