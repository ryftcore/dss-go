// Produces the known-answer fixtures the CMSSRC chunk tests read:
//
//   cades-full.p7s        a CAdES-like CMS SignedData carrying, in one document, everything
//                         CMSCertificateSource / CMSCRLSource / CMSOCSPSource extract:
//                         SignedData.certificates and SignedData.crls (a CertificateList plus
//                         an OtherRevocationInfoFormat under id-ri-ocsp-response and another
//                         under id-pkix-ocsp-basic), the signingCertificate and
//                         signingCertificateV2 signed attributes, and the certValues,
//                         certificateRefs, attrCertificateRefs, revocationValues,
//                         revocationRefs and attrRevocationRefs unsigned attributes;
//   kat_cms_sources.txt   what the REAL DSS 6.5.RC1 sources extract from it - this generator
//                         instantiates eu.europa.esig.dss.spi.x509.CMSCertificateSource,
//                         CMSCRLSource and CMSOCSPSource and dumps their state, so the Go
//                         tests compare against upstream itself rather than against a
//                         hand-written expectation.
//
// Every collection is sorted before it is written, because the Java sources keep their state
// in HashMaps whose iteration order is unspecified.
//
// Run with (BouncyCastle 1.84 + DSS 6.5.RC1 from the local maven repository):
//
//   M2=$HOME/.m2/repository
//   CP=$M2/org/bouncycastle/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar
//   CP=$CP:$M2/org/bouncycastle/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar
//   CP=$CP:$M2/org/bouncycastle/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-spi/6.5.RC1/dss-spi-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-model/6.5.RC1/dss-model-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-enumerations/6.5.RC1/dss-enumerations-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-alert/6.5.RC1/dss-alert-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-utils/6.5.RC1/dss-utils-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-utils-apache-commons/6.5.RC1/dss-utils-apache-commons-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-crl-parser/6.5.RC1/dss-crl-parser-6.5.RC1.jar
//   CP=$CP:$M2/eu/europa/ec/joinup/sd-dss/dss-crl-parser-x509crl/6.5.RC1/dss-crl-parser-x509crl-6.5.RC1.jar
//   CP=$CP:$M2/org/apache/commons/commons-lang3/3.20.0/commons-lang3-3.20.0.jar
//   CP=$CP:$M2/org/apache/commons/commons-collections4/4.5.0/commons-collections4-4.5.0.jar
//   CP=$CP:$M2/commons-codec/commons-codec/1.18.0/commons-codec-1.18.0.jar
//   CP=$CP:$M2/commons-io/commons-io/2.22.0/commons-io-2.22.0.jar
//   CP=$CP:$M2/org/slf4j/slf4j-api/2.0.18/slf4j-api-2.0.18.jar
//   javac -cp "$CP" -d /tmp/cmssrckat CmsSourceKatGen.java
//   java -cp "$CP:/tmp/cmssrckat" CmsSourceKatGen ..
//
// The key material is generated on the fly from a fixed seed, so a re-run reproduces the same
// bytes. Nothing here is a secret.

