/*
 * Java oracle driver for the 83 eu.europa.esig.dss.validation.process.bbb.xcv
 * checks/sub.checks/sub.checks.pseudo/rfc.checks classes ported by the XCVB
 * half of phase 8d (see the manifest in
 * scratchpad/s8d_XCVB.txt of the porting session).
 *
 * Every check is driven directly through a chain of exactly one item at
 * Level.FAIL, over REAL certificates (and, for the rfc.checks, real revocation
 * data) drawn from the marshal-parity corpus in
 * dss/diagnostic/jaxb/testdata/oracle. Driving each check over every used
 * certificate of every dump is what produces both the happy and the failure
 * row wherever the corpus varies (CA vs leaf certificates, self-signed roots,
 * QC-statement-bearing certificates, certificates with/without name
 * constraints, etc). A handful of checks whose failure branch the corpus never
 * reaches (a certificate carrying a QC statement whose retention period is
 * below an artificially high threshold, an XmlAOV/XmlCRS/XmlRFC "wrapper"
 * result at each conclusion shape, a name-constraint violation) are
 * additionally driven over synthetic literals/wrappers built in this file, in
 * the same spirit as FcSavDirectOracle's byte-range/duplicated-reference
 * synthetics.
 *
 * Rows are written to
 *   dss/validation/process/bbb/xcv/testdata/oracle/xcv_direct.jsonl
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle XcvOracle.java
 *   java  -cp "$CP:/tmp/oracle" XcvOracle <diagnostic-dump-dir> <dss-repo-root>
 */

import eu.europa.esig.dss.detailedreport.jaxb.XmlAOV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlBlockType;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCRS;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCertificateChainCryptographicValidation;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraint;
import eu.europa.esig.dss.detailedreport.jaxb.XmlConstraintsConclusion;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicAlgorithm;
import eu.europa.esig.dss.detailedreport.jaxb.XmlCryptographicValidation;
import eu.europa.esig.dss.detailedreport.jaxb.XmlMessage;
import eu.europa.esig.dss.detailedreport.jaxb.XmlRFC;
import eu.europa.esig.dss.detailedreport.jaxb.XmlSubXCV;
import eu.europa.esig.dss.detailedreport.jaxb.XmlXCV;
import eu.europa.esig.dss.diagnostic.CertificateRevocationWrapper;
import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.RevocationWrapper;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.enumerations.Context;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.Level;
import eu.europa.esig.dss.enumerations.SubContext;
import eu.europa.esig.dss.enumerations.SubIndication;
import eu.europa.esig.dss.i18n.I18nProvider;
import eu.europa.esig.dss.model.policy.CertificateApplicabilityRule;
import eu.europa.esig.dss.model.policy.DurationRule;
import eu.europa.esig.dss.model.policy.LevelRule;
import eu.europa.esig.dss.model.policy.MultiValuesRule;
import eu.europa.esig.dss.model.policy.NumericValueRule;
import eu.europa.esig.dss.model.policy.ValueRule;
import eu.europa.esig.dss.validation.process.Chain;
import eu.europa.esig.dss.validation.process.ChainItem;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.CertificateValidationBeforeSunsetDateCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.CertificateValidationBeforeSunsetDateWithIdCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.CheckSubXCVResult;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.ProspectiveCertificateChainAtValidationTimeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.ProspectiveCertificateChainCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.TrustServiceStatusCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.checks.TrustServiceTypeIdentifierCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rfc.checks.AcceptableRevocationDataAvailableCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rfc.checks.NextUpdateCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rfc.checks.RevocationDataFreshCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rfc.checks.RevocationDataFreshCheckWithNullConstraint;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.AuthorityInfoAccessPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.AuthorityKeyIdentifierPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.BasicConstraintsCACheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.BasicConstraintsMaxPathLengthCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateAlgorithmObsolescenceValidationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateForbiddenExtensionsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateIssuedToLegalPersonCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateIssuedToNaturalPersonCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateIssuerNameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateMinQcEuRetentionPeriodCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateMinQcTransactionLimitCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateNameConstraintsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateNotOnHoldCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateNotRevokedCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateNotSelfSignedCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.rfc.checks.RevocationDataAvailableCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePS2DQcCompetentAuthorityIdCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePS2DQcCompetentAuthorityNameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePS2DQcRolesOfPSPCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePolicyIdsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePolicyQualifiedIdsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePolicySupportedByQSCDIdsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificatePolicyTreeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcCCLegislationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcComplianceCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcEuLimitValueCurrencyCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcEuPDSLocationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcIdentificationMethodCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcPSBAuthSourceIdentificationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcPSBCountryOfLegislationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcPSBLegislationIdentificationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcQSCDLegislationCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcSSCDCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateQcTypeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateRevocationSelectorResultCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateSelfSignedCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateSemanticsIdentifierCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateSignatureValidCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateSupportedCriticalExtensionsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CertificateValidityRangeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CommonNameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.CountryCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.EmailCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.ExtendedKeyUsageCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.GivenNameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.KeyUsageCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.LocalityCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.NoRevAvailCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.OrganizationIdentifierCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.OrganizationNameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.OrganizationUnitCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.OtherTrustAnchorExistsCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.PseudoUsageCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.PseudonymCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.RevocationDataRequiredCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.RevocationFreshnessCheckerResultCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.RevocationInfoAccessPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.RevocationIssuerTrustedCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.RevocationIssuerValidityRangeCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.SerialNumberCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.StateCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.SubjectKeyIdentifierPresentCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.SurnameCheck;
import eu.europa.esig.dss.validation.process.bbb.xcv.sub.checks.TitleCheck;

