package xmldom

import "strconv"

// Kind discriminates the node types. There is deliberately no interface hierarchy:
// canonicalization walks by pointer and uses pointer identity as node-set membership,
// so a single concrete struct keeps *Node a stable, comparable identity and keeps the
// parent/sibling links uniform.
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

var kindNames = [...]string{
	Document:  "Document",
	Element:   "Element",
	Attribute: "Attribute",
	Text:      "Text",
	CDATA:     "CDATA",
	Comment:   "Comment",
	ProcInst:  "ProcInst",
}

func (k Kind) String() string {
	if int(k) < len(kindNames) && kindNames[k] != "" {
		return kindNames[k]
	}
	return "Kind(" + strconv.Itoa(int(k)) + ")"
}

// Node is every node type. Pointer identity is node identity: NodeSet membership,
// subset visibility and Exclude comparisons are all pointer comparisons.
//
// Attrs holds an element's attributes, including namespace declarations, in DOCUMENT
// ORDER. Callers must not reorder or mutate it directly; use SetAttr/RemoveAttr.
//
// Field usage per kind:
//
//	Document   Name zero      Value ""                       Attrs nil   children: prolog/epilog + root
//	Element    Name set       Value ""                       Attrs set   children: yes
//	Attribute  Name set       Value normalized value         Attrs nil   children: no
//	Text       Name zero      Value character data           Attrs nil   children: no
//	CDATA      Name zero      Value character data           Attrs nil   children: no
//	Comment    Name zero      Value comment data, no markers Attrs nil   children: no
//	ProcInst   Name.Local set Value instruction data         Attrs nil   children: no
//
// An Attribute node's Parent is the element carrying it. org.w3c.dom reports null
// there and offers ownerElement instead; we follow XPath, because canonicalization
// routinely has an attribute node in hand and needs its element.
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

	// doc is non-nil only on a Document node. Ownership of every other node is
	// derived positionally, by walking Parent (see Document).
	doc *docState

	// owner remembers the document a node belonged to when it was detached from it,
	// and is meaningful only on the root of a detached subtree (OwnerDocument walks
	// up to that root before reading it).
	//
	// org.w3c.dom assigns ownerDocument at creation - every node comes out of a
	// Document.createXxx factory - and removeChild does NOT clear it, so a detached
	// node still answers getOwnerDocument(). DomUtils.serializeNode reads exactly that
	// to pick the output encoding, so a subtree lifted out of an ISO-8859-1 document
	// and serialized on its own is written in ISO-8859-1 by Java. Deriving ownership
	// purely positionally answered nil there and fell back to UTF-8, which changed the
	// bytes a ds:Reference digests. See OwnerDocument.
	owner *Node
}

// docState is the per-document bookkeeping that hangs off a Document node.
type docState struct {
	// gen counts structural and attribute mutations anywhere in the document. The
	// ID index is rebuilt lazily whenever idGen falls behind it.
	gen uint64

	idGen   uint64
	built   bool
	scanAll bool // RegisterIDs was called and must be replayed on rebuild
	manual  []manualID

	ids     map[string]*Node
	idAttrs map[*Node][]*Node
	dups    map[string]struct{}

	// The XML declaration of the parsed source, kept because the serializer needs it:
	// Java reads all three back off the DOM (Document.getXmlEncoding /
	// getXmlStandalone / getXmlVersion) and they decide the declaration it writes and
	// the charset it writes the document in. Zero values for a document built in
	// memory, which is what Xerces reports for one too.
	xmlEncoding   string
	xmlStandalone bool
}

type manualID struct {
	elem *Node
	attr *Node
}

func newDocState() *docState {
	return &docState{
		ids:     make(map[string]*Node),
		idAttrs: make(map[*Node][]*Node),
		dups:    make(map[string]struct{}),
	}
}

// ---------------------------------------------------------------- construction

// NewDocument returns an empty document node.
func NewDocument() *Node { return &Node{Kind: Document, doc: newDocState()} }

// NewElement returns a detached element node.
func NewElement(name Name) *Node { return &Node{Kind: Element, Name: name} }

// NewAttr returns a detached attribute node.
func NewAttr(name Name, value string) *Node {
	return &Node{Kind: Attribute, Name: name, Value: value}
}

// NewText returns a detached text node.
func NewText(data string) *Node { return &Node{Kind: Text, Value: data} }

// NewCDATA returns a detached CDATA section node.
func NewCDATA(data string) *Node { return &Node{Kind: CDATA, Value: data} }

// NewComment returns a detached comment node.
func NewComment(data string) *Node { return &Node{Kind: Comment, Value: data} }

// NewProcInst returns a detached processing-instruction node.
func NewProcInst(target, data string) *Node {
	return &Node{Kind: ProcInst, Name: Name{Local: target}, Value: data}
}

// ---------------------------------------------------------------- tree mutation

// canContain reports whether a parent of kind p may hold a child of kind c.
func canContain(p, c Kind) bool {
	switch p {
	case Document:
		return c == Element || c == Comment || c == ProcInst
	case Element:
		return c == Element || c == Text || c == CDATA || c == Comment || c == ProcInst
	}
	return false
}

// AppendChild appends c to n and returns c. It panics if c already has a parent or if
// the parent/child kind combination is illegal.
func (n *Node) AppendChild(c *Node) *Node { return n.InsertBefore(c, nil) }

// InsertBefore inserts c before ref (a child of n) and returns c. A nil ref appends.
func (n *Node) InsertBefore(c, ref *Node) *Node {
	n.checkInsert(c, ref)
	if ref != nil && ref.Parent != n {
		panic("xmldom: InsertBefore: ref is not a child of n")
	}
	n.link(c, ref)
	n.bump()
	return c
}