import eu.europa.esig.dss.enumerations.CertificateOrigin;
import eu.europa.esig.dss.enumerations.CertificateRefOrigin;
import eu.europa.esig.dss.enumerations.RevocationOrigin;
import eu.europa.esig.dss.enumerations.RevocationRefOrigin;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.x509.CMSCRLSource;
import eu.europa.esig.dss.spi.x509.CMSCertificateSource;
import eu.europa.esig.dss.spi.x509.CMSOCSPSource;
import eu.europa.esig.dss.spi.x509.CertificateRef;
import eu.europa.esig.dss.spi.x509.SignerIdentifier;
import eu.europa.esig.dss.spi.x509.revocation.crl.CRLRef;
import eu.europa.esig.dss.spi.x509.revocation.ocsp.OCSPRef;
import eu.europa.esig.dss.spi.x509.revocation.ocsp.OCSPResponseBinary;
import eu.europa.esig.dss.model.identifier.EncapsulatedRevocationTokenIdentifier;
import eu.europa.esig.dss.model.x509.revocation.crl.CRL;
import eu.europa.esig.dss.model.x509.revocation.ocsp.OCSP;
import eu.europa.esig.dss.spi.x509.revocation.RevocationRef;
import org.bouncycastle.asn1.ASN1Encodable;
import org.bouncycastle.asn1.ASN1EncodableVector;
import org.bouncycastle.asn1.ASN1Encoding;
import org.bouncycastle.asn1.ASN1ObjectIdentifier;
import org.bouncycastle.asn1.DEROctetString;
import org.bouncycastle.asn1.DERSequence;
import org.bouncycastle.asn1.DERTaggedObject;
import org.bouncycastle.asn1.cms.Attribute;
import org.bouncycastle.asn1.cms.AttributeTable;
import org.bouncycastle.asn1.cms.CMSObjectIdentifiers;
import org.bouncycastle.asn1.esf.CrlIdentifier;
import org.bouncycastle.asn1.esf.CrlListID;
import org.bouncycastle.asn1.esf.CrlOcspRef;
import org.bouncycastle.asn1.esf.CrlValidatedID;
import org.bouncycastle.asn1.esf.OcspIdentifier;
import org.bouncycastle.asn1.esf.OcspListID;
import org.bouncycastle.asn1.esf.OcspResponsesID;
import org.bouncycastle.asn1.esf.OtherHash;
import org.bouncycastle.asn1.esf.OtherHashAlgAndValue;
import org.bouncycastle.asn1.ess.ESSCertID;
import org.bouncycastle.asn1.ess.ESSCertIDv2;
import org.bouncycastle.asn1.ess.OtherCertID;
import org.bouncycastle.asn1.ess.SigningCertificate;
import org.bouncycastle.asn1.ess.SigningCertificateV2;
import org.bouncycastle.asn1.nist.NISTObjectIdentifiers;
import org.bouncycastle.asn1.ocsp.BasicOCSPResponse;
import org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers;
import org.bouncycastle.asn1.oiw.OIWObjectIdentifiers;
import org.bouncycastle.asn1.pkcs.PKCSObjectIdentifiers;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.GeneralName;
import org.bouncycastle.asn1.x509.GeneralNames;
import org.bouncycastle.asn1.x509.IssuerSerial;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.X509CRLHolder;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v2CRLBuilder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
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
import org.bouncycastle.cms.DefaultSignedAttributeTableGenerator;
import org.bouncycastle.cms.SignerInformation;
import org.bouncycastle.cms.SignerInformationStore;
import org.bouncycastle.cms.jcajce.JcaSignerInfoGeneratorBuilder;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.DigestCalculator;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.operator.jcajce.JcaDigestCalculatorProviderBuilder;
import org.bouncycastle.util.CollectionStore;
import org.bouncycastle.util.encoders.Hex;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.MessageDigest;
import java.security.SecureRandom;
import java.security.Security;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.Date;
import java.util.List;
import java.util.Set;
import java.util.TreeSet;

public final class CmsSourceKatGen {

    private static final byte[] CONTENT = "Hello DSS Go port.\n".getBytes(StandardCharsets.UTF_8);
    // A fixed instant: 2026-08-12T12:23:42Z, so the fixture is reproducible.
    private static final Date NOW = new Date(1786537422000L);
    private static final Date NOT_BEFORE = new Date(NOW.getTime() - 86400000L);
    private static final Date NOT_AFTER = new Date(NOW.getTime() + 86400000L * 3650);

    private static final ASN1ObjectIdentifier ATTRIBUTE_CERTIFICATE_REFS =
            new ASN1ObjectIdentifier("1.2.840.113549.1.9.16.2.44");
    private static final ASN1ObjectIdentifier ATTRIBUTE_REVOCATION_REFS =
            new ASN1ObjectIdentifier("1.2.840.113549.1.9.16.2.45");

