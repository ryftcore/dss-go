// Produces the RFC 3161 fixtures of dss/spi/validation and the BouncyCastle oracle the
// TimestampToken and KeyEntityTSPSource known-answer tests compare against.
//
// The token is issued with exactly the wiring eu.europa.esig.dss.spi.x509.tsp.KeyEntityTSPSource
// uses (JcaContentSignerBuilder + SignerInfoGeneratorBuilder(BcDigestCalculatorProvider) with a
// pinned signingTime + TimeStampTokenGenerator + TimeStampResponseGenerator), so that the Go
// port can be compared with it byte for byte. Run with:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/tspfix TspFixtures.java && java -cp "$CP:/tmp/tspfix" TspFixtures ..
//
// The key material is deterministic (fixed RSA key, fixed serials, fixed validity), so a re-run
// reproduces these files byte for byte. Nothing here is a secret.

import org.bouncycastle.asn1.ASN1EncodableVector;
import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.DERSet;
import org.bouncycastle.asn1.cms.Attribute;
import org.bouncycastle.asn1.cms.AttributeTable;
import org.bouncycastle.asn1.cms.CMSAttributes;
import org.bouncycastle.asn1.cms.Time;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.ExtendedKeyUsage;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.KeyPurposeId;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.jcajce.JcaCertStore;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.cms.CMSAttributeTableGenerator;
import org.bouncycastle.cms.DefaultSignedAttributeTableGenerator;
import org.bouncycastle.cms.SignerInformation;
import org.bouncycastle.cms.SignerInfoGenerator;
import org.bouncycastle.cms.SignerInfoGeneratorBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.DigestCalculator;
import org.bouncycastle.operator.bc.BcDigestCalculatorProvider;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.operator.jcajce.JcaDigestCalculatorProviderBuilder;
import org.bouncycastle.tsp.TimeStampRequest;
import org.bouncycastle.tsp.TimeStampRequestGenerator;
import org.bouncycastle.tsp.TimeStampResponse;
import org.bouncycastle.tsp.TimeStampResponseGenerator;
import org.bouncycastle.tsp.TimeStampToken;
import org.bouncycastle.tsp.TimeStampTokenGenerator;

import java.io.PrintWriter;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.MessageDigest;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.Security;
import java.security.cert.X509Certificate;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Date;
import java.util.HashSet;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.Set;

public final class TspFixtures {

    private static final byte[] CONTENT = "Hello DSS Go port.\n".getBytes(StandardCharsets.UTF_8);

    /** 2021-01-02T03:04:05Z, the pinned generation and signing time. */
    private static final Date GEN_TIME = new Date(1609556645000L);
    /** notBefore / notAfter of the generated certificates. */
    private static final Date NOT_BEFORE = new Date(1577836800000L); // 2020-01-01T00:00:00Z
    private static final Date NOT_AFTER = new Date(2524608000000L);  // 2050-01-01T00:00:00Z

    private static final BigInteger TST_SERIAL = new BigInteger("123456789012345678901234567890");
    private static final String TSA_POLICY = "1.2.3.4.1";

    public static void main(String[] args) throws Exception {
        Security.addProvider(new BouncyCastleProvider());
        Path out = Path.of(args.length > 0 ? args[0] : ".");

        // Deterministic keys: a fixed seed makes the whole corpus reproducible.
        KeyPair caKeyPair = deterministicRsa("ca");
        KeyPair tsaKeyPair = deterministicRsa("tsa");

        X509Certificate caCertificate = issue(caKeyPair, caKeyPair, "CN=Test CA,O=DSS Go Port,C=LU",
                "CN=Test CA,O=DSS Go Port,C=LU", BigInteger.valueOf(1), true, false);
        X509Certificate tsaCertificate = issue(caKeyPair, tsaKeyPair, "CN=Test CA,O=DSS Go Port,C=LU",
                "CN=Test TSA,O=DSS Go Port,C=LU", BigInteger.valueOf(0x7473), false, true);
        // A certificate with the same key as the TSA but no timeStamping EKU, so the port can
        // check that TSPUtil.validateCertificate rejects it.
        X509Certificate plainCertificate = issue(caKeyPair, tsaKeyPair, "CN=Test CA,O=DSS Go Port,C=LU",
                "CN=Plain Signer,O=DSS Go Port,C=LU", BigInteger.valueOf(0x7474), false, false);

        Files.write(out.resolve("ca.crt"), caCertificate.getEncoded());
        Files.write(out.resolve("tsa.crt"), tsaCertificate.getEncoded());
        Files.write(out.resolve("plain.crt"), plainCertificate.getEncoded());
        Files.write(out.resolve("tsa.key"), tsaKeyPair.getPrivate().getEncoded());
        Files.write(out.resolve("content.bin"), CONTENT);

        byte[] sha256 = MessageDigest.getInstance("SHA-256").digest(CONTENT);
        byte[] token = generate(tsaKeyPair.getPrivate(), tsaCertificate,
                List.of(tsaCertificate, caCertificate), "SHA-256", "SHA512withRSA", sha256);
        Files.write(out.resolve("timestamp-token.tst"), token);

        // A second token over SHA-512, signed with SHA-256, whose chain holds the TSA only.
        byte[] sha512 = MessageDigest.getInstance("SHA-512").digest(CONTENT);
        byte[] tokenSha512 = generate(tsaKeyPair.getPrivate(), tsaCertificate, List.of(tsaCertificate),
                "SHA-512", "SHA256withRSA", sha512);
        Files.write(out.resolve("timestamp-token-sha512.tst"), tokenSha512);

        try (PrintWriter oracle = new PrintWriter(Files.newBufferedWriter(out.resolve("bc-oracle.txt")))) {
            describe(oracle, "timestamp-token.tst", token, sha256);
            describe(oracle, "timestamp-token-sha512.tst", tokenSha512, sha512);
        }
    }

