# isc / vci / cv building-block oracle corpora

`gen/BbbBlocksOracle.java` drives the upstream building blocks

* `eu.europa.esig.dss.validation.process.bbb.isc.IdentificationOfTheSigningCertificate`
* `eu.europa.esig.dss.validation.process.bbb.vci.ValidationContextInitialization`
* `eu.europa.esig.dss.validation.process.bbb.cv.CryptographicVerification`

and writes their results as JSON rows into the three packages' own
`testdata/oracle` directories:

| file | rows | input |
| --- | --- | --- |
| `../isc/testdata/oracle/isc_blocks.jsonl` | 184 | every signature, time-stamp and revocation of the 50 usable diagnostic-data dumps |
| `../isc/testdata/oracle/isc_direct.jsonl` | 8 | the four isc checks alone, at `Level.FAIL`, over synthetic signing-certificate references |
| `../vci/testdata/oracle/vci_blocks.jsonl` | 53 | every signature of the same dumps |
| `../vci/testdata/oracle/vci_direct.jsonl` | 10 | the five vci checks alone, at `Level.FAIL`, over synthetic signature policies |
| `../cv/testdata/oracle/cv_blocks.jsonl` | 184 | every signature, time-stamp and revocation of the same dumps |
| `../cv/testdata/oracle/cv_direct.jsonl` | 13 | the cv checks the chain never wires, plus the branches the corpus does not reach |

The real inputs are the marshal-parity corpus in
`dss/diagnostic/jaxb/testdata/oracle` - the Go tests read the very same files, so
neither side gets a private fixture. The validation policy is the default ETSI
one (`EtsiValidationPolicyFactory#loadDefaultValidationPolicy`).

Four dumps (`model-deep.xml`, `model-empty.xml`, `model-full.xml`,
`model-specials.xml`) are schema-coverage fixtures whose wrapper graph is
incomplete; upstream's own `CertificateWrapper` throws on them, so they carry no
rows. The Go tests iterate the dumps named by the rows, so they skip exactly the
same four.

## Why some checks are driven directly

`CryptographicVerification` never wires `AtLeastOneReferenceDataObjectFoundCheck`,
`ReferenceDataGroupCheck` or `SignatureIntactWithIdCheck` - they belong to the
evidence-record and archival blocks. Those, plus the branches the
50 dumps happen not to reach (a failing manifest-entry group, a name mismatch, a
`CERTIFICATE`-context signature-intact message, a covering HashTree-renewal
time-stamp, and every vci failure - the default policy runs the vci checks below
`FAIL`), are driven through a one-item chain at `Level.FAIL` so that each check
class has both a happy and a failure row.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle gen/BbbBlocksOracle.java
    java  -cp "$CP:/tmp/oracle" BbbBlocksOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>
