// Package tsl ports dss-tsl-validation, flattening its Java subpackages
// (alerts, alerts/detections, alerts/handlers/log, cache, cache/access, definition/mra,
// download, dto, dto/builder, dto/condition, function, function/converter, job, parsing,
// runnable, sha2, source, summary, sync, validation) into a single Go package, as instructed
// for this phase.
//
// The flattening was collision-checked: snake_case-ing every upstream file name produces no
// clash across the subpackages, so no subpackage-name filename prefixing was needed. The one
// name that repeats between dss-tsl-validation and dss-validation-job - CacheCleaner - lands
// in a different Go package (dss/validation/job) and therefore does not collide either.
//
// # What this package is
//
// It is the ETSI TS 119 612 trusted-list machinery: TLSource/LOTLSource describe where a
// Trusted List (TL) or List of Trusted Lists (LOTL) is fetched from and how its content is
// filtered; the parsing tasks turn the unmarshalled trusted-list XML model into the
// dss/model/tsl DTO records the qualification blocks consume; TLValidatorTask validates a
// trusted list's XAdES signature against the certificates announced for it; the dto/condition
// types implement the certificate-matching criteria that gate qualification trust propagation;
// and the sha2 types implement the ".sha2" digest side-file publication scheme of TS 119 612
// clause 6.1.
//
// # Ported from
//
// dss-tsl-validation (DSS 6.5.RC1). Every file names its own upstream source in its header.
package tsl
