// The BouncyCastle oracle behind testdata/bc-oracle.txt and
// testdata/adversarial/bc-oracle.txt: it loads every fixture through BouncyCastle's own CMS,
// TimeStampToken and TimeStampResponse classes and dumps, as "key=value" lines, every field
// internal/cmscore exposes. TestBouncyCastleOracle produces the same lines from the Go side
// and compares them, so the checked-in golden files are BouncyCastle's answers, not the Go
// port's own.
//
// Run it with OpenJDK 21 and BouncyCastle 1.84 from the local Maven repository:
//
//   BC=$HOME/.m2/repository/org/bouncycastle
//   CP=$BC/bcprov-jdk18on/1.84/bcprov-jdk18on-1.84.jar:$BC/bcpkix-jdk18on/1.84/bcpkix-jdk18on-1.84.jar:$BC/bcutil-jdk18on/1.84/bcutil-jdk18on-1.84.jar
//   javac -cp "$CP" -d /tmp/cmsoracle CmsOracle.java
//   java -cp "$CP:/tmp/cmsoracle" CmsOracle <directory> > bc-oracle.txt
//
// Values longer than 64 octets are pinned by their SHA-256 so the golden file stays readable;
// the Go side abbreviates them the same way.
// Independent BouncyCastle oracle: dumps every field internal/cmscore exposes, as
// "key=value" lines, for mechanical diffing against the Go dumper.
import org.bouncycastle.asn1.*;
import org.bouncycastle.asn1.cms.*;
import org.bouncycastle.asn1.cmp.PKIStatusInfo;
import org.bouncycastle.asn1.tsp.*;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.GeneralName;
import org.bouncycastle.cms.*;
import org.bouncycastle.tsp.*;
import org.bouncycastle.util.encoders.Hex;

import java.io.File;
import java.nio.file.Files;
import java.util.*;

public class CmsOracle {
    static StringBuilder out = new StringBuilder();

    static void p(String k, Object v) { out.append(k).append('=').append(v).append('\n'); }
    // Long values are pinned by their digest, so the golden file stays readable.
    static void ph(String k, byte[] v) {
        if (v == null) { p(k, "<nil>"); return; }
        if (v.length <= 64) { p(k, Hex.toHexString(v)); return; }
        try {
            p(k, "sha256:" + Hex.toHexString(
                java.security.MessageDigest.getInstance("SHA-256").digest(v)));
        } catch (Exception e) { throw new RuntimeException(e); }
    }

    public static void main(String[] args) throws Exception {
        String dir = args[0];
        List<String> names = new ArrayList<>(Arrays.asList(new File(dir).list()));
        Collections.sort(names);
        for (String n : names) {
            File f = new File(dir, n);
            if (!f.isFile()) continue;
            byte[] data = Files.readAllBytes(f.toPath());
            int mark = out.length();
            boolean compact = System.getenv("COMPACT") != null;
            try {
                if (n.endsWith(".p7s")) dumpCms(n, data);
                else if (n.endsWith(".tst")) dumpTst(n, data);
                else if (n.endsWith(".tsr")) dumpTsr(n, data);
                else continue;
                if (compact) {
                    String body = out.substring(mark);
                    out.setLength(mark);
                    p(n + ".sha", Hex.toHexString(
                        java.security.MessageDigest.getInstance("SHA-256").digest(body.getBytes())));
                }
                p(n + ".parseOK", true);
            } catch (Throwable t) {
                out.setLength(mark);
                p(n + ".parseOK", false);
                System.err.println(n + " rejected: " + t);
            }
        }
        System.out.print(out);
    }

    static String algId(AlgorithmIdentifier a) throws Exception {
        if (a == null) return "<nil>";
        return a.getAlgorithm().getId() + "/" +
            (a.getParameters() == null ? "<absent>"
                : Hex.toHexString(a.getParameters().toASN1Primitive().getEncoded(ASN1Encoding.DER)));
    }

