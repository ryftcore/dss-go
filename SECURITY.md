# Security Policy

esig is a cryptographic library used to create and validate digital
signatures. We take security reports seriously and ask that you report
suspected vulnerabilities privately so that a fix can be prepared before
any public disclosure.

## Reporting a vulnerability

**Do not open a public GitHub issue for a suspected vulnerability.**

Please report security issues privately using GitHub Security Advisories:

1. Go to https://github.com/ryftcore/dss-go/security/advisories/new
2. Fill in as much detail as you can:
   - Affected package(s)/version(s) or commit
   - A description of the issue and its potential impact (e.g. signature
     forgery, validation bypass, denial of service, information
     disclosure)
   - Steps to reproduce, or a minimal proof-of-concept (a crafted
     document/signature/certificate is often the clearest reproduction
     for this kind of library)
   - Any suggested fix or mitigation, if you have one

If you are unable to use GitHub Security Advisories for any reason,
contact the repository owner through their GitHub profile
(https://github.com/utain) to arrange another private channel before
sharing details.

Please do not test for vulnerabilities against infrastructure or
services you do not own or have permission to test (for example, do not
probe production Trusted List / LOTL endpoints beyond normal read
access).

## What to expect

- We will acknowledge receipt of a report as soon as we are able to.
- We will investigate, ask follow-up questions if needed, and work with
  you on a fix and a coordinated disclosure timeline.
- Credit is given to reporters in the advisory unless you ask to remain
  anonymous.

This is a community open-source project without a dedicated security
team or a formal SLA; response times are best-effort.

## Supported versions

esig has not yet made a tagged `v1.0.0` (or later) release. Until a
first tagged release is published, security fixes are made against the
`main` branch only.

| Version         | Supported          |
| ---------------- | ------------------ |
| `main` (pre-1.0)  | :white_check_mark: |

This table will be updated with concrete version ranges once tagged
releases begin.

## Scope notes

esig is a Go port of [esig/dss](https://github.com/esig/dss). A
vulnerability that originates in the upstream Java implementation's
*design* (e.g. a weakness in a signature format or validation algorithm
as specified) is best reported upstream as well, since it likely affects
both projects; a vulnerability specific to this Go implementation
(a porting bug, a Go-specific memory/logic issue, a divergence from the
upstream algorithm) should be reported here.
