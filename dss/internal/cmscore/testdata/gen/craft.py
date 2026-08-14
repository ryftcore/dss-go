#!/usr/bin/env python3
"""Generates the testdata/adversarial/craft-*.p7s fixtures: BER and DER pathologies no
ordinary CMS producer emits - and that BouncyCastle itself will not write - built by taking
rsa-sha256-attached.p7s apart and reassembling it with deliberately awkward encodings.

Non-minimal length octets, high tag numbers, constructed OCTET STRING content and signature
values, indefinite lengths in every position, an Attribute SEQUENCE with a third component,
and five documents that have to be rejected. What BouncyCastle makes of each is recorded in
testdata/adversarial/bc-oracle.txt by gen/CmsOracle.java.

    python3 gen/craft.py testdata/adversarial
"""
import os, sys

SRC = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "rsa-sha256-attached.p7s")
OUT = sys.argv[1] if len(sys.argv) > 1 else "craft"


# --------------------------------------------------------------------------- tiny ASN.1 kit
def dlen(n):
    if n < 0x80:
        return bytes([n])
    b = n.to_bytes((n.bit_length() + 7) // 8, "big")
    return bytes([0x80 | len(b)]) + b


def padded_len(n, extra):
    """Non-minimal long-form length: `extra` leading zero octets."""
    b = n.to_bytes((n.bit_length() + 7) // 8, "big") or b"\x00"
    b = b"\x00" * extra + b
    return bytes([0x80 | len(b)]) + b


def tlv(tag, content):
    return bytes(tag if isinstance(tag, (bytes, bytearray)) else [tag]) + dlen(len(content)) + content


def tlv_padded(tag, content, extra=2):
    return bytes(tag if isinstance(tag, (bytes, bytearray)) else [tag]) + padded_len(len(content), extra) + content


def indef(tag, content):
    return bytes(tag if isinstance(tag, (bytes, bytearray)) else [tag]) + b"\x80" + content + b"\x00\x00"


def high_tag(cls_constructed, number):
    """Identifier octets for a high tag number (>= 31)."""
    out = bytearray([cls_constructed | 0x1F])
    stack = []
    v = number
    while True:
        stack.append(v & 0x7F)
        v >>= 7
        if v == 0:
            break
    for i, part in enumerate(reversed(stack)):
        out.append(part | (0x80 if i < len(stack) - 1 else 0x00))
    return bytes(out)


# --------------------------------------------------------------------------- tiny parser
def parse(buf, i=0):
    start = i
    first = buf[i]
    i += 1
    if first & 0x1F == 0x1F:
        num = 0
        while True:
            b = buf[i]
            i += 1
            num = num << 7 | (b & 0x7F)
            if not b & 0x80:
                break
    else:
        num = first & 0x1F
    ln = buf[i]
    i += 1
    if ln == 0x80:
        children = []
        while buf[i:i + 2] != b"\x00\x00":
            child, i = parse(buf, i)
            children.append(child)
        i += 2
        return dict(id=first, num=num, cons=True, indef=True, children=children,
                    body=None, enc=buf[start:i]), i
    if ln & 0x80:
        n = ln & 0x7F
        length = int.from_bytes(buf[i:i + n], "big")
        i += n
    else:
        length = ln
    body = buf[i:i + length]
    i += length
    node = dict(id=first, num=num, cons=bool(first & 0x20), indef=False, body=body,
                enc=buf[start:i], children=[])
    if node["cons"]:
        j = 0
        while j < len(body):
            child, j = parse(body, j)
            node["children"].append(child)
    return node, i


data = open(SRC, "rb").read()
ci, _ = parse(data)
sd = ci["children"][1]["children"][0]                  # SignedData
version, digestAlgs, encap, certs, signerInfos = (sd["children"] + [None] * 5)[:5]
signerInfo = signerInfos["children"][0]
si_children = signerInfo["children"]
si_version, sid, si_digest, signedAttrs, sigAlg, sigValue = si_children[:6]
content_octets = encap["children"][1]["children"][0]["body"]
econtent_type = encap["children"][0]["enc"]

os.makedirs(OUT, exist_ok=True)
for f in os.listdir(OUT):
    if f.startswith("craft-"):
        os.remove(os.path.join(OUT, f))


def emit(name, signed_data_body, wrap=tlv, outer=tlv):
    sd_enc = wrap(0x30, signed_data_body)
    body = bytes.fromhex("06092a864886f70d010702") + outer(0xA0, sd_enc)
    open(os.path.join(OUT, name), "wb").write(outer(0x30, body))


def signer(children):
    return tlv(0x30, b"".join(c["enc"] if isinstance(c, dict) else c for c in children))


def sd_body(encap_enc, certs_enc=None, si_enc=None, algs_enc=None):
    return (version["enc"] + (algs_enc or digestAlgs["enc"]) + encap_enc
            + (certs["enc"] if certs_enc is None else certs_enc)
            + (si_enc if si_enc is not None else tlv(0x31, signer(si_children))))


# h1 non-minimal lengths at every level -------------------------------------------------
emit("craft-nonminimal-length.p7s",
     sd_body(tlv_padded(0x30, encap["children"][0]["enc"]
                        + tlv_padded(0xA0, tlv_padded(0x04, content_octets)))),
     wrap=tlv_padded, outer=tlv_padded)

# h2 high tag numbers inside an attribute value ------------------------------------------
high_value = high_tag(0xC0, 100) + dlen(3) + b"\x01\x02\x03"
extra_attr = tlv(0x30, bytes.fromhex("06052a03040506") + tlv(0x31, high_value))
attrs = [c["enc"] for c in signedAttrs["children"]] + [extra_attr]
si_high = list(si_children)
si_high[3] = tlv(0xA0, b"".join(attrs))
emit("craft-high-tag-number.p7s", sd_body(encap["enc"], si_enc=tlv(0x31, signer(si_high))))

# h3 constructed, definite-length OCTET STRING eContent, three segments -------------------
segments = b"".join(tlv(0x04, content_octets[i:i + 7]) for i in range(0, len(content_octets), 7))
emit("craft-constructed-econtent.p7s",
     sd_body(tlv(0x30, econtent_type + tlv(0xA0, tlv(0x24, segments)))))

# h4 indefinite [0] EXPLICIT eContent wrapper around a primitive OCTET STRING -------------
emit("craft-indefinite-explicit.p7s",
     sd_body(tlv(0x30, econtent_type + indef(0xA0, tlv(0x04, content_octets)))))

# h5 indefinite SignedData inside a definite ContentInfo ----------------------------------
emit("craft-mixed-indefinite.p7s", sd_body(encap["enc"]), wrap=indef)

# h6 signed attributes in reverse order, with a non-minimal length ------------------------
rev = [c["enc"] for c in reversed(signedAttrs["children"])]
si_rev = list(si_children)
si_rev[3] = tlv_padded(0xA0, b"".join(rev))
emit("craft-reversed-attrs.p7s", sd_body(encap["enc"], si_enc=tlv(0x31, signer(si_rev))))

# h7 BER BOOLEAN TRUE (0x01) in an attribute value: DER must normalise it to 0xFF ---------
bool_attr = tlv(0x30, bytes.fromhex("06052a03040507") + tlv(0x31, tlv(0x01, b"\x01")))
si_bool = list(si_children)
si_bool[3] = tlv(0xA0, b"".join(c["enc"] for c in signedAttrs["children"]) + bool_attr)
emit("craft-ber-boolean-attr.p7s", sd_body(encap["enc"], si_enc=tlv(0x31, signer(si_bool))))

# h8 indefinite-length certificates [0] ---------------------------------------------------
emit("craft-indefinite-certs.p7s",
     sd_body(encap["enc"], certs_enc=indef(0xA0, b"".join(c["enc"] for c in certs["children"]))))

# h9 constructed OCTET STRING signature value ---------------------------------------------
si_sig = list(si_children)
sig = sigValue["body"]
si_sig[5] = tlv(0x24, tlv(0x04, sig[:100]) + tlv(0x04, sig[100:]))
emit("craft-constructed-signature.p7s", sd_body(encap["enc"], si_enc=tlv(0x31, signer(si_sig))))

# h10 unsorted digestAlgorithms with a duplicate ------------------------------------------
sha512 = tlv(0x30, bytes.fromhex("0609608648016503040203"))
emit("craft-unsorted-digest-algs.p7s",
     sd_body(encap["enc"], algs_enc=tlv(0x31, sha512 + digestAlgs["children"][0]["enc"])))

# h11 an Attribute SEQUENCE carrying a third component, which a rebuild would silently drop --
third = tlv(0x30, bytes.fromhex("06052a03040508") + tlv(0x31, tlv(0x04, b"\x01\x02"))
            + tlv(0x17, b"240102030405Z"))
si_third = list(si_children)
si_third[3] = tlv(0xA0, b"".join(c["enc"] for c in signedAttrs["children"]) + third)
emit("craft-attribute-third-component.p7s",
     sd_body(encap["enc"], si_enc=tlv(0x31, signer(si_third))))

# h12 a definite-length constructed eContent inside an indefinite SignedData ----------------
emit("craft-definite-string-in-indefinite.p7s",
     sd_body(tlv(0x30, econtent_type + tlv(0xA0, tlv(0x24, segments)))),
     wrap=indef, outer=indef)

# --------------------------------------------------------------------------- rejections
bad = {}
bad["craft-bad-truncated.p7s"] = data[:len(data) - 5]
bad["craft-bad-trailing.p7s"] = data + b"\x00\x00\x00"
# indefinite length on a primitive element
bad["craft-bad-indefinite-primitive.p7s"] = (
    data[:0] + bytes([0x30]) + b"\x80" + bytes.fromhex("06092a864886f70d010702")
    + b"\x04\x80\x01\x02\x00\x00" + b"\x00\x00")
# missing end-of-contents octets
truncated_indef = indef(0x30, bytes.fromhex("06092a864886f70d010702"))[:-2]
bad["craft-bad-missing-eoc.p7s"] = truncated_indef
# length running past the buffer
bad["craft-bad-overlong-length.p7s"] = bytes([0x30, 0x82, 0x7F, 0xFF]) + data[4:40]
for name, blob in bad.items():
    open(os.path.join(OUT, name), "wb").write(blob)

print("\n".join(sorted(os.listdir(OUT))))
