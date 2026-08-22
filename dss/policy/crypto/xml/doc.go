// Package cryptoxml ports dss-policy-crypto-xml (DSS 6.5.RC1), an
// implementation of the ETSI TS 119 322 cryptographic suite catalogue over
// an RFC 5698 ("dssc") XML document (schema: xsd/rfc5698.xsd, extended by
// xsd/19322algocatxmlschema.xsd's "MoreDetails" element).
//
// Unlike dss-policy-jaxb (every generated JAXB class under
// target/generated-sources/xjc is ported, per PORTING.md's generated-JAXB
// rule, into dss/policy/jaxb), this package implements only the four
// hand-written classes (CryptographicSuiteXmlCatalogue/-Facade/-Factory/-Utils)
// - no generated eu.europa.esig.dss.policy.crypto.xml.jaxb.* sources exist in
// this port. Rather than leaving CryptographicSuiteXmlCatalogue unbuildable
// on an unported JAXB model, this port collapses the minimal slice of that
// generated model its one caller (CryptographicSuiteXmlCatalogue) actually
// reads into xml_types.go, as plain encoding/xml-tagged Go structs derived
// directly from rfc5698.xsd/19322algocatxmlschema.xsd - the same collapse
// strategy dss/diagnostic/diagnostic_data_facade.go and the sibling
// cryptojson package use for their own unported dependencies (see those
// files' headers). Fields/elements the XSD declares but no in-scope code
// path reads (the trailing xs:any wildcards, ds:Signature) are not
// modelled, matching the same "port what's needed" pragmatism.
package cryptoxml
