#!/usr/bin/env python3
"""Writes internal/xmlc14n/testdata/corpus/ - the adversarial C14N corpus - and gen/cases.txt,
the list of (document, scope) pairs the Java oracle canonicalizes.

Every document is written as raw bytes so that BOMs, lone CRs, CRLFs and the ISO-8859-1
document survive source control unchanged; nothing here may be reformatted by hand or by an
editor. Each document targets a named rule of the design document (scratchpad/XML_DESIGN.md
sections 2.4 to 2.8) and is deliberately small so that a golden diff is readable.

Apex-scope cases name their apex by an attribute whose local name is "Id"/"ID"/"id"
(case-insensitively) or by xml:id; both the oracle (C14nOracle.java) and the Go test select
"the first element in document order carrying such an attribute with that value", which needs
no XPath on either side. Apex values are unique across a document.

    python3 gen/corpus.py testdata/corpus gen/cases.txt

corpus/negative/ holds documents that must be rejected at parse time; the oracle records what
Java makes of each so the Go parser (internal/xmldom) has an oracle for its own rejections.
"""
import os
import sys

CORPUS = sys.argv[1] if len(sys.argv) > 1 else "corpus"
CASES = sys.argv[2] if len(sys.argv) > 2 else "cases.txt"

DSS_UPSTREAM = os.environ.get("DSS", "")  # set DSS to your upstream DSS checkout path

docs = {}   # name -> bytes
cases = []  # (doc, scope-slug, apex, prefixlist)


def doc(name, data):
    docs[name] = data if isinstance(data, bytes) else data.encode("utf-8")


def case(name, slug, apex="", prefixes=""):
    cases.append((name, slug, apex, prefixes))


def whole(name):
    """Canonicalize the whole document (XMLCanonicalizer.canonicalize(byte[]))."""
    case(name, "doc")


def apex(name, *ids):
    """Canonicalize the subtree rooted at each apex (XMLCanonicalizer.canonicalize(Node))."""
    for i in ids:
        case(name, "apex-" + i, i)


def prefixlist(name, apexid, prefixes, slug):
    case(name, "apex-%s-incl-%s" % (apexid, slug), apexid, prefixes)


# --------------------------------------------------------------------------- 1-4 structure
doc("minimal.xml", "<r/>")
whole("minimal.xml")

# Positional newlines around prolog and epilog nodes (section 2.8).
doc("prolog-epilog.xml", '<!--a--><?p1 x?><r Id="r"><x Id="x"/></r><?p2 y?><!--b-->')
whole("prolog-epilog.xml")
apex("prolog-epilog.xml", "r", "x")

# The Santuario epilog-drop quirk: a childless document element loses the whole epilog.
doc("prolog-epilog-empty-root.xml", '<!--a--><r Id="r"/><?pi d?><!--b-->')
whole("prolog-epilog-empty-root.xml")
apex("prolog-epilog-empty-root.xml", "r")

# PI with no data at all versus PI with data; "&#13;" inside PI data is six literal
# characters, because PIs recognize no references.
doc("prolog-pi-nodata.xml", '<?a?><?b c?><r Id="r"><?p a&#13;b?><?q?>t</r><?z?>')
whole("prolog-pi-nodata.xml")
apex("prolog-pi-nodata.xml", "r")

# --------------------------------------------------------------------------- 5-8 attributes
doc("attr-order.xml",
    '<r xmlns:b="urn:b" xmlns:a="urn:a" b:z="1" a:z="2" z="3" a="4" xmlns="urn:d" Id="r"/>')
whole("attr-order.xml")
apex("attr-order.xml", "r")

# 100 attributes in anti-sorted order across five prefixes plus no-namespace names.
attrs = []
for i in range(5):
    attrs.append('xmlns:p%d="urn:p%d"' % (4 - i, 4 - i))
for i in range(99, -1, -1):
    if i % 3 == 0:
        attrs.append('p%d:a%02d="%d"' % (i % 5, i, i))
    elif i % 3 == 1:
        attrs.append('a%02d="%d"' % (i, i))
    else:
        attrs.append('p%d:z%02d="%d"' % (4 - i % 5, i, i))
doc("attr-many.xml", '<r Id="r" ' + " ".join(attrs) + "/>")
whole("attr-many.xml")
apex("attr-many.xml", "r")

# XML 1.0 3.3.3 attribute-value normalization: literal whitespace collapses to #x20,
# referenced whitespace does not (section 1.4).
doc("attr-ws-normalize.xml",
    b'<r Id="r" a="x\ty\nz\r\nw\rv" b="&#9;&#10;&#13;" c=" lead trail "'
    b' d="a  b" e="&#x9;&#xA;&#xD;"/>')
