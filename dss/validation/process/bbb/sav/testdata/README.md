# sav (signature acceptance validation) oracle corpora

Both files are pure Java dumps, produced by the drivers in `../../testdata/gen`.

| file | rows | driver | input |
| --- | --- | --- | --- |
| `oracle/sav_blocks.jsonl` | 188 | `FcSavOracle.java` | `SignatureAcceptanceValidation` / `TimestampAcceptanceValidation` / `RevocationAcceptanceValidation` over every signature, time-stamp and revocation of the marshal-parity diagnostic-data corpus, default ETSI policy, validation time 2024-01-01T00:00:00Z |
| `oracle/sav_direct.jsonl` | 2235 | `FcSavDirectOracle.java` | all 34 sav check classes driven alone at `Level.FAIL` - the ones the three chains never run under the default policy (constraint undefined, or only wired for a signature form the corpus does not carry) and the ones they do run but never fail - over the same dumps plus synthetic signatures for the branches no dump reaches |

The `XmlAOV` the chains consume is built by the `aov` package; the block corpus
feeds them a PASSED one (no errors, warnings or infos), which is what
`AlgorithmObsolescenceValidationCheck` reads to pick its own `Level`. The
detailed-report `XmlTimestamp`s `TLevelTimeStampCheck` / `LTALevelTimeStampCheck`
read are likewise synthesised, at both a PASSED and a FAILED basic-timestamp
conclusion.

The real inputs are the marshal-parity corpus in `dss/diagnostic/jaxb/testdata/oracle`;
the Go tests (`../sav_blocks_oracle_test.go`, `../sav_direct_oracle_test.go`) read the
very same files.

Every check class covered here has both an `OK` and a `NOT OK` row; the direct test
fails if any check class ever loses one of the two. `TimestampMessageImprintWithIdCheck`
lives in `validation/process/vpfltvd` (over `vpftspwatsp`'s base class) but is wired by
`SignatureAcceptanceValidation`, so its rows sit here.

The four `model-*.xml` schema-coverage fixtures are excluded from this corpus, for the
reason given in `../../fc/testdata/README.md`.

## Regenerating

See `../../fc/testdata/README.md`; the two drivers write both packages' corpora.

## Known mapping

A single-item chain defines no title `MessageTag`, so Java leaves the `Title`
attribute null; the generated Go model carries `Title` as a plain string. The
direct test normalises that one field and nothing else.