import java.io.File;
import java.io.PrintWriter;
import java.util.Arrays;
import java.util.Collections;
import java.util.Comparator;
import java.util.Date;
import java.util.List;
import java.util.function.BiFunction;

public class XcvOracle {

    private static final LevelRule FAIL = () -> Level.FAIL;

    /** Fixed validation/current time: 2024-01-01T00:00:00Z. */
    private static final Date CURRENT_TIME = new Date(1704067200000L);

    /** A FAIL MultiValuesRule accepting any present value. */
    private static final MultiValuesRule ANY = multiValues("*");
    /** A FAIL MultiValuesRule accepting nothing the corpus carries. */
    private static final MultiValuesRule NONE = multiValues("NO-SUCH-VALUE");
    /** A FAIL NumericValueRule with a low threshold most present values clear. */
    private static final NumericValueRule NUM_LOW = numeric(0);
    /** A FAIL NumericValueRule with a threshold no corpus value clears. */
    private static final NumericValueRule NUM_HIGH = numeric(999999999);
    /** A FAIL ValueRule accepting any present value. */
    private static final ValueRule VALUE_ANY = value("*");
    /** A FAIL ValueRule accepting nothing the corpus carries. */
    private static final ValueRule VALUE_NONE = value("NO-SUCH-VALUE");
    /** A FAIL CertificateApplicabilityRule matching nothing (revocation never skipped). */
    private static final CertificateApplicabilityRule APPLIES_NEVER = new CertificateApplicabilityRule() {
        @Override public Level getLevel() { return Level.FAIL; }
        @Override public MultiValuesRule getCertificateExtensions() { return null; }
        @Override public MultiValuesRule getCertificatePolicies() { return null; }
    };
    /** A FAIL DurationRule with a small duration most revocation data will not satisfy. */
    private static final DurationRule DURATION_TIGHT = duration(1000L);
    /** A FAIL DurationRule with a huge duration almost all revocation data satisfies. */
    private static final DurationRule DURATION_LOOSE = duration(1000L * 60 * 60 * 24 * 365 * 50L);

    public static void main(String[] args) throws Exception {
        File inputDir = new File(args[0]);
        File repoRoot = new File(args[1]);

        File[] files = inputDir.listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(files, Comparator.comparing(File::getName));

        // XCVA's committed synthetic dumps (trust anchors, sunset dates, ...) carry
        // shapes the marshal-parity corpus has none of; reused here rather than
        // duplicated, per the file's own "neither side gets a private fixture" spirit.
        File[] synthetic = new File(repoRoot, "validation/process/bbb/xcv/testdata/dd")
                .listFiles((dir, name) -> name.endsWith(".xml"));
        Arrays.sort(synthetic, Comparator.comparing(File::getName));

        I18nProvider i18n = new I18nProvider();

        try (PrintWriter out = writer(repoRoot,
                "validation/process/bbb/xcv/testdata/oracle/xcv_direct.jsonl")) {

            for (File file : synthetic) {
                DiagnosticData dd = new DiagnosticData(
                        DiagnosticDataFacade.newFacade().unmarshall(file, false));
                emitDump(out, i18n, "dd/" + file.getName(), dd);
            }

            for (File file : files) {
                String name = file.getName();
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
                emitDump(out, i18n, name, dd);
            }

            emitSynthetic(out, i18n);
        }
    }

    // ------------------------------------------------------------- real corpus

    private static void emitDump(PrintWriter out, I18nProvider i18n, String file, DiagnosticData dd) {
        for (CertificateWrapper certificate : dd.getUsedCertificates()) {
            emitCertificateChecks(out, i18n, file, certificate);

            for (CertificateRevocationWrapper revocation : certificate.getCertificateRevocationData()) {
                String token = certificate.getId() + "|" + revocation.getId();
                rfcRow(out, i18n, file, token, "NextUpdateCheck",
                        (r, l) -> new NextUpdateCheck(i18n, r, revocation, l));
                rfcRow(out, i18n, file, token, "RevocationDataFreshCheck",
                        (r, l) -> new RevocationDataFreshCheck(i18n, r, revocation, CURRENT_TIME, DURATION_LOOSE));
                rfcRow(out, i18n, file, token, "RevocationDataFreshCheck",
                        (r, l) -> new RevocationDataFreshCheck(i18n, r, revocation, CURRENT_TIME, DURATION_TIGHT));
                rfcRow(out, i18n, file, token, "RevocationDataFreshCheckWithNullConstraint",
                        (r, l) -> new RevocationDataFreshCheckWithNullConstraint(i18n, r, revocation, CURRENT_TIME, FAIL));
                subRow(out, i18n, file, token, "AcceptableRevocationDataAvailableCheck",
                        (r, l) -> new AcceptableRevocationDataAvailableCheck<>(i18n, r, revocation, l), true);
                subRow(out, i18n, file, token, "RevocationIssuerTrustedCheck",
                        (r, l) -> new RevocationIssuerTrustedCheck<>(i18n, r, revocation.getSigningCertificate(), CURRENT_TIME, FAIL, l), true);
                subRow(out, i18n, file, token, "RevocationIssuerValidityRangeCheck",
                        (r, l) -> new RevocationIssuerValidityRangeCheck<>(i18n, r, revocation, CURRENT_TIME, l), true);
                for (SubContext subContext : new SubContext[] { SubContext.SIGNING_CERT, SubContext.CA_CERTIFICATE }) {
                    subRow(out, i18n, file, token + "|" + subContext, "CertificateNotOnHoldCheck",
                            (r, l) -> new CertificateNotOnHoldCheck(i18n, r, revocation, CURRENT_TIME, l), true);
                    subRow(out, i18n, file, token + "|" + subContext, "CertificateNotRevokedCheck",
                            (r, l) -> new CertificateNotRevokedCheck(i18n, r, revocation, CURRENT_TIME, l, subContext), true);
                }
            }
            subRow(out, i18n, file, certificate.getId(), "AcceptableRevocationDataAvailableCheck-none",
                    (r, l) -> new AcceptableRevocationDataAvailableCheck<>(i18n, r, null, l), true);
            subRow(out, i18n, file, certificate.getId(), "RevocationDataAvailableCheck",
                    (r, l) -> new RevocationDataAvailableCheck<>(i18n, r, certificate, l), true);
        }
    }