    private static final StringBuilder OUT = new StringBuilder();

    private static void put(String key, String value) {
        OUT.append(key).append('|').append(value).append('\n');
    }

    private static void putList(String key, Collection<String> values) {
        put(key, String.join(",", new TreeSet<>(values)));
    }

    public static void main(String[] args) throws Exception {
        Security.addProvider(new BouncyCastleProvider());
        Path out = Path.of(args.length > 0 ? args[0] : ".");

        SecureRandom random = SecureRandom.getInstance("SHA1PRNG");
        random.setSeed("dss-go-port-cmssrc".getBytes(StandardCharsets.UTF_8));

        KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA", "BC");
        generator.initialize(2048, random);

        KeyPair caKeyPair = generator.generateKeyPair();
        KeyPair signerKeyPair = generator.generateKeyPair();
        KeyPair otherKeyPair = generator.generateKeyPair();

        X500Name caName = new X500Name("C=LU,O=DSS Go Port,CN=CMS Source Test CA");
        X500Name signerName = new X500Name("C=LU,O=DSS Go Port,CN=CMS Source Test Signer");
        X500Name otherName = new X500Name("C=LU,O=DSS Go Port,CN=CMS Source Test Other");

        X509CertificateHolder caCert = issue(caName, caKeyPair, caName, caKeyPair, BigInteger.valueOf(1), true);
        X509CertificateHolder signerCert =
                issue(signerName, signerKeyPair, caName, caKeyPair, BigInteger.valueOf(0x5167), false);
        X509CertificateHolder otherCert =
                issue(otherName, otherKeyPair, caName, caKeyPair, BigInteger.valueOf(0x0777), false);

        X509CRLHolder crl = issueCRL(caName, caKeyPair);
        BasicOCSPResp basicOCSPResp = issueOCSP(caCert, signerCert, caKeyPair);
        OCSPResp ocspResp = new OCSPRespBuilder().build(OCSPRespBuilder.SUCCESSFUL, basicOCSPResp);
        // A second, different response, so that the id-pkix-ocsp-basic store is not a
        // duplicate of the id-ri-ocsp-response one and both code paths are observable.
        BasicOCSPResp otherBasicOCSPResp = issueOCSP(caCert, otherCert, caKeyPair);

        // ------------------------------------------------------------------ signed attributes
        byte[] signerSha1 = digest("SHA-1", signerCert.getEncoded());
        byte[] signerSha256 = digest("SHA-256", signerCert.getEncoded());
        IssuerSerial signerIssuerSerial = new IssuerSerial(
                new GeneralNames(new GeneralName(caName)), signerCert.getSerialNumber());

        SigningCertificate signingCertificate =
                new SigningCertificate(new ESSCertID(signerSha1, signerIssuerSerial));
        SigningCertificateV2 signingCertificateV2 = new SigningCertificateV2(
                new ESSCertIDv2[] { new ESSCertIDv2(
                        new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256), signerSha256, signerIssuerSerial) });

        ASN1EncodableVector signedAttributes = new ASN1EncodableVector();
        signedAttributes.add(new Attribute(PKCSObjectIdentifiers.id_aa_signingCertificate,
                new org.bouncycastle.asn1.DERSet(signingCertificate)));
        signedAttributes.add(new Attribute(PKCSObjectIdentifiers.id_aa_signingCertificateV2,
                new org.bouncycastle.asn1.DERSet(signingCertificateV2)));

