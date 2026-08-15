/*
 * Java oracle driver for the fc / sav check classes the top-level building blocks
 * of FcSavOracle do not exercise under the default ETSI validation policy - either
 * because the policy leaves their constraint undefined (Level null, the item is
 * skipped) or because the chain only wires them for a signature form the 50-dump
 * marshal-parity corpus does not contain.
 *
 * Every such check is driven directly through a chain of exactly one item at
 * Level.FAIL, over REAL corpus inputs: each dump's signatures and time-stamps in
 * turn. Driving the same check over every token is what produces both the happy
 * and the failure row wherever the corpus varies; a few checks that read a value
 * no dump carries are additionally driven over synthetic literals.
 *
 * Rows are written to
 *   dss/validation/process/bbb/fc/testdata/oracle/fc_direct.jsonl
 *   dss/validation/process/bbb/sav/testdata/oracle/sav_direct.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle FcSavDirectOracle.java
 *   java  -cp "$CP:/tmp/oracle" FcSavDirectOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlBasicBuildingBlocks;
import eu.europa.esig.dss.detailedreport.jaxb.XmlAOV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicAlgorithm;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicValidation;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlFC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSAV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlTimestamp;
import eu.europa.esig.dss.detailedreport.jaxb.XmlValidationProcessBasicTimestamp;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.PDFRevisionWrapper;
import eu.europa.esig.dss.diagnostic.SignatureWrapper;
import eu.europa.esig.dss.diagnostic.TimestampWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.model.policy.MultiValuesRule;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.AcceptableZipCommentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.AnnotationChangesCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ByteRangeAllDocumentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.DocMDPCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.EllipticCurveKeySizeCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.FormFillChangesCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.FullScopeCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.PDFAComplianceCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.PDFAProfileCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ZipCommentPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ByteRangeCollisionCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.CAdESV3HashIndexCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.PdfVisualDifferenceCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.SignedFilesPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.KeyIdentifierMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.MessageDigestOrSignedPropertiesCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SigningCertificateAttributePresentCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SigningCertificateReferencesValidityCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SigningTimeInCertificateValidityRangeCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.StructuralValidationCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.TSAGeneralNameValueMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.UnicitySigningCertificateAttributeCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.AcceptableMimetypeFileContentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.AllFilesSignedCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ContainerTypeCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.FieldMDPCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.FormatCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ManifestFilePresentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.MimeTypeFilePresentCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.PdfAnnotationOverlapCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.PdfSignatureDictionaryCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.ReferencesNotAmbiguousCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.SigFieldLockCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.SignatureNotAmbiguousCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.SignerInformationStoreCheck;
import eu.europa.esig.dss.validation.process.bbb.fc.checks.UndefinedChangesCheck;
import eu.europa.esig.dss.diagnostic.jaxb.XmlContainerInfo;

import eu.europa.esig.dss.validation.process.bbb.sav.checks.AllCertificatesInPathReferencedCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ArchiveTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.CertifiedRolesCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ClaimedRolesCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.CommitmentTypeIndicationsCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ContentHintsCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ContentIdentifierCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ContentTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ContentTimestampBasicValidationCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ContentTypeCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.CounterSignatureCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.DocumentTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.KeyIdentifierPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.LTALevelTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SignatureTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SignatureTypeCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.SignerLocationCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.TLevelTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.TSAGeneralNameFieldPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.TSAGeneralNameOrderMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ValidationDataRefsOnlyTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.ValidationDataTimeStampCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.X509UrlMatchCheck;
import eu.europa.esig.dss.validation.process.bbb.sav.checks.X509UrlPresentCheck;
import eu.europa.esig.dss.validation.process.vpfltvd.checks.TimestampMessageImprintWithIdCheck;
import eu.europa.esig.dss.validation.process.bbb.aov.checks.AlgorithmObsolescenceValidationCheck;
import eu.europa.esig.dss.i18n.MessageTag;

import java.io.File;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.function.BiFunction;

public class FcSavDirectOracle {

    private static final LevelRule FAIL = () -> Level.FAIL;

    /** A FAIL MultiValuesRule accepting any present value. */
    private static final MultiValuesRule ANY = new MultiValuesRule() {
        @Override public Level getLevel() { return Level.FAIL; }
        @Override public List<String> getValues() { return Collections.singletonList("*"); }
    };

    /** A FAIL MultiValuesRule accepting nothing the corpus carries. */
    private static final MultiValuesRule NONE = new MultiValuesRule() {
        @Override public Level getLevel() { return Level.FAIL; }
        @Override public List<String> getValues() { return Collections.singletonList("NO-SUCH-VALUE"); }
    };

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();

        try (PrintWriter fc = writer(repoRoot, "validation/process/bbb/fc/testdata/oracle/fc_direct.jsonl");
             PrintWriter sav = writer(repoRoot, "validation/process/bbb/sav/testdata/oracle/sav_direct.jsonl")) {

            for (File file : files) {
                DiagnosticData dd;
                String name = file.getName();
                // The four model-*.xml dumps are schema-coverage fixtures, not validation
                // output: every property of the generated model is filled by reflection from a
                // counter. Their IDREF graph is dangling (e.g. <SigningCertificate> pointing at
                // no <Certificate>) and their attribute values carry raw control characters,
                // and the two runtimes resolve those edges differently - upstream's
                // getSigningCertificate() answers null where the Go wrapper resolves an object,
                // and the JCA parser normalises a literal TAB in an attribute value where
                // encoding/xml keeps it. Rows taken over them would compare two different
                // inputs, so the direct corpora skip them; the block corpora skip them already
                // (upstream's own wrappers throw). Every row here comes from a real dump or
                // from an explicitly built synthetic below.
                if (name.startsWith("model-")) {
                    System.err.println("SKIPPED-FIXTURE " + name);
                    continue;
                }
                try {
                    XmlDiagnosticData jaxb = DiagnosticDataFacade.newFacade().unmarshall(file, false);
                    dd = new DiagnosticData(jaxb);
                    dd.getSignatures();
                    dd.getTimestampList();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED-FILE " + name + ": " + e);
                    continue;
                }
                final DiagnosticData diagnosticData = dd;

                for (SignatureWrapper s : diagnosticData.getSignatures()) {
                    final SignatureWrapper signature = s;
                    final PDFRevisionWrapper rev = signature.getPDFRevision();

                    // --- fc, over the signature
                    fcRow(fc, i18n, name, signature.getId(), "EllipticCurveKeySizeCheck",
                            (r, rule) -> new EllipticCurveKeySizeCheck(i18n, r, signature, rule));
                    fcRow(fc, i18n, name, signature.getId(), "FullScopeCheck",
                            (r, rule) -> new FullScopeCheck(i18n, r, signature.getSignatureScopes(), rule));
                    fcRow(fc, i18n, name, signature.getId(), "ByteRangeAllDocumentCheck",
                            (r, rule) -> new ByteRangeAllDocumentCheck(i18n, r, diagnosticData, rule));
                    if (rev != null) {
                        fcRow(fc, i18n, name, signature.getId(), "AnnotationChangesCheck",
                                (r, rule) -> new AnnotationChangesCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "DocMDPCheck",
                                (r, rule) -> new DocMDPCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "FormFillChangesCheck",
                                (r, rule) -> new FormFillChangesCheck(i18n, r, rev, rule));
                    }

                    // --- fc checks the chain does run, but which the corpus never fails:
                    // driven alone so that both branches are recorded.
                    fcRow(fc, i18n, name, signature.getId(), "ReferencesNotAmbiguousCheck",
                            (r, rule) -> new ReferencesNotAmbiguousCheck(i18n, r, signature, rule));
                    fcRow(fc, i18n, name, signature.getId(), "SignatureNotAmbiguousCheck",
                            (r, rule) -> new SignatureNotAmbiguousCheck(i18n, r, signature, rule));
                    fcRow(fc, i18n, name, signature.getId(), "SignerInformationStoreCheck",
                            (r, rule) -> new SignerInformationStoreCheck(i18n, r, signature, rule));
                    fcRowV(fc, i18n, name, signature.getId(), "FormatCheck-any",
                            (r, rule) -> new FormatCheck(i18n, r, signature, ANY));
                    fcRowV(fc, i18n, name, signature.getId(), "FormatCheck-none",
                            (r, rule) -> new FormatCheck(i18n, r, signature, NONE));
                    if (rev != null) {
                        fcRow(fc, i18n, name, signature.getId(), "PdfSignatureDictionaryCheck",
                                (r, rule) -> new PdfSignatureDictionaryCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "PdfAnnotationOverlapCheck",
                                (r, rule) -> new PdfAnnotationOverlapCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "UndefinedChangesCheck",
                                (r, rule) -> new UndefinedChangesCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "SigFieldLockCheck",
                                (r, rule) -> new SigFieldLockCheck(i18n, r, rev, rule));
                        fcRow(fc, i18n, name, signature.getId(), "FieldMDPCheck",
                                (r, rule) -> new FieldMDPCheck(i18n, r, rev, rule));
                    }
                    final XmlContainerInfo containerInfo = diagnosticData.getContainerInfo();
                    if (containerInfo != null) {
                        fcRow(fc, i18n, name, signature.getId(), "AllFilesSignedCheck",
                                (r, rule) -> new AllFilesSignedCheck(i18n, r, signature, containerInfo, rule));
                        fcRow(fc, i18n, name, signature.getId(), "SignedFilesPresentCheck",
                                (r, rule) -> new SignedFilesPresentCheck(i18n, r, containerInfo, rule));
                        fcRow(fc, i18n, name, signature.getId(), "ManifestFilePresentCheck",
                                (r, rule) -> new ManifestFilePresentCheck(i18n, r, containerInfo, rule));
                        fcRowV(fc, i18n, name, signature.getId(), "ContainerTypeCheck-any",
                                (r, rule) -> new ContainerTypeCheck(i18n, r, containerInfo.getContainerType(), ANY));
                        fcRowV(fc, i18n, name, signature.getId(), "ContainerTypeCheck-none",
                                (r, rule) -> new ContainerTypeCheck(i18n, r, containerInfo.getContainerType(), NONE));
                        fcRowV(fc, i18n, name, signature.getId(), "AcceptableMimetypeFileContentCheck-any",
                                (r, rule) -> new AcceptableMimetypeFileContentCheck(
                                        i18n, r, containerInfo.getMimeTypeContent(), ANY));
                        fcRowV(fc, i18n, name, signature.getId(), "AcceptableMimetypeFileContentCheck-none",
                                (r, rule) -> new AcceptableMimetypeFileContentCheck(
                                        i18n, r, containerInfo.getMimeTypeContent(), NONE));
                    }

                    fcRow(fc, i18n, name, signature.getId(), "ByteRangeCollisionCheck",
                            (r, rule) -> new ByteRangeCollisionCheck(i18n, r, signature, diagnosticData, rule));
                    if (rev != null) {
                        fcRow(fc, i18n, name, signature.getId(), "PdfVisualDifferenceCheck",
                                (r, rule) -> new PdfVisualDifferenceCheck(i18n, r, rev, rule));
                    }

                    // --- sav, over the signature
                    savRow(sav, i18n, name, signature.getId(), "KeyIdentifierMatchCheck",
                            (r, rule) -> new KeyIdentifierMatchCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "MessageDigestOrSignedPropertiesCheck",
                            (r, rule) -> new MessageDigestOrSignedPropertiesCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "SigningCertificateAttributePresentCheck",
                            (r, rule) -> new SigningCertificateAttributePresentCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "SigningCertificateReferencesValidityCheck",
                            (r, rule) -> new SigningCertificateReferencesValidityCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "SigningTimeInCertificateValidityRangeCheck",
                            (r, rule) -> new SigningTimeInCertificateValidityRangeCheck<>(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "StructuralValidationCheck",
                            (r, rule) -> new StructuralValidationCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "UnicitySigningCertificateAttributeCheck",
                            (r, rule) -> new UnicitySigningCertificateAttributeCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "AllCertificatesInPathReferencedCheck",
                            (r, rule) -> new AllCertificatesInPathReferencedCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "ArchiveTimeStampCheck",
                            (r, rule) -> new ArchiveTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "ContentTimeStampCheck",
                            (r, rule) -> new ContentTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "CounterSignatureCheck",
                            (r, rule) -> new CounterSignatureCheck(i18n, r, diagnosticData, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "DocumentTimeStampCheck",
                            (r, rule) -> new DocumentTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "KeyIdentifierPresentCheck",
                            (r, rule) -> new KeyIdentifierPresentCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "SignatureTimeStampCheck",
                            (r, rule) -> new SignatureTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "SignerLocationCheck",
                            (r, rule) -> new SignerLocationCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "ValidationDataRefsOnlyTimeStampCheck",
                            (r, rule) -> new ValidationDataRefsOnlyTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "ValidationDataTimeStampCheck",
                            (r, rule) -> new ValidationDataTimeStampCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "X509UrlPresentCheck",
                            (r, rule) -> new X509UrlPresentCheck(i18n, r, signature, rule));
                    savRow(sav, i18n, name, signature.getId(), "X509UrlMatchCheck",
                            (r, rule) -> new X509UrlMatchCheck(i18n, r, signature, rule));
                    // T/LTA-level presence: the detailed-report XmlTimestamps these read
                    // are an 8d product, so they are synthesised here at both outcomes.
                    for (Indication indication : new Indication[] { Indication.PASSED, Indication.FAILED }) {
                        final List<XmlTimestamp> reportTimestamps = reportTimestamps(diagnosticData, indication);
                        savRow(sav, i18n, name, signature.getId() + "/" + indication.name(), "TLevelTimeStampCheck",
                                (r, rule) -> new TLevelTimeStampCheck<>(i18n, r, signature,
                                        Collections.<String, XmlBasicBuildingBlocks>emptyMap(), reportTimestamps, rule));
                        savRow(sav, i18n, name, signature.getId() + "/" + indication.name(), "LTALevelTimeStampCheck",
                                (r, rule) -> new LTALevelTimeStampCheck<>(i18n, r, signature,
                                        Collections.<String, XmlBasicBuildingBlocks>emptyMap(), reportTimestamps, rule));
                    }

                    // multi-values checks: the "*" rule (accept whatever is there) and a
                    // rule accepting nothing, so both branches are recorded
                    for (String variant : new String[] { "any", "none" }) {
                        final MultiValuesRule values = "any".equals(variant) ? ANY : NONE;
                        savRowV(sav, i18n, name, signature.getId(), "CertifiedRolesCheck-" + variant,
                                (r, rule) -> new CertifiedRolesCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "ClaimedRolesCheck-" + variant,
                                (r, rule) -> new ClaimedRolesCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "CommitmentTypeIndicationsCheck-" + variant,
                                (r, rule) -> new CommitmentTypeIndicationsCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "ContentHintsCheck-" + variant,
                                (r, rule) -> new ContentHintsCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "ContentIdentifierCheck-" + variant,
                                (r, rule) -> new ContentIdentifierCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "ContentTypeCheck-" + variant,
                                (r, rule) -> new ContentTypeCheck(i18n, r, signature, values));
                        savRowV(sav, i18n, name, signature.getId(), "SignatureTypeCheck-" + variant,
                                (r, rule) -> new SignatureTypeCheck(i18n, r, signature, values));
                    }

                    List<TimestampWrapper> contentTimestamps;
                    try {
                        contentTimestamps = signature.getContentTimestamps();
                    } catch (RuntimeException e) {
                        System.err.println("SKIPPED " + name + " " + signature.getId() + " contentTimestamps: " + e);
                        contentTimestamps = Collections.emptyList();
                    }
                    for (TimestampWrapper contentTimestamp : contentTimestamps) {
                        final TimestampWrapper ct = contentTimestamp;
                        for (Indication indication : new Indication[] { Indication.PASSED, Indication.FAILED }) {
                            final XmlConclusion conclusion = new XmlConclusion();
                            conclusion.setIndication(indication);
                            savRow(sav, i18n, name, signature.getId() + "/" + ct.getId() + "/" + indication.name(),
                                    "ContentTimestampBasicValidationCheck",
                                    (r, rule) -> new ContentTimestampBasicValidationCheck(i18n, r, ct, conclusion, rule));
                        }
                    }
                }

                List<TimestampWrapper> allTimestamps;
                try {
                    allTimestamps = diagnosticData.getTimestampList();
                } catch (RuntimeException e) {
                    System.err.println("SKIPPED " + name + " timestampList: " + e);
                    allTimestamps = Collections.emptyList();
                }
                for (TimestampWrapper t : allTimestamps) {
                    final TimestampWrapper timestamp = t;
                    // vpfltvd/vpftspwatsp: SignatureAcceptanceValidation wires
                    // TimestampMessageImprintWithIdCheck for every content time-stamp, but the
                    // default policy leaves its constraint undefined, so the corpus never runs it.
                    savRow(sav, i18n, name, timestamp.getId(), "TimestampMessageImprintWithIdCheck",
                            (r, rule) -> new TimestampMessageImprintWithIdCheck<>(i18n, r, timestamp, rule));
                    fcRow(fc, i18n, name, timestamp.getId(), "CAdESV3HashIndexCheck",
                            (r, rule) -> new CAdESV3HashIndexCheck(i18n, r, timestamp, rule));
                    savRow(sav, i18n, name, timestamp.getId(), "TSAGeneralNameValueMatchCheck",
                            (r, rule) -> new TSAGeneralNameValueMatchCheck(i18n, r, timestamp, rule));
                    savRow(sav, i18n, name, timestamp.getId(), "TSAGeneralNameFieldPresentCheck",
                            (r, rule) -> new TSAGeneralNameFieldPresentCheck(i18n, r, timestamp, rule));
                    savRow(sav, i18n, name, timestamp.getId(), "TSAGeneralNameOrderMatchCheck",
                            (r, rule) -> new TSAGeneralNameOrderMatchCheck(i18n, r, timestamp, rule));
                }
            }

            // --- fc checks reading a literal no dump carries
            for (String comment : new String[] { null, "", "mimetype=application/vnd.etsi.asic-e+zip" }) {
                final String zipComment = comment;
                String label = comment == null ? "null" : (comment.isEmpty() ? "empty" : "asice");
                fcRow(fc, i18n, "synthetic", "zip-comment-" + label, "ZipCommentPresentCheck",
                        (r, rule) -> new ZipCommentPresentCheck(i18n, r, zipComment, rule));
                fcRowV(fc, i18n, "synthetic", "zip-comment-" + label + "-any", "AcceptableZipCommentCheck",
                        (r, rule) -> new AcceptableZipCommentCheck(i18n, r, zipComment, ANY));
                fcRowV(fc, i18n, "synthetic", "zip-comment-" + label + "-none", "AcceptableZipCommentCheck",
                        (r, rule) -> new AcceptableZipCommentCheck(i18n, r, zipComment, NONE));
            }
            // --- the branches no corpus dump reaches
            {
                DiagnosticData badByteRangeData = failingByteRangeDiagnosticData();
                fcRow(fc, i18n, "synthetic", "byte-range-invalid", "ByteRangeAllDocumentCheck",
                        (r, rule) -> new ByteRangeAllDocumentCheck(i18n, r, badByteRangeData, rule));

                final SignatureWrapper duplicatedRef = signatureWithDuplicatedReference();
                fcRow(fc, i18n, "synthetic", "duplicated-reference", "ReferencesNotAmbiguousCheck",
                        (r, rule) -> new ReferencesNotAmbiguousCheck(i18n, r, duplicatedRef, rule));

                final SignatureWrapper duplicatedSig = signatureWithDuplicatedSignature();
                fcRow(fc, i18n, "synthetic", "duplicated-signature", "SignatureNotAmbiguousCheck",
                        (r, rule) -> new SignatureNotAmbiguousCheck(i18n, r, duplicatedSig, rule));

                final PDFRevisionWrapper failingRevision = failingPdfRevision();
                fcRow(fc, i18n, "synthetic", "failing-revision", "DocMDPCheck",
                        (r, rule) -> new DocMDPCheck(i18n, r, failingRevision, rule));
                fcRow(fc, i18n, "synthetic", "failing-revision", "SigFieldLockCheck",
                        (r, rule) -> new SigFieldLockCheck(i18n, r, failingRevision, rule));
                fcRow(fc, i18n, "synthetic", "failing-revision", "FieldMDPCheck",
                        (r, rule) -> new FieldMDPCheck(i18n, r, failingRevision, rule));
                fcRow(fc, i18n, "synthetic", "failing-revision", "PdfSignatureDictionaryCheck",
                        (r, rule) -> new PdfSignatureDictionaryCheck(i18n, r, failingRevision, rule));
                fcRow(fc, i18n, "synthetic", "failing-revision", "PdfAnnotationOverlapCheck",
                        (r, rule) -> new PdfAnnotationOverlapCheck(i18n, r, failingRevision, rule));

                final SignatureWrapper contentAttrs = signatureWithContentAttributes();
                savRowV(sav, i18n, "synthetic", "content-attrs", "ContentHintsCheck-any",
                        (r, rule) -> new ContentHintsCheck(i18n, r, contentAttrs, ANY));
                savRowV(sav, i18n, "synthetic", "content-attrs", "ContentIdentifierCheck-any",
                        (r, rule) -> new ContentIdentifierCheck(i18n, r, contentAttrs, ANY));
            }
            {
                XmlContainerInfo emptyContainer = new XmlContainerInfo();
                emptyContainer.setContainerType(eu.europa.esig.dss.enumerations.ASiCContainerType.ASiC_E);
                fcRow(fc, i18n, "synthetic", "no-content-files", "SignedFilesPresentCheck",
                        (r, rule) -> new SignedFilesPresentCheck(i18n, r, emptyContainer, rule));

                final SignatureWrapper kidSignature = signatureWithMismatchedKeyIdentifier();
                savRow(sav, i18n, "synthetic", "kid-mismatch", "KeyIdentifierMatchCheck",
                        (r, rule) -> new KeyIdentifierMatchCheck(i18n, r, kidSignature, rule));

                DiagnosticData collidingData = collidingByteRangeDiagnosticData();
                final SignatureWrapper collidingSignature = collidingData.getSignatures().get(0);
                fcRow(fc, i18n, "synthetic", "byte-range-collision", "ByteRangeCollisionCheck",
                        (r, rule) -> new ByteRangeCollisionCheck(i18n, r, collidingSignature, collidingData, rule));
            }

            // --- aov: AlgorithmObsolescenceValidationCheck at each conclusion shape. The
            // corpus SAV rows only ever feed it a clean PASSED XmlAOV; these pin the Level it
            // derives from errors/warnings/infos, its additional-info rendering (which is the
            // only caller of ValidationProcessUtils#getFormattedDate reachable in phase 8c),
            // and the error message it lifts out of the block.
            for (String shape : new String[] { "passed", "passed-with-algo", "passed-with-algo-nokeysize",
                                               "error", "warning", "info" }) {
                final XmlAOV aov = aovOfShape(shape);
                savRow(sav, i18n, "synthetic", "aov-" + shape, "AlgorithmObsolescenceValidationCheck",
                        (r, rule) -> new AlgorithmObsolescenceValidationCheck<>(
                                i18n, r, aov, CURRENT_TIME, MessageTag.ACCM_POS_SIG_SIG, "T-AOV"));
            }

            // --- synthetic signatures for the branches no corpus dump reaches
            {
                final SignatureWrapper roleSignature = signatureWithRoles();
                savRowV(sav, i18n, "synthetic", "certified-roles", "CertifiedRolesCheck-any",
                        (r, rule) -> new CertifiedRolesCheck(i18n, r, roleSignature, ANY));
                savRowV(sav, i18n, "synthetic", "claimed-roles", "ClaimedRolesCheck-any",
                        (r, rule) -> new ClaimedRolesCheck(i18n, r, roleSignature, ANY));

                final SignatureWrapper commitmentSignature = signatureWithCommitment();
                savRowV(sav, i18n, "synthetic", "commitment", "CommitmentTypeIndicationsCheck-any",
                        (r, rule) -> new CommitmentTypeIndicationsCheck(i18n, r, commitmentSignature, ANY));
                savRowV(sav, i18n, "synthetic", "commitment", "CommitmentTypeIndicationsCheck-match",
                        (r, rule) -> new CommitmentTypeIndicationsCheck(i18n, r, commitmentSignature, MATCH_COMMITMENT));

                final SignatureWrapper x509Signature = signatureWithX509Url();
                savRow(sav, i18n, "synthetic", "x509url", "X509UrlPresentCheck",
                        (r, rule) -> new X509UrlPresentCheck(i18n, r, x509Signature, rule));
                savRow(sav, i18n, "synthetic", "x509url", "X509UrlMatchCheck",
                        (r, rule) -> new X509UrlMatchCheck(i18n, r, x509Signature, rule));

                final SignatureWrapper vdRefsSignature = signatureWithValidationDataRefsOnlyTimestamp();
                savRow(sav, i18n, "synthetic", "vd-refs-only", "ValidationDataRefsOnlyTimeStampCheck",
                        (r, rule) -> new ValidationDataRefsOnlyTimeStampCheck(i18n, r, vdRefsSignature, rule));

                final SignatureWrapper ecSignature = signatureWithEcdsaKeySizeMismatch();
                fcRow(fc, i18n, "synthetic", "ecdsa-mismatch", "EllipticCurveKeySizeCheck",
                        (r, rule) -> new EllipticCurveKeySizeCheck(i18n, r, ecSignature, rule));
            }

            for (boolean present : new boolean[] { true, false }) {
                final boolean mimetypePresent = present;
                fcRow(fc, i18n, "synthetic", "mimetype-present-" + present, "MimeTypeFilePresentCheck",
                        (r, rule) -> new MimeTypeFilePresentCheck(i18n, r, mimetypePresent, rule));
            }
            for (boolean compliant : new boolean[] { true, false }) {
                final boolean pdfaCompliant = compliant;
                fcRow(fc, i18n, "synthetic", "pdfa-compliant-" + compliant, "PDFAComplianceCheck",
                        (r, rule) -> new PDFAComplianceCheck(i18n, r, pdfaCompliant, rule));
            }
            for (String profile : new String[] { null, "PDF/A-2A" }) {
                final String pdfaProfile = profile;
                String label = profile == null ? "null" : profile;
                fcRowV(fc, i18n, "synthetic", "pdfa-profile-" + label + "-any", "PDFAProfileCheck",
                        (r, rule) -> new PDFAProfileCheck(i18n, r, pdfaProfile, ANY));
                fcRowV(fc, i18n, "synthetic", "pdfa-profile-" + label + "-none", "PDFAProfileCheck",
                        (r, rule) -> new PDFAProfileCheck(i18n, r, pdfaProfile, NONE));
            }
        }
    }

    /** Fixed validation time, shared with the AOV additional-info rendering. */
    private static final java.util.Date CURRENT_TIME = new java.util.Date(1704067200000L);

    /** A dump whose single PAdES signature carries an invalid /ByteRange. */
    private static DiagnosticData failingByteRangeDiagnosticData() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData jaxb =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData();
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlSignature();
        xml.setId("S-BAD-BYTERANGE");
        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary dictionary =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary();
        eu.europa.esig.dss.diagnostic.jaxb.XmlByteRange byteRange =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlByteRange();
        byteRange.getValue().addAll(java.util.Arrays.asList(
                java.math.BigInteger.ZERO, java.math.BigInteger.valueOf(100),
                java.math.BigInteger.valueOf(200), java.math.BigInteger.valueOf(300)));
        byteRange.setValid(false);
        dictionary.setSignatureByteRange(byteRange);
        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision revision =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision();
        revision.setPDFSignatureDictionary(dictionary);
        xml.setPDFRevision(revision);
        jaxb.getSignatures().add(xml);
        return new DiagnosticData(jaxb);
    }

    /** A signature carrying a duplicated digest-matcher reference. */
    private static SignatureWrapper signatureWithDuplicatedReference() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlDigestMatcher matcher =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlDigestMatcher();
        matcher.setType(eu.europa.esig.dss.enumerations.DigestMatcherType.REFERENCE);
        matcher.setUri("#r-id");
        matcher.setDuplicated(true);
        xml.getDigestMatchers().add(matcher);
        return new SignatureWrapper(xml);
    }

    /** A signature flagged as duplicated. */
    private static SignatureWrapper signatureWithDuplicatedSignature() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        xml.setDuplicated(true);
        return new SignatureWrapper(xml);
    }

    /** A PDF revision that fails every lock / consistency / overlap check. */
    private static PDFRevisionWrapper failingPdfRevision() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision revision =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision();

        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary dictionary =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary();
        dictionary.setConsistent(false);
        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFLockDictionary lock =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFLockDictionary();
        lock.setAction(eu.europa.esig.dss.enumerations.PdfLockAction.ALL);
        lock.setPermissions(eu.europa.esig.dss.enumerations.CertificationPermission.NO_CHANGE_PERMITTED);
        eu.europa.esig.dss.diagnostic.jaxb.XmlDocMDP docMDP =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlDocMDP();
        docMDP.setPermissions(eu.europa.esig.dss.enumerations.CertificationPermission.NO_CHANGE_PERMITTED);
        dictionary.setDocMDP(docMDP);
        dictionary.setFieldMDP(lock);
        revision.setPDFSignatureDictionary(dictionary);

        eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureField field =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureField();
        field.setSigFieldLock(lock);
        revision.getFields().add(field);

        eu.europa.esig.dss.diagnostic.jaxb.XmlObjectModification modification =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlObjectModification();
        modification.setFieldName("field-1");
        eu.europa.esig.dss.diagnostic.jaxb.XmlObjectModifications modifications =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlObjectModifications();
        modifications.getUndefined().add(modification);
        revision.setModificationDetection(new eu.europa.esig.dss.diagnostic.jaxb.XmlModificationDetection());
        revision.getModificationDetection().setObjectModifications(modifications);
        eu.europa.esig.dss.diagnostic.jaxb.XmlModification overlap =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlModification();
        overlap.setPage(java.math.BigInteger.ONE);
        revision.getModificationDetection().getAnnotationOverlap().add(overlap);

        return new PDFRevisionWrapper(revision);
    }

    /** A signature carrying content-type / content-hints / content-identifier attributes. */
    private static SignatureWrapper signatureWithContentAttributes() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        xml.setContentType("1.2.840.113549.1.7.1");
        xml.setContentHints("content-hints");
        xml.setContentIdentifier("content-identifier");
        return new SignatureWrapper(xml);
    }

    /** A signature whose key-identifier reference does not match the issuer serial. */
    private static SignatureWrapper signatureWithMismatchedKeyIdentifier() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef ref =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef();
        ref.setOrigin(eu.europa.esig.dss.enumerations.CertificateRefOrigin.KEY_IDENTIFIER);
        eu.europa.esig.dss.diagnostic.jaxb.XmlIssuerSerial issuerSerial =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlIssuerSerial();
        issuerSerial.setValue(new byte[] { 1, 2, 3 });
        issuerSerial.setMatch(false);
        ref.setIssuerSerial(issuerSerial);
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate certificate =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate();
        certificate.setId("C-SYNTHETIC");
        eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate related =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate();
        related.setCertificate(certificate);
        related.getCertificateRefs().add(ref);
        eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates found =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates();
        found.getRelatedCertificates().add(related);
        xml.setFoundCertificates(found);
        return new SignatureWrapper(xml);
    }

    /** Two PAdES signatures whose /ByteRange intervals overlap. */
    private static DiagnosticData collidingByteRangeDiagnosticData() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData jaxb =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData();
        for (int i = 0; i < 2; i++) {
            eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml =
                    new eu.europa.esig.dss.diagnostic.jaxb.XmlSignature();
            xml.setId("S-COLLIDE-" + i);
            eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary dictionary =
                    new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFSignatureDictionary();
            eu.europa.esig.dss.diagnostic.jaxb.XmlByteRange byteRange =
                    new eu.europa.esig.dss.diagnostic.jaxb.XmlByteRange();
            byteRange.getValue().addAll(java.util.Arrays.asList(
                    java.math.BigInteger.ZERO, java.math.BigInteger.valueOf(100),
                    java.math.BigInteger.valueOf(200), java.math.BigInteger.valueOf(300)));
            dictionary.setSignatureByteRange(byteRange);
            eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision revision =
                    new eu.europa.esig.dss.diagnostic.jaxb.XmlPDFRevision();
            revision.setPDFSignatureDictionary(dictionary);
            xml.setPDFRevision(revision);
            jaxb.getSignatures().add(xml);
        }
        return new DiagnosticData(jaxb);
    }

    /** An XmlAOV whose conclusion carries the requested message shape. */
    private static XmlAOV aovOfShape(String shape) {
        XmlAOV aov = new XmlAOV();
        XmlConclusion conclusion = new XmlConclusion();
        switch (shape) {
            case "error":
                conclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.INDETERMINATE);
                conclusion.setSubIndication(eu.europa.esig.dss.enumerations.SubIndication.CRYPTO_CONSTRAINTS_FAILURE);
                conclusion.getErrors().add(message("ASCCM_AR_ANS_ANR", "The algorithm is no longer reliable!"));
                break;
            case "warning":
                conclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
                conclusion.getWarnings().add(message("ASCCM_AR_ANS_AKSNR", "The key size is no longer reliable!"));
                break;
            case "info":
                conclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
                conclusion.getInfos().add(message("ASCCM_AR_ANS_ANR", "The algorithm expires soon."));
                break;
            default:
                conclusion.setIndication(eu.europa.esig.dss.enumerations.Indication.PASSED);
        }
        aov.setConclusion(conclusion);
        if ("passed-with-algo".equals(shape) || "passed-with-algo-nokeysize".equals(shape)) {
            XmlCryptographicAlgorithm algorithm = new XmlCryptographicAlgorithm();
            algorithm.setName("RSA with SHA256");
            if ("passed-with-algo".equals(shape)) {
                algorithm.setKeyLength("2048");
            }
            XmlCryptographicValidation validation = new XmlCryptographicValidation();
            validation.setAlgorithm(algorithm);
            aov.setSignatureCryptographicValidation(validation);
        }
        return aov;
    }

    private static XmlMessage message(String key, String value) {
        XmlMessage m = new XmlMessage();
        m.setKey(key);
        m.setValue(value);
        return m;
    }

    /** A FAIL MultiValuesRule accepting exactly the synthetic commitment identifier. */
    private static final MultiValuesRule MATCH_COMMITMENT = new MultiValuesRule() {
        @Override public Level getLevel() { return Level.FAIL; }
        @Override public List<String> getValues() { return Collections.singletonList("1.2.840.113549.1.9.16.6.1"); }
    };

    /** A signature carrying one claimed and one certified signer role. */
    private static SignatureWrapper signatureWithRoles() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignerRole claimed = new eu.europa.esig.dss.diagnostic.jaxb.XmlSignerRole();
        claimed.setRole("claimed-role");
        claimed.setCategory(eu.europa.esig.dss.enumerations.EndorsementType.CLAIMED);
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignerRole certified = new eu.europa.esig.dss.diagnostic.jaxb.XmlSignerRole();
        certified.setRole("certified-role");
        certified.setCategory(eu.europa.esig.dss.enumerations.EndorsementType.CERTIFIED);
        xml.getSignerRole().add(claimed);
        xml.getSignerRole().add(certified);
        return new SignatureWrapper(xml);
    }

    /** A signature carrying one commitment-type indication. */
    private static SignatureWrapper signatureWithCommitment() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlCommitmentTypeIndication indication =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCommitmentTypeIndication();
        indication.setIdentifier("1.2.840.113549.1.9.16.6.1");
        xml.getCommitmentTypeIndications().add(indication);
        return new SignatureWrapper(xml);
    }

    /** A signature whose signing certificate is also referenced by an X509_URL reference. */
    private static SignatureWrapper signatureWithX509Url() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate certificate =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate();
        certificate.setId("C-SYNTHETIC");

        eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate signingCertificate =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate();
        signingCertificate.setCertificate(certificate);
        xml.setSigningCertificate(signingCertificate);

        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef ref =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRef();
        ref.setOrigin(eu.europa.esig.dss.enumerations.CertificateRefOrigin.X509_URL);
        eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate related =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlRelatedCertificate();
        related.setCertificate(certificate);
        related.getCertificateRefs().add(ref);
        eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates found =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlFoundCertificates();
        found.getRelatedCertificates().add(related);
        xml.setFoundCertificates(found);
        return new SignatureWrapper(xml);
    }

    /** A signature carrying one VALIDATION_DATA_REFS_ONLY_TIMESTAMP. */
    private static SignatureWrapper signatureWithValidationDataRefsOnlyTimestamp() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlTimestamp timestamp =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlTimestamp();
        timestamp.setId("T-SYNTHETIC");
        timestamp.setType(eu.europa.esig.dss.enumerations.TimestampType.VALIDATION_DATA_REFSONLY_TIMESTAMP);
        eu.europa.esig.dss.diagnostic.jaxb.XmlFoundTimestamp found =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlFoundTimestamp();
        found.setTimestamp(timestamp);
        xml.getFoundTimestamps().add(found);
        return new SignatureWrapper(xml);
    }

    /** A signature signed with ECDSA whose key size does not match its digest algorithm. */
    private static SignatureWrapper signatureWithEcdsaKeySizeMismatch() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = baseSignature();
        eu.europa.esig.dss.diagnostic.jaxb.XmlBasicSignature basic =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlBasicSignature();
        basic.setEncryptionAlgoUsedToSignThisToken(eu.europa.esig.dss.enumerations.EncryptionAlgorithm.ECDSA);
        basic.setDigestAlgoUsedToSignThisToken(eu.europa.esig.dss.enumerations.DigestAlgorithm.SHA512);
        basic.setKeyLengthUsedToSignThisToken("256");
        xml.setBasicSignature(basic);
        return new SignatureWrapper(xml);
    }

    private static eu.europa.esig.dss.diagnostic.jaxb.XmlSignature baseSignature() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignature xml = new eu.europa.esig.dss.diagnostic.jaxb.XmlSignature();
        xml.setId("S-SYNTHETIC");
        return xml;
    }

    /**
     * Detailed-report XmlTimestamps for every used time-stamp of the dump, each
     * carrying a basic-validation conclusion with the given indication.
     */
    private static List<XmlTimestamp> reportTimestamps(DiagnosticData diagnosticData, Indication indication) {
        List<XmlTimestamp> timestamps = new java.util.ArrayList<>();
        List<TimestampWrapper> wrappers;
        try {
            wrappers = diagnosticData.getTimestampList();
        } catch (RuntimeException e) {
            return timestamps;
        }
        for (TimestampWrapper timestamp : wrappers) {
            XmlConclusion conclusion = new XmlConclusion();
            conclusion.setIndication(indication);
            XmlValidationProcessBasicTimestamp basic = new XmlValidationProcessBasicTimestamp();
            basic.setConclusion(conclusion);
            XmlTimestamp xmlTimestamp = new XmlTimestamp();
            xmlTimestamp.setId(timestamp.getId());
            xmlTimestamp.setValidationProcessBasicTimestamp(basic);
            timestamps.add(xmlTimestamp);
        }
        return timestamps;
    }

    // --------------------------------------------------------------- plumbing

    /** A chain of exactly one item over an XmlFC. */
    static class SingleFCChain extends Chain<XmlFC> {
        private final BiFunction<XmlFC, LevelRule, ChainItem<XmlFC>> factory;

        SingleFCChain(I18nProvider i18nProvider, BiFunction<XmlFC, LevelRule, ChainItem<XmlFC>> factory) {
            super(i18nProvider, new XmlFC());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    /** A chain of exactly one item over an XmlSAV. */
    static class SingleSAVChain extends Chain<XmlSAV> {
        private final BiFunction<XmlSAV, LevelRule, ChainItem<XmlSAV>> factory;

        SingleSAVChain(I18nProvider i18nProvider, BiFunction<XmlSAV, LevelRule, ChainItem<XmlSAV>> factory) {
            super(i18nProvider, new XmlSAV());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    private static void fcRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                              BiFunction<XmlFC, LevelRule, ChainItem<XmlFC>> factory) {
        fcRowV(out, i18n, file, token, check, factory);
    }

    private static void fcRowV(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                               BiFunction<XmlFC, LevelRule, ChainItem<XmlFC>> factory) {
        try {
            XmlFC result = new SingleFCChain(i18n, factory).execute();
            out.println(row(file, token, check, Context.SIGNATURE, "FC", result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    private static void savRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                               BiFunction<XmlSAV, LevelRule, ChainItem<XmlSAV>> factory) {
        savRowV(out, i18n, file, token, check, factory);
    }

    private static void savRowV(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                                BiFunction<XmlSAV, LevelRule, ChainItem<XmlSAV>> factory) {
        try {
            XmlSAV result = new SingleSAVChain(i18n, factory).execute();
            out.println(row(file, token, check, Context.SIGNATURE, "SAV", result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    // ------------------------------------------------------------------- JSON

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
    }

    private static String row(String file, String tokenId, String check, Context context, String block,
                              XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(tokenId)).append(',');
        key(out, "check").append(str(check)).append(',');
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
