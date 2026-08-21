#!/usr/bin/env python3
"""Writes testdata/adversarial/*.xml - hand-built documents whose ds:XPath and
xpf:XPath elements carry expressions the upstream dss-xades corpus does not contain.

The corpus-derived known answers in transform-kat.txt only ever exercise the handful of
expressions real signers emit. These fixtures exist to ask the questions an attacker or a
careless signer would: here(), which Xalan provided and Santuario 3.0 no longer does;
ancestor-or-self on a tree deep enough that a recursive evaluator would be tempted to
memoise it; a ds:Signature nested inside another signature's ds:Object, where
"not(ancestor-or-self::ds:Signature)" has two possible answers; id() over duplicated Ids;
prefixes the ds:XPath element does not declare; and the position()/last()/count() family
whose context depends on how the evaluator builds a node-set.

gen/TransformXPathOracle.java then records the JDK XPath's answer for every one of them,
so nothing here decides what is correct - it only decides what gets asked.

    python3 gen/make-adversarial.py <adversarial dir>
"""
import os
import sys

DS = "http://www.w3.org/2000/09/xmldsig#"
F2 = "http://www.w3.org/2002/06/xmldsig-filter2"


def sig(inner_refs, body="", sig_id="sig1", extra_root_ns=""):
    """A ds:Signature carrying the given ds:Transform bodies."""
    return (
        '<ds:Signature xmlns:ds="%s" Id="%s"%s><ds:SignedInfo>'
        '<ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/>'
        '<ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/>'
        '<ds:Reference URI=""><ds:Transforms>%s</ds:Transforms>'
        '<ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/>'
        '<ds:DigestValue>AAAA</ds:DigestValue></ds:Reference></ds:SignedInfo>'
        '<ds:SignatureValue>AAAA</ds:SignatureValue>%s</ds:Signature>'
        % (DS, sig_id, extra_root_ns, inner_refs, body)
    )


def xpath_transform(expr, extra_ns=""):
    return ('<ds:Transform Algorithm="http://www.w3.org/TR/1999/REC-xpath-19991116">'
            '<ds:XPath%s>%s</ds:XPath></ds:Transform>' % (extra_ns, expr))


def filter2_transform(expr, filt="intersect", extra_ns=""):
    return ('<ds:Transform Algorithm="http://www.w3.org/2002/06/xmldsig-filter2">'
            '<xpf:XPath xmlns:xpf="%s" Filter="%s"%s>%s</xpf:XPath></ds:Transform>'
            % (F2, filt, extra_ns, expr))


DECL = '<?xml version="1.0" encoding="UTF-8"?>'


