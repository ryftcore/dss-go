/*
 * Regenerates spi/testdata/certificate_extensions/kat.tsv, the known answers of
 * CertificateExtensionsUtils, QcStatementUtils and SignerIdentifier.
 *
 * It runs the upstream DSS 6.5.RC1 algorithms verbatim on top of BouncyCastle and the JDK
 * X500Principal. Certificates are loaded through the BouncyCastle JCE provider because DSSUtils
 * does (DSSCertificateTokenSecurityFactory -> DSSSecurityProvider), which matters: BouncyCastle's
 * X509Certificate#getSubjectAlternativeNames() differs from the JDK's.
 *
 *   javac -cp bcprov-jdk18on-1.78.1.jar -d classes KatGen.java
 *   java  -cp classes:bcprov-jdk18on-1.78.1.jar KatGen .. ../kat.tsv
 *
 * The certificates in .. are the ones upstream's CertificateExtensionUtilsTest and
 * QcStatementsUtilsTest exercise, the remaining dss-spi test resources, and the synthetic
 * certificates built by synthetic_certificates.go.
 */
import org.bouncycastle.asn1.*;
import org.bouncycastle.asn1.x500.RDN;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x500.style.IETFUtils;
import org.bouncycastle.asn1.x500.style.RFC4519Style;
import org.bouncycastle.asn1.x509.*;
import org.bouncycastle.asn1.x509.qualified.*;
import org.bouncycastle.asn1.ocsp.OCSPObjectIdentifiers;

import javax.security.auth.x500.X500Principal;
import java.io.*;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.security.cert.*;
import java.util.*;

/** Regenerates known answers for CertificateExtensionsUtils / QcStatementUtils / SignerIdentifier. */
public class KatGen {

    static PrintStream out;
    static int idx;

    static final String OID_VALASSURED = "0.4.0.194121.2.1";
    static final String OID_QC_CCLEG = "0.4.0.1862.1.7";
    static final String OID_QC_IDENTMETHOD = "0.4.0.1862.1.8";
    static final String OID_QC_QSCDLEG = "0.4.0.1862.1.9";
    static final String OID_QC_PSB = "0.4.0.194126.1.3";
    static final String OID_PSD2 = "0.4.0.19495.2";

    public static void main(String[] args) throws Exception {
        Path dir = Paths.get(args[0]);
        out = new PrintStream(new FileOutputStream(args[1]), false, "UTF-8");
        List<Path> certs = new ArrayList<>();
        try (DirectoryStream<Path> ds = Files.newDirectoryStream(dir, "*.der")) { for (Path p : ds) certs.add(p); }
        Collections.sort(certs);
        java.security.Provider bc = new org.bouncycastle.jce.provider.BouncyCastleProvider();
        CertificateFactory cf = CertificateFactory.getInstance("X.509", bc);
        for (Path p : certs) {
            X509Certificate c;
            try (InputStream in = Files.newInputStream(p)) { c = (X509Certificate) cf.generateCertificate(in); }
            emit("file", p.getFileName().toString());
            dump(c);
            idx++;
        }
        out.flush();
        out.close();
    }

    static void emit(String key, String value) {
        String v = value == null ? "-" : "b" + Base64.getEncoder().encodeToString(value.getBytes(StandardCharsets.UTF_8));
        out.println(idx + "\t" + key + "\t" + v);
    }
    static void emitBin(String key, byte[] value) { emit(key, value == null ? null : hex(value)); }
    static void emitInt(String key, int v) { emit(key, Integer.toString(v)); }
    static void emitBool(String key, boolean v) { emit(key, Boolean.toString(v)); }
    static String hex(byte[] b) { StringBuilder sb = new StringBuilder(); for (byte x : b) sb.append(String.format("%02x", x)); return sb.toString(); }

