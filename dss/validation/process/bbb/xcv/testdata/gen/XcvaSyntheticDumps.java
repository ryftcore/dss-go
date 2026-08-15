/*
 * Builds the synthetic XmlDiagnosticData dumps the XCVA block oracle needs.
 *
 * The 50-dump marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle carries
 * no trusted certificate at all (grep Trusted: every certificate answers false), so
 * X509CertificateValidation stops on its very first check for every token of it and
 * none of the trust-anchor walk, the SubXCV loop, the validation-model date
 * progression or the SUB_XCV_TA message filtering is ever reached. These dumps
 * supply exactly those inputs.
 *
 * They are written once, with the upstream DiagnosticDataFacade, to
 *   dss/validation/process/bbb/xcv/testdata/dd/*.xml
 * and committed. Both XcvaOracle and the Go test read those very files, so neither
 * side gets a private fixture.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle XcvaSyntheticDumps.java
 *   java  -cp "$CP:/tmp/oracle" XcvaSyntheticDumps <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlBasicSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef;
import eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrusted;
import eu.europa.esig.dss.enumerations.CertificateRefOrigin;
import eu.europa.esig.dss.enumerations.CertificateStatus;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.EncryptionAlgorithm;
import eu.europa.esig.dss.enumerations.RevocationType;
import eu.europa.esig.dss.enumerations.RevocationOrigin;

import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

public class XcvaSyntheticDumps {

    public static void main(String[] args) throws Exception {
        File repoRoot = new File(args[0]);
        File outDir = new File(repoRoot, "validation/process/bbb/xcv/testdata/dd");
        outDir.mkdirs();

        write(outDir, "trusted-root.xml", trustedRoot());
        write(outDir, "trusted-signing-cert.xml", trustedSigningCert());
        write(outDir, "sunset-past.xml", sunsetPast());
        write(outDir, "sunset-future.xml", sunsetFuture());
        write(outDir, "trust-start-date.xml", trustStartDate());
        write(outDir, "untrusted-chain.xml", untrustedChain());
        write(outDir, "trusted-root-revocation.xml", trustedRootWithRevocation());
        write(outDir, "rac-inputs.xml", racInputs());
        write(outDir, "rac-passed.xml", racPassed());
        write(outDir, "crs-two-acceptable.xml", crsTwoAcceptable());
        write(outDir, "crs-mixed-acceptance.xml", crsMixedAcceptance());
        write(outDir, "rac-anchor-mid-chain.xml", racAnchorMidChain());
    }

    private static void write(File dir, String name, XmlDiagnosticData data) throws Exception {
        try (OutputStream out = new FileOutputStream(new File(dir, name))) {
            DiagnosticDataFacade.newFacade().marshall(data, out, false);
        }
    }

    // --------------------------------------------------------------- fixtures

    /** Signing certificate under an intermediate CA under a trusted root. */
    private static XmlDiagnosticData trustedRoot() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, null);
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);
        return data("trusted-root", signer, Arrays.asList(signer, ca, root));
    }

    /** The signing certificate is itself the trust anchor: the early-return path. */
    private static XmlDiagnosticData trustedSigningCert() {
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        trust(signer, null, null);
        return data("trusted-signing-cert", signer, Arrays.asList(signer));
    }

    /** Trusted root whose sunset date is before the validation time. */
    private static XmlDiagnosticData sunsetPast() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, "2020-06-01T00:00:00Z");
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);
        return data("sunset-past", signer, Arrays.asList(signer, ca, root));
    }

    /** Trusted root whose sunset date is after the validation time. */
    private static XmlDiagnosticData sunsetFuture() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, "2030-06-01T00:00:00Z");
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);
        return data("sunset-future", signer, Arrays.asList(signer, ca, root));
    }

    /** Trusted root carrying a trust start date but no sunset date. */
    private static XmlDiagnosticData trustStartDate() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, "2016-01-01T00:00:00Z", null);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        return data("trust-start-date", signer, Arrays.asList(signer, root));
    }

    /** No trust anchor at all: the SubXCV loop runs over the whole chain. */
    private static XmlDiagnosticData untrustedChain() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);
        return data("untrusted-chain", signer, Arrays.asList(signer, ca, root));
    }

    /** Trusted root, with an OCSP response over the signing certificate. */
    private static XmlDiagnosticData trustedRootWithRevocation() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, null);
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);

        XmlRevocation ocsp = new XmlRevocation();
        ocsp.setId("R-OCSP");
        ocsp.setOrigin(RevocationOrigin.EXTERNAL);
        ocsp.setType(RevocationType.OCSP);
        ocsp.setProductionDate(date("2023-01-01T00:00:00Z"));
        ocsp.setThisUpdate(date("2023-01-01T00:00:00Z"));
        ocsp.setCertHashExtensionPresent(true);
        ocsp.setCertHashExtensionMatch(true);
        XmlSigningCertificate ocspSigner = new XmlSigningCertificate();
        ocspSigner.setCertificate(ca);
        ocsp.setSigningCertificate(ocspSigner);
        XmlChainItem ocspChainCa = new XmlChainItem();
        ocspChainCa.setCertificate(ca);
        XmlChainItem ocspChainRoot = new XmlChainItem();
        ocspChainRoot.setCertificate(root);
        ocsp.getCertificateChain().add(ocspChainCa);
        ocsp.getCertificateChain().add(ocspChainRoot);
        XmlBasicSignature ocspSignature = new XmlBasicSignature();
        ocspSignature.setSignatureIntact(true);
        ocspSignature.setSignatureValid(true);
        ocsp.setBasicSignature(ocspSignature);

        XmlCertificateRevocation certificateRevocation = new XmlCertificateRevocation();
        certificateRevocation.setRevocation(ocsp);
        certificateRevocation.setStatus(CertificateStatus.GOOD);
        signer.getRevocations().add(certificateRevocation);

        XmlDiagnosticData data = data("trusted-root-revocation", signer, Arrays.asList(signer, ca, root));
        data.getUsedRevocations().add(ocsp);
        return data;
    }

    /**
     * The one shape that makes RevocationAcceptanceChecker conclude PASSED: an OCSP
     * signed by the trust anchor itself, so the revocation-chain loop breaks on its
     * first certificate, with every preceding check satisfied.
     */
    /**
     * A signing certificate carrying TWO acceptable (PASSED-RAC) OCSP responses with
     * different production dates, so that CertificateRevocationSelector's
     * "keep the latest acceptable one" comparison actually has to choose between two
     * candidates. Not one certificate of the marshal-parity corpus, nor of the other
     * dumps here, has more than one revocation whose RAC passes, so without this
     * fixture that comparison is never exercised (an audit mutation reversing it
     * survived the corpus).
     *
     * The two responses are listed newest-first, so a selector that simply kept the
     * last acceptable one - or that reversed the comparison - would pick the older.
     */
    private static XmlDiagnosticData crsTwoAcceptable() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, null);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);

        XmlRevocation newer = ocsp("R-OCSP-NEWER", root, "2023-06-01T00:00:00Z", "2023-06-01T00:00:00Z");
        newer.setCertHashExtensionPresent(true);
        newer.setCertHashExtensionMatch(true);
        responderIdRef(newer, root);
        revoke(signer, newer, CertificateStatus.GOOD);

        XmlRevocation older = ocsp("R-OCSP-OLDER", root, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        older.setCertHashExtensionPresent(true);
        older.setCertHashExtensionMatch(true);
        responderIdRef(older, root);
        revoke(signer, older, CertificateStatus.GOOD);

        XmlDiagnosticData data = data("crs-two-acceptable", signer, Arrays.asList(signer, root));
        data.getUsedRevocations().add(newer);
        data.getUsedRevocations().add(older);
        return data;
    }

    /**
     * A signing certificate with one acceptable and one UNACCEPTABLE revocation, so
     * that CertificateRevocationSelector reaches its "the CRS itself is valid but one
     * of its RACs is not" state - the only state in which its two overridden message
     * collectors branch: collectMessages() then drops the RAC-blockType constraint's
     * messages, and collectAdditionalMessages() takes its else branch and collects
     * only the valid RAC's messages.
     *
     * Neither the marshal-parity corpus nor the other dumps here reach it: every
     * certificate of theirs has either only-failing revocations (CRS invalid) or
     * only-passing ones (no RAC warning to filter). Audit mutations removing either
     * filter survived until this fixture existed.
     *
     * The unacceptable one is a CRL published before the certificate was issued
     * (RevocationAfterCertificateIssuanceCheck fails at FAIL level inside the RAC),
     * and it is the newer of the two, so a selector that ignored acceptability would
     * also pick the wrong latest revocation.
     */
    private static XmlDiagnosticData crsMixedAcceptance() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, null);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);

        XmlRevocation good = ocsp("R-OCSP-GOOD", root, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        good.setCertHashExtensionPresent(true);
        good.setCertHashExtensionMatch(true);
        responderIdRef(good, root);
        revoke(signer, good, CertificateStatus.GOOD);

        // thisUpdate before the signing certificate's notBefore: the revocation
        // cannot carry information about a certificate that did not exist yet.
        XmlRevocation stale = crl("R-CRL-BEFORE-ISSUANCE", root, "2015-01-01T00:00:00Z", "2023-06-01T00:00:00Z");
        revoke(signer, stale, CertificateStatus.GOOD);

        XmlDiagnosticData data = data("crs-mixed-acceptance", signer, Arrays.asList(signer, root));
        data.getUsedRevocations().add(good);
        data.getUsedRevocations().add(stale);
        return data;
    }

    /**
     * A revocation whose own certificate chain carries a certificate BEHIND its trust
     * anchor: it is signed by a trusted intermediate CA that is itself issued by an
     * untrusted root, so the chain reads [C-CA (trusted), C-ROOT (untrusted)].
     *
     * RevocationAcceptanceChecker#initChain() breaks out of its certificate-chain walk
     * at the first trust anchor. Every other dump here (and the whole marshal-parity
     * corpus) has the trust anchor last in the chain, where breaking and merely
     * skipping the entry are indistinguishable - an audit mutation turning that
     * break into a continue survived the corpus until this fixture existed.
     */
    private static XmlDiagnosticData racAnchorMidChain() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        XmlCertificate ca = certificate("C-CA", "2015-01-01T00:00:00Z", "2032-01-01T00:00:00Z");
        chain(ca, root);
        issuedBy(ca, root);
        // trust() marks the certificate self-signed; the CA is not, and its issuer
        // must stay visible behind it, so the trust flag is set on its own here.
        XmlTrusted trusted = new XmlTrusted();
        trusted.setValue(true);
        ca.setTrusted(trusted);

        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, ca, root);
        issuedBy(signer, ca);

        XmlRevocation ocsp = ocsp("R-OCSP-MIDANCHOR", ca, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocsp.getCertificateChain().clear();
        XmlChainItem caItem = new XmlChainItem();
        caItem.setCertificate(ca);
        ocsp.getCertificateChain().add(caItem);
        XmlChainItem rootItem = new XmlChainItem();
        rootItem.setCertificate(root);
        ocsp.getCertificateChain().add(rootItem);
        responderIdRef(ocsp, ca);
        revoke(signer, ocsp, CertificateStatus.GOOD);

        XmlDiagnosticData data = data("rac-anchor-mid-chain", signer, Arrays.asList(signer, ca, root));
        data.getUsedRevocations().add(ocsp);
        return data;
    }

    private static XmlDiagnosticData racPassed() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root, null, null);
        XmlCertificate signer = certificate("C-SIGNER", "2020-01-01T00:00:00Z", "2028-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);

        XmlRevocation ocsp = ocsp("R-OCSP-PASSED", root, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocsp.setCertHashExtensionPresent(true);
        ocsp.setCertHashExtensionMatch(true);
        responderIdRef(ocsp, root);
        revoke(signer, ocsp, CertificateStatus.GOOD);

        XmlDiagnosticData data = data("rac-passed", signer, Arrays.asList(signer, root));
        data.getUsedRevocations().add(ocsp);
        return data;
    }

    /**
     * Revocation shapes the marshal-parity corpus does not carry, for the rac checks:
     * a matching and a mismatching certHash extension, a responder-id reference that
     * resolves, a production date outside the issuer validity, an absent thisUpdate
     * and an absent production date, an expiredCertsOnCRL both before and after
     * thisUpdate, and an archiveCutoff.
     */
    private static XmlDiagnosticData racInputs() {
        XmlCertificate issuer = certificate("C-ISSUER", "2020-01-01T00:00:00Z", "2030-01-01T00:00:00Z");
        XmlCertificate subject = certificate("C-SUBJECT", "2022-01-01T00:00:00Z", "2026-01-01T00:00:00Z");
        chain(subject, issuer);
        issuedBy(subject, issuer);
        XmlCertificate expired = certificate("C-EXPIRED", "2015-01-01T00:00:00Z", "2018-01-01T00:00:00Z");
        chain(expired, issuer);
        issuedBy(expired, issuer);

        List<XmlRevocation> used = new ArrayList<>();

        // certHash present and matching, responder-id reference resolving, issuer valid
        // at production time, thisUpdate inside the subject validity.
        XmlRevocation ocspOk = ocsp("R-OCSP-OK", issuer, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocspOk.setCertHashExtensionPresent(true);
        ocspOk.setCertHashExtensionMatch(true);
        responderIdRef(ocspOk, issuer);
        used.add(ocspOk);
        revoke(subject, ocspOk, CertificateStatus.GOOD);

        // certHash present but not matching, unknown status, production date before the
        // issuer's notBefore, and the subject itself inside the revocation's chain.
        XmlRevocation ocspKo = ocsp("R-OCSP-KO", issuer, "2019-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        ocspKo.setCertHashExtensionPresent(true);
        ocspKo.setCertHashExtensionMatch(false);
        XmlChainItem self = new XmlChainItem();
        self.setCertificate(subject);
        ocspKo.getCertificateChain().add(self);
        used.add(ocspKo);
        revoke(subject, ocspKo, CertificateStatus.UNKNOWN);

        // No signing certificate, no thisUpdate, no production date.
        XmlRevocation bare = new XmlRevocation();
        bare.setId("R-BARE");
        bare.setOrigin(RevocationOrigin.EXTERNAL);
        bare.setType(RevocationType.CRL);
        XmlBasicSignature bareSignature = new XmlBasicSignature();
        bareSignature.setSignatureIntact(false);
        bareSignature.setSignatureValid(false);
        bare.setBasicSignature(bareSignature);
        used.add(bare);
        revoke(subject, bare, CertificateStatus.GOOD);

        // expiredCertsOnCRL before thisUpdate, over an expired certificate.
        XmlRevocation crlExpiredCerts = crl("R-CRL-EXPIREDCERTS", issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        crlExpiredCerts.setExpiredCertsOnCRL(date("2014-01-01T00:00:00Z"));
        used.add(crlExpiredCerts);
        revoke(expired, crlExpiredCerts, CertificateStatus.GOOD);

        // expiredCertsOnCRL NOT before thisUpdate: the logged-only branch.
        XmlRevocation crlLate = crl("R-CRL-LATE", issuer, "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        crlLate.setExpiredCertsOnCRL(date("2024-06-01T00:00:00Z"));
        used.add(crlLate);
        revoke(expired, crlLate, CertificateStatus.GOOD);

        // archiveCutoff before thisUpdate, over the same expired certificate.
        XmlRevocation ocspArchiveCutoff = ocsp("R-OCSP-ARCHIVECUTOFF", issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocspArchiveCutoff.setArchiveCutOff(date("2014-01-01T00:00:00Z"));
        used.add(ocspArchiveCutoff);
        revoke(expired, ocspArchiveCutoff, CertificateStatus.GOOD);

        XmlDiagnosticData data = data("rac-inputs", subject, Arrays.asList(subject, expired, issuer));
        data.getUsedRevocations().addAll(used);
        return data;
    }

    private static XmlRevocation ocsp(String id, XmlCertificate issuer, String thisUpdate, String productionDate) {
        return revocation(id, RevocationType.OCSP, issuer, thisUpdate, productionDate);
    }

    private static XmlRevocation crl(String id, XmlCertificate issuer, String thisUpdate, String productionDate) {
        return revocation(id, RevocationType.CRL, issuer, thisUpdate, productionDate);
    }

    private static XmlRevocation revocation(String id, RevocationType type, XmlCertificate issuer,
                                            String thisUpdate, String productionDate) {
        XmlRevocation revocation = new XmlRevocation();
        revocation.setId(id);
        revocation.setOrigin(RevocationOrigin.EXTERNAL);
        revocation.setType(type);
        if (thisUpdate != null) {
            revocation.setThisUpdate(date(thisUpdate));
        }
        if (productionDate != null) {
            revocation.setProductionDate(date(productionDate));
        }
        XmlSigningCertificate signing = new XmlSigningCertificate();
        signing.setCertificate(issuer);
        revocation.setSigningCertificate(signing);
        XmlChainItem item = new XmlChainItem();
        item.setCertificate(issuer);
        revocation.getCertificateChain().add(item);
        XmlBasicSignature basicSignature = new XmlBasicSignature();
        basicSignature.setSignatureIntact(true);
        basicSignature.setSignatureValid(true);
        revocation.setBasicSignature(basicSignature);
        return revocation;
    }

    private static void responderIdRef(XmlRevocation revocation, XmlCertificate certificate) {
        XmlRelatedCertificate related = new XmlRelatedCertificate();
        related.setCertificate(certificate);
        XmlCertificateRef ref = new XmlCertificateRef();
        ref.setOrigin(CertificateRefOrigin.SIGNING_CERTIFICATE);
        related.getCertificateRefs().add(ref);
        XmlFoundCertificates found = new XmlFoundCertificates();
        found.getRelatedCertificates().add(related);
        revocation.setFoundCertificates(found);
    }

    private static void revoke(XmlCertificate certificate, XmlRevocation revocation, CertificateStatus status) {
        XmlCertificateRevocation certificateRevocation = new XmlCertificateRevocation();
        certificateRevocation.setRevocation(revocation);
        certificateRevocation.setStatus(status);
        certificate.getRevocations().add(certificateRevocation);
    }

    // ----------------------------------------------------------------- model

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
        signature.setBasicSignature(basicSignature);
        data.getSignatures().add(signature);

        data.getUsedCertificates().addAll(usedCertificates);
        return data;
    }

    private static XmlCertificate certificate(String id, String notBefore, String notAfter) {
        XmlCertificate certificate = new XmlCertificate();
        certificate.setId(id);
        certificate.setNotBefore(date(notBefore));
        certificate.setNotAfter(date(notAfter));
        certificate.setSelfSigned(false);
        certificate.setPublicKeySize(2048);
        certificate.setPublicKeyEncryptionAlgo(EncryptionAlgorithm.RSA);
        certificate.setEntityKey(id + "-KEY");
        XmlBasicSignature basicSignature = new XmlBasicSignature();
        basicSignature.setSignatureIntact(true);
        basicSignature.setSignatureValid(true);
        basicSignature.setDigestAlgoUsedToSignThisToken(DigestAlgorithm.SHA256);
        basicSignature.setEncryptionAlgoUsedToSignThisToken(EncryptionAlgorithm.RSA);
        basicSignature.setKeyLengthUsedToSignThisToken("2048");
        certificate.setBasicSignature(basicSignature);
        return certificate;
    }

    /** Marks the certificate trusted, optionally with a trust start and sunset date. */
    private static void trust(XmlCertificate certificate, String startDate, String sunsetDate) {
        XmlTrusted trusted = new XmlTrusted();
        trusted.setValue(true);
        if (startDate != null) {
            trusted.setStartDate(date(startDate));
        }
        if (sunsetDate != null) {
            trusted.setSunsetDate(date(sunsetDate));
        }
        certificate.setTrusted(trusted);
        certificate.setSelfSigned(true);
    }

    private static void chain(XmlCertificate certificate, XmlCertificate... parents) {
        List<XmlChainItem> items = new ArrayList<>();
        for (XmlCertificate parent : parents) {
            XmlChainItem item = new XmlChainItem();
            item.setCertificate(parent);
            items.add(item);
        }
        certificate.getCertificateChain().addAll(items);
    }

    private static void issuedBy(XmlCertificate certificate, XmlCertificate issuer) {
        XmlSigningCertificate signing = new XmlSigningCertificate();
        signing.setCertificate(issuer);
        certificate.setSigningCertificate(signing);
    }

    private static Date date(String iso) {
        return Date.from(Instant.parse(iso));
    }

}
