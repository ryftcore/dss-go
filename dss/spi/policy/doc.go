// Package policy ports the dss-spi policy subpackage
// (eu.europa.esig.dss.spi.policy), resolution and validation of the
// explicit signature policy a signature references (ASN.1- or
// non-ASN.1-encoded), used by the validation engine's EPES-level checks.
//
// The main entry types are SignaturePolicyProvider (locates a policy
// document by identifier/URL), and the SignaturePolicyValidator
// implementations built on AbstractSignaturePolicyValidator
// (BasicASN1SignaturePolicyValidator, NonASN1SignaturePolicyValidator,
// EmptySignaturePolicyValidator).
package policy
