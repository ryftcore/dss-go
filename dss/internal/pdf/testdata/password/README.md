# Password-charset fixtures

Five one-page PDFs protected by pdfbox 3.0.6's `StandardSecurityHandler`
(`PasswordFixtures.java` in this directory, run by hand - see its header),
pinning how the password's *text* becomes the bytes the handler hashes:
`String.getBytes(ISO_8859_1)` for `/R 2-4`, `SaslPrep` then
`String.getBytes(UTF_8)` for `/R 6`. `crypt.go`'s `passwordBytes` is the port
of that step and `password_kat_test.go` the known-answer test. The pdfbox
version is 3.0.6 rather than the 3.0.7 the package mirrors elsewhere
(`DESIGN.md`); the charset code is identical in both (PDFBOX-4155).

| File | `/V` `/R` | Cipher | User password | Owner password | Pins |
|---|---|---|---|---|---|
| `rc4_r3_latin1.pdf` | 2 3 | RC4 128 | `café` | `ownér` | ISO-8859-1 for RC4 |
| `aes128_r4_latin1.pdf` | 4 4 | AES-128 `/AESV2` | `café` | `ownér` | ISO-8859-1 for AES-128, the common case |
| `aes128_r4_unmappable.pdf` | 4 4 | AES-128 `/AESV2` | `caf€` (U+20AC) | `owner` | a code point above U+00FF hashes as `?`, so `caf?` opens it too |
| `aes256_r6_utf8.pdf` | 5 6 | AES-256 `/AESV3` | `café` | `ownér` | UTF-8 for R6, no transcoding |
| `aes256_r6_saslprep.pdf` | 5 6 | AES-256 `/AESV3` | `ﬁsh` U+00AD `ca` U+00A0 `fé` (U+FB01 ligature, soft hyphen, no-break space) | `owner` | SASLprep: NFKC, map-to-nothing, non-ASCII space - stored as `fishca fé`, which opens it too |

The passwords are test values only.
