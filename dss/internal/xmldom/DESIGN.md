# XML stack binding design: `internal/xmldom` + `internal/xmlc14n`

**Status: BINDING.** Two implementers work in parallel against this document — one owns
`internal/xmldom`, one owns `internal/xmlc14n`. Every API in §1.9 and §2.9 is pinned. If an
implementer believes a pinned signature is wrong, they raise it with the tech lead and the doc
changes first; they do not change the code and reconcile later.

Upstream baseline: DSS **6.5.RC1** (`4c212986`), Apache Santuario **xmlsec 3.0.6** (the version
resolved in `~/.m2` for `dss-xml-utils`), OpenJDK 21 / Xerces as shipped in the JDK.

Sources read for this design (all verified, not recalled):

| What | Where |
|---|---|
| `XMLCanonicalizer`, `DomUtils`, `SantuarioInitializer`, `DocumentBuilderFactoryBuilder` | `/home/user/dss-upstream/dss-xml-utils`, `/home/user/dss-upstream/dss-xml-common` |
| `CanonicalizerBase`, `Canonicalizer20010315`, `Canonicalizer20010315Excl`, `CanonicalizerPhysical`, `NameSpaceSymbTable`, `XmlAttrStack`, `AttrCompare`, `C14nHelper`, `UtfHelpper`, `InclusiveNamespaces` | `xmlsec-3.0.6-sources.jar` (fetched to the maven cache; extracted at `scratchpad/xmlsec-src/`) |
| XAdES c14n call sites, ID registration, transform pipeline | `/home/user/dss-upstream/dss-xades` (`DSSXMLUtils`, `XAdESDOMDocument`, `reference/*`) |
| Live behaviour | `scratchpad/Probe.java`, `Probe2.java`, `Probe3.java`, `Probe4.java` (Santuario 3.0.6 oracle runs) and `scratchpad/xmlprobe/` (Go `encoding/xml` probes) |

Every "Santuario does X" claim below is backed by an executed probe, quoted inline.

---

## 0. Scope, layering and non-goals

```
internal/xmldom     minimal namespace-aware DOM + parser + serializer   (stdlib only)
internal/xmlc14n    C14N 1.0 / 1.1 / exclusive / physical               (stdlib + xmldom)
xmldsig  (phase 4b) Reference/Transform pipeline, resolvers, SignedInfo (+ spi, model)
xades    (phase 4c) XAdES B/T/LT/LTA                                    (+ everything)
```

Per `PORTING.md`, both `internal/` packages replace machinery that upstream gets from a third
party (Java DOM, Apache Santuario) and therefore have no Java class to mirror one-to-one. Each
carries a `doc.go` stating provenance. Neither imports a DSS package; `xmlc14n` imports only
`xmldom` and the standard library.

**`xmldom` is not a general XML library.** It exists to feed c14n, XML-DSig reference
processing and XAdES element construction. Explicit non-goals: DTD/entity declarations, schema
validation, XPath, XSLT, XInclude, namespace fixup on mutation, `Document.normalize()`,
live NodeLists, mutation events, `xml:space` interpretation, pretty-printing.

**One-canonicalization-per-call.** Santuario's inclusive canonicalizers are single-use: the
`firstCall` field of `Canonicalizer20010315` is never reset, so a reused instance silently drops
inherited namespaces and `xml:*` attributes (SANTUARIO-463). `Probe4` proves it:

```
reuse#0 => <p:c xmlns:p="urn:1" a="1" xml:lang="en"><t></t></p:c>
reuse#1 => <p:c a="1"><t></t></p:c>          <-- wrong
fresh#0 => <p:c xmlns:p="urn:1" a="1" xml:lang="en"><t></t></p:c>
fresh#1 => <p:c xmlns:p="urn:1" a="1" xml:lang="en"><t></t></p:c>
```

Every DSS call site is `XMLCanonicalizer.createInstance(m).canonicalize(...)` inline (verified by
grep over `dss-xades`), so DSS always gets the `fresh#` behaviour. **Our implementation is
stateless per call: all mutable state is created inside `Canonicalize` and dies with it.** The
KAT oracle must construct a fresh `Canonicalizer` per document (§3.4).

---

## 1. `internal/xmldom`

### 1.1 Why not `encoding/xml`'s object model, and what we do use

`encoding/xml` has no tree model at all, so the tree is ours. Its **tokenizer**, however, is
worth reusing, and its `RawToken()` mode is exactly right for c14n. Probed behaviour
(`scratchpad/xmlprobe`):

| Behaviour | `RawToken()` result | Verdict |
|---|---|---|
| Attribute order | preserved verbatim in `StartElement.Attr` | **good** — c14n needs document order for the "physical" method and for stable tie-breaks |
| Namespace declarations | delivered as ordinary attributes: `xmlns:z="u"` → `{Space:"xmlns", Local:"z"}`, `xmlns="u"` → `{Space:"", Local:"xmlns"}` | **good** — Java DOM models them the same way |
| Prefixes | `RawToken` leaves `Name.Space` as the **literal prefix** (`{"a","root"}` for `<a:root>`); `Token()` would rewrite it to the URI and destroy the prefix | **must use `RawToken`** — c14n emits QNames verbatim |
| Predefined entities | decoded (`&amp;` → `&`) | **good** — c14n re-escapes (§2.6) |
| Unknown entities | hard error `invalid character entity &foo;` | **good** — fail closed |
| DOCTYPE | surfaced as `xml.Directive`, never processed | we reject it (§1.5) |
| Line endings in character data | `\r\n` and lone `\r` → `\n` | **good** — XML §2.11, and c14n operates post-normalization |
| CDATA | delivered as `CharData`, one token per section, content verbatim | **good** — c14n replaces CDATA with escaped text |
| `]]>` outside CDATA | hard error | correct |
| `<` inside an attribute value | hard error | correct |
| `&#0;` and other illegal chars | hard error | correct |
| **Attribute-value normalization** | **NOT performed** | **trap — see §1.4** |
| Element start/end matching | **not** checked (`<r></s>` accepted) | we check (§1.5) |
| Single root element | **not** enforced; trailing content accepted | we enforce (§1.5) |
| Undeclared prefix | **not** rejected | we reject (§1.5) |
| `xmlns:xml="urn:bogus"`, `xmlns:p=""` | **accepted** | we reject (§1.5) — Xerces rejects both |
| Non-UTF-8 declared encoding | error unless `CharsetReader` set | we supply one (§1.6) |
| `Decoder.InputOffset()` | returns the exact end offset of the token just returned; for `<r/>` the `StartElement` span is the whole tag and the synthetic `EndElement` span is empty | **load-bearing** — this is how we recover raw attribute text (§1.4) |

**Decision D1.** The parser is `encoding/xml.Decoder` in `RawToken()` mode over an in-memory
`[]byte`, wrapped by our own namespace resolution, well-formedness checks and attribute-value
normalizer. We do not write a tokenizer from scratch, and we do not use `Token()`.

### 1.2 Node model

**Decision D2.** One concrete struct with a `Kind` discriminator — not an interface hierarchy.
Rationale: c14n walks by pointer and uses pointer identity as node-set membership; a single
struct makes `*Node` a stable, comparable, allocation-free identity, keeps parent/sibling links
uniform, and removes type-switch/interface overhead from the hot traversal loop.

```
Document ── children: ProcInst | Comment | Element (exactly one)
Element  ── Attrs: []*Node (Kind Attribute, DOCUMENT ORDER, namespace decls included)
         ── children: Element | Text | CDATA | Comment | ProcInst
```

Field usage per kind:

| Kind | `Name` | `Value` | `Attrs` | children |
|---|---|---|---|---|
| `Document` | zero | `""` | nil | prolog/epilog + root |
| `Element` | `{Space,Local,Prefix}` | `""` | attributes | yes |
| `Attribute` | `{Space,Local,Prefix}` | normalized value | nil | no |
| `Text` | zero | character data | nil | no |
| `CDATA` | zero | character data | nil | no |
| `Comment` | zero | comment data (no `<!--`/`-->`) | nil | no |
| `ProcInst` | `Local` = target | instruction data (may be `""`) | nil | no |

**Namespace declarations are attribute nodes**, modelled exactly as Java DOM does — this is what
makes Santuario's `AttrCompare` port trivially:

| Source | `Name.Space` | `Name.Local` | `Name.Prefix` | `QName()` |
|---|---|---|---|---|
| `xmlns:p="u"` | `http://www.w3.org/2000/xmlns/` | `p` | `xmlns` | `xmlns:p` |
| `xmlns="u"` | `http://www.w3.org/2000/xmlns/` | `xmlns` | `""` | `xmlns` |
| `p:a="v"` | resolved URI of `p` | `a` | `p` | `p:a` |
| `a="v"` | `""` (no namespace) | `a` | `""` | `a` |

Java DOM reports `namespaceURI == null` for an unprefixed attribute; we use `""`. `AttrCompare`
sorts null before any non-null URI, and `""` is lexicographically least among all strings, so the
two are equivalent — **no separate "has namespace" flag is needed**, and none shall be added.

