package xmldom

// NodeSet is a document-subset membership set for canonicalization. Namespace
// declarations participate as their attribute nodes, matching Santuario.
type NodeSet map[*Node]struct{}

// NewNodeSet returns a set containing nodes.
func NewNodeSet(nodes ...*Node) NodeSet {
	s := make(NodeSet, len(nodes))
	s.Add(nodes...)
	return s
}

// Add inserts nodes into s. Nil nodes are ignored.
func (s NodeSet) Add(nodes ...*Node) {
	for _, n := range nodes {
		if n != nil {
			s[n] = struct{}{}
		}
	}
}

// Remove deletes nodes from s.
func (s NodeSet) Remove(nodes ...*Node) {
	for _, n := range nodes {
		delete(s, n)
	}
}

// Has reports whether n is a member of s.
func (s NodeSet) Has(n *Node) bool {
	_, ok := s[n]
	return ok
}

// Len returns the number of members.
func (s NodeSet) Len() int { return len(s) }

// AddSubtree adds n, all its descendants, and all their attributes to s.
func (s NodeSet) AddSubtree(n *Node) {
	n.Walk(func(x *Node) bool {
		s[x] = struct{}{}
		for _, a := range x.Attrs {
			s[a] = struct{}{}
		}
		return true
	})
}
