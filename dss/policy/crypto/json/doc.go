// Package cryptojson ports dss-policy-crypto-json (DSS 6.5.RC1), an
// implementation of the ETSI TS 119 322 cryptographic suite catalogue over a
// JSON document.
//
// Upstream's classes depend on dss-json-common
// (eu.europa.esig.json.JsonObjectWrapper/JSONParser/RFC3339DateUtils/
// JSONSchemaAbstractUtils), a module this port does not otherwise provide.
// As with dss/diagnostic/diagnostic_data_facade.go (which collapses
// dss-jaxb-common's AbstractJaxbFacade directly using encoding/xml) and
// dss/diagnostic/diagnostic_data_xml_definer.go (which stubs the
// schema-validation half as deferred, no stdlib JSON-schema validator
// existing without a new third-party dependency needing tech-lead sign-off
// per PORTING.md's "Dependency policy"), this package collapses the
// minimal slice of JsonObjectWrapper/JSONParser/RFC3339DateUtils that
// CryptographicSuiteJsonCatalogue/-Factory actually use directly into
// json_object.go, using Go's stdlib encoding/json.
package cryptojson
