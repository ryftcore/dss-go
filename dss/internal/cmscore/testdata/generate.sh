#!/bin/sh
# Regenerates the CMS and RFC 3161 fixtures of internal/cmscore.
#
# Everything here is produced by OpenSSL 3.0.13 alone, except ocsp-crl.p7s and
# ed25519-attached.p7s, which OpenSSL 3.0 cannot write; those two come from
# gen/CmsFixtures.java (BouncyCastle 1.84), whose header states how to run it.
#
# The other corpora have their own generators, each documented in its own header:
#
#   adversarial/adv-*.p7s     gen/Adversarial.java   CMS documents only BouncyCastle writes
#   adversarial/craft-*.p7s   gen/craft.py           BER pathologies nothing writes
#   build/*.der               gen/BuildOracle.java   SignedData assembled by BouncyCastle
#   bc-oracle.txt             gen/CmsOracle.java     BouncyCastle's reading of a whole corpus
#
# The keys and certificates are throw-away test material with fixed serials; they exist only
# so the signatures are real and the structures well formed. Nothing here is a secret.
set -eu

cd "$(dirname "$0")"
rm -rf work
mkdir -p work

# ---------------------------------------------------------------------------- keys & certs
openssl req -x509 -newkey rsa:2048 -keyout work/ca.key -out work/ca.crt -days 7300 -nodes \
    -subj "/C=LU/O=DSS Go Port/CN=Test CA" -set_serial 1 -sha256 2>/dev/null

issue() { # issue <name> <keyspec> <digest>
    name=$1; keyspec=$2; digest=$3
    case $keyspec in
    rsa*) openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out "work/$name.key" 2>/dev/null ;;
    ec*) openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out "work/$name.key" 2>/dev/null ;;
    esac
    openssl req -new -key "work/$name.key" -out "work/$name.csr" \
        -subj "/C=LU/O=DSS Go Port/CN=$name" 2>/dev/null
    openssl x509 -req -in "work/$name.csr" -CA work/ca.crt -CAkey work/ca.key \
        -set_serial "0x$(printf %s "$name" | openssl dgst -sha256 -r | cut -c1-8)" \
        -days 7300 "-$digest" -out "work/$name.crt" \
        -extfile /dev/stdin <<EOF 2>/dev/null
basicConstraints=CA:FALSE
keyUsage=digitalSignature,nonRepudiation
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid
EOF
}

issue rsa-signer rsa sha256
issue ec-signer ec sha256
# The TSA needs its own EKU, so it is issued separately.
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out work/tsa.key 2>/dev/null
openssl req -new -key work/tsa.key -out work/tsa.csr -subj "/C=LU/O=DSS Go Port/CN=Test TSA" 2>/dev/null
openssl x509 -req -in work/tsa.csr -CA work/ca.crt -CAkey work/ca.key -set_serial 0x7473 \
    -days 7300 -sha256 -out work/tsa.crt -extfile /dev/stdin <<EOF 2>/dev/null
basicConstraints=CA:FALSE
keyUsage=nonRepudiation
extendedKeyUsage=critical,timeStamping
subjectKeyIdentifier=hash
EOF

printf 'Hello DSS Go port.\n' > work/content.bin

# ---------------------------------------------------------------------------- CMS SignedData
sign() { # sign <output> <signer> <extra openssl cms args...>
    out=$1; signer=$2; shift 2
    openssl cms -sign -binary -noindef -outform DER -in work/content.bin \
        -signer "work/$signer.crt" -inkey "work/$signer.key" -out "$out" "$@"
}

