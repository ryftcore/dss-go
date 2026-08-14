// Signature placement: the /Sig field, its widget, the AcroForm wiring, the
// /Contents and /ByteRange placeholders, DocMDP and FieldMDP. See DESIGN.md
// §3.3 and rules R17, R18, R19.
//
// Provenance: PDDocument.addSignature and PDSignatureField of pdfbox 3.0.7,
// driven exactly as PdfBoxSignatureService.createSignatureDictionary /
// setFieldMDP / setMDPPermission drive them, plus PAdESUtils.replaceSignature
// for the cached to-be-signed path.

package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// DefaultContentSize is the reserved /Contents size in bytes (R17), matching
// both PAdESSignatureParameters.signatureSize and
// SignatureOptions.DEFAULT_SIGNATURE_SIZE upstream.
const DefaultContentSize = 9472

// reservedByteRange is pdfbox's RESERVE_BYTE_RANGE (R18). Serialized it is
// "[0 1000000000 1000000000 1000000000]": 34 content bytes plus the ']' give
// the 35 bytes the real value is later formatted into.
func reservedByteRange() Array {
	return Array{Integer(0), Integer(1000000000), Integer(1000000000), Integer(1000000000)}
}

// SignatureOptions describes the signature dictionary and the field it goes in.
type SignatureOptions struct {
	Type        Name // "Sig" (default) or "DocTimeStamp"
	Filter      Name // default "Adobe.PPKLite"
	SubFilter   Name // e.g. "ETSI.CAdES.detached", "ETSI.RFC3161"
	ContentSize int  // reserved /Contents bytes; default DefaultContentSize
	SignerName  string
	Reason      string
	Location    string
	ContactInfo string
	SigningTime time.Time // /M; the zero value omits the key
	AppName     string    // /Prop_Build /App /Name
	FieldID     string    // fill this existing empty field; "" creates a new one
	Page        int       // 1-based; used only when creating a field
	Rect        Rect      // the zero Rect creates an invisible field
	Appearance  *Stream   // /AP /N; nil for an invisible field
	DocMDP      int       // 1..3; 0 = none
	Lock        *Dict     // /Lock of the target field, drives FieldMDP
	DocumentID  []byte    // second element of /ID
}

// Placeholder names the objects AddSignature created.
type Placeholder struct {
	SigKey   ObjectKey
	FieldKey ObjectKey
}

