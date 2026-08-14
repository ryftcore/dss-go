import base64
import random
import sys

random.seed(int(sys.argv[1]))

PFX = ["", "p", "q", "ds", "n1", "x"]
NSU = ["urn:1", "urn:2", "urn:3", "http://www.w3.org/2000/09/xmldsig#"]
TEXT = ["a", "b", "&amp;", "&lt;", "&gt;", "&quot;", "&apos;", "&#13;", "&#9;", "\n", " ", "\t",
        "é", "中", "&#128512;", "]]", "]]&gt;", "--", "?>", "&#127;", "&#160;"]
ENCS = ['', ' encoding="UTF-8"', ' encoding="utf-8"', ' encoding="ISO-8859-1"',
        ' encoding="US-ASCII"', ' encoding="UTF8"']
SA = ['', ' standalone="yes"', ' standalone="no"']

paths = []


def elem(depth, path, declared):
    paths.append(path)
    pfx = random.choice(PFX)
    local = random.choice(["a", "b", "Sig", "Data"])
    decls = []
    mine = dict(declared)
    for _ in range(random.randint(0, 2)):
        dp = random.choice(PFX)
        uri = random.choice(NSU)
        if dp == "":
            decls.append('xmlns="%s"' % random.choice(NSU + [""]))
            mine[""] = uri
        else:
            decls.append('xmlns:%s="%s"' % (dp, uri))
            mine[dp] = uri
    if pfx and pfx not in mine:
        decls.append('xmlns:%s="%s"' % (pfx, random.choice(NSU)))
        mine[pfx] = "x"
    qname = (pfx + ":" + local) if pfx else local
    attrs = list(decls)
    for _ in range(random.randint(0, 3)):
        ap = random.choice([p for p in mine if p] + [""])
        aloc = random.choice(["x", "y", "Id", "v"])
        aq = (ap + ":" + aloc) if ap else aloc
        if any(a.startswith(aq + "=") for a in attrs):
            continue
        attrs.append('%s="%s"' % (aq, "".join(random.choice(TEXT) for _ in range(random.randint(0, 3)))))
    body = []
    kid = 0
    for _ in range(random.randint(0, 3)):
        k = random.random()
        if k < 0.30 and depth < 3:
            body.append(elem(depth + 1, "%s.%d" % (path, kid), mine))
        elif k < 0.55:
            body.append("".join(random.choice(TEXT) for _ in range(random.randint(1, 4))))
        elif k < 0.70:
            body.append("<![CDATA[%s]]>" % "".join(
                random.choice(["a", "]", "]]", "<", "&", "é", "中"]) for _ in range(random.randint(0, 4))))
        elif k < 0.85:
            body.append("<!--%s-->" % random.choice(["", "c", " - ", "a b"]))
        else:
            body.append("<?pi%d %s?>" % (random.randint(0, 9), random.choice(["", "d", " sp  d "])))
        kid += 1
    open_tag = "<" + qname + ("".join(" " + a for a in attrs))
    if not body:
        return open_tag + "/>"
    return open_tag + ">" + "".join(body) + "</" + qname + ">"


rows = []
n = int(sys.argv[2])
for i in range(n):
    del paths[:]
    doc = elem(0, "0", {})
    decl = '<?xml version="1.0"%s%s?>' % (random.choice(ENCS), random.choice(SA))
    src = decl + doc
    enc = "utf-8"
    if 'ISO-8859-1' in decl:
        enc = "latin-1"
    if 'US-ASCII' in decl:
        enc = "ascii"
    try:
        raw = src.encode(enc)
    except UnicodeEncodeError:
        continue
    sel = random.choice(["."] + paths)
    rows.append("case pr-%03d\nparseb64 %s\nsel %s\n" % (i, base64.b64encode(raw).decode(), sel))

header = """
# ---------------------------------------------- seeded randomized differential, parsed
# %d pseudo-random XML documents (Python's Mersenne twister, seed %s), each parsed and
# then serialized from a randomly chosen node. Where the built trees of genrand.py cover
# DOM shapes no parser produces, these cover what a parser does produce: inherited
# namespace scopes, whitespace, entity references, mixed content and every spelling of
# the encoding declaration. Documents Xerces rejects are recorded as errors and must be
# rejected here too. Regenerate with gen/genparse.py rather than editing.
""" % (len(rows), sys.argv[1])
with open(sys.argv[3], "a") as f:
    f.write(header + "\n".join(rows))
print("wrote", len(rows))
