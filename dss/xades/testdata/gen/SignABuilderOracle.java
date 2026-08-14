// Generates testdata/sign-a-builder.txt: the exact bytes upstream DSS 6.5.RC1 produces from the
// XAdES signature-building core ported by the SIGN-A chunk - XAdESSignatureBuilder and its four
// packaging subclasses, ManifestBuilder, PrettyPrintTransformer and CounterSignatureBuilder.
//
// The Go port must reproduce these bytes exactly; that is the whole point of the chunk, so the
// expectations are upstream's own output rather than anything derived by hand. For every case it
// dumps, where applicable:
//
//   signedinfo-c14n  XAdESSignatureBuilder#build()  - the canonicalized ds:SignedInfo, i.e. the
//                    data-to-be-signed itself
//   signedinfo-raw   the pre-canonicalization serialization of the same ds:SignedInfo element
//   signedprops      the pre-canonicalization serialization of the xades:SignedProperties element
//   document         XAdESSignatureBuilder#signDocument(byte[]) with a fixed dummy signature
//                    value - the complete signature XML, which pins attribute order, namespace
//                    placement, prefix conventions and (absence of) indentation
//
// This program lives in package eu.europa.esig.dss.xades.signature because
// CounterSignatureBuilder's constructor and the four packaging builders are package-private.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies. Note that
// dss-xades cannot be built by maven in this environment (specs-trusted-list's JAXB generation
// needs schemas it cannot fetch), so its main sources are compiled directly, minus the tsl
// package - the only one that depends on specs-trusted-list:
//
//   cd /home/user/dss-upstream
//   mvn -o -pl dss-cades dependency:build-classpath -Dmdep.outputFile=/tmp/cades_cp.txt \
//       -Dmdep.includeScope=test
//   M2=$HOME/.m2/repository
//   CP="dss-xml-common/target/classes:dss-xml-utils/target/classes:\
//       $M2/eu/europa/ec/joinup/sd-dss/specs-xades/6.5.RC1/specs-xades-6.5.RC1.jar:\
//       $(cat /tmp/cades_cp.txt)"
//   find dss-xades/src/main/java -name '*.java' ! -path '*/xades/tsl/*' > /tmp/xades_srcs.txt
//   javac -nowarn -cp "$CP" -d /tmp/xadesbuild @/tmp/xades_srcs.txt
//   javac -nowarn -cp "$CP:/tmp/xadesbuild" -d /tmp/xadesbuild SignABuilderOracle.java
//   java -cp "$CP:/tmp/xadesbuild:dss-xades/src/main/resources" \
//        eu.europa.esig.dss.xades.signature.SignABuilderOracle <testdata directory>
//
// Every input is fixed - the checked-in testdata/signer.crt ("Plain Signer", taken verbatim from
// spi/validation/testdata/plain.crt), testdata/content-timestamp.tst (the RFC 3161 token of
// spi/validation/testdata/timestamp-token.tst), a frozen signing date and a constant dummy
// signature value - so a re-run reproduces the same bytes.
package eu.europa.esig.dss.xades.signature;

import eu.europa.esig.dss.enumerations.CommitmentTypeEnum;
import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.MimeTypeEnum;
import eu.europa.esig.dss.enumerations.ObjectIdentifierQualifier;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.enumerations.SignaturePackaging;
import eu.europa.esig.dss.enumerations.SigningOperation;
import eu.europa.esig.dss.model.CommitmentQualifier;
import eu.europa.esig.dss.model.CommonCommitmentType;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.Policy;
import eu.europa.esig.dss.model.SignerLocation;
import eu.europa.esig.dss.model.SpDocSpecification;
import eu.europa.esig.dss.model.UserNotice;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.validation.CertificateVerifier;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.tsp.TimestampInclude;
import eu.europa.esig.dss.spi.x509.tsp.TimestampToken;
import eu.europa.esig.dss.enumerations.TimestampType;
import eu.europa.esig.dss.xades.DSSObject;
import eu.europa.esig.dss.xades.XAdESSignatureParameters;
import eu.europa.esig.dss.xades.dataobject.DSSDataObjectFormat;
import eu.europa.esig.dss.xades.definition.XAdESNamespace;
import eu.europa.esig.dss.xades.reference.DSSReference;
import eu.europa.esig.dss.xades.reference.DSSTransform;
import eu.europa.esig.dss.xml.common.definition.DSSNamespace;
import eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigNamespace;
import eu.europa.esig.dss.xml.utils.DomUtils;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Calendar;
import java.util.Collections;
import java.util.Date;
import java.util.List;
import java.util.TimeZone;