`Text` vs `CDATA` is kept only for `Serialize` and for parity of node counts with Java DOM
(`setCoalescing(false)` is the JAXP default, so Xerces also produces one node per CDATA section
and one per contiguous character run). **c14n treats the two identically** — Santuario's
`canonicalizeSubTree` falls through `TEXT_NODE`/`CDATA_SECTION_NODE` to the same
`outputTextToWriter`. Probe: `<r><![CDATA[a < b & c]]></r>` → `<r>a &lt; b &amp; c</r>` under
every algorithm including physical.

### 1.3 Whitespace, CDATA and entity preservation

- **Nothing is ever stripped.** No whitespace trimming, no "ignorable whitespace" concept (that
  needs a DTD/schema, which we do not have), no coalescing of adjacent text nodes.
- `xml:space` is **stored as an ordinary attribute and never interpreted**. It matters only as an
  inherited `xml:*` attribute in C14N 1.0/1.1 (§2.5).
- CDATA text is stored **exactly as it appeared** between `<![CDATA[` and `]]>`, after XML §2.11
  line-ending normalization only. It is not escaped, not unescaped, not re-wrapped.
- The five predefined entities are decoded at parse time; c14n re-escapes (§2.6). A character
  reference and the character it denotes are indistinguishable in text nodes **and this is
  correct**: literal CR is normalized to LF by §2.11, so a `\r` surviving in a text node can only
  have come from `&#xD;`/`&#13;`, which c14n must emit as `&#xD;`. Probe: `<r>a&#13;b</r>` →
  `<r>a&#xD;b</r>`. In **attribute** values the same reasoning does not hold — see §1.4.

### 1.4 Attribute-value normalization — the one place `encoding/xml` is wrong

XML 1.0 §3.3.3 requires, for a CDATA-typed attribute (all of ours, since DTDs are banned):

1. line-ending normalization (§2.11) is applied first;
2. each **literal** `#x20`, `#x9`, `#xA`, `#xD` in the value becomes `#x20`;
3. characters produced by a **character or entity reference** are appended as-is and are *not*
   whitespace-normalized.

Go performs step 1 but **not** step 2:

```
input : <r a="x<TAB>y<LF>z<CRLF>w<CR>v" b="&#9;&#10;&#13;"/>
Go    : a = "x\ty\nz\nw\nv"      b = "\t\n\r"
Java  : a = "x y z w v"          b = "\t\n\r"
```

and Santuario therefore emits (probe, all seven algorithms agree):

```
<r a="x y z" b="&#x9;&#xA;&#xD;"></r>
```

Because Go decodes references in place, the literal-vs-reference distinction is destroyed by the
time `Attr.Value` is available. **`xml.Attr.Value` is unusable for c14n and must never be read.**

**Decision D3.** The parser keeps the whole document in memory and, for each `StartElement`,
recovers the raw byte span of the start tag via `InputOffset()` bracketing:

```go
start := d.InputOffset()
tok, err := d.RawToken()          // xml.StartElement
end := d.InputOffset()            // src[start:end] is exactly "<a:r ... >" or "<a:r ... />"
```

It then re-scans `src[start:end]` with a small dedicated scanner — `Name (S Eq AttValue)* S? '/'? '>'`,
no nesting, no CDATA, no comments — and computes each attribute value by:

- normalizing `\r\n` and lone `\r` to `\n` over the raw value;
- walking the value: `&#...;` / `&#x...;` / `&amp;` / `&lt;` / `&gt;` / `&quot;` / `&apos;` →
  append the referenced character **verbatim**; any other `&` → error; literal `\t`, `\n`
  (and any `\r` that survived, which cannot) → append `#x20`; everything else → append verbatim.

The scanner uses the token's `Attr` slice only to cross-check count and QName ordering (a cheap
self-test that the two agree); on mismatch it returns an internal error rather than guessing.
Character references to non-XML-1.0 characters and to surrogate code points are rejected, matching
Xerces.

Cost: one extra pass over start tags. Benefit: exact §3.3.3 semantics, and the raw span is also
what lets `Serialize` and future work reproduce source detail. This is the single most important
correctness decision in the parser; it is not optional and not deferrable.

### 1.5 Well-formedness and security checks the parser adds

Go's `RawToken` is deliberately permissive. `xmldom.Parse` rejects, with a `*SyntaxError`
carrying line/column:

1. **Any `DOCTYPE`** (`xml.Directive` beginning `DOCTYPE`) unless `ParseOptions.AllowDoctype`,
   which nothing in DSS sets. This replicates
   `DocumentBuilderFactoryBuilder`'s `disallow-doctype-decl=true` and Santuario's
   `XMLParserImpl.createDocumentBuilder(disallowDocTypeDeclarations=true)`. Probe:
   `<!DOCTYPE r><r/>` → `XMLParserException: Error parsing the inputstream`. Consequence: no
   external entities, no parameter entities, no internal subset, no attribute defaulting, no
   DTD-declared ID attributes, no billion-laughs. `AllowDoctype` still never *processes* the
   subset — it only tolerates the declaration — so entity references other than the five
   predefined ones remain a hard error either way.
2. Start/end tag mismatch, and unclosed elements at EOF.
3. More or fewer than exactly one element child of `Document`.
4. Any non-whitespace character data outside the document element.
5. Duplicate attributes — both by raw QName and by resolved `(Space, Local)`.
6. Undeclared prefix on an element or attribute name.
7. `xmlns:p=""` (empty binding for a prefixed declaration — forbidden by Namespaces in XML 1.0
   3rd ed.). Xerces: *"Prefixed namespace bindings may not be empty"*; Go accepts it silently.
8. `xmlns:xml` bound to anything other than `http://www.w3.org/XML/1998/namespace`; any
   declaration of the `xmlns` prefix; any binding of another prefix to the `xmlns` URI.
9. `<?xml version="1.1"?>`. **Decision D4:** XML 1.1 is rejected. Xerces accepts it and then
   applies XML 1.1 line-ending normalization (NEL `U+0085`, LSEP `U+2028`), which would silently
   diverge from us. Nothing in XAdES, ETSI TS 119 612 trusted lists or eIDAS profiles uses XML
   1.1. Failing closed beats diverging. Recorded as an accepted gap in `PORTING_PLAN.md`.
   (Under XML 1.0, NEL and LSEP are *not* normalized — probe confirms `<r>a<U+0085>b</r>` and
   `<r>a<U+2028>b</r>` survive c14n unchanged.)
10. Nesting deeper than `ParseOptions.MaxDepth` (default 500) and documents larger than
    `ParseOptions.MaxBytes` (default 0 = unlimited; callers on untrusted input set it).

The XML declaration itself is **not** a node: a leading `ProcInst` with target `xml` at offset 0
is consumed for its `version`/`encoding`/`standalone` pseudo-attributes, exactly as DOM does. It
is not discarded outright, though: the encoding and standalone flag are kept on the `Document`
and readable through `XMLEncoding`/`XMLStandalone`, because `Document.getXmlEncoding()` and
`getXmlStandalone()` are what `DomUtils.serializeNode` and `DOM2TO` read back to decide the
declaration they write and the charset they write it in (§1.8). A leading UTF-8 BOM is consumed
and discarded (Go surfaces it as `CharData "﻿"`).
Probe confirms Santuario emits nothing for the declaration or the BOM.

### 1.6 Encoding

Java honours the encoding declaration: feeding UTF-8 bytes under
`<?xml version="1.0" encoding="ISO-8859-1"?>` produces mojibake (`aé` → `aÃ©`), proving
the declaration wins. We match.

**Decision D5.** `xmldom` ships a built-in decoder table covering `UTF-8`, `US-ASCII`,
`ISO-8859-1` / `latin1` (byte → rune) and, by BOM detection, `UTF-16BE`/`UTF-16LE`
(`unicode/utf16`), all matched case-insensitively. Anything else is an error unless the caller
supplies `ParseOptions.CharsetReader`. No `golang.org/x/text` dependency.

### 1.7 ID attributes

DSS registers IDs itself because, with DTDs banned, the parser knows no attribute types.
`XAdESDOMDocument.setIDIdentifier` (and the deprecated `DSSXMLUtils.setIDIdentifier`) walks the
attributes of each element in document order and, for the **first** attribute whose *local name*
equals `"Id"` **case-insensitively**, calls `setIdAttribute(nodeName, true)` and `break`s.
`recursiveIdBrowse` applies this to every element from the document element down.

**Decision D6.** `xmldom` reproduces that rule exactly and adds `xml:id`:

- `RegisterIDs()` — DSS parity: per element, in `Attrs` order, the first attribute with
  `strings.EqualFold(Local, "Id")` is registered; then stop for that element. This matches `Id`,
  `ID`, `id`, `iD`, prefixed or not.
- Additionally, an attribute with `Space == XMLNamespace && Local == "id"` is always registered
  (the `xml:id` Recommendation; Xerces honours it independently of the DSS scan). If an element
  has both, the DSS-rule attribute wins for `ElementByID`, and both are reported by `IDAttrs`.