    private static void emitCertificateChecks(PrintWriter out, I18nProvider i18n, String file, CertificateWrapper certificate) {
        String id = certificate.getId();

        subRow(out, i18n, file, id, "BasicConstraintsCACheck",
                (r, l) -> new BasicConstraintsCACheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "BasicConstraintsMaxPathLengthCheck",
                (r, l) -> new BasicConstraintsMaxPathLengthCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "AuthorityInfoAccessPresentCheck",
                (r, l) -> new AuthorityInfoAccessPresentCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "AuthorityKeyIdentifierPresentCheck",
                (r, l) -> new AuthorityKeyIdentifierPresentCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "SubjectKeyIdentifierPresentCheck",
                (r, l) -> new SubjectKeyIdentifierPresentCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "RevocationInfoAccessPresentCheck",
                (r, l) -> new RevocationInfoAccessPresentCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "NoRevAvailCheck",
                (r, l) -> new NoRevAvailCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateSelfSignedCheck",
                (r, l) -> new CertificateSelfSignedCheck<>(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateNotSelfSignedCheck",
                (r, l) -> new CertificateNotSelfSignedCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateSignatureValidCheck",
                (r, l) -> new CertificateSignatureValidCheck<>(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateIssuerNameCheck",
                (r, l) -> new CertificateIssuerNameCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "OtherTrustAnchorExistsCheck",
                (r, l) -> new OtherTrustAnchorExistsCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificatePolicyTreeCheck",
                (r, l) -> new CertificatePolicyTreeCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateNameConstraintsCheck",
                (r, l) -> new CertificateNameConstraintsCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "SerialNumberCheck",
                (r, l) -> new SerialNumberCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "PseudoUsageCheck",
                (r, l) -> new PseudoUsageCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificatePolicyQualifiedIdsCheck",
                (r, l) -> new CertificatePolicyQualifiedIdsCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificatePolicySupportedByQSCDIdsCheck",
                (r, l) -> new CertificatePolicySupportedByQSCDIdsCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateQcComplianceCheck",
                (r, l) -> new CertificateQcComplianceCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateQcSSCDCheck",
                (r, l) -> new CertificateQcSSCDCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateIssuedToNaturalPersonCheck",
                (r, l) -> new CertificateIssuedToNaturalPersonCheck(i18n, r, certificate, l), true);
        subRow(out, i18n, file, id, "CertificateIssuedToLegalPersonCheck",
                (r, l) -> new CertificateIssuedToLegalPersonCheck(i18n, r, certificate, l), true);

        for (String variant : new String[] { "any", "none" }) {
            MultiValuesRule mv = "any".equals(variant) ? ANY : NONE;
            multiSubRow(out, i18n, file, id, "CommonNameCheck-" + variant,
                    (r, l) -> new CommonNameCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CountryCheck-" + variant,
                    (r, l) -> new CountryCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "EmailCheck-" + variant,
                    (r, l) -> new EmailCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "GivenNameCheck-" + variant,
                    (r, l) -> new GivenNameCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "LocalityCheck-" + variant,
                    (r, l) -> new LocalityCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "OrganizationIdentifierCheck-" + variant,
                    (r, l) -> new OrganizationIdentifierCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "OrganizationNameCheck-" + variant,
                    (r, l) -> new OrganizationNameCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "OrganizationUnitCheck-" + variant,
                    (r, l) -> new OrganizationUnitCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "StateCheck-" + variant,
                    (r, l) -> new StateCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "SurnameCheck-" + variant,
                    (r, l) -> new SurnameCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "TitleCheck-" + variant,
                    (r, l) -> new TitleCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "PseudonymCheck-" + variant,
                    (r, l) -> new PseudonymCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificatePolicyIdsCheck-" + variant,
                    (r, l) -> new CertificatePolicyIdsCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateForbiddenExtensionsCheck-" + variant,
                    (r, l) -> new CertificateForbiddenExtensionsCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateSupportedCriticalExtensionsCheck-" + variant,
                    (r, l) -> new CertificateSupportedCriticalExtensionsCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcEuPDSLocationCheck-" + variant,
                    (r, l) -> new CertificateQcEuPDSLocationCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcTypeCheck-" + variant,
                    (r, l) -> new CertificateQcTypeCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcCCLegislationCheck-" + variant,
                    (r, l) -> new CertificateQcCCLegislationCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcQSCDLegislationCheck-" + variant,
                    (r, l) -> new CertificateQcQSCDLegislationCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcIdentificationMethodCheck-" + variant,
                    (r, l) -> new CertificateQcIdentificationMethodCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateSemanticsIdentifierCheck-" + variant,
                    (r, l) -> new CertificateSemanticsIdentifierCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificatePS2DQcCompetentAuthorityIdCheck-" + variant,
                    (r, l) -> new CertificatePS2DQcCompetentAuthorityIdCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificatePS2DQcCompetentAuthorityNameCheck-" + variant,
                    (r, l) -> new CertificatePS2DQcCompetentAuthorityNameCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificatePS2DQcRolesOfPSPCheck-" + variant,
                    (r, l) -> new CertificatePS2DQcRolesOfPSPCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcPSBCountryOfLegislationCheck-" + variant,
                    (r, l) -> new CertificateQcPSBCountryOfLegislationCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcPSBAuthSourceIdentificationCheck-" + variant,
                    (r, l) -> new CertificateQcPSBAuthSourceIdentificationCheck(i18n, r, certificate, mv));
            multiSubRow(out, i18n, file, id, "CertificateQcPSBLegislationIdentificationCheck-" + variant,
                    (r, l) -> new CertificateQcPSBLegislationIdentificationCheck(i18n, r, certificate, mv));
        }

        for (String variant : new String[] { "low", "high" }) {
            NumericValueRule nv = "low".equals(variant) ? NUM_LOW : NUM_HIGH;
            subRow(out, i18n, file, id, "CertificateMinQcEuRetentionPeriodCheck-" + variant,
                    (r, l) -> new CertificateMinQcEuRetentionPeriodCheck(i18n, r, certificate, nv), true);
            subRow(out, i18n, file, id, "CertificateMinQcTransactionLimitCheck-" + variant,
                    (r, l) -> new CertificateMinQcTransactionLimitCheck(i18n, r, certificate, nv), true);
        }
        for (String variant : new String[] { "any", "none" }) {
            ValueRule vr = "any".equals(variant) ? VALUE_ANY : VALUE_NONE;
            subRow(out, i18n, file, id, "CertificateQcEuLimitValueCurrencyCheck-" + variant,
                    (r, l) -> new CertificateQcEuLimitValueCurrencyCheck(i18n, r, certificate, vr), true);
        }

        for (SubContext subContext : new SubContext[] { SubContext.SIGNING_CERT, SubContext.CA_CERTIFICATE }) {
            String suffix = subContext.name();
            for (String variant : new String[] { "any", "none" }) {
                MultiValuesRule mv = "any".equals(variant) ? ANY : NONE;
                multiSubRow(out, i18n, file, id, "KeyUsageCheck-" + suffix + "-" + variant,
                        (r, l) -> new KeyUsageCheck(i18n, r, certificate, Context.SIGNATURE, subContext, mv));
                multiSubRow(out, i18n, file, id, "ExtendedKeyUsageCheck-" + suffix + "-" + variant,
                        (r, l) -> new ExtendedKeyUsageCheck(i18n, r, certificate, Context.SIGNATURE, subContext, mv));
            }
            subRow(out, i18n, file, id, "CertificateValidationBeforeSunsetDateCheck-" + suffix,
                    (r, l) -> new CertificateValidationBeforeSunsetDateCheck<>(i18n, r, certificate, CURRENT_TIME, l), true);
            subRow(out, i18n, file, id, "CertificateValidationBeforeSunsetDateWithIdCheck-" + suffix,
                    (r, l) -> new CertificateValidationBeforeSunsetDateWithIdCheck<>(i18n, r, certificate, CURRENT_TIME, l), true);
            subRow(out, i18n, file, id, "RevocationDataRequiredCheck-" + suffix,
                    (r, l) -> new RevocationDataRequiredCheck<>(i18n, r, certificate, CURRENT_TIME, FAIL, APPLIES_NEVER), true);
        }

        subRow(out, i18n, file, id, "CertificateValidityRangeCheck",
                (r, l) -> new CertificateValidityRangeCheck<>(i18n, r, certificate, null, false, false, false, CURRENT_TIME, l), true);

        xcvRow(out, i18n, file, id, "ProspectiveCertificateChainCheck",
                (r, l) -> new ProspectiveCertificateChainCheck<>(i18n, r, certificate, Context.SIGNATURE, l));
        xcvRow(out, i18n, file, id, "ProspectiveCertificateChainAtValidationTimeCheck",
                (r, l) -> new ProspectiveCertificateChainAtValidationTimeCheck(i18n, r, certificate, CURRENT_TIME, l));

        Date usageTime = certificate.getNotBefore();
        if (usageTime != null) {
            for (String variant : new String[] { "any", "none" }) {
                MultiValuesRule mv = "any".equals(variant) ? ANY : NONE;
                xcvRow(out, i18n, file, id, "TrustServiceStatusCheck-" + variant,
                        (r, l) -> new TrustServiceStatusCheck(i18n, r, certificate, usageTime, Context.SIGNATURE, mv));
                xcvRow(out, i18n, file, id, "TrustServiceTypeIdentifierCheck-" + variant,
                        (r, l) -> new TrustServiceTypeIdentifierCheck(i18n, r, certificate, usageTime, Context.SIGNATURE, mv));
            }
        }

        // --- checks that consume a synthetic wrapper result keyed to this real certificate,
        // so the token-lookup / additional-info logic runs against a real certificate id.
        for (String shape : new String[] { "passed", "passed-with-algo", "error", "warning" }) {
            XmlAOV aov = aovOfShape(shape, id);
            subRow(out, i18n, file, id + "/" + shape, "CertificateAlgorithmObsolescenceValidationCheck",
                    (r, l) -> new CertificateAlgorithmObsolescenceValidationCheck<>(i18n, r, aov, CURRENT_TIME,
                            eu.europa.esig.dss.i18n.MessageTag.SIGNING_CERTIFICATE, id), true);
        }
        for (Indication indication : new Indication[] { Indication.PASSED, Indication.INDETERMINATE }) {
            XmlCRS crs = crsOfIndication(indication, id);
            subRow(out, i18n, file, id + "/" + indication, "CertificateRevocationSelectorResultCheck",
                    (r, l) -> new CertificateRevocationSelectorResultCheck<>(i18n, r, crs, l), true);

            XmlRFC rfc = rfcOfIndication(indication, id);
            subRow(out, i18n, file, id + "/" + indication, "RevocationFreshnessCheckerResultCheck",
                    (r, l) -> new RevocationFreshnessCheckerResultCheck<>(i18n, r, rfc, l), true);

            XmlSubXCV subXCV = subXCVOfIndication(indication, id);
            xcvRow(out, i18n, file, id + "/" + indication, "CheckSubXCVResult",
                    (r, l) -> new CheckSubXCVResult(i18n, r, subXCV, l));
        }
    }

    // ----------------------------------------------------------------- synthetic

    private static void emitSynthetic(PrintWriter out, I18nProvider i18n) {
        // The branches below are the ones the real corpus (plus XCVA's committed
        // dd/*.xml dumps) does not reach: see the per-check comments.

        // CertificateNameConstraintsCheck: no corpus CA restricts its subject's
        // DN with a name constraint the corpus leaf actually violates.
        CertificateWrapper ncLeaf = nameConstraintViolationLeaf();
        subRow(out, i18n, "synthetic", ncLeaf.getId(), "CertificateNameConstraintsCheck",
                (r, l) -> new CertificateNameConstraintsCheck(i18n, r, ncLeaf, l), true);

        // NoRevAvailCheck: no corpus certificate carries the noRevAvail extension
        // together with a conflicting OCSP access point.
        CertificateWrapper noRevAvailViolation = noRevAvailViolation();
        subRow(out, i18n, "synthetic", noRevAvailViolation.getId(), "NoRevAvailCheck",
                (r, l) -> new NoRevAvailCheck(i18n, r, noRevAvailViolation, l), true);

        // SerialNumberCheck: every corpus certificate carries a serial number.
        CertificateWrapper blankSerial = blankSerialNumberCertificate();
        subRow(out, i18n, "synthetic", blankSerial.getId(), "SerialNumberCheck",
                (r, l) -> new SerialNumberCheck(i18n, r, blankSerial, l), true);

        // The QC-statement / PSD2 / pseudonym checks below only ever see the
        // corpus's "absent" branch (no corpus certificate carries these QC
        // statements), so ANY never matches; this certificate supplies every
        // value at once to produce the "matched" OK branch too.
        CertificateWrapper qc = qcStatementCertificate();
        multiSubRow(out, i18n, "synthetic", qc.getId(), "PseudonymCheck-any",
                (r, l) -> new PseudonymCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcCCLegislationCheck-any",
                (r, l) -> new CertificateQcCCLegislationCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcQSCDLegislationCheck-any",
                (r, l) -> new CertificateQcQSCDLegislationCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcIdentificationMethodCheck-any",
                (r, l) -> new CertificateQcIdentificationMethodCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificatePS2DQcCompetentAuthorityIdCheck-any",
                (r, l) -> new CertificatePS2DQcCompetentAuthorityIdCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificatePS2DQcCompetentAuthorityNameCheck-any",
                (r, l) -> new CertificatePS2DQcCompetentAuthorityNameCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificatePS2DQcRolesOfPSPCheck-any",
                (r, l) -> new CertificatePS2DQcRolesOfPSPCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcPSBCountryOfLegislationCheck-any",
                (r, l) -> new CertificateQcPSBCountryOfLegislationCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcPSBAuthSourceIdentificationCheck-any",
                (r, l) -> new CertificateQcPSBAuthSourceIdentificationCheck(i18n, r, qc, ANY));
        multiSubRow(out, i18n, "synthetic", qc.getId(), "CertificateQcPSBLegislationIdentificationCheck-any",
                (r, l) -> new CertificateQcPSBLegislationIdentificationCheck(i18n, r, qc, ANY));

        // TrustServiceStatusCheck / TrustServiceTypeIdentifierCheck: no corpus
        // certificate carries an associated TrustServiceProvider entry.
        Date usageTime = new Date(1700000000000L);
        CertificateWrapper trustService = trustServiceCertificate(usageTime);
        xcvRow(out, i18n, "synthetic", trustService.getId(), "TrustServiceStatusCheck-any",
                (r, l) -> new TrustServiceStatusCheck(i18n, r, trustService, usageTime, Context.SIGNATURE, ANY));
        xcvRow(out, i18n, "synthetic", trustService.getId(), "TrustServiceTypeIdentifierCheck-any",
                (r, l) -> new TrustServiceTypeIdentifierCheck(i18n, r, trustService, usageTime, Context.SIGNATURE, ANY));

        // CertificateNotOnHoldCheck: no corpus revocation carries reason CERTIFICATE_HOLD.
        CertificateRevocationWrapper onHold = onHoldRevocation();
        subRow(out, i18n, "synthetic", onHold.getId(), "CertificateNotOnHoldCheck",
                (r, l) -> new CertificateNotOnHoldCheck(i18n, r, onHold, CURRENT_TIME, l), true);
    }

    /** A certificate revocation whose reason is CERTIFICATE_HOLD, at a past revocation date. */
    private static CertificateRevocationWrapper onHoldRevocation() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlRevocation revocation =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlRevocation();
        revocation.setId("R-SYNTH-ON-HOLD");
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRevocation certRevocation =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRevocation();
        certRevocation.setRevocation(revocation);
        certRevocation.setStatus(eu.europa.esig.dss.enumerations.CertificateStatus.REVOKED);
        certRevocation.setReason(eu.europa.esig.dss.enumerations.RevocationReason.CERTIFICATE_HOLD);
        certRevocation.setRevocationDate(new Date(CURRENT_TIME.getTime() - 1000L * 60 * 60 * 24));
        return new CertificateRevocationWrapper(certRevocation);
    }

    /** A leaf certificate whose DN falls outside its issuing CA's permitted name-constraint subtree. */
    private static CertificateWrapper nameConstraintViolationLeaf() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate ca = baseCertificate("C-SYNTH-NC-CA");
        eu.europa.esig.dss.diagnostic.jaxb.XmlNameConstraints nameConstraints =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlNameConstraints();
        eu.europa.esig.dss.diagnostic.jaxb.XmlGeneralSubtree permitted =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlGeneralSubtree();
        permitted.setType(eu.europa.esig.dss.enumerations.GeneralNameType.DIRECTORY_NAME);
        permitted.setValue("O=Allowed");
        nameConstraints.getPermittedSubtrees().add(permitted);
        nameConstraints.setOID(eu.europa.esig.dss.enumerations.CertificateExtensionEnum.NAME_CONSTRAINTS.getOid());
        ca.getCertificateExtensions().add(nameConstraints);

        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate leaf = baseCertificate("C-SYNTH-NC-LEAF");
        eu.europa.esig.dss.diagnostic.jaxb.XmlDistinguishedName leafDN =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlDistinguishedName();
        leafDN.setFormat("RFC2253");
        leafDN.setValue("CN=Leaf,O=Excluded");
        leaf.getSubjectDistinguishedName().add(leafDN);
        eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem chainItem =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem();
        chainItem.setCertificate(ca);
        leaf.getCertificateChain().add(chainItem);
        return new CertificateWrapper(leaf);
    }

