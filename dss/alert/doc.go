// Package alert ports dss-alert (eu.europa.esig.dss.alert), the generic
// detect-and-handle mechanism DSS uses to surface non-fatal validation and
// signing conditions (e.g. an expired certificate, a revoked token) as a
// configurable log, exception, or silent response instead of a hard-coded
// choice.
//
// The main entry types are the generic Alert[T] interface, AbstractAlert[T]
// (detector + handler composition), the Detector[T]/Handler[T]
// interfaces, and the ready-made handlers (LogHandler, ExceptionOnStatusAlert,
// SilentOnStatusAlert, CompositeAlertHandler).
package alert
