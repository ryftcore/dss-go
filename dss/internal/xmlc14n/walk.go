// Ported from org.apache.xml.security.c14n.implementations.CanonicalizerBase
// (Apache Santuario xmlsec 3.0.6): the traversal, shared by every algorithm, plus the
// ancestor-context gathering that feeds it.
package xmlc14n

import (
	"bufio"
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
)

// Where the node being written sits relative to the document element. The values decide the
// positional newline around a prolog or epilog comment or processing instruction.
const (
	nodeBeforeDocumentElement           = -1
	nodeNotBeforeOrAfterDocumentElement = 0
	nodeAfterDocumentElement            = 1
)

// engine holds one canonicalization. Every field is created by newEngine and dies with the
// call, so nothing here is shared between calls and the package has no mutable global state.
type engine struct {
	includeComments   bool
	exclusive         bool
	c14n11            bool
	physical          bool
	inclusivePrefixes []string // exclusive only, sorted; from ParsePrefixList

	w        *bufio.Writer
	subset   xmldom.NodeSet // explicit node-set membership; nil means "no explicit set"
	exclude  *xmldom.Node   // subtree mode only
	xmlAttrs *xmlAttrStack

	// filters is CanonicalizerBase.nodeFilter, the NodeFilter list the XML-DSig transform
	// pipeline attaches to its XMLSignatureInput. filterErr is the first failure any of them
	// reported; see noteFilterErr.
	filters   []NodeFilter
	filterErr error

	// firstCall is Canonicalizer20010315.firstCall: the flag that flushes the ancestors'
	// namespace and xml:* context onto the apex, exactly once. Santuario never resets it,
	// which is why a reused canonicalizer silently drops that context (SANTUARIO-463); here
	// it lives and dies with the call.
	firstCall bool
}

// canonicalizeSubTree ports CanonicalizerBase.canonicalizeSubTree, iteratively and literally,
// including the branch that produces the epilog-drop quirk: when the document element has no
// children, the childless-element branch overwrites sibling with its (nil) first child and
// only restores it when parentNode != nil, which is false for a walk started at the Document
// node, so the loop returns and the whole epilog is dropped. That is a Santuario bug, and
// reproducing it bit-for-bit is the interoperability contract with Java DSS (design section
// 2.8; KAT prolog-epilog-empty-root). Do not "fix" it, and do not restructure the loop.
func (e *engine) canonicalizeSubTree(cur *xmldom.Node, ns *nsStack, endnode *xmldom.Node, documentLevel int) error {
	if cur == nil {
		return nil
	}
	var sibling, parentNode *xmldom.Node
	for {
		switch cur.Kind {
		case xmldom.Attribute:
			return fmt.Errorf("xmlc14n: illegal node type during traversal: %s", cur.Kind)

		case xmldom.Document:
			ns.outputNodePush()
			sibling = cur.FirstChild

		case xmldom.Comment:
			if e.includeComments {
				e.outputComment(cur.Value, documentLevel)
			}

		case xmldom.ProcInst:
			e.outputPI(cur.Name.Local, cur.Value, documentLevel)

		case xmldom.Text, xmldom.CDATA:
			writeTextEscaped(e.w, cur.Value)

		case xmldom.Element:
			documentLevel = nodeNotBeforeOrAfterDocumentElement
			if e.exclude != nil && cur == e.exclude {
				break
			}
			ns.outputNodePush()
			name := cur.Name.QName()
			e.w.WriteByte('<')
			e.w.WriteString(name)
			if err := e.outputAttributesSubtree(cur, ns); err != nil {
				return err
			}
			e.w.WriteByte('>')
			sibling = cur.FirstChild
			if sibling == nil {
				e.w.WriteString("</")
				e.w.WriteString(name)
				e.w.WriteByte('>')
				ns.outputNodePop()
				if parentNode != nil {
					sibling = cur.NextSibling
				}
			} else {
				parentNode = cur
			}
		}

		for sibling == nil && parentNode != nil {
			e.w.WriteString("</")
			e.w.WriteString(parentNode.Name.QName())
			e.w.WriteByte('>')
			ns.outputNodePop()
			if parentNode == endnode {
				return nil
			}
			sibling = parentNode.NextSibling
			parentNode = parentNode.Parent
			if parentNode == nil || parentNode.Kind != xmldom.Element {
				documentLevel = nodeAfterDocumentElement
				parentNode = nil
			}
		}
		if sibling == nil {
			return nil
		}
		cur = sibling
		sibling = cur.NextSibling
	}
}