    static void dump(X509Certificate c) {
        // ---- getCertificateExtensions() ordering ----
        List<String> ordered = new ArrayList<>();
        Set<String> crit = c.getCriticalExtensionOIDs();
        Set<String> ncrit = c.getNonCriticalExtensionOIDs();
        if (crit != null) ordered.addAll(crit);
        if (ncrit != null) ordered.addAll(ncrit);
        emitInt("allOids.count", ordered.size());
        for (int i = 0; i < ordered.size(); i++) emit("allOids." + i, ordered.get(i));
        emit("critSetClass", crit == null ? null : crit.getClass().getName());

        // ---- SubjectAlternativeNames ----
        subjectAlternativeNames(c);
        // ---- AuthorityInformationAccess ----
        authorityInformationAccess(c);
        // ---- AuthorityKeyIdentifier ----
        authorityKeyIdentifier(c);
        // ---- SubjectKeyIdentifier ----
        subjectKeyIdentifier(c);
        // ---- CRLDistributionPoints / FreshestCRL ----
        crlDistributionPoints(c, "2.5.29.31", "crldp");
        crlDistributionPoints(c, "2.5.29.46", "freshest");
        // ---- BasicConstraints ----
        emitBin("bc.octets", c.getExtensionValue("2.5.29.19"));
        emitBool("bc.critical", isCritical(c, "2.5.29.19"));
        emitInt("bc.value", c.getBasicConstraints());
        // ---- NameConstraints ----
        nameConstraints(c);
        // ---- PolicyConstraints ----
        policyConstraints(c);
        // ---- InhibitAnyPolicy ----
        inhibitAnyPolicy(c);
        // ---- KeyUsage ----
        keyUsage(c);
        // ---- ExtendedKeyUsage ----
        extendedKeyUsage(c);
        // ---- CertificatePolicies ----
        certificatePolicies(c);
        // ---- ocsp-nocheck / valassured / noRevAvail ----
        nullFlag(c, OCSPObjectIdentifiers.id_pkix_ocsp_nocheck.getId(), "onc");
        nullFlag(c, OID_VALASSURED, "vast");
        nullFlag(c, org.bouncycastle.asn1.x509.Extension.noRevAvail.getId(), "nra");
        // ---- QcStatements ----
        qcStatements(c);
        // ---- SignerIdentifier ----
        signerIdentifier(c);
    }

    static boolean isCritical(X509Certificate c, String oid) {
        Set<String> s = c.getCriticalExtensionOIDs();
        return s != null && s.contains(oid);
    }

    // ============ SubjectAlternativeNames ============
    static void subjectAlternativeNames(X509Certificate c) {
        try {
            emitBin("san.octets", c.getExtensionValue("2.5.29.17"));
            List<String[]> result = new ArrayList<>();
            Collection<List<?>> sans = c.getSubjectAlternativeNames();
            if (sans != null && !sans.isEmpty()) {
                for (List<?> gn : sans) {
                    String[] g = generalName(gn);
                    if (g != null) result.add(g);
                }
            }
            emit("san", "present");
            emitBool("san.critical", isCritical(c, "2.5.29.17"));
            emitInt("san.count", result.size());
            for (int i = 0; i < result.size(); i++) {
                emit("san." + i + ".type", result.get(i)[0]);
                emit("san." + i + ".value", result.get(i)[1]);
            }
        } catch (Exception e) {
            emit("san", null);
        }
    }

    static final String[] GN_TYPES = {"OTHER_NAME","RFC822_NAME","DNS_NAME","X400_ADDRESS","DIRECTORY_NAME",
            "EDI_PARTY_NAME","UNIFORM_RESOURCE_IDENTIFIER","IP_ADDRESS","REGISTERED_ID"};

    static String[] generalName(List<?> gn) {
        if (gn != null && gn.size() == 2) {
            try {
                if (!(gn.get(0) instanceof Integer)) return null;
                int t = (Integer) gn.get(0);
                String type = (t >= 0 && t < GN_TYPES.length) ? GN_TYPES[t] : null;
                Object value = gn.get(1);
                if (value instanceof String) {
                    String s = (String) value;
                    if ("DIRECTORY_NAME".equals(type)) s = toRFC2253RDN(s);
                    return new String[]{type, s};
                } else if (value instanceof byte[]) {
                    return new String[]{type, "#" + hex((byte[]) value)};
                }
                return null;
            } catch (Exception e) { return null; }
        }
        return null;
    }

    static String toRFC2253RDN(String value) {
        try {
            RDN[] rdns = RFC4519Style.INSTANCE.fromString(value);
            X500Principal p = new X500Principal(new X500Name(rdns).getEncoded());
            return p.getName(X500Principal.RFC2253);
        } catch (Exception e) { return value; }
    }

