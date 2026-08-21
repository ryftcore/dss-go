// Produces the two cmscore fixtures OpenSSL 3.0.13 cannot write:
//
//   ocsp-crl.p7s       a SignedData whose RevocationInfoChoices holds a CertificateList AND an
//                      OtherRevocationInfoFormat carrying an OCSPResponse (id-ri-ocsp-response),
//                      plus a second one carrying a BasicOCSPResponse (id-pkix-ocsp-basic);
//   ed25519-attached.p7s  an RFC 8419 Ed25519 signature (OpenSSL refuses: EdDSA CMS support
//                      landed in OpenSSL 3.2).
//
// Everything else in testdata/ comes from testdata/generate.sh. Run with:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/cmsfix CmsFixtures.java && java -cp "$CP:/tmp/cmsfix" CmsFixtures ..
//
// The key material is generated on the fly and thrown away; the fixtures are checked in, so a
// re-run produces different (equally valid) bytes rather than reproducing these files.

import org.bouncycastle.asn1.ASN1Encoding;
import org.bouncycastle.asn1.cms.CMSObjectIdentifiers;
import org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers;
import org.bouncycastle.asn1.oiw.OIWObjectIdentifiers;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.X509CRLHolder;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v2CRLBuilder;
import org.bouncycastle.cert.jcajce.JcaCRLStore;
import org.bouncycastle.cert.jcajce.JcaCertStore;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.cert.ocsp.BasicOCSPResp;
import org.bouncycastle.cert.ocsp.BasicOCSPRespBuilder;
import org.bouncycastle.cert.ocsp.CertificateID;
import org.bouncycastle.cert.ocsp.CertificateStatus;
import org.bouncycastle.cert.ocsp.OCSPResp;
import org.bouncycastle.cert.ocsp.OCSPRespBuilder;
import org.bouncycastle.cert.ocsp.RespID;
import org.bouncycastle.cms.CMSProcessableByteArray;
import org.bouncycastle.cms.CMSSignedData;
import org.bouncycastle.cms.CMSSignedDataGenerator;
import org.bouncycastle.cms.jcajce.JcaSignerInfoGeneratorBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.DigestCalculator;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.operator.jcajce.JcaDigestCalculatorProviderBuilder;
import org.bouncycastle.util.CollectionStore;

