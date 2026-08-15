/*
 * Java oracle driver for the twelve
 * eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks classes (the XCVA half
 * of phase 8d).
 *
 * Every check is driven directly through a chain of exactly one item at
 * Level.FAIL, over REAL corpus inputs: each dump's used certificates crossed with
 * each of their certificate revocation data. Driving the same check over every
 * (certificate, revocation) pair is what produces both the happy and the failure
 * row wherever the corpus varies; the branches the 50-dump marshal-parity corpus
 * never reaches are additionally driven over synthetic wrappers built here from
 * the generated diagnostic model, so that every check keeps an OK and a NOT OK
 * row.
 *
 * Rows are written to
 *   dss/validation/process/bbb/xcv/testdata/oracle/xcva_direct.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle XcvaDirectOracle.java
 *   java  -cp "$CP:/tmp/oracle" XcvaDirectOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlRAC;
import eu.europa.esig.dss.diagnostic.CertificateRevocationWrapper;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate;
import eu.europa.esig.dss.enumerations.CertificateRefOrigin;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.enumerations.RevocationType;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationAcceptanceCheckerResultCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationAfterCertificateIssuanceCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationCertHashMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationCertHashPresenceCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationDataKnownCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationHasInformationAboutCertificateCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerKnownCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerRevocationDataAvailableCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationIssuerValidAtProductionTimeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.RevocationResponderIdMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.SelfIssuedOCSPCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rac.checks.ThisUpdatePresenceCheck;

import javax.xml.datatype.DatatypeFactory;
import javax.xml.datatype.XMLGregorianCalendar;
import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.GregorianCalendar;
import java.util.List;
import java.util.function.BiFunction;

public class XcvaDirectOracle {

    private static final LevelRule FAIL = () -> Level.FAIL;

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();

        try (PrintWriter out = writer(repoRoot,
                "validation/process/bbb/xcv/testdata/oracle/xcva_direct.jsonl")) {

            for (File file : files) {
                String name = file.getName();
                // See ../../../fc/testdata/README.md: the four model-*.xml dumps are
                // schema-coverage fixtures whose IDREF graph is dangling and whose
                // attribute values carry raw control characters, so the two runtimes
                // would be reading two different inputs. Every row below comes from a
                // real dump or from an explicit synthetic.
                if (name.startsWith("model-")) {
                    System.err.println("SKIPPED-FIXTURE " + name);
                    continue;
                }
                DiagnosticData dd;
                try {
                    XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                    dd = new DiagnosticData(jaxb);
                    dd.getUsedCertificates();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED-FILE " + name + ": " + e);
                    continue;
                }

                for (CertificateWrapper certificate : dd.getUsedCertificates()) {
                    row(out, i18n, name, certificate.getId(), "RevocationIssuerRevocationDataAvailableCheck",
                            (r, l) -> new RevocationIssuerRevocationDataAvailableCheck(i18n, r, certificate, l));

                    for (CertificateRevocationWrapper revocation : certificate.getCertificateRevocationData()) {
                        String token = certificate.getId() + "|" + revocation.getId();
                        emitAll(out, i18n, name, token, certificate, revocation);
                    }
                }
            }

            // ------------------------------------------------------------- synthetics
            //
            // The corpus carries no revocation with a certHash extension, none with an
            // absent thisUpdate, none whose responder-id reference resolves, and none
            // whose production date falls outside the issuer validity, so the OK/NOT OK
            // pair of the checks reading those is completed from wrappers built here.
            for (Synthetic synthetic : synthetics()) {
                emitAll(out, i18n, "synthetic", synthetic.name, synthetic.certificate, synthetic.revocation);
            }

            // RevocationAcceptanceCheckerResultCheck over hand-built XmlRAC results:
            // a PASSED one, an INDETERMINATE one carrying errors and both dates, and one
            // with no production date at all (the null-additionalInfo branch).
            racResultRows(out, i18n);
        }
    }

    /** Drives every rac check that reads a (certificate, revocation) pair. */
    private static void emitAll(PrintWriter out, I18nProvider i18n, String file, String token,
                                CertificateWrapper certificate, CertificateRevocationWrapper revocation) {
        row(out, i18n, file, token, "RevocationDataKnownCheck",
                (r, l) -> new RevocationDataKnownCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationIssuerKnownCheck",
                (r, l) -> new RevocationIssuerKnownCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "ThisUpdatePresenceCheck",
                (r, l) -> new ThisUpdatePresenceCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationIssuerValidAtProductionTimeCheck",
                (r, l) -> new RevocationIssuerValidAtProductionTimeCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationResponderIdMatchCheck",
                (r, l) -> new RevocationResponderIdMatchCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationCertHashPresenceCheck",
                (r, l) -> new RevocationCertHashPresenceCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "RevocationCertHashMatchCheck",
                (r, l) -> new RevocationCertHashMatchCheck(i18n, r, revocation, l));
        row(out, i18n, file, token, "SelfIssuedOCSPCheck",
                (r, l) -> new SelfIssuedOCSPCheck(i18n, r, certificate, revocation, l));
        row(out, i18n, file, token, "RevocationAfterCertificateIssuanceCheck",
                (r, l) -> new RevocationAfterCertificateIssuanceCheck(i18n, r, certificate, revocation, l));
        row(out, i18n, file, token, "RevocationHasInformationAboutCertificateCheck",
                (r, l) -> new RevocationHasInformationAboutCertificateCheck(i18n, r, certificate, revocation, l));
    }

    private static void racResultRows(PrintWriter out, I18nProvider i18n) {
        XmlRAC passed = new XmlRAC();
        passed.setId("R-PASSED");
        passed.setRevocationThisUpdate(date("2023-06-01T00:00:00Z"));
        passed.setRevocationProductionDate(date("2023-06-02T00:00:00Z"));
        XmlConclusion passedConclusion = new XmlConclusion();
        passedConclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
        passed.setConclusion(passedConclusion);
        row(out, i18n, "synthetic", "rac-passed", "RevocationAcceptanceCheckerResultCheck",
                (r, l) -> new RevocationAcceptanceCheckerResultCheck<>(i18n, r, passed, l));

        XmlRAC failed = new XmlRAC();
        failed.setId("R-FAILED");
        failed.setRevocationThisUpdate(date("2023-06-01T00:00:00Z"));
        failed.setRevocationProductionDate(date("2023-06-02T00:00:00Z"));
        XmlConclusion failedConclusion = new XmlConclusion();
        failedConclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.INDETERMINATE);
        failedConclusion.setSubIndication(eu.europa.esig.dss.enumerations.SubIndication.TRY_LATER);
        XmlMessage error = new XmlMessage();
        error.setKey("BBB_XCV_IARDPFC_ANS");
        error.setValue("No acceptable revocation data found");
        failedConclusion.getErrors().add(error);
        failed.setConclusion(failedConclusion);
        row(out, i18n, "synthetic", "rac-failed", "RevocationAcceptanceCheckerResultCheck",
                (r, l) -> new RevocationAcceptanceCheckerResultCheck<>(i18n, r, failed, l));

        XmlRAC noDates = new XmlRAC();
        noDates.setId("R-NO-DATES");
        XmlConclusion noDatesConclusion = new XmlConclusion();
        noDatesConclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
        noDates.setConclusion(noDatesConclusion);
        row(out, i18n, "synthetic", "rac-no-dates", "RevocationAcceptanceCheckerResultCheck",
                (r, l) -> new RevocationAcceptanceCheckerResultCheck<>(i18n, r, noDates, l));
    }

    // ------------------------------------------------------------------ synthetics

    /** A synthetic (certificate, revocation) pair built straight on the generated model. */
    static class Synthetic {
        final String name;
        final CertificateWrapper certificate;
        final CertificateRevocationWrapper revocation;

        Synthetic(String name, CertificateWrapper certificate, CertificateRevocationWrapper revocation) {
            this.name = name;
            this.certificate = certificate;
            this.revocation = revocation;
        }
    }

    private static List<Synthetic> synthetics() {
        List<Synthetic> out = new ArrayList<>();

        // 1. OCSP, certHash present and matching, issuer valid at production time,
        //    responder-id reference resolving to the signing certificate, thisUpdate
        //    inside the certificate validity, issuer not the certificate itself.
        XmlCertificate issuer = certificate("C-ISSUER", "2020-01-01T00:00:00Z", "2030-01-01T00:00:00Z");
        XmlCertificate subject = certificate("C-SUBJECT", "2022-01-01T00:00:00Z", "2026-01-01T00:00:00Z");
        XmlRevocation ocsp = revocation("R-OCSP-OK", RevocationType.OCSP, issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocsp.setCertHashExtensionPresent(true);
        ocsp.setCertHashExtensionMatch(true);
        addResponderIdRef(ocsp, issuer);
        out.add(new Synthetic("ocsp-certhash-ok",
                new CertificateWrapper(subject), certificateRevocation(ocsp, true)));

        // 2. Same, with a certHash that does not match, an unknown status, a production
        //    date before the issuer's notBefore, a thisUpdate before the certificate's
        //    notBefore and the certificate itself in the revocation's chain.
        XmlRevocation ocspKo = revocation("R-OCSP-KO", RevocationType.OCSP, issuer,
                "2019-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        ocspKo.setCertHashExtensionPresent(true);
        ocspKo.setCertHashExtensionMatch(false);
        chain(ocspKo, subject);
        out.add(new Synthetic("ocsp-certhash-ko",
                new CertificateWrapper(subject), certificateRevocation(ocspKo, false)));

        // 3. A revocation with no signing certificate and no thisUpdate at all.
        XmlRevocation bare = revocation("R-BARE", RevocationType.CRL, null, null, null);
        out.add(new Synthetic("crl-bare",
                new CertificateWrapper(subject), certificateRevocation(bare, true)));

        // 4. A CRL carrying expiredCertsOnCRL before its thisUpdate, over an expired
        //    certificate: the expiredCertsOnCRL branch of
        //    RevocationHasInformationAboutCertificateCheck.
        XmlCertificate expired = certificate("C-EXPIRED", "2015-01-01T00:00:00Z", "2018-01-01T00:00:00Z");
        XmlRevocation crlExpiredCerts = revocation("R-CRL-EXPIREDCERTS", RevocationType.CRL, issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        crlExpiredCerts.setExpiredCertsOnCRL(date("2014-01-01T00:00:00Z"));
        out.add(new Synthetic("crl-expiredcertsoncrl",
                new CertificateWrapper(expired), certificateRevocation(crlExpiredCerts, true)));

        // 5. An OCSP carrying archiveCutoff before its thisUpdate, over the same expired
        //    certificate: the archiveCutoff branch.
        XmlRevocation ocspArchiveCutoff = revocation("R-OCSP-ARCHIVECUTOFF", RevocationType.OCSP, issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        ocspArchiveCutoff.setArchiveCutOff(date("2014-01-01T00:00:00Z"));
        out.add(new Synthetic("ocsp-archivecutoff",
                new CertificateWrapper(expired), certificateRevocation(ocspArchiveCutoff, true)));

        // 6. A CRL whose expiredCertsOnCRL is NOT before thisUpdate (the logged-only
        //    branch), over a certificate that has already expired.
        XmlRevocation crlLateExpiredCerts = revocation("R-CRL-LATE", RevocationType.CRL, issuer,
                "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z");
        crlLateExpiredCerts.setExpiredCertsOnCRL(date("2024-01-01T00:00:00Z"));
        out.add(new Synthetic("crl-expiredcertsoncrl-late",
                new CertificateWrapper(expired), certificateRevocation(crlLateExpiredCerts, true)));

        return out;
    }

    private static XmlCertificate certificate(String id, String notBefore, String notAfter) {
        XmlCertificate certificate = new XmlCertificate();
        certificate.setId(id);
        certificate.setNotBefore(date(notBefore));
        certificate.setNotAfter(date(notAfter));
        return certificate;
    }

    private static XmlRevocation revocation(String id, RevocationType type, XmlCertificate signingCertificate,
                                            String thisUpdate, String productionDate) {
        XmlRevocation revocation = new XmlRevocation();
        revocation.setId(id);
        revocation.setType(type);
        if (thisUpdate != null) {
            revocation.setThisUpdate(date(thisUpdate));
        }
        if (productionDate != null) {
            revocation.setProductionDate(date(productionDate));
        }
        if (signingCertificate != null) {
            XmlSigningCertificate signing = new XmlSigningCertificate();
            signing.setCertificate(signingCertificate);
            revocation.setSigningCertificate(signing);
            chain(revocation, signingCertificate);
        }
        return revocation;
    }

    private static void chain(XmlRevocation revocation, XmlCertificate certificate) {
        XmlChainItem item = new XmlChainItem();
        item.setCertificate(certificate);
        revocation.getCertificateChain().add(item);
    }

    private static void addResponderIdRef(XmlRevocation revocation, XmlCertificate certificate) {
        XmlRelatedCertificate related = new XmlRelatedCertificate();
        related.setCertificate(certificate);
        XmlCertificateRef ref = new XmlCertificateRef();
        ref.setOrigin(CertificateRefOrigin.SIGNING_CERTIFICATE);
        related.getCertificateRefs().add(ref);
        XmlFoundCertificates found = new XmlFoundCertificates();
        found.getRelatedCertificates().add(related);
        revocation.setFoundCertificates(found);
    }

    private static CertificateRevocationWrapper certificateRevocation(XmlRevocation revocation, boolean known) {
        XmlCertificateRevocation certificateRevocation = new XmlCertificateRevocation();
        certificateRevocation.setRevocation(revocation);
        certificateRevocation.setStatus(known
                ? eu.europa.esig.dss.enumerations.CertificateStatus.GOOD
                : eu.europa.esig.dss.enumerations.CertificateStatus.UNKNOWN);
        return new CertificateRevocationWrapper(certificateRevocation);
    }

    // ------------------------------------------------------------------ harness

    /** A chain of exactly one item over an XmlRAC. */
    static class SingleRACChain extends Chain<XmlRAC> {
        private final BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory;

        SingleRACChain(I18nProvider i18nProvider, BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory) {
            super(i18nProvider, new XmlRAC());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    private static void row(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                            BiFunction<XmlRAC, LevelRule, ChainItem<XmlRAC>> factory) {
        try {
            XmlRAC result = new SingleRACChain(i18n, factory).execute();
            out.println(json(file, token, check, result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    // --------------------------------------------------------------------- JSON

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static String json(String file, String token, String check, XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(token)).append(',');
        key(out, "check").append(str(check)).append(',');
        key(out, "context").append(str(Context.REVOCATION.name())).append(',');
        key(out, "block").append(str("RAC")).append(',');
        key(out, "title").append(str(result.getTitle())).append(',');
        key(out, "conclusion").append(conclusion(result.getConclusion())).append(',');
        key(out, "constraints").append('[');
        List<XmlConstraint> constraints = result.getConstraint();
        for (int i = 0; i < constraints.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(constraint(constraints.get(i)));
        }
        out.append(']');
        return out.append('}').toString();
    }

    private static StringBuilder key(StringBuilder out, String name) {
        return out.append(str(name)).append(':');
    }

    private static String conclusion(XmlConclusion conclusion) {
        if (conclusion == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "indication").append(str(conclusion.getIndication() == null ? null : conclusion.getIndication().name())).append(',');
        key(out, "subIndication").append(str(conclusion.getSubIndication() == null ? null : conclusion.getSubIndication().name())).append(',');
        key(out, "errors").append(messages(conclusion.getErrors())).append(',');
        key(out, "warnings").append(messages(conclusion.getWarnings())).append(',');
        key(out, "infos").append(messages(conclusion.getInfos()));
        return out.append('}').toString();
    }

    private static String messages(List<XmlMessage> messages) {
        StringBuilder out = new StringBuilder("[");
        for (int i = 0; i < messages.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(message(messages.get(i)));
        }
        return out.append(']').toString();
    }

    private static String message(XmlMessage message) {
        if (message == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "key").append(str(message.getKey())).append(',');
        key(out, "value").append(str(message.getValue()));
        return out.append('}').toString();
    }

    private static String constraint(XmlConstraint constraint) {
        StringBuilder out = new StringBuilder("{");
        key(out, "name").append(message(constraint.getName())).append(',');
        key(out, "status").append(str(constraint.getStatus() == null ? null : constraint.getStatus().value())).append(',');
        key(out, "error").append(message(constraint.getError())).append(',');
        key(out, "warning").append(message(constraint.getWarning())).append(',');
        key(out, "info").append(message(constraint.getInfo())).append(',');
        key(out, "additionalInfo").append(str(constraint.getAdditionalInfo())).append(',');
        key(out, "id").append(str(constraint.getId())).append(',');
        key(out, "blockType").append(str(constraint.getBlockType() == null ? null : constraint.getBlockType().value()));
        return out.append('}').toString();
    }

    private static String str(String value) {
        if (value == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("\"");
        for (int i = 0; i < value.length(); i++) {
            char c = value.charAt(i);
            switch (c) {
                case '"': out.append("\\\""); break;
                case '\\': out.append("\\\\"); break;
                case '\n': out.append("\\n"); break;
                case '\r': out.append("\\r"); break;
                case '\t': out.append("\\t"); break;
                default:
                    if (c < 0x20) {
                        out.append(String.format("\\u%04x", (int) c));
                    } else {
                        out.append(c);
                    }
            }
        }
        return out.append('"').toString();
    }

    private static Date date(String iso) {
        return Date.from(java.time.Instant.parse(iso));
    }

}
