// Package alerts ports the dss-spi alerts subpackage
// (eu.europa.esig.dss.spi.alerts), a ready-made alert.Alert wiring for one
// specific SPI-layer condition: an external resource (e.g. an AIA/OCSP/CRL
// endpoint) failing in a way that should be surfaced through the alert
// package rather than silently ignored or hard-failing.
//
// The main entry type is DSSExternalResourceExceptionAlert.
package alerts
