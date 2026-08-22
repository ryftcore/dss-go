// Ported from
// dss-asic-xades/src/main/java/eu/europa/esig/dss/asic/xades/definition/ManifestNamespace.java
// (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/xml/common"

// ManifestNS is the OpenDocument Manifest namespace. Ports ManifestNamespace.NS;
// ManifestNamespace itself was a namespace-only holder class with a private constructor, so it
// has no other surviving members.
//
// See http://docs.oasis-open.org/office/v1.2/OpenDocument-v1.2-part3.pdf (Open Document Format
// for Office Applications (OpenDocument) Version 1.2; Part 3: Packages).
//
// KAT-verified against a Java oracle run over dss-asic-xades 6.5.RC1 + dss-xml-common 6.5.RC1
// (java.lang.reflect over ManifestNamespace.class.getDeclaredFields()) - see
// manifest_namespace_test.go.
var ManifestNS = common.NewDSSNamespace("urn:oasis:names:tc:opendocument:xmlns:manifest:1.0", "manifest")