// AddSignature installs a signature dictionary with placeholder /Contents and
// /ByteRange, wires it into a signature field, the AcroForm and (for a new
// field) the page's /Annots. It returns ErrSignatureAlreadyAdded on a second
// call (R19), matching PDDocument.addSignature's IllegalStateException:
// multiple signatures are multiple increments.
func (u *Updater) AddSignature(opts SignatureOptions) (*Placeholder, error) {
	if u.sigAdded {
		return nil, ErrSignatureAlreadyAdded
	}
	typ := opts.Type
	if typ == "" {
		typ = "Sig"
	}
	filter := opts.Filter
	if filter == "" {
		filter = "Adobe.PPKLite"
	}
	size := opts.ContentSize
	if size <= 0 {
		size = DefaultContentSize
	}
	isTimeStamp := typ == "DocTimeStamp"

	// The signature object number is reserved first so that the field can
	// reference it; the dictionary itself is scheduled at the end.
	sigKey := u.Alloc()

	// Key order is fixed and is DESIGN.md §3.3's: /Type /Filter /SubFilter
	// /Name /Location /Reason /ContactInfo /M /Prop_Build /Reference /Contents
	// /ByteRange. /M and /Name are omitted for a DocTimeStamp.
	sig := NewDict()
	sig.Set("Type", typ)
	sig.Set("Filter", filter)
	if opts.SubFilter != "" {
		sig.Set("SubFilter", opts.SubFilter)
	}
	if !isTimeStamp && opts.SignerName != "" {
		sig.Set("Name", String{Bytes: []byte(opts.SignerName)})
	}
	if opts.Location != "" {
		sig.Set("Location", String{Bytes: []byte(opts.Location)})
	}
	if opts.Reason != "" {
		sig.Set("Reason", String{Bytes: []byte(opts.Reason)})
	}
	if opts.ContactInfo != "" {
		sig.Set("ContactInfo", String{Bytes: []byte(opts.ContactInfo)})
	}
	if !isTimeStamp && !opts.SigningTime.IsZero() {
		sig.Set("M", String{Bytes: []byte(FormatDate(opts.SigningTime))})
	}
	if opts.AppName != "" {
		sig.Set("Prop_Build", DictOf(
			Name("App"), DictOf(Name("Name"), Name(opts.AppName)),
		))
	}

	cat := u.Catalog()

	// DocMDP (§3.3 step 7): only for a /Sig, only with a permission of 1..3,
	// and only when no filled signature exists yet — "a document can contain
	// only one signature field that contains a DocMDP transform method; it
	// shall be the first signed field in the document".
	if !isTimeStamp && opts.DocMDP >= 1 && opts.DocMDP <= 3 {
		filled, err := u.hasFilledSignature()
		if err != nil {
			return nil, err
		}
		if !filled {
			sig.Set("Reference", Array{DictOf(
				Name("Type"), Name("SigRef"),
				Name("TransformMethod"), Name("DocMDP"),
				Name("TransformParams"), DictOf(
					Name("Type"), Name("TransformParams"),
					Name("P"), Integer(opts.DocMDP),
					Name("V"), Name("1.2"),
				),
			)})
			perms := NewDict()
			if p, ok := cat.GetRaw("Perms").(*Dict); ok {
				perms = p.Clone()
			}
			perms.Set("DocMDP", Ref{Num: sigKey.Num, Gen: sigKey.Gen})
			cat.Set("Perms", perms)
		}
	}

	fieldKey, lock, err := u.placeSignatureField(&opts, sigKey)
	if err != nil {
		return nil, err
	}

	// FieldMDP from a /Lock (§3.3 step 6). setFieldMDP runs after
	// createSignatureDictionary upstream and overwrites /Reference, so a field
	// /Lock wins over DocMDP here too.
	if lock != nil {
		tp := lock.Clone()
		tp.Set("Type", Name("TransformParams"))
		tp.Set("V", Name("1.2"))
		sigRef := DictOf(
			Name("Type"), Name("SigRef"),
			Name("TransformMethod"), Name("FieldMDP"),
			Name("TransformParams"), tp,
		)
		if !u.catalogKey.IsZero() {
			sigRef.Set("Data", Ref{Num: u.catalogKey.Num, Gen: u.catalogKey.Gen})
		}
		sig.Set("Reference", Array{sigRef})
	}

	// R17: "<" + 2*ContentSize ASCII '0' + ">". The reserved bytes are a hex
	// string of ContentSize zero bytes; the CMS is written over them later.
	sig.Set("Contents", String{Bytes: make([]byte, size), Hex: true})
	sig.Set("ByteRange", reservedByteRange()) // R18

	u.Put(sigKey, sig)
	u.sigAdded = true
	u.sig = &signaturePlan{dict: sig, key: sigKey, fieldKey: fieldKey}
	if opts.DocumentID != nil {
		u.documentID = opts.DocumentID
	}
	return &Placeholder{SigKey: sigKey, FieldKey: fieldKey}, nil
}

