// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pdf/modifications/PdfObjectTree.java
// (DSS 6.5.RC1).
//
// eu.europa.esig.dss.pdf.modifications is part of the eu.europa.esig.dss.pdf module that landed
// in no s5b manifest (see pdf_object.go's header). This file names the type PdfObjectTree
// (Java's own name), rather than the bare "ObjectTree" a landed sibling's forward-dependency
// header speculatively assumed (pdf_signature_dictionary.go); see that file's comment, updated
// during integration to match.
package pades

import (
	"strconv"
	"strings"
)

// PdfObjectTree represents a PDF object chain from a root to the current object.
// Port of the PdfObjectTree class.
type PdfObjectTree struct {
	// keyChain is the key chain.
	keyChain []string

	// refChain holds the processed references.
	refChain []PdfObjectKey

	// sb is used to build a user-friendly string.
	sb strings.Builder
}

// NewPdfObjectTree builds an object tree, optionally starting with a root key. Java splits this
// into a no-arg constructor and a String-arg constructor; Go collapses them, an empty key
// standing in for "no starting key".
// Port of PdfObjectTree() / PdfObjectTree(String).
func NewPdfObjectTree(key string) *PdfObjectTree {
	t := &PdfObjectTree{}
	if key != "" {
		t.AddKey(key)
	}
	return t
}

// Copy creates a copy of the object tree (changes within a copy will not affect the original).
// Port of #copy.
func (t *PdfObjectTree) Copy() *PdfObjectTree {
	c := &PdfObjectTree{
		keyChain: append([]string(nil), t.keyChain...),
		refChain: append([]PdfObjectKey(nil), t.refChain...),
	}
	c.sb.WriteString(t.sb.String())
	return c
}

// AddKey adds a key. Port of #addKey.
func (t *PdfObjectTree) AddKey(key string) {
	t.keyChain = append(t.keyChain, key)
	if t.sb.Len() != 0 {
		t.sb.WriteString(" ")
	}
	t.sb.WriteString("/")
	t.sb.WriteString(key)
}

// AddReference adds a PDF object key. Port of #addReference.
func (t *PdfObjectTree) AddReference(objectKey PdfObjectKey) {
	t.refChain = append(t.refChain, objectKey)
	if t.sb.Len() != 0 {
		t.sb.WriteString(" ")
	}
	t.sb.WriteString(strconv.FormatInt(objectKey.Number(), 10))
	t.sb.WriteString(" 0 R")
}

// SetStream specifies that a stream has been processed. Port of #setStream.
func (t *PdfObjectTree) SetStream() {
	if t.sb.Len() != 0 {
		t.sb.WriteString(" ")
	}
	t.sb.WriteString("stream")
}

// KeyChain gets the complete key chain. Port of #getKeyChain.
func (t *PdfObjectTree) KeyChain() []string { return t.keyChain }

// ChainDeepness returns the deepness of the current objects chain. Port of #getChainDeepness.
func (t *PdfObjectTree) ChainDeepness() int { return len(t.keyChain) }

// LastKey returns the last key. Port of #getLastKey.
func (t *PdfObjectTree) LastKey() string {
	if len(t.keyChain) == 0 {
		return ""
	}
	return t.keyChain[len(t.keyChain)-1]
}

// IsProcessedReference checks whether a reference to the given object has already been
// processed in this tree. Port of #isProcessedReference.
func (t *PdfObjectTree) IsProcessedReference(objectKey PdfObjectKey) bool {
	if objectKey == nil {
		return false
	}
	for _, ref := range t.refChain {
		if ref == objectKey {
			return true
		}
	}
	return false
}

// Equals ports #equals: structural equality of the key chain, the reference chain and the
// built string.
func (t *PdfObjectTree) Equals(other *PdfObjectTree) bool {
	if t == other {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	if len(t.keyChain) != len(other.keyChain) || len(t.refChain) != len(other.refChain) {
		return false
	}
	for i, key := range t.keyChain {
		if key != other.keyChain[i] {
			return false
		}
	}
	for i, ref := range t.refChain {
		if ref != other.refChain[i] {
			return false
		}
	}
	return t.sb.String() == other.sb.String()
}

// String ports #toString.
func (t *PdfObjectTree) String() string {
	return t.sb.String()
}
