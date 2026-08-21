/*
 * Java oracle driver for the isc / vci / cv building blocks.
 *
 * It runs the UPSTREAM building blocks
 *
 *   eu.europa.esig.dss.validation.process.bbb.isc.IdentificationOfTheSigningCertificate
 *   eu.europa.esig.dss.validation.process.bbb.vci.ValidationContextInitialization
 *   eu.europa.esig.dss.validation.process.bbb.cv.CryptographicVerification
 *
 * over REAL inputs - every XmlDiagnosticData dump of the marshal-parity corpus in
 * dss/diagnostic/jaxb/testdata/oracle - with the default ETSI validation policy,
 * and dumps each produced XmlISC / XmlVCI / XmlCV as one JSON object per line.
 * The Go tests in ../../{isc,vci,cv} replay the same inputs against the Go port
 * and compare.
 *
 * A fourth file (cv_direct.jsonl) covers the three checks of the cv package that
 * CryptographicVerification itself never wires (they belong to the evidence
 * record / archival blocks of later phases): AtLeastOneReferenceDataObjectFound,
 * ReferenceDataGroup and SignatureIntactWithId. Those are driven directly, over
 * synthetic digest matchers and over the corpus' first signature, through a
 * single-item chain.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"   # mvn dependency:build-classpath -pl dss-validation
 *   javac -cp "$CP" -d /tmp/oracle BbbBlocksOracle.java
 *   java  -cp "$CP:/tmp/oracle" BbbBlocksOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlCV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlISC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlChainItem;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlISC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlVCI;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.TokenProxy;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlEvidenceRecord;
import eu.europa.esig.dss.diagnostic.jaxb.XmlFoundTimestamp;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTimestamp;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDigestAlgoAndValue;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDigestMatcher;
import eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates;
import eu.europa.esig.dss.diagnostic.jaxb.XmlIssuerSerial;
import eu.europa.esig.dss.diagnostic.jaxb.XmlPolicy;
import eu.europa.esig.dss.diagnostic.jaxb.XmlPolicyDigestAlgoAndValue;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignaturePolicyStore;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.CertificateRefOrigin;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.DigestMatcherType;
import eu.europa.esig.dss.enumerations.EvidenceRecordTimestampType;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.model.policy.MultiValuesRule;
import eu.europa.esig.dss.model.policy.ValidationPolicy;
import eu.europa.esig.dss.policy.EtsiValidationPolicyFactory;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.bbb.cv.CryptographicVerification;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.AtLeastOneReferenceDataObjectFoundCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.ReferenceDataGroupCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.EvidenceRecordHashTreeRenewalTimestampCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.ManifestEntryExistenceCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.ManifestEntryGroupCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.ReferenceDataNameMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.SignatureIntactCheck;
import eu.europa.esig.dss.validation.process.bbb.cv.checks.SignatureIntactWithIdCheck;
import eu.europa.esig.dss.validation.process.bbb.isc.checks.DigestValueMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.isc.checks.DigestValuePresentCheck;
import eu.europa.esig.dss.validation.process.bbb.isc.checks.IssuerSerialMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.isc.checks.SigningCertificateRecognitionCheck;
import eu.europa.esig.dss.validation.process.bbb.vci.checks.SignaturePolicyHashValidCheck;
import eu.europa.esig.dss.validation.process.bbb.vci.checks.SignaturePolicyIdentifiedCheck;
import eu.europa.esig.dss.validation.process.bbb.vci.checks.SignaturePolicyIdentifierCheck;
import eu.europa.esig.dss.validation.process.bbb.vci.checks.SignaturePolicyStoreCheck;
import eu.europa.esig.dss.validation.process.bbb.vci.checks.SignaturePolicyZeroHashCheck;
import eu.europa.esig.dss.validation.process.bbb.isc.IdentificationOfTheSigningCertificate;
import eu.europa.esig.dss.validation.process.bbb.vci.ValidationContextInitialization;

import java.io.File;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.function.BiFunction;

public class BbbBlocksOracle {

    private static final LevelRule FAIL = () -> Level.FAIL;

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        I18nProvider i18nProvider = new I18nProvider();
        ValidationPolicy policy = new EtsiValidationPolicyFactory().loadDefaultValidationPolicy();

        try (PrintWriter isc = writer(repoRoot, "validation/process/bbb/isc/testdata/oracle/isc_blocks.jsonl");
             PrintWriter vci = writer(repoRoot, "validation/process/bbb/vci/testdata/oracle/vci_blocks.jsonl");
             PrintWriter cv = writer(repoRoot, "validation/process/bbb/cv/testdata/oracle/cv_blocks.jsonl");
             PrintWriter direct = writer(repoRoot, "validation/process/bbb/cv/testdata/oracle/cv_direct.jsonl");
             PrintWriter iscDirect = writer(repoRoot, "validation/process/bbb/isc/testdata/oracle/isc_direct.jsonl");
             PrintWriter vciDirect = writer(repoRoot, "validation/process/bbb/vci/testdata/oracle/vci_direct.jsonl")) {

            SignatureWrapper firstSignature = null;
            String firstSignatureFile = null;

            for (File file : files) {
              try {
                XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                DiagnosticData diagnosticData = new DiagnosticData(jaxb);
                String name = file.getName();

                for (SignatureWrapper signature : diagnosticData.getSignatures()) {
                    Context context = signature.isCounterSignature() ? Context.COUNTER_SIGNATURE : Context.SIGNATURE;
                    isc.println(iscRow(name, signature, context,
                            new IdentificationOfTheSigningCertificate(i18nProvider, signature, context, policy).execute()));
                    vci.println(row(name, signature.getId(), context, "VCI",
                            new ValidationContextInitialization(i18nProvider, signature, context, policy).execute()));
                    cv.println(row(name, signature.getId(), context, "CV",
                            new CryptographicVerification(i18nProvider, diagnosticData, signature, context, policy).execute()));
                    if (firstSignature == null) {
                        firstSignature = signature;
                        firstSignatureFile = name;
                    }
                }

                for (TimestampWrapper timestamp : diagnosticData.getTimestampList()) {
                    isc.println(iscRow(name, timestamp, Context.TIMESTAMP,
                            new IdentificationOfTheSigningCertificate(i18nProvider, timestamp, Context.TIMESTAMP, policy).execute()));
                    cv.println(row(name, timestamp.getId(), Context.TIMESTAMP, "CV",
                            new CryptographicVerification(i18nProvider, diagnosticData, timestamp, Context.TIMESTAMP, policy).execute()));
                }

                List<RevocationWrapper> revocations = new ArrayList<>(diagnosticData.getAllRevocationData());
                revocations.sort(Comparator.comparing(RevocationWrapper::getId));
                for (RevocationWrapper revocation : revocations) {
                    isc.println(iscRow(name, revocation, Context.REVOCATION,
                            new IdentificationOfTheSigningCertificate(i18nProvider, revocation, Context.REVOCATION, policy).execute()));
                    cv.println(row(name, revocation.getId(), Context.REVOCATION, "CV",
                            new CryptographicVerification(i18nProvider, diagnosticData, revocation, Context.REVOCATION, policy).execute()));
                }
              } catch (RuntimeException e) {
                // A few dumps of the marshal-parity corpus are schema-coverage
                // fixtures whose wrapper graph is incomplete (e.g. a related
                // certificate reference without the certificate); upstream's own
                // wrappers throw on them. Such a dump carries no oracle rows, and
                // the Go test skips it the same way - it iterates the dumps named
                // by the rows, not the directory.
                System.err.println("SKIPPED " + file.getName() + ": " + e);
              }
            }

            // --- checks of the cv package that CryptographicVerification does not wire

            direct.println(row("synthetic", "at-least-one-found-ok", Context.EVIDENCE_RECORD, "CV",
                    single(i18nProvider, (result, rule) -> new AtLeastOneReferenceDataObjectFoundCheck<>(
                            i18nProvider, result, digestMatchers(true, false), rule))));
            direct.println(row("synthetic", "at-least-one-found-ko", Context.EVIDENCE_RECORD, "CV",
                    single(i18nProvider, (result, rule) -> new AtLeastOneReferenceDataObjectFoundCheck<>(
                            i18nProvider, result, digestMatchers(false, false), rule))));
            direct.println(row("synthetic", "reference-data-group-ok", Context.EVIDENCE_RECORD, "CV",
                    single(i18nProvider, (result, rule) -> new ReferenceDataGroupCheck<>(
                            i18nProvider, result, digestMatchers(true, false), rule))));
            direct.println(row("synthetic", "reference-data-group-ko", Context.EVIDENCE_RECORD, "CV",
                    single(i18nProvider, (result, rule) -> new ReferenceDataGroupCheck<>(
                            i18nProvider, result, digestMatchers(true, true), rule))));

            final SignatureWrapper signature = firstSignature;
            direct.println(row(firstSignatureFile, signature.getId(), Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new SignatureIntactWithIdCheck<>(
                            i18nProvider, result, signature, Context.SIGNATURE, rule))));

            direct.println(row("synthetic", "manifest-entry-existence-ok", Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new ManifestEntryExistenceCheck(
                            i18nProvider, result, manifestEntries(true, true), rule))));
            direct.println(row("synthetic", "manifest-entry-existence-ko", Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new ManifestEntryExistenceCheck(
                            i18nProvider, result, manifestEntries(false, false), rule))));
            direct.println(row("synthetic", "manifest-entry-group-ok", Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new ManifestEntryGroupCheck(
                            i18nProvider, result, manifestEntries(true, true), rule))));
            direct.println(row("synthetic", "manifest-entry-group-ko", Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new ManifestEntryGroupCheck(
                            i18nProvider, result, manifestEntries(true, false), rule))));
            direct.println(row("synthetic", "reference-name-match-ko", Context.SIGNATURE, "CV",
                    single(i18nProvider, (result, rule) -> new ReferenceDataNameMatchCheck<>(
                            i18nProvider, result, manifestEntries(true, false).get(1), rule))));
            for (boolean covered : new boolean[] { true, false }) {
                XmlDiagnosticData synthetic = evidenceRecordDiagnosticData(covered);
                DiagnosticData syntheticData = new DiagnosticData(synthetic);
                TimestampWrapper renewal = syntheticData.getTimestampList().get(0);
                direct.println(row("synthetic",
                        covered ? "er-hash-tree-renewal-ok" : "er-hash-tree-renewal-ko",
                        Context.TIMESTAMP, "CV",
                        single(i18nProvider, (result, rule) -> new EvidenceRecordHashTreeRenewalTimestampCheck(
                                i18nProvider, result, syntheticData, renewal, rule))));
            }

            direct.println(row("synthetic", "signature-intact-certificate-context", Context.CERTIFICATE, "CV",
                    single(i18nProvider, (result, rule) -> new SignatureIntactCheck<>(
                            i18nProvider, result, syntheticSignature(false, false, false, false, true),
                            Context.CERTIFICATE, rule))));

            // --- isc checks, over synthetic signing-certificate references

            iscDirect.println(iscRowDirect("recognition-ok", i18nProvider,
                    syntheticSignature(true, true, true, true, true),
                    (result, token, rule) -> new SigningCertificateRecognitionCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("recognition-ko", i18nProvider,
                    syntheticSignature(true, true, true, true, false),
                    (result, token, rule) -> new SigningCertificateRecognitionCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("digest-value-present-ok", i18nProvider,
                    syntheticSignature(true, true, true, true, true),
                    (result, token, rule) -> new DigestValuePresentCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("digest-value-present-ko", i18nProvider,
                    syntheticSignature(false, false, true, true, true),
                    (result, token, rule) -> new DigestValuePresentCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("digest-value-match-ok", i18nProvider,
                    syntheticSignature(true, true, true, true, true),
                    (result, token, rule) -> new DigestValueMatchCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("digest-value-match-ko", i18nProvider,
                    syntheticSignature(true, false, true, true, true),
                    (result, token, rule) -> new DigestValueMatchCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("issuer-serial-match-ok", i18nProvider,
                    syntheticSignature(true, true, true, true, true),
                    (result, token, rule) -> new IssuerSerialMatchCheck(i18nProvider, result, token, rule)));
            iscDirect.println(iscRowDirect("issuer-serial-match-ko", i18nProvider,
                    syntheticSignature(true, true, true, false, true),
                    (result, token, rule) -> new IssuerSerialMatchCheck(i18nProvider, result, token, rule)));

            // --- vci checks, over synthetic signature policies

            MultiValuesRule anyPolicy = new MultiValuesRule() {
                @Override public Level getLevel() { return Level.FAIL; }
                @Override public List<String> getValues() { return Arrays.asList("ANY_POLICY"); }
            };
            vciDirect.println(vciRowDirect("policy-identifier-ok", i18nProvider,
                    policySignature("1.2.3.4", true, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyIdentifierCheck(i18nProvider, result, sig, anyPolicy)));
            vciDirect.println(vciRowDirect("policy-identifier-ko", i18nProvider,
                    policySignature(null, true, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyIdentifierCheck(i18nProvider, result, sig, anyPolicy)));
            vciDirect.println(vciRowDirect("policy-identified-ok", i18nProvider,
                    policySignature("1.2.3.4", true, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyIdentifiedCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-identified-ko", i18nProvider,
                    policySignature("1.2.3.4", false, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyIdentifiedCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-store-ok", i18nProvider,
                    policySignature("1.2.3.4", true, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyStoreCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-store-ko", i18nProvider,
                    policySignature("1.2.3.4", true, true, false, false),
                    (result, sig, rule) -> new SignaturePolicyStoreCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-hash-valid-ok", i18nProvider,
                    policySignature("1.2.3.4", true, true, false, true),
                    (result, sig, rule) -> new SignaturePolicyHashValidCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-hash-valid-ko", i18nProvider,
                    policySignature("1.2.3.4", true, false, false, true),
                    (result, sig, rule) -> new SignaturePolicyHashValidCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-zero-hash-ok", i18nProvider,
                    policySignature("1.2.3.4", true, false, true, true),
                    (result, sig, rule) -> new SignaturePolicyZeroHashCheck(i18nProvider, result, sig, rule)));
            vciDirect.println(vciRowDirect("policy-zero-hash-ko", i18nProvider,
                    policySignature("1.2.3.4", true, false, false, true),
                    (result, sig, rule) -> new SignaturePolicyZeroHashCheck(i18nProvider, result, sig, rule)));
        }
    }

    /**
     * A diagnostic data holding one evidence record covering "doc.xml" and one
     * HashTree-renewal archive time-stamp which does, or does not, cover it too.
     */
    private static XmlDiagnosticData evidenceRecordDiagnosticData(boolean covered) {
        XmlDigestMatcher erMatcher = new XmlDigestMatcher();
        erMatcher.setType(DigestMatcherType.EVIDENCE_RECORD_ARCHIVE_OBJECT);
        erMatcher.setDocumentName("doc.xml");
        erMatcher.setDataFound(true);
        erMatcher.setDataIntact(true);

        XmlDigestMatcher tstMatcher = new XmlDigestMatcher();
        tstMatcher.setType(DigestMatcherType.EVIDENCE_RECORD_ARCHIVE_OBJECT);
        tstMatcher.setDocumentName(covered ? "doc.xml" : "other.xml");
        tstMatcher.setDataFound(true);
        tstMatcher.setDataIntact(true);

        XmlTimestamp timestamp = new XmlTimestamp();
        timestamp.setId("T-SYNTHETIC");
        timestamp.setType(TimestampType.EVIDENCE_RECORD_TIMESTAMP);
        timestamp.setEvidenceRecordTimestampType(EvidenceRecordTimestampType.HASH_TREE_RENEWAL_ARCHIVE_TIMESTAMP);
        timestamp.getDigestMatchers().add(tstMatcher);

        XmlFoundTimestamp foundTimestamp = new XmlFoundTimestamp();
        foundTimestamp.setTimestamp(timestamp);

        XmlEvidenceRecord evidenceRecord = new XmlEvidenceRecord();
        evidenceRecord.setId("ER-SYNTHETIC");
        evidenceRecord.getDigestMatchers().add(erMatcher);
        evidenceRecord.getEvidenceRecordTimestamps().add(foundTimestamp);

        XmlDiagnosticData diagnosticData = new XmlDiagnosticData();
        diagnosticData.getEvidenceRecords().add(evidenceRecord);
        diagnosticData.getUsedTimestamps().add(timestamp);
        return diagnosticData;
    }

    /** A chain of exactly one item over an XmlISC. */
    static class SingleISCChain extends Chain<XmlISC> {
        private final BiFunction<XmlISC, LevelRule, ChainItem<XmlISC>> factory;

        SingleISCChain(I18nProvider i18nProvider, BiFunction<XmlISC, LevelRule, ChainItem<XmlISC>> factory) {
            super(i18nProvider, new XmlISC());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    /** A chain of exactly one item over an XmlVCI. */
    static class SingleVCIChain extends Chain<XmlVCI> {
        private final BiFunction<XmlVCI, LevelRule, ChainItem<XmlVCI>> factory;

        SingleVCIChain(I18nProvider i18nProvider, BiFunction<XmlVCI, LevelRule, ChainItem<XmlVCI>> factory) {
            super(i18nProvider, new XmlVCI());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    interface ISCCheckFactory {
        ChainItem<XmlISC> create(XmlISC result, TokenProxy token, LevelRule rule);
    }

    interface VCICheckFactory {
        ChainItem<XmlVCI> create(XmlVCI result, SignatureWrapper signature, LevelRule rule);
    }

    private static String iscRowDirect(String name, I18nProvider i18nProvider, SignatureWrapper token,
                                       ISCCheckFactory factory) {
        XmlISC result = new SingleISCChain(i18nProvider, (r, rule) -> factory.create(r, token, rule)).execute();
        return row("synthetic", name, Context.SIGNATURE, "ISC", result);
    }

    private static String vciRowDirect(String name, I18nProvider i18nProvider, SignatureWrapper signature,
                                       VCICheckFactory factory) {
        XmlVCI result = new SingleVCIChain(i18nProvider, (r, rule) -> factory.create(r, signature, rule)).execute();
        return row("synthetic", name, Context.SIGNATURE, "VCI", result);
    }

    /** A signature carrying one signing-certificate reference with the requested flags. */
    private static SignatureWrapper syntheticSignature(boolean digestPresent, boolean digestMatch,
                                                       boolean issuerSerialPresent, boolean issuerSerialMatch,
                                                       boolean withSigningCertificate) {
        XmlSignature xmlSignature = new XmlSignature();
        xmlSignature.setId("S-SYNTHETIC");

        XmlCertificate xmlCertificate = new XmlCertificate();
        xmlCertificate.setId("C-SYNTHETIC");

        XmlCertificateRef ref = new XmlCertificateRef();
        ref.setOrigin(CertificateRefOrigin.SIGNING_CERTIFICATE);
        if (digestPresent) {
            XmlDigestAlgoAndValue digest = new XmlDigestAlgoAndValue();
            digest.setDigestMethod(DigestAlgorithm.SHA256);
            digest.setDigestValue(new byte[] { 1, 2, 3 });
            digest.setMatch(digestMatch);
            ref.setDigestAlgoAndValue(digest);
        }
        if (issuerSerialPresent) {
            XmlIssuerSerial issuerSerial = new XmlIssuerSerial();
            issuerSerial.setValue(new byte[] { 4, 5, 6 });
            issuerSerial.setMatch(issuerSerialMatch);
            ref.setIssuerSerial(issuerSerial);
        }

        XmlRelatedCertificate related = new XmlRelatedCertificate();
        related.setCertificate(xmlCertificate);
        related.getCertificateRefs().add(ref);

        XmlFoundCertificates foundCertificates = new XmlFoundCertificates();
        foundCertificates.getRelatedCertificates().add(related);
        xmlSignature.setFoundCertificates(foundCertificates);

        if (withSigningCertificate) {
            XmlSigningCertificate signingCertificate = new XmlSigningCertificate();
            signingCertificate.setCertificate(xmlCertificate);
            xmlSignature.setSigningCertificate(signingCertificate);
        }
        return new SignatureWrapper(xmlSignature);
    }

    /** A signature carrying a signature policy with the requested flags. */
    private static SignatureWrapper policySignature(String policyId, Boolean identified, Boolean digestMatch,
                                                    Boolean zeroHash, boolean storePresent) {
        XmlSignature xmlSignature = new XmlSignature();
        xmlSignature.setId("S-SYNTHETIC-POLICY");

        XmlPolicy policy = new XmlPolicy();
        policy.setId(policyId);
        policy.setIdentified(identified);
        XmlPolicyDigestAlgoAndValue digest = new XmlPolicyDigestAlgoAndValue();
        digest.setDigestMethod(DigestAlgorithm.SHA256);
        digest.setDigestValue(new byte[] { 7, 8, 9 });
        digest.setMatch(digestMatch);
        digest.setZeroHash(zeroHash);
        policy.setDigestAlgoAndValue(digest);
        xmlSignature.setPolicy(policy);

        if (storePresent) {
            XmlSignaturePolicyStore store = new XmlSignaturePolicyStore();
            store.setId("SPS-SYNTHETIC");
            xmlSignature.setSignaturePolicyStore(store);
        }
        return new SignatureWrapper(xmlSignature);
    }

    /** Two MANIFEST_ENTRY matchers, the second one optionally not found. */
    private static List<XmlDigestMatcher> manifestEntries(boolean firstFound, boolean secondFound) {
        List<XmlDigestMatcher> digestMatchers = new ArrayList<>();
        XmlDigestMatcher first = new XmlDigestMatcher();
        first.setType(DigestMatcherType.MANIFEST_ENTRY);
        first.setUri("doc.xml");
        first.setDocumentName("doc.xml");
        first.setDataFound(firstFound);
        first.setDataIntact(firstFound);
        digestMatchers.add(first);
        XmlDigestMatcher second = new XmlDigestMatcher();
        second.setType(DigestMatcherType.MANIFEST_ENTRY);
        second.setUri("other.xml");
        second.setDocumentName("renamed.xml");
        second.setDataFound(secondFound);
        second.setDataIntact(secondFound);
        digestMatchers.add(second);
        return digestMatchers;
    }

    /** A chain of exactly one item, used to drive a check in isolation. */
    static class SingleCheckChain extends Chain<XmlCV> {
        private final BiFunction<XmlCV, LevelRule, ChainItem<XmlCV>> factory;

        SingleCheckChain(I18nProvider i18nProvider, BiFunction<XmlCV, LevelRule, ChainItem<XmlCV>> factory) {
            super(i18nProvider, new XmlCV());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    private static XmlCV single(I18nProvider i18nProvider, BiFunction<XmlCV, LevelRule, ChainItem<XmlCV>> factory) {
        return new SingleCheckChain(i18nProvider, factory).execute();
    }

    /** Two EVIDENCE_RECORD_ARCHIVE_OBJECT matchers plus, optionally, an orphan reference. */
    private static List<XmlDigestMatcher> digestMatchers(boolean dataFound, boolean withOrphan) {
        List<XmlDigestMatcher> digestMatchers = new ArrayList<>();
        XmlDigestMatcher first = new XmlDigestMatcher();
        first.setType(DigestMatcherType.EVIDENCE_RECORD_ARCHIVE_OBJECT);
        first.setDocumentName("doc.xml");
        first.setDataFound(dataFound);
        first.setDataIntact(dataFound);
        digestMatchers.add(first);
        XmlDigestMatcher second = new XmlDigestMatcher();
        second.setType(DigestMatcherType.EVIDENCE_RECORD_ARCHIVE_OBJECT);
        second.setDocumentName("other.xml");
        second.setDataFound(false);
        second.setDataIntact(false);
        digestMatchers.add(second);
        if (withOrphan) {
            XmlDigestMatcher orphan = new XmlDigestMatcher();
            orphan.setType(DigestMatcherType.EVIDENCE_RECORD_ORPHAN_REFERENCE);
            orphan.setDataFound(true);
            orphan.setDataIntact(true);
            digestMatchers.add(orphan);
        }
        return digestMatchers;
    }

    // ------------------------------------------------------------------- JSON

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static String iscRow(String file, TokenProxy token, Context context, XmlISC result) {
        StringBuilder out = new StringBuilder(row(file, token.getId(), context, "ISC", result));
        out.setLength(out.length() - 1); // drop the closing brace
        out.append(',');
        key(out, "certificateChain");
        if (result.getCertificateChain() == null) {
            out.append("null");
        } else {
            out.append('[');
            List<XmlChainItem> chainItems = result.getCertificateChain().getChainItem();
            for (int i = 0; i < chainItems.size(); i++) {
                if (i > 0) {
                    out.append(',');
                }
                out.append('{');
                key(out, "id").append(str(chainItems.get(i).getId())).append(',');
                key(out, "source").append(str(chainItems.get(i).getSource() == null
                        ? null : chainItems.get(i).getSource().name()));
                out.append('}');
            }
            out.append(']');
        }
        return out.append('}').toString();
    }

    private static String row(String file, String tokenId, Context context, String block,
                              XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(tokenId)).append(',');
        key(out, "context").append(str(context.name())).append(',');
        key(out, "block").append(str(block)).append(',');
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

}