// placeSignatureField either fills the existing field named by FieldID or
// creates a merged field+widget, and wires the AcroForm and the page in either
// case. It returns the field's key and the /Lock that drives FieldMDP.
func (u *Updater) placeSignatureField(opts *SignatureOptions, sigKey ObjectKey) (ObjectKey, *Dict, error) {
	sigRef := Ref{Num: sigKey.Num, Gen: sigKey.Gen}
	lock := opts.Lock

	if opts.FieldID != "" {
		fields, err := u.doc.SignatureFields()
		if err != nil {
			return ObjectKey{}, nil, err
		}
		for _, f := range fields {
			if f.Name != opts.FieldID {
				continue
			}
			if f.Value != nil {
				// Upstream's wording, kept verbatim so a caller matching on it
				// keeps matching.
				return ObjectKey{}, nil, fmt.Errorf(
					"pdf: The signature field '%s' can not be signed since its already signed.", opts.FieldID)
			}
			obj, err := u.Update(f.Key)
			if err != nil {
				return ObjectKey{}, nil, err
			}
			fd, ok := obj.(*Dict)
			if !ok {
				return ObjectKey{}, nil, fmt.Errorf("pdf: signature field '%s' is not a dictionary", opts.FieldID)
			}
			fd.Set("V", sigRef)
			if lock == nil {
				lock = f.Lock
			}
			if err := u.setPrintFlag(fd, f); err != nil {
				return ObjectKey{}, nil, err
			}
			// The field is already in /Fields; pdfbox only marks the array as
			// needing an update, which we do by re-emitting it.
			if err := u.wireAcroForm(ObjectKey{}, false); err != nil {
				return ObjectKey{}, nil, err
			}
			return f.Key, lock, nil
		}
		return ObjectKey{}, nil, fmt.Errorf("pdf: The signature field '%s' does not exist.", opts.FieldID)
	}

	// A new merged field+widget dictionary (§3.3 step 3).
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	if n := u.doc.NumberOfPages(); n > 0 && page > n {
		page = n
	}
	_, pageKey, err := u.doc.Page(page)
	if err != nil {
		return ObjectKey{}, nil, err
	}

	name, err := u.generateFieldName()
	if err != nil {
		return ObjectKey{}, nil, err
	}

	fieldKey := u.Alloc()
	fd := NewDict()
	fd.Set("FT", Name("Sig"))
	fd.Set("Type", Name("Annot"))
	fd.Set("Subtype", Name("Widget"))
	fd.Set("F", Integer(4)) // Print, per PDF/A-1: Hidden/Invisible/NoView clear
	fd.Set("T", String{Bytes: []byte(name)})
	fd.Set("Rect", opts.Rect.Array())
	if !pageKey.IsZero() {
		fd.Set("P", Ref{Num: pageKey.Num, Gen: pageKey.Gen})
	}
	fd.Set("V", sigRef)
	if opts.Appearance != nil {
		apKey := u.Add(opts.Appearance)
		fd.Set("AP", DictOf(Name("N"), Ref{Num: apKey.Num, Gen: apKey.Gen}))
	}
	u.Put(fieldKey, fd)

	if err := u.wireAcroForm(fieldKey, true); err != nil {
		return ObjectKey{}, nil, err
	}
	if err := u.appendAnnotation(pageKey, fieldKey); err != nil {
		return ObjectKey{}, nil, err
	}
	return fieldKey, lock, nil
}

// setPrintFlag sets the Print bit of the field's first widget, which is what
// PDDocument.addSignature does unconditionally (firstWidget.setPrinted(true)).
func (u *Updater) setPrintFlag(fieldDict *Dict, f SignatureField) error {
	target := fieldDict
	if st, ok := fieldDict.GetRaw("Subtype").(Name); !ok || st != "Widget" {
		if len(f.WidgetKeys) == 0 {
			return nil
		}
		obj, err := u.Update(f.WidgetKeys[0])
		if err != nil {
			return err
		}
		wd, ok := obj.(*Dict)
		if !ok {
			return nil
		}
		target = wd
	}
	flags := int64(0)
	if v, ok := target.GetRaw("F").(Integer); ok {
		flags = int64(v)
	}
	target.Set("F", Integer(flags|4))
	return nil
}

