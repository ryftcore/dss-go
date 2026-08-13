// Generates the ground truth for cms/custom_content_signer.go's
// signatureAlgorithmIdentifierHex table (every JCE signature algorithm name DSS's
// enumerations.SignatureAlgorithm.JCEID() can produce, mapped to the DER of the signature
// AlgorithmIdentifier org.bouncycastle.operator.DefaultSignatureAlgorithmIdentifierFinder#find
// resolves it to) and for cms_signed_attribute_table_generator.go's cmsAlgorithmProtectionDER
// (the CMSAlgorithmProtection encoding BouncyCastle's own class produces, used to confirm the
// [1] signatureAlgorithm field is IMPLICIT rather than EXPLICIT tagged) and for
// signer_attribute_v2.go's DER (confirming a plain two-argument
// `new DERTaggedObject(tagNo, obj)` is EXPLICIT).
//
// Run it with OpenJDK 21 and BouncyCastle 1.84 - the exact version dss-upstream's pom.xml pins
// (<bouncycastle.version>1.84</bouncycastle.version>) - from the local Maven repository:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/cmsapioracle GenSignatureAlgorithmIdentifiers.java
//   java -cp "$CP:/tmp/cmsapioracle" GenSignatureAlgorithmIdentifiers
//
// bcpkix-jdk18on additionally needs bcutil-jdk18on on the classpath at runtime (not just
// compile time): DefaultSignatureAlgorithmIdentifierFinder's static initialiser reaches a
// class (org.bouncycastle.asn1.rosstandart.RosstandartObjectIdentifiers) that 1.84 only ships
// in bcutil's copy of the multi-release JAR, not bcprov's.
import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.DERNull;
import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.DERTaggedObject;
import org.bouncycastle.asn1.cms.CMSAlgorithmProtection;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.operator.DefaultSignatureAlgorithmIdentifierFinder;

public class GenSignatureAlgorithmIdentifiers {

    // Every JCE name enumerations.SignatureAlgorithm.JCEID() (dss-enumerations) can produce for
    // a family CAdES signing actually uses (RSA, RSASSA-PSS, ECDSA, PLAIN-ECDSA, EdDSA, DSA).
    // HMAC (a MAC, not a signature) and the four "NONEwith*" raw algorithms are deliberately
    // exercised too, to confirm they are the ones BouncyCastle's finder itself rejects.
    static final String[] JCE_IDS = {
        "NONEwithRSA",
        "SHA1withRSA", "SHA224withRSA", "SHA256withRSA", "SHA384withRSA", "SHA512withRSA",
        "SHA3-224withRSA", "SHA3-256withRSA", "SHA3-384withRSA", "SHA3-512withRSA",
        "NONEwithRSAandMGF1",
        "SHA1withRSAandMGF1", "SHA224withRSAandMGF1", "SHA256withRSAandMGF1",
        "SHA384withRSAandMGF1", "SHA512withRSAandMGF1",
        "SHA3-224withRSAandMGF1", "SHA3-256withRSAandMGF1", "SHA3-384withRSAandMGF1", "SHA3-512withRSAandMGF1",
        "RIPEMD160withRSA",
        "MD5withRSA", "MD2withRSA",
        "NONEwithECDSA",
        "SHA1withECDSA", "SHA224withECDSA", "SHA256withECDSA", "SHA384withECDSA", "SHA512withECDSA",
        "RIPEMD160withECDSA",
        "SHA3-224withECDSA", "SHA3-256withECDSA", "SHA3-384withECDSA", "SHA3-512withECDSA",
        "SHA1withPLAIN-ECDSA", "SHA224withPLAIN-ECDSA", "SHA256withPLAIN-ECDSA",
        "SHA384withPLAIN-ECDSA", "SHA512withPLAIN-ECDSA",
        "RIPEMD160withPLAIN-ECDSA",
        "SHA3-224withPLAIN-ECDSA", "SHA3-256withPLAIN-ECDSA", "SHA3-384withPLAIN-ECDSA", "SHA3-512withPLAIN-ECDSA",
        "Ed25519", "Ed448",
        "NONEwithDSA",
        "SHA1withDSA", "SHA224withDSA", "SHA256withDSA", "SHA384withDSA", "SHA512withDSA",
        "SHA3-224withDSA", "SHA3-256withDSA", "SHA3-384withDSA", "SHA3-512withDSA",
    };

    public static void main(String[] args) throws Exception {
        DefaultSignatureAlgorithmIdentifierFinder finder = new DefaultSignatureAlgorithmIdentifierFinder();

        System.out.println("# signatureAlgorithmIdentifierHex (DefaultSignatureAlgorithmIdentifierFinder#find)");
        for (String id : JCE_IDS) {
            try {
                AlgorithmIdentifier algId = finder.find(id);
                System.out.println(id + "=" + toHex(algId.getEncoded("DER")));
            } catch (Throwable t) {
                System.out.println(id + "=ERROR:" + t.getClass().getSimpleName() + ":" + t.getMessage());
            }
        }

        System.out.println();
        System.out.println("# CMSAlgorithmProtection tag shape (SIGNATURE alternative)");
        AlgorithmIdentifier digAlg = new AlgorithmIdentifier(new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.1"));
        AlgorithmIdentifier sigAlg = new AlgorithmIdentifier(new ASN1ObjectIdentifier("1.2.840.113549.1.1.11"));
        CMSAlgorithmProtection prot = new CMSAlgorithmProtection(digAlg, CMSAlgorithmProtection.SIGNATURE, sigAlg);
        System.out.println("no-params=" + toHex(prot.getEncoded("DER")));

        AlgorithmIdentifier digAlgWithNull = new AlgorithmIdentifier(
            new ASN1ObjectIdentifier("2.16.840.1.101.3.4.2.1"), DERNull.INSTANCE);
        CMSAlgorithmProtection prot2 = new CMSAlgorithmProtection(digAlgWithNull, CMSAlgorithmProtection.SIGNATURE, sigAlg);
        System.out.println("digest-null-param=" + toHex(prot2.getEncoded("DER")));

        System.out.println();
        System.out.println("# DERTaggedObject(int, ASN1Encodable) two-argument default explicitness");
        DERSequence seq = new DERSequence(new ASN1ObjectIdentifier("1.2.3.4"));
        System.out.println("two-arg=" + toHex(new DERTaggedObject(0, seq).getEncoded("DER")));
        System.out.println("explicit=" + toHex(new DERTaggedObject(true, 0, seq).getEncoded("DER")));
        System.out.println("implicit=" + toHex(new DERTaggedObject(false, 0, seq).getEncoded("DER")));
    }

    static String toHex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) sb.append(String.format("%02x", x));
        return sb.toString();
    }
}