// checkInsert panics unless c may be inserted into n.
func (n *Node) checkInsert(c, ref *Node) {
	if n == nil || c == nil {
		panic("xmldom: nil node")
	}
	if c.Parent != nil {
		panic("xmldom: node already has a parent")
	}
	if !canContain(n.Kind, c.Kind) {
		panic("xmldom: cannot insert " + c.Kind.String() + " into " + n.Kind.String())
	}
	if c == n || c.Contains(n) {
		panic("xmldom: insertion would create a cycle")
	}
	if n.Kind == Document && c.Kind == Element {
		if e := n.DocumentElement(); e != nil && e != ref {
			panic("xmldom: a document may have only one element child")
		}
	}
}

// link splices c into n's child list before ref (nil ref appends). No validation.
func (n *Node) link(c, ref *Node) {
	c.Parent = n
	if ref == nil {
		c.PrevSibling = n.LastChild
		c.NextSibling = nil
		if n.LastChild != nil {
			n.LastChild.NextSibling = c
		} else {
			n.FirstChild = c
		}
		n.LastChild = c
		return
	}
	c.PrevSibling = ref.PrevSibling
	c.NextSibling = ref
	if ref.PrevSibling != nil {
		ref.PrevSibling.NextSibling = c
	} else {
		n.FirstChild = c
	}
	ref.PrevSibling = c
}

// rememberOwner records the document c is being detached from, so that
// OwnerDocument keeps answering it the way org.w3c.dom's ownerDocument does.
func rememberOwner(c *Node) {
	if c.Kind == Document {
		return
	}
	if d := c.Document(); d != nil {
		c.owner = d
	}
}

// unlink removes c from its parent's child list. No validation.
func (n *Node) unlink(c *Node) {
	rememberOwner(c)
	if c.PrevSibling != nil {
		c.PrevSibling.NextSibling = c.NextSibling
	} else {
		n.FirstChild = c.NextSibling
	}
	if c.NextSibling != nil {
		c.NextSibling.PrevSibling = c.PrevSibling
	} else {
		n.LastChild = c.PrevSibling
	}
	c.Parent = nil
	c.PrevSibling = nil
	c.NextSibling = nil
}

// RemoveChild detaches c from n and returns c.
func (n *Node) RemoveChild(c *Node) *Node {
	if c == nil || c.Parent != n {
		panic("xmldom: RemoveChild: not a child of n")
	}
	n.bump()
	n.unlink(c)
	return c
}

// ReplaceChild puts newC where oldC was and returns oldC.
//
// The splice is done in place rather than as remove-then-insert so that replacing a
// document's root element never transiently violates the one-element-child rule.
func (n *Node) ReplaceChild(newC, oldC *Node) *Node {
	if oldC == nil || oldC.Parent != n {
		panic("xmldom: ReplaceChild: oldC is not a child of n")
	}
	n.checkInsert(newC, oldC)
	n.bump()
	rememberOwner(oldC)
	newC.Parent = n
	newC.PrevSibling = oldC.PrevSibling
	newC.NextSibling = oldC.NextSibling
	if oldC.PrevSibling != nil {
		oldC.PrevSibling.NextSibling = newC
	} else {
		n.FirstChild = newC
	}
	if oldC.NextSibling != nil {
		oldC.NextSibling.PrevSibling = newC
	} else {
		n.LastChild = newC
	}
	oldC.Parent = nil
	oldC.PrevSibling = nil
	oldC.NextSibling = nil
	return oldC
}

// Clone deep- or shallow-copies n. The copy has no parent, and - like
// org.w3c.dom's cloneNode - keeps n's owning document as its ownerDocument until it
// is inserted somewhere.
func (n *Node) Clone(deep bool) *Node {
	if n == nil {
		return nil
	}
	c := &Node{Kind: n.Kind, Name: n.Name, Value: n.Value, owner: n.OwnerDocument()}
	if n.Kind == Document {
		c.doc = newDocState()
		// Xerces' Document.cloneNode copies the declaration fields, and
		// DSSXMLUtils.createDeepCopy round-trips a document through them before
		// canonicalizing, so a copy has to serialize the way its original does.
		c.doc.xmlEncoding = n.doc.xmlEncoding
		c.doc.xmlStandalone = n.doc.xmlStandalone
	}
	if len(n.Attrs) > 0 {
		c.Attrs = make([]*Node, 0, len(n.Attrs))
		for _, a := range n.Attrs {
			c.Attrs = append(c.Attrs, &Node{Kind: a.Kind, Name: a.Name, Value: a.Value, Parent: c})
		}
	}
	if deep {
		for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
			c.link(ch.Clone(true), nil)
		}
	}
	return c
}

// Import returns a deep or shallow copy of src owned by n's document, ready to insert.
//
// Ownership in this model is positional whenever the node is in a tree - Document
// walks the Parent chain - so an imported node becomes owned the moment it is
// inserted. Until then it answers n's document, which is what importNode does. The
// method exists because callers ported from org.w3c.dom expect importNode to be the
// thing they call before adopting foreign nodes.
func (n *Node) Import(src *Node, deep bool) *Node {
	c := src.Clone(deep)
	if c != nil {
		c.owner = nil
		if n.Kind == Document {
			c.owner = n
		} else {
			c.owner = n.OwnerDocument()
		}
	}
	return c
}

// bump invalidates the owning document's derived indexes.
func (n *Node) bump() {
	if d := n.Document(); d != nil && d.doc != nil {
		d.doc.gen++
	}
}
