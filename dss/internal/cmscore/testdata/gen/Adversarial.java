// Generates testdata/adversarial/adv-*.p7s: CMS documents BouncyCastle can write but no
// ordinary producer emits - unsorted signed attributes, empty certificate and CRL sets,
// several SignerInfos, every CertificateChoices and RevocationInfoChoice alternative, and a
// fully indefinite-length BER document. They exist so the Go parser can be diffed against
// BouncyCastle's reading of them, not against the Go port's own output.
//
// Run it with OpenJDK 21 and BouncyCastle 1.84 from the local Maven repository:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/cmsoracle Adversarial.java
//   java -cp "$CP:/tmp/cmsoracle" Adversarial <output directory> <testdata directory>
//
// The certificates, CRL and OCSP response are taken from rsa-sha256-chain.p7s and
// ocsp-crl.p7s, so a re-run reproduces the same bytes.
// Produces adversarial CMS documents with BouncyCastle, so the Go parser can be diffed
// against BC's interpretation of encodings no ordinary producer emits.
import org.bouncycastle.asn1.*;
import org.bouncycastle.asn1.cms.*;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.Certificate;
import org.bouncycastle.asn1.x509.CertificateList;

import java.io.File;
import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

public class Adversarial {
    static Path outDir;
    static byte[] content = "adversarial content".getBytes();
    static byte[] signature = new byte[64];
    static byte[] ski = new byte[]{0x0a, 0x0b, 0x0c, 0x0d, 0x0e};
    static ASN1ObjectIdentifier idData = CMSObjectIdentifiers.data;
    static ASN1ObjectIdentifier idTst = new ASN1ObjectIdentifier("1.2.840.113549.1.9.16.1.4");
    static AlgorithmIdentifier sha256 = new AlgorithmIdentifier(
        new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.1"));
    static AlgorithmIdentifier rsa = new AlgorithmIdentifier(
        new ASN1ObjectIdentifier("1.2.840.113549.1.1.1"), DERNull.INSTANCE);
    static Certificate cert1, cert2;
    static CertificateList crl;
    static ASN1Sequence ocspResponse;

    public static void main(String[] args) throws Exception {
        outDir = Paths.get(args[0]);
        Files.createDirectories(outDir);
        Path fixtures = Paths.get(args[1]);
        for (int i = 0; i < signature.length; i++) signature[i] = (byte) (i + 1);

        // Reuse real certificates / CRL / OCSP response from the checked-in fixtures.
        SignedData chain = SignedData.getInstance(
            ContentInfo.getInstance(ASN1Primitive.fromByteArray(
                Files.readAllBytes(fixtures.resolve("rsa-sha256-chain.p7s")))).getContent());
        cert1 = Certificate.getInstance(chain.getCertificates().getObjectAt(0));
        cert2 = Certificate.getInstance(chain.getCertificates().getObjectAt(1));
        SignedData ocspCrl = SignedData.getInstance(
            ContentInfo.getInstance(ASN1Primitive.fromByteArray(
                Files.readAllBytes(fixtures.resolve("ocsp-crl.p7s")))).getContent());
        for (int i = 0; i < ocspCrl.getCRLs().size(); i++) {
            ASN1Encodable o = ocspCrl.getCRLs().getObjectAt(i);
            if (o instanceof ASN1TaggedObject) {
                OtherRevocationInfoFormat other = OtherRevocationInfoFormat.getInstance(
                    ASN1Sequence.getInstance((ASN1TaggedObject) o, false));
                if (other.getInfoFormat().getId().equals("1.3.6.1.5.5.7.16.2"))
                    ocspResponse = ASN1Sequence.getInstance(other.getInfo());
            } else {
                crl = CertificateList.getInstance(o);
            }
        }

        // ---- attributes, deliberately supplied out of DER order -------------------------
        ASN1EncodableVector reversed = new ASN1EncodableVector();
        reversed.add(smimeCapLike());                 // longest encoding, sorts last
        reversed.add(messageDigestAttr());
        reversed.add(signingTimeAttr());
        reversed.add(contentTypeAttr(idData));        // shortest, sorts first
        ASN1Set unsortedSignedAttrs = new BERSet(reversed);

        // 1. unsorted signed attributes, one signer, one certificate
        write("adv-unsorted-attrs.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null,
            new DERSet(signerInfo(1, iasSid(), unsortedSignedAttrs, null))));

        // 2. two signers, one with SKI (version 3), the other with issuerAndSerial
        ASN1EncodableVector two = new ASN1EncodableVector();
        two.add(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null));
        two.add(signerInfo(3, skiSid(), new DERSet(sortedAttrs()), null));
        write("adv-two-signers.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(new ASN1Encodable[]{cert1, cert2}), null, new DERSet(two)));

        // 3. certificates present but empty; crls present but empty
        write("adv-empty-sets.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(), new DERSet(),
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 4. signedAttrs present but empty, and unsignedAttrs present
        ASN1EncodableVector unsigned = new ASN1EncodableVector();
        unsigned.add(new Attribute(new ASN1ObjectIdentifier("1.2.840.113549.1.9.6"),
            new DERSet(new DEROctetString(new byte[]{1, 2, 3}))));
        write("adv-empty-signed-attrs.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(), new DERSet(unsigned)))));

        // 5. CRL + OtherRevocationInfoFormat: version 5
        ASN1EncodableVector revocation = new ASN1EncodableVector();
        revocation.add(crl);
        revocation.add(new DERTaggedObject(false, 1, new OtherRevocationInfoFormat(
            new ASN1ObjectIdentifier("1.3.6.1.5.5.7.16.2"), ocspResponse)));
        write("adv-crl-and-other.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), new DERSet(revocation),
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 6. eContentType other than id-data: version 3
        write("adv-other-econtent-type.p7s", signedData(
            new DERSet(sha256), eci(idTst, content),
            new DERSet(cert1), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 7. a [2] IMPLICIT AttributeCertificateV2 in the CertificateSet: version 4
        ASN1EncodableVector v2certs = new ASN1EncodableVector();
        v2certs.add(cert1);
        v2certs.add(new DERTaggedObject(false, 2, fakeAttrCert()));
        write("adv-attr-cert-v2.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(v2certs), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 8. a [1] IMPLICIT AttributeCertificateV1: version 3
        ASN1EncodableVector v1certs = new ASN1EncodableVector();
        v1certs.add(cert1);
        v1certs.add(new DERTaggedObject(false, 1, fakeAttrCert()));
        write("adv-attr-cert-v1.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(v1certs), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 9. a [3] IMPLICIT OtherCertificateFormat: version 5
        ASN1EncodableVector otherCerts = new ASN1EncodableVector();
        otherCerts.add(cert1);
        otherCerts.add(new DERTaggedObject(false, 3, fakeAttrCert()));
        write("adv-other-cert-format.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(otherCerts), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        // 10. detached, no certificates, no signed attributes
        write("adv-detached-noattr.p7s", signedData(
            new DERSet(sha256), new ContentInfo(idData, null), null, null,
            new DERSet(signerInfo(1, iasSid(), null, null))));

        // 11. an attribute carrying three values, supplied out of order
        ASN1EncodableVector values = new ASN1EncodableVector();
        values.add(new DEROctetString(new byte[]{(byte) 0xff, (byte) 0xff}));
        values.add(new DEROctetString(new byte[]{0x00}));
        values.add(new DEROctetString(new byte[]{0x00, 0x01}));
        ASN1EncodableVector multi = new ASN1EncodableVector();
        multi.add(contentTypeAttr(idData));
        multi.add(messageDigestAttr());
        multi.add(new Attribute(new ASN1ObjectIdentifier("1.2.3.4.99"), new BERSet(values)));
        write("adv-multivalued-attr.p7s", signedData(
            new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null,
            new DERSet(signerInfo(1, iasSid(), new BERSet(multi), null))));

        // 12. everything indefinite: BER sequences and a segmented eContent
        ContentInfo berEci = new ContentInfo(idData,
            new BEROctetString(new ASN1OctetString[]{
                new DEROctetString(java.util.Arrays.copyOfRange(content, 0, 5)),
                new DEROctetString(java.util.Arrays.copyOfRange(content, 5, content.length))}));
        SignedData berSd = new SignedData(new BERSet(sha256), berEci,
            new BERSet(cert1), null,
            new BERSet(signerInfo(1, iasSid(), new BERSet(reversedCopy()), null)));
        Files.write(outDir.resolve("adv-ber-indefinite.p7s"),
            new ContentInfo(CMSObjectIdentifiers.signedData, berSd).getEncoded(ASN1Encoding.BER));

        // 13. digestAlgorithms holding two members, out of order
        ASN1EncodableVector algs = new ASN1EncodableVector();
        algs.add(new AlgorithmIdentifier(new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.3")));
        algs.add(sha256);
        write("adv-two-digest-algs.p7s", signedData(
            new BERSet(algs), eci(idData, content),
            new DERSet(cert1), null,
            new DERSet(signerInfo(1, iasSid(), new DERSet(sortedAttrs()), null))));

        System.out.println("written");
    }

    static ASN1Encodable[] reversedCopy() {
        return new ASN1Encodable[]{smimeCapLike(), messageDigestAttr(), signingTimeAttr(),
            contentTypeAttr(idData)};
    }

    static ASN1Encodable[] sortedAttrs() {
        return new ASN1Encodable[]{contentTypeAttr(idData), messageDigestAttr()};
    }

    static Attribute contentTypeAttr(ASN1ObjectIdentifier type) {
        return new Attribute(CMSAttributes.contentType, new DERSet(type));
    }

    static Attribute messageDigestAttr() {
        byte[] digest = new byte[32];
        for (int i = 0; i < 32; i++) digest[i] = (byte) (0x40 + i);
        return new Attribute(CMSAttributes.messageDigest, new DERSet(new DEROctetString(digest)));
    }

    static Attribute signingTimeAttr() {
        return new Attribute(CMSAttributes.signingTime,
            new DERSet(new DERUTCTime("240102030405Z")));
    }

    // A long attribute whose encoding sorts last.
    static Attribute smimeCapLike() {
        ASN1EncodableVector caps = new ASN1EncodableVector();
        for (int i = 0; i < 6; i++)
            caps.add(new DERSequence(new ASN1ObjectIdentifier("1.2.840.113549.3." + (i + 2))));
        return new Attribute(new ASN1ObjectIdentifier("1.2.840.113549.1.9.15"),
            new DERSet(new DERSequence(caps)));
    }

    // A stand-in for an attribute certificate: only the tag number matters to the parser.
    static ASN1Sequence fakeAttrCert() {
        return new DERSequence(new ASN1Encodable[]{new ASN1Integer(1),
            new DEROctetString(new byte[]{9, 9, 9})});
    }

    static ContentInfo eci(ASN1ObjectIdentifier type, byte[] data) {
        return new ContentInfo(type, new DEROctetString(data));
    }

    static SignerIdentifier iasSid() {
        return new SignerIdentifier(new IssuerAndSerialNumber(
            cert1.getIssuer(), cert1.getSerialNumber().getValue()));
    }

    static SignerIdentifier skiSid() {
        return new SignerIdentifier(new DEROctetString(ski));
    }

    static SignerInfo signerInfo(int version, SignerIdentifier sid, ASN1Set signed, ASN1Set unsigned) {
        return new SignerInfo(sid, sha256, signed, rsa, new DEROctetString(signature), unsigned);
    }

    static SignedData signedData(ASN1Set algs, ContentInfo eci, ASN1Set certs, ASN1Set crls,
        ASN1Set signers) {
        return new SignedData(algs, eci, certs, crls, signers);
    }

    static void write(String name, SignedData sd) throws Exception {
        Files.write(outDir.resolve(name),
            new ContentInfo(CMSObjectIdentifiers.signedData, sd).getEncoded(ASN1Encoding.DER));
    }
}
