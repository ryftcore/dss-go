// Package i18n ports dss-i18n (eu.europa.esig.dss.i18n), the message
// catalog and MessageFormat-style formatter behind every human-readable
// string the EN 319 102-1 validation engine and its reports produce: the
// ETSI validation-process building-block messages (e.g. "the certificate is
// valid") and their positive/negative variants, keyed by MessageTag and
// resolved from an embedded copy of dss-messages.properties.
//
// The main entry types are Provider (looks up and formats a message for
// a MessageTag) and MessageTag itself (one constant per message key defined
// in dss-messages.properties).
package i18n
