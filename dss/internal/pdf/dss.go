// The /DSS dictionary, its /VRI sub-dictionaries and the token streams they
// point at. See DESIGN.md §3.5.
//
// Provenance: PdfBoxSignatureService.buildDSSDictionary / getPdfObjectForToken
// and PdfBoxDocumentReader.createCOSStream of dss-pades-pdfbox, which is the
// only thing upstream that writes this dictionary.
//
// Two rules are worth restating because they are easy to "improve" wrongly:
//
//   - Ordering inside /Certs, /CRLs, /OCSPs and inside /VRI is the caller's
//     slice order. internal/pdf never sorts and never deduplicates: upstream's
//     dedup is an indexOf check over token identity, and token identity is
//     knowledge the pades layer has and this package does not.
//   - A token whose TokenRef.Key is non-zero is not re-serialised; the object
//     already in the document is referenced instead. That is why the reader has
//     to expose the object key of every array element (§4.3 IndexRef).

package pdf

import "time"

// TokenRef is one validation-data token: either a reference to an object the
// document already carries, or the DER bytes of a new one.
type TokenRef struct {
	Key  ObjectKey // non-zero reuses an existing object; zero writes a new stream
	Data []byte    // DER; ignored when Key is non-zero
}

// VRIEntry is one /VRI sub-dictionary.
type VRIEntry struct {
	Name  string // uppercase base-16 SHA-1 of the signature, computed by pades
	Certs []TokenRef
	CRLs  []TokenRef
	OCSPs []TokenRef
	TU    time.Time
	TS    []byte
}

// DSSDictionary is the document security store to write.
type DSSDictionary struct {
	Certs []TokenRef
	CRLs  []TokenRef
	OCSPs []TokenRef
	VRI   []VRIEntry // omitted from output when empty
}

// SetDSSDictionary writes /Root /DSS and marks the catalog updated. Order is
// the caller's; nothing is sorted or deduplicated here.
func (u *Updater) SetDSSDictionary(dss DSSDictionary) error {
	cat := u.Catalog()

	d := NewDict()
	if a := u.tokenArray(dss.Certs); len(a) > 0 {
		d.Set("Certs", a)
	}
	if a := u.tokenArray(dss.CRLs); len(a) > 0 {
		d.Set("CRLs", a)
	}
	if a := u.tokenArray(dss.OCSPs); len(a) > 0 {
		d.Set("OCSPs", a)
	}
	if len(dss.VRI) > 0 {
		vri := NewDict()
		for _, e := range dss.VRI {
			// Direct dictionary, as sigVriDictionary.setDirect(true) upstream:
			// the /VRI entries are inline in the /DSS object.
			ed := NewDict()
			if a := u.tokenArray(e.Certs); len(a) > 0 {
				ed.Set("Cert", a)
			}
			if a := u.tokenArray(e.CRLs); len(a) > 0 {
				ed.Set("CRL", a)
			}
			if a := u.tokenArray(e.OCSPs); len(a) > 0 {
				ed.Set("OCSP", a)
			}
			if !e.TU.IsZero() {
				ed.Set("TU", String{Bytes: []byte(FormatDate(e.TU))})
			}
			if len(e.TS) > 0 {
				ed.Set("TS", u.addTokenStream(e.TS))
			}
			vri.Set(Name(e.Name), ed)
		}
		d.Set("VRI", vri)
	}

	// pdfbox writes a freshly built COSDictionary as its own indirect object
	// (COSBase.direct defaults to false), which is what makes /DSS addressable
	// by object key on the reading side.
	k := u.Add(d)
	cat.Set("DSS", Ref{Num: k.Num, Gen: k.Gen})
	return nil
}

// tokenArray turns token references into an array of indirect references,
// writing a new stream object for every token that does not already exist.
func (u *Updater) tokenArray(tokens []TokenRef) Array {
	if len(tokens) == 0 {
		return nil
	}
	out := make(Array, 0, len(tokens))
	for _, t := range tokens {
		if !t.Key.IsZero() {
			out = append(out, Ref{Num: t.Key.Num, Gen: t.Key.Gen})
			continue
		}
		out = append(out, u.addTokenStream(t.Data))
	}
	return out
}

// addTokenStream writes the DER bytes as a new plain stream object, unfiltered
// — no /Filter — which is what PdfBoxDocumentReader.createCOSStream produces
// and what keeps the increment readable.
func (u *Updater) addTokenStream(der []byte) Ref {
	k := u.Add(NewStream(NewDict(), append([]byte(nil), der...)))
	return Ref{Num: k.Num, Gen: k.Gen}
}
