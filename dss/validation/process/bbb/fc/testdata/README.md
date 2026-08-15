# fc (format checking) oracle corpora

Both files are pure Java dumps, produced by the drivers in `../../testdata/gen`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/fc_blocks.jsonl` | 121 | `FcSavOracle.java` | `SignatureFormatChecking` / `TimestampFormatChecking` over every signature and time-stamp of the marshal-parity diagnostic-data corpus, default ETSI policy |
| `oracle/fc_direct.jsonl` | 771 | `FcSavDirectOracle.java` | all 28 fc check classes driven alone at `Level.FAIL` - the ones the two chains never run under the default policy, and the ones the chains do run but the corpus never fails - over the same dumps plus explicitly built synthetic inputs |

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`;
the Go tests (`../fc_blocks_oracle_test.go`, `../fc_direct_oracle_test.go`) read the
very same files, so neither side gets a private fixture.

Four dumps (`model-deep.xml`, `model-empty.xml`, `model-full.xml`,
`model-specials.xml`) are schema-coverage fixtures, not validation output: every
property of the generated model is filled by reflection from a counter. The blocks
corpus carries no rows for the tokens upstream's own wrappers throw on, and the
**direct corpus skips those four files entirely** - their IDREF graph is dangling
(a `<SigningCertificate>` pointing at no `<Certificate>`: upstream's
`getSigningCertificate()` answers null where the Go wrapper resolves an object) and
their attribute values carry raw control characters (a literal TAB, which XML
attribute-value normalisation turns into a space for the JCA parser but which
`encoding/xml` keeps). Rows taken over them would compare two different inputs. The
branches they used to reach are covered by explicit synthetic inputs instead.

Every check class covered here has both an `OK` and a `NOT OK` row; the direct test
fails if any check class ever loses one of the two.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle ../../testdata/gen/FcSavOracle.java ../../testdata/gen/FcSavDirectOracle.java
    java  -cp "$CP:/tmp/oracle" FcSavOracle       <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>
    java  -cp "$CP:/tmp/oracle" FcSavDirectOracle <dss-repo>/diagnostic/jaxb/testdata/oracle <dss-repo>

## Known mapping

A single-item chain defines no title `MessageTag`, so Java leaves the `Title`
attribute null; the generated Go model carries `Title` as a plain string. The
direct test normalises that one field and nothing else.
