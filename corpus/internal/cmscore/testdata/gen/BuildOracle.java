// Generates testdata/build/*.der: SignedData structures assembled with BouncyCastle's
// low-level org.bouncycastle.asn1.cms classes, which derive both CMSVersions themselves.
// TestBuilderMatchesBouncyCastle rebuilds each of them through cmscore's public builders from
// the same inputs and requires the DER to be byte-identical, which pins the SET OF ordering,
// the RFC 5652 clause 5.1/5.3 version computation and the NULL-versus-absent handling of
// algorithm parameters against an independent implementation.
//
// Run it with OpenJDK 21 and BouncyCastle 1.84 from the local Maven repository:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/cmsoracle BuildOracle.java
//   java -cp "$CP:/tmp/cmsoracle" BuildOracle <output directory> <testdata directory>
//
// manifest.txt records the versions BouncyCastle computed, for reference.
// Builds SignedData structures with BouncyCastle's low-level ASN.1 classes - which derive both
// CMSVersions themselves - from inputs the Go side can reconstruct byte for byte, and writes
// the DER of each. cmscore's builders must reproduce these exactly.
import org.bouncycastle.asn1.*;
import org.bouncycastle.asn1.cms.*;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.Certificate;
import org.bouncycastle.asn1.x509.CertificateList;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

public class BuildOracle {
    static Path out;
    static byte[] content = "adversarial content".getBytes();
    static byte[] signature = new byte[64];
    static byte[] ski = new byte[]{0x0a, 0x0b, 0x0c, 0x0d, 0x0e};
    static ASN1ObjectIdentifier idData = CMSObjectIdentifiers.data;
    static ASN1ObjectIdentifier idTst = new ASN1ObjectIdentifier("1.2.840.113549.1.9.16.1.4");
    static AlgorithmIdentifier sha256 = new AlgorithmIdentifier(new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.1"));
    static AlgorithmIdentifier sha512 = new AlgorithmIdentifier(new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.3"));
    static AlgorithmIdentifier rsa = new AlgorithmIdentifier(
        new ASN1ObjectIdentifier("1.2.840.113549.1.1.1"), DERNull.INSTANCE);
    static AlgorithmIdentifier ed25519 = new AlgorithmIdentifier(new ASN1ObjectIdentifier("1.3.101.112"));
    static Certificate cert1, cert2;
    static CertificateList crl;
    static ASN1Sequence ocsp;
    static StringBuilder manifest = new StringBuilder();

    public static void main(String[] args) throws Exception {
        out = Paths.get(args[0]);
        Files.createDirectories(out);
        Path fixtures = Paths.get(args[1]);
        for (int i = 0; i < signature.length; i++) signature[i] = (byte) (i + 1);

        SignedData chain = SignedData.getInstance(ContentInfo.getInstance(
            ASN1Primitive.fromByteArray(Files.readAllBytes(fixtures.resolve("rsa-sha256-chain.p7s")))).getContent());
        cert1 = Certificate.getInstance(chain.getCertificates().getObjectAt(0));
        cert2 = Certificate.getInstance(chain.getCertificates().getObjectAt(1));
        SignedData oc = SignedData.getInstance(ContentInfo.getInstance(
            ASN1Primitive.fromByteArray(Files.readAllBytes(fixtures.resolve("ocsp-crl.p7s")))).getContent());
        for (int i = 0; i < oc.getCRLs().size(); i++) {
            ASN1Encodable o = oc.getCRLs().getObjectAt(i);
            if (o instanceof ASN1TaggedObject) {
                OtherRevocationInfoFormat other = OtherRevocationInfoFormat.getInstance(
                    ASN1Sequence.getInstance((ASN1TaggedObject) o, false));
                if (other.getInfoFormat().getId().equals("1.3.6.1.5.5.7.16.2"))
                    ocsp = ASN1Sequence.getInstance(other.getInfo());
            } else {
                crl = CertificateList.getInstance(o);
            }
        }

        sidIas = IssuerAndSerialNumber.getInstance(
            SignerInfo.getInstance(chain.getSignerInfos().getObjectAt(0)).getSID().getId());

        // Signed attributes, deliberately in reverse DER order.
        ASN1EncodableVector rev = new ASN1EncodableVector();
        rev.add(smime());
        rev.add(digestAttr());
        rev.add(timeAttr());
        rev.add(typeAttr(idData));

        SignerInfo si = new SignerInfo(ias(), sha256, new DERSet(rev), rsa,
            new DEROctetString(signature), null);
        emit("b1-basic", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null, new DERSet(si)));

        SignerInfo siSki = new SignerInfo(new SignerIdentifier(new DEROctetString(ski)), sha256,
            new DERSet(rev), rsa, new DEROctetString(signature), null);
        emit("b2-ski", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null, new DERSet(siSki)));

        emit("b3-detached", new SignedData(new DERSet(sha256), new ContentInfo(idData, null),
            new DERSet(cert1), null, new DERSet(si)));

        // Two certificates supplied in reverse order: the DER SET OF has to sort them.
        emit("b4-two-certs", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(new ASN1Encodable[]{cert2, cert1}), null, new DERSet(si)));

