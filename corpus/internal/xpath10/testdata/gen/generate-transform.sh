#!/bin/sh
# Regenerates testdata/transform-kat.txt: the JDK XPath's answers for every XML-DSig transform
# expression in internal/xmldsig/testdata/corpus, which is itself a verbatim copy of upstream
# dss-xades test fixtures. The corpus is shared rather than duplicated - these are the same
# documents whose reference digests internal/xmldsig pins.
#
#   sh gen/generate-transform.sh
#
# Requires OpenJDK 21 only: the engine under test is the JDK's own XPath, which is what Apache
# Santuario 3.0 evaluates transform expressions with.
set -eu
cd "$(dirname "$0")/.."   # testdata/

CORPUS=${CORPUS:-../../xmldsig/testdata/corpus}
[ -d "$CORPUS" ] || { echo "missing corpus: $CORPUS" >&2; exit 1; }

OUT=$(mktemp -d)
javac -encoding UTF-8 -d "$OUT" gen/TransformXPathOracle.java
java -cp "$OUT" TransformXPathOracle "$CORPUS" transform-kat.txt
rm -rf "$OUT"
wc -l transform-kat.txt
