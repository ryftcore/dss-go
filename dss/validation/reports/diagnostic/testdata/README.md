# Certificate diagnostic-data builder oracle

`oracle/<name>.xml` is the marshalled `XmlDiagnosticData` upstream's
`CertificateDiagnosticDataBuilder` produces for the single certificate
`dss/spi/testdata/certificate_extensions/<name>.der`, at the pinned validation date
2024-01-01T00:00:00Z, with `DigestAlgorithm.SHA256` and
`TokenExtractionStrategy.NONE`. 49 documents, produced by `gen/CertificateDiagnosticDataOracle.java`.

`../certificate_diagnostic_data_builder_kat_test.go` builds the same document in Go
and byte-compares the marshalled XML. That is what pins the element order of every
certificate, certificate-extension and QcStatements builder in this package - the
part of `DiagnosticDataBuilder` reachable without a live signature.

One document per certificate keeps the comparison independent of the `java.util.HashSet`
iteration order over the used-certificate set, which Go cannot reproduce.

## Marshalling overload

The driver uses `DiagnosticDataFacade#marshall(Object, boolean)` - the `StringWriter`
overload - the same one `dss/diagnostic/jaxb/testdata/gen/DiagnosticDataOracle.java`
uses for the marshal-parity corpus, and the one `dss/diagnostic/jaxb.Marshal`
reproduces.

This matters: the `OutputStream` overload takes a different JAXB output class
(`org.glassfish.jaxb.runtime...output.IndentingUTF8XmlOutput`) whose `printIndent()`
computes `i = depth % 8` and then shifts *that* value right by three instead of
shifting `depth`, so its "whole blocks of eight" loop never runs and the indentation
wraps back to column 0 at depth 8, 16, ... The two overloads therefore do NOT produce
the same bytes for a document nested deeper than eight elements (a certificate with
PSD2 `RoleOfPSP` entries is exactly such a document). Any future end-to-end byte-parity
oracle has to pin which overload it compares against.

## Known deviations

Two certificates are not compared, and two more are compared as a multiset of lines;
`katKnownDeviation` / `katOrderDeviation` in the test carry the reason for each, and
the test fails if a listed file stops deviating.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle gen/CertificateDiagnosticDataOracle.java
    java  -cp "$CP:/tmp/oracle" CertificateDiagnosticDataOracle \
        <dss-repo>/spi/testdata/certificate_extensions oracle
