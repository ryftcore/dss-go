package xmldom

// Attr returns the attribute node with the given namespace URI and local name.
// Namespace declarations are found with space == XMLNSNamespace; the default
// declaration xmlns="..." has local name "xmlns".
func (n *Node) Attr(space, local string) *Node {
	if n == nil {
		return nil
	}
	for _, a := range n.Attrs {
		if a.Name.Space == space && a.Name.Local == local {
			return a
		}
	}
	return nil
}

// AttrValue is Attr(...).Value, or "" when absent.
func (n *Node) AttrValue(space, local string) string {
	if a := n.Attr(space, local); a != nil {
		return a.Value
	}
	return ""
}

// SetAttr sets or replaces an attribute, preserving its position in Attrs when it
// already exists and appending otherwise. It returns the attribute node.
func (n *Node) SetAttr(name Name, value string) *Node {
	if n == nil || n.Kind != Element {
		panic("xmldom: SetAttr on a non-element")
	}
	n.bump()
	if a := n.Attr(name.Space, name.Local); a != nil {
		// The expanded name is unchanged by definition; the literal prefix may not
		// be, and it is load-bearing for canonicalization, so it is overwritten too.
		a.Name = name
		a.Value = value
		return a
	}
	a := &Node{Kind: Attribute, Name: name, Value: value, Parent: n}
	n.Attrs = append(n.Attrs, a)
	return a
}

// RemoveAttr removes an attribute and reports whether one was removed.
func (n *Node) RemoveAttr(space, local string) bool {
	if n == nil {
		return false
	}
	for i, a := range n.Attrs {
		if a.Name.Space == space && a.Name.Local == local {
			n.bump()
			rememberOwner(a)
			a.Parent = nil
			n.Attrs = append(n.Attrs[:i], n.Attrs[i+1:]...)
			return true
		}
	}
	return false
}
