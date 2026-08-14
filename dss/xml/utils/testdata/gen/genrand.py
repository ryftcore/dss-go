import random

random.seed(20260813)

NS = ["urn:1", "urn:2", "urn:3", "http://www.w3.org/2000/09/xmldsig#", "-"]
PFX = ["", "p", "q", "ds", "n1", "xmlfoo", "x"]
CHARS = list("ab<>&\"']>-?") + ["\t", "\r", "\n", "é", "ü", "中", " ", " "]


def esc_dsl(s):
    out = []
    for ch in s:
        cp = ord(ch)
        if ch == "\\":
            out.append("\\\\")
        elif ch == " ":
            out.append("\\s")
        elif ch == "\n":
            out.append("\\n")
        elif ch == "\r":
            out.append("\\r")
        elif ch == "\t":
            out.append("\\t")
        elif ch == ";":
            out.append("\\u003B")
        elif cp < 0x20 or cp > 0x7E:
            out.append("\\u%04X" % cp)
        else:
            out.append(ch)
    return "".join(out) or "~"


def rand_text(n=6):
    return "".join(random.choice(CHARS) for _ in range(random.randint(1, n)))


PATHS = []


def gen(depth=0, path="0"):
    PATHS.append(path)
    kid = 0
    ops = []
    ns = random.choice(NS)
    pfx = random.choice(PFX)
    local = random.choice(["a", "b", "Sig", "Data", "é"])
    qname = (pfx + ":" + local) if pfx else local
    if ns == "-":
        qname = local  # a prefix without a namespace is not constructible
    ops.append("e %s %s" % (ns, qname))
    for _ in range(random.randint(0, 3)):
        ans = random.choice(NS)
        apfx = random.choice(PFX)
        aloc = random.choice(["x", "y", "Id", "lang"])
        aq = (apfx + ":" + aloc) if (apfx and ans != "-") else aloc
        ops.append("a %s %s %s" % (ans, aq, esc_dsl(rand_text())))
    for _ in range(random.randint(0, 3)):
        k = random.random()
        if k < 0.30 and depth < 3:
            ops.extend(gen(depth + 1, "%s.%d" % (path, kid)))
            ops.append("/")
        elif k < 0.55:
            ops.append("t " + esc_dsl(rand_text()))
        elif k < 0.70:
            ops.append("d " + esc_dsl(rand_text()))
        elif k < 0.85:
            ops.append("c " + esc_dsl(rand_text()))
        else:
            ops.append("p pi%d %s" % (random.randint(0, 9), esc_dsl(rand_text())))
        kid += 1
    return ops


rows = []
for i in range(140):
    del PATHS[:]
    ops = gen()
    sel = random.choice(["."] + PATHS)
    rows.append("case rand-%03d\nbuild %s\nsel %s\n" % (i, "; ".join(ops), sel))

header = """
# ------------------------------------------------------ seeded randomized differential
# 140 pseudo-random trees (Python's Mersenne twister, seed 20260813) built from the same
# alphabet of namespaces, prefixes, node kinds and awkward characters as the hand-written
# cases, each serialized from a randomly chosen node. Selectors that do not resolve are
# recorded as errors on both sides. Regenerate with gen/genrand.py rather than editing.
"""
with open("xml/utils/testdata/serialize/corpus.txt", "a") as f:
    f.write(header + "\n".join(rows))
print("added", len(rows))