public class SignABuilderOracle {

    /** The frozen signing date every case uses: 2021-01-01T00:00:00Z. */
    static final Date SIGNING_DATE = utc(2021, Calendar.JANUARY, 1, 0, 0, 0);

    /** The constant dummy ds:SignatureValue: 256 bytes 0x00..0xFF, an RSA-2048-sized block. */
    static final byte[] SIGNATURE_VALUE = new byte[256];

    static final String TEXT_CONTENT = "Hello World!";
    static final String XML_CONTENT =
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root xmlns=\"http://sample.com\" Id=\"root-id\">"
                    + "<child Id=\"child-id\">text</child></root>";

    static Path testdata;
    static CertificateToken signer;
    static StringBuilder out = new StringBuilder();

    public static void main(String[] args) throws Exception {
        testdata = Paths.get(args[0]);
        signer = loadCertificate(testdata.resolve("signer.crt"));
        for (int i = 0; i < SIGNATURE_VALUE.length; i++) {
            SIGNATURE_VALUE[i] = (byte) i;
        }

        out.append("# XAdES signature-building core, upstream DSS 6.5.RC1. Generated by gen/SignABuilderOracle.java.\n");
        out.append("# Records are 'case <name>' followed by '<key> <base64 of the bytes>' lines.\n");

        // ------------------------------------------------------------------ packaging variants
        emitBuild("enveloping-b", envelopingParams(), text());
        emitBuild("enveloped-b", envelopedParams(), xml());
        emitBuild("detached-b", detachedParams(), text());
        emitBuild("internally-detached-b", internallyDetachedParams(), xml());

        // multiple documents
        emitBuildAll("enveloping-two-documents", envelopingParams(),
                Arrays.asList(text(), new InMemoryDocument("second".getBytes(StandardCharsets.UTF_8),
                        "second.txt", MimeTypeEnum.TEXT)));
        emitBuildAll("detached-two-documents", detachedParams(),
                Arrays.asList(text(), new InMemoryDocument("second".getBytes(StandardCharsets.UTF_8),
                        "second.txt", MimeTypeEnum.TEXT)));

        // ------------------------------------------------------- EnvelopingSignatureBuilder paths
        XAdESSignatureParameters embedXml = envelopingParams();
        embedXml.setEmbedXML(true);
        emitBuild("enveloping-embed-xml", embedXml, xml());

        XAdESSignatureParameters manifest = envelopingParams();
        manifest.setManifestSignature(true);
        DSSDocument manifestDocument = new ManifestBuilder(DigestAlgorithm.SHA256,
                Collections.singletonList(text())).build();
        emitBuild("enveloping-manifest", manifest, manifestDocument);

        // ------------------------------------------------- XPathPlacementSignatureBuilder paths
        XAdESSignatureParameters xpathAfter = envelopedParams();
        xpathAfter.setXPathLocationString("//*[local-name()='child']");
        xpathAfter.setXPathElementPlacement(XAdESSignatureParameters.XPathElementPlacement.XPathAfter);
        emitBuild("enveloped-xpath-after", xpathAfter, xml());

        XAdESSignatureParameters xpathFirstChild = envelopedParams();
        xpathFirstChild.setXPathLocationString("//*[local-name()='root']");
        xpathFirstChild.setXPathElementPlacement(XAdESSignatureParameters.XPathElementPlacement.XPathFirstChildOf);
        emitBuild("enveloped-xpath-first-child", xpathFirstChild, xml());

        XAdESSignatureParameters internallyDetachedXPath = internallyDetachedParams();
        internallyDetachedXPath.setXPathLocationString("//*[local-name()='root']");
        emitBuild("internally-detached-xpath", internallyDetachedXPath, xml());

        // -------------------------------------------------------------------------- namespaces
        XAdESSignatureParameters noPrefix = envelopingParams();
        noPrefix.setXmldsigNamespace(new DSSNamespace(XMLDSigNamespace.NS.getUri(), ""));
        emitBuild("enveloping-default-xmldsig-prefix", noPrefix, text());

        XAdESSignatureParameters samePrefix = envelopingParams();
        samePrefix.setXmldsigNamespace(new DSSNamespace(XMLDSigNamespace.NS.getUri(), "ns"));
        samePrefix.setXadesNamespace(new DSSNamespace(XAdESNamespace.XADES_132.getUri(), "ns"));
        emitBuild("enveloping-same-prefix", samePrefix, text());

        XAdESSignatureParameters xades111 = envelopingParams();
        xades111.setXadesNamespace(new DSSNamespace(XAdESNamespace.XADES_111.getUri(), "xades111"));
        xades111.setEn319132(false);
        emitBuild("enveloping-xades111", xades111, text());

        XAdESSignatureParameters xades122 = envelopingParams();
        xades122.setXadesNamespace(new DSSNamespace(XAdESNamespace.XADES_122.getUri(), "xades122"));
        xades122.setEn319132(false);
        emitBuild("enveloping-xades122", xades122, text());

        // ------------------------------------------------------------------------- KeyInfo path
        XAdESSignatureParameters signKeyInfo = envelopingParams();
        signKeyInfo.setSignKeyInfo(true);
        emitBuild("enveloping-sign-key-info", signKeyInfo, text());

        XAdESSignatureParameters subjectName = envelopingParams();
        subjectName.setAddX509SubjectName(true);
        emitBuild("enveloping-x509-subject-name", subjectName, text());

        XAdESSignatureParameters noCert = new XAdESSignatureParameters();
        noCert.bLevel().setSigningDate(SIGNING_DATE);
        noCert.setSignatureLevel(SignatureLevel.XAdES_BASELINE_B);
        noCert.setSignaturePackaging(SignaturePackaging.ENVELOPING);
        noCert.setGenerateTBSWithoutCertificate(true);
        noCert.setDigestAlgorithm(DigestAlgorithm.SHA256);
        noCert.setEncryptionAlgorithm(eu.europa.esig.dss.enumerations.EncryptionAlgorithm.RSA);
        emitBuild("enveloping-no-signing-certificate", noCert, text());

        XAdESSignatureParameters sha1Cert = envelopingParams();
        sha1Cert.setSigningCertificateDigestMethod(DigestAlgorithm.SHA1);
        emitBuild("enveloping-signing-cert-sha1", sha1Cert, text());

        XAdESSignatureParameters v1 = envelopingParams();
        v1.setEn319132(false);
        emitBuild("enveloping-signing-certificate-v1", v1, text());

        // ------------------------------------------------------------------ signed properties
        XAdESSignatureParameters implicitPolicy = envelopingParams();
        Policy implied = new Policy();
        implicitPolicy.bLevel().setSignaturePolicy(implied);
        emitBuild("enveloping-policy-implied", implicitPolicy, text());

        XAdESSignatureParameters explicitPolicy = envelopingParams();
        explicitPolicy.bLevel().setSignaturePolicy(fullPolicy());
        emitBuild("enveloping-policy-explicit", explicitPolicy, text());

        XAdESSignatureParameters oidPolicy = envelopingParams();
        Policy urnPolicy = new Policy();
        urnPolicy.setId("1.2.3.4.5");
        urnPolicy.setQualifier(ObjectIdentifierQualifier.OID_AS_URN);
        urnPolicy.setDigestAlgorithm(DigestAlgorithm.SHA256);
        urnPolicy.setDigestValue(new byte[32]);
        oidPolicy.bLevel().setSignaturePolicy(urnPolicy);
        emitBuild("enveloping-policy-urn-oid", oidPolicy, text());

        XAdESSignatureParameters roles = envelopingParams();
        roles.bLevel().setClaimedSignerRoles(Arrays.asList("Manager", "Director"));
        roles.bLevel().setSignedAssertions(Collections.singletonList(
                "<Assertion xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\">assertion-body</Assertion>"));
        emitBuild("enveloping-signer-role-v2", roles, text());

        XAdESSignatureParameters rolesV1 = envelopingParams();
        rolesV1.setEn319132(false);
        rolesV1.bLevel().setClaimedSignerRoles(Arrays.asList("Manager", "Director"));
        emitBuild("enveloping-signer-role-v1", rolesV1, text());

        XAdESSignatureParameters place = envelopingParams();
        place.bLevel().setSignerLocation(signerLocation());
        emitBuild("enveloping-production-place-v2", place, text());

        XAdESSignatureParameters placeV1 = envelopingParams();
        placeV1.setEn319132(false);
        placeV1.bLevel().setSignerLocation(signerLocation());
        emitBuild("enveloping-production-place-v1", placeV1, text());

        XAdESSignatureParameters commitments = envelopingParams();
        commitments.bLevel().setCommitmentTypeIndications(
                Collections.singletonList(CommitmentTypeEnum.ProofOfOrigin));
        emitBuild("enveloping-commitment-enum", commitments, text());

        XAdESSignatureParameters qualifiedCommitments = envelopingParams();
        qualifiedCommitments.bLevel().setCommitmentTypeIndications(
                Collections.singletonList(qualifiedCommitmentType()));
        emitBuild("enveloping-commitment-qualified", qualifiedCommitments, text());

        XAdESSignatureParameters dataObjectFormats = envelopingParams();
        dataObjectFormats.setDataObjectFormatList(Collections.singletonList(customDataObjectFormat()));
        emitBuild("enveloping-data-object-format", dataObjectFormats, text());

        // ----------------------------------------------------------------------- ds:Object paths
        XAdESSignatureParameters objects = envelopingParams();
        objects.setObjects(Arrays.asList(xmlObject(), textObject()));
        emitBuild("enveloping-custom-objects", objects, text());

        // --------------------------------------------------------------------- content timestamp
        XAdESSignatureParameters contentTst = envelopingParams();
        contentTst.setContentTimestamps(Collections.singletonList(allDataObjectsTimestamp()));
        emitBuild("enveloping-content-timestamp", contentTst, text());

        XAdESSignatureParameters individualTst = envelopingParams();
        individualTst.setContentTimestamps(Collections.singletonList(individualDataObjectsTimestamp()));
        emitBuild("enveloping-individual-content-timestamp", individualTst, text());

        // ---------------------------------------------------------------------- pretty printing
        XAdESSignatureParameters pretty = envelopingParams();
        pretty.setPrettyPrint(true);
        pretty.getContext().setOperationKind(SigningOperation.SIGN);
        emitBuild("enveloping-pretty-print", pretty, text());

        XAdESSignatureParameters prettyEnveloped = envelopedParams();
        prettyEnveloped.setPrettyPrint(true);
        prettyEnveloped.getContext().setOperationKind(SigningOperation.SIGN);
        emitBuild("enveloped-pretty-print", prettyEnveloped, xml());

        // ----------------------------------------------------------------- explicit ds:SignedInfo
        XAdESSignatureParameters signedData = envelopingParams();
        XAdESSignatureParameters signedDataSource = envelopingParams();
        XAdESSignatureBuilder sourceBuilder = XAdESSignatureBuilder.getSignatureBuilder(
                signedDataSource, text(), verifier());
        sourceBuilder.build();
        byte[] signedInfoBytes = DomUtils.serializeNode(sourceBuilder.signedInfoDom);
        signedData.setSignedData(signedInfoBytes);
        signedData.setSignedAdESObject(DomUtils.serializeNode(
                sourceBuilder.qualifyingPropertiesDom.getParentNode()));
        emitBuild("enveloping-explicit-signed-data", signedData, text());

        // ---------------------------------------------------------------------- ManifestBuilder
        emitManifest("manifest-default", new ManifestBuilder(DigestAlgorithm.SHA256,
                Arrays.asList(text(), xml())));
        emitManifest("manifest-named-sha512", new ManifestBuilder("my-manifest", DigestAlgorithm.SHA512,
                Collections.singletonList(text())));
        emitManifest("manifest-custom-namespace", new ManifestBuilder("prefixed", DigestAlgorithm.SHA256,
                Collections.singletonList(text()), new DSSNamespace(XMLDSigNamespace.NS.getUri(), "dsig")));

        // ---------------------------------------------------------------- PrettyPrintTransformer
        emitPrettyPrint("pretty-keep-indents", true, 4,
                "<a><b>text</b><c/></a>");
        emitPrettyPrint("pretty-keep-indents-preindented", true, 4,
                "<a>\n  <b>text</b>\n  <c/>\n</a>");
        emitPrettyPrint("pretty-drop-indents-preindented", false, 4,
                "<a>\n  <b>text</b>\n  <c/>\n</a>");
        emitPrettyPrint("pretty-indent-2", true, 2,
                "<a><b><c>deep</c></b></a>");
        emitPrettyPrint("pretty-mixed-content", true, 4,
                "<a>lead<b>text</b>tail</a>");

        // --------------------------------------------------------------- CounterSignatureBuilder
        emitCounterSignature();

        Files.write(testdata.resolve("sign-a-builder.txt"), out.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("wrote " + testdata.resolve("sign-a-builder.txt"));
    }

    // ------------------------------------------------------------------------------- emitters

    static void emitBuild(String name, XAdESSignatureParameters params, DSSDocument document) {
        emitBuildAll(name, params, Collections.singletonList(document));
    }

    static void emitBuildAll(String name, XAdESSignatureParameters params, List<DSSDocument> documents) {
        XAdESSignatureBuilder builder =
                XAdESSignatureBuilder.getSignatureBuilder(params, documents, verifier());
        byte[] canonicalizedSignedInfo = builder.build();
        record(name, "signedinfo-c14n", canonicalizedSignedInfo);
        record(name, "signedinfo-raw", DomUtils.serializeNode(builder.signedInfoDom));
        if (builder.signedPropertiesDom != null) {
            record(name, "signedprops", DomUtils.serializeNode(builder.signedPropertiesDom));
        }
        DSSDocument signed = builder.signDocument(SIGNATURE_VALUE);
        record(name, "document", toByteArray(signed));
    }

    static void emitManifest(String name, ManifestBuilder builder) {
        record(name, "document", toByteArray(builder.build()));
    }

    static void emitPrettyPrint(String name, boolean keepOriginalIndents, int indentAmount, String xml) {
        org.w3c.dom.Document document = DomUtils.buildDOM(xml);
        org.w3c.dom.Node result = new PrettyPrintTransformer()
                .setKeepOriginalIndents(keepOriginalIndents)
                .setIndentAmount(indentAmount)
                .transform(document.getDocumentElement());
        record(name, "node", DomUtils.serializeNode(result));
    }

    static void emitCounterSignature() {
        // A signature to be counter-signed: built exactly as the enveloping-b case above.
        XAdESSignatureParameters masterParams = envelopingParams();
        XAdESSignatureBuilder masterBuilder =
                XAdESSignatureBuilder.getSignatureBuilder(masterParams, text(), verifier());
        masterBuilder.build();
        DSSDocument master = masterBuilder.signDocument(SIGNATURE_VALUE);
        record("counter-signature", "master-document", toByteArray(master));

        XAdESCounterSignatureParameters counterParams = new XAdESCounterSignatureParameters();
        counterParams.setSigningCertificate(signer);
        counterParams.setCertificateChain(Collections.singletonList(signer));
        counterParams.bLevel().setSigningDate(SIGNING_DATE);
        counterParams.setSignatureLevel(SignatureLevel.XAdES_BASELINE_B);
        counterParams.setSignaturePackaging(SignaturePackaging.ENVELOPED);
        counterParams.setDigestAlgorithm(DigestAlgorithm.SHA256);
        counterParams.setSignatureIdToCounterSign(masterParams.getDeterministicId());

        record("counter-signature", "deterministic-id",
                counterParams.getDeterministicId().getBytes(StandardCharsets.UTF_8));

        CounterSignatureBuilder builder = new CounterSignatureBuilder(verifier());
        DSSDocument canonicalizedSignatureValue =
                builder.getCanonicalizedSignatureValue(master, counterParams);
        record("counter-signature", "canonicalized-signature-value",
                toByteArray(canonicalizedSignatureValue));

        // Same builder instance, as XAdESService#counterSignSignature does: only
        // getCanonicalizedSignatureValue assigns params, and buildCounterSignatureDSSReference
        // relies on that having happened.
        DSSReference reference = builder.buildCounterSignatureDSSReference(master, counterParams);
        record("counter-signature", "reference-id", reference.getId().getBytes(StandardCharsets.UTF_8));
        record("counter-signature", "reference-uri", reference.getUri().getBytes(StandardCharsets.UTF_8));
        record("counter-signature", "reference-type", reference.getType().getBytes(StandardCharsets.UTF_8));
        record("counter-signature", "reference-digest-method",
                reference.getDigestMethodAlgorithm().name().getBytes(StandardCharsets.UTF_8));
        record("counter-signature", "reference-contents", toByteArray(reference.getContents()));
        DSSTransform transform = reference.getTransforms().get(0);
        record("counter-signature", "reference-transform",
                transform.getAlgorithm().getBytes(StandardCharsets.UTF_8));

        // The counter signature itself: an enveloped signature over the canonicalized
        // SignatureValue, built with the counter-signature reference above.
        counterParams.setReferences(Collections.singletonList(reference));
        XAdESSignatureParameters embedded = new XAdESSignatureParameters();
        embedded.setSigningCertificate(signer);
        embedded.setCertificateChain(Collections.singletonList(signer));
        embedded.bLevel().setSigningDate(SIGNING_DATE);
        embedded.setSignatureLevel(SignatureLevel.XAdES_BASELINE_B);
        embedded.setSignaturePackaging(SignaturePackaging.ENVELOPING);
        embedded.setDigestAlgorithm(DigestAlgorithm.SHA256);
        XAdESSignatureBuilder counterBuilder = XAdESSignatureBuilder.getSignatureBuilder(
                embedded, canonicalizedSignatureValue, verifier());
        counterBuilder.build();
        DSSDocument counterSignature = counterBuilder.signDocument(SIGNATURE_VALUE);
        record("counter-signature", "counter-signature-document", toByteArray(counterSignature));

        CounterSignatureBuilder embedBuilder = new CounterSignatureBuilder(verifier());
        DSSDocument embeddedResult =
                embedBuilder.buildEmbeddedCounterSignature(master, counterSignature, counterParams);
        record("counter-signature", "embedded-document", toByteArray(embeddedResult));
    }

    static void record(String name, String key, byte[] value) {
        out.append("case ").append(name).append('\n');
        out.append(key).append(' ').append(Base64.getEncoder().encodeToString(value)).append('\n');
    }

    // -------------------------------------------------------------------------------- fixtures

    static XAdESSignatureParameters baseParams() {
        XAdESSignatureParameters params = new XAdESSignatureParameters();
        params.setSigningCertificate(signer);
        params.setCertificateChain(Collections.singletonList(signer));
        params.bLevel().setSigningDate(SIGNING_DATE);
        params.setSignatureLevel(SignatureLevel.XAdES_BASELINE_B);
        params.setDigestAlgorithm(DigestAlgorithm.SHA256);
        return params;
    }

    static XAdESSignatureParameters envelopingParams() {
        XAdESSignatureParameters params = baseParams();
        params.setSignaturePackaging(SignaturePackaging.ENVELOPING);
        return params;
    }

    static XAdESSignatureParameters envelopedParams() {
        XAdESSignatureParameters params = baseParams();
        params.setSignaturePackaging(SignaturePackaging.ENVELOPED);
        return params;
    }

    static XAdESSignatureParameters detachedParams() {
        XAdESSignatureParameters params = baseParams();
        params.setSignaturePackaging(SignaturePackaging.DETACHED);
        return params;
    }

    static XAdESSignatureParameters internallyDetachedParams() {
        XAdESSignatureParameters params = baseParams();
        params.setSignaturePackaging(SignaturePackaging.INTERNALLY_DETACHED);
        return params;
    }

    static DSSDocument text() {
        return new InMemoryDocument(TEXT_CONTENT.getBytes(StandardCharsets.UTF_8), "hello.txt", MimeTypeEnum.TEXT);
    }

    static DSSDocument xml() {
        return new InMemoryDocument(XML_CONTENT.getBytes(StandardCharsets.UTF_8), "sample.xml", MimeTypeEnum.XML);
    }

    static Policy fullPolicy() {
        Policy policy = new Policy();
        policy.setId("http://spuri.test/policy");
        policy.setDescription("Test policy description");
        policy.setDocumentationReferences("http://doc.ref/1", "http://doc.ref/2");
        policy.setDigestAlgorithm(DigestAlgorithm.SHA256);
        policy.setDigestValue(new byte[]{1, 2, 3, 4});
        policy.setSpuri("http://spuri.test");
        UserNotice userNotice = new UserNotice();
        userNotice.setOrganization("ACME Ltd");
        userNotice.setNoticeNumbers(new int[]{1, 2});
        userNotice.setExplicitText("This is a test policy");
        policy.setUserNotice(userNotice);
        SpDocSpecification spDocSpecification = new SpDocSpecification();
        spDocSpecification.setId("1.2.3.4.6");
        spDocSpecification.setDescription("Spec description");
        spDocSpecification.setDocumentationReferences("http://spec.ref/1");
        policy.setSpDocSpecification(spDocSpecification);
        return policy;
    }

    static SignerLocation signerLocation() {
        SignerLocation signerLocation = new SignerLocation();
        signerLocation.setLocality("Luxembourg");
        signerLocation.setStreetAddress("1 Main Street");
        signerLocation.setStateOrProvince("Luxembourg");
        signerLocation.setPostalCode("L-1234");
        signerLocation.setCountry("LU");
        return signerLocation;
    }

    static CommonCommitmentType qualifiedCommitmentType() {
        CommonCommitmentType commitmentType = new CommonCommitmentType();
        commitmentType.setOid("1.2.840.113549.1.9.16.6.1");
        commitmentType.setQualifier(ObjectIdentifierQualifier.OID_AS_URN);
        commitmentType.setDescription("Commitment description");
        commitmentType.setDocumentationReferences("http://commitment.ref/1");
        commitmentType.setSignedDataObjects("r-id-1", "#r-id-2");
        CommitmentQualifier xmlQualifier = new CommitmentQualifier();
        xmlQualifier.setContent(new InMemoryDocument(
                "<qualifier xmlns=\"http://sample.com\">value</qualifier>".getBytes(StandardCharsets.UTF_8)));
        CommitmentQualifier textQualifier = new CommitmentQualifier();
        textQualifier.setContent(new InMemoryDocument("plain qualifier".getBytes(StandardCharsets.UTF_8)));
        commitmentType.setCommitmentTypeQualifiers(xmlQualifier, textQualifier);
        return commitmentType;
    }

    static DSSDataObjectFormat customDataObjectFormat() {
        DSSDataObjectFormat dataObjectFormat = new DSSDataObjectFormat();
        dataObjectFormat.setDescription("Custom description");
        dataObjectFormat.setMimeType(MimeTypeEnum.TEXT.getMimeTypeString());
        dataObjectFormat.setEncoding("http://www.w3.org/2000/09/xmldsig#base64");
        dataObjectFormat.setObjectReference("#r-id-1");
        return dataObjectFormat;
    }

    static DSSObject xmlObject() {
        DSSObject object = new DSSObject();
        object.setContent(new InMemoryDocument(
                "<data xmlns=\"http://sample.com\">object-content</data>".getBytes(StandardCharsets.UTF_8)));
        object.setId("custom-xml-object");
        object.setMimeType(MimeTypeEnum.XML.getMimeTypeString());
        return object;
    }

    static DSSObject textObject() {
        DSSObject object = new DSSObject();
        object.setContent(new InMemoryDocument("plain object".getBytes(StandardCharsets.UTF_8)));
        object.setId("custom-text-object");
        object.setMimeType(MimeTypeEnum.TEXT.getMimeTypeString());
        object.setEncodingAlgorithm("http://www.w3.org/2000/09/xmldsig#base64");
        return object;
    }

    static TimestampToken allDataObjectsTimestamp() throws Exception {
        byte[] der = Files.readAllBytes(testdata.resolve("content-timestamp.tst"));
        TimestampToken token = new TimestampToken(der, TimestampType.ALL_DATA_OBJECTS_TIMESTAMP);
        token.setCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#");
        return token;
    }

    static TimestampToken individualDataObjectsTimestamp() throws Exception {
        byte[] der = Files.readAllBytes(testdata.resolve("content-timestamp.tst"));
        TimestampToken token = new TimestampToken(der, TimestampType.INDIVIDUAL_DATA_OBJECTS_TIMESTAMP);
        token.setCanonicalizationMethod("http://www.w3.org/2001/10/xml-exc-c14n#");
        List<TimestampInclude> includes = new ArrayList<>();
        includes.add(new TimestampInclude("r-id-1", true));
        token.setTimestampIncludes(includes);
        return token;
    }

    static CertificateVerifier verifier() {
        return new CommonCertificateVerifier();
    }

    // ---------------------------------------------------------------------------------- helpers

    static byte[] toByteArray(DSSDocument document) {
        try (java.io.InputStream is = document.openStream();
             java.io.ByteArrayOutputStream baos = new java.io.ByteArrayOutputStream()) {
            byte[] buffer = new byte[8192];
            int read;
            while ((read = is.read(buffer)) != -1) {
                baos.write(buffer, 0, read);
            }
            return baos.toByteArray();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    static CertificateToken loadCertificate(Path path) throws Exception {
        CertificateFactory factory = CertificateFactory.getInstance("X.509");
        try (java.io.InputStream is = Files.newInputStream(path)) {
            return new CertificateToken((X509Certificate) factory.generateCertificate(is));
        }
    }

    static Date utc(int year, int month, int day, int hour, int minute, int second) {
        Calendar calendar = Calendar.getInstance(TimeZone.getTimeZone("UTC"));
        calendar.clear();
        calendar.set(year, month, day, hour, minute, second);
        return calendar.getTime();
    }
}
