// Package timedependent ports the dss-model timedependent subpackage
// (eu.europa.esig.dss.model.timedependent), a small generic container for
// values that change over time (e.g. a certificate qualification that
// differs before/after a given date), letting callers look up the value in
// effect at an arbitrary instant.
//
// The main entry types are the TimeDependent interface,
// Values[T] (an ordered, queryable collection of TimeDependent
// values), MutableValues[T], and BaseTimeDependent.
package timedependent
