//go:build ignore

// Command synthetic_certificates writes spi/testdata/certificate_extensions/synthetic_*.der:
//
//	go run synthetic_certificates.go ..
//
// Re-run KatGen.java afterwards to refresh the known answers.
package main

// Builds the synthetic certificates that round out the known-answer corpus: general name types,
// attribute value re-typing, freshestCRL, inhibitAnyPolicy, policyConstraints and name-constraint
// base distances that the real dss-spi test certificates do not exercise.

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/cryptobyte"
	cbasn1 "golang.org/x/crypto/cryptobyte/asn1"
)

func ctx(tag int) cbasn1.Tag { return cbasn1.Tag(tag).ContextSpecific() }
func ctxc(tag int) cbasn1.Tag {
	return cbasn1.Tag(tag).ContextSpecific().Constructed()
}

func build(f func(*cryptobyte.Builder)) []byte {
	b := cryptobyte.NewBuilder(nil)
	f(b)
	out, err := b.Bytes()
	if err != nil {
		panic(err)
	}
	return out
}

func oid(v ...int) asn1.ObjectIdentifier { return asn1.ObjectIdentifier(v) }

// directoryName builds an RDNSequence mixing string types and an unregistered attribute type.
func directoryName() []byte {
	rdn := func(b *cryptobyte.Builder, t asn1.ObjectIdentifier, tag cbasn1.Tag, value []byte) {
		b.AddASN1(cbasn1.SET, func(set *cryptobyte.Builder) {
			set.AddASN1(cbasn1.SEQUENCE, func(ava *cryptobyte.Builder) {
				ava.AddASN1ObjectIdentifier(t)
				ava.AddASN1(tag, func(v *cryptobyte.Builder) { v.AddBytes(value) })
			})
		})
	}
	return build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			rdn(seq, oid(2, 5, 4, 6), cbasn1.PrintableString, []byte("BE"))
			rdn(seq, oid(1, 2, 3, 4, 5), cbasn1.PrintableString, []byte("Foo Bar"))
			rdn(seq, oid(2, 5, 4, 3), cbasn1.UTF8String, []byte("Ünïcode, Inc."))
			rdn(seq, oid(0, 9, 2342, 19200300, 100, 1, 25), cbasn1.IA5String, []byte("example"))
			rdn(seq, oid(2, 5, 4, 5), cbasn1.UTF8String, []byte("12345"))
		})
	})
}

func generalNames(entries ...func(*cryptobyte.Builder)) []byte {
	return build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			for _, entry := range entries {
				entry(seq)
			}
		})
	})
}

