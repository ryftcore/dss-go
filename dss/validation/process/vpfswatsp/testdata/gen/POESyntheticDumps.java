/*
 * Builds the synthetic XmlDiagnosticData dumps the POE oracle needs.
 *
 * The 54-dump marshal-parity corpus in dss/diagnostic/jaxb/testdata/oracle carries
 * no evidence record at all, no time-stamp pair sharing a production time, and no
 * time-stamp whose message imprint is broken, so POEComparator's compareByType,
 * compareByTimestampType and compareByTimestampedReferences tie-breakers, and
 * POEExtraction#extractPOE's message-imprint gate, are never reached by it. These
 * dumps supply exactly those inputs, plus the three scenarios EN 319 102-1 5.6.2
 * turns on: an expired revocation revived by a POE, a chain of archive
 * time-stamps each covering the previous one, and an evidence record covering the
 * signature.
 *
 * They are written once, with the upstream DiagnosticDataFacade, to
 *   dss/validation/process/vpfswatsp/testdata/dd/*.xml
 * and committed. Both POEOracle and the Go test read those very files, so neither
 * side gets a private fixture.
 *
 * Regenerate (from a dss-upstream checkout with the modules built):
 *
 *   CP="dss-validation/target/classes:$(cat cp.txt)"
 *   javac -cp "$CP" -d /tmp/oracle POESyntheticDumps.java
 *   java  -cp "$CP:/tmp/oracle" POESyntheticDumps <dss-repo-root>
 */

import eu.europa.esig.dss.diagnostic.DiagnosticDataFacade;
import eu.europa.esig.dss.diagnostic.jaxb.XmlBasicSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlCertificateRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlChainItem;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDiagnosticData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlDigestMatcher;
import eu.europa.esig.dss.diagnostic.jaxb.XmlEvidenceRecord;
import eu.europa.esig.dss.diagnostic.jaxb.XmlFoundTimestamp;
import eu.europa.esig.dss.diagnostic.jaxb.XmlOrphanCertificateToken;
import eu.europa.esig.dss.diagnostic.jaxb.XmlOrphanTokens;
import eu.europa.esig.dss.diagnostic.jaxb.XmlRevocation;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignature;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSignerData;
import eu.europa.esig.dss.diagnostic.jaxb.XmlSigningCertificate;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTimestamp;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTimestampedObject;
import eu.europa.esig.dss.diagnostic.jaxb.XmlTrusted;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.DigestMatcherType;
import eu.europa.esig.dss.enumerations.EncryptionAlgorithm;
import eu.europa.esig.dss.enumerations.EvidenceRecordOrigin;
import eu.europa.esig.dss.enumerations.EvidenceRecordTypeEnum;
import eu.europa.esig.dss.enumerations.RevocationType;
import eu.europa.esig.dss.enumerations.CertificateStatus;
import eu.europa.esig.dss.enumerations.RevocationOrigin;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.enumerations.TimestampedObjectType;

import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

public class POESyntheticDumps {

    public static void main(String[] args) throws Exception {
        File repoRoot = new File(args[0]);
        File outDir = new File(repoRoot, "validation/process/vpfswatsp/testdata/dd");
        outDir.mkdirs();

        write(outDir, "poe-timestamp-chain.xml", timestampChain());
        write(outDir, "poe-evidence-record.xml", evidenceRecord());
        write(outDir, "poe-revoked-revived.xml", revokedRevived());
        write(outDir, "poe-tie-timestamp-type.xml", tieTimestampType());
        write(outDir, "poe-tie-covered-count.xml", tieCoveredCount());
        write(outDir, "poe-broken-message-imprint.xml", brokenMessageImprint());
        write(outDir, "poe-orphans-and-signer-data.xml", orphansAndSignerData());
    }

    private static void write(File dir, String name, XmlDiagnosticData data) throws Exception {
        try (OutputStream out = new FileOutputStream(new File(dir, name))) {
            DiagnosticDataFacade.newFacade().marshall(data, out, false);
        }
    }