// wireAcroForm creates /AcroForm when it is absent, appends fieldKey to
// /Fields when add is set, and sets /SigFlags 3 (SignaturesExist|AppendOnly).
// /AcroForm is written direct inside the catalog when it was direct and
// indirect when it was indirect (§3.3 step 4).
func (u *Updater) wireAcroForm(fieldKey ObjectKey, add bool) error {
	cat := u.Catalog()
	var af *Dict
	switch v := cat.GetRaw("AcroForm").(type) {
	case Ref:
		obj, err := u.Update(v.Key())
		if err != nil {
			return err
		}
		d, ok := obj.(*Dict)
		if !ok {
			return errors.New("pdf: /AcroForm is not a dictionary")
		}
		af = d
	case *Dict:
		af = v.Clone()
		cat.Set("AcroForm", af)
	default:
		af = NewDict()
		cat.Set("AcroForm", af)
	}

	switch fv := af.GetRaw("Fields").(type) {
	case Ref:
		obj, err := u.Update(fv.Key())
		if err != nil {
			return err
		}
		arr, _ := obj.(Array)
		if add && !fieldKey.IsZero() {
			arr = append(arr, Ref{Num: fieldKey.Num, Gen: fieldKey.Gen})
		}
		u.Put(fv.Key(), arr)
		u.updated[fv.Key()] = arr
	case Array:
		arr := append(Array(nil), fv...)
		if add && !fieldKey.IsZero() {
			arr = append(arr, Ref{Num: fieldKey.Num, Gen: fieldKey.Gen})
		}
		af.Set("Fields", arr)
	default:
		arr := Array{}
		if add && !fieldKey.IsZero() {
			arr = append(arr, Ref{Num: fieldKey.Num, Gen: fieldKey.Gen})
		}
		af.Set("Fields", arr)
	}
	af.Set("SigFlags", Integer(3))
	return nil
}

// appendAnnotation appends the widget to the page's /Annots and re-emits that
// array direct. pdfbox does this deliberately (page.setAnnotations): an
// indirect /Annots that is not itself rewritten makes Adobe Reader report the
// document as modified.
func (u *Updater) appendAnnotation(pageKey, widgetKey ObjectKey) error {
	if pageKey.IsZero() {
		return nil
	}
	obj, err := u.Update(pageKey)
	if err != nil {
		return err
	}
	pd, ok := obj.(*Dict)
	if !ok {
		return errors.New("pdf: page object is not a dictionary")
	}
	var annots Array
	switch v := pd.GetRaw("Annots").(type) {
	case Array:
		annots = append(annots, v...)
	case Ref:
		if a, ok := u.doc.Resolve(v).(Array); ok {
			annots = append(annots, a...)
		}
	}
	annots = append(annots, Ref{Num: widgetKey.Num, Gen: widgetKey.Gen})
	pd.Set("Annots", annots)
	return nil
}

// generateFieldName reproduces PDSignatureField.generatePartialName: "Signature"
// followed by the smallest positive integer that no field's /T already uses.
func (u *Updater) generateFieldName() (string, error) {
	taken := map[string]bool{}
	cat, err := u.doc.Catalog()
	if err == nil && cat != nil {
		if af, ok := u.doc.GetDict(cat, "AcroForm"); ok {
			if fields, ok := u.doc.GetArray(af, "Fields"); ok {
				u.collectFieldNames(fields, taken, 0)
			}
		}
	}
	for i := 1; ; i++ {
		n := "Signature" + strconv.Itoa(i)
		if !taken[n] {
			return n, nil
		}
	}
}

func (u *Updater) collectFieldNames(fields Array, taken map[string]bool, depth int) {
	if depth > 64 {
		return
	}
	for _, e := range fields {
		d, ok := u.doc.Resolve(e).(*Dict)
		if !ok {
			continue
		}
		if t, ok := u.doc.GetString(d, "T"); ok {
			taken[string(t)] = true
		}
		if kids, ok := u.doc.GetArray(d, "Kids"); ok {
			u.collectFieldNames(kids, taken, depth+1)
		}
	}
}