# Attached (opaque) and detached, SHA-256, RSA.
sign rsa-sha256-attached.p7s rsa-signer -nodetach -md sha256
sign rsa-sha256-detached.p7s rsa-signer -md sha256
# SHA-384, RSA, attached.
sign rsa-sha384-attached.p7s rsa-signer -nodetach -md sha384
# ECDSA P-256 / SHA-256, attached and detached.
sign ec-sha256-attached.p7s ec-signer -nodetach -md sha256
sign ec-sha256-detached.p7s ec-signer -md sha256
# Ed25519 (RFC 8419) is NOT produced here: OpenSSL 3.0.13 refuses to sign CMS with an
# EdDSA key ("eddsa_digest_signverify_init:invalid digest"), RFC 8419 support having landed
# in OpenSSL 3.2. ed25519-attached.p7s comes from gen/CmsFixtures.java instead.
# Two certificates in the CertificateSet, exercising the DER SET OF ordering.
sign rsa-sha256-chain.p7s rsa-signer -nodetach -md sha256 -certfile work/ca.crt
# No certificates at all.
sign rsa-sha256-nocerts.p7s rsa-signer -nodetach -md sha256 -nocerts
# No signed attributes: the signature is computed over the content itself.
sign rsa-sha256-noattr.p7s rsa-signer -nodetach -md sha256 -noattr
# subjectKeyIdentifier SignerIdentifier: SignerInfo version 3, SignedData version 3.
sign rsa-sha256-keyid.p7s rsa-signer -nodetach -md sha256 -keyid
# BER: streaming turns eContent into an indefinite-length constructed OCTET STRING and the
# outer ContentInfo/SignedData into indefinite-length SEQUENCEs.
openssl cms -sign -binary -stream -outform DER -in work/content.bin \
    -signer work/rsa-signer.crt -inkey work/rsa-signer.key -nodetach -md sha256 \
    -out ber-stream-attached.p7s

# ---------------------------------------------------------------------------- RFC 3161
cat > work/tsa.cnf <<'EOF'
[ tsa ]
default_tsa = tsa_config
[ tsa_config ]
serial = ./work/tsa-serial
crypto_device = builtin
signer_cert = ./work/tsa.crt
certs = ./work/tsa.crt
signer_key = ./work/tsa.key
signer_digest = sha256
default_policy = 1.2.3.4.1
other_policies = 1.2.3.4.5, 1.2.3.4.6
digests = sha1, sha256, sha384, sha512
accuracy = secs:1, millisecs:500, microsecs:100
clock_precision_digits = 3
ordering = yes
tsa_name = yes
ess_cert_id_chain = no
EOF
echo 01 > work/tsa-serial

# A request with a nonce and certReq, so the token carries a nonce and the TSA certificate.
openssl ts -query -data work/content.bin -sha256 -cert -out work/request-nonce.tsq 2>/dev/null
openssl ts -reply -config work/tsa.cnf -section tsa_config -queryfile work/request-nonce.tsq \
    -out timestamp-response.tsr 2>/dev/null
openssl ts -reply -in timestamp-response.tsr -token_out -out timestamp-token.tst 2>/dev/null

# A request without nonce and without certReq, over SHA-512.
openssl ts -query -data work/content.bin -sha512 -no_nonce -out work/request-plain.tsq 2>/dev/null
openssl ts -reply -config work/tsa.cnf -section tsa_config -queryfile work/request-plain.tsq \
    -out timestamp-response-nononce.tsr 2>/dev/null
openssl ts -reply -in timestamp-response-nononce.tsr -token_out \
    -out timestamp-token-nononce.tst 2>/dev/null

# A rejected response: no token, only the PKIStatusInfo.
openssl ts -query -data work/content.bin -md5 -out work/request-md5.tsq 2>/dev/null
openssl ts -reply -config work/tsa.cnf -section tsa_config -queryfile work/request-md5.tsq \
    -out timestamp-response-rejected.tsr 2>/dev/null || true

# ---------------------------------------------------------------------------- companions
# The signed content, so the tests can check the eContent and the message-digest attribute.
# The certificates are not copied: every fixture but rsa-sha256-nocerts.p7s embeds its own.
cp work/content.bin content.bin

rm -rf work
