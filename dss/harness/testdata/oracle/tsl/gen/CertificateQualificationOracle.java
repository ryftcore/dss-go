import eu.europa.esig.dss.diagnostic.CertificateWrapper;
import eu.europa.esig.dss.diagnostic.DiagnosticData;
import eu.europa.esig.dss.enumerations.Indication;
import eu.europa.esig.dss.enumerations.SubIndication;
import eu.europa.esig.dss.enumerations.CertificateQualification;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.DSSException;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.simplecertificatereport.SimpleCertificateReport;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.client.http.DSSFileLoader;
import eu.europa.esig.dss.spi.tsl.TrustedListsCertificateSource;
import eu.europa.esig.dss.spi.validation.CertificateVerifier;
import eu.europa.esig.dss.spi.validation.CommonCertificateVerifier;
import eu.europa.esig.dss.spi.x509.CommonCertificateSource;
import eu.europa.esig.dss.tsl.job.TLValidationJob;
import eu.europa.esig.dss.tsl.source.TLSource;
import eu.europa.esig.dss.validation.CertificateValidator;
import eu.europa.esig.dss.validation.reports.CertificateReports;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

/**
 * Phase 9 harness contract item (D): produces testdata/oracle/tsl/certificate_qualification.json,
 * the Java-oracle dump certificate_qualification_oracle_test.go cross-validates against - see
 * that file's header for the exact scenario (ported verbatim from
 * eu.europa.esig.dss.tsl.validation.SKCertificateTest#skTLTest).
 *
 * Usage: java CertificateQualificationOracle &lt;output-file&gt;
 * (run from dss/harness so "../tsl/testdata/sk-tl-sn-95.xml" resolves)
 */
public class CertificateQualificationOracle {

	private static final String TL_ISSUER =
			"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";

	private static final String CERTIFICATE =
			"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";

	public static void main(String[] args) throws Exception {
		Locale.setDefault(Locale.ENGLISH);
		Path outFile = Paths.get(args[0]);

		String tlUrl = "sk-tl.xml";
		Map<String, DSSDocument> urlMap = new HashMap<>();
		urlMap.put(tlUrl, new FileDocument("../tsl/testdata/sk-tl-sn-95.xml"));
		DSSFileLoader loader = url -> {
			DSSDocument document = urlMap.get(url);
			if (document == null) {
				throw new DSSException("no fixture mapped for url " + url);
			}
			return document;
		};

		CommonCertificateSource issuerSource = new CommonCertificateSource();
		issuerSource.addCertificate(DSSUtils.loadCertificateFromBase64EncodedString(TL_ISSUER));

		TLSource tlSource = new TLSource();
		tlSource.setUrl(tlUrl);
		tlSource.setTLVersions(Arrays.asList(5, 6));
		tlSource.setCertificateSource(issuerSource);

		TLValidationJob tlValidationJob = new TLValidationJob();
		tlValidationJob.setTrustedListSources(tlSource);
		tlValidationJob.setOfflineDataLoader(loader);
		TrustedListsCertificateSource trustedCertificateSource = new TrustedListsCertificateSource();
		tlValidationJob.setTrustedListCertificateSource(trustedCertificateSource);

		tlValidationJob.offlineRefresh();

		CertificateToken certificate = DSSUtils.loadCertificateFromBase64EncodedString(CERTIFICATE);

		CertificateValidator certificateValidator = CertificateValidator.fromCertificate(certificate);
		CertificateVerifier certificateVerifier = new CommonCertificateVerifier();
		certificateVerifier.setTrustedCertSources(trustedCertificateSource);
		certificateValidator.setCertificateVerifier(certificateVerifier);

		CertificateReports reports = certificateValidator.validate();
		SimpleCertificateReport simpleReport = reports.getSimpleReport();
		String certId = certificate.getDSSIdAsString();

		Indication indication = simpleReport.getCertificateIndication(certId);
		SubIndication subIndication = simpleReport.getCertificateSubIndication(certId);
		CertificateQualification qualification = simpleReport.getQualificationAtCertificateIssuance();

		DiagnosticData diagnosticData = reports.getDiagnosticData();
		CertificateWrapper certificateWrapper = diagnosticData.getCertificateById(certId);
		int numberOfTrustServiceProviders = certificateWrapper != null
				? certificateWrapper.getTrustServiceProviders().size() : 0;

		StringBuilder out = new StringBuilder();
		out.append("{");
		out.append("\"numberOfTrustedCertificates\":").append(trustedCertificateSource.getCertificates().size()).append(",");
		out.append("\"indication\":").append(json(indication != null ? indication.name() : "")).append(",");
		out.append("\"subIndication\":").append(json(subIndication != null ? subIndication.name() : "")).append(",");
		out.append("\"qualificationAtIssuance\":").append(json(qualification != null ? qualification.name() : "")).append(",");
		out.append("\"numberOfTrustServiceProviders\":").append(numberOfTrustServiceProviders);
		out.append("}");

		Files.write(outFile, (out.toString() + "\n").getBytes(StandardCharsets.UTF_8));
		System.out.println("wrote " + outFile);
	}

	private static String json(String s) {
		return "\"" + s.replace("\\", "\\\\").replace("\"", "\\\"") + "\"";
	}
}
