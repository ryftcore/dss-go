/*
 * Java oracle driver for the AOV (algorithm-obsolescence validation) building
 * blocks and the cc (cryptographic checker) family of phase 8d:
 *
 *   eu.europa.esig.dss.validation.process.bbb.aov.*                (10 concrete blocks)
 *   eu.europa.esig.dss.validation.process.bbb.aov.cc.*             (2 checkers)
 *   eu.europa.esig.dss.validation.process.bbb.aov.cc.checks.*      (6 leaf checks, wired)
 *   eu.europa.esig.dss.validation.process.bbb.aov.checks.*         (4 result checks)
 *
 * It runs the UPSTREAM classes over REAL inputs - every XmlDiagnosticData dump of
 * the marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle plus the nine
 * synthetic dumps XCVA's XcvaSyntheticDumps writes into
 * ../../xcv/testdata/dd - with the default ETSI validation policy and the fixed
 * validation time 2024-01-01T00:00:00Z, and dumps one JSON object per line.
 *
 * Three corpora are written into ../oracle:
 *
 *   aov_blocks.jsonl - every concrete AlgorithmObsolescenceValidation subclass over
 *                      every token of every dump; the row carries the whole XmlAOV
 *                      (constraints, conclusion, and all four XmlCryptographicValidation
 *                      members) so the Go replay compares the produced tree and not
 *                      only its conclusion.
 *   aov_cc.jsonl     - SignatureAlgorithmCryptographicChecker and
 *                      DigestAlgorithmCryptographicChecker over an algorithm x key
 *                      length x validation time matrix chosen so that every one of
 *                      the six wired cc/checks classes gets an OK and a NOT OK row,
 *                      including the boundary instants of the policy's expiration
 *                      dates (expiry-1ms, expiry, expiry+1ms).
 *   aov_direct.jsonl - the four aov/checks result checks driven alone through a
 *                      one-item chain, over hand-built XmlAOV / XmlCC shapes at
 *                      every branch they take.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle AovOracle.java
 *   java  -cp "$CP:/tmp/oracle" AovOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlAOV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCertificateChainCryptographicValidation;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicAlgorithm;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicValidation;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSAV;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.EvidenceRecordWrapper;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.enumerations.SignatureAlgorithm;
import eu.europa.esig.dss.enumerations.SubContext;
import eu.europa.esig.dss.enumerations.SubIndication;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.i18n.MessageTag;
import eu.europa.esig.dss.model.policy.CryptographicSuite;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.ValidationProcessUtils;
import eu.europa.esig.dss.validation.process.bbb.aov.CertificateAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.CertificateAndChainAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.EvidenceRecordAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.RevocationDataAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.SignatureAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.SignatureSignedDataAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.TimestampAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.TokenCertificateChainAlgorithmObsolescenceValidation;
import eu.europa.esig.dss.validation.process.bbb.aov.cc.DigestAlgorithmCryptographicChecker;
import eu.europa.esig.dss.validation.process.bbb.aov.cc.SignatureAlgorithmCryptographicChecker;
import eu.europa.esig.dss.validation.process.bbb.aov.checks.AlgorithmObsolescenceValidationCheck;
import eu.europa.esig.dss.validation.process.bbb.aov.checks.AlgorithmObsolescenceValidationCheckWithId;
import eu.europa.esig.dss.validation.process.bbb.aov.checks.DigestAlgorithmCryptographicCheckerResultCheck;
import eu.europa.esig.dss.validation.process.bbb.aov.checks.SignatureAlgorithmCryptographicCheckerResultCheck;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.io.PrintWriter;
import java.text.SimpleDateFormat;
import java.util.Arrays;
import java.util.Comparator;
import java.util.Date;
import java.util.List;
import java.util.TimeZone;

public class AovOracle {

    /** Fixed validation time: 2024-01-01T00:00:00Z. */
    private static final Date CURRENT_TIME = new Date(1704067200000L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        File[] synthetic = new File(repoRoot, "validation/process/bbb/xcv/testdata/dd")
                .listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(synthetic, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();
        ValidationPolicy policy = new EtsiValidationPolicyFactory().loadDefaultValidationPolicy();

        try (PrintWriter out = writer(repoRoot, "validation/process/bbb/aov/testdata/oracle/aov_blocks.jsonl")) {
            for (File synth : synthetic) {
                DiagnosticData dd = new DiagnosticData(
                        DiagnosticDataFacade.newFacade().unmarshall(synth, false));
                emitBlocks(out, i18n, policy, "dd/" + synth.getName(), dd);
            }
            for (File file : files) {
                String name = file.getName();
                // See ../../fc/testdata/README.md for why the four model-*.xml
                // schema-coverage fixtures are excluded.
                if (name.startsWith("model-")) {
                    System.err.println("SKIPPED-FIXTURE " + name);
                    continue;
                }
                XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                emitBlocks(out, i18n, policy, name, new DiagnosticData(jaxb));
            }
        }

        try (PrintWriter out = writer(repoRoot, "validation/process/bbb/aov/testdata/oracle/aov_cc.jsonl")) {
            emitCC(out, i18n, policy);
        }

        try (PrintWriter out = writer(repoRoot, "validation/process/bbb/aov/testdata/oracle/aov_direct.jsonl")) {
            emitDirect(out, i18n);
        }
    }

    // ------------------------------------------------------------------ blocks

    private static void emitBlocks(PrintWriter out, I18nProvider i18n, ValidationPolicy policy,
                                   String file, DiagnosticData dd) {
        for (SignatureWrapper signature : dd.getSignatures()) {
            Context context = signature.isCounterSignature() ? Context.COUNTER_SIGNATURE : Context.SIGNATURE;
            String id = signature.getId();
            block(out, file, "SIG|" + id, "SignatureAlgorithmObsolescenceValidation", context,
                    () -> new SignatureAlgorithmObsolescenceValidation<>(i18n, signature, context, CURRENT_TIME, policy).execute());
            block(out, file, "SIG|" + id, "SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation", context,
                    () -> new SignatureValueAndSignedAttributesAlgorithmObsolescenceValidation<>(i18n, signature, context, CURRENT_TIME, policy).execute());
            block(out, file, "SIG|" + id, "SignatureSignedDataAlgorithmObsolescenceValidation", context,
                    () -> new SignatureSignedDataAlgorithmObsolescenceValidation(i18n, signature, context, CURRENT_TIME, policy).execute());
            block(out, file, "SIG|" + id, "TokenCertificateChainAlgorithmObsolescenceValidation", context,
                    () -> new TokenCertificateChainAlgorithmObsolescenceValidation<>(i18n, signature, context, CURRENT_TIME, policy).execute());
        }
        for (TimestampWrapper timestamp : dd.getTimestampList()) {
            block(out, file, "TST|" + timestamp.getId(), "TimestampAlgorithmObsolescenceValidation", Context.TIMESTAMP,
                    () -> new TimestampAlgorithmObsolescenceValidation(i18n, timestamp, CURRENT_TIME, policy).execute());
            block(out, file, "TST|" + timestamp.getId(), "TokenCertificateChainAlgorithmObsolescenceValidation", Context.TIMESTAMP,
                    () -> new TokenCertificateChainAlgorithmObsolescenceValidation<>(i18n, timestamp, Context.TIMESTAMP, CURRENT_TIME, policy).execute());
        }
        for (RevocationWrapper revocation : sortedRevocations(dd)) {
            block(out, file, "REV|" + revocation.getId(), "RevocationDataAlgorithmObsolescenceValidation", Context.REVOCATION,
                    () -> new RevocationDataAlgorithmObsolescenceValidation(i18n, revocation, CURRENT_TIME, policy).execute());
        }
        for (EvidenceRecordWrapper evidenceRecord : dd.getEvidenceRecords()) {
            block(out, file, "ER|" + evidenceRecord.getId(), "EvidenceRecordAlgorithmObsolescenceValidation", Context.EVIDENCE_RECORD,
                    () -> new EvidenceRecordAlgorithmObsolescenceValidation(i18n, evidenceRecord, CURRENT_TIME, policy).execute());
        }
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            for (SubContext subContext : new SubContext[] { SubContext.SIGNING_CERT, SubContext.CA_CERTIFICATE }) {
                block(out, file, "CERT|" + certificate.getId() + "|" + subContext.name(),
                        "CertificateAlgorithmObsolescenceValidation", Context.CERTIFICATE,
                        () -> new CertificateAlgorithmObsolescenceValidation(i18n, certificate, Context.CERTIFICATE,
                                subContext, CURRENT_TIME, policy).execute());
            }
            block(out, file, "CERT|" + certificate.getId(), "CertificateAndChainAlgorithmObsolescenceValidation",
                    Context.CERTIFICATE,
                    () -> new CertificateAndChainAlgorithmObsolescenceValidation(i18n, certificate, Context.CERTIFICATE,
                            CURRENT_TIME, policy).execute());
        }
    }

    private static List<RevocationWrapper> sortedRevocations(DiagnosticData dd) {
        List<RevocationWrapper> revocations = new java.util.ArrayList<>(dd.getAllRevocationData());
        revocations.sort(Comparator.comparing(RevocationWrapper::getId));
        return revocations;
    }

    private interface AovRun {
        XmlAOV run();
    }

    private static void block(PrintWriter out, String file, String token, String blockName, Context context, AovRun run) {
        XmlAOV result;
        try {
            result = run.run();
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + blockName + ": " + e);
            return;
        }
        StringBuilder row = new StringBuilder("{");
        key(row, "file").append(str(file)).append(',');
        key(row, "token").append(str(token)).append(',');
        key(row, "block").append(str(blockName)).append(',');
        key(row, "context").append(str(context.name())).append(',');
        key(row, "validationTime").append(str(iso(result.getValidationTime()))).append(',');
        body(row, result);
        row.append(',');
        key(row, "signatureCryptographicValidation").append(cv(result.getSignatureCryptographicValidation())).append(',');
        key(row, "signedAttributesValidation").append(cv(result.getSignedAttributesValidation())).append(',');
        key(row, "digestMatchersValidation").append(cv(result.getDigestMatchersValidation())).append(',');
        key(row, "certificateChainCryptographicValidation").append(chain(result.getCertificateChainCryptographicValidation()));
        out.println(row.append('}').toString());
    }

    // ---------------------------------------------------------------------- cc

    /**
     * The cc matrix. Chosen so that each of the six wired leaf checks reaches both
     * of its branches, and so that the reliability/expiration comparisons are probed
     * at the exact instant the policy's expiration date names, one millisecond before
     * it and one millisecond after it.
     */
    private static void emitCC(PrintWriter out, I18nProvider i18n, ValidationPolicy policy) {
        CryptographicSuite suite = policy.getSignatureCryptographicConstraint(Context.SIGNATURE);
        MessageTag position = MessageTag.ACCM_POS_SIG_SIG;

        for (DigestAlgorithm digestAlgorithm : new DigestAlgorithm[] {
                DigestAlgorithm.SHA1, DigestAlgorithm.SHA224, DigestAlgorithm.SHA256, DigestAlgorithm.SHA512,
                DigestAlgorithm.MD5, DigestAlgorithm.MD2, DigestAlgorithm.RIPEMD160, DigestAlgorithm.SHA3_256, null }) {
            for (Date date : datesFor(suite, digestAlgorithm)) {
                final DigestAlgorithm da = digestAlgorithm;
                final Date d = date;
                cc(out, "DIGEST|" + (da == null ? "null" : da.name()) + "|" + iso(d),
                        "DigestAlgorithmCryptographicChecker",
                        () -> new DigestAlgorithmCryptographicChecker(i18n, da, d, position, suite).execute());
            }
        }

        // No "" entry: TokenProxy#getKeyLengthUsedToSignThisToken is nullable and the
        // Go port's wrapper collapses that null to "", so "" is the port's spelling of
        // null and an empty-but-non-null Java key length is unreachable from any caller.
        String[] keyLengths = { "1024", "1900", "2048", "3072", "0", null, "256", "not-a-number" };
        for (SignatureAlgorithm signatureAlgorithm : new SignatureAlgorithm[] {
                SignatureAlgorithm.RSA_SHA256, SignatureAlgorithm.RSA_SHA1, SignatureAlgorithm.DSA_SHA256,
                SignatureAlgorithm.ECDSA_SHA256, SignatureAlgorithm.ECDSA_SHA1, SignatureAlgorithm.RSA_SSA_PSS_SHA256_MGF1,
                SignatureAlgorithm.ED25519, SignatureAlgorithm.HMAC_SHA256, null }) {
            for (String keyLength : keyLengths) {
                for (Date date : datesFor(suite, signatureAlgorithm, keyLength)) {
                    final SignatureAlgorithm sa = signatureAlgorithm;
                    final String kl = keyLength;
                    final Date d = date;
                    cc(out, "SIG|" + (sa == null ? "null" : sa.name()) + "|" + (kl == null ? "null" : kl) + "|" + iso(d),
                            "SignatureAlgorithmCryptographicChecker",
                            () -> new SignatureAlgorithmCryptographicChecker(i18n, sa, kl, d, position, suite).execute());
                }
            }
        }
    }

    /** validation time, plus the exact expiration instant of this algorithm and its two neighbours. */
    private static Date[] datesFor(CryptographicSuite suite, DigestAlgorithm digestAlgorithm) {
        Date expiration = null;
        try {
            expiration = eu.europa.esig.dss.validation.policy.CryptographicSuiteUtils
                    .getExpirationDate(suite, digestAlgorithm);
        } catch (RuntimeException e) {
            // no expiration defined
        }
        return boundary(expiration);
    }

    private static Date[] datesFor(CryptographicSuite suite, SignatureAlgorithm signatureAlgorithm, String keyLength) {
        Date expiration = null;
        try {
            expiration = eu.europa.esig.dss.validation.policy.CryptographicSuiteUtils
                    .getExpirationDate(suite, signatureAlgorithm, keyLength);
        } catch (RuntimeException e) {
            // no expiration defined
        }
        return boundary(expiration);
    }

    private static Date[] boundary(Date expiration) {
        if (expiration == null) {
            return new Date[] { CURRENT_TIME };
        }
        return new Date[] {
                CURRENT_TIME,
                new Date(expiration.getTime() - 1L),
                new Date(expiration.getTime()),
                new Date(expiration.getTime() + 1L),
        };
    }

    private interface CcRun {
        XmlCC run();
    }

    private static void cc(PrintWriter out, String token, String blockName, CcRun run) {
        XmlCC result;
        try {
            result = run.run();
        } catch (RuntimeException e) {
            System.err.println("SKIPPED cc " + token + ": " + e);
            return;
        }
        StringBuilder row = new StringBuilder("{");
        key(row, "token").append(str(token)).append(',');
        key(row, "block").append(str(blockName)).append(',');
        body(row, result);
        row.append(',');
        key(row, "cryptographicValidation").append(cv(result.getCryptographicValidation()));
        out.println(row.append('}').toString());
    }

    // ------------------------------------------------------------------ direct

    /** A one-item chain, the phase 8c FcSavDirectOracle pattern. */
    private static class SingleSAVChain extends Chain<XmlSAV> {
        private final java.util.function.Function<XmlSAV, ChainItem<XmlSAV>> factory;

        SingleSAVChain(I18nProvider i18nProvider, java.util.function.Function<XmlSAV, ChainItem<XmlSAV>> factory) {
            super(i18nProvider, new XmlSAV());
            this.factory = factory;
        }

        @Override
        protected void initChain() {
            firstItem = factory.apply(result);
        }
    }

    private static void emitDirect(PrintWriter out, I18nProvider i18n) {
        LevelRule fail = ValidationProcessUtils.getLevelRule(Level.FAIL);
        MessageTag position = MessageTag.ACCM_POS_SIG_SIG;

        for (String shape : new String[] { "passed", "passed-with-algo", "error", "warning", "info" }) {
            XmlAOV aov = aovOfShape(shape);
            direct(out, "AlgorithmObsolescenceValidationCheck|" + shape, i18n,
                    result -> new AlgorithmObsolescenceValidationCheck<>(i18n, result, aov, CURRENT_TIME, position, "T-1"));
            direct(out, "AlgorithmObsolescenceValidationCheckWithId|" + shape, i18n,
                    result -> new AlgorithmObsolescenceValidationCheckWithId<>(i18n, result, aov, CURRENT_TIME, position, "T-1"));
        }

        for (String shape : new String[] { "passed", "passed-no-keylength", "failed" }) {
            XmlCC cc = ccOfShape(shape);
            direct(out, "SignatureAlgorithmCryptographicCheckerResultCheck|" + shape, i18n,
                    result -> new SignatureAlgorithmCryptographicCheckerResultCheck<>(i18n, result, CURRENT_TIME,
                            position, cc, fail));
            direct(out, "SignatureAlgorithmCryptographicCheckerResultCheck|cert|" + shape, i18n,
                    result -> new SignatureAlgorithmCryptographicCheckerResultCheck<>(i18n, result, CURRENT_TIME,
                            Context.CERTIFICATE, position, cc, fail, "C-1"));
            direct(out, "DigestAlgorithmCryptographicCheckerResultCheck|" + shape, i18n,
                    result -> new DigestAlgorithmCryptographicCheckerResultCheck<>(i18n, result, CURRENT_TIME,
                            position, cc, fail));
        }
    }

    private static void direct(PrintWriter out, String token, I18nProvider i18n,
                               java.util.function.Function<XmlSAV, ChainItem<XmlSAV>> factory) {
        XmlSAV result;
        try {
            result = new SingleSAVChain(i18n, factory).execute();
        } catch (RuntimeException e) {
            System.err.println("SKIPPED direct " + token + ": " + e);
            return;
        }
        StringBuilder row = new StringBuilder("{");
        key(row, "token").append(str(token)).append(',');
        key(row, "block").append(str("direct")).append(',');
        body(row, result);
        out.println(row.append('}').toString());
    }

    private static XmlAOV aovOfShape(String shape) {
        XmlAOV aov = new XmlAOV();
        XmlConclusion conclusion = new XmlConclusion();
        switch (shape) {
            case "passed":
                conclusion.setIndication(Indication.PASSED);
                break;
            case "passed-with-algo":
                conclusion.setIndication(Indication.PASSED);
                XmlCryptographicValidation cv = new XmlCryptographicValidation();
                XmlCryptographicAlgorithm algorithm = new XmlCryptographicAlgorithm();
                algorithm.setName("RSA");
                algorithm.setUri("http://www.w3.org/2001/04/xmldsig-more#rsa-sha256");
                algorithm.setKeyLength("2048");
                cv.setAlgorithm(algorithm);
                XmlConclusion cvConclusion = new XmlConclusion();
                cvConclusion.setIndication(Indication.PASSED);
                cv.setConclusion(cvConclusion);
                aov.setSignatureCryptographicValidation(cv);
                break;
            case "error":
                conclusion.setIndication(Indication.INDETERMINATE);
                conclusion.setSubIndication(SubIndication.CRYPTO_CONSTRAINTS_FAILURE_NO_POE);
                conclusion.getErrors().add(message("BBB_SAV_ASCCM_ANS", "the algorithm is no longer reliable"));
                break;
            case "warning":
                conclusion.setIndication(Indication.PASSED);
                conclusion.getWarnings().add(message("BBB_SAV_ASCCM_ANS", "the algorithm is about to expire"));
                break;
            case "info":
                conclusion.setIndication(Indication.PASSED);
                conclusion.getInfos().add(message("BBB_SAV_ASCCM_ANS", "informative note"));
                break;
            default:
                throw new IllegalArgumentException(shape);
        }
        aov.setConclusion(conclusion);
        return aov;
    }

    private static XmlCC ccOfShape(String shape) {
        XmlCC cc = new XmlCC();
        XmlConclusion conclusion = new XmlConclusion();
        XmlCryptographicValidation cv = new XmlCryptographicValidation();
        XmlCryptographicAlgorithm algorithm = new XmlCryptographicAlgorithm();
        algorithm.setName("SHA256");
        algorithm.setUri("http://www.w3.org/2001/04/xmlenc#sha256");
        switch (shape) {
            case "passed":
                conclusion.setIndication(Indication.PASSED);
                algorithm.setKeyLength("2048");
                break;
            case "passed-no-keylength":
                conclusion.setIndication(Indication.PASSED);
                break;
            case "failed":
                conclusion.setIndication(Indication.INDETERMINATE);
                conclusion.setSubIndication(SubIndication.CRYPTO_CONSTRAINTS_FAILURE_NO_POE);
                conclusion.getErrors().add(message("ASCCM_AR_ANS_ANR", "SHA1 is not reliable"));
                break;
            default:
                throw new IllegalArgumentException(shape);
        }
        cv.setAlgorithm(algorithm);
        XmlConclusion cvConclusion = new XmlConclusion();
        cvConclusion.setIndication(conclusion.getIndication());
        cvConclusion.setSubIndication(conclusion.getSubIndication());
        cvConclusion.getErrors().addAll(conclusion.getErrors());
        cv.setConclusion(cvConclusion);
        cc.setCryptographicValidation(cv);
        cc.setConclusion(conclusion);
        return cc;
    }

    private static XmlMessage message(String key, String value) {
        XmlMessage message = new XmlMessage();
        message.setKey(key);
        message.setValue(value);
        return message;
    }

    // --------------------------------------------------------------------- JSON

    /** title + conclusion + constraints, without the enclosing braces. */
    private static void body(StringBuilder out, XmlConstraintsConclusion result) {
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
    }

    private static String cv(XmlCryptographicValidation validation) {
        if (validation == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "algorithm").append(algorithm(validation.getAlgorithm())).append(',');
        key(out, "notAfter").append(str(iso(validation.getNotAfter()))).append(',');
        key(out, "concernedMaterialDescription").append(str(validation.getConcernedMaterialDescription())).append(',');
        key(out, "tokenId").append(str(validation.getTokenId())).append(',');
        key(out, "conclusion").append(conclusion(validation.getConclusion()));
        return out.append('}').toString();
    }

    private static String algorithm(XmlCryptographicAlgorithm algorithm) {
        if (algorithm == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("{");
        key(out, "name").append(str(algorithm.getName())).append(',');
        key(out, "uri").append(str(algorithm.getUri())).append(',');
        key(out, "keyLength").append(str(algorithm.getKeyLength()));
        return out.append('}').toString();
    }

    private static String chain(XmlCertificateChainCryptographicValidation validation) {
        if (validation == null) {
            return "null";
        }
        StringBuilder out = new StringBuilder("[");
        List<XmlCryptographicValidation> items = validation.getCertificateCryptographicValidation();
        for (int i = 0; i < items.size(); i++) {
            if (i > 0) {
                out.append(',');
            }
            out.append(cv(items.get(i)));
        }
        return out.append(']').toString();
    }

    private static String iso(Date date) {
        if (date == null) {
            return null;
        }
        SimpleDateFormat format = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'");
        format.setTimeZone(TimeZone.getTimeZone("UTC"));
        return format.format(date);
    }

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
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

}
