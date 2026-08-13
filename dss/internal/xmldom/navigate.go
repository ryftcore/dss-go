package xmldom

import "strings"

// Document returns the owning document node, or nil for a detached subtree.
func (n *Node) Document() *Node {
	for c := n; c != nil; c = c.Parent {
		if c.Kind == Document {
			return c
		}
	}
	return nil
}

// DocumentElement returns the single element child of a document node, else nil.
func (n *Node) DocumentElement() *Node {
	if n == nil || n.Kind != Document {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind == Element {
			return c
		}
	}
	return nil
}

// Children returns a snapshot of n's children.
func (n *Node) Children() []*Node {
	if n == nil || n.FirstChild == nil {
		return nil
	}
	var out []*Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		out = append(out, c)
	}
	return out
}

// Elements returns a snapshot of n's element children.
func (n *Node) Elements() []*Node {
	if n == nil {
		return nil
	}
	var out []*Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind == Element {
			out = append(out, c)
		}
	}
	return out
}

// FirstElementChild returns the first element child, or nil.
func (n *Node) FirstElementChild() *Node {
	if n == nil {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind == Element {
			return c
		}
	}
	return nil
}

// NextElementSibling returns the next element sibling, or nil.
func (n *Node) NextElementSibling() *Node {
	if n == nil {
		return nil
	}
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Kind == Element {
			return c
		}
	}
	return nil
}

// Ancestors returns n's ancestors, nearest first, excluding n.
func (n *Node) Ancestors() []*Node {
	if n == nil {
		return nil
	}
	var out []*Node
	for c := n.Parent; c != nil; c = c.Parent {
		out = append(out, c)
	}
	return out
}

// Depth returns the number of ancestors of n; a document node has depth 0.
func (n *Node) Depth() int {
	if n == nil {
		return 0
	}
	d := 0
	for c := n.Parent; c != nil; c = c.Parent {
		d++
	}
	return d
}

// Contains reports whether other is n or a descendant of n. An attribute node counts
// as a descendant of the element carrying it.
func (n *Node) Contains(other *Node) bool {
	if n == nil || other == nil {
		return false
	}
	for c := other; c != nil; c = c.Parent {
		if c == n {
			return true
		}
	}
	return false
}

// TextContent concatenates the data of all Text and CDATA descendants in document order.
// For a Text, CDATA or Attribute node it is the node's own value.
func (n *Node) TextContent() string {
	if n == nil {
		return ""
	}
	switch n.Kind {
	case Text, CDATA, Attribute:
		return n.Value
	}
	var b strings.Builder
	n.Walk(func(x *Node) bool {
		if x.Kind == Text || x.Kind == CDATA {
			b.WriteString(x.Value)
		}
		return true
	})
	return b.String()
}

// SetTextContent replaces all children of an element with a single text node. An empty
// data string leaves the element childless, matching org.w3c.dom's setTextContent.
func (n *Node) SetTextContent(data string) {
	if n == nil || n.Kind != Element {
		panic("xmldom: SetTextContent on a non-element")
	}
	n.bump()
	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		c.Parent = nil
		c.PrevSibling = nil
		c.NextSibling = nil
		c = next
	}
	n.FirstChild = nil
	n.LastChild = nil
	if data != "" {
		n.link(NewText(data), nil)
	}
}

// Walk calls fn for n and every descendant in document order, attributes excluded.
// Returning false from fn skips that node's subtree.
func (n *Node) Walk(fn func(*Node) bool) {
	if n == nil {
		return
	}
	if !fn(n) {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		c.Walk(fn)
	}
}

// LookupNamespaceURI resolves prefix ("" means the default namespace) against the
// declarations in scope at n. It reports whether a binding was found.
//
// An explicit xmlns="" is a binding: it reports ("", true), because the default
// namespace has been deliberately undeclared. Only the absence of any declaration
// reports ("", false).
func (n *Node) LookupNamespaceURI(prefix string) (string, bool) {
	switch prefix {
	case "xml":
		return XMLNamespace, true
	case "xmlns":
		return XMLNSNamespace, true
	}
	local := prefix
	if prefix == "" {
		local = "xmlns"
	}
	for el := elementOf(n); el != nil; el = el.Parent {
		if el.Kind != Element {
			continue
		}
		for _, a := range el.Attrs {
			if a.Name.Space == XMLNSNamespace && a.Name.Local == local {
				return a.Value, true
			}
		}
	}
	return "", false
}

// LookupPrefix returns the innermost prefix bound to uri that is in scope at n.
//
// The default declaration is never a candidate: an unprefixed name is not a way to
// name uri for an attribute, and callers ask this question to build a QName. A prefix
// rebound closer to n shadows its outer binding and is skipped.
func (n *Node) LookupPrefix(uri string) (string, bool) {
	if uri == "" {
		return "", false
	}
	if uri == XMLNamespace {
		return "xml", true
	}
	shadowed := make(map[string]struct{})
	for el := elementOf(n); el != nil; el = el.Parent {
		if el.Kind != Element {
			continue
		}
		for _, a := range el.Attrs {
			if a.Name.Space != XMLNSNamespace || a.Name.Prefix != "xmlns" {
				continue
			}
			p := a.Name.Local
			if _, dup := shadowed[p]; dup {
				continue
			}
			shadowed[p] = struct{}{}
			if a.Value == uri {
				return p, true
			}
		}
	}
	return "", false
}

// elementOf returns n itself, or the element carrying n when n is an attribute.
func elementOf(n *Node) *Node {
	if n != nil && n.Kind == Attribute {
		return n.Parent
	}
	return n
}