// canonicalizeXPathNodeSet ports CanonicalizerBase.canonicalizeXPathNodeSet: the same walk
// over a document subset, where an element can be invisible while its descendants are visible.
// Reached through Input.Subset, i.e. from the XPath, XPath2-filter and enveloped-signature
// transforms of the layer above.
func (e *engine) canonicalizeXPathNodeSet(cur *xmldom.Node, endnode *xmldom.Node) error {
	ns := newNSStack()
	if cur == nil {
		return nil
	}
	if e.isVisibleInt(cur) == -1 {
		return e.filterErr
	}
	if cur.Kind == xmldom.Element {
		e.getParentNameSpaces(cur, ns)
	}
	var sibling, parentNode *xmldom.Node
	documentLevel := nodeBeforeDocumentElement
	currentNodeIsVisible := false
	for {
		switch cur.Kind {
		case xmldom.Attribute:
			return fmt.Errorf("xmlc14n: illegal node type during traversal: %s", cur.Kind)

		case xmldom.Document:
			ns.outputNodePush()
			sibling = cur.FirstChild

		case xmldom.Comment:
			if e.includeComments && e.isVisibleDO(cur, ns.level()) == 1 {
				e.outputComment(cur.Value, documentLevel)
			}

		case xmldom.ProcInst:
			if e.isVisible(cur) {
				e.outputPI(cur.Name.Local, cur.Value, documentLevel)
			}

		case xmldom.Text, xmldom.CDATA:
			if e.isVisible(cur) {
				writeTextEscaped(e.w, cur.Value)
				// Santuario consumes the whole run of adjacent character data here, which
				// matters because the run's later nodes may be invisible in the node set and
				// are emitted anyway.
				for next := cur.NextSibling; next != nil && (next.Kind == xmldom.Text || next.Kind == xmldom.CDATA); next = next.NextSibling {
					writeTextEscaped(e.w, next.Value)
					cur = next
					sibling = cur.NextSibling
				}
			}

		case xmldom.Element:
			documentLevel = nodeNotBeforeOrAfterDocumentElement
			var name string
			// A -1 answer prunes the subtree outright: no namespace frame is pushed and
			// outputAttributes never runs, so nothing below can be reached or declared. That
			// is how the enveloped-signature filter removes a ds:Signature.
			if i := e.isVisibleDO(cur, ns.level()); i == -1 {
				sibling = cur.NextSibling
				break
			} else {
				currentNodeIsVisible = i == 1
			}
			if currentNodeIsVisible {
				ns.outputNodePush()
				name = cur.Name.QName()
				e.w.WriteByte('<')
				e.w.WriteString(name)
			} else {
				ns.push()
			}

			if err := e.outputAttributes(cur, ns); err != nil {
				return err
			}

			if currentNodeIsVisible {
				e.w.WriteByte('>')
			}
			sibling = cur.FirstChild
			if sibling == nil {
				if currentNodeIsVisible {
					e.w.WriteString("</")
					e.w.WriteString(name)
					e.w.WriteByte('>')
					ns.outputNodePop()
				} else {
					ns.pop()
				}
				if parentNode != nil {
					sibling = cur.NextSibling
				}
			} else {
				parentNode = cur
			}
		}

		for sibling == nil && parentNode != nil {
			if e.isVisible(parentNode) {
				e.w.WriteString("</")
				e.w.WriteString(parentNode.Name.QName())
				e.w.WriteByte('>')
				ns.outputNodePop()
			} else {
				ns.pop()
			}
			if parentNode == endnode {
				return nil
			}
			sibling = parentNode.NextSibling
			parentNode = parentNode.Parent
			if parentNode == nil || parentNode.Kind != xmldom.Element {
				documentLevel = nodeAfterDocumentElement
				parentNode = nil
			}
		}
		if sibling == nil {
			return nil
		}
		cur = sibling
		sibling = cur.NextSibling
	}
}