    // ============ AuthorityInformationAccess ============
    static void authorityInformationAccess(X509Certificate c) {
        byte[] v = c.getExtensionValue("1.3.6.1.5.5.7.1.1");
        if (v == null || v.length == 0) { emit("aia", null); emit("caIssuersUrls.count", "0"); emit("ocspUrls.count", "0"); return; }
        try {
            ASN1Sequence seq = seqFromOctet(v);
            if (seq == null || seq.size() == 0) { emit("aia", null); emit("caIssuersUrls.count", "0"); emit("ocspUrls.count", "0"); return; }
            org.bouncycastle.asn1.x509.AuthorityInformationAccess aia =
                    org.bouncycastle.asn1.x509.AuthorityInformationAccess.getInstance(seq);
            AccessDescription[] ads = aia.getAccessDescriptions();
            List<String> ca = accessUrls(ads, X509ObjectIdentifiers.id_ad_caIssuers);
            List<String> ocsp = accessUrls(ads, X509ObjectIdentifiers.id_ad_ocsp);
            emit("aia", "present");
            emitBin("aia.octets", v);
            emitBool("aia.critical", isCritical(c, "1.3.6.1.5.5.7.1.1"));
            emitList("aia.caIssuers", ca);
            emitList("aia.ocsp", ocsp);
            emitList("caIssuersUrls", ca);
            emitList("ocspUrls", ocsp);
        } catch (Exception e) {
            emit("aia", null); emit("caIssuersUrls.count", "0"); emit("ocspUrls.count", "0");
        }
    }

    static void emitList(String prefix, List<String> l) {
        emitInt(prefix + ".count", l.size());
        for (int i = 0; i < l.size(); i++) emit(prefix + "." + i, l.get(i));
    }

    static List<String> accessUrls(AccessDescription[] ads, ASN1ObjectIdentifier oid) {
        List<String> urls = new ArrayList<>();
        for (AccessDescription ad : ads) {
            if (oid.equals(ad.getAccessMethod())) {
                String s = parseGn(ad.getAccessLocation());
                if (s != null) urls.add(s);
            }
        }
        return urls;
    }

    static String parseGn(GeneralName gn) {
        try {
            if (GeneralName.uniformResourceIdentifier == gn.getTagNo()) {
                ASN1String s = (ASN1String) ((DERTaggedObject) gn.toASN1Primitive()).getBaseObject();
                return s.getString();
            }
        } catch (Exception e) { /* ignore */ }
        return null;
    }

    // ============ AKI / SKI ============
    static void authorityKeyIdentifier(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.35");
        if (v == null || v.length == 0) { emit("aki", null); return; }
        try {
            ASN1Primitive prim = parseExtensionValue(v);
            org.bouncycastle.asn1.x509.AuthorityKeyIdentifier aki =
                    org.bouncycastle.asn1.x509.AuthorityKeyIdentifier.getInstance(prim);
            emit("aki", "present");
            emitBin("aki.octets", v);
            emitBool("aki.critical", isCritical(c, "2.5.29.35"));
            emitBin("aki.keyIdentifier", aki.getKeyIdentifier());
            if (aki.getAuthorityCertIssuer() != null && aki.getAuthorityCertSerialNumber() != null) {
                IssuerSerial is = new IssuerSerial(aki.getAuthorityCertIssuer(), aki.getAuthorityCertSerialNumber());
                emitBin("aki.issuerSerial", is.toASN1Primitive().getEncoded(ASN1Encoding.DER));
            } else {
                emit("aki.issuerSerial", null);
            }
        } catch (Exception e) {
            emit("aki", "ERROR");
        }
    }

    static void subjectKeyIdentifier(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.14");
        if (v == null || v.length == 0) { emit("ski", null); return; }
        try {
            ASN1Primitive prim = parseExtensionValue(v);
            org.bouncycastle.asn1.x509.SubjectKeyIdentifier ski =
                    org.bouncycastle.asn1.x509.SubjectKeyIdentifier.getInstance(prim);
            emit("ski", "present");
            emitBin("ski.octets", v);
            emitBool("ski.critical", isCritical(c, "2.5.29.14"));
            emitBin("ski.ski", ski.getKeyIdentifier());
        } catch (Exception e) { emit("ski", "ERROR"); }
    }