    static void dumpCms(String n, byte[] data) throws Exception {
        ASN1Primitive top = ASN1Primitive.fromByteArray(data);
        ContentInfo ci = ContentInfo.getInstance(top);
        p(n + ".contentType", ci.getContentType().getId());
        // Definite length: BC hands back a DLSequence for a definite-length input and a
        // BERSequence for an indefinite-length one.
        p(n + ".definiteLength", !(top instanceof BERSequence));
        if (!ci.getContentType().equals(CMSObjectIdentifiers.signedData))
            throw new IllegalArgumentException("not id-signedData: " + ci.getContentType());
        SignedData sd = SignedData.getInstance(ci.getContent());
        p(n + ".version", sd.getVersion().getValue());
        for (int i = 0; i < sd.getDigestAlgorithms().size(); i++)
            p(n + ".digestAlgorithms[" + i + "]",
                algId(AlgorithmIdentifier.getInstance(sd.getDigestAlgorithms().getObjectAt(i))));
        ContentInfo eci = sd.getEncapContentInfo();
        p(n + ".eContentType", eci.getContentType().getId());
        p(n + ".detached", eci.getContent() == null);
        CMSSignedData cms = new CMSSignedData(data);
        Object sc = cms.getSignedContent();
        ph(n + ".signedContent", sc == null ? null : (byte[]) ((CMSProcessable) sc).getContent());

        // certificates / attribute certificates
        ASN1Set certs = sd.getCertificates();
        p(n + ".certificatesPresent", certs != null);
        int ci509 = 0, ciAttr = 0;
        if (certs != null) {
            for (int i = 0; i < certs.size(); i++) {
                ASN1Encodable o = certs.getObjectAt(i);
                if (o instanceof ASN1TaggedObject) {
                    ASN1TaggedObject t = ASN1TaggedObject.getInstance(o);
                    p(n + ".certChoiceTag[" + i + "]", t.getTagNo());
                    if (t.getTagNo() == 2)
                        ph(n + ".attrCert[" + (ciAttr++) + "]", t.getEncoded(ASN1Encoding.DER));
                } else {
                    p(n + ".certChoiceTag[" + i + "]", -1);
                    ph(n + ".cert[" + (ci509++) + "]", o.toASN1Primitive().getEncoded(ASN1Encoding.DER));
                }
            }
        }
        p(n + ".certCount", ci509);
        p(n + ".attrCertCount", ciAttr);

        // crls / OCSP
        ASN1Set crls = sd.getCRLs();
        p(n + ".crlsPresent", crls != null);
        int nc = 0, nOcsp = 0, nBasic = 0;
        if (crls != null) {
            for (int i = 0; i < crls.size(); i++) {
                ASN1Encodable o = crls.getObjectAt(i);
                if (o instanceof ASN1TaggedObject) {
                    ASN1Sequence seq = ASN1Sequence.getInstance((ASN1TaggedObject) o, false);
                    OtherRevocationInfoFormat other = OtherRevocationInfoFormat.getInstance(seq);
                    String fmt = other.getInfoFormat().getId();
                    p(n + ".otherRevFormat[" + i + "]", fmt);
                    byte[] info = other.getInfo().toASN1Primitive().getEncoded(ASN1Encoding.DER);
                    if (fmt.equals("1.3.6.1.5.5.7.16.2")) ph(n + ".ocspResponse[" + (nOcsp++) + "]", info);
                    else if (fmt.equals("1.3.6.1.5.5.7.48.1.1")) ph(n + ".ocspBasic[" + (nBasic++) + "]", info);
                } else {
                    ph(n + ".crl[" + (nc++) + "]", o.toASN1Primitive().getEncoded(ASN1Encoding.DER));
                }
            }
        }
        p(n + ".crlCount", nc);
        p(n + ".ocspResponseCount", nOcsp);
        p(n + ".ocspBasicCount", nBasic);

        // signer infos, in the order the SET carries them
        ASN1Set sis = sd.getSignerInfos();
        p(n + ".signerCount", sis.size());
        for (int i = 0; i < sis.size(); i++) {
            SignerInfo si = SignerInfo.getInstance(sis.getObjectAt(i));
            String k = n + ".signer[" + i + "]";
            p(k + ".version", si.getVersion().getValue());
            SignerIdentifier sid = si.getSID();
            p(k + ".sidIsSKI", sid.isTagged());
            if (sid.isTagged()) {
                ph(k + ".ski", ASN1OctetString.getInstance(sid.getId()).getOctets());
            } else {
                IssuerAndSerialNumber ias = IssuerAndSerialNumber.getInstance(sid.getId());
                ph(k + ".issuer", ias.getName().toASN1Primitive().getEncoded(ASN1Encoding.DER));
                p(k + ".serial", ias.getSerialNumber().getValue());
            }
            p(k + ".digestAlgorithm", algId(si.getDigestAlgorithm()));
            p(k + ".signatureAlgorithm", algId(si.getDigestEncryptionAlgorithm()));
            ph(k + ".signature", si.getEncryptedDigest().getOctets());
            ASN1Set sa = si.getAuthenticatedAttributes();
            p(k + ".hasSignedAttrs", sa != null);
            if (sa != null) {
                ph(k + ".signedAttrsDER", sa.getEncoded(ASN1Encoding.DER));
                p(k + ".signedAttrCount", sa.size());
                for (int j = 0; j < sa.size(); j++) {
                    Attribute at = Attribute.getInstance(sa.getObjectAt(j));
                    p(k + ".signedAttr[" + j + "].type", at.getAttrType().getId());
                    for (int v = 0; v < at.getAttrValues().size(); v++)
                        ph(k + ".signedAttr[" + j + "].value[" + v + "]",
                            at.getAttrValues().getObjectAt(v).toASN1Primitive().getEncoded(ASN1Encoding.DER));
                }
            }
            ASN1Set ua = si.getUnauthenticatedAttributes();
            p(k + ".hasUnsignedAttrs", ua != null);
            if (ua != null) ph(k + ".unsignedAttrsDER", ua.getEncoded(ASN1Encoding.DER));
        }
        ph(n + ".derEncoded", top.getEncoded(ASN1Encoding.DER));
        ph(n + ".dlEncoded", top.getEncoded(ASN1Encoding.DL));
        ph(n + ".berEncoded", top.getEncoded(ASN1Encoding.BER));
        // The five subtree encodings CMSObjectUtils writes for an ATSv3 message imprint.
        java.io.ByteArrayOutputStream bos = new java.io.ByteArrayOutputStream();
        sd.getDigestAlgorithms().encodeTo(bos);
        ph(n + ".ats.digestAlgorithms", bos.toByteArray());
        ph(n + ".ats.contentInfo", eci.getContent() instanceof BEROctetString
            ? eci.toASN1Primitive().getEncoded(ASN1Encoding.BER)
            : eci.toASN1Primitive().getEncoded(ASN1Encoding.DER));
        if (certs != null) {
            ph(n + ".ats.certificates", certs instanceof BERSet
                ? new BERTaggedObject(false, 0, new BERSequence(certs.toArray())).getEncoded()
                : new DERTaggedObject(false, 0, new DERSequence(certs.toArray())).getEncoded());
        }
        if (crls != null) {
            ph(n + ".ats.crls", crls instanceof BERSet
                ? new BERTaggedObject(false, 1, new BERSequence(crls.toArray())).getEncoded()
                : new DERTaggedObject(false, 1, new DERSequence(crls.toArray())).getEncoded());
        }
        ph(n + ".ats.signerInfos", sis instanceof BERSet
            ? new BERSet(sis.toArray()).getEncoded()
            : new DERSet(sis.toArray()).getEncoded());
        // RFC 5652 5.1 version recomputation, straight out of BC's SignedData constructor.
        SignedData rebuilt = new SignedData(sd.getDigestAlgorithms(), sd.getEncapContentInfo(),
            sd.getCertificates(), sd.getCRLs(), sd.getSignerInfos());
        p(n + ".recomputedVersion", rebuilt.getVersion().getValue());
    }