- The index is a `map[string]*Node` on the `Document` node. **First registration in document
  order wins**; later collisions are recorded and returned by `DuplicateIDs()`, mirroring
  `DSSXMLUtils.isDuplicateIdsDetected`. Registration never errors — the caller decides whether a
  duplicate is fatal (XAdES validation treats it as an attack indicator).
- The index is invalidated (lazily rebuilt) by any mutation of an element's attributes or by
  tree-structure changes. Implementers: bump a `gen` counter on the owning document in
  `SetAttr`/`RemoveAttr`/`AppendChild`/`InsertBefore`/`RemoveChild`/`ReplaceChild`.

`ElementByID` is what the phase-4b `SameDocumentResolver` uses for `#foo` references, and what
`#xpointer(id('foo'))` reduces to.

### 1.8 Serialization

`Serialize` is **not** c14n and must never be used where c14n is required: it keeps prefixes and
QNames verbatim from `Name`, does not sort attributes into c14n order, and does not import an
ancestor's namespace axis the way c14n does. Anything covered by a `ds:Transform` goes through
`internal/xmlc14n`.

Byte-parity with Java's identity `Transformer` **is** a requirement — an earlier revision of this
note said it was not, on the grounds that `DomUtils.serializeNode` output is "never itself
signed". That is false. `DSSXMLUtils.applyTransforms(Node, List<DSSTransform>)` (dss-xades,
~line 1049) returns `DomUtils.getNodeBytes(node)` when the reference carries no transforms, and
that byte array *is* what the `ds:Reference` `DigestValue` is computed over. Audit item D6.

`Serialize` therefore reproduces the two JDK classes the Transformer runs — `DOM2TO` (the DOM
walk and namespace fixup) and `ToStream`/`ToXMLStream` (the bytes). The consequences that most
often surprise:

- The XML declaration is written for **every** node kind, not just `Document`: `DomUtils` never
  sets `OMIT_XML_DECLARATION`. `standalone="no"` is added only when the serialized node *is* the
  `Document` and the source did not declare `standalone="yes"`.
- Empty elements are written `<e/>`.
- Attributes come out in Xerces' `NamedNodeMap` order — sorted by node name, **not** source order
  — with namespace declarations moved in front, and a fixup declaration for the element's own
  prefix appended last when nothing else declared it. This is what lets a subtree serialized on
  its own regain the declarations it inherited.
- A declaration that rebinds a prefix to the URI it already has, and any declaration of a prefix
  starting with `xml`, are dropped.
- Text escapes `&`, `<`, `>` and writes CR as `&#13;`; attribute values additionally escape `"`
  and write TAB/LF/CR as `&#9;`/`&#10;`/`&#13;`. Character references are **decimal**.
- C0/C1 controls are character references in text and are written literally in attribute values;
  supplementary characters are character references in both, but literal inside `CDATA`.
- The declared encoding decides the declaration text, which characters are literal, and the
  charset of the bytes; `xmldom` keeps it (and `standalone`) on the `Document` for that reason.

