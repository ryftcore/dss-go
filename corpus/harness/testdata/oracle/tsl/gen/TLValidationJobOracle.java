import eu.europa.esig.dss.alert.detector.AlertDetector;
import eu.europa.esig.dss.alert.handler.AlertHandler;
import eu.europa.esig.dss.model.DSSDocument;
import eu.europa.esig.dss.model.DSSException;
import eu.europa.esig.dss.model.FileDocument;
import eu.europa.esig.dss.model.job.DownloadInfoRecord;
import eu.europa.esig.dss.model.job.InfoRecord;
import eu.europa.esig.dss.model.job.ParsingInfoRecord;
import eu.europa.esig.dss.model.job.ValidationInfoRecord;
import eu.europa.esig.dss.model.tsl.LOTLInfo;
import eu.europa.esig.dss.model.tsl.TLInfo;
import eu.europa.esig.dss.model.tsl.TrustProperties;
import eu.europa.esig.dss.model.tsl.TrustServiceProvider;
import eu.europa.esig.dss.model.tsl.TrustServiceStatusAndInformationExtensions;
import eu.europa.esig.dss.model.tsl.TLValidationJobSummary;
import eu.europa.esig.dss.model.x509.CertificateToken;
import eu.europa.esig.dss.spi.DSSUtils;
import eu.europa.esig.dss.spi.client.http.DSSFileLoader;
import eu.europa.esig.dss.spi.tsl.TrustedListsCertificateSource;
import eu.europa.esig.dss.spi.x509.CertificateSource;
import eu.europa.esig.dss.spi.x509.CommonCertificateSource;
import eu.europa.esig.dss.tsl.alerts.LOTLAlert;
import eu.europa.esig.dss.tsl.alerts.TLAlert;
import eu.europa.esig.dss.tsl.alerts.detections.LOTLLocationChangeDetection;
import eu.europa.esig.dss.tsl.alerts.detections.OJUrlChangeDetection;
import eu.europa.esig.dss.tsl.alerts.detections.TLParsingErrorDetection;
import eu.europa.esig.dss.tsl.alerts.detections.TLSignatureErrorDetection;
import eu.europa.esig.dss.tsl.function.OtherTSLPointerPredicate;
import eu.europa.esig.dss.tsl.job.TLValidationJob;
import eu.europa.esig.dss.tsl.source.LOTLSource;
import eu.europa.esig.dss.tsl.source.TLSource;
import eu.europa.esig.trustedlist.jaxb.tsl.OtherTSLPointerType;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

/**
 * Phase 9 harness contract item (C): produces testdata/oracle/tsl/tl_validation_job.json, the
 * Java-oracle dump tl_validation_job_oracle_test.go cross-validates against - see that file's
 * header for the exact scenario (one LOTLSource + three independently-configured TLSources, all
 * offline, all fixtures shared with contract item (A)).
 *
 * Usage: java TLValidationJobOracle &lt;output-file&gt;
 * (run from dss/harness so the "testdata/oracle/tsl/..." fixture paths below resolve)
 */
public class TLValidationJobOracle {

	private static final String SIGNING_CERTIFICATE =
			"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";

