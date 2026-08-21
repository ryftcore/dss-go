// Package trustedlist ports specs-trusted-list, the ETSI TS 119 612
// trusted-list XML facade layer: (un)marshaling and helper accessors over
// the generated JAXB-equivalent model in trustedlist/jaxb, plus the
// eIDAS Mutual Recognition Agreement (MRA) extension types.
//
// # Main entry types
//
// TrustedListFacade marshals and unmarshals a trustedlist/jaxb.TrustStatusListType
// to and from its ETSI TS 119 612 XML representation. MRAFacade does the
// same for the MRA extension XML fragments. The mra_* and
// trusted_list_utils.go helpers implement the equivalence-context, status
// and general lookups the qualification process (dss/validation/process/qualification)
// and the tsl package build on.
//
// Ported from specs-trusted-list (DSS 6.5.RC1); every file names its own
// upstream source in its header.
package trustedlist
