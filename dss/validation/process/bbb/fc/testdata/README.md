# fc (format checking) oracle corpora

Both files are pure Java dumps, produced by the drivers in `../../testdata/gen`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/fc_blocks.jsonl` | 121 | `FcSavOracle.java` | `SignatureFormatChecking` / `TimestampFormatChecking` over every signature and time-stamp of the marshal-parity diagnostic-data corpus, default ETSI policy |
| `oracle/fc_direct.jsonl` | 240 | `FcSavDirectOracle.java` | the ten fc checks the two chains never run under the default policy, each driven alone at `Level.FAIL` over the same corpus plus a few synthetic literals |

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`;
the Go tests (`../fc_blocks_oracle_test.go`, `../fc_direct_oracle_test.go`) read the
very same files, so neither side gets a private fixture.

Four dumps (`model-deep.xml`, `model-empty.xml`, `model-full.xml`,
`model-specials.xml`) are schema-coverage fixtures whose wrapper graph is
incomplete; upstream's own wrappers throw on them, so they carry no rows for the
tokens that throw. The Go tests iterate the rows, not the directory, so they skip
exactly the same tokens.

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
