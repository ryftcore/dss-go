// Package qwacvalidator ports
// dss-validation/src/main/java/eu/europa/esig/dss/validation/qwac/QWACValidator.java
// (DSS 6.5.RC1).
//
// Package split / import-cycle break: upstream places QWACValidator in
// the SAME Java package as QWACUtils and LinkHeaderParser
// (eu.europa.esig.dss.validation.qwac); a literal Go port would put all
// three in dss/validation/qwac too. That layout is impossible in Go: QWACValidator
// needs the root dss/validation package (SignedDocumentValidator, AbstractCertificateValidator,
// ...), but dss/validation transitively imports dss/validation/qwac already - by design, not by
// accident - because validation/process/qualification's QWAC qualification checks
// (TLSCertificateBindingPresentInSignatureCheck and friends, ported in an earlier phase) call
// qwac.GetIdentifiedTLSCertificates. Put QWACValidator in the qwac package too and the two
// directions close a real cycle:
//
//	dss/validation (root)
//	  -> dss/validation/executor (report builders)
//	  -> dss/validation/process/eaa -> .../process/qualification (QWAC qualification block)
//	  -> dss/validation/qwac                              (for QWACUtils)
//	  -> dss/validation (root)                            (only if QWACValidator lived here too)
//
// The same technique is used for process/eaa's checks subpackage, for the same
// reason. Here the split runs the other way: the
// dependency-free utility code (QWACUtils, LinkHeaderParser) stays at dss/validation/qwac -
// exactly where validation/process/qualification needs to find it - and QWACValidator, the one
// member of the Java package that actually needs the root validation package, moves one level
// down into this qwacvalidator package instead. Nothing in dss/validation or
// dss/validation/process needs QWACValidator itself (only the reverse: QWACValidator needs
// them), so this direction carries no cycle.
package qwacvalidator
