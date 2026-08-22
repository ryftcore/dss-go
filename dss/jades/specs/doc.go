// Package specs is the Go port of the specs-jades Maven module (eu.europa.esig.jades), plus the
// minimal machinery from dss-json-common and specs-jws that specs-jades' own classes needed and
// that PORTING.md routes into this package rather than into internal/. specs is dependency-closed:
// it imports only the standard library, never a jades/* (or any other DSS) package, so that
// dss/jades (which imports specs) never forms a cycle.
//
// schema/ holds the ETSI TS 119 182-1 JSON schemas plus their RFC 7515/7517/7519/7797 dependency
// schemas, embedded via go:embed (schema_fs.go) and interpreted by the evidence-bounded JSON
// Schema engine in schema_engine.go (see that file's header for the scope and deviations of that
// engine).
package specs
