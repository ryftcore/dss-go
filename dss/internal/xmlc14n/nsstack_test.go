package xmlc14n

import "testing"

// render drives addMappingAndRender the way the inclusive emitter does and reports what would
// be written, or "" for "nothing to render".
func render(s *nsStack, prefix, uri string) string {
	a, ok := s.addMappingAndRender(prefix, uri, true)
	if !ok {
		return ""
	}
	return a.qname + "=" + a.value
}

func TestNSStackSuppressesSuperfluousDeclarations(t *testing.T) {
	// corpus/ns-superfluous.xml declares xmlns:p="urn:1" at three depths and
	// golden/ns-superfluous.c14n10.doc emits it once:
	//   <r xmlns:p="urn:1" Id="r"><p:c Id="c"><p:d Id="d"></p:d></p:c></r>
	ns := newNSStack()
	ns.outputNodePush()
	if got := render(ns, "p", "urn:1"); got != "xmlns:p=urn:1" {
		t.Fatalf("first declaration rendered %q", got)
	}
	ns.outputNodePush()
	if got := render(ns, "p", "urn:1"); got != "" {
		t.Errorf("redeclaration at depth 2 rendered %q, want nothing", got)
	}
	ns.outputNodePush()
	if got := render(ns, "p", "urn:1"); got != "" {
		t.Errorf("redeclaration at depth 3 rendered %q, want nothing", got)
	}
	ns.outputNodePop()
	ns.outputNodePop()
	ns.outputNodePop()
}

func TestNSStackRebind(t *testing.T) {
	// corpus/ns-rebind.xml rebinds p to urn:2 and back to urn:1; every change is rendered,
	// and the frames restore the outer binding on pop.
	ns := newNSStack()
	ns.outputNodePush()
	if got := render(ns, "p", "urn:1"); got != "xmlns:p=urn:1" {
		t.Fatalf("r rendered %q", got)
	}
	ns.outputNodePush()
	if got := render(ns, "p", "urn:2"); got != "xmlns:p=urn:2" {
		t.Errorf("rebind rendered %q", got)
	}
	ns.outputNodePush()
	if got := render(ns, "p", "urn:1"); got != "xmlns:p=urn:1" {
		t.Errorf("rebind back rendered %q", got)
	}
	ns.outputNodePop()
	ns.outputNodePop()
	if e := ns.get("p"); e == nil || e.uri != "urn:1" || !e.rendered {
		t.Errorf("after two pops the binding is %+v, want the outer urn:1, rendered", e)
	}
	ns.outputNodePop()
}

func TestNSStackApexUndeclaresDefaultNamespace(t *testing.T) {
	// The apex of a subtree inherits nothing, so the empty default namespace is already in
	// force and no xmlns="" is emitted. The initial entry - uri "", rendered, lastrendered ""
	// - is the whole mechanism (design 2.4, golden ns-default-apex.c14n10.apex-c).
	ns := newNSStack()
	if _, ok := ns.getMappingWithoutRendered(xmlnsPrefix); ok {
		t.Error("the initial xmlns entry must be rendered, so getMappingWithoutRendered reports nothing")
	}
	ns.outputNodePush()
	if got := render(ns, xmlnsPrefix, ""); got != "" {
		t.Errorf("apex undeclaration rendered %q, want nothing", got)
	}
}

func TestNSStackInheritedUndeclarationIsNotFlushed(t *testing.T) {
	// The mirror case: an ancestor chain whose innermost default binding is "". After
	// getParentNameSpaces injects the null node, the apex must emit nothing - golden
	// ns-default-undeclare.c14n10.apex-d is <d Id="d"></d> - while an apex under a re-declared
	// default namespace does emit it (apex-f is <f xmlns="urn:a" Id="f"></f>).
	ns := newNSStack()
	ns.addMapping(xmlnsPrefix, "urn:a", true) // <r xmlns="urn:a">
	ns.addMapping(xmlnsPrefix, "", true)      // <c xmlns="">
	// The entry is already marked rendered: addMapping carries lastrendered forward from the
	// initial xmlns entry, whose value is "", and a new binding whose URI equals lastrendered
	// counts as rendered. getParentNameSpaces' nullNode injection is therefore a no-op on this
	// path - it can only fire if something rendered the default namespace first, which nothing
	// does before the apex - and it is ported for fidelity rather than for effect.
	if uri, ok := ns.getMappingWithoutRendered(xmlnsPrefix); ok {
		t.Fatalf("getMappingWithoutRendered = %q, %v; want nothing, the entry counts as rendered", uri, ok)
	}
	if got := render(ns, xmlnsPrefix, ""); got != "" {
		t.Errorf("the nullNode injection rendered %q, want nothing", got)
	}

	var set attrSet
	ns.getUnrenderedNodes(&set)
	if set.len() != 0 {
		t.Errorf("flushed %d declarations, want none", set.len())
	}

	ns2 := newNSStack()
	ns2.addMapping(xmlnsPrefix, "urn:a", true)
	ns2.addMapping(xmlnsPrefix, "", true)
	ns2.addMapping(xmlnsPrefix, "urn:a", true) // <e xmlns="urn:a">
	var set2 attrSet
	ns2.getUnrenderedNodes(&set2)
	if got := emit(t, &set2); got != ` xmlns="urn:a"` {
		t.Errorf("flushed %q, want the re-declared default namespace", got)
	}
}

