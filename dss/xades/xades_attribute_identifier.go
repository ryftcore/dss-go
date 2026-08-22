// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESAttributeIdentifier.java
// (DSS 6.5.RC1).
//
// DomUtils.serializeNode(Node) is xml/utils.DomUtilsSerializeNode(*xmldom.Node) ([]byte, error);
// unlike Java (which never fails here because Santuario's XMLUtils never throws on a live DOM
// node), the Go transformer can fail, so a serialization error is treated the same way DSSXMLUtils
// treats other "cannot happen in practice" DOM failures elsewhere in this package: it panics with
// the same DSSException message Java would have raised had serializeNode declared a checked
// exception, per PORTING.md's requireNonNull/throws-in-constructor convention.
package xades

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/spi/validation/identifier"
	"github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESAttributeIdentifier represents an identifier of a XAdES Attribute. Port of the class
// AttributeIdentifier, extending identifier.SignatureAttributeIdentifier.
type AttributeIdentifier struct {
	identifier.SignatureAttributeIdentifier
}

// newXAdESAttributeIdentifier is the port of the package-private XAdESAttributeIdentifier(byte[])
// constructor.
func newXAdESAttributeIdentifier(data []byte) *AttributeIdentifier {
	return &AttributeIdentifier{
		SignatureAttributeIdentifier: identifier.NewSignatureAttributeIdentifierBase("XAdESAttributeIdentifier", data),
	}
}

// AttributeIdentifierBuild builds the AttributeIdentifier from the given property
// node. Port of the static build(Node).
//
// Panics with the Java DSSException message when the node cannot be serialized (Java's
// try-with-resources catches only IOException from the in-memory ByteArrayOutputStream/
// DataOutputStream pair, which cannot actually fail; DomUtils.serializeNode(node) itself does not
// declare a checked exception in Java, so this mirrors an unchecked failure here too).
func AttributeIdentifierBuild(node *xmldom.Node) *AttributeIdentifier {
	binaries, err := xadesAttributeIdentifierGetBinaries(node)
	if err != nil {
		panic(fmt.Sprintf("Unable to build a XAdES Attribute Identifier : %s", err.Error()))
	}
	order := xadesAttributeIdentifierGetOrder(node)

	var buf bytes.Buffer
	buf.Write(binaries)
	var orderBytes [4]byte
	binary.BigEndian.PutUint32(orderBytes[:], uint32(order))
	buf.Write(orderBytes[:])

	return newXAdESAttributeIdentifier(buf.Bytes())
}

// xadesAttributeIdentifierGetBinaries ports the private static getBinaries(Node).
func xadesAttributeIdentifierGetBinaries(node *xmldom.Node) ([]byte, error) {
	return utils.DomUtilsSerializeNode(node)
}

// xadesAttributeIdentifierGetOrder ports the private static getOrder(Node): the zero-based
// position of node among its parent's children, 0 when node has no parent.
func xadesAttributeIdentifierGetOrder(node *xmldom.Node) int {
	parent := node.Parent
	if parent != nil {
		ii := 0
		for child := parent.FirstChild; child != nil; child = child.NextSibling {
			if child == node {
				return ii
			}
			ii++
		}
	}
	return 0
}