func main() {
	out := os.Args[1]

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// --- synthetic_00: every general name type in a subjectAlternativeName ---
	sanEveryType := generalNames(
		func(b *cryptobyte.Builder) { // [0] otherName
			b.AddASN1(ctxc(0), func(other *cryptobyte.Builder) {
				other.AddASN1ObjectIdentifier(oid(1, 2, 3, 4))
				other.AddASN1(ctxc(0), func(v *cryptobyte.Builder) {
					v.AddASN1(cbasn1.UTF8String, func(s *cryptobyte.Builder) { s.AddBytes([]byte("other")) })
				})
			})
		},
		func(b *cryptobyte.Builder) { // [1] rfc822Name
			b.AddASN1(ctx(1), func(v *cryptobyte.Builder) { v.AddBytes([]byte("a@b.example")) })
		},
		func(b *cryptobyte.Builder) { // [2] dNSName
			b.AddASN1(ctx(2), func(v *cryptobyte.Builder) { v.AddBytes([]byte("example.org")) })
		},
		func(b *cryptobyte.Builder) { // [3] x400Address (an empty SEQUENCE body)
			b.AddASN1(ctxc(3), func(v *cryptobyte.Builder) {})
		},
		func(b *cryptobyte.Builder) { // [4] directoryName, EXPLICIT
			b.AddASN1(ctxc(4), func(v *cryptobyte.Builder) { v.AddBytes(directoryName()) })
		},
		func(b *cryptobyte.Builder) { // [5] ediPartyName
			b.AddASN1(ctxc(5), func(v *cryptobyte.Builder) {})
		},
		func(b *cryptobyte.Builder) { // [6] uniformResourceIdentifier, with an '=' to escape
			b.AddASN1(ctx(6), func(v *cryptobyte.Builder) { v.AddBytes([]byte("http://example.org/x?a=b")) })
		},
		func(b *cryptobyte.Builder) { // [7] iPAddress, IPv4
			b.AddASN1(ctx(7), func(v *cryptobyte.Builder) { v.AddBytes([]byte{192, 0, 2, 1}) })
		},
		func(b *cryptobyte.Builder) { // [7] iPAddress, IPv6
			b.AddASN1(ctx(7), func(v *cryptobyte.Builder) {
				v.AddBytes([]byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01})
			})
		},
		func(b *cryptobyte.Builder) { // [7] iPAddress, IPv4-mapped IPv6
			b.AddASN1(ctx(7), func(v *cryptobyte.Builder) {
				v.AddBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 10, 1, 2, 3})
			})
		},
		func(b *cryptobyte.Builder) { // [8] registeredID
			b.AddASN1(ctx(8), func(v *cryptobyte.Builder) { v.AddBytes([]byte{0x2b, 0x06, 0x01, 0x04, 0x01, 0x86, 0x8d, 0x1f}) })
		},
	)

	// --- synthetic_01: a subjectAlternativeName with an unsupported tag ---
	sanBadTag := generalNames(
		func(b *cryptobyte.Builder) {
			b.AddASN1(ctx(2), func(v *cryptobyte.Builder) { v.AddBytes([]byte("good.example")) })
		},
		func(b *cryptobyte.Builder) { // [9] does not exist
			b.AddASN1(ctx(9), func(v *cryptobyte.Builder) { v.AddBytes([]byte("bad")) })
		},
	)

	// --- synthetic_02: freshestCRL, inhibitAnyPolicy, policyConstraints, name constraints ---
	fullNameDistributionPoint := func(seq *cryptobyte.Builder) {
		{
			seq.AddASN1(cbasn1.SEQUENCE, func(dp *cryptobyte.Builder) {
				dp.AddASN1(ctxc(0), func(name *cryptobyte.Builder) { // EXPLICIT DistributionPointName
					name.AddASN1(ctxc(0), func(full *cryptobyte.Builder) { // [0] fullName
						full.AddASN1(ctx(6), func(v *cryptobyte.Builder) {
							v.AddBytes([]byte("http://example.org/delta.crl"))
						})
						full.AddASN1(ctx(2), func(v *cryptobyte.Builder) { // not a URI: skipped
							v.AddBytes([]byte("example.org"))
						})
					})
				})
			})
		}
	}
	// cRLDistributionPoints: crypto/x509 only accepts fullName distribution points.
	distributionPoints := build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, fullNameDistributionPoint)
	})
	// freshestCRL is not parsed by crypto/x509, so it can also carry a nameRelativeToCRLIssuer
	// distribution point, which upstream skips.
	deltaDistributionPoints := build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			fullNameDistributionPoint(seq)
			seq.AddASN1(cbasn1.SEQUENCE, func(dp *cryptobyte.Builder) {
				dp.AddASN1(ctxc(0), func(name *cryptobyte.Builder) {
					name.AddASN1(ctxc(1), func(relative *cryptobyte.Builder) { // nameRelativeToCRLIssuer
						relative.AddASN1(cbasn1.SEQUENCE, func(ava *cryptobyte.Builder) {
							ava.AddASN1ObjectIdentifier(oid(2, 5, 4, 3))
							ava.AddASN1(cbasn1.UTF8String, func(v *cryptobyte.Builder) { v.AddBytes([]byte("crl")) })
						})
					})
				})
			})
		})
	})
	nameConstraints := build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			seq.AddASN1(ctxc(0), func(permitted *cryptobyte.Builder) {
				permitted.AddASN1(cbasn1.SEQUENCE, func(subtree *cryptobyte.Builder) {
					subtree.AddASN1(ctx(2), func(v *cryptobyte.Builder) { v.AddBytes([]byte("example.org")) })
					subtree.AddASN1(ctx(0), func(v *cryptobyte.Builder) { v.AddBytes([]byte{2}) })  // minimum
					subtree.AddASN1(ctx(1), func(v *cryptobyte.Builder) { v.AddBytes([]byte{10}) }) // maximum
				})
				permitted.AddASN1(cbasn1.SEQUENCE, func(subtree *cryptobyte.Builder) {
					subtree.AddASN1(ctxc(4), func(v *cryptobyte.Builder) { v.AddBytes(directoryName()) })
				})
				permitted.AddASN1(cbasn1.SEQUENCE, func(subtree *cryptobyte.Builder) {
					subtree.AddASN1(ctx(8), func(v *cryptobyte.Builder) {
						v.AddBytes([]byte{0x2b, 0x06, 0x01, 0x04, 0x01, 0x86, 0x8d, 0x1f})
					})
				})
			})
			seq.AddASN1(ctxc(1), func(excluded *cryptobyte.Builder) {
				excluded.AddASN1(cbasn1.SEQUENCE, func(subtree *cryptobyte.Builder) {
					subtree.AddASN1(ctx(1), func(v *cryptobyte.Builder) { v.AddBytes([]byte(".bad.example")) })
					subtree.AddASN1(ctx(1), func(v *cryptobyte.Builder) { v.AddBytes([]byte{5}) }) // maximum only
				})
			})
		})
	})
	policyConstraints := build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			seq.AddASN1(ctx(0), func(v *cryptobyte.Builder) { v.AddBytes([]byte{3}) })
			seq.AddASN1(ctx(1), func(v *cryptobyte.Builder) { v.AddBytes([]byte{7}) })
		})
	})
	inhibitAnyPolicy := build(func(b *cryptobyte.Builder) { b.AddASN1Int64(5) })

	// --- synthetic_03: QCStatements exercising the branches the real certificates miss ---
	qcStatements := build(func(b *cryptobyte.Builder) {
		b.AddASN1(cbasn1.SEQUENCE, func(seq *cryptobyte.Builder) {
			// QcCompliance
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 1))
			})
			// QcLimitValue with an alphabetic currency and a non-zero amount/exponent
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 2))
				st.AddASN1(cbasn1.SEQUENCE, func(mv *cryptobyte.Builder) {
					mv.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("USD")) })
					mv.AddASN1Int64(1500)
					mv.AddASN1Int64(-2)
				})
			})
			// QcRetentionPeriod
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 3))
				st.AddASN1Int64(15)
			})
			// QcSSCD
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 4))
			})
			// QcPds, with a URL carrying characters IETFUtils escapes
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 5))
				st.AddASN1(cbasn1.SEQUENCE, func(locations *cryptobyte.Builder) {
					locations.AddASN1(cbasn1.SEQUENCE, func(l *cryptobyte.Builder) {
						l.AddASN1(cbasn1.IA5String, func(v *cryptobyte.Builder) {
							v.AddBytes([]byte("https://example.org/pds?lang=en,v=1"))
						})
						l.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("en")) })
					})
					locations.AddASN1(cbasn1.SEQUENCE, func(l *cryptobyte.Builder) {
						l.AddASN1(cbasn1.IA5String, func(v *cryptobyte.Builder) { v.AddBytes([]byte("  spaced  ")) })
						l.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("fr")) })
					})
				})
			})
			// QcType with a registered and an unregistered type
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 6))
				st.AddASN1(cbasn1.SEQUENCE, func(types *cryptobyte.Builder) {
					types.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 6, 1))
					types.AddASN1ObjectIdentifier(oid(1, 2, 3, 4, 5, 6))
				})
			})
			// QcCClegislation
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 7))
				st.AddASN1(cbasn1.SEQUENCE, func(codes *cryptobyte.Builder) {
					codes.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("BE")) })
					codes.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("LU")) })
				})
			})
			// QcIdentMethod
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 8))
				st.AddASN1(cbasn1.SEQUENCE, func(methods *cryptobyte.Builder) {
					methods.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 8, 3))
				})
			})
			// QcQSCDlegislation
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 9))
				st.AddASN1(cbasn1.SEQUENCE, func(codes *cryptobyte.Builder) {
					codes.AddASN1(cbasn1.PrintableString, func(v *cryptobyte.Builder) { v.AddBytes([]byte("FR")) })
				})
			})
			// QcSemanticsIdentifier (RFC 3739 pkixQCSyntax-v2)
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(1, 3, 6, 1, 5, 5, 7, 11, 2))
				st.AddASN1(cbasn1.SEQUENCE, func(si *cryptobyte.Builder) {
					si.AddASN1ObjectIdentifier(oid(0, 4, 0, 194121, 1, 3))
				})
			})
			// Two unsupported statements, one of which carries no info at all
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(1, 2, 3, 4, 5, 6, 7))
			})
			seq.AddASN1(cbasn1.SEQUENCE, func(st *cryptobyte.Builder) {
				st.AddASN1ObjectIdentifier(oid(1, 2, 3, 4, 5, 6, 8))
				st.AddASN1(cbasn1.SEQUENCE, func(info *cryptobyte.Builder) { info.AddASN1Int64(1) })
			})
			// ocsp-nocheck style DER NULL info on a QcStatement is not a QcStatement branch, but a
			// bare OBJECT IDENTIFIER element is skipped by QCStatement.getInstance.
			seq.AddASN1ObjectIdentifier(oid(0, 4, 0, 1862, 1, 1))
		})
	})

	templates := []struct {
		name     string
		template *x509.Certificate
	}{
		{"synthetic_00.der", &x509.Certificate{
			ExtraExtensions: []pkix.Extension{
				{Id: oid(2, 5, 29, 17), Value: sanEveryType},
			},
		}},
		{"synthetic_01.der", &x509.Certificate{
			ExtraExtensions: []pkix.Extension{
				{Id: oid(2, 5, 29, 17), Critical: true, Value: sanBadTag},
			},
		}},
		{"synthetic_02.der", &x509.Certificate{
			ExtraExtensions: []pkix.Extension{
				{Id: oid(2, 5, 29, 46), Value: deltaDistributionPoints},             // freshestCRL
				{Id: oid(2, 5, 29, 31), Value: distributionPoints},                  // cRLDistributionPoints
				{Id: oid(2, 5, 29, 30), Critical: true, Value: nameConstraints},     // nameConstraints
				{Id: oid(2, 5, 29, 36), Critical: true, Value: policyConstraints},   // policyConstraints
				{Id: oid(2, 5, 29, 54), Critical: true, Value: inhibitAnyPolicy},    // inhibitAnyPolicy
				{Id: oid(1, 3, 6, 1, 5, 5, 7, 48, 1, 5), Value: []byte{0x05, 0x00}}, // ocsp-nocheck
				{Id: oid(0, 4, 0, 194121, 2, 1), Value: []byte{0x05, 0x00}},         // valassured-ST-certs
				{Id: oid(2, 5, 29, 56), Value: []byte{0x05, 0x00}},                  // noRevAvail
			},
			IsCA:                  true,
			BasicConstraintsValid: true,
			MaxPathLen:            3,
		}},
		{"synthetic_03.der", &x509.Certificate{
			ExtraExtensions: []pkix.Extension{
				{Id: oid(1, 3, 6, 1, 5, 5, 7, 1, 3), Critical: true, Value: qcStatements},
			},
		}},
	}

	for i, entry := range templates {
		template := entry.template
		template.SerialNumber = big.NewInt(int64(1000 + i))
		template.Subject = pkix.Name{CommonName: "synthetic-" + string(rune('a'+i))}
		template.NotBefore = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		template.NotAfter = time.Date(2034, 1, 1, 0, 0, 0, 0, time.UTC)
		der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
		if err != nil {
			panic(entry.name + ": " + err.Error())
		}
		if _, err := x509.ParseCertificate(der); err != nil {
			panic(entry.name + ": reparse: " + err.Error())
		}
		if err := os.WriteFile(filepath.Join(out, entry.name), der, 0o644); err != nil {
			panic(err)
		}
	}
}
