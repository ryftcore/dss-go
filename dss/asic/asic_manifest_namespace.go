// Ported from dss-asic-common/src/main/java/eu/europa/esig/dss/asic/common/definition/ASiCManifestNamespace.java (DSS 6.5.RC1).
package asic

import "github.com/utain/esig/dss/xml/common"

// ASiCManifestNS is the namespace for ASiC and its schema. Ports ASiCManifestNamespace.NS;
// ASiCManifestNamespace itself was a namespace-only holder class with a private constructor,
// so its other constants survive below as package-level constants.
var ASiCManifestNS = common.NewDSSNamespace("http://uri.etsi.org/02918/v1.2.1#", "asic")

const (
	// ASiCManifestNamespace_LIBREOFFICE_SIGNATURES is the LibreOffice signatures tag name.
	ASiCManifestNamespace_LIBREOFFICE_SIGNATURES = "document-signatures"

	// ASiCManifestNamespace_LIBREOFFICE_NS is the LibreOffice namespace URI.
	ASiCManifestNamespace_LIBREOFFICE_NS = "urn:oasis:names:tc:opendocument:xmlns:digitalsignature:1.0"
)
