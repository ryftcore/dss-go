// Generates testdata/refs.txt: the exact bytes upstream DSS 6.5.RC1 produces from the XAdES
// reference-and-transform core ported by the REFS chunk - the DSSTransform hierarchy,
// DSSTransformOutput, DSSReference, ReferenceIdProvider, ReferenceBuilder, ReferenceProcessor
// and ReferenceVerifier.
//
// The Go port must reproduce these bytes exactly, so every expectation is upstream's own output
// rather than anything derived by hand. The keys dumped are:
//
//   transforms       DSSXMLUtils.incorporateTransforms(...) - the ds:Transforms subtree a
//                    transform WRITES into a ds:Reference, which pins element namespaces,
//                    prefixes, attribute order and the namespace declarations the Filter 2.0
//                    transform puts on its XPath element
//   algorithm        DSSTransform#getAlgorithm
//   output           DSSXMLUtils.applyTransforms(node, transforms) - the octets a transform
//                    chain EXECUTES to, i.e. what the reference digest is taken over. This is
//                    the Santuario path the Go port routes through internal/xmldsig.
//   error            the message of the exception a rejected setup throws ("" when none)
//   ref-*            the fields ReferenceBuilder#build() fills in, per reference
//   references       ReferenceProcessor#incorporateReferences(...) - the complete ds:Reference
//                    list, digests included
//   signedinfo-raw   XAdESSignatureBuilder's ds:SignedInfo element before canonicalization
//   signedinfo-c14n  XAdESSignatureBuilder#build() - the canonicalized ds:SignedInfo, i.e. the
//                    data-to-be-signed
//
// This program lives in package eu.europa.esig.dss.xades.signature because
// XAdESSignatureBuilder#getSignatureBuilder and ManifestBuilder are package-private; everything
// it exercises from eu.europa.esig.dss.xades.reference is public.
//
// Run it with OpenJDK 21 against the built upstream DSS 6.5.RC1 and its dependencies. Note that
// dss-xades cannot be built by maven in this environment (specs-trusted-list's JAXB generation
// needs schemas it cannot fetch), so its main sources are compiled directly, minus the tsl
// package - the only one that depends on specs-trusted-list:
//
//   cd $DSS_UPSTREAM_HOME  (your built upstream DSS 6.5.RC1 checkout)
//   mvn -o -pl dss-cades dependency:build-classpath -Dmdep.outputFile=/tmp/cades_cp.txt \
//       -Dmdep.includeScope=test
//   M2=$HOME/.m2/repository
//   CP="dss-xml-common/target/classes:dss-xml-utils/target/classes:\
//       $M2/eu/europa/ec/joinup/sd-dss/specs-xades/6.5.RC1/specs-xades-6.5.RC1.jar:\
//       $(cat /tmp/cades_cp.txt)"
//   find dss-xades/src/main/java -name '*.java' ! -path '*/xades/tsl/*' > /tmp/xades_srcs.txt
//   javac -nowarn -cp "$CP" -d /tmp/xadesbuild @/tmp/xades_srcs.txt
//   javac -nowarn -cp "$CP:/tmp/xadesbuild" -d /tmp/xadesbuild RefsOracle.java
//   java -cp "$CP:/tmp/xadesbuild:dss-xades/src/main/resources" \
//        eu.europa.esig.dss.xades.signature.RefsOracle <testdata directory>
//
// Every input is fixed - the checked-in testdata/signer.crt, a frozen signing date and a
// constant dummy signature value - so a re-run reproduces the same bytes.
package eu.europa.esig.dss.xades.signature;

