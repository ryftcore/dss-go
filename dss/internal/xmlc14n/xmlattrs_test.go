package xmlc14n

import "testing"

func xmlAttr(local, value string) outAttr {
	return outAttr{space: xmlNamespace, local: local, qname: "xml:" + local, value: value}
}

func TestXMLAttrStackSubtreeOutermostWins(t *testing.T) {
	// corpus/xmlattrs-lang-space.xml with apex <t>:
	//   <r xml:lang="en" xml:space="preserve"><m xml:lang="fr"><t a="1"/></m>...
	// golden/xmlattrs-lang-space.c14n10.apex-t is
	//   <t Id="t" a="1" xml:lang="en" xml:space="preserve"></t>
	// - the OUTERMOST xml:lang wins. handleParent pushes level -1 for every ancestor, so all
	// of them share one level, their attributes are appended outermost-first, and the loa map
	// keeps the first occurrence per QName.
	s := &xmlAttrStack{}
	s.push(-1)
	s.addXmlnsAttr(xmlAttr("lang", "en"))
	s.addXmlnsAttr(xmlAttr("space", "preserve"))
	s.push(-1)
	s.addXmlnsAttr(xmlAttr("lang", "fr"))

	var set attrSet
	s.getXmlnsAttr(&set)
	if got := emit(t, &set); got != ` xml:lang="en" xml:space="preserve"` {
		t.Errorf("flushed %q, want the outermost xml:lang and xml:space", got)
	}
}

func TestXMLAttrStackNothingToFlush(t *testing.T) {
	// An apex with no xml:* ancestors: golden/attr-order.c14n10.apex-r carries none.
	s := &xmlAttrStack{}
	s.push(-1)
	var set attrSet
	s.getXmlnsAttr(&set)
	if set.len() != 0 {
		t.Errorf("flushed %d attributes, want none", set.len())
	}
}

func TestXMLAttrStackParentRenderedShortcut(t *testing.T) {
	// In the node-set walk the levels are real depths. Once the immediately enclosing level
	// has been rendered, its attributes are already in the output and a child flushes only
	// its own - here none, so nothing at all.
	s := &xmlAttrStack{}
	s.push(1)
	s.addXmlnsAttr(xmlAttr("lang", "en"))
	var first attrSet
	s.getXmlnsAttr(&first) // renders level 1
	if got := emit(t, &first); got != ` xml:lang="en"` {
		t.Fatalf("level 1 flushed %q", got)
	}

	s.push(2)
	var second attrSet
	s.getXmlnsAttr(&second)
	if second.len() != 0 {
		t.Errorf("level 2 flushed %q, want nothing: the parent already rendered it", emit(t, &second))
	}

	// A child that carries its own xml:* attribute is not on the shortcut path, because the
	// level getXmlnsAttr inspects is then the child's own unrendered one: it flushes the whole
	// chain of unrendered ancestors again, which is Santuario's behaviour and not an
	// optimization to reproduce selectively.
	s2 := &xmlAttrStack{}
	s2.push(1)
	s2.addXmlnsAttr(xmlAttr("lang", "en"))
	s2.push(2)
	s2.addXmlnsAttr(xmlAttr("space", "preserve"))
	var third attrSet
	s2.getXmlnsAttr(&third)
	if got := emit(t, &third); got != ` xml:lang="en" xml:space="preserve"` {
		t.Errorf("level 2 with its own attribute flushed %q", got)
	}
}

func TestXMLAttrStackPushUnwindsLevels(t *testing.T) {
	// push(level) drops every level at or below the new one, so a sibling does not inherit its
	// predecessor's xml:* attributes.
	s := &xmlAttrStack{}
	s.push(1)
	s.addXmlnsAttr(xmlAttr("lang", "en"))
	s.push(2)
	s.addXmlnsAttr(xmlAttr("space", "preserve"))
	s.push(2) // the next element at the same depth
	if len(s.levels) != 1 || s.levels[0].level != 1 {
		t.Fatalf("levels = %v, want only the level-1 frame", s.levels)
	}
	var set attrSet
	s.getXmlnsAttr(&set)
	if got := emit(t, &set); got != ` xml:lang="en"` {
		t.Errorf("flushed %q, want only the enclosing xml:lang", got)
	}
}
