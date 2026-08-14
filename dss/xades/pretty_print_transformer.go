// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/signature/PrettyPrintTransformer.java (DSS 6.5.RC1).
//
// A self-contained DOM transformation: no Santuario, no canonicalization, no serialization. The
// only DSS dependency is DSSXMLUtils.TRANSFORMER_INDENT_NUMBER, the default indent amount.
//
// The indent walk keeps upstream's exact shape, including the detail that indent() advances
// childNode with getNextSibling() after recursing, so an inserted indent text node is stepped over
// rather than revisited, and the "skip" flag that suppresses indentation next to pre-existing text
// when keepOriginalIndents is on. Java's Node#getOwnerDocument on the clone is xmldom's
// OwnerDocument, which answers the document a detached subtree was cloned out of - so the indent
// text nodes are created against the same document Java creates them against.
package xades

import (
	"strings"

	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/utils"
)

const (
	// prettyPrintTransformerNewLine is the new line character. Port of the private NEW_LINE.
	prettyPrintTransformerNewLine = "\n"

	// prettyPrintTransformerSpace is the whitespace character. Port of the private SPACE.
	prettyPrintTransformerSpace = " "
)

// PrettyPrintTransformer performs pretty-print transformations on an XML signature.
type PrettyPrintTransformer struct {
	// ownerDocument is the parent document.
	ownerDocument *xmldom.Node

	// indentAmount is the indent amount (4 by default).
	indentAmount int

	// keepOriginalIndents indicates whether original indents are to be preserved.
	keepOriginalIndents bool
}

// NewPrettyPrintTransformer is the default constructor. Port of the empty constructor, including
// its two field initializers (indentAmount = DSSXMLUtils.TRANSFORMER_INDENT_NUMBER,
// keepOriginalIndents = true).
func NewPrettyPrintTransformer() *PrettyPrintTransformer {
	return &PrettyPrintTransformer{
		indentAmount:        DSSXMLUtilsTransformerIndentNumber,
		keepOriginalIndents: true,
	}
}

// SetIndentAmount configures the amount of spaces to add. Port of #setIndentAmount.
func (t *PrettyPrintTransformer) SetIndentAmount(indentAmount int) *PrettyPrintTransformer {
	t.indentAmount = indentAmount
	return t
}

// SetKeepOriginalIndents sets whether the original indents are to be kept.
// Default: true (original indents are not modified). Port of #setKeepOriginalIndents.
func (t *PrettyPrintTransformer) SetKeepOriginalIndents(keepOriginalIndents bool) *PrettyPrintTransformer {
	t.keepOriginalIndents = keepOriginalIndents
	return t
}

// Transform indents the provided node and returns the indented copy. Port of #transform.
func (t *PrettyPrintTransformer) Transform(nodeToTransform *xmldom.Node) *xmldom.Node {
	clonedNode := nodeToTransform.Clone(true)
	t.ownerDocument = clonedNode.OwnerDocument()
	return t.indent(clonedNode, 1)
}

// indent ports the private indent.
func (t *PrettyPrintTransformer) indent(nodeToTransform *xmldom.Node, level int) *xmldom.Node {
	if prettyPrintTransformerHasElementChilds(nodeToTransform) {
		indentString := t.indentString(level)
		skip := false
		childNode := nodeToTransform.FirstChild
		for childNode != nil {
			if xmldom.Text == childNode.Kind {
				if t.keepOriginalIndents {
					skip = true
				} else {
					// Upstream DISCARDS this return value, and childNode stays pointing at the
					// text node just removed. Xerces clears a removed node's sibling links, so
					// the childNode.getNextSibling() below answers null and the loop stops after
					// the first stripped text run - which is exactly what the Java oracle emits
					// (testdata/sign-a-builder.txt, case pretty-drop-indents-preindented, keeps
					// the tail of the child list unindented). xmldom.unlink clears the same links,
					// so the Go loop stops at the same place. Advancing to the returned node
					// instead would strip every indent and change the output.
					prettyPrintTransformerRemoveSiblingIndents(childNode)
				}

			} else {
				if !skip && utils.IsStringNotEmpty(indentString) {
					indentNode := t.indentNode(indentString)
					nodeToTransform.InsertBefore(indentNode, childNode)
				}
				skip = false
				childNode = t.indent(childNode, level+1)
			}
			childNode = childNode.NextSibling
		}
		if !skip {
			indentString = t.indentString(level - 1)
			indentNode := t.indentNode(indentString)
			nodeToTransform.AppendChild(indentNode)
		}
	}
	return nodeToTransform
}

// prettyPrintTransformerRemoveSiblingIndents ports the private removeSiblingIndents. A trailing
// text node makes the recursion reach a nil node and dereference it, which is upstream's own
// NullPointerException in the same situation, reproduced rather than papered over.
func prettyPrintTransformerRemoveSiblingIndents(node *xmldom.Node) *xmldom.Node {
	parentNode := node.Parent
	if xmldom.Text == node.Kind {
		nextSibling := node.NextSibling
		parentNode.RemoveChild(node)
		return prettyPrintTransformerRemoveSiblingIndents(nextSibling)
	}
	return node
}

// indentNode ports the private getIndentNode.
func (t *PrettyPrintTransformer) indentNode(indentString string) *xmldom.Node {
	// Java's ownerDocument.createTextNode; xmldom text nodes carry no document of their own until
	// they are inserted, which is what the caller does immediately.
	_ = t.ownerDocument
	return xmldom.NewText(indentString)
}

// indentString ports the private getIndentString.
func (t *PrettyPrintTransformer) indentString(level int) string {
	spacesExpected := level * t.indentAmount
	var stringBuilder strings.Builder
	stringBuilder.WriteString(prettyPrintTransformerNewLine)
	for ii := 0; ii < spacesExpected; ii++ {
		stringBuilder.WriteString(prettyPrintTransformerSpace)
	}
	return stringBuilder.String()
}

// prettyPrintTransformerHasElementChilds ports the private hasElementChilds.
func prettyPrintTransformerHasElementChilds(node *xmldom.Node) bool {
	if node == nil {
		return false
	}
	for _, item := range node.Children() {
		if xmldom.Element == item.Kind {
			return true
		}
	}
	return false
}