import eu.europa.esig.dss.enumerations.DigestAlgorithm;
import eu.europa.esig.dss.enumerations.MimeTypeEnum;
import eu.europa.esig.dss.enumerations.SignatureLevel;
import eu.europa.esig.dss.enumerations.SignaturePackaging;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.InMemoryDocument;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.validation.CertificateVerifier;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.xades.DSSXMLUtils;
import eu.europa.esig.dss.xades.XAdESSignatureParameters;
import eu.europa.esig.dss.xades.definition.XAdESNamespace;
import eu.europa.esig.dss.xades.reference.Base64Transform;
import eu.europa.esig.dss.xades.reference.CanonicalizationTransform;
import eu.europa.esig.dss.xades.reference.DSSReference;
import eu.europa.esig.dss.xades.reference.DSSTransform;
import eu.europa.esig.dss.xades.reference.EnvelopedSignatureTransform;
import eu.europa.esig.dss.xades.reference.ReferenceBuilder;
import eu.europa.esig.dss.xades.reference.ReferenceIdProvider;
import eu.europa.esig.dss.xades.reference.ReferenceProcessor;
import eu.europa.esig.dss.xades.reference.ReferenceVerifier;
import eu.europa.esig.dss.xades.reference.SPDocDigestAsInSpecificationTransform;
import eu.europa.esig.dss.xades.reference.XPath2FilterEnvelopedSignatureTransform;
import eu.europa.esig.dss.xades.reference.XPath2FilterTransform;
import eu.europa.esig.dss.xades.reference.XPathEnvelopedSignatureTransform;
import eu.europa.esig.dss.xades.reference.XPathTransform;
import eu.europa.esig.dss.xades.reference.XsltTransform;
import eu.europa.esig.dss.xml.common.definition.DSSNamespace;
import eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigElement;
import eu.europa.esig.dss.xml.common.definition.xmldsig.XMLDSigNamespace;
import eu.europa.esig.dss.xml.utils.DomUtils;
import eu.europa.esig.dss.xml.utils.XMLCanonicalizer;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Calendar;
import java.util.Collections;
import java.util.Date;
import java.util.GregorianCalendar;
import java.util.List;
import java.util.TimeZone;

public class RefsOracle {

    /** The frozen signing date every case uses: 2021-01-01T00:00:00Z. */
    static final Date SIGNING_DATE = utc(2021, Calendar.JANUARY, 1, 0, 0, 0);

    /** The constant dummy ds:SignatureValue: 256 bytes 0x00..0xFF, an RSA-2048-sized block. */
    static final byte[] SIGNATURE_VALUE = new byte[256];

    static final String TEXT_CONTENT = "Hello World!";

    static final String XML_CONTENT =
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root xmlns=\"http://sample.com\" Id=\"root-id\">"
                    + "<child Id=\"child-id\">text</child></root>";

    /** An XML document that already carries a ds:Signature, so the enveloped filters have work. */
    static final String XML_SIGNED_CONTENT =
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?><root xmlns=\"http://sample.com\" Id=\"root-id\">"
                    + "<child Id=\"child-id\">text</child><!-- a comment -->"
                    + "<ds:Signature xmlns:ds=\"http://www.w3.org/2000/09/xmldsig#\" Id=\"sig-1\">"
                    + "<ds:SignedInfo><ds:Reference URI=\"\"/></ds:SignedInfo></ds:Signature></root>";

    static final String XSLT_CONTENT =
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
                    + "<xsl:stylesheet xmlns:xsl=\"http://www.w3.org/1999/XSL/Transform\" version=\"1.0\">"
                    + "<xsl:template match=\"/\"><out/></xsl:template></xsl:stylesheet>";

    /** Santuario's "physical" canonicalization: a canonicalizer, but never a TransformSpi. */
    static final String PHYSICAL_C14N = "http://santuario.apache.org/c14n/physical";

    static final DSSNamespace DS_NS = XMLDSigNamespace.NS;
    static final DSSNamespace DEFAULT_NS = new DSSNamespace(XMLDSigNamespace.NS.getUri(), "");
    static final DSSNamespace DSIG_NS = new DSSNamespace(XMLDSigNamespace.NS.getUri(), "dsig");

    static Path testdata;
    static CertificateToken signer;
    static StringBuilder out = new StringBuilder();

    public static void main(String[] args) throws Exception {
        testdata = Paths.get(args[0]);
        signer = DSSUtils.loadCertificate(Files.readAllBytes(testdata.resolve("signer.crt")));
        for (int i = 0; i < SIGNATURE_VALUE.length; i++) {
            SIGNATURE_VALUE[i] = (byte) i;
        }

        out.append("# XAdES reference/transform core, upstream DSS 6.5.RC1. Generated by gen/RefsOracle.java.\n");
        out.append("# Records are 'case <name>' followed by '<key> <base64 of the bytes>' lines.\n");

        emitTransformDoms();
        emitTransformOutputs();
        emitReferenceIdProvider();
        emitReferenceBuilder();
        emitIncorporateReferences();
        emitReferenceVerifier();
        emitSignedInfo();

        Files.write(testdata.resolve("refs.txt"), out.toString().getBytes(StandardCharsets.UTF_8));
        System.out.println("wrote " + testdata.resolve("refs.txt"));
    }

    // ------------------------------------------------------------------ 1. the ds:Transforms DOM