    static ASN1Primitive parseExtensionValue(byte[] enc) throws IOException {
        return ASN1Primitive.fromByteArray(ASN1OctetString.getInstance(enc).getOctets());
    }

    // ============ CRLDistributionPoints ============
    static void crlDistributionPoints(X509Certificate c, String oid, String key) {
        byte[] v = c.getExtensionValue(oid);
        if (v == null) { emit(key, null); if ("crldp".equals(key)) emit("crlUrls.count", "0"); return; }
        List<String> urls = crlUrls(v);
        emit(key, "present");
        emitBin(key + ".octets", v);
        emitBool(key + ".critical", isCritical(c, oid));
        emitList(key + ".urls", urls);
        if ("crldp".equals(key)) emitList("crlUrls", urls);
    }

    static List<String> crlUrls(byte[] bytes) {
        try {
            List<String> urls = new ArrayList<>();
            ASN1Sequence seq = seqFromOctet(bytes);
            CRLDistPoint dp = CRLDistPoint.getInstance(seq);
            for (DistributionPoint d : dp.getDistributionPoints()) {
                DistributionPointName dpn = d.getDistributionPoint();
                if (DistributionPointName.FULL_NAME != dpn.getType()) continue;
                GeneralNames gns = (GeneralNames) dpn.getName();
                for (GeneralName n : gns.getNames()) {
                    String s = parseGn(n);
                    if (s != null) urls.add(s);
                }
            }
            return urls;
        } catch (Exception e) { return Collections.emptyList(); }
    }

    // ============ NameConstraints ============
    static void nameConstraints(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.30");
        if (v == null || v.length == 0) { emit("nc", null); return; }
        try {
            ASN1Sequence seq = seqFromOctet(v);
            org.bouncycastle.asn1.x509.NameConstraints nc = org.bouncycastle.asn1.x509.NameConstraints.getInstance(seq);
            emit("nc", "present");
            emitBin("nc.octets", v);
            emitBool("nc.critical", isCritical(c, "2.5.29.30"));
            subtrees("nc.permitted", nc.getPermittedSubtrees());
            subtrees("nc.excluded", nc.getExcludedSubtrees());
        } catch (Exception e) { emit("nc", null); }
    }

    static void subtrees(String prefix, org.bouncycastle.asn1.x509.GeneralSubtree[] sts) {
        if (sts == null || sts.length == 0) { emitInt(prefix + ".count", 0); return; }
        List<String[]> rows = new ArrayList<>();
        for (org.bouncycastle.asn1.x509.GeneralSubtree st : sts) {
            GeneralName gn = st.getBase();
            int t = gn.getTagNo();
            if (t < 0 || t >= GN_TYPES.length) continue;
            String type = GN_TYPES[t];
            BigInteger min = st.getMinimum();
            BigInteger max = st.getMaximum();
            rows.add(new String[]{type, min == null ? null : min.toString(), max == null ? null : max.toString(),
                    stringValue(type, gn.getName())});
        }
        emitInt(prefix + ".count", rows.size());
        for (int i = 0; i < rows.size(); i++) {
            emit(prefix + "." + i + ".type", rows.get(i)[0]);
            emit(prefix + "." + i + ".min", rows.get(i)[1]);
            emit(prefix + "." + i + ".max", rows.get(i)[2]);
            emit(prefix + "." + i + ".value", rows.get(i)[3]);
        }
    }

    static String stringValue(String type, ASN1Encodable v) {
        try {
            switch (type) {
                case "OTHER_NAME": case "EDI_PARTY_NAME": case "X400_ADDRESS":
                    return "#" + hex(v.toASN1Primitive().getEncoded(ASN1Encoding.DER));
                case "RFC822_NAME": case "DNS_NAME": case "UNIFORM_RESOURCE_IDENTIFIER":
                    if (v instanceof ASN1String) return ((ASN1String) v).getString();
                    return "#" + hex(v.toASN1Primitive().getEncoded(ASN1Encoding.DER));
                case "DIRECTORY_NAME": {
                    X500Principal p = new X500Principal(v.toASN1Primitive().getEncoded(ASN1Encoding.DER));
                    return p.getName(X500Principal.RFC2253);
                }
                case "IP_ADDRESS":
                    return "#" + hex(ASN1OctetString.getInstance(v).getOctets());
                case "REGISTERED_ID":
                    return ASN1ObjectIdentifier.getInstance(v).getId();
                default:
                    return "#" + hex(v.toASN1Primitive().getEncoded(ASN1Encoding.DER));
            }
        } catch (Exception e) {
            try { return "#" + hex(v.toASN1Primitive().getEncoded(ASN1Encoding.DER)); } catch (Exception e2) { return null; }
        }
    }

