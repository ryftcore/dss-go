# Validation-root oracle corpus

| file | rows | generator |
| --- | --- | --- |
| `oracle/verifier_factories.tsv` | 8 | `oracle/gen/VerifierFactoryOracle.java` |

A pure Java dump of what `TrustAnchorVerifierFactory` and
`RevocationDataVerifierFactory` build from each of the four validation policies
DSS ships: the default `/policy/constraint.xml` plus the certificate, QWAC and
EAA constraint files, the last three of which are embedded byte-identically
under `../resources/`. Validation time `1700000000000` (2023-11-14T22:13:20Z).

Two lines per policy:

* `tav.*` - every value `TrustAnchorVerifier` carries. Asserted in full by
  `../trust_anchor_verifier_factory_oracle_test.go`.
* `rdv.*` - `RevocationDataVerifier`'s configuration, which the Java generator
  reads BY REFLECTION: upstream keeps those fields package-private and exposes
  no getters, and the port keeps them unexported for the same reason. The row is
  shipped as the recorded contract of what the factory must produce; the Go test
  can only assert that the factory builds. Closing that gap needs behavioural
  assertions driven through `RevocationDataVerifier`'s public surface.

## Regenerating

```
CP=<dss-validation runtime classpath>:<dss-validation target/classes>
javac -cp "$CP" -d /tmp/oracle oracle/gen/VerifierFactoryOracle.java
java  -cp "$CP:/tmp/oracle" VerifierFactoryOracle 2>/dev/null \
      | grep -E '^(default|/policy)' > oracle/verifier_factories.tsv
```
