# Diagnostic-data marshal-parity oracles

Ground truth for `dss/diagnostic/jaxb`: every file under `oracle/` was produced
by the JAXB reference implementation of upstream DSS 6.5.RC1, and the KAT in
`../xml_kat_test.go` requires each of them to unmarshal into the Go model and
marshal back byte for byte.

- `DiagnosticDataOracle.java` — runs upstream validation over a signed fixture
  and marshals the resulting `XmlDiagnosticData` with `DiagnosticDataFacade`.
  50 dumps, taken from the `cades/`, `xades/`, `pades/`, `jades/` and `asic/`
  testdata trees, at the pinned validation time 2025-06-15T12:00:00Z.
  Prefix a fixture with `none:` to skip token extraction (keeps a large dump
  small); append `::<file>` to supply detached content.
- `DiagnosticDataFillOracle.java` — fills *every* property of every reachable
  generated class by reflection and marshals that. Real dumps only exercise a
  fraction of the schema; these exercise the model itself.
  - `oracle/model-full.xml` — every property set.
  - `oracle/model-empty.xml` — every string and byte[] set to its empty lexical
    form, which pins the `<X></X>` versus `<X/>` spelling across the model.
  - `oracle/model-specials.xml` — every string carries the characters whose XML
    spelling differs between the RI and `encoding/xml`.
  - `oracle/model-deep.xml` — `model-full.xml` at recursion depth 7 instead of
    6. The extra level is what reaches `DrivingPrivilegeCodeClaim`'s own
    `Code`/`Sign`/`Value` children, the last two element bindings of the model
    that no depth-6 or real dump exercises; with it the corpus covers all 476
    element and all 74 attribute names the model binds.
- `EnumLexicalOracle.java` — the lexical form each generated adapter prints for
  every constant of the enumeration it binds, as `adapters.tsv`
  (`adapter <TAB> enumeration <TAB> constant <TAB> lexical`, `\N` = no lexical
  form). Consumed by `../jaxb_adapters_test.go`.
- `xsd/DiagnosticData.xsd`, `xsd/bindings.xml` — byte-identical copies of
  `dss-diagnostic-jaxb/src/main/resources/xsd`, swept against the model by
  `../jaxb_schema_test.go`.
- `xslt/svg/diagnostic-data.xslt` — copied for the record only; the SVG
  rendering of `DiagnosticDataFacade` is deferred, no XSLT engine is ported.

## Regenerating

Build upstream DSS with maven, then compose a classpath of the installed
`eu.europa.ec.joinup.sd-dss` jars, their dependencies, and
`dss-pades-pdfbox/target/classes` plus the pdfbox jars (the PAdES validator is
loaded through `ServiceLoader`, so the implementation module has to be on the
classpath):

    javac -cp "$CP" -d /tmp/oracle *.java
    java -cp "$CP:/tmp/oracle" DiagnosticDataOracle oracle <fixture>...
    java -cp "$CP:/tmp/oracle" DiagnosticDataFillOracle oracle/model-full.xml 6 full
    java -cp "$CP:/tmp/oracle" DiagnosticDataFillOracle oracle/model-empty.xml 6 empty
    java -cp "$CP:/tmp/oracle" DiagnosticDataFillOracle oracle/model-specials.xml 6 specials
    java -cp "$CP:/tmp/oracle" DiagnosticDataFillOracle oracle/model-deep.xml 7 full
    java -cp "$CP:/tmp/oracle" EnumLexicalOracle adapters.tsv

`DSS_DIAGNOSTIC_ORACLE_DIR` replays the KAT over a larger local corpus without
committing it.
