# A custom validation policy

Almost none of the "is this acceptable?" decisions in validation are hard-coded.
Which digest algorithms are allowed, how large an RSA key must be, how fresh
revocation data has to be, how long after signing a time-stamp may arrive,
whether a missing signing time is fatal or merely worth a mention — all of it
lives in a **validation policy**: an XML constraint document.

The library ships the ETSI policy and uses it whenever you do not say
otherwise. This page is about saying otherwise.

## Why you would

Legitimate reasons, roughly in descending order of how often they are the right
call:

- **A scheme mandates its own policy.** An e-invoicing regime, a national eID
  profile, a procurement platform: they publish constraints and you must apply
  them. This is the common case.
- **Your risk appetite differs on one axis.** You want SHA-1 to be a hard
  failure today rather than on the ETSI schedule, or you need a longer
  revocation-freshness window because your OCSP responder is slow.
- **You are validating historical material** under the rules that applied at the
  time.

And the reason that is usually *not* good: making a verdict you dislike go away.
Relaxing a constraint does not change the facts — it changes what the process
does about them. That trade should be written down somewhere with a name
attached to it.

## How constraints work

Each constraint carries a **level**, and the level is the whole design:

| Level | Effect when the constraint is not met |
|---|---|
| `FAIL` | The check fails, and it propagates into the indication. |
| `WARN` | Recorded as a warning. The verdict is unaffected. |
| `INFORM` | Recorded as information. Purely advisory. |
| `IGNORE` | Not evaluated at all. |

```xml
<ProspectiveCertificateChain Level="FAIL" />
```

Change that one attribute to `WARN` and an unanchored chain stops being a
blocking problem and becomes a note in the report. Nothing else changes — the
chain is just as unanchored as before.

## The smallest useful recipe

Start from the policy the library actually uses, change one thing, pass it in.
Do **not** hand-write a policy from scratch: it is a large document and the
default encodes a great deal of considered judgement.

```go
import (
	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/policy"
	"github.com/ryftcore/dss-go/dss/policy/jaxb"
)

// Load the library's own default policy.
facade := policy.NewValidationPolicyFacade()
constraints, err := facade.UnmarshalFile("constraint.xml")
if err != nil {
	log.Fatal(err)
}

// Change exactly one thing.
constraints.SignatureConstraints.BasicSignatureConstraints.
	ProspectiveCertificateChain.Level = jaxb.LevelValue(enumerations.LevelWarn)

// Marshal it back and validate against it.
relaxed, err := facade.Marshal(constraints)
if err != nil {
	log.Fatal(err)
}

reports, err := dss.Validate(doc, dss.ValidateOptions{
	Policy: dss.NewDocument("relaxed-policy.xml", relaxed),
})
```

The default document ships in the module at `policy/resources/constraint.xml`.
A complete runnable version of the above — including printing the before and
after verdicts side by side — is
[`examples/08-custom-policy`](https://github.com/ryftcore/dss-go/tree/main/dss/examples/08-custom-policy).

From the CLI, a policy is just a file:

```sh
esig validate contract.pdf -policy my-policy.xml -trust ca.cer
```

## What is in there

A tour of the sections you are most likely to touch, by the question each
answers:

| Section | Governs |
|---|---|
| `ContainerConstraints` | ASiC: acceptable container types, manifest presence, whether every file must be signed. |
| `SignatureConstraints` | The signature itself: acceptable formats and levels, the signing certificate, the chain, signed attributes, signing time. |
| `CounterSignatureConstraints` | The same, for counter-signatures. |
| `Timestamp` constraints | Time-stamp acceptance, and the permitted delay between signing and time-stamping. |
| `Revocation` constraints | Freshness, and what to do when revocation data is unavailable. |
| `Cryptographic` constraints | Acceptable digest and encryption algorithms, minimum key sizes, and expiry dates per algorithm. |
| `eIDAS` constraints | Trusted-list and qualification determination rules. |

The cryptographic section deserves a note, because it works differently from
the rest: algorithms have **expiry dates**, so a policy encodes not "SHA-1 is
bad" but "SHA-1 stopped being acceptable on this date". That is what makes it
possible to judge a 2009 signature by 2009's standards while refusing to accept
a new one today — and it is why the archival validation process can slide the
validation time backwards and still reach a defensible answer.

## The cryptographic suite, separately

Algorithm and key-size constraints can also come from a separate **cryptographic
suite** document (ETSI TS 119 312, in XML or JSON), which is the format
regulators actually publish those tables in.

```go
reports, err := dss.Validate(doc, dss.ValidateOptions{
	Policy:             myPolicy,
	CryptographicSuite: mySuite, // requires Policy to be set too
})
```

`CryptographicSuite` requires `Policy`, mirroring the underlying two-document
entry point.

## Advice

**Version your policy alongside your code.** Which constraints applied is part
of what a validation outcome *means*. Storing the outcome without the policy
that produced it makes it unreproducible.

**Prefer `WARN` over `IGNORE`.** A warning still shows up in the report, so the
next person can see the decision you made. `IGNORE` erases the evidence that
there was ever a question.

**Change one constraint at a time, and check the reports.** The
[DetailedReport](../concepts/reports.md) names the exact check that fired, so
you can confirm the change did what you meant and nothing else.

**Do not relax cryptographic constraints to accept old signatures.** That is
what the archival validation process and proofs of existence are for — they let
you accept a 2009 signature *because it is provably from 2009*, rather than by
pretending SHA-1 is still fine today.

## Next

- [How validation works](../concepts/validation.md) — where in the process each
  constraint is evaluated.
- [The four reports](../concepts/reports.md) — how to see which constraint fired.
