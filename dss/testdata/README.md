# Facade test fixtures

Small, self-contained fixtures for the root `dss` package's example and
regression tests. They are deliberately tiny so that the module zip stays small
for consumers; anything heavy belongs in the repository-level corpus, not here.

| File | What it is |
|---|---|
| `signer_rsa.p12` | PKCS#12 key store holding one self-signed RSA test key (`CN=Go Port Test RSA`), password `testpassword`. Same key store the CAdES/XAdES/PAdES/JAdES cross-validation generators use. Test key only - it anchors nothing and is not secret. |
| `signer_rsa.cer` | The DER encoding of the certificate inside `signer_rsa.p12`, for the certificate-loading examples. |
| `tsa_ec.p12` | PKCS#12 key store holding one EC test key, password `testpassword`, used as a self-hosted RFC 3161 time-stamp authority through `spi/validation.KeyEntityTSPSource` so the T-level examples need no network. TSA policy OID `1.2.3.4.5.6.7.8.9` is an unregistered placeholder. |
| `sample.pdf` | A minimal one-page PDF 1.4 document, written for these tests, so the PAdES examples have something to sign. |

The examples sign these fixtures and validate the result in the same process,
which is why they need no signed documents checked in.
