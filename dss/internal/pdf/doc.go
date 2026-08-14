// Package pdf implements the subset of ISO 32000-1/2 that DSS's PAdES support
// needs: a lenient reader for the whole corpus of signed PDFs DSS validates, and
// an append-only incremental writer for the ones it produces.
//
// # Provenance
//
// Upstream DSS gets its PDF object model from a third party — Apache PDFBox
// (org.apache.pdfbox:pdfbox:3.0.7, driven by dss-pades-pdfbox's
// PdfBoxDocumentReader / PdfBoxSignatureService) — so, exactly as PORTING.md
// prescribes for BouncyCastle-replacement machinery, this package has no Java
// class to mirror one-to-one and therefore lives under internal/. It imports the
// standard library only and never imports a DSS package.
//
// Everything here is scoped by what dss-pades-pdfbox actually calls; see
// internal/pdf/DESIGN.md §0.1 for the enumerated contract. It is deliberately not
// a general-purpose PDF library: no rasterisation, no content-stream
// interpretation, no image filters, no full-rewrite serializer, no PDF/A,
// linearization or tagged PDF.
//
// The reader's leniency is not "best effort": every tolerance is a copy of a
// specific pdfbox 3.0.7 recovery path (pdfbox runs isLenient=true by default and
// Loader.loadPDF never turns it off, so lenient is the only mode), enumerated in
// DESIGN.md §2.7 and pinned by a unit test named for its rule ID. Behavioural
// divergence from pdfbox is a bug even when pdfbox is the one behaving oddly —
// the /ByteRange arithmetic DSS performs downstream depends on reproducing it.
//
// # Layering
//
//	pades  (ported Java classes: PAdESUtils, PdfSigDictWrapper, SingleDssDict, ByteRange, …)
//	  │  imports
//	  ▼
//	internal/pdf      ← this package. stdlib only.
//
// This package knows nothing about CMS, certificates, OCSP or ETSI; it hands
// []byte upward. In particular it parses /ByteRange into []int64 and stops there:
// eu.europa.esig.dss.pades.validation.ByteRange owns validate() and getLength().
package pdf