        // ------------------------------------------------------------------ the CMS document
        CMSSignedDataGenerator cmsGenerator = new CMSSignedDataGenerator();
        DigestCalculator sha256Calculator =
                new JcaDigestCalculatorProviderBuilder().setProvider("BC").build().get(
                        new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256));
        ContentSigner contentSigner =
                new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(signerKeyPair.getPrivate());
        cmsGenerator.addSignerInfoGenerator(new JcaSignerInfoGeneratorBuilder(
                new JcaDigestCalculatorProviderBuilder().setProvider("BC").build())
                .setSignedAttributeGenerator(new DefaultSignedAttributeTableGenerator(new AttributeTable(signedAttributes)))
                .build(contentSigner, signerCert));

        List<X509CertificateHolder> certificates = new ArrayList<>(Arrays.asList(signerCert, caCert));
        cmsGenerator.addCertificates(new CollectionStore<>(certificates));
        cmsGenerator.addCRLs(new CollectionStore<>(List.of(crl)));
        cmsGenerator.addOtherRevocationInfo(CMSObjectIdentifiers.id_ri_ocsp_response,
                ocspResp.toASN1Structure());
        cmsGenerator.addOtherRevocationInfo(OCSPObjectIdentifiers.id_pkix_ocsp_basic,
                BasicOCSPResponse.getInstance(otherBasicOCSPResp.getEncoded()));

        CMSSignedData cms = cmsGenerator.generate(new CMSProcessableByteArray(CONTENT), true);

        // ------------------------------------------------------------------ unsigned attributes
        byte[] otherSha256 = digest("SHA-256", otherCert.getEncoded());
        byte[] caSha1 = digest("SHA-1", caCert.getEncoded());
        IssuerSerial otherIssuerSerial = new IssuerSerial(
                new GeneralNames(new GeneralName(caName)), otherCert.getSerialNumber());
        IssuerSerial caIssuerSerial = new IssuerSerial(
                new GeneralNames(new GeneralName(caName)), caCert.getSerialNumber());

        // certValues: SEQUENCE OF Certificate
        Attribute certValues = new Attribute(PKCSObjectIdentifiers.id_aa_ets_certValues,
                new org.bouncycastle.asn1.DERSet(new DERSequence(otherCert.toASN1Structure())));

        // certificateRefs: SEQUENCE OF OtherCertID, with the OtherHashAlgAndValue alternative
        OtherCertID otherCertIDSha256 = new OtherCertID(
                new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256), otherSha256, otherIssuerSerial);
        Attribute certificateRefs = new Attribute(PKCSObjectIdentifiers.id_aa_ets_certificateRefs,
                new org.bouncycastle.asn1.DERSet(new DERSequence(otherCertIDSha256)));

        // attrCertificateRefs: SEQUENCE OF OtherCertID, with the bare sha1Hash alternative
        // BouncyCastle has no constructor for the bare sha1Hash alternative of OtherHash, so
        // the OtherCertID is assembled by hand: SEQUENCE { OCTET STRING, IssuerSerial }.
        DERSequence otherCertIDSha1 = new DERSequence(new ASN1Encodable[] {
                new DEROctetString(caSha1), caIssuerSerial });
        Attribute attrCertificateRefs = new Attribute(ATTRIBUTE_CERTIFICATE_REFS,
                new org.bouncycastle.asn1.DERSet(new DERSequence(otherCertIDSha1)));

        // revocationValues: crlVals + ocspVals
        ASN1EncodableVector crlValsVector = new ASN1EncodableVector();
        crlValsVector.add(crl.toASN1Structure());
        ASN1EncodableVector ocspValsVector = new ASN1EncodableVector();
        ocspValsVector.add(BasicOCSPResponse.getInstance(basicOCSPResp.getEncoded()));
        ASN1EncodableVector revocationValuesVector = new ASN1EncodableVector();
        revocationValuesVector.add(new DERTaggedObject(true, 0, new DERSequence(crlValsVector)));
        revocationValuesVector.add(new DERTaggedObject(true, 1, new DERSequence(ocspValsVector)));
        Attribute revocationValues = new Attribute(PKCSObjectIdentifiers.id_aa_ets_revocationValues,
                new org.bouncycastle.asn1.DERSet(new DERSequence(revocationValuesVector)));

        // revocationRefs: a CrlOcspRef holding both a CRLListID and an OcspListID
        byte[] crlSha256 = digest("SHA-256", crl.getEncoded());
        CrlValidatedID crlValidatedID = new CrlValidatedID(
                new OtherHash(new OtherHashAlgAndValue(new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256),
                        new DEROctetString(crlSha256))),
                new CrlIdentifier(caName, new org.bouncycastle.asn1.ASN1UTCTime(NOW), BigInteger.valueOf(3)));
        byte[] ocspSha256 = digest("SHA-256", basicOCSPResp.getEncoded());
        OcspResponsesID ocspResponsesID = new OcspResponsesID(
                new OcspIdentifier(new RespID(caName).toASN1Primitive(),
                        new org.bouncycastle.asn1.ASN1GeneralizedTime(basicOCSPResp.getProducedAt())),
                new OtherHash(new OtherHashAlgAndValue(new AlgorithmIdentifier(NISTObjectIdentifiers.id_sha256),
                        new DEROctetString(ocspSha256))));
        CrlOcspRef bothRef = new CrlOcspRef(new CrlListID(new CrlValidatedID[] { crlValidatedID }),
                new OcspListID(new OcspResponsesID[] { ocspResponsesID }), null);
        Attribute revocationRefs = new Attribute(PKCSObjectIdentifiers.id_aa_ets_revocationRefs,
                new org.bouncycastle.asn1.DERSet(new DERSequence(bothRef)));

        // attrRevocationRefs: a CrlOcspRef with a CRLListID only, whose CrlValidatedID has no
        // CrlIdentifier and uses the bare sha1Hash alternative.
        CrlValidatedID hashOnly = new CrlValidatedID(new OtherHash(digest("SHA-1", crl.getEncoded())));
        CrlOcspRef crlOnlyRef = new CrlOcspRef(new CrlListID(new CrlValidatedID[] { hashOnly }), null, null);
        Attribute attrRevocationRefs = new Attribute(ATTRIBUTE_REVOCATION_REFS,
                new org.bouncycastle.asn1.DERSet(new DERSequence(crlOnlyRef)));

        ASN1EncodableVector unsignedVector = new ASN1EncodableVector();
        unsignedVector.add(certValues);
        unsignedVector.add(certificateRefs);
        unsignedVector.add(attrCertificateRefs);
        unsignedVector.add(revocationValues);
        unsignedVector.add(revocationRefs);
        unsignedVector.add(attrRevocationRefs);
        AttributeTable unsignedAttributes = new AttributeTable(unsignedVector);

        SignerInformation signerInformation = cms.getSignerInfos().getSigners().iterator().next();
        SignerInformation withUnsigned =
                SignerInformation.replaceUnsignedAttributes(signerInformation, unsignedAttributes);
        cms = CMSSignedData.replaceSigners(cms, new SignerInformationStore(withUnsigned));

        byte[] encoded = cms.toASN1Structure().getEncoded(ASN1Encoding.DER);
        Files.write(out.resolve("cades-full.p7s"), encoded);

        // ------------------------------------------------------------------ the oracle
        CMSSignedData parsed = new CMSSignedData(encoded);
        SignerInformation signer = parsed.getSignerInfos().getSigners().iterator().next();

        dumpCertificateSource(new CMSCertificateSource(parsed.getSignerInfos(), parsed.getCertificates(), signer) {});
        dumpCRLSource(new CMSCRLSource(parsed.getCRLs(), signer.getUnsignedAttributes()) {});
        dumpOCSPSource(new CMSOCSPSource(parsed.getOtherRevocationInfo(CMSObjectIdentifiers.id_ri_ocsp_response),
                parsed.getOtherRevocationInfo(OCSPObjectIdentifiers.id_pkix_ocsp_basic),
                signer.getUnsignedAttributes()) {});

        Files.writeString(out.resolve("kat_cms_sources.txt"), OUT.toString());
    }

    // ---------------------------------------------------------------------- oracle dumping

    private static void dumpCertificateSource(CMSCertificateSource source) {
        putList("cert.all", ids(source.getCertificates()));
        putList("cert.signedData", ids(source.getSignedDataCertificates()));
        putList("cert.certificateValues", ids(source.getCertificateValues()));

        List<String> identifiers = new ArrayList<>();
        for (SignerIdentifier identifier : source.getAllCertificateIdentifiers()) {
            identifiers.add(signerIdentifier(identifier));
        }
        putList("cert.identifiers", identifiers);
        put("cert.currentIdentifier", signerIdentifier(source.getCurrentCertificateIdentifier()));

        putList("ref.signingCertificate", certRefs(source.getSigningCertificateRefs()));
        putList("ref.completeCertificateRefs", certRefs(source.getCompleteCertificateRefs()));
        putList("ref.attributeCertificateRefs", certRefs(source.getAttributeCertificateRefs()));
    }

    private static void dumpCRLSource(CMSCRLSource source) {
        putList("crl.cmsSignedData", binaries(source.getCMSSignedDataRevocationBinaries()));
        putList("crl.revocationValues", binaries(source.getRevocationValuesBinaries()));
        List<String> complete = new ArrayList<>();
        for (RevocationRef<CRL> ref : source.getCompleteRevocationRefs()) {
            complete.add(crlRef((CRLRef) ref));
        }
        putList("crl.completeRefs", complete);
        List<String> attribute = new ArrayList<>();
        for (RevocationRef<CRL> ref : source.getAttributeRevocationRefs()) {
            attribute.add(crlRef((CRLRef) ref));
        }
        putList("crl.attributeRefs", attribute);
    }

    private static void dumpOCSPSource(CMSOCSPSource source) {
        List<String> fromSignedData = new ArrayList<>();
        for (EncapsulatedRevocationTokenIdentifier<?> binary : source.getCMSSignedDataRevocationBinaries()) {
            OCSPResponseBinary ocspBinary = (OCSPResponseBinary) binary;
            fromSignedData.add(ocspBinary.asXmlId() + ":" + ocspBinary.getAsn1ObjectIdentifier().getId());
        }
        putList("ocsp.cmsSignedData", fromSignedData);
        putList("ocsp.revocationValues", binaries(source.getRevocationValuesBinaries()));
        List<String> complete = new ArrayList<>();
        for (RevocationRef<OCSP> ref : source.getCompleteRevocationRefs()) {
            complete.add(ocspRef((OCSPRef) ref));
        }
        putList("ocsp.completeRefs", complete);
        List<String> attribute = new ArrayList<>();
        for (RevocationRef<OCSP> ref : source.getAttributeRevocationRefs()) {
            attribute.add(ocspRef((OCSPRef) ref));
        }
        putList("ocsp.attributeRefs", attribute);
    }

    private static List<String> ids(Collection<CertificateToken> tokens) {
        List<String> values = new ArrayList<>();
        for (CertificateToken token : tokens) {
            values.add(token.getDSSIdAsString());
        }
        return values;
    }

    private static List<String> binaries(Collection<? extends EncapsulatedRevocationTokenIdentifier<?>> binaries) {
        List<String> values = new ArrayList<>();
        for (EncapsulatedRevocationTokenIdentifier<?> binary : binaries) {
            values.add(binary.asXmlId());
        }
        return values;
    }

    private static List<String> certRefs(Collection<CertificateRef> refs) {
        List<String> values = new ArrayList<>();
        for (CertificateRef ref : refs) {
            StringBuilder builder = new StringBuilder();
            builder.append(ref.getDSSIdAsString()).append(':');
            if (ref.getCertDigest() != null) {
                builder.append(ref.getCertDigest().getAlgorithm().name()).append(':')
                        .append(Hex.toHexString(ref.getCertDigest().getValue()));
            }
            builder.append(':').append(signerIdentifier(ref.getCertificateIdentifier()));
            values.add(builder.toString());
        }
        return values;
    }

    private static String signerIdentifier(SignerIdentifier identifier) {
        if (identifier == null) {
            return "<null>";
        }
        return (identifier.getIssuerName() == null ? "" : identifier.getIssuerName().getName())
                + ";" + (identifier.getSerialNumber() == null ? "" : identifier.getSerialNumber().toString())
                + ";" + (identifier.getSki() == null ? "" : Hex.toHexString(identifier.getSki()))
                + ";" + identifier.isCurrent();
    }

    private static String crlRef(CRLRef ref) {
        return ref.getDSSIdAsString()
                + ":" + ref.getDigest().getAlgorithm().name()
                + ":" + Hex.toHexString(ref.getDigest().getValue())
                + ":" + (ref.getCrlIssuer() == null ? "" : ref.getCrlIssuer().getName())
                + ":" + (ref.getCrlIssueTime() == null ? "" : Long.toString(ref.getCrlIssueTime().getTime()))
                + ":" + (ref.getCrlNumber() == null ? "" : ref.getCrlNumber().toString());
    }

    private static String ocspRef(OCSPRef ref) {
        return ref.getDSSIdAsString()
                + ":" + ref.getDigest().getAlgorithm().name()
                + ":" + Hex.toHexString(ref.getDigest().getValue())
                + ":" + (ref.getProducedAt() == null ? "" : Long.toString(ref.getProducedAt().getTime()))
                + ":" + (ref.getResponderId() == null ? "" : (ref.getResponderId().getX500Principal() == null ? "" : ref.getResponderId().getX500Principal().getName()));
    }

    // ---------------------------------------------------------------------- material

    private static X509CertificateHolder issue(X500Name subject, KeyPair subjectKeyPair, X500Name issuer,
            KeyPair issuerKeyPair, BigInteger serial, boolean ca) throws Exception {
        X509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(
                issuer, serial, NOT_BEFORE, NOT_AFTER, subject, subjectKeyPair.getPublic());
        builder.addExtension(Extension.basicConstraints, true, new BasicConstraints(ca));
        builder.addExtension(Extension.keyUsage, true, new KeyUsage(
                ca ? KeyUsage.keyCertSign | KeyUsage.cRLSign : KeyUsage.digitalSignature | KeyUsage.nonRepudiation));
        ContentSigner signer =
                new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(issuerKeyPair.getPrivate());
        return builder.build(signer);
    }

    private static X509CRLHolder issueCRL(X500Name issuer, KeyPair issuerKeyPair) throws Exception {
        X509v2CRLBuilder builder = new X509v2CRLBuilder(issuer, NOW);
        builder.setNextUpdate(NOT_AFTER);
        builder.addCRLEntry(BigInteger.valueOf(0x0999), NOW, 0);
        ContentSigner signer =
                new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(issuerKeyPair.getPrivate());
        return builder.build(signer);
    }

    private static BasicOCSPResp issueOCSP(X509CertificateHolder caCert, X509CertificateHolder signerCert,
            KeyPair caKeyPair) throws Exception {
        DigestCalculator sha1 = new JcaDigestCalculatorProviderBuilder().setProvider("BC").build()
                .get(new AlgorithmIdentifier(OIWObjectIdentifiers.idSHA1));
        CertificateID certificateID = new CertificateID(sha1, caCert, signerCert.getSerialNumber());
        BasicOCSPRespBuilder builder = new BasicOCSPRespBuilder(new RespID(caCert.getSubject()));
        builder.addResponse(certificateID, CertificateStatus.GOOD, NOW, null, null);
        ContentSigner signer =
                new JcaContentSignerBuilder("SHA256withRSA").setProvider("BC").build(caKeyPair.getPrivate());
        return builder.build(signer, new X509CertificateHolder[] { caCert }, NOW);
    }

    private static byte[] digest(String algorithm, byte[] data) throws Exception {
        return MessageDigest.getInstance(algorithm).digest(data);
    }
}