func TestNSStackGetUnrenderedNodesSkipsEmptyURIs(t *testing.T) {
	// SymbMap.entrySet() drops entries whose URI is empty, so an inherited xmlns="" can never
	// be flushed onto an apex.
	ns := newNSStack()
	ns.addMapping("p", "urn:1", true)
	ns.addMapping(xmlnsPrefix, "", true)
	var set attrSet
	ns.getUnrenderedNodes(&set)
	if got := emit(t, &set); got != ` xmlns:p="urn:1"` {
		t.Errorf("flushed %q, want only xmlns:p", got)
	}
	// Everything flushed is now rendered, so a second flush writes nothing.
	var again attrSet
	ns.getUnrenderedNodes(&again)
	if again.len() != 0 {
		t.Errorf("second flush produced %d declarations", again.len())
	}
}

func TestNSStackGetMapping(t *testing.T) {
	// Exclusive c14n renders through getMapping: unknown prefixes (an InclusiveNamespaces
	// PrefixList entry that names nothing) report nothing, and a prefix already rendered with
	// the same URI is not rendered twice.
	ns := newNSStack()
	ns.outputNodePush()
	if _, ok := ns.getMapping("zz"); ok {
		t.Error("unknown prefix reported a mapping")
	}
	ns.addMapping("p", "urn:1", true)
	a, ok := ns.getMapping("p")
	if !ok || a.qname != "xmlns:p" || a.value != "urn:1" {
		t.Fatalf("getMapping(p) = %+v, %v", a, ok)
	}
	if _, ok := ns.getMapping("p"); ok {
		t.Error("second getMapping(p) rendered again")
	}
	ns.outputNodePush()
	ns.addMapping("p", "urn:1", true) // superfluous redeclaration deeper down
	if _, ok := ns.getMapping("p"); ok {
		t.Error("redeclaration to the same URI rendered again")
	}
}

func TestNSStackRemoveMappings(t *testing.T) {
	ns := newNSStack()
	ns.push()
	ns.addMapping("p", "urn:1", true)
	ns.removeMapping("p")
	if ns.get("p") != nil {
		t.Error("removeMapping left the binding in place")
	}
	ns.pop()

	ns.push()
	ns.addMapping("q", "urn:2", true)
	ns.removeMappingIfNotRender("q")
	if ns.get("q") != nil {
		t.Error("removeMappingIfNotRender kept an unrendered binding")
	}
	ns.addMapping("q", "urn:2", true)
	ns.getMapping("q") // now rendered
	ns.removeMappingIfNotRender("q")
	if ns.get("q") == nil {
		t.Error("removeMappingIfNotRender dropped a rendered binding")
	}
	// removeMappingIfRender always reports false upstream, whatever it does to the table.
	if ns.removeMappingIfRender("q") {
		t.Error("removeMappingIfRender reported true; Santuario always returns false")
	}
	if ns.get("q") != nil {
		t.Error("removeMappingIfRender kept a rendered binding")
	}
	ns.pop()
}

func TestNSStackFramesUndoEveryWrite(t *testing.T) {
	ns := newNSStack()
	ns.addMapping("p", "urn:outer", true)
	ns.push()
	ns.addMapping("p", "urn:inner", true)
	ns.addMapping("q", "urn:q", true)
	ns.removeMapping("p")
	if ns.level() != 1 {
		t.Fatalf("level = %d, want 1", ns.level())
	}
	ns.pop()
	if e := ns.get("p"); e == nil || e.uri != "urn:outer" {
		t.Errorf("p = %+v after pop, want urn:outer", e)
	}
	if ns.get("q") != nil {
		t.Error("q survived the pop")
	}
	if ns.level() != 0 {
		t.Errorf("level = %d after pop, want 0", ns.level())
	}
}
