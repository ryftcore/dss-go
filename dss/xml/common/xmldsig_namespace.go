// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigNamespace.java (DSS 6.5.RC1).
package common

// XMLDSigNS is the namespace of
// https://www.w3.org/TR/xmldsig-core/xmldsig-core-schema.xsd. Ports XMLDSigNamespace.NS;
// XMLDSigNamespace itself was a namespace-only holder class with a private constructor and
// no other members, so only the constant survives the port (there is nothing else to name
// "XMLDSigNamespace" in Go).
var XMLDSigNS = NewDSSNamespace("http://www.w3.org/2000/09/xmldsig#", "ds")
