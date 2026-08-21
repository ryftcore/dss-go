// Package dss is a Go port of the EU Digital Signature Services (DSS)
// library (https://github.com/esig/dss), targeting interoperability
// compatibility with upstream DSS 6.5: signatures produced here validate in
// Java DSS and vice versa, and validation verdicts and report schemas match.
//
// This root package is a thin, idiomatic facade over the ported packages. It
// covers the two paths most applications need - create a signature, validate a
// signed document - and nothing else. Every function here delegates to the
// same services the port's own cross-validation harness drives; the facade
// holds no signature, canonicalization, or validation logic of its own.
//
// # Signing
//
// [Sign] runs the two-step DSS signing dance (compute the data to be signed,
// then embed the signature value) for one of the six [Format] values, at one of
// the four baseline [Level] values:
//
//	signer, err := dss.OpenPKCS12("keystore.p12", "password")
//	if err != nil { return err }
//	defer signer.Close()
//
//	doc, err := dss.OpenDocument("contract.pdf")
//	if err != nil { return err }
//
//	signed, err := dss.Sign(doc, signer, dss.SignOptions{
//	    Format: dss.FormatPAdES,
//	    Level:  dss.LevelB,
//	})
//	if err != nil { return err }
//	err = signed.Save("contract-signed.pdf")
//
// [Extend] raises an existing signature to a higher level (B to T, T to LT,
// and so on).
//
// # Validation
//
// [Validate] auto-detects the document format and runs the full ETSI EN
// 319 102-1 validation process, returning the DSS reports:
//
//	reports, err := dss.Validate(doc, dss.ValidateOptions{})
//	if err != nil { return err }
//	for _, v := range reports.Verdicts() {
//	    fmt.Println(v.ID, v.Indication, v.SubIndication, v.Qualification)
//	}
//
// [Reports] embeds the upstream *reports.Reports, so the full DSS report API -
// diagnostic data, detailed report, ETSI validation report, and the XML
// marshalling of each - is available alongside the convenience accessors.
//
// # What the facade deliberately leaves out
//
// The facade exposes the common paths only. Anything else - visible PAdES
// signature appearances, counter-signatures, evidence records, XAdES
// references and transforms, ASiC filename factories, signature policy stores,
// custom certificate/revocation sources, trusted-list refresh jobs - is
// reached through the packages the facade delegates to, which stay fully
// exported and are documented on their own:
// [github.com/utain/esig/dss/cades], [github.com/utain/esig/dss/xades],
// [github.com/utain/esig/dss/pades], [github.com/utain/esig/dss/jades],
// [github.com/utain/esig/dss/asic], [github.com/utain/esig/dss/validation],
// [github.com/utain/esig/dss/tsl] and [github.com/utain/esig/dss/token].
// Facade types are plain aliases of the underlying ones wherever possible, so
// mixing the two levels needs no conversion.
//
// # Network access
//
// Nothing in this package performs a network request unless you ask for it.
// The upstream OnlineTSPSource, OnlineCRLSource and OnlineOCSPSource of the
// Java dss-service module are NOT part of this port, so:
//
//   - Levels T, LT and LTA require a [TSPSource] you supply
//     (SignOptions.TSPSource). The port ships
//     [github.com/utain/esig/dss/spi/validation.KeyEntityTSPSource], which
//     issues RFC 3161 tokens from a local key - enough for tests and for a
//     self-hosted TSA, but not an HTTP TSA client.
//   - Revocation data for LT and LTA must likewise come from a CRL/OCSP source
//     you set on a [CertificateVerifier] of your own.
//   - Certificate retrieval over AIA is available and is off by default; set
//     ValidateOptions.EnableAIA to turn it on.
//
// # Errors
//
// The ported services follow Java DSS and raise unchecked exceptions, which the
// port turns into panics. Every facade function recovers them and returns them
// as an error wrapping the original value, so errors.As against the port's
// error types (for example [github.com/utain/esig/dss/model.DSSError])
// keeps working.
//
// # Registration
//
// Java DSS discovers format validators, validation policies and cryptographic
// suites through java.util.ServiceLoader. Go has no such mechanism, so
// importing this package registers all six format families, the ETSI
// validation policy and the XML cryptographic suite - which is what makes
// [Validate]'s format auto-detection work out of the box. Applications that
// use the underlying packages directly must perform that registration
// themselves; see [github.com/utain/esig/dss/validation.RegisterDocumentValidatorFactory]
// and [github.com/utain/esig/dss/validation/policy.RegisterValidationPolicyFactory].
//
// # Where to look next
//
// The examples/ directory holds nine runnable programs, one story each, built
// on this facade; cmd/esig is a command-line front end built on it too, and is
// the largest worked example of the API below. The repository README states the
// feature matrix and the accepted gaps, PORTING.md the porting conventions, and
// /PORTING_PLAN.md the module mapping, phase roadmap and compatibility record.
package dss
