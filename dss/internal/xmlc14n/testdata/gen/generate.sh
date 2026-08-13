#!/bin/sh
# Regenerates internal/xmlc14n/testdata: the corpus, the goldens and the manifest.
#
#   corpus/*.xml, corpus/negative/*.xml   gen/corpus.py        the adversarial documents
#   gen/cases.txt                         gen/corpus.py        (document, scope) pairs to run
#   golden/*, manifest.txt                gen/C14nOracle.java  Apache Santuario's answers
#
# The goldens are Java's answers, produced through the DSS call path; nothing here ever runs
# the Go implementation. Regenerating is a deliberate, reviewed act: a golden that changes in
# a pull request is a red flag, exactly as for the BouncyCastle oracle in internal/cmscore.
#
#   sh gen/generate.sh
#
# Requires OpenJDK 21, python3, the maven-built DSS 6.5.RC1 jars under $DSS and Apache
# Santuario 3.0.6 plus slf4j-api in the local maven repository. Override with
# DSS=... M2=... sh gen/generate.sh
set -eu

cd "$(dirname "$0")/.."   # testdata/

DSS=${DSS:-/home/user/dss-upstream}
M2=${M2:-$HOME/.m2/repository}
VER=6.5.RC1
XMLSEC=${XMLSEC:-3.0.6}
SLF4J=${SLF4J:-2.0.18}
LANG3=${LANG3:-3.20.0}
COLL4=${COLL4:-4.5.0}
COMMONSIO=${COMMONSIO:-2.22.0}
CODEC=${CODEC:-1.18.0}
JAXB=${JAXB:-3.0.1}
WOODSTOX=${WOODSTOX:-6.5.1}
STAX2=${STAX2:-4.2.1}

# The runtime dependencies of dss-xml-utils (mvn -o dependency:build-classpath), plus
# dss-utils-apache-commons and its own dependencies: dss-utils dispatches to an implementation
# module through a ServiceLoader and refuses to initialize without exactly one on the path.
CP="$DSS/dss-xml-utils/target/dss-xml-utils-$VER.jar"
for m in dss-xml-common dss-model dss-enumerations dss-alert dss-utils dss-utils-apache-commons; do
    CP="$CP:$DSS/$m/target/$m-$VER.jar"
done
CP="$CP:$M2/org/apache/santuario/xmlsec/$XMLSEC/xmlsec-$XMLSEC.jar"
CP="$CP:$M2/org/slf4j/slf4j-api/$SLF4J/slf4j-api-$SLF4J.jar"
CP="$CP:$M2/jakarta/xml/bind/jakarta.xml.bind-api/$JAXB/jakarta.xml.bind-api-$JAXB.jar"
CP="$CP:$M2/com/fasterxml/woodstox/woodstox-core/$WOODSTOX/woodstox-core-$WOODSTOX.jar"
CP="$CP:$M2/org/codehaus/woodstox/stax2-api/$STAX2/stax2-api-$STAX2.jar"
CP="$CP:$M2/commons-codec/commons-codec/$CODEC/commons-codec-$CODEC.jar"
CP="$CP:$M2/org/apache/commons/commons-lang3/$LANG3/commons-lang3-$LANG3.jar"
CP="$CP:$M2/org/apache/commons/commons-collections4/$COLL4/commons-collections4-$COLL4.jar"
CP="$CP:$M2/commons-io/commons-io/$COMMONSIO/commons-io-$COMMONSIO.jar"

for jar in $(echo "$CP" | tr ':' ' '); do
    [ -f "$jar" ] || { echo "missing jar: $jar" >&2; exit 1; }
done

echo "--- corpus"
rm -rf corpus
DSS="$DSS" python3 gen/corpus.py corpus gen/cases.txt

echo "--- oracle"
OUT=$(mktemp -d)
trap 'rm -rf "$OUT"' EXIT
javac -encoding UTF-8 -cp "$CP" -d "$OUT" gen/C14nOracle.java
java -cp "$CP:$OUT" C14nOracle corpus golden manifest.txt gen/cases.txt

echo "--- summary"
echo "corpus:  $(find corpus -name '*.xml' | wc -l) documents"
echo "goldens: $(ls golden | wc -l) files"
echo "errors:  $(awk -F'\t' '$8 == "ERROR"' manifest.txt | wc -l) KATs Java refuses"
