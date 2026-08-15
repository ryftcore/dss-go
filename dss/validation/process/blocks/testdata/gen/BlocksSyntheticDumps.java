/*
 * Builds the synthetic XmlDiagnosticData dump the BasicBuildingBlocks oracle needs
 * on top of the marshal-parity corpus and XCVA's dd/ dumps.
 *
 * BasicBuildingBlocks#addAdditionalInfo(XmlXCV) hangs the cross-certificate and
 * equivalent-certificate id lists off every XmlSubXCV. Neither the 50-dump
 * marshal-parity corpus nor XCVA's dd/ dumps contain two used certificates sharing
 * an EntityKey, so getCrossCertificates()/getEquivalentCertificates() answer empty
 * for all of them and that whole method is dead weight in the corpus - an audit
 * mutation dropping its `equivalentCertificates.removeAll(crossCertificates)` step
 * survived. This dump supplies the missing shape.
 *
 * It is written once, with the upstream DiagnosticDataFacade, to
 *   dss/validation/process/blocks/testdata/dd/*.xml
 * and committed. Both BlocksOracle and the Go test read those very files, so
 * neither side gets a private fixture.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle BlocksSyntheticDumps.java
 *   java  -cp "$CP:/tmp/oracle" BlocksSyntheticDumps <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlBasicSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDistinguishedName;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrusted;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.EncryptionAlgorithm;

import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.time.Instant;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

public class BlocksSyntheticDumps {

    public static void main(String[] args) throws Exception {
        File repoRoot = new File(args[0]);
        File outDir = new File(repoRoot, "validation/process/blocks/testdata/dd");
        outDir.mkdirs();
        write(outDir, "cross-equivalent.xml", crossEquivalent());
    }

    private static void write(File dir, String name, XmlDiagnosticData data) throws Exception {
        try (OutputStream out = new FileOutputStream(new File(dir, name))) {
            DiagnosticDataFacade.newFacade().marshall(data, out, false);
        }
    }

    /**
     * A trusted root, a signing certificate, and two further used certificates that
     * share the signing certificate's EntityKey:
     *
     *  - C-CROSS carries a different subject DN, so getCrossCertificates() returns it
     *    (and getEquivalentCertificates() does too - which is exactly what
     *    addAdditionalInfo()'s removeAll() step has to subtract again);
     *  - C-EQUIV carries the very same subject and issuer DN, so it is equivalent but
     *    not cross.
     */
    private static XmlDiagnosticData crossEquivalent() {
        XmlCertificate root = certificate("C-ROOT", "CN=Root", "CN=Root", "C-ROOT-KEY");
        XmlTrusted trusted = new XmlTrusted();
        trusted.setValue(true);
        root.setTrusted(trusted);
        root.setSelfSigned(true);

        XmlCertificate signer = certificate("C-SIGNER", "CN=Signer", "CN=Root", "C-SHARED-KEY");
        chain(signer, root);

        XmlCertificate cross = certificate("C-CROSS", "CN=Signer Cross", "CN=Other Root", "C-SHARED-KEY");
        XmlCertificate equivalent = certificate("C-EQUIV", "CN=Signer", "CN=Root", "C-SHARED-KEY");

        return data("cross-equivalent", signer, Arrays.asList(signer, cross, equivalent, root));
    }

    // ----------------------------------------------------------------- helpers

    private static XmlCertificate certificate(String id, String subjectDN, String issuerDN, String entityKey) {
        XmlCertificate certificate = new XmlCertificate();
        certificate.setId(id);
        certificate.setNotBefore(date("2020-01-01T00:00:00Z"));
        certificate.setNotAfter(date("2028-01-01T00:00:00Z"));
        certificate.setSelfSigned(false);
        certificate.setPublicKeySize(2048);
        certificate.setPublicKeyEncryptionAlgo(EncryptionAlgorithm.RSA);
        certificate.setEntityKey(entityKey);
        certificate.getSubjectDistinguishedName().add(distinguishedName(subjectDN));
        certificate.getIssuerDistinguishedName().add(distinguishedName(issuerDN));
        XmlBasicSignature basicSignature = new XmlBasicSignature();
        basicSignature.setSignatureIntact(true);
        basicSignature.setSignatureValid(true);
        basicSignature.setDigestAlgoUsedToSignThisToken(DigestAlgorithm.SHA256);
        basicSignature.setEncryptionAlgoUsedToSignThisToken(EncryptionAlgorithm.RSA);
        basicSignature.setKeyLengthUsedToSignThisToken("2048");
        certificate.setBasicSignature(basicSignature);
        return certificate;
    }

    private static XmlDistinguishedName distinguishedName(String value) {
        XmlDistinguishedName name = new XmlDistinguishedName();
        name.setFormat("RFC2253");
        name.setValue(value);
        return name;
    }

    private static void chain(XmlCertificate certificate, XmlCertificate... parents) {
        for (XmlCertificate parent : parents) {
            XmlChainItem item = new XmlChainItem();
            item.setCertificate(parent);
            certificate.getCertificateChain().add(item);
        }
        XmlSigningCertificate signing = new XmlSigningCertificate();
        signing.setCertificate(parents[0]);
        certificate.setSigningCertificate(signing);
    }

    private static XmlDiagnosticData data(String name, XmlCertificate signingCertificate,
                                          List<XmlCertificate> usedCertificates) {
        XmlDiagnosticData data = new XmlDiagnosticData();
        data.setDocumentName(name);
        data.setValidationDate(date("2024-01-01T00:00:00Z"));

        XmlSignature signature = new XmlSignature();
        signature.setId("S-1");
        signature.setSignatureFilename(name);
        XmlSigningCertificate signing = new XmlSigningCertificate();
        signing.setCertificate(signingCertificate);
        signature.setSigningCertificate(signing);
        for (XmlChainItem item : signingCertificate.getCertificateChain()) {
            XmlChainItem copy = new XmlChainItem();
            copy.setCertificate(item.getCertificate());
            signature.getCertificateChain().add(copy);
        }
        XmlBasicSignature basicSignature = new XmlBasicSignature();
        basicSignature.setSignatureIntact(true);
        basicSignature.setSignatureValid(true);
        basicSignature.setDigestAlgoUsedToSignThisToken(DigestAlgorithm.SHA256);
        basicSignature.setEncryptionAlgoUsedToSignThisToken(EncryptionAlgorithm.RSA);
        basicSignature.setKeyLengthUsedToSignThisToken("2048");
        signature.setBasicSignature(basicSignature);
        data.getSignatures().add(signature);

        data.getUsedCertificates().addAll(usedCertificates);
        return data;
    }

    private static Date date(String instant) {
        return Date.from(Instant.parse(instant));
    }
}
