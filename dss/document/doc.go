// Package document ports the shared, format-independent parts of dss-document
// (eu.europa.esig.dss.signature / eu.europa.esig.dss.extension /
// eu.europa.esig.dss.evidencerecord): the generic building blocks that
// CAdES, XAdES, PAdES, JAdES and ASiC signature and extension services are
// built from, so that logic common to every format (parameter validation,
// resources handling, profile bookkeeping) is written once.
//
// The main entry types are AbstractSignatureService and
// AbstractDocumentExtender (the generic bases per-format services embed),
// the SignatureService/CounterSignatureService/
// MultipleDocumentsSignatureService/EvidenceRecordIncorporationService
// interfaces those services implement, AbstractSignatureParameters (the
// parameter base type extended by e.g. SignatureParameters), and the
// ResourcesHandler family (InMemoryResourcesHandler, TempFileResourcesHandler)
// used to materialize large intermediate signing artifacts.
//
// This package is mostly consumed indirectly, through the per-format
// packages that embed its types; direct use is for advanced callers writing
// a new format binding.
package document