    // ============ PolicyConstraints ============
    static void policyConstraints(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.36");
        if (v == null || v.length == 0) { emit("pc", null); return; }
        try {
            ASN1Sequence seq = seqFromOctet(v);
            org.bouncycastle.asn1.x509.PolicyConstraints pc = org.bouncycastle.asn1.x509.PolicyConstraints.getInstance(seq);
            int inhibit = -1, require = -1;
            if (pc.getInhibitPolicyMapping() != null) inhibit = pc.getInhibitPolicyMapping().intValue();
            if (pc.getRequireExplicitPolicyMapping() != null) require = pc.getRequireExplicitPolicyMapping().intValue();
            emit("pc", "present");
            emitBin("pc.octets", v);
            emitBool("pc.critical", isCritical(c, "2.5.29.36"));
            emitInt("pc.requireExplicitPolicy", require);
            emitInt("pc.inhibitPolicyMapping", inhibit);
        } catch (Exception e) { emit("pc", null); }
    }

    // ============ InhibitAnyPolicy ============
    static void inhibitAnyPolicy(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.54");
        if (v == null || v.length == 0) { emit("iap", null); return; }
        try {
            ASN1Integer i = (ASN1Integer) ASN1Primitive.fromByteArray(ASN1OctetString.getInstance(v).getOctets());
            if (i != null && i.getValue() != null) {
                emit("iap", "present");
                emitBin("iap.octets", v);
                emitBool("iap.critical", isCritical(c, "2.5.29.54"));
                emitInt("iap.value", i.getValue().intValue());
                return;
            }
        } catch (Exception e) { /* fallthrough */ }
        emit("iap", null);
    }

    // ============ KeyUsage ============
    static void keyUsage(X509Certificate c) {
        boolean[] ku = c.getKeyUsage();
        if (ku == null) { emit("ku", null); return; }
        String[] names = {"DIGITAL_SIGNATURE","NON_REPUDIATION","KEY_ENCIPHERMENT","DATA_ENCIPHERMENT",
                "KEY_AGREEMENT","KEY_CERT_SIGN","CRL_SIGN","ENCIPHER_ONLY","DECIPHER_ONLY"};
        List<String> bits = new ArrayList<>();
        for (int i = 0; i < names.length; i++) if (ku[i]) bits.add(names[i]);
        emit("ku", "present");
        emitBin("ku.octets", c.getExtensionValue("2.5.29.15"));
        emitBool("ku.critical", isCritical(c, "2.5.29.15"));
        emitList("ku.bits", bits);
    }

    // ============ ExtendedKeyUsage ============
    static void extendedKeyUsage(X509Certificate c) {
        try {
            List<String> oids = c.getExtendedKeyUsage();
            emit("eku", "present");
            emitBin("eku.octets", c.getExtensionValue("2.5.29.37"));
            emitBool("eku.critical", isCritical(c, "2.5.29.37"));
            if (oids == null) emit("eku.oids.count", null);
            else emitList("eku.oids", oids);
        } catch (CertificateParsingException e) { emit("eku", null); }
    }

