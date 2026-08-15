# Chain / ChainItem oracle corpus

`oracle/chain_semantics.jsonl` holds one JSON object per synthetic chain
scenario: the `XmlConstraintsConclusion` upstream's
`eu.europa.esig.dss.validation.process.Chain` / `ChainItem` /
`UninterruptedChainItem` produce when driven over that scenario. The rows are
produced by `gen/ChainSemanticsOracle.java`, which subclasses the upstream
classes directly - it is the same check-execution algorithm the Go port has to
reproduce, not a re-implementation of it.

`../chain_test.go` carries the identical scenario table and compares its own
`Chain`/`ChainItem` output against these rows.

## What the scenarios pin

* every `Level` branch: `FAIL`, `WARN`, `INFORM`, `IGNORE`, and an undefined
  constraint (Java's null `LevelRule`, which skips the check and calls the next
  item);
* the first-failure short circuit at `Level.FAIL`, and its opposite, an
  `UninterruptedChainItem` continuing the chain past a failure (including the
  conclusion reuse that appends the second error to the first conclusion);
* the conclusion built on failure: indication, sub-indication, and the choice
  between `getPreviousErrors()` and the check's own error message;
* the custom success conclusion (`getSuccessIndication()` non-null), which
  records the conclusion instead of calling the next item;
* the warning / info collection into the chain conclusion, in constraint order;
* the population of every `XmlConstraint` member - `Name`, `Status`, `Error`,
  `Warning`, `Info`, `AdditionalInfo`, `Id`, `BlockType` - and the rule that an
  `IGNORED` constraint carries no additional info;
* the chain title, including a chain with no `MessageTag` title.

## Regenerating

From a `dss-upstream` checkout with the modules built:

    cd dss-upstream
    mvn -o dependency:build-classpath -pl dss-validation -Dmdep.outputFile=/tmp/cp.txt
    CP="dss-validation/target/classes:$(cat /tmp/cp.txt)"
    javac -cp "$CP" -d /tmp/oracle gen/ChainSemanticsOracle.java
    java  -cp "$CP:/tmp/oracle" ChainSemanticsOracle > oracle/chain_semantics.jsonl

## Known mapping

Java leaves the `Title` attribute null when a chain defines no title
`MessageTag`; the generated Go model carries `Title` as a plain string, so the
rows' `null` and the port's `""` are the same document. The test normalises that
one field and nothing else.
