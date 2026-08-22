// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestNamespace.java (DSS 6.5.RC1).
package asic

import "github.com/ryftcore/dss-go/dss/xml/common"

// ASiCManifestNS is the namespace for ASiC and its schema. Ports ASiCManifestNamespace.NS;
// ASiCManifestNamespace itself was a namespace-only holder class with a private constructor,
// so its other constants survive below as package-level constants.
var ASiCManifestNS = common.NewDSSNamespace("http://uri.etsi.org/02918/v1.2.1#", "asic")

const (
	// ASiCManifestNamespaceLibreOfficeSignatures is the LibreOffice signatures tag name.
	ASiCManifestNamespaceLibreOfficeSignatures = "document-signatures"

	// ASiCManifestNamespaceLibreOfficeNS is the LibreOffice namespace URI.
	ASiCManifestNamespaceLibreOfficeNS = "urn:oasis:names:tc:opendocument:xmlns:digitalsignature:1.0"
)
