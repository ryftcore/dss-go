#!/usr/bin/env python3
"""Builds the deliberately-corrupted fixtures under testdata/corpus/mutants/.

A known-answer corpus made only of well-formed upstream signatures proves that the port
agrees with Santuario when everything is intact. It does not prove the two disagree in the
same places. These mutants are the negative half: each one takes a real fixture and makes
exactly one textual change, and gen/XmlDsigOracle.java then records whatever verdict
Santuario reaches. The Go port has to reach the same one - including the cases where the
verdict does NOT flip, which are the interesting ones (canonicalization sorts attributes,
so reordering them must leave every digest alone).

Driven by gen/mutants.txt, whose columns are

    <mutant path>  <base fixture>  <kind>  <find>  <replace>

with \\n, \\t and \\\\ escapes in the last two fields. The substring in <find> must occur
EXACTLY ONCE in the base fixture; anything else aborts, so a mutant can never silently
become a no-op or hit the wrong element when the upstream fixture changes.

    python3 gen/mutate.py <corpus dir> <mutants.txt> <resources dir>
"""
import os
import sys


def unescape(s):
    out = []
    i = 0
    while i < len(s):
        c = s[i]
        if c == "\\" and i + 1 < len(s):
            n = s[i + 1]
            if n == "n":
                out.append("\n")
                i += 2
                continue
            if n == "t":
                out.append("\t")
                i += 2
                continue
            if n == "\\":
                out.append("\\")
                i += 2
                continue
        out.append(c)
        i += 1
    return "".join(out)


def main():
    corpus, spec, resources = sys.argv[1], sys.argv[2], sys.argv[3]
    n = 0
    for lineno, line in enumerate(open(spec, encoding="utf-8"), 1):
        line = line.rstrip("\n")
        if not line.strip() or line.startswith("#"):
            continue
        parts = line.split("\t")
        if len(parts) != 5:
            sys.exit("mutants.txt:%d: want 5 tab-separated fields, got %d" % (lineno, len(parts)))
        target, base, _kind, find, repl = parts
        find, repl = unescape(find), unescape(repl)
        if not target.startswith("mutants/"):
            sys.exit("mutants.txt:%d: mutant path must live under mutants/" % lineno)

        src = os.path.join(resources, base)
        with open(src, "rb") as fh:
            data = fh.read()
        fb, rb = find.encode("utf-8"), repl.encode("utf-8")
        count = data.count(fb)
        if count != 1:
            sys.exit("mutants.txt:%d: %r occurs %d times in %s, want exactly 1"
                     % (lineno, find, count, base))
        if fb == rb:
            sys.exit("mutants.txt:%d: the mutation is a no-op" % lineno)
        out = os.path.join(corpus, target)
        os.makedirs(os.path.dirname(out), exist_ok=True)
        with open(out, "wb") as fh:
            fh.write(data.replace(fb, rb, 1))
        n += 1
    print("mutate.py: wrote %d mutants" % n)


if __name__ == "__main__":
    main()