// hasFilledSignature ports containsFilledSignature: a signature dictionary that
// carries a /ByteRange has been signed.
func (u *Updater) hasFilledSignature() (bool, error) {
	sigs, err := u.doc.SignatureDictionaries()
	if err != nil {
		return false, err
	}
	for _, s := range sigs {
		if s.Dict != nil && s.Dict.Has("ByteRange") {
			return true, nil
		}
	}
	return false, nil
}

// FormatDate renders t as a PDF date string "D:YYYYMMDDHHmmSS+HH'mm'", the
// format pdfbox's DateConverter.toString produces. A zero UTC offset is written
// "+00'00'", never "Z".
func FormatDate(t time.Time) string {
	_, offset := t.Zone()
	sign := byte('+')
	if offset < 0 {
		sign = '-'
		offset = -offset
	}
	return fmt.Sprintf("D:%04d%02d%02d%02d%02d%02d%c%02d'%02d'",
		t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(),
		sign, offset/3600, (offset%3600)/60)
}

// InsertContents writes cms as uppercase hex into the reserved /Contents span
// (R17), leaving the unused reserved bytes as ASCII '0'. It returns
// ErrContentsTooLarge when cms does not fit.
//
// Upstream is not self-consistent here and we deliberately are not either:
// COSWriter.writeExternalSignature uses pdfbox's Hex.getBytes (uppercase) while
// PAdESUtils.replaceSignature uses Utils.toHex (lowercase). Both are legal PDF
// hex strings, both parse, and each function below matches the upstream path it
// ports. See DESIGN.md §3.6.
func (r *Result) InsertContents(cms []byte) error {
	if r.ContentsLength <= 2 {
		return errors.New("pdf: no /Contents placeholder in this result")
	}
	if int64(2*len(cms)) > r.ContentsLength-2 {
		return fmt.Errorf("%w: %d hex bytes into %d reserved",
			ErrContentsTooLarge, 2*len(cms), r.ContentsLength-2)
	}
	p := r.ContentsOffset + 1 // just past the '<'
	for _, b := range cms {
		r.Bytes[p] = hexUpper[b>>4]
		r.Bytes[p+1] = hexUpper[b&0x0f]
		p += 2
	}
	return nil
}

const hexLower = "0123456789abcdef"

// ReplaceContents is the port of PAdESUtils.replaceSignature for the cached
// to-be-signed path: it finds the single all-zero hex placeholder in doc and
// substitutes cms, lowercase-hex encoded. It errors when zero or more than one
// placeholder is present.
//
// The scanner is upstream's, byte for byte: a '<' arms a suspicion, a run of
// exactly len(hex(cms)) '0' bytes after it is the placeholder, any other byte
// disarms it, and the reserved bytes beyond the CMS keep their '0's — so the
// document's length never changes.
func ReplaceContents(doc []byte, cms []byte) ([]byte, error) {
	if len(cms) == 0 {
		return nil, errors.New("pdf: cmsSignedData cannot be empty")
	}
	sig := make([]byte, 0, 2*len(cms))
	for _, b := range cms {
		sig = append(sig, hexLower[b>>4], hexLower[b&0x0f])
	}

	var out bytes.Buffer
	out.Grow(len(doc))
	var temp []byte
	suspicion, pasted := false, false
	for _, b := range doc {
		if suspicion {
			if b == '0' {
				temp = append(temp, b)
				if len(temp) == len(sig) {
					if pasted {
						return nil, errors.New("pdf: PDF document contains more than one empty signature!")
					}
					out.Write(sig)
					temp = nil
					suspicion = false
					pasted = true
				}
				continue
			}
			out.Write(temp)
			temp = nil
			suspicion = false
		}
		out.WriteByte(b)
		if b == '<' {
			temp = nil
			suspicion = true
		}
	}
	if suspicion {
		// Upstream drops a run that is still open at EOF; we flush it, because
		// silently truncating the tail of a document is never the better
		// failure mode. No PDF can reach EOF mid-run anyway: %%EOF follows.
		out.Write(temp)
	}
	if !pasted {
		return nil, errors.New("pdf: Reserved space to insert a signature was not found!")
	}
	return out.Bytes(), nil
}