import java.io.OutputStream;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.Security;
import java.security.cert.X509CRL;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public final class CmsFixtures {

    private static final byte[] CONTENT = "Hello DSS Go port.\n".getBytes(StandardCharsets.UTF_8);

    public static void main(String[] args) throws Exception {
        Security.addProvider(new BouncyCastleProvider());
        Path out = Path.of(args.length > 0 ? args[0] : ".");

        writeOcspCrlFixture(out.resolve("ocsp-crl.p7s"));
        writeEd25519Fixture(out.resolve("ed25519-attached.p7s"));
    }

    private static void writeOcspCrlFixture(Path target) throws Exception {
        KeyPairGenerator rsa = KeyPairGenerator.getInstance("RSA", "BC");
        rsa.initialize(2048);
        KeyPair caPair = rsa.generateKeyPair();
        KeyPair signerPair = rsa.generateKeyPair();

        Date notBefore = new Date(1600000000000L);
        Date notAfter = new Date(2500000000000L);

        X500Name caName = new X500Name("C=LU,O=DSS Go Port,CN=BC Test CA");
        ContentSigner caSigner = new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(caPair.getPrivate());
        X509CertificateHolder caHolder = new JcaX509v3CertificateBuilder(
                caName, BigInteger.valueOf(1), notBefore, notAfter, caName, caPair.getPublic())
                .addExtension(Extension.basicConstraints, true, new BasicConstraints(0))
                .addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.keyCertSign | KeyUsage.cRLSign | KeyUsage.digitalSignature))
                .build(caSigner);

        X500Name signerName = new X500Name("C=LU,O=DSS Go Port,CN=BC Test Signer");
        X509CertificateHolder signerHolder = new JcaX509v3CertificateBuilder(
                caName, BigInteger.valueOf(0x1234), notBefore, notAfter, signerName, signerPair.getPublic())
                .addExtension(Extension.basicConstraints, true, new BasicConstraints(false))
                .addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.digitalSignature | KeyUsage.nonRepudiation))
                .build(caSigner);

        JcaX509CertificateConverter converter = new JcaX509CertificateConverter().setProvider("BC");
        X509Certificate caCertificate = converter.getCertificate(caHolder);
        X509Certificate signerCertificate = converter.getCertificate(signerHolder);

        // A CRL revoking an unrelated serial, so SignedData.crls holds a real CertificateList.
        X509v2CRLBuilder crlBuilder = new X509v2CRLBuilder(caName, new Date(1700000000000L));
        crlBuilder.setNextUpdate(new Date(2400000000000L));
        crlBuilder.addCRLEntry(BigInteger.valueOf(0xDEAD), new Date(1650000000000L), 1);
        X509CRLHolder crlHolder = crlBuilder.build(caSigner);
        X509CRL crl = new org.bouncycastle.cert.jcajce.JcaX509CRLConverter().setProvider("BC").getCRL(crlHolder);

        // An OCSP response saying the signer certificate is good, wrapped both as an
        // OCSPResponse (id-ri-ocsp-response) and as a bare BasicOCSPResponse (id-pkix-ocsp-basic).
        DigestCalculator sha1 = new JcaDigestCalculatorProviderBuilder().setProvider("BC").build()
                .get(new AlgorithmIdentifier(OIWObjectIdentifiers.idSHA1));
        CertificateID certificateID = new CertificateID(sha1, caHolder, signerHolder.getSerialNumber());
        BasicOCSPResp basicResponse = new BasicOCSPRespBuilder(new RespID(caName))
                .addResponse(certificateID, CertificateStatus.GOOD)
                .build(caSigner, new X509CertificateHolder[] { caHolder }, new Date(1700000001000L));
        OCSPResp ocspResponse = new OCSPRespBuilder().build(OCSPRespBuilder.SUCCESSFUL, basicResponse);

        List<X509Certificate> certificates = new ArrayList<>();
        certificates.add(signerCertificate);
        certificates.add(caCertificate);

        CMSSignedDataGenerator generator = new CMSSignedDataGenerator();
        generator.addSignerInfoGenerator(new JcaSignerInfoGeneratorBuilder(
                new JcaDigestCalculatorProviderBuilder().setProvider("BC").build())
                .build(new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(signerPair.getPrivate()),
                        signerCertificate));
        generator.addCertificates(new JcaCertStore(certificates));
        generator.addCRLs(new JcaCRLStore(List.of(crl)));
        generator.addOtherRevocationInfo(CMSObjectIdentifiers.id_ri_ocsp_response,
                new CollectionStore<>(List.of(ocspResponse.toASN1Structure())));
        generator.addOtherRevocationInfo(OCSPObjectIdentifiers.id_pkix_ocsp_basic,
                new CollectionStore<>(List.of(
                        org.bouncycastle.asn1.ocsp.BasicOCSPResponse.getInstance(basicResponse.getEncoded()))));

        CMSSignedData signedData = generator.generate(new CMSProcessableByteArray(CONTENT), true);
        write(target, signedData);
    }

    private static void writeEd25519Fixture(Path target) throws Exception {
        KeyPair pair = KeyPairGenerator.getInstance("Ed25519", "BC").generateKeyPair();
        X500Name name = new X500Name("C=LU,O=DSS Go Port,CN=BC Ed25519 Signer");
        ContentSigner signer = new JcaContentSignerBuilder("Ed25519").setProvider("BC").build(pair.getPrivate());
        X509CertificateHolder holder = new JcaX509v3CertificateBuilder(
                name, BigInteger.valueOf(0x25519), new Date(1600000000000L), new Date(2500000000000L), name, pair.getPublic())
                .addExtension(Extension.basicConstraints, true, new BasicConstraints(false))
                .addExtension(Extension.keyUsage, true, new KeyUsage(KeyUsage.digitalSignature | KeyUsage.nonRepudiation))
                .build(signer);
        X509Certificate certificate = new JcaX509CertificateConverter().setProvider("BC").getCertificate(holder);

        CMSSignedDataGenerator generator = new CMSSignedDataGenerator();
        generator.addSignerInfoGenerator(new JcaSignerInfoGeneratorBuilder(
                new JcaDigestCalculatorProviderBuilder().setProvider("BC").build()).build(signer, certificate));
        generator.addCertificates(new JcaCertStore(List.of(certificate)));

        CMSSignedData signedData = generator.generate(new CMSProcessableByteArray(CONTENT), true);
        write(target, signedData);
    }

    private static void write(Path target, CMSSignedData signedData) throws Exception {
        try (OutputStream stream = Files.newOutputStream(target)) {
            stream.write(signedData.toASN1Structure().getEncoded(ASN1Encoding.DER));
        }
        System.out.println("wrote " + target + " (" + Files.size(target) + " bytes)");
    }
}