        ASN1EncodableVector revoc = new ASN1EncodableVector();
        revoc.add(new DERTaggedObject(false, 1, new OtherRevocationInfoFormat(
            new ASN1ObjectIdentifier("1.3.6.1.5.5.7.16.2"), ocsp)));
        revoc.add(crl);
        emit("b5-crl-other", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), new DERSet(revoc), new DERSet(si)));

        emit("b6-attr-cert-v2", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(new ASN1Encodable[]{cert1, new DERTaggedObject(false, 2, fake())}),
            null, new DERSet(si)));
        emit("b7-attr-cert-v1", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(new ASN1Encodable[]{cert1, new DERTaggedObject(false, 1, fake())}),
            null, new DERSet(si)));
        emit("b8-other-cert", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(new ASN1Encodable[]{cert1, new DERTaggedObject(false, 3, fake())}),
            null, new DERSet(si)));
        emit("b9-tstinfo-type", new SignedData(new DERSet(sha256), eci(idTst, content),
            new DERSet(cert1), null, new DERSet(si)));

        ASN1EncodableVector unsigned = new ASN1EncodableVector();
        unsigned.add(new Attribute(new ASN1ObjectIdentifier("1.2.840.113549.1.9.6"),
            new DERSet(new DEROctetString(new byte[]{1, 2, 3}))));
        unsigned.add(new Attribute(new ASN1ObjectIdentifier("1.2.3.4.5"),
            new DERSet(new ASN1Integer(42))));
        SignerInfo siUnsigned = new SignerInfo(ias(), sha256, new DERSet(rev), rsa,
            new DEROctetString(signature), new DERSet(unsigned));
        emit("b10-unsigned-attrs", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null, new DERSet(siUnsigned)));

        // One attribute holding three values, supplied out of order.
        ASN1EncodableVector values = new ASN1EncodableVector();
        values.add(new DEROctetString(new byte[]{(byte) 0xff, (byte) 0xff}));
        values.add(new DEROctetString(new byte[]{0x00}));
        values.add(new DEROctetString(new byte[]{0x00, 0x01}));
        ASN1EncodableVector multi = new ASN1EncodableVector();
        multi.add(typeAttr(idData));
        multi.add(new Attribute(new ASN1ObjectIdentifier("1.2.3.4.99"), new DERSet(values)));
        SignerInfo siMulti = new SignerInfo(ias(), sha256, new DERSet(multi), rsa,
            new DEROctetString(signature), null);
        emit("b11-multivalue", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null, new DERSet(siMulti)));

        // Two signers with different digest algorithms; digestAlgorithms derived from them.
        SignerInfo si512 = new SignerInfo(new SignerIdentifier(new DEROctetString(ski)), sha512,
            new DERSet(rev), ed25519, new DEROctetString(signature), null);
        ASN1EncodableVector signers = new ASN1EncodableVector();
        signers.add(si);
        signers.add(si512);
        emit("b12-two-signers", new SignedData(new DERSet(new ASN1Encodable[]{sha256, sha512}),
            eci(idData, content), new DERSet(cert1), null, new DERSet(signers)));

        // RFC 8419: both algorithm identifiers carry no parameters at all.
        SignerInfo siEd = new SignerInfo(ias(), sha512, new DERSet(rev), ed25519,
            new DEROctetString(signature), null);
        emit("b13-ed25519", new SignedData(new DERSet(sha512), eci(idData, content),
            new DERSet(cert1), null, new DERSet(siEd)));

        // No signed attributes at all.
        SignerInfo siNoAttr = new SignerInfo(ias(), sha256, (ASN1Set) null, rsa,
            new DEROctetString(signature), (ASN1Set) null);
        emit("b14-no-signed-attrs", new SignedData(new DERSet(sha256), eci(idData, content),
            new DERSet(cert1), null, new DERSet(siNoAttr)));

        Files.write(out.resolve("manifest.txt"), manifest.toString().getBytes());
    }

    static void emit(String name, SignedData sd) throws Exception {
        Files.write(out.resolve(name + ".der"),
            new ContentInfo(CMSObjectIdentifiers.signedData, sd).getEncoded(ASN1Encoding.DER));
        manifest.append(name).append(" version=").append(sd.getVersion().getValue());
        ASN1Set sis = sd.getSignerInfos();
        for (int i = 0; i < sis.size(); i++)
            manifest.append(" signer").append(i).append('=')
                .append(SignerInfo.getInstance(sis.getObjectAt(i)).getVersion().getValue());
        manifest.append('\n');
    }

    static IssuerAndSerialNumber sidIas;

    static SignerIdentifier ias() {
        return new SignerIdentifier(sidIas);
    }

    static ContentInfo eci(ASN1ObjectIdentifier t, byte[] c) {
        return new ContentInfo(t, new DEROctetString(c));
    }

    static Attribute typeAttr(ASN1ObjectIdentifier t) {
        return new Attribute(CMSAttributes.contentType, new DERSet(t));
    }

    static Attribute digestAttr() {
        byte[] d = new byte[32];
        for (int i = 0; i < 32; i++) d[i] = (byte) (0x40 + i);
        return new Attribute(CMSAttributes.messageDigest, new DERSet(new DEROctetString(d)));
    }

    static Attribute timeAttr() {
        return new Attribute(CMSAttributes.signingTime, new DERSet(new DERUTCTime("240102030405Z")));
    }

    static Attribute smime() {
        ASN1EncodableVector caps = new ASN1EncodableVector();
        for (int i = 0; i < 6; i++)
            caps.add(new DERSequence(new ASN1ObjectIdentifier("1.2.840.113549.3." + (i + 2))));
        return new Attribute(new ASN1ObjectIdentifier("1.2.840.113549.1.9.15"),
            new DERSet(new DERSequence(caps)));
    }

    static ASN1Sequence fake() {
        return new DERSequence(new ASN1Encodable[]{new ASN1Integer(1),
            new DEROctetString(new byte[]{9, 9, 9})});
    }
}