    /** Reproduces KeyEntityTSPSource#getTimeStampResponse for one request. */
    private static byte[] generate(PrivateKey privateKey, X509Certificate certificate,
                                   List<X509Certificate> chain, String messageImprintDigest,
                                   String signatureAlgorithm, byte[] digest) throws Exception {
        ASN1ObjectIdentifier digestOid = new ASN1ObjectIdentifier(oidOf(messageImprintDigest));

        // KeyEntityTSPSource#createRequest
        TimeStampRequestGenerator requestGenerator = new TimeStampRequestGenerator();
        requestGenerator.setCertReq(true);
        TimeStampRequest request = requestGenerator.generate(digestOid, digest);

        // KeyEntityTSPSource#initResponseGenerator
        ContentSigner signer = new JcaContentSignerBuilder(signatureAlgorithm).build(privateKey);
        X509CertificateHolder certificateHolder = new X509CertificateHolder(certificate.getEncoded());
        SignerInfoGenerator infoGenerator = new SignerInfoGeneratorBuilder(new BcDigestCalculatorProvider())
                .setSignedAttributeGenerator(signedAttributeGenerator()).build(signer, certificateHolder);
        AlgorithmIdentifier digestAlgorithmIdentifier = new AlgorithmIdentifier(digestOid);
        DigestCalculator digestCalculator = new JcaDigestCalculatorProviderBuilder().build()
                .get(digestAlgorithmIdentifier);
        TimeStampTokenGenerator tokenGenerator = new TimeStampTokenGenerator(infoGenerator, digestCalculator,
                new ASN1ObjectIdentifier(TSA_POLICY));
        tokenGenerator.addCertificates(new JcaCertStore(chain));
        Set<ASN1ObjectIdentifier> accepted = new HashSet<>();
        for (String name : new String[] {"SHA-224", "SHA-256", "SHA-384", "SHA-512"}) {
            accepted.add(new ASN1ObjectIdentifier(oidOf(name)));
        }
        TimeStampResponseGenerator responseGenerator = new TimeStampResponseGenerator(tokenGenerator, accepted);

        TimeStampResponse response = responseGenerator.generate(request, TST_SERIAL, GEN_TIME);
        return response.getTimeStampToken().getEncoded();
    }

    /** KeyEntityTSPSource#getSignedAttributeGenerator, with the production time pinned. */
    @SuppressWarnings({"unchecked", "rawtypes"})
    private static CMSAttributeTableGenerator signedAttributeGenerator() {
        return new DefaultSignedAttributeTableGenerator() {
            @Override
            protected Hashtable createStandardAttributeTable(Map map) {
                Hashtable hashtable = super.createStandardAttributeTable(map);
                Attribute attr = new Attribute(CMSAttributes.signingTime, new DERSet(new Time(GEN_TIME)));
                hashtable.put(CMSAttributes.signingTime, attr);
                return hashtable;
            }
        };
    }

