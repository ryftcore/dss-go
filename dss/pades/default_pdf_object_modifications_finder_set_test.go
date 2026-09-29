package pades

import (
	"math/rand"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/pdf"
)

// linearObjectModificationSet is the collection objectModificationSet replaced: the same
// insertion-ordered, duplicate-free contract, with a scan of every earlier entry per Add.
type linearObjectModificationSet struct {
	items []ObjectModification
}

func (s *linearObjectModificationSet) Add(objectModification ObjectModification) {
	for _, existing := range s.items {
		if existing.Equals(objectModification) {
			return
		}
	}
	s.items = append(s.items, objectModification)
}

func testObjectTree(keys []string, refs ...pdf.ObjectKey) *PdfObjectTree {
	tree := NewPdfObjectTree("")
	for _, key := range keys {
		tree.AddKey(key)
	}
	for _, ref := range refs {
		tree.AddReference(NewNativePdfObjectKey(ref))
	}
	return tree
}

// Equal object trees and types collapse to one entry, in first-seen order; anything that
// ObjectModification.Equals tells apart stays, including two trees that render to the same text
// (a reference's generation is not part of it).
func TestObjectModificationSetKeepsEqualsSemantics(t *testing.T) {
	set := &objectModificationSet{}
	first := NewObjectModificationCreate(testObjectTree([]string{"Catalog", "Pages"}), nil)
	set.Add(first)
	set.Add(NewObjectModificationCreate(testObjectTree([]string{"Catalog", "Pages"}), nil))               // duplicate
	set.Add(NewObjectModificationDelete(testObjectTree([]string{"Catalog", "Pages"}), nil))               // other type
	set.Add(NewObjectModificationCreate(testObjectTree([]string{"Catalog", "Pages", "Kids"}), nil))       // other tree
	set.Add(NewObjectModificationCreate(testObjectTree([]string{"Catalog"}, pdf.ObjectKey{Num: 4}), nil)) // ref, gen 0
	set.Add(NewObjectModificationCreate(testObjectTree([]string{"Catalog"}, pdf.ObjectKey{Num: 4, Gen: 1}), nil))
	set.Add(NewObjectModificationCreate(testObjectTree([]string{"Catalog"}, pdf.ObjectKey{Num: 4}), nil)) // duplicate
	set.Add(NewObjectModificationCreate(nil, nil))
	set.Add(NewObjectModificationCreate(nil, nil)) // duplicate

	if got := len(set.items); got != 6 {
		t.Fatalf("len(items) = %d, want 6", got)
	}
	if set.items[0].ObjectTree() != first.ObjectTree() {
		t.Error("the first-seen entry was not kept in front")
	}
	if set.items[1].ActionType() != enumerations.PdfObjectModificationTypeDeletion {
		t.Errorf("second entry = %v, want the deletion", set.items[1].ActionType())
	}
}

func TestObjectModificationSetMatchesLinearScan(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	keys := []string{"Catalog", "Pages", "Kids", "AcroForm", "Fields"}
	types := []func(*PdfObjectTree, PdfObject) ObjectModification{
		NewObjectModificationCreate,
		func(tree *PdfObjectTree, o PdfObject) ObjectModification { return NewObjectModificationDelete(tree, o) },
		func(tree *PdfObjectTree, o PdfObject) ObjectModification {
			return NewObjectModificationModify(tree, o, o)
		},
	}

	set := &objectModificationSet{}
	reference := &linearObjectModificationSet{}
	for i := 0; i < 5000; i++ {
		tree := NewPdfObjectTree("")
		for depth := rng.Intn(4); depth >= 0; depth-- {
			if rng.Intn(3) == 0 {
				tree.AddReference(NewNativePdfObjectKey(pdf.ObjectKey{Num: int64(1 + rng.Intn(3)), Gen: uint16(rng.Intn(2))}))
			} else {
				tree.AddKey(keys[rng.Intn(len(keys))])
			}
		}
		modification := types[rng.Intn(len(types))](tree, nil)
		set.Add(modification)
		reference.Add(modification)
	}

	if len(set.items) != len(reference.items) || len(set.items) < 20 {
		t.Fatalf("%d entries, the linear scan keeps %d (want a non-trivial number)", len(set.items), len(reference.items))
	}
	for i := range set.items {
		if set.items[i].ObjectTree() != reference.items[i].ObjectTree() ||
			set.items[i].ActionType() != reference.items[i].ActionType() {
			t.Fatalf("entry %d differs from the linear scan", i)
		}
	}
}