`SerializeOptions.XMLDeclaration` (default true) is `OMIT_XML_DECLARATION`; `.Encoding` (default:
the document's own) is `OutputKeys.ENCODING`.

The gate is `xml/utils.TestSerializeAgainstJavaTransformerOracle`: an adversarial corpus
(`xml/utils/testdata/serialize/corpus.txt`) replayed through the real `DomUtils.serializeNode`
and `DomUtils.getNodeBytes` on OpenJDK 21 by `testdata/gen/SerializeOracle.java`, demanding byte
equality on both.

### 1.9 Exported API — **PINNED**

```go
// Package xmldom is a minimal, namespace-aware XML document model sufficient for
// XML-DSig canonicalization, reference processing and XAdES construction. It replaces
// the org.w3c.dom API that upstream DSS obtains from the JDK. It is not a general
// purpose XML library: no DTDs, no schema validation, no XPath, no XSLT.
package xmldom

import "io"

// ---------------------------------------------------------------- namespaces

const (
	XMLNamespace   = "http://www.w3.org/XML/1998/namespace" // the "xml" prefix
	XMLNSNamespace = "http://www.w3.org/2000/xmlns/"        // namespace declarations
)

// Name is an expanded name plus the literal prefix as it appeared in the source.
// Canonicalization emits QNames verbatim, so Prefix is load-bearing and must never
// be regenerated. Space is "" for a name in no namespace.
type Name struct {
	Space  string
	Local  string
	Prefix string
}

// QName returns "prefix:local", or "local" when Prefix is empty.
func (n Name) QName() string

// ---------------------------------------------------------------- node model

type Kind uint8

const (
	Document Kind = iota + 1
	Element
	Attribute
	Text
	CDATA
	Comment
	ProcInst
)

func (k Kind) String() string

// Node is every node type. Pointer identity is node identity: NodeSet membership,
// subset visibility and Exclude comparisons are all pointer comparisons.
//
// Attrs holds an element's attributes, including namespace declarations, in DOCUMENT
// ORDER. Callers must not reorder or mutate it directly; use SetAttr/RemoveAttr.
type Node struct {
	Kind  Kind
	Name  Name   // Element, Attribute; ProcInst target in Name.Local
	Value string // Attribute value, Text/CDATA data, Comment data, ProcInst data
	Attrs []*Node

	Parent      *Node
	FirstChild  *Node
	LastChild   *Node
	PrevSibling *Node
	NextSibling *Node
	// contains filtered or unexported fields
}

// ---------------------------------------------------------------- parsing

type ParseOptions struct {
	AllowDoctype  bool // default false: any DOCTYPE declaration is a SyntaxError
	MaxDepth      int  // default 500 when zero
	MaxBytes      int64 // default 0 = unlimited
	CharsetReader func(charset string, input io.Reader) (io.Reader, error)
}

// Parse builds a document from src. opts may be nil, which selects the DSS-equivalent
// secure defaults (no DOCTYPE, no external anything, UTF-8/US-ASCII/ISO-8859-1/UTF-16).
// The returned node has Kind Document.
func Parse(src []byte, opts *ParseOptions) (*Node, error)

// ParseReader reads r to completion and calls Parse.
func ParseReader(r io.Reader, opts *ParseOptions) (*Node, error)

type SyntaxError struct {
	Line, Column int
	Offset       int64
	Msg          string
}

func (e *SyntaxError) Error() string

// ---------------------------------------------------------------- construction

func NewDocument() *Node
func NewElement(name Name) *Node
func NewAttr(name Name, value string) *Node
func NewText(data string) *Node
func NewCDATA(data string) *Node
func NewComment(data string) *Node
func NewProcInst(target, data string) *Node

// ---------------------------------------------------------------- tree mutation

// AppendChild appends c to n and returns c. It panics if c already has a parent or if
// the parent/child kind combination is illegal.
func (n *Node) AppendChild(c *Node) *Node

// InsertBefore inserts c before ref (a child of n) and returns c. A nil ref appends.
func (n *Node) InsertBefore(c, ref *Node) *Node

// RemoveChild detaches c from n and returns c.
func (n *Node) RemoveChild(c *Node) *Node

// ReplaceChild puts newC where oldC was and returns oldC.
func (n *Node) ReplaceChild(newC, oldC *Node) *Node

// Clone deep- or shallow-copies n. The copy has no parent and no document.
func (n *Node) Clone(deep bool) *Node

// Import returns a deep or shallow copy of src owned by n's document, ready to insert.
func (n *Node) Import(src *Node, deep bool) *Node

// ---------------------------------------------------------------- attributes

// Attr returns the attribute node with the given namespace URI and local name.
// Namespace declarations are found with space == XMLNSNamespace; the default
// declaration xmlns="..." has local name "xmlns".
func (n *Node) Attr(space, local string) *Node

// AttrValue is Attr(...).Value, or "" when absent.
func (n *Node) AttrValue(space, local string) string

// SetAttr sets or replaces an attribute, preserving its position in Attrs when it
// already exists and appending otherwise. It returns the attribute node.
func (n *Node) SetAttr(name Name, value string) *Node

// RemoveAttr removes an attribute and reports whether one was removed.
func (n *Node) RemoveAttr(space, local string) bool

// ---------------------------------------------------------------- navigation

// Document returns the owning document node, or nil for a detached subtree.
func (n *Node) Document() *Node

// DocumentElement returns the single element child of a document node, else nil.
func (n *Node) DocumentElement() *Node

// Children returns a snapshot of n's children.
func (n *Node) Children() []*Node

// Elements returns a snapshot of n's element children.
func (n *Node) Elements() []*Node

// FirstElementChild returns the first element child, or nil.
func (n *Node) FirstElementChild() *Node

// NextElementSibling returns the next element sibling, or nil.
func (n *Node) NextElementSibling() *Node

// Ancestors returns n's ancestors, nearest first, excluding n.
func (n *Node) Ancestors() []*Node

// Depth returns the number of ancestors of n; a document node has depth 0.
func (n *Node) Depth() int

// Contains reports whether other is n or a descendant of n.
func (n *Node) Contains(other *Node) bool

// TextContent concatenates the data of all Text and CDATA descendants in document order.
func (n *Node) TextContent() string

// SetTextContent replaces all children of an element with a single text node.
func (n *Node) SetTextContent(data string)

// Walk calls fn for n and every descendant in document order, attributes excluded.
// Returning false from fn skips that node's subtree.
func (n *Node) Walk(fn func(*Node) bool)

// LookupNamespaceURI resolves prefix ("" means the default namespace) against the
// declarations in scope at n. It reports whether a binding was found.
func (n *Node) LookupNamespaceURI(prefix string) (string, bool)

// LookupPrefix returns the innermost prefix bound to uri that is in scope at n.
func (n *Node) LookupPrefix(uri string) (string, bool)

// ---------------------------------------------------------------- node sets

// NodeSet is a document-subset membership set for canonicalization. Namespace
// declarations participate as their attribute nodes, matching Santuario.
type NodeSet map[*Node]struct{}

func NewNodeSet(nodes ...*Node) NodeSet
func (s NodeSet) Add(nodes ...*Node)
func (s NodeSet) Remove(nodes ...*Node)
func (s NodeSet) Has(n *Node) bool
func (s NodeSet) Len() int

// AddSubtree adds n, all its descendants, and all their attributes to s.
func (s NodeSet) AddSubtree(n *Node)

// ---------------------------------------------------------------- ID attributes

// RegisterIDs indexes ID attributes over the whole document, reproducing
// XAdESDOMDocument.recursiveIdBrowse: per element, in attribute order, the first
// attribute whose local name equals "Id" case-insensitively; plus any xml:id.
// It must be called on a document node.
func (n *Node) RegisterIDs()

// RegisterIDAttr indexes a single attribute of elem as an ID attribute.
func (n *Node) RegisterIDAttr(elem, attr *Node)

// ElementByID returns the element carrying the given ID value, or nil. A leading '#'
// is not accepted; strip the fragment delimiter before calling.
func (n *Node) ElementByID(id string) *Node

// IDAttrs returns the ID attribute nodes registered for elem, in attribute order.
func (n *Node) IDAttrs(elem *Node) []*Node

// DuplicateIDs returns the ID values registered more than once, sorted.
func (n *Node) DuplicateIDs() []string

// ---------------------------------------------------------------- serialization

type SerializeOptions struct {
	XMLDeclaration bool   // OMIT_XML_DECLARATION; default true via nil options
	Encoding       string // OutputKeys.ENCODING; default: the document's own, else UTF-8
}

// Serialize writes n as XML, byte for byte as OpenJDK's identity Transformer does -
// which is what DomUtils.serializeNode runs, and what a no-transform ds:Reference
// digests. It is still NOT canonicalization: use internal/xmlc14n for a ds:Transform.
func (n *Node) Serialize(w io.Writer, opts *SerializeOptions) error

// Bytes is Serialize into a buffer.
func (n *Node) Bytes(opts *SerializeOptions) ([]byte, error)

// XMLEncoding and XMLStandalone report the parsed source's XML declaration, which the
// serializer needs: Document.getXmlEncoding / getXmlStandalone.
func (n *Node) XMLEncoding() string
func (n *Node) XMLStandalone() bool
```

---

## 2. `internal/xmlc14n`

### 2.1 Algorithms and their upstream registration

`XMLCanonicalizer.registerDefaultCanonicalizers()` registers exactly these seven, and
`XMLCanonicalizer` refuses anything else with `IllegalArgumentException`:

| URI | Constant | Santuario class |
|---|---|---|
| `http://www.w3.org/TR/2001/REC-xml-c14n-20010315` | `C14N10` | `Canonicalizer20010315OmitComments` |
| `…-20010315#WithComments` | `C14N10WithComments` | `Canonicalizer20010315WithComments` |
| `http://www.w3.org/2006/12/xml-c14n11` | `C14N11` | `Canonicalizer11_OmitComments` |
| `http://www.w3.org/2006/12/xml-c14n11#WithComments` | `C14N11WithComments` | `Canonicalizer11_WithComments` |
| `http://www.w3.org/2001/10/xml-exc-c14n#` | `C14NExclusive` | `Canonicalizer20010315ExclOmitComments` |
| `http://www.w3.org/2001/10/xml-exc-c14n#WithComments` | `C14NExclusiveWithComments` | `Canonicalizer20010315ExclWithComments` |
| `http://santuario.apache.org/c14n/physical` | `C14NPhysical` | `CanonicalizerPhysical` |

Defaults, copied verbatim from `XMLCanonicalizer`:
`DefaultDSS = C14NExclusive` (`DEFAULT_DSS_C14N_METHOD`),
`DefaultXMLDSig = C14N10` (`DEFAULT_XMLDSIG_C14N_METHOD`, used when a signature omits the method,
DSS-2208). An empty algorithm string resolves to `DefaultXMLDSig`.

### 2.2 Internal architecture

```
xmlc14n
├── algorithm.go     Algorithm type, URI constants, registry, Supported()
├── c14n.go          Canonicalize / CanonicalizeToBytes / CanonicalizeBytes; Input; errors
├── walk.go          the shared traversal (port of CanonicalizerBase.canonicalizeSubTree
│                    and canonicalizeXPathNodeSet) parameterized by an attrEmitter
├── nsstack.go       nsStack — port of NameSpaceSymbTable (+ NameSpaceSymbEntry)
├── xmlattrs.go      xmlAttrStack — port of XmlAttrStack, incl. the 1.1 xml:base join
├── emit_inclusive.go  1.0 and 1.1 attribute emission (Canonicalizer20010315)
├── emit_exclusive.go  exclusive attribute emission (Canonicalizer20010315Excl)
├── emit_physical.go   physical attribute emission (CanonicalizerPhysical)
├── attrsort.go      attrLess — port of AttrCompare
└── escape.go        the three escaping tables + UTF-8 writer
```

The traversal is written **once**. The three emitters differ only in
`emitAttrsSubtree(el, ns, w)` and `emitAttrsSubset(el, ns, w)`, and in three behavioural flags:

```go
type engine struct {
	includeComments   bool
	c14n11            bool
	physical          bool   // suppresses positional newlines around prolog/epilog nodes
	inclusivePrefixes map[string]struct{} // exclusive only
}
```

`nsStack` is a faithful port of `NameSpaceSymbTable`, including the `lastrendered` field — that
field, not any spec text, is what suppresses superfluous declarations and the redundant `xmlns=""`
(worked example in §2.4). Do **not** "simplify" it. Santuario's `SymbMap` open-addressing table is
an optimization only; a plain `map[string]*nsEntry` with an explicit undo journal per frame is the
required Go shape (deterministic, no hash-order dependence: `getUnrenderedNodes` feeds a sorted
set, so map iteration order never reaches the output).

### 2.3 Document-subset model: subtree vs node-set

Santuario has two entry points and they are **not** the same algorithm:

- **Subtree** (`engineCanonicalizeSubTree`): the node set is "the apex plus everything below it".
  Ancestor namespace and `xml:*` context is gathered up-front by `getParentNameSpaces` →
  `handleParent`, and flushed onto the apex element by the `firstCall` branch. This is what
  `XMLCanonicalizer.canonicalize(Node)` and `canonicalize(byte[])` use — i.e. **everything DSS
  does directly**.
- **Node-set** (`engineCanonicalizeXPathNodeSet`): arbitrary membership, elements can be
  "not visible" while their descendants are, `xmlns=""` can have to be synthesized, and
  `xml:*` attributes float up from unselected ancestors. Reached only through Santuario's
  Transform pipeline (XPath / XPath2 filter / enveloped-signature transforms), which DSS drives
  via `DSSXMLUtils.applyTransforms` → `ReferenceProcessor`.

**Decision D7.** `xmlc14n` implements **both** from day one — the API and the `walk.go`
parameterization are designed for it — but the phase-4a golden gate covers **subtree only**.
Node-set KATs are generated in phase 4b through the real transform pipeline, because building a
Java-side node set by hand runs into Xerces' namespace-node identity problem (the very issue
Santuario's `circumventBug2650` exists for) and would pin an artefact of the test harness rather
than of DSS. `Input.Subset == nil` selects subtree mode; non-nil selects node-set mode.

`Input.Exclude` ports `engineCanonicalizeSubTree(root, excludeNode, writer)`: the named element
and its subtree are skipped. Used by the enveloped-signature transform.

`propagateDefaultNamespace` (exclusive only, Santuario ≥ 2.2, for XML Encryption) is **not**
implemented; DSS never sets it. `Canonicalize` has no knob for it.

### 2.4 Namespace rules

**Inclusive 1.0/1.1 (subtree).** The apex inherits every prefix binding in scope from its
ancestors; each element then emits only the bindings that are *newly rendered* at that element.
`addMappingAndRender` returns nil — meaning "emit nothing" — when the prefix is already bound to
the same URI *and already rendered*, which is the superfluous-declaration rule:

```
<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"><p:d xmlns:p="urn:1"/></p:c></r>
C14N10 → <r xmlns:p="urn:1"><p:c><p:d></p:d></p:c></r>
```

Unused declarations **are** emitted by 1.0 and 1.1 (that is the whole point of exclusive c14n):

```
<r xmlns:unused="urn:u"><c/></r>
C14N10 → <r xmlns:unused="urn:u"><c></c></r>
EXCL   → <r><c></c></r>
```

`xmlns:xml="http://www.w3.org/XML/1998/namespace"` is never emitted by 1.0/1.1/exclusive
(`!(XML.equals(NName) && XML_LANG_URI.equals(NValue))` guards every path); physical does emit it.
Probe confirms both.

**Default-namespace undeclaration (`xmlns=""`).** Emitted when, and only when, the default
namespace is non-empty in the rendered context and the element undeclares it:

```
<r xmlns="urn:a"><c xmlns=""><d/></c></r>          → <r xmlns="urn:a"><c xmlns=""><d></d></c></r>
```

But when the *apex of the subtree* is the undeclaring element, no `xmlns=""` is emitted, because
the initial `nsStack` entry for `xmlns` has `uri=""` and `lastrendered=""`, so
`addMappingAndRender("xmlns","",…)` finds `ob.lastrendered == uri` and returns nil:

```
<r xmlns="urn:d"><c xmlns=""><t/></c></r>, apex = <c>   → <c><t></t></c>
```

This is spec-correct (the apex inherits nothing, so the empty default namespace is already in
force) and is exactly the mechanism the `lastrendered` field exists for. The mirror case in
`getParentNameSpaces` — an ancestor chain whose innermost default binding is `""` — is handled by
the explicit `nullNode` injection at the end of `getParentNameSpaces`.

**Exclusive 1.0.** No inheritance. At each *output* element the emitter computes the
**visibly-utilized** prefix set:

- the prefix of the element name, or `"xmlns"` when the element name is unprefixed
  (`prefix = XMLNS` in `Canonicalizer20010315Excl`);
- the prefix of every *output* attribute, excluding `xml` and `xmlns`;
- plus every entry of the `InclusiveNamespaces PrefixList`.

Each visibly-utilized prefix is looked up with `ns.getMapping(prefix)`, which returns nil when
already rendered with the same URI, and the survivors are emitted. Note that **an unprefixed
attribute does not utilize the default namespace** (unprefixed attributes are in no namespace),
which is why:

```
<r xmlns:p="urn:1" xmlns:q="urn:2" xmlns="urn:d"><p:c q:a="1" b="2"/></r>, apex = <c>
EXCL → <p:c xmlns:p="urn:1" xmlns:q="urn:2" b="2" q:a="1"></p:c>
```

and why the default declaration migrates down to the element that actually uses it:

```
<r xmlns:p="urn:1" xmlns="urn:d" xmlns:un="urn:u"><p:c a="1"><t/></p:c></r>, apex = <c>
EXCL   → <p:c xmlns:p="urn:1" a="1"><t xmlns="urn:d"></t></p:c>
C14N10 → <p:c xmlns="urn:d" xmlns:p="urn:1" xmlns:un="urn:u" a="1"><t></t></p:c>
```

**PrefixList.** `InclusiveNamespaces.prefixStr2Set`: split on `\s`, map `#default` → `xmlns`,
collect into a sorted set; null or empty string → empty set. Unknown prefixes in the list are
silently ignored (`ns.getMapping` returns nil). Probes:

```
incl="q"        <r xmlns="urn:d" xmlns:p="urn:1" xmlns:q="urn:2"><p:c/></r>  apex=<c>
                EXCL → <p:c xmlns:p="urn:1" xmlns:q="urn:2"></p:c>
incl="#default q"  same doc, <p:c><t2/></p:c>
                EXCL → <p:c xmlns="urn:d" xmlns:p="urn:1" xmlns:q="urn:2"><t2></t2></p:c>
incl="zz"       <r xmlns:p="urn:1"><p:c/></r>  apex=<c>
                EXCL → <p:c xmlns:p="urn:1"></p:c>          (no effect)
```

**Relative namespace URIs.** `C14nHelper.namespaceIsAbsolute` = `value == "" || indexOf(':') > 0`.
Anything else raises `CanonicalizationException("Element r has a relative namespace: p=...")` for
1.0, 1.1 and exclusive. **Physical accepts it.** Port as `*RelativeNamespaceError`. Note the
check fires *only when the declaration is actually rendered*, so an unused relative declaration
survives exclusive c14n silently — reproduce that, do not "improve" it.

**Physical.** No inheritance, no suppression, no relative-URI check: every attribute physically
present on the element is emitted, sorted by the same comparator. Probes:

```
<r xmlns:p="urn:1"><p:c xmlns:p="urn:1"><p:d xmlns:p="urn:1"/></p:c></r>
PHYS → <r xmlns:p="urn:1"><p:c xmlns:p="urn:1"><p:d xmlns:p="urn:1"></p:d></p:c></r>
apex=<c> of <r xmlns:p="urn:1" xmlns="urn:d"><p:c a="1"><t/></p:c></r>
PHYS → <p:c a="1"><t></t></p:c>          (undeclared prefix in the output — by design)
```

### 2.5 `xml:*` attribute inheritance — 1.0 vs 1.1

Only **inclusive** c14n inherits `xml:*` attributes from unselected ancestors; exclusive and
physical never do (`<r xml:space="preserve"><t a="1"/></r>`, apex `<t>`: `C14N10` →
`<t a="1" xml:space="preserve"></t>`, `EXCL`/`PHYS` → `<t a="1"></t>`).

`XmlAttrStack` collects the ancestors' `xml:*` attributes and flushes them onto the first output
element. Rules, ported exactly:

| | C14N 1.0 | C14N 1.1 |
|---|---|---|
| `xml:lang`, `xml:space` | nearest-in-source-iteration-order wins; `loa` is keyed by QName and the **first** occurrence encountered wins | identical |
| `xml:id` | inherited like any other `xml:*` attribute | **never inherited**; treated as an ordinary attribute (emitted only when the element itself carries it) |
| `xml:base` | first occurrence wins (no joining) | the chain of unrendered ancestors' `xml:base` values is **joined** with `XmlAttrStack.joinURI` |
| omitted-ancestor tracking | none | `successiveOmitted` stops the walk at the first already-rendered level |

Probes that pin the two divergences:

```
<r xml:id="top" xml:lang="en"><m xml:id="mid"><t/></m></r>, apex = <t>
C14N10 → <t xml:id="top" xml:lang="en"></t>
C14N11 → <t xml:lang="en"></t>

<r xml:base="b/"><m xml:base="c/"><t/></m></r>, apex = <t>
C14N10 → <t xml:base="b/"></t>
C14N11 → <t xml:base="c/b/"></t>

<r xml:base="http://x/a/"><m xml:base="b/"><n xml:base="c/"><t/></n></m></r>, apex = <t>
C14N10 → <t xml:base="http://x/a/"></t>
C14N11 → <t xml:base="http://x/a/"></t>       (outer absolute base wins)
```

**`joinURI` is Santuario's, not RFC 3986's, and the argument order is counter-intuitive**: the
outermost base is seeded first, and each further-in ancestor is folded in as
`base = joinURI(inner.value, base)` — i.e. the *accumulated* value is passed as the
`relativeURI` parameter. Port `joinURI` and `removeDotSegments` **line for line**, including the
`//`-collapsing loop, the `endsWith("..")` fix-ups and the `URI(scheme, authority, path, query, null)`
recomposition. Do not substitute `net/url.ResolveReference`; it produces different answers.
`URISyntaxException` is swallowed (`LOG.debug`) and the previous value is kept — reproduce that,
including keeping the *previous* accumulated value on failure. When the joined result is empty,
no `xml:base` is emitted.

An ordinary-attribute detail from the 1.1 branch of `outputAttributes`: `xml:id` on a *visible*
element is added to the result set like any other attribute, whereas every other `xml:*`
attribute goes to the `xmlAttrStack`. Probe: `<r><t xml:id="me"/></r>`, apex `<t>` →
`<t xml:id="me"></t>` under all four algorithms.

### 2.6 Escaping tables — exact, per context

Ported from `CanonicalizerBase`. Uppercase hex digits, lowercase `x`, no leading zeroes.

| Character | Attribute value | Text / CDATA | Comment | PI target | PI data |
|---|---|---|---|---|---|
| `&` | `&amp;` | `&amp;` | — | — | — |
| `<` | `&lt;` | `&lt;` | — | — | — |
| `>` | — | `&gt;` | — | — | — |
| `"` | `&quot;` | — | — | — | — |
| `#x9` | `&#x9;` | — | — | — | — |
| `#xA` | `&#xA;` | — | — | — | — |
| `#xD` | `&#xD;` | `&#xD;` | `&#xD;` | `&#xD;` | `&#xD;` |

Everything else is written as UTF-8. `'` is never escaped. Attribute values are always delimited
by `"`.

Comment and PI content is otherwise **verbatim**, and note *why*: comments and PIs have no markup
or reference recognition, so `<r><!--a&#13;b--></r>` puts the six literal characters `&#13;` in
the comment data and c14n copies them through — `<r><!--a&#13;b--></r>` under `#WithComments`.
The only way a CR reaches comment/PI data is a literal CR in the source, which XML §2.11 turns
into LF; the `&#xD;` column above is therefore reachable only through DOM construction, not
through parsing. Contrast the text case, where references *are* recognized:
`<r>a&#13;b</r>` → `<r>a&#xD;b</r>`. The escaper always emits `&#xD;`, never `&#13;`.

Empty elements are always written as a start tag followed by an end tag: `<r/>` → `<r></r>`.

### 2.7 Attribute ordering

Port of `AttrCompare`, applied to the *emitted* set (not to `Attrs`):

1. namespace declarations (`Space == XMLNSNamespace`) sort before all other attributes;
2. among namespace declarations: by local name, with `"xmlns"` (the default declaration)
   mapped to `""` so it sorts first;
3. among other attributes: no-namespace (`Space == ""`) before any namespaced attribute;
   within no-namespace, by **QName**; otherwise by namespace URI, then by local name.

All comparisons are byte-wise on the UTF-8/UTF-16-code-unit-agnostic Go string ordering; Java's
`String.compareTo` compares UTF-16 code units, Go compares bytes. **These differ for code points
above U+FFFF versus U+E000–U+FFFF** (surrogate pairs sort below `U+E000` in UTF-16 but above it in
UTF-8). Namespace URIs and NCNames in this range do not occur in practice; the KAT corpus includes
a probe document for it, and if the golden disagrees we add an explicit UTF-16-order comparator.
Flag this in review.

Probe:

```
<r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d"/>
all seven → <r xmlns="urn:d" xmlns:a="urn:a" xmlns:b="urn:b" a="4" z="3" a:z="2" b:z="1"></r>
```

### 2.8 Prolog/epilog nodes, and the Santuario epilog quirk

For comments and PIs outside the document element, `documentLevel` decides a newline:

- before the document element: node, then `\n`;
- after the document element: `\n`, then node;
- inside: neither.

`CanonicalizerPhysical` overrides both writers to force "inside", so physical never adds newlines.

```
<!--a--><?p1 x?><r><x/></r><?p2 y?><!--b-->
C14N10WC → <!--a-->\n<?p1 x?>\n<r><x></x></r>\n<?p2 y?>\n<!--b-->
PHYS     → <!--a--><?p1 x?><r><x></x></r><?p2 y?><!--b-->
```

**Quirk (reproduce, do not fix).** When the document element has **no children**, Santuario drops
the entire epilog. In `canonicalizeSubTree`, the childless-element branch overwrites `sibling`
with `firstChild` (nil) and only restores it `if (parentNode != null)` — which is false when the
walk started at the `Document` node, so the loop returns immediately:

```
<!--a--><r/><!--b-->        C14N10WC → <!--a-->\n<r></r>            ("b" lost)
<r/><?pi d?><!--c-->        C14N10WC → <r></r>                      (both lost)
<!--a--><r><x/></r><!--b-->  C14N10WC → <!--a-->\n<r><x></x></r>\n<!--b-->   (fine)
<r>t</r><?pi d?>            C14N10WC → <r>t</r>\n<?pi d?>           (fine)
<r/><?pi d?>                PHYS     → <r></r>                      (physical too)
```

Interoperability with Java DSS is the contract (`PORTING_PLAN.md` §"What 100% compatibility
means"), so `xmlc14n` reproduces this bit-for-bit. It is implemented as the same structural
consequence, not as a special case: port the traversal loop literally, including the
`if parentNode != nil { sibling = cur.NextSibling }` guard, and the behaviour falls out. A
dedicated KAT (`prolog-epilog-empty-root`) locks it, and a comment at the call site records that
it is a deliberate bug-for-bug port with a link to this section.

Whitespace-only character data outside the document element never reaches the tree (§1.5 rejects
non-whitespace, and whitespace produces no node), so `  <r/>\n\n` → `<r></r>`.

### 2.9 Exported API — **PINNED**

```go
package xmlc14n

import (
	"errors"
	"io"

	"<module>/internal/xmldom"
)

type Algorithm string

const (
	C14N10                    Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	C14N10WithComments        Algorithm = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"
	C14N11                    Algorithm = "http://www.w3.org/2006/12/xml-c14n11"
	C14N11WithComments        Algorithm = "http://www.w3.org/2006/12/xml-c14n11#WithComments"
	C14NExclusive             Algorithm = "http://www.w3.org/2001/10/xml-exc-c14n#"
	C14NExclusiveWithComments Algorithm = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"
	C14NPhysical              Algorithm = "http://santuario.apache.org/c14n/physical"

	DefaultDSS     = C14NExclusive // XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD
	DefaultXMLDSig = C14N10        // XMLCanonicalizer.DEFAULT_XMLDSIG_C14N_METHOD
)

// Supported reports whether alg is one of the seven registered algorithms.
func Supported(alg Algorithm) bool

// Resolve maps "" to DefaultXMLDSig and validates alg.
func Resolve(alg Algorithm) (Algorithm, error)

// Input describes what to canonicalize.
//
// Node is the apex: a Document node or an Element node. When Subset is nil the input
// is the subtree rooted at Node and ancestor namespace and xml:* context is inherited
// onto the apex. When Subset is non-nil the input is that document subset and Node is
// the traversal root.
//
// Exclude, when set, skips that element and its subtree (subtree mode only).
// InclusivePrefixes is the exclusive-c14n InclusiveNamespaces PrefixList, already
// split into tokens; "#default" is accepted and normalized. It is ignored by the
// non-exclusive algorithms.
type Input struct {
	Node              *xmldom.Node
	Subset            xmldom.NodeSet
	Exclude           *xmldom.Node
	InclusivePrefixes []string

	// Added in phase 4b - see the amendment note below.
	NodeSet bool
	Filters []NodeFilter
}

// NodeFilter is org.apache.xml.security.signature.NodeFilter: the three-valued node-set
// membership test the XML-DSig transform pipeline attaches to its input.
//
//	 1  include the node
//	 0  exclude the node but keep walking into its subtree
//	-1  exclude the node AND its whole subtree, with no namespace bookkeeping
type NodeFilter interface {
	IsNodeInclude(n *xmldom.Node) (int, error)
	IsNodeIncludeDO(n *xmldom.Node, level int) (int, error)
}

// Canonicalize writes the canonical form of in to w.
func Canonicalize(alg Algorithm, in Input, w io.Writer) error

// CanonicalizeToBytes is Canonicalize into a buffer.
func CanonicalizeToBytes(alg Algorithm, in Input) ([]byte, error)

// CanonicalizeNode is CanonicalizeToBytes over the subtree rooted at n. It is the
// equivalent of XMLCanonicalizer.canonicalize(Node).
func CanonicalizeNode(alg Algorithm, n *xmldom.Node) ([]byte, error)

// CanonicalizeBytes parses src with xmldom's secure defaults and canonicalizes the
// whole document. It is the equivalent of XMLCanonicalizer.canonicalize(byte[]).
func CanonicalizeBytes(alg Algorithm, src []byte) ([]byte, error)

// ParsePrefixList splits an InclusiveNamespaces PrefixList attribute value on
// whitespace and maps "#default" to "xmlns", matching
// org.apache.xml.security.transforms.params.InclusiveNamespaces.prefixStr2Set.
func ParsePrefixList(s string) []string

// ErrUnsupportedAlgorithm is returned by Resolve and Canonicalize.
var ErrUnsupportedAlgorithm = errors.New("xmlc14n: unsupported canonicalization algorithm")

// RelativeNamespaceError reports a rendered namespace declaration whose URI is
// relative. C14N 1.0, 1.1 and exclusive reject these; physical accepts them.
type RelativeNamespaceError struct {
	Element string // the element's QName
	Prefix  string // "xmlns" for the default declaration
	URI     string
}

func (e *RelativeNamespaceError) Error() string
```

**Amendment (phase 4b, agreed with the tech lead).** `Input` gained `NodeSet` and `Filters`, and
`NodeFilter` was added. The reason is that §4's sketch assumed the transform layer could hand
`xmlc14n` a *materialized* `NodeSet`, and it cannot:

- Santuario's `isVisibleDO` answer `-1` prunes a whole subtree **without** pushing a namespace
  frame or running `outputAttributes`, while `0` excludes one node and keeps the frame. A
  materialized set can only express `0`.
- `XPath2NodeFilter` is **stateful**: `isNodeIncludeDO` records the symbol-table level at which
  the subtree it is inside began, and `Canonicalizer20010315.outputAttributes` asks it a second
  time at the *deeper* level after the frame is pushed. Materializing would have to replay that
  traversal exactly, i.e. reimplement the walk it is trying to avoid.

The change is additive: `Subset` still selects node-set mode on its own and every existing KAT
is unaffected, because a materialized subset never answers `-1`. Inside the engine,
`isVisible`/`isVisibleDO`/`isVisibleInt` are now the literal `CanonicalizerBase` trio, the
node-set traversal honours the `-1` prune, and the two node-set emitters ask `isVisibleDO(el,
ns.level())` where they previously asked a boolean.

---

## 3. Known-answer-test strategy

### 3.1 Shape

Same pattern as `internal/cmscore/testdata/gen/`: a Java oracle program, checked in beside the
goldens, with its exact build/run command in the file header; goldens are the **Java** answers,
never the Go port's. Layout:

```
internal/xmlc14n/testdata/
├── corpus/                     the ~32 adversarial documents (§3.3)
├── manifest.txt                one line per KAT: doc, algorithm, scope, apex, prefixlist
├── golden/<doc>.<alg>.<scope>  raw canonical bytes, or a single line "!ERROR <JavaClass>"
└── gen/
    ├── C14nOracle.java         the oracle
    ├── corpus.py               emits corpus/ (byte-exact control over BOMs, CRs, encodings)
    └── generate.sh             build + run, regenerates manifest.txt and golden/
```

`<alg>` is a short slug (`c14n10`, `c14n10wc`, `c14n11`, `c14n11wc`, `excl`, `exclwc`, `phys`);
`<scope>` is `doc` or `apex-<id>` or `apex-<id>-incl-<slug>`.

### 3.2 Test entry points

- `TestC14nGolden` — for every manifest line, parse `corpus/<doc>` with `xmldom`, canonicalize,
  compare bytes to `golden/…`. `!ERROR` lines assert that Go also fails, and that the failure maps
  to the right Go error type (`RelativeNamespaceError`, `*xmldom.SyntaxError`).
- `TestC14nAgainstXMLCanonicalizerCallPath` — a smaller set driven through
  `CanonicalizeBytes`, i.e. exactly what `XMLCanonicalizer.canonicalize(byte[])` does, so the
  parse-then-canonicalize seam is covered end to end.
- `TestDOMRoundTrip` — `Parse` → `Serialize` → `Parse` → `CanonicalizeToBytes` must equal
  `Parse` → `CanonicalizeToBytes`. The semantic complement to the byte-parity oracle in
  `xml/utils`: whatever the Transformer-shaped serializer rewrites must not move the canonical
  form.
- `FuzzCanonicalize` — `Parse` must never panic; when it succeeds, all seven algorithms must run
  without panicking and `Serialize`+re-`Parse`+c14n must be stable. Seeded from `corpus/`. The
  physical method is excluded from the round-trip leg: it is defined to reproduce the namespace
  declarations the serializer drops (redundant rebinds, `xmlns:xml`), so the round trip
  legitimately moves its output — in Java exactly as here.

### 3.3 The corpus (~32 documents)

Each is small, hand-built, and targets a named rule. `corpus.py` writes them byte-exactly so BOMs,
lone CRs, CRLFs and non-UTF-8 encodings survive source control.

| # | File | Targets |
|---|---|---|
| 1 | `minimal.xml` | `<r/>` → `<r></r>` |
| 2 | `prolog-epilog.xml` | comments+PIs before/after a **non-empty** root; positional newlines |
| 3 | `prolog-epilog-empty-root.xml` | the epilog-drop quirk (§2.8) |
| 4 | `prolog-pi-nodata.xml` | `<?a?>` with no data vs `<?b c?>`; PI target with `&#xD;` |
| 5 | `attr-order.xml` | anti-sorted mix of no-ns, prefixed, default-ns-declaring attributes |
| 6 | `attr-many.xml` | 100 attributes, shuffled prefixes and URIs |
| 7 | `attr-ws-normalize.xml` | literal TAB/LF/CR/CRLF vs `&#9;&#10;&#13;` (§1.4) |
| 8 | `attr-quotes.xml` | `'` and `"` mixed delimiters, `&quot;`, `&apos;` |
| 9 | `text-escapes.xml` | `& < > "` literal and referenced; `&#13;` in text |
| 10 | `cdata.xml` | CDATA with `<`/`&`, split `]]]]><![CDATA[>`, adjacent text/CDATA runs |
| 11 | `comments-pi-nested.xml` | comments/PIs at every depth, `#WithComments` on/off |
| 12 | `ns-superfluous.xml` | same prefix redeclared to the same URI at three depths |
| 13 | `ns-rebind.xml` | prefix rebound to a different URI, then back |
| 14 | `ns-unused.xml` | declared-but-unused prefixes (1.0/1.1 keep, exclusive drops) |
| 15 | `ns-default-undeclare.xml` | `xmlns="u"` → `xmlns=""` → redeclare |
| 16 | `ns-default-apex.xml` | apex is the undeclaring element (§2.4 suppression case) |
| 17 | `ns-xml-prefix.xml` | explicit `xmlns:xml=".../XML/1998/namespace"` |
| 18 | `ns-relative.xml` | relative URI, used and unused — errors for six, output for physical |
| 19 | `excl-visibly-utilized.xml` | element prefix, attribute prefix, unprefixed attr, default-ns child |
| 20 | `excl-prefixlist.xml` | apexes run with `""`, `"q"`, `"#default q"`, `"zz"` |
| 21 | `xmlattrs-lang-space.xml` | `xml:lang`/`xml:space` on several ancestors, deep apex |
| 22 | `xmlbase-absolute.xml` | `xml:base` chain with an absolute outermost |
| 23 | `xmlbase-relative.xml` | all-relative chain — 1.0 vs 1.1 `joinURI` divergence |
| 24 | `xmlbase-dotsegments.xml` | `../`, `./`, `//`, trailing `..` — exercises `removeDotSegments` |
| 25 | `xml-id.xml` | `xml:id` on ancestors and on the apex — 1.0 vs 1.1 divergence |
| 26 | `unicode.xml` | astral plane, combining marks, NEL U+0085, LSEP U+2028, accented names |
| 27 | `sort-astral.xml` | attribute local names straddling U+FFFF (UTF-16 vs UTF-8 order, §2.7) |
| 28 | `bom-decl.xml` | UTF-8 BOM + XML declaration + `standalone` |
| 29 | `encoding-latin1.xml` | declared `ISO-8859-1` with high bytes |
| 30 | `entities.xml` | all five predefined entities in text and attribute values |
| 31 | `empty-elements.xml` | `<a></a>` vs `<a/>`, empty text node, whitespace-only children |
| 32 | `deep-nesting.xml` | 200 levels, apex at 100 |
| 33 | `xades-signature.xml` | a real upstream XAdES-B `ds:Signature` (from `dss-xades/src/test/resources`) |
| 34 | `trusted-list.xml` | trimmed EU trusted list — heavy default-ns, `Id` attributes |

Apexes per document are named in `manifest.txt` by an `Id`/`ID`/`xml:id` attribute so the oracle
and the Go test select the same element without XPath.

Negative corpus (parse must fail on both sides, asserted by `TestParseRejects`):
`<!DOCTYPE …>`, internal entity subset, `xmlns:p=""`, `xmlns:xml="urn:bogus"`, undeclared prefix,
duplicate attribute, `<r></s>`, two root elements, trailing text, `&#0;`, `&foo;`, `]]>` in text,
`<` in an attribute value, `<?xml version="1.1"?>`.

### 3.4 Oracle program

`internal/xmlc14n/testdata/gen/C14nOracle.java`, header comment carrying the exact command:

```
DSS=/home/user/dss-upstream
M2=$HOME/.m2/repository
CP=$DSS/dss-xml-utils/target/dss-xml-utils-6.5.RC1.jar\
:$M2/org/apache/santuario/xmlsec/3.0.6/xmlsec-3.0.6.jar\
:$M2/org/slf4j/slf4j-api/2.0.18/slf4j-api-2.0.18.jar\
:<dss-model, dss-utils, dss-xml-common jars from the maven-built upstream>
javac -encoding UTF-8 -cp "$CP" -d /tmp/c14noracle C14nOracle.java
java -cp "$CP:/tmp/c14noracle" C14nOracle <corpus dir> <golden dir> <manifest file>
```

Rules the oracle obeys:

1. **A fresh canonicalizer per KAT** (`XMLCanonicalizer.createInstance(uri)` /
   `Canonicalizer.getInstance(uri)`), because of SANTUARIO-463 (§0).
2. Document-scope KATs go through `XMLCanonicalizer.createInstance(uri).canonicalize(byte[])` —
   the exact DSS call path, including Santuario's own secure parser.
3. Apex-scope KATs parse with `DomUtils.buildDOM(byte[])` — the DSS secure
   `DocumentBuilderFactory` — then call `XMLCanonicalizer.createInstance(uri).canonicalize(Node)`.
4. PrefixList KATs need a parameter `XMLCanonicalizer` does not expose, so they drop to
   `Canonicalizer.getInstance(uri).canonicalizeSubtree(node, prefixList, out)`, with a comment
   stating that `XMLCanonicalizer` is a thin, parameterless wrapper over exactly this call.
5. Failures are recorded as `!ERROR <SimpleClassName>` in the golden file and the manifest,
   never as a missing file, so a Go-side success on a Java-side failure is a test failure.
6. `manifest.txt` also records the SHA-256 of each golden, so a corrupted regeneration is visible
   in review.
7. The oracle never reads a Go source file and the Go tests never invoke Java. Regeneration is a
   deliberate, reviewed act — a golden that changes in a PR is a red flag, exactly as for the
   BouncyCastle oracle.

### 3.5 Cross-validation beyond KATs

Phase 4b/4c gates, listed here so nobody designs them away:

- Go-produced XAdES-B/T signatures validate in Java DSS (both `SignedInfo` and
  `SignedProperties` digests reproduce, which is a c14n equality proof over real inputs).
- Every upstream `dss-xades/src/test/resources` signature validates in the Go port, and for each
  its recomputed reference digests match the `ds:DigestValue` in the file — a large, free
  c14n corpus with independently-derived expected answers.

---

## 4. What the `xmldsig` layer above will need (interface sketch only)

**Status: superseded by the code.** `internal/xmldsig` exists; its `doc.go` carries the
authoritative Go-name-to-Santuario-member table and its own deviations. The sketch below is
kept because it records what the two lower packages were designed to be sufficient for, and
because three of its guesses turned out wrong in ways worth remembering:

- `Data` is a struct, not an interface. Santuario's `XMLSignatureInput` is a union whose
  discriminator the canonicalizer dispatches on (`isOctetStream`, then `isElement`, then
  `isNodeSet`, in that order), and an interface hierarchy cannot express states such as "an
  element subtree that also carries an exclude node and two node filters".
- `NodeSetData` with a materialized `Nodes` set is not enough - see the §2.9 amendment.
- The transform signature needs the `ds:Transform` element *and* the base URI *and* the
  secure-validation flag, because a transform resolves its own namespace prefixes against the
  element that carries it.

Not binding beyond the fact that `xmldom`/`xmlc14n` must be sufficient for it. Written now so the
two implementers can see the consumers.

```go
package xmldsig

// Data is a transform pipeline value: octets or a document subset. It mirrors
// org.apache.xml.security.signature.XMLSignatureInput.
type Data interface{ isData() }

type OctetData struct{ Octets []byte }

type NodeSetData struct {
	Root            *xmldom.Node   // traversal root (document or element)
	Nodes           xmldom.NodeSet // nil ⇒ the whole subtree under Root
	ExcludeComments bool           // XMLDSIG 4.4.3.3 step 4
}

// Transform is one ds:Transform. Params is the ds:Transform element itself, so a
// transform can read its own children (XPath expressions, InclusiveNamespaces, …).
type Transform interface {
	Algorithm() string
	Transform(in Data, params *xmldom.Node) (Data, error)
}

// TransformRegistry maps algorithm URIs to implementations; the built-ins are
// enveloped-signature, XPath, XPath2 filter, base64, XSLT (rejected by default), and
// the seven canonicalization methods (each a Transform as well as a c14n method).
type TransformRegistry interface {
	Lookup(alg string) (Transform, bool)
	Register(t Transform)
}

// URIResolver dereferences a ds:Reference URI. SameDocumentResolver handles "",
// "#id" (via xmldom.ElementByID) and "#xpointer(...)"; others are supplied by the
// caller for detached references.
type URIResolver interface {
	CanResolve(uri, baseURI string) bool
	Resolve(uri, baseURI string, ctx *xmldom.Node) (Data, error)
}

// Reference is one ds:Reference: dereference, run the transform chain, canonicalize
// any surviving node set with the default method, then digest.
type Reference struct {
	URI        string
	Type       string
	ID         string
	Transforms []Transform
	Digest     enumerations.DigestAlgorithm
}

func (r *Reference) Process(ctx *xmldom.Node, res URIResolver, reg TransformRegistry) (Data, []byte, error)

// SignedInfo canonicalizes with CanonicalizationMethod and returns the octets to sign.
type SignedInfo struct {
	CanonicalizationMethod xmlc14n.Algorithm
	SignatureMethod        string
	References             []*Reference
}

func (si *SignedInfo) CanonicalOctets(el *xmldom.Node) ([]byte, error)
```

Two consequences that constrain §1 and §2 and are therefore binding:

- A transform that yields a node set must be able to hand `xmlc14n` a `NodeSet` plus a root —
  hence `Input.Subset` and `xmldom.NodeSet` exist from the start (§2.3).
- `DSSXMLUtils`'s enveloped-signature handling removes comment nodes before transforming
  (`DomUtils.excludeComments`, per XMLDSIG 4.4.3.3 step 4) — that is a tree edit on a **clone**,
  so `Clone(deep)` and `Import` must be correct and cheap.

---

## 5. Review checklist for the two implementers

- [ ] `xml.Attr.Value` appears nowhere in `xmldom` except in the start-tag cross-check (§1.4).
- [ ] `RawToken()`, never `Token()`.
- [ ] Every one of the ten checks in §1.5 has a negative test.
- [ ] `nsStack` keeps `lastrendered`; no map iteration reaches the output.
- [ ] `joinURI`/`removeDotSegments` are a line-for-line port; `net/url` is not used for them.
- [ ] The epilog-drop quirk (§2.8) is reproduced, tested and commented as deliberate.
- [ ] `RelativeNamespaceError` fires only for *rendered* declarations, and never for physical.
- [ ] Canonicalization allocates all state per call; no package-level mutable state.
- [ ] Goldens are Java's answers; no golden was ever produced by the Go implementation.

---

## 6. Amendments from the phase 4c audit

Three divergences from Xerces/Xalan found by an independent re-run of the oracles, and the
decisions taken on them. Each is now pinned by a test whose failure was verified by mutation.

### 6.1 `ownerDocument` survives detachment (`OwnerDocument`, `Node.owner`)

`org.w3c.dom` assigns a node's owning document at creation and never clears it, so
`removeChild` leaves `getOwnerDocument()` intact. `DomUtils.serializeNode` reads exactly that
to choose its output encoding. Deriving ownership purely positionally answered `nil` for a
detached subtree and fell back to UTF-8, so an element lifted out of an ISO-8859-1 document
and serialized on its own came out in UTF-8 here and in ISO-8859-1 in Java — a different
byte string under the digest of a no-transform `ds:Reference`.

`Document()` keeps its positional meaning, which is what tree-scope decisions want.
`OwnerDocument()` is the DOM one: it walks to the topmost ancestor and, when that is not a
Document, reads an `owner` recorded at the moment of detachment (`unlink`, `ReplaceChild`,
`RemoveAttr`, `SetTextContent`). `Clone` carries the source's owner, as `cloneNode` does;
`Import` assigns the importing document, as `importNode` does. `XMLEncoding` and
`XMLStandalone` read `OwnerDocument()` — for a Document node, itself.

### 6.2 A duplicated `Id` resolves to the LAST element, not the first (`ids.go`)

`CoreDocumentImpl.putIdentifier` is `identifiers.put(id, element)` into a `HashMap`, and
`XAdESDOMDocument.recursiveIdBrowse` registers in document order, so each duplicate replaces
the one before it. Probed against OpenJDK 21 with DSS's own registration loop, and again
through the XPath `id()` function.

This is not cosmetic. Two elements sharing an `Id` is the shape of an XML signature wrapping
attack, and resolving `#id` to the first element makes the forged digest match:
`internal/xmldsig/testdata/corpus/validation/dss2329/xades-with-manifest-with-duplicated-reference.xml`
is that document, Santuario answers `false` on it, and this port answered `true` until the
index was changed to last-wins. `DuplicateIDs()` still reports the collision for
`DSSXMLUtils.isDuplicateIdsDetected`.

### 6.3 ISO-8859-2 (`codepage.go`)

Two upstream XAdES fixtures declare it (`validation/Signature-X-HU_MIC-1.xml`,
`validation/BaselineBWithCertificateValues.xml`) and Xerces reads them, so refusing the
encoding meant being unable to validate a Central-European signature at all. Added on both
sides — decode in `decodeSource`, `writeLatin2`/`latin2InEncoding` in the output table — with
the alias set and the declaration-echo rules probed against OpenJDK 21 rather than guessed
(`ISO8859_2`, `ISO8859-2` and `iso-8859-2` are rewritten to `ISO-8859-2`; `latin2`,
`csISOLatin2`, `ISO_8859-2` and `cp912` are echoed unchanged; `8859_2` is rejected outright,
as `latin-1` is).

The table is pinned character-for-character against `new String(bytes, "ISO-8859-2")`, not
only through the serializer goldens: decode and encode read the same table, so a wrong entry
round-trips to the same byte and is invisible in the output bytes while changing the
character canonicalization writes out as UTF-8.

### 6.4 Encodings still refused

`windows-1252` and XML 1.1 remain fail-closed refusals at `Parse` (D4 for 1.1). No upstream
fixture uses either; both are refusals, never wrong answers.
