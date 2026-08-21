// Generates the ESF / Common PKI known-answer fixtures the CRLOCSP chunk tests read.
// Compile & run with BouncyCastle on the classpath; writes kat_esf.txt next to itself.
import org.bouncycastle.asn1.*;
import org.bouncycastle.asn1.esf.*;
import org.bouncycastle.asn1.isismtt.ocsp.CertHash;
import org.bouncycastle.asn1.nist.NISTObjectIdentifiers;
import org.bouncycastle.asn1.ocsp.ResponderID;
import org.bouncycastle.asn1.oiw.OIWObjectIdentifiers;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.util.encoders.Hex;

import java.io.PrintWriter;
import java.math.BigInteger;
import java.nio.file.*;
import java.util.Date;

public class EsfKatGen {

    static StringBuilder out = new StringBuilder();

    static void put(String k, String v) { out.append(k).append('|').append(v).append('\n'); }
    static void putDer(String k, ASN1Object o) throws Exception { put(k, Hex.toHexString(o.getEncoded(ASN1Encoding.DER))); }

    public static void main(String[] args) throws Exception {
        byte[] sha256 = new byte[32];
        for (int i = 0; i < 32; i++) sha256[i] = (byte) (i + 1);
        byte[] sha1 = new byte[20];
        for (int i = 0; i < 20; i++) sha1[i] = (byte) (0xA0 + i);
        byte[] ski = Hex.decode("fc6c3388b0c280647c753f9e892e9ad705619bf7");

        X500Name issuer = new X500Name("C=LU,O=Nowina Test,CN=Test CRL Issuer");
        // A fixed instant: 2026-08-12T12:23:42Z
        Date issuedTime = new Date(1786537422000L);

        OtherHashAlgAndValue algAndValue =
                new OtherHashAlgAndValue(new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256, DERNull.INSTANCE),
                        new DEROctetString(sha256));
        OtherHash otherHashSha256 = new OtherHash(algAndValue);
        OtherHash otherHashSha1 = new OtherHash(sha1);

        // ---- CrlValidatedID with a full CrlIdentifier -------------------------------
        CrlIdentifier crlIdentifier = new CrlIdentifier(issuer, new ASN1UTCTime(issuedTime), BigInteger.valueOf(0x1234));
        CrlValidatedID full = new CrlValidatedID(otherHashSha256, crlIdentifier);
        putDer("crl.full.der", full);
        put("crl.full.hashalg", NISTObjectIdentifiers.id_sha256.getId());
        put("crl.full.hashvalue", Hex.toHexString(sha256));
        put("crl.full.issuer", Hex.toHexString(issuer.getEncoded(ASN1Encoding.DER)));
        put("crl.full.issuedTime", Long.toString(crlIdentifier.getCrlIssuedTime().getDate().getTime()));
        put("crl.full.number", crlIdentifier.getCrlNumber().toString());

        // ---- CrlValidatedID with a CrlIdentifier without crlNumber -------------------
        CrlIdentifier noNumber = new CrlIdentifier(issuer, new ASN1UTCTime(issuedTime));
        CrlValidatedID noNum = new CrlValidatedID(otherHashSha256, noNumber);
        putDer("crl.nonumber.der", noNum);

        // ---- CrlValidatedID with the sha1Hash alternative and no identifier ----------
        CrlValidatedID hashOnly = new CrlValidatedID(otherHashSha1);
        putDer("crl.hashonly.der", hashOnly);
        put("crl.hashonly.hashalg", OIWObjectIdentifiers.idSHA1.getId());
        put("crl.hashonly.hashvalue", Hex.toHexString(sha1));

        // ---- OcspResponsesID, responder by name --------------------------------------
        X500Name responder = new X500Name("C=LU,O=Nowina Test,CN=Test OCSP CA");
        ResponderID byName = new ResponderID(responder);
        OcspIdentifier byNameId = new OcspIdentifier(byName, new ASN1GeneralizedTime(issuedTime));
        OcspResponsesID byNameResp = new OcspResponsesID(byNameId, otherHashSha256);
        putDer("ocsp.byname.der", byNameResp);
        put("ocsp.byname.producedAt", Long.toString(byNameId.getProducedAt().getDate().getTime()));
        put("ocsp.byname.name", Hex.toHexString(responder.getEncoded(ASN1Encoding.DER)));
        put("ocsp.byname.hashalg", NISTObjectIdentifiers.id_sha256.getId());
        put("ocsp.byname.hashvalue", Hex.toHexString(sha256));

        // ---- OcspResponsesID, responder by key, no hash ------------------------------
        ResponderID byKey = new ResponderID(new DEROctetString(ski));
        OcspIdentifier byKeyId = new OcspIdentifier(byKey, new ASN1GeneralizedTime(issuedTime));
        OcspResponsesID byKeyResp = new OcspResponsesID(byKeyId);
        putDer("ocsp.bykey.der", byKeyResp);
        put("ocsp.bykey.keyhash", Hex.toHexString(ski));
        put("ocsp.bykey.producedAt", Long.toString(byKeyId.getProducedAt().getDate().getTime()));

        // ---- CertHash (Common PKI private OCSP extension) ----------------------------
        CertHash certHash = new CertHash(new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256, DERNull.INSTANCE), sha256);
        putDer("certhash.der", certHash);
        put("certhash.hashalg", NISTObjectIdentifiers.id_sha256.getId());
        put("certhash.hashvalue", Hex.toHexString(sha256));

        Path target = Paths.get(args[0]);
        try (PrintWriter writer = new PrintWriter(Files.newBufferedWriter(target))) {
            writer.print(out);
        }
        System.out.println("written " + target);
    }
}