    /** A certificate declaring noRevAvail while still publishing an OCSP access point (RFC 9608 violation). */
    private static CertificateWrapper noRevAvailViolation() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate xml = baseCertificate("C-SYNTH-NORA-VIOLATION");
        eu.europa.esig.dss.diagnostic.jaxb.XmlNoRevAvail noRevAvail =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlNoRevAvail();
        noRevAvail.setPresent(true);
        noRevAvail.setOID(eu.europa.esig.dss.enumerations.CertificateExtensionEnum.NO_REVOCATION_AVAILABLE.getOid());
        xml.getCertificateExtensions().add(noRevAvail);
        eu.europa.esig.dss.diagnostic.jaxb.XmlAuthorityInformationAccess aia =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlAuthorityInformationAccess();
        aia.getOcspUrls().add("http://example.org/ocsp");
        aia.setOID(eu.europa.esig.dss.enumerations.CertificateExtensionEnum.AUTHORITY_INFORMATION_ACCESS.getOid());
        xml.getCertificateExtensions().add(aia);
        return new CertificateWrapper(xml);
    }

    /** A certificate with no serial number set. */
    private static CertificateWrapper blankSerialNumberCertificate() {
        return new CertificateWrapper(baseCertificate("C-SYNTH-NO-SERIAL"));
    }

    /** A certificate carrying every QC/PSD2/pseudonym attribute the corpus has none of. */
    private static CertificateWrapper qcStatementCertificate() {
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate xml = baseCertificate("C-SYNTH-QC");
        xml.setPseudonym("Synthetic Pseudonym");

        eu.europa.esig.dss.diagnostic.jaxb.XmlQcStatements qc =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlQcStatements();
        qc.getQcCClegislation().add("EU");
        qc.getQcQSCDlegislation().add("EU");
        eu.europa.esig.dss.diagnostic.jaxb.XmlOID identMethod = new eu.europa.esig.dss.diagnostic.jaxb.XmlOID();
        identMethod.setValue("0.4.0.1862.1.6.1");
        identMethod.setDescription("physical presence");
        qc.setQcIdentMethod(identMethod);
        eu.europa.esig.dss.diagnostic.jaxb.XmlOID semantics = new eu.europa.esig.dss.diagnostic.jaxb.XmlOID();
        semantics.setValue("0.4.0.194121.1.1");
        semantics.setDescription("natural person");
        qc.setSemanticsIdentifier(semantics);

        eu.europa.esig.dss.diagnostic.jaxb.XmlQcPSB psb = new eu.europa.esig.dss.diagnostic.jaxb.XmlQcPSB();
        psb.setCountryOfLegislation("EU");
        psb.setAuthSourceIdentification("Synthetic Authentic Source");
        psb.setLegislationIdentification("Synthetic Legislation");
        qc.setQcPSB(psb);

        eu.europa.esig.dss.diagnostic.jaxb.XmlPSD2QcInfo psd2 = new eu.europa.esig.dss.diagnostic.jaxb.XmlPSD2QcInfo();
        psd2.setNcaId("EU-NCA-Synthetic");
        psd2.setNcaName("Synthetic National Competent Authority");
        eu.europa.esig.dss.diagnostic.jaxb.XmlRoleOfPSP role = new eu.europa.esig.dss.diagnostic.jaxb.XmlRoleOfPSP();
        role.setName("PSP_AS");
        eu.europa.esig.dss.diagnostic.jaxb.XmlOID roleOid = new eu.europa.esig.dss.diagnostic.jaxb.XmlOID();
        roleOid.setValue("0.4.0.19495.1.1");
        role.setOid(roleOid);
        psd2.getRolesOfPSP().add(role);
        qc.setPSD2QcInfo(psd2);

        qc.setOID(eu.europa.esig.dss.enumerations.CertificateExtensionEnum.QC_STATEMENTS.getOid());
        xml.getCertificateExtensions().add(qc);
        return new CertificateWrapper(xml);
    }

    /** A certificate associated with a TrustServiceProvider entry covering usageTime. */
    private static CertificateWrapper trustServiceCertificate(Date usageTime) {
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate xml = baseCertificate("C-SYNTH-TRUST-SERVICE");

        eu.europa.esig.dss.diagnostic.jaxb.XmlTrustService service =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlTrustService();
        service.setServiceType("http://uri.etsi.org/TrstSvc/Svctype/CA/QC");
        service.setStatus("http://uri.etsi.org/TrstSvc/TrustedList/Svcstatus/granted");
        service.setStartDate(new Date(usageTime.getTime() - 1000L * 60 * 60 * 24 * 365));
        service.setEndDate(null);
        service.setServiceDigitalIdentifier(xml);

        eu.europa.esig.dss.diagnostic.jaxb.XmlTrustServiceProvider provider =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlTrustServiceProvider();
        provider.getTrustServices().add(service);
        xml.getTrustServiceProviders().add(provider);
        return new CertificateWrapper(xml);
    }

    private static eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate baseCertificate(String id) {
        eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate xml = new eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate();
        xml.setId(id);
        return xml;
    }

    private static XmlAOV aovOfShape(String shape, String certificateId) {
        XmlAOV aov = new XmlAOV();
        XmlConclusion conclusion = new XmlConclusion();
        switch (shape) {
            case "error":
                conclusion.setIndication(Indication.INDETERMINATE);
                conclusion.setSubIndication(SubIndication.CRYPTO_CONSTRAINTS_FAILURE);
                conclusion.getErrors().add(message("ASCCM_AR_ANS_ANR", "The algorithm is no longer reliable!"));
                break;
            case "warning":
                conclusion.setIndication(Indication.PASSED);
                conclusion.getWarnings().add(message("ASCCM_AR_ANS_AKSNR", "The key size is no longer reliable!"));
                break;
            default:
                conclusion.setIndication(Indication.PASSED);
        }
        aov.setConclusion(conclusion);

        XmlCryptographicAlgorithm algorithm = new XmlCryptographicAlgorithm();
        algorithm.setName("RSA with SHA256");
        if ("passed-with-algo".equals(shape)) {
            algorithm.setKeyLength("2048");
        }
        XmlCryptographicValidation certValidation = new XmlCryptographicValidation();
        certValidation.setAlgorithm(algorithm);
        certValidation.setConclusion(conclusion);
        certValidation.setTokenId(certificateId);

        XmlCertificateChainCryptographicValidation chain = new XmlCertificateChainCryptographicValidation();
        chain.getCertificateCryptographicValidation().add(certValidation);
        aov.setCertificateChainCryptographicValidation(chain);
        return aov;
    }

    private static XmlMessage message(String key, String value) {
        XmlMessage m = new XmlMessage();
        m.setKey(key);
        m.setValue(value);
        return m;
    }

    private static XmlCRS crsOfIndication(Indication indication, String id) {
        XmlCRS crs = new XmlCRS();
        crs.setId(id);
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(indication);
        if (indication != Indication.PASSED) {
            conclusion.setSubIndication(SubIndication.CERTIFICATE_CHAIN_GENERAL_FAILURE);
            conclusion.getErrors().add(message("BBB_XCV_IARDPFC_ANS", "No acceptable revocation data found."));
        } else {
            crs.setLatestAcceptableRevocationId("R-" + id);
        }
        crs.setConclusion(conclusion);
        return crs;
    }

    private static XmlRFC rfcOfIndication(Indication indication, String id) {
        XmlRFC rfc = new XmlRFC();
        rfc.setId(id);
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(indication);
        if (indication != Indication.PASSED) {
            conclusion.setSubIndication(SubIndication.TRY_LATER);
            conclusion.getErrors().add(message("BBB_RFC_IRIF_ANS", "The revocation is not considered as 'fresh'."));
        }
        rfc.setConclusion(conclusion);
        return rfc;
    }

    private static XmlSubXCV subXCVOfIndication(Indication indication, String id) {
        XmlSubXCV subXCV = new XmlSubXCV();
        subXCV.setId(id);
        XmlConclusion conclusion = new XmlConclusion();
        conclusion.setIndication(indication);
        if (indication != Indication.PASSED) {
            conclusion.setSubIndication(SubIndication.CHAIN_CONSTRAINTS_FAILURE);
            conclusion.getErrors().add(message("BBB_XCV_SUB_ANS", "The certificate validation is not conclusive!"));
        }
        subXCV.setConclusion(conclusion);
        return subXCV;
    }

    // ------------------------------------------------------------------ helpers

    private static MultiValuesRule multiValues(String value) {
        return new MultiValuesRule() {
            @Override public Level getLevel() { return Level.FAIL; }
            @Override public List<String> getValues() { return Collections.singletonList(value); }
        };
    }

    private static NumericValueRule numeric(long value) {
        return new NumericValueRule() {
            @Override public Level getLevel() { return Level.FAIL; }
            @Override public Number getValue() { return value; }
        };
    }

    private static ValueRule value(String v) {
        return new ValueRule() {
            @Override public Level getLevel() { return Level.FAIL; }
            @Override public String getValue() { return v; }
        };
    }

    private static DurationRule duration(long millis) {
        return new DurationRule() {
            @Override public Level getLevel() { return Level.FAIL; }
            @Override public long getDuration() { return millis; }
        };
    }

    // --------------------------------------------------------------- plumbing

    /** A chain of exactly one item over an XmlSubXCV. */
    static class SingleSubXCVChain extends Chain<XmlSubXCV> {
        private final BiFunction<XmlSubXCV, LevelRule, ChainItem<XmlSubXCV>> factory;

        SingleSubXCVChain(I18nProvider i18nProvider, BiFunction<XmlSubXCV, LevelRule, ChainItem<XmlSubXCV>> factory) {
            super(i18nProvider, new XmlSubXCV());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    /** A chain of exactly one item over an XmlXCV. */
    static class SingleXCVChain extends Chain<XmlXCV> {
        private final BiFunction<XmlXCV, LevelRule, ChainItem<XmlXCV>> factory;

        SingleXCVChain(I18nProvider i18nProvider, BiFunction<XmlXCV, LevelRule, ChainItem<XmlXCV>> factory) {
            super(i18nProvider, new XmlXCV());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    /** A chain of exactly one item over an XmlRFC. */
    static class SingleRFCChain extends Chain<XmlRFC> {
        private final BiFunction<XmlRFC, LevelRule, ChainItem<XmlRFC>> factory;

        SingleRFCChain(I18nProvider i18nProvider, BiFunction<XmlRFC, LevelRule, ChainItem<XmlRFC>> factory) {
            super(i18nProvider, new XmlRFC());
            this.factory = factory;
        }

        @Override protected void initChain() {
            firstItem = factory.apply(result, FAIL);
        }
    }

    private static void subRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                               BiFunction<XmlSubXCV, LevelRule, ChainItem<XmlSubXCV>> factory, boolean unused) {
        try {
            XmlSubXCV result = new SingleSubXCVChain(i18n, factory).execute();
            out.println(row(file, token, check, "SUB_XCV", result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    private static void multiSubRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                                    BiFunction<XmlSubXCV, LevelRule, ChainItem<XmlSubXCV>> factory) {
        subRow(out, i18n, file, token, check, factory, true);
    }

    private static void xcvRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                               BiFunction<XmlXCV, LevelRule, ChainItem<XmlXCV>> factory) {
        try {
            XmlXCV result = new SingleXCVChain(i18n, factory).execute();
            out.println(row(file, token, check, "XCV", result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    private static void rfcRow(PrintWriter out, I18nProvider i18n, String file, String token, String check,
                               BiFunction<XmlRFC, LevelRule, ChainItem<XmlRFC>> factory) {
        try {
            XmlRFC result = new SingleRFCChain(i18n, factory).execute();
            out.println(row(file, token, check, "RFC", result));
        } catch (RuntimeException e) {
            System.err.println("SKIPPED " + file + " " + token + " " + check + ": " + e);
        }
    }

    // ------------------------------------------------------------------- JSON

    private static String row(String file, String token, String check, String block, XmlConstraintsConclusion result) {
        StringBuilder out = new StringBuilder("{");
        key(out, "file").append(str(file)).append(',');
        key(out, "token").append(str(token)).append(',');
        key(out, "check").append(str(check)).append(',');
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

    private static PrintWriter writer(File repoRoot, String relative) throws Exception {
        File out = new File(repoRoot, relative);
        out.getParentFile().mkdirs();
        return new PrintWriter(out, "UTF-8");
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