    // --------------------------------------------------------------- fixtures

    /**
     * Three archive time-stamps, each covering the signature, the signing
     * certificate and the previous time-stamp: the POE of the signature is the
     * oldest one, and the POE of every intermediate time-stamp comes from its
     * successor.
     */
    private static XmlDiagnosticData timestampChain() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp t1 = timestamp("T-1", "2016-06-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(t1, signature, TimestampedObjectType.SIGNATURE);
        covers(t1, signer, TimestampedObjectType.CERTIFICATE);

        XmlTimestamp t2 = timestamp("T-2", "2017-06-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(t2, signature, TimestampedObjectType.SIGNATURE);
        covers(t2, signer, TimestampedObjectType.CERTIFICATE);
        covers(t2, t1, TimestampedObjectType.TIMESTAMP);

        XmlTimestamp t3 = timestamp("T-3", "2018-06-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(t3, signature, TimestampedObjectType.SIGNATURE);
        covers(t3, signer, TimestampedObjectType.CERTIFICATE);
        covers(t3, t2, TimestampedObjectType.TIMESTAMP);

        XmlDiagnosticData data = data("poe-timestamp-chain", signature, Arrays.asList(signer, root));
        data.getUsedTimestamps().addAll(Arrays.asList(t1, t2, t3));
        return data;
    }

    /**
     * An evidence record covering the signature, its certificate and a revocation,
     * with two time-stamps of its own: EvidenceRecordPOE takes the production time
     * of the FIRST of them.
     */
    private static XmlDiagnosticData evidenceRecord() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlRevocation revocation = revocation("R-1", "2016-01-01T00:00:00Z", "2016-01-01T00:00:00Z", root);
        revoke(signer, revocation, null);
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp erT1 = timestamp("ER-T-1", "2016-03-01T00:00:00Z", TimestampType.EVIDENCE_RECORD_TIMESTAMP);
        XmlTimestamp erT2 = timestamp("ER-T-2", "2019-03-01T00:00:00Z", TimestampType.EVIDENCE_RECORD_TIMESTAMP);

        XmlEvidenceRecord er = new XmlEvidenceRecord();
        er.setId("ER-1");
        er.setType(EvidenceRecordTypeEnum.ASN1_EVIDENCE_RECORD);
        er.setOrigin(EvidenceRecordOrigin.EXTERNAL);
        er.setDocumentName("er.ers");
        XmlDigestMatcher dm = new XmlDigestMatcher();
        dm.setType(DigestMatcherType.EVIDENCE_RECORD_ARCHIVE_OBJECT);
        dm.setDataFound(true);
        dm.setDataIntact(true);
        dm.setDigestMethod(DigestAlgorithm.SHA256);
        dm.setDigestValue(new byte[] { 1, 2, 3 });
        dm.setDocumentName("doc.bin");
        er.getDigestMatchers().add(dm);
        for (XmlTimestamp t : Arrays.asList(erT1, erT2)) {
            XmlFoundTimestamp ft = new XmlFoundTimestamp();
            ft.setTimestamp(t);
            er.getEvidenceRecordTimestamps().add(ft);
        }
        covers(er, signature, TimestampedObjectType.SIGNATURE);
        covers(er, signer, TimestampedObjectType.CERTIFICATE);
        covers(er, revocation, TimestampedObjectType.REVOCATION);

        XmlDiagnosticData data = data("poe-evidence-record", signature, Arrays.asList(signer, root));
        data.getUsedRevocations().add(revocation);
        data.getUsedTimestamps().addAll(Arrays.asList(erT1, erT2));
        data.getEvidenceRecords().add(er);
        return data;
    }

    /**
     * A signing certificate revoked in 2017 with a revocation issued in 2017, and
     * an archive time-stamp from 2016 covering the signature, the certificate and
     * the revocation: the POE for all three predates the revocation time. This is
     * the "expired/revoked revived by POE" shape 5.6.2.4 step 3) turns on.
     */
    private static XmlDiagnosticData revokedRevived() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlRevocation revocation = revocation("R-1", "2017-06-01T00:00:00Z", "2017-06-01T00:00:00Z", root);
        revoke(signer, revocation, "2017-01-01T00:00:00Z");
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp t1 = timestamp("T-1", "2016-06-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(t1, signature, TimestampedObjectType.SIGNATURE);
        covers(t1, signer, TimestampedObjectType.CERTIFICATE);
        covers(t1, revocation, TimestampedObjectType.REVOCATION);

        XmlDiagnosticData data = data("poe-revoked-revived", signature, Arrays.asList(signer, root));
        data.getUsedRevocations().add(revocation);
        data.getUsedTimestamps().add(t1);
        return data;
    }

    /**
     * Two time-stamps with the SAME production time and DIFFERENT types, both
     * covering the signature: POEComparator falls through compareByTime and
     * compareByType into compareByTimestampType.
     */
    private static XmlDiagnosticData tieTimestampType() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp archive = timestamp("T-ARCHIVE", "2017-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(archive, signature, TimestampedObjectType.SIGNATURE);
        XmlTimestamp sigTst = timestamp("T-SIGNATURE", "2017-01-01T00:00:00Z", TimestampType.SIGNATURE_TIMESTAMP);
        covers(sigTst, signature, TimestampedObjectType.SIGNATURE);
        XmlTimestamp content = timestamp("T-CONTENT", "2017-01-01T00:00:00Z", TimestampType.CONTENT_TIMESTAMP);
        covers(content, signature, TimestampedObjectType.SIGNATURE);

        XmlDiagnosticData data = data("poe-tie-timestamp-type", signature, Arrays.asList(signer, root));
        data.getUsedTimestamps().addAll(Arrays.asList(archive, sigTst, content));
        return data;
    }

    /**
     * Two time-stamps with the same production time AND the same type, covering a
     * different number of objects: POEComparator falls all the way through to
     * compareByTimestampedReferences.
     */
    private static XmlDiagnosticData tieCoveredCount() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp narrow = timestamp("T-NARROW", "2017-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(narrow, signature, TimestampedObjectType.SIGNATURE);

        XmlTimestamp wide = timestamp("T-WIDE", "2017-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(wide, signature, TimestampedObjectType.SIGNATURE);
        covers(wide, signer, TimestampedObjectType.CERTIFICATE);
        covers(wide, root, TimestampedObjectType.CERTIFICATE);

        XmlDiagnosticData data = data("poe-tie-covered-count", signature, Arrays.asList(signer, root));
        data.getUsedTimestamps().addAll(Arrays.asList(narrow, wide));
        return data;
    }

    /**
     * Three time-stamps covering the signature, one with a message imprint that
     * was not found, one with a message imprint that is not intact, one valid:
     * extractPOE's gate accepts only the third.
     */
    private static XmlDiagnosticData brokenMessageImprint() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlSignature signature = signature("SIG-1", signer);

        XmlTimestamp notFound = timestamp("T-NOT-FOUND", "2016-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        notFound.getDigestMatchers().get(0).setDataFound(false);
        notFound.getDigestMatchers().get(0).setDataIntact(false);
        covers(notFound, signature, TimestampedObjectType.SIGNATURE);

        XmlTimestamp broken = timestamp("T-BROKEN", "2017-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        broken.getDigestMatchers().get(0).setDataIntact(false);
        covers(broken, signature, TimestampedObjectType.SIGNATURE);

        XmlTimestamp valid = timestamp("T-VALID", "2018-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(valid, signature, TimestampedObjectType.SIGNATURE);

        // covers nothing at all: the Utils.isCollectionNotEmpty guard
        XmlTimestamp empty = timestamp("T-EMPTY", "2015-01-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);

        XmlDiagnosticData data = data("poe-broken-message-imprint", signature, Arrays.asList(signer, root));
        data.getUsedTimestamps().addAll(Arrays.asList(notFound, broken, valid, empty));
        return data;
    }

    /**
     * Orphan certificate tokens and signer data, so that POEExtraction#init walks
     * the four getAllOrphan... lists and getAllSignerDocuments, which the corpus
     * leaves empty.
     */
    private static XmlDiagnosticData orphansAndSignerData() {
        XmlCertificate root = certificate("C-ROOT", "2010-01-01T00:00:00Z", "2035-01-01T00:00:00Z");
        trust(root);
        XmlCertificate signer = certificate("C-SIGNER", "2015-01-01T00:00:00Z", "2019-01-01T00:00:00Z");
        chain(signer, root);
        issuedBy(signer, root);
        XmlSignature signature = signature("SIG-1", signer);

        XmlSignerData signerData = new XmlSignerData();
        signerData.setId("D-1");
        signerData.setReferencedName("doc.bin");
        signature.getSignatureScopes().add(scope(signerData));

        XmlOrphanCertificateToken orphanCert = new XmlOrphanCertificateToken();
        orphanCert.setId("OC-1");
        XmlOrphanTokens orphanTokens = new XmlOrphanTokens();
        orphanTokens.getOrphanCertificates().add(orphanCert);

        XmlTimestamp t1 = timestamp("T-1", "2016-06-01T00:00:00Z", TimestampType.ARCHIVE_TIMESTAMP);
        covers(t1, signature, TimestampedObjectType.SIGNATURE);
        covers(t1, signerData, TimestampedObjectType.SIGNED_DATA);
        covers(t1, orphanCert, TimestampedObjectType.ORPHAN_CERTIFICATE);

        XmlDiagnosticData data = data("poe-orphans-and-signer-data", signature, Arrays.asList(signer, root));
        data.getOriginalDocuments().add(signerData);
        data.setOrphanTokens(orphanTokens);
        data.getUsedTimestamps().add(t1);
        return data;
    }

    // ---------------------------------------------------------------- helpers

    private static eu.europa.esig.dss.diagnostic.jaxb.XmlSignatureScope scope(XmlSignerData signerData) {
        eu.europa.esig.dss.diagnostic.jaxb.XmlSignatureScope s =
                new eu.europa.esig.dss.diagnostic.jaxb.XmlSignatureScope();
        s.setSignerData(signerData);
        s.setName(signerData.getReferencedName());
        return s;
    }

    private static XmlDiagnosticData data(String name, XmlSignature signature, List<XmlCertificate> certificates) {
        XmlDiagnosticData data = new XmlDiagnosticData();
        data.setDocumentName(name);
        data.setValidationDate(date("2024-01-01T00:00:00Z"));
        data.getSignatures().add(signature);
        data.getUsedCertificates().addAll(certificates);
        return data;
    }

    private static XmlSignature signature(String id, XmlCertificate signingCertificate) {
        XmlSignature signature = new XmlSignature();
        signature.setId(id);
        signature.setSigningCertificate(signingCertificateRef(signingCertificate));
        signature.setBasicSignature(basicSignature());
        signature.getCertificateChain().addAll(chainItems(signingCertificate));
        return signature;
    }

    private static XmlTimestamp timestamp(String id, String productionTime, TimestampType type) {
        XmlTimestamp timestamp = new XmlTimestamp();
        timestamp.setId(id);
        timestamp.setProductionTime(date(productionTime));
        timestamp.setType(type);
        XmlDigestMatcher messageImprint = new XmlDigestMatcher();
        messageImprint.setType(DigestMatcherType.MESSAGE_IMPRINT);
        messageImprint.setDataFound(true);
        messageImprint.setDataIntact(true);
        messageImprint.setDigestMethod(DigestAlgorithm.SHA256);
        messageImprint.setDigestValue(new byte[] { 4, 5, 6 });
        timestamp.getDigestMatchers().add(messageImprint);
        return timestamp;
    }

    private static void covers(XmlTimestamp timestamp,
            eu.europa.esig.dss.diagnostic.jaxb.XmlAbstractToken token, TimestampedObjectType category) {
        timestamp.getTimestampedObjects().add(timestampedObject(token, category));
    }

    private static void covers(XmlEvidenceRecord evidenceRecord,
            eu.europa.esig.dss.diagnostic.jaxb.XmlAbstractToken token, TimestampedObjectType category) {
        evidenceRecord.getTimestampedObjects().add(timestampedObject(token, category));
    }

    private static XmlTimestampedObject timestampedObject(
            eu.europa.esig.dss.diagnostic.jaxb.XmlAbstractToken token, TimestampedObjectType category) {
        XmlTimestampedObject o = new XmlTimestampedObject();
        o.setToken(token);
        o.setCategory(category);
        return o;
    }

    private static XmlCertificate certificate(String id, String notBefore, String notAfter) {
        XmlCertificate certificate = new XmlCertificate();
        certificate.setId(id);
        certificate.setNotBefore(date(notBefore));
        certificate.setNotAfter(date(notAfter));
        certificate.setBasicSignature(basicSignature());
        certificate.setSelfSigned(false);
        return certificate;
    }

    private static void trust(XmlCertificate certificate) {
        XmlTrusted trusted = new XmlTrusted();
        trusted.setValue(true);
        certificate.setTrusted(trusted);
        certificate.setSelfSigned(true);
    }

    private static void chain(XmlCertificate certificate, XmlCertificate... parents) {
        for (XmlCertificate parent : parents) {
            XmlChainItem item = new XmlChainItem();
            item.setCertificate(parent);
            certificate.getCertificateChain().add(item);
        }
    }

    private static List<XmlChainItem> chainItems(XmlCertificate certificate) {
        List<XmlChainItem> items = new ArrayList<>();
        XmlChainItem first = new XmlChainItem();
        first.setCertificate(certificate);
        items.add(first);
        for (XmlChainItem parent : certificate.getCertificateChain()) {
            XmlChainItem item = new XmlChainItem();
            item.setCertificate(parent.getCertificate());
            items.add(item);
        }
        return items;
    }

    private static void issuedBy(XmlCertificate certificate, XmlCertificate issuer) {
        certificate.setSigningCertificate(signingCertificateRef(issuer));
    }

    private static XmlSigningCertificate signingCertificateRef(XmlCertificate certificate) {
        XmlSigningCertificate signingCertificate = new XmlSigningCertificate();
        signingCertificate.setCertificate(certificate);
        return signingCertificate;
    }

    private static XmlRevocation revocation(String id, String productionDate, String thisUpdate,
            XmlCertificate issuer) {
        XmlRevocation revocation = new XmlRevocation();
        revocation.setId(id);
        revocation.setProductionDate(date(productionDate));
        revocation.setThisUpdate(date(thisUpdate));
        revocation.setOrigin(RevocationOrigin.EXTERNAL);
        revocation.setType(RevocationType.CRL);
        revocation.setBasicSignature(basicSignature());
        revocation.setSigningCertificate(signingCertificateRef(issuer));
        return revocation;
    }

    private static void revoke(XmlCertificate certificate, XmlRevocation revocation, String revocationDate) {
        XmlCertificateRevocation certificateRevocation = new XmlCertificateRevocation();
        certificateRevocation.setRevocation(revocation);
        if (revocationDate != null) {
            certificateRevocation.setStatus(CertificateStatus.REVOKED);
            certificateRevocation.setRevocationDate(date(revocationDate));
        } else {
            certificateRevocation.setStatus(CertificateStatus.GOOD);
        }
        certificate.getRevocations().add(certificateRevocation);
    }

    private static XmlBasicSignature basicSignature() {
        XmlBasicSignature basicSignature = new XmlBasicSignature();
        basicSignature.setDigestAlgoUsedToSignThisToken(DigestAlgorithm.SHA256);
        basicSignature.setEncryptionAlgoUsedToSignThisToken(EncryptionAlgorithm.RSA);
        basicSignature.setKeyLengthUsedToSignThisToken("2048");
        basicSignature.setSignatureIntact(true);
        basicSignature.setSignatureValid(true);
        return basicSignature;
    }

    private static Date date(String isoInstant) {
        return Date.from(Instant.parse(isoInstant));
    }

}
