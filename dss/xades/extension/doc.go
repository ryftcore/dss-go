// Package extension ports the dss-xades extension subpackage
// (eu.europa.esig.dss.xades.extension), XAdES-specific extend-to-a-higher-level
// logic (adding a signature timestamp, revocation data, or an archive
// timestamp to an existing XAdES signature).
//
// The main entry types are XAdESDocumentExtender and
// XAdESDocumentExtenderFactory; most callers reach this through
// xades.Service rather than directly.
package extension