	public static void main(String[] args) throws Exception {
		Locale.setDefault(Locale.ENGLISH);
		Path outFile = Paths.get(args[0]);

		Map<String, DSSDocument> urlMap = new HashMap<>();
		urlMap.put("LOTL_URL", new FileDocument("testdata/oracle/tsl/eu-lotl.xml"));
		urlMap.put("DE_TL_URL", new FileDocument("testdata/oracle/tsl/de-tl.xml"));
		urlMap.put("BROKEN_TL_URL", new FileDocument("testdata/oracle/tsl/eu-lotl-broken-sig.xml"));
		urlMap.put("BAD_TL_URL", new FileDocument("testdata/oracle/tsl/eu-lotl-not-parseable.xml"));
		DSSFileLoader loader = url -> {
			DSSDocument document = urlMap.get(url);
			if (document == null) {
				throw new DSSException("no fixture mapped for url " + url);
			}
			return document;
		};

		LOTLSource lotlSource = new LOTLSource();
		lotlSource.setUrl("LOTL_URL");
		lotlSource.setCertificateSource(certificateSource(SIGNING_CERTIFICATE));
		lotlSource.setTlPredicate((OtherTSLPointerPredicate) (OtherTSLPointerType o) -> false);

		TLSource deTLSource = new TLSource();
		deTLSource.setUrl("DE_TL_URL");
		deTLSource.setCertificateSource(certificateSource());

		TLSource brokenTLSource = new TLSource();
		brokenTLSource.setUrl("BROKEN_TL_URL");
		brokenTLSource.setCertificateSource(certificateSource());

		TLSource badTLSource = new TLSource();
		badTLSource.setUrl("BAD_TL_URL");
		badTLSource.setCertificateSource(certificateSource());

		TLValidationJob job = new TLValidationJob();
		job.setListOfTrustedListSources(lotlSource);
		job.setTrustedListSources(deTLSource, brokenTLSource, badTLSource);
		job.setOfflineDataLoader(loader);
		TrustedListsCertificateSource certSource = new TrustedListsCertificateSource();
		job.setTrustedListCertificateSource(certSource);

		List<String> locationFired = new ArrayList<>();
		List<String> ojUrlFired = new ArrayList<>();
		List<String> tlSigFired = new ArrayList<>();
		List<String> tlParseFired = new ArrayList<>();

		job.setLOTLAlerts(java.util.Arrays.asList(
				new LOTLAlert(new LOTLLocationChangeDetection(lotlSource), recorder(locationFired)),
				new LOTLAlert(new OJUrlChangeDetection(lotlSource), recorder(ojUrlFired))));
		job.setTLAlerts(java.util.Arrays.asList(
				new TLAlert(new TLSignatureErrorDetection(), recorder(tlSigFired)),
				new TLAlert(new TLParsingErrorDetection(), recorder(tlParseFired))));

		job.offlineRefresh();

		TLValidationJobSummary summary = job.getSummary();

		StringBuilder out = new StringBuilder();
		out.append("{");
		out.append("\"numberOfProcessedLOTLs\":").append(summary.getNumberOfProcessedLOTLs()).append(",");
		out.append("\"numberOfProcessedTLs\":").append(summary.getNumberOfProcessedTLs()).append(",");

		out.append("\"lotls\":[");
		boolean first = true;
		for (LOTLInfo l : summary.getLOTLInfos()) {
			if (!first) out.append(",");
			first = false;
			out.append(dumpTLInfo(l));
		}
		out.append("],");

		out.append("\"tls\":[");
		first = true;
		for (TLInfo t : summary.getOtherTLInfos()) {
			if (!first) out.append(",");
			first = false;
			out.append(dumpTLInfo(t));
		}
		out.append("],");

		List<String> allFired = new ArrayList<>();
		for (String u : locationFired) allFired.add("LOTL_LOCATION:" + u);
		for (String u : ojUrlFired) allFired.add("OJ_URL:" + u);
		for (String u : tlSigFired) allFired.add("TL_SIG:" + u);
		for (String u : tlParseFired) allFired.add("TL_PARSE:" + u);
		Collections.sort(allFired);
		out.append("\"firedAlerts\":[");
		first = true;
		for (String s : allFired) {
			if (!first) out.append(",");
			first = false;
			out.append(json(s));
		}
		out.append("],");

		out.append("\"trustedCertificates\":").append(dumpCertSource(certSource));
		out.append("}");

		Files.write(outFile, (out.toString() + "\n").getBytes(StandardCharsets.UTF_8));
		System.out.println("wrote " + outFile);
	}

	private static <T extends TLInfo> AlertHandler<T> recorder(List<String> sink) {
		return currentInfo -> sink.add(currentInfo.getUrl());
	}

	private static CertificateSource certificateSource(String... base64Certificates) {
		CommonCertificateSource source = new CommonCertificateSource();
		for (String b64 : base64Certificates) {
			source.addCertificate(DSSUtils.loadCertificateFromBase64EncodedString(b64));
		}
		return source;
	}