    static void dumpTst(String n, byte[] data) throws Exception {
        TimeStampToken tst = new TimeStampToken(ContentInfo.getInstance(ASN1Primitive.fromByteArray(data)));
        dumpTstInfo(n, tst);
        dumpCms(n + ".cms", data);
    }

    static void dumpTstInfo(String n, TimeStampToken tst) throws Exception {
        TimeStampTokenInfo info = tst.getTimeStampInfo();
        TSTInfo raw = info.toASN1Structure();
        p(n + ".tst.version", raw.getVersion().getValue());
        p(n + ".tst.policy", info.getPolicy().getId());
        p(n + ".tst.serial", info.getSerialNumber());
        p(n + ".tst.genTimeMillis", info.getGenTime().getTime());
        p(n + ".tst.genTimeString", raw.getGenTime().getTimeString());
        p(n + ".tst.imprintAlg", algId(info.getHashAlgorithm()));
        ph(n + ".tst.imprintDigest", info.getMessageImprintDigest());
        GenTimeAccuracy acc = info.getGenTimeAccuracy();
        p(n + ".tst.accuracyPresent", acc != null);
        if (acc != null) {
            Accuracy a = raw.getAccuracy();
            p(n + ".tst.accSeconds", a.getSeconds() == null ? "<nil>" : a.getSeconds().getValue());
            p(n + ".tst.accMillis", a.getMillis() == null ? "<nil>" : a.getMillis().getValue());
            p(n + ".tst.accMicros", a.getMicros() == null ? "<nil>" : a.getMicros().getValue());
        }
        p(n + ".tst.ordering", info.isOrdered());
        p(n + ".tst.nonce", info.getNonce() == null ? "<nil>" : info.getNonce());
        GeneralName tsa = raw.getTsa();
        p(n + ".tst.tsaPresent", tsa != null);
        if (tsa != null) {
            p(n + ".tst.tsaTag", tsa.getTagNo());
            ph(n + ".tst.tsaDER", tsa.toASN1Primitive().getEncoded(ASN1Encoding.DER));
        }
        p(n + ".tst.extensionsPresent", raw.getExtensions() != null);
        ph(n + ".tst.encoded", tst.getEncoded());
    }

    static void dumpTsr(String n, byte[] data) throws Exception {
        TimeStampResponse resp = new TimeStampResponse(data);
        p(n + ".status", resp.getStatus());
        p(n + ".statusString", resp.getStatusString() == null ? "<nil>" : "[" + resp.getStatusString() + "]");
        ASN1BitString fi = org.bouncycastle.asn1.cmp.PKIStatusInfo.getInstance(
                ASN1Sequence.getInstance(ASN1Primitive.fromByteArray(data)).getObjectAt(0)).getFailInfo();
        p(n + ".failInfoPresent", fi != null);
        if (fi != null) { ph(n + ".failInfoBits", fi.getBytes()); p(n + ".failInfoPadBits", fi.getPadBits()); }
        p(n + ".tokenPresent", resp.getTimeStampToken() != null);
        if (resp.getTimeStampToken() != null) {
            dumpTstInfo(n, resp.getTimeStampToken());
            dumpCms(n + ".cms", resp.getTimeStampToken().getEncoded());
        }
    }
}