    static void emitTransformDoms() {
        emitTransformDom("tf-base64", DS_NS, new Base64Transform());
        emitTransformDom("tf-base64-dsig", DSIG_NS, new Base64Transform(DSIG_NS));
        emitTransformDom("tf-base64-default-prefix", DEFAULT_NS, new Base64Transform(DEFAULT_NS));

        emitTransformDom("tf-c14n-exclusive", DS_NS,
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD));
        emitTransformDom("tf-c14n-inclusive", DS_NS,
                new CanonicalizationTransform("http://www.w3.org/TR/2001/REC-xml-c14n-20010315"));
        emitTransformDom("tf-c14n-11-comments", DS_NS,
                new CanonicalizationTransform("http://www.w3.org/2006/12/xml-c14n11#WithComments"));
        emitTransformDom("tf-c14n-exclusive-dsig", DSIG_NS,
                new CanonicalizationTransform(DSIG_NS, XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD));

        emitTransformDom("tf-enveloped", DS_NS, new EnvelopedSignatureTransform());
        emitTransformDom("tf-enveloped-default-prefix", DEFAULT_NS,
                new EnvelopedSignatureTransform(DEFAULT_NS));

        emitTransformDom("tf-xpath", DS_NS, new XPathTransform("//*[local-name()='child']"));
        emitTransformDom("tf-xpath-dsig", DSIG_NS, new XPathTransform(DSIG_NS, "//*[local-name()='child']"));
        emitTransformDom("tf-xpath-enveloped", DS_NS, new XPathEnvelopedSignatureTransform());
        emitTransformDom("tf-xpath-enveloped-dsig", DSIG_NS, new XPathEnvelopedSignatureTransform(DSIG_NS));

        emitTransformDom("tf-xpath2filter", DS_NS,
                new XPath2FilterTransform("/root/child", "intersect"));
        emitTransformDom("tf-xpath2filter-enveloped", DS_NS, new XPath2FilterEnvelopedSignatureTransform());
        emitTransformDom("tf-xpath2filter-enveloped-dsig", DSIG_NS,
                new XPath2FilterEnvelopedSignatureTransform(DSIG_NS));
        // The one case that exercises the "xmlns" (no prefix) branch of XPath2FilterTransform.
        emitTransformDom("tf-xpath2filter-enveloped-default-prefix", DEFAULT_NS,
                new XPath2FilterEnvelopedSignatureTransform(DEFAULT_NS));

        emitTransformDom("tf-xslt", DS_NS, new XsltTransform(DomUtils.buildDOM(XSLT_CONTENT)));
        emitTransformDom("tf-spdoc-digest", DS_NS, new SPDocDigestAsInSpecificationTransform());

