// Package extension ports the dss-cades extension subpackage
// (eu.europa.esig.dss.cades.extension), CAdES-specific extend-to-a-higher-level
// logic (adding a signature timestamp, revocation data, or an archive
// timestamp to an existing CAdES signature).
//
// The main entry types are CAdESDocumentExtender and
// CAdESDocumentExtenderFactory; most callers reach this through
// cades.CAdESService rather than directly.
package extension