def main():
    out = sys.argv[1]
    os.makedirs(out, exist_ok=True)
    files = {}

    # --- here(), the Xalan extension Santuario 3.0 dropped ------------------------------
    files["here-bare.xml"] = DECL + "<doc><d Id=\"data\">x</d>" + sig(
        xpath_transform("here()")) + "</doc>"
    files["here-ancestor.xml"] = DECL + "<doc><d Id=\"data\">x</d>" + sig(
        xpath_transform("here()/ancestor::ds:Signature[1]")) + "</doc>"
    files["here-count.xml"] = DECL + "<doc><d Id=\"data\">x</d>" + sig(
        xpath_transform("count(here()) &gt; 0")) + "</doc>"
    files["here-filter2.xml"] = DECL + "<doc><d Id=\"data\">x</d>" + sig(
        filter2_transform("here()/ancestor::ds:Signature[1]//. ")) + "</doc>"

    # --- ancestor-or-self over a deep tree ----------------------------------------------
    depth = 300
    deep = "".join("<n%d>" % i for i in range(depth)) + "leaf" + \
           "".join("</n%d>" % i for i in reversed(range(depth)))
    files["deep-ancestor.xml"] = DECL + "<doc>" + deep + sig(
        xpath_transform("not(ancestor-or-self::ds:Signature)")) + "</doc>"
    files["deep-ancestor-count.xml"] = DECL + "<doc>" + deep + sig(
        xpath_transform("count(ancestor-or-self::*) &lt; 10")) + "</doc>"
    files["deep-ancestor-name.xml"] = DECL + "<doc>" + deep + sig(
        xpath_transform("ancestor-or-self::n7 and not(ancestor-or-self::ds:Signature)")) + "</doc>"

    # --- a ds:Signature nested inside another signature's ds:Object ----------------------
    inner = sig(xpath_transform("not(ancestor-or-self::ds:Signature)"), sig_id="inner")
    files["nested-signature.xml"] = DECL + "<doc><d>x</d>" + sig(
        xpath_transform("not(ancestor-or-self::ds:Signature)"),
        body='<ds:Object Id="o1">' + inner + "</ds:Object>", sig_id="outer") + "</doc>"
    files["nested-signature-descendant.xml"] = DECL + "<doc><d>x</d>" + sig(
        xpath_transform("/descendant::ds:Signature"),
        body='<ds:Object Id="o1">' + inner + "</ds:Object>", sig_id="outer") + "</doc>"
    files["nested-signature-self.xml"] = DECL + "<doc><d>x</d>" + sig(
        xpath_transform("count(ancestor-or-self::ds:Signature) = 1"),
        body='<ds:Object Id="o1">' + inner + "</ds:Object>", sig_id="outer") + "</doc>"

    # --- id() over duplicated and missing Ids --------------------------------------------
    files["id-duplicate.xml"] = DECL + '<doc><a Id="dup">first</a><b Id="dup">second</b>' \
        '<c Id="dup">third</c>' + sig(xpath_transform("id('dup')")) + "</doc>"
    files["id-duplicate-nodeset.xml"] = DECL + '<doc><a Id="dup">first</a><b Id="dup">second</b>' \
        + sig(xpath_transform("id('dup')/node()")) + "</doc>"
    files["id-missing.xml"] = DECL + '<doc><a Id="present"/>' + sig(
        xpath_transform("id('absent')")) + "</doc>"
    files["id-multiarg.xml"] = DECL + '<doc><a Id="p"/><b Id="q"/>' + sig(
        xpath_transform("id('p q')")) + "</doc>"

    # --- prefixes the ds:XPath element does not declare ----------------------------------
    files["prefix-undeclared.xml"] = DECL + '<doc xmlns:zz="urn:Z"><zz:a/>' + sig(
        xpath_transform("not(ancestor-or-self::zz:a)")) + "</doc>"
    files["prefix-on-xpath-element.xml"] = DECL + '<doc xmlns:zz="urn:Z"><zz:a/>' + sig(
        xpath_transform("not(ancestor-or-self::q:a)", extra_ns=' xmlns:q="urn:Z"')) + "</doc>"
    files["prefix-shadowed.xml"] = DECL + '<doc xmlns:q="urn:OUTER"><q:a/>' + sig(
        xpath_transform("ancestor-or-self::q:a", extra_ns=' xmlns:q="urn:INNER"')) + "</doc>"

    # --- context-dependent functions ------------------------------------------------------
    files["position-last.xml"] = DECL + "<doc><a/><a/><a/><a/>" + sig(
        xpath_transform("position() = last()")) + "</doc>"
    files["self-position.xml"] = DECL + "<doc><a/><a/><a/>" + sig(
        xpath_transform("count(preceding-sibling::*) &gt; 1")) + "</doc>"
    files["namespace-axis.xml"] = DECL + '<doc xmlns:p="urn:P" xmlns="urn:D"><a/>' + sig(
        xpath_transform("count(namespace::*) &gt; 2")) + "</doc>"
    files["name-functions.xml"] = DECL + '<doc xmlns:p="urn:P"><p:a/><b/>' + sig(
        xpath_transform("starts-with(name(), 'p:') or local-name() = 'b'")) + "</doc>"
    files["text-predicates.xml"] = DECL + "<doc><a>alpha</a><a>beta</a><a></a>" + sig(
        xpath_transform("string-length(normalize-space(.)) &gt; 4")) + "</doc>"
    files["numeric-edges.xml"] = DECL + "<doc><a>1</a><a>-0</a><a>NaN</a><a>1e3</a>" + sig(
        xpath_transform("number(.) &gt; 0 or number(.) = 0")) + "</doc>"

    # --- Filter 2.0 shapes ----------------------------------------------------------------
    files["filter2-union.xml"] = DECL + '<doc><a Id="x"/><b/>' + sig(
        filter2_transform("//a | //b", "union")) + "</doc>"
    files["filter2-subtract.xml"] = DECL + "<doc><a><b/><c/></a>" + sig(
        filter2_transform("//c", "subtract")) + "</doc>"
    files["filter2-root.xml"] = DECL + "<doc><a/>" + sig(
        filter2_transform("/", "intersect")) + "</doc>"
    files["filter2-attribute.xml"] = DECL + '<doc><a k="v"/>' + sig(
        filter2_transform("//@k", "intersect")) + "</doc>"
    files["filter2-empty.xml"] = DECL + "<doc><a/>" + sig(
        filter2_transform("//nothing", "intersect")) + "</doc>"

    # --- malformed and hostile expressions -------------------------------------------------
    files["expr-unbalanced.xml"] = DECL + "<doc><a/>" + sig(
        xpath_transform("count(//a")) + "</doc>"
    files["expr-unknown-function.xml"] = DECL + "<doc><a/>" + sig(
        xpath_transform("no-such-function()")) + "</doc>"
    files["expr-empty.xml"] = DECL + "<doc><a/>" + sig(xpath_transform("")) + "</doc>"
    files["expr-whitespace.xml"] = DECL + "<doc><a/>" + sig(
        xpath_transform("\n\t  not(ancestor-or-self::ds:Signature)  \n ")) + "</doc>"
    files["expr-comment-syntax.xml"] = DECL + "<doc><a/>" + sig(
        xpath_transform("//a[1] (: not an XPath 1.0 comment :)")) + "</doc>"
    files["expr-variable.xml"] = DECL + "<doc><a/>" + sig(
        xpath_transform("$undefined")) + "</doc>"

    for name, body in files.items():
        with open(os.path.join(out, name), "wb") as fh:
            fh.write(body.encode("utf-8"))
    print("make-adversarial.py: wrote %d fixtures" % len(files))


if __name__ == "__main__":
    main()