whole("attr-ws-normalize.xml")
apex("attr-ws-normalize.xml", "r")

doc("attr-quotes.xml",
    b'<r Id="r" a=\'he said "hi"\' b="it&apos;s" c="&quot;&apos;" d=\'&#34;&#39;\''
    b" e='&lt;&amp;&gt;'/>")
whole("attr-quotes.xml")
apex("attr-quotes.xml", "r")

# --------------------------------------------------------------------------- 9-11 character data
doc("text-escapes.xml",
    b'<r Id="r">a&amp;b&lt;c&gt;d"e]]&gt;f&#13;g&#x0D;h\ti\nj&#9;k</r>')
whole("text-escapes.xml")
apex("text-escapes.xml", "r")

doc("cdata.xml",
    b'<r Id="r"><![CDATA[a < b & c]]>mid<![CDATA[]]]]><![CDATA[>]]>tail<![CDATA[]]>'
    b'<e Id="e">x<![CDATA[y]]>z</e></r>')
whole("cdata.xml")
apex("cdata.xml", "r", "e")

doc("comments-pi-nested.xml",
    '<r Id="r"><!--top--><a Id="a"><?pi d?><!--in a--><b Id="b">t<!--after t--></b></a>'
    "<!--tail--></r>")
whole("comments-pi-nested.xml")
apex("comments-pi-nested.xml", "r", "a", "b")

# --------------------------------------------------------------------------- 12-18 namespaces
doc("ns-superfluous.xml",
    '<r xmlns:p="urn:1" Id="r"><p:c xmlns:p="urn:1" Id="c">'
    '<p:d xmlns:p="urn:1" Id="d"/></p:c></r>')
whole("ns-superfluous.xml")
apex("ns-superfluous.xml", "r", "c", "d")

doc("ns-rebind.xml",
    '<r xmlns:p="urn:1" Id="r"><p:c xmlns:p="urn:2" Id="c">'
    '<p:d xmlns:p="urn:1" Id="d"><p:e Id="e"/></p:d></p:c></r>')
whole("ns-rebind.xml")
apex("ns-rebind.xml", "r", "c", "d", "e")

doc("ns-unused.xml",
    '<r xmlns:unused="urn:u" xmlns:u2="urn:u2" Id="r"><c Id="c"><u2:d Id="d"/></c></r>')
whole("ns-unused.xml")
apex("ns-unused.xml", "r", "c", "d")

doc("ns-default-undeclare.xml",
    '<r xmlns="urn:a" Id="r"><c xmlns="" Id="c"><d Id="d"/>'
    '<e xmlns="urn:a" Id="e"><f Id="f"/></e></c></r>')
whole("ns-default-undeclare.xml")
apex("ns-default-undeclare.xml", "r", "c", "d", "e", "f")

# The apex is the undeclaring element: no xmlns="" is emitted (section 2.4).
doc("ns-default-apex.xml", '<r xmlns="urn:d" Id="r"><c xmlns="" Id="c"><t Id="t"/></c></r>')
whole("ns-default-apex.xml")
apex("ns-default-apex.xml", "r", "c", "t")

doc("ns-xml-prefix.xml",
    '<r xmlns:xml="http://www.w3.org/XML/1998/namespace" Id="r">'
    '<c Id="c" xml:lang="en"><t Id="t"/></c></r>')
whole("ns-xml-prefix.xml")
apex("ns-xml-prefix.xml", "r", "c", "t")

# Relative namespace URIs: an error for 1.0/1.1/exclusive when the declaration is actually
# rendered, accepted by physical, silently survived by exclusive when unused.
doc("ns-relative.xml",
    '<r xmlns:rel="relative" xmlns:un="alsorelative" Id="r">'
    '<rel:c Id="c"><t Id="t"/></rel:c></r>')
whole("ns-relative.xml")
apex("ns-relative.xml", "r", "c", "t")

# --------------------------------------------------------------------------- 19-20 exclusive
doc("excl-visibly-utilized.xml",
    '<r xmlns:p="urn:1" xmlns:q="urn:2" xmlns="urn:d" xmlns:un="urn:u" Id="r">'
    '<p:c q:a="1" b="2" Id="c"><t Id="t"/></p:c></r>')
whole("excl-visibly-utilized.xml")
apex("excl-visibly-utilized.xml", "r", "c", "t")

doc("excl-prefixlist.xml",
    '<r xmlns="urn:d" xmlns:p="urn:1" xmlns:q="urn:2" Id="r">'
    '<p:c Id="c"><t2 Id="t2"/></p:c></r>')