    // ============ CertificatePolicies ============
    static void certificatePolicies(X509Certificate c) {
        byte[] v = c.getExtensionValue("2.5.29.32");
        if (v == null || v.length == 0) { emit("cp", null); return; }
        try {
            ASN1Sequence seq = seqFromOctet(v);
            List<String[]> pols = new ArrayList<>();
            for (int i = 0; i < seq.size(); i++) {
                PolicyInformation pi = PolicyInformation.getInstance(seq.getObjectAt(i));
                String oid = pi.getPolicyIdentifier().getId();
                String cps = null;
                ASN1Sequence q = pi.getPolicyQualifiers();
                if (q != null) {
                    for (int j = 0; j < q.size(); j++) {
                        org.bouncycastle.asn1.x509.PolicyQualifierInfo pqi = org.bouncycastle.asn1.x509.PolicyQualifierInfo.getInstance(q.getObjectAt(j));
                        if (PolicyQualifierId.id_qt_cps.equals(pqi.getPolicyQualifierId())) cps = getString(pqi.getQualifier());
                    }
                }
                pols.add(new String[]{oid, cps});
            }
            emit("cp", "present");
            emitBin("cp.octets", v);
            emitBool("cp.critical", isCritical(c, "2.5.29.32"));
            emitInt("cp.count", pols.size());
            for (int i = 0; i < pols.size(); i++) {
                emit("cp." + i + ".oid", pols.get(i)[0]);
                emit("cp." + i + ".cpsUrl", pols.get(i)[1]);
            }
        } catch (Exception e) { emit("cp", null); }
    }

    // ============ null-identified flags ============
    static void nullFlag(X509Certificate c, String oid, String key) {
        byte[] v = c.getExtensionValue(oid);
        if (v == null) { emit(key, null); return; }
        boolean present = false;
        try {
            ASN1Primitive p = ASN1Primitive.fromByteArray(v);
            if (p instanceof DEROctetString) {
                ASN1Primitive inner = ASN1Primitive.fromByteArray(((DEROctetString) p).getOctets());
                present = DERNull.INSTANCE.equals(inner);
            }
        } catch (Exception e) { /* false */ }
        emit(key, "present");
        emitBin(key + ".octets", v);
        emitBool(key + ".critical", isCritical(c, oid));
        emitBool(key + ".value", present);
    }

    // ============ QcStatements ============
    static void qcStatements(X509Certificate c) {
        byte[] v = c.getExtensionValue(org.bouncycastle.asn1.x509.Extension.qCStatements.getId());
        if (v == null || v.length == 0) { emit("qc", null); return; }
        ASN1Sequence seq;
        try { seq = seqFromOctet(v); } catch (Exception e) { emit("qc", null); return; }
        try {
            emit("qc", "present");
            emitBin("qc.octets", seq.toASN1Primitive().getEncoded(ASN1Encoding.DER));
            emitBool("qc.critical", isCritical(c, org.bouncycastle.asn1.x509.Extension.qCStatements.getId()));

            boolean qcCompliance = false, qcQSCD = false;
            String[] limit = null; Integer retention = null;
            List<String[]> pds = new ArrayList<>();
            List<String> types = new ArrayList<>();
            List<String> ccleg = new ArrayList<>(), qscdleg = new ArrayList<>();
            String semantics = null, identMethod = null;
            String[] psd2 = null; List<String[]> roles = null;
            String[] psb = null;
            List<String> others = new ArrayList<>();

            for (int i = 0; i < seq.size(); i++) {
                QCStatement st;
                try { st = QCStatement.getInstance(seq.getObjectAt(i)); } catch (Exception e) { continue; }
                String oid = st.getStatementId().getId();
                ASN1Encodable info = st.getStatementInfo();
                if (ETSIQCObjectIdentifiers.id_etsi_qcs_QcCompliance.getId().equals(oid)) qcCompliance = true;
                else if (ETSIQCObjectIdentifiers.id_etsi_qcs_LimiteValue.getId().equals(oid)) limit = qcLimit(info);
                else if (ETSIQCObjectIdentifiers.id_etsi_qcs_RetentionPeriod.getId().equals(oid)) retention = qcRetention(info);
                else if (ETSIQCObjectIdentifiers.id_etsi_qcs_QcSSCD.getId().equals(oid)) qcQSCD = true;
                else if (ETSIQCObjectIdentifiers.id_etsi_qcs_QcPds.getId().equals(oid)) pds = qcPds(info);
                else if (ETSIQCObjectIdentifiers.id_etsi_qcs_QcType.getId().equals(oid)) types = qcTypes(info);
                else if (OID_QC_CCLEG.equals(oid)) ccleg = qcLegislation(info);
                else if (RFC3739QCObjectIdentifiers.id_qcs_pkixQCSyntax_v2.getId().equals(oid)) semantics = qcSemantics(info);
                else if (OID_PSD2.equals(oid)) { String[][] r = psd2(info); if (r != null) { psd2 = r[0]; roles = new ArrayList<>(Arrays.asList(r).subList(1, r.length)); } }
                else if (OID_QC_QSCDLEG.equals(oid)) qscdleg = qcLegislation(info);
                else if (OID_QC_IDENTMETHOD.equals(oid)) identMethod = qcIdentMethod(info);
                else if (OID_QC_PSB.equals(oid)) psb = qcPSB(info);
                else others.add(oid);
            }
            emitBool("qc.qcCompliance", qcCompliance);
            emitBool("qc.qcQSCD", qcQSCD);
            if (limit == null) emit("qc.limit", null);
            else { emit("qc.limit", "present"); emit("qc.limit.currency", limit[0]); emit("qc.limit.amount", limit[1]); emit("qc.limit.exponent", limit[2]); }
            emit("qc.retentionPeriod", retention == null ? null : retention.toString());
            emitInt("qc.pds.count", pds.size());
            for (int i = 0; i < pds.size(); i++) { emit("qc.pds." + i + ".url", pds.get(i)[0]); emit("qc.pds." + i + ".language", pds.get(i)[1]); }
            emitList("qc.types", types);
            emitList("qc.ccLegislation", ccleg);
            emitList("qc.qscdLegislation", qscdleg);
            emit("qc.semanticsIdentifier", semantics);
            emit("qc.identMethod", identMethod);
            if (psd2 == null) emit("qc.psd2", null);
            else {
                emit("qc.psd2", "present");
                emit("qc.psd2.ncaName", psd2[0]);
                emit("qc.psd2.ncaId", psd2[1]);
                emitInt("qc.psd2.roles.count", roles.size());
                for (int i = 0; i < roles.size(); i++) { emit("qc.psd2.roles." + i + ".oid", roles.get(i)[0]); emit("qc.psd2.roles." + i + ".name", roles.get(i)[1]); }
            }
            if (psb == null) emit("qc.psb", null);
            else { emit("qc.psb", "present"); emit("qc.psb.country", psb[0]); emit("qc.psb.authSource", psb[1]); emit("qc.psb.legislation", psb[2]); }
            emitList("qc.otherOids", others);
        } catch (Exception e) { emit("qc", "ERROR:" + e); }
    }