    private static void describe(PrintWriter oracle, String name, byte[] encoded, byte[] imprint) throws Exception {
        TimeStampToken token = new TimeStampToken(new org.bouncycastle.cms.CMSSignedData(encoded));
        oracle.printf("%s.length=%d%n", name, encoded.length);
        oracle.printf("%s.sha256=%s%n", name, hex(MessageDigest.getInstance("SHA-256").digest(encoded)));
        oracle.printf("%s.genTime=%d%n", name, token.getTimeStampInfo().getGenTime().getTime());
        oracle.printf("%s.serialNumber=%s%n", name, token.getTimeStampInfo().getSerialNumber());
        oracle.printf("%s.policy=%s%n", name, token.getTimeStampInfo().getPolicy().getId());
        oracle.printf("%s.messageImprintAlgOID=%s%n", name, token.getTimeStampInfo().getMessageImprintAlgOID().getId());
        oracle.printf("%s.messageImprintDigest=%s%n", name, hex(token.getTimeStampInfo().getMessageImprintDigest()));
        oracle.printf("%s.expectedImprint=%s%n", name, hex(imprint));
        oracle.printf("%s.tsa=%s%n", name, token.getTimeStampInfo().getTsa());
        oracle.printf("%s.nonce=%s%n", name, token.getTimeStampInfo().getNonce());
        oracle.printf("%s.certificateCount=%d%n", name, token.getCertificates().getMatches(null).size());
        SignerInformation signerInformation = (SignerInformation) token.toCMSSignedData().getSignerInfos()
                .getSigners().iterator().next();
        oracle.printf("%s.signerDigestAlgOID=%s%n", name, signerInformation.getDigestAlgOID());
        oracle.printf("%s.signerEncryptionAlgOID=%s%n", name, signerInformation.getEncryptionAlgOID());
        AttributeTable signed = signerInformation.getSignedAttributes();
        List<String> oids = new ArrayList<>();
        ASN1EncodableVector vector = signed.toASN1EncodableVector();
        for (int index = 0; index < vector.size(); index++) {
            oids.add(Attribute.getInstance(vector.get(index)).getAttrType().getId());
        }
        oracle.printf("%s.signedAttributeOIDs=%s%n", name, String.join(",", oids));
        oracle.printf("%s.unsignedAttributes=%s%n", name, signerInformation.getUnsignedAttributes());
        oracle.printf("%s.signedAttributesDER=%s%n", name,
                hex(MessageDigest.getInstance("SHA-256").digest(signerInformation.getEncodedSignedAttributes())));
        oracle.printf("%s.tstInfoDER.sha256=%s%n", name, hex(MessageDigest.getInstance("SHA-256")
                .digest(token.getTimeStampInfo().toASN1Structure().getEncoded("DER"))));
        oracle.printf("%s.base64=%s%n", name, Base64.getEncoder().encodeToString(encoded));
    }

    private static KeyPair deterministicRsa(String seed) throws Exception {
        Path cached = Path.of(seed + ".pk8");
        if (Files.exists(cached)) {
            return toKeyPair(Files.readAllBytes(cached));
        }
        SecureRandom random = SecureRandom.getInstance("SHA1PRNG");
        random.setSeed(seed.getBytes(StandardCharsets.UTF_8));
        KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA", "BC");
        generator.initialize(2048, random);
        return generator.generateKeyPair();
    }

    private static KeyPair toKeyPair(byte[] pkcs8) throws Exception {
        PrivateKey key = KeyFactory.getInstance("RSA", "BC").generatePrivate(new PKCS8EncodedKeySpec(pkcs8));
        java.security.interfaces.RSAPrivateCrtKey crt = (java.security.interfaces.RSAPrivateCrtKey) key;
        return new KeyPair(KeyFactory.getInstance("RSA", "BC").generatePublic(
                new java.security.spec.RSAPublicKeySpec(crt.getModulus(), crt.getPublicExponent())), key);
    }

    private static X509Certificate issue(KeyPair issuerKeyPair, KeyPair subjectKeyPair, String issuer,
                                         String subject, BigInteger serial, boolean ca, boolean tsa)
            throws Exception {
        JcaX509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(new X500Name(issuer), serial,
                NOT_BEFORE, NOT_AFTER, new X500Name(subject), subjectKeyPair.getPublic());
        builder.addExtension(Extension.basicConstraints, true, new BasicConstraints(ca));
        if (ca) {
            builder.addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.keyCertSign | KeyUsage.cRLSign));
        } else {
            builder.addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.nonRepudiation));
        }
        if (tsa) {
            builder.addExtension(Extension.extendedKeyUsage, true,
                    new ExtendedKeyUsage(KeyPurposeId.id_kp_timeStamping));
        }
        ContentSigner signer = new JcaContentSignerBuilder("SHA256withRSA").build(issuerKeyPair.getPrivate());
        return new JcaX509CertificateConverter().setProvider("BC").getCertificate(builder.build(signer));
    }

    private static String oidOf(String digestName) {
        switch (digestName) {
        case "SHA-224": return "2.16.840.1.101.3.4.2.4";
        case "SHA-256": return "2.16.840.1.101.3.4.2.1";
        case "SHA-384": return "2.16.840.1.101.3.4.2.2";
        case "SHA-512": return "2.16.840.1.101.3.4.2.3";
        default: throw new IllegalArgumentException(digestName);
        }
    }

    private static String hex(byte[] value) {
        StringBuilder builder = new StringBuilder();
        for (byte b : value) {
            builder.append(String.format("%02x", b));
        }
        return builder.toString();
    }
}