	private static String dumpTLInfo(TLInfo info) {
		StringBuilder sb = new StringBuilder();
		sb.append("{");
		sb.append("\"url\":").append(json(info.getUrl())).append(",");
		sb.append("\"download\":").append(dumpInfoRecord(info.getDownloadCacheInfo())).append(",");
		sb.append("\"parsing\":").append(dumpInfoRecord(info.getParsingCacheInfo())).append(",");

		ValidationInfoRecord v = info.getValidationCacheInfo();
		sb.append("\"validation\":{");
		sb.append(dumpInfoRecordFields(v));
		sb.append(",\"indication\":").append(json(v != null && v.getIndication() != null ? v.getIndication().name() : ""));
		sb.append(",\"subIndication\":").append(json(v != null && v.getSubIndication() != null ? v.getSubIndication().name() : ""));
		sb.append("}");
		sb.append("}");
		return sb.toString();
	}

	private static String dumpInfoRecord(InfoRecord r) {
		return "{" + dumpInfoRecordFields(r) + "}";
	}

	private static String dumpInfoRecordFields(InfoRecord r) {
		if (r == null) {
			return "\"statusName\":\"\",\"synchronized\":false,\"desynchronized\":false,\"error\":false,\"resultExist\":false";
		}
		return "\"statusName\":" + json(r.getStatusName()) +
				",\"synchronized\":" + r.isSynchronized() +
				",\"desynchronized\":" + r.isDesynchronized() +
				",\"error\":" + r.isError() +
				",\"resultExist\":" + r.isResultExist();
	}

	private static String dumpCertSource(TrustedListsCertificateSource source) throws Exception {
		MessageDigest sha256 = MessageDigest.getInstance("SHA-256");
		TreeSet<String> byHash = new TreeSet<>();
		Map<String, String> entryByHash = new HashMap<>();
		for (CertificateToken cert : source.getCertificates()) {
			String hash = toHex(sha256.digest(cert.getEncoded()));
			sha256.reset();

			Set<String> tlUrls = new LinkedHashSet<>();
			Set<String> territories = new LinkedHashSet<>();
			Set<String> typeStatus = new LinkedHashSet<>();
			for (TrustProperties tp : source.getTrustServices(cert)) {
				if (tp.getTLInfo() != null) {
					tlUrls.add(tp.getTLInfo().getUrl());
				}
				TrustServiceProvider tsp = tp.getTrustServiceProvider();
				if (tsp != null) {
					territories.add(tsp.getTerritory());
				}
				if (tp.getTrustService() != null) {
					for (TrustServiceStatusAndInformationExtensions entry : tp.getTrustService()) {
						typeStatus.add(entry.getType() + "|" + entry.getStatus());
					}
				}
			}

			StringBuilder sb = new StringBuilder();
			sb.append("{\"sha256\":").append(json(hash));
			sb.append(",\"tlUrls\":").append(jsonArray(new TreeSet<>(tlUrls)));
			sb.append(",\"tspTerritories\":").append(jsonArray(new TreeSet<>(territories)));
			sb.append(",\"serviceTypeStatus\":").append(jsonArray(new TreeSet<>(typeStatus)));
			sb.append("}");

			byHash.add(hash);
			entryByHash.put(hash, sb.toString());
		}

		StringBuilder out = new StringBuilder("[");
		boolean first = true;
		for (String hash : byHash) {
			if (!first) out.append(",");
			first = false;
			out.append(entryByHash.get(hash));
		}
		out.append("]");
		return out.toString();
	}

	private static String jsonArray(Iterable<String> values) {
		StringBuilder sb = new StringBuilder("[");
		boolean first = true;
		for (String v : values) {
			if (!first) sb.append(",");
			first = false;
			sb.append(json(v));
		}
		sb.append("]");
		return sb.toString();
	}

	private static String toHex(byte[] bytes) {
		StringBuilder sb = new StringBuilder(bytes.length * 2);
		for (byte b : bytes) {
			sb.append(String.format("%02x", b));
		}
		return sb.toString();
	}

	private static String json(String s) {
		if (s == null) {
			return "null";
		}
		StringBuilder sb = new StringBuilder("\"");
		for (int i = 0; i < s.length(); i++) {
			char c = s.charAt(i);
			switch (c) {
				case '"': sb.append("\\\""); break;
				case '\\': sb.append("\\\\"); break;
				case '\n': sb.append("\\n"); break;
				case '\r': sb.append("\\r"); break;
				case '\t': sb.append("\\t"); break;
				default:
					if (c < 0x20) {
						sb.append(String.format("\\u%04x", (int) c));
					} else {
						sb.append(c);
					}
			}
		}
		sb.append("\"");
		return sb.toString();
	}
}