whole("excl-prefixlist.xml")
apex("excl-prefixlist.xml", "r", "c", "t2")
prefixlist("excl-prefixlist.xml", "c", "q", "q")
prefixlist("excl-prefixlist.xml", "c", "#default q", "default-q")
prefixlist("excl-prefixlist.xml", "c", "#default", "default")
prefixlist("excl-prefixlist.xml", "c", "zz", "zz")
prefixlist("excl-prefixlist.xml", "c", "p q  zz", "p-q-zz")
prefixlist("excl-prefixlist.xml", "t2", "#default q", "default-q")

# --------------------------------------------------------------------------- 21-25 xml:*
doc("xmlattrs-lang-space.xml",
    '<r xml:lang="en" xml:space="preserve" Id="r"><m xml:lang="fr" Id="m"><t a="1" Id="t"/></m>'
    '<n Id="n"><u xml:space="default" Id="u"><v Id="v"/></u></n></r>')
whole("xmlattrs-lang-space.xml")
apex("xmlattrs-lang-space.xml", "r", "m", "t", "n", "u", "v")

doc("xmlbase-absolute.xml",
    '<r xml:base="http://x/a/" Id="r"><m xml:base="b/" Id="m">'
    '<n xml:base="c/" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-absolute.xml")
apex("xmlbase-absolute.xml", "r", "m", "n", "t")

doc("xmlbase-relative.xml",
    '<r xml:base="b/" Id="r"><m xml:base="c/" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-relative.xml")
apex("xmlbase-relative.xml", "r", "m", "t")

doc("xmlbase-dotsegments.xml",
    '<r xml:base="http://x/a/b/../c/./d//e" Id="r"><m xml:base="../f/.." Id="m">'
    '<n xml:base="./g/../h" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-dotsegments.xml")
apex("xmlbase-dotsegments.xml", "r", "m", "n", "t")

doc("xml-id.xml",
    '<r xml:id="top" xml:lang="en" Id="r"><m xml:id="mid" Id="m"><t Id="t"/>'
    '<u xml:id="me" Id="u"/></m></r>')
whole("xml-id.xml")
apex("xml-id.xml", "r", "m", "t", "u")

# --------------------------------------------------------------------------- 26-29 encoding
doc("unicode.xml",
    '<r Id="r" a="\U0001D11E \u0301\u00E8 \u0085 \u2028 caf\u00E9">'
    "\U0001D11E combining: e\u0301 NEL:\u0085 LSEP:\u2028 caf\u00E9"
    "</r>")
whole("unicode.xml")
apex("unicode.xml", "r")

# Attribute sort keys that straddle U+FFFF: UTF-16 order (Java) puts the astral URI first,
# UTF-8 order (Go) puts U+FDF0 first (section 2.7).
doc("sort-astral.xml",
    '<r Id="r" xmlns:a="urn:x\uFDF0" xmlns:b="urn:x\U00010000" a:z="1" b:z="2"/>')
whole("sort-astral.xml")
apex("sort-astral.xml", "r")

doc("bom-decl.xml",
    b"\xef\xbb\xbf" +
    b'<?xml version="1.0" encoding="UTF-8" standalone="no"?>\n<r Id="r">t</r>\n')
whole("bom-decl.xml")
apex("bom-decl.xml", "r")

doc("encoding-latin1.xml",
    ('<?xml version="1.0" encoding="ISO-8859-1"?><r Id="r" a="caf\u00E9">'
     "na\u00EFve r\u00E9sum\u00E9</r>").encode("iso-8859-1"))
whole("encoding-latin1.xml")
apex("encoding-latin1.xml", "r")

# --------------------------------------------------------------------------- 30-32 misc
doc("entities.xml",
    '<r Id="r" a="&amp;&lt;&gt;&quot;&apos;">&amp;&lt;&gt;&quot;&apos;'
    '<c Id="c" b="&#38;&#60;&#62;&#34;&#39;">&#38;&#60;&#62;&#34;&#39;</c></r>')
whole("entities.xml")
apex("entities.xml", "r", "c")

doc("empty-elements.xml",
    b'<r Id="r"><a Id="a"></a><b Id="b"/><c Id="c">   </c><d Id="d">\n</d>'
    b'<e Id="e"><f Id="f"/></e></r>')
whole("empty-elements.xml")
apex("empty-elements.xml", "r", "a", "b", "c", "d", "e")

deep = []
for i in range(200):
    ns = ""
    if i % 50 == 0:
        ns = ' xmlns:p%d="urn:d%d"' % (i, i)
    if i % 37 == 0:
        ns += ' xml:lang="l%d"' % i
    deep.append('<l%03d Id="l%03d"%s>' % (i, i, ns))
deep.append("text")
for i in range(199, -1, -1):
    deep.append("</l%03d>" % i)
doc("deep-nesting.xml", "".join(deep))
whole("deep-nesting.xml")
apex("deep-nesting.xml", "l000", "l100", "l199")

# --------------------------------------------------------------------------- 33-34 real world
upstream = [
    ("xades-signature.xml",
     "dss-signature-remote/src/test/resources/xades-signed.xml"),
    ("trusted-list.xml",
     "dss-signature-remote/src/test/resources/trusted-list-v6.xml"),
]
for name, rel in upstream:
    path = os.path.join(DSS_UPSTREAM, rel)
    with open(path, "rb") as fh:
        docs[name] = fh.read()
whole("xades-signature.xml")
apex("xades-signature.xml", "id-910825ec07149183c174c83fce12ac93", "r-id-1", "o-id-1")
prefixlist("xades-signature.xml", "r-id-1", "ds", "ds")
whole("trusted-list.xml")
apex("trusted-list.xml", "TL12345")

# =========================================================================== ADVERSARIAL
# Added by the phase-4a audit. Each block targets a rule the original corpus left unpinned;
# several of them found real defects, noted inline.

# --- nested default-namespace redeclarations -------------------------------------------
# Undeclare, redeclare to a third URI, undeclare again. Every apex exercises a different
# initial state of the nsStack "xmlns" entry, which is what lastrendered exists to track.
doc("ns-default-nested.xml",
    '<r xmlns="urn:a" Id="r"><b xmlns="urn:b" Id="b"><c xmlns="" Id="c">'
    '<d xmlns="urn:a" Id="d"><e xmlns="" Id="e"><f Id="f"/></e></d></c></b></r>')
whole("ns-default-nested.xml")
apex("ns-default-nested.xml", "r", "b", "c", "d", "e", "f")

# A default namespace declared and undeclared on the SAME element as a prefixed one, with
# the element itself unprefixed: exclusive c14n must decide "visibly utilized" per element.
doc("ns-default-mixed.xml",
    '<r xmlns="urn:d" xmlns:p="urn:1" Id="r"><p:a Id="a"><b xmlns="" Id="b">'
    '<p:c Id="c"><d Id="d"/></p:c></b></p:a></r>')
whole("ns-default-mixed.xml")
apex("ns-default-mixed.xml", "r", "a", "b", "c", "d")

# --- same local name, different namespace URIs ------------------------------------------
# AttrCompare sorts namespaced attributes by URI then local name, never by prefix. Here the
# prefix order and the URI order are deliberately opposed: prefix "z" is bound to "urn:a" and
# prefix "a" to "urn:z", so a correct sort emits z:k before a:k.
doc("attr-same-local-diff-uri.xml",
    '<r xmlns:z="urn:a" xmlns:a="urn:z" xmlns:m="urn:m" Id="r"'
    ' a:k="from-urn-z" z:k="from-urn-a" m:k="from-urn-m" k="no-namespace"/>')
whole("attr-same-local-diff-uri.xml")
apex("attr-same-local-diff-uri.xml", "r")

# The same collision inside a namespaced element, plus one attribute per URI sharing every
# local name, so a comparator that ties on either field alone produces a different order.
doc("attr-collide-deep.xml",
    '<r xmlns:p="urn:1" xmlns:q="urn:2" Id="r">'
    '<p:e Id="e" q:a="qa" p:a="pa" q:b="qb" p:b="pb" a="plain-a" b="plain-b">'
    '<q:f Id="f" p:z="pz" q:z="qz" xmlns:r2="urn:0" r2:z="rz"/></p:e></r>')
whole("attr-collide-deep.xml")
apex("attr-collide-deep.xml", "r", "e", "f")

# --- xml:* on ancestors outside the canonicalized subtree -------------------------------
# The apex inherits lang/space/base from ancestors that are not themselves output. 1.0 takes
# the first occurrence of each QName; 1.1 joins the base chain and drops xml:id.
doc("xmlattrs-outside-subtree.xml",
    '<r xml:lang="en" xml:space="preserve" xml:base="http://x/a/" xml:id="rid" Id="r">'
    '<m xml:lang="fr" xml:base="b/" Id="m"><n xml:base="c/" xml:space="default" Id="n">'
    '<t Id="t"><u Id="u"/></t></n></m></r>')
whole("xmlattrs-outside-subtree.xml")
apex("xmlattrs-outside-subtree.xml", "r", "m", "n", "t", "u")

# The apex carries its OWN xml:base, which seeds the 1.1 accumulator ahead of every ancestor.
doc("xmlbase-apex-own.xml",
    '<r xml:base="http://x/a/" Id="r"><m xml:base="b/" Id="m">'
    '<t xml:base="own/" Id="t"><u Id="u"/></t></m></r>')
whole("xmlbase-apex-own.xml")
apex("xmlbase-apex-own.xml", "r", "m", "t", "u")

# --- xml:base shapes that stress joinURI and java.net.URI recomposition ------------------
# Percent escapes: java.net.URI's five-argument constructor quotes '%' itself, so an already
# escaped path comes back double-escaped. Reproduce, do not "fix".
doc("xmlbase-escaped.xml",
    '<r xml:base="http://x/a%20b/" Id="r"><m xml:base="c%2Fd/" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-escaped.xml")
apex("xmlbase-escaped.xml", "r", "m", "t")

# Characters java.net.URI refuses: a literal space and a double quote make the constructor
# throw URISyntaxException, which Santuario swallows, keeping the previously accumulated
# value. The fold is skipped, not the whole canonicalization.
doc("xmlbase-badchars.xml",
    '<r xml:base="http://x/a/" Id="r"><m xml:base="has space/" Id="m">'
    '<n xml:base="ok/" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-badchars.xml")
apex("xmlbase-badchars.xml", "r", "m", "n", "t")

# Query strings and authorities, including the "authority present, path empty" branch.
doc("xmlbase-query.xml",
    '<r xml:base="http://x/a/?q=1" Id="r"><m xml:base="b/?r=2" Id="m">'
    '<n xml:base="//other" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-query.xml")
apex("xmlbase-query.xml", "r", "m", "n", "t")

# An opaque URI (scheme whose scheme-specific part is not hierarchical): java.net.URI reports
# a null path, and Santuario dereferences it.
doc("xmlbase-opaque.xml",
    '<r xml:base="urn:example:a" Id="r"><m xml:base="b/" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-opaque.xml")
apex("xmlbase-opaque.xml", "r", "m", "t")

# Non-ASCII in xml:base: allowed unescaped by java.net.URI's scanEscape, so it survives.
doc("xmlbase-unicode.xml",
    '<r xml:base="http://x/caf\u00e9/" Id="r"><m xml:base="r\u00e9sum\u00e9/" Id="m">'
    '<t Id="t"/></m></r>')
whole("xmlbase-unicode.xml")
apex("xmlbase-unicode.xml", "r", "m", "t")

# Santuario looks for the accumulator seed by LOCAL NAME alone, so an ordinary unprefixed
# base="..." attribute on the apex is picked up as though it were xml:base.
doc("xmlbase-plain-base.xml",
    '<r xml:base="http://x/a/" Id="r"><m xml:base="b/" Id="m">'
    '<t base="plain/" Id="t"><u Id="u"/></t></m></r>')
whole("xmlbase-plain-base.xml")
apex("xmlbase-plain-base.xml", "r", "m", "t", "u")

# removeDotSegments' own branches: leading "../", a bare "..", "/../" against a short output
# buffer, and the trailing-".." fix-up that appends a slash.
doc("xmlbase-dotsegments2.xml",
    '<r xml:base="http://x/a/b/c/" Id="r"><m xml:base="../../d" Id="m">'
    '<n xml:base=".." Id="n"><o xml:base="/../../e/" Id="o">'
    '<p2 xml:base="./././f" Id="p2"><t Id="t"/></p2></o></n></m></r>')
whole("xmlbase-dotsegments2.xml")
apex("xmlbase-dotsegments2.xml", "r", "m", "n", "o", "p2", "t")

# An all-relative chain with no scheme anywhere, which keeps joinURI in its relative branch
# for every fold.
doc("xmlbase-norelscheme.xml",
    '<r xml:base="a/b" Id="r"><m xml:base="c/d" Id="m"><n xml:base="../e" Id="n">'
    '<t Id="t"/></n></m></r>')
whole("xmlbase-norelscheme.xml")
apex("xmlbase-norelscheme.xml", "r", "m", "n", "t")

# --- xml:base branches only reachable with unusual component shapes -------------------
# An empty xml:base, and one that is a bare query: both give joinURI a relative reference
# whose path is empty, which is the branch that inherits the base's path and query.
doc("xmlbase-empty.xml",
    '<r xml:base="http://x/a/b?q=1" Id="r"><m xml:base="" Id="m">'
    '<n xml:base="?r=2" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-empty.xml")
apex("xmlbase-empty.xml", "r", "m", "n", "t")

# A base with no slash at all, so bpath.lastIndexOf('/') is -1 and the relative path replaces
# it outright rather than being appended to a directory.
doc("xmlbase-noslash.xml",
    '<r xml:base="abc" Id="r"><m xml:base="def" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-noslash.xml")
apex("xmlbase-noslash.xml", "r", "m", "t")

# Authority present, path empty: "http://host" + "x" takes the "/" + rpath branch.
doc("xmlbase-authority-empty.xml",
    '<r xml:base="http://host" Id="r"><m xml:base="x" Id="m">'
    '<n xml:base="y/z" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-authority-empty.xml")
apex("xmlbase-authority-empty.xml", "r", "m", "n", "t")

# A fragment on an xml:base. joinURI recomposes with a null fragment, so it is dropped.
doc("xmlbase-fragment.xml",
    '<r xml:base="http://x/a/#top" Id="r"><m xml:base="b/#mid" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-fragment.xml")
apex("xmlbase-fragment.xml", "r", "m", "t")

# removeDotSegments' remaining arms: a leading "./", a lone ".", a lone "/.", and a "../"
# run long enough to exhaust the output buffer and start accumulating "../" segments.
doc("xmlbase-dotsegments3.xml",
    '<r xml:base="http://x/a/b/c/d/" Id="r"><m xml:base="./e" Id="m">'
    '<n xml:base="." Id="n"><o xml:base="/." Id="o">'
    '<p3 xml:base="../../../../../f" Id="p3"><q3 xml:base="g/./h/../i" Id="q3">'
    '<t Id="t"/></q3></p3></o></n></m></r>')
whole("xmlbase-dotsegments3.xml")
apex("xmlbase-dotsegments3.xml", "r", "m", "n", "o", "p3", "q3", "t")

# An all-relative chain of "../" with no scheme, which is where the output buffer legitimately
# ends up holding "../../" and the endsWith("..") fix-ups fire.
doc("xmlbase-dotdot-relative.xml",
    '<r xml:base="../a" Id="r"><m xml:base="../../b" Id="m"><n xml:base=".." Id="n">'
    '<t Id="t"/></n></m></r>')
whole("xmlbase-dotdot-relative.xml")
apex("xmlbase-dotdot-relative.xml", "r", "m", "n", "t")

# Malformed percent escape and a bare scheme: both make java.net.URI throw
# URISyntaxException, which is swallowed, so the accumulated value survives untouched.
doc("xmlbase-badescape.xml",
    '<r xml:base="http://x/a/" Id="r"><m xml:base="b%zz/" Id="m">'
    '<n xml:base="foo:" Id="n"><t Id="t"/></n></m></r>')
whole("xmlbase-badescape.xml")
apex("xmlbase-badescape.xml", "r", "m", "n", "t")

# A no-break space (U+00A0) is a Unicode space character, so java.net.URI refuses it
# unescaped on the way in and percent-encodes it as UTF-8 on the way out.
doc("xmlbase-nbsp.xml",
    '<r xml:base="http://x/a\u00a0b/" Id="r"><m xml:base="c/" Id="m"><t Id="t"/></m></r>')
whole("xmlbase-nbsp.xml")
apex("xmlbase-nbsp.xml", "r", "m", "t")

# --- PrefixList corners -----------------------------------------------------------------
# "#empty" is the sentinel for the empty-string PrefixList overload (see C14nOracle.runCase);
# "xml" and a repeated token check that prefixStr2Set's TreeSet dedups and that the xml
# prefix is never rendered even when named explicitly.
doc("excl-prefixlist2.xml",
    '<r xmlns="urn:d" xmlns:p="urn:1" xmlns:q="urn:2" xmlns:un="urn:u" xml:lang="en" Id="r">'
    '<p:c q:a="1" b="2" Id="c"><t Id="t"><q:d Id="d"/></t></p:c></r>')
whole("excl-prefixlist2.xml")
apex("excl-prefixlist2.xml", "r", "c", "t", "d")
prefixlist("excl-prefixlist2.xml", "c", "#empty", "empty")
prefixlist("excl-prefixlist2.xml", "t", "#empty", "empty")
prefixlist("excl-prefixlist2.xml", "c", "#default", "default")
prefixlist("excl-prefixlist2.xml", "t", "#default", "default")
prefixlist("excl-prefixlist2.xml", "c", "xml", "xml")
prefixlist("excl-prefixlist2.xml", "c", "p p q q", "dup")
prefixlist("excl-prefixlist2.xml", "c", "un", "un")
prefixlist("excl-prefixlist2.xml", "d", "#default p un", "default-p-un")

# --- whitespace-only text nodes between elements -----------------------------------------
# Nothing is stripped: there is no DTD, so there is no ignorable whitespace. Mixed CR, CRLF
# and lone LF also pin XML 2.11 line-ending normalization through c14n.
doc("ws-between-elements.xml",
    b'<r Id="r">  <a Id="a"/>\n\t<b Id="b"> </b>\r\n  <c Id="c">\r</c>\n'
    b'<d Id="d"><e Id="e"/>   <f Id="f"/></d>\n</r>')
whole("ws-between-elements.xml")
apex("ws-between-elements.xml", "r", "a", "b", "c", "d", "e", "f")

# Whitespace and comments in the prolog and epilog around a root that HAS children, so the
# epilog-drop quirk does not mask the positional-newline rule.
doc("ws-prolog-epilog.xml",
    b'\n\n  <!--c1-->\n<?pi1 a?>\n<r Id="r"> <x Id="x"/> </r>\n<?pi2 b?>\n<!--c2-->\n\n')
whole("ws-prolog-epilog.xml")
apex("ws-prolog-epilog.xml", "r", "x")

# --- BOM handling ------------------------------------------------------------------------
# A UTF-8 BOM with no XML declaration at all.
doc("bom-nodecl.xml", b"\xef\xbb\xbf" + b'<r Id="r">t</r>')
whole("bom-nodecl.xml")
apex("bom-nodecl.xml", "r")

# UTF-16LE and UTF-16BE, detected by BOM. The declaration names the encoding too.
doc("bom-utf16le.xml",
    b"\xff\xfe" +
    ('<?xml version="1.0" encoding="UTF-16"?><r Id="r" a="caf\u00e9">t\u00e9xt</r>')
    .encode("utf-16-le"))
whole("bom-utf16le.xml")
apex("bom-utf16le.xml", "r")

doc("bom-utf16be.xml",
    b"\xfe\xff" +
    ('<?xml version="1.0" encoding="UTF-16"?><r Id="r" a="caf\u00e9">t\u00e9xt</r>')
    .encode("utf-16-be"))
whole("bom-utf16be.xml")
apex("bom-utf16be.xml", "r")

# --- NEL, LSEP and friends ---------------------------------------------------------------
# XML 1.0 does NOT normalize NEL (U+0085) or LSEP (U+2028); XML 1.1 would, which is why
# xmldom rejects 1.1 outright. Placed in text, attribute values, a comment and PI data.
doc("nel-lsep.xml",
    '<r Id="r" a="nel\u0085lsep\u2028end" b="&#x85;&#x2028;">'
    'text\u0085nel\u2028lsep<!--comment\u0085\u2028here--><?pi data\u0085\u2028?>'
    '<c Id="c">\u0085\u2028</c></r>')
whole("nel-lsep.xml")
apex("nel-lsep.xml", "r", "c")

# --- codepoints above 127 in attribute values and names ----------------------------------
# Accented and astral characters in attribute VALUES, in attribute LOCAL NAMES, in element
# names and in namespace URIs, so the UTF-16 comparator is exercised on real sort keys.
doc("attr-high-codepoints.xml",
    '<caf\u00e9 Id="r" xmlns:\u00e9="urn:\u00e9" xmlns:\u4e2d="urn:\u4e2d"'
    ' \u00e9:v="e-acute" \u4e2d:v="zhong" na\u00efve="value-\u00e9\u4e2d\U0001d11e"'
    ' plain="\U0001d11e\u0301"><r\u00e9sum\u00e9 Id="c">\u4e2d\u6587</r\u00e9sum\u00e9>'
    '</caf\u00e9>')
whole("attr-high-codepoints.xml")
apex("attr-high-codepoints.xml", "r", "c")

# Namespace URIs straddling U+FFFF on BOTH sides of the comparison, plus a third in the BMP,
# so a byte-wise sort and a UTF-16 sort disagree on more than one pair.
doc("sort-astral2.xml",
    '<r Id="r" xmlns:a="urn:\uff00" xmlns:b="urn:\U00010000" xmlns:c="urn:\ue000"'
    ' xmlns:d="urn:\U0010ffff" a:z="a" b:z="b" c:z="c" d:z="d"/>')
whole("sort-astral2.xml")
apex("sort-astral2.xml", "r")

# --------------------------------------------------------------------------- negatives
negative = {
    "doctype.xml": b'<!DOCTYPE r><r/>',
    "entity-subset.xml": b'<!DOCTYPE r [<!ENTITY e "x">]><r>&e;</r>',
    "empty-prefix-binding.xml": b'<r xmlns:p=""><p:c/></r>',
    "xmlns-xml-bogus.xml": b'<r xmlns:xml="urn:bogus"/>',
    "undeclared-prefix.xml": b'<r><p:c/></r>',
    "dup-attr.xml": b'<r a="1" a="2"/>',
    "dup-attr-ns.xml": b'<r xmlns:p="urn:1" xmlns:q="urn:1" p:a="1" q:a="2"/>',
    "tag-mismatch.xml": b'<r></s>',
    "two-roots.xml": b'<r/><s/>',
    "trailing-text.xml": b'<r/>tail',
    "nul-charref.xml": b'<r>&#0;</r>',
    "unknown-entity.xml": b'<r>&foo;</r>',
    "cdata-end-in-text.xml": b'<r>a]]>b</r>',
    "lt-in-attr.xml": b'<r a="<"/>',
    "xml11.xml": b'<?xml version="1.1"?><r/>',
    "xmlns-prefix-declared.xml": b'<r xmlns:xmlns="http://www.w3.org/2000/xmlns/"/>',

    # --- added by the phase-4a audit ---
    # DOCTYPE in every spelling a parser might treat differently. disallow-doctype-decl is
    # what stops external entities, parameter entities and billion-laughs, so each shape has
    # to be refused, not just the canonical one.
    "doctype-external.xml": b'<!DOCTYPE r SYSTEM "r.dtd"><r/>',
    "doctype-public.xml": b'<!DOCTYPE r PUBLIC "-//X//DTD//EN" "http://x/r.dtd"><r/>',
    "doctype-newlines.xml": b'<!DOCTYPE\n   r\n   [\n<!ELEMENT r EMPTY>\n]\n><r/>',
    "doctype-after-comment.xml": b'<!--lead--><!DOCTYPE r><r/>',
    "doctype-after-decl.xml": b'<?xml version="1.0"?><!DOCTYPE r><r/>',
    "doctype-param-entity.xml": b'<!DOCTYPE r [<!ENTITY % p "<!ENTITY e \'x\'>">%p;]><r>&e;</r>',
    "doctype-billion-laughs.xml":
        b'<!DOCTYPE r [<!ENTITY a "xx"><!ENTITY b "&a;&a;"><!ENTITY c "&b;&b;">]><r>&c;</r>',
    "doctype-notation.xml": b'<!DOCTYPE r [<!NOTATION n SYSTEM "n">]><r/>',

    # Character references Go's encoding/xml is known to mishandle or to accept where XML 1.0
    # forbids them. Surrogates in particular are silently folded to U+FFFD by encoding/xml.
    "surrogate-charref.xml": b'<r>&#xD800;</r>',
    "surrogate-charref-low.xml": b'<r>&#xDFFF;</r>',
    "surrogate-charref-attr.xml": b'<r a="&#xD800;"/>',
    "surrogate-charref-decimal.xml": b'<r>&#55296;</r>',
    "uppercase-x-charref.xml": b'<r>&#X41;</r>',
    "noncharacter-fffe.xml": b'<r>&#xFFFE;</r>',
    "noncharacter-ffff.xml": b'<r>&#xFFFF;</r>',
    "charref-out-of-range.xml": b'<r>&#x110000;</r>',
    "charref-control.xml": b'<r>&#x1;</r>',
    "charref-empty.xml": b'<r>&#;</r>',
    "charref-unterminated.xml": b'<r>&#x41</r>',

    # Namespace rules Go accepts silently.
    "xmlns-uri-to-other-prefix.xml": b'<r xmlns:p="http://www.w3.org/2000/xmlns/"/>',
    "xml-prefix-other-uri.xml": b'<r xmlns:foo="http://www.w3.org/XML/1998/namespace"/>',
    "undeclared-attr-prefix.xml": b'<r p:a="1"/>',
    "default-ns-empty-redeclare.xml": b'<r xmlns:p=""><c/></r>',

    # Structure.
    "unclosed-at-eof.xml": b'<r><c></r>',
    "nested-mismatch.xml": b'<r><a></b></a></r>',
    "text-before-root.xml": b'lead<r/>',
    "no-root.xml": b'<!--only a comment-->',
    "empty-document.xml": b'',
    "cdata-unterminated.xml": b'<r><![CDATA[a</r>',
    "comment-double-hyphen.xml": b'<r><!--a--b--></r>',
    "pi-xml-target-late.xml": b'<r><?xml version="1.0"?></r>',
}

# --------------------------------------------------------------------------- write
os.makedirs(CORPUS, exist_ok=True)
os.makedirs(os.path.join(CORPUS, "negative"), exist_ok=True)
for name in sorted(docs):
    with open(os.path.join(CORPUS, name), "wb") as fh:
        fh.write(docs[name])
for name in sorted(negative):
    with open(os.path.join(CORPUS, "negative", name), "wb") as fh:
        fh.write(negative[name])

os.makedirs(os.path.dirname(os.path.abspath(CASES)), exist_ok=True)
with open(CASES, "w", encoding="utf-8") as fh:
    fh.write("# document\tscope-slug\tapex\tprefixlist   (written by gen/corpus.py)\n")
    for name, slug, ap, pl in cases:
        fh.write("%s\t%s\t%s\t%s\n" % (name, slug, ap, pl))

print("corpus: %d documents, %d negative, %d cases"
      % (len(docs), len(negative), len(cases)))
