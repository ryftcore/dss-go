// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/vpfswatsp/POE.java (DSS 6.5.RC1).
//
// Package flattening: Java's
// eu.europa.esig.dss.validation.process.vpfswatsp splits into five packages -
// the root, checks, checks/pcv(+checks), checks/psv(+checks), checks/vts(+checks)
// and evidencerecord(+checks) - which import each other in both directions
// (PastSignatureValidation in checks/psv builds a PastCertificateValidation from
// checks/pcv, which builds a ValidationTimeSliding from checks/vts, which reads
// the root's POEExtraction; the root's ValidationProcessForSignaturesWithArchivalData
// in turn builds a PastSignatureValidation). Go forbids import cycles, so the
// whole tree is one package `vpfswatsp` here; the class names are collision-free
// across the six Java packages, so every file keeps its upstream name.
//
// POE HIERARCHY. Java's POE is a concrete class with two subclasses that
// override three of its four methods; POEComparator and POEExtraction hold and
// compare the base type polymorphically, and POEComparator additionally tests
// `instanceof TimestampPOE`. The Go form is therefore an interface POE plus a
// POEBase struct carrying the base state and the base method bodies, which
// TimestampPOE and EvidenceRecordPOE embed and shadow (see timestamp_poe.go and
// evidence_record_poe.go). No base method calls an overridable one, so the
// InitXxx overrides-registration pattern the Chain hierarchy needs is not
// required here.
package vpfswatsp

import (
	"time"

	diagnosticjaxb "github.com/ryftcore/dss-go/dss/diagnostic/jaxb"
)

// POE contains Proof Of Existence for validation objects. It is the Go form of
// the Java class POE seen through its overridable API: the concrete base
// behaviour lives in POEBase.
type POE interface {
	// Time returns time of the POE. Port of getTime().
	Time() time.Time
	// POEProviderId returns an Id of a token providing the POE (e.g. a
	// time-stamp Id). Port of getPOEProviderId(): Java's documented NULL for a
	// POE not provided by a token is nil here, since the one reader
	// (ValidationProcessForSignaturesWithArchivalData#toXmlProofOfExistence)
	// copies it into a nullable XmlProofOfExistence member where absent and
	// empty differ.
	POEProviderId() *string
	// POEObjects returns a list of objects covered by the POE if applicable.
	// Port of getPOEObjects().
	POEObjects() []*diagnosticjaxb.XmlTimestampedObject
	// IsTokenProvided returns whether the POE is provided by a token (i.e. a
	// time-stamp or an evidence record). Port of isTokenProvided().
	IsTokenProvided() bool
}

// POEBase is the concrete half of the Java class POE: a POE defined by a
// control/validation time, and the base of TimestampPOE and EvidenceRecordPOE.
type POEBase struct {
	// poeTime is the POE time.
	poeTime time.Time
}

// NewPOE is the constructor to instantiate a global POE by a control/validation
// time. NOTE: the POE will be applied for all tokens. Port of POE(Date).
//
// Java guards the argument with Objects.requireNonNull(controlTime, "The
// controlTime must be defined!"); a Go time.Time is a value that cannot be nil,
// so the guard has no counterpart here. The two call sites that can reach the
// Java constructor with a null - TimestampPOE with a time-stamp carrying no
// production time, and EvidenceRecordPOE whose first time-stamp carries none -
// raise that very message themselves (see the two files).
func NewPOE(controlTime time.Time) *POEBase {
	return &POEBase{poeTime: controlTime}
}

// Time returns time of the POE. Port of getTime().
func (p *POEBase) Time() time.Time {
	return p.poeTime
}

// POEProviderId returns an Id of a token providing the POE. Port of
// getPOEProviderId(), whose base implementation returns NULL.
func (p *POEBase) POEProviderId() *string {
	return nil
}

// POEObjects returns a list of objects covered by the POE if applicable. Port
// of getPOEObjects(), whose base implementation returns Collections.emptyList().
func (p *POEBase) POEObjects() []*diagnosticjaxb.XmlTimestampedObject {
	return nil
}

// IsTokenProvided returns whether the POE is provided by a token. Port of
// isTokenProvided(), whose base implementation returns FALSE.
func (p *POEBase) IsTokenProvided() bool {
	return false
}

// timestampedObjectId reads XmlTimestampedObject#getToken().getId().
//
// The generated Go model splits Java's resolved IDREF into XmlTokenRef.Token
// (the linked object, populated by Unmarshal) and XmlTokenRef.ID (the raw
// idref). Java only ever sees the linked object, so its id is preferred here;
// the raw idref is the fallback for a graph that was never linked, which is the
// case Java would answer with a NullPointerException.
func timestampedObjectId(timestampedObject *diagnosticjaxb.XmlTimestampedObject) string {
	tokenRef := timestampedObject.Token
	if tokenRef.Token != nil {
		return tokenRef.Token.TokenID()
	}
	return tokenRef.ID
}
