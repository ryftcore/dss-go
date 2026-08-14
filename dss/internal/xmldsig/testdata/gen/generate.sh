#!/bin/sh
# Regenerates internal/xmldsig/testdata:
#
#   corpus/**            copied verbatim from dss-xades/src/test/resources per gen/cases.txt
#   golden/*, manifest.txt   gen/XmlDsigOracle.java - Apache Santuario 3.0.6's answers
#
# The goldens are Java's answers. Nothing here ever runs the Go implementation, and a golden
# that changes in a pull request is a red flag, exactly as for internal/xmlc14n and
# internal/cmscore.
#
#   sh gen/generate.sh
#
# Requires OpenJDK 21, the upstream DSS checkout under $DSS and Apache Santuario 3.0.6 plus
# slf4j-api in the local maven repository. Override with DSS=... M2=... sh gen/generate.sh
set -eu

cd "$(dirname "$0")/.."   # testdata/

DSS=${DSS:-/home/user/dss-upstream}
M2=${M2:-$HOME/.m2/repository}
XMLSEC=${XMLSEC:-3.0.6}
SLF4J=${SLF4J:-2.0.18}

RES="$DSS/dss-xades/src/test/resources"

CP="$M2/org/apache/santuario/xmlsec/$XMLSEC/xmlsec-$XMLSEC.jar"
CP="$CP:$M2/org/slf4j/slf4j-api/$SLF4J/slf4j-api-$SLF4J.jar"
for jar in $(echo "$CP" | tr ':' ' '); do
    [ -f "$jar" ] || { echo "missing jar: $jar" >&2; exit 1; }
done

echo "--- corpus"
rm -rf corpus
mkdir -p corpus
sed 's/#.*//' gen/cases.txt | while IFS="$(printf '\t')" read -r fixture detached; do
    [ -n "${fixture:-}" ] || continue
    fixture=$(printf '%s' "$fixture" | tr -d ' ')
    [ -n "$fixture" ] || continue
    # mutants/* are not upstream files; gen/mutate.py derives them below.
    case "$fixture" in mutants/*) continue ;; esac
    mkdir -p "corpus/$(dirname "$fixture")"
    cp "$RES/$fixture" "corpus/$fixture"
    [ -n "${detached:-}" ] || continue
    for d in $(printf '%s' "$detached" | tr ',' ' '); do
        mkdir -p "corpus/$(dirname "$d")"
        cp "$RES/$d" "corpus/$d"
    done
done

echo "--- mutants"
python3 gen/mutate.py corpus gen/mutants.txt "$RES"

echo "--- oracle"
OUT=$(mktemp -d)
javac -encoding UTF-8 -cp "$CP" -d "$OUT" gen/XmlDsigOracle.java
rm -rf golden
mkdir -p golden
java -cp "$CP:$OUT" XmlDsigOracle corpus golden manifest.txt gen/cases.txt
rm -rf "$OUT"

echo "--- done"
wc -l manifest.txt
ls golden | wc -l
