# ETSI Validation Report marshal-parity oracles

Ground truth for `dss/validationreport/jaxb`: every file under `oracle/` was
produced by the JAXB reference implementation of upstream DSS 6.5.RC1, and
the KAT in `../xml_kat_test.go` requires each of them to unmarshal into the
Go model and marshal back byte for byte.

- `ValidationReportOracle.java` — runs upstream validation (default
  validation policy) over a signed fixture and marshals the resulting
  `ValidationReportType` via `Reports.getXmlValidationReport()`
  (`ValidationReportFacade` under the hood), at the pinned validation time
  2025-06-15T12:00:00Z. 18 dumps, taken from the `xades/`, `cades/`,
  `pades/`, `jades/` and `asic/` testdata trees, chosen to spread across
  every signature format this repo carries fixtures for plus revocation
  data, timestamps, evidence records and counter-signatures (the schema
  areas that exercise `POEType`/`RevocationStatusInformationType`/
  `CryptoInformationType`/`SignatureReferenceType`'s XAdESSignaturePtr
  branch).
- `xsd/1910202xmlSchema.xsd`, `xsd/bindings.xml`, `xsd/catalog.cat` —
  byte-identical copies of `specs-validation-report/src/main/resources/xsd`,
  swept against the model by `../xml_schema_test.go`.

## Regenerating

Build upstream DSS with maven so every module is installed to the local
`.m2` repository, then compose a classpath of the installed
`eu.europa.ec.joinup.sd-dss` jars, their dependencies, and
`dss-pades-pdfbox/target/classes` plus the pdfbox/fontbox jars (the PAdES
validator, and every format's `DocumentValidatorFactory`, are loaded through
`ServiceLoader`, so each implementation module's own jar - not just its
dependencies - has to be on the classpath: `dependency:build-classpath` only
lists a module's dependencies, never the module's own artifact):

    javac -cp "$CP" -d /tmp/oracle ValidationReportOracle.java
    java -cp "$CP:/tmp/oracle" ValidationReportOracle oracle <fixture>...

`DSS_VALIDATIONREPORT_ORACLE_DIR` replays the KAT over a larger local corpus
without committing it.