// getParentNameSpaces ports CanonicalizerBase.getParentNameSpaces: fill the symbol table with
// the declarations in scope at el, outermost ancestor first, then inject the explicit xmlns=""
// entry when the innermost default binding an ancestor left behind is empty.
func (e *engine) getParentNameSpaces(el *xmldom.Node, ns *nsStack) {
	p := el.Parent
	if p == nil || p.Kind != xmldom.Element {
		return
	}
	var parents []*xmldom.Node
	for parent := p; parent != nil && parent.Kind == xmldom.Element; parent = parent.Parent {
		parents = append(parents, parent)
	}
	for i := len(parents) - 1; i >= 0; i-- {
		e.handleParent(parents[i], ns)
	}
	if uri, ok := ns.getMappingWithoutRendered(xmlnsPrefix); ok && uri == "" {
		ns.addMappingAndRender(xmlnsPrefix, "", true)
	}
}

// handleParent ports CanonicalizerBase.handleParent and, for inclusive c14n, the
// Canonicalizer20010315 override that also feeds the xml:* stack.
func (e *engine) handleParent(el *xmldom.Node, ns *nsStack) {
	if e.physical {
		// CanonicalizerPhysical overrides handleParent to do nothing: the physical method
		// inherits no ancestor context whatsoever. getParentNameSpaces still runs, but with
		// an empty symbol table its trailing xmlns="" injection cannot fire either.
		return
	}
	if len(el.Attrs) == 0 && el.Name.Space == "" {
		return
	}
	inclusive := !e.exclusive && !e.physical
	if inclusive {
		e.xmlAttrs.push(-1)
	}
	for _, attr := range el.Attrs {
		name, value := attr.Name.Local, attr.Value
		switch attr.Name.Space {
		case xmldom.XMLNSNamespace:
			// The default mapping for the xml prefix is never rendered.
			if name != xmlPrefix || value != xmldom.XMLNamespace {
				ns.addMapping(name, value, true)
			}
		case xmldom.XMLNamespace:
			if inclusive && (!e.c14n11 || name != "id") {
				e.xmlAttrs.addXmlnsAttr(attrOf(attr))
			}
		}
	}
	// Santuario synthesizes the declaration for the element's own namespace, so that a DOM
	// built by hand - without the declaration attribute - still contributes its binding. For a
	// parsed document the binding is already present and addMapping is a no-op.
	if el.Name.Space != "" {
		prefix := el.Name.Prefix
		if prefix == "" {
			prefix = xmlnsPrefix
		}
		ns.addMapping(prefix, el.Name.Space, true)
	}
}

// outputComment ports CanonicalizerBase.outputCommentToWriter.
func (e *engine) outputComment(data string, position int) {
	position = e.position(position)
	if position == nodeAfterDocumentElement {
		e.w.WriteByte('\n')
	}
	e.w.WriteString("<!--")
	writeCarriageReturnEscaped(e.w, data)
	e.w.WriteString("-->")
	if position == nodeBeforeDocumentElement {
		e.w.WriteByte('\n')
	}
}

// outputPI ports CanonicalizerBase.outputPItoWriter. A PI with no data gets no separating
// space, which is why "<?a?>" and "<?b c?>" are not interchangeable.
func (e *engine) outputPI(target, data string, position int) {
	position = e.position(position)
	if position == nodeAfterDocumentElement {
		e.w.WriteByte('\n')
	}
	e.w.WriteString("<?")
	writeCarriageReturnEscaped(e.w, target)
	if len(data) > 0 {
		e.w.WriteByte(' ')
		writeCarriageReturnEscaped(e.w, data)
	}
	e.w.WriteString("?>")
	if position == nodeBeforeDocumentElement {
		e.w.WriteByte('\n')
	}
}

// position applies CanonicalizerPhysical's override of both writers, which forces "inside the
// document element" and so suppresses every positional newline.
func (e *engine) position(position int) int {
	if e.physical {
		return nodeNotBeforeOrAfterDocumentElement
	}
	return position
}

// attrOf converts an xmldom attribute node into the emitter's value type.
func attrOf(attr *xmldom.Node) outAttr {
	return outAttr{
		space: attr.Name.Space,
		local: attr.Name.Local,
		qname: attr.Name.QName(),
		value: attr.Value,
	}
}