    static String[] qcLimit(ASN1Encodable info) {
        try {
            MonetaryValue mv = MonetaryValue.getInstance(info);
            return new String[]{mv.getCurrency().getAlphabetic(), Integer.toString(mv.getAmount().intValue()), Integer.toString(mv.getExponent().intValue())};
        } catch (Exception e) { return null; }
    }
    static Integer qcRetention(ASN1Encodable info) {
        try { return ASN1Integer.getInstance(info).intValueExact(); } catch (Exception e) { return null; }
    }
    static List<String[]> qcPds(ASN1Encodable info) {
        List<String[]> r = new ArrayList<>();
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            for (int i = 0; i < s.size(); i++) {
                ASN1Encodable e1 = s.getObjectAt(i);
                if (e1 instanceof ASN1Sequence) {
                    ASN1Sequence q = (ASN1Sequence) e1;
                    r.add(new String[]{getString(q.getObjectAt(0)), getString(q.getObjectAt(1))});
                }
            }
        } catch (Exception e) { /* partial */ }
        return r;
    }
    static List<String> qcTypes(ASN1Encodable info) {
        List<String> oids = new ArrayList<>();
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            for (int i = 0; i < s.size(); i++) {
                ASN1Encodable e1 = s.getObjectAt(i);
                if (e1 instanceof ASN1ObjectIdentifier) oids.add(((ASN1ObjectIdentifier) e1).getId());
            }
        } catch (Exception e) { /* partial */ }
        return oids;
    }
    static List<String> qcLegislation(ASN1Encodable info) {
        List<String> r = new ArrayList<>();
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            for (int i = 0; i < s.size(); i++) { String cc = getString(s.getObjectAt(i)); if (cc != null) r.add(cc); }
        } catch (Exception e) { /* partial */ }
        return r;
    }
    static String qcSemantics(ASN1Encodable info) {
        try {
            SemanticsInformation si = SemanticsInformation.getInstance(info);
            if (si != null && si.getSemanticsIdentifier() != null) return si.getSemanticsIdentifier().getId();
        } catch (Exception e) { /* null */ }
        return null;
    }
    static String[][] psd2(ASN1Encodable info) {
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            ASN1Sequence rolesSeq = ASN1Sequence.getInstance(s.getObjectAt(0));
            List<String[]> rows = new ArrayList<>();
            for (int i = 0; i < rolesSeq.size(); i++) {
                ASN1Sequence one = ASN1Sequence.getInstance(rolesSeq.getObjectAt(i));
                ASN1ObjectIdentifier oid = (ASN1ObjectIdentifier) one.getObjectAt(0);
                rows.add(new String[]{oid.getId(), getString(one.getObjectAt(1))});
            }
            String nca = getString(s.getObjectAt(1));
            String ncaId = getString(s.getObjectAt(2));
            String[][] out = new String[rows.size() + 1][];
            out[0] = new String[]{nca, ncaId};
            for (int i = 0; i < rows.size(); i++) out[i + 1] = rows.get(i);
            return out;
        } catch (Exception e) { return null; }
    }
    static String qcIdentMethod(ASN1Encodable info) {
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            if (s.size() != 1) return null;
            ASN1Encodable e1 = s.getObjectAt(0);
            if (e1 instanceof ASN1ObjectIdentifier) return ((ASN1ObjectIdentifier) e1).getId();
        } catch (Exception e) { /* null */ }
        return null;
    }
    static String[] qcPSB(ASN1Encodable info) {
        try {
            ASN1Sequence s = ASN1Sequence.getInstance(info);
            if (s.size() != 3) return null;
            return new String[]{getString(s.getObjectAt(0)), getString(s.getObjectAt(1)), getString(s.getObjectAt(2))};
        } catch (Exception e) { return null; }
    }

    static String getString(ASN1Encodable v) {
        if (v == null) return null;
        try { return IETFUtils.valueToString(v).trim(); } catch (Exception e) { return null; }
    }

    // ============ SignerIdentifier ============
    static void signerIdentifier(X509Certificate c) {
        try {
            X500Principal issuer = c.getIssuerX500Principal();
            BigInteger serial = c.getSerialNumber();
            X500Name n = X500Name.getInstance(issuer.getEncoded());
            GeneralNames gns = new GeneralNames(new GeneralName(n));
            IssuerSerial is = new IssuerSerial(gns, serial);
            emitBin("signerId.issuerSerial", is.toASN1Primitive().getEncoded(ASN1Encoding.DER));
            emit("signerId.toString", "IssuerSerialInfo [issuerName=" + issuer.getName(X500Principal.RFC2253) + ", serialNumber=" + serial + "]");
            // DSSASN1Utils.get(X500Principal) attribute map, sorted for stability
            emitList("signerId.issuerMap", attrMap(issuer));
        } catch (Exception e) { emit("signerId.issuerSerial", null); }
    }

    static List<String> attrMap(X500Principal p) {
        Map<String, String> m = new HashMap<>();
        ASN1Sequence seq = ASN1Sequence.getInstance(p.getEncoded());
        for (ASN1Encodable e : seq.toArray()) {
            ASN1Set set = ASN1Set.getInstance(e);
            for (int i = 0; i < set.size(); i++) {
                ASN1Sequence s = ASN1Sequence.getInstance(set.getObjectAt(i));
                if (s.size() != 2) throw new RuntimeException("The DLSequence must contains exactly 2 elements.");
                m.put(getString(s.getObjectAt(0)), getString(s.getObjectAt(1)));
            }
        }
        List<String> rows = new ArrayList<>();
        for (Map.Entry<String, String> en : m.entrySet()) rows.add(en.getKey() + "=" + en.getValue());
        Collections.sort(rows);
        return rows;
    }

    static ASN1Sequence seqFromOctet(byte[] bytes) throws IOException {
        byte[] content = ((DEROctetString) ASN1Primitive.fromByteArray(bytes)).getOctets();
        return (ASN1Sequence) ASN1Primitive.fromByteArray(content);
    }
}