        emitTransformDoms("tf-chain-enveloped-c14n", DS_NS,
                Arrays.asList(new XPath2FilterEnvelopedSignatureTransform(),
                        new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
    }

    static void emitTransformDom(String name, DSSNamespace namespace, DSSTransform transform) {
        emitTransformDoms(name, namespace, Collections.singletonList(transform));
    }

    static void emitTransformDoms(String name, DSSNamespace namespace, List<DSSTransform> transforms) {
        Document document = DomUtils.buildDOM();
        Element referenceDom = DomUtils.createElementNS(document, namespace, XMLDSigElement.REFERENCE);
        document.appendChild(referenceDom);
        DSSXMLUtils.incorporateTransforms(referenceDom, transforms, namespace);
        record(name, "transforms", DomUtils.serializeNode(referenceDom));
        StringBuilder algorithms = new StringBuilder();
        for (DSSTransform transform : transforms) {
            if (algorithms.length() > 0) {
                algorithms.append('\n');
            }
            algorithms.append(transform.getAlgorithm());
        }
        record(name, "algorithm", algorithms.toString().getBytes(StandardCharsets.UTF_8));
    }

    // -------------------------------------------------------------- 2. the executed transforms

    static void emitTransformOutputs() {
        emitApply("apply-no-transforms", XML_SIGNED_CONTENT, Collections.emptyList());

        emitApply("apply-c14n-exclusive", XML_SIGNED_CONTENT,
                Collections.singletonList(new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
        emitApply("apply-c14n-inclusive", XML_SIGNED_CONTENT, Collections.singletonList(
                new CanonicalizationTransform("http://www.w3.org/TR/2001/REC-xml-c14n-20010315")));
        emitApply("apply-c14n-inclusive-comments", XML_SIGNED_CONTENT, Collections.singletonList(
                new CanonicalizationTransform("http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments")));
        emitApply("apply-c14n-11", XML_SIGNED_CONTENT, Collections.singletonList(
                new CanonicalizationTransform("http://www.w3.org/2006/12/xml-c14n11")));
        emitApply("apply-c14n-exclusive-comments", XML_SIGNED_CONTENT, Collections.singletonList(
                new CanonicalizationTransform("http://www.w3.org/2001/10/xml-exc-c14n#WithComments")));

        // The enveloped-signature transform is a DSS no-op on the creation path.
        emitApply("apply-enveloped-then-c14n", XML_SIGNED_CONTENT, Arrays.asList(
                new EnvelopedSignatureTransform(),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));

        // The two filters that really run: Santuario subtracts / excludes the ds:Signature.
        emitApply("apply-xpath2filter-enveloped-then-c14n", XML_SIGNED_CONTENT, Arrays.asList(
                new XPath2FilterEnvelopedSignatureTransform(),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
        emitApply("apply-xpath2filter-enveloped", XML_SIGNED_CONTENT,
                Collections.singletonList(new XPath2FilterEnvelopedSignatureTransform()));
        emitApply("apply-xpath-enveloped-then-c14n", XML_SIGNED_CONTENT, Arrays.asList(
                new XPathEnvelopedSignatureTransform(),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
        emitApply("apply-xpath-enveloped", XML_SIGNED_CONTENT,
                Collections.singletonList(new XPathEnvelopedSignatureTransform()));
        emitApply("apply-xpath-select-child", XML_SIGNED_CONTENT, Arrays.asList(
                new XPathTransform("ancestor-or-self::*[local-name()='child']"),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
        emitApply("apply-xpath2filter-intersect", XML_SIGNED_CONTENT, Arrays.asList(
                new XPath2FilterTransform("//*[local-name()='child']", "intersect"),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));

        // Base64Transform#performTransform is the identity; the bytes come out of the node.
        emitApply("apply-base64", XML_SIGNED_CONTENT, Collections.singletonList(new Base64Transform()));

        // Rejections.
        emitApply("apply-spdoc-digest", XML_SIGNED_CONTENT,
                Collections.singletonList(new SPDocDigestAsInSpecificationTransform()));
        emitApply("apply-xslt", XML_SIGNED_CONTENT,
                Collections.singletonList(new XsltTransform(DomUtils.buildDOM(XSLT_CONTENT))));
        // "physical" canonicalizes but is not a registered transform algorithm.
        emitApply("apply-c14n-physical", XML_SIGNED_CONTENT,
                Collections.singletonList(new CanonicalizationTransform(PHYSICAL_C14N)));
    }

    static void emitApply(String name, String xml, List<DSSTransform> transforms) {
        Document document = DomUtils.buildDOM(xml);
        try {
            record(name, "output", DSSXMLUtils.applyTransforms(document, transforms));
            record(name, "error", new byte[0]);
        } catch (RuntimeException e) {
            record(name, "error", message(e).getBytes(StandardCharsets.UTF_8));
        }
    }

    // ---------------------------------------------------------------- 3. the reference id provider

    static void emitReferenceIdProvider() {
        ReferenceIdProvider plain = new ReferenceIdProvider();
        record("refid-no-parameters", "ids", ids(plain, 3));

        ReferenceIdProvider withParams = new ReferenceIdProvider();
        XAdESSignatureParameters params = envelopingParams();
        withParams.setSignatureParameters(params);
        record("refid-with-parameters", "deterministic-id",
                params.getDeterministicId().getBytes(StandardCharsets.UTF_8));
        record("refid-with-parameters", "ids", ids(withParams, 3));

        ReferenceIdProvider prefixed = new ReferenceIdProvider();
        prefixed.setReferenceIdPrefix("m-id");
        record("refid-custom-prefix", "ids", ids(prefixed, 2));

        try {
            new ReferenceIdProvider().setReferenceIdPrefix("  ");
            record("refid-blank-prefix", "error", new byte[0]);
        } catch (RuntimeException e) {
            record("refid-blank-prefix", "error", message(e).getBytes(StandardCharsets.UTF_8));
        }
    }

    static byte[] ids(ReferenceIdProvider provider, int count) {
        StringBuilder builder = new StringBuilder();
        for (int i = 0; i < count; i++) {
            if (i > 0) {
                builder.append('\n');
            }
            builder.append(provider.getReferenceId());
        }
        return builder.toString().getBytes(StandardCharsets.UTF_8);
    }

    // --------------------------------------------------------------------- 4. ReferenceBuilder

    static void emitReferenceBuilder() {
        emitBuild("refbuild-enveloped", envelopedParams(), Collections.singletonList(xml()));
        emitBuild("refbuild-enveloping", envelopingParams(), Collections.singletonList(text()));
        emitBuild("refbuild-detached", detachedParams(), Collections.singletonList(text()));
        emitBuild("refbuild-detached-unnamed", detachedParams(), Collections.singletonList(
                new InMemoryDocument(TEXT_CONTENT.getBytes(StandardCharsets.UTF_8))));
        emitBuild("refbuild-internally-detached", internallyDetachedParams(), Collections.singletonList(xml()));

        XAdESSignatureParameters embedXml = envelopingParams();
        embedXml.setEmbedXML(true);
        emitBuild("refbuild-enveloping-embed-xml", embedXml, Collections.singletonList(xml()));

        XAdESSignatureParameters manifest = envelopingParams();
        manifest.setManifestSignature(true);
        emitBuild("refbuild-enveloping-manifest", manifest, Collections.singletonList(manifestDocument()));

        emitBuild("refbuild-two-documents", detachedParams(), Arrays.asList(text(), second()));

        // The detached constructor, i.e. the ManifestBuilder path: no signature parameters.
        ReferenceIdProvider provider = new ReferenceIdProvider();
        ReferenceBuilder builder = new ReferenceBuilder(Arrays.asList(text(), xml()),
                DigestAlgorithm.SHA256, provider);
        recordReferences("refbuild-no-parameters", builder.build());

        // Rejections.
        emitBuildError("refbuild-enveloped-not-xml", envelopedParams(), Collections.singletonList(text()));
        XAdESSignatureParameters manifestNoId = envelopingParams();
        manifestNoId.setManifestSignature(true);
        emitBuildError("refbuild-manifest-without-id", manifestNoId, Collections.singletonList(xmlWithoutId()));
        XAdESSignatureParameters embedNotXml = envelopingParams();
        embedNotXml.setEmbedXML(true);
        emitBuildError("refbuild-embed-xml-not-xml", embedNotXml, Collections.singletonList(text()));
    }

    static void emitBuild(String name, XAdESSignatureParameters params, List<DSSDocument> documents) {
        ReferenceIdProvider provider = new ReferenceIdProvider();
        provider.setSignatureParameters(params);
        recordReferences(name, new ReferenceBuilder(documents, params, provider).build());
    }

    static void emitBuildError(String name, XAdESSignatureParameters params, List<DSSDocument> documents) {
        ReferenceIdProvider provider = new ReferenceIdProvider();
        provider.setSignatureParameters(params);
        try {
            new ReferenceBuilder(documents, params, provider).build();
            record(name, "error", new byte[0]);
        } catch (RuntimeException e) {
            record(name, "error", message(e).getBytes(StandardCharsets.UTF_8));
        }
    }

    static void recordReferences(String name, List<DSSReference> references) {
        record(name, "count", String.valueOf(references.size()).getBytes(StandardCharsets.UTF_8));
        for (int i = 0; i < references.size(); i++) {
            DSSReference reference = references.get(i);
            record(name, "ref" + i + "-id", nullable(reference.getId()));
            record(name, "ref" + i + "-uri-present",
                    String.valueOf(reference.getUri() != null).getBytes(StandardCharsets.UTF_8));
            record(name, "ref" + i + "-uri", nullable(reference.getUri()));
            record(name, "ref" + i + "-type", nullable(reference.getType()));
            record(name, "ref" + i + "-digest-method",
                    reference.getDigestMethodAlgorithm().name().getBytes(StandardCharsets.UTF_8));
            StringBuilder algorithms = new StringBuilder();
            if (reference.getTransforms() != null) {
                for (DSSTransform transform : reference.getTransforms()) {
                    if (algorithms.length() > 0) {
                        algorithms.append('\n');
                    }
                    algorithms.append(transform.getAlgorithm());
                }
            }
            record(name, "ref" + i + "-transforms", algorithms.toString().getBytes(StandardCharsets.UTF_8));
        }
    }

    // ------------------------------------------------------------------- 5. ReferenceProcessor

    static void emitIncorporateReferences() {
        emitIncorporate("incorporate-enveloped", envelopedParams(), Collections.singletonList(xml()), DS_NS);
        emitIncorporate("incorporate-enveloping", envelopingParams(), Collections.singletonList(text()), DS_NS);
        emitIncorporate("incorporate-detached", detachedParams(), Collections.singletonList(text()), DS_NS);
        emitIncorporate("incorporate-detached-unnamed", detachedParams(), Collections.singletonList(
                new InMemoryDocument(TEXT_CONTENT.getBytes(StandardCharsets.UTF_8))), DS_NS);
        emitIncorporate("incorporate-internally-detached", internallyDetachedParams(),
                Collections.singletonList(xml()), DS_NS);
        emitIncorporate("incorporate-two-documents", detachedParams(), Arrays.asList(text(), second()), DS_NS);
        emitIncorporate("incorporate-enveloped-dsig-prefix", envelopedParams(),
                Collections.singletonList(xml()), DSIG_NS);
        emitIncorporate("incorporate-enveloped-default-prefix", envelopedParams(),
                Collections.singletonList(xml()), DEFAULT_NS);

        XAdESSignatureParameters embedXml = envelopingParams();
        embedXml.setEmbedXML(true);
        emitIncorporate("incorporate-embed-xml", embedXml, Collections.singletonList(xml()), DS_NS);

        XAdESSignatureParameters manifest = envelopingParams();
        manifest.setManifestSignature(true);
        emitIncorporate("incorporate-manifest", manifest,
                Collections.singletonList(manifestDocument()), DS_NS);

        // The ManifestBuilder path: no signature parameters at all.
        ReferenceIdProvider provider = new ReferenceIdProvider();
        List<DSSReference> references =
                new ReferenceBuilder(Arrays.asList(text(), xml()), DigestAlgorithm.SHA256, provider).build();
        recordIncorporated("incorporate-no-parameters", new ReferenceProcessor(), references, DS_NS);

        // A digest-only document short-circuits getReferenceOutput.
        DSSReference digestReference = new DSSReference();
        digestReference.setId("r-digest");
        digestReference.setUri("digest.bin");
        digestReference.setDigestMethodAlgorithm(DigestAlgorithm.SHA256);
        digestReference.setContents(digestDocument());
        recordIncorporated("incorporate-digest-document", new ReferenceProcessor(),
                Collections.singletonList(digestReference), DS_NS);

        // An explicit reference over an already-signed XML, filtering the ds:Signature away.
        DSSReference filtered = new DSSReference();
        filtered.setId("r-filtered");
        filtered.setUri("");
        filtered.setDigestMethodAlgorithm(DigestAlgorithm.SHA256);
        filtered.setContents(signedXml());
        filtered.setTransforms(Arrays.asList(new XPath2FilterEnvelopedSignatureTransform(),
                new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)));
        recordIncorporated("incorporate-filtered-signature", new ReferenceProcessor(),
                Collections.singletonList(filtered), DS_NS);
    }

    static void emitIncorporate(String name, XAdESSignatureParameters params, List<DSSDocument> documents,
                                DSSNamespace namespace) {
        ReferenceIdProvider provider = new ReferenceIdProvider();
        provider.setSignatureParameters(params);
        List<DSSReference> references = new ReferenceBuilder(documents, params, provider).build();
        recordIncorporated(name, new ReferenceProcessor(params), references, namespace);
    }

    static void recordIncorporated(String name, ReferenceProcessor processor, List<DSSReference> references,
                                   DSSNamespace namespace) {
        Document document = DomUtils.buildDOM();
        Element signedInfoDom = DomUtils.createElementNS(document, namespace, XMLDSigElement.SIGNED_INFO);
        document.appendChild(signedInfoDom);
        processor.incorporateReferences(signedInfoDom, references, namespace);
        record(name, "references", DomUtils.serializeNode(signedInfoDom));
        for (int i = 0; i < references.size(); i++) {
            DSSDocument output = processor.getReferenceOutput(references.get(i));
            if (output instanceof eu.europa.esig.dss.model.DigestDocument) {
                // A DigestDocument is returned as-is and has no octets to stream.
                record(name, "output" + i + "-digest", output.getDigestValue(DigestAlgorithm.SHA256));
            } else {
                record(name, "output" + i, toByteArray(output));
            }
        }
    }

    // --------------------------------------------------------------------- 6. ReferenceVerifier

    static void emitReferenceVerifier() {
        // Valid: an enveloped signature carrying the filter transform.
        emitVerify("verify-enveloped-ok", envelopedParams(), Collections.singletonList(xml()));
        emitVerify("verify-enveloping-ok", envelopingParams(), Collections.singletonList(text()));
        emitVerify("verify-detached-ok", detachedParams(), Collections.singletonList(text()));

        // Enveloped without any transform.
        XAdESSignatureParameters envelopedNoTransform = envelopedParams();
        envelopedNoTransform.setReferences(Collections.singletonList(
                reference("r-1", "", DigestAlgorithm.SHA256, xml(), null)));
        emitVerifyParams("verify-enveloped-without-transform", envelopedNoTransform);

        // base64 with embedXML.
        XAdESSignatureParameters base64EmbedXml = envelopingParams();
        base64EmbedXml.setEmbedXML(true);
        base64EmbedXml.setReferences(Collections.singletonList(reference("r-1", "#o-r-1",
                DigestAlgorithm.SHA256, xml(), Collections.singletonList(new Base64Transform()))));
        emitVerifyParams("verify-base64-embed-xml", base64EmbedXml);

        // base64 with a manifest signature.
        XAdESSignatureParameters base64Manifest = envelopingParams();
        base64Manifest.setManifestSignature(true);
        base64Manifest.setReferences(Collections.singletonList(reference("r-1", "#o-r-1",
                DigestAlgorithm.SHA256, xml(), Collections.singletonList(new Base64Transform()))));
        emitVerifyParams("verify-base64-manifest", base64Manifest);

        // base64 with a non-enveloping packaging.
        XAdESSignatureParameters base64Detached = detachedParams();
        base64Detached.setReferences(Collections.singletonList(reference("r-1", "hello.txt",
                DigestAlgorithm.SHA256, text(), Collections.singletonList(new Base64Transform()))));
        emitVerifyParams("verify-base64-detached", base64Detached);

        // base64 with a second transform.
        XAdESSignatureParameters base64Chain = envelopingParams();
        base64Chain.setReferences(Collections.singletonList(reference("r-1", "#o-r-1",
                DigestAlgorithm.SHA256, text(), Arrays.asList(new Base64Transform(),
                        new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)))));
        emitVerifyParams("verify-base64-with-other-transform", base64Chain);

        // An element reference pointing at an Id the content does not carry.
        XAdESSignatureParameters missingId = internallyDetachedParams();
        missingId.setReferences(Collections.singletonList(reference("r-1", "#absent-id",
                DigestAlgorithm.SHA256, xml(),
                Collections.singletonList(new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)))));
        emitVerifyParams("verify-missing-element-id", missingId);

        // A reference without an Id gets a deterministic one assigned in place.
        XAdESSignatureParameters noId = envelopingParams();
        DSSReference withoutId = reference(null, "#o-r-1", DigestAlgorithm.SHA256, text(),
                Collections.singletonList(new Base64Transform()));
        noId.setReferences(Collections.singletonList(withoutId));
        emitVerifyParams("verify-generates-missing-id", noId);
        record("verify-generates-missing-id", "assigned-id", nullable(withoutId.getId()));
    }

    static void emitVerify(String name, XAdESSignatureParameters params, List<DSSDocument> documents) {
        ReferenceIdProvider provider = new ReferenceIdProvider();
        provider.setSignatureParameters(params);
        params.setReferences(new ReferenceBuilder(documents, params, provider).build());
        emitVerifyParams(name, params);
    }

    static void emitVerifyParams(String name, XAdESSignatureParameters params) {
        try {
            new ReferenceVerifier(params).checkReferencesValidity();
            record(name, "error", new byte[0]);
        } catch (RuntimeException e) {
            record(name, "error", message(e).getBytes(StandardCharsets.UTF_8));
        }
    }

    static DSSReference reference(String id, String uri, DigestAlgorithm digestAlgorithm, DSSDocument contents,
                                  List<DSSTransform> transforms) {
        DSSReference reference = new DSSReference();
        reference.setId(id);
        reference.setUri(uri);
        reference.setDigestMethodAlgorithm(digestAlgorithm);
        reference.setContents(contents);
        reference.setTransforms(transforms);
        return reference;
    }

    // ------------------------------------------------------ 7. ds:SignedInfo, end to end

    static void emitSignedInfo() {
        emitSignedInfo("si-enveloped", envelopedParams(), Collections.singletonList(xml()));
        emitSignedInfo("si-enveloping", envelopingParams(), Collections.singletonList(text()));
        emitSignedInfo("si-detached", detachedParams(), Collections.singletonList(text()));
        emitSignedInfo("si-internally-detached", internallyDetachedParams(), Collections.singletonList(xml()));

        XAdESSignatureParameters embedXml = envelopingParams();
        embedXml.setEmbedXML(true);
        emitSignedInfo("si-enveloping-embed-xml", embedXml, Collections.singletonList(xml()));

        XAdESSignatureParameters manifest = envelopingParams();
        manifest.setManifestSignature(true);
        emitSignedInfo("si-enveloping-manifest", manifest, Collections.singletonList(manifestDocument()));

        // Explicit references with the two other enveloped-filter flavours.
        XAdESSignatureParameters xpathEnveloped = envelopedParams();
        xpathEnveloped.setReferences(Collections.singletonList(reference("r-1", "", DigestAlgorithm.SHA256, xml(),
                Arrays.asList(new XPathEnvelopedSignatureTransform(),
                        new CanonicalizationTransform(XMLCanonicalizer.DEFAULT_DSS_C14N_METHOD)))));
        emitSignedInfo("si-enveloped-xpath-transform", xpathEnveloped, Collections.singletonList(xml()));

        XAdESSignatureParameters envelopedTransform = envelopedParams();
        envelopedTransform.setReferences(Collections.singletonList(reference("r-1", "", DigestAlgorithm.SHA256,
                xml(), Arrays.asList(new EnvelopedSignatureTransform(),
                        new CanonicalizationTransform("http://www.w3.org/TR/2001/REC-xml-c14n-20010315")))));
        emitSignedInfo("si-enveloped-enveloped-transform", envelopedTransform, Collections.singletonList(xml()));

        // A non-default xmldsig prefix reaches the transforms too.
        XAdESSignatureParameters dsigPrefix = envelopedParams();
        dsigPrefix.setXmldsigNamespace(DSIG_NS);
        emitSignedInfo("si-enveloped-dsig-prefix", dsigPrefix, Collections.singletonList(xml()));

        XAdESSignatureParameters defaultPrefix = envelopedParams();
        defaultPrefix.setXmldsigNamespace(DEFAULT_NS);
        emitSignedInfo("si-enveloped-default-prefix", defaultPrefix, Collections.singletonList(xml()));

        // Two documents, so the reference index runs past one.
        emitSignedInfo("si-detached-two-documents", detachedParams(), Arrays.asList(text(), second()));
    }

    static void emitSignedInfo(String name, XAdESSignatureParameters params, List<DSSDocument> documents) {
        XAdESSignatureBuilder builder =
                XAdESSignatureBuilder.getSignatureBuilder(params, documents, verifier());
        record(name, "signedinfo-c14n", builder.build());
        record(name, "signedinfo-raw", DomUtils.serializeNode(builder.signedInfoDom));
        record(name, "document", toByteArray(builder.signDocument(SIGNATURE_VALUE)));
    }

    // ------------------------------------------------------------------------------- fixtures

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

    static DSSDocument second() {
        return new InMemoryDocument("second".getBytes(StandardCharsets.UTF_8), "second.txt", MimeTypeEnum.TEXT);
    }

    static DSSDocument xml() {
        return new InMemoryDocument(XML_CONTENT.getBytes(StandardCharsets.UTF_8), "sample.xml", MimeTypeEnum.XML);
    }

    static DSSDocument signedXml() {
        return new InMemoryDocument(XML_SIGNED_CONTENT.getBytes(StandardCharsets.UTF_8), "signed.xml",
                MimeTypeEnum.XML);
    }

    static DSSDocument xmlWithoutId() {
        return new InMemoryDocument(
                ("<?xml version=\"1.0\" encoding=\"UTF-8\"?><root xmlns=\"http://sample.com\">"
                        + "<child>text</child></root>").getBytes(StandardCharsets.UTF_8),
                "no-id.xml", MimeTypeEnum.XML);
    }

    static DSSDocument manifestDocument() {
        return new ManifestBuilder(DigestAlgorithm.SHA256, Collections.singletonList(text())).build();
    }

    static DSSDocument digestDocument() {
        eu.europa.esig.dss.model.DigestDocument document = new eu.europa.esig.dss.model.DigestDocument(
                DigestAlgorithm.SHA256, DSSUtils.digest(DigestAlgorithm.SHA256,
                TEXT_CONTENT.getBytes(StandardCharsets.UTF_8)), "digest.bin");
        return document;
    }

    static CertificateVerifier verifier() {
        return new CommonCertificateVerifier();
    }

    // -------------------------------------------------------------------------------- helpers

    static void record(String name, String key, byte[] value) {
        out.append("case ").append(name).append('\n');
        out.append(key).append(' ').append(Base64.getEncoder().encodeToString(value)).append('\n');
    }

    /** Java's null becomes the literal "\0NULL"; an empty string stays empty. */
    static byte[] nullable(String value) {
        return value == null ? "\0NULL".getBytes(StandardCharsets.UTF_8)
                : value.getBytes(StandardCharsets.UTF_8);
    }

    static String message(Throwable e) {
        return e.getClass().getSimpleName() + ": " + e.getMessage();
    }

    static byte[] toByteArray(DSSDocument document) {
        try (InputStream is = document.openStream(); ByteArrayOutputStream os = new ByteArrayOutputStream()) {
            byte[] buffer = new byte[4096];
            int read;
            while ((read = is.read(buffer)) != -1) {
                os.write(buffer, 0, read);
            }
            return os.toByteArray();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    static Date utc(int year, int month, int day, int hour, int minute, int second) {
        GregorianCalendar calendar = new GregorianCalendar(TimeZone.getTimeZone("UTC"));
        calendar.clear();
        calendar.set(year, month, day, hour, minute, second);
        return calendar.getTime();
    }
}
