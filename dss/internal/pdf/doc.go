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
// class to mirror one-to-one and therefore lives under internal/. It never
// imports a DSS package. Beyond the standard library it imports exactly one
// module, golang.org/x/text, and from one file: saslprep.go needs NFKC
// normalisation and Unicode bidirectional classes to reproduce what pdfbox's
// SaslPrep gets from java.text.Normalizer and Character.getDirectionality (see
// DESIGN.md §2.6). Nothing else here reaches outside the standard library.
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
// DESIGN.md §2.7 and pinned by a unit test named for its rule ID. On those
// recovery paths behavioural divergence from pdfbox is a bug even when pdfbox is
// the one behaving oddly — the /ByteRange arithmetic DSS performs downstream
// depends on reproducing it.
//
// The one carve-out is refusal. Where pdfbox recovers by accepting a document
// whose own integrity check has just failed, this package may reject instead,
// because it is the network-facing side of a signature validator and a document
// under validation is hostile input. Such a refusal never widens what is
// accepted, is listed in DESIGN.md §2.6, and carries a "DIVERGENCE, deliberate:"
// comment naming the pdfbox method it departs from and the concrete input that
// motivates it.
//
// # Layering
//
//	pades  (ported Java classes: PAdESUtils, PdfSigDictWrapper, SingleDssDict, ByteRange, …)
//	  │  imports
//	  ▼
//	internal/pdf      ← this package. stdlib, plus golang.org/x/text in saslprep.go.
//
// This package knows nothing about CMS, certificates, OCSP or ETSI; it hands
// []byte upward. In particular it parses /ByteRange into []int64 and stops there:
// eu.europa.esig.dss.pades.validation.ByteRange owns validate() and getLength().
package pdf
